package models

import "time"

// ServiceShop represents a service shop location.
type ServiceShop struct {
	ID            string    `json:"id"`
	TenantID      string    `json:"tenantId"`
	CountyCode    string    `json:"countyCode"`
	CountyName    string    `json:"countyName"`
	SubCountyCode string    `json:"subCountyCode"`
	SubCountyName string    `json:"subCountyName"`
	CoverageLevel string    `json:"coverageLevel"`
	Name          string    `json:"name"`
	Location      string    `json:"location"`
	Active        bool      `json:"active"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// StaffRole represents the role of a staff member.
type StaffRole string

const (
	StaffRoleLeadTechnician      StaffRole = "lead_technician"
	StaffRoleAssistantTechnician StaffRole = "assistant_technician"
	StaffRoleStorekeeper         StaffRole = "storekeeper"
)

// ServiceStaff represents a staff member at a service shop.
type ServiceStaff struct {
	ID            string    `json:"id"`
	TenantID      string    `json:"tenantId"`
	ServiceShopID string    `json:"serviceShopId"`
	UserID        string    `json:"userId"`
	Role          StaffRole `json:"role"`
	Phone         string    `json:"phone"`
	Active        bool      `json:"active"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// Part represents a spare part.
type Part struct {
	ID            string    `json:"id"`
	TenantID      string    `json:"tenantId"`
	SKU           string    `json:"sku"`
	Name          string    `json:"name"`
	Category      string    `json:"category"`
	Description   string    `json:"description"`
	UnitCostCents int       `json:"unitCostCents"`
	Supplier      string    `json:"supplier"`
	SupplierSku   string    `json:"supplierSku"`
	Active        bool      `json:"active"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// InventoryItem represents inventory at a service shop.
type InventoryItem struct {
	ID               string    `json:"id"`
	TenantID         string    `json:"tenantId"`
	ServiceShopID    string    `json:"serviceShopId"`
	PartID           string    `json:"partId"`
	QtyAvailable     int64     `json:"qtyAvailable"`
	QtyReserved      int64     `json:"qtyReserved"`
	ReorderThreshold int64     `json:"reorderThreshold"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

// SchoolContact represents a contact at a school.
type SchoolContact struct {
	ID        string    `json:"id"`
	TenantID  string    `json:"tenantId"`
	SchoolID  string    `json:"schoolId"`
	UserID    string    `json:"userId"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	IsPrimary bool      `json:"isPrimary"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
