package models

import "time"

// WorkOrderPart represents a part in a work order BOM.
type WorkOrderPart struct {
	ID            string `json:"id"`
	TenantID      string `json:"tenantId"`
	SchoolID      string `json:"schoolId"`
	WorkOrderID   string `json:"workOrderId"`
	ServiceShopID string `json:"serviceShopId"`
	PartID        string `json:"partId"`

	// Denormalized part/device fields
	PartName      string `json:"partName"`
	PartPUK       string `json:"partPuk"`
	PartCategory  string `json:"partCategory"`
	DeviceModelID string `json:"deviceModelId"`
	IsCompatible  bool   `json:"isCompatible"`

	QtyPlanned int64     `json:"qtyPlanned"`
	QtyUsed    int64     `json:"qtyUsed"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

// BOQItem represents an item in a bill of quantities.
type BOQItem struct {
	ID                 string    `json:"id"`
	TenantID           string    `json:"tenantId"`
	ProjectID          string    `json:"projectId"`
	Category           string    `json:"category"`
	Description        string    `json:"description"`
	PartID             string    `json:"partId"`
	Qty                int64     `json:"qty"`
	Unit               string    `json:"unit"`
	EstimatedCostCents int64     `json:"estimatedCostCents"`
	Approved           bool      `json:"approved"`
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"`
}
