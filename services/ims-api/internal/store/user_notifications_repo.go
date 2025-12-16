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

type UserNotificationsRepo struct{ pool *pgxpool.Pool }

// CreateNotification creates a new notification for a user.
func (r *UserNotificationsRepo) CreateNotification(ctx context.Context, n models.UserNotification) error {
	metadataJSON, err := json.Marshal(n.Metadata)
	if err != nil {
		return err
	}
	_, err = r.pool.Exec(ctx, `
		INSERT INTO user_notifications (
			id, tenant_id, user_id, notification_type, entity_type, entity_id,
			project_id, title, body, metadata, is_read, created_at
		) VALUES ($1,$2,$3,$4,$5,$6,NULLIF($7,''),$8,$9,$10,$11,$12)
	`, n.ID, n.TenantID, n.UserID, n.NotificationType, n.EntityType, n.EntityID,
		n.ProjectID, n.Title, n.Body, metadataJSON, n.IsRead, n.CreatedAt)
	return err
}

// CreateBulkNotifications creates notifications for multiple users.
func (r *UserNotificationsRepo) CreateBulkNotifications(ctx context.Context, notifications []models.UserNotification) error {
	if len(notifications) == 0 {
		return nil
	}

	batch := &pgx.Batch{}
	for _, n := range notifications {
		metadataJSON, err := json.Marshal(n.Metadata)
		if err != nil {
			return err
		}
		batch.Queue(`
			INSERT INTO user_notifications (
				id, tenant_id, user_id, notification_type, entity_type, entity_id,
				project_id, title, body, metadata, is_read, created_at
			) VALUES ($1,$2,$3,$4,$5,$6,NULLIF($7,''),$8,$9,$10,$11,$12)
		`, n.ID, n.TenantID, n.UserID, n.NotificationType, n.EntityType, n.EntityID,
			n.ProjectID, n.Title, n.Body, metadataJSON, n.IsRead, n.CreatedAt)
	}

	br := r.pool.SendBatch(ctx, batch)
	defer br.Close()

	for range notifications {
		if _, err := br.Exec(); err != nil {
			return err
		}
	}
	return nil
}

// GetNotification retrieves a specific notification.
func (r *UserNotificationsRepo) GetNotification(ctx context.Context, tenantID, notificationID string) (models.UserNotification, error) {
	var n models.UserNotification
	var metadataJSON []byte
	var projectID *string
	row := r.pool.QueryRow(ctx, `
		SELECT id, tenant_id, user_id, notification_type, entity_type, entity_id,
			project_id, title, body, metadata, is_read, read_at, created_at
		FROM user_notifications
		WHERE tenant_id = $1 AND id = $2
	`, tenantID, notificationID)
	if err := row.Scan(&n.ID, &n.TenantID, &n.UserID, &n.NotificationType, &n.EntityType, &n.EntityID,
		&projectID, &n.Title, &n.Body, &metadataJSON, &n.IsRead, &n.ReadAt, &n.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.UserNotification{}, errors.New("not found")
		}
		return models.UserNotification{}, err
	}
	if projectID != nil {
		n.ProjectID = *projectID
	}
	if err := json.Unmarshal(metadataJSON, &n.Metadata); err != nil {
		n.Metadata = make(map[string]any)
	}
	return n, nil
}

// NotificationListParams holds parameters for listing notifications.
type NotificationListParams struct {
	TenantID   string
	UserID     string
	UnreadOnly bool
	ProjectID  string
	Limit      int
	HasCursor  bool
	CursorTime time.Time
	CursorID   string
}

// ListNotifications lists notifications for a user.
func (r *UserNotificationsRepo) ListNotifications(ctx context.Context, p NotificationListParams) ([]models.UserNotification, string, error) {
	conds := []string{"tenant_id=$1", "user_id=$2"}
	args := []any{p.TenantID, p.UserID}
	argN := 3

	if p.UnreadOnly {
		conds = append(conds, "is_read=FALSE")
	}
	if p.ProjectID != "" {
		conds = append(conds, "project_id=$"+itoa(argN))
		args = append(args, p.ProjectID)
		argN++
	}
	if p.HasCursor {
		conds = append(conds, "(created_at, id) < ($"+itoa(argN)+", $"+itoa(argN+1)+")")
		args = append(args, p.CursorTime, p.CursorID)
		argN += 2
	}

	limitPlus := p.Limit + 1
	args = append(args, limitPlus)

	sql := `
		SELECT id, tenant_id, user_id, notification_type, entity_type, entity_id,
			project_id, title, body, metadata, is_read, read_at, created_at
		FROM user_notifications
		WHERE ` + strings.Join(conds, " AND ") + `
		ORDER BY created_at DESC, id DESC
		LIMIT $` + itoa(argN)

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	var notifications []models.UserNotification
	for rows.Next() {
		var n models.UserNotification
		var metadataJSON []byte
		var projectID *string
		if err := rows.Scan(&n.ID, &n.TenantID, &n.UserID, &n.NotificationType, &n.EntityType, &n.EntityID,
			&projectID, &n.Title, &n.Body, &metadataJSON, &n.IsRead, &n.ReadAt, &n.CreatedAt); err != nil {
			return nil, "", err
		}
		if projectID != nil {
			n.ProjectID = *projectID
		}
		if err := json.Unmarshal(metadataJSON, &n.Metadata); err != nil {
			n.Metadata = make(map[string]any)
		}
		notifications = append(notifications, n)
	}

	next := ""
	if len(notifications) > p.Limit {
		last := notifications[p.Limit-1]
		next = EncodeCursor(last.CreatedAt, last.ID)
		notifications = notifications[:p.Limit]
	}
	return notifications, next, nil
}

// MarkAsRead marks a single notification as read.
func (r *UserNotificationsRepo) MarkAsRead(ctx context.Context, tenantID, userID, notificationID string) error {
	now := time.Now()
	_, err := r.pool.Exec(ctx, `
		UPDATE user_notifications
		SET is_read = TRUE, read_at = $4
		WHERE tenant_id = $1 AND user_id = $2 AND id = $3 AND is_read = FALSE
	`, tenantID, userID, notificationID, now)
	return err
}

// MarkAllAsRead marks all notifications for a user as read.
func (r *UserNotificationsRepo) MarkAllAsRead(ctx context.Context, tenantID, userID string) (int64, error) {
	now := time.Now()
	result, err := r.pool.Exec(ctx, `
		UPDATE user_notifications
		SET is_read = TRUE, read_at = $3
		WHERE tenant_id = $1 AND user_id = $2 AND is_read = FALSE
	`, tenantID, userID, now)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

// MarkMultipleAsRead marks specific notifications as read.
func (r *UserNotificationsRepo) MarkMultipleAsRead(ctx context.Context, tenantID, userID string, notificationIDs []string) (int64, error) {
	if len(notificationIDs) == 0 {
		return 0, nil
	}
	now := time.Now()
	result, err := r.pool.Exec(ctx, `
		UPDATE user_notifications
		SET is_read = TRUE, read_at = $3
		WHERE tenant_id = $1 AND user_id = $2 AND id = ANY($4) AND is_read = FALSE
	`, tenantID, userID, now, notificationIDs)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

// GetUnreadCount returns the count of unread notifications for a user.
func (r *UserNotificationsRepo) GetUnreadCount(ctx context.Context, tenantID, userID string) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM user_notifications
		WHERE tenant_id = $1 AND user_id = $2 AND is_read = FALSE
	`, tenantID, userID).Scan(&count)
	return count, err
}

// DeleteOldNotifications deletes notifications older than a specified duration.
func (r *UserNotificationsRepo) DeleteOldNotifications(ctx context.Context, tenantID string, olderThan time.Time) (int64, error) {
	result, err := r.pool.Exec(ctx, `
		DELETE FROM user_notifications
		WHERE tenant_id = $1 AND created_at < $2
	`, tenantID, olderThan)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}
