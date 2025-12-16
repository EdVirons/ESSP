-- +goose Up

-- Snapshots (cache) from SSOT systems
CREATE TABLE IF NOT EXISTS schools_snapshot (
  tenant_id TEXT NOT NULL,
  school_id TEXT NOT NULL,
  name TEXT NOT NULL DEFAULT '',
  county_code TEXT NOT NULL DEFAULT '',
  county_name TEXT NOT NULL DEFAULT '',
  sub_county_code TEXT NOT NULL DEFAULT '',
  sub_county_name TEXT NOT NULL DEFAULT '',
  updated_at TIMESTAMPTZ NOT NULL,
  PRIMARY KEY (tenant_id, school_id)
);

CREATE INDEX IF NOT EXISTS idx_schools_snapshot_geo
  ON schools_snapshot (tenant_id, county_code, sub_county_code);

CREATE TABLE IF NOT EXISTS devices_snapshot (
  tenant_id TEXT NOT NULL,
  device_id TEXT NOT NULL,
  school_id TEXT NOT NULL,
  model TEXT NOT NULL DEFAULT '',
  serial TEXT NOT NULL DEFAULT '',
  asset_tag TEXT NOT NULL DEFAULT '',
  status TEXT NOT NULL DEFAULT 'active',
  updated_at TIMESTAMPTZ NOT NULL,
  PRIMARY KEY (tenant_id, device_id)
);

CREATE INDEX IF NOT EXISTS idx_devices_snapshot_school
  ON devices_snapshot (tenant_id, school_id);

CREATE TABLE IF NOT EXISTS parts_snapshot (
  tenant_id TEXT NOT NULL,
  part_id TEXT NOT NULL,
  puk TEXT NOT NULL DEFAULT '',
  name TEXT NOT NULL DEFAULT '',
  category TEXT NOT NULL DEFAULT '',
  unit TEXT NOT NULL DEFAULT '',
  updated_at TIMESTAMPTZ NOT NULL,
  PRIMARY KEY (tenant_id, part_id)
);

CREATE INDEX IF NOT EXISTS idx_parts_snapshot_puk
  ON parts_snapshot (tenant_id, puk);

-- Service shops: optional sub-county coverage + coverage level
ALTER TABLE IF EXISTS service_shops
  ADD COLUMN IF NOT EXISTS sub_county_code TEXT NOT NULL DEFAULT '',
  ADD COLUMN IF NOT EXISTS sub_county_name TEXT NOT NULL DEFAULT '',
  ADD COLUMN IF NOT EXISTS coverage_level TEXT NOT NULL DEFAULT 'county'; -- county|sub_county|cluster

CREATE INDEX IF NOT EXISTS idx_service_shops_geo
  ON service_shops (tenant_id, county_code, sub_county_code, coverage_level, active);

-- +goose Down
ALTER TABLE IF EXISTS service_shops
  DROP COLUMN IF EXISTS sub_county_code,
  DROP COLUMN IF EXISTS sub_county_name,
  DROP COLUMN IF EXISTS coverage_level;

DROP TABLE IF EXISTS parts_snapshot;
DROP TABLE IF EXISTS devices_snapshot;
DROP TABLE IF EXISTS schools_snapshot;
