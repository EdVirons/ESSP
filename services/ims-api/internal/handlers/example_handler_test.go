package handlers_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/edvirons/ssp/ims/internal/mocks"
	"github.com/edvirons/ssp/ims/internal/testutil"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

// This is an example test demonstrating how to use the HTTP test client
// and mock implementations for handler testing.
//
// NOTE: This is a reference example. Actual handlers should be tested
// in their respective packages.

// Example handler that uses Redis for caching
func exampleCacheHandler(redis *mocks.MockRedisClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		key := "example-key"

		// Try to get from cache
		result := redis.Get(ctx, key)
		if result.Err() == nil {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("X-Cache", "HIT")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"data":"` + result.Val() + `","cached":true}`))
			return
		}

		// Cache miss
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Cache", "MISS")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":"fresh","cached":false}`))
	}
}

// Example handler that publishes events via NATS
func exampleEventHandler(nats *mocks.MockNATSPublisher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tenantID := r.Header.Get("X-Tenant-Id")
		schoolID := r.Header.Get("X-School-Id")

		// Publish event
		event := map[string]string{
			"tenant_id": tenantID,
			"school_id": schoolID,
			"action":    "test",
		}

		err := nats.PublishJSON("test.event", event)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"failed to publish event"}`))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"status":"event published"}`))
	}
}

func TestExampleHandler_WithHTTPClient(t *testing.T) {
	// Create mock dependencies
	redis := mocks.NewMockRedisClient()
	nats := mocks.NewMockNATSPublisher()

	// Set up router
	r := chi.NewRouter()
	r.Get("/cache", exampleCacheHandler(redis))
	r.Post("/event", exampleEventHandler(nats))

	// Create HTTP test client
	client := testutil.NewHTTPTestClient(t, r)

	t.Run("Cache miss", func(t *testing.T) {
		// Make request
		resp := client.Get("/cache")

		// Assert response
		resp.AssertStatus(http.StatusOK).
			AssertContentType("application/json").
			AssertHeader("X-Cache", "MISS").
			AssertBodyContains("fresh")

		// Verify Redis was called
		assert.Equal(t, 1, redis.GetCallCount("Get"))
	})

	t.Run("Cache hit", func(t *testing.T) {
		// Set value in cache
		ctx := context.Background()
		redis.Set(ctx, "example-key", "cached-value", 0)

		// Make request
		resp := client.Get("/cache")

		// Assert response
		resp.AssertStatus(http.StatusOK).
			AssertContentType("application/json").
			AssertHeader("X-Cache", "HIT").
			AssertBodyContains("cached-value")
	})

	t.Run("Publish event with headers", func(t *testing.T) {
		// Make request with tenant/school headers
		resp := client.
			WithTenant("tenant-123").
			WithSchool("school-456").
			Post("/event", nil)

		// Assert response
		resp.AssertStatus(http.StatusCreated).
			AssertBodyContains("event published")

		// Verify event was published
		assert.True(t, nats.WasPublished("test.event"))
		assert.Equal(t, 1, nats.GetMessageCount("test.event"))

		// Verify event content
		msg, ok := nats.GetLastMessage("test.event")
		assert.True(t, ok)

		eventMap, ok := msg.(map[string]string)
		assert.True(t, ok)
		assert.Equal(t, "tenant-123", eventMap["tenant_id"])
		assert.Equal(t, "school-456", eventMap["school_id"])
	})
}

// Example table-driven test with HTTP client
func TestExampleHandler_TableDriven(t *testing.T) {
	redis := mocks.NewMockRedisClient()
	r := chi.NewRouter()
	r.Get("/cache", exampleCacheHandler(redis))
	client := testutil.NewHTTPTestClient(t, r)

	tests := []struct {
		name          string
		setupCache    bool
		expectedCache string
		expectedBody  string
	}{
		{
			name:          "with cached data",
			setupCache:    true,
			expectedCache: "HIT",
			expectedBody:  "cached-value",
		},
		{
			name:          "without cached data",
			setupCache:    false,
			expectedCache: "MISS",
			expectedBody:  "fresh",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset cache
			redis.Reset()

			// Setup
			if tt.setupCache {
				ctx := context.Background()
				redis.Set(ctx, "example-key", "cached-value", 0)
			}

			// Execute
			resp := client.Get("/cache")

			// Assert
			resp.AssertStatus(http.StatusOK).
				AssertHeader("X-Cache", tt.expectedCache).
				AssertBodyContains(tt.expectedBody)
		})
	}
}

// Example test showing how to parse JSON responses
func TestExampleHandler_JSONResponse(t *testing.T) {
	redis := mocks.NewMockRedisClient()
	r := chi.NewRouter()
	r.Get("/cache", exampleCacheHandler(redis))
	client := testutil.NewHTTPTestClient(t, r)

	t.Run("Parse JSON response", func(t *testing.T) {
		resp := client.Get("/cache")

		// Define response structure
		type Response struct {
			Data   string `json:"data"`
			Cached bool   `json:"cached"`
		}

		var result Response
		resp.AssertJSON(&result)

		// Assert fields
		assert.Equal(t, "fresh", result.Data)
		assert.False(t, result.Cached)
	})
}

// Example test demonstrating parallel test execution
func TestExampleHandler_Parallel(t *testing.T) {
	// Tests that don't share state can run in parallel
	t.Run("Test 1", func(t *testing.T) {
		t.Parallel()

		redis := mocks.NewMockRedisClient()
		r := chi.NewRouter()
		r.Get("/cache", exampleCacheHandler(redis))
		client := testutil.NewHTTPTestClient(t, r)

		resp := client.Get("/cache")
		resp.AssertStatus(http.StatusOK)
	})

	t.Run("Test 2", func(t *testing.T) {
		t.Parallel()

		redis := mocks.NewMockRedisClient()
		r := chi.NewRouter()
		r.Get("/cache", exampleCacheHandler(redis))
		client := testutil.NewHTTPTestClient(t, r)

		resp := client.Get("/cache")
		resp.AssertStatus(http.StatusOK)
	})
}
