package admin

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"go.uber.org/zap"
)

// Notification represents a notification for the dashboard
type Notification struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Action    string                 `json:"action"`
	Actor     string                 `json:"actor"`
	Target    string                 `json:"target"`
	Summary   string                 `json:"summary"`
	Timestamp string                 `json:"timestamp"`
	Metadata  map[string]interface{} `json:"metadata"`
	Read      bool                   `json:"read"`
}

// NotificationsResponse is the response for the notifications list
type NotificationsResponse struct {
	Items       []Notification `json:"items"`
	UnreadCount int            `json:"unreadCount"`
	Total       int            `json:"total"`
}

// UnreadCountResponse is the response for unread count
type UnreadCountResponse struct {
	Count int `json:"count"`
}

// MarkReadRequest is the request body for marking notifications as read
type MarkReadRequest struct {
	IDs string `json:"ids"` // comma-separated IDs or "all"
}

// GetNotifications returns recent notifications from audit logs
func (h *Handler) GetNotifications(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	// Parse query parameters
	limitStr := r.URL.Query().Get("limit")
	limit := 50
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 200 {
			limit = l
		}
	}

	pool := h.pg.RawPool()

	// Query recent audit logs
	rows, err := pool.Query(ctx, `
		SELECT
			id,
			entity_type,
			action,
			COALESCE(actor_name, 'System'),
			entity_id,
			created_at
		FROM audit_logs
		ORDER BY created_at DESC
		LIMIT $1
	`, limit)
	if err != nil {
		h.logger.Error("failed to query notifications", zap.Error(err))
		http.Error(w, "failed to fetch notifications", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var notifications []Notification
	for rows.Next() {
		var (
			id         string
			entityType string
			action     string
			actorName  string
			entityID   string
			createdAt  time.Time
		)

		if err := rows.Scan(&id, &entityType, &action, &actorName, &entityID, &createdAt); err != nil {
			continue
		}

		// Generate a human-readable summary
		summary := generateSummary(entityType, action, entityID)

		notifications = append(notifications, Notification{
			ID:        id,
			Type:      entityType,
			Action:    action,
			Actor:     actorName,
			Target:    entityID,
			Summary:   summary,
			Timestamp: createdAt.UTC().Format(time.RFC3339),
			Metadata:  map[string]interface{}{},
			Read:      false, // For simplicity, all notifications are unread initially
		})
	}

	if notifications == nil {
		notifications = []Notification{}
	}

	response := NotificationsResponse{
		Items:       notifications,
		UnreadCount: len(notifications), // Simplified: all are unread
		Total:       len(notifications),
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}

// GetUnreadCount returns the count of unread notifications
func (h *Handler) GetUnreadCount(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	pool := h.pg.RawPool()

	// Count recent audit logs (simplified: count last 24 hours)
	var count int
	err := pool.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM audit_logs
		WHERE created_at > NOW() - INTERVAL '24 hours'
	`).Scan(&count)
	if err != nil {
		h.logger.Error("failed to count notifications", zap.Error(err))
		http.Error(w, "failed to count notifications", http.StatusInternalServerError)
		return
	}

	// Cap at 99 for display
	if count > 99 {
		count = 99
	}

	response := UnreadCountResponse{Count: count}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}

// MarkNotificationsRead marks notifications as read (placeholder - stores in memory/Redis in production)
func (h *Handler) MarkNotificationsRead(w http.ResponseWriter, r *http.Request) {
	// For simplicity, just acknowledge the request
	// In production, store read status in Redis or a separate table
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

// generateSummary creates a human-readable summary for a notification
func generateSummary(entityType, action, entityID string) string {
	actionVerb := "updated"
	switch action {
	case "create":
		actionVerb = "created"
	case "delete":
		actionVerb = "deleted"
	}

	entityLabel := entityType
	switch entityType {
	case "incident":
		entityLabel = "Incident"
	case "work_order":
		entityLabel = "Work Order"
	case "program":
		entityLabel = "Program"
	case "service_shop":
		entityLabel = "Service Shop"
	case "device":
		entityLabel = "Device"
	}

	// Shorten the ID for display
	shortID := entityID
	if len(entityID) > 12 {
		shortID = entityID[:12] + "..."
	}

	return entityLabel + " " + shortID + " was " + actionVerb
}
