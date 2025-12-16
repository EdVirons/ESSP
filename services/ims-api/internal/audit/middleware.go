package audit

import (
	"context"
	"net/http"
	"strings"

	"github.com/edvirons/ssp/ims/internal/middleware"
)

type auditCtxKey string

const (
	ctxAuditKey auditCtxKey = "audit"
)

// Middleware captures audit context from the request
func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auditCtx := AuditContext{
				TenantID:  middleware.TenantID(r.Context()),
				UserID:    extractUserID(r),
				UserEmail: extractUserEmail(r),
				IPAddress: extractIPAddress(r),
				UserAgent: r.Header.Get("User-Agent"),
				RequestID: r.Header.Get("X-Request-Id"),
			}

			ctx := context.WithValue(r.Context(), ctxAuditKey, auditCtx)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetAuditContext retrieves the audit context from the request context
func GetAuditContext(ctx context.Context) AuditContext {
	if auditCtx, ok := ctx.Value(ctxAuditKey).(AuditContext); ok {
		return auditCtx
	}
	// Return empty context if not found
	return AuditContext{}
}

// extractUserID extracts the user ID from JWT claims or headers
func extractUserID(r *http.Request) string {
	// Try to get from X-User-Id header (set by auth middleware)
	if userID := r.Header.Get("X-User-Id"); userID != "" {
		return userID
	}

	// Try to get from JWT sub claim (would need to parse token or middleware)
	// For now, we'll rely on the auth middleware to set the header
	return ""
}

// extractUserEmail extracts the user email from JWT claims or headers
func extractUserEmail(r *http.Request) string {
	// Try to get from X-User-Email header (set by auth middleware)
	if email := r.Header.Get("X-User-Email"); email != "" {
		return email
	}
	return ""
}

// extractIPAddress extracts the client IP address from the request
func extractIPAddress(r *http.Request) string {
	// Check X-Forwarded-For header first (for proxied requests)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// X-Forwarded-For can contain multiple IPs, take the first one
		parts := strings.Split(xff, ",")
		if len(parts) > 0 {
			return strings.TrimSpace(parts[0])
		}
	}

	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Fall back to RemoteAddr
	ip := r.RemoteAddr
	// Remove port if present
	if idx := strings.LastIndex(ip, ":"); idx != -1 {
		ip = ip[:idx]
	}
	return ip
}
