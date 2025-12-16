-- ssot-devices init
CREATE TABLE IF NOT EXISTS device_models (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL,
  make TEXT NOT NULL,
  model TEXT NOT NULL,
  category TEXT NOT NULL DEFAULT 'other',
  spec_json TEXT NOT NULL DEFAULT '{}',
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_device_models ON device_models(tenant_id, make, model);

CREATE TABLE IF NOT EXISTS devices (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL,
  serial TEXT NOT NULL DEFAULT '',
  asset_tag TEXT NOT NULL DEFAULT '',
  device_model_id TEXT NOT NULL,
  school_id TEXT NOT NULL DEFAULT '',
  assigned_to TEXT NOT NULL DEFAULT '',
  lifecycle TEXT NOT NULL DEFAULT 'in_stock',
  enrolled BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_devices_school ON devices(tenant_id, school_id, lifecycle);
CREATE UNIQUE INDEX IF NOT EXISTS ux_devices_serial ON devices(tenant_id, serial) WHERE serial <> '';
