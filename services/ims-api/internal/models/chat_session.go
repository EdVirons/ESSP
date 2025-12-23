package models

import "time"

// ChatSessionStatus represents the status of a chat session
type ChatSessionStatus string

const (
	ChatStatusAIActive ChatSessionStatus = "ai_active" // AI is handling the chat
	ChatStatusWaiting  ChatSessionStatus = "waiting"   // Waiting for human agent
	ChatStatusActive   ChatSessionStatus = "active"    // Human agent is handling
	ChatStatusEnded    ChatSessionStatus = "ended"     // Chat ended
)

// ChatSession represents a livechat session
type ChatSession struct {
	ID                   string            `json:"id"`
	TenantID             string            `json:"tenantId"`
	SchoolID             string            `json:"schoolId"`
	ThreadID             string            `json:"threadId"`
	SchoolContactID      string            `json:"schoolContactId"`
	SchoolContactName    string            `json:"schoolContactName"`
	AssignedAgentID      *string           `json:"assignedAgentId,omitempty"`
	AssignedAgentName    *string           `json:"assignedAgentName,omitempty"`
	Status               ChatSessionStatus `json:"status"`
	QueuePosition        *int              `json:"queuePosition,omitempty"`
	StartedAt            time.Time         `json:"startedAt"`
	AgentJoinedAt        *time.Time        `json:"agentJoinedAt,omitempty"`
	EndedAt              *time.Time        `json:"endedAt,omitempty"`
	FirstResponseSeconds *int              `json:"firstResponseSeconds,omitempty"`
	TotalMessages        int               `json:"totalMessages"`
	Rating               *int              `json:"rating,omitempty"`
	Feedback             *string           `json:"feedback,omitempty"`
	CreatedAt            time.Time         `json:"createdAt"`
	UpdatedAt            time.Time         `json:"updatedAt"`

	// AI Support fields
	AIHandled         bool           `json:"aiHandled"`
	AIResolved        *bool          `json:"aiResolved,omitempty"`
	AITurns           int            `json:"aiTurns"`
	EscalationReason  *string        `json:"escalationReason,omitempty"`
	EscalationSummary map[string]any `json:"escalationSummary,omitempty"`
	IssueCategory     *string        `json:"issueCategory,omitempty"`
	IssueSeverity     *string        `json:"issueSeverity,omitempty"`
	CollectedInfo     map[string]any `json:"collectedInfo,omitempty"`

	// Associated thread (for API responses)
	Thread *MessageThread `json:"thread,omitempty"`
}

// AgentAvailability represents an agent's availability status
type AgentAvailability struct {
	TenantID           string    `json:"tenantId"`
	UserID             string    `json:"userId"`
	IsAvailable        bool      `json:"isAvailable"`
	MaxConcurrentChats int       `json:"maxConcurrentChats"`
	CurrentChatCount   int       `json:"currentChatCount"`
	LastSeenAt         time.Time `json:"lastSeenAt"`
	UpdatedAt          time.Time `json:"updatedAt"`

	// Denormalized for display
	UserName string `json:"userName,omitempty"`
}

// StartSessionRequest represents the request to start a chat session
type StartSessionRequest struct {
	Subject    string  `json:"subject,omitempty"`
	IncidentID *string `json:"incidentId,omitempty"`
}

// StartSessionResponse represents the response when starting a session
type StartSessionResponse struct {
	Session       ChatSession   `json:"session"`
	Thread        MessageThread `json:"thread"`
	QueuePosition *int          `json:"queuePosition,omitempty"`
}

// QueuePositionResponse represents the queue position for a session
type QueuePositionResponse struct {
	Position             int `json:"position"`
	EstimatedWaitMinutes int `json:"estimatedWaitMinutes"`
}

// EndSessionRequest represents the request to end a chat session
type EndSessionRequest struct {
	Rating   *int    `json:"rating,omitempty"`
	Feedback *string `json:"feedback,omitempty"`
}

// AcceptChatResponse represents the response when accepting a chat
type AcceptChatResponse struct {
	Session ChatSession   `json:"session"`
	Thread  MessageThread `json:"thread"`
}

// TransferChatRequest represents the request to transfer a chat
type TransferChatRequest struct {
	TargetAgentID string  `json:"targetAgentId"`
	Reason        *string `json:"reason,omitempty"`
}

// SetAvailabilityRequest represents the request to set agent availability
type SetAvailabilityRequest struct {
	Available          bool `json:"available"`
	MaxConcurrentChats *int `json:"maxConcurrentChats,omitempty"`
}

// ChatQueueItem represents an item in the chat queue
type ChatQueueItem struct {
	SessionID         string    `json:"sessionId"`
	SchoolID          string    `json:"schoolId"`
	SchoolContactName string    `json:"schoolContactName"`
	Subject           string    `json:"subject"`
	WaitingTime       int       `json:"waitingTimeSeconds"`
	QueuePosition     int       `json:"queuePosition"`
	StartedAt         time.Time `json:"startedAt"`
}

// ChatQueueResponse represents the chat queue
type ChatQueueResponse struct {
	Items           []ChatQueueItem `json:"items"`
	TotalWaiting    int             `json:"totalWaiting"`
	AvailableAgents int             `json:"availableAgents"`
}

// ActiveChatItem represents an active chat for an agent
type ActiveChatItem struct {
	Session       ChatSession   `json:"session"`
	Thread        MessageThread `json:"thread"`
	UnreadCount   int           `json:"unreadCount"`
	LastMessageAt *time.Time    `json:"lastMessageAt,omitempty"`
}

// ActiveChatsResponse represents active chats for an agent
type ActiveChatsResponse struct {
	Items []ActiveChatItem `json:"items"`
	Total int              `json:"total"`
}

// ChatMetrics represents chat metrics for analytics
type ChatMetrics struct {
	TotalSessions       int     `json:"totalSessions"`
	AverageWaitTime     float64 `json:"averageWaitTimeSeconds"`
	AverageResponseTime float64 `json:"averageResponseTimeSeconds"`
	AverageRating       float64 `json:"averageRating"`
	SessionsWithRating  int     `json:"sessionsWithRating"`
	ActiveSessions      int     `json:"activeSessions"`
	WaitingSessions     int     `json:"waitingSessions"`
	EndedSessions       int     `json:"endedSessions"`
}

// AI Chat types

// AIChatMessageRequest represents a message sent to the AI
type AIChatMessageRequest struct {
	Content      string  `json:"content"`
	DeviceSerial *string `json:"deviceSerial,omitempty"`
}

// AIChatMessageResponse represents the AI's response
type AIChatMessageResponse struct {
	Message        Message           `json:"message"`
	AITyping       bool              `json:"aiTyping"`
	ShouldEscalate bool              `json:"shouldEscalate"`
	EscalateReason *string           `json:"escalateReason,omitempty"`
	SessionStatus  ChatSessionStatus `json:"sessionStatus"`
}

// AIEscalationRequest represents a request to escalate to human agent
type AIEscalationRequest struct {
	Reason string `json:"reason,omitempty"`
}

// AIEscalationResponse represents the response after escalation
type AIEscalationResponse struct {
	Session       ChatSession `json:"session"`
	QueuePosition *int        `json:"queuePosition,omitempty"`
}

// AIConversationContext represents the context for agent handoff
type AIConversationContext struct {
	SessionID           string         `json:"sessionId"`
	TurnCount           int            `json:"turnCount"`
	Category            *string        `json:"category,omitempty"`
	Severity            *string        `json:"severity,omitempty"`
	EscalationReason    *string        `json:"escalationReason,omitempty"`
	Summary             map[string]any `json:"summary,omitempty"`
	CollectedInfo       map[string]any `json:"collectedInfo,omitempty"`
	ConversationHistory []Message      `json:"conversationHistory"`
	DeviceContext       map[string]any `json:"deviceContext,omitempty"`
	SchoolContext       map[string]any `json:"schoolContext,omitempty"`
}
