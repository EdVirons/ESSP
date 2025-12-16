package api

import (
	"github.com/edvirons/ssp/ims/internal/auth"
	"github.com/edvirons/ssp/ims/internal/handlers"
	"github.com/edvirons/ssp/ims/internal/middleware"
	"github.com/edvirons/ssp/ims/internal/store"
	"github.com/edvirons/ssp/ims/internal/ws"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// RegisterMessagingRoutes registers the messaging API routes
func RegisterMessagingRoutes(r chi.Router, log *zap.Logger, pg *store.Postgres, hub *ws.Hub) {
	h := handlers.NewMessagingHandler(log, pg, hub)

	r.Route("/messages", func(r chi.Router) {
		// Read operations - require messages:read permission
		r.Group(func(r chi.Router) {
			r.Use(middleware.RequirePermission(auth.PermMessagesRead, log))
			r.Get("/threads", h.ListThreads)
			r.Get("/threads/{id}", h.GetThread)
			r.Get("/unread", h.GetUnreadCounts)
			r.Get("/search", h.SearchMessages)
		})

		// Create operations - require messages:create permission
		r.Group(func(r chi.Router) {
			r.Use(middleware.RequirePermission(auth.PermMessagesCreate, log))
			r.Post("/threads", h.CreateThread)
			r.Post("/threads/{id}/messages", h.CreateMessage)
			r.Post("/threads/{id}/read", h.MarkRead)
		})

		// Management operations - require messages:manage permission
		r.Group(func(r chi.Router) {
			r.Use(middleware.RequirePermission(auth.PermMessagesManage, log))
			r.Patch("/threads/{id}/status", h.UpdateThreadStatus)
		})

		// Analytics (admin only)
		r.Group(func(r chi.Router) {
			r.Use(middleware.RequireRoles("ssp_admin"))
			r.Get("/analytics", h.GetAnalytics)
		})
	})
}
