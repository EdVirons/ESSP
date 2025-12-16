package models

import "time"

// IncidentStatus represents the status of an incident.
type IncidentStatus string

const (
	IncidentNew          IncidentStatus = "new"
	IncidentAcknowledged IncidentStatus = "acknowledged"
	IncidentInProgress   IncidentStatus = "in_progress"
	IncidentEscalated    IncidentStatus = "escalated"
	IncidentResolved     IncidentStatus = "resolved"
	IncidentClosed       IncidentStatus = "closed"
)

// Severity represents the severity level of an incident.
type Severity string

const (
	SeverityLow      Severity = "low"
	SeverityMedium   Severity = "medium"
	SeverityHigh     Severity = "high"
	SeverityCritical Severity = "critical"
)

// Incident represents a device incident.
type Incident struct {
	ID       string `json:"id"`
	TenantID string `json:"tenantId"`
	SchoolID string `json:"schoolId"`
	DeviceID string `json:"deviceId"`

	// Denormalized lookup fields (from SSOT snapshots)
	SchoolName     string `json:"schoolName"`
	CountyID       string `json:"countyId"`
	CountyName     string `json:"countyName"`
	SubCountyID    string `json:"subCountyId"`
	SubCountyName  string `json:"subCountyName"`
	ContactName    string `json:"contactName"`
	ContactPhone   string `json:"contactPhone"`
	ContactEmail   string `json:"contactEmail"`
	DeviceSerial   string `json:"deviceSerial"`
	DeviceAssetTag string `json:"deviceAssetTag"`
	DeviceModelID  string `json:"deviceModelId"`
	DeviceMake     string `json:"deviceMake"`
	DeviceModel    string `json:"deviceModel"`
	DeviceCategory string `json:"deviceCategory"`

	Category    string         `json:"category"`
	Severity    Severity       `json:"severity"`
	Status      IncidentStatus `json:"status"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	ReportedBy  string         `json:"reportedBy"`

	SLADueAt    time.Time `json:"slaDueAt"`
	SLABreached bool      `json:"slaBreached"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
