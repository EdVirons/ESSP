package admin

import (
	"net/http"
	"time"

	"github.com/edvirons/ssp/ims/internal/config"
	"github.com/edvirons/ssp/ims/internal/store"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// Handler holds dependencies for admin endpoints
type Handler struct {
	cfg       config.Config
	logger    *zap.Logger
	pg        *store.Postgres
	adminAuth *AdminAuth
}

// NewHandler creates a new admin handler
func NewHandler(cfg config.Config, logger *zap.Logger, pg *store.Postgres) *Handler {
	authCfg := AdminAuthConfig{
		AdminUsername: cfg.AdminUsername,
		AdminPassword: cfg.AdminPassword,
		TokenExpiry:   time.Duration(cfg.AdminJWTExpiry) * time.Hour,
		CookieSecure:  cfg.AdminCookieSecure,
	}

	adminAuth, err := NewAdminAuth(authCfg, logger)
	if err != nil {
		logger.Fatal("failed to create admin auth", zap.Error(err))
	}

	return &Handler{
		cfg:       cfg,
		logger:    logger,
		pg:        pg,
		adminAuth: adminAuth,
	}
}

// RegisterAuthRoutes registers authentication routes (no auth required)
func (h *Handler) RegisterAuthRoutes(r chi.Router) {
	r.Post("/auth/login", h.adminAuth.Login)
	r.Post("/auth/logout", h.adminAuth.Logout)
	r.Get("/auth/me", h.adminAuth.Me)
	r.Post("/auth/refresh", h.adminAuth.Refresh)
	// Profile endpoint returns extended SSO user profile
	r.Get("/auth/profile", h.adminAuth.Profile)
}

// RegisterAPIRoutes registers admin API routes (used with /admin/v1 prefix)
// These routes require authentication
func (h *Handler) RegisterAPIRoutes(r chi.Router) {
	// Health aggregation
	r.Get("/health/services", h.GetServicesHealth)

	// Metrics summary
	r.Get("/metrics/summary", h.GetMetricsSummary)

	// Activity feed
	r.Get("/activity", h.GetActivityFeed)

	// Notifications
	r.Get("/notifications", h.GetNotifications)
	r.Get("/notifications/unread-count", h.GetUnreadCount)
	r.Post("/notifications/mark-read", h.MarkNotificationsRead)
}

// AdminAuthMiddleware returns the admin authentication middleware
func (h *Handler) AdminAuthMiddleware() func(http.Handler) http.Handler {
	return h.adminAuth.AdminAuthMiddleware
}

// Login handles admin login requests
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	h.adminAuth.Login(w, r)
}

// Logout handles admin logout requests
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	h.adminAuth.Logout(w, r)
}

// Me returns the current authenticated user
func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	h.adminAuth.Me(w, r)
}

// Refresh refreshes the auth token
func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	h.adminAuth.Refresh(w, r)
}

// Profile returns the extended SSO user profile
func (h *Handler) Profile(w http.ResponseWriter, r *http.Request) {
	h.adminAuth.Profile(w, r)
}
