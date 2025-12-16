-- +goose Up

CREATE TABLE IF NOT EXISTS school_contacts (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL,
  school_id TEXT NOT NULL,
  user_id TEXT NOT NULL DEFAULT '',
  name TEXT NOT NULL,
  phone TEXT NOT NULL DEFAULT '',
  email TEXT NOT NULL DEFAULT '',
  role TEXT NOT NULL DEFAULT 'point_of_contact',
  is_primary BOOLEAN NOT NULL DEFAULT FALSE,
  active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_school_contacts_school
  ON school_contacts (tenant_id, school_id, is_primary DESC, active DESC, created_at DESC);

ALTER TABLE IF EXISTS incidents
  ADD COLUMN IF NOT EXISTS reporter_contact_id TEXT NOT NULL DEFAULT '';

ALTER TABLE IF EXISTS work_orders
  ADD COLUMN IF NOT EXISTS program_id TEXT NOT NULL DEFAULT '',
  ADD COLUMN IF NOT EXISTS phase_id TEXT NOT NULL DEFAULT '',
  ADD COLUMN IF NOT EXISTS onsite_contact_id TEXT NOT NULL DEFAULT '',
  ADD COLUMN IF NOT EXISTS approval_status TEXT NOT NULL DEFAULT 'not_required';

CREATE INDEX IF NOT EXISTS idx_work_orders_program
  ON work_orders (tenant_id, school_id, program_id, created_at DESC);

CREATE TABLE IF NOT EXISTS work_order_schedules (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL,
  school_id TEXT NOT NULL,
  work_order_id TEXT NOT NULL,
  scheduled_start TIMESTAMPTZ,
  scheduled_end TIMESTAMPTZ,
  timezone TEXT NOT NULL DEFAULT 'Africa/Nairobi',
  notes TEXT NOT NULL DEFAULT '',
  created_by_user_id TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_work_order_schedules_wo
  ON work_order_schedules (tenant_id, school_id, work_order_id, created_at DESC);

CREATE TABLE IF NOT EXISTS work_order_deliverables (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL,
  school_id TEXT NOT NULL,
  work_order_id TEXT NOT NULL,
  title TEXT NOT NULL,
  description TEXT NOT NULL DEFAULT '',
  status TEXT NOT NULL DEFAULT 'pending',
  evidence_attachment_id TEXT NOT NULL DEFAULT '',
  submitted_by_user_id TEXT NOT NULL DEFAULT '',
  submitted_at TIMESTAMPTZ,
  reviewed_by_user_id TEXT NOT NULL DEFAULT '',
  reviewed_at TIMESTAMPTZ,
  review_notes TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_work_order_deliverables_wo
  ON work_order_deliverables (tenant_id, school_id, work_order_id, created_at DESC);

CREATE TABLE IF NOT EXISTS work_order_approvals (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL,
  school_id TEXT NOT NULL,
  work_order_id TEXT NOT NULL,
  approval_type TEXT NOT NULL DEFAULT 'school_signoff',
  requested_by_user_id TEXT NOT NULL DEFAULT '',
  requested_at TIMESTAMPTZ NOT NULL,
  status TEXT NOT NULL DEFAULT 'pending',
  decided_by_user_id TEXT NOT NULL DEFAULT '',
  decided_at TIMESTAMPTZ,
  decision_notes TEXT NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_work_order_approvals_wo
  ON work_order_approvals (tenant_id, school_id, work_order_id, requested_at DESC);

-- +goose Down
DROP TABLE IF EXISTS work_order_approvals;
DROP TABLE IF EXISTS work_order_deliverables;
DROP TABLE IF EXISTS work_order_schedules;

ALTER TABLE IF EXISTS work_orders
  DROP COLUMN IF EXISTS program_id,
  DROP COLUMN IF EXISTS phase_id,
  DROP COLUMN IF EXISTS onsite_contact_id,
  DROP COLUMN IF EXISTS approval_status;

ALTER TABLE IF EXISTS incidents
  DROP COLUMN IF EXISTS reporter_contact_id;

DROP TABLE IF EXISTS school_contacts;
