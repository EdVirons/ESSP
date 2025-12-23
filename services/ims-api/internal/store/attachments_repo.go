package store

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/edvirons/ssp/ims/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AttachmentRepo struct {
	pool *pgxpool.Pool
}

func (r *AttachmentRepo) Create(ctx context.Context, a models.Attachment) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO attachments (
			id, tenant_id, school_id, entity_type, entity_id,
			file_name, content_type, size_bytes, object_key, created_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
	`, a.ID, a.TenantID, a.SchoolID, a.EntityType, a.EntityID, a.FileName, a.ContentType, a.SizeBytes, a.ObjectKey, a.CreatedAt)
	return err
}

func (r *AttachmentRepo) GetByID(ctx context.Context, tenantID, schoolID, id string) (models.Attachment, error) {
	var a models.Attachment
	row := r.pool.QueryRow(ctx, `
		SELECT id, tenant_id, school_id, entity_type, entity_id,
		       file_name, content_type, size_bytes, object_key, created_at
		FROM attachments
		WHERE tenant_id=$1 AND school_id=$2 AND id=$3
	`, tenantID, schoolID, id)

	err := row.Scan(&a.ID, &a.TenantID, &a.SchoolID, &a.EntityType, &a.EntityID,
		&a.FileName, &a.ContentType, &a.SizeBytes, &a.ObjectKey, &a.CreatedAt)
	if err != nil {
		return models.Attachment{}, errors.New("not found")
	}
	return a, nil
}

type AttachmentListParams struct {
	TenantID   string
	SchoolID   string
	EntityType string
	EntityID   string
	Limit      int

	HasCursor       bool
	CursorCreatedAt time.Time
	CursorID        string
}

func (r *AttachmentRepo) List(ctx context.Context, p AttachmentListParams) ([]models.Attachment, string, error) {
	conds := []string{"tenant_id=$1", "school_id=$2"}
	args := []any{p.TenantID, p.SchoolID}
	argN := 3

	if p.EntityType != "" {
		conds = append(conds, "entity_type=$"+itoa(argN))
		args = append(args, p.EntityType)
		argN++
	}
	if p.EntityID != "" {
		conds = append(conds, "entity_id=$"+itoa(argN))
		args = append(args, p.EntityID)
		argN++
	}
	if p.HasCursor {
		conds = append(conds, "(created_at, id) < ($"+itoa(argN)+", $"+itoa(argN+1)+")")
		args = append(args, p.CursorCreatedAt, p.CursorID)
		argN += 2
	}

	limitPlus := p.Limit + 1
	args = append(args, limitPlus)

	sql := `
		SELECT id, tenant_id, school_id, entity_type, entity_id,
		       file_name, content_type, size_bytes, object_key, created_at
		FROM attachments
		WHERE ` + strings.Join(conds, " AND ") + `
		ORDER BY created_at DESC, id DESC
		LIMIT $` + itoa(argN)

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	out := []models.Attachment{}
	for rows.Next() {
		var a models.Attachment
		if err := rows.Scan(&a.ID, &a.TenantID, &a.SchoolID, &a.EntityType, &a.EntityID,
			&a.FileName, &a.ContentType, &a.SizeBytes, &a.ObjectKey, &a.CreatedAt); err != nil {
			return nil, "", err
		}
		out = append(out, a)
	}

	next := ""
	if len(out) > p.Limit {
		last := out[p.Limit-1]
		next = EncodeCursor(last.CreatedAt, last.ID)
		out = out[:p.Limit]
	}
	return out, next, nil
}
