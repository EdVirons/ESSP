-- +goose Up
-- Knowledge Base Articles table
CREATE TABLE IF NOT EXISTS kb_articles (
    id TEXT PRIMARY KEY,
    tenant_id TEXT NOT NULL,

    -- Content
    title TEXT NOT NULL,
    slug TEXT NOT NULL,
    summary TEXT NOT NULL DEFAULT '',
    content TEXT NOT NULL,
    content_type TEXT NOT NULL DEFAULT 'runbook',

    -- Organization
    module TEXT NOT NULL DEFAULT 'general',
    lifecycle_stage TEXT NOT NULL DEFAULT 'support',
    tags TEXT[] DEFAULT '{}',

    -- Version tracking
    version INT NOT NULL DEFAULT 1,

    -- Status
    status TEXT NOT NULL DEFAULT 'draft',

    -- Authorship
    created_by_id TEXT NOT NULL,
    created_by_name TEXT NOT NULL DEFAULT '',
    updated_by_id TEXT NOT NULL,
    updated_by_name TEXT NOT NULL DEFAULT '',

    -- Timestamps
    published_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,

    -- Unique slug per tenant
    UNIQUE (tenant_id, slug)
);

-- Indexes for common queries
CREATE INDEX idx_kb_articles_tenant_status ON kb_articles(tenant_id, status, updated_at DESC);
CREATE INDEX idx_kb_articles_module ON kb_articles(tenant_id, module, status);
CREATE INDEX idx_kb_articles_lifecycle ON kb_articles(tenant_id, lifecycle_stage, status);
CREATE INDEX idx_kb_articles_content_type ON kb_articles(tenant_id, content_type, status);
CREATE INDEX idx_kb_articles_tags ON kb_articles USING gin(tags);

-- Full-text search index
CREATE INDEX idx_kb_articles_search ON kb_articles USING gin(to_tsvector('english', title || ' ' || summary || ' ' || content));

-- +goose Down
DROP TABLE IF EXISTS kb_articles;
