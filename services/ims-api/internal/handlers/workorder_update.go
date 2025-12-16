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

// WorkOrderUpdateHandler handles PATCH updates to work orders.
type WorkOrderUpdateHandler struct {
	log   *zap.Logger
	pg    *store.Postgres
	audit audit.AuditLogger
}

// NewWorkOrderUpdateHandler creates a new update handler.
func NewWorkOrderUpdateHandler(log *zap.Logger, pg *store.Postgres, auditLogger audit.AuditLogger) *WorkOrderUpdateHandler {
	return &WorkOrderUpdateHandler{log: log, pg: pg, audit: auditLogger}
}

// Update handles PATCH /work-orders/{id}
func (h *WorkOrderUpdateHandler) Update(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())
	school := middleware.SchoolID(r.Context())
	id := chi.URLParam(r, "id")
	userID := middleware.UserID(r.Context())
	userName := middleware.UserName(r.Context())

	// Check feature flag
	enabled, err := h.pg.FeatureConfig().IsFeatureEnabled(r.Context(), tenant, models.FeatureWorkOrderUpdate)
	if err != nil {
		h.log.Error("failed to check feature flag", zap.Error(err))
	}
	if !enabled {
		http.Error(w, "work order update feature is disabled", http.StatusForbidden)
		return
	}

	// Parse request
	var req models.WorkOrderUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if !req.HasUpdates() {
		http.Error(w, "no fields to update", http.StatusBadRequest)
		return
	}

	// Get current state for audit
	current, err := h.pg.WorkOrders().GetByID(r.Context(), tenant, school, id)
	if err != nil {
		http.Error(w, "work order not found", http.StatusNotFound)
		return
	}

	// Build previous values for audit
	previousValues := make(map[string]any)
	updatedFields := []string{}
	updates := req.ToMap()

	for field := range updates {
		switch field {
		case "assigned_staff_id":
			previousValues[field] = current.AssignedStaffID
			updatedFields = append(updatedFields, "assignedStaffId")
		case "service_shop_id":
			previousValues[field] = current.ServiceShopID
			updatedFields = append(updatedFields, "serviceShopId")
		case "cost_estimate_cents":
			previousValues[field] = current.CostEstimateCents
			updatedFields = append(updatedFields, "costEstimateCents")
		case "notes":
			previousValues[field] = current.Notes
			updatedFields = append(updatedFields, "notes")
		case "repair_location":
			previousValues[field] = current.RepairLocation
			updatedFields = append(updatedFields, "repairLocation")
		case "onsite_contact_id":
			previousValues[field] = current.OnsiteContactID
			updatedFields = append(updatedFields, "onsiteContactId")
		}
	}

	// Perform update
	now := time.Now().UTC()
	updated, err := h.pg.WorkOrders().UpdateFields(r.Context(), tenant, school, id, updates, now)
	if err != nil {
		h.log.Error("failed to update work order", zap.Error(err), zap.String("id", id))
		http.Error(w, "failed to update work order", http.StatusInternalServerError)
		return
	}

	// Audit log
	if err := h.audit.LogUpdate(r.Context(), "work_order", id,
		previousValues,
		map[string]any{
			"fields":   updatedFields,
			"new":      updates,
			"userId":   userID,
			"userName": userName,
		}); err != nil {
		h.log.Error("failed to log update audit", zap.Error(err))
	}

	// Return response
	resp := models.WorkOrderUpdateResponse{
		WorkOrder:      updated,
		UpdatedFields:  updatedFields,
		PreviousValues: previousValues,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
