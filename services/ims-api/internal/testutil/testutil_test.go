package testutil_test

import (
	"testing"
	"time"

	"github.com/edvirons/ssp/ims/internal/testutil"
)

// Example test demonstrating testutil usage
func TestGenerateTestID(t *testing.T) {
	id1 := testutil.GenerateTestID(t)
	time.Sleep(10 * time.Millisecond)
	id2 := testutil.GenerateTestID(t)

	if id1 == id2 {
		t.Error("GenerateTestID should generate unique IDs")
	}
}

func TestWaitForCondition(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping in short mode")
	}

	counter := 0
	condition := func() bool {
		counter++
		return counter >= 3
	}

	testutil.WaitForCondition(t, condition, 1*time.Second, "counter to reach 3")

	if counter < 3 {
		t.Errorf("expected counter >= 3, got %d", counter)
	}
}

func TestGetEnvOrDefault(t *testing.T) {
	// Test default value when env var doesn't exist
	val := testutil.GetEnvOrDefault("NON_EXISTENT_VAR_12345", "default")
	if val != "default" {
		t.Errorf("expected 'default', got '%s'", val)
	}

	// Set an env var and test it's retrieved
	t.Setenv("TEST_VAR_12345", "test-value")
	val = testutil.GetEnvOrDefault("TEST_VAR_12345", "default")
	if val != "test-value" {
		t.Errorf("expected 'test-value', got '%s'", val)
	}
}
