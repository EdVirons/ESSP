package api

import (
	"github.com/edvirons/ssp/ims/internal/auth"
	"github.com/edvirons/ssp/ims/internal/handlers"
	"github.com/edvirons/ssp/ims/internal/middleware"
	"github.com/go-chi/chi/v5"
)

// mountKBRoutes registers all knowledge base routes.
func (s *Server) mountKBRoutes(r chi.Router, kb *handlers.KBArticlesHandler) {
	// KB Articles - read operations (available to support roles)
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequirePermission(auth.PermKBRead, s.logger))
		r.Get("/kb/articles", kb.List)
		r.Get("/kb/articles/{id}", kb.GetByID)
		r.Get("/kb/articles/slug/{slug}", kb.GetBySlug)
		r.Get("/kb/search", kb.Search)
		r.Get("/kb/stats", kb.GetStats)
	})

	// KB Articles - create operations (admin only)
	r.Group(func(r chi.Router) {
		r.Use(s.writeRateLimitMiddleware())
		r.Use(middleware.RequirePermission(auth.PermKBCreate, s.logger))
		r.Post("/kb/articles", kb.Create)
	})

	// KB Articles - update operations (admin only)
	r.Group(func(r chi.Router) {
		r.Use(s.writeRateLimitMiddleware())
		r.Use(middleware.RequirePermission(auth.PermKBUpdate, s.logger))
		r.Patch("/kb/articles/{id}", kb.Update)
		r.Post("/kb/articles/{id}/publish", kb.Publish)
	})

	// KB Articles - delete operations (admin only)
	r.Group(func(r chi.Router) {
		r.Use(s.writeRateLimitMiddleware())
		r.Use(middleware.RequirePermission(auth.PermKBDelete, s.logger))
		r.Delete("/kb/articles/{id}", kb.Delete)
	})
}
