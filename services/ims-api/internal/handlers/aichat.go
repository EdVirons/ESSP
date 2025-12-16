package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/edvirons/ssp/ims/internal/claude"
	"github.com/edvirons/ssp/ims/internal/config"
	"github.com/edvirons/ssp/ims/internal/middleware"
	"github.com/edvirons/ssp/ims/internal/models"
	"github.com/edvirons/ssp/ims/internal/store"
	"github.com/edvirons/ssp/ims/internal/ws"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// AIChatHandler handles AI chat endpoints
type AIChatHandler struct {
	log            *zap.Logger
	pg             *store.Postgres
	hub            *ws.Hub
	claude         *claude.Client
	contextBuilder *claude.ContextBuilder
	escalation     *claude.EscalationAnalyzer
	cfg            config.Config
}

// NewAIChatHandler creates a new AI chat handler
func NewAIChatHandler(log *zap.Logger, pg *store.Postgres, hub *ws.Hub, cfg config.Config) *AIChatHandler {
	// Create Claude client
	claudeClient := claude.NewClient(claude.ClientConfig{
		APIKey:         cfg.ClaudeAPIKey,
		Model:          cfg.ClaudeModel,
		MaxTokens:      cfg.ClaudeMaxTokens,
		TimeoutSeconds: cfg.ClaudeTimeoutSeconds,
	}, log)

	// Create context builder
	contextBuilder := claude.NewContextBuilder(pg.RawPool(), log)

	// Create escalation analyzer
	escalation := claude.NewEscalationAnalyzer(cfg.AIMaxTurns, cfg.AIFrustrationThreshold)

	return &AIChatHandler{
		log:            log,
		pg:             pg,
		hub:            hub,
		claude:         claudeClient,
		contextBuilder: contextBuilder,
		escalation:     escalation,
		cfg:            cfg,
	}
}

// HandleAIMessage handles POST /v1/chat/ai/sessions/{id}/message
func (h *AIChatHandler) HandleAIMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID := middleware.TenantID(ctx)
	userID := middleware.UserID(ctx)
	userName := middleware.UserName(ctx)
	sessionID := chi.URLParam(r, "id")

	// Check if AI is enabled
	if !h.cfg.AIEnabled || !h.claude.IsEnabled() {
		http.Error(w, "AI support is not enabled", http.StatusServiceUnavailable)
		return
	}

	var req models.AIChatMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if req.Content == "" {
		http.Error(w, "content is required", http.StatusBadRequest)
		return
	}

	// Get session
	session, err := h.pg.ChatSessions().GetSessionByID(ctx, tenantID, sessionID)
	if err != nil {
		http.Error(w, "session not found", http.StatusNotFound)
		return
	}

	// Verify user owns this session
	if session.SchoolContactID != userID {
		http.Error(w, "access denied", http.StatusForbidden)
		return
	}

	// Check session is in AI mode
	if session.Status != models.ChatStatusAIActive {
		http.Error(w, "session is not in AI mode", http.StatusBadRequest)
		return
	}

	// Save user message
	userMessage := models.Message{
		ID:          store.NewID("msg"),
		TenantID:    tenantID,
		ThreadID:    session.ThreadID,
		SenderID:    userID,
		SenderName:  userName,
		SenderRole:  "ssp_school_contact",
		Content:     req.Content,
		ContentType: models.ContentTypeText,
		CreatedAt:   time.Now().UTC(),
	}

	if err := h.pg.Messaging().CreateMessage(ctx, userMessage); err != nil {
		h.log.Error("failed to save user message", zap.Error(err))
		http.Error(w, "failed to save message", http.StatusInternalServerError)
		return
	}

	// Broadcast user message via WebSocket
	h.broadcastMessage(tenantID, session.ThreadID, userMessage)

	// Analyze message for escalation signals
	signals := h.escalation.AnalyzeMessage(req.Content)

	// Get conversation history
	history, _ := h.pg.Messaging().GetThreadMessages(ctx, tenantID, session.ThreadID, 0, 20)
	conversationHistory := h.buildClaudeHistory(history)

	// Build SSOT context
	ssotContext, _ := h.contextBuilder.BuildContext(ctx, tenantID, session.SchoolID, req.DeviceSerial)

	// Build system prompt
	promptData := claude.PromptData{
		Context:    ssotContext,
		TurnNumber: session.AITurns + 1,
		MaxTurns:   h.cfg.AIMaxTurns,
		UserName:   userName,
		SessionID:  sessionID,
	}
	systemPrompt, err := claude.BuildSystemPrompt(promptData)
	if err != nil {
		h.log.Error("failed to build system prompt", zap.Error(err))
		systemPrompt = claude.SystemPromptTemplate
	}

	// Send typing indicator
	h.broadcastTyping(tenantID, session.ThreadID, true)

	// Call Claude API
	aiResponse, err := h.claude.ChatWithRetry(ctx, systemPrompt, conversationHistory, 2)
	if err != nil {
		h.log.Error("Claude API error", zap.Error(err))
		h.broadcastTyping(tenantID, session.ThreadID, false)
		// Return a fallback message
		aiResponse = &claude.AIResponse{
			Content: "I apologize, but I'm experiencing technical difficulties. Let me connect you with a support agent who can help.",
		}
		signals.UnresolvedIssue = true
	}

	// Stop typing indicator
	h.broadcastTyping(tenantID, session.ThreadID, false)

	// Parse AI decision from response
	decision, cleanContent := claude.ParseAIDecision(aiResponse.Content)

	// Check for escalation
	shouldEscalate, escalationReason := h.escalation.ShouldEscalate(signals, session.AITurns+1, decision)

	// Update session AI turns and collected info
	h.updateSessionAIData(ctx, tenantID, sessionID, decision, session.AITurns+1)

	// Save AI response message
	aiMessage := models.Message{
		ID:          store.NewID("msg"),
		TenantID:    tenantID,
		ThreadID:    session.ThreadID,
		SenderID:    "ai_assistant",
		SenderName:  "ESSP Support Assistant",
		SenderRole:  "ai",
		Content:     cleanContent,
		ContentType: models.ContentTypeText,
		CreatedAt:   time.Now().UTC(),
	}

	if err := h.pg.Messaging().CreateMessage(ctx, aiMessage); err != nil {
		h.log.Error("failed to save AI message", zap.Error(err))
	}

	// Broadcast AI message
	h.broadcastMessage(tenantID, session.ThreadID, aiMessage)

	// Handle escalation
	if shouldEscalate {
		h.handleEscalation(ctx, tenantID, sessionID, escalationReason, signals, decision, conversationHistory)
	}

	// Check if AI resolved the issue
	if decision != nil && decision.Resolved {
		resolved := true
		h.pg.ChatSessions().UpdateAIResolution(ctx, tenantID, sessionID, resolved)
	}

	// Log conversation turn
	h.logConversationTurn(ctx, tenantID, sessionID, session.AITurns+1, req.Content, aiResponse, signals, ssotContext)

	// Get updated session status
	updatedSession, _ := h.pg.ChatSessions().GetSessionByID(ctx, tenantID, sessionID)

	writeJSON(w, http.StatusOK, models.AIChatMessageResponse{
		Message:        aiMessage,
		AITyping:       false,
		ShouldEscalate: shouldEscalate,
		EscalateReason: &escalationReason,
		SessionStatus:  updatedSession.Status,
	})
}

// RequestEscalation handles POST /v1/chat/ai/sessions/{id}/escalate
func (h *AIChatHandler) RequestEscalation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID := middleware.TenantID(ctx)
	userID := middleware.UserID(ctx)
	sessionID := chi.URLParam(r, "id")

	var req models.AIEscalationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		req = models.AIEscalationRequest{}
	}

	// Get session
	session, err := h.pg.ChatSessions().GetSessionByID(ctx, tenantID, sessionID)
	if err != nil {
		http.Error(w, "session not found", http.StatusNotFound)
		return
	}

	// Verify user owns this session
	if session.SchoolContactID != userID {
		http.Error(w, "access denied", http.StatusForbidden)
		return
	}

	// Check session is in AI mode
	if session.Status != models.ChatStatusAIActive {
		http.Error(w, "session is not in AI mode", http.StatusBadRequest)
		return
	}

	// Build escalation context
	history, _ := h.pg.Messaging().GetThreadMessages(ctx, tenantID, session.ThreadID, 0, 20)
	conversationHistory := h.buildClaudeHistory(history)

	reason := "user_request"
	if req.Reason != "" {
		reason = req.Reason
	}

	h.handleEscalation(ctx, tenantID, sessionID, reason, claude.EscalationSignals{ExplicitRequest: true}, nil, conversationHistory)

	// Send system message
	h.sendSystemMessage(ctx, tenantID, session.ThreadID, "Connecting you with a support agent...")

	// Get updated session
	updatedSession, _ := h.pg.ChatSessions().GetSessionByID(ctx, tenantID, sessionID)
	position, _ := h.pg.ChatSessions().GetQueuePosition(ctx, tenantID, sessionID)

	writeJSON(w, http.StatusOK, models.AIEscalationResponse{
		Session:       updatedSession,
		QueuePosition: &position,
	})
}

// GetConversationContext handles GET /v1/chat/ai/sessions/{id}/context
func (h *AIChatHandler) GetConversationContext(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID := middleware.TenantID(ctx)
	roles := middleware.Roles(ctx)
	sessionID := chi.URLParam(r, "id")

	// Only support agents can view context
	if !h.isSupportAgent(roles) && !h.isAdmin(roles) {
		http.Error(w, "access denied", http.StatusForbidden)
		return
	}

	// Get session
	session, err := h.pg.ChatSessions().GetSessionByID(ctx, tenantID, sessionID)
	if err != nil {
		http.Error(w, "session not found", http.StatusNotFound)
		return
	}

	// Get conversation history
	messages, _ := h.pg.Messaging().GetThreadMessages(ctx, tenantID, session.ThreadID, 0, 100)

	// Build SSOT context
	ssotContext, _ := h.contextBuilder.BuildContext(ctx, tenantID, session.SchoolID, nil)

	context := models.AIConversationContext{
		SessionID:           sessionID,
		TurnCount:           session.AITurns,
		Category:            session.IssueCategory,
		Severity:            session.IssueSeverity,
		EscalationReason:    session.EscalationReason,
		Summary:             session.EscalationSummary,
		CollectedInfo:       session.CollectedInfo,
		ConversationHistory: messages,
	}

	if ssotContext != nil {
		if ssotContext.School != nil {
			context.SchoolContext = map[string]any{
				"name":       ssotContext.School.Name,
				"countyName": ssotContext.School.CountyName,
				"district":   ssotContext.School.District,
			}
		}
		if ssotContext.Device != nil {
			context.DeviceContext = map[string]any{
				"make":           ssotContext.Device.Make,
				"model":          ssotContext.Device.Model,
				"serialNumber":   ssotContext.Device.SerialNumber,
				"warrantyStatus": ssotContext.Device.WarrantyStatus,
			}
		}
	}

	writeJSON(w, http.StatusOK, context)
}

// Helper methods

func (h *AIChatHandler) buildClaudeHistory(messages []models.Message) []claude.Message {
	var history []claude.Message
	for _, msg := range messages {
		role := "user"
		if msg.SenderRole == "ai" || msg.SenderID == "ai_assistant" {
			role = "assistant"
		}
		// Skip system messages
		if msg.SenderRole == "system" {
			continue
		}
		history = append(history, claude.Message{
			Role:    role,
			Content: msg.Content,
		})
	}
	return history
}

func (h *AIChatHandler) handleEscalation(ctx context.Context, tenantID, sessionID, reason string, signals claude.EscalationSignals, decision *claude.AIDecisionData, history []claude.Message) {
	// Build escalation summary
	summary := claude.BuildEscalationSummary(signals, decision, 0, history)

	// Update session status to waiting
	h.pg.ChatSessions().EscalateToHuman(ctx, tenantID, sessionID, reason, summary)

	// Update queue positions
	h.pg.ChatSessions().UpdateQueuePositions(ctx, tenantID)

	// Get session for thread ID
	session, _ := h.pg.ChatSessions().GetSessionByID(ctx, tenantID, sessionID)

	// Notify agents of new chat waiting
	if h.hub != nil {
		h.hub.Broadcast(tenantID, &ws.Message{
			Type: ws.MessageTypeNewChatWaiting,
			Payload: map[string]any{
				"sessionId":         sessionID,
				"schoolContactName": session.SchoolContactName,
				"escalatedFromAI":   true,
				"escalationReason":  reason,
				"aiSummary":         summary,
			},
		})
	}

	// Broadcast status update
	h.broadcastSessionUpdate(tenantID, sessionID, models.ChatStatusWaiting)
}

func (h *AIChatHandler) updateSessionAIData(ctx context.Context, tenantID, sessionID string, decision *claude.AIDecisionData, turns int) {
	if decision == nil {
		h.pg.ChatSessions().IncrementAITurns(ctx, tenantID, sessionID)
		return
	}

	h.pg.ChatSessions().UpdateAISessionData(ctx, tenantID, sessionID, turns, decision.Category, decision.Severity, decision.CollectedInfo)
}

func (h *AIChatHandler) logConversationTurn(ctx context.Context, tenantID, sessionID string, turnNumber int, userMsg string, aiResp *claude.AIResponse, signals claude.EscalationSignals, ssotCtx *claude.SSOTContext) {
	turn := claude.ConversationTurn{
		ID:                    store.NewID("turn"),
		TenantID:              tenantID,
		SessionID:             sessionID,
		TurnNumber:            turnNumber,
		UserMessage:           userMsg,
		AIResponse:            aiResp.Content,
		InputTokens:           aiResp.InputTokens,
		OutputTokens:          aiResp.OutputTokens,
		ResponseTimeMs:        int(aiResp.ResponseTime.Milliseconds()),
		EscalationRecommended: signals.ExplicitRequest || signals.SensitiveTopic || signals.FrustrationScore > 0.7,
		EscalationSignals:     signals,
		ContextUsed:           ssotCtx,
		CreatedAt:             time.Now().UTC(),
	}

	h.pg.ChatSessions().LogAIConversationTurn(ctx, turn)
}

func (h *AIChatHandler) broadcastMessage(tenantID, threadID string, message models.Message) {
	if h.hub == nil {
		return
	}
	h.hub.Broadcast(tenantID, &ws.Message{
		Type: ws.MessageTypeChatMessage,
		Payload: map[string]any{
			"threadId": threadID,
			"message":  message,
		},
	})
}

func (h *AIChatHandler) broadcastTyping(tenantID, threadID string, isTyping bool) {
	if h.hub == nil {
		return
	}
	h.hub.Broadcast(tenantID, &ws.Message{
		Type: ws.MessageTypeTypingIndicator,
		Payload: map[string]any{
			"threadId": threadID,
			"userId":   "ai_assistant",
			"userName": "ESSP Support Assistant",
			"isTyping": isTyping,
		},
	})
}

func (h *AIChatHandler) broadcastSessionUpdate(tenantID, sessionID string, status models.ChatSessionStatus) {
	if h.hub == nil {
		return
	}
	h.hub.Broadcast(tenantID, &ws.Message{
		Type: ws.MessageTypeChatSessionUpdate,
		Payload: map[string]any{
			"sessionId": sessionID,
			"status":    status,
		},
	})
}

func (h *AIChatHandler) sendSystemMessage(ctx context.Context, tenantID, threadID, content string) {
	message := models.Message{
		ID:          store.NewID("msg"),
		TenantID:    tenantID,
		ThreadID:    threadID,
		SenderID:    "system",
		SenderName:  "System",
		SenderRole:  "system",
		Content:     content,
		ContentType: models.ContentTypeSystem,
		CreatedAt:   time.Now().UTC(),
	}

	if err := h.pg.Messaging().CreateMessage(ctx, message); err != nil {
		h.log.Warn("failed to send system message", zap.Error(err))
		return
	}

	h.broadcastMessage(tenantID, threadID, message)
}

func (h *AIChatHandler) isSupportAgent(roles []string) bool {
	for _, r := range roles {
		if r == "ssp_support_agent" {
			return true
		}
	}
	return false
}

func (h *AIChatHandler) isAdmin(roles []string) bool {
	for _, r := range roles {
		if r == "ssp_admin" {
			return true
		}
	}
	return false
}
