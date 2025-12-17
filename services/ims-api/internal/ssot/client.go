package ssot

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/edvirons/ssp/ims/internal/models"
)

type Client struct {
	BaseURL  string
	TenantID string
	HTTP     *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTP:    &http.Client{Timeout: 15 * time.Second},
	}
}

func (c *Client) WithTenant(tenantID string) *Client {
	return &Client{
		BaseURL:  c.BaseURL,
		TenantID: tenantID,
		HTTP:     c.HTTP,
	}
}

type SchoolsPage struct {
	Items []models.SchoolSnapshot `json:"items"`
	Next  string                 `json:"nextCursor"`
}

type DevicesPage struct {
	Items []models.DeviceSnapshot `json:"items"`
	Next  string                 `json:"nextCursor"`
}

type PartsPage struct {
	Items []models.PartSnapshot `json:"items"`
	Next  string               `json:"nextCursor"`
}

func (c *Client) fetch(path string, q url.Values, out any) error {
	if c.BaseURL == "" {
		return fmt.Errorf("base url not set")
	}
	u, err := url.Parse(c.BaseURL)
	if err != nil { return err }
	u.Path = path
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil { return err }
	if c.TenantID != "" {
		req.Header.Set("X-Tenant-Id", c.TenantID)
	}
	resp, err := c.HTTP.Do(req)
	if err != nil { return err }
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("ssot fetch failed: %s", resp.Status)
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

func (c *Client) ListSchools(updatedSince time.Time, cursor string, limit int) (SchoolsPage, error) {
	var out SchoolsPage
	q := url.Values{}
	if !updatedSince.IsZero() { q.Set("updatedSince", updatedSince.Format(time.RFC3339)) }
	if cursor != "" { q.Set("cursor", cursor) }
	if limit > 0 { q.Set("limit", fmt.Sprintf("%d", limit)) }
	err := c.fetch("/v1/schools", q, &out)
	return out, err
}

func (c *Client) ListDevices(updatedSince time.Time, cursor string, limit int) (DevicesPage, error) {
	var out DevicesPage
	q := url.Values{}
	if !updatedSince.IsZero() { q.Set("updatedSince", updatedSince.Format(time.RFC3339)) }
	if cursor != "" { q.Set("cursor", cursor) }
	if limit > 0 { q.Set("limit", fmt.Sprintf("%d", limit)) }
	err := c.fetch("/v1/devices", q, &out)
	return out, err
}

func (c *Client) ListParts(updatedSince time.Time, cursor string, limit int) (PartsPage, error) {
	var out PartsPage
	q := url.Values{}
	if !updatedSince.IsZero() { q.Set("updatedSince", updatedSince.Format(time.RFC3339)) }
	if cursor != "" { q.Set("cursor", cursor) }
	if limit > 0 { q.Set("limit", fmt.Sprintf("%d", limit)) }
	err := c.fetch("/v1/parts", q, &out)
	return out, err
}

// HR SSOT types and methods

type HRExportPayload struct {
	Version         string           `json:"version"`
	GeneratedAt     time.Time        `json:"generatedAt"`
	OrgUnits        []HROrgUnit      `json:"orgUnits"`
	People          []HRPerson       `json:"people"`
	Teams           []HRTeam         `json:"teams"`
	TeamMemberships []HRMembership   `json:"teamMemberships"`
}

type HROrgUnit struct {
	ID        string    `json:"id"`
	TenantID  string    `json:"tenantId"`
	ParentID  string    `json:"parentId"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	Kind      string    `json:"kind"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type HRPerson struct {
	ID         string    `json:"id"`
	TenantID   string    `json:"tenantId"`
	OrgUnitID  string    `json:"orgUnitId"`
	Status     string    `json:"status"`
	GivenName  string    `json:"givenName"`
	FamilyName string    `json:"familyName"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	Title      string    `json:"title"`
	AvatarURL  string    `json:"avatarUrl"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type HRTeam struct {
	ID          string    `json:"id"`
	TenantID    string    `json:"tenantId"`
	OrgUnitID   string    `json:"orgUnitId"`
	Key         string    `json:"key"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type HRMembership struct {
	ID        string     `json:"id"`
	TenantID  string     `json:"tenantId"`
	TeamID    string     `json:"teamId"`
	PersonID  string     `json:"personId"`
	Role      string     `json:"role"`
	Status    string     `json:"status"`
	StartedAt *time.Time `json:"startedAt"`
	EndedAt   *time.Time `json:"endedAt"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}

// FetchHRExport fetches the full HR export from the ssot-hr service
func (c *Client) FetchHRExport() (HRExportPayload, error) {
	var out HRExportPayload
	q := url.Values{}
	err := c.fetch("/v1/export", q, &out)
	return out, err
}
