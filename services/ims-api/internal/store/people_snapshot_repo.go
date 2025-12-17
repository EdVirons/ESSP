package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PersonSnapshot struct {
	TenantID   string    `json:"tenantId"`
	PersonID   string    `json:"personId"`
	OrgUnitID  string    `json:"orgUnitId,omitempty"`
	Status     string    `json:"status"`
	GivenName  string    `json:"givenName"`
	FamilyName string    `json:"familyName"`
	FullName   string    `json:"fullName"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone,omitempty"`
	Title      string    `json:"title,omitempty"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type PeopleSnapshotRepo struct {
	pool *pgxpool.Pool
}

type PersonSnapshotListParams struct {
	TenantID  string
	Query     string
	Status    string
	OrgUnitID string
	Limit     int
	Offset    int
}

func (r *PeopleSnapshotRepo) Upsert(ctx context.Context, p PersonSnapshot) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO people_snapshot (tenant_id, person_id, org_unit_id, status, given_name, family_name, full_name, email, phone, title, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (tenant_id, person_id) DO UPDATE SET
			org_unit_id = EXCLUDED.org_unit_id,
			status = EXCLUDED.status,
			given_name = EXCLUDED.given_name,
			family_name = EXCLUDED.family_name,
			full_name = EXCLUDED.full_name,
			email = EXCLUDED.email,
			phone = EXCLUDED.phone,
			title = EXCLUDED.title,
			updated_at = EXCLUDED.updated_at
	`, p.TenantID, p.PersonID, p.OrgUnitID, p.Status, p.GivenName, p.FamilyName, p.FullName, p.Email, p.Phone, p.Title, p.UpdatedAt)
	return err
}

func (r *PeopleSnapshotRepo) UpsertBatch(ctx context.Context, people []PersonSnapshot) (int, error) {
	if len(people) == 0 {
		return 0, nil
	}

	batch := &pgx.Batch{}
	for _, p := range people {
		batch.Queue(`
			INSERT INTO people_snapshot (tenant_id, person_id, org_unit_id, status, given_name, family_name, full_name, email, phone, title, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
			ON CONFLICT (tenant_id, person_id) DO UPDATE SET
				org_unit_id = EXCLUDED.org_unit_id,
				status = EXCLUDED.status,
				given_name = EXCLUDED.given_name,
				family_name = EXCLUDED.family_name,
				full_name = EXCLUDED.full_name,
				email = EXCLUDED.email,
				phone = EXCLUDED.phone,
				title = EXCLUDED.title,
				updated_at = EXCLUDED.updated_at
		`, p.TenantID, p.PersonID, p.OrgUnitID, p.Status, p.GivenName, p.FamilyName, p.FullName, p.Email, p.Phone, p.Title, p.UpdatedAt)
	}

	br := r.pool.SendBatch(ctx, batch)
	defer br.Close()

	for range people {
		if _, err := br.Exec(); err != nil {
			return 0, err
		}
	}
	return len(people), nil
}

func (r *PeopleSnapshotRepo) List(ctx context.Context, params PersonSnapshotListParams) ([]PersonSnapshot, int, error) {
	if params.Limit == 0 {
		params.Limit = 50
	}

	countQuery := `SELECT COUNT(*) FROM people_snapshot WHERE tenant_id = $1`
	args := []any{params.TenantID}
	argIdx := 2

	if params.Query != "" {
		countQuery += ` AND (full_name ILIKE $` + itoa(argIdx) + ` OR email ILIKE $` + itoa(argIdx) + `)`
		args = append(args, "%"+params.Query+"%")
		argIdx++
	}
	if params.Status != "" {
		countQuery += ` AND status = $` + itoa(argIdx)
		args = append(args, params.Status)
		argIdx++
	}
	if params.OrgUnitID != "" {
		countQuery += ` AND org_unit_id = $` + itoa(argIdx)
		args = append(args, params.OrgUnitID)
		argIdx++
	}

	var total int
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	listQuery := `SELECT tenant_id, person_id, org_unit_id, status, given_name, family_name, full_name, email, phone, title, updated_at
		FROM people_snapshot WHERE tenant_id = $1`
	args = []any{params.TenantID}
	argIdx = 2

	if params.Query != "" {
		listQuery += ` AND (full_name ILIKE $` + itoa(argIdx) + ` OR email ILIKE $` + itoa(argIdx) + `)`
		args = append(args, "%"+params.Query+"%")
		argIdx++
	}
	if params.Status != "" {
		listQuery += ` AND status = $` + itoa(argIdx)
		args = append(args, params.Status)
		argIdx++
	}
	if params.OrgUnitID != "" {
		listQuery += ` AND org_unit_id = $` + itoa(argIdx)
		args = append(args, params.OrgUnitID)
		argIdx++
	}

	listQuery += ` ORDER BY full_name ASC LIMIT $` + itoa(argIdx) + ` OFFSET $` + itoa(argIdx+1)
	args = append(args, params.Limit, params.Offset)

	rows, err := r.pool.Query(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var items []PersonSnapshot
	for rows.Next() {
		var p PersonSnapshot
		if err := rows.Scan(&p.TenantID, &p.PersonID, &p.OrgUnitID, &p.Status, &p.GivenName, &p.FamilyName, &p.FullName, &p.Email, &p.Phone, &p.Title, &p.UpdatedAt); err != nil {
			return nil, 0, err
		}
		items = append(items, p)
	}
	return items, total, rows.Err()
}

func (r *PeopleSnapshotRepo) Get(ctx context.Context, tenantID, personID string) (PersonSnapshot, error) {
	var p PersonSnapshot
	err := r.pool.QueryRow(ctx, `
		SELECT tenant_id, person_id, org_unit_id, status, given_name, family_name, full_name, email, phone, title, updated_at
		FROM people_snapshot WHERE tenant_id = $1 AND person_id = $2
	`, tenantID, personID).Scan(&p.TenantID, &p.PersonID, &p.OrgUnitID, &p.Status, &p.GivenName, &p.FamilyName, &p.FullName, &p.Email, &p.Phone, &p.Title, &p.UpdatedAt)
	return p, err
}

func (r *PeopleSnapshotRepo) GetByEmail(ctx context.Context, tenantID, email string) (PersonSnapshot, error) {
	var p PersonSnapshot
	err := r.pool.QueryRow(ctx, `
		SELECT tenant_id, person_id, org_unit_id, status, given_name, family_name, full_name, email, phone, title, updated_at
		FROM people_snapshot WHERE tenant_id = $1 AND email = $2
	`, tenantID, email).Scan(&p.TenantID, &p.PersonID, &p.OrgUnitID, &p.Status, &p.GivenName, &p.FamilyName, &p.FullName, &p.Email, &p.Phone, &p.Title, &p.UpdatedAt)
	return p, err
}

func (r *PeopleSnapshotRepo) Count(ctx context.Context, tenantID string) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM people_snapshot WHERE tenant_id = $1`, tenantID).Scan(&count)
	return count, err
}

func (r *PeopleSnapshotRepo) Delete(ctx context.Context, tenantID, personID string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM people_snapshot WHERE tenant_id = $1 AND person_id = $2`, tenantID, personID)
	return err
}

func (r *PeopleSnapshotRepo) DeleteAll(ctx context.Context, tenantID string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM people_snapshot WHERE tenant_id = $1`, tenantID)
	return err
}
