package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/edvirons/ssp/ims/internal/middleware"
	"github.com/edvirons/ssp/ims/internal/models"
	"github.com/edvirons/ssp/ims/internal/store"
	"github.com/edvirons/ssp/ims/internal/ws"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// LivechatHandler handles livechat endpoints
type LivechatHandler struct {
	log *zap.Logger
	pg  *store.Postgres
	hub *ws.Hub
}

// NewLivechatHandler creates a new livechat handler
func NewLivechatHandler(log *zap.Logger, pg *store.Postgres, hub *ws.Hub) *LivechatHandler {
	return &LivechatHandler{log: log, pg: pg, hub: hub}
}

// StartSession handles POST /v1/chat/sessions
func (h *LivechatHandler) StartSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID := middleware.TenantID(ctx)
	userID := middleware.UserID(ctx)
	userName := middleware.UserName(ctx)
	roles := middleware.Roles(ctx)

	// Only school contacts can start chat sessions
	if !h.isSchoolContact(roles) {
		http.Error(w, "only school contacts can start chat sessions", http.StatusForbidden)
		return
	}

	// Check if user already has an active session
	existingSession, _ := h.pg.ChatSessions().GetActiveSessionForUser(ctx, tenantID, userID)
	if existingSession.ID != "" {
		// Return existing session
		thread, _ := h.pg.Messaging().GetThreadByID(ctx, tenantID, existingSession.ThreadID)
		position, _ := h.pg.ChatSessions().GetQueuePosition(ctx, tenantID, existingSession.ID)
		existingSession.QueuePosition = &position

		writeJSON(w, http.StatusOK, models.StartSessionResponse{
			Session:       existingSession,
			Thread:        thread,
			QueuePosition: &position,
		})
		return
	}

	var req models.StartSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// Allow empty body
		req = models.StartSessionRequest{}
	}

	// Get school ID
	schools := middleware.AssignedSchools(ctx)
	if len(schools) == 0 {
		http.Error(w, "no school assigned", http.StatusBadRequest)
		return
	}
	schoolID := schools[0]

	now := time.Now().UTC()

	// Create thread for the chat session
	subject := req.Subject
	if subject == "" {
		subject = "Live Chat - " + userName
	}

	thread := models.MessageThread{
		ID:                 store.NewID("thr"),
		TenantID:           tenantID,
		SchoolID:           schoolID,
		Subject:            subject,
		ThreadType:         models.ThreadTypeLivechat,
		Status:             models.ThreadStatusOpen,
		IncidentID:         req.IncidentID,
		CreatedBy:          userID,
		CreatedByRole:      "ssp_school_contact",
		CreatedByName:      userName,
		MessageCount:       0,
		UnreadCountSchool:  0,
		UnreadCountSupport: 0,
		LastMessageAt:      nil,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	if err := h.pg.Messaging().CreateThread(ctx, thread); err != nil {
		h.log.Error("failed to create chat thread", zap.Error(err))
		http.Error(w, "failed to start session", http.StatusInternalServerError)
		return
	}

	// Create chat session - Start in AI mode by default
	session := models.ChatSession{
		ID:                store.NewID("chat"),
		TenantID:          tenantID,
		SchoolID:          schoolID,
		ThreadID:          thread.ID,
		SchoolContactID:   userID,
		SchoolContactName: userName,
		Status:            models.ChatStatusAIActive, // AI-first mode
		AIHandled:         true,
		AITurns:           0,
		StartedAt:         now,
		TotalMessages:     0,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	if err := h.pg.ChatSessions().CreateSession(ctx, session); err != nil {
		h.log.Error("failed to create chat session", zap.Error(err))
		http.Error(w, "failed to start session", http.StatusInternalServerError)
		return
	}

	// Add user as participant
	h.pg.Messaging().AddThreadParticipant(ctx, models.ThreadParticipant{
		ThreadID: thread.ID,
		UserID:   userID,
		UserName: userName,
		UserRole: "ssp_school_contact",
		JoinedAt: now,
	})

	// Send AI welcome message
	welcomeMessage := models.Message{
		ID:          store.NewID("msg"),
		TenantID:    tenantID,
		ThreadID:    thread.ID,
		SenderID:    "ai_assistant",
		SenderName:  "ESSP Support Assistant",
		SenderRole:  "ai",
		Content:     "Hello! I'm ESSP Support Assistant. I'm here to help you with any device or technology issues.\n\nHow can I assist you today?",
		ContentType: models.ContentTypeText,
		CreatedAt:   now,
	}

	if err := h.pg.Messaging().CreateMessage(ctx, welcomeMessage); err != nil {
		h.log.Warn("failed to save AI welcome message", zap.Error(err))
	}

	// Broadcast welcome message
	if h.hub != nil {
		h.hub.Broadcast(tenantID, &ws.Message{
			Type: ws.MessageTypeChatMessage,
			Payload: map[string]any{
				"threadId": thread.ID,
				"message":  welcomeMessage,
			},
		})
	}

	writeJSON(w, http.StatusOK, models.StartSessionResponse{
		Session:       session,
		Thread:        thread,
		QueuePosition: nil, // No queue position in AI mode
	})
}

// GetQueuePosition handles GET /v1/chat/sessions/{id}/queue
func (h *LivechatHandler) GetQueuePosition(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID := middleware.TenantID(ctx)
	sessionID := chi.URLParam(r, "id")

	session, err := h.pg.ChatSessions().GetSessionByID(ctx, tenantID, sessionID)
	if err != nil {
		http.Error(w, "session not found", http.StatusNotFound)
		return
	}

	if session.Status != models.ChatStatusWaiting {
		writeJSON(w, http.StatusOK, models.QueuePositionResponse{
			Position:             0,
			EstimatedWaitMinutes: 0,
		})
		return
	}

	position, err := h.pg.ChatSessions().GetQueuePosition(ctx, tenantID, sessionID)
	if err != nil {
		h.log.Error("failed to get queue position", zap.Error(err))
		http.Error(w, "failed to get queue position", http.StatusInternalServerError)
		return
	}

	// Estimate wait time (roughly 2 minutes per position)
	estimatedWait := position * 2

	writeJSON(w, http.StatusOK, models.QueuePositionResponse{
		Position:             position,
		EstimatedWaitMinutes: estimatedWait,
	})
}

// EndSession handles POST /v1/chat/sessions/{id}/end
func (h *LivechatHandler) EndSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID := middleware.TenantID(ctx)
	userID := middleware.UserID(ctx)
	sessionID := chi.URLParam(r, "id")

	var req models.EndSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// Allow empty body
		req = models.EndSessionRequest{}
	}

	session, err := h.pg.ChatSessions().GetSessionByID(ctx, tenantID, sessionID)
	if err != nil {
		http.Error(w, "session not found", http.StatusNotFound)
		return
	}

	// Check if user can end this session (participant or admin)
	if session.SchoolContactID != userID && (session.AssignedAgentID == nil || *session.AssignedAgentID != userID) {
		roles := middleware.Roles(ctx)
		if !h.isAdmin(roles) && !h.isSupportAgent(roles) {
			http.Error(w, "access denied", http.StatusForbidden)
			return
		}
	}

	if err := h.pg.ChatSessions().EndSession(ctx, tenantID, sessionID, req.Rating, req.Feedback); err != nil {
		h.log.Error("failed to end session", zap.Error(err))
		http.Error(w, "failed to end session", http.StatusInternalServerError)
		return
	}

	// Decrement agent chat count if assigned
	if session.AssignedAgentID != nil {
		h.pg.ChatSessions().DecrementAgentChatCount(ctx, tenantID, *session.AssignedAgentID)
	}

	// Close the thread
	h.pg.Messaging().UpdateThreadStatus(ctx, tenantID, session.ThreadID, models.ThreadStatusClosed)

	// Update queue positions
	h.pg.ChatSessions().UpdateQueuePositions(ctx, tenantID)

	// Broadcast session ended
	h.broadcastSessionUpdate(tenantID, session.ID, models.ChatStatusEnded, nil, nil)

	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

// AcceptChat handles POST /v1/chat/accept
func (h *LivechatHandler) AcceptChat(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID := middleware.TenantID(ctx)
	userID := middleware.UserID(ctx)
	userName := middleware.UserName(ctx)
	roles := middleware.Roles(ctx)

	// Only agents can accept chats
	if !h.isSupportAgent(roles) && !h.isAdmin(roles) {
		http.Error(w, "only agents can accept chats", http.StatusForbidden)
		return
	}

	// Check agent availability
	availability, _ := h.pg.ChatSessions().GetAgentAvailability(ctx, tenantID, userID)
	if availability.CurrentChatCount >= availability.MaxConcurrentChats {
		http.Error(w, "you have reached maximum concurrent chats", http.StatusBadRequest)
		return
	}

	// Get next session in queue
	session, err := h.pg.ChatSessions().GetNextInQueue(ctx, tenantID)
	if err != nil || session.ID == "" {
		http.Error(w, "no chats waiting", http.StatusNotFound)
		return
	}

	// Assign agent to session
	if err := h.pg.ChatSessions().AssignAgent(ctx, tenantID, session.ID, userID, userName); err != nil {
		h.log.Error("failed to assign agent", zap.Error(err))
		http.Error(w, "failed to accept chat", http.StatusInternalServerError)
		return
	}

	// Increment agent chat count
	h.pg.ChatSessions().IncrementAgentChatCount(ctx, tenantID, userID)

	// Add agent as participant
	now := time.Now().UTC()
	h.pg.Messaging().AddThreadParticipant(ctx, models.ThreadParticipant{
		ThreadID: session.ThreadID,
		UserID:   userID,
		UserName: userName,
		UserRole: h.getPrimaryRole(roles),
		JoinedAt: now,
	})

	// Update queue positions for remaining sessions
	h.pg.ChatSessions().UpdateQueuePositions(ctx, tenantID)

	// Get updated session
	session, _ = h.pg.ChatSessions().GetSessionByID(ctx, tenantID, session.ID)

	// Get thread
	thread, _ := h.pg.Messaging().GetThreadByID(ctx, tenantID, session.ThreadID)

	// Broadcast session update
	h.broadcastSessionUpdate(tenantID, session.ID, models.ChatStatusActive, &userID, &userName)

	// Send system message
	h.sendSystemMessage(ctx, tenantID, session.ThreadID, userName+" has joined the chat")

	writeJSON(w, http.StatusOK, models.AcceptChatResponse{
		Session: session,
		Thread:  thread,
	})
}

// TransferChat handles POST /v1/chat/sessions/{id}/transfer
func (h *LivechatHandler) TransferChat(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID := middleware.TenantID(ctx)
	userID := middleware.UserID(ctx)
	userName := middleware.UserName(ctx)
	roles := middleware.Roles(ctx)
	sessionID := chi.URLParam(r, "id")

	// Only agents can transfer chats
	if !h.isSupportAgent(roles) && !h.isAdmin(roles) {
		http.Error(w, "only agents can transfer chats", http.StatusForbidden)
		return
	}

	var req models.TransferChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if req.TargetAgentID == "" {
		http.Error(w, "targetAgentId is required", http.StatusBadRequest)
		return
	}

	session, err := h.pg.ChatSessions().GetSessionByID(ctx, tenantID, sessionID)
	if err != nil {
		http.Error(w, "session not found", http.StatusNotFound)
		return
	}

	// Check if current user is the assigned agent
	if session.AssignedAgentID == nil || *session.AssignedAgentID != userID {
		if !h.isAdmin(roles) {
			http.Error(w, "you are not the assigned agent", http.StatusForbidden)
			return
		}
	}

	// Check target agent availability
	targetAvailability, _ := h.pg.ChatSessions().GetAgentAvailability(ctx, tenantID, req.TargetAgentID)
	if !targetAvailability.IsAvailable || targetAvailability.CurrentChatCount >= targetAvailability.MaxConcurrentChats {
		http.Error(w, "target agent is not available", http.StatusBadRequest)
		return
	}

	// Get target agent name (in production, this would come from user service)
	targetAgentName := "Support Agent"

	// Transfer the session
	if err := h.pg.ChatSessions().TransferSession(ctx, tenantID, sessionID, req.TargetAgentID, targetAgentName); err != nil {
		h.log.Error("failed to transfer session", zap.Error(err))
		http.Error(w, "failed to transfer chat", http.StatusInternalServerError)
		return
	}

	// Update chat counts (best-effort, errors logged but not blocking)
	_ = h.pg.ChatSessions().DecrementAgentChatCount(ctx, tenantID, userID)
	_ = h.pg.ChatSessions().IncrementAgentChatCount(ctx, tenantID, req.TargetAgentID)

	// Add new agent as participant (best-effort)
	now := time.Now().UTC()
	_ = h.pg.Messaging().AddThreadParticipant(ctx, models.ThreadParticipant{
		ThreadID: session.ThreadID,
		UserID:   req.TargetAgentID,
		UserName: targetAgentName,
		UserRole: "ssp_support_agent",
		JoinedAt: now,
	})

	// Send system message
	reason := ""
	if req.Reason != nil {
		reason = " Reason: " + *req.Reason
	}
	h.sendSystemMessage(ctx, tenantID, session.ThreadID, userName+" transferred the chat to "+targetAgentName+"."+reason)

	// Broadcast session update
	h.broadcastSessionUpdate(tenantID, sessionID, models.ChatStatusActive, &req.TargetAgentID, &targetAgentName)

	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

// SetAvailability handles PUT /v1/chat/availability
func (h *LivechatHandler) SetAvailability(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID := middleware.TenantID(ctx)
	userID := middleware.UserID(ctx)
	roles := middleware.Roles(ctx)

	// Only agents can set availability
	if !h.isSupportAgent(roles) && !h.isAdmin(roles) {
		http.Error(w, "only agents can set availability", http.StatusForbidden)
		return
	}

	var req models.SetAvailabilityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	maxChats := 3
	if req.MaxConcurrentChats != nil {
		maxChats = *req.MaxConcurrentChats
	} else {
		// Get existing max chats setting
		existing, _ := h.pg.ChatSessions().GetAgentAvailability(ctx, tenantID, userID)
		if existing.MaxConcurrentChats > 0 {
			maxChats = existing.MaxConcurrentChats
		}
	}

	if err := h.pg.ChatSessions().SetAgentAvailability(ctx, tenantID, userID, req.Available, maxChats); err != nil {
		h.log.Error("failed to set availability", zap.Error(err))
		http.Error(w, "failed to set availability", http.StatusInternalServerError)
		return
	}

	// Broadcast presence update
	h.broadcastPresenceUpdate(tenantID, userID, req.Available)

	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "available": req.Available})
}

// GetAvailability handles GET /v1/chat/availability
func (h *LivechatHandler) GetAvailability(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID := middleware.TenantID(ctx)
	userID := middleware.UserID(ctx)

	availability, err := h.pg.ChatSessions().GetAgentAvailability(ctx, tenantID, userID)
	if err != nil {
		h.log.Error("failed to get availability", zap.Error(err))
		http.Error(w, "failed to get availability", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, availability)
}

// GetQueue handles GET /v1/chat/queue
func (h *LivechatHandler) GetQueue(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID := middleware.TenantID(ctx)
	roles := middleware.Roles(ctx)

	// Only agents can view queue
	if !h.isSupportAgent(roles) && !h.isAdmin(roles) {
		http.Error(w, "access denied", http.StatusForbidden)
		return
	}

	sessions, err := h.pg.ChatSessions().GetWaitingSessions(ctx, tenantID)
	if err != nil {
		h.log.Error("failed to get queue", zap.Error(err))
		http.Error(w, "failed to get queue", http.StatusInternalServerError)
		return
	}

	now := time.Now().UTC()
	var items []models.ChatQueueItem
	for i, s := range sessions {
		thread, _ := h.pg.Messaging().GetThreadByID(ctx, tenantID, s.ThreadID)
		items = append(items, models.ChatQueueItem{
			SessionID:         s.ID,
			SchoolID:          s.SchoolID,
			SchoolContactName: s.SchoolContactName,
			Subject:           thread.Subject,
			WaitingTime:       int(now.Sub(s.StartedAt).Seconds()),
			QueuePosition:     i + 1,
			StartedAt:         s.StartedAt,
		})
	}

	// Get available agents count
	agents, _ := h.pg.ChatSessions().GetAvailableAgents(ctx, tenantID)

	writeJSON(w, http.StatusOK, models.ChatQueueResponse{
		Items:           items,
		TotalWaiting:    len(items),
		AvailableAgents: len(agents),
	})
}

// GetActiveChats handles GET /v1/chat/active
func (h *LivechatHandler) GetActiveChats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID := middleware.TenantID(ctx)
	userID := middleware.UserID(ctx)
	roles := middleware.Roles(ctx)

	// Only agents can view their active chats
	if !h.isSupportAgent(roles) && !h.isAdmin(roles) {
		http.Error(w, "access denied", http.StatusForbidden)
		return
	}

	sessions, err := h.pg.ChatSessions().GetActiveSessionsForAgent(ctx, tenantID, userID)
	if err != nil {
		h.log.Error("failed to get active chats", zap.Error(err))
		http.Error(w, "failed to get active chats", http.StatusInternalServerError)
		return
	}

	var items []models.ActiveChatItem
	for _, s := range sessions {
		thread, _ := h.pg.Messaging().GetThreadByID(ctx, tenantID, s.ThreadID)
		items = append(items, models.ActiveChatItem{
			Session:       s,
			Thread:        thread,
			UnreadCount:   thread.UnreadCountSupport,
			LastMessageAt: thread.LastMessageAt,
		})
	}

	writeJSON(w, http.StatusOK, models.ActiveChatsResponse{
		Items: items,
		Total: len(items),
	})
}

// GetChatMetrics handles GET /v1/chat/metrics
func (h *LivechatHandler) GetChatMetrics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID := middleware.TenantID(ctx)
	roles := middleware.Roles(ctx)

	// Only admins can view metrics
	if !h.isAdmin(roles) {
		http.Error(w, "access denied", http.StatusForbidden)
		return
	}

	// Parse date range
	now := time.Now().UTC()
	from := now.AddDate(0, 0, -30)
	to := now

	if fromStr := r.URL.Query().Get("from"); fromStr != "" {
		if t, err := time.Parse(time.RFC3339, fromStr); err == nil {
			from = t
		}
	}
	if toStr := r.URL.Query().Get("to"); toStr != "" {
		if t, err := time.Parse(time.RFC3339, toStr); err == nil {
			to = t
		}
	}

	metrics, err := h.pg.ChatSessions().GetChatMetrics(ctx, tenantID, from, to)
	if err != nil {
		h.log.Error("failed to get metrics", zap.Error(err))
		http.Error(w, "failed to get metrics", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, metrics)
}

// Helper methods

func (h *LivechatHandler) isSchoolContact(roles []string) bool {
	for _, r := range roles {
		if r == "ssp_school_contact" {
			return true
		}
	}
	return false
}

func (h *LivechatHandler) isSupportAgent(roles []string) bool {
	for _, r := range roles {
		if r == "ssp_support_agent" {
			return true
		}
	}
	return false
}

func (h *LivechatHandler) isAdmin(roles []string) bool {
	for _, r := range roles {
		if r == "ssp_admin" {
			return true
		}
	}
	return false
}

func (h *LivechatHandler) getPrimaryRole(roles []string) string {
	priority := []string{"ssp_admin", "ssp_support_agent", "ssp_school_contact"}
	for _, p := range priority {
		for _, r := range roles {
			if r == p {
				return r
			}
		}
	}
	if len(roles) > 0 {
		return roles[0]
	}
	return "unknown"
}

//nolint:unused // reserved for future auto-assignment feature
func (h *LivechatHandler) tryAutoAssign(ctx context.Context, tenantID, sessionID string) {
	agents, err := h.pg.ChatSessions().GetAvailableAgents(ctx, tenantID)
	if err != nil || len(agents) == 0 {
		return
	}

	// Pick the first available agent (least busy)
	agent := agents[0]

	// In production, we'd get the agent name from user service
	agentName := "Support Agent"

	if err := h.pg.ChatSessions().AssignAgent(ctx, tenantID, sessionID, agent.UserID, agentName); err != nil {
		h.log.Warn("auto-assign failed", zap.Error(err))
		return
	}

	_ = h.pg.ChatSessions().IncrementAgentChatCount(ctx, tenantID, agent.UserID)

	// Get session for thread ID
	session, _ := h.pg.ChatSessions().GetSessionByID(ctx, tenantID, sessionID)

	// Add agent as participant
	now := time.Now().UTC()
	_ = h.pg.Messaging().AddThreadParticipant(ctx, models.ThreadParticipant{
		ThreadID: session.ThreadID,
		UserID:   agent.UserID,
		UserName: agentName,
		UserRole: "ssp_support_agent",
		JoinedAt: now,
	})

	// Broadcast session update
	h.broadcastSessionUpdate(tenantID, sessionID, models.ChatStatusActive, &agent.UserID, &agentName)

	// Send system message
	h.sendSystemMessage(ctx, tenantID, session.ThreadID, agentName+" has joined the chat")
}

//nolint:unused // reserved for future notification feature
func (h *LivechatHandler) notifyAgentsNewChat(tenantID string, session models.ChatSession, thread models.MessageThread) {
	if h.hub == nil {
		return
	}

	h.hub.Broadcast(tenantID, &ws.Message{
		Type: ws.MessageTypeNewChatWaiting,
		Payload: map[string]any{
			"sessionId":         session.ID,
			"schoolContactName": session.SchoolContactName,
			"subject":           thread.Subject,
			"startedAt":         session.StartedAt,
		},
	})
}

func (h *LivechatHandler) broadcastSessionUpdate(tenantID, sessionID string, status models.ChatSessionStatus, agentID, agentName *string) {
	if h.hub == nil {
		return
	}

	h.hub.Broadcast(tenantID, &ws.Message{
		Type: ws.MessageTypeChatSessionUpdate,
		Payload: map[string]any{
			"sessionId": sessionID,
			"status":    status,
			"agentId":   agentID,
			"agentName": agentName,
		},
	})
}

func (h *LivechatHandler) broadcastPresenceUpdate(tenantID, userID string, online bool) {
	if h.hub == nil {
		return
	}

	status := "offline"
	if online {
		status = "online"
	}

	h.hub.Broadcast(tenantID, &ws.Message{
		Type: ws.MessageTypePresenceUpdate,
		Payload: map[string]any{
			"userId": userID,
			"status": status,
		},
	})
}

func (h *LivechatHandler) sendSystemMessage(ctx context.Context, tenantID, threadID, content string) {
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

	_ = h.pg.Messaging().UpdateThreadLastMessage(ctx, tenantID, threadID, false)

	// Broadcast
	if h.hub != nil {
		h.hub.Broadcast(tenantID, &ws.Message{
			Type: ws.MessageTypeChatMessage,
			Payload: map[string]any{
				"threadId": threadID,
				"message":  message,
			},
		})
	}
}
