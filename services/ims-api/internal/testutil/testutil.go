package testutil

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// MustGetEnv retrieves an environment variable or fails the test if not found.
func MustGetEnv(t *testing.T, key string) string {
	t.Helper()
	val := os.Getenv(key)
	if val == "" {
		t.Fatalf("environment variable %s is required for this test", key)
	}
	return val
}

// GetEnvOrDefault retrieves an environment variable or returns a default value.
func GetEnvOrDefault(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

// WaitForCondition polls a condition function until it returns true or times out.
func WaitForCondition(t *testing.T, condition func() bool, timeout time.Duration, message string) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			t.Fatalf("timeout waiting for condition: %s", message)
		case <-ticker.C:
			if condition() {
				return
			}
		}
	}
}

// RequireNoError is a helper that fails the test if err is not nil.
func RequireNoError(t *testing.T, err error, msgAndArgs ...interface{}) {
	t.Helper()
	require.NoError(t, err, msgAndArgs...)
}

// RequireEqual is a helper that fails the test if expected != actual.
func RequireEqual(t *testing.T, expected, actual interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	require.Equal(t, expected, actual, msgAndArgs...)
}

// AssertEventually retries an assertion for a specified duration.
func AssertEventually(t *testing.T, condition func() bool, timeout time.Duration, tick time.Duration, msg string) {
	t.Helper()
	require.Eventually(t, condition, timeout, tick, msg)
}

// SkipIfShort skips the test if testing.Short() is true.
func SkipIfShort(t *testing.T) {
	t.Helper()
	if testing.Short() {
		t.Skip("Skipping test in short mode")
	}
}

// SkipIfNoIntegration skips the test if INTEGRATION_TEST env var is not set.
func SkipIfNoIntegration(t *testing.T) {
	t.Helper()
	if os.Getenv("INTEGRATION_TEST") == "" {
		t.Skip("Skipping integration test (set INTEGRATION_TEST=1 to run)")
	}
}

// GenerateTestID generates a unique test ID based on the test name and timestamp.
func GenerateTestID(t *testing.T) string {
	t.Helper()
	return t.Name() + "_" + time.Now().Format("20060102150405")
}

// Cleanup registers a cleanup function that will be called when the test completes.
func Cleanup(t *testing.T, fn func()) {
	t.Helper()
	t.Cleanup(fn)
}
