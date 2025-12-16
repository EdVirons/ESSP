package api

import (
	"github.com/edvirons/ssp/ims/internal/config"
	"github.com/edvirons/ssp/ims/internal/handlers"
	"github.com/edvirons/ssp/ims/internal/middleware"
	"github.com/edvirons/ssp/ims/internal/store"
	"github.com/edvirons/ssp/ims/internal/ws"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// RegisterAIChatRoutes registers the AI chat API routes
func RegisterAIChatRoutes(r chi.Router, log *zap.Logger, pg *store.Postgres, hub *ws.Hub, cfg config.Config) {
	aiHandler := handlers.NewAIChatHandler(log, pg, hub, cfg)
	livechatHandler := handlers.NewLivechatHandler(log, pg, hub)

	r.Route("/chat", func(r chi.Router) {
		// Session management (school contacts)
		r.Post("/sessions", livechatHandler.StartSession)
		r.Get("/sessions/{id}/queue", livechatHandler.GetQueuePosition)
		r.Post("/sessions/{id}/end", livechatHandler.EndSession)

		// AI Chat routes - school contacts can send messages to AI
		r.Route("/ai", func(r chi.Router) {
			r.Post("/sessions/{id}/message", aiHandler.HandleAIMessage)
			r.Post("/sessions/{id}/escalate", aiHandler.RequestEscalation)

			// Agent-only: get AI conversation context for handoff
			r.Group(func(r chi.Router) {
				r.Use(middleware.RequireRoles("ssp_admin", "ssp_support_agent"))
				r.Get("/sessions/{id}/context", aiHandler.GetConversationContext)
			})
		})

		// Agent operations
		r.Group(func(r chi.Router) {
			r.Use(middleware.RequireRoles("ssp_admin", "ssp_support_agent"))

			r.Post("/accept", livechatHandler.AcceptChat)
			r.Post("/sessions/{id}/transfer", livechatHandler.TransferChat)
			r.Get("/queue", livechatHandler.GetQueue)
			r.Get("/active", livechatHandler.GetActiveChats)

			// Availability
			r.Put("/availability", livechatHandler.SetAvailability)
			r.Get("/availability", livechatHandler.GetAvailability)
		})

		// Admin only - metrics
		r.Group(func(r chi.Router) {
			r.Use(middleware.RequireRoles("ssp_admin"))
			r.Get("/metrics", livechatHandler.GetChatMetrics)
		})
	})
}
