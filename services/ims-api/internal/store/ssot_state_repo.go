package store

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SSOTResource string

const (
	SSOTSchools  SSOTResource = "schools"
	SSOTDevices  SSOTResource = "devices"
	SSOTParts    SSOTResource = "parts"
	SSOTPeople   SSOTResource = "people"
	SSOTTeams    SSOTResource = "teams"
	SSOTOrgUnits SSOTResource = "org-units"
	SSOTTeamMemberships SSOTResource = "team-memberships"
)

type SSOTSyncState struct {
	TenantID          string
	Resource          SSOTResource
	LastUpdatedSince  time.Time
	LastCursor        string
	UpdatedAt         time.Time
}

type SSOTStateRepo struct{ pool *pgxpool.Pool }

func (r *SSOTStateRepo) Get(ctx context.Context, tenantID string, res SSOTResource) (SSOTSyncState, error) {
	var s SSOTSyncState
	row := r.pool.QueryRow(ctx, `
		SELECT tenant_id, resource, last_updated_since, last_cursor, updated_at
		FROM ssot_sync_state
		WHERE tenant_id=$1 AND resource=$2
	`, tenantID, string(res))
	if err := row.Scan(&s.TenantID, &s.Resource, &s.LastUpdatedSince, &s.LastCursor, &s.UpdatedAt); err != nil {
		return SSOTSyncState{}, errors.New("not found")
	}
	return s, nil
}

func (r *SSOTStateRepo) Upsert(ctx context.Context, s SSOTSyncState) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO ssot_sync_state (tenant_id, resource, last_updated_since, last_cursor, updated_at)
		VALUES ($1,$2,$3,$4,$5)
		ON CONFLICT (tenant_id, resource)
		DO UPDATE SET last_updated_since=EXCLUDED.last_updated_since, last_cursor=EXCLUDED.last_cursor, updated_at=EXCLUDED.updated_at
	`, s.TenantID, string(s.Resource), s.LastUpdatedSince, s.LastCursor, s.UpdatedAt)
	return err
}

func NewSSOTSyncState(tenantID string, res SSOTResource) SSOTSyncState {
	return SSOTSyncState{
		TenantID: tenantID,
		Resource: res,
		LastUpdatedSince: time.Unix(0, 0).UTC(),
		LastCursor: "",
		UpdatedAt: time.Now().UTC(),
	}
}
