package store

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/edvirons/ssp/ims/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WorkOrderDeliverablesRepo struct{ pool *pgxpool.Pool }

func (r *WorkOrderDeliverablesRepo) Create(ctx context.Context, d models.WorkOrderDeliverable) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO work_order_deliverables (
			id, tenant_id, school_id, work_order_id, phase_id, title, description, status, evidence_attachment_id,
			submitted_by_user_id, submitted_at, reviewed_by_user_id, reviewed_at, review_notes, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)
	`, d.ID, d.TenantID, d.SchoolID, d.WorkOrderID, d.PhaseID, d.Title, d.Description, d.Status, d.EvidenceAttachmentID,
		d.SubmittedByUserID, d.SubmittedAt, d.ReviewedByUserID, d.ReviewedAt, d.ReviewNotes, d.CreatedAt, d.UpdatedAt)
	return err
}

func (r *WorkOrderDeliverablesRepo) GetByID(ctx context.Context, tenantID, schoolID, id string) (models.WorkOrderDeliverable, error) {
	var d models.WorkOrderDeliverable
	row := r.pool.QueryRow(ctx, `
		SELECT id, tenant_id, school_id, work_order_id, phase_id, title, description, status, evidence_attachment_id,
			submitted_by_user_id, submitted_at, reviewed_by_user_id, reviewed_at, review_notes, created_at, updated_at
		FROM work_order_deliverables
		WHERE tenant_id=$1 AND school_id=$2 AND id=$3
	`, tenantID, schoolID, id)
	if err := row.Scan(&d.ID, &d.TenantID, &d.SchoolID, &d.WorkOrderID, &d.PhaseID, &d.Title, &d.Description, &d.Status, &d.EvidenceAttachmentID,
		&d.SubmittedByUserID, &d.SubmittedAt, &d.ReviewedByUserID, &d.ReviewedAt, &d.ReviewNotes, &d.CreatedAt, &d.UpdatedAt); err != nil {
		return models.WorkOrderDeliverable{}, errors.New("not found")
	}
	return d, nil
}

type DeliverableListParams struct {
	TenantID        string
	SchoolID        string
	WorkOrderID     string
	Status          string
	Limit           int
	HasCursor       bool
	CursorCreatedAt time.Time
	CursorID        string
}

func (r *WorkOrderDeliverablesRepo) List(ctx context.Context, p DeliverableListParams) ([]models.WorkOrderDeliverable, string, error) {
	conds := []string{"tenant_id=$1", "school_id=$2", "work_order_id=$3"}
	args := []any{p.TenantID, p.SchoolID, p.WorkOrderID}
	argN := 4
	if p.Status != "" {
		conds = append(conds, "status=$"+itoa(argN))
		args = append(args, p.Status)
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
		SELECT id, tenant_id, school_id, work_order_id, phase_id, title, description, status, evidence_attachment_id,
			submitted_by_user_id, submitted_at, reviewed_by_user_id, reviewed_at, review_notes, created_at, updated_at
		FROM work_order_deliverables
		WHERE ` + strings.Join(conds, " AND ") + `
		ORDER BY created_at DESC, id DESC
		LIMIT $` + itoa(argN)

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	out := []models.WorkOrderDeliverable{}
	for rows.Next() {
		var x models.WorkOrderDeliverable
		if err := rows.Scan(&x.ID, &x.TenantID, &x.SchoolID, &x.WorkOrderID, &x.PhaseID, &x.Title, &x.Description, &x.Status, &x.EvidenceAttachmentID,
			&x.SubmittedByUserID, &x.SubmittedAt, &x.ReviewedByUserID, &x.ReviewedAt, &x.ReviewNotes, &x.CreatedAt, &x.UpdatedAt); err != nil {
			return nil, "", err
		}
		out = append(out, x)
	}
	next := ""
	if len(out) > p.Limit {
		last := out[p.Limit-1]
		next = EncodeCursor(last.CreatedAt, last.ID)
		out = out[:p.Limit]
	}
	return out, next, nil
}

func (r *WorkOrderDeliverablesRepo) MarkSubmitted(ctx context.Context, tenantID, schoolID, id, userID, evidence, notes string) error {
	now := time.Now().UTC()
	_, err := r.pool.Exec(ctx, `
		UPDATE work_order_deliverables
		SET status='submitted', evidence_attachment_id=$4, submitted_by_user_id=$5, submitted_at=$6, updated_at=$6, description=$7
		WHERE tenant_id=$1 AND school_id=$2 AND id=$3
	`, tenantID, schoolID, id, evidence, userID, now, notes)
	return err
}

func (r *WorkOrderDeliverablesRepo) Review(ctx context.Context, tenantID, schoolID, id, reviewerID, status, notes string) error {
	now := time.Now().UTC()
	status = strings.TrimSpace(status)
	if status != "approved" && status != "rejected" {
		return errors.New("invalid status")
	}
	_, err := r.pool.Exec(ctx, `
		UPDATE work_order_deliverables
		SET status=$4, reviewed_by_user_id=$5, reviewed_at=$6, review_notes=$7, updated_at=$6
		WHERE tenant_id=$1 AND school_id=$2 AND id=$3
	`, tenantID, schoolID, id, status, reviewerID, now, notes)
	return err
}

func (r *WorkOrderDeliverablesRepo) CountNotApprovedByWorkOrder(ctx context.Context, tenantID, schoolID, workOrderID string) (int64, error) {
	var c int64
	err := r.pool.QueryRow(ctx, `
		SELECT COUNT(1)
		FROM work_order_deliverables
		WHERE tenant_id=$1 AND school_id=$2 AND work_order_id=$3 AND status <> 'approved'
	`, tenantID, schoolID, workOrderID).Scan(&c)
	return c, err
}
