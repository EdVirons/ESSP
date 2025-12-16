-- +goose Up
ALTER TABLE IF EXISTS inventory
  ADD COLUMN IF NOT EXISTS qty_reserved BIGINT NOT NULL DEFAULT 0;

CREATE TABLE IF NOT EXISTS work_order_parts (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL,
  school_id TEXT NOT NULL,
  work_order_id TEXT NOT NULL,
  service_shop_id TEXT NOT NULL,
  part_id TEXT NOT NULL,
  qty_planned BIGINT NOT NULL DEFAULT 0,
  qty_used BIGINT NOT NULL DEFAULT 0,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_work_order_parts_wo
  ON work_order_parts (tenant_id, school_id, work_order_id, created_at DESC, id DESC);

CREATE INDEX IF NOT EXISTS idx_work_order_parts_shop_part
  ON work_order_parts (tenant_id, service_shop_id, part_id);

-- +goose Down
DROP TABLE IF EXISTS work_order_parts;
ALTER TABLE IF EXISTS inventory DROP COLUMN IF EXISTS qty_reserved;
