package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/edvirons/ssp/ims/internal/middleware"
	"github.com/edvirons/ssp/ims/internal/models"
	"github.com/edvirons/ssp/ims/internal/store"
	"github.com/edvirons/ssp/ims/internal/ws"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// MessagingHandler handles messaging endpoints
type MessagingHandler struct {
	log *zap.Logger
	pg  *store.Postgres
	hub *ws.Hub
}

// NewMessagingHandler creates a new messaging handler
func NewMessagingHandler(log *zap.Logger, pg *store.Postgres, hub *ws.Hub) *MessagingHandler {
	return &MessagingHandler{log: log, pg: pg, hub: hub}
}

// ListThreads handles GET /v1/messages/threads
func (h *MessagingHandler) ListThreads(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID := middleware.TenantID(ctx)
	userID := middleware.UserID(ctx)
	roles := middleware.Roles(ctx)

	status := strings.TrimSpace(r.URL.Query().Get("status"))
	incidentID := strings.TrimSpace(r.URL.Query().Get("incidentId"))
	query := strings.TrimSpace(r.URL.Query().Get("q"))
	limit := parseLimit(r.URL.Query().Get("limit"), 50, 200)
	curT, curID, hasCur := decodeCursor(strings.TrimSpace(r.URL.Query().Get("cursor")))

	// Determine school ID filter based on role
	schoolID := ""
	if h.isSchoolContact(roles) {
		// School contacts only see their school's threads
		schools := middleware.AssignedSchools(ctx)
		if len(schools) > 0 {
			schoolID = schools[0]
		} else {
			writeJSON(w, http.StatusOK, models.ThreadsListResponse{Items: []models.MessageThread{}})
			return
		}
	}

	threads, next, err := h.pg.Messaging().ListThreads(ctx, store.ThreadListParams{
		TenantID:        tenantID,
		SchoolID:        schoolID,
		UserID:          userID,
		UserRoles:       roles,
		Status:          status,
		IncidentID:      incidentID,
		Query:           query,
		Limit:           limit,
		HasCursor:       hasCur,
		CursorTimestamp: curT,
		CursorID:        curID,
	})
	if err != nil {
		h.log.Error("failed to list threads", zap.Error(err))
		http.Error(w, "failed to list threads", http.StatusInternalServerError)
		return
	}

	// Fetch last message for each thread
	for i := range threads {
		if threads[i].MessageCount > 0 {
			lastMsg, err := h.pg.Messaging().GetLastMessage(ctx, threads[i].ID)
			if err == nil {
				threads[i].LastMessage = &lastMsg
			}
		}
	}

	writeJSON(w, http.StatusOK, models.ThreadsListResponse{
		Items:      threads,
		NextCursor: next,
	})
}

// GetThread handles GET /v1/messages/threads/{id}
func (h *MessagingHandler) GetThread(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID := middleware.TenantID(ctx)
	roles := middleware.Roles(ctx)
	threadID := chi.URLParam(r, "id")

	thread, err := h.pg.Messaging().GetThreadByID(ctx, tenantID, threadID)
	if err != nil {
		http.Error(w, "thread not found", http.StatusNotFound)
		return
	}

	// Check access for school contacts
	if h.isSchoolContact(roles) {
		schools := middleware.AssignedSchools(ctx)
		if !contains(schools, thread.SchoolID) {
			http.Error(w, "access denied", http.StatusForbidden)
			return
		}
	}

	// Get messages
	messages, _, err := h.pg.Messaging().ListMessages(ctx, store.MessageListParams{
		ThreadID: threadID,
		Limit:    100,
	})
	if err != nil {
		h.log.Error("failed to list messages", zap.Error(err))
		http.Error(w, "failed to get messages", http.StatusInternalServerError)
		return
	}

	// Get attachments for each message
	for i := range messages {
		attachments, _ := h.pg.Messaging().GetAttachmentsByMessageID(ctx, messages[i].ID)
		if len(attachments) > 0 {
			messages[i].Attachments = attachments
		}
	}

	// Get participants
	participants, _ := h.pg.Messaging().GetThreadParticipants(ctx, threadID)

	writeJSON(w, http.StatusOK, models.ThreadDetailResponse{
		Thread:       thread,
		Messages:     messages,
		Participants: participants,
	})
}

// CreateThread handles POST /v1/messages/threads
func (h *MessagingHandler) CreateThread(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID := middleware.TenantID(ctx)
	userID := middleware.UserID(ctx)
	userName := middleware.UserName(ctx)
	roles := middleware.Roles(ctx)

	var req models.CreateThreadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Subject) == "" {
		http.Error(w, "subject is required", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(req.InitialMessage) == "" {
		http.Error(w, "initialMessage is required", http.StatusBadRequest)
		return
	}

	// Determine school ID
	schoolID := req.SchoolID
	if h.isSchoolContact(roles) {
		schools := middleware.AssignedSchools(ctx)
		if len(schools) > 0 {
			schoolID = schools[0]
		} else {
			http.Error(w, "no school assigned", http.StatusBadRequest)
			return
		}
	}
	if schoolID == "" {
		http.Error(w, "schoolId is required", http.StatusBadRequest)
		return
	}

	// Determine primary role
	primaryRole := h.getPrimaryRole(roles)

	now := time.Now().UTC()

	// Create thread
	thread := models.MessageThread{
		ID:                 store.NewID("thr"),
		TenantID:           tenantID,
		SchoolID:           schoolID,
		Subject:            strings.TrimSpace(req.Subject),
		ThreadType:         models.ThreadTypeGeneral,
		Status:             models.ThreadStatusOpen,
		IncidentID:         req.IncidentID,
		CreatedBy:          userID,
		CreatedByRole:      primaryRole,
		CreatedByName:      userName,
		MessageCount:       1,
		UnreadCountSchool:  0,
		UnreadCountSupport: 0,
		LastMessageAt:      &now,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	// Set thread type based on incident link
	if req.IncidentID != nil && *req.IncidentID != "" {
		thread.ThreadType = models.ThreadTypeIncident
	}

	// Set unread count for the opposite side
	if h.isSchoolContact(roles) {
		thread.UnreadCountSupport = 1
	} else {
		thread.UnreadCountSchool = 1
	}

	if err := h.pg.Messaging().CreateThread(ctx, thread); err != nil {
		h.log.Error("failed to create thread", zap.Error(err))
		http.Error(w, "failed to create thread", http.StatusInternalServerError)
		return
	}

	// Create initial message
	message := models.Message{
		ID:          store.NewID("msg"),
		TenantID:    tenantID,
		ThreadID:    thread.ID,
		SenderID:    userID,
		SenderName:  userName,
		SenderRole:  primaryRole,
		Content:     strings.TrimSpace(req.InitialMessage),
		ContentType: models.ContentTypeText,
		CreatedAt:   now,
	}

	if err := h.pg.Messaging().CreateMessage(ctx, message); err != nil {
		h.log.Error("failed to create message", zap.Error(err))
		http.Error(w, "failed to create message", http.StatusInternalServerError)
		return
	}

	// Add creator as participant
	_ = h.pg.Messaging().AddThreadParticipant(ctx, models.ThreadParticipant{
		ThreadID: thread.ID,
		UserID:   userID,
		UserName: userName,
		UserRole: primaryRole,
		JoinedAt: now,
	})

	// Broadcast new thread via WebSocket
	h.broadcastNewMessage(tenantID, thread.ID, message)

	thread.LastMessage = &message
	writeJSON(w, http.StatusOK, models.CreateThreadResponse{
		Thread:  thread,
		Message: message,
	})
}

// UpdateThreadStatus handles PATCH /v1/messages/threads/{id}/status
func (h *MessagingHandler) UpdateThreadStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID := middleware.TenantID(ctx)
	roles := middleware.Roles(ctx)
	threadID := chi.URLParam(r, "id")

	var req models.UpdateThreadStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	// Validate status
	if req.Status != models.ThreadStatusOpen && req.Status != models.ThreadStatusClosed {
		http.Error(w, "invalid status", http.StatusBadRequest)
		return
	}

	// Get thread to check access
	thread, err := h.pg.Messaging().GetThreadByID(ctx, tenantID, threadID)
	if err != nil {
		http.Error(w, "thread not found", http.StatusNotFound)
		return
	}

	// Check access for school contacts
	if h.isSchoolContact(roles) {
		schools := middleware.AssignedSchools(ctx)
		if !contains(schools, thread.SchoolID) {
			http.Error(w, "access denied", http.StatusForbidden)
			return
		}
	}

	if err := h.pg.Messaging().UpdateThreadStatus(ctx, tenantID, threadID, req.Status); err != nil {
		h.log.Error("failed to update thread status", zap.Error(err))
		http.Error(w, "failed to update status", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

// CreateMessage handles POST /v1/messages/threads/{id}/messages
func (h *MessagingHandler) CreateMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID := middleware.TenantID(ctx)
	userID := middleware.UserID(ctx)
	userName := middleware.UserName(ctx)
	roles := middleware.Roles(ctx)
	threadID := chi.URLParam(r, "id")

	var req models.CreateMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Content) == "" {
		http.Error(w, "content is required", http.StatusBadRequest)
		return
	}

	// Get thread to check access
	thread, err := h.pg.Messaging().GetThreadByID(ctx, tenantID, threadID)
	if err != nil {
		http.Error(w, "thread not found", http.StatusNotFound)
		return
	}

	// Check access for school contacts
	if h.isSchoolContact(roles) {
		schools := middleware.AssignedSchools(ctx)
		if !contains(schools, thread.SchoolID) {
			http.Error(w, "access denied", http.StatusForbidden)
			return
		}
	}

	// Don't allow messages to closed threads
	if thread.Status == models.ThreadStatusClosed {
		http.Error(w, "thread is closed", http.StatusBadRequest)
		return
	}

	primaryRole := h.getPrimaryRole(roles)
	now := time.Now().UTC()

	message := models.Message{
		ID:          store.NewID("msg"),
		TenantID:    tenantID,
		ThreadID:    threadID,
		SenderID:    userID,
		SenderName:  userName,
		SenderRole:  primaryRole,
		Content:     strings.TrimSpace(req.Content),
		ContentType: models.ContentTypeText,
		CreatedAt:   now,
	}

	if err := h.pg.Messaging().CreateMessage(ctx, message); err != nil {
		h.log.Error("failed to create message", zap.Error(err))
		http.Error(w, "failed to create message", http.StatusInternalServerError)
		return
	}

	// Update thread stats
	isSchoolSender := h.isSchoolContact(roles)
	if err := h.pg.Messaging().UpdateThreadLastMessage(ctx, tenantID, threadID, isSchoolSender); err != nil {
		h.log.Warn("failed to update thread stats", zap.Error(err))
	}

	// Add sender as participant if not already
	_ = h.pg.Messaging().AddThreadParticipant(ctx, models.ThreadParticipant{
		ThreadID: threadID,
		UserID:   userID,
		UserName: userName,
		UserRole: primaryRole,
		JoinedAt: now,
	})

	// Update chat session message count if this is a livechat thread
	if thread.ThreadType == models.ThreadTypeLivechat {
		session, err := h.pg.ChatSessions().GetSessionByThreadID(ctx, tenantID, threadID)
		if err == nil && session.ID != "" {
			_ = h.pg.ChatSessions().IncrementMessageCount(ctx, tenantID, session.ID)

			// Set first response time if this is the first agent message
			if !isSchoolSender && session.FirstResponseSeconds == nil && session.AgentJoinedAt != nil {
				seconds := int(now.Sub(*session.AgentJoinedAt).Seconds())
				_ = h.pg.ChatSessions().SetFirstResponseTime(ctx, tenantID, session.ID, seconds)
			}
		}
	}

	// Broadcast message via WebSocket
	h.broadcastNewMessage(tenantID, threadID, message)

	writeJSON(w, http.StatusOK, models.CreateMessageResponse{
		Message: message,
	})
}

// MarkRead handles POST /v1/messages/threads/{id}/read
func (h *MessagingHandler) MarkRead(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID := middleware.TenantID(ctx)
	userID := middleware.UserID(ctx)
	roles := middleware.Roles(ctx)
	threadID := chi.URLParam(r, "id")

	var req models.MarkReadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	// Get thread to check access
	thread, err := h.pg.Messaging().GetThreadByID(ctx, tenantID, threadID)
	if err != nil {
		http.Error(w, "thread not found", http.StatusNotFound)
		return
	}

	// Check access for school contacts
	if h.isSchoolContact(roles) {
		schools := middleware.AssignedSchools(ctx)
		if !contains(schools, thread.SchoolID) {
			http.Error(w, "access denied", http.StatusForbidden)
			return
		}
	}

	isSchoolContact := h.isSchoolContact(roles)
	if err := h.pg.Messaging().MarkThreadRead(ctx, threadID, userID, req.LastMessageID, isSchoolContact); err != nil {
		h.log.Error("failed to mark thread read", zap.Error(err))
		http.Error(w, "failed to mark read", http.StatusInternalServerError)
		return
	}

	// Broadcast read receipt via WebSocket
	h.broadcastReadReceipt(tenantID, threadID, userID, req.LastMessageID)

	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

// GetUnreadCounts handles GET /v1/messages/unread
func (h *MessagingHandler) GetUnreadCounts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID := middleware.TenantID(ctx)
	userID := middleware.UserID(ctx)
	roles := middleware.Roles(ctx)

	isSchoolContact := h.isSchoolContact(roles)
	schoolID := ""
	if isSchoolContact {
		schools := middleware.AssignedSchools(ctx)
		if len(schools) > 0 {
			schoolID = schools[0]
		}
	}

	counts, err := h.pg.Messaging().GetUnreadCounts(ctx, tenantID, userID, isSchoolContact, schoolID)
	if err != nil {
		h.log.Error("failed to get unread counts", zap.Error(err))
		http.Error(w, "failed to get unread counts", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, counts)
}

// SearchMessages handles GET /v1/messages/search
func (h *MessagingHandler) SearchMessages(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID := middleware.TenantID(ctx)
	roles := middleware.Roles(ctx)

	query := strings.TrimSpace(r.URL.Query().Get("q"))
	if query == "" {
		http.Error(w, "q is required", http.StatusBadRequest)
		return
	}

	limit := parseLimit(r.URL.Query().Get("limit"), 20, 100)

	// Determine school filter
	schoolID := ""
	if h.isSchoolContact(roles) {
		schools := middleware.AssignedSchools(ctx)
		if len(schools) > 0 {
			schoolID = schools[0]
		}
	}

	results, err := h.pg.Messaging().SearchMessages(ctx, tenantID, schoolID, query, limit)
	if err != nil {
		h.log.Error("failed to search messages", zap.Error(err))
		http.Error(w, "failed to search", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, models.SearchResponse{
		Items: results,
		Total: len(results),
	})
}

// GetAnalytics handles GET /v1/messages/analytics
func (h *MessagingHandler) GetAnalytics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID := middleware.TenantID(ctx)
	roles := middleware.Roles(ctx)

	// Only admins can access analytics
	if !h.isAdmin(roles) {
		http.Error(w, "access denied", http.StatusForbidden)
		return
	}

	// Parse date range
	now := time.Now().UTC()
	from := now.AddDate(0, 0, -30) // Default: last 30 days
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

	analytics, err := h.pg.Messaging().GetMessagingAnalytics(ctx, tenantID, from, to)
	if err != nil {
		h.log.Error("failed to get analytics", zap.Error(err))
		http.Error(w, "failed to get analytics", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, analytics)
}

// Helper methods

func (h *MessagingHandler) isSchoolContact(roles []string) bool {
	for _, r := range roles {
		if r == "ssp_school_contact" {
			return true
		}
	}
	return false
}

func (h *MessagingHandler) isAdmin(roles []string) bool {
	for _, r := range roles {
		if r == "ssp_admin" {
			return true
		}
	}
	return false
}

func (h *MessagingHandler) getPrimaryRole(roles []string) string {
	// Priority order for determining primary role
	priority := []string{"ssp_admin", "ssp_support_agent", "ssp_school_contact", "ssp_lead_tech", "ssp_field_tech"}
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

func (h *MessagingHandler) broadcastNewMessage(tenantID, threadID string, message models.Message) {
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

func (h *MessagingHandler) broadcastReadReceipt(tenantID, threadID, userID, lastMessageID string) {
	if h.hub == nil {
		return
	}

	h.hub.Broadcast(tenantID, &ws.Message{
		Type: ws.MessageTypeChatRead,
		Payload: map[string]any{
			"threadId":      threadID,
			"userId":        userID,
			"lastMessageId": lastMessageID,
			"readAt":        time.Now().UTC().Format(time.RFC3339),
		},
	})
}

func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}
