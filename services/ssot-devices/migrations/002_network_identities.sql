-- Device Network Identities (MAC addresses)
-- A device can have multiple NICs, each with a MAC address

CREATE TABLE IF NOT EXISTS device_network_identities (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL,
  device_id TEXT NOT NULL REFERENCES devices(id) ON DELETE CASCADE,
  mac_address TEXT NOT NULL,                      -- Normalized: lowercase, colon-separated (aa:bb:cc:dd:ee:ff)
  interface_name TEXT NOT NULL DEFAULT '',        -- e.g., 'eth0', 'wlan0', 'en0'
  interface_type TEXT NOT NULL DEFAULT 'unknown', -- 'ethernet', 'wifi', 'bluetooth', 'unknown'
  is_primary BOOLEAN NOT NULL DEFAULT FALSE,
  first_seen_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  last_seen_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Unique MAC per tenant (MAC addresses should be globally unique but multi-tenant safety)
CREATE UNIQUE INDEX IF NOT EXISTS ux_device_mac ON device_network_identities(tenant_id, mac_address);

-- Find all MACs for a device
CREATE INDEX IF NOT EXISTS idx_device_network_device ON device_network_identities(tenant_id, device_id);

-- Find device by MAC (network discovery use case)
CREATE INDEX IF NOT EXISTS idx_device_network_mac ON device_network_identities(mac_address);
