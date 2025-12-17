package store

import (
	"context"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TeamSnapshot struct {
	TenantID    string    `json:"tenantId"`
	TeamID      string    `json:"teamId"`
	OrgUnitID   string    `json:"orgUnitId,omitempty"`
	Key         string    `json:"key"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type TeamsSnapshotRepo struct {
	pool *pgxpool.Pool
}

type TeamSnapshotListParams struct {
	TenantID  string
	Query     string
	Key       string
	OrgUnitID string
	Limit     int
	Offset    int
}

func (r *TeamsSnapshotRepo) Upsert(ctx context.Context, t TeamSnapshot) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO teams_snapshot (tenant_id, team_id, org_unit_id, key, name, description, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (tenant_id, team_id) DO UPDATE SET
			org_unit_id = EXCLUDED.org_unit_id,
			key = EXCLUDED.key,
			name = EXCLUDED.name,
			description = EXCLUDED.description,
			updated_at = EXCLUDED.updated_at
	`, t.TenantID, t.TeamID, t.OrgUnitID, t.Key, t.Name, t.Description, t.UpdatedAt)
	return err
}

func (r *TeamsSnapshotRepo) UpsertBatch(ctx context.Context, teams []TeamSnapshot) (int, error) {
	if len(teams) == 0 {
		return 0, nil
	}

	batch := &pgx.Batch{}
	for _, t := range teams {
		batch.Queue(`
			INSERT INTO teams_snapshot (tenant_id, team_id, org_unit_id, key, name, description, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			ON CONFLICT (tenant_id, team_id) DO UPDATE SET
				org_unit_id = EXCLUDED.org_unit_id,
				key = EXCLUDED.key,
				name = EXCLUDED.name,
				description = EXCLUDED.description,
				updated_at = EXCLUDED.updated_at
		`, t.TenantID, t.TeamID, t.OrgUnitID, t.Key, t.Name, t.Description, t.UpdatedAt)
	}

	br := r.pool.SendBatch(ctx, batch)
	defer br.Close()

	for range teams {
		if _, err := br.Exec(); err != nil {
			return 0, err
		}
	}
	return len(teams), nil
}

func (r *TeamsSnapshotRepo) List(ctx context.Context, params TeamSnapshotListParams) ([]TeamSnapshot, int, error) {
	if params.Limit == 0 {
		params.Limit = 50
	}

	countQuery := `SELECT COUNT(*) FROM teams_snapshot WHERE tenant_id = $1`
	args := []any{params.TenantID}
	argIdx := 2

	if params.Query != "" {
		countQuery += ` AND (name ILIKE $` + strconv.Itoa(argIdx) + ` OR key ILIKE $` + strconv.Itoa(argIdx) + `)`
		args = append(args, "%"+params.Query+"%")
		argIdx++
	}
	if params.Key != "" {
		countQuery += ` AND key = $` + strconv.Itoa(argIdx)
		args = append(args, params.Key)
		argIdx++
	}
	if params.OrgUnitID != "" {
		countQuery += ` AND org_unit_id = $` + strconv.Itoa(argIdx)
		args = append(args, params.OrgUnitID)
		argIdx++
	}

	var total int
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	listQuery := `SELECT tenant_id, team_id, org_unit_id, key, name, description, updated_at
		FROM teams_snapshot WHERE tenant_id = $1`
	args = []any{params.TenantID}
	argIdx = 2

	if params.Query != "" {
		listQuery += ` AND (name ILIKE $` + strconv.Itoa(argIdx) + ` OR key ILIKE $` + strconv.Itoa(argIdx) + `)`
		args = append(args, "%"+params.Query+"%")
		argIdx++
	}
	if params.Key != "" {
		listQuery += ` AND key = $` + strconv.Itoa(argIdx)
		args = append(args, params.Key)
		argIdx++
	}
	if params.OrgUnitID != "" {
		listQuery += ` AND org_unit_id = $` + strconv.Itoa(argIdx)
		args = append(args, params.OrgUnitID)
		argIdx++
	}

	listQuery += ` ORDER BY name ASC LIMIT $` + strconv.Itoa(argIdx) + ` OFFSET $` + strconv.Itoa(argIdx+1)
	args = append(args, params.Limit, params.Offset)

	rows, err := r.pool.Query(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var items []TeamSnapshot
	for rows.Next() {
		var t TeamSnapshot
		if err := rows.Scan(&t.TenantID, &t.TeamID, &t.OrgUnitID, &t.Key, &t.Name, &t.Description, &t.UpdatedAt); err != nil {
			return nil, 0, err
		}
		items = append(items, t)
	}
	return items, total, rows.Err()
}

func (r *TeamsSnapshotRepo) Get(ctx context.Context, tenantID, teamID string) (TeamSnapshot, error) {
	var t TeamSnapshot
	err := r.pool.QueryRow(ctx, `
		SELECT tenant_id, team_id, org_unit_id, key, name, description, updated_at
		FROM teams_snapshot WHERE tenant_id = $1 AND team_id = $2
	`, tenantID, teamID).Scan(&t.TenantID, &t.TeamID, &t.OrgUnitID, &t.Key, &t.Name, &t.Description, &t.UpdatedAt)
	return t, err
}

func (r *TeamsSnapshotRepo) GetByKey(ctx context.Context, tenantID, key string) (TeamSnapshot, error) {
	var t TeamSnapshot
	err := r.pool.QueryRow(ctx, `
		SELECT tenant_id, team_id, org_unit_id, key, name, description, updated_at
		FROM teams_snapshot WHERE tenant_id = $1 AND key = $2
	`, tenantID, key).Scan(&t.TenantID, &t.TeamID, &t.OrgUnitID, &t.Key, &t.Name, &t.Description, &t.UpdatedAt)
	return t, err
}

func (r *TeamsSnapshotRepo) Count(ctx context.Context, tenantID string) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM teams_snapshot WHERE tenant_id = $1`, tenantID).Scan(&count)
	return count, err
}

func (r *TeamsSnapshotRepo) Delete(ctx context.Context, tenantID, teamID string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM teams_snapshot WHERE tenant_id = $1 AND team_id = $2`, tenantID, teamID)
	return err
}

func (r *TeamsSnapshotRepo) DeleteAll(ctx context.Context, tenantID string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM teams_snapshot WHERE tenant_id = $1`, tenantID)
	return err
}
