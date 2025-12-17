package api

import (
	"context"
	"net/http"

	"github.com/edvirons/ssp/ims/internal/admin"
	"github.com/edvirons/ssp/ims/internal/audit"
	"github.com/edvirons/ssp/ims/internal/auth"
	"github.com/edvirons/ssp/ims/internal/blob"
	"github.com/edvirons/ssp/ims/internal/config"
	"github.com/edvirons/ssp/ims/internal/handlers"
	"github.com/edvirons/ssp/ims/internal/metrics"
	"github.com/edvirons/ssp/ims/internal/middleware"
	"github.com/edvirons/ssp/ims/internal/store"
	"github.com/edvirons/ssp/ims/internal/ws"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// Server represents the HTTP API server.
type Server struct {
	cfg    config.Config
	logger *zap.Logger
	r      chi.Router

	pg  *store.Postgres
	rdb *redis.Client

	authVerifier *auth.Verifier
	wsHub        *ws.Hub
}

// NewServer creates a new HTTP server with all routes configured.
func NewServer(cfg config.Config, logger *zap.Logger, pg *store.Postgres, rdb *redis.Client) *Server {
	s := &Server{cfg: cfg, logger: logger, pg: pg, rdb: rdb}
	s.r = chi.NewRouter()

	// Initialize WebSocket hub
	s.wsHub = ws.NewHub(logger)
	go s.wsHub.Run()

	s.setupMiddleware()
	s.setupHealthAndMetrics()

	blobClient := s.initBlobClient()
	s.setupAPIRoutes(blobClient)
	s.setupAdminRoutes()

	return s
}

// WSHub returns the WebSocket hub for broadcasting events
func (s *Server) WSHub() *ws.Hub {
	return s.wsHub
}

// Router returns the HTTP handler for the server.
func (s *Server) Router() http.Handler { return s.r }

// setupMiddleware configures all middleware for the server.
func (s *Server) setupMiddleware() {
	// Security middleware - applied first
	s.r.Use(middleware.SecurityHeaders())
	s.r.Use(middleware.MaxBodySize(10 * 1024 * 1024)) // 10MB limit

	s.r.Use(middleware.RequestID())
	s.r.Use(middleware.Recoverer(s.logger))
	s.r.Use(middleware.Logger(s.logger))
	s.r.Use(middleware.MetricsMiddleware())

	s.r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   splitCSV(s.cfg.CORSAllowedOrigins),
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Request-Id", s.cfg.TenantHeader, s.cfg.SchoolHeader, "X-Impersonate-User", "X-Impersonate-Reason"},
		ExposedHeaders:   []string{"X-Request-Id", "X-Impersonation-Active", "X-Impersonated-User", "X-Impersonated-Schools"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	if s.cfg.AuthEnabled {
		s.authVerifier = auth.NewVerifier(s.cfg.AuthIssuer, s.cfg.AuthJWKSURL, s.cfg.AuthAudience)
		s.r.Use(middleware.AuthJWT(s.authVerifier, s.logger))
	}
	s.r.Use(middleware.Tenancy(s.cfg))

	// Audit middleware - captures request context for audit logging
	s.r.Use(audit.Middleware())
}

// setupHealthAndMetrics configures health check and metrics endpoints.
func (s *Server) setupHealthAndMetrics() {
	// Health - no rate limiting on health checks
	s.r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	s.r.Get("/readyz", handlers.ReadyzHandler(s.pg, s.rdb))

	// Metrics endpoint
	s.r.Get("/metrics", func(w http.ResponseWriter, r *http.Request) {
		metrics.Handler().ServeHTTP(w, r)
	})
}

// initBlobClient initializes the blob storage client.
func (s *Server) initBlobClient() *blob.MinIO {
	blobClient, err := blob.NewMinIO(s.cfg)
	if err != nil {
		s.logger.Fatal("minio client init failed", zap.Error(err))
	}
	if err := blobClient.EnsureBucket(context.Background()); err != nil {
		s.logger.Fatal("minio ensure bucket failed", zap.Error(err))
	}
	return blobClient
}

// setupAPIRoutes configures all v1 API routes.
func (s *Server) setupAPIRoutes(blobClient *blob.MinIO) {
	s.r.Route("/v1", func(r chi.Router) {
		// Apply rate limiting to all v1 routes if enabled
		if s.cfg.RateLimitEnabled {
			defaultRateLimit := middleware.RateLimitConfig{
				RequestsPerMinute: s.cfg.RateLimitReadRPM,
				BurstSize:         s.cfg.RateLimitReadRPM + s.cfg.RateLimitBurst,
				KeyPrefix:         "ratelimit",
			}
			r.Use(middleware.RateLimit(s.rdb, defaultRateLimit, s.logger))
		}

		// Initialize audit store and logger
		auditStore := audit.NewStore(s.pg.AuditStorePool())
		auditLogger := audit.NewLogger(auditStore)

		// Initialize handlers
		inc := handlers.NewIncidentHandler(s.cfg, s.logger, s.pg, s.rdb, auditLogger)
		wo := handlers.NewWorkOrderHandler(s.logger, s.pg, s.rdb, auditLogger)
		att := handlers.NewAttachmentHandler(s.cfg, s.logger, s.pg, blobClient)
		tel := handlers.NewTelemetryHandler(s.cfg, s.logger, s.pg)

		sch := handlers.NewSchoolHandler(s.logger, s.pg)
		contacts := handlers.NewSchoolContactsHandler(s.logger, s.pg)
		shops := handlers.NewServiceShopHandler(s.logger, s.pg)
		staff := handlers.NewServiceStaffHandler(s.logger, s.pg)
		parts := handlers.NewPartsHandler(s.logger, s.pg)
		inv := handlers.NewInventoryHandler(s.logger, s.pg)
		whDash := handlers.NewWarehouseDashboardHandler(s.logger, s.pg)
		ltDash := handlers.NewLeadTechDashboardHandler(s.logger, s.pg)
		saDash := handlers.NewSupportAgentDashboardHandler(s.logger, s.pg)
		bom := handlers.NewBOMHandler(s.logger, s.pg)
		woops := handlers.NewWorkOrderOpsHandler(s.logger, s.pg)
		sync := handlers.NewSSOTSyncHandler(s.cfg, s.logger, s.pg)
		ssotList := handlers.NewSSOTListHandler(s.cfg, s.logger, s.pg)
		wh := handlers.NewSSOTWebhookHandler(s.cfg, s.logger, s.pg)

		proj := handlers.NewProjectsHandler(s.logger, s.pg)
		ph := handlers.NewPhasesHandler(s.logger, s.pg)
		surv := handlers.NewSurveysHandler(s.logger, s.pg)
		boq := handlers.NewBOQHandler(s.logger, s.pg)
		auditLogs := handlers.NewAuditLogsHandler(s.logger, auditStore)

		// Project collaboration handlers
		projectTeam := handlers.NewProjectTeamHandler(s.logger, s.pg)
		projectActivities := handlers.NewProjectActivitiesHandler(s.logger, s.pg)
		projectWorkOrders := handlers.NewProjectWorkOrdersHandler(s.logger, s.pg, projectActivities)
		userNotifications := handlers.NewUserNotificationsHandler(s.logger, s.pg)

		// Work order enhancement handlers
		woUpdate := handlers.NewWorkOrderUpdateHandler(s.logger, s.pg, auditLogger)
		woRework := handlers.NewWorkOrderReworkHandler(s.logger, s.pg, auditLogger)
		woBulk := handlers.NewWorkOrderBulkHandler(s.logger, s.pg, auditLogger)

		// Reports handler
		rpt := handlers.NewReportsHandler(s.logger, s.pg)

		// EdTech profiles handler
		edtech := handlers.NewEdTechProfilesHandler(s.logger, s.pg, s.cfg)

		// Sales/Marketing handlers
		demoPipeline := handlers.NewDemoPipelineHandler(s.logger, s.pg)
		presentations := handlers.NewPresentationsHandler(s.logger, s.pg, blobClient)
		salesMetrics := handlers.NewSalesMetricsHandler(s.logger, s.pg)

		// Knowledge Base handler
		kbArticles := handlers.NewKBArticlesHandler(s.logger, s.pg)

		// Marketing Knowledge Base handler
		marketingKB := handlers.NewMarketingKBHandler(s.logger, s.pg)

		// Device inventory handler
		deviceInv := handlers.NewDeviceInventoryHandler(s.logger, s.pg, auditLogger)

		// Impersonation handler
		impersonation := handlers.NewImpersonationHandler(s.logger, s.pg)

		// Add impersonation middleware - must be after auth middleware
		r.Use(middleware.Impersonation(s.logger, impersonation.LoadImpersonationTarget))

		// Mount routes from separate files
		s.mountIncidentRoutes(r, inc)
		s.mountWorkOrderRoutes(r, wo, woops, woUpdate, woRework, woBulk)
		s.mountBOMRoutes(r, bom)
		s.mountSSOTRoutes(r, sync, ssotList, wh)
		s.mountProjectRoutes(r, proj, ph, surv, boq, projectTeam, projectActivities, projectWorkOrders)
		s.mountServiceShopRoutes(r, shops, staff, parts, inv, whDash)
		s.mountLeadTechDashboardRoutes(r, ltDash)
		s.mountSupportAgentDashboardRoutes(r, saDash)
		s.mountAdminRoutes(r, auditLogs, sch, contacts, att, tel)
		s.mountNotificationRoutes(r, userNotifications)
		s.mountReportRoutes(r, rpt)
		s.mountEdTechRoutes(r, edtech)
		s.mountSalesRoutes(r, demoPipeline, presentations, salesMetrics)
		s.mountKBRoutes(r, kbArticles)
		s.mountMarketingKBRoutes(r, marketingKB)
		s.mountDeviceInventoryRoutes(r, deviceInv)
		s.mountImpersonationRoutes(r, impersonation)

		// Messaging routes
		RegisterMessagingRoutes(r, s.logger, s.pg, s.wsHub)

		// AI Chat and Livechat routes
		RegisterAIChatRoutes(r, s.logger, s.pg, s.wsHub, s.cfg)
	})
}

// setupAdminRoutes configures admin dashboard routes.
func (s *Server) setupAdminRoutes() {
	adminHandler := admin.NewHandler(s.cfg, s.logger, s.pg)

	// Auth routes - no authentication required (directly on main router under /v1)
	s.r.Post("/v1/auth/login", adminHandler.Login)
	s.r.Post("/v1/auth/logout", adminHandler.Logout)
	s.r.Get("/v1/auth/me", adminHandler.Me)
	s.r.Post("/v1/auth/refresh", adminHandler.Refresh)
	s.r.Get("/v1/auth/profile", adminHandler.Profile)

	// Protected admin API routes - require admin authentication
	s.r.Group(func(r chi.Router) {
		r.Use(adminHandler.AdminAuthMiddleware())
		r.Get("/v1/health/services", adminHandler.GetServicesHealth)
		r.Get("/v1/metrics/summary", adminHandler.GetMetricsSummary)
		r.Get("/v1/activity", adminHandler.GetActivityFeed)
		r.Get("/v1/notifications", adminHandler.GetNotifications)
		r.Get("/v1/notifications/unread-count", adminHandler.GetUnreadCount)
		r.Post("/v1/notifications/mark-read", adminHandler.MarkNotificationsRead)
	})

	// WebSocket endpoint for real-time notifications
	wsHandler := ws.NewHandler(s.wsHub, s.logger)
	s.r.Get("/ws", wsHandler.ServeWS)

	// Dashboard static file serving from root
	// Static files don't require auth - served for all
	// The dashboard itself will check auth via API calls
	admin.ServeDashboard(s.r)
}

// writeRateLimitMiddleware returns a rate limit middleware for write operations.
func (s *Server) writeRateLimitMiddleware() func(http.Handler) http.Handler {
	if !s.cfg.RateLimitEnabled {
		return func(next http.Handler) http.Handler {
			return next
		}
	}
	writeRateLimit := middleware.RateLimitConfig{
		RequestsPerMinute: s.cfg.RateLimitWriteRPM,
		BurstSize:         s.cfg.RateLimitWriteRPM + (s.cfg.RateLimitBurst / 2),
		KeyPrefix:         "ratelimit",
	}
	return middleware.RateLimit(s.rdb, writeRateLimit, s.logger)
}

// bulkRateLimitMiddleware returns a stricter rate limit middleware for bulk operations.
func (s *Server) bulkRateLimitMiddleware() func(http.Handler) http.Handler {
	if !s.cfg.RateLimitEnabled {
		return func(next http.Handler) http.Handler {
			return next
		}
	}
	// Bulk operations have stricter limits: 10 requests per minute
	bulkRateLimit := middleware.RateLimitConfig{
		RequestsPerMinute: 10,
		BurstSize:         3,
		KeyPrefix:         "ratelimit_bulk",
	}
	return middleware.RateLimit(s.rdb, bulkRateLimit, s.logger)
}

// splitCSV splits a comma-separated string into a slice.
func splitCSV(s string) []string {
	out := []string{}
	cur := ""
	for _, ch := range s {
		if ch == ',' {
			if t := trim(cur); t != "" {
				out = append(out, t)
			}
			cur = ""
			continue
		}
		cur += string(ch)
	}
	if t := trim(cur); t != "" {
		out = append(out, t)
	}
	if len(out) == 0 {
		return []string{"*"}
	}
	return out
}

// trim removes leading and trailing whitespace.
func trim(s string) string {
	i, j := 0, len(s)-1
	for i <= j && (s[i] == ' ' || s[i] == '\n' || s[i] == '\t' || s[i] == '\r') {
		i++
	}
	for j >= i && (s[j] == ' ' || s[j] == '\n' || s[j] == '\t' || s[j] == '\r') {
		j--
	}
	if i > j {
		return ""
	}
	return s[i : j+1]
}
