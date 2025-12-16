package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/edvirons/ssp/ims/internal/audit"
	"github.com/edvirons/ssp/ims/internal/middleware"
	"github.com/edvirons/ssp/ims/internal/models"
	"github.com/edvirons/ssp/ims/internal/store"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// WorkOrderReworkHandler handles work order rejection/rework operations.
type WorkOrderReworkHandler struct {
	log   *zap.Logger
	pg    *store.Postgres
	audit audit.AuditLogger
}

// NewWorkOrderReworkHandler creates a new rework handler.
func NewWorkOrderReworkHandler(log *zap.Logger, pg *store.Postgres, auditLogger audit.AuditLogger) *WorkOrderReworkHandler {
	return &WorkOrderReworkHandler{log: log, pg: pg, audit: auditLogger}
}

// Reject handles POST /work-orders/{id}/reject
func (h *WorkOrderReworkHandler) Reject(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())
	school := middleware.SchoolID(r.Context())
	id := chi.URLParam(r, "id")
	userID := middleware.UserID(r.Context())
	userName := middleware.UserName(r.Context())

	// Check feature flag
	enabled, err := h.pg.FeatureConfig().IsFeatureEnabled(r.Context(), tenant, models.FeatureWorkOrderRework)
	if err != nil {
		h.log.Error("failed to check feature flag", zap.Error(err))
	}
	if !enabled {
		http.Error(w, "work order rework feature is disabled", http.StatusForbidden)
		return
	}

	// Get rework configuration
	cfg, err := h.pg.FeatureConfig().GetReworkConfig(r.Context(), tenant)
	if err != nil {
		h.log.Error("failed to get rework config", zap.Error(err))
		cfg = models.DefaultReworkConfig()
	}

	// Parse request
	var req models.RejectWorkOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	// Validate reason if required
	if cfg.RequireReason && req.Reason == "" {
		http.Error(w, "rejection reason is required", http.StatusBadRequest)
		return
	}

	// Validate category
	if req.Category != "" && !models.IsValidRejectionCategory(string(req.Category)) {
		http.Error(w, "invalid rejection category", http.StatusBadRequest)
		return
	}
	if req.Category == "" {
		req.Category = models.RejectionOther
	}

	// Get current work order
	current, err := h.pg.WorkOrders().GetByID(r.Context(), tenant, school, id)
	if err != nil {
		http.Error(w, "work order not found", http.StatusNotFound)
		return
	}

	// Validate transition
	if !models.CanReworkTo(current.Status, req.TargetStatus) {
		http.Error(w, "invalid rework transition from "+string(current.Status)+" to "+string(req.TargetStatus), http.StatusBadRequest)
		return
	}

	// Check max rework count
	if current.ReworkCount >= cfg.MaxReworkCount {
		http.Error(w, "maximum rework count exceeded", http.StatusBadRequest)
		return
	}

	now := time.Now().UTC()

	// Get next sequence number
	nextSeq, err := h.pg.WorkOrderRework().GetNextReworkSequence(r.Context(), tenant, school, id)
	if err != nil {
		h.log.Error("failed to get rework sequence", zap.Error(err))
		nextSeq = current.ReworkCount + 1
	}

	// Create rework history entry
	history := models.WorkOrderReworkHistory{
		ID:                store.NewID("rw"),
		TenantID:          tenant,
		SchoolID:          school,
		WorkOrderID:       id,
		FromStatus:        current.Status,
		ToStatus:          req.TargetStatus,
		RejectionReason:   req.Reason,
		RejectionCategory: req.Category,
		RejectedByUserID:  userID,
		RejectedByName:    userName,
		ReworkSequence:    nextSeq,
		CreatedAt:         now,
	}

	if err := h.pg.WorkOrderRework().Create(r.Context(), history); err != nil {
		h.log.Error("failed to create rework history", zap.Error(err))
		http.Error(w, "failed to record rework history", http.StatusInternalServerError)
		return
	}

	// Update work order status
	updated, err := h.pg.WorkOrders().UpdateStatusWithRework(r.Context(), tenant, school, id, req.TargetStatus, req.Reason, now)
	if err != nil {
		h.log.Error("failed to update work order", zap.Error(err))
		http.Error(w, "failed to update work order status", http.StatusInternalServerError)
		return
	}

	// Audit log
	if err := h.audit.LogUpdate(r.Context(), "work_order", id,
		map[string]any{"status": current.Status},
		map[string]any{
			"action":      "reject",
			"status":      req.TargetStatus,
			"reason":      req.Reason,
			"category":    req.Category,
			"reworkCount": updated.ReworkCount,
			"userId":      userID,
			"userName":    userName,
		}); err != nil {
		h.log.Error("failed to log rework audit", zap.Error(err))
	}

	// Return response
	resp := models.RejectWorkOrderResponse{
		WorkOrder:     updated,
		ReworkHistory: history,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GetReworkHistory handles GET /work-orders/{id}/rework-history
func (h *WorkOrderReworkHandler) GetReworkHistory(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())
	school := middleware.SchoolID(r.Context())
	id := chi.URLParam(r, "id")

	limit := 50
	cursor := r.URL.Query().Get("cursor")

	params := store.ReworkHistoryListParams{
		TenantID:    tenant,
		SchoolID:    school,
		WorkOrderID: id,
		Limit:       limit,
	}

	if cursor != "" {
		t, cid, ok := store.DecodeCursor(cursor)
		if ok {
			params.HasCursor = true
			params.CursorTime = t
			params.CursorID = cid
		}
	}

	history, nextCursor, err := h.pg.WorkOrderRework().List(r.Context(), params)
	if err != nil {
		h.log.Error("failed to list rework history", zap.Error(err))
		http.Error(w, "failed to fetch rework history", http.StatusInternalServerError)
		return
	}

	resp := struct {
		History    []models.WorkOrderReworkHistory `json:"history"`
		NextCursor string                          `json:"nextCursor,omitempty"`
	}{
		History:    history,
		NextCursor: nextCursor,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
