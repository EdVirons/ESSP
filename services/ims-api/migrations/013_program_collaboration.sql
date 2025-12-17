-- +goose Up
-- Migration: Program Multi-Team Collaboration
-- Adds team management, activity feed, and notifications for programs

-- ===========================================
-- 1. Program Team Members
-- ===========================================
CREATE TABLE IF NOT EXISTS program_team_members (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL,
  program_id TEXT NOT NULL REFERENCES school_service_programs(id) ON DELETE CASCADE,
  user_id TEXT NOT NULL,
  user_email TEXT NOT NULL DEFAULT '',
  user_name TEXT NOT NULL DEFAULT '',
  role TEXT NOT NULL DEFAULT 'collaborator',  -- owner|collaborator|viewer
  assigned_phases TEXT[] NOT NULL DEFAULT '{}',
  responsibility TEXT NOT NULL DEFAULT '',
  assigned_by_user_id TEXT NOT NULL DEFAULT '',
  assigned_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  removed_at TIMESTAMPTZ,                      -- soft delete
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Unique constraint for active members only
CREATE UNIQUE INDEX IF NOT EXISTS idx_program_team_unique_active
  ON program_team_members(tenant_id, program_id, user_id)
  WHERE removed_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_program_team_program ON program_team_members(tenant_id, program_id) WHERE removed_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_program_team_user ON program_team_members(tenant_id, user_id) WHERE removed_at IS NULL;

-- ===========================================
-- 2. Program Activities (Activity Feed)
-- ===========================================
CREATE TABLE IF NOT EXISTS program_activities (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL,
  program_id TEXT NOT NULL REFERENCES school_service_programs(id) ON DELETE CASCADE,
  phase_id TEXT,                               -- optional, for phase-specific activities
  work_order_id TEXT,                          -- optional, for WO-related activities
  activity_type TEXT NOT NULL,                 -- comment|note|file_upload|status_change|
                                               -- assignment|work_order|phase_transition|mention
  actor_user_id TEXT NOT NULL,
  actor_email TEXT NOT NULL DEFAULT '',
  actor_name TEXT NOT NULL DEFAULT '',
  content TEXT NOT NULL DEFAULT '',            -- for comments/notes
  metadata JSONB NOT NULL DEFAULT '{}',        -- flexible data for different activity types
  attachment_ids TEXT[] NOT NULL DEFAULT '{}', -- for file uploads
  visibility TEXT NOT NULL DEFAULT 'team',     -- team|public|private
  is_pinned BOOLEAN NOT NULL DEFAULT FALSE,
  edited_at TIMESTAMPTZ,
  deleted_at TIMESTAMPTZ,                      -- soft delete
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_program_activities_program ON program_activities(tenant_id, program_id, created_at DESC) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_program_activities_phase ON program_activities(tenant_id, program_id, phase_id, created_at DESC) WHERE deleted_at IS NULL AND phase_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_program_activities_type ON program_activities(tenant_id, program_id, activity_type, created_at DESC) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_program_activities_user ON program_activities(tenant_id, actor_user_id, created_at DESC) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_program_activities_pinned ON program_activities(tenant_id, program_id, is_pinned, created_at DESC) WHERE deleted_at IS NULL AND is_pinned = TRUE;

-- ===========================================
-- 3. Program Attachments
-- ===========================================
CREATE TABLE IF NOT EXISTS program_attachments (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL,
  program_id TEXT NOT NULL REFERENCES school_service_programs(id) ON DELETE CASCADE,
  phase_id TEXT,
  activity_id TEXT REFERENCES program_activities(id) ON DELETE SET NULL,
  file_name TEXT NOT NULL,
  content_type TEXT NOT NULL,
  size_bytes BIGINT NOT NULL DEFAULT 0,
  object_key TEXT NOT NULL,                    -- MinIO/S3 key
  uploaded_by_user_id TEXT NOT NULL,
  uploaded_by_user_name TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_program_attachments_program ON program_attachments(tenant_id, program_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_program_attachments_activity ON program_attachments(activity_id) WHERE activity_id IS NOT NULL;

-- ===========================================
-- 4. User Notifications
-- ===========================================
CREATE TABLE IF NOT EXISTS user_notifications (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL,
  user_id TEXT NOT NULL,
  notification_type TEXT NOT NULL,             -- assignment|mention|status_change|comment|work_order
  entity_type TEXT NOT NULL,                   -- program|phase|work_order
  entity_id TEXT NOT NULL,
  program_id TEXT,                             -- for quick filtering
  title TEXT NOT NULL,
  body TEXT NOT NULL DEFAULT '',
  metadata JSONB NOT NULL DEFAULT '{}',
  is_read BOOLEAN NOT NULL DEFAULT FALSE,
  read_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_notifications_user_unread ON user_notifications(tenant_id, user_id, is_read, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_notifications_user_all ON user_notifications(tenant_id, user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_notifications_program ON user_notifications(tenant_id, program_id, created_at DESC) WHERE program_id IS NOT NULL;

-- ===========================================
-- 5. Modify service_phases for user ownership
-- ===========================================
ALTER TABLE service_phases
  ADD COLUMN IF NOT EXISTS owner_user_id TEXT NOT NULL DEFAULT '',
  ADD COLUMN IF NOT EXISTS owner_user_name TEXT NOT NULL DEFAULT '',
  ADD COLUMN IF NOT EXISTS status_changed_at TIMESTAMPTZ,
  ADD COLUMN IF NOT EXISTS status_changed_by_user_id TEXT NOT NULL DEFAULT '',
  ADD COLUMN IF NOT EXISTS status_changed_by_user_name TEXT NOT NULL DEFAULT '';

-- ===========================================
-- 6. Modify work_orders for program creation tracking
-- ===========================================
ALTER TABLE work_orders
  ADD COLUMN IF NOT EXISTS created_from_program BOOLEAN NOT NULL DEFAULT FALSE,
  ADD COLUMN IF NOT EXISTS created_by_user_id TEXT NOT NULL DEFAULT '',
  ADD COLUMN IF NOT EXISTS created_by_user_name TEXT NOT NULL DEFAULT '';

-- Index for work orders by program
CREATE INDEX IF NOT EXISTS idx_work_orders_program_id ON work_orders(tenant_id, program_id) WHERE program_id IS NOT NULL AND program_id != '';

-- ===========================================
-- 7. Phase user assignments (for multiple assignees per phase)
-- ===========================================
CREATE TABLE IF NOT EXISTS phase_user_assignments (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL,
  phase_id TEXT NOT NULL REFERENCES service_phases(id) ON DELETE CASCADE,
  program_id TEXT NOT NULL,                    -- denormalized for queries
  user_id TEXT NOT NULL,
  user_email TEXT NOT NULL DEFAULT '',
  user_name TEXT NOT NULL DEFAULT '',
  assignment_type TEXT NOT NULL DEFAULT 'collaborator', -- owner|collaborator|reviewer
  assigned_by_user_id TEXT NOT NULL DEFAULT '',
  assigned_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  completed_at TIMESTAMPTZ,
  removed_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_phase_assignments_unique_active
  ON phase_user_assignments(tenant_id, phase_id, user_id, assignment_type)
  WHERE removed_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_phase_assignments_phase ON phase_user_assignments(tenant_id, phase_id) WHERE removed_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_phase_assignments_user ON phase_user_assignments(tenant_id, user_id) WHERE removed_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_phase_assignments_program ON phase_user_assignments(tenant_id, program_id) WHERE removed_at IS NULL;

-- +goose Down
DROP TABLE IF EXISTS phase_user_assignments;
ALTER TABLE work_orders DROP COLUMN IF EXISTS created_by_user_name;
ALTER TABLE work_orders DROP COLUMN IF EXISTS created_by_user_id;
ALTER TABLE work_orders DROP COLUMN IF EXISTS created_from_program;
ALTER TABLE service_phases DROP COLUMN IF EXISTS status_changed_by_user_name;
ALTER TABLE service_phases DROP COLUMN IF EXISTS status_changed_by_user_id;
ALTER TABLE service_phases DROP COLUMN IF EXISTS status_changed_at;
ALTER TABLE service_phases DROP COLUMN IF EXISTS owner_user_name;
ALTER TABLE service_phases DROP COLUMN IF EXISTS owner_user_id;
DROP TABLE IF EXISTS user_notifications;
DROP TABLE IF EXISTS program_attachments;
DROP TABLE IF EXISTS program_activities;
DROP TABLE IF EXISTS program_team_members;
