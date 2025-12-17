package api

import (
	"github.com/edvirons/ssp/ims/internal/audit"
	"github.com/edvirons/ssp/ims/internal/auth"
	"github.com/edvirons/ssp/ims/internal/handlers"
	"github.com/edvirons/ssp/ims/internal/middleware"
	"github.com/go-chi/chi/v5"
)

// mountAdminRoutes registers routes that require admin privileges.
func (s *Server) mountAdminRoutes(r chi.Router, auditLogs *handlers.AuditLogsHandler, sch *handlers.SchoolHandler, contacts *handlers.SchoolContactsHandler, att *handlers.AttachmentHandler, tel *handlers.TelemetryHandler) {
	// Audit logs - admin only
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireRole("ssp_admin", s.logger))
		r.Get("/audit-logs", auditLogs.List)
		r.Get("/audit-logs/{id}", auditLogs.GetByID)
	})

	// Schools - upsert (admin only)
	r.Group(func(r chi.Router) {
		r.Use(s.writeRateLimitMiddleware())
		r.Use(middleware.RequireAnyPermission(s.logger, auth.PermSchoolCreate, auth.PermSchoolUpdate))
		r.Post("/schools/upsert", sch.Upsert)
	})

	// School Contacts - read operations
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequirePermission(auth.PermSchoolContactRead, s.logger))
		r.Get("/schools/{schoolId}/contacts", contacts.List)
	})

	// School Contacts - create/update operations
	r.Group(func(r chi.Router) {
		r.Use(s.writeRateLimitMiddleware())
		r.Use(middleware.RequireAnyPermission(s.logger, auth.PermSchoolContactCreate, auth.PermSchoolContactUpdate))
		r.Post("/schools/{schoolId}/contacts", contacts.Create)
		r.Patch("/schools/{schoolId}/contacts/primary", contacts.SetPrimary)
	})

	// Attachments - read operations
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequirePermission(auth.PermAttachmentRead, s.logger))
		r.Get("/attachments/{id}", att.GetByID)
		r.Get("/attachments", att.List)
		r.Get("/attachments/{id}/download-url", att.DownloadURL)
	})

	// Attachments - create operations
	r.Group(func(r chi.Router) {
		r.Use(s.writeRateLimitMiddleware())
		r.Use(middleware.RequirePermission(auth.PermAttachmentCreate, s.logger))
		r.Post("/attachments", att.Create)
		r.Post("/attachments/{id}/upload-url", att.UploadURL)
	})

	// Telemetry - ingest
	r.Group(func(r chi.Router) {
		r.Use(s.writeRateLimitMiddleware())
		r.Use(middleware.RequirePermission(auth.PermTelemetryIngest, s.logger))
		r.Post("/telemetry/events", tel.Ingest)
	})
}

// mountLeadTechDashboardRoutes registers lead tech dashboard routes.
func (s *Server) mountLeadTechDashboardRoutes(r chi.Router, ltDash *handlers.LeadTechDashboardHandler) {
	// Lead Tech Dashboard - read operations (for lead tech dashboard)
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequirePermission(auth.PermWorkOrderApproval, s.logger))
		r.Get("/leadtech/dashboard", ltDash.GetDashboardSummary)
		r.Get("/leadtech/approvals", ltDash.GetPendingApprovals)
		r.Get("/leadtech/schedule", ltDash.GetTodaysSchedule)
		r.Get("/leadtech/team-metrics", ltDash.GetTeamMetrics)
	})
}

// mountSupportAgentDashboardRoutes registers support agent dashboard routes.
func (s *Server) mountSupportAgentDashboardRoutes(r chi.Router, saDash *handlers.SupportAgentDashboardHandler) {
	// Support Agent Dashboard - read operations (for support agent dashboard)
	// Uses chat:accept permission as the key permission for support agents
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequirePermission(auth.PermChatAccept, s.logger))
		r.Get("/supportagent/dashboard", saDash.GetDashboardSummary)
		r.Get("/supportagent/incidents", saDash.GetIncidentQueue)
		r.Get("/supportagent/chats", saDash.GetChatQueue)
		r.Get("/supportagent/work-orders", saDash.GetWorkOrderQueue)
		r.Get("/supportagent/metrics", saDash.GetIncidentMetrics)
	})
}

// mountServiceShopRoutes registers service shop, staff, parts, inventory, and warehouse dashboard routes.
func (s *Server) mountServiceShopRoutes(r chi.Router, shops *handlers.ServiceShopHandler, staff *handlers.ServiceStaffHandler, parts *handlers.PartsHandler, inv *handlers.InventoryHandler, whDash *handlers.WarehouseDashboardHandler) {
	// Service Shops - read operations
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequirePermission(auth.PermServiceShopRead, s.logger))
		r.Get("/service-shops/{id}", shops.GetByID)
		r.Get("/service-shops", shops.List)
	})

	// Service Shops - create operations
	r.Group(func(r chi.Router) {
		r.Use(s.writeRateLimitMiddleware())
		r.Use(middleware.RequirePermission(auth.PermServiceShopCreate, s.logger))
		r.Post("/service-shops", shops.Create)
	})

	// Service Staff - read operations
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequirePermission(auth.PermServiceStaffRead, s.logger))
		r.Get("/service-staff/stats", staff.GetStats)
		r.Get("/service-staff/{id}", staff.GetByID)
		r.Get("/service-staff", staff.List)
	})

	// Service Staff - create operations
	r.Group(func(r chi.Router) {
		r.Use(s.writeRateLimitMiddleware())
		r.Use(middleware.RequirePermission(auth.PermServiceStaffCreate, s.logger))
		r.Post("/service-staff", staff.Create)
	})

	// Service Staff - update operations
	r.Group(func(r chi.Router) {
		r.Use(s.writeRateLimitMiddleware())
		r.Use(middleware.RequirePermission(auth.PermServiceStaffUpdate, s.logger))
		r.Patch("/service-staff/{id}", staff.Update)
	})

	// Service Staff - delete operations
	r.Group(func(r chi.Router) {
		r.Use(s.writeRateLimitMiddleware())
		r.Use(middleware.RequirePermission(auth.PermServiceStaffUpdate, s.logger))
		r.Delete("/service-staff/{id}", staff.Delete)
	})

	// Parts - read operations
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequirePermission(auth.PermPartsRead, s.logger))
		r.Get("/parts/categories", parts.GetCategories)
		r.Get("/parts/stats", parts.GetStats)
		r.Get("/parts/{id}", parts.GetByID)
		r.Get("/parts", parts.List)
	})

	// Parts - create operations
	r.Group(func(r chi.Router) {
		r.Use(s.writeRateLimitMiddleware())
		r.Use(middleware.RequirePermission(auth.PermPartsCreate, s.logger))
		r.Post("/parts", parts.Create)
		r.Post("/parts/import", parts.Import)
	})

	// Parts - update operations
	r.Group(func(r chi.Router) {
		r.Use(s.writeRateLimitMiddleware())
		r.Use(middleware.RequirePermission(auth.PermPartsUpdate, s.logger))
		r.Patch("/parts/{id}", parts.Update)
	})

	// Parts - delete operations
	r.Group(func(r chi.Router) {
		r.Use(s.writeRateLimitMiddleware())
		r.Use(middleware.RequirePermission(auth.PermPartsDelete, s.logger))
		r.Delete("/parts/{id}", parts.Delete)
	})

	// Parts - export (read permission)
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequirePermission(auth.PermPartsRead, s.logger))
		r.Post("/parts/export", parts.Export)
	})

	// Inventory - read operations
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequirePermission(auth.PermInventoryRead, s.logger))
		r.Get("/inventory", inv.List)
	})

	// Inventory - create/update operations
	r.Group(func(r chi.Router) {
		r.Use(s.writeRateLimitMiddleware())
		r.Use(middleware.RequireAnyPermission(s.logger, auth.PermInventoryCreate, auth.PermInventoryUpdate))
		r.Post("/inventory/upsert", inv.Upsert)
	})

	// Warehouse Dashboard - read operations (for warehouse manager dashboard)
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequirePermission(auth.PermInventoryRead, s.logger))
		r.Get("/warehouse/dashboard", whDash.GetDashboardSummary)
		r.Get("/warehouse/low-stock", whDash.GetLowStockItems)
		r.Get("/warehouse/pending-issues", whDash.GetPendingPartIssues)
		r.Get("/warehouse/movements", whDash.GetStockMovements)
	})
}

// newAuditLogger creates a new audit logger for the given audit store.
func newAuditLogger(auditStore *audit.Store) *audit.Logger {
	return audit.NewLogger(auditStore)
}
