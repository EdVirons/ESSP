-- +goose Up
-- Messaging and Livechat System
-- Message threads (conversations)
CREATE TABLE IF NOT EXISTS message_threads (
    id TEXT PRIMARY KEY,
    tenant_id TEXT NOT NULL,
    school_id TEXT NOT NULL,

    -- Thread metadata
    subject TEXT NOT NULL,
    thread_type TEXT NOT NULL DEFAULT 'general', -- 'general', 'incident', 'livechat'
    status TEXT NOT NULL DEFAULT 'open',         -- 'open', 'closed', 'archived'

    -- Linking to incidents (optional)
    incident_id TEXT,                            -- NULL if not linked to incident

    -- Participants tracking
    created_by TEXT NOT NULL,                    -- User ID who started thread
    created_by_role TEXT NOT NULL,               -- Role of creator
    created_by_name TEXT NOT NULL DEFAULT '',    -- Name of creator for display

    -- Stats
    message_count INT NOT NULL DEFAULT 0,
    unread_count_school INT NOT NULL DEFAULT 0,  -- Unread by school contact
    unread_count_support INT NOT NULL DEFAULT 0, -- Unread by support team

    -- Timestamps
    last_message_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    closed_at TIMESTAMPTZ
);

CREATE INDEX idx_threads_tenant_school ON message_threads(tenant_id, school_id, last_message_at DESC);
CREATE INDEX idx_threads_incident ON message_threads(tenant_id, incident_id) WHERE incident_id IS NOT NULL;
CREATE INDEX idx_threads_status ON message_threads(tenant_id, status, last_message_at DESC);
CREATE INDEX idx_threads_created_by ON message_threads(created_by, created_at DESC);

-- Individual messages
CREATE TABLE IF NOT EXISTS messages (
    id TEXT PRIMARY KEY,
    tenant_id TEXT NOT NULL,
    thread_id TEXT NOT NULL REFERENCES message_threads(id) ON DELETE CASCADE,

    -- Sender info
    sender_id TEXT NOT NULL,
    sender_name TEXT NOT NULL,
    sender_role TEXT NOT NULL,

    -- Content
    content TEXT NOT NULL,
    content_type TEXT NOT NULL DEFAULT 'text', -- 'text', 'system', 'attachment'

    -- Metadata for system messages
    metadata JSONB DEFAULT '{}',

    -- Edit tracking
    edited_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,

    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_messages_thread ON messages(thread_id, created_at ASC);
CREATE INDEX idx_messages_sender ON messages(sender_id, created_at DESC);
CREATE INDEX idx_messages_tenant ON messages(tenant_id, created_at DESC);
CREATE INDEX idx_messages_search ON messages USING gin(to_tsvector('english', content));

-- Message attachments (extends existing attachment system)
CREATE TABLE IF NOT EXISTS message_attachments (
    id TEXT PRIMARY KEY,
    tenant_id TEXT NOT NULL,
    message_id TEXT NOT NULL REFERENCES messages(id) ON DELETE CASCADE,

    -- File info
    file_name TEXT NOT NULL,
    content_type TEXT NOT NULL,
    size_bytes BIGINT NOT NULL DEFAULT 0,
    object_key TEXT NOT NULL,

    -- Thumbnail for images
    thumbnail_key TEXT,

    created_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_msg_attachments_message ON message_attachments(message_id);

-- Read receipts tracking
CREATE TABLE IF NOT EXISTS message_read_receipts (
    thread_id TEXT NOT NULL REFERENCES message_threads(id) ON DELETE CASCADE,
    user_id TEXT NOT NULL,
    last_read_message_id TEXT NOT NULL REFERENCES messages(id),
    last_read_at TIMESTAMPTZ NOT NULL,
    PRIMARY KEY (thread_id, user_id)
);

-- Livechat sessions
CREATE TABLE IF NOT EXISTS chat_sessions (
    id TEXT PRIMARY KEY,
    tenant_id TEXT NOT NULL,
    school_id TEXT NOT NULL,
    thread_id TEXT NOT NULL REFERENCES message_threads(id),

    -- Participants
    school_contact_id TEXT NOT NULL,
    school_contact_name TEXT NOT NULL DEFAULT '',
    assigned_agent_id TEXT,
    assigned_agent_name TEXT,

    -- Session state
    status TEXT NOT NULL DEFAULT 'waiting', -- 'waiting', 'active', 'ended'
    queue_position INT,

    -- Timing
    started_at TIMESTAMPTZ NOT NULL,
    agent_joined_at TIMESTAMPTZ,
    ended_at TIMESTAMPTZ,

    -- Metrics
    first_response_seconds INT,
    total_messages INT NOT NULL DEFAULT 0,

    -- Rating (optional post-chat)
    rating INT,                             -- 1-5 stars
    feedback TEXT,

    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_chat_sessions_tenant_status ON chat_sessions(tenant_id, status, started_at);
CREATE INDEX idx_chat_sessions_agent ON chat_sessions(assigned_agent_id, status);
CREATE INDEX idx_chat_sessions_school ON chat_sessions(school_id, created_at DESC);
CREATE INDEX idx_chat_sessions_thread ON chat_sessions(thread_id);

-- Thread participants (for multi-party threads)
CREATE TABLE IF NOT EXISTS thread_participants (
    thread_id TEXT NOT NULL REFERENCES message_threads(id) ON DELETE CASCADE,
    user_id TEXT NOT NULL,
    user_name TEXT NOT NULL,
    user_role TEXT NOT NULL,
    joined_at TIMESTAMPTZ NOT NULL,
    left_at TIMESTAMPTZ,
    PRIMARY KEY (thread_id, user_id)
);

CREATE INDEX idx_participants_user ON thread_participants(user_id, thread_id);

-- Agent availability tracking
CREATE TABLE IF NOT EXISTS agent_availability (
    tenant_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    is_available BOOLEAN NOT NULL DEFAULT false,
    max_concurrent_chats INT NOT NULL DEFAULT 3,
    current_chat_count INT NOT NULL DEFAULT 0,
    last_seen_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    PRIMARY KEY (tenant_id, user_id)
);

CREATE INDEX idx_agent_availability_available ON agent_availability(tenant_id, is_available, current_chat_count);

-- +goose Down
DROP TABLE IF EXISTS agent_availability;
DROP TABLE IF EXISTS thread_participants;
DROP TABLE IF EXISTS chat_sessions;
DROP TABLE IF EXISTS message_read_receipts;
DROP TABLE IF EXISTS message_attachments;
DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS message_threads;
