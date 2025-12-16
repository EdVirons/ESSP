package store

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/edvirons/ssp/ims/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// MessagingRepo handles database operations for messaging
type MessagingRepo struct {
	pool *pgxpool.Pool
}

// ThreadListParams contains parameters for listing threads
type ThreadListParams struct {
	TenantID   string
	SchoolID   string // Optional: if empty, list all threads (for support/admin)
	UserID     string
	UserRoles  []string
	Status     string
	IncidentID string
	Query      string
	Limit      int

	HasCursor       bool
	CursorTimestamp time.Time
	CursorID        string
}

// MessageListParams contains parameters for listing messages
type MessageListParams struct {
	ThreadID string
	Limit    int

	HasCursor       bool
	CursorTimestamp time.Time
	CursorID        string
}

// CreateThread creates a new message thread
func (r *MessagingRepo) CreateThread(ctx context.Context, t models.MessageThread) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO message_threads (
			id, tenant_id, school_id, subject, thread_type, status,
			incident_id, created_by, created_by_role, created_by_name,
			message_count, unread_count_school, unread_count_support,
			last_message_at, created_at, updated_at, closed_at
		) VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9, $10,
			$11, $12, $13,
			$14, $15, $16, $17
		)
	`, t.ID, t.TenantID, t.SchoolID, t.Subject, t.ThreadType, t.Status,
		t.IncidentID, t.CreatedBy, t.CreatedByRole, t.CreatedByName,
		t.MessageCount, t.UnreadCountSchool, t.UnreadCountSupport,
		t.LastMessageAt, t.CreatedAt, t.UpdatedAt, t.ClosedAt)
	return err
}

// GetThreadByID retrieves a thread by ID
func (r *MessagingRepo) GetThreadByID(ctx context.Context, tenantID, threadID string) (models.MessageThread, error) {
	var t models.MessageThread
	row := r.pool.QueryRow(ctx, `
		SELECT id, tenant_id, school_id, subject, thread_type, status,
			   incident_id, created_by, created_by_role, created_by_name,
			   message_count, unread_count_school, unread_count_support,
			   last_message_at, created_at, updated_at, closed_at
		FROM message_threads
		WHERE tenant_id = $1 AND id = $2
	`, tenantID, threadID)

	err := row.Scan(&t.ID, &t.TenantID, &t.SchoolID, &t.Subject, &t.ThreadType, &t.Status,
		&t.IncidentID, &t.CreatedBy, &t.CreatedByRole, &t.CreatedByName,
		&t.MessageCount, &t.UnreadCountSchool, &t.UnreadCountSupport,
		&t.LastMessageAt, &t.CreatedAt, &t.UpdatedAt, &t.ClosedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.MessageThread{}, errors.New("thread not found")
		}
		return models.MessageThread{}, err
	}
	return t, nil
}

// ListThreads lists threads based on parameters
func (r *MessagingRepo) ListThreads(ctx context.Context, p ThreadListParams) ([]models.MessageThread, string, error) {
	conds := []string{"tenant_id = $1"}
	args := []any{p.TenantID}
	argN := 2

	// School ID filter (for school contacts)
	if p.SchoolID != "" {
		conds = append(conds, "school_id = $"+itoa(argN))
		args = append(args, p.SchoolID)
		argN++
	}

	// Status filter
	if p.Status != "" && p.Status != "all" {
		conds = append(conds, "status = $"+itoa(argN))
		args = append(args, p.Status)
		argN++
	}

	// Incident ID filter
	if p.IncidentID != "" {
		conds = append(conds, "incident_id = $"+itoa(argN))
		args = append(args, p.IncidentID)
		argN++
	}

	// Search query
	if p.Query != "" {
		conds = append(conds, "subject ILIKE $"+itoa(argN))
		args = append(args, "%"+p.Query+"%")
		argN++
	}

	// Cursor pagination
	if p.HasCursor {
		conds = append(conds, "(COALESCE(last_message_at, created_at), id) < ($"+itoa(argN)+", $"+itoa(argN+1)+")")
		args = append(args, p.CursorTimestamp, p.CursorID)
		argN += 2
	}

	// Limit
	limitPlus := p.Limit + 1
	args = append(args, limitPlus)

	sql := `
		SELECT id, tenant_id, school_id, subject, thread_type, status,
			   incident_id, created_by, created_by_role, created_by_name,
			   message_count, unread_count_school, unread_count_support,
			   last_message_at, created_at, updated_at, closed_at
		FROM message_threads
		WHERE ` + strings.Join(conds, " AND ") + `
		ORDER BY COALESCE(last_message_at, created_at) DESC, id DESC
		LIMIT $` + itoa(argN)

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	var threads []models.MessageThread
	for rows.Next() {
		var t models.MessageThread
		if err := rows.Scan(&t.ID, &t.TenantID, &t.SchoolID, &t.Subject, &t.ThreadType, &t.Status,
			&t.IncidentID, &t.CreatedBy, &t.CreatedByRole, &t.CreatedByName,
			&t.MessageCount, &t.UnreadCountSchool, &t.UnreadCountSupport,
			&t.LastMessageAt, &t.CreatedAt, &t.UpdatedAt, &t.ClosedAt); err != nil {
			return nil, "", err
		}
		threads = append(threads, t)
	}

	next := ""
	if len(threads) > p.Limit {
		last := threads[p.Limit-1]
		ts := last.CreatedAt
		if last.LastMessageAt != nil {
			ts = *last.LastMessageAt
		}
		next = EncodeCursor(ts, last.ID)
		threads = threads[:p.Limit]
	}

	return threads, next, nil
}

// UpdateThreadStatus updates the status of a thread
func (r *MessagingRepo) UpdateThreadStatus(ctx context.Context, tenantID, threadID string, status models.ThreadStatus) error {
	now := time.Now().UTC()
	var closedAt *time.Time
	if status == models.ThreadStatusClosed {
		closedAt = &now
	}

	_, err := r.pool.Exec(ctx, `
		UPDATE message_threads
		SET status = $3, updated_at = $4, closed_at = $5
		WHERE tenant_id = $1 AND id = $2
	`, tenantID, threadID, status, now, closedAt)
	return err
}

// UpdateThreadLastMessage updates the thread's last message timestamp and counts
func (r *MessagingRepo) UpdateThreadLastMessage(ctx context.Context, tenantID, threadID string, isSchoolSender bool) error {
	now := time.Now().UTC()

	// Increment the appropriate unread counter
	var sql string
	if isSchoolSender {
		sql = `
			UPDATE message_threads
			SET message_count = message_count + 1,
				unread_count_support = unread_count_support + 1,
				last_message_at = $3,
				updated_at = $3
			WHERE tenant_id = $1 AND id = $2
		`
	} else {
		sql = `
			UPDATE message_threads
			SET message_count = message_count + 1,
				unread_count_school = unread_count_school + 1,
				last_message_at = $3,
				updated_at = $3
			WHERE tenant_id = $1 AND id = $2
		`
	}

	_, err := r.pool.Exec(ctx, sql, tenantID, threadID, now)
	return err
}

// CreateMessage creates a new message
func (r *MessagingRepo) CreateMessage(ctx context.Context, m models.Message) error {
	metadataBytes, _ := json.Marshal(m.Metadata)
	if m.Metadata == nil {
		metadataBytes = []byte("{}")
	}

	_, err := r.pool.Exec(ctx, `
		INSERT INTO messages (
			id, tenant_id, thread_id, sender_id, sender_name, sender_role,
			content, content_type, metadata, edited_at, deleted_at, created_at
		) VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9, $10, $11, $12
		)
	`, m.ID, m.TenantID, m.ThreadID, m.SenderID, m.SenderName, m.SenderRole,
		m.Content, m.ContentType, metadataBytes, m.EditedAt, m.DeletedAt, m.CreatedAt)
	return err
}

// GetMessageByID retrieves a message by ID
func (r *MessagingRepo) GetMessageByID(ctx context.Context, tenantID, messageID string) (models.Message, error) {
	var m models.Message
	var metadataBytes []byte

	row := r.pool.QueryRow(ctx, `
		SELECT id, tenant_id, thread_id, sender_id, sender_name, sender_role,
			   content, content_type, metadata, edited_at, deleted_at, created_at
		FROM messages
		WHERE tenant_id = $1 AND id = $2
	`, tenantID, messageID)

	err := row.Scan(&m.ID, &m.TenantID, &m.ThreadID, &m.SenderID, &m.SenderName, &m.SenderRole,
		&m.Content, &m.ContentType, &metadataBytes, &m.EditedAt, &m.DeletedAt, &m.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Message{}, errors.New("message not found")
		}
		return models.Message{}, err
	}

	if len(metadataBytes) > 0 {
		json.Unmarshal(metadataBytes, &m.Metadata)
	}

	return m, nil
}

// ListMessages lists messages in a thread
func (r *MessagingRepo) ListMessages(ctx context.Context, p MessageListParams) ([]models.Message, string, error) {
	conds := []string{"thread_id = $1", "deleted_at IS NULL"}
	args := []any{p.ThreadID}
	argN := 2

	// Cursor pagination (for messages, we paginate in ascending order)
	if p.HasCursor {
		conds = append(conds, "(created_at, id) > ($"+itoa(argN)+", $"+itoa(argN+1)+")")
		args = append(args, p.CursorTimestamp, p.CursorID)
		argN += 2
	}

	limitPlus := p.Limit + 1
	args = append(args, limitPlus)

	sql := `
		SELECT id, tenant_id, thread_id, sender_id, sender_name, sender_role,
			   content, content_type, metadata, edited_at, deleted_at, created_at
		FROM messages
		WHERE ` + strings.Join(conds, " AND ") + `
		ORDER BY created_at ASC, id ASC
		LIMIT $` + itoa(argN)

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var m models.Message
		var metadataBytes []byte
		if err := rows.Scan(&m.ID, &m.TenantID, &m.ThreadID, &m.SenderID, &m.SenderName, &m.SenderRole,
			&m.Content, &m.ContentType, &metadataBytes, &m.EditedAt, &m.DeletedAt, &m.CreatedAt); err != nil {
			return nil, "", err
		}
		if len(metadataBytes) > 0 {
			json.Unmarshal(metadataBytes, &m.Metadata)
		}
		messages = append(messages, m)
	}

	next := ""
	if len(messages) > p.Limit {
		last := messages[p.Limit-1]
		next = EncodeCursor(last.CreatedAt, last.ID)
		messages = messages[:p.Limit]
	}

	return messages, next, nil
}

// GetLastMessage gets the last message in a thread
func (r *MessagingRepo) GetLastMessage(ctx context.Context, threadID string) (models.Message, error) {
	var m models.Message
	var metadataBytes []byte

	row := r.pool.QueryRow(ctx, `
		SELECT id, tenant_id, thread_id, sender_id, sender_name, sender_role,
			   content, content_type, metadata, edited_at, deleted_at, created_at
		FROM messages
		WHERE thread_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT 1
	`, threadID)

	err := row.Scan(&m.ID, &m.TenantID, &m.ThreadID, &m.SenderID, &m.SenderName, &m.SenderRole,
		&m.Content, &m.ContentType, &metadataBytes, &m.EditedAt, &m.DeletedAt, &m.CreatedAt)
	if err != nil {
		return models.Message{}, err
	}

	if len(metadataBytes) > 0 {
		json.Unmarshal(metadataBytes, &m.Metadata)
	}

	return m, nil
}

// UpdateMessage updates a message's content
func (r *MessagingRepo) UpdateMessage(ctx context.Context, tenantID, messageID, content string) error {
	now := time.Now().UTC()
	_, err := r.pool.Exec(ctx, `
		UPDATE messages
		SET content = $3, edited_at = $4
		WHERE tenant_id = $1 AND id = $2 AND deleted_at IS NULL
	`, tenantID, messageID, content, now)
	return err
}

// DeleteMessage soft-deletes a message
func (r *MessagingRepo) DeleteMessage(ctx context.Context, tenantID, messageID string) error {
	now := time.Now().UTC()
	_, err := r.pool.Exec(ctx, `
		UPDATE messages
		SET deleted_at = $3, content = '[Message deleted]'
		WHERE tenant_id = $1 AND id = $2
	`, tenantID, messageID, now)
	return err
}

// MarkThreadRead marks a thread as read for a user
func (r *MessagingRepo) MarkThreadRead(ctx context.Context, threadID, userID, messageID string, isSchoolContact bool) error {
	now := time.Now().UTC()

	// Upsert read receipt
	_, err := r.pool.Exec(ctx, `
		INSERT INTO message_read_receipts (thread_id, user_id, last_read_message_id, last_read_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (thread_id, user_id)
		DO UPDATE SET last_read_message_id = $3, last_read_at = $4
	`, threadID, userID, messageID, now)
	if err != nil {
		return err
	}

	// Reset unread counter for this user's side
	if isSchoolContact {
		_, err = r.pool.Exec(ctx, `
			UPDATE message_threads
			SET unread_count_school = 0, updated_at = $2
			WHERE id = $1
		`, threadID, now)
	} else {
		_, err = r.pool.Exec(ctx, `
			UPDATE message_threads
			SET unread_count_support = 0, updated_at = $2
			WHERE id = $1
		`, threadID, now)
	}

	return err
}

// GetUnreadCounts gets unread counts for a user
func (r *MessagingRepo) GetUnreadCounts(ctx context.Context, tenantID, userID string, isSchoolContact bool, schoolID string) (models.UnreadCounts, error) {
	var counts models.UnreadCounts

	var sql string
	var args []any

	if isSchoolContact {
		// School contacts see their school's unread count
		sql = `
			SELECT
				COUNT(*) FILTER (WHERE unread_count_school > 0),
				COALESCE(SUM(unread_count_school), 0)
			FROM message_threads
			WHERE tenant_id = $1 AND school_id = $2 AND status != 'archived'
		`
		args = []any{tenantID, schoolID}
	} else {
		// Support staff see all threads' support unread count
		sql = `
			SELECT
				COUNT(*) FILTER (WHERE unread_count_support > 0),
				COALESCE(SUM(unread_count_support), 0)
			FROM message_threads
			WHERE tenant_id = $1 AND status != 'archived'
		`
		args = []any{tenantID}
	}

	row := r.pool.QueryRow(ctx, sql, args...)
	err := row.Scan(&counts.Threads, &counts.Messages)
	return counts, err
}

// CreateAttachment creates a message attachment
func (r *MessagingRepo) CreateAttachment(ctx context.Context, a models.MessageAttachment) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO message_attachments (
			id, tenant_id, message_id, file_name, content_type,
			size_bytes, object_key, thumbnail_key, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, a.ID, a.TenantID, a.MessageID, a.FileName, a.ContentType,
		a.SizeBytes, a.ObjectKey, a.ThumbnailKey, a.CreatedAt)
	return err
}

// GetAttachmentsByMessageID gets all attachments for a message
func (r *MessagingRepo) GetAttachmentsByMessageID(ctx context.Context, messageID string) ([]models.MessageAttachment, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, tenant_id, message_id, file_name, content_type,
			   size_bytes, object_key, thumbnail_key, created_at
		FROM message_attachments
		WHERE message_id = $1
	`, messageID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attachments []models.MessageAttachment
	for rows.Next() {
		var a models.MessageAttachment
		if err := rows.Scan(&a.ID, &a.TenantID, &a.MessageID, &a.FileName, &a.ContentType,
			&a.SizeBytes, &a.ObjectKey, &a.ThumbnailKey, &a.CreatedAt); err != nil {
			return nil, err
		}
		attachments = append(attachments, a)
	}

	return attachments, nil
}

// AddThreadParticipant adds a participant to a thread
func (r *MessagingRepo) AddThreadParticipant(ctx context.Context, p models.ThreadParticipant) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO thread_participants (thread_id, user_id, user_name, user_role, joined_at, left_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (thread_id, user_id) DO UPDATE
		SET left_at = NULL, joined_at = $5
	`, p.ThreadID, p.UserID, p.UserName, p.UserRole, p.JoinedAt, p.LeftAt)
	return err
}

// GetThreadParticipants gets all participants of a thread
func (r *MessagingRepo) GetThreadParticipants(ctx context.Context, threadID string) ([]models.ThreadParticipant, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT thread_id, user_id, user_name, user_role, joined_at, left_at
		FROM thread_participants
		WHERE thread_id = $1 AND left_at IS NULL
		ORDER BY joined_at ASC
	`, threadID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var participants []models.ThreadParticipant
	for rows.Next() {
		var p models.ThreadParticipant
		if err := rows.Scan(&p.ThreadID, &p.UserID, &p.UserName, &p.UserRole, &p.JoinedAt, &p.LeftAt); err != nil {
			return nil, err
		}
		participants = append(participants, p)
	}

	return participants, nil
}

// SearchMessages searches for messages
func (r *MessagingRepo) SearchMessages(ctx context.Context, tenantID, schoolID, query string, limit int) ([]models.MessageSearchResult, error) {
	conds := []string{"m.tenant_id = $1"}
	args := []any{tenantID}
	argN := 2

	if schoolID != "" {
		conds = append(conds, "t.school_id = $"+itoa(argN))
		args = append(args, schoolID)
		argN++
	}

	conds = append(conds, "m.content ILIKE $"+itoa(argN))
	args = append(args, "%"+query+"%")
	argN++

	args = append(args, limit)

	sql := `
		SELECT m.id, m.tenant_id, m.thread_id, m.sender_id, m.sender_name, m.sender_role,
			   m.content, m.content_type, m.metadata, m.edited_at, m.deleted_at, m.created_at,
			   t.id, t.tenant_id, t.school_id, t.subject, t.thread_type, t.status,
			   t.incident_id, t.created_by, t.created_by_role, t.created_by_name,
			   t.message_count, t.unread_count_school, t.unread_count_support,
			   t.last_message_at, t.created_at, t.updated_at, t.closed_at
		FROM messages m
		JOIN message_threads t ON m.thread_id = t.id
		WHERE ` + strings.Join(conds, " AND ") + ` AND m.deleted_at IS NULL
		ORDER BY m.created_at DESC
		LIMIT $` + itoa(argN)

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.MessageSearchResult
	for rows.Next() {
		var m models.Message
		var t models.MessageThread
		var metadataBytes []byte

		if err := rows.Scan(
			&m.ID, &m.TenantID, &m.ThreadID, &m.SenderID, &m.SenderName, &m.SenderRole,
			&m.Content, &m.ContentType, &metadataBytes, &m.EditedAt, &m.DeletedAt, &m.CreatedAt,
			&t.ID, &t.TenantID, &t.SchoolID, &t.Subject, &t.ThreadType, &t.Status,
			&t.IncidentID, &t.CreatedBy, &t.CreatedByRole, &t.CreatedByName,
			&t.MessageCount, &t.UnreadCountSchool, &t.UnreadCountSupport,
			&t.LastMessageAt, &t.CreatedAt, &t.UpdatedAt, &t.ClosedAt,
		); err != nil {
			return nil, err
		}

		if len(metadataBytes) > 0 {
			json.Unmarshal(metadataBytes, &m.Metadata)
		}

		results = append(results, models.MessageSearchResult{
			Message: m,
			Thread:  t,
		})
	}

	return results, nil
}

// GetThreadMessages retrieves messages from a thread by offset/limit
func (r *MessagingRepo) GetThreadMessages(ctx context.Context, tenantID, threadID string, offset, limit int) ([]models.Message, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, tenant_id, thread_id, sender_id, sender_name, sender_role,
			   content, content_type, metadata, edited_at, deleted_at, created_at
		FROM messages
		WHERE tenant_id = $1 AND thread_id = $2 AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`, tenantID, threadID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var m models.Message
		var metadataBytes []byte
		if err := rows.Scan(&m.ID, &m.TenantID, &m.ThreadID, &m.SenderID, &m.SenderName, &m.SenderRole,
			&m.Content, &m.ContentType, &metadataBytes, &m.EditedAt, &m.DeletedAt, &m.CreatedAt); err != nil {
			return nil, err
		}
		if len(metadataBytes) > 0 {
			json.Unmarshal(metadataBytes, &m.Metadata)
		}
		messages = append(messages, m)
	}

	return messages, nil
}

// GetMessagingAnalytics gets messaging analytics for a tenant
func (r *MessagingRepo) GetMessagingAnalytics(ctx context.Context, tenantID string, from, to time.Time) (models.MessagingAnalytics, error) {
	var analytics models.MessagingAnalytics

	// Get thread stats
	row := r.pool.QueryRow(ctx, `
		SELECT
			COUNT(*) FILTER (WHERE created_at BETWEEN $2 AND $3),
			COUNT(*) FILTER (WHERE status = 'open'),
			COUNT(*) FILTER (WHERE status = 'closed'),
			COALESCE(AVG(message_count), 0)
		FROM message_threads
		WHERE tenant_id = $1
	`, tenantID, from, to)

	err := row.Scan(&analytics.ThreadsCreated, &analytics.ActiveThreads, &analytics.ClosedThreads, &analytics.AvgMessagesThread)
	if err != nil {
		return analytics, err
	}

	// Get message count
	row = r.pool.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM messages
		WHERE tenant_id = $1 AND created_at BETWEEN $2 AND $3
	`, tenantID, from, to)
	row.Scan(&analytics.MessagesSent)

	// Get chat session stats
	row = r.pool.QueryRow(ctx, `
		SELECT
			COUNT(*),
			COALESCE(AVG(rating), 0)
		FROM chat_sessions
		WHERE tenant_id = $1 AND created_at BETWEEN $2 AND $3
	`, tenantID, from, to)
	row.Scan(&analytics.ChatSessions, &analytics.AvgChatRating)

	return analytics, nil
}
