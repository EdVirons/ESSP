package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// RateLimitConfig holds rate limit configuration
type RateLimitConfig struct {
	RequestsPerMinute int
	BurstSize         int
	KeyPrefix         string
}

// rateLimitResponse is the JSON response for rate limit exceeded
type rateLimitResponse struct {
	Error      string `json:"error"`
	Message    string `json:"message"`
	RetryAfter int    `json:"retryAfter"`
}

// RateLimit creates a rate limiting middleware using Redis
// It limits requests per tenant to prevent abuse
// Uses a sliding window counter algorithm for accurate rate limiting
func RateLimit(rdb *redis.Client, cfg RateLimitConfig, logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tenantID := TenantID(r.Context())
			if tenantID == "" {
				logger.Warn("rate limit skipped: no tenant ID in context")
				next.ServeHTTP(w, r)
				return
			}

			// Create rate limit key with tenant and endpoint
			endpoint := sanitizeEndpoint(r.URL.Path)
			key := fmt.Sprintf("%s:tenant:%s:%s", cfg.KeyPrefix, tenantID, endpoint)

			allowed, remaining, resetTime, err := checkRateLimit(r.Context(), rdb, key, cfg, logger)
			if err != nil {
				logger.Error("rate limit check failed",
					zap.Error(err),
					zap.String("tenant_id", tenantID),
					zap.String("endpoint", endpoint),
				)
				// On error, allow request to proceed (fail open)
				next.ServeHTTP(w, r)
				return
			}

			// Set rate limit headers
			setRateLimitHeaders(w, cfg.RequestsPerMinute, remaining, resetTime)

			if !allowed {
				retryAfter := int(time.Until(resetTime).Seconds())
				if retryAfter < 0 {
					retryAfter = 60
				}

				w.Header().Set("Retry-After", strconv.Itoa(retryAfter))
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)

				resp := rateLimitResponse{
					Error:      "rate_limit_exceeded",
					Message:    "Too many requests. Please try again later.",
					RetryAfter: retryAfter,
				}
				_ = json.NewEncoder(w).Encode(resp)

				logger.Warn("rate limit exceeded",
					zap.String("tenant_id", tenantID),
					zap.String("endpoint", endpoint),
					zap.Int("retry_after", retryAfter),
				)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RateLimitByIP creates a rate limiter keyed by client IP
// Useful for unauthenticated endpoints
func RateLimitByIP(rdb *redis.Client, cfg RateLimitConfig, logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			clientIP := getClientIP(r)
			if clientIP == "" {
				logger.Warn("rate limit skipped: no client IP found")
				next.ServeHTTP(w, r)
				return
			}

			// Create rate limit key with IP and endpoint
			endpoint := sanitizeEndpoint(r.URL.Path)
			key := fmt.Sprintf("%s:ip:%s:%s", cfg.KeyPrefix, clientIP, endpoint)

			allowed, remaining, resetTime, err := checkRateLimit(r.Context(), rdb, key, cfg, logger)
			if err != nil {
				logger.Error("rate limit check failed",
					zap.Error(err),
					zap.String("client_ip", clientIP),
					zap.String("endpoint", endpoint),
				)
				// On error, allow request to proceed (fail open)
				next.ServeHTTP(w, r)
				return
			}

			// Set rate limit headers
			setRateLimitHeaders(w, cfg.RequestsPerMinute, remaining, resetTime)

			if !allowed {
				retryAfter := int(time.Until(resetTime).Seconds())
				if retryAfter < 0 {
					retryAfter = 60
				}

				w.Header().Set("Retry-After", strconv.Itoa(retryAfter))
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)

				resp := rateLimitResponse{
					Error:      "rate_limit_exceeded",
					Message:    "Too many requests. Please try again later.",
					RetryAfter: retryAfter,
				}
				_ = json.NewEncoder(w).Encode(resp)

				logger.Warn("rate limit exceeded",
					zap.String("client_ip", clientIP),
					zap.String("endpoint", endpoint),
					zap.Int("retry_after", retryAfter),
				)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// checkRateLimit performs the rate limit check using sliding window counter algorithm
// Returns: allowed, remaining, resetTime, error
func checkRateLimit(ctx context.Context, rdb *redis.Client, key string, cfg RateLimitConfig, logger *zap.Logger) (bool, int, time.Time, error) {
	now := time.Now()
	windowStart := now.Truncate(time.Minute)
	resetTime := windowStart.Add(time.Minute)

	// Use sliding window with current and previous minute
	currentKey := fmt.Sprintf("%s:%d", key, windowStart.Unix())
	previousKey := fmt.Sprintf("%s:%d", key, windowStart.Add(-time.Minute).Unix())

	pipe := rdb.Pipeline()
	currentCount := pipe.Get(ctx, currentKey)
	previousCount := pipe.Get(ctx, previousKey)
	_, err := pipe.Exec(ctx)

	if err != nil && err != redis.Nil {
		return false, 0, resetTime, fmt.Errorf("redis pipeline exec failed: %w", err)
	}

	// Calculate weighted count using sliding window
	var current, previous int64

	if val, err := currentCount.Int64(); err == nil {
		current = val
	}
	if val, err := previousCount.Int64(); err == nil {
		previous = val
	}

	// Calculate the position in the current window (0.0 to 1.0)
	elapsed := now.Sub(windowStart)
	windowProgress := float64(elapsed) / float64(time.Minute)

	// Sliding window count: current + (previous * (1 - progress))
	slidingCount := float64(current) + (float64(previous) * (1.0 - windowProgress))

	// Check if limit is exceeded
	limit := int64(cfg.RequestsPerMinute)
	if cfg.BurstSize > 0 && cfg.BurstSize > cfg.RequestsPerMinute {
		limit = int64(cfg.BurstSize)
	}

	allowed := slidingCount < float64(limit)

	if allowed {
		// Increment current window counter
		pipe := rdb.Pipeline()
		pipe.Incr(ctx, currentKey)
		pipe.Expire(ctx, currentKey, 2*time.Minute) // Keep for 2 minutes for sliding window
		_, err := pipe.Exec(ctx)
		if err != nil {
			return false, 0, resetTime, fmt.Errorf("redis incr failed: %w", err)
		}
	}

	remaining := int(float64(limit) - slidingCount - 1)
	if remaining < 0 {
		remaining = 0
	}

	return allowed, remaining, resetTime, nil
}

// setRateLimitHeaders sets the standard rate limit headers
func setRateLimitHeaders(w http.ResponseWriter, limit, remaining int, resetTime time.Time) {
	w.Header().Set("X-RateLimit-Limit", strconv.Itoa(limit))
	w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(remaining))
	w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(resetTime.Unix(), 10))
}

// sanitizeEndpoint normalizes endpoint paths for rate limiting
// Removes IDs and parameters to group similar endpoints together
func sanitizeEndpoint(path string) string {
	// Split path into segments
	segments := strings.Split(strings.Trim(path, "/"), "/")

	var sanitized []string
	for i, seg := range segments {
		// Keep first segment (version)
		if i == 0 {
			sanitized = append(sanitized, seg)
			continue
		}

		// Replace UUIDs and numeric IDs with placeholder
		if isID(seg) {
			sanitized = append(sanitized, "{id}")
		} else {
			sanitized = append(sanitized, seg)
		}
	}

	result := strings.Join(sanitized, "/")
	if result == "" {
		return "root"
	}
	return result
}

// isID checks if a segment looks like an ID (UUID or numeric)
func isID(s string) bool {
	// Check if it's a number
	if _, err := strconv.ParseInt(s, 10, 64); err == nil {
		return true
	}

	// Check if it's a UUID (simple check: contains hyphens and hex chars)
	if len(s) == 36 && strings.Count(s, "-") == 4 {
		return true
	}

	// Check if it's a hex string (like MongoDB ObjectID)
	if len(s) == 24 || len(s) == 32 {
		for _, ch := range s {
			if !((ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F')) {
				return false
			}
		}
		return true
	}

	return false
}

// getClientIP extracts the client IP from the request
// Checks X-Forwarded-For, X-Real-IP headers, then falls back to RemoteAddr
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// Take the first IP in the chain
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
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
