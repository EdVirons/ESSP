-- +goose Up

-- User notification preferences for work order alerts
CREATE TABLE IF NOT EXISTS user_notification_preferences (
    id TEXT PRIMARY KEY,
    tenant_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    enabled_types TEXT[] DEFAULT ARRAY['status_change','assignment','approval_requested','approval_decided','deliverable_review','rework_required'],
    in_app_enabled BOOLEAN DEFAULT TRUE,
    email_enabled BOOLEAN DEFAULT FALSE,
    quiet_hours_start TEXT DEFAULT '',
    quiet_hours_end TEXT DEFAULT '',
    quiet_hours_timezone TEXT DEFAULT 'Africa/Nairobi',
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    UNIQUE (tenant_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_user_notification_prefs_tenant_user
    ON user_notification_preferences (tenant_id, user_id);

-- Work order rework/rejection history
CREATE TABLE IF NOT EXISTS work_order_rework_history (
    id TEXT PRIMARY KEY,
    tenant_id TEXT NOT NULL,
    school_id TEXT NOT NULL,
    work_order_id TEXT NOT NULL,
    from_status TEXT NOT NULL,
    to_status TEXT NOT NULL,
    rejection_reason TEXT NOT NULL,
    rejection_category TEXT DEFAULT 'quality',
    rejected_by_user_id TEXT NOT NULL,
    rejected_by_name TEXT DEFAULT '',
    rework_sequence INT DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_work_order_rework_history_wo
    ON work_order_rework_history (tenant_id, school_id, work_order_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_work_order_rework_history_user
    ON work_order_rework_history (tenant_id, rejected_by_user_id, created_at DESC);

-- Add rework tracking columns to work_orders
ALTER TABLE IF EXISTS work_orders
    ADD COLUMN IF NOT EXISTS rework_count INT DEFAULT 0,
    ADD COLUMN IF NOT EXISTS last_rework_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS last_rework_reason TEXT DEFAULT '';

-- Bulk operation tracking for audit and progress
CREATE TABLE IF NOT EXISTS bulk_operation_log (
    id TEXT PRIMARY KEY,
    tenant_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    operation_type TEXT NOT NULL,
    entity_type TEXT DEFAULT 'work_order',
    requested_ids TEXT[] NOT NULL,
    successful_ids TEXT[] DEFAULT '{}',
    failed_ids TEXT[] DEFAULT '{}',
    errors JSONB DEFAULT '[]',
    started_at TIMESTAMPTZ NOT NULL,
    completed_at TIMESTAMPTZ,
    total_count INT DEFAULT 0,
    success_count INT DEFAULT 0,
    failure_count INT DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_bulk_operation_log_tenant_user
    ON bulk_operation_log (tenant_id, user_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_bulk_operation_log_type
    ON bulk_operation_log (tenant_id, operation_type, created_at DESC);

-- Feature flags for tenant-level feature control
CREATE TABLE IF NOT EXISTS feature_config (
    id TEXT PRIMARY KEY,
    tenant_id TEXT NOT NULL,
    feature_key TEXT NOT NULL,
    enabled BOOLEAN DEFAULT TRUE,
    config_value JSONB DEFAULT '{}',
    description TEXT DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    UNIQUE (tenant_id, feature_key)
);

CREATE INDEX IF NOT EXISTS idx_feature_config_tenant_key
    ON feature_config (tenant_id, feature_key);

-- Insert default feature configurations (tenant_id '*' means global default)
INSERT INTO feature_config (id, tenant_id, feature_key, enabled, config_value, description, created_at, updated_at)
VALUES
    ('feat_wo_notif', '*', 'work_order_notifications', true, '{"batch_delay_seconds": 5}', 'Work order status/assignment notifications', NOW(), NOW()),
    ('feat_wo_rework', '*', 'work_order_rework', true, '{"max_rework_count": 5, "require_reason": true}', 'Work order rework/rejection flow', NOW(), NOW()),
    ('feat_wo_bulk', '*', 'work_order_bulk_operations', true, '{"max_batch_size": 100, "rate_limit_per_minute": 10}', 'Bulk work order operations', NOW(), NOW()),
    ('feat_wo_update', '*', 'work_order_update', true, '{}', 'Work order PATCH update endpoint', NOW(), NOW())
ON CONFLICT (tenant_id, feature_key) DO NOTHING;

-- +goose Down

DELETE FROM feature_config WHERE id IN ('feat_wo_notif', 'feat_wo_rework', 'feat_wo_bulk', 'feat_wo_update');

DROP INDEX IF EXISTS idx_feature_config_tenant_key;
DROP TABLE IF EXISTS feature_config;

DROP INDEX IF EXISTS idx_bulk_operation_log_type;
DROP INDEX IF EXISTS idx_bulk_operation_log_tenant_user;
DROP TABLE IF EXISTS bulk_operation_log;

ALTER TABLE IF EXISTS work_orders
    DROP COLUMN IF EXISTS rework_count,
    DROP COLUMN IF EXISTS last_rework_at,
    DROP COLUMN IF EXISTS last_rework_reason;

DROP INDEX IF EXISTS idx_work_order_rework_history_user;
DROP INDEX IF EXISTS idx_work_order_rework_history_wo;
DROP TABLE IF EXISTS work_order_rework_history;

DROP INDEX IF EXISTS idx_user_notification_prefs_tenant_user;
DROP TABLE IF EXISTS user_notification_preferences;
