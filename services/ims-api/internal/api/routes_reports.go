package api

import (
	"github.com/edvirons/ssp/ims/internal/auth"
	"github.com/edvirons/ssp/ims/internal/handlers"
	"github.com/edvirons/ssp/ims/internal/middleware"
	"github.com/go-chi/chi/v5"
)

// mountReportRoutes registers all report-related routes.
func (s *Server) mountReportRoutes(r chi.Router, rpt *handlers.ReportsHandler) {
	r.Route("/reports", func(r chi.Router) {
		// All report routes require reports:read permission
		r.Use(middleware.RequirePermission(auth.PermReportsRead, s.logger))

		r.Get("/work-orders", rpt.WorkOrdersReport)
		r.Get("/incidents", rpt.IncidentsReport)
		r.Get("/schools", rpt.SchoolsReport)
		r.Get("/executive", rpt.ExecutiveDashboard)

		// Inventory report - accessible with reports:read permission
		r.Get("/inventory", rpt.InventoryReport)
	})
}
