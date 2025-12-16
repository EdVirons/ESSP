-- +goose Up
-- AI-centric support system schema

-- Add AI fields to chat_sessions table
ALTER TABLE chat_sessions
ADD COLUMN ai_handled BOOLEAN NOT NULL DEFAULT false,
ADD COLUMN ai_resolved BOOLEAN DEFAULT NULL,
ADD COLUMN ai_turns INT NOT NULL DEFAULT 0,
ADD COLUMN escalation_reason TEXT,
ADD COLUMN escalation_summary JSONB DEFAULT '{}',
ADD COLUMN issue_category TEXT,
ADD COLUMN issue_severity TEXT,
ADD COLUMN collected_info JSONB DEFAULT '{}';

-- Update status constraint to include 'ai_active' status
ALTER TABLE chat_sessions DROP CONSTRAINT IF EXISTS chat_sessions_status_check;
ALTER TABLE chat_sessions ADD CONSTRAINT chat_sessions_status_check
    CHECK (status IN ('ai_active', 'waiting', 'active', 'ended'));

-- AI conversation logs for debugging and analytics
CREATE TABLE ai_conversation_logs (
    id TEXT PRIMARY KEY,
    tenant_id TEXT NOT NULL,
    session_id TEXT NOT NULL REFERENCES chat_sessions(id) ON DELETE CASCADE,
    turn_number INT NOT NULL,
    user_message TEXT NOT NULL,
    ai_response TEXT NOT NULL,
    input_tokens INT DEFAULT 0,
    output_tokens INT DEFAULT 0,
    response_time_ms INT,
    escalation_recommended BOOLEAN DEFAULT false,
    escalation_signals JSONB DEFAULT '{}',
    context_used JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_ai_conversation_logs_session ON ai_conversation_logs(session_id);
CREATE INDEX idx_ai_conversation_logs_tenant_date ON ai_conversation_logs(tenant_id, created_at);

-- Daily AI support metrics aggregation
CREATE TABLE ai_support_metrics (
    tenant_id TEXT NOT NULL,
    date DATE NOT NULL,
    total_sessions INT DEFAULT 0,
    ai_resolved INT DEFAULT 0,
    escalated_to_human INT DEFAULT 0,
    avg_turns_to_resolution DECIMAL(5,2) DEFAULT 0,
    avg_response_time_ms INT DEFAULT 0,
    total_input_tokens BIGINT DEFAULT 0,
    total_output_tokens BIGINT DEFAULT 0,
    escalation_reasons JSONB DEFAULT '{}',
    issue_categories JSONB DEFAULT '{}',
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (tenant_id, date)
);

-- AI escalation rules configuration per tenant
CREATE TABLE ai_escalation_rules (
    tenant_id TEXT PRIMARY KEY,
    max_turns INT NOT NULL DEFAULT 10,
    frustration_threshold DECIMAL(3,2) NOT NULL DEFAULT 0.70,
    auto_escalate_categories TEXT[] DEFAULT ARRAY['billing', 'complaint', 'legal'],
    sensitive_keywords TEXT[] DEFAULT ARRAY['refund', 'lawsuit', 'manager', 'supervisor'],
    enabled BOOLEAN NOT NULL DEFAULT true,
    custom_prompts JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Insert default escalation rules
INSERT INTO ai_escalation_rules (tenant_id) VALUES ('default');

-- +goose Down
DROP TABLE IF EXISTS ai_escalation_rules;
DROP TABLE IF EXISTS ai_support_metrics;
DROP TABLE IF EXISTS ai_conversation_logs;

ALTER TABLE chat_sessions DROP CONSTRAINT IF EXISTS chat_sessions_status_check;
ALTER TABLE chat_sessions ADD CONSTRAINT chat_sessions_status_check
    CHECK (status IN ('waiting', 'active', 'ended'));

ALTER TABLE chat_sessions
DROP COLUMN IF EXISTS ai_handled,
DROP COLUMN IF EXISTS ai_resolved,
DROP COLUMN IF EXISTS ai_turns,
DROP COLUMN IF EXISTS escalation_reason,
DROP COLUMN IF EXISTS escalation_summary,
DROP COLUMN IF EXISTS issue_category,
DROP COLUMN IF EXISTS issue_severity,
DROP COLUMN IF EXISTS collected_info;
