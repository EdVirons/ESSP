package store

import (
	"context"
	"time"

	"github.com/edvirons/ssp/ims/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

// WorkOrderReworkRepo handles work order rework history operations.
type WorkOrderReworkRepo struct {
	pool *pgxpool.Pool
}

// Create inserts a new rework history entry.
func (r *WorkOrderReworkRepo) Create(ctx context.Context, entry models.WorkOrderReworkHistory) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO work_order_rework_history (
			id, tenant_id, school_id, work_order_id, from_status, to_status,
			rejection_reason, rejection_category, rejected_by_user_id, rejected_by_name,
			rework_sequence, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`, entry.ID, entry.TenantID, entry.SchoolID, entry.WorkOrderID, entry.FromStatus, entry.ToStatus,
		entry.RejectionReason, entry.RejectionCategory, entry.RejectedByUserID, entry.RejectedByName,
		entry.ReworkSequence, entry.CreatedAt)
	return err
}

// GetReworkCount returns the current rework count for a work order.
func (r *WorkOrderReworkRepo) GetReworkCount(ctx context.Context, tenantID, schoolID, workOrderID string) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM work_order_rework_history
		WHERE tenant_id = $1 AND school_id = $2 AND work_order_id = $3
	`, tenantID, schoolID, workOrderID).Scan(&count)
	return count, err
}

// GetNextReworkSequence returns the next sequence number for a work order.
func (r *WorkOrderReworkRepo) GetNextReworkSequence(ctx context.Context, tenantID, schoolID, workOrderID string) (int, error) {
	var maxSeq *int
	err := r.pool.QueryRow(ctx, `
		SELECT MAX(rework_sequence)
		FROM work_order_rework_history
		WHERE tenant_id = $1 AND school_id = $2 AND work_order_id = $3
	`, tenantID, schoolID, workOrderID).Scan(&maxSeq)
	if err != nil {
		return 1, err
	}
	if maxSeq == nil {
		return 1, nil
	}
	return *maxSeq + 1, nil
}

// ReworkHistoryListParams holds parameters for listing rework history.
type ReworkHistoryListParams struct {
	TenantID    string
	SchoolID    string
	WorkOrderID string
	Limit       int
	HasCursor   bool
	CursorTime  time.Time
	CursorID    string
}

// List retrieves rework history for a work order.
func (r *WorkOrderReworkRepo) List(ctx context.Context, p ReworkHistoryListParams) ([]models.WorkOrderReworkHistory, string, error) {
	args := []any{p.TenantID, p.SchoolID, p.WorkOrderID}
	argN := 4

	cursorCond := ""
	if p.HasCursor {
		cursorCond = " AND (created_at, id) < ($" + itoa(argN) + ", $" + itoa(argN+1) + ")"
		args = append(args, p.CursorTime, p.CursorID)
		argN += 2
	}

	limitPlus := p.Limit + 1
	args = append(args, limitPlus)

	query := `
		SELECT id, tenant_id, school_id, work_order_id, from_status, to_status,
			rejection_reason, rejection_category, rejected_by_user_id, rejected_by_name,
			rework_sequence, created_at
		FROM work_order_rework_history
		WHERE tenant_id = $1 AND school_id = $2 AND work_order_id = $3` + cursorCond + `
		ORDER BY created_at DESC, id DESC
		LIMIT $` + itoa(argN)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	out := []models.WorkOrderReworkHistory{}
	for rows.Next() {
		var h models.WorkOrderReworkHistory
		if err := rows.Scan(&h.ID, &h.TenantID, &h.SchoolID, &h.WorkOrderID, &h.FromStatus, &h.ToStatus,
			&h.RejectionReason, &h.RejectionCategory, &h.RejectedByUserID, &h.RejectedByName,
			&h.ReworkSequence, &h.CreatedAt); err != nil {
			return nil, "", err
		}
		out = append(out, h)
	}

	next := ""
	if len(out) > p.Limit {
		last := out[p.Limit-1]
		next = EncodeCursor(last.CreatedAt, last.ID)
		out = out[:p.Limit]
	}
	return out, next, nil
}
