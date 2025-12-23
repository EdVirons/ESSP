package middleware

import (
	"net/http"

	"github.com/edvirons/ssp/ims/internal/auth"
	"go.uber.org/zap"
)

// RequireRoles is a simplified middleware that checks if the user has any of the specified roles
// This version doesn't require a logger to be passed (uses nop logger)
func RequireRoles(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRoles := Roles(r.Context())
			if len(userRoles) == 0 {
				http.Error(w, "forbidden: no roles assigned", http.StatusForbidden)
				return
			}

			hasAnyRole := false
			for _, userRole := range userRoles {
				for _, requiredRole := range roles {
					if userRole == requiredRole {
						hasAnyRole = true
						break
					}
				}
				if hasAnyRole {
					break
				}
			}

			if !hasAnyRole {
				http.Error(w, "forbidden: required role not assigned", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequirePermission returns a middleware that checks if the user has the specified permission
func RequirePermission(permission string, logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			roles := Roles(r.Context())
			if len(roles) == 0 {
				logger.Warn("no roles found in context",
					zap.String("path", r.URL.Path),
					zap.String("method", r.Method))
				http.Error(w, "forbidden: no roles assigned", http.StatusForbidden)
				return
			}

			if !auth.UserHasPermission(roles, permission) {
				logger.Warn("permission denied",
					zap.Strings("roles", roles),
					zap.String("required_permission", permission),
					zap.String("path", r.URL.Path),
					zap.String("method", r.Method))
				http.Error(w, "forbidden: insufficient permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireAnyPermission returns a middleware that checks if the user has any of the specified permissions
func RequireAnyPermission(logger *zap.Logger, permissions ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			roles := Roles(r.Context())
			if len(roles) == 0 {
				logger.Warn("no roles found in context",
					zap.String("path", r.URL.Path),
					zap.String("method", r.Method))
				http.Error(w, "forbidden: no roles assigned", http.StatusForbidden)
				return
			}

			if !auth.UserHasAnyPermission(roles, permissions...) {
				logger.Warn("permission denied",
					zap.Strings("roles", roles),
					zap.Strings("required_any_permissions", permissions),
					zap.String("path", r.URL.Path),
					zap.String("method", r.Method))
				http.Error(w, "forbidden: insufficient permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireRole returns a middleware that checks if the user has the specified role
func RequireRole(role string, logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			roles := Roles(r.Context())
			if len(roles) == 0 {
				logger.Warn("no roles found in context",
					zap.String("path", r.URL.Path),
					zap.String("method", r.Method))
				http.Error(w, "forbidden: no roles assigned", http.StatusForbidden)
				return
			}

			hasRole := false
			for _, r := range roles {
				if r == role {
					hasRole = true
					break
				}
			}

			if !hasRole {
				logger.Warn("role check failed",
					zap.Strings("user_roles", roles),
					zap.String("required_role", role),
					zap.String("path", r.URL.Path),
					zap.String("method", r.Method))
				http.Error(w, "forbidden: required role not assigned", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireAnyRole returns a middleware that checks if the user has any of the specified roles
func RequireAnyRole(logger *zap.Logger, roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRoles := Roles(r.Context())
			if len(userRoles) == 0 {
				logger.Warn("no roles found in context",
					zap.String("path", r.URL.Path),
					zap.String("method", r.Method))
				http.Error(w, "forbidden: no roles assigned", http.StatusForbidden)
				return
			}

			hasAnyRole := false
			for _, userRole := range userRoles {
				for _, requiredRole := range roles {
					if userRole == requiredRole {
						hasAnyRole = true
						break
					}
				}
				if hasAnyRole {
					break
				}
			}

			if !hasAnyRole {
				logger.Warn("role check failed",
					zap.Strings("user_roles", userRoles),
					zap.Strings("required_any_roles", roles),
					zap.String("path", r.URL.Path),
					zap.String("method", r.Method))
				http.Error(w, "forbidden: required role not assigned", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireSchoolAccess returns a middleware that validates school-scoped access
// It checks that the user has access to the school specified in the URL or context
func RequireSchoolAccess(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get the school ID from the request context
			requestedSchoolID := SchoolID(r.Context())
			if requestedSchoolID == "" {
				// No school ID in request, allow access (tenant-wide access)
				next.ServeHTTP(w, r)
				return
			}

			// Get user's assigned schools from JWT claims
			assignedSchools := AssignedSchools(r.Context())

			// Admin roles have access to all schools
			roles := Roles(r.Context())
			for _, role := range roles {
				if role == "ssp_admin" {
					next.ServeHTTP(w, r)
					return
				}
			}

			// If user has no school assignments, deny access
			if len(assignedSchools) == 0 {
				logger.Warn("school access denied: no schools assigned",
					zap.String("requested_school", requestedSchoolID),
					zap.Strings("roles", roles),
					zap.String("path", r.URL.Path))
				http.Error(w, "forbidden: no school access", http.StatusForbidden)
				return
			}

			// Check if user has access to the requested school
			hasAccess := false
			for _, schoolID := range assignedSchools {
				if schoolID == requestedSchoolID {
					hasAccess = true
					break
				}
			}

			if !hasAccess {
				logger.Warn("school access denied",
					zap.String("requested_school", requestedSchoolID),
					zap.Strings("assigned_schools", assignedSchools),
					zap.Strings("roles", roles),
					zap.String("path", r.URL.Path))
				http.Error(w, "forbidden: school access denied", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
