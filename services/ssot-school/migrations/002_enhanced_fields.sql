-- Migration: Add enhanced school fields
-- These fields come from the MOE scrapper database and provide additional school metadata

-- Add KNEC code (Kenya National Examinations Council identifier)
ALTER TABLE schools ADD COLUMN IF NOT EXISTS knec_code TEXT NOT NULL DEFAULT '';

-- Add UIC (Unique Institution Code)
ALTER TABLE schools ADD COLUMN IF NOT EXISTS uic TEXT NOT NULL DEFAULT '';

-- Add sex (boys, girls, mixed)
ALTER TABLE schools ADD COLUMN IF NOT EXISTS sex TEXT NOT NULL DEFAULT '';

-- Add cluster (school cluster identifier)
ALTER TABLE schools ADD COLUMN IF NOT EXISTS cluster TEXT NOT NULL DEFAULT '';

-- Add accommodation (day, boarding, day_and_boarding)
ALTER TABLE schools ADD COLUMN IF NOT EXISTS accommodation TEXT NOT NULL DEFAULT '';

-- Add geographic coordinates
ALTER TABLE schools ADD COLUMN IF NOT EXISTS latitude DOUBLE PRECISION NOT NULL DEFAULT 0.0;
ALTER TABLE schools ADD COLUMN IF NOT EXISTS longitude DOUBLE PRECISION NOT NULL DEFAULT 0.0;

-- Create index for KNEC code lookups
CREATE INDEX IF NOT EXISTS idx_schools_knec ON schools(tenant_id, knec_code) WHERE knec_code != '';

-- Create index for UIC lookups
CREATE INDEX IF NOT EXISTS idx_schools_uic ON schools(tenant_id, uic) WHERE uic != '';

-- Create index for sex filter
CREATE INDEX IF NOT EXISTS idx_schools_sex ON schools(tenant_id, sex) WHERE sex != '';

-- Create index for level filter (optimize common queries)
CREATE INDEX IF NOT EXISTS idx_schools_level ON schools(tenant_id, level);

-- Create index for type filter
CREATE INDEX IF NOT EXISTS idx_schools_type ON schools(tenant_id, type);
