package models

import (
	"encoding/json"
	"time"
)

// ThreadType represents the type of message thread
type ThreadType string

const (
	ThreadTypeGeneral  ThreadType = "general"
	ThreadTypeIncident ThreadType = "incident"
	ThreadTypeLivechat ThreadType = "livechat"
)

// ThreadStatus represents the status of a message thread
type ThreadStatus string

const (
	ThreadStatusOpen     ThreadStatus = "open"
	ThreadStatusClosed   ThreadStatus = "closed"
	ThreadStatusArchived ThreadStatus = "archived"
)

// ContentType represents the type of message content
type ContentType string

const (
	ContentTypeText       ContentType = "text"
	ContentTypeSystem     ContentType = "system"
	ContentTypeAttachment ContentType = "attachment"
)

// MessageThread represents a conversation thread
type MessageThread struct {
	ID                 string       `json:"id"`
	TenantID           string       `json:"tenantId"`
	SchoolID           string       `json:"schoolId"`
	Subject            string       `json:"subject"`
	ThreadType         ThreadType   `json:"threadType"`
	Status             ThreadStatus `json:"status"`
	IncidentID         *string      `json:"incidentId,omitempty"`
	CreatedBy          string       `json:"createdBy"`
	CreatedByRole      string       `json:"createdByRole"`
	CreatedByName      string       `json:"createdByName"`
	MessageCount       int          `json:"messageCount"`
	UnreadCountSchool  int          `json:"unreadCountSchool"`
	UnreadCountSupport int          `json:"unreadCountSupport"`
	LastMessageAt      *time.Time   `json:"lastMessageAt,omitempty"`
	CreatedAt          time.Time    `json:"createdAt"`
	UpdatedAt          time.Time    `json:"updatedAt"`
	ClosedAt           *time.Time   `json:"closedAt,omitempty"`

	// Denormalized fields for display
	SchoolName  string   `json:"schoolName,omitempty"`
	LastMessage *Message `json:"lastMessage,omitempty"`
}

// Message represents an individual message in a thread
type Message struct {
	ID          string            `json:"id"`
	TenantID    string            `json:"tenantId"`
	ThreadID    string            `json:"threadId"`
	SenderID    string            `json:"senderId"`
	SenderName  string            `json:"senderName"`
	SenderRole  string            `json:"senderRole"`
	Content     string            `json:"content"`
	ContentType ContentType       `json:"contentType"`
	Metadata    map[string]any    `json:"metadata,omitempty"`
	EditedAt    *time.Time        `json:"editedAt,omitempty"`
	DeletedAt   *time.Time        `json:"deletedAt,omitempty"`
	CreatedAt   time.Time         `json:"createdAt"`
	Attachments []MessageAttachment `json:"attachments,omitempty"`
}

// MessageAttachment represents a file attached to a message
type MessageAttachment struct {
	ID           string    `json:"id"`
	TenantID     string    `json:"tenantId"`
	MessageID    string    `json:"messageId"`
	FileName     string    `json:"fileName"`
	ContentType  string    `json:"contentType"`
	SizeBytes    int64     `json:"sizeBytes"`
	ObjectKey    string    `json:"objectKey"`
	ThumbnailKey *string   `json:"thumbnailKey,omitempty"`
	CreatedAt    time.Time `json:"createdAt"`

	// Presigned URLs (populated on read)
	DownloadURL  string `json:"downloadUrl,omitempty"`
	ThumbnailURL string `json:"thumbnailUrl,omitempty"`
}

// ThreadParticipant represents a user participating in a thread
type ThreadParticipant struct {
	ThreadID string     `json:"threadId"`
	UserID   string     `json:"userId"`
	UserName string     `json:"userName"`
	UserRole string     `json:"userRole"`
	JoinedAt time.Time  `json:"joinedAt"`
	LeftAt   *time.Time `json:"leftAt,omitempty"`
}

// ReadReceipt tracks when a user last read a thread
type ReadReceipt struct {
	ThreadID          string    `json:"threadId"`
	UserID            string    `json:"userId"`
	LastReadMessageID string    `json:"lastReadMessageId"`
	LastReadAt        time.Time `json:"lastReadAt"`
}

// UnreadCounts holds unread message counts for a user
type UnreadCounts struct {
	Threads  int `json:"threads"`
	Messages int `json:"messages"`
}

// CreateThreadRequest represents the request to create a new thread
type CreateThreadRequest struct {
	Subject        string   `json:"subject"`
	IncidentID     *string  `json:"incidentId,omitempty"`
	InitialMessage string   `json:"initialMessage"`
	Attachments    []string `json:"attachments,omitempty"`
	SchoolID       string   `json:"schoolId,omitempty"` // Optional, for support agents creating threads
}

// CreateMessageRequest represents the request to send a message
type CreateMessageRequest struct {
	Content     string   `json:"content"`
	Attachments []string `json:"attachments,omitempty"`
}

// UpdateThreadStatusRequest represents the request to update thread status
type UpdateThreadStatusRequest struct {
	Status ThreadStatus `json:"status"`
}

// MarkReadRequest represents the request to mark a thread as read
type MarkReadRequest struct {
	LastMessageID string `json:"lastMessageId"`
}

// ThreadsListResponse represents a paginated list of threads
type ThreadsListResponse struct {
	Items      []MessageThread `json:"items"`
	NextCursor string          `json:"nextCursor,omitempty"`
}

// MessagesListResponse represents a paginated list of messages
type MessagesListResponse struct {
	Items      []Message `json:"items"`
	NextCursor string    `json:"nextCursor,omitempty"`
}

// ThreadDetailResponse represents a thread with its messages
type ThreadDetailResponse struct {
	Thread       MessageThread       `json:"thread"`
	Messages     []Message           `json:"messages"`
	Participants []ThreadParticipant `json:"participants,omitempty"`
}

// CreateThreadResponse represents the response when creating a thread
type CreateThreadResponse struct {
	Thread  MessageThread `json:"thread"`
	Message Message       `json:"message"`
}

// CreateMessageResponse represents the response when creating a message
type CreateMessageResponse struct {
	Message Message `json:"message"`
}

// UploadURLRequest represents the request to get an upload URL
type UploadURLRequest struct {
	FileName    string `json:"fileName"`
	ContentType string `json:"contentType"`
	SizeBytes   int64  `json:"sizeBytes"`
}

// UploadURLResponse represents the response with an upload URL
type UploadURLResponse struct {
	UploadURL     string `json:"uploadUrl"`
	AttachmentRef string `json:"attachmentRef"`
}

// MessageSearchResult represents a message in search results
type MessageSearchResult struct {
	Message    Message       `json:"message"`
	Thread     MessageThread `json:"thread"`
	Highlights []string      `json:"highlights,omitempty"`
}

// SearchResponse represents search results
type SearchResponse struct {
	Items []MessageSearchResult `json:"items"`
	Total int                   `json:"total"`
}

// MessagingAnalytics represents messaging analytics data
type MessagingAnalytics struct {
	ThreadsCreated    int     `json:"threadsCreated"`
	MessagesSent      int     `json:"messagesSent"`
	AvgResponseTime   float64 `json:"avgResponseTimeMinutes"`
	ChatSessions      int     `json:"chatSessions"`
	AvgChatRating     float64 `json:"avgChatRating"`
	ActiveThreads     int     `json:"activeThreads"`
	ClosedThreads     int     `json:"closedThreads"`
	AvgMessagesThread float64 `json:"avgMessagesPerThread"`
}

// MetadataBytes returns metadata as JSON bytes
func (m *Message) MetadataBytes() []byte {
	if m.Metadata == nil {
		return []byte("{}")
	}
	b, _ := json.Marshal(m.Metadata)
	return b
}

// IsSchoolRole returns true if the role is a school contact
func IsSchoolRole(role string) bool {
	return role == "ssp_school_contact"
}

// IsSupportRole returns true if the role is a support/admin role
func IsSupportRole(role string) bool {
	return role == "ssp_admin" || role == "ssp_support_agent"
}
