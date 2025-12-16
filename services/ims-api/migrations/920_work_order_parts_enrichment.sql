-- +goose Up
ALTER TABLE IF EXISTS work_order_parts
  ADD COLUMN IF NOT EXISTS part_name TEXT NOT NULL DEFAULT '',
  ADD COLUMN IF NOT EXISTS part_puk TEXT NOT NULL DEFAULT '',
  ADD COLUMN IF NOT EXISTS part_category TEXT NOT NULL DEFAULT '',
  ADD COLUMN IF NOT EXISTS device_model_id TEXT NOT NULL DEFAULT '',
  ADD COLUMN IF NOT EXISTS is_compatible BOOLEAN NOT NULL DEFAULT TRUE;

CREATE INDEX IF NOT EXISTS idx_work_order_parts_compat
  ON work_order_parts (tenant_id, work_order_id, device_model_id, is_compatible);

-- +goose Down
DROP INDEX IF EXISTS idx_work_order_parts_compat;

ALTER TABLE IF EXISTS work_order_parts
  DROP COLUMN IF EXISTS part_name,
  DROP COLUMN IF EXISTS part_puk,
  DROP COLUMN IF EXISTS part_category,
  DROP COLUMN IF EXISTS device_model_id,
  DROP COLUMN IF EXISTS is_compatible;
