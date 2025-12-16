package service

import "github.com/edvirons/ssp/ims/internal/models"

// Valid status transitions (simple but useful).
func CanTransitionIncident(from, to models.IncidentStatus) bool {
	switch from {
	case models.IncidentNew:
		return to == models.IncidentAcknowledged || to == models.IncidentEscalated
	case models.IncidentAcknowledged:
		return to == models.IncidentInProgress || to == models.IncidentEscalated
	case models.IncidentInProgress:
		return to == models.IncidentResolved || to == models.IncidentEscalated
	case models.IncidentEscalated:
		return to == models.IncidentInProgress || to == models.IncidentResolved
	case models.IncidentResolved:
		return to == models.IncidentClosed
	default:
		return false
	}
}

func CanTransitionWorkOrder(from, to models.WorkOrderStatus) bool {
	switch from {
	case models.WorkOrderDraft:
		return to == models.WorkOrderAssigned
	case models.WorkOrderAssigned:
		return to == models.WorkOrderInRepair
	case models.WorkOrderInRepair:
		return to == models.WorkOrderQA || to == models.WorkOrderCompleted
	case models.WorkOrderQA:
		return to == models.WorkOrderCompleted
	case models.WorkOrderCompleted:
		return to == models.WorkOrderApproved
	default:
		return false
	}
}
