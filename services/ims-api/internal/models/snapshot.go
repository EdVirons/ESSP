package models

import "time"

// County represents a county.
type County struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// SchoolProfile represents a school profile.
type SchoolProfile struct {
	TenantID   string    `json:"tenantId"`
	SchoolID   string    `json:"schoolId"`
	CountyCode string    `json:"countyCode"`
	CountyName string    `json:"countyName"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

// SchoolSnapshot represents a snapshot of school data.
type SchoolSnapshot struct {
	TenantID      string    `json:"tenantId"`
	SchoolID      string    `json:"schoolId"`
	Name          string    `json:"name"`
	CountyCode    string    `json:"countyCode"`
	CountyName    string    `json:"countyName"`
	SubCountyCode string    `json:"subCountyCode"`
	SubCountyName string    `json:"subCountyName"`
	Level         string    `json:"level,omitempty"`
	Type          string    `json:"type,omitempty"`
	KnecCode      string    `json:"knecCode,omitempty"`
	Uic           string    `json:"uic,omitempty"`
	Sex           string    `json:"sex,omitempty"`
	Cluster       string    `json:"cluster,omitempty"`
	Accommodation string    `json:"accommodation,omitempty"`
	Latitude      float64   `json:"latitude,omitempty"`
	Longitude     float64   `json:"longitude,omitempty"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// DeviceSnapshot represents a snapshot of device data.
type DeviceSnapshot struct {
	TenantID  string    `json:"tenantId"`
	DeviceID  string    `json:"deviceId"`
	SchoolID  string    `json:"schoolId"`
	Model     string    `json:"model"`
	Serial    string    `json:"serial"`
	AssetTag  string    `json:"assetTag"`
	Status    string    `json:"status"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// PartSnapshot represents a snapshot of part data.
type PartSnapshot struct {
	TenantID  string    `json:"tenantId"`
	PartID    string    `json:"partId"`
	PUK       string    `json:"puk"`
	Name      string    `json:"name"`
	Category  string    `json:"category"`
	Unit      string    `json:"unit"`
	UpdatedAt time.Time `json:"updatedAt"`
}
