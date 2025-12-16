-- +goose Up

CREATE TABLE IF NOT EXISTS school_service_programs (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL,
  school_id TEXT NOT NULL,
  status TEXT NOT NULL DEFAULT 'active',              -- active|paused|completed
  current_phase TEXT NOT NULL DEFAULT 'demo',         -- demo|survey|install|integrate|commission|ops
  start_date DATE,
  go_live_date DATE,
  account_manager_user_id TEXT NOT NULL DEFAULT '',
  notes TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_programs_school
  ON school_service_programs (tenant_id, school_id, created_at DESC);

CREATE TABLE IF NOT EXISTS service_phases (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL,
  program_id TEXT NOT NULL,
  phase_type TEXT NOT NULL,                           -- demo|survey|install|integrate|commission|ops
  status TEXT NOT NULL DEFAULT 'pending',             -- pending|in_progress|blocked|done
  owner_role TEXT NOT NULL DEFAULT '',                -- sales|engineer|lead_tech|admin
  start_date DATE,
  end_date DATE,
  notes TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_phases_program
  ON service_phases (tenant_id, program_id, created_at DESC);

CREATE TABLE IF NOT EXISTS site_surveys (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL,
  program_id TEXT NOT NULL,
  status TEXT NOT NULL DEFAULT 'draft',               -- draft|submitted|approved
  conducted_by_user_id TEXT NOT NULL DEFAULT '',
  conducted_at TIMESTAMPTZ,
  summary TEXT NOT NULL DEFAULT '',
  risks TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_surveys_program
  ON site_surveys (tenant_id, program_id, created_at DESC);

CREATE TABLE IF NOT EXISTS survey_rooms (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL,
  survey_id TEXT NOT NULL,
  name TEXT NOT NULL,
  room_type TEXT NOT NULL DEFAULT '',                 -- classroom|lab|office|server_room|store|other
  floor TEXT NOT NULL DEFAULT '',
  power_notes TEXT NOT NULL DEFAULT '',
  network_notes TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_survey_rooms
  ON survey_rooms (tenant_id, survey_id, created_at DESC);

CREATE TABLE IF NOT EXISTS survey_photos (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL,
  survey_id TEXT NOT NULL,
  room_id TEXT NOT NULL DEFAULT '',
  attachment_id TEXT NOT NULL,
  caption TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_survey_photos
  ON survey_photos (tenant_id, survey_id, created_at DESC);

CREATE TABLE IF NOT EXISTS boq_items (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL,
  program_id TEXT NOT NULL,
  category TEXT NOT NULL,                             -- cabling|switches|racks|power|labor|devices|other
  description TEXT NOT NULL,
  part_id TEXT NOT NULL DEFAULT '',                   -- Parts SSOT id (optional)
  qty BIGINT NOT NULL DEFAULT 0,
  unit TEXT NOT NULL DEFAULT '',
  estimated_cost_cents BIGINT NOT NULL DEFAULT 0,
  approved BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_boq_program
  ON boq_items (tenant_id, program_id, created_at DESC);

-- +goose Down
DROP TABLE IF EXISTS boq_items;
DROP TABLE IF EXISTS survey_photos;
DROP TABLE IF EXISTS survey_rooms;
DROP TABLE IF EXISTS site_surveys;
DROP TABLE IF EXISTS service_phases;
DROP TABLE IF EXISTS school_service_programs;
