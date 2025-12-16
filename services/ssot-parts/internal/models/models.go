package models

import "time"

type Part struct {
	ID string `json:"id"`
	TenantID string `json:"tenantId"`
	Name string `json:"name"`
	Category string `json:"category"`
	Puk string `json:"puk"`
	SpecJSON string `json:"specJson"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type PartCompatibility struct {
	ID string `json:"id"`
	TenantID string `json:"tenantId"`
	PartID string `json:"partId"`
	DeviceModelID string `json:"deviceModelId"`
	CreatedAt time.Time `json:"createdAt"`
}

type VendorSKU struct {
	ID string `json:"id"`
	TenantID string `json:"tenantId"`
	PartID string `json:"partId"`
	VendorID string `json:"vendorId"`
	SKU string `json:"sku"`
	UnitPriceCents int64 `json:"unitPriceCents"`
	Currency string `json:"currency"`
	LeadTimeDays int `json:"leadTimeDays"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type ExportPayload struct {
	Version string `json:"version"`
	GeneratedAt time.Time `json:"generatedAt"`
	Parts []Part `json:"parts"`
	Compatibility []PartCompatibility `json:"compatibility"`
	VendorSKUs []VendorSKU `json:"vendorSkus"`
}
