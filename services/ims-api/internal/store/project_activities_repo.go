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
	"github.com/lib/pq"
)

type ProjectActivitiesRepo struct{ pool *pgxpool.Pool }

// CreateActivity creates a new activity in the project's activity feed.
func (r *ProjectActivitiesRepo) CreateActivity(ctx context.Context, a models.ProjectActivity) error {
	metadataJSON, err := json.Marshal(a.Metadata)
	if err != nil {
		return err
	}
	_, err = r.pool.Exec(ctx, `
		INSERT INTO project_activities (
			id, tenant_id, project_id, phase_id, work_order_id, activity_type,
			actor_user_id, actor_email, actor_name, content, metadata,
			attachment_ids, visibility, is_pinned, created_at
		) VALUES ($1,$2,$3,NULLIF($4,''),NULLIF($5,''),$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)
	`, a.ID, a.TenantID, a.ProjectID, a.PhaseID, a.WorkOrderID, a.ActivityType,
		a.ActorUserID, a.ActorEmail, a.ActorName, a.Content, metadataJSON,
		pq.Array(a.AttachmentIDs), a.Visibility, a.IsPinned, a.CreatedAt)
	return err
}

// UpdateActivity updates an activity's content.
func (r *ProjectActivitiesRepo) UpdateActivity(ctx context.Context, tenantID, activityID, content string) error {
	now := time.Now()
	_, err := r.pool.Exec(ctx, `
		UPDATE project_activities
		SET content = $3, edited_at = $4
		WHERE tenant_id = $1 AND id = $2 AND deleted_at IS NULL
	`, tenantID, activityID, content, now)
	return err
}

// DeleteActivity soft-deletes an activity.
func (r *ProjectActivitiesRepo) DeleteActivity(ctx context.Context, tenantID, activityID string) error {
	now := time.Now()
	_, err := r.pool.Exec(ctx, `
		UPDATE project_activities
		SET deleted_at = $3
		WHERE tenant_id = $1 AND id = $2 AND deleted_at IS NULL
	`, tenantID, activityID, now)
	return err
}

// TogglePin toggles the pin status of an activity.
func (r *ProjectActivitiesRepo) TogglePin(ctx context.Context, tenantID, activityID string, pinned bool) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE project_activities
		SET is_pinned = $3
		WHERE tenant_id = $1 AND id = $2 AND deleted_at IS NULL
	`, tenantID, activityID, pinned)
	return err
}

// GetActivity retrieves a specific activity.
func (r *ProjectActivitiesRepo) GetActivity(ctx context.Context, tenantID, activityID string) (models.ProjectActivity, error) {
	var a models.ProjectActivity
	var metadataJSON []byte
	var attachmentIDs []string
	var phaseID, workOrderID *string
	row := r.pool.QueryRow(ctx, `
		SELECT id, tenant_id, project_id, phase_id, work_order_id, activity_type,
			actor_user_id, actor_email, actor_name, content, metadata,
			attachment_ids, visibility, is_pinned, edited_at, deleted_at, created_at
		FROM project_activities
		WHERE tenant_id = $1 AND id = $2
	`, tenantID, activityID)
	if err := row.Scan(&a.ID, &a.TenantID, &a.ProjectID, &phaseID, &workOrderID, &a.ActivityType,
		&a.ActorUserID, &a.ActorEmail, &a.ActorName, &a.Content, &metadataJSON,
		pq.Array(&attachmentIDs), &a.Visibility, &a.IsPinned, &a.EditedAt, &a.DeletedAt, &a.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.ProjectActivity{}, errors.New("not found")
		}
		return models.ProjectActivity{}, err
	}
	if phaseID != nil {
		a.PhaseID = *phaseID
	}
	if workOrderID != nil {
		a.WorkOrderID = *workOrderID
	}
	if err := json.Unmarshal(metadataJSON, &a.Metadata); err != nil {
		a.Metadata = make(map[string]any)
	}
	a.AttachmentIDs = attachmentIDs
	return a, nil
}

// ActivityListParams holds parameters for listing activities.
type ActivityListParams struct {
	TenantID     string
	ProjectID    string
	PhaseID      string
	ActivityType string
	ActorUserID  string
	PinnedOnly   bool
	Limit        int
	HasCursor    bool
	CursorTime   time.Time
	CursorID     string
}

// ListActivities lists activities for a project with optional filters.
func (r *ProjectActivitiesRepo) ListActivities(ctx context.Context, p ActivityListParams) ([]models.ProjectActivity, string, error) {
	conds := []string{"tenant_id=$1", "project_id=$2", "deleted_at IS NULL"}
	args := []any{p.TenantID, p.ProjectID}
	argN := 3

	if p.PhaseID != "" {
		conds = append(conds, "phase_id=$"+itoa(argN))
		args = append(args, p.PhaseID)
		argN++
	}
	if p.ActivityType != "" {
		conds = append(conds, "activity_type=$"+itoa(argN))
		args = append(args, p.ActivityType)
		argN++
	}
	if p.ActorUserID != "" {
		conds = append(conds, "actor_user_id=$"+itoa(argN))
		args = append(args, p.ActorUserID)
		argN++
	}
	if p.PinnedOnly {
		conds = append(conds, "is_pinned=TRUE")
	}
	if p.HasCursor {
		conds = append(conds, "(created_at, id) < ($"+itoa(argN)+", $"+itoa(argN+1)+")")
		args = append(args, p.CursorTime, p.CursorID)
		argN += 2
	}

	limitPlus := p.Limit + 1
	args = append(args, limitPlus)

	sql := `
		SELECT id, tenant_id, project_id, phase_id, work_order_id, activity_type,
			actor_user_id, actor_email, actor_name, content, metadata,
			attachment_ids, visibility, is_pinned, edited_at, deleted_at, created_at
		FROM project_activities
		WHERE ` + strings.Join(conds, " AND ") + `
		ORDER BY created_at DESC, id DESC
		LIMIT $` + itoa(argN)

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	var activities []models.ProjectActivity
	for rows.Next() {
		var a models.ProjectActivity
		var metadataJSON []byte
		var attachmentIDs []string
		var phaseID, workOrderID *string
		if err := rows.Scan(&a.ID, &a.TenantID, &a.ProjectID, &phaseID, &workOrderID, &a.ActivityType,
			&a.ActorUserID, &a.ActorEmail, &a.ActorName, &a.Content, &metadataJSON,
			pq.Array(&attachmentIDs), &a.Visibility, &a.IsPinned, &a.EditedAt, &a.DeletedAt, &a.CreatedAt); err != nil {
			return nil, "", err
		}
		if phaseID != nil {
			a.PhaseID = *phaseID
		}
		if workOrderID != nil {
			a.WorkOrderID = *workOrderID
		}
		if err := json.Unmarshal(metadataJSON, &a.Metadata); err != nil {
			a.Metadata = make(map[string]any)
		}
		a.AttachmentIDs = attachmentIDs
		activities = append(activities, a)
	}

	next := ""
	if len(activities) > p.Limit {
		last := activities[p.Limit-1]
		next = EncodeCursor(last.CreatedAt, last.ID)
		activities = activities[:p.Limit]
	}
	return activities, next, nil
}

// Project Attachments

// CreateAttachment creates a new project attachment.
func (r *ProjectActivitiesRepo) CreateAttachment(ctx context.Context, a models.ProjectAttachment) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO project_attachments (
			id, tenant_id, project_id, phase_id, activity_id, file_name,
			content_type, size_bytes, object_key, uploaded_by_user_id,
			uploaded_by_user_name, created_at
		) VALUES ($1,$2,$3,NULLIF($4,''),NULLIF($5,''),$6,$7,$8,$9,$10,$11,$12)
	`, a.ID, a.TenantID, a.ProjectID, a.PhaseID, a.ActivityID, a.FileName,
		a.ContentType, a.SizeBytes, a.ObjectKey, a.UploadedByUserID,
		a.UploadedByUserName, a.CreatedAt)
	return err
}

// GetAttachment retrieves an attachment by ID.
func (r *ProjectActivitiesRepo) GetAttachment(ctx context.Context, tenantID, attachmentID string) (models.ProjectAttachment, error) {
	var a models.ProjectAttachment
	var phaseID, activityID *string
	row := r.pool.QueryRow(ctx, `
		SELECT id, tenant_id, project_id, phase_id, activity_id, file_name,
			content_type, size_bytes, object_key, uploaded_by_user_id,
			uploaded_by_user_name, created_at
		FROM project_attachments
		WHERE tenant_id = $1 AND id = $2
	`, tenantID, attachmentID)
	if err := row.Scan(&a.ID, &a.TenantID, &a.ProjectID, &phaseID, &activityID, &a.FileName,
		&a.ContentType, &a.SizeBytes, &a.ObjectKey, &a.UploadedByUserID,
		&a.UploadedByUserName, &a.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.ProjectAttachment{}, errors.New("not found")
		}
		return models.ProjectAttachment{}, err
	}
	if phaseID != nil {
		a.PhaseID = *phaseID
	}
	if activityID != nil {
		a.ActivityID = *activityID
	}
	return a, nil
}

// ListAttachments lists attachments for a project.
func (r *ProjectActivitiesRepo) ListAttachments(ctx context.Context, tenantID, projectID string, limit int) ([]models.ProjectAttachment, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, tenant_id, project_id, phase_id, activity_id, file_name,
			content_type, size_bytes, object_key, uploaded_by_user_id,
			uploaded_by_user_name, created_at
		FROM project_attachments
		WHERE tenant_id = $1 AND project_id = $2
		ORDER BY created_at DESC
		LIMIT $3
	`, tenantID, projectID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attachments []models.ProjectAttachment
	for rows.Next() {
		var a models.ProjectAttachment
		var phaseID, activityID *string
		if err := rows.Scan(&a.ID, &a.TenantID, &a.ProjectID, &phaseID, &activityID, &a.FileName,
			&a.ContentType, &a.SizeBytes, &a.ObjectKey, &a.UploadedByUserID,
			&a.UploadedByUserName, &a.CreatedAt); err != nil {
			return nil, err
		}
		if phaseID != nil {
			a.PhaseID = *phaseID
		}
		if activityID != nil {
			a.ActivityID = *activityID
		}
		attachments = append(attachments, a)
	}
	return attachments, nil
}

// DeleteAttachment deletes an attachment.
func (r *ProjectActivitiesRepo) DeleteAttachment(ctx context.Context, tenantID, attachmentID string) error {
	_, err := r.pool.Exec(ctx, `
		DELETE FROM project_attachments
		WHERE tenant_id = $1 AND id = $2
	`, tenantID, attachmentID)
	return err
}
