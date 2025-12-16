-- +goose Up
-- Migration: Rename all "program" references to "project"

-- ===========================================
-- 1. Rename main tables
-- ===========================================
ALTER TABLE school_service_programs RENAME TO school_service_projects;
ALTER TABLE program_team_members RENAME TO project_team_members;
ALTER TABLE program_activities RENAME TO project_activities;
ALTER TABLE program_attachments RENAME TO project_attachments;

-- ===========================================
-- 2. Rename columns in school_service_projects
-- ===========================================
ALTER TABLE school_service_projects RENAME COLUMN program_type TO project_type;

-- ===========================================
-- 3. Rename columns in service_phases
-- ===========================================
ALTER TABLE service_phases RENAME COLUMN program_id TO project_id;

-- ===========================================
-- 4. Rename columns in site_surveys
-- ===========================================
ALTER TABLE site_surveys RENAME COLUMN program_id TO project_id;

-- ===========================================
-- 5. Rename columns in boq_items
-- ===========================================
ALTER TABLE boq_items RENAME COLUMN program_id TO project_id;

-- ===========================================
-- 6. Rename columns in project_team_members (formerly program_team_members)
-- ===========================================
ALTER TABLE project_team_members RENAME COLUMN program_id TO project_id;

-- ===========================================
-- 7. Rename columns in project_activities (formerly program_activities)
-- ===========================================
ALTER TABLE project_activities RENAME COLUMN program_id TO project_id;

-- ===========================================
-- 8. Rename columns in project_attachments (formerly program_attachments)
-- ===========================================
ALTER TABLE project_attachments RENAME COLUMN program_id TO project_id;

-- ===========================================
-- 9. Rename columns in user_notifications
-- ===========================================
ALTER TABLE user_notifications RENAME COLUMN program_id TO project_id;

-- ===========================================
-- 10. Rename columns in phase_user_assignments
-- ===========================================
ALTER TABLE phase_user_assignments RENAME COLUMN program_id TO project_id;

-- ===========================================
-- 11. Rename columns in work_orders
-- ===========================================
ALTER TABLE work_orders RENAME COLUMN program_id TO project_id;
ALTER TABLE work_orders RENAME COLUMN created_from_program TO created_from_project;

-- ===========================================
-- 12. Rename columns in phase_checklist_templates
-- ===========================================
ALTER TABLE phase_checklist_templates RENAME COLUMN program_type TO project_type;

-- ===========================================
-- 13. Drop old indexes and create new ones with updated names
-- ===========================================

-- school_service_projects indexes
DROP INDEX IF EXISTS idx_programs_school;
CREATE INDEX IF NOT EXISTS idx_projects_school
  ON school_service_projects (tenant_id, school_id, created_at DESC);

DROP INDEX IF EXISTS idx_programs_type;
CREATE INDEX IF NOT EXISTS idx_projects_type
  ON school_service_projects (tenant_id, project_type, created_at DESC);

-- service_phases indexes
DROP INDEX IF EXISTS idx_phases_program;
CREATE INDEX IF NOT EXISTS idx_phases_project
  ON service_phases (tenant_id, project_id, created_at DESC);

-- site_surveys indexes
DROP INDEX IF EXISTS idx_surveys_program;
CREATE INDEX IF NOT EXISTS idx_surveys_project
  ON site_surveys (tenant_id, project_id, created_at DESC);

-- boq_items indexes
DROP INDEX IF EXISTS idx_boq_program;
CREATE INDEX IF NOT EXISTS idx_boq_project
  ON boq_items (tenant_id, project_id, created_at DESC);

-- project_team_members indexes
DROP INDEX IF EXISTS idx_program_team_unique_active;
CREATE UNIQUE INDEX IF NOT EXISTS idx_project_team_unique_active
  ON project_team_members(tenant_id, project_id, user_id)
  WHERE removed_at IS NULL;

DROP INDEX IF EXISTS idx_program_team_program;
CREATE INDEX IF NOT EXISTS idx_project_team_project
  ON project_team_members(tenant_id, project_id) WHERE removed_at IS NULL;

DROP INDEX IF EXISTS idx_program_team_user;
CREATE INDEX IF NOT EXISTS idx_project_team_user
  ON project_team_members(tenant_id, user_id) WHERE removed_at IS NULL;

-- project_activities indexes
DROP INDEX IF EXISTS idx_program_activities_program;
CREATE INDEX IF NOT EXISTS idx_project_activities_project
  ON project_activities(tenant_id, project_id, created_at DESC) WHERE deleted_at IS NULL;

DROP INDEX IF EXISTS idx_program_activities_phase;
CREATE INDEX IF NOT EXISTS idx_project_activities_phase
  ON project_activities(tenant_id, project_id, phase_id, created_at DESC) WHERE deleted_at IS NULL AND phase_id IS NOT NULL;

DROP INDEX IF EXISTS idx_program_activities_type;
CREATE INDEX IF NOT EXISTS idx_project_activities_type
  ON project_activities(tenant_id, project_id, activity_type, created_at DESC) WHERE deleted_at IS NULL;

DROP INDEX IF EXISTS idx_program_activities_user;
CREATE INDEX IF NOT EXISTS idx_project_activities_user
  ON project_activities(tenant_id, actor_user_id, created_at DESC) WHERE deleted_at IS NULL;

DROP INDEX IF EXISTS idx_program_activities_pinned;
CREATE INDEX IF NOT EXISTS idx_project_activities_pinned
  ON project_activities(tenant_id, project_id, is_pinned, created_at DESC) WHERE deleted_at IS NULL AND is_pinned = TRUE;

-- project_attachments indexes
DROP INDEX IF EXISTS idx_program_attachments_program;
CREATE INDEX IF NOT EXISTS idx_project_attachments_project
  ON project_attachments(tenant_id, project_id, created_at DESC);

DROP INDEX IF EXISTS idx_program_attachments_activity;
CREATE INDEX IF NOT EXISTS idx_project_attachments_activity
  ON project_attachments(activity_id) WHERE activity_id IS NOT NULL;

-- user_notifications indexes
DROP INDEX IF EXISTS idx_notifications_program;
CREATE INDEX IF NOT EXISTS idx_notifications_project
  ON user_notifications(tenant_id, project_id, created_at DESC) WHERE project_id IS NOT NULL;

-- phase_user_assignments indexes
DROP INDEX IF EXISTS idx_phase_assignments_program;
CREATE INDEX IF NOT EXISTS idx_phase_assignments_project
  ON phase_user_assignments(tenant_id, project_id) WHERE removed_at IS NULL;

-- work_orders indexes
DROP INDEX IF EXISTS idx_work_orders_program_id;
CREATE INDEX IF NOT EXISTS idx_work_orders_project_id
  ON work_orders(tenant_id, project_id) WHERE project_id IS NOT NULL AND project_id != '';

-- phase_checklist_templates indexes
DROP INDEX IF EXISTS idx_phase_checklists_v2;
CREATE INDEX IF NOT EXISTS idx_phase_checklists_v3
  ON phase_checklist_templates (tenant_id, project_type, phase_type, required DESC, created_at DESC);

-- +goose Down
-- Reverse all renames

-- Rename indexes back
DROP INDEX IF EXISTS idx_phase_checklists_v3;
CREATE INDEX IF NOT EXISTS idx_phase_checklists_v2
  ON phase_checklist_templates (tenant_id, program_type, phase_type, required DESC, created_at DESC);

DROP INDEX IF EXISTS idx_work_orders_project_id;
CREATE INDEX IF NOT EXISTS idx_work_orders_program_id
  ON work_orders(tenant_id, program_id) WHERE program_id IS NOT NULL AND program_id != '';

DROP INDEX IF EXISTS idx_phase_assignments_project;
CREATE INDEX IF NOT EXISTS idx_phase_assignments_program
  ON phase_user_assignments(tenant_id, program_id) WHERE removed_at IS NULL;

DROP INDEX IF EXISTS idx_notifications_project;
CREATE INDEX IF NOT EXISTS idx_notifications_program
  ON user_notifications(tenant_id, program_id, created_at DESC) WHERE program_id IS NOT NULL;

DROP INDEX IF EXISTS idx_project_attachments_activity;
DROP INDEX IF EXISTS idx_project_attachments_project;
CREATE INDEX IF NOT EXISTS idx_program_attachments_program
  ON program_attachments(tenant_id, program_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_program_attachments_activity
  ON program_attachments(activity_id) WHERE activity_id IS NOT NULL;

DROP INDEX IF EXISTS idx_project_activities_pinned;
DROP INDEX IF EXISTS idx_project_activities_user;
DROP INDEX IF EXISTS idx_project_activities_type;
DROP INDEX IF EXISTS idx_project_activities_phase;
DROP INDEX IF EXISTS idx_project_activities_project;
CREATE INDEX IF NOT EXISTS idx_program_activities_program
  ON program_activities(tenant_id, program_id, created_at DESC) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_program_activities_phase
  ON program_activities(tenant_id, program_id, phase_id, created_at DESC) WHERE deleted_at IS NULL AND phase_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_program_activities_type
  ON program_activities(tenant_id, program_id, activity_type, created_at DESC) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_program_activities_user
  ON program_activities(tenant_id, actor_user_id, created_at DESC) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_program_activities_pinned
  ON program_activities(tenant_id, program_id, is_pinned, created_at DESC) WHERE deleted_at IS NULL AND is_pinned = TRUE;

DROP INDEX IF EXISTS idx_project_team_user;
DROP INDEX IF EXISTS idx_project_team_project;
DROP INDEX IF EXISTS idx_project_team_unique_active;
CREATE UNIQUE INDEX IF NOT EXISTS idx_program_team_unique_active
  ON program_team_members(tenant_id, program_id, user_id) WHERE removed_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_program_team_program
  ON program_team_members(tenant_id, program_id) WHERE removed_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_program_team_user
  ON program_team_members(tenant_id, user_id) WHERE removed_at IS NULL;

DROP INDEX IF EXISTS idx_boq_project;
CREATE INDEX IF NOT EXISTS idx_boq_program
  ON boq_items (tenant_id, program_id, created_at DESC);

DROP INDEX IF EXISTS idx_surveys_project;
CREATE INDEX IF NOT EXISTS idx_surveys_program
  ON site_surveys (tenant_id, program_id, created_at DESC);

DROP INDEX IF EXISTS idx_phases_project;
CREATE INDEX IF NOT EXISTS idx_phases_program
  ON service_phases (tenant_id, program_id, created_at DESC);

DROP INDEX IF EXISTS idx_projects_type;
DROP INDEX IF EXISTS idx_projects_school;
CREATE INDEX IF NOT EXISTS idx_programs_school
  ON school_service_programs (tenant_id, school_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_programs_type
  ON school_service_programs (tenant_id, program_type, created_at DESC);

-- Rename columns back
ALTER TABLE phase_checklist_templates RENAME COLUMN project_type TO program_type;
ALTER TABLE work_orders RENAME COLUMN created_from_project TO created_from_program;
ALTER TABLE work_orders RENAME COLUMN project_id TO program_id;
ALTER TABLE phase_user_assignments RENAME COLUMN project_id TO program_id;
ALTER TABLE user_notifications RENAME COLUMN project_id TO program_id;
ALTER TABLE project_attachments RENAME COLUMN project_id TO program_id;
ALTER TABLE project_activities RENAME COLUMN project_id TO program_id;
ALTER TABLE project_team_members RENAME COLUMN project_id TO program_id;
ALTER TABLE boq_items RENAME COLUMN project_id TO program_id;
ALTER TABLE site_surveys RENAME COLUMN project_id TO program_id;
ALTER TABLE service_phases RENAME COLUMN project_id TO program_id;
ALTER TABLE school_service_projects RENAME COLUMN project_type TO program_type;

-- Rename tables back
ALTER TABLE project_attachments RENAME TO program_attachments;
ALTER TABLE project_activities RENAME TO program_activities;
ALTER TABLE project_team_members RENAME TO program_team_members;
ALTER TABLE school_service_projects RENAME TO school_service_programs;
