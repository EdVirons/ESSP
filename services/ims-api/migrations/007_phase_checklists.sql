-- +goose Up

-- Phase checklist templates (SSOT-ish per tenant; can later be global SSOT)
CREATE TABLE IF NOT EXISTS phase_checklist_templates (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL,
  phase_type TEXT NOT NULL,                -- demo|survey|install|integrate|commission|ops
  title TEXT NOT NULL,
  description TEXT NOT NULL DEFAULT '',
  required BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_phase_checklists
  ON phase_checklist_templates (tenant_id, phase_type, required DESC, created_at DESC);

-- Link work-order deliverables to a phase (optional; derived from WO.phase_id, but useful for queries)
ALTER TABLE IF EXISTS work_order_deliverables
  ADD COLUMN IF NOT EXISTS phase_id TEXT NOT NULL DEFAULT '';

CREATE INDEX IF NOT EXISTS idx_deliverables_phase
  ON work_order_deliverables (tenant_id, school_id, phase_id, created_at DESC);

-- +goose Down
ALTER TABLE IF EXISTS work_order_deliverables
  DROP COLUMN IF EXISTS phase_id;

DROP TABLE IF EXISTS phase_checklist_templates;
