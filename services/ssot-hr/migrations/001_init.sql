-- ssot-hr init migration
-- Uses TEXT primary keys with prefixed ULIDs (e.g., org_01ARZ..., person_01ARZ...)
-- All tables include tenant_id for multi-tenancy
-- No foreign keys (SSOT isolation pattern)

-- org_units: organizational hierarchy (departments, divisions, teams, etc.)
CREATE TABLE IF NOT EXISTS org_units (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL,
  parent_id TEXT NOT NULL DEFAULT '',
  code TEXT NOT NULL,
  name TEXT NOT NULL,
  kind TEXT NOT NULL DEFAULT 'department',
  spec_json TEXT NOT NULL DEFAULT '{}',
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS ux_org_units_code ON org_units(tenant_id, code);
CREATE INDEX IF NOT EXISTS idx_org_units_tenant ON org_units(tenant_id, kind);
CREATE INDEX IF NOT EXISTS idx_org_units_parent ON org_units(tenant_id, parent_id);

-- people: person/employee records
CREATE TABLE IF NOT EXISTS people (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL,
  org_unit_id TEXT NOT NULL DEFAULT '',
  status TEXT NOT NULL DEFAULT 'active',
  given_name TEXT NOT NULL,
  family_name TEXT NOT NULL,
  email TEXT NOT NULL,
  phone TEXT NOT NULL DEFAULT '',
  title TEXT NOT NULL DEFAULT '',
  avatar_url TEXT NOT NULL DEFAULT '',
  spec_json TEXT NOT NULL DEFAULT '{}',
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS ux_people_email ON people(tenant_id, email);
CREATE INDEX IF NOT EXISTS idx_people_tenant ON people(tenant_id, status);
CREATE INDEX IF NOT EXISTS idx_people_org ON people(tenant_id, org_unit_id);
CREATE INDEX IF NOT EXISTS idx_people_name ON people(tenant_id, family_name, given_name);

-- teams: team/workgroup records
CREATE TABLE IF NOT EXISTS teams (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL,
  org_unit_id TEXT NOT NULL DEFAULT '',
  key TEXT NOT NULL,
  name TEXT NOT NULL,
  description TEXT NOT NULL DEFAULT '',
  spec_json TEXT NOT NULL DEFAULT '{}',
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS ux_teams_key ON teams(tenant_id, key);
CREATE INDEX IF NOT EXISTS idx_teams_tenant ON teams(tenant_id);
CREATE INDEX IF NOT EXISTS idx_teams_org ON teams(tenant_id, org_unit_id);

-- team_memberships: many-to-many people to teams
CREATE TABLE IF NOT EXISTS team_memberships (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL,
  team_id TEXT NOT NULL,
  person_id TEXT NOT NULL,
  role TEXT NOT NULL DEFAULT 'member',
  status TEXT NOT NULL DEFAULT 'active',
  started_at TIMESTAMPTZ,
  ended_at TIMESTAMPTZ,
  spec_json TEXT NOT NULL DEFAULT '{}',
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS ux_team_memberships ON team_memberships(tenant_id, team_id, person_id) WHERE status = 'active';
CREATE INDEX IF NOT EXISTS idx_team_memberships_team ON team_memberships(tenant_id, team_id, status);
CREATE INDEX IF NOT EXISTS idx_team_memberships_person ON team_memberships(tenant_id, person_id, status);

-- hr_audit_log: immutable audit trail for compliance
CREATE TABLE IF NOT EXISTS hr_audit_log (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL,
  actor_person_id TEXT,
  action TEXT NOT NULL,
  entity_type TEXT NOT NULL,
  entity_id TEXT NOT NULL,
  request_id TEXT,
  ip_address TEXT,
  user_agent TEXT,
  before_json JSONB,
  after_json JSONB,
  diff_json JSONB,
  created_at TIMESTAMPTZ NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_hr_audit_entity ON hr_audit_log(tenant_id, entity_type, entity_id);
CREATE INDEX IF NOT EXISTS idx_hr_audit_actor ON hr_audit_log(tenant_id, actor_person_id);

-- hr_outbox_events: transactional outbox for event delivery
CREATE TABLE IF NOT EXISTS hr_outbox_events (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL,
  topic TEXT NOT NULL,
  payload JSONB NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  processed_at TIMESTAMPTZ,
  failed_at TIMESTAMPTZ,
  fail_count INT NOT NULL DEFAULT 0
);
CREATE INDEX IF NOT EXISTS idx_hr_outbox_unprocessed ON hr_outbox_events(processed_at) WHERE processed_at IS NULL;
