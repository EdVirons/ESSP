package models

import "time"

// AttachmentEntityType represents the type of entity an attachment is linked to.
type AttachmentEntityType string

const (
	AttachmentIncident  AttachmentEntityType = "incident"
	AttachmentWorkOrder AttachmentEntityType = "work_order"
)

// Attachment represents a file attachment.
type Attachment struct {
	ID         string               `json:"id"`
	TenantID   string               `json:"tenantId"`
	SchoolID   string               `json:"schoolId"`
	EntityType AttachmentEntityType `json:"entityType"`
	EntityID   string               `json:"entityId"`

	FileName    string `json:"fileName"`
	ContentType string `json:"contentType"`
	SizeBytes   int64  `json:"sizeBytes"`

	ObjectKey string    `json:"objectKey"` // S3/MinIO key
	CreatedAt time.Time `json:"createdAt"`
}
