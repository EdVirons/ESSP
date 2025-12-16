-- +goose Up
-- Add program_type column to support multiple service program types
-- Existing programs will default to 'full_installation'

ALTER TABLE school_service_programs
  ADD COLUMN IF NOT EXISTS program_type TEXT NOT NULL DEFAULT 'full_installation';

-- Create index for efficient filtering by program type
CREATE INDEX IF NOT EXISTS idx_programs_type
  ON school_service_programs (tenant_id, program_type, created_at DESC);

-- Add program_type to phase_checklist_templates for type-specific checklists
ALTER TABLE phase_checklist_templates
  ADD COLUMN IF NOT EXISTS program_type TEXT NOT NULL DEFAULT 'full_installation';

-- Update checklist index to include program_type
DROP INDEX IF EXISTS idx_phase_checklists;
CREATE INDEX IF NOT EXISTS idx_phase_checklists_v2
  ON phase_checklist_templates (tenant_id, program_type, phase_type, required DESC, created_at DESC);

-- +goose Down
ALTER TABLE phase_checklist_templates
  DROP COLUMN IF EXISTS program_type;

DROP INDEX IF EXISTS idx_phase_checklists_v2;
CREATE INDEX IF NOT EXISTS idx_phase_checklists
  ON phase_checklist_templates (tenant_id, phase_type, required DESC, created_at DESC);

DROP INDEX IF EXISTS idx_programs_type;
ALTER TABLE school_service_programs
  DROP COLUMN IF EXISTS program_type;
