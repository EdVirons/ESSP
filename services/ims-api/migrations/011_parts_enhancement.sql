-- Parts Enhancement Migration
-- Add new fields for price, supplier, and status tracking

-- Add new columns to parts table
ALTER TABLE parts
  ADD COLUMN IF NOT EXISTS unit_cost_cents INTEGER DEFAULT 0,
  ADD COLUMN IF NOT EXISTS supplier TEXT DEFAULT '',
  ADD COLUMN IF NOT EXISTS supplier_sku TEXT DEFAULT '',
  ADD COLUMN IF NOT EXISTS description TEXT DEFAULT '',
  ADD COLUMN IF NOT EXISTS active BOOLEAN DEFAULT true,
  ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ DEFAULT NOW();

-- Add indexes for filtering
CREATE INDEX IF NOT EXISTS idx_parts_category ON parts (tenant_id, category, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_parts_active ON parts (tenant_id, active, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_parts_supplier ON parts (tenant_id, supplier) WHERE supplier != '';
