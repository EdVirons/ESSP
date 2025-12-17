-- +goose Up
-- Device Inventory Management
-- Supports school-contact inventory views, device assignments, and grouping

-- ============================================
-- Locations (Hierarchical school spaces)
-- ============================================
CREATE TABLE IF NOT EXISTS locations (
    id TEXT PRIMARY KEY,
    tenant_id TEXT NOT NULL,
    school_id TEXT NOT NULL,
    parent_id TEXT REFERENCES locations(id) ON DELETE SET NULL,
    name TEXT NOT NULL,
    location_type TEXT NOT NULL DEFAULT 'room',  -- block, floor, room, lab, storage, office
    code TEXT NOT NULL DEFAULT '',               -- Short code: B1-R101, LAB-A
    capacity INTEGER NOT NULL DEFAULT 0,
    metadata JSONB NOT NULL DEFAULT '{}',        -- Extra attributes (contact, phone, etc)
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Find locations by school
CREATE INDEX IF NOT EXISTS idx_locations_school ON locations(tenant_id, school_id, active);

-- Unique code per school (when code is set)
CREATE UNIQUE INDEX IF NOT EXISTS ux_locations_code ON locations(tenant_id, school_id, code) WHERE code != '';

-- Find children of a parent location
CREATE INDEX IF NOT EXISTS idx_locations_parent ON locations(parent_id) WHERE parent_id IS NOT NULL;

-- ============================================
-- Device Assignments (Temporal tracking)
-- ============================================
CREATE TABLE IF NOT EXISTS device_assignments (
    id TEXT PRIMARY KEY,
    tenant_id TEXT NOT NULL,
    device_id TEXT NOT NULL,                     -- References device in ssot-devices
    location_id TEXT REFERENCES locations(id) ON DELETE SET NULL,
    assigned_to_user TEXT NOT NULL DEFAULT '',   -- User/student ID if assigned to person
    assignment_type TEXT NOT NULL DEFAULT 'permanent', -- permanent, temporary, loan, repair, storage
    effective_from TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    effective_to TIMESTAMPTZ,                    -- NULL = current assignment
    notes TEXT NOT NULL DEFAULT '',
    created_by TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Find current assignment for a device
CREATE INDEX IF NOT EXISTS idx_assignments_device_current ON device_assignments(tenant_id, device_id)
    WHERE effective_to IS NULL;

-- Find all devices at a location (current assignments only)
CREATE INDEX IF NOT EXISTS idx_assignments_location ON device_assignments(tenant_id, location_id)
    WHERE effective_to IS NULL;

-- Find assignments by type
CREATE INDEX IF NOT EXISTS idx_assignments_type ON device_assignments(tenant_id, assignment_type, effective_to);

-- ============================================
-- Device Groups (For control & policies)
-- ============================================
CREATE TABLE IF NOT EXISTS device_groups (
    id TEXT PRIMARY KEY,
    tenant_id TEXT NOT NULL,
    school_id TEXT,                              -- NULL = tenant-wide group
    name TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    group_type TEXT NOT NULL DEFAULT 'manual',   -- manual, location, dynamic
    location_id TEXT REFERENCES locations(id) ON DELETE SET NULL,  -- For location-based groups
    selector JSONB,                              -- For dynamic groups: {"model": "Chromebook", "lifecycle": ["active"]}
    policies JSONB NOT NULL DEFAULT '{}',        -- Future: exam_mode, restrictions, etc
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Find groups by school
CREATE INDEX IF NOT EXISTS idx_groups_school ON device_groups(tenant_id, school_id, active);

-- Find location-based groups
CREATE INDEX IF NOT EXISTS idx_groups_location ON device_groups(location_id) WHERE location_id IS NOT NULL;

-- ============================================
-- Group Members (Manual group membership)
-- ============================================
CREATE TABLE IF NOT EXISTS group_members (
    id TEXT PRIMARY KEY,
    tenant_id TEXT NOT NULL,
    group_id TEXT NOT NULL REFERENCES device_groups(id) ON DELETE CASCADE,
    device_id TEXT NOT NULL,
    added_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    added_by TEXT NOT NULL DEFAULT '',
    UNIQUE(group_id, device_id)
);

-- Find all members of a group
CREATE INDEX IF NOT EXISTS idx_group_members_group ON group_members(group_id);

-- Find all groups a device belongs to
CREATE INDEX IF NOT EXISTS idx_group_members_device ON group_members(tenant_id, device_id);

-- ============================================
-- Device Network Snapshot (Cache from ssot-devices)
-- ============================================
CREATE TABLE IF NOT EXISTS device_network_snapshot (
    tenant_id TEXT NOT NULL,
    device_id TEXT NOT NULL,
    mac_address TEXT NOT NULL,
    interface_type TEXT NOT NULL DEFAULT 'unknown',  -- ethernet, wifi, bluetooth
    is_primary BOOLEAN NOT NULL DEFAULT FALSE,
    last_seen_at TIMESTAMPTZ,
    synced_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (tenant_id, device_id, mac_address)
);

-- Find device by MAC address
CREATE INDEX IF NOT EXISTS idx_network_snapshot_mac ON device_network_snapshot(tenant_id, mac_address);

-- ============================================
-- Helper function for location path
-- ============================================
CREATE OR REPLACE FUNCTION get_location_path(loc_id TEXT)
RETURNS TEXT AS $$
DECLARE
    result TEXT;
BEGIN
    WITH RECURSIVE path AS (
        SELECT id, name, parent_id, name::TEXT as full_path
        FROM locations WHERE id = loc_id
        UNION ALL
        SELECT l.id, l.name, l.parent_id, l.name || ' > ' || p.full_path
        FROM locations l
        JOIN path p ON l.id = p.parent_id
    )
    SELECT full_path INTO result FROM path WHERE parent_id IS NULL;
    RETURN COALESCE(result, '');
END;
$$ LANGUAGE plpgsql STABLE;
