package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/edvirons/ssp/ims/internal/auth"
	"go.uber.org/zap"
)

// Impersonation context keys
const (
	ctxImpersonation ctxKey = "impersonation"
)

// ImpersonationContext holds information about an impersonation session
type ImpersonationContext struct {
	Active        bool     `json:"active"`
	ActorUserID   string   `json:"actorUserId"`   // The ops manager performing the action
	ActorEmail    string   `json:"actorEmail"`    // Email of ops manager
	TargetUserID  string   `json:"targetUserId"`  // The school contact being impersonated
	TargetEmail   string   `json:"targetEmail"`   // Email of school contact
	TargetSchools []string `json:"targetSchools"` // Schools the school contact has access to
	Reason        string   `json:"reason"`        // Reason for impersonation
}

// ImpersonationLoader is a function that loads target user info
type ImpersonationLoader func(ctx context.Context, tenantID, targetUserID string) (*ImpersonationTarget, error)

// ImpersonationTarget represents the user being impersonated
type ImpersonationTarget struct {
	UserID   string
	Email    string
	Roles    []string
	Schools  []string
	TenantID string
}

// Impersonation middleware checks for impersonation headers and validates permissions
func Impersonation(logger *zap.Logger, loader ImpersonationLoader) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			targetUserID := r.Header.Get("X-Impersonate-User")

			// No impersonation requested
			if targetUserID == "" {
				// Set empty impersonation context
				ctx := WithImpersonation(r.Context(), ImpersonationContext{Active: false})
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			// Impersonation requested - validate
			roles := Roles(r.Context())
			tenantID := TenantID(r.Context())
			actorUserID := UserID(r.Context())

			// Check if user has impersonation permission
			if !auth.UserHasPermission(roles, auth.PermImpersonate) && !auth.UserHasPermission(roles, auth.PermAll) {
				logger.Warn("impersonation attempt without permission",
					zap.String("actorUserId", actorUserID),
					zap.String("targetUserId", targetUserID),
				)
				http.Error(w, "Forbidden: impersonation permission required", http.StatusForbidden)
				return
			}

			// Prevent impersonation chaining
			existingImp := GetImpersonation(r.Context())
			if existingImp.Active {
				logger.Warn("attempted to chain impersonation",
					zap.String("actorUserId", actorUserID),
					zap.String("existingTargetUserId", existingImp.TargetUserID),
					zap.String("newTargetUserId", targetUserID),
				)
				http.Error(w, "Forbidden: cannot impersonate while already impersonating", http.StatusForbidden)
				return
			}

			// Load target user info
			target, err := loader(r.Context(), tenantID, targetUserID)
			if err != nil {
				logger.Warn("failed to load impersonation target",
					zap.String("targetUserId", targetUserID),
					zap.Error(err),
				)
				http.Error(w, "Invalid impersonation target", http.StatusBadRequest)
				return
			}

			// Validate target is a school contact
			if !containsRole(target.Roles, "ssp_school_contact") {
				logger.Warn("attempted to impersonate non-school-contact user",
					zap.String("actorUserId", actorUserID),
					zap.String("targetUserId", targetUserID),
					zap.Strings("targetRoles", target.Roles),
				)
				http.Error(w, "Forbidden: can only impersonate school contacts", http.StatusForbidden)
				return
			}

			// Validate same tenant
			if target.TenantID != tenantID {
				logger.Warn("attempted to impersonate user in different tenant",
					zap.String("actorTenantId", tenantID),
					zap.String("targetTenantId", target.TenantID),
				)
				http.Error(w, "Forbidden: cannot impersonate users in different tenant", http.StatusForbidden)
				return
			}

			// Get actor email from claims
			claims := Claims(r.Context())
			actorEmail := ""
			if email, ok := claims["email"].(string); ok {
				actorEmail = email
			}

			// Build impersonation context
			impCtx := ImpersonationContext{
				Active:        true,
				ActorUserID:   actorUserID,
				ActorEmail:    actorEmail,
				TargetUserID:  target.UserID,
				TargetEmail:   target.Email,
				TargetSchools: target.Schools,
				Reason:        r.Header.Get("X-Impersonate-Reason"),
			}

			logger.Info("impersonation session started",
				zap.String("actorUserId", actorUserID),
				zap.String("actorEmail", actorEmail),
				zap.String("targetUserId", target.UserID),
				zap.String("targetEmail", target.Email),
				zap.Strings("targetSchools", target.Schools),
				zap.String("reason", impCtx.Reason),
			)

			// Set response headers to inform frontend
			w.Header().Set("X-Impersonation-Active", "true")
			w.Header().Set("X-Impersonated-User", target.UserID)
			w.Header().Set("X-Impersonated-Schools", strings.Join(target.Schools, ","))

			ctx := WithImpersonation(r.Context(), impCtx)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// WithImpersonation stores impersonation context
func WithImpersonation(ctx context.Context, imp ImpersonationContext) context.Context {
	return context.WithValue(ctx, ctxImpersonation, imp)
}

// GetImpersonation retrieves impersonation context
func GetImpersonation(ctx context.Context) ImpersonationContext {
	v, _ := ctx.Value(ctxImpersonation).(ImpersonationContext)
	return v
}

// EffectiveSchools returns the schools to use for authorization
// If impersonating, returns the target's schools; otherwise returns the user's assigned schools
func EffectiveSchools(ctx context.Context) []string {
	imp := GetImpersonation(ctx)
	if imp.Active {
		return imp.TargetSchools
	}
	return AssignedSchools(ctx)
}

// IsSchoolAllowed checks if a school is in the effective schools list
func IsSchoolAllowed(ctx context.Context, schoolID string) bool {
	schools := EffectiveSchools(ctx)
	// Empty list means no restrictions (admin)
	if len(schools) == 0 {
		return true
	}
	for _, s := range schools {
		if s == schoolID {
			return true
		}
	}
	return false
}

// containsRole checks if a role is in the list
func containsRole(roles []string, target string) bool {
	for _, r := range roles {
		if r == target {
			return true
		}
	}
	return false
}
