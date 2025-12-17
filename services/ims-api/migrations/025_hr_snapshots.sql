-- +goose Up
-- HR SSOT Snapshot Tables for IMS-API
-- These are read-optimized caches of HR data from the ssot-hr service

-- people_snapshot: person data from HR SSOT
CREATE TABLE IF NOT EXISTS people_snapshot (
  tenant_id TEXT NOT NULL,
  person_id TEXT NOT NULL,
  org_unit_id TEXT NOT NULL DEFAULT '',
  status TEXT NOT NULL DEFAULT 'active',
  given_name TEXT NOT NULL DEFAULT '',
  family_name TEXT NOT NULL DEFAULT '',
  full_name TEXT NOT NULL DEFAULT '',
  email TEXT NOT NULL DEFAULT '',
  phone TEXT NOT NULL DEFAULT '',
  title TEXT NOT NULL DEFAULT '',
  updated_at TIMESTAMPTZ NOT NULL,
  PRIMARY KEY (tenant_id, person_id)
);
CREATE INDEX IF NOT EXISTS idx_people_snapshot_email ON people_snapshot(tenant_id, email) WHERE email != '';
CREATE INDEX IF NOT EXISTS idx_people_snapshot_org ON people_snapshot(tenant_id, org_unit_id) WHERE org_unit_id != '';
CREATE INDEX IF NOT EXISTS idx_people_snapshot_status ON people_snapshot(tenant_id, status);

-- teams_snapshot: team data from HR SSOT
CREATE TABLE IF NOT EXISTS teams_snapshot (
  tenant_id TEXT NOT NULL,
  team_id TEXT NOT NULL,
  org_unit_id TEXT NOT NULL DEFAULT '',
  key TEXT NOT NULL DEFAULT '',
  name TEXT NOT NULL DEFAULT '',
  description TEXT NOT NULL DEFAULT '',
  updated_at TIMESTAMPTZ NOT NULL,
  PRIMARY KEY (tenant_id, team_id)
);
CREATE INDEX IF NOT EXISTS idx_teams_snapshot_key ON teams_snapshot(tenant_id, key) WHERE key != '';
CREATE INDEX IF NOT EXISTS idx_teams_snapshot_org ON teams_snapshot(tenant_id, org_unit_id) WHERE org_unit_id != '';

-- org_units_snapshot: org unit data from HR SSOT
CREATE TABLE IF NOT EXISTS org_units_snapshot (
  tenant_id TEXT NOT NULL,
  org_unit_id TEXT NOT NULL,
  parent_id TEXT NOT NULL DEFAULT '',
  code TEXT NOT NULL DEFAULT '',
  name TEXT NOT NULL DEFAULT '',
  kind TEXT NOT NULL DEFAULT '',
  updated_at TIMESTAMPTZ NOT NULL,
  PRIMARY KEY (tenant_id, org_unit_id)
);
CREATE INDEX IF NOT EXISTS idx_org_units_snapshot_code ON org_units_snapshot(tenant_id, code) WHERE code != '';
CREATE INDEX IF NOT EXISTS idx_org_units_snapshot_parent ON org_units_snapshot(tenant_id, parent_id) WHERE parent_id != '';

-- +goose Down
DROP TABLE IF EXISTS org_units_snapshot;
DROP TABLE IF EXISTS teams_snapshot;
DROP TABLE IF EXISTS people_snapshot;
