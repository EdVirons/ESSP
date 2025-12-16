package api

import (
	"github.com/edvirons/ssp/ims/internal/handlers"
	"github.com/edvirons/ssp/ims/internal/middleware"
	"github.com/edvirons/ssp/ims/internal/store"
	"github.com/edvirons/ssp/ims/internal/ws"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// RegisterLivechatRoutes registers the livechat API routes
func RegisterLivechatRoutes(r chi.Router, log *zap.Logger, pg *store.Postgres, hub *ws.Hub) {
	h := handlers.NewLivechatHandler(log, pg, hub)

	r.Route("/chat", func(r chi.Router) {
		// Session management
		r.Post("/sessions", h.StartSession)
		r.Get("/sessions/{id}/queue", h.GetQueuePosition)
		r.Post("/sessions/{id}/end", h.EndSession)

		// Agent operations
		r.Group(func(r chi.Router) {
			r.Use(middleware.RequireRoles("ssp_admin", "ssp_support_agent"))

			r.Post("/accept", h.AcceptChat)
			r.Post("/sessions/{id}/transfer", h.TransferChat)
			r.Get("/queue", h.GetQueue)
			r.Get("/active", h.GetActiveChats)

			// Availability
			r.Put("/availability", h.SetAvailability)
			r.Get("/availability", h.GetAvailability)
		})

		// Admin only - metrics
		r.Group(func(r chi.Router) {
			r.Use(middleware.RequireRoles("ssp_admin"))
			r.Get("/metrics", h.GetChatMetrics)
		})
	})
}
