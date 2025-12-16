-- +goose Up
-- Sales presentations and marketing materials

CREATE TABLE IF NOT EXISTS presentations (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL,

  -- Content Info
  title TEXT NOT NULL,
  description TEXT,
  type TEXT NOT NULL,  -- presentation, brochure, case_study, video, roi_calculator, template, other
  category TEXT,  -- general, product_overview, technical, pricing, onboarding, training

  -- File Storage (MinIO)
  file_key TEXT,  -- MinIO object key for the file
  file_name TEXT,
  file_size BIGINT,
  file_type TEXT,  -- MIME type: application/pdf, video/mp4, etc.

  -- Preview
  thumbnail_key TEXT,  -- MinIO key for thumbnail/preview image
  preview_type TEXT,  -- image, pdf_preview, video_thumbnail

  -- Metadata
  tags TEXT[] DEFAULT '{}',
  version INTEGER DEFAULT 1,
  is_active BOOLEAN DEFAULT TRUE,
  is_featured BOOLEAN DEFAULT FALSE,

  -- Usage Stats
  view_count INTEGER DEFAULT 0,
  download_count INTEGER DEFAULT 0,
  last_viewed_at TIMESTAMPTZ,

  -- Audit
  created_by TEXT NOT NULL,
  updated_by TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_presentations_tenant ON presentations(tenant_id);
CREATE INDEX idx_presentations_type ON presentations(tenant_id, type);
CREATE INDEX idx_presentations_category ON presentations(tenant_id, category);
CREATE INDEX idx_presentations_active ON presentations(tenant_id, is_active);
CREATE INDEX idx_presentations_featured ON presentations(tenant_id, is_featured) WHERE is_featured = TRUE;

-- Version history for presentations
CREATE TABLE IF NOT EXISTS presentation_versions (
  id TEXT PRIMARY KEY,
  presentation_id TEXT NOT NULL REFERENCES presentations(id) ON DELETE CASCADE,

  version INTEGER NOT NULL,
  file_key TEXT NOT NULL,
  file_name TEXT,
  file_size BIGINT,

  change_notes TEXT,

  created_by TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_presentation_versions_pres ON presentation_versions(presentation_id, version DESC);

-- Track who viewed what
CREATE TABLE IF NOT EXISTS presentation_views (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL,
  presentation_id TEXT NOT NULL REFERENCES presentations(id) ON DELETE CASCADE,

  viewed_by TEXT NOT NULL,  -- user_id
  viewed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  -- Optional context
  context TEXT,  -- where they viewed from: 'dashboard', 'search', 'share_link'
  duration_seconds INTEGER  -- how long they viewed (for videos/PDFs)
);

CREATE INDEX idx_presentation_views_pres ON presentation_views(presentation_id, viewed_at DESC);
CREATE INDEX idx_presentation_views_user ON presentation_views(tenant_id, viewed_by, viewed_at DESC);

-- Sales metrics aggregation (for dashboard)
CREATE TABLE IF NOT EXISTS sales_metrics_daily (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL,
  metric_date DATE NOT NULL,

  -- Pipeline metrics
  new_leads INTEGER DEFAULT 0,
  leads_contacted INTEGER DEFAULT 0,
  demos_scheduled INTEGER DEFAULT 0,
  demos_completed INTEGER DEFAULT 0,
  proposals_sent INTEGER DEFAULT 0,
  deals_won INTEGER DEFAULT 0,
  deals_lost INTEGER DEFAULT 0,

  -- Value metrics
  pipeline_value DECIMAL(14, 2) DEFAULT 0,
  won_value DECIMAL(14, 2) DEFAULT 0,
  lost_value DECIMAL(14, 2) DEFAULT 0,

  -- Activity metrics
  calls_made INTEGER DEFAULT 0,
  emails_sent INTEGER DEFAULT 0,
  meetings_held INTEGER DEFAULT 0,

  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  UNIQUE (tenant_id, metric_date)
);

CREATE INDEX idx_sales_metrics_daily_tenant_date ON sales_metrics_daily(tenant_id, metric_date DESC);

-- +goose Down
DROP TABLE IF EXISTS sales_metrics_daily;
DROP TABLE IF EXISTS presentation_views;
DROP TABLE IF EXISTS presentation_versions;
DROP TABLE IF EXISTS presentations;
