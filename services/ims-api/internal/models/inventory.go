package models

import (
	"encoding/json"
	"time"
)

// LocationType represents the type of a location.
type LocationType string

const (
	LocationTypeBlock   LocationType = "block"
	LocationTypeFloor   LocationType = "floor"
	LocationTypeRoom    LocationType = "room"
	LocationTypeLab     LocationType = "lab"
	LocationTypeStorage LocationType = "storage"
	LocationTypeOffice  LocationType = "office"
)

// Location represents a physical location within a school (room, lab, block, etc.)
type Location struct {
	ID           string          `json:"id"`
	TenantID     string          `json:"tenantId"`
	SchoolID     string          `json:"schoolId"`
	ParentID     *string         `json:"parentId,omitempty"`
	Name         string          `json:"name"`
	LocationType LocationType    `json:"locationType"`
	Code         string          `json:"code"`
	Capacity     int             `json:"capacity"`
	Metadata     json.RawMessage `json:"metadata,omitempty"`
	Active       bool            `json:"active"`
	CreatedAt    time.Time       `json:"createdAt"`
	UpdatedAt    time.Time       `json:"updatedAt"`

	// Computed fields (not stored)
	Path        string `json:"path,omitempty"`        // Full hierarchical path: "Block A > Floor 1 > Lab 101"
	DeviceCount int    `json:"deviceCount,omitempty"` // Number of devices at this location
}

// AssignmentType represents the type of device assignment.
type AssignmentType string

const (
	AssignmentTypePermanent AssignmentType = "permanent"
	AssignmentTypeTemporary AssignmentType = "temporary"
	AssignmentTypeLoan      AssignmentType = "loan"
	AssignmentTypeRepair    AssignmentType = "repair"
	AssignmentTypeStorage   AssignmentType = "storage"
)

// DeviceAssignment represents a device's assignment to a location or user.
type DeviceAssignment struct {
	ID             string         `json:"id"`
	TenantID       string         `json:"tenantId"`
	DeviceID       string         `json:"deviceId"`
	LocationID     *string        `json:"locationId,omitempty"`
	AssignedToUser string         `json:"assignedToUser,omitempty"`
	AssignmentType AssignmentType `json:"assignmentType"`
	EffectiveFrom  time.Time      `json:"effectiveFrom"`
	EffectiveTo    *time.Time     `json:"effectiveTo,omitempty"` // nil = current
	Notes          string         `json:"notes,omitempty"`
	CreatedBy      string         `json:"createdBy,omitempty"`
	CreatedAt      time.Time      `json:"createdAt"`

	// Computed/joined fields
	Location *Location `json:"location,omitempty"`
}

// GroupType represents the type of device group.
type GroupType string

const (
	GroupTypeManual   GroupType = "manual"   // Manually managed membership
	GroupTypeLocation GroupType = "location" // Auto-membership by location
	GroupTypeDynamic  GroupType = "dynamic"  // Auto-membership by selector criteria
)

// GroupSelector defines criteria for dynamic group membership.
type GroupSelector struct {
	Model        string   `json:"model,omitempty"`        // Device model contains
	Lifecycle    []string `json:"lifecycle,omitempty"`    // Match any lifecycle status
	LocationType string   `json:"locationType,omitempty"` // Match location type
	LocationID   string   `json:"locationId,omitempty"`   // Match specific location
	SchoolID     string   `json:"schoolId,omitempty"`     // Match specific school
}

// DeviceGroup represents a group of devices for control/policy purposes.
type DeviceGroup struct {
	ID          string          `json:"id"`
	TenantID    string          `json:"tenantId"`
	SchoolID    *string         `json:"schoolId,omitempty"` // nil = tenant-wide
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	GroupType   GroupType       `json:"groupType"`
	LocationID  *string         `json:"locationId,omitempty"` // For location-based groups
	Selector    *GroupSelector  `json:"selector,omitempty"`   // For dynamic groups
	Policies    json.RawMessage `json:"policies,omitempty"`   // Future: exam_mode, restrictions
	Active      bool            `json:"active"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`

	// Computed fields
	MemberCount int `json:"memberCount,omitempty"`
}

// GroupMember represents a device's membership in a group.
type GroupMember struct {
	ID       string    `json:"id"`
	TenantID string    `json:"tenantId"`
	GroupID  string    `json:"groupId"`
	DeviceID string    `json:"deviceId"`
	AddedAt  time.Time `json:"addedAt"`
	AddedBy  string    `json:"addedBy,omitempty"`
}

// DeviceNetworkSnapshot is a cached copy of MAC address data from ssot-devices.
type DeviceNetworkSnapshot struct {
	TenantID      string     `json:"tenantId"`
	DeviceID      string     `json:"deviceId"`
	MACAddress    string     `json:"macAddress"`
	InterfaceType string     `json:"interfaceType"`
	IsPrimary     bool       `json:"isPrimary"`
	LastSeenAt    *time.Time `json:"lastSeenAt,omitempty"`
	SyncedAt      time.Time  `json:"syncedAt"`
}

// InventoryDevice represents a device with its inventory context.
type InventoryDevice struct {
	ID         string     `json:"id"`
	TenantID   string     `json:"tenantId"`
	Serial     string     `json:"serial"`
	AssetTag   string     `json:"assetTag"`
	Model      string     `json:"model"` // Denormalized from device_model
	Make       string     `json:"make"`
	SchoolID   string     `json:"schoolId"`
	Lifecycle  string     `json:"lifecycle"`
	Enrolled   bool       `json:"enrolled"`
	LastSeenAt *time.Time `json:"lastSeenAt,omitempty"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`

	// Joined data
	Location     *Location `json:"location,omitempty"`
	LocationPath string    `json:"locationPath,omitempty"`
	MACAddresses []string  `json:"macAddresses,omitempty"`
	Groups       []string  `json:"groups,omitempty"` // Group IDs
}

// InventorySummary provides aggregate statistics for inventory views.
type InventorySummary struct {
	TotalDevices int            `json:"totalDevices"`
	ByStatus     map[string]int `json:"byStatus"`   // lifecycle status -> count
	ByLocation   map[string]int `json:"byLocation"` // "assigned" vs "unassigned"
	ByModel      map[string]int `json:"byModel,omitempty"`
}

// SchoolInventoryResponse is the response for GET /v1/schools/{schoolId}/inventory.
type SchoolInventoryResponse struct {
	School struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"school"`
	Summary   InventorySummary  `json:"summary"`
	Devices   []InventoryDevice `json:"devices"`
	Locations []Location        `json:"locations"`
}

// LocationTreeNode represents a location with its children for hierarchical views.
type LocationTreeNode struct {
	Location
	Children []LocationTreeNode `json:"children,omitempty"`
}
