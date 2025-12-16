package api

import (
	"github.com/edvirons/ssp/ims/internal/auth"
	"github.com/edvirons/ssp/ims/internal/handlers"
	"github.com/edvirons/ssp/ims/internal/middleware"
	"github.com/go-chi/chi/v5"
)

// mountBOMRoutes registers all BOM (Bill of Materials) related routes.
func (s *Server) mountBOMRoutes(r chi.Router, bom *handlers.BOMHandler) {
	// BOM - read operations
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequirePermission(auth.PermBOMRead, s.logger))
		r.Get("/work-orders/{id}/bom", bom.List)
		r.Get("/work-orders/{id}/bom/suggest", bom.Suggest)
	})

	// BOM - create/update operations
	r.Group(func(r chi.Router) {
		r.Use(s.writeRateLimitMiddleware())
		r.Use(middleware.RequireAnyPermission(s.logger, auth.PermBOMCreate, auth.PermBOMUpdate))
		r.Post("/work-orders/{id}/bom/items", bom.AddItem)
	})

	// BOM - consume operations
	r.Group(func(r chi.Router) {
		r.Use(s.writeRateLimitMiddleware())
		r.Use(middleware.RequirePermission(auth.PermBOMConsume, s.logger))
		r.Patch("/work-orders/{id}/bom/items/{itemId}/consume", bom.Consume)
		r.Patch("/work-orders/{id}/bom/items/{itemId}/release", bom.Release)
	})
}
