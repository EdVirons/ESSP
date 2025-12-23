package store

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/edvirons/ssp/ims/internal/claude"
	"github.com/edvirons/ssp/ims/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ChatSessionsRepo handles database operations for chat sessions
type ChatSessionsRepo struct {
	pool *pgxpool.Pool
}

// ChatSessionListParams contains parameters for listing chat sessions
type ChatSessionListParams struct {
	TenantID        string
	AgentID         string
	SchoolID        string
	Status          string
	Limit           int
	HasCursor       bool
	CursorTimestamp time.Time
	CursorID        string
}

// CreateSession creates a new chat session
func (r *ChatSessionsRepo) CreateSession(ctx context.Context, s models.ChatSession) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO chat_sessions (
			id, tenant_id, school_id, thread_id, school_contact_id, school_contact_name,
			assigned_agent_id, assigned_agent_name, status, queue_position,
			started_at, agent_joined_at, ended_at, first_response_seconds,
			total_messages, rating, feedback, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9, $10,
			$11, $12, $13, $14,
			$15, $16, $17, $18, $19
		)
	`, s.ID, s.TenantID, s.SchoolID, s.ThreadID, s.SchoolContactID, s.SchoolContactName,
		s.AssignedAgentID, s.AssignedAgentName, s.Status, s.QueuePosition,
		s.StartedAt, s.AgentJoinedAt, s.EndedAt, s.FirstResponseSeconds,
		s.TotalMessages, s.Rating, s.Feedback, s.CreatedAt, s.UpdatedAt)
	return err
}

// GetSessionByID retrieves a session by ID
func (r *ChatSessionsRepo) GetSessionByID(ctx context.Context, tenantID, sessionID string) (models.ChatSession, error) {
	var s models.ChatSession
	row := r.pool.QueryRow(ctx, `
		SELECT id, tenant_id, school_id, thread_id, school_contact_id, school_contact_name,
			   assigned_agent_id, assigned_agent_name, status, queue_position,
			   started_at, agent_joined_at, ended_at, first_response_seconds,
			   total_messages, rating, feedback, created_at, updated_at
		FROM chat_sessions
		WHERE tenant_id = $1 AND id = $2
	`, tenantID, sessionID)

	err := row.Scan(&s.ID, &s.TenantID, &s.SchoolID, &s.ThreadID, &s.SchoolContactID, &s.SchoolContactName,
		&s.AssignedAgentID, &s.AssignedAgentName, &s.Status, &s.QueuePosition,
		&s.StartedAt, &s.AgentJoinedAt, &s.EndedAt, &s.FirstResponseSeconds,
		&s.TotalMessages, &s.Rating, &s.Feedback, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.ChatSession{}, errors.New("session not found")
		}
		return models.ChatSession{}, err
	}
	return s, nil
}

// GetSessionByThreadID retrieves a session by thread ID
func (r *ChatSessionsRepo) GetSessionByThreadID(ctx context.Context, tenantID, threadID string) (models.ChatSession, error) {
	var s models.ChatSession
	row := r.pool.QueryRow(ctx, `
		SELECT id, tenant_id, school_id, thread_id, school_contact_id, school_contact_name,
			   assigned_agent_id, assigned_agent_name, status, queue_position,
			   started_at, agent_joined_at, ended_at, first_response_seconds,
			   total_messages, rating, feedback, created_at, updated_at
		FROM chat_sessions
		WHERE tenant_id = $1 AND thread_id = $2
		ORDER BY created_at DESC
		LIMIT 1
	`, tenantID, threadID)

	err := row.Scan(&s.ID, &s.TenantID, &s.SchoolID, &s.ThreadID, &s.SchoolContactID, &s.SchoolContactName,
		&s.AssignedAgentID, &s.AssignedAgentName, &s.Status, &s.QueuePosition,
		&s.StartedAt, &s.AgentJoinedAt, &s.EndedAt, &s.FirstResponseSeconds,
		&s.TotalMessages, &s.Rating, &s.Feedback, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return models.ChatSession{}, err
	}
	return s, nil
}

// GetActiveSessionForUser gets the active session for a school contact
func (r *ChatSessionsRepo) GetActiveSessionForUser(ctx context.Context, tenantID, userID string) (models.ChatSession, error) {
	var s models.ChatSession
	row := r.pool.QueryRow(ctx, `
		SELECT id, tenant_id, school_id, thread_id, school_contact_id, school_contact_name,
			   assigned_agent_id, assigned_agent_name, status, queue_position,
			   started_at, agent_joined_at, ended_at, first_response_seconds,
			   total_messages, rating, feedback, created_at, updated_at,
			   COALESCE(ai_handled, false), ai_resolved, COALESCE(ai_turns, 0),
			   escalation_reason, COALESCE(escalation_summary, '{}'),
			   issue_category, issue_severity, COALESCE(collected_info, '{}')
		FROM chat_sessions
		WHERE tenant_id = $1 AND school_contact_id = $2 AND status IN ('ai_active', 'waiting', 'active')
		ORDER BY created_at DESC
		LIMIT 1
	`, tenantID, userID)

	var escalationSummaryJSON, collectedInfoJSON []byte
	err := row.Scan(&s.ID, &s.TenantID, &s.SchoolID, &s.ThreadID, &s.SchoolContactID, &s.SchoolContactName,
		&s.AssignedAgentID, &s.AssignedAgentName, &s.Status, &s.QueuePosition,
		&s.StartedAt, &s.AgentJoinedAt, &s.EndedAt, &s.FirstResponseSeconds,
		&s.TotalMessages, &s.Rating, &s.Feedback, &s.CreatedAt, &s.UpdatedAt,
		&s.AIHandled, &s.AIResolved, &s.AITurns,
		&s.EscalationReason, &escalationSummaryJSON,
		&s.IssueCategory, &s.IssueSeverity, &collectedInfoJSON)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.ChatSession{}, nil // No active session
		}
		return models.ChatSession{}, err
	}
	// Parse JSON fields
	if len(escalationSummaryJSON) > 0 {
		_ = json.Unmarshal(escalationSummaryJSON, &s.EscalationSummary)
	}
	if len(collectedInfoJSON) > 0 {
		_ = json.Unmarshal(collectedInfoJSON, &s.CollectedInfo)
	}
	return s, nil
}

// GetWaitingSessions gets all waiting sessions (queue)
func (r *ChatSessionsRepo) GetWaitingSessions(ctx context.Context, tenantID string) ([]models.ChatSession, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, tenant_id, school_id, thread_id, school_contact_id, school_contact_name,
			   assigned_agent_id, assigned_agent_name, status, queue_position,
			   started_at, agent_joined_at, ended_at, first_response_seconds,
			   total_messages, rating, feedback, created_at, updated_at
		FROM chat_sessions
		WHERE tenant_id = $1 AND status = 'waiting'
		ORDER BY started_at ASC
	`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []models.ChatSession
	for rows.Next() {
		var s models.ChatSession
		if err := rows.Scan(&s.ID, &s.TenantID, &s.SchoolID, &s.ThreadID, &s.SchoolContactID, &s.SchoolContactName,
			&s.AssignedAgentID, &s.AssignedAgentName, &s.Status, &s.QueuePosition,
			&s.StartedAt, &s.AgentJoinedAt, &s.EndedAt, &s.FirstResponseSeconds,
			&s.TotalMessages, &s.Rating, &s.Feedback, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		sessions = append(sessions, s)
	}

	return sessions, nil
}

// GetActiveSessionsForAgent gets all active sessions for an agent
func (r *ChatSessionsRepo) GetActiveSessionsForAgent(ctx context.Context, tenantID, agentID string) ([]models.ChatSession, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, tenant_id, school_id, thread_id, school_contact_id, school_contact_name,
			   assigned_agent_id, assigned_agent_name, status, queue_position,
			   started_at, agent_joined_at, ended_at, first_response_seconds,
			   total_messages, rating, feedback, created_at, updated_at
		FROM chat_sessions
		WHERE tenant_id = $1 AND assigned_agent_id = $2 AND status = 'active'
		ORDER BY started_at DESC
	`, tenantID, agentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []models.ChatSession
	for rows.Next() {
		var s models.ChatSession
		if err := rows.Scan(&s.ID, &s.TenantID, &s.SchoolID, &s.ThreadID, &s.SchoolContactID, &s.SchoolContactName,
			&s.AssignedAgentID, &s.AssignedAgentName, &s.Status, &s.QueuePosition,
			&s.StartedAt, &s.AgentJoinedAt, &s.EndedAt, &s.FirstResponseSeconds,
			&s.TotalMessages, &s.Rating, &s.Feedback, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		sessions = append(sessions, s)
	}

	return sessions, nil
}

// AssignAgent assigns an agent to a session
func (r *ChatSessionsRepo) AssignAgent(ctx context.Context, tenantID, sessionID, agentID, agentName string) error {
	now := time.Now().UTC()
	_, err := r.pool.Exec(ctx, `
		UPDATE chat_sessions
		SET assigned_agent_id = $3, assigned_agent_name = $4,
			status = 'active', queue_position = NULL,
			agent_joined_at = $5, updated_at = $5
		WHERE tenant_id = $1 AND id = $2 AND status = 'waiting'
	`, tenantID, sessionID, agentID, agentName, now)
	return err
}

// EndSession ends a chat session
func (r *ChatSessionsRepo) EndSession(ctx context.Context, tenantID, sessionID string, rating *int, feedback *string) error {
	now := time.Now().UTC()
	_, err := r.pool.Exec(ctx, `
		UPDATE chat_sessions
		SET status = 'ended', ended_at = $3, rating = $4, feedback = $5, updated_at = $3
		WHERE tenant_id = $1 AND id = $2
	`, tenantID, sessionID, now, rating, feedback)
	return err
}

// TransferSession transfers a session to another agent
func (r *ChatSessionsRepo) TransferSession(ctx context.Context, tenantID, sessionID, newAgentID, newAgentName string) error {
	now := time.Now().UTC()
	_, err := r.pool.Exec(ctx, `
		UPDATE chat_sessions
		SET assigned_agent_id = $3, assigned_agent_name = $4, updated_at = $5
		WHERE tenant_id = $1 AND id = $2 AND status = 'active'
	`, tenantID, sessionID, newAgentID, newAgentName, now)
	return err
}

// IncrementMessageCount increments the message count for a session
func (r *ChatSessionsRepo) IncrementMessageCount(ctx context.Context, tenantID, sessionID string) error {
	now := time.Now().UTC()
	_, err := r.pool.Exec(ctx, `
		UPDATE chat_sessions
		SET total_messages = total_messages + 1, updated_at = $3
		WHERE tenant_id = $1 AND id = $2
	`, tenantID, sessionID, now)
	return err
}

// SetFirstResponseTime sets the first response time for a session
func (r *ChatSessionsRepo) SetFirstResponseTime(ctx context.Context, tenantID, sessionID string, seconds int) error {
	now := time.Now().UTC()
	_, err := r.pool.Exec(ctx, `
		UPDATE chat_sessions
		SET first_response_seconds = $3, updated_at = $4
		WHERE tenant_id = $1 AND id = $2 AND first_response_seconds IS NULL
	`, tenantID, sessionID, seconds, now)
	return err
}

// GetQueuePosition gets the queue position for a session
func (r *ChatSessionsRepo) GetQueuePosition(ctx context.Context, tenantID, sessionID string) (int, error) {
	var position int
	row := r.pool.QueryRow(ctx, `
		SELECT COUNT(*) + 1
		FROM chat_sessions
		WHERE tenant_id = $1 AND status = 'waiting'
		AND started_at < (SELECT started_at FROM chat_sessions WHERE id = $2)
	`, tenantID, sessionID)
	err := row.Scan(&position)
	return position, err
}

// UpdateQueuePositions updates queue positions for all waiting sessions
func (r *ChatSessionsRepo) UpdateQueuePositions(ctx context.Context, tenantID string) error {
	_, err := r.pool.Exec(ctx, `
		WITH ranked AS (
			SELECT id, ROW_NUMBER() OVER (ORDER BY started_at ASC) as pos
			FROM chat_sessions
			WHERE tenant_id = $1 AND status = 'waiting'
		)
		UPDATE chat_sessions cs
		SET queue_position = r.pos
		FROM ranked r
		WHERE cs.id = r.id
	`, tenantID)
	return err
}

// GetNextInQueue gets the next session in the queue
func (r *ChatSessionsRepo) GetNextInQueue(ctx context.Context, tenantID string) (models.ChatSession, error) {
	var s models.ChatSession
	row := r.pool.QueryRow(ctx, `
		SELECT id, tenant_id, school_id, thread_id, school_contact_id, school_contact_name,
			   assigned_agent_id, assigned_agent_name, status, queue_position,
			   started_at, agent_joined_at, ended_at, first_response_seconds,
			   total_messages, rating, feedback, created_at, updated_at
		FROM chat_sessions
		WHERE tenant_id = $1 AND status = 'waiting'
		ORDER BY started_at ASC
		LIMIT 1
	`, tenantID)

	err := row.Scan(&s.ID, &s.TenantID, &s.SchoolID, &s.ThreadID, &s.SchoolContactID, &s.SchoolContactName,
		&s.AssignedAgentID, &s.AssignedAgentName, &s.Status, &s.QueuePosition,
		&s.StartedAt, &s.AgentJoinedAt, &s.EndedAt, &s.FirstResponseSeconds,
		&s.TotalMessages, &s.Rating, &s.Feedback, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.ChatSession{}, nil
		}
		return models.ChatSession{}, err
	}
	return s, nil
}

// Agent Availability

// GetAgentAvailability gets an agent's availability
func (r *ChatSessionsRepo) GetAgentAvailability(ctx context.Context, tenantID, userID string) (models.AgentAvailability, error) {
	var a models.AgentAvailability
	row := r.pool.QueryRow(ctx, `
		SELECT tenant_id, user_id, is_available, max_concurrent_chats,
			   current_chat_count, last_seen_at, updated_at
		FROM agent_availability
		WHERE tenant_id = $1 AND user_id = $2
	`, tenantID, userID)

	err := row.Scan(&a.TenantID, &a.UserID, &a.IsAvailable, &a.MaxConcurrentChats,
		&a.CurrentChatCount, &a.LastSeenAt, &a.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Return default availability
			return models.AgentAvailability{
				TenantID:           tenantID,
				UserID:             userID,
				IsAvailable:        false,
				MaxConcurrentChats: 3,
				CurrentChatCount:   0,
				LastSeenAt:         time.Now().UTC(),
				UpdatedAt:          time.Now().UTC(),
			}, nil
		}
		return models.AgentAvailability{}, err
	}
	return a, nil
}

// SetAgentAvailability sets an agent's availability
func (r *ChatSessionsRepo) SetAgentAvailability(ctx context.Context, tenantID, userID string, available bool, maxChats int) error {
	now := time.Now().UTC()
	_, err := r.pool.Exec(ctx, `
		INSERT INTO agent_availability (tenant_id, user_id, is_available, max_concurrent_chats, current_chat_count, last_seen_at, updated_at)
		VALUES ($1, $2, $3, $4, 0, $5, $5)
		ON CONFLICT (tenant_id, user_id)
		DO UPDATE SET is_available = $3, max_concurrent_chats = $4, last_seen_at = $5, updated_at = $5
	`, tenantID, userID, available, maxChats, now)
	return err
}

// IncrementAgentChatCount increments the current chat count for an agent
func (r *ChatSessionsRepo) IncrementAgentChatCount(ctx context.Context, tenantID, userID string) error {
	now := time.Now().UTC()
	_, err := r.pool.Exec(ctx, `
		UPDATE agent_availability
		SET current_chat_count = current_chat_count + 1, last_seen_at = $3, updated_at = $3
		WHERE tenant_id = $1 AND user_id = $2
	`, tenantID, userID, now)
	return err
}

// DecrementAgentChatCount decrements the current chat count for an agent
func (r *ChatSessionsRepo) DecrementAgentChatCount(ctx context.Context, tenantID, userID string) error {
	now := time.Now().UTC()
	_, err := r.pool.Exec(ctx, `
		UPDATE agent_availability
		SET current_chat_count = GREATEST(current_chat_count - 1, 0), last_seen_at = $3, updated_at = $3
		WHERE tenant_id = $1 AND user_id = $2
	`, tenantID, userID, now)
	return err
}

// GetAvailableAgents gets all available agents
func (r *ChatSessionsRepo) GetAvailableAgents(ctx context.Context, tenantID string) ([]models.AgentAvailability, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT tenant_id, user_id, is_available, max_concurrent_chats,
			   current_chat_count, last_seen_at, updated_at
		FROM agent_availability
		WHERE tenant_id = $1 AND is_available = true
		AND current_chat_count < max_concurrent_chats
		ORDER BY current_chat_count ASC, last_seen_at ASC
	`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var agents []models.AgentAvailability
	for rows.Next() {
		var a models.AgentAvailability
		if err := rows.Scan(&a.TenantID, &a.UserID, &a.IsAvailable, &a.MaxConcurrentChats,
			&a.CurrentChatCount, &a.LastSeenAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		agents = append(agents, a)
	}

	return agents, nil
}

// GetChatMetrics gets chat metrics for analytics
func (r *ChatSessionsRepo) GetChatMetrics(ctx context.Context, tenantID string, from, to time.Time) (models.ChatMetrics, error) {
	var m models.ChatMetrics

	row := r.pool.QueryRow(ctx, `
		SELECT
			COUNT(*),
			COALESCE(AVG(EXTRACT(EPOCH FROM (agent_joined_at - started_at))), 0),
			COALESCE(AVG(first_response_seconds), 0),
			COALESCE(AVG(rating), 0),
			COUNT(*) FILTER (WHERE rating IS NOT NULL),
			COUNT(*) FILTER (WHERE status = 'active'),
			COUNT(*) FILTER (WHERE status = 'waiting'),
			COUNT(*) FILTER (WHERE status = 'ended')
		FROM chat_sessions
		WHERE tenant_id = $1 AND created_at BETWEEN $2 AND $3
	`, tenantID, from, to)

	err := row.Scan(&m.TotalSessions, &m.AverageWaitTime, &m.AverageResponseTime,
		&m.AverageRating, &m.SessionsWithRating, &m.ActiveSessions,
		&m.WaitingSessions, &m.EndedSessions)

	return m, err
}

// AI Support Methods

// EscalateToHuman escalates a session from AI to human queue
func (r *ChatSessionsRepo) EscalateToHuman(ctx context.Context, tenantID, sessionID, reason string, summary map[string]any) error {
	now := time.Now().UTC()
	summaryJSON, _ := json.Marshal(summary)
	_, err := r.pool.Exec(ctx, `
		UPDATE chat_sessions
		SET status = 'waiting', escalation_reason = $3, escalation_summary = $4,
			ai_resolved = false, updated_at = $5
		WHERE tenant_id = $1 AND id = $2 AND status = 'ai_active'
	`, tenantID, sessionID, reason, summaryJSON, now)
	return err
}

// IncrementAITurns increments the AI turn counter for a session
func (r *ChatSessionsRepo) IncrementAITurns(ctx context.Context, tenantID, sessionID string) error {
	now := time.Now().UTC()
	_, err := r.pool.Exec(ctx, `
		UPDATE chat_sessions
		SET ai_turns = ai_turns + 1, updated_at = $3
		WHERE tenant_id = $1 AND id = $2
	`, tenantID, sessionID, now)
	return err
}

// UpdateAISessionData updates AI-collected information for a session
func (r *ChatSessionsRepo) UpdateAISessionData(ctx context.Context, tenantID, sessionID string, turns int, category, severity string, collectedInfo map[string]any) error {
	now := time.Now().UTC()
	infoJSON, _ := json.Marshal(collectedInfo)

	var catPtr, sevPtr *string
	if category != "" {
		catPtr = &category
	}
	if severity != "" {
		sevPtr = &severity
	}

	_, err := r.pool.Exec(ctx, `
		UPDATE chat_sessions
		SET ai_turns = $3, issue_category = COALESCE($4, issue_category),
			issue_severity = COALESCE($5, issue_severity),
			collected_info = $6, updated_at = $7
		WHERE tenant_id = $1 AND id = $2
	`, tenantID, sessionID, turns, catPtr, sevPtr, infoJSON, now)
	return err
}

// UpdateAIResolution marks a session as resolved by AI
func (r *ChatSessionsRepo) UpdateAIResolution(ctx context.Context, tenantID, sessionID string, resolved bool) error {
	now := time.Now().UTC()
	_, err := r.pool.Exec(ctx, `
		UPDATE chat_sessions
		SET ai_resolved = $3, updated_at = $4
		WHERE tenant_id = $1 AND id = $2
	`, tenantID, sessionID, resolved, now)
	return err
}

// LogAIConversationTurn logs a conversation turn for analytics
func (r *ChatSessionsRepo) LogAIConversationTurn(ctx context.Context, turn claude.ConversationTurn) error {
	signalsJSON, _ := json.Marshal(turn.EscalationSignals)
	contextJSON, _ := json.Marshal(turn.ContextUsed)

	_, err := r.pool.Exec(ctx, `
		INSERT INTO ai_conversation_logs (
			id, tenant_id, session_id, turn_number, user_message, ai_response,
			input_tokens, output_tokens, response_time_ms,
			escalation_recommended, escalation_signals, context_used, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`, turn.ID, turn.TenantID, turn.SessionID, turn.TurnNumber, turn.UserMessage, turn.AIResponse,
		turn.InputTokens, turn.OutputTokens, turn.ResponseTimeMs,
		turn.EscalationRecommended, signalsJSON, contextJSON, turn.CreatedAt)
	return err
}

// GetAIConversationLogs gets conversation logs for a session
func (r *ChatSessionsRepo) GetAIConversationLogs(ctx context.Context, tenantID, sessionID string) ([]claude.ConversationTurn, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, tenant_id, session_id, turn_number, user_message, ai_response,
			   input_tokens, output_tokens, response_time_ms,
			   escalation_recommended, escalation_signals, context_used, created_at
		FROM ai_conversation_logs
		WHERE tenant_id = $1 AND session_id = $2
		ORDER BY turn_number ASC
	`, tenantID, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var turns []claude.ConversationTurn
	for rows.Next() {
		var t claude.ConversationTurn
		var signalsJSON, contextJSON []byte
		if err := rows.Scan(&t.ID, &t.TenantID, &t.SessionID, &t.TurnNumber, &t.UserMessage, &t.AIResponse,
			&t.InputTokens, &t.OutputTokens, &t.ResponseTimeMs,
			&t.EscalationRecommended, &signalsJSON, &contextJSON, &t.CreatedAt); err != nil {
			return nil, err
		}
		_ = json.Unmarshal(signalsJSON, &t.EscalationSignals)
		_ = json.Unmarshal(contextJSON, &t.ContextUsed)
		turns = append(turns, t)
	}

	return turns, nil
}
