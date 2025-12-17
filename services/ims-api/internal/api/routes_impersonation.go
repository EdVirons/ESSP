package api

import (
	"github.com/edvirons/ssp/ims/internal/auth"
	"github.com/edvirons/ssp/ims/internal/handlers"
	"github.com/edvirons/ssp/ims/internal/middleware"
	"github.com/go-chi/chi/v5"
)

// mountImpersonationRoutes registers impersonation-related routes.
// These routes require the PermImpersonate permission (ops managers and admins only).
func (s *Server) mountImpersonationRoutes(r chi.Router, imp *handlers.ImpersonationHandler) {
	r.Route("/impersonate", func(r chi.Router) {
		// All impersonation routes require PermImpersonate
		r.Use(middleware.RequirePermission(auth.PermImpersonate, s.logger))

		// List users that can be impersonated
		r.Get("/users", imp.ListImpersonatableUsers)

		// Validate impersonation target before starting session
		r.Post("/validate", imp.ValidateImpersonation)
	})
}
