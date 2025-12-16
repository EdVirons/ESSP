package store

import (
	"context"
	"strings"
	"time"

	"github.com/edvirons/ssp/ims/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BOQRepo struct{ pool *pgxpool.Pool }

func (r *BOQRepo) Create(ctx context.Context, b models.BOQItem) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO boq_items (
			id, tenant_id, project_id, category, description, part_id, qty, unit, estimated_cost_cents, approved, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
	`, b.ID, b.TenantID, b.ProjectID, b.Category, b.Description, b.PartID, b.Qty, b.Unit, b.EstimatedCostCents, b.Approved, b.CreatedAt, b.UpdatedAt)
	return err
}

type BOQListParams struct {
	TenantID string
	ProjectID string
	Approved *bool
	Limit int
	HasCursor bool
	CursorCreatedAt time.Time
	CursorID string
}

func (r *BOQRepo) List(ctx context.Context, p BOQListParams) ([]models.BOQItem, string, error) {
	conds := []string{"tenant_id=$1","project_id=$2"}
	args := []any{p.TenantID,p.ProjectID}
	argN := 3
	if p.Approved != nil {
		conds = append(conds, "approved=$"+itoa(argN)); args=append(args,*p.Approved); argN++
	}
	if p.HasCursor {
		conds = append(conds, "(created_at, id) < ($"+itoa(argN)+", $"+itoa(argN+1)+")")
		args = append(args, p.CursorCreatedAt, p.CursorID); argN += 2
	}
	limitPlus := p.Limit+1
	args = append(args, limitPlus)

	sql := `
		SELECT id, tenant_id, project_id, category, description, part_id, qty, unit, estimated_cost_cents, approved, created_at, updated_at
		FROM boq_items
		WHERE ` + strings.Join(conds," AND ") + `
		ORDER BY created_at DESC, id DESC
		LIMIT $` + itoa(argN)

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil { return nil,"",err }
	defer rows.Close()
	out := []models.BOQItem{}
	for rows.Next() {
		var x models.BOQItem
		if err := rows.Scan(&x.ID,&x.TenantID,&x.ProjectID,&x.Category,&x.Description,&x.PartID,&x.Qty,&x.Unit,&x.EstimatedCostCents,&x.Approved,&x.CreatedAt,&x.UpdatedAt); err != nil { return nil,"",err }
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
