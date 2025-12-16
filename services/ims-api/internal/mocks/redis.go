package mocks

import (
	"context"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// MockRedisClient is a mock implementation of Redis client for testing.
// It provides an in-memory store that mimics basic Redis operations.
type MockRedisClient struct {
	mu    sync.RWMutex
	data  map[string]string
	ttls  map[string]time.Time
	calls map[string]int // Track method calls for verification
}

// NewMockRedisClient creates a new mock Redis client.
func NewMockRedisClient() *MockRedisClient {
	return &MockRedisClient{
		data:  make(map[string]string),
		ttls:  make(map[string]time.Time),
		calls: make(map[string]int),
	}
}

// Get retrieves a value from the mock store.
func (m *MockRedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	m.mu.RLock()
	defer m.mu.RUnlock()

	m.trackCall("Get")

	// Check if key has expired
	if expiry, exists := m.ttls[key]; exists && time.Now().After(expiry) {
		delete(m.data, key)
		delete(m.ttls, key)
		return redis.NewStringResult("", redis.Nil)
	}

	val, exists := m.data[key]
	if !exists {
		return redis.NewStringResult("", redis.Nil)
	}
	return redis.NewStringResult(val, nil)
}

// Set stores a value in the mock store.
func (m *MockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.trackCall("Set")

	m.data[key] = value.(string)
	if expiration > 0 {
		m.ttls[key] = time.Now().Add(expiration)
	} else {
		delete(m.ttls, key)
	}
	return redis.NewStatusResult("OK", nil)
}

// Del deletes one or more keys from the mock store.
func (m *MockRedisClient) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.trackCall("Del")

	deleted := int64(0)
	for _, key := range keys {
		if _, exists := m.data[key]; exists {
			delete(m.data, key)
			delete(m.ttls, key)
			deleted++
		}
	}
	return redis.NewIntResult(deleted, nil)
}

// Exists checks if one or more keys exist.
func (m *MockRedisClient) Exists(ctx context.Context, keys ...string) *redis.IntCmd {
	m.mu.RLock()
	defer m.mu.RUnlock()

	m.trackCall("Exists")

	count := int64(0)
	for _, key := range keys {
		// Check if key has expired
		if expiry, exists := m.ttls[key]; exists && time.Now().After(expiry) {
			continue
		}
		if _, exists := m.data[key]; exists {
			count++
		}
	}
	return redis.NewIntResult(count, nil)
}

// Expire sets a timeout on a key.
func (m *MockRedisClient) Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.trackCall("Expire")

	if _, exists := m.data[key]; exists {
		m.ttls[key] = time.Now().Add(expiration)
		return redis.NewBoolResult(true, nil)
	}
	return redis.NewBoolResult(false, nil)
}

// TTL returns the remaining time to live of a key.
func (m *MockRedisClient) TTL(ctx context.Context, key string) *redis.DurationCmd {
	m.mu.RLock()
	defer m.mu.RUnlock()

	m.trackCall("TTL")

	expiry, exists := m.ttls[key]
	if !exists {
		// Key doesn't have TTL or doesn't exist
		if _, keyExists := m.data[key]; keyExists {
			return redis.NewDurationResult(-1*time.Second, nil) // No TTL
		}
		return redis.NewDurationResult(-2*time.Second, nil) // Key doesn't exist
	}

	ttl := time.Until(expiry)
	if ttl < 0 {
		delete(m.data, key)
		delete(m.ttls, key)
		return redis.NewDurationResult(-2*time.Second, nil)
	}
	return redis.NewDurationResult(ttl, nil)
}

// Ping tests the connection (always succeeds for mock).
func (m *MockRedisClient) Ping(ctx context.Context) *redis.StatusCmd {
	m.trackCall("Ping")
	return redis.NewStatusResult("PONG", nil)
}

// Close closes the client (no-op for mock).
func (m *MockRedisClient) Close() error {
	m.trackCall("Close")
	return nil
}

// Reset clears all data in the mock store.
func (m *MockRedisClient) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.data = make(map[string]string)
	m.ttls = make(map[string]time.Time)
	m.calls = make(map[string]int)
}

// GetCallCount returns the number of times a method was called.
func (m *MockRedisClient) GetCallCount(method string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.calls[method]
}

// GetAllData returns all stored data (for testing/debugging).
func (m *MockRedisClient) GetAllData() map[string]string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make(map[string]string)
	for k, v := range m.data {
		// Skip expired keys
		if expiry, exists := m.ttls[k]; exists && time.Now().After(expiry) {
			continue
		}
		result[k] = v
	}
	return result
}

// trackCall increments the call counter for a method (must be called with lock held).
func (m *MockRedisClient) trackCall(method string) {
	m.calls[method]++
}

// SetWithoutLock sets a value without acquiring a lock (for internal use in tests).
func (m *MockRedisClient) SetWithoutLock(key, value string) {
	m.data[key] = value
}

// GetWithoutLock gets a value without acquiring a lock (for internal use in tests).
func (m *MockRedisClient) GetWithoutLock(key string) (string, bool) {
	val, exists := m.data[key]
	return val, exists
}
