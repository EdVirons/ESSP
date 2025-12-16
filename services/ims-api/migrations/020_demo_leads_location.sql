-- +goose Up
-- Add county and sub-county fields to demo_leads

ALTER TABLE demo_leads
ADD COLUMN IF NOT EXISTS county_code TEXT,
ADD COLUMN IF NOT EXISTS county_name TEXT,
ADD COLUMN IF NOT EXISTS sub_county_code TEXT,
ADD COLUMN IF NOT EXISTS sub_county_name TEXT;

-- Add index for filtering by county
CREATE INDEX IF NOT EXISTS idx_demo_leads_county ON demo_leads(tenant_id, county_code) WHERE county_code IS NOT NULL;

-- +goose Down
DROP INDEX IF EXISTS idx_demo_leads_county;

ALTER TABLE demo_leads
DROP COLUMN IF EXISTS county_code,
DROP COLUMN IF EXISTS county_name,
DROP COLUMN IF EXISTS sub_county_code,
DROP COLUMN IF EXISTS sub_county_name;
