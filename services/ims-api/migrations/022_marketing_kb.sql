-- +goose Up

-- Marketing KB Articles
CREATE TABLE IF NOT EXISTS marketing_kb_articles (
    id TEXT PRIMARY KEY,
    tenant_id TEXT NOT NULL,
    title TEXT NOT NULL,
    slug TEXT NOT NULL,
    summary TEXT NOT NULL DEFAULT '',
    content TEXT NOT NULL,
    content_type TEXT NOT NULL DEFAULT 'messaging',
    personas TEXT[] NOT NULL DEFAULT '{}',
    context_tags TEXT[] NOT NULL DEFAULT '{}',
    tags TEXT[] DEFAULT '{}',
    version INT NOT NULL DEFAULT 1,
    status TEXT NOT NULL DEFAULT 'draft',
    usage_count INT NOT NULL DEFAULT 0,
    last_used_at TIMESTAMPTZ,
    created_by_id TEXT NOT NULL,
    created_by_name TEXT NOT NULL DEFAULT '',
    updated_by_id TEXT NOT NULL,
    updated_by_name TEXT NOT NULL DEFAULT '',
    approved_at TIMESTAMPTZ,
    approved_by_id TEXT,
    approved_by_name TEXT,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    UNIQUE (tenant_id, slug)
);

CREATE INDEX idx_mkb_tenant_status ON marketing_kb_articles(tenant_id, status, updated_at DESC);
CREATE INDEX idx_mkb_content_type ON marketing_kb_articles(tenant_id, content_type, status);
CREATE INDEX idx_mkb_personas ON marketing_kb_articles USING gin(personas);
CREATE INDEX idx_mkb_context_tags ON marketing_kb_articles USING gin(context_tags);
CREATE INDEX idx_mkb_search ON marketing_kb_articles USING gin(to_tsvector('english', title || ' ' || summary || ' ' || content));

-- Pitch Kits (saved collections of content blocks)
CREATE TABLE IF NOT EXISTS marketing_pitch_kits (
    id TEXT PRIMARY KEY,
    tenant_id TEXT NOT NULL,
    name TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    target_persona TEXT NOT NULL DEFAULT 'director',
    context_tags TEXT[] NOT NULL DEFAULT '{}',
    article_ids TEXT[] NOT NULL DEFAULT '{}',
    is_template BOOLEAN NOT NULL DEFAULT false,
    created_by_id TEXT NOT NULL,
    created_by_name TEXT NOT NULL DEFAULT '',
    updated_by_id TEXT NOT NULL,
    updated_by_name TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_pitch_kits_tenant ON marketing_pitch_kits(tenant_id, updated_at DESC);
CREATE INDEX idx_pitch_kits_persona ON marketing_pitch_kits(tenant_id, target_persona);

-- +goose Down
DROP TABLE IF EXISTS marketing_pitch_kits;
DROP TABLE IF EXISTS marketing_kb_articles;
