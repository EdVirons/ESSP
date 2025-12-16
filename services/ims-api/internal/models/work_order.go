package models

import "time"

// WorkOrderStatus represents the status of a work order.
type WorkOrderStatus string

const (
	WorkOrderDraft     WorkOrderStatus = "draft"
	WorkOrderAssigned  WorkOrderStatus = "assigned"
	WorkOrderInRepair  WorkOrderStatus = "in_repair"
	WorkOrderQA        WorkOrderStatus = "qa"
	WorkOrderCompleted WorkOrderStatus = "completed"
	WorkOrderApproved  WorkOrderStatus = "approved"
)

// RepairLocation represents where the repair is performed.
type RepairLocation string

const (
	RepairLocationServiceShop RepairLocation = "service_shop"
	RepairLocationOnSite      RepairLocation = "on_site"
)

// WorkOrder represents a repair work order.
type WorkOrder struct {
	ID         string `json:"id"`
	IncidentID string `json:"incidentId"`
	TenantID   string `json:"tenantId"`
	SchoolID   string `json:"schoolId"`
	DeviceID   string `json:"deviceId"`

	// Denormalized lookup fields (from SSOT snapshots)
	SchoolName     string `json:"schoolName"`
	ContactName    string `json:"contactName"`
	ContactPhone   string `json:"contactPhone"`
	ContactEmail   string `json:"contactEmail"`
	DeviceSerial   string `json:"deviceSerial"`
	DeviceAssetTag string `json:"deviceAssetTag"`
	DeviceModelID  string `json:"deviceModelId"`
	DeviceMake     string `json:"deviceMake"`
	DeviceModel    string `json:"deviceModel"`
	DeviceCategory string `json:"deviceCategory"`

	Status            WorkOrderStatus `json:"status"`
	ServiceShopID     string          `json:"serviceShopId"`
	AssignedStaffID   string          `json:"assignedStaffId"`
	RepairLocation    RepairLocation  `json:"repairLocation"`
	AssignedTo        string          `json:"assignedTo"`
	TaskType          string          `json:"taskType"`
	ProjectID         string          `json:"projectId"`
	PhaseID           string          `json:"phaseId"`
	OnsiteContactID   string          `json:"onsiteContactId"`
	ApprovalStatus    string          `json:"approvalStatus"`
	CostEstimateCents int64           `json:"costEstimateCents"`
	Notes             string          `json:"notes"`

	// Project-originated work orders
	CreatedFromProject bool   `json:"createdFromProject"`
	CreatedByUserID    string `json:"createdByUserId"`
	CreatedByUserName  string `json:"createdByUserName"`

	// Rework tracking
	ReworkCount      int        `json:"reworkCount"`
	LastReworkAt     *time.Time `json:"lastReworkAt,omitempty"`
	LastReworkReason string     `json:"lastReworkReason"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// WorkOrderSchedule represents a scheduled work order session.
type WorkOrderSchedule struct {
	ID              string     `json:"id"`
	TenantID        string     `json:"tenantId"`
	SchoolID        string     `json:"schoolId"`
	WorkOrderID     string     `json:"workOrderId"`
	PhaseID         string     `json:"phaseId"`
	ScheduledStart  *time.Time `json:"scheduledStart"`
	ScheduledEnd    *time.Time `json:"scheduledEnd"`
	Timezone        string     `json:"timezone"`
	Notes           string     `json:"notes"`
	CreatedByUserID string     `json:"createdByUserId"`
	CreatedAt       time.Time  `json:"createdAt"`
}

// DeliverableStatus represents the status of a deliverable.
type DeliverableStatus string

const (
	DeliverablePending   DeliverableStatus = "pending"
	DeliverableSubmitted DeliverableStatus = "submitted"
	DeliverableApproved  DeliverableStatus = "approved"
	DeliverableRejected  DeliverableStatus = "rejected"
)

// WorkOrderDeliverable represents a deliverable for a work order.
type WorkOrderDeliverable struct {
	ID                   string            `json:"id"`
	TenantID             string            `json:"tenantId"`
	SchoolID             string            `json:"schoolId"`
	WorkOrderID          string            `json:"workOrderId"`
	PhaseID              string            `json:"phaseId"`
	Title                string            `json:"title"`
	Description          string            `json:"description"`
	Status               DeliverableStatus `json:"status"`
	EvidenceAttachmentID string            `json:"evidenceAttachmentId"`
	SubmittedByUserID    string            `json:"submittedByUserId"`
	SubmittedAt          *time.Time        `json:"submittedAt"`
	ReviewedByUserID     string            `json:"reviewedByUserId"`
	ReviewedAt           *time.Time        `json:"reviewedAt"`
	ReviewNotes          string            `json:"reviewNotes"`
	CreatedAt            time.Time         `json:"createdAt"`
	UpdatedAt            time.Time         `json:"updatedAt"`
}
