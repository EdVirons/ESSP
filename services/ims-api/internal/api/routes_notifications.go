package api

import (
	"github.com/edvirons/ssp/ims/internal/auth"
	"github.com/edvirons/ssp/ims/internal/handlers"
	"github.com/edvirons/ssp/ims/internal/middleware"
	"github.com/go-chi/chi/v5"
)

// mountNotificationRoutes registers user notification routes.
func (s *Server) mountNotificationRoutes(r chi.Router, notifications *handlers.UserNotificationsHandler) {
	// Notification read operations
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequirePermission(auth.PermNotificationRead, s.logger))
		r.Get("/users/me/notifications", notifications.ListNotifications)
		r.Get("/users/me/notifications/unread-count", notifications.GetUnreadCount)
	})

	// Notification update operations
	r.Group(func(r chi.Router) {
		r.Use(s.writeRateLimitMiddleware())
		r.Use(middleware.RequirePermission(auth.PermNotificationUpdate, s.logger))
		r.Post("/users/me/notifications/{notificationId}/read", notifications.MarkAsRead)
		r.Post("/users/me/notifications/read-all", notifications.MarkAllAsRead)
		r.Post("/users/me/notifications/read", notifications.MarkMultipleAsRead)
	})
}
