package store

import (
	"context"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TeamMembershipSnapshot struct {
	TenantID     string     `json:"tenantId"`
	MembershipID string     `json:"membershipId"`
	TeamID       string     `json:"teamId"`
	PersonID     string     `json:"personId"`
	Role         string     `json:"role"`
	Status       string     `json:"status"`
	StartedAt    *time.Time `json:"startedAt,omitempty"`
	EndedAt      *time.Time `json:"endedAt,omitempty"`
	UpdatedAt    time.Time  `json:"updatedAt"`
}

type TeamMembershipsSnapshotRepo struct {
	pool *pgxpool.Pool
}

type TeamMembershipSnapshotListParams struct {
	TenantID string
	TeamID   string
	PersonID string
	Role     string
	Status   string
	Limit    int
	Offset   int
}

func (r *TeamMembershipsSnapshotRepo) Upsert(ctx context.Context, m TeamMembershipSnapshot) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO team_memberships_snapshot (tenant_id, membership_id, team_id, person_id, role, status, started_at, ended_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (tenant_id, membership_id) DO UPDATE SET
			team_id = EXCLUDED.team_id,
			person_id = EXCLUDED.person_id,
			role = EXCLUDED.role,
			status = EXCLUDED.status,
			started_at = EXCLUDED.started_at,
			ended_at = EXCLUDED.ended_at,
			updated_at = EXCLUDED.updated_at
	`, m.TenantID, m.MembershipID, m.TeamID, m.PersonID, m.Role, m.Status, m.StartedAt, m.EndedAt, m.UpdatedAt)
	return err
}

func (r *TeamMembershipsSnapshotRepo) UpsertBatch(ctx context.Context, memberships []TeamMembershipSnapshot) (int, error) {
	if len(memberships) == 0 {
		return 0, nil
	}

	batch := &pgx.Batch{}
	for _, m := range memberships {
		batch.Queue(`
			INSERT INTO team_memberships_snapshot (tenant_id, membership_id, team_id, person_id, role, status, started_at, ended_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			ON CONFLICT (tenant_id, membership_id) DO UPDATE SET
				team_id = EXCLUDED.team_id,
				person_id = EXCLUDED.person_id,
				role = EXCLUDED.role,
				status = EXCLUDED.status,
				started_at = EXCLUDED.started_at,
				ended_at = EXCLUDED.ended_at,
				updated_at = EXCLUDED.updated_at
		`, m.TenantID, m.MembershipID, m.TeamID, m.PersonID, m.Role, m.Status, m.StartedAt, m.EndedAt, m.UpdatedAt)
	}

	br := r.pool.SendBatch(ctx, batch)
	defer br.Close()

	for range memberships {
		if _, err := br.Exec(); err != nil {
			return 0, err
		}
	}
	return len(memberships), nil
}

func (r *TeamMembershipsSnapshotRepo) List(ctx context.Context, params TeamMembershipSnapshotListParams) ([]TeamMembershipSnapshot, int, error) {
	if params.Limit == 0 {
		params.Limit = 50
	}

	countQuery := `SELECT COUNT(*) FROM team_memberships_snapshot WHERE tenant_id = $1`
	args := []any{params.TenantID}
	argIdx := 2

	if params.TeamID != "" {
		countQuery += ` AND team_id = $` + strconv.Itoa(argIdx)
		args = append(args, params.TeamID)
		argIdx++
	}
	if params.PersonID != "" {
		countQuery += ` AND person_id = $` + strconv.Itoa(argIdx)
		args = append(args, params.PersonID)
		argIdx++
	}
	if params.Role != "" {
		countQuery += ` AND role = $` + strconv.Itoa(argIdx)
		args = append(args, params.Role)
		argIdx++
	}
	if params.Status != "" {
		countQuery += ` AND status = $` + strconv.Itoa(argIdx)
		args = append(args, params.Status)
		argIdx++
	}

	var total int
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	listQuery := `SELECT tenant_id, membership_id, team_id, person_id, role, status, started_at, ended_at, updated_at
		FROM team_memberships_snapshot WHERE tenant_id = $1`
	args = []any{params.TenantID}
	argIdx = 2

	if params.TeamID != "" {
		listQuery += ` AND team_id = $` + strconv.Itoa(argIdx)
		args = append(args, params.TeamID)
		argIdx++
	}
	if params.PersonID != "" {
		listQuery += ` AND person_id = $` + strconv.Itoa(argIdx)
		args = append(args, params.PersonID)
		argIdx++
	}
	if params.Role != "" {
		listQuery += ` AND role = $` + strconv.Itoa(argIdx)
		args = append(args, params.Role)
		argIdx++
	}
	if params.Status != "" {
		listQuery += ` AND status = $` + strconv.Itoa(argIdx)
		args = append(args, params.Status)
		argIdx++
	}

	listQuery += ` ORDER BY updated_at DESC LIMIT $` + strconv.Itoa(argIdx) + ` OFFSET $` + strconv.Itoa(argIdx+1)
	args = append(args, params.Limit, params.Offset)

	rows, err := r.pool.Query(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var items []TeamMembershipSnapshot
	for rows.Next() {
		var m TeamMembershipSnapshot
		if err := rows.Scan(&m.TenantID, &m.MembershipID, &m.TeamID, &m.PersonID, &m.Role, &m.Status, &m.StartedAt, &m.EndedAt, &m.UpdatedAt); err != nil {
			return nil, 0, err
		}
		items = append(items, m)
	}
	return items, total, rows.Err()
}

func (r *TeamMembershipsSnapshotRepo) Get(ctx context.Context, tenantID, membershipID string) (TeamMembershipSnapshot, error) {
	var m TeamMembershipSnapshot
	err := r.pool.QueryRow(ctx, `
		SELECT tenant_id, membership_id, team_id, person_id, role, status, started_at, ended_at, updated_at
		FROM team_memberships_snapshot WHERE tenant_id = $1 AND membership_id = $2
	`, tenantID, membershipID).Scan(&m.TenantID, &m.MembershipID, &m.TeamID, &m.PersonID, &m.Role, &m.Status, &m.StartedAt, &m.EndedAt, &m.UpdatedAt)
	return m, err
}

func (r *TeamMembershipsSnapshotRepo) ListByTeam(ctx context.Context, tenantID, teamID string) ([]TeamMembershipSnapshot, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT tenant_id, membership_id, team_id, person_id, role, status, started_at, ended_at, updated_at
		FROM team_memberships_snapshot WHERE tenant_id = $1 AND team_id = $2 ORDER BY role, updated_at DESC
	`, tenantID, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []TeamMembershipSnapshot
	for rows.Next() {
		var m TeamMembershipSnapshot
		if err := rows.Scan(&m.TenantID, &m.MembershipID, &m.TeamID, &m.PersonID, &m.Role, &m.Status, &m.StartedAt, &m.EndedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, m)
	}
	return items, rows.Err()
}

func (r *TeamMembershipsSnapshotRepo) ListByPerson(ctx context.Context, tenantID, personID string) ([]TeamMembershipSnapshot, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT tenant_id, membership_id, team_id, person_id, role, status, started_at, ended_at, updated_at
		FROM team_memberships_snapshot WHERE tenant_id = $1 AND person_id = $2 ORDER BY updated_at DESC
	`, tenantID, personID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []TeamMembershipSnapshot
	for rows.Next() {
		var m TeamMembershipSnapshot
		if err := rows.Scan(&m.TenantID, &m.MembershipID, &m.TeamID, &m.PersonID, &m.Role, &m.Status, &m.StartedAt, &m.EndedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, m)
	}
	return items, rows.Err()
}

func (r *TeamMembershipsSnapshotRepo) Count(ctx context.Context, tenantID string) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM team_memberships_snapshot WHERE tenant_id = $1`, tenantID).Scan(&count)
	return count, err
}

func (r *TeamMembershipsSnapshotRepo) CountByTeam(ctx context.Context, tenantID, teamID string) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM team_memberships_snapshot WHERE tenant_id = $1 AND team_id = $2`, tenantID, teamID).Scan(&count)
	return count, err
}

func (r *TeamMembershipsSnapshotRepo) Delete(ctx context.Context, tenantID, membershipID string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM team_memberships_snapshot WHERE tenant_id = $1 AND membership_id = $2`, tenantID, membershipID)
	return err
}

func (r *TeamMembershipsSnapshotRepo) DeleteByTeam(ctx context.Context, tenantID, teamID string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM team_memberships_snapshot WHERE tenant_id = $1 AND team_id = $2`, tenantID, teamID)
	return err
}

func (r *TeamMembershipsSnapshotRepo) DeleteAll(ctx context.Context, tenantID string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM team_memberships_snapshot WHERE tenant_id = $1`, tenantID)
	return err
}
