-- ssot-parts init
CREATE TABLE IF NOT EXISTS parts (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL,
  name TEXT NOT NULL,
  category TEXT NOT NULL DEFAULT 'misc',
  puk TEXT NOT NULL DEFAULT '',
  spec_json TEXT NOT NULL DEFAULT '{}',
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS ux_parts_puk ON parts(tenant_id, puk) WHERE puk <> '';

CREATE TABLE IF NOT EXISTS part_compatibility (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL,
  part_id TEXT NOT NULL,
  device_model_id TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_part_compat_part ON part_compatibility(tenant_id, part_id);

CREATE TABLE IF NOT EXISTS vendor_skus (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL,
  part_id TEXT NOT NULL,
  vendor_id TEXT NOT NULL,
  sku TEXT NOT NULL,
  unit_price_cents BIGINT NOT NULL DEFAULT 0,
  currency TEXT NOT NULL DEFAULT 'KES',
  lead_time_days INT NOT NULL DEFAULT 7,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_vendor_skus_part ON vendor_skus(tenant_id, part_id);
