package admin

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// AdminAuthConfig holds configuration for admin authentication
type AdminAuthConfig struct {
	// AdminUsername is the admin login username
	AdminUsername string
	// AdminPassword is the admin login password (should be bcrypt hashed in production)
	AdminPassword string
	// JWTSecretKey is used for signing JWTs when external OIDC is not configured
	JWTSecretKey string
	// TokenExpiry is how long tokens are valid
	TokenExpiry time.Duration
	// CookieDomain is the domain for auth cookies
	CookieDomain string
	// CookieSecure determines if cookies require HTTPS
	CookieSecure bool
}

// LoginRequest represents the login request body
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse represents the login response body
type LoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	User    *User  `json:"user,omitempty"`
}

// User represents the authenticated user info
type User struct {
	Username    string   `json:"username"`
	Roles       []string `json:"roles"`
	Email       string   `json:"email,omitempty"`
	DisplayName string   `json:"displayName,omitempty"`
	AvatarURL   string   `json:"avatarUrl,omitempty"`
	TenantID    string   `json:"tenantId,omitempty"`
}

// MeResponse represents the current user response
type MeResponse struct {
	Authenticated bool  `json:"authenticated"`
	User          *User `json:"user,omitempty"`
}

// AdminClaims represents the JWT claims for admin tokens
// Extended to support SSO user profile data from Edvirons ecosystem
type AdminClaims struct {
	jwt.RegisteredClaims
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	TenantID string   `json:"tenantId"`

	// Extended SSO profile fields
	Email            string `json:"email,omitempty"`
	DisplayName      string `json:"displayName,omitempty"`
	FirstName        string `json:"firstName,omitempty"`
	LastName         string `json:"lastName,omitempty"`
	AvatarURL        string `json:"avatarUrl,omitempty"`
	OrganizationID   string `json:"organizationId,omitempty"`
	OrganizationName string `json:"organizationName,omitempty"`
	OrganizationType string `json:"organizationType,omitempty"`
}
