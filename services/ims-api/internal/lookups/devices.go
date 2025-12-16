package lookups

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// Device export types

type DevicesExport struct {
	Version     string        `json:"version"`
	GeneratedAt time.Time     `json:"generatedAt"`
	Models      []DeviceModel `json:"models"`
	Devices     []Device      `json:"devices"`
}

type DeviceModel struct {
	ID       string `json:"id"`
	Make     string `json:"make"`
	Model    string `json:"model"`
	Category string `json:"category"`
	SpecJSON string `json:"specJson"`
}

type Device struct {
	ID            string `json:"id"`
	Serial        string `json:"serial"`
	AssetTag      string `json:"assetTag"`
	DeviceModelID string `json:"deviceModelId"`
	SchoolID      string `json:"schoolId"`
	AssignedTo    string `json:"assignedTo"`
	Lifecycle     string `json:"lifecycle"`
	Enrolled      bool   `json:"enrolled"`
}

// Summary types for lookup results

type DeviceSummary struct {
	ID        string
	Serial    string
	AssetTag  string
	SchoolID  string
	ModelID   string
	Make      string
	Model     string
	Category  string
	Lifecycle string
	Enrolled  bool
}

// LoadDevicesExport loads the devices export snapshot
func (s *Store) LoadDevicesExport(ctx context.Context, tenant string) (*DevicesExport, error) {
	snap, err := s.GetSnapshot(ctx, tenant, KindDevices)
	if err != nil {
		return nil, err
	}
	var ex DevicesExport
	if err := json.Unmarshal(snap.Payload, &ex); err != nil {
		return nil, fmt.Errorf("unmarshal devices export: %w", err)
	}
	return &ex, nil
}

// DeviceByID looks up a device by its ID
func (s *Store) DeviceByID(ctx context.Context, tenant, deviceID string) (*DeviceSummary, error) {
	ex, err := s.LoadDevicesExport(ctx, tenant)
	if err != nil {
		return nil, ErrSnapshotMissing
	}

	modelByID := map[string]DeviceModel{}
	for _, m := range ex.Models {
		modelByID[m.ID] = m
	}

	for _, d := range ex.Devices {
		if d.ID == deviceID {
			m := modelByID[d.DeviceModelID]
			return &DeviceSummary{
				ID: d.ID, Serial: d.Serial, AssetTag: d.AssetTag,
				SchoolID: d.SchoolID, ModelID: d.DeviceModelID,
				Make: m.Make, Model: m.Model, Category: m.Category,
				Lifecycle: d.Lifecycle, Enrolled: d.Enrolled,
			}, nil
		}
	}
	return nil, ErrNotFound
}

// DeviceBySerial looks up a device by its serial number
func (s *Store) DeviceBySerial(ctx context.Context, tenant, serial string) (*DeviceSummary, error) {
	ex, err := s.LoadDevicesExport(ctx, tenant)
	if err != nil {
		return nil, ErrSnapshotMissing
	}

	modelByID := map[string]DeviceModel{}
	for _, m := range ex.Models {
		modelByID[m.ID] = m
	}

	for _, d := range ex.Devices {
		if d.Serial != "" && d.Serial == serial {
			m := modelByID[d.DeviceModelID]
			return &DeviceSummary{
				ID: d.ID, Serial: d.Serial, AssetTag: d.AssetTag,
				SchoolID: d.SchoolID, ModelID: d.DeviceModelID,
				Make: m.Make, Model: m.Model, Category: m.Category,
				Lifecycle: d.Lifecycle, Enrolled: d.Enrolled,
			}, nil
		}
	}
	return nil, ErrNotFound
}
