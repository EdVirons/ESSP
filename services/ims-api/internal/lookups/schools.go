package lookups

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// School export types

type SchoolExport struct {
	Version     string          `json:"version"`
	GeneratedAt time.Time       `json:"generatedAt"`
	Counties    []County        `json:"counties"`
	SubCounties []SubCounty     `json:"subCounties"`
	Schools     []School        `json:"schools"`
	Contacts    []SchoolContact `json:"contacts"`
}

type County struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

type SubCounty struct {
	ID       string `json:"id"`
	CountyID string `json:"countyId"`
	Name     string `json:"name"`
	Code     string `json:"code"`
}

type School struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	CountyID    string `json:"countyId"`
	SubCountyID string `json:"subCountyId"`
	Level       string `json:"level"`
	Type        string `json:"type"`
	Active      bool   `json:"active"`
}

type SchoolContact struct {
	ID        string `json:"id"`
	SchoolID  string `json:"schoolId"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	IsPrimary bool   `json:"isPrimary"`
	Active    bool   `json:"active"`
}

// Summary types for lookup results

type SchoolSummary struct {
	ID            string
	Name          string
	Code          string
	CountyID      string
	CountyName    string
	SubCountyID   string
	SubCountyName string
	Level         string
	Type          string
	Active        bool
}

type ContactSummary struct {
	Name  string
	Phone string
	Email string
	Role  string
}

// LoadSchoolExport loads the school export snapshot
func (s *Store) LoadSchoolExport(ctx context.Context, tenant string) (*SchoolExport, error) {
	snap, err := s.GetSnapshot(ctx, tenant, KindSchool)
	if err != nil {
		return nil, err
	}
	var ex SchoolExport
	if err := json.Unmarshal(snap.Payload, &ex); err != nil {
		return nil, fmt.Errorf("unmarshal school export: %w", err)
	}
	return &ex, nil
}

// SchoolByID looks up a school by its ID
func (s *Store) SchoolByID(ctx context.Context, tenant, schoolID string) (*SchoolSummary, error) {
	ex, err := s.LoadSchoolExport(ctx, tenant)
	if err != nil {
		return nil, ErrSnapshotMissing
	}
	countyNameByID := map[string]string{}
	for _, c := range ex.Counties {
		countyNameByID[c.ID] = c.Name
	}
	subNameByID := map[string]string{}
	for _, scn := range ex.SubCounties {
		subNameByID[scn.ID] = scn.Name
	}
	for _, sc := range ex.Schools {
		if sc.ID == schoolID {
			return &SchoolSummary{
				ID: sc.ID, Name: sc.Name, Code: sc.Code,
				CountyID: sc.CountyID, CountyName: countyNameByID[sc.CountyID],
				SubCountyID: sc.SubCountyID, SubCountyName: subNameByID[sc.SubCountyID],
				Level: sc.Level, Type: sc.Type, Active: sc.Active,
			}, nil
		}
	}
	return nil, ErrNotFound
}

// PrimaryContactBySchoolID looks up the primary contact for a school
func (s *Store) PrimaryContactBySchoolID(ctx context.Context, tenant, schoolID string) (*ContactSummary, error) {
	ex, err := s.LoadSchoolExport(ctx, tenant)
	if err != nil {
		return nil, ErrSnapshotMissing
	}
	var best *SchoolContact
	for i := range ex.Contacts {
		c := &ex.Contacts[i]
		if c.SchoolID != schoolID || !c.Active {
			continue
		}
		if c.IsPrimary {
			best = c
			break
		}
		if best == nil {
			best = c
		}
	}
	if best == nil {
		return nil, ErrNotFound
	}
	return &ContactSummary{Name: best.Name, Phone: best.Phone, Email: best.Email, Role: best.Role}, nil
}
