package api

import (
	"github.com/edvirons/ssp/ims/internal/auth"
	"github.com/edvirons/ssp/ims/internal/handlers"
	"github.com/edvirons/ssp/ims/internal/middleware"
	"github.com/go-chi/chi/v5"
)

// mountMarketingKBRoutes registers all marketing knowledge base routes.
func (s *Server) mountMarketingKBRoutes(r chi.Router, mkb *handlers.MarketingKBHandler) {
	// Marketing KB Articles - read operations
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequirePermission(auth.PermMKBRead, s.logger))
		r.Get("/marketing-kb/articles", mkb.ListArticles)
		r.Get("/marketing-kb/articles/{id}", mkb.GetArticleByID)
		r.Get("/marketing-kb/articles/slug/{slug}", mkb.GetArticleBySlug)
		r.Get("/marketing-kb/search", mkb.SearchArticles)
		r.Get("/marketing-kb/stats", mkb.GetArticleStats)
		r.Get("/marketing-kb/pitch-kits", mkb.ListPitchKits)
		r.Get("/marketing-kb/pitch-kits/{id}", mkb.GetPitchKitByID)
	})

	// Marketing KB Articles - create operations
	r.Group(func(r chi.Router) {
		r.Use(s.writeRateLimitMiddleware())
		r.Use(middleware.RequirePermission(auth.PermMKBCreate, s.logger))
		r.Post("/marketing-kb/articles", mkb.CreateArticle)
		r.Post("/marketing-kb/pitch-kits", mkb.CreatePitchKit)
	})

	// Marketing KB Articles - update operations
	r.Group(func(r chi.Router) {
		r.Use(s.writeRateLimitMiddleware())
		r.Use(middleware.RequirePermission(auth.PermMKBUpdate, s.logger))
		r.Patch("/marketing-kb/articles/{id}", mkb.UpdateArticle)
		r.Post("/marketing-kb/articles/{id}/submit-review", mkb.SubmitForReview)
		r.Post("/marketing-kb/articles/{id}/usage", mkb.RecordUsage)
		r.Patch("/marketing-kb/pitch-kits/{id}", mkb.UpdatePitchKit)
	})

	// Marketing KB Articles - approve operations (admin only)
	r.Group(func(r chi.Router) {
		r.Use(s.writeRateLimitMiddleware())
		r.Use(middleware.RequirePermission(auth.PermMKBApprove, s.logger))
		r.Post("/marketing-kb/articles/{id}/approve", mkb.ApproveArticle)
	})

	// Marketing KB Articles - delete operations (admin only)
	r.Group(func(r chi.Router) {
		r.Use(s.writeRateLimitMiddleware())
		r.Use(middleware.RequirePermission(auth.PermMKBDelete, s.logger))
		r.Delete("/marketing-kb/articles/{id}", mkb.DeleteArticle)
		r.Delete("/marketing-kb/pitch-kits/{id}", mkb.DeletePitchKit)
	})
}
