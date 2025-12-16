package middleware

import (
	"net/http"
	"strings"
)

// SecurityHeaders adds essential security headers to all responses
func SecurityHeaders() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Prevent MIME type sniffing
			w.Header().Set("X-Content-Type-Options", "nosniff")

			// Prevent clickjacking attacks
			w.Header().Set("X-Frame-Options", "DENY")

			// Enable XSS protection in older browsers
			w.Header().Set("X-XSS-Protection", "1; mode=block")

			// Enforce HTTPS for 1 year
			w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

			// CSP: Use relaxed policy for dashboard, strict for API
			// API paths: /v1/*, /healthz, /readyz, /metrics, /ws
			isAPIPath := strings.HasPrefix(r.URL.Path, "/v1/") ||
				r.URL.Path == "/healthz" ||
				r.URL.Path == "/readyz" ||
				r.URL.Path == "/metrics" ||
				r.URL.Path == "/ws"

			if isAPIPath {
				// Strict CSP for API (no content loading)
				w.Header().Set("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'")
			} else {
				// Dashboard CSP - allows scripts, styles, images, and API connections
				w.Header().Set("Content-Security-Policy",
					"default-src 'self'; "+
					"script-src 'self'; "+
					"style-src 'self' 'unsafe-inline'; "+
					"img-src 'self' data: blob:; "+
					"font-src 'self'; "+
					"connect-src 'self' ws: wss:; "+
					"frame-ancestors 'none'")
			}

			// Control referrer information
			w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

			next.ServeHTTP(w, r)
		})
	}
}

// MaxBodySize limits the size of request bodies to prevent DoS attacks
func MaxBodySize(maxBytes int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ContentLength > maxBytes {
				http.Error(w, "request body too large", http.StatusRequestEntityTooLarge)
				return
			}

			// Limit the request body reader to prevent memory exhaustion
			r.Body = http.MaxBytesReader(w, r.Body, maxBytes)

			next.ServeHTTP(w, r)
		})
	}
}
