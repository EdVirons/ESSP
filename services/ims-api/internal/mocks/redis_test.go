package mocks_test

import (
	"context"
	"testing"
	"time"

	"github.com/edvirons/ssp/ims/internal/mocks"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Example test demonstrating MockRedisClient usage
func TestMockRedisClient_GetSet(t *testing.T) {
	ctx := context.Background()
	client := mocks.NewMockRedisClient()

	// Set a value
	cmd := client.Set(ctx, "key1", "value1", 0)
	require.NoError(t, cmd.Err())

	// Get the value
	getCmd := client.Get(ctx, "key1")
	require.NoError(t, getCmd.Err())
	assert.Equal(t, "value1", getCmd.Val())
}

func TestMockRedisClient_Expiration(t *testing.T) {
	ctx := context.Background()
	client := mocks.NewMockRedisClient()

	// Set a value with short expiration
	cmd := client.Set(ctx, "temp-key", "temp-value", 50*time.Millisecond)
	require.NoError(t, cmd.Err())

	// Value should exist immediately
	getCmd := client.Get(ctx, "temp-key")
	require.NoError(t, getCmd.Err())
	assert.Equal(t, "temp-value", getCmd.Val())

	// Wait for expiration
	time.Sleep(100 * time.Millisecond)

	// Value should be expired
	getCmd = client.Get(ctx, "temp-key")
	assert.Equal(t, redis.Nil, getCmd.Err())
}

func TestMockRedisClient_Delete(t *testing.T) {
	ctx := context.Background()
	client := mocks.NewMockRedisClient()

	// Set multiple values
	client.Set(ctx, "key1", "value1", 0)
	client.Set(ctx, "key2", "value2", 0)

	// Delete one key
	delCmd := client.Del(ctx, "key1")
	require.NoError(t, delCmd.Err())
	assert.Equal(t, int64(1), delCmd.Val())

	// Verify deletion
	getCmd := client.Get(ctx, "key1")
	assert.Equal(t, redis.Nil, getCmd.Err())

	// Other key should still exist
	getCmd = client.Get(ctx, "key2")
	require.NoError(t, getCmd.Err())
	assert.Equal(t, "value2", getCmd.Val())
}

func TestMockRedisClient_CallTracking(t *testing.T) {
	ctx := context.Background()
	client := mocks.NewMockRedisClient()

	// Perform operations
	client.Get(ctx, "key1")
	client.Set(ctx, "key1", "value1", 0)
	client.Get(ctx, "key1")
	client.Del(ctx, "key1")

	// Verify call counts
	assert.Equal(t, 2, client.GetCallCount("Get"))
	assert.Equal(t, 1, client.GetCallCount("Set"))
	assert.Equal(t, 1, client.GetCallCount("Del"))
}

func TestMockRedisClient_Reset(t *testing.T) {
	ctx := context.Background()
	client := mocks.NewMockRedisClient()

	// Set some data
	client.Set(ctx, "key1", "value1", 0)
	client.Get(ctx, "key1")

	// Reset
	client.Reset()

	// Data should be cleared
	getCmd := client.Get(ctx, "key1")
	assert.Equal(t, redis.Nil, getCmd.Err())

	// Call count should be reset (Get count includes the one after Reset)
	assert.Equal(t, 1, client.GetCallCount("Get"))
}
