package models

import "time"

type County struct {
	ID string `json:"id"`
	TenantID string `json:"tenantId"`
	Name string `json:"name"`
	Code string `json:"code"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type SubCounty struct {
	ID string `json:"id"`
	TenantID string `json:"tenantId"`
	CountyID string `json:"countyId"`
	Name string `json:"name"`
	Code string `json:"code"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type School struct {
	ID            string    `json:"id"`
	TenantID      string    `json:"tenantId"`
	Name          string    `json:"name"`
	Code          string    `json:"code"`
	CountyID      string    `json:"countyId"`
	SubCountyID   string    `json:"subCountyId"`
	Level         string    `json:"level"`
	Type          string    `json:"type"`
	Active        bool      `json:"active"`
	KnecCode      string    `json:"knecCode,omitempty"`
	Uic           string    `json:"uic,omitempty"`
	Sex           string    `json:"sex,omitempty"`
	Cluster       string    `json:"cluster,omitempty"`
	Accommodation string    `json:"accommodation,omitempty"`
	Latitude      float64   `json:"latitude,omitempty"`
	Longitude     float64   `json:"longitude,omitempty"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type Contact struct {
	ID string `json:"id"`
	TenantID string `json:"tenantId"`
	SchoolID string `json:"schoolId"`
	Name string `json:"name"`
	Phone string `json:"phone"`
	Email string `json:"email"`
	Role string `json:"role"`
	IsPrimary bool `json:"isPrimary"`
	Active bool `json:"active"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type ExportPayload struct {
	Version string `json:"version"`
	GeneratedAt time.Time `json:"generatedAt"`
	Counties []County `json:"counties"`
	SubCounties []SubCounty `json:"subCounties"`
	Schools []School `json:"schools"`
	Contacts []Contact `json:"contacts"`
}
