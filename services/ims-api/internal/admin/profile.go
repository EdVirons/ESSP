package admin

import (
	"encoding/json"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// SSOUserProfile represents the full user profile from SSO
// This is designed to work with Edvirons ecosystem SSO
type SSOUserProfile struct {
	// Core identity
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`

	// Display information
	DisplayName string `json:"displayName"`
	FirstName   string `json:"firstName,omitempty"`
	LastName    string `json:"lastName,omitempty"`
	AvatarURL   string `json:"avatarUrl,omitempty"`

	// Organization context
	Organization *Organization `json:"organization,omitempty"`
	TenantID     string        `json:"tenantId"`

	// Authorization
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions,omitempty"`

	// SSO metadata
	SSOProvider   string     `json:"ssoProvider,omitempty"`
	SSOSubject    string     `json:"ssoSubject,omitempty"`
	EmailVerified bool       `json:"emailVerified"`
	LastLoginAt   *time.Time `json:"lastLoginAt,omitempty"`
	CreatedAt     *time.Time `json:"createdAt,omitempty"`

	// Preferences
	Preferences *UserPreferences `json:"preferences,omitempty"`
}

// Organization represents the user's organization in the Edvirons ecosystem
type Organization struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName,omitempty"`
	Type        string `json:"type,omitempty"` // e.g., "district", "school", "service_provider"
	LogoURL     string `json:"logoUrl,omitempty"`
}

// UserPreferences holds user-specific preferences
type UserPreferences struct {
	Theme           string                   `json:"theme,omitempty"`    // "light", "dark", "system"
	Language        string                   `json:"language,omitempty"` // e.g., "en", "es"
	Timezone        string                   `json:"timezone,omitempty"` // e.g., "America/New_York"
	SidebarCollapse bool                     `json:"sidebarCollapsed,omitempty"`
	Notifications   *NotificationPreferences `json:"notifications,omitempty"`
}

// NotificationPreferences holds notification settings
type NotificationPreferences struct {
	EmailEnabled    bool `json:"emailEnabled"`
	BrowserEnabled  bool `json:"browserEnabled"`
	IncidentAlerts  bool `json:"incidentAlerts"`
	WorkOrderAlerts bool `json:"workOrderAlerts"`
}

// ProfileResponse is returned by the profile endpoint
type ProfileResponse struct {
	Profile *SSOUserProfile `json:"profile"`
}

// Profile returns the full user profile for the authenticated user
func (a *AdminAuth) Profile(w http.ResponseWriter, r *http.Request) {
	// Get token from cookie
	cookie, err := r.Cookie("essp_admin_token")
	if err != nil {
		sendJSONError(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	claims, err := a.validateToken(cookie.Value)
	if err != nil {
		a.logger.Debug("invalid token in profile request", logError(err))
		sendJSONError(w, http.StatusUnauthorized, "invalid token")
		return
	}

	// Build profile from claims
	// In a full SSO implementation, this would fetch additional data from the SSO provider
	// or from a local user database that syncs with SSO
	profile := buildProfileFromClaims(claims)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ProfileResponse{Profile: profile})
}

// buildProfileFromClaims constructs a user profile from JWT claims
// This is used when SSO data is embedded in the token claims
func buildProfileFromClaims(claims *AdminClaims) *SSOUserProfile {
	now := time.Now()

	profile := &SSOUserProfile{
		ID:          claims.Subject,
		Username:    claims.Username,
		DisplayName: claims.Username, // Default to username
		TenantID:    claims.TenantID,
		Roles:       claims.Roles,
		Permissions: derivePermissionsFromRoles(claims.Roles),
		LastLoginAt: &now,
	}

	// Check for extended claims (if SSO token contains more data)
	if claims.Email != "" {
		profile.Email = claims.Email
		profile.EmailVerified = true
	}

	if claims.DisplayName != "" {
		profile.DisplayName = claims.DisplayName
	}

	if claims.FirstName != "" {
		profile.FirstName = claims.FirstName
	}

	if claims.LastName != "" {
		profile.LastName = claims.LastName
		if claims.FirstName != "" {
			profile.DisplayName = claims.FirstName + " " + claims.LastName
		}
	}

	if claims.AvatarURL != "" {
		profile.AvatarURL = claims.AvatarURL
	}

	if claims.OrganizationID != "" {
		profile.Organization = &Organization{
			ID:          claims.OrganizationID,
			Name:        claims.OrganizationName,
			DisplayName: claims.OrganizationName,
			Type:        claims.OrganizationType,
		}
	}

	// Set default preferences
	profile.Preferences = &UserPreferences{
		Theme:    "light",
		Language: "en",
		Timezone: "America/New_York",
		Notifications: &NotificationPreferences{
			EmailEnabled:    true,
			BrowserEnabled:  false,
			IncidentAlerts:  true,
			WorkOrderAlerts: true,
		},
	}

	return profile
}

// derivePermissionsFromRoles maps roles to permissions
// This provides a basic RBAC implementation
func derivePermissionsFromRoles(roles []string) []string {
	permissionMap := map[string][]string{
		"ssp_admin": {
			"incident:read", "incident:create", "incident:update", "incident:delete",
			"workorder:read", "workorder:create", "workorder:update", "workorder:delete",
			"program:read", "program:create", "program:update", "program:delete",
			"school:read", "school:update",
			"device:read", "device:update",
			"serviceshop:read", "serviceshop:create", "serviceshop:update",
			"inventory:read", "inventory:update",
			"audit:read",
			"settings:read", "settings:update",
			"user:read",
		},
		"ssp_operator": {
			"incident:read", "incident:create", "incident:update",
			"workorder:read", "workorder:create", "workorder:update",
			"school:read",
			"device:read",
			"serviceshop:read",
			"inventory:read",
		},
		"ssp_viewer": {
			"incident:read",
			"workorder:read",
			"school:read",
			"device:read",
		},
	}

	seen := make(map[string]bool)
	var permissions []string

	for _, role := range roles {
		if perms, ok := permissionMap[role]; ok {
			for _, p := range perms {
				if !seen[p] {
					seen[p] = true
					permissions = append(permissions, p)
				}
			}
		}
	}

	return permissions
}

// sendJSONError sends a JSON error response
func sendJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// logError creates a zap field for an error
func logError(err error) zap.Field {
	return zap.Error(err)
}
