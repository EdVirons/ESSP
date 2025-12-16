package models

import "time"

// RejectionCategory represents the reason category for work order rejection.
type RejectionCategory string

const (
	RejectionQuality    RejectionCategory = "quality"
	RejectionIncomplete RejectionCategory = "incomplete"
	RejectionWrongParts RejectionCategory = "wrong_parts"
	RejectionSafety     RejectionCategory = "safety"
	RejectionOther      RejectionCategory = "other"
)

// ValidRejectionCategories contains all valid rejection categories.
var ValidRejectionCategories = []RejectionCategory{
	RejectionQuality,
	RejectionIncomplete,
	RejectionWrongParts,
	RejectionSafety,
	RejectionOther,
}

// IsValidRejectionCategory checks if a category is valid.
func IsValidRejectionCategory(c string) bool {
	for _, v := range ValidRejectionCategories {
		if string(v) == c {
			return true
		}
	}
	return false
}

// ValidReworkTransitions defines which backward transitions are allowed.
// Key is current status, value is list of allowed target statuses.
var ValidReworkTransitions = map[WorkOrderStatus][]WorkOrderStatus{
	WorkOrderAssigned:  {WorkOrderDraft},
	WorkOrderInRepair:  {WorkOrderAssigned, WorkOrderDraft},
	WorkOrderQA:        {WorkOrderInRepair, WorkOrderAssigned},
	WorkOrderCompleted: {WorkOrderQA, WorkOrderInRepair},
	WorkOrderApproved:  {WorkOrderCompleted},
}

// CanReworkTo checks if a work order can be sent back to a target status.
func CanReworkTo(currentStatus, targetStatus WorkOrderStatus) bool {
	allowed, exists := ValidReworkTransitions[currentStatus]
	if !exists {
		return false
	}
	for _, s := range allowed {
		if s == targetStatus {
			return true
		}
	}
	return false
}

// WorkOrderReworkHistory records a rejection/rework event.
type WorkOrderReworkHistory struct {
	ID                string            `json:"id"`
	TenantID          string            `json:"tenantId"`
	SchoolID          string            `json:"schoolId"`
	WorkOrderID       string            `json:"workOrderId"`
	FromStatus        WorkOrderStatus   `json:"fromStatus"`
	ToStatus          WorkOrderStatus   `json:"toStatus"`
	RejectionReason   string            `json:"rejectionReason"`
	RejectionCategory RejectionCategory `json:"rejectionCategory"`
	RejectedByUserID  string            `json:"rejectedByUserId"`
	RejectedByName    string            `json:"rejectedByName"`
	ReworkSequence    int               `json:"reworkSequence"`
	CreatedAt         time.Time         `json:"createdAt"`
}

// RejectWorkOrderRequest is the request body for rejecting a work order.
type RejectWorkOrderRequest struct {
	TargetStatus WorkOrderStatus   `json:"targetStatus"`
	Reason       string            `json:"reason"`
	Category     RejectionCategory `json:"category"`
}

// RejectWorkOrderResponse is the response after rejecting a work order.
type RejectWorkOrderResponse struct {
	WorkOrder     WorkOrder              `json:"workOrder"`
	ReworkHistory WorkOrderReworkHistory `json:"reworkHistory"`
}
