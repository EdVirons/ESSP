package models

import (
	"encoding/json"
	"time"
)

// BulkOperationType represents the type of bulk operation.
type BulkOperationType string

const (
	BulkOpStatusUpdate BulkOperationType = "status_update"
	BulkOpAssignment   BulkOperationType = "assignment"
	BulkOpApproval     BulkOperationType = "approval"
)

// BulkOperationLog tracks bulk operations for audit purposes.
type BulkOperationLog struct {
	ID            string            `json:"id"`
	TenantID      string            `json:"tenantId"`
	UserID        string            `json:"userId"`
	OperationType BulkOperationType `json:"operationType"`
	EntityType    string            `json:"entityType"`
	RequestedIDs  []string          `json:"requestedIds"`
	SuccessfulIDs []string          `json:"successfulIds"`
	FailedIDs     []string          `json:"failedIds"`
	Errors        json.RawMessage   `json:"errors"`
	StartedAt     time.Time         `json:"startedAt"`
	CompletedAt   *time.Time        `json:"completedAt,omitempty"`
	TotalCount    int               `json:"totalCount"`
	SuccessCount  int               `json:"successCount"`
	FailureCount  int               `json:"failureCount"`
	CreatedAt     time.Time         `json:"createdAt"`
}

// BulkOperationError represents an error for a single item in a bulk operation.
type BulkOperationError struct {
	ID      string `json:"id"`
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
}

// BulkStatusUpdateRequest is the request for bulk status updates.
type BulkStatusUpdateRequest struct {
	WorkOrderIDs []string        `json:"workOrderIds"`
	Status       WorkOrderStatus `json:"status"`
}

// BulkAssignmentRequest is the request for bulk assignment updates.
type BulkAssignmentRequest struct {
	WorkOrderIDs    []string `json:"workOrderIds"`
	AssignedStaffID *string  `json:"assignedStaffId,omitempty"`
	ServiceShopID   *string  `json:"serviceShopId,omitempty"`
}

// BulkApprovalRequest is the request for bulk approvals.
type BulkApprovalRequest struct {
	WorkOrderIDs []string `json:"workOrderIds"`
	Decision     string   `json:"decision"` // "approved" or "rejected"
	Notes        string   `json:"notes,omitempty"`
}

// BulkOperationResult is the response for bulk operations.
type BulkOperationResult struct {
	OperationID  string               `json:"operationId"`
	Succeeded    []string             `json:"succeeded"`
	Failed       []BulkOperationError `json:"failed"`
	TotalCount   int                  `json:"totalCount"`
	SuccessCount int                  `json:"successCount"`
	FailureCount int                  `json:"failureCount"`
}

// BulkOperationConfig holds configuration for bulk operations.
type BulkOperationConfig struct {
	MaxBatchSize       int `json:"maxBatchSize"`
	RateLimitPerMinute int `json:"rateLimitPerMinute"`
}

// DefaultBulkConfig returns default bulk operation configuration.
func DefaultBulkConfig() BulkOperationConfig {
	return BulkOperationConfig{
		MaxBatchSize:       100,
		RateLimitPerMinute: 10,
	}
}

// Validate checks if a bulk status update request is valid.
func (r *BulkStatusUpdateRequest) Validate(maxBatchSize int) error {
	if len(r.WorkOrderIDs) == 0 {
		return ErrBulkNoIDs
	}
	if len(r.WorkOrderIDs) > maxBatchSize {
		return ErrBulkTooManyIDs
	}
	if r.Status == "" {
		return ErrBulkNoStatus
	}
	return nil
}

// Validate checks if a bulk assignment request is valid.
func (r *BulkAssignmentRequest) Validate(maxBatchSize int) error {
	if len(r.WorkOrderIDs) == 0 {
		return ErrBulkNoIDs
	}
	if len(r.WorkOrderIDs) > maxBatchSize {
		return ErrBulkTooManyIDs
	}
	if r.AssignedStaffID == nil && r.ServiceShopID == nil {
		return ErrBulkNoAssignment
	}
	return nil
}

// Validate checks if a bulk approval request is valid.
func (r *BulkApprovalRequest) Validate(maxBatchSize int) error {
	if len(r.WorkOrderIDs) == 0 {
		return ErrBulkNoIDs
	}
	if len(r.WorkOrderIDs) > maxBatchSize {
		return ErrBulkTooManyIDs
	}
	if r.Decision != "approved" && r.Decision != "rejected" {
		return ErrBulkInvalidDecision
	}
	return nil
}
