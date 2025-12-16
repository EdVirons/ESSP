-- +goose Up
-- EdTech profile assessment for schools

CREATE TABLE IF NOT EXISTS edtech_profiles (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL,
  school_id TEXT NOT NULL,

  -- Infrastructure Section
  total_devices INTEGER DEFAULT 0,
  device_types JSONB DEFAULT '{}',
  network_quality TEXT,
  internet_speed TEXT,
  lms_platform TEXT,
  existing_software JSONB DEFAULT '[]',
  it_staff_count INTEGER DEFAULT 0,
  device_age TEXT,

  -- Pain Points Section
  pain_points JSONB DEFAULT '[]',
  support_satisfaction INTEGER,
  biggest_challenges TEXT[] DEFAULT '{}',
  support_frequency TEXT,
  avg_resolution_time TEXT,
  biggest_frustration TEXT,
  wish_list TEXT,

  -- Goals Section
  strategic_goals TEXT[] DEFAULT '{}',
  budget_range TEXT,
  timeline TEXT,
  expansion_plans TEXT,
  priority_ranking JSONB DEFAULT '[]',
  decision_makers TEXT[] DEFAULT '{}',

  -- AI Section
  ai_summary TEXT,
  ai_recommendations JSONB DEFAULT '[]',
  follow_up_questions JSONB DEFAULT '[]',
  follow_up_responses JSONB DEFAULT '{}',

  -- Metadata
  status TEXT DEFAULT 'draft',
  completed_at TIMESTAMPTZ,
  completed_by TEXT,
  version INTEGER DEFAULT 1,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  UNIQUE (tenant_id, school_id)
);

CREATE INDEX idx_edtech_profiles_tenant_school ON edtech_profiles(tenant_id, school_id);
CREATE INDEX idx_edtech_profiles_status ON edtech_profiles(tenant_id, status);

-- Profile version history for tracking changes
CREATE TABLE IF NOT EXISTS edtech_profile_history (
  id TEXT PRIMARY KEY,
  profile_id TEXT NOT NULL REFERENCES edtech_profiles(id) ON DELETE CASCADE,
  snapshot JSONB NOT NULL,
  changed_by TEXT NOT NULL,
  change_reason TEXT,
  changed_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_edtech_profile_history_profile ON edtech_profile_history(profile_id, changed_at DESC);

-- +goose Down
DROP TABLE IF EXISTS edtech_profile_history;
DROP TABLE IF EXISTS edtech_profiles;
