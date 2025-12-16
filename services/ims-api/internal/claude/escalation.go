package claude

import (
	"encoding/json"
	"regexp"
	"strings"
)

// EscalationAnalyzer analyzes messages for escalation triggers
type EscalationAnalyzer struct {
	maxTurns             int
	frustrationThreshold float64
	sensitiveCategories  []string
	sensitiveKeywords    []string
}

// NewEscalationAnalyzer creates a new escalation analyzer
func NewEscalationAnalyzer(maxTurns int, frustrationThreshold float64) *EscalationAnalyzer {
	return &EscalationAnalyzer{
		maxTurns:             maxTurns,
		frustrationThreshold: frustrationThreshold,
		sensitiveCategories:  SensitiveCategories,
		sensitiveKeywords:    EscalationKeywords,
	}
}

// AnalyzeMessage analyzes a user message for escalation signals
func (ea *EscalationAnalyzer) AnalyzeMessage(message string) EscalationSignals {
	signals := EscalationSignals{}
	lowerMsg := strings.ToLower(message)

	// Check for explicit escalation request
	signals.ExplicitRequest = ea.containsEscalationRequest(lowerMsg)
	if signals.ExplicitRequest {
		signals.Keywords = append(signals.Keywords, "explicit_request")
	}

	// Check for sensitive topics
	signals.SensitiveTopic = ea.containsSensitiveTopic(lowerMsg)
	if signals.SensitiveTopic {
		signals.Keywords = append(signals.Keywords, "sensitive_topic")
	}

	// Calculate frustration score
	signals.FrustrationScore = ea.calculateFrustrationScore(message)
	if signals.FrustrationScore > ea.frustrationThreshold {
		signals.Keywords = append(signals.Keywords, "frustrated")
	}

	// Assess technical complexity
	signals.TechnicalComplexity = ea.assessTechnicalComplexity(lowerMsg)

	return signals
}

// ShouldEscalate determines if the conversation should be escalated
func (ea *EscalationAnalyzer) ShouldEscalate(signals EscalationSignals, turnNumber int, aiDecision *AIDecisionData) (bool, string) {
	// Explicit user request always escalates
	if signals.ExplicitRequest {
		return true, "user_request"
	}

	// Sensitive topics always escalate
	if signals.SensitiveTopic {
		return true, "sensitive"
	}

	// High frustration escalates
	if signals.FrustrationScore > ea.frustrationThreshold {
		return true, "frustration"
	}

	// Max turns reached
	if turnNumber >= ea.maxTurns {
		return true, "max_turns"
	}

	// Check AI's own recommendation
	if aiDecision != nil && aiDecision.Escalate {
		reason := aiDecision.EscalateReason
		if reason == "" {
			reason = "ai_recommendation"
		}
		return true, reason
	}

	// High complexity issues
	if signals.TechnicalComplexity > 0.8 {
		return true, "complexity"
	}

	return false, ""
}

// containsEscalationRequest checks for explicit human agent requests
func (ea *EscalationAnalyzer) containsEscalationRequest(msg string) bool {
	patterns := []string{
		`\bhuman\b`,
		`\bagent\b`,
		`\breal person\b`,
		`\bspeak to someone\b`,
		`\btalk to someone\b`,
		`\bsupervisor\b`,
		`\bmanager\b`,
		`\bescalate\b`,
	}

	for _, pattern := range patterns {
		if matched, _ := regexp.MatchString(pattern, msg); matched {
			return true
		}
	}
	return false
}

// containsSensitiveTopic checks for sensitive topics
func (ea *EscalationAnalyzer) containsSensitiveTopic(msg string) bool {
	sensitivePatterns := []string{
		`\bbilling\b`,
		`\brefund\b`,
		`\bpayment\b`,
		`\bcharge\b`,
		`\binvoice\b`,
		`\blawyer\b`,
		`\blegal\b`,
		`\blawsuit\b`,
		`\bcomplaint\b`,
		`\bharassment\b`,
	}

	for _, pattern := range sensitivePatterns {
		if matched, _ := regexp.MatchString(pattern, msg); matched {
			return true
		}
	}
	return false
}

// calculateFrustrationScore estimates user frustration from message
func (ea *EscalationAnalyzer) calculateFrustrationScore(msg string) float64 {
	score := 0.0
	lowerMsg := strings.ToLower(msg)

	// Check for frustration keywords
	for _, keyword := range FrustrationKeywords {
		if strings.Contains(lowerMsg, keyword) {
			score += 0.2
		}
	}

	// Check for excessive caps (shouting)
	capsCount := 0
	for _, c := range msg {
		if c >= 'A' && c <= 'Z' {
			capsCount++
		}
	}
	if len(msg) > 10 && float64(capsCount)/float64(len(msg)) > 0.5 {
		score += 0.3
	}

	// Check for excessive punctuation
	exclamations := strings.Count(msg, "!")
	questions := strings.Count(msg, "?")
	if exclamations > 2 || questions > 3 {
		score += 0.2
	}

	// Check for negative sentiment words
	negativeWords := []string{
		"not working", "doesn't work", "broken", "failed", "error",
		"problem", "issue", "wrong", "bad", "worst",
	}
	for _, word := range negativeWords {
		if strings.Contains(lowerMsg, word) {
			score += 0.1
		}
	}

	// Cap at 1.0
	if score > 1.0 {
		score = 1.0
	}

	return score
}

// assessTechnicalComplexity estimates issue complexity
func (ea *EscalationAnalyzer) assessTechnicalComplexity(msg string) float64 {
	complexity := 0.0

	// Hardware issues requiring physical inspection
	hardwareIndicators := []string{
		"screen crack", "broken screen", "physical damage",
		"won't turn on", "not charging", "battery", "keyboard broken",
		"hinge", "water damage", "dropped",
	}
	for _, indicator := range hardwareIndicators {
		if strings.Contains(msg, indicator) {
			complexity += 0.3
		}
	}

	// Network issues often need investigation
	networkIndicators := []string{
		"network", "wifi", "internet", "connection", "firewall",
		"dns", "proxy", "vpn",
	}
	for _, indicator := range networkIndicators {
		if strings.Contains(msg, indicator) {
			complexity += 0.1
		}
	}

	// Account/access issues may need admin intervention
	accountIndicators := []string{
		"locked out", "password reset", "can't login", "access denied",
		"permission", "admin", "blocked",
	}
	for _, indicator := range accountIndicators {
		if strings.Contains(msg, indicator) {
			complexity += 0.15
		}
	}

	// Cap at 1.0
	if complexity > 1.0 {
		complexity = 1.0
	}

	return complexity
}

// AIDecisionData represents the structured decision from AI response
type AIDecisionData struct {
	Category       string         `json:"category"`
	Severity       string         `json:"severity"`
	Escalate       bool           `json:"escalate"`
	EscalateReason string         `json:"escalate_reason"`
	CollectedInfo  map[string]any `json:"collected_info"`
	Resolved       bool           `json:"resolved"`
	NeedsMoreInfo  bool           `json:"needs_more_info"`
}

// ParseAIDecision extracts JSON decision data from AI response
func ParseAIDecision(response string) (*AIDecisionData, string) {
	// Look for JSON block in response
	jsonStart := strings.Index(response, "```json")
	if jsonStart == -1 {
		return nil, response
	}

	jsonEnd := strings.Index(response[jsonStart+7:], "```")
	if jsonEnd == -1 {
		return nil, response
	}

	jsonStr := response[jsonStart+7 : jsonStart+7+jsonEnd]
	jsonStr = strings.TrimSpace(jsonStr)

	var decision AIDecisionData
	if err := json.Unmarshal([]byte(jsonStr), &decision); err != nil {
		return nil, response
	}

	// Return clean message (without JSON block)
	cleanMsg := strings.TrimSpace(response[:jsonStart])
	if remaining := response[jsonStart+7+jsonEnd+3:]; len(strings.TrimSpace(remaining)) > 0 {
		cleanMsg += "\n" + strings.TrimSpace(remaining)
	}

	return &decision, cleanMsg
}

// BuildEscalationSummary creates a summary for agent handoff
func BuildEscalationSummary(
	signals EscalationSignals,
	decision *AIDecisionData,
	turnCount int,
	conversationHistory []Message,
) map[string]any {
	summary := map[string]any{
		"turn_count":         turnCount,
		"frustration_score":  signals.FrustrationScore,
		"escalation_signals": signals.Keywords,
	}

	if decision != nil {
		summary["category"] = decision.Category
		summary["severity"] = decision.Severity
		summary["ai_assessment"] = decision.EscalateReason
		summary["collected_info"] = decision.CollectedInfo
	}

	// Add brief conversation summary
	if len(conversationHistory) > 0 {
		var userMessages []string
		for _, msg := range conversationHistory {
			if msg.Role == "user" && len(msg.Content) > 0 {
				// Truncate long messages
				content := msg.Content
				if len(content) > 200 {
					content = content[:200] + "..."
				}
				userMessages = append(userMessages, content)
			}
		}
		summary["user_messages"] = userMessages
	}

	return summary
}
