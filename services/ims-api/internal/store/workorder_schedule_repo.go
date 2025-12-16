package store

import (
	"context"
	"strings"
	"time"

	"github.com/edvirons/ssp/ims/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WorkOrderScheduleRepo struct{ pool *pgxpool.Pool }

func (r *WorkOrderScheduleRepo) Create(ctx context.Context, s models.WorkOrderSchedule) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO work_order_schedules (
			id, tenant_id, school_id, work_order_id, scheduled_start, scheduled_end, timezone, notes, created_by_user_id, created_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
	`, s.ID,s.TenantID,s.SchoolID,s.WorkOrderID,s.ScheduledStart,s.ScheduledEnd,s.Timezone,s.Notes,s.CreatedByUserID,s.CreatedAt)
	return err
}

type ScheduleListParams struct{
	TenantID string
	SchoolID string
	WorkOrderID string
	Limit int
	HasCursor bool
	CursorCreatedAt time.Time
	CursorID string
}

func (r *WorkOrderScheduleRepo) List(ctx context.Context, p ScheduleListParams) ([]models.WorkOrderSchedule, string, error) {
	conds := []string{"tenant_id=$1","school_id=$2","work_order_id=$3"}
	args := []any{p.TenantID,p.SchoolID,p.WorkOrderID}
	argN := 4
	if p.HasCursor {
		conds = append(conds, "(created_at, id) < ($"+itoa(argN)+", $"+itoa(argN+1)+")")
		args = append(args, p.CursorCreatedAt, p.CursorID); argN += 2
	}
	limitPlus := p.Limit + 1
	args = append(args, limitPlus)

	sql := `
		SELECT id, tenant_id, school_id, work_order_id, scheduled_start, scheduled_end, timezone, notes, created_by_user_id, created_at
		FROM work_order_schedules
		WHERE ` + strings.Join(conds," AND ") + `
		ORDER BY created_at DESC, id DESC
		LIMIT $` + itoa(argN)

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil { return nil,"",err }
	defer rows.Close()

	out := []models.WorkOrderSchedule{}
	for rows.Next() {
		var x models.WorkOrderSchedule
		if err := rows.Scan(&x.ID,&x.TenantID,&x.SchoolID,&x.WorkOrderID,&x.ScheduledStart,&x.ScheduledEnd,&x.Timezone,&x.Notes,&x.CreatedByUserID,&x.CreatedAt); err != nil {
			return nil,"",err
		}
		out = append(out,x)
	}
	next := ""
	if len(out) > p.Limit {
		last := out[p.Limit-1]
		next = EncodeCursor(last.CreatedAt, last.ID)
		out = out[:p.Limit]
	}
	return out,next,nil
}
