package lookups

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// Parts export types

type PartsExport struct {
	Version       string              `json:"version"`
	GeneratedAt   time.Time           `json:"generatedAt"`
	Parts         []Part              `json:"parts"`
	Compatibility []PartCompatibility `json:"compatibility"`
	VendorSKUs    []VendorSKU         `json:"vendorSkus"`
}

type Part struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
	PUK      string `json:"puk"`
	SpecJSON string `json:"specJson"`
}

type PartCompatibility struct {
	ID            string `json:"id"`
	PartID        string `json:"partId"`
	DeviceModelID string `json:"deviceModelId"`
}

type VendorSKU struct {
	ID             string `json:"id"`
	PartID         string `json:"partId"`
	VendorID       string `json:"vendorId"`
	SKU            string `json:"sku"`
	UnitPriceCents int64  `json:"unitPriceCents"`
	Currency       string `json:"currency"`
	LeadTimeDays   int    `json:"leadTimeDays"`
}

// Summary types for lookup results

type PartSummary struct {
	ID       string
	Name     string
	Category string
	PUK      string
}

// LoadPartsExport loads the parts export snapshot
func (s *Store) LoadPartsExport(ctx context.Context, tenant string) (*PartsExport, error) {
	snap, err := s.GetSnapshot(ctx, tenant, KindParts)
	if err != nil {
		return nil, err
	}
	var ex PartsExport
	if err := json.Unmarshal(snap.Payload, &ex); err != nil {
		return nil, fmt.Errorf("unmarshal parts export: %w", err)
	}
	return &ex, nil
}

// PartByID looks up a part by its ID
func (s *Store) PartByID(ctx context.Context, tenant, partID string) (*PartSummary, error) {
	ex, err := s.LoadPartsExport(ctx, tenant)
	if err != nil {
		return nil, ErrSnapshotMissing
	}
	for _, p := range ex.Parts {
		if p.ID == partID {
			return &PartSummary{ID: p.ID, Name: p.Name, Category: p.Category, PUK: p.PUK}, nil
		}
	}
	return nil, ErrNotFound
}

// PartByPUK looks up a part by its PUK
func (s *Store) PartByPUK(ctx context.Context, tenant, puk string) (*PartSummary, error) {
	ex, err := s.LoadPartsExport(ctx, tenant)
	if err != nil {
		return nil, ErrSnapshotMissing
	}
	for _, p := range ex.Parts {
		if p.PUK != "" && p.PUK == puk {
			return &PartSummary{ID: p.ID, Name: p.Name, Category: p.Category, PUK: p.PUK}, nil
		}
	}
	return nil, ErrNotFound
}

// IsPartCompatibleWithDeviceModel checks if a part is compatible with a device model
func (s *Store) IsPartCompatibleWithDeviceModel(ctx context.Context, tenant, partID, deviceModelID string) (bool, error) {
	ex, err := s.LoadPartsExport(ctx, tenant)
	if err != nil {
		return false, ErrSnapshotMissing
	}
	for _, c := range ex.Compatibility {
		if c.PartID == partID && c.DeviceModelID == deviceModelID {
			return true, nil
		}
	}
	return false, nil
}
