package middleware

import (
	"net/http"
	"strings"

	"github.com/edvirons/ssp/ims/internal/auth"
	"go.uber.org/zap"
)

func AuthJWT(v *auth.Verifier, logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h := r.Header.Get("Authorization")
			if h == "" {
				http.Error(w, "missing authorization", http.StatusUnauthorized)
				return
			}
			parts := strings.SplitN(h, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				http.Error(w, "invalid authorization header", http.StatusUnauthorized)
				return
			}

			claims, err := v.Verify(r.Context(), parts[1])
			if err != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			// Store claims in context
			ctx := WithClaims(r.Context(), claims)

			// Extract and store tenant ID
			if tenant, ok := claims["tenantId"].(string); ok && tenant != "" {
				r.Header.Set("X-Tenant-Id", tenant)
			}

			// Extract and store school ID (for backward compatibility)
			if school, ok := claims["schoolId"].(string); ok && school != "" {
				r.Header.Set("X-School-Id", school)
			}

			// Extract and store roles
			roles := extractRolesFromClaims(claims)
			ctx = WithRoles(ctx, roles)

			// Extract and store assigned schools
			schools := extractSchoolsFromClaims(claims)
			ctx = WithAssignedSchools(ctx, schools)

			// Log authentication details for debugging
			logger.Debug("authenticated user",
				zap.Strings("roles", roles),
				zap.Strings("schools", schools),
				zap.String("path", r.URL.Path))

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// extractRolesFromClaims extracts roles from JWT claims
// Supports multiple claim formats: "roles", "realm_access.roles", "resource_access"
func extractRolesFromClaims(claims map[string]any) []string {
	var roles []string

	// Check for direct "roles" claim (simple array)
	if rolesVal, ok := claims["roles"]; ok {
		if rolesArr, ok := rolesVal.([]interface{}); ok {
			for _, r := range rolesArr {
				if roleStr, ok := r.(string); ok {
					roles = append(roles, roleStr)
				}
			}
		}
	}

	// Check for Keycloak "realm_access.roles" format
	if realmAccess, ok := claims["realm_access"].(map[string]interface{}); ok {
		if realmRoles, ok := realmAccess["roles"].([]interface{}); ok {
			for _, r := range realmRoles {
				if roleStr, ok := r.(string); ok {
					// Only include SSP roles
					if strings.HasPrefix(roleStr, "ssp_") {
						roles = append(roles, roleStr)
					}
				}
			}
		}
	}

	// Check for Keycloak "resource_access" format
	if resourceAccess, ok := claims["resource_access"].(map[string]interface{}); ok {
		for _, resource := range resourceAccess {
			if resourceMap, ok := resource.(map[string]interface{}); ok {
				if resourceRoles, ok := resourceMap["roles"].([]interface{}); ok {
					for _, r := range resourceRoles {
						if roleStr, ok := r.(string); ok {
							// Only include SSP roles
							if strings.HasPrefix(roleStr, "ssp_") {
								roles = append(roles, roleStr)
							}
						}
					}
				}
			}
		}
	}

	return roles
}

// extractSchoolsFromClaims extracts school assignments from JWT claims
func extractSchoolsFromClaims(claims map[string]any) []string {
	var schools []string

	// Check for "schools" claim (array of school IDs)
	if schoolsVal, ok := claims["schools"]; ok {
		if schoolsArr, ok := schoolsVal.([]interface{}); ok {
			for _, s := range schoolsArr {
				if schoolStr, ok := s.(string); ok {
					schools = append(schools, schoolStr)
				}
			}
		}
	}

	// Check for single "schoolId" claim (backward compatibility)
	if schoolID, ok := claims["schoolId"].(string); ok && schoolID != "" {
		schools = append(schools, schoolID)
	}

	return schools
}
