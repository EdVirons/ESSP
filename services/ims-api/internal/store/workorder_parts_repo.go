package store

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/edvirons/ssp/ims/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WorkOrderPartRepo struct{ pool *pgxpool.Pool }

func (r *WorkOrderPartRepo) Create(ctx context.Context, p models.WorkOrderPart) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO work_order_parts (
			id, tenant_id, school_id, work_order_id, service_shop_id, part_id,
			part_name, part_puk, part_category, device_model_id, is_compatible,
			qty_planned, qty_used, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)
	`, p.ID, p.TenantID, p.SchoolID, p.WorkOrderID, p.ServiceShopID, p.PartID,
		p.PartName, p.PartPUK, p.PartCategory, p.DeviceModelID, p.IsCompatible,
		p.QtyPlanned, p.QtyUsed, p.CreatedAt, p.UpdatedAt)
	return err
}

func (r *WorkOrderPartRepo) GetByID(ctx context.Context, tenantID, schoolID, id string) (models.WorkOrderPart, error) {
	var p models.WorkOrderPart
	row := r.pool.QueryRow(ctx, `
		SELECT id, tenant_id, school_id, work_order_id, service_shop_id, part_id,
		       part_name, part_puk, part_category, device_model_id, is_compatible,
		       qty_planned, qty_used, created_at, updated_at
		FROM work_order_parts
		WHERE tenant_id=$1 AND school_id=$2 AND id=$3
	`, tenantID, schoolID, id)
	if err := row.Scan(&p.ID, &p.TenantID, &p.SchoolID, &p.WorkOrderID, &p.ServiceShopID, &p.PartID,
		&p.PartName, &p.PartPUK, &p.PartCategory, &p.DeviceModelID, &p.IsCompatible,
		&p.QtyPlanned, &p.QtyUsed, &p.CreatedAt, &p.UpdatedAt); err != nil {
		return models.WorkOrderPart{}, errors.New("not found")
	}
	return p, nil
}

type WorkOrderPartListParams struct {
	TenantID        string
	SchoolID        string
	WorkOrderID     string
	Limit           int
	HasCursor       bool
	CursorCreatedAt time.Time
	CursorID        string
}

func (r *WorkOrderPartRepo) List(ctx context.Context, p WorkOrderPartListParams) ([]models.WorkOrderPart, string, error) {
	conds := []string{"tenant_id=$1", "school_id=$2", "work_order_id=$3"}
	args := []any{p.TenantID, p.SchoolID, p.WorkOrderID}
	argN := 4

	if p.HasCursor {
		conds = append(conds, "(created_at, id) < ($"+itoa(argN)+", $"+itoa(argN+1)+")")
		args = append(args, p.CursorCreatedAt, p.CursorID)
		argN += 2
	}

	limitPlus := p.Limit + 1
	args = append(args, limitPlus)

	sql := `
		SELECT id, tenant_id, school_id, work_order_id, service_shop_id, part_id,
		       part_name, part_puk, part_category, device_model_id, is_compatible,
		       qty_planned, qty_used, created_at, updated_at
		FROM work_order_parts
		WHERE ` + strings.Join(conds, " AND ") + `
		ORDER BY created_at DESC, id DESC
		LIMIT $` + itoa(argN)

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	out := []models.WorkOrderPart{}
	for rows.Next() {
		var x models.WorkOrderPart
		if err := rows.Scan(&x.ID, &x.TenantID, &x.SchoolID, &x.WorkOrderID, &x.ServiceShopID, &x.PartID,
			&x.PartName, &x.PartPUK, &x.PartCategory, &x.DeviceModelID, &x.IsCompatible,
			&x.QtyPlanned, &x.QtyUsed, &x.CreatedAt, &x.UpdatedAt); err != nil {
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

// Transactional helpers

type Tx interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

func CreateWorkOrderPartTx(ctx context.Context, tx Tx, p models.WorkOrderPart) error {
	_, err := tx.Exec(ctx, `
		INSERT INTO work_order_parts (
			id, tenant_id, school_id, work_order_id, service_shop_id, part_id,
			part_name, part_puk, part_category, device_model_id, is_compatible,
			qty_planned, qty_used, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)
	`, p.ID, p.TenantID, p.SchoolID, p.WorkOrderID, p.ServiceShopID, p.PartID,
		p.PartName, p.PartPUK, p.PartCategory, p.DeviceModelID, p.IsCompatible,
		p.QtyPlanned, p.QtyUsed, p.CreatedAt, p.UpdatedAt)
	return err
}

func UpdateWorkOrderPartUsedTx(ctx context.Context, tx Tx, tenantID, schoolID, id string, addUsed int64, now time.Time) error {
	_, err := tx.Exec(ctx, `
		UPDATE work_order_parts
		SET qty_used = qty_used + $4, updated_at=$5
		WHERE tenant_id=$1 AND school_id=$2 AND id=$3
		  AND (qty_planned - qty_used) >= $4
	`, tenantID, schoolID, id, addUsed, now)
	return err
}

func UpdateWorkOrderPartPlannedTx(ctx context.Context, tx Tx, tenantID, schoolID, id string, newPlanned int64, now time.Time) error {
	_, err := tx.Exec(ctx, `
		UPDATE work_order_parts
		SET qty_planned = $4, updated_at=$5
		WHERE tenant_id=$1 AND school_id=$2 AND id=$3
	`, tenantID, schoolID, id, newPlanned, now)
	return err
}
