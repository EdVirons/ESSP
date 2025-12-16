-- +goose Up
-- Demo pipeline for sales lead tracking

CREATE TABLE IF NOT EXISTS demo_leads (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL,

  -- School/Lead Info
  school_id TEXT,  -- NULL if lead is not yet a school in system
  school_name TEXT NOT NULL,
  contact_name TEXT,
  contact_email TEXT,
  contact_phone TEXT,
  contact_role TEXT,

  -- Pipeline Info
  stage TEXT NOT NULL DEFAULT 'new_lead',  -- new_lead, contacted, demo_scheduled, demo_completed, proposal_sent, negotiation, won, lost
  stage_changed_at TIMESTAMPTZ DEFAULT NOW(),

  -- Deal Info
  estimated_value DECIMAL(12, 2),
  estimated_devices INTEGER,
  probability INTEGER DEFAULT 0,  -- 0-100
  expected_close_date DATE,

  -- Source & Attribution
  lead_source TEXT,  -- website, referral, event, cold_outreach, inbound
  assigned_to TEXT,  -- user_id of sales rep

  -- Notes & Details
  notes TEXT,
  tags TEXT[] DEFAULT '{}',

  -- Lost reason (if stage = lost)
  lost_reason TEXT,
  lost_notes TEXT,

  -- Metadata
  created_by TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_demo_leads_tenant ON demo_leads(tenant_id);
CREATE INDEX idx_demo_leads_stage ON demo_leads(tenant_id, stage);
CREATE INDEX idx_demo_leads_assigned ON demo_leads(tenant_id, assigned_to);
CREATE INDEX idx_demo_leads_school ON demo_leads(tenant_id, school_id) WHERE school_id IS NOT NULL;

-- Activity tracking for leads
CREATE TABLE IF NOT EXISTS demo_lead_activities (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL,
  lead_id TEXT NOT NULL REFERENCES demo_leads(id) ON DELETE CASCADE,

  activity_type TEXT NOT NULL,  -- note, call, email, meeting, demo, stage_change, created, updated
  description TEXT,

  -- For stage changes
  from_stage TEXT,
  to_stage TEXT,

  -- For scheduled activities
  scheduled_at TIMESTAMPTZ,
  completed_at TIMESTAMPTZ,

  created_by TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_demo_lead_activities_lead ON demo_lead_activities(lead_id, created_at DESC);
CREATE INDEX idx_demo_lead_activities_tenant ON demo_lead_activities(tenant_id, created_at DESC);

-- Demo scheduling
CREATE TABLE IF NOT EXISTS demo_schedules (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL,
  lead_id TEXT NOT NULL REFERENCES demo_leads(id) ON DELETE CASCADE,

  scheduled_date DATE NOT NULL,
  scheduled_time TIME,
  duration_minutes INTEGER DEFAULT 60,

  location TEXT,  -- 'virtual', 'on-site', or address
  meeting_link TEXT,

  attendees JSONB DEFAULT '[]',  -- [{name, email, role}]

  status TEXT DEFAULT 'scheduled',  -- scheduled, completed, cancelled, rescheduled
  outcome TEXT,  -- positive, neutral, negative
  outcome_notes TEXT,

  reminder_sent BOOLEAN DEFAULT FALSE,

  created_by TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_demo_schedules_lead ON demo_schedules(lead_id);
CREATE INDEX idx_demo_schedules_date ON demo_schedules(tenant_id, scheduled_date);

-- +goose Down
DROP TABLE IF EXISTS demo_schedules;
DROP TABLE IF EXISTS demo_lead_activities;
DROP TABLE IF EXISTS demo_leads;
