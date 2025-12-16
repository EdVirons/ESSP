package api

import (
	"github.com/edvirons/ssp/ims/internal/auth"
	"github.com/edvirons/ssp/ims/internal/handlers"
	"github.com/edvirons/ssp/ims/internal/middleware"
	"github.com/go-chi/chi/v5"
)

// mountIncidentRoutes registers all incident-related routes.
func (s *Server) mountIncidentRoutes(r chi.Router, inc *handlers.IncidentHandler) {
	// Incidents - read operations
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequirePermission(auth.PermIncidentRead, s.logger))
		r.Get("/incidents/{id}", inc.GetByID)
		r.Get("/incidents", inc.List)
	})

	// Incidents - create operations
	r.Group(func(r chi.Router) {
		r.Use(s.writeRateLimitMiddleware())
		r.Use(middleware.RequirePermission(auth.PermIncidentCreate, s.logger))
		r.Post("/incidents", inc.Create)
	})

	// Incidents - update operations
	r.Group(func(r chi.Router) {
		r.Use(s.writeRateLimitMiddleware())
		r.Use(middleware.RequirePermission(auth.PermIncidentUpdate, s.logger))
		r.Patch("/incidents/{id}/status", inc.UpdateStatus)
	})
}
