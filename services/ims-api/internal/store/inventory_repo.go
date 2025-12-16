package store

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/edvirons/ssp/ims/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type InventoryRepo struct{ pool *pgxpool.Pool }

func (r *InventoryRepo) Upsert(ctx context.Context, i models.InventoryItem) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO inventory (id, tenant_id, service_shop_id, part_id, qty_available, qty_reserved, reorder_threshold, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		ON CONFLICT (tenant_id, service_shop_id, part_id)
		DO UPDATE SET qty_available=EXCLUDED.qty_available, qty_reserved=EXCLUDED.qty_reserved, reorder_threshold=EXCLUDED.reorder_threshold, updated_at=EXCLUDED.updated_at
	`, i.ID, i.TenantID, i.ServiceShopID, i.PartID, i.QtyAvailable, i.QtyReserved, i.ReorderThreshold, i.UpdatedAt)
	return err
}

func (r *InventoryRepo) Get(ctx context.Context, tenantID, shopID, partID string) (models.InventoryItem, error) {
	var i models.InventoryItem
	row := r.pool.QueryRow(ctx, `
		SELECT id, tenant_id, service_shop_id, part_id, qty_available, qty_reserved, reorder_threshold, updated_at
		FROM inventory
		WHERE tenant_id=$1 AND service_shop_id=$2 AND part_id=$3
	`, tenantID, shopID, partID)
	if err := row.Scan(&i.ID,&i.TenantID,&i.ServiceShopID,&i.PartID,&i.QtyAvailable,&i.QtyReserved,&i.ReorderThreshold,&i.UpdatedAt); err != nil {
		return models.InventoryItem{}, errors.New("not found")
	}
	return i, nil
}

type InventoryListParams struct {
	TenantID string
	ShopID string
	PartID string
	Limit int
	HasCursor bool
	CursorUpdatedAt time.Time
	CursorID string
}

func (r *InventoryRepo) List(ctx context.Context, p InventoryListParams) ([]models.InventoryItem, string, error) {
	conds := []string{"tenant_id=$1"}
	args := []any{p.TenantID}
	argN := 2
	if p.ShopID != "" {
		conds = append(conds, "service_shop_id=$"+itoa(argN)); args = append(args, p.ShopID); argN++
	}
	if p.PartID != "" {
		conds = append(conds, "part_id=$"+itoa(argN)); args = append(args, p.PartID); argN++
	}
	if p.HasCursor {
		conds = append(conds, "(updated_at, id) < ($"+itoa(argN)+", $"+itoa(argN+1)+")")
		args = append(args, p.CursorUpdatedAt, p.CursorID); argN += 2
	}
	limitPlus := p.Limit+1
	args = append(args, limitPlus)

	sql := `
		SELECT id, tenant_id, service_shop_id, part_id, qty_available, qty_reserved, reorder_threshold, updated_at
		FROM inventory
		WHERE ` + strings.Join(conds, " AND ") + `
		ORDER BY updated_at DESC, id DESC
		LIMIT $` + itoa(argN)
	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil { return nil,"",err }
	defer rows.Close()
	out := []models.InventoryItem{}
	for rows.Next() {
		var x models.InventoryItem
		if err := rows.Scan(&x.ID,&x.TenantID,&x.ServiceShopID,&x.PartID,&x.QtyAvailable,&x.QtyReserved,&x.ReorderThreshold,&x.UpdatedAt); err != nil {
			return nil,"",err
		}
		out = append(out, x)
	}
	next := ""
	if len(out) > p.Limit {
		last := out[p.Limit-1]
		next = EncodeCursor(last.UpdatedAt, last.ID)
		out = out[:p.Limit]
	}
	return out,next,nil
}
func (r *InventoryRepo) Reserve(ctx context.Context, tenantID, shopID, partID string, qty int64) error {
	// Ensures (qty_available - qty_reserved) >= qty
	_, err := r.pool.Exec(ctx, `
		UPDATE inventory
		SET qty_reserved = qty_reserved + $4, updated_at = $5
		WHERE tenant_id=$1 AND service_shop_id=$2 AND part_id=$3
		  AND (qty_available - qty_reserved) >= $4
	`, tenantID, shopID, partID, qty, time.Now().UTC())
	return err
}

func (r *InventoryRepo) Release(ctx context.Context, tenantID, shopID, partID string, qty int64) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE inventory
		SET qty_reserved = GREATEST(qty_reserved - $4, 0), updated_at = $5
		WHERE tenant_id=$1 AND service_shop_id=$2 AND part_id=$3
	`, tenantID, shopID, partID, qty, time.Now().UTC())
	return err
}

func (r *InventoryRepo) Consume(ctx context.Context, tenantID, shopID, partID string, qty int64) error {
	// Decrease both reserved and available (on-hand)
	_, err := r.pool.Exec(ctx, `
		UPDATE inventory
		SET qty_reserved = GREATEST(qty_reserved - $4, 0),
		    qty_available = GREATEST(qty_available - $4, 0),
		    updated_at = $5
		WHERE tenant_id=$1 AND service_shop_id=$2 AND part_id=$3
	`, tenantID, shopID, partID, qty, time.Now().UTC())
	return err
}

