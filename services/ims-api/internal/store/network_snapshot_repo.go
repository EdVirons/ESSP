package store

import (
	"context"
	"time"

	"github.com/edvirons/ssp/ims/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NetworkSnapshotRepo struct{ pool *pgxpool.Pool }

// Upsert updates or inserts a network snapshot entry
func (r *NetworkSnapshotRepo) Upsert(ctx context.Context, snap models.DeviceNetworkSnapshot) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO device_network_snapshot (tenant_id, device_id, mac_address, interface_type, is_primary, last_seen_at, synced_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
		ON CONFLICT (tenant_id, device_id, mac_address) DO UPDATE SET
			interface_type=EXCLUDED.interface_type, is_primary=EXCLUDED.is_primary, last_seen_at=EXCLUDED.last_seen_at, synced_at=EXCLUDED.synced_at
	`, snap.TenantID, snap.DeviceID, snap.MACAddress, snap.InterfaceType, snap.IsPrimary, snap.LastSeenAt, snap.SyncedAt)
	return err
}

// Delete removes a network snapshot entry
func (r *NetworkSnapshotRepo) Delete(ctx context.Context, tenantID, deviceID, macAddress string) error {
	_, err := r.pool.Exec(ctx, `
		DELETE FROM device_network_snapshot WHERE tenant_id=$1 AND device_id=$2 AND mac_address=$3
	`, tenantID, deviceID, macAddress)
	return err
}

// DeleteByDevice removes all network snapshot entries for a device
func (r *NetworkSnapshotRepo) DeleteByDevice(ctx context.Context, tenantID, deviceID string) error {
	_, err := r.pool.Exec(ctx, `
		DELETE FROM device_network_snapshot WHERE tenant_id=$1 AND device_id=$2
	`, tenantID, deviceID)
	return err
}

// ListByDevice returns all MAC addresses for a device
func (r *NetworkSnapshotRepo) ListByDevice(ctx context.Context, tenantID, deviceID string) ([]models.DeviceNetworkSnapshot, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT tenant_id, device_id, mac_address, interface_type, is_primary, last_seen_at, synced_at
		FROM device_network_snapshot
		WHERE tenant_id=$1 AND device_id=$2
		ORDER BY is_primary DESC, mac_address
	`, tenantID, deviceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []models.DeviceNetworkSnapshot{}
	for rows.Next() {
		var s models.DeviceNetworkSnapshot
		if err := rows.Scan(&s.TenantID, &s.DeviceID, &s.MACAddress, &s.InterfaceType, &s.IsPrimary, &s.LastSeenAt, &s.SyncedAt); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, nil
}

// LookupByMAC finds a device by MAC address
func (r *NetworkSnapshotRepo) LookupByMAC(ctx context.Context, tenantID, macAddress string) (models.DeviceNetworkSnapshot, error) {
	var s models.DeviceNetworkSnapshot
	err := r.pool.QueryRow(ctx, `
		SELECT tenant_id, device_id, mac_address, interface_type, is_primary, last_seen_at, synced_at
		FROM device_network_snapshot
		WHERE tenant_id=$1 AND mac_address=$2
	`, tenantID, macAddress).Scan(&s.TenantID, &s.DeviceID, &s.MACAddress, &s.InterfaceType, &s.IsPrimary, &s.LastSeenAt, &s.SyncedAt)
	return s, err
}

// GetMACAddressesForDevices returns MAC addresses for multiple devices (for inventory view)
func (r *NetworkSnapshotRepo) GetMACAddressesForDevices(ctx context.Context, tenantID string, deviceIDs []string) (map[string][]string, error) {
	if len(deviceIDs) == 0 {
		return make(map[string][]string), nil
	}

	rows, err := r.pool.Query(ctx, `
		SELECT device_id, mac_address
		FROM device_network_snapshot
		WHERE tenant_id=$1 AND device_id = ANY($2)
		ORDER BY device_id, is_primary DESC, mac_address
	`, tenantID, deviceIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string][]string)
	for rows.Next() {
		var deviceID, mac string
		if err := rows.Scan(&deviceID, &mac); err != nil {
			return nil, err
		}
		result[deviceID] = append(result[deviceID], mac)
	}
	return result, nil
}

// SyncFromSSOT replaces all MAC addresses for a device with new data from ssot-devices
func (r *NetworkSnapshotRepo) SyncFromSSOT(ctx context.Context, tenantID, deviceID string, macs []models.DeviceNetworkSnapshot) error {
	now := time.Now().UTC()

	// Delete existing
	if _, err := r.pool.Exec(ctx, `DELETE FROM device_network_snapshot WHERE tenant_id=$1 AND device_id=$2`, tenantID, deviceID); err != nil {
		return err
	}

	// Insert new
	for _, mac := range macs {
		mac.TenantID = tenantID
		mac.DeviceID = deviceID
		mac.SyncedAt = now
		if err := r.Upsert(ctx, mac); err != nil {
			return err
		}
	}
	return nil
}
