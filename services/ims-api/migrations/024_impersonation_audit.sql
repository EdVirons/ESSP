-- +goose Up
-- Migration 024: Impersonation Audit Enhancement
-- Adds columns to track when actions are performed on behalf of another user

-- Add impersonation columns to audit_logs
ALTER TABLE audit_logs
    ADD COLUMN IF NOT EXISTS impersonated_user_id TEXT,
    ADD COLUMN IF NOT EXISTS impersonated_user_email TEXT,
    ADD COLUMN IF NOT EXISTS impersonation_reason TEXT;

-- Index for efficient querying of impersonated actions
CREATE INDEX IF NOT EXISTS idx_audit_logs_impersonation ON audit_logs(impersonated_user_id)
    WHERE impersonated_user_id IS NOT NULL;

-- Add comments
COMMENT ON COLUMN audit_logs.impersonated_user_id IS 'User ID of the school contact being impersonated (null if not impersonating)';
COMMENT ON COLUMN audit_logs.impersonated_user_email IS 'Email of the school contact being impersonated';
COMMENT ON COLUMN audit_logs.impersonation_reason IS 'Reason provided for the impersonation session';

-- +goose Down
DROP INDEX IF EXISTS idx_audit_logs_impersonation;
ALTER TABLE audit_logs
    DROP COLUMN IF EXISTS impersonation_reason,
    DROP COLUMN IF EXISTS impersonated_user_email,
    DROP COLUMN IF EXISTS impersonated_user_id;
