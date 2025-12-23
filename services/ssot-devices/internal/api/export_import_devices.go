package api

import (
	"context"
	"errors"
	"time"

	"github.com/edvirons/ssp/ssot_devices/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func exportAll(ctx context.Context, db *pgxpool.Pool, tenant string) (models.ExportPayload, error) {
	p := models.ExportPayload{Version: "2", GeneratedAt: time.Now().UTC()}
	rows, err := db.Query(ctx, `SELECT id, tenant_id, make, model, category, spec_json, created_at, updated_at FROM device_models WHERE tenant_id=$1 ORDER BY make, model`, tenant)
	if err != nil {
		return p, err
	}
	for rows.Next() {
		var x models.DeviceModel
		if err := rows.Scan(&x.ID, &x.TenantID, &x.Make, &x.Model, &x.Category, &x.SpecJSON, &x.CreatedAt, &x.UpdatedAt); err != nil {
			rows.Close()
			return p, err
		}
		p.Models = append(p.Models, x)
	}
	rows.Close()

	rows, err = db.Query(ctx, `SELECT id, tenant_id, serial, asset_tag, device_model_id, school_id, assigned_to, lifecycle, enrolled, created_at, updated_at FROM devices WHERE tenant_id=$1 ORDER BY created_at DESC`, tenant)
	if err != nil {
		return p, err
	}
	for rows.Next() {
		var x models.Device
		if err := rows.Scan(&x.ID, &x.TenantID, &x.Serial, &x.AssetTag, &x.DeviceModelID, &x.SchoolID, &x.AssignedTo, &x.Lifecycle, &x.Enrolled, &x.CreatedAt, &x.UpdatedAt); err != nil {
			rows.Close()
			return p, err
		}
		p.Devices = append(p.Devices, x)
	}
	rows.Close()

	// Export network identities (MAC addresses)
	rows, err = db.Query(ctx, `SELECT id, tenant_id, device_id, mac_address, interface_name, interface_type, is_primary, first_seen_at, last_seen_at, created_at, updated_at FROM device_network_identities WHERE tenant_id=$1 ORDER BY device_id, is_primary DESC`, tenant)
	if err != nil {
		return p, err
	}
	for rows.Next() {
		var x models.DeviceNetworkIdentity
		var ifType string
		if err := rows.Scan(&x.ID, &x.TenantID, &x.DeviceID, &x.MACAddress, &x.InterfaceName, &ifType, &x.IsPrimary, &x.FirstSeenAt, &x.LastSeenAt, &x.CreatedAt, &x.UpdatedAt); err != nil {
			rows.Close()
			return p, err
		}
		x.InterfaceType = models.InterfaceType(ifType)
		p.NetworkIdentities = append(p.NetworkIdentities, x)
	}
	rows.Close()
	return p, nil
}

func importAll(ctx context.Context, db *pgxpool.Pool, tenant string, body map[string]any) (map[string]any, error) {
	modelsList, _ := body["models"].([]any)
	devices, _ := body["devices"].([]any)
	networkIds, _ := body["networkIdentities"].([]any)
	if len(modelsList) == 0 && len(devices) == 0 && len(networkIds) == 0 {
		return nil, errors.New("no ssot data provided")
	}
	res := map[string]any{"models": 0, "devices": 0, "networkIdentities": 0}

	err := withTx(ctx, db, func(tx pgx.Tx) error {
		now := time.Now().UTC()

		for _, it := range modelsList {
			m, ok := it.(map[string]any)
			if !ok {
				continue
			}
			id := trim(m["id"])
			if id == "" {
				id = newID("dmodel")
			}
			make := trim(m["make"])
			model := trim(m["model"])
			if make == "" || model == "" {
				continue
			}
			cat := trim(m["category"])
			if cat == "" {
				cat = "other"
			}
			spec := trim(m["specJson"])
			if spec == "" {
				spec = "{}"
			}
			_, err := tx.Exec(ctx, `
			INSERT INTO device_models (id, tenant_id, make, model, category, spec_json, created_at, updated_at)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$7)
			ON CONFLICT (id) DO UPDATE SET make=EXCLUDED.make, model=EXCLUDED.model, category=EXCLUDED.category, spec_json=EXCLUDED.spec_json, updated_at=$7
		`, id, tenant, make, model, cat, spec, now)
			if err != nil {
				return err
			}
			if c, ok := res["models"].(int); ok {
				res["models"] = c + 1
			}
		}

		for _, it := range devices {
			m, ok := it.(map[string]any)
			if !ok {
				continue
			}
			id := trim(m["id"])
			if id == "" {
				id = newID("dev")
			}
			serial := trim(m["serial"])
			asset := trim(m["assetTag"])
			dmid := trim(m["deviceModelId"])
			school := trim(m["schoolId"])
			assigned := trim(m["assignedTo"])
			life := trim(m["lifecycle"])
			if life == "" {
				life = "in_stock"
			}
			enrolled := false
			if v, ok := m["enrolled"].(bool); ok {
				enrolled = v
			}
			_, err := tx.Exec(ctx, `
			INSERT INTO devices (id, tenant_id, serial, asset_tag, device_model_id, school_id, assigned_to, lifecycle, enrolled, created_at, updated_at)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$10)
			ON CONFLICT (id) DO UPDATE SET serial=EXCLUDED.serial, asset_tag=EXCLUDED.asset_tag, device_model_id=EXCLUDED.device_model_id,
				school_id=EXCLUDED.school_id, assigned_to=EXCLUDED.assigned_to, lifecycle=EXCLUDED.lifecycle, enrolled=EXCLUDED.enrolled, updated_at=$10
		`, id, tenant, serial, asset, dmid, school, assigned, life, enrolled, now)
			if err != nil {
				return err
			}
			if c, ok := res["devices"].(int); ok {
				res["devices"] = c + 1
			}
		}

		// Import network identities (MAC addresses)
		for _, it := range networkIds {
			m, ok := it.(map[string]any)
			if !ok {
				continue
			}
			id := trim(m["id"])
			if id == "" {
				id = newID("netid")
			}
			deviceID := trim(m["deviceId"])
			mac := trim(m["macAddress"])
			if deviceID == "" || mac == "" {
				continue
			}
			mac = models.NormalizeMACAddress(mac)
			ifName := trim(m["interfaceName"])
			ifType := trim(m["interfaceType"])
			if ifType == "" {
				ifType = "unknown"
			}
			isPrimary := false
			if v, ok := m["isPrimary"].(bool); ok {
				isPrimary = v
			}
			_, err := tx.Exec(ctx, `
			INSERT INTO device_network_identities (id, tenant_id, device_id, mac_address, interface_name, interface_type, is_primary, first_seen_at, last_seen_at, created_at, updated_at)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$8,$8,$8)
			ON CONFLICT (tenant_id, mac_address) DO UPDATE SET device_id=EXCLUDED.device_id, interface_name=EXCLUDED.interface_name, interface_type=EXCLUDED.interface_type, is_primary=EXCLUDED.is_primary, last_seen_at=$8, updated_at=$8
		`, id, tenant, deviceID, mac, ifName, ifType, isPrimary, now)
			if err != nil {
				return err
			}
			if c, ok := res["networkIdentities"].(int); ok {
				res["networkIdentities"] = c + 1
			}
		}
		return nil
	})
	return res, err
}
