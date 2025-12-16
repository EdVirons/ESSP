-- +goose Up
CREATE TABLE IF NOT EXISTS ssot_sync_state (
  tenant_id TEXT NOT NULL,
  resource TEXT NOT NULL, -- schools|devices|parts
  last_updated_since TIMESTAMPTZ NOT NULL DEFAULT '1970-01-01T00:00:00Z',
  last_cursor TEXT NOT NULL DEFAULT '',
  updated_at TIMESTAMPTZ NOT NULL,
  PRIMARY KEY (tenant_id, resource)
);

-- +goose Down
DROP TABLE IF EXISTS ssot_sync_state;
