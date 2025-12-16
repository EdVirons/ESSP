package models

import "time"

// NotificationType represents the type of work order notification.
type NotificationType string

const (
	NotificationStatusChange      NotificationType = "status_change"
	NotificationAssignment        NotificationType = "assignment"
	NotificationApprovalRequested NotificationType = "approval_requested"
	NotificationApprovalDecided   NotificationType = "approval_decided"
	NotificationDeliverableReview NotificationType = "deliverable_review"
	NotificationReworkRequired    NotificationType = "rework_required"
)

// AllNotificationTypes contains all notification types.
var AllNotificationTypes = []NotificationType{
	NotificationStatusChange,
	NotificationAssignment,
	NotificationApprovalRequested,
	NotificationApprovalDecided,
	NotificationDeliverableReview,
	NotificationReworkRequired,
}

// UserNotificationPreferences stores a user's notification preferences.
type UserNotificationPreferences struct {
	ID                 string             `json:"id"`
	TenantID           string             `json:"tenantId"`
	UserID             string             `json:"userId"`
	EnabledTypes       []NotificationType `json:"enabledTypes"`
	InAppEnabled       bool               `json:"inAppEnabled"`
	EmailEnabled       bool               `json:"emailEnabled"`
	QuietHoursStart    string             `json:"quietHoursStart"`
	QuietHoursEnd      string             `json:"quietHoursEnd"`
	QuietHoursTimezone string             `json:"quietHoursTimezone"`
	CreatedAt          time.Time          `json:"createdAt"`
	UpdatedAt          time.Time          `json:"updatedAt"`
}

// IsNotificationEnabled checks if a notification type is enabled.
func (p *UserNotificationPreferences) IsNotificationEnabled(t NotificationType) bool {
	for _, enabled := range p.EnabledTypes {
		if enabled == t {
			return true
		}
	}
	return false
}

// WorkOrderNotificationEvent represents a notification event to be sent.
type WorkOrderNotificationEvent struct {
	Type           NotificationType `json:"type"`
	TenantID       string           `json:"tenantId"`
	SchoolID       string           `json:"schoolId"`
	WorkOrderID    string           `json:"workOrderId"`
	WorkOrderTitle string           `json:"workOrderTitle"`
	ActorUserID    string           `json:"actorUserId"`
	ActorName      string           `json:"actorName"`
	TargetUserIDs  []string         `json:"targetUserIds"`
	TargetRoles    []string         `json:"targetRoles"`
	Payload        map[string]any   `json:"payload"`
	CreatedAt      time.Time        `json:"createdAt"`
}

// NotificationTargetsByEvent defines which roles receive notifications by event type.
var NotificationTargetsByEvent = map[NotificationType][]string{
	NotificationStatusChange:      {"ssp_lead_tech", "ssp_support_agent"},
	NotificationAssignment:        {"ssp_field_tech", "ssp_lead_tech"},
	NotificationApprovalRequested: {"ssp_admin", "ssp_lead_tech"},
	NotificationApprovalDecided:   {"ssp_school_contact", "ssp_support_agent"},
	NotificationDeliverableReview: {"ssp_lead_tech"},
	NotificationReworkRequired:    {"ssp_field_tech", "ssp_lead_tech"},
}

// UpdateNotificationPreferencesRequest is the request to update preferences.
type UpdateNotificationPreferencesRequest struct {
	EnabledTypes       []NotificationType `json:"enabledTypes"`
	InAppEnabled       *bool              `json:"inAppEnabled,omitempty"`
	EmailEnabled       *bool              `json:"emailEnabled,omitempty"`
	QuietHoursStart    *string            `json:"quietHoursStart,omitempty"`
	QuietHoursEnd      *string            `json:"quietHoursEnd,omitempty"`
	QuietHoursTimezone *string            `json:"quietHoursTimezone,omitempty"`
}
