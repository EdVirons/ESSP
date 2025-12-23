package store

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/edvirons/ssp/ims/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PartsSnapshotRepo struct{ pool *pgxpool.Pool }

func (r *PartsSnapshotRepo) Upsert(ctx context.Context, p models.PartSnapshot) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO parts_snapshot (
			tenant_id, part_id, puk, name, category, unit, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7)
		ON CONFLICT (tenant_id, part_id)
		DO UPDATE SET
		  puk=EXCLUDED.puk,
		  name=EXCLUDED.name,
		  category=EXCLUDED.category,
		  unit=EXCLUDED.unit,
		  updated_at=EXCLUDED.updated_at
	`, p.TenantID, p.PartID, p.PUK, p.Name, p.Category, p.Unit, p.UpdatedAt)
	return err
}

func (r *PartsSnapshotRepo) Get(ctx context.Context, tenantID, partID string) (models.PartSnapshot, error) {
	var p models.PartSnapshot
	row := r.pool.QueryRow(ctx, `
		SELECT tenant_id, part_id, puk, name, category, unit, updated_at
		FROM parts_snapshot
		WHERE tenant_id=$1 AND part_id=$2
	`, tenantID, partID)
	if err := row.Scan(&p.TenantID, &p.PartID, &p.PUK, &p.Name, &p.Category, &p.Unit, &p.UpdatedAt); err != nil {
		return models.PartSnapshot{}, errors.New("not found")
	}
	return p, nil
}

func NewPartSnapshot(tenantID, partID string) models.PartSnapshot {
	return models.PartSnapshot{TenantID: tenantID, PartID: partID, UpdatedAt: time.Now().UTC()}
}

type PartSnapshotListParams struct {
	TenantID string
	Category string
	Query    string
	Limit    int
	Offset   int
}

func (r *PartsSnapshotRepo) List(ctx context.Context, p PartSnapshotListParams) ([]models.PartSnapshot, int, error) {
	conds := []string{"tenant_id=$1"}
	args := []any{p.TenantID}
	argN := 2

	if p.Category != "" {
		conds = append(conds, "category=$"+itoa(argN))
		args = append(args, p.Category)
		argN++
	}
	if p.Query != "" {
		conds = append(conds, "(name ILIKE $"+itoa(argN)+" OR puk ILIKE $"+itoa(argN)+")")
		args = append(args, "%"+p.Query+"%")
		argN++
	}

	// Count total
	countSQL := "SELECT COUNT(*) FROM parts_snapshot WHERE " + strings.Join(conds, " AND ")
	var total int
	if err := r.pool.QueryRow(ctx, countSQL, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Fetch items
	limit := p.Limit
	if limit <= 0 {
		limit = 50
	}
	offset := p.Offset
	if offset < 0 {
		offset = 0
	}

	args = append(args, limit, offset)
	sql := `
		SELECT tenant_id, part_id, puk, name, category, unit, updated_at
		FROM parts_snapshot
		WHERE ` + strings.Join(conds, " AND ") + `
		ORDER BY name ASC
		LIMIT $` + itoa(argN) + ` OFFSET $` + itoa(argN+1)

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var items []models.PartSnapshot
	for rows.Next() {
		var pt models.PartSnapshot
		if err := rows.Scan(&pt.TenantID, &pt.PartID, &pt.PUK, &pt.Name, &pt.Category, &pt.Unit, &pt.UpdatedAt); err != nil {
			return nil, 0, err
		}
		items = append(items, pt)
	}

	return items, total, nil
}

func (r *PartsSnapshotRepo) Count(ctx context.Context, tenantID string) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM parts_snapshot WHERE tenant_id=$1", tenantID).Scan(&count)
	return count, err
}
