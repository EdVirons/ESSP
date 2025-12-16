-- Migration 008: Audit Logs
-- Comprehensive audit logging for all critical entity operations

CREATE TABLE audit_logs (
    id TEXT PRIMARY KEY,
    tenant_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    user_email TEXT,
    action TEXT NOT NULL,  -- create, update, delete
    entity_type TEXT NOT NULL,  -- incident, work_order, etc.
    entity_id TEXT NOT NULL,
    before_state JSONB,  -- state before change (for updates/deletes)
    after_state JSONB,   -- state after change (for creates/updates)
    ip_address TEXT,
    user_agent TEXT,
    request_id TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes for efficient querying
CREATE INDEX idx_audit_logs_tenant_entity ON audit_logs(tenant_id, entity_type, entity_id);
CREATE INDEX idx_audit_logs_tenant_user ON audit_logs(tenant_id, user_id);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at);
CREATE INDEX idx_audit_logs_tenant_created ON audit_logs(tenant_id, created_at DESC);

-- Add comment for documentation
COMMENT ON TABLE audit_logs IS 'Comprehensive audit trail for all critical entity operations in the IMS system';
COMMENT ON COLUMN audit_logs.action IS 'Type of action performed: create, update, delete';
COMMENT ON COLUMN audit_logs.entity_type IS 'Type of entity being audited: incident, work_order, service_shop, etc.';
COMMENT ON COLUMN audit_logs.before_state IS 'Complete entity state before the change (for updates/deletes)';
COMMENT ON COLUMN audit_logs.after_state IS 'Complete entity state after the change (for creates/updates)';
