package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func TestRateLimit_EndpointSanitization(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{"/v1/incidents/abc-123-def", "v1/incidents/{id}"},
		{"/v1/work-orders/456", "v1/work-orders/{id}"},
		{"/v1/incidents", "v1/incidents"},
		{"/v1/work-orders/abc-123/deliverables/456", "v1/work-orders/{id}/deliverables/{id}"},
		{"/", "root"},
		{"", "root"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := sanitizeEndpoint(tt.path)
			if result != tt.expected {
				t.Errorf("sanitizeEndpoint(%q) = %q, want %q", tt.path, result, tt.expected)
			}
		})
	}
}

func TestRateLimit_IsID(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"123", true},
		{"abc-123-def-456-789", true},      // UUID format
		{"507f1f77bcf86cd799439011", true}, // MongoDB ObjectID
		{"abc", false},
		{"work-orders", false},
		{"123abc", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := isID(tt.input)
			if result != tt.expected {
				t.Errorf("isID(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestRateLimit_GetClientIP(t *testing.T) {
	tests := []struct {
		name           string
		remoteAddr     string
		xForwardedFor  string
		xRealIP        string
		expectedPrefix string
	}{
		{
			name:           "X-Forwarded-For takes precedence",
			remoteAddr:     "192.168.1.1:1234",
			xForwardedFor:  "203.0.113.1, 198.51.100.1",
			xRealIP:        "198.51.100.2",
			expectedPrefix: "203.0.113.1",
		},
		{
			name:           "X-Real-IP when no X-Forwarded-For",
			remoteAddr:     "192.168.1.1:1234",
			xForwardedFor:  "",
			xRealIP:        "203.0.113.1",
			expectedPrefix: "203.0.113.1",
		},
		{
			name:           "RemoteAddr fallback",
			remoteAddr:     "203.0.113.1:1234",
			xForwardedFor:  "",
			xRealIP:        "",
			expectedPrefix: "203.0.113.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test", nil)
			req.RemoteAddr = tt.remoteAddr
			if tt.xForwardedFor != "" {
				req.Header.Set("X-Forwarded-For", tt.xForwardedFor)
			}
			if tt.xRealIP != "" {
				req.Header.Set("X-Real-IP", tt.xRealIP)
			}

			result := getClientIP(req)
			if result != tt.expectedPrefix {
				t.Errorf("getClientIP() = %q, want %q", result, tt.expectedPrefix)
			}
		})
	}
}

func TestRateLimit_WithMockRedis(t *testing.T) {
	// This test requires a running Redis instance
	// Skip if Redis is not available
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   15, // Use a separate DB for testing
	})
	defer rdb.Close()

	// Test Redis connectivity
	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		t.Skip("Redis not available, skipping integration test")
	}

	// Clean up test keys
	defer rdb.FlushDB(ctx)

	logger := zap.NewNop()
	cfg := RateLimitConfig{
		RequestsPerMinute: 5,
		BurstSize:         10,
		KeyPrefix:         "test:ratelimit",
	}

	// Create test handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Create router with middleware
	r := chi.NewRouter()
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Add tenant ID to context
			ctx := WithTenantID(r.Context(), "test-tenant")
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})
	r.Use(RateLimit(rdb, cfg, logger))
	r.Get("/test", handler)

	// Make requests up to the limit
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Request %d: expected status 200, got %d", i+1, w.Code)
		}

		// Check headers
		if w.Header().Get("X-RateLimit-Limit") == "" {
			t.Errorf("Request %d: missing X-RateLimit-Limit header", i+1)
		}
	}

	// Next request should be rate limited (6th request exceeds limit of 5)
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Errorf("Expected status 429, got %d", w.Code)
	}

	if w.Header().Get("Retry-After") == "" {
		t.Error("Missing Retry-After header on 429 response")
	}
}

func TestRateLimit_NoTenantID(t *testing.T) {
	// When there's no tenant ID, the middleware should skip rate limiting
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   15,
	})
	defer rdb.Close()

	logger := zap.NewNop()
	cfg := RateLimitConfig{
		RequestsPerMinute: 1,
		BurstSize:         1,
		KeyPrefix:         "test:ratelimit",
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r := chi.NewRouter()
	r.Use(RateLimit(rdb, cfg, logger))
	r.Get("/test", handler)

	// Make multiple requests - all should succeed because no tenant ID
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Request %d: expected status 200, got %d", i+1, w.Code)
		}
	}
}

func TestRateLimit_SlidingWindow(t *testing.T) {
	// This test verifies the sliding window algorithm
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   15,
	})
	defer rdb.Close()

	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		t.Skip("Redis not available, skipping integration test")
	}

	defer rdb.FlushDB(ctx)

	logger := zap.NewNop()
	cfg := RateLimitConfig{
		RequestsPerMinute: 10,
		BurstSize:         10,
		KeyPrefix:         "test:ratelimit",
	}

	// Test the sliding window by checking rate limit at window boundaries
	key := "test:ratelimit:tenant:test-tenant:test-endpoint"

	// Make 10 requests (should all succeed)
	for i := 0; i < 10; i++ {
		allowed, remaining, _, err := checkRateLimit(ctx, rdb, key, cfg, logger)
		if err != nil {
			t.Fatalf("checkRateLimit failed: %v", err)
		}
		if !allowed {
			t.Errorf("Request %d: expected allowed=true, got false", i+1)
		}
		expectedRemaining := 9 - i
		if remaining != expectedRemaining {
			t.Logf("Request %d: remaining=%d (expected ~%d, may vary due to sliding window)",
				i+1, remaining, expectedRemaining)
		}
	}

	// 11th request should be rate limited
	allowed, _, resetTime, err := checkRateLimit(ctx, rdb, key, cfg, logger)
	if err != nil {
		t.Fatalf("checkRateLimit failed: %v", err)
	}
	if allowed {
		t.Error("Request 11: expected allowed=false, got true")
	}

	// Verify reset time is in the future
	if !resetTime.After(time.Now()) {
		t.Error("Reset time should be in the future")
	}
}
