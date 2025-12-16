package models

import "time"

// ActivityType represents the type of activity in a project.
type ActivityType string

const (
	ActivityComment         ActivityType = "comment"
	ActivityNote            ActivityType = "note"
	ActivityFileUpload      ActivityType = "file_upload"
	ActivityStatusChange    ActivityType = "status_change"
	ActivityAssignment      ActivityType = "assignment"
	ActivityWorkOrder       ActivityType = "work_order"
	ActivityPhaseTransition ActivityType = "phase_transition"
	ActivityMention         ActivityType = "mention"
)

// ValidActivityTypes returns all valid activity types.
func ValidActivityTypes() []ActivityType {
	return []ActivityType{
		ActivityComment,
		ActivityNote,
		ActivityFileUpload,
		ActivityStatusChange,
		ActivityAssignment,
		ActivityWorkOrder,
		ActivityPhaseTransition,
		ActivityMention,
	}
}

// IsValidActivityType checks if an activity type is valid.
func IsValidActivityType(t string) bool {
	switch ActivityType(t) {
	case ActivityComment, ActivityNote, ActivityFileUpload, ActivityStatusChange,
		ActivityAssignment, ActivityWorkOrder, ActivityPhaseTransition, ActivityMention:
		return true
	default:
		return false
	}
}

// ActivityVisibility represents the visibility of an activity.
type ActivityVisibility string

const (
	VisibilityTeam    ActivityVisibility = "team"
	VisibilityPublic  ActivityVisibility = "public"
	VisibilityPrivate ActivityVisibility = "private"
)

// ProjectActivity represents an activity in a project's timeline.
type ProjectActivity struct {
	ID            string             `json:"id"`
	TenantID      string             `json:"tenantId"`
	ProjectID     string             `json:"projectId"`
	PhaseID       string             `json:"phaseId,omitempty"`
	WorkOrderID   string             `json:"workOrderId,omitempty"`
	ActivityType  ActivityType       `json:"activityType"`
	ActorUserID   string             `json:"actorUserId"`
	ActorEmail    string             `json:"actorEmail"`
	ActorName     string             `json:"actorName"`
	Content       string             `json:"content"`
	Metadata      map[string]any     `json:"metadata"`
	AttachmentIDs []string           `json:"attachmentIds"`
	Visibility    ActivityVisibility `json:"visibility"`
	IsPinned      bool               `json:"isPinned"`
	EditedAt      *time.Time         `json:"editedAt,omitempty"`
	DeletedAt     *time.Time         `json:"deletedAt,omitempty"`
	CreatedAt     time.Time          `json:"createdAt"`
}

// ProjectAttachment represents a file attachment in a project.
type ProjectAttachment struct {
	ID                 string    `json:"id"`
	TenantID           string    `json:"tenantId"`
	ProjectID          string    `json:"projectId"`
	PhaseID            string    `json:"phaseId,omitempty"`
	ActivityID         string    `json:"activityId,omitempty"`
	FileName           string    `json:"fileName"`
	ContentType        string    `json:"contentType"`
	SizeBytes          int64     `json:"sizeBytes"`
	ObjectKey          string    `json:"objectKey"`
	UploadedByUserID   string    `json:"uploadedByUserId"`
	UploadedByUserName string    `json:"uploadedByUserName"`
	CreatedAt          time.Time `json:"createdAt"`
}

// ProjectNotificationType represents the type of project notification.
type ProjectNotificationType string

const (
	ProjectNotificationAssignment   ProjectNotificationType = "assignment"
	ProjectNotificationMention      ProjectNotificationType = "mention"
	ProjectNotificationStatusChange ProjectNotificationType = "status_change"
	ProjectNotificationComment      ProjectNotificationType = "comment"
	ProjectNotificationWorkOrder    ProjectNotificationType = "work_order"
)

// UserNotification represents a notification for a user.
type UserNotification struct {
	ID               string                  `json:"id"`
	TenantID         string                  `json:"tenantId"`
	UserID           string                  `json:"userId"`
	NotificationType ProjectNotificationType `json:"notificationType"`
	EntityType       string           `json:"entityType"` // project|phase|work_order
	EntityID         string           `json:"entityId"`
	ProjectID        string           `json:"projectId,omitempty"`
	Title            string           `json:"title"`
	Body             string           `json:"body"`
	Metadata         map[string]any   `json:"metadata"`
	IsRead           bool             `json:"isRead"`
	ReadAt           *time.Time       `json:"readAt,omitempty"`
	CreatedAt        time.Time        `json:"createdAt"`
}

// StatusChangeMetadata holds metadata for status change activities.
type StatusChangeMetadata struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Reason string `json:"reason,omitempty"`
}

// AssignmentMetadata holds metadata for assignment activities.
type AssignmentMetadata struct {
	UserID   string `json:"userId"`
	UserName string `json:"userName"`
	Action   string `json:"action"` // added|removed
	Role     string `json:"role"`
}

// WorkOrderMetadata holds metadata for work order activities.
type WorkOrderMetadata struct {
	WorkOrderID string `json:"workOrderId"`
	Action      string `json:"action"` // created|updated|completed
	Status      string `json:"status,omitempty"`
}

// PhaseTransitionMetadata holds metadata for phase transition activities.
type PhaseTransitionMetadata struct {
	PhaseID     string `json:"phaseId"`
	PhaseType   string `json:"phaseType"`
	From        string `json:"from"`
	To          string `json:"to"`
	CompletedBy string `json:"completedBy,omitempty"`
}

// FileUploadMetadata holds metadata for file upload activities.
type FileUploadMetadata struct {
	FileID      string `json:"fileId"`
	FileName    string `json:"fileName"`
	FileSize    int64  `json:"fileSize"`
	ContentType string `json:"contentType"`
}

// CommentMetadata holds metadata for comment activities.
type CommentMetadata struct {
	Mentions []string `json:"mentions,omitempty"`
}
