package store

import (
	"context"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrgUnitSnapshot struct {
	TenantID  string    `json:"tenantId"`
	OrgUnitID string    `json:"orgUnitId"`
	ParentID  string    `json:"parentId,omitempty"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	Kind      string    `json:"kind"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type OrgTreeNode struct {
	OrgUnitSnapshot
	Children []OrgTreeNode `json:"children,omitempty"`
}

type OrgUnitsSnapshotRepo struct {
	pool *pgxpool.Pool
}

type OrgUnitSnapshotListParams struct {
	TenantID string
	Query    string
	Kind     string
	ParentID string
	Limit    int
	Offset   int
}

func (r *OrgUnitsSnapshotRepo) Upsert(ctx context.Context, o OrgUnitSnapshot) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO org_units_snapshot (tenant_id, org_unit_id, parent_id, code, name, kind, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (tenant_id, org_unit_id) DO UPDATE SET
			parent_id = EXCLUDED.parent_id,
			code = EXCLUDED.code,
			name = EXCLUDED.name,
			kind = EXCLUDED.kind,
			updated_at = EXCLUDED.updated_at
	`, o.TenantID, o.OrgUnitID, o.ParentID, o.Code, o.Name, o.Kind, o.UpdatedAt)
	return err
}

func (r *OrgUnitsSnapshotRepo) UpsertBatch(ctx context.Context, units []OrgUnitSnapshot) (int, error) {
	if len(units) == 0 {
		return 0, nil
	}

	batch := &pgx.Batch{}
	for _, o := range units {
		batch.Queue(`
			INSERT INTO org_units_snapshot (tenant_id, org_unit_id, parent_id, code, name, kind, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			ON CONFLICT (tenant_id, org_unit_id) DO UPDATE SET
				parent_id = EXCLUDED.parent_id,
				code = EXCLUDED.code,
				name = EXCLUDED.name,
				kind = EXCLUDED.kind,
				updated_at = EXCLUDED.updated_at
		`, o.TenantID, o.OrgUnitID, o.ParentID, o.Code, o.Name, o.Kind, o.UpdatedAt)
	}

	br := r.pool.SendBatch(ctx, batch)
	defer br.Close()

	for range units {
		if _, err := br.Exec(); err != nil {
			return 0, err
		}
	}
	return len(units), nil
}

func (r *OrgUnitsSnapshotRepo) List(ctx context.Context, params OrgUnitSnapshotListParams) ([]OrgUnitSnapshot, int, error) {
	if params.Limit == 0 {
		params.Limit = 50
	}

	countQuery := `SELECT COUNT(*) FROM org_units_snapshot WHERE tenant_id = $1`
	args := []any{params.TenantID}
	argIdx := 2

	if params.Query != "" {
		countQuery += ` AND (name ILIKE $` + strconv.Itoa(argIdx) + ` OR code ILIKE $` + strconv.Itoa(argIdx) + `)`
		args = append(args, "%"+params.Query+"%")
		argIdx++
	}
	if params.Kind != "" {
		countQuery += ` AND kind = $` + strconv.Itoa(argIdx)
		args = append(args, params.Kind)
		argIdx++
	}
	if params.ParentID != "" {
		countQuery += ` AND parent_id = $` + strconv.Itoa(argIdx)
		args = append(args, params.ParentID)
		argIdx++
	}

	var total int
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	listQuery := `SELECT tenant_id, org_unit_id, parent_id, code, name, kind, updated_at
		FROM org_units_snapshot WHERE tenant_id = $1`
	args = []any{params.TenantID}
	argIdx = 2

	if params.Query != "" {
		listQuery += ` AND (name ILIKE $` + strconv.Itoa(argIdx) + ` OR code ILIKE $` + strconv.Itoa(argIdx) + `)`
		args = append(args, "%"+params.Query+"%")
		argIdx++
	}
	if params.Kind != "" {
		listQuery += ` AND kind = $` + strconv.Itoa(argIdx)
		args = append(args, params.Kind)
		argIdx++
	}
	if params.ParentID != "" {
		listQuery += ` AND parent_id = $` + strconv.Itoa(argIdx)
		args = append(args, params.ParentID)
		argIdx++
	}

	listQuery += ` ORDER BY name ASC LIMIT $` + strconv.Itoa(argIdx) + ` OFFSET $` + strconv.Itoa(argIdx+1)
	args = append(args, params.Limit, params.Offset)

	rows, err := r.pool.Query(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var items []OrgUnitSnapshot
	for rows.Next() {
		var o OrgUnitSnapshot
		if err := rows.Scan(&o.TenantID, &o.OrgUnitID, &o.ParentID, &o.Code, &o.Name, &o.Kind, &o.UpdatedAt); err != nil {
			return nil, 0, err
		}
		items = append(items, o)
	}
	return items, total, rows.Err()
}

func (r *OrgUnitsSnapshotRepo) Get(ctx context.Context, tenantID, orgUnitID string) (OrgUnitSnapshot, error) {
	var o OrgUnitSnapshot
	err := r.pool.QueryRow(ctx, `
		SELECT tenant_id, org_unit_id, parent_id, code, name, kind, updated_at
		FROM org_units_snapshot WHERE tenant_id = $1 AND org_unit_id = $2
	`, tenantID, orgUnitID).Scan(&o.TenantID, &o.OrgUnitID, &o.ParentID, &o.Code, &o.Name, &o.Kind, &o.UpdatedAt)
	return o, err
}

func (r *OrgUnitsSnapshotRepo) GetTree(ctx context.Context, tenantID string) ([]OrgTreeNode, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT tenant_id, org_unit_id, parent_id, code, name, kind, updated_at
		FROM org_units_snapshot WHERE tenant_id = $1 ORDER BY name
	`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var allUnits []OrgUnitSnapshot
	for rows.Next() {
		var o OrgUnitSnapshot
		if err := rows.Scan(&o.TenantID, &o.OrgUnitID, &o.ParentID, &o.Code, &o.Name, &o.Kind, &o.UpdatedAt); err != nil {
			return nil, err
		}
		allUnits = append(allUnits, o)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Build tree structure
	nodeMap := make(map[string]*OrgTreeNode)
	var roots []OrgTreeNode

	// First pass: create nodes
	for _, u := range allUnits {
		node := OrgTreeNode{OrgUnitSnapshot: u}
		nodeMap[u.OrgUnitID] = &node
	}

	// Second pass: link children
	for _, u := range allUnits {
		node := nodeMap[u.OrgUnitID]
		if u.ParentID == "" {
			roots = append(roots, *node)
		} else if parent, ok := nodeMap[u.ParentID]; ok {
			parent.Children = append(parent.Children, *node)
		} else {
			// Parent not found, treat as root
			roots = append(roots, *node)
		}
	}

	return roots, nil
}

func (r *OrgUnitsSnapshotRepo) Count(ctx context.Context, tenantID string) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM org_units_snapshot WHERE tenant_id = $1`, tenantID).Scan(&count)
	return count, err
}

func (r *OrgUnitsSnapshotRepo) Delete(ctx context.Context, tenantID, orgUnitID string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM org_units_snapshot WHERE tenant_id = $1 AND org_unit_id = $2`, tenantID, orgUnitID)
	return err
}

func (r *OrgUnitsSnapshotRepo) DeleteAll(ctx context.Context, tenantID string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM org_units_snapshot WHERE tenant_id = $1`, tenantID)
	return err
}
