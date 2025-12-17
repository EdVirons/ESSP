package models

import "time"

// ID Prefixes for ESSP standard prefixed ULIDs
const (
	PrefixPerson     = "person"
	PrefixTeam       = "team"
	PrefixMembership = "tmem"
	PrefixOrgUnit    = "org"
	PrefixAudit      = "hraud"
)

// OrgUnit represents an organizational unit (department, division, team, etc.)
type OrgUnit struct {
	ID        string         `json:"id"`
	TenantID  string         `json:"tenantId"`
	ParentID  string         `json:"parentId,omitempty"`
	Code      string         `json:"code"`
	Name      string         `json:"name"`
	Kind      string         `json:"kind"` // company|division|department|team|group
	SpecJSON  map[string]any `json:"specJson,omitempty"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
}

// Person represents a person/employee record
type Person struct {
	ID         string         `json:"id"`
	TenantID   string         `json:"tenantId"`
	OrgUnitID  string         `json:"orgUnitId,omitempty"`
	Status     string         `json:"status"` // active|inactive|onboarding|offboarding
	GivenName  string         `json:"givenName"`
	FamilyName string         `json:"familyName"`
	Email      string         `json:"email"`
	Phone      string         `json:"phone,omitempty"`
	Title      string         `json:"title,omitempty"`
	AvatarURL  string         `json:"avatarUrl,omitempty"`
	SpecJSON   map[string]any `json:"specJson,omitempty"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
}

// Team represents a team/workgroup
type Team struct {
	ID          string         `json:"id"`
	TenantID    string         `json:"tenantId"`
	OrgUnitID   string         `json:"orgUnitId,omitempty"`
	Key         string         `json:"key"` // stable slug for integrations
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	SpecJSON    map[string]any `json:"specJson,omitempty"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
}

// TeamMembership represents a person's membership in a team
type TeamMembership struct {
	ID        string         `json:"id"`
	TenantID  string         `json:"tenantId"`
	TeamID    string         `json:"teamId"`
	PersonID  string         `json:"personId"`
	Role      string         `json:"role"`   // lead|member|observer
	Status    string         `json:"status"` // active|inactive
	StartedAt *time.Time     `json:"startedAt,omitempty"`
	EndedAt   *time.Time     `json:"endedAt,omitempty"`
	SpecJSON  map[string]any `json:"specJson,omitempty"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
}

// AuditLog represents an audit log entry
type AuditLog struct {
	ID            string         `json:"id"`
	TenantID      string         `json:"tenantId"`
	ActorPersonID string         `json:"actorPersonId,omitempty"`
	Action        string         `json:"action"`
	EntityType    string         `json:"entityType"`
	EntityID      string         `json:"entityId"`
	RequestID     string         `json:"requestId,omitempty"`
	IPAddress     string         `json:"ipAddress,omitempty"`
	UserAgent     string         `json:"userAgent,omitempty"`
	BeforeJSON    map[string]any `json:"beforeJson,omitempty"`
	AfterJSON     map[string]any `json:"afterJson,omitempty"`
	DiffJSON      map[string]any `json:"diffJson,omitempty"`
	CreatedAt     time.Time      `json:"createdAt"`
}

// ExportPayload for SSOT export/import
type ExportPayload struct {
	Version         string           `json:"version"`
	GeneratedAt     time.Time        `json:"generatedAt"`
	OrgUnits        []OrgUnit        `json:"orgUnits"`
	People          []Person         `json:"people"`
	Teams           []Team           `json:"teams"`
	TeamMemberships []TeamMembership `json:"teamMemberships"`
}
