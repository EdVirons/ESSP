package models

import "time"

type DeviceModel struct {
	ID string `json:"id"`
	TenantID string `json:"tenantId"`
	Make string `json:"make"`
	Model string `json:"model"`
	Category string `json:"category"`
	SpecJSON string `json:"specJson"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Device struct {
	ID string `json:"id"`
	TenantID string `json:"tenantId"`
	Serial string `json:"serial"`
	AssetTag string `json:"assetTag"`
	DeviceModelID string `json:"deviceModelId"`
	SchoolID string `json:"schoolId"`
	AssignedTo string `json:"assignedTo"`
	Lifecycle string `json:"lifecycle"`
	Enrolled bool `json:"enrolled"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type ExportPayload struct {
	Version string `json:"version"`
	GeneratedAt time.Time `json:"generatedAt"`
	Models []DeviceModel `json:"models"`
	Devices []Device `json:"devices"`
}
