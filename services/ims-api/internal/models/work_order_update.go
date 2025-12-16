package models

// WorkOrderUpdateRequest contains optional fields for partial work order updates.
// Fields that are nil will not be updated.
type WorkOrderUpdateRequest struct {
	AssignedStaffID   *string         `json:"assignedStaffId,omitempty"`
	ServiceShopID     *string         `json:"serviceShopId,omitempty"`
	CostEstimateCents *int64          `json:"costEstimateCents,omitempty"`
	Notes             *string         `json:"notes,omitempty"`
	RepairLocation    *RepairLocation `json:"repairLocation,omitempty"`
	OnsiteContactID   *string         `json:"onsiteContactId,omitempty"`
}

// HasUpdates returns true if at least one field is set.
func (r *WorkOrderUpdateRequest) HasUpdates() bool {
	return r.AssignedStaffID != nil ||
		r.ServiceShopID != nil ||
		r.CostEstimateCents != nil ||
		r.Notes != nil ||
		r.RepairLocation != nil ||
		r.OnsiteContactID != nil
}

// ToMap converts the update request to a map of field names to values.
// Only non-nil fields are included.
func (r *WorkOrderUpdateRequest) ToMap() map[string]any {
	updates := make(map[string]any)
	if r.AssignedStaffID != nil {
		updates["assigned_staff_id"] = *r.AssignedStaffID
	}
	if r.ServiceShopID != nil {
		updates["service_shop_id"] = *r.ServiceShopID
	}
	if r.CostEstimateCents != nil {
		updates["cost_estimate_cents"] = *r.CostEstimateCents
	}
	if r.Notes != nil {
		updates["notes"] = *r.Notes
	}
	if r.RepairLocation != nil {
		updates["repair_location"] = string(*r.RepairLocation)
	}
	if r.OnsiteContactID != nil {
		updates["onsite_contact_id"] = *r.OnsiteContactID
	}
	return updates
}

// WorkOrderUpdateResponse is the response after updating a work order.
type WorkOrderUpdateResponse struct {
	WorkOrder      WorkOrder         `json:"workOrder"`
	UpdatedFields  []string          `json:"updatedFields"`
	PreviousValues map[string]any    `json:"previousValues,omitempty"`
}
