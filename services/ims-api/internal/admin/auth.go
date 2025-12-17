package admin

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

// DemoUser represents a demo user for testing different roles
type DemoUser struct {
	Username    string
	Password    string
	Roles       []string
	Email       string
	DisplayName string
}

// demoUsers contains predefined demo users for testing role-based UI
// In production, these would be managed via external SSO (Keycloak)
var demoUsers = map[string]DemoUser{
	"admin": {
		Username:    "admin",
		Password:    "admin123",
		Roles:       []string{"ssp_admin"},
		Email:       "admin@essp.local",
		DisplayName: "System Admin",
	},
	"school_contact": {
		Username:    "school_contact",
		Password:    "school123",
		Roles:       []string{"ssp_school_contact"},
		Email:       "contact@greenwood.edu",
		DisplayName: "Mary Johnson",
	},
	"support_agent": {
		Username:    "support_agent",
		Password:    "support123",
		Roles:       []string{"ssp_support_agent"},
		Email:       "support@essp.local",
		DisplayName: "James Wilson",
	},
	"lead_tech": {
		Username:    "lead_tech",
		Password:    "lead123",
		Roles:       []string{"ssp_lead_tech"},
		Email:       "lead.tech@essp.local",
		DisplayName: "Robert Chen",
	},
	"ops_manager": {
		Username:    "ops_manager",
		Password:    "ops123",
		Roles:       []string{"ssp_ops_manager"},
		Email:       "ops.manager@essp.local",
		DisplayName: "Michael Thompson",
	},
	"field_tech": {
		Username:    "field_tech",
		Password:    "tech123",
		Roles:       []string{"ssp_field_tech"},
		Email:       "field.tech@essp.local",
		DisplayName: "Sarah Martinez",
	},
	"warehouse": {
		Username:    "warehouse",
		Password:    "warehouse123",
		Roles:       []string{"ssp_warehouse_manager"},
		Email:       "warehouse@essp.local",
		DisplayName: "David Kim",
	},
	"sales_marketing": {
		Username:    "sales_marketing",
		Password:    "sales123",
		Roles:       []string{"ssp_sales_marketing"},
		Email:       "sales@essp.local",
		DisplayName: "Emily Taylor",
	},
}

// AdminAuth handles authentication for the admin dashboard
type AdminAuth struct {
	cfg        AdminAuthConfig
	logger     *zap.Logger
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

// NewAdminAuth creates a new admin authentication handler
func NewAdminAuth(cfg AdminAuthConfig, logger *zap.Logger) (*AdminAuth, error) {
	// Generate RSA key pair for JWT signing
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	return &AdminAuth{
		cfg:        cfg,
		logger:     logger,
		privateKey: privateKey,
		publicKey:  &privateKey.PublicKey,
	}, nil
}

// Login handles admin login requests
func (a *AdminAuth) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		a.logger.Debug("failed to decode login request", zap.Error(err))
		a.sendLoginResponse(w, http.StatusBadRequest, false, "invalid request body", nil)
		return
	}

	// Validate credentials
	if req.Username == "" || req.Password == "" {
		a.sendLoginResponse(w, http.StatusBadRequest, false, "username and password are required", nil)
		return
	}

	// Check demo users first
	var authenticatedUser *DemoUser
	if demoUser, exists := demoUsers[req.Username]; exists {
		if demoUser.Password == req.Password {
			authenticatedUser = &demoUser
		}
	}

	// Fallback to config-based admin (for backwards compatibility)
	if authenticatedUser == nil {
		if req.Username == a.cfg.AdminUsername && req.Password == a.cfg.AdminPassword {
			authenticatedUser = &DemoUser{
				Username:    req.Username,
				Password:    req.Password,
				Roles:       []string{"ssp_admin"},
				Email:       "admin@essp.local",
				DisplayName: "Admin",
			}
		}
	}

	if authenticatedUser == nil {
		a.logger.Info("failed login attempt",
			zap.String("username", req.Username),
			zap.String("remote_addr", r.RemoteAddr))
		a.sendLoginResponse(w, http.StatusUnauthorized, false, "invalid credentials", nil)
		return
	}

	// Generate JWT token with user info
	token, err := a.generateTokenForUser(authenticatedUser)
	if err != nil {
		a.logger.Error("failed to generate token", zap.Error(err))
		a.sendLoginResponse(w, http.StatusInternalServerError, false, "authentication failed", nil)
		return
	}

	// Set auth cookie
	a.setAuthCookie(w, token)

	a.logger.Info("login successful",
		zap.String("username", authenticatedUser.Username),
		zap.Strings("roles", authenticatedUser.Roles),
		zap.String("remote_addr", r.RemoteAddr))

	a.sendLoginResponse(w, http.StatusOK, true, "", &User{
		Username:    authenticatedUser.Username,
		Roles:       authenticatedUser.Roles,
		Email:       authenticatedUser.Email,
		DisplayName: authenticatedUser.DisplayName,
	})
}

// Logout handles admin logout requests
func (a *AdminAuth) Logout(w http.ResponseWriter, r *http.Request) {
	// Clear the auth cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "essp_admin_token",
		Value:    "",
		Path:     "/",
		Domain:   a.cfg.CookieDomain,
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   a.cfg.CookieSecure,
		SameSite: http.SameSiteStrictMode,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

// Me returns the current authenticated user
func (a *AdminAuth) Me(w http.ResponseWriter, r *http.Request) {
	// Try to get token from cookie
	cookie, err := r.Cookie("essp_admin_token")
	if err != nil {
		a.sendMeResponse(w, false, nil)
		return
	}

	claims, err := a.validateToken(cookie.Value)
	if err != nil {
		a.logger.Debug("invalid token in me request", zap.Error(err))
		a.sendMeResponse(w, false, nil)
		return
	}

	// Build display name from available fields
	displayName := claims.DisplayName
	if displayName == "" && claims.FirstName != "" {
		displayName = claims.FirstName
		if claims.LastName != "" {
			displayName += " " + claims.LastName
		}
	}
	if displayName == "" {
		displayName = claims.Username
	}

	a.sendMeResponse(w, true, &User{
		Username:    claims.Username,
		Roles:       claims.Roles,
		Email:       claims.Email,
		DisplayName: displayName,
		AvatarURL:   claims.AvatarURL,
		TenantID:    claims.TenantID,
	})
}

// Refresh refreshes the auth token
func (a *AdminAuth) Refresh(w http.ResponseWriter, r *http.Request) {
	// Get existing token from cookie
	cookie, err := r.Cookie("essp_admin_token")
	if err != nil {
		http.Error(w, "no token found", http.StatusUnauthorized)
		return
	}

	claims, err := a.validateToken(cookie.Value)
	if err != nil {
		a.logger.Debug("invalid token in refresh request", zap.Error(err))
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	// Generate new token
	token, err := a.generateToken(claims.Username)
	if err != nil {
		a.logger.Error("failed to refresh token", zap.Error(err))
		http.Error(w, "failed to refresh token", http.StatusInternalServerError)
		return
	}

	// Set new auth cookie
	a.setAuthCookie(w, token)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

// AdminAuthMiddleware is a middleware that validates admin authentication
// It accepts requests with either a valid cookie or a valid Authorization header
func (a *AdminAuth) AdminAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var tokenString string

		// First try to get token from cookie
		if cookie, err := r.Cookie("essp_admin_token"); err == nil {
			tokenString = cookie.Value
		}

		// If no cookie, try Authorization header
		if tokenString == "" {
			authHeader := r.Header.Get("Authorization")
			if authHeader != "" {
				parts := strings.SplitN(authHeader, " ", 2)
				if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
					tokenString = parts[1]
				}
			}
		}

		if tokenString == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		claims, err := a.validateToken(tokenString)
		if err != nil {
			a.logger.Debug("invalid token in request", zap.Error(err))
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// Add claims to request headers for downstream handlers
		r.Header.Set("X-Admin-Username", claims.Username)
		r.Header.Set("X-Tenant-Id", claims.TenantID)

		next.ServeHTTP(w, r)
	})
}

// sendLoginResponse sends a JSON login response
func (a *AdminAuth) sendLoginResponse(w http.ResponseWriter, status int, success bool, message string, user *User) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(LoginResponse{
		Success: success,
		Message: message,
		User:    user,
	})
}

// sendMeResponse sends a JSON me response
func (a *AdminAuth) sendMeResponse(w http.ResponseWriter, authenticated bool, user *User) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(MeResponse{
		Authenticated: authenticated,
		User:          user,
	})
}
