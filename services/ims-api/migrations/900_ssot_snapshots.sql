-- IMS SSOT snapshot cache
CREATE TABLE IF NOT EXISTS ims_ssot_snapshots (
  tenant_id TEXT NOT NULL,
  kind TEXT NOT NULL, -- school|devices|parts
  version TEXT NOT NULL DEFAULT '1',
  payload JSONB NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  PRIMARY KEY (tenant_id, kind)
);

CREATE INDEX IF NOT EXISTS idx_ims_ssot_snapshots_updated ON ims_ssot_snapshots(updated_at DESC);
