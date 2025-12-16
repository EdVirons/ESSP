-- +goose Up

-- Add enhanced school fields to the schools_snapshot table
-- These fields mirror the SSOT school data for quick access

ALTER TABLE schools_snapshot ADD COLUMN IF NOT EXISTS level TEXT NOT NULL DEFAULT '';
ALTER TABLE schools_snapshot ADD COLUMN IF NOT EXISTS type TEXT NOT NULL DEFAULT '';
ALTER TABLE schools_snapshot ADD COLUMN IF NOT EXISTS knec_code TEXT NOT NULL DEFAULT '';
ALTER TABLE schools_snapshot ADD COLUMN IF NOT EXISTS uic TEXT NOT NULL DEFAULT '';
ALTER TABLE schools_snapshot ADD COLUMN IF NOT EXISTS sex TEXT NOT NULL DEFAULT '';
ALTER TABLE schools_snapshot ADD COLUMN IF NOT EXISTS cluster TEXT NOT NULL DEFAULT '';
ALTER TABLE schools_snapshot ADD COLUMN IF NOT EXISTS accommodation TEXT NOT NULL DEFAULT '';
ALTER TABLE schools_snapshot ADD COLUMN IF NOT EXISTS latitude DOUBLE PRECISION NOT NULL DEFAULT 0.0;
ALTER TABLE schools_snapshot ADD COLUMN IF NOT EXISTS longitude DOUBLE PRECISION NOT NULL DEFAULT 0.0;

-- Create indexes for common filters
CREATE INDEX IF NOT EXISTS idx_schools_snapshot_level ON schools_snapshot(tenant_id, level) WHERE level != '';
CREATE INDEX IF NOT EXISTS idx_schools_snapshot_type ON schools_snapshot(tenant_id, type) WHERE type != '';
CREATE INDEX IF NOT EXISTS idx_schools_snapshot_knec ON schools_snapshot(tenant_id, knec_code) WHERE knec_code != '';

-- +goose Down
ALTER TABLE schools_snapshot DROP COLUMN IF EXISTS level;
ALTER TABLE schools_snapshot DROP COLUMN IF EXISTS type;
ALTER TABLE schools_snapshot DROP COLUMN IF EXISTS knec_code;
ALTER TABLE schools_snapshot DROP COLUMN IF EXISTS uic;
ALTER TABLE schools_snapshot DROP COLUMN IF EXISTS sex;
ALTER TABLE schools_snapshot DROP COLUMN IF EXISTS cluster;
ALTER TABLE schools_snapshot DROP COLUMN IF EXISTS accommodation;
ALTER TABLE schools_snapshot DROP COLUMN IF EXISTS latitude;
ALTER TABLE schools_snapshot DROP COLUMN IF EXISTS longitude;

DROP INDEX IF EXISTS idx_schools_snapshot_level;
DROP INDEX IF EXISTS idx_schools_snapshot_type;
DROP INDEX IF EXISTS idx_schools_snapshot_knec;
