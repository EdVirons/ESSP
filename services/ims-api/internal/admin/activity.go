package admin

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"go.uber.org/zap"
)

// ActivityEvent represents a single activity event
type ActivityEvent struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Action    string                 `json:"action"`
	Actor     string                 `json:"actor"`
	Target    string                 `json:"target"`
	Timestamp string                 `json:"timestamp"`
	Metadata  map[string]interface{} `json:"metadata"`
}

// GetActivityFeed returns recent activity events
func (h *Handler) GetActivityFeed(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	// Parse limit parameter
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
			user_email,
			entity_id,
			created_at
		FROM audit_logs
		ORDER BY created_at DESC
		LIMIT $1
	`, limit)
	if err != nil {
		h.logger.Error("failed to query activity feed", zap.Error(err))
		http.Error(w, "failed to fetch activity", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var events []ActivityEvent
	for rows.Next() {
		var (
			id         string
			entityType string
			action     string
			userEmail  string
			entityID   string
			createdAt  time.Time
		)

		if err := rows.Scan(&id, &entityType, &action, &userEmail, &entityID, &createdAt); err != nil {
			continue
		}

		events = append(events, ActivityEvent{
			ID:        id,
			Type:      entityType,
			Action:    action,
			Actor:     userEmail,
			Target:    entityID,
			Timestamp: createdAt.UTC().Format(time.RFC3339),
			Metadata:  map[string]interface{}{},
		})
	}

	if events == nil {
		events = []ActivityEvent{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}
