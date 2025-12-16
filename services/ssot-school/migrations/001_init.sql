-- ssot-school init
CREATE TABLE IF NOT EXISTS counties (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL,
  name TEXT NOT NULL,
  code TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS ux_counties_code ON counties(tenant_id, code);

CREATE TABLE IF NOT EXISTS sub_counties (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL,
  county_id TEXT NOT NULL,
  name TEXT NOT NULL,
  code TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_sub_counties_county ON sub_counties(tenant_id, county_id);

CREATE TABLE IF NOT EXISTS schools (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL,
  name TEXT NOT NULL,
  code TEXT NOT NULL DEFAULT '',
  county_id TEXT NOT NULL,
  sub_county_id TEXT NOT NULL,
  level TEXT NOT NULL DEFAULT 'Other',
  type TEXT NOT NULL DEFAULT 'public',
  active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_schools_region ON schools(tenant_id, county_id, sub_county_id);

CREATE TABLE IF NOT EXISTS school_contacts (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL,
  school_id TEXT NOT NULL,
  name TEXT NOT NULL,
  phone TEXT NOT NULL DEFAULT '',
  email TEXT NOT NULL DEFAULT '',
  role TEXT NOT NULL DEFAULT 'point_of_contact',
  is_primary BOOLEAN NOT NULL DEFAULT FALSE,
  active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_contacts_school ON school_contacts(tenant_id, school_id, is_primary DESC, active DESC);
