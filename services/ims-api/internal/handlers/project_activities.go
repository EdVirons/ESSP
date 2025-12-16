package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/edvirons/ssp/ims/internal/middleware"
	"github.com/edvirons/ssp/ims/internal/models"
	"github.com/edvirons/ssp/ims/internal/store"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type ProjectActivitiesHandler struct {
	log *zap.Logger
	pg  *store.Postgres
}

func NewProjectActivitiesHandler(log *zap.Logger, pg *store.Postgres) *ProjectActivitiesHandler {
	return &ProjectActivitiesHandler{log: log, pg: pg}
}

type createActivityReq struct {
	PhaseID      string   `json:"phaseId"`
	ActivityType string   `json:"activityType"`
	Content      string   `json:"content"`
	Visibility   string   `json:"visibility"`
	Mentions     []string `json:"mentions"`
}

func (h *ProjectActivitiesHandler) CreateActivity(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "id")
	tenant := middleware.TenantID(r.Context())
	actorID := middleware.UserID(r.Context())
	actorName := middleware.UserName(r.Context())

	var req createActivityReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	content := strings.TrimSpace(req.Content)
	if content == "" {
		http.Error(w, "content required", http.StatusBadRequest)
		return
	}

	// Default activity type to comment
	activityType := models.ActivityType(req.ActivityType)
	if activityType == "" {
		activityType = models.ActivityComment
	}
	if !models.IsValidActivityType(string(activityType)) {
		activityType = models.ActivityComment
	}

	// Only allow comment and note types from this endpoint
	if activityType != models.ActivityComment && activityType != models.ActivityNote {
		activityType = models.ActivityComment
	}

	visibility := models.ActivityVisibility(req.Visibility)
	if visibility == "" {
		visibility = models.VisibilityTeam
	}

	metadata := make(map[string]any)
	if len(req.Mentions) > 0 {
		metadata["mentions"] = req.Mentions
	}

	now := time.Now().UTC()
	activity := models.ProjectActivity{
		ID:            store.NewID("act"),
		TenantID:      tenant,
		ProjectID:     projectID,
		PhaseID:       strings.TrimSpace(req.PhaseID),
		ActivityType:  activityType,
		ActorUserID:   actorID,
		ActorName:     actorName,
		Content:       content,
		Metadata:      metadata,
		AttachmentIDs: []string{},
		Visibility:    visibility,
		IsPinned:      false,
		CreatedAt:     now,
	}

	if err := h.pg.ProjectActivities().CreateActivity(r.Context(), activity); err != nil {
		h.log.Error("failed to create activity", zap.Error(err))
		http.Error(w, "failed to create activity", http.StatusInternalServerError)
		return
	}

	// Create notifications for mentions
	if len(req.Mentions) > 0 {
		h.createMentionNotifications(r, projectID, activity.ID, req.Mentions, actorName, content)
	}

	writeJSON(w, http.StatusCreated, activity)
}

type updateActivityReq struct {
	Content string `json:"content"`
}

func (h *ProjectActivitiesHandler) UpdateActivity(w http.ResponseWriter, r *http.Request) {
	activityID := chi.URLParam(r, "activityId")
	tenant := middleware.TenantID(r.Context())
	actorID := middleware.UserID(r.Context())

	var req updateActivityReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	content := strings.TrimSpace(req.Content)
	if content == "" {
		http.Error(w, "content required", http.StatusBadRequest)
		return
	}

	// Get existing activity
	activity, err := h.pg.ProjectActivities().GetActivity(r.Context(), tenant, activityID)
	if err != nil {
		http.Error(w, "activity not found", http.StatusNotFound)
		return
	}

	// Only the author can edit
	if activity.ActorUserID != actorID {
		http.Error(w, "not authorized to edit this activity", http.StatusForbidden)
		return
	}

	// Only comments and notes can be edited
	if activity.ActivityType != models.ActivityComment && activity.ActivityType != models.ActivityNote {
		http.Error(w, "this activity type cannot be edited", http.StatusBadRequest)
		return
	}

	if err := h.pg.ProjectActivities().UpdateActivity(r.Context(), tenant, activityID, content); err != nil {
		h.log.Error("failed to update activity", zap.Error(err))
		http.Error(w, "failed to update activity", http.StatusInternalServerError)
		return
	}

	// Get updated activity
	updated, _ := h.pg.ProjectActivities().GetActivity(r.Context(), tenant, activityID)
	writeJSON(w, http.StatusOK, updated)
}

func (h *ProjectActivitiesHandler) DeleteActivity(w http.ResponseWriter, r *http.Request) {
	activityID := chi.URLParam(r, "activityId")
	tenant := middleware.TenantID(r.Context())
	actorID := middleware.UserID(r.Context())

	// Get existing activity
	activity, err := h.pg.ProjectActivities().GetActivity(r.Context(), tenant, activityID)
	if err != nil {
		http.Error(w, "activity not found", http.StatusNotFound)
		return
	}

	// Only the author can delete (or admin - could add role check here)
	if activity.ActorUserID != actorID {
		http.Error(w, "not authorized to delete this activity", http.StatusForbidden)
		return
	}

	if err := h.pg.ProjectActivities().DeleteActivity(r.Context(), tenant, activityID); err != nil {
		h.log.Error("failed to delete activity", zap.Error(err))
		http.Error(w, "failed to delete activity", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ProjectActivitiesHandler) TogglePin(w http.ResponseWriter, r *http.Request) {
	activityID := chi.URLParam(r, "activityId")
	tenant := middleware.TenantID(r.Context())

	// Get existing activity
	activity, err := h.pg.ProjectActivities().GetActivity(r.Context(), tenant, activityID)
	if err != nil {
		http.Error(w, "activity not found", http.StatusNotFound)
		return
	}

	// Toggle pin status
	newPinned := !activity.IsPinned
	if err := h.pg.ProjectActivities().TogglePin(r.Context(), tenant, activityID, newPinned); err != nil {
		h.log.Error("failed to toggle pin", zap.Error(err))
		http.Error(w, "failed to toggle pin", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"id":       activityID,
		"isPinned": newPinned,
	})
}

func (h *ProjectActivitiesHandler) ListActivities(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "id")
	tenant := middleware.TenantID(r.Context())

	phaseID := strings.TrimSpace(r.URL.Query().Get("phaseId"))
	activityType := strings.TrimSpace(r.URL.Query().Get("type"))
	actorUserID := strings.TrimSpace(r.URL.Query().Get("userId"))
	pinnedOnly := r.URL.Query().Get("pinned") == "true"
	limit := parseLimit(r.URL.Query().Get("limit"), 50, 200)
	curT, curID, hasCur := decodeCursor(strings.TrimSpace(r.URL.Query().Get("cursor")))

	activities, next, err := h.pg.ProjectActivities().ListActivities(r.Context(), store.ActivityListParams{
		TenantID:     tenant,
		ProjectID:    projectID,
		PhaseID:      phaseID,
		ActivityType: activityType,
		ActorUserID:  actorUserID,
		PinnedOnly:   pinnedOnly,
		Limit:        limit,
		HasCursor:    hasCur,
		CursorTime:   curT,
		CursorID:     curID,
	})
	if err != nil {
		h.log.Error("failed to list activities", zap.Error(err))
		http.Error(w, "failed to list activities", http.StatusInternalServerError)
		return
	}
	if activities == nil {
		activities = []models.ProjectActivity{}
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"items":      activities,
		"nextCursor": next,
	})
}

// Attachments

func (h *ProjectActivitiesHandler) ListAttachments(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "id")
	tenant := middleware.TenantID(r.Context())
	limit := parseLimit(r.URL.Query().Get("limit"), 50, 200)

	attachments, err := h.pg.ProjectActivities().ListAttachments(r.Context(), tenant, projectID, limit)
	if err != nil {
		h.log.Error("failed to list attachments", zap.Error(err))
		http.Error(w, "failed to list attachments", http.StatusInternalServerError)
		return
	}
	if attachments == nil {
		attachments = []models.ProjectAttachment{}
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"items": attachments,
		"total": len(attachments),
	})
}

func (h *ProjectActivitiesHandler) DeleteAttachment(w http.ResponseWriter, r *http.Request) {
	attachmentID := chi.URLParam(r, "attachmentId")
	tenant := middleware.TenantID(r.Context())

	if err := h.pg.ProjectActivities().DeleteAttachment(r.Context(), tenant, attachmentID); err != nil {
		h.log.Error("failed to delete attachment", zap.Error(err))
		http.Error(w, "failed to delete attachment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Helper to create notifications for mentions
func (h *ProjectActivitiesHandler) createMentionNotifications(r *http.Request, projectID, activityID string, userIDs []string, actorName, content string) {
	if len(userIDs) == 0 {
		return
	}

	tenant := middleware.TenantID(r.Context())
	now := time.Now().UTC()

	// Truncate content for notification
	preview := content
	if len(preview) > 100 {
		preview = preview[:100] + "..."
	}

	notifications := make([]models.UserNotification, len(userIDs))
	for i, userID := range userIDs {
		notifications[i] = models.UserNotification{
			ID:               store.NewID("ntf"),
			TenantID:         tenant,
			UserID:           userID,
			NotificationType: models.ProjectNotificationMention,
			EntityType:       "project",
			EntityID:         activityID,
			ProjectID:        projectID,
			Title:            actorName + " mentioned you",
			Body:             preview,
			Metadata: map[string]any{
				"activityId": activityID,
				"projectId":  projectID,
			},
			IsRead:    false,
			CreatedAt: now,
		}
	}

	if err := h.pg.UserNotifications().CreateBulkNotifications(r.Context(), notifications); err != nil {
		h.log.Warn("failed to create mention notifications", zap.Error(err))
	}
}

// LogStatusChange creates a status change activity (called from phases handler)
func (h *ProjectActivitiesHandler) LogStatusChange(ctx context.Context, tenant, projectID, phaseID, actorID, actorName, from, to, reason string) {
	now := time.Now().UTC()
	activity := models.ProjectActivity{
		ID:           store.NewID("act"),
		TenantID:     tenant,
		ProjectID:    projectID,
		PhaseID:      phaseID,
		ActivityType: models.ActivityStatusChange,
		ActorUserID:  actorID,
		ActorName:    actorName,
		Content:      "",
		Metadata: map[string]any{
			"from":   from,
			"to":     to,
			"reason": reason,
		},
		AttachmentIDs: []string{},
		Visibility:    models.VisibilityTeam,
		IsPinned:      false,
		CreatedAt:     now,
	}

	if err := h.pg.ProjectActivities().CreateActivity(ctx, activity); err != nil {
		h.log.Warn("failed to log status change activity", zap.Error(err))
	}
}

// LogPhaseTransition creates a phase transition activity
func (h *ProjectActivitiesHandler) LogPhaseTransition(ctx context.Context, tenant, projectID, phaseID, phaseType, actorID, actorName, from, to string) {
	now := time.Now().UTC()
	activity := models.ProjectActivity{
		ID:           store.NewID("act"),
		TenantID:     tenant,
		ProjectID:    projectID,
		PhaseID:      phaseID,
		ActivityType: models.ActivityPhaseTransition,
		ActorUserID:  actorID,
		ActorName:    actorName,
		Content:      "",
		Metadata: map[string]any{
			"phaseId":     phaseID,
			"phaseType":   phaseType,
			"from":        from,
			"to":          to,
			"completedBy": actorName,
		},
		AttachmentIDs: []string{},
		Visibility:    models.VisibilityTeam,
		IsPinned:      false,
		CreatedAt:     now,
	}

	if err := h.pg.ProjectActivities().CreateActivity(ctx, activity); err != nil {
		h.log.Warn("failed to log phase transition activity", zap.Error(err))
	}
}

// LogWorkOrderActivity creates a work order activity
func (h *ProjectActivitiesHandler) LogWorkOrderActivity(ctx context.Context, tenant, projectID, phaseID, workOrderID, actorID, actorName, action, status string) {
	now := time.Now().UTC()
	activity := models.ProjectActivity{
		ID:           store.NewID("act"),
		TenantID:     tenant,
		ProjectID:    projectID,
		PhaseID:      phaseID,
		WorkOrderID:  workOrderID,
		ActivityType: models.ActivityWorkOrder,
		ActorUserID:  actorID,
		ActorName:    actorName,
		Content:      "",
		Metadata: map[string]any{
			"workOrderId": workOrderID,
			"action":      action,
			"status":      status,
		},
		AttachmentIDs: []string{},
		Visibility:    models.VisibilityTeam,
		IsPinned:      false,
		CreatedAt:     now,
	}

	if err := h.pg.ProjectActivities().CreateActivity(ctx, activity); err != nil {
		h.log.Warn("failed to log work order activity", zap.Error(err))
	}
}
