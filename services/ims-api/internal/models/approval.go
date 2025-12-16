package models

import "time"

// ApprovalStatus represents the status of an approval.
type ApprovalStatus string

const (
	ApprovalPending  ApprovalStatus = "pending"
	ApprovalApproved ApprovalStatus = "approved"
	ApprovalRejected ApprovalStatus = "rejected"
)

// WorkOrderApproval represents an approval for a work order.
type WorkOrderApproval struct {
	ID                string         `json:"id"`
	TenantID          string         `json:"tenantId"`
	SchoolID          string         `json:"schoolId"`
	WorkOrderID       string         `json:"workOrderId"`
	PhaseID           string         `json:"phaseId"`
	ApprovalType      string         `json:"approvalType"`
	RequestedByUserID string         `json:"requestedByUserId"`
	RequestedAt       time.Time      `json:"requestedAt"`
	Status            ApprovalStatus `json:"status"`
	DecidedByUserID   string         `json:"decidedByUserId"`
	DecidedAt         *time.Time     `json:"decidedAt"`
	DecisionNotes     string         `json:"decisionNotes"`
}
