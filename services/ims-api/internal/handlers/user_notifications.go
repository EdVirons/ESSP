package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/edvirons/ssp/ims/internal/middleware"
	"github.com/edvirons/ssp/ims/internal/models"
	"github.com/edvirons/ssp/ims/internal/store"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type UserNotificationsHandler struct {
	log *zap.Logger
	pg  *store.Postgres
}

func NewUserNotificationsHandler(log *zap.Logger, pg *store.Postgres) *UserNotificationsHandler {
	return &UserNotificationsHandler{log: log, pg: pg}
}

func (h *UserNotificationsHandler) ListNotifications(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())
	userID := middleware.UserID(r.Context())

	unreadOnly := r.URL.Query().Get("unread") == "true"
	projectID := strings.TrimSpace(r.URL.Query().Get("projectId"))
	limit := parseLimit(r.URL.Query().Get("limit"), 50, 200)
	curT, curID, hasCur := decodeCursor(strings.TrimSpace(r.URL.Query().Get("cursor")))

	notifications, next, err := h.pg.UserNotifications().ListNotifications(r.Context(), store.NotificationListParams{
		TenantID:   tenant,
		UserID:     userID,
		UnreadOnly: unreadOnly,
		ProjectID:  projectID,
		Limit:      limit,
		HasCursor:  hasCur,
		CursorTime: curT,
		CursorID:   curID,
	})
	if err != nil {
		h.log.Error("failed to list notifications", zap.Error(err))
		http.Error(w, "failed to list notifications", http.StatusInternalServerError)
		return
	}
	if notifications == nil {
		notifications = []models.UserNotification{}
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"items":      notifications,
		"nextCursor": next,
	})
}

func (h *UserNotificationsHandler) GetUnreadCount(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())
	userID := middleware.UserID(r.Context())

	count, err := h.pg.UserNotifications().GetUnreadCount(r.Context(), tenant, userID)
	if err != nil {
		h.log.Error("failed to get unread count", zap.Error(err))
		http.Error(w, "failed to get unread count", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"unreadCount": count,
	})
}

func (h *UserNotificationsHandler) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	notificationID := chi.URLParam(r, "notificationId")
	tenant := middleware.TenantID(r.Context())
	userID := middleware.UserID(r.Context())

	if err := h.pg.UserNotifications().MarkAsRead(r.Context(), tenant, userID, notificationID); err != nil {
		h.log.Error("failed to mark notification as read", zap.Error(err))
		http.Error(w, "failed to mark as read", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"id":     notificationID,
		"isRead": true,
	})
}

func (h *UserNotificationsHandler) MarkAllAsRead(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())
	userID := middleware.UserID(r.Context())

	count, err := h.pg.UserNotifications().MarkAllAsRead(r.Context(), tenant, userID)
	if err != nil {
		h.log.Error("failed to mark all notifications as read", zap.Error(err))
		http.Error(w, "failed to mark all as read", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"markedCount": count,
	})
}

type markMultipleReadReq struct {
	NotificationIDs []string `json:"notificationIds"`
}

func (h *UserNotificationsHandler) MarkMultipleAsRead(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())
	userID := middleware.UserID(r.Context())

	var req markMultipleReadReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if len(req.NotificationIDs) == 0 {
		http.Error(w, "notificationIds required", http.StatusBadRequest)
		return
	}

	count, err := h.pg.UserNotifications().MarkMultipleAsRead(r.Context(), tenant, userID, req.NotificationIDs)
	if err != nil {
		h.log.Error("failed to mark notifications as read", zap.Error(err))
		http.Error(w, "failed to mark as read", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"markedCount": count,
	})
}
