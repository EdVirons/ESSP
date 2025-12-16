package store

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/edvirons/ssp/ims/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WorkOrderApprovalsRepo struct{ pool *pgxpool.Pool }

func (r *WorkOrderApprovalsRepo) Request(ctx context.Context, a models.WorkOrderApproval) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO work_order_approvals (
			id, tenant_id, school_id, work_order_id, approval_type, requested_by_user_id, requested_at, status,
			decided_by_user_id, decided_at, decision_notes
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
	`, a.ID,a.TenantID,a.SchoolID,a.WorkOrderID,a.ApprovalType,a.RequestedByUserID,a.RequestedAt,a.Status,a.DecidedByUserID,a.DecidedAt,a.DecisionNotes)
	return err
}

func (r *WorkOrderApprovalsRepo) Decide(ctx context.Context, tenantID, schoolID, approvalID, deciderID, status, notes string) error {
	now := time.Now().UTC()
	status = strings.TrimSpace(status)
	if status != "approved" && status != "rejected" { return errors.New("invalid status") }
	_, err := r.pool.Exec(ctx, `
		UPDATE work_order_approvals
		SET status=$4, decided_by_user_id=$5, decided_at=$6, decision_notes=$7
		WHERE tenant_id=$1 AND school_id=$2 AND id=$3
	`, tenantID, schoolID, approvalID, status, deciderID, now, notes)
	return err
}

func (r *WorkOrderApprovalsRepo) GetByID(ctx context.Context, tenantID, schoolID, id string) (models.WorkOrderApproval, error) {
	var a models.WorkOrderApproval
	row := r.pool.QueryRow(ctx, `
		SELECT id, tenant_id, school_id, work_order_id, approval_type, requested_by_user_id, requested_at, status,
			decided_by_user_id, decided_at, decision_notes
		FROM work_order_approvals
		WHERE tenant_id=$1 AND school_id=$2 AND id=$3
	`, tenantID, schoolID, id)
	if err := row.Scan(&a.ID,&a.TenantID,&a.SchoolID,&a.WorkOrderID,&a.ApprovalType,&a.RequestedByUserID,&a.RequestedAt,&a.Status,
		&a.DecidedByUserID,&a.DecidedAt,&a.DecisionNotes); err != nil {
		return models.WorkOrderApproval{}, errors.New("not found")
	}
	return a, nil
}
