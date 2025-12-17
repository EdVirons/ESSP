package models

import (
	"strings"
	"time"
)

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
	ID            string    `json:"id"`
	TenantID      string    `json:"tenantId"`
	Serial        string    `json:"serial"`
	AssetTag      string    `json:"assetTag"`
	DeviceModelID string    `json:"deviceModelId"`
	SchoolID      string    `json:"schoolId"`
	AssignedTo    string    `json:"assignedTo"`
	Lifecycle     string    `json:"lifecycle"`
	Enrolled      bool      `json:"enrolled"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// InterfaceType represents the type of network interface
type InterfaceType string

const (
	InterfaceEthernet  InterfaceType = "ethernet"
	InterfaceWiFi      InterfaceType = "wifi"
	InterfaceBluetooth InterfaceType = "bluetooth"
	InterfaceUnknown   InterfaceType = "unknown"
)

// DeviceNetworkIdentity represents a MAC address associated with a device
type DeviceNetworkIdentity struct {
	ID            string        `json:"id"`
	TenantID      string        `json:"tenantId"`
	DeviceID      string        `json:"deviceId"`
	MACAddress    string        `json:"macAddress"`    // Normalized: lowercase, colon-separated (aa:bb:cc:dd:ee:ff)
	InterfaceName string        `json:"interfaceName"` // e.g., "eth0", "wlan0", "en0"
	InterfaceType InterfaceType `json:"interfaceType"` // ethernet, wifi, bluetooth, unknown
	IsPrimary     bool          `json:"isPrimary"`
	FirstSeenAt   time.Time     `json:"firstSeenAt"`
	LastSeenAt    time.Time     `json:"lastSeenAt"`
	CreatedAt     time.Time     `json:"createdAt"`
	UpdatedAt     time.Time     `json:"updatedAt"`
}

// NormalizeMACAddress converts a MAC address to lowercase colon-separated format
func NormalizeMACAddress(mac string) string {
	// Remove common separators and convert to lowercase
	mac = strings.ToLower(mac)
	mac = strings.ReplaceAll(mac, "-", "")
	mac = strings.ReplaceAll(mac, ":", "")
	mac = strings.ReplaceAll(mac, ".", "")

	// Reformat as aa:bb:cc:dd:ee:ff
	if len(mac) == 12 {
		return mac[0:2] + ":" + mac[2:4] + ":" + mac[4:6] + ":" + mac[6:8] + ":" + mac[8:10] + ":" + mac[10:12]
	}
	return mac
}

type ExportPayload struct {
	Version           string                  `json:"version"`
	GeneratedAt       time.Time               `json:"generatedAt"`
	Models            []DeviceModel           `json:"models"`
	Devices           []Device                `json:"devices"`
	NetworkIdentities []DeviceNetworkIdentity `json:"networkIdentities,omitempty"`
}
