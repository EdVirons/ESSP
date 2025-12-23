package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/edvirons/ssp/ims/internal/audit"
	"github.com/edvirons/ssp/ims/internal/middleware"
	"github.com/edvirons/ssp/ims/internal/models"
	"github.com/edvirons/ssp/ims/internal/store"
	"go.uber.org/zap"
)

// WorkOrderBulkHandler handles bulk work order operations.
type WorkOrderBulkHandler struct {
	log   *zap.Logger
	pg    *store.Postgres
	audit audit.AuditLogger
}

// NewWorkOrderBulkHandler creates a new bulk handler.
func NewWorkOrderBulkHandler(log *zap.Logger, pg *store.Postgres, auditLogger audit.AuditLogger) *WorkOrderBulkHandler {
	return &WorkOrderBulkHandler{log: log, pg: pg, audit: auditLogger}
}

// BulkStatusUpdate handles POST /work-orders/bulk/status
func (h *WorkOrderBulkHandler) BulkStatusUpdate(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())
	school := middleware.SchoolID(r.Context())
	userID := middleware.UserID(r.Context())

	// Check feature flag
	enabled, err := h.pg.FeatureConfig().IsFeatureEnabled(r.Context(), tenant, models.FeatureWorkOrderBulkOperations)
	if err != nil {
		h.log.Error("failed to check feature flag", zap.Error(err))
	}
	if !enabled {
		http.Error(w, "bulk operations feature is disabled", http.StatusForbidden)
		return
	}

	// Get bulk config
	cfg, err := h.pg.FeatureConfig().GetBulkConfig(r.Context(), tenant)
	if err != nil {
		h.log.Error("failed to get bulk config", zap.Error(err))
		cfg = models.DefaultBulkConfig()
	}

	// Parse request
	var req models.BulkStatusUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if err := req.Validate(cfg.MaxBatchSize); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	now := time.Now().UTC()

	// Create operation log
	opLog := models.BulkOperationLog{
		ID:            store.NewID("bulk"),
		TenantID:      tenant,
		UserID:        userID,
		OperationType: models.BulkOpStatusUpdate,
		EntityType:    "work_order",
		RequestedIDs:  req.WorkOrderIDs,
		SuccessfulIDs: []string{},
		FailedIDs:     []string{},
		StartedAt:     now,
		TotalCount:    len(req.WorkOrderIDs),
		CreatedAt:     now,
	}

	if err := h.pg.BulkOperations().Create(r.Context(), opLog); err != nil {
		h.log.Error("failed to create bulk operation log", zap.Error(err))
	}

	// Validate all work orders exist and can transition
	workOrders, err := h.pg.WorkOrders().GetByIDs(r.Context(), tenant, school, req.WorkOrderIDs)
	if err != nil {
		h.log.Error("failed to fetch work orders", zap.Error(err))
		http.Error(w, "failed to fetch work orders", http.StatusInternalServerError)
		return
	}

	// Build ID map for validation
	woMap := make(map[string]models.WorkOrder)
	for _, wo := range workOrders {
		woMap[wo.ID] = wo
	}

	succeeded := []string{}
	failed := []models.BulkOperationError{}

	// Validate each work order
	validIDs := []string{}
	for _, id := range req.WorkOrderIDs {
		wo, exists := woMap[id]
		if !exists {
			failed = append(failed, models.BulkOperationError{
				ID:      id,
				Message: "work order not found",
				Code:    "not_found",
			})
			continue
		}

		// Check if status transition is valid (forward transitions only for bulk)
		if !isValidForwardTransition(wo.Status, req.Status) {
			failed = append(failed, models.BulkOperationError{
				ID:      id,
				Message: "invalid status transition from " + string(wo.Status) + " to " + string(req.Status),
				Code:    "invalid_transition",
			})
			continue
		}

		validIDs = append(validIDs, id)
	}

	// Perform bulk update for valid IDs
	if len(validIDs) > 0 {
		if err := h.pg.WorkOrders().BulkUpdateStatus(r.Context(), tenant, school, validIDs, req.Status, now); err != nil {
			h.log.Error("failed to bulk update status", zap.Error(err))
			// Mark all as failed
			for _, id := range validIDs {
				failed = append(failed, models.BulkOperationError{
					ID:      id,
					Message: "database error: " + err.Error(),
					Code:    "db_error",
				})
			}
		} else {
			succeeded = validIDs

			// Audit log each successful update
			for _, id := range validIDs {
				if wo, ok := woMap[id]; ok {
					_ = h.audit.LogUpdate(r.Context(), "work_order", id,
						map[string]any{"status": wo.Status},
						map[string]any{"status": req.Status, "bulkOpId": opLog.ID, "userId": userID})
				}
			}
		}
	}

	// Update operation log
	failedIDs := make([]string, len(failed))
	for i, f := range failed {
		failedIDs[i] = f.ID
	}
	_ = h.pg.BulkOperations().Update(r.Context(), opLog.ID, succeeded, failedIDs, failed, now)

	// Return response
	resp := models.BulkOperationResult{
		OperationID:  opLog.ID,
		Succeeded:    succeeded,
		Failed:       failed,
		TotalCount:   len(req.WorkOrderIDs),
		SuccessCount: len(succeeded),
		FailureCount: len(failed),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// BulkAssignment handles POST /work-orders/bulk/assignment
func (h *WorkOrderBulkHandler) BulkAssignment(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())
	school := middleware.SchoolID(r.Context())
	userID := middleware.UserID(r.Context())

	// Check feature flag
	enabled, err := h.pg.FeatureConfig().IsFeatureEnabled(r.Context(), tenant, models.FeatureWorkOrderBulkOperations)
	if err != nil {
		h.log.Error("failed to check feature flag", zap.Error(err))
	}
	if !enabled {
		http.Error(w, "bulk operations feature is disabled", http.StatusForbidden)
		return
	}

	// Get bulk config
	cfg, err := h.pg.FeatureConfig().GetBulkConfig(r.Context(), tenant)
	if err != nil {
		h.log.Error("failed to get bulk config", zap.Error(err))
		cfg = models.DefaultBulkConfig()
	}

	// Parse request
	var req models.BulkAssignmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if err := req.Validate(cfg.MaxBatchSize); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	now := time.Now().UTC()

	// Create operation log
	opLog := models.BulkOperationLog{
		ID:            store.NewID("bulk"),
		TenantID:      tenant,
		UserID:        userID,
		OperationType: models.BulkOpAssignment,
		EntityType:    "work_order",
		RequestedIDs:  req.WorkOrderIDs,
		SuccessfulIDs: []string{},
		FailedIDs:     []string{},
		StartedAt:     now,
		TotalCount:    len(req.WorkOrderIDs),
		CreatedAt:     now,
	}

	if err := h.pg.BulkOperations().Create(r.Context(), opLog); err != nil {
		h.log.Error("failed to create bulk operation log", zap.Error(err))
	}

	// Validate all work orders exist
	workOrders, err := h.pg.WorkOrders().GetByIDs(r.Context(), tenant, school, req.WorkOrderIDs)
	if err != nil {
		h.log.Error("failed to fetch work orders", zap.Error(err))
		http.Error(w, "failed to fetch work orders", http.StatusInternalServerError)
		return
	}

	// Build ID map
	woMap := make(map[string]models.WorkOrder)
	for _, wo := range workOrders {
		woMap[wo.ID] = wo
	}

	succeeded := []string{}
	failed := []models.BulkOperationError{}

	// Validate each work order exists
	validIDs := []string{}
	for _, id := range req.WorkOrderIDs {
		if _, exists := woMap[id]; !exists {
			failed = append(failed, models.BulkOperationError{
				ID:      id,
				Message: "work order not found",
				Code:    "not_found",
			})
			continue
		}
		validIDs = append(validIDs, id)
	}

	// Perform bulk assignment
	if len(validIDs) > 0 {
		if err := h.pg.WorkOrders().BulkUpdateAssignment(r.Context(), tenant, school, validIDs, req.AssignedStaffID, req.ServiceShopID, now); err != nil {
			h.log.Error("failed to bulk update assignment", zap.Error(err))
			for _, id := range validIDs {
				failed = append(failed, models.BulkOperationError{
					ID:      id,
					Message: "database error: " + err.Error(),
					Code:    "db_error",
				})
			}
		} else {
			succeeded = validIDs

			// Audit log each successful update
			for _, id := range validIDs {
				_ = h.audit.LogUpdate(r.Context(), "work_order", id, nil,
					map[string]any{
						"action":          "bulk_assignment",
						"assignedStaffId": req.AssignedStaffID,
						"serviceShopId":   req.ServiceShopID,
						"bulkOpId":        opLog.ID,
						"userId":          userID,
					})
			}
		}
	}

	// Update operation log
	failedIDs := make([]string, len(failed))
	for i, f := range failed {
		failedIDs[i] = f.ID
	}
	_ = h.pg.BulkOperations().Update(r.Context(), opLog.ID, succeeded, failedIDs, failed, now)

	// Return response
	resp := models.BulkOperationResult{
		OperationID:  opLog.ID,
		Succeeded:    succeeded,
		Failed:       failed,
		TotalCount:   len(req.WorkOrderIDs),
		SuccessCount: len(succeeded),
		FailureCount: len(failed),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// BulkApproval handles POST /work-orders/bulk/approval
func (h *WorkOrderBulkHandler) BulkApproval(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())
	school := middleware.SchoolID(r.Context())
	userID := middleware.UserID(r.Context())

	// Check feature flag
	enabled, err := h.pg.FeatureConfig().IsFeatureEnabled(r.Context(), tenant, models.FeatureWorkOrderBulkOperations)
	if err != nil {
		h.log.Error("failed to check feature flag", zap.Error(err))
	}
	if !enabled {
		http.Error(w, "bulk operations feature is disabled", http.StatusForbidden)
		return
	}

	// Get bulk config
	cfg, err := h.pg.FeatureConfig().GetBulkConfig(r.Context(), tenant)
	if err != nil {
		h.log.Error("failed to get bulk config", zap.Error(err))
		cfg = models.DefaultBulkConfig()
	}

	// Parse request
	var req models.BulkApprovalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if err := req.Validate(cfg.MaxBatchSize); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	now := time.Now().UTC()

	// Create operation log
	opLog := models.BulkOperationLog{
		ID:            store.NewID("bulk"),
		TenantID:      tenant,
		UserID:        userID,
		OperationType: models.BulkOpApproval,
		EntityType:    "work_order",
		RequestedIDs:  req.WorkOrderIDs,
		SuccessfulIDs: []string{},
		FailedIDs:     []string{},
		StartedAt:     now,
		TotalCount:    len(req.WorkOrderIDs),
		CreatedAt:     now,
	}

	if err := h.pg.BulkOperations().Create(r.Context(), opLog); err != nil {
		h.log.Error("failed to create bulk operation log", zap.Error(err))
	}

	// Validate all work orders exist and are in completable state
	workOrders, err := h.pg.WorkOrders().GetByIDs(r.Context(), tenant, school, req.WorkOrderIDs)
	if err != nil {
		h.log.Error("failed to fetch work orders", zap.Error(err))
		http.Error(w, "failed to fetch work orders", http.StatusInternalServerError)
		return
	}

	// Build ID map
	woMap := make(map[string]models.WorkOrder)
	for _, wo := range workOrders {
		woMap[wo.ID] = wo
	}

	succeeded := []string{}
	failed := []models.BulkOperationError{}

	// Determine target status based on decision
	var targetStatus models.WorkOrderStatus
	if req.Decision == "approved" {
		targetStatus = models.WorkOrderApproved
	} else {
		// For rejected, we send back to QA
		targetStatus = models.WorkOrderQA
	}

	// Validate each work order
	validIDs := []string{}
	for _, id := range req.WorkOrderIDs {
		wo, exists := woMap[id]
		if !exists {
			failed = append(failed, models.BulkOperationError{
				ID:      id,
				Message: "work order not found",
				Code:    "not_found",
			})
			continue
		}

		// Only completed work orders can be approved
		if wo.Status != models.WorkOrderCompleted {
			failed = append(failed, models.BulkOperationError{
				ID:      id,
				Message: "work order must be in completed status to approve",
				Code:    "invalid_status",
			})
			continue
		}

		validIDs = append(validIDs, id)
	}

	// Perform bulk status update
	if len(validIDs) > 0 {
		if err := h.pg.WorkOrders().BulkUpdateStatus(r.Context(), tenant, school, validIDs, targetStatus, now); err != nil {
			h.log.Error("failed to bulk update approval", zap.Error(err))
			for _, id := range validIDs {
				failed = append(failed, models.BulkOperationError{
					ID:      id,
					Message: "database error: " + err.Error(),
					Code:    "db_error",
				})
			}
		} else {
			succeeded = validIDs

			// Audit log each successful update
			for _, id := range validIDs {
				_ = h.audit.LogUpdate(r.Context(), "work_order", id, nil,
					map[string]any{
						"action":   "bulk_approval",
						"decision": req.Decision,
						"notes":    req.Notes,
						"bulkOpId": opLog.ID,
						"userId":   userID,
					})
			}
		}
	}

	// Update operation log
	failedIDs := make([]string, len(failed))
	for i, f := range failed {
		failedIDs[i] = f.ID
	}
	_ = h.pg.BulkOperations().Update(r.Context(), opLog.ID, succeeded, failedIDs, failed, now)

	// Return response
	resp := models.BulkOperationResult{
		OperationID:  opLog.ID,
		Succeeded:    succeeded,
		Failed:       failed,
		TotalCount:   len(req.WorkOrderIDs),
		SuccessCount: len(succeeded),
		FailureCount: len(failed),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// isValidForwardTransition checks if a status transition is valid (forward only).
func isValidForwardTransition(from, to models.WorkOrderStatus) bool {
	validTransitions := map[models.WorkOrderStatus][]models.WorkOrderStatus{
		models.WorkOrderDraft:     {models.WorkOrderAssigned},
		models.WorkOrderAssigned:  {models.WorkOrderInRepair},
		models.WorkOrderInRepair:  {models.WorkOrderQA},
		models.WorkOrderQA:        {models.WorkOrderCompleted},
		models.WorkOrderCompleted: {models.WorkOrderApproved},
	}

	allowed, exists := validTransitions[from]
	if !exists {
		return false
	}
	for _, s := range allowed {
		if s == to {
			return true
		}
	}
	return false
}
