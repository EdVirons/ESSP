-- +goose Up
CREATE TABLE IF NOT EXISTS incidents (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL,
  school_id TEXT NOT NULL,
  device_id TEXT NOT NULL,
  category TEXT NOT NULL,
  severity TEXT NOT NULL,
  status TEXT NOT NULL,
  service_shop_id TEXT NOT NULL DEFAULT '',
  assigned_staff_id TEXT NOT NULL DEFAULT '',
  repair_location TEXT NOT NULL DEFAULT 'service_shop',
  title TEXT NOT NULL,
  description TEXT NOT NULL,
  reported_by TEXT NOT NULL,
  sla_due_at TIMESTAMPTZ NOT NULL,
  sla_breached BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_incidents_tenant_school_created
  ON incidents (tenant_id, school_id, created_at DESC, id DESC);

CREATE INDEX IF NOT EXISTS idx_incidents_sla_due
  ON incidents (sla_breached, sla_due_at);

CREATE TABLE IF NOT EXISTS work_orders (
  id TEXT PRIMARY KEY,
  incident_id TEXT NOT NULL,
  tenant_id TEXT NOT NULL,
  school_id TEXT NOT NULL,
  device_id TEXT NOT NULL,
  status TEXT NOT NULL,
  service_shop_id TEXT NOT NULL DEFAULT '',
  assigned_staff_id TEXT NOT NULL DEFAULT '',
  repair_location TEXT NOT NULL DEFAULT 'service_shop',
  assigned_to TEXT NOT NULL,
  task_type TEXT NOT NULL,
  cost_estimate_cents BIGINT NOT NULL DEFAULT 0,
  notes TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_work_orders_tenant_school_created
  ON work_orders (tenant_id, school_id, created_at DESC, id DESC);

CREATE TABLE IF NOT EXISTS attachments (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL,
  school_id TEXT NOT NULL,
  entity_type TEXT NOT NULL,
  entity_id TEXT NOT NULL,
  file_name TEXT NOT NULL,
  content_type TEXT NOT NULL,
  size_bytes BIGINT NOT NULL DEFAULT 0,
  object_key TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_attachments_entity
  ON attachments (tenant_id, school_id, entity_type, entity_id, created_at DESC, id DESC);


CREATE TABLE IF NOT EXISTS schools (
  tenant_id TEXT NOT NULL,
  school_id TEXT NOT NULL,
  county_code TEXT NOT NULL,
  county_name TEXT NOT NULL DEFAULT '',
  updated_at TIMESTAMPTZ NOT NULL,
  PRIMARY KEY (tenant_id, school_id)
);

CREATE TABLE IF NOT EXISTS service_shops (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL,
  county_code TEXT NOT NULL,
  county_name TEXT NOT NULL DEFAULT '',
  name TEXT NOT NULL,
  location TEXT NOT NULL DEFAULT '',
  active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  UNIQUE (tenant_id, county_code)
);

CREATE INDEX IF NOT EXISTS idx_service_shops_tenant_created
  ON service_shops (tenant_id, created_at DESC, id DESC);

CREATE TABLE IF NOT EXISTS service_staff (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL,
  service_shop_id TEXT NOT NULL,
  user_id TEXT NOT NULL,
  role TEXT NOT NULL,
  phone TEXT NOT NULL DEFAULT '',
  active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_service_staff_shop_role
  ON service_staff (tenant_id, service_shop_id, role, active, created_at);

CREATE TABLE IF NOT EXISTS parts (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL,
  sku TEXT NOT NULL,
  name TEXT NOT NULL,
  category TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL,
  UNIQUE (tenant_id, sku)
);

CREATE INDEX IF NOT EXISTS idx_parts_tenant_created
  ON parts (tenant_id, created_at DESC, id DESC);

CREATE TABLE IF NOT EXISTS inventory (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL,
  service_shop_id TEXT NOT NULL,
  part_id TEXT NOT NULL,
  qty_available BIGINT NOT NULL DEFAULT 0,
  reorder_threshold BIGINT NOT NULL DEFAULT 0,
  updated_at TIMESTAMPTZ NOT NULL,
  UNIQUE (tenant_id, service_shop_id, part_id)
);

CREATE INDEX IF NOT EXISTS idx_inventory_shop
  ON inventory (tenant_id, service_shop_id, updated_at DESC, id DESC);

-- +goose Down
DROP TABLE IF EXISTS inventory;
DROP TABLE IF EXISTS parts;
DROP TABLE IF EXISTS service_staff;
DROP TABLE IF EXISTS service_shops;
DROP TABLE IF EXISTS schools;
DROP TABLE IF EXISTS attachments;
DROP TABLE IF EXISTS work_orders;
DROP TABLE IF EXISTS incidents;
