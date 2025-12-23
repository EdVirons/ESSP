package claude

import "time"

// Message represents a single message in the conversation
type Message struct {
	Role    string `json:"role"` // "user" or "assistant"
	Content string `json:"content"`
}

// ChatRequest represents a request to Claude API
type ChatRequest struct {
	Model       string    `json:"model"`
	MaxTokens   int       `json:"max_tokens"`
	System      string    `json:"system,omitempty"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature,omitempty"`
}

// ChatResponse represents Claude API response
type ChatResponse struct {
	ID           string         `json:"id"`
	Type         string         `json:"type"`
	Role         string         `json:"role"`
	Content      []ContentBlock `json:"content"`
	Model        string         `json:"model"`
	StopReason   string         `json:"stop_reason"`
	StopSequence *string        `json:"stop_sequence"`
	Usage        Usage          `json:"usage"`
}

// ContentBlock represents a content block in the response
type ContentBlock struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
}

// Usage tracks token usage
type Usage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

// AIResponse represents processed Claude response for our handlers
type AIResponse struct {
	Content      string
	InputTokens  int
	OutputTokens int
	ResponseTime time.Duration
}

// EscalationDecision contains the AI's assessment of whether to escalate
type EscalationDecision struct {
	ShouldEscalate bool           `json:"should_escalate"`
	Reason         string         `json:"reason"`         // "user_request", "frustration", "complexity", "max_turns", "sensitive"
	Category       string         `json:"category"`       // "hardware", "software", "network", "account", "billing", "other"
	Severity       string         `json:"severity"`       // "low", "medium", "high", "critical"
	Summary        string         `json:"summary"`        // Brief summary for agent handoff
	CollectedInfo  map[string]any `json:"collected_info"` // Structured data collected during conversation
	Confidence     float64        `json:"confidence"`     // 0-1 confidence in the decision
}

// EscalationSignals captures signals detected during conversation analysis
type EscalationSignals struct {
	FrustrationScore    float64  `json:"frustration_score"`    // 0-1 based on sentiment
	ExplicitRequest     bool     `json:"explicit_request"`     // User asked for human
	SensitiveTopic      bool     `json:"sensitive_topic"`      // Billing, complaints, legal
	TechnicalComplexity float64  `json:"technical_complexity"` // 0-1 based on issue type
	UnresolvedIssue     bool     `json:"unresolved_issue"`     // AI couldn't help
	Keywords            []string `json:"keywords"`             // Detected escalation keywords
}

// SSOTContext contains device and school context for AI prompts
type SSOTContext struct {
	School  *SchoolContext  `json:"school,omitempty"`
	Device  *DeviceContext  `json:"device,omitempty"`
	History *HistoryContext `json:"history,omitempty"`
}

// SchoolContext contains school information
type SchoolContext struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	CountyName string `json:"county_name"`
	District   string `json:"district,omitempty"`
	State      string `json:"state,omitempty"`
	Type       string `json:"type,omitempty"` // elementary, middle, high
}

// DeviceContext contains device information
type DeviceContext struct {
	ID             string `json:"id"`
	SerialNumber   string `json:"serial_number"`
	Make           string `json:"make"`
	Model          string `json:"model"`
	DeviceType     string `json:"device_type"` // chromebook, laptop, tablet
	WarrantyStatus string `json:"warranty_status"`
	WarrantyExpiry string `json:"warranty_expiry,omitempty"`
	AssignedTo     string `json:"assigned_to,omitempty"`
	LastRepairDate string `json:"last_repair_date,omitempty"`
}

// HistoryContext contains recent support history
type HistoryContext struct {
	RecentIncidents  int      `json:"recent_incidents"`
	LastIncidentDate string   `json:"last_incident_date,omitempty"`
	CommonIssues     []string `json:"common_issues,omitempty"`
	TotalRepairs     int      `json:"total_repairs"`
}

// ConversationTurn represents a single turn in the AI conversation
type ConversationTurn struct {
	ID                    string            `json:"id"`
	TenantID              string            `json:"tenant_id"`
	SessionID             string            `json:"session_id"`
	TurnNumber            int               `json:"turn_number"`
	UserMessage           string            `json:"user_message"`
	AIResponse            string            `json:"ai_response"`
	InputTokens           int               `json:"input_tokens"`
	OutputTokens          int               `json:"output_tokens"`
	ResponseTimeMs        int               `json:"response_time_ms"`
	EscalationRecommended bool              `json:"escalation_recommended"`
	EscalationSignals     EscalationSignals `json:"escalation_signals"`
	ContextUsed           *SSOTContext      `json:"context_used,omitempty"`
	CreatedAt             time.Time         `json:"created_at"`
}

// AIMetrics tracks daily AI support metrics
type AIMetrics struct {
	TenantID             string         `json:"tenant_id"`
	Date                 string         `json:"date"` // YYYY-MM-DD
	TotalSessions        int            `json:"total_sessions"`
	AIResolved           int            `json:"ai_resolved"`
	EscalatedToHuman     int            `json:"escalated_to_human"`
	AvgTurnsToResolution float64        `json:"avg_turns_to_resolution"`
	AvgResponseTimeMs    int            `json:"avg_response_time_ms"`
	TotalInputTokens     int64          `json:"total_input_tokens"`
	TotalOutputTokens    int64          `json:"total_output_tokens"`
	EscalationReasons    map[string]int `json:"escalation_reasons"`
	IssueCategories      map[string]int `json:"issue_categories"`
}

// EscalationRules configures AI escalation behavior per tenant
type EscalationRules struct {
	TenantID               string            `json:"tenant_id"`
	MaxTurns               int               `json:"max_turns"`
	FrustrationThreshold   float64           `json:"frustration_threshold"`
	AutoEscalateCategories []string          `json:"auto_escalate_categories"`
	SensitiveKeywords      []string          `json:"sensitive_keywords"`
	Enabled                bool              `json:"enabled"`
	CustomPrompts          map[string]string `json:"custom_prompts"`
}

// ChatSessionStatus represents the status of a chat session
type ChatSessionStatus string

const (
	ChatStatusAIActive ChatSessionStatus = "ai_active"
	ChatStatusWaiting  ChatSessionStatus = "waiting"
	ChatStatusActive   ChatSessionStatus = "active"
	ChatStatusEnded    ChatSessionStatus = "ended"
)
