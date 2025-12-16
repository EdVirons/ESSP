package store

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/edvirons/ssp/ims/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PartRepo struct{ pool *pgxpool.Pool }

func (r *PartRepo) Create(ctx context.Context, p models.Part) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO parts (id, tenant_id, sku, name, category, description, unit_cost_cents, supplier, supplier_sku, active, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
	`, p.ID, p.TenantID, p.SKU, p.Name, p.Category, p.Description, p.UnitCostCents, p.Supplier, p.SupplierSku, p.Active, p.CreatedAt, p.UpdatedAt)
	return err
}

func (r *PartRepo) GetByID(ctx context.Context, tenantID, id string) (models.Part, error) {
	var p models.Part
	row := r.pool.QueryRow(ctx, `
		SELECT id, tenant_id, sku, name, category, description, unit_cost_cents, supplier, supplier_sku, active, created_at, updated_at
		FROM parts WHERE tenant_id=$1 AND id=$2
	`, tenantID, id)
	if err := row.Scan(&p.ID, &p.TenantID, &p.SKU, &p.Name, &p.Category, &p.Description, &p.UnitCostCents, &p.Supplier, &p.SupplierSku, &p.Active, &p.CreatedAt, &p.UpdatedAt); err != nil {
		return models.Part{}, errors.New("not found")
	}
	return p, nil
}

func (r *PartRepo) Update(ctx context.Context, p models.Part) error {
	result, err := r.pool.Exec(ctx, `
		UPDATE parts SET
			name = $3,
			category = $4,
			description = $5,
			unit_cost_cents = $6,
			supplier = $7,
			supplier_sku = $8,
			active = $9,
			updated_at = $10
		WHERE tenant_id = $1 AND id = $2
	`, p.TenantID, p.ID, p.Name, p.Category, p.Description, p.UnitCostCents, p.Supplier, p.SupplierSku, p.Active, p.UpdatedAt)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errors.New("not found")
	}
	return nil
}

func (r *PartRepo) Delete(ctx context.Context, tenantID, id string) error {
	result, err := r.pool.Exec(ctx, `
		DELETE FROM parts WHERE tenant_id = $1 AND id = $2
	`, tenantID, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errors.New("not found")
	}
	return nil
}

type PartListParams struct {
	TenantID        string
	Q               string
	Category        string
	Active          *bool
	Limit           int
	HasCursor       bool
	CursorCreatedAt time.Time
	CursorID        string
}

func (r *PartRepo) List(ctx context.Context, p PartListParams) ([]models.Part, string, error) {
	conds := []string{"tenant_id=$1"}
	args := []any{p.TenantID}
	argN := 2

	if p.Q != "" {
		conds = append(conds, "(sku ILIKE $"+itoa(argN)+" OR name ILIKE $"+itoa(argN)+" OR supplier ILIKE $"+itoa(argN)+")")
		args = append(args, "%"+p.Q+"%")
		argN++
	}
	if p.Category != "" {
		conds = append(conds, "category = $"+itoa(argN))
		args = append(args, p.Category)
		argN++
	}
	if p.Active != nil {
		conds = append(conds, "active = $"+itoa(argN))
		args = append(args, *p.Active)
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
		SELECT id, tenant_id, sku, name, category, description, unit_cost_cents, supplier, supplier_sku, active, created_at, updated_at
		FROM parts
		WHERE ` + strings.Join(conds, " AND ") + `
		ORDER BY created_at DESC, id DESC
		LIMIT $` + itoa(argN)

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	out := []models.Part{}
	for rows.Next() {
		var x models.Part
		if err := rows.Scan(&x.ID, &x.TenantID, &x.SKU, &x.Name, &x.Category, &x.Description, &x.UnitCostCents, &x.Supplier, &x.SupplierSku, &x.Active, &x.CreatedAt, &x.UpdatedAt); err != nil {
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

// Count returns total count of parts for tenant
func (r *PartRepo) Count(ctx context.Context, tenantID string) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM parts WHERE tenant_id = $1`, tenantID).Scan(&count)
	return count, err
}

// CountByCategory returns count of parts grouped by category
func (r *PartRepo) CountByCategory(ctx context.Context, tenantID string) (map[string]int, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT COALESCE(category, ''), COUNT(*) as count
		FROM parts
		WHERE tenant_id = $1
		GROUP BY category
	`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	counts := make(map[string]int)
	for rows.Next() {
		var cat string
		var count int
		if err := rows.Scan(&cat, &count); err != nil {
			return nil, err
		}
		counts[cat] = count
	}
	return counts, nil
}

// GetCategories returns distinct categories
func (r *PartRepo) GetCategories(ctx context.Context, tenantID string) ([]string, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT DISTINCT category FROM parts WHERE tenant_id = $1 AND category != '' ORDER BY category
	`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var cat string
		if err := rows.Scan(&cat); err != nil {
			return nil, err
		}
		categories = append(categories, cat)
	}
	return categories, nil
}
