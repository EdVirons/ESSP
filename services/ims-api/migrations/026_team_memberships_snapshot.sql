-- +goose Up
-- Team Memberships Snapshot Table for IMS-API
-- Read-optimized cache of team membership data from the ssot-hr service

CREATE TABLE IF NOT EXISTS team_memberships_snapshot (
  tenant_id TEXT NOT NULL,
  membership_id TEXT NOT NULL,
  team_id TEXT NOT NULL,
  person_id TEXT NOT NULL,
  role TEXT NOT NULL DEFAULT 'member',
  status TEXT NOT NULL DEFAULT 'active',
  started_at TIMESTAMPTZ,
  ended_at TIMESTAMPTZ,
  updated_at TIMESTAMPTZ NOT NULL,
  PRIMARY KEY (tenant_id, membership_id)
);

-- Index for finding all memberships of a team
CREATE INDEX IF NOT EXISTS idx_team_memberships_snapshot_team ON team_memberships_snapshot(tenant_id, team_id);

-- Index for finding all teams a person belongs to
CREATE INDEX IF NOT EXISTS idx_team_memberships_snapshot_person ON team_memberships_snapshot(tenant_id, person_id);

-- Index for filtering by status
CREATE INDEX IF NOT EXISTS idx_team_memberships_snapshot_status ON team_memberships_snapshot(tenant_id, status);

-- Composite index for team + person uniqueness lookups
CREATE INDEX IF NOT EXISTS idx_team_memberships_snapshot_team_person ON team_memberships_snapshot(tenant_id, team_id, person_id);

-- +goose Down
DROP TABLE IF EXISTS team_memberships_snapshot;
