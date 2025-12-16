package ids_test

import (
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/edvirons/ssp/shared/pkg/ids"
)

func TestNew_Format(t *testing.T) {
	id := ids.New("wo")

	// Should have prefix
	if !strings.HasPrefix(id, "wo_") {
		t.Errorf("expected prefix 'wo_', got %s", id)
	}

	// ULID part should be 26 characters
	parts := strings.Split(id, "_")
	if len(parts) != 2 {
		t.Fatalf("expected 2 parts, got %d", len(parts))
	}
	if len(parts[1]) != 26 {
		t.Errorf("expected ULID length 26, got %d for id %s", len(parts[1]), id)
	}
}

func TestNew_DifferentPrefixes(t *testing.T) {
	prefixes := []string{"wo", "proj", "inc", "part", "inv", "msg"}

	for _, prefix := range prefixes {
		id := ids.New(prefix)
		if !strings.HasPrefix(id, prefix+"_") {
			t.Errorf("expected prefix '%s_', got %s", prefix, id)
		}
	}
}

func TestNew_Uniqueness(t *testing.T) {
	seen := make(map[string]bool)
	for i := 0; i < 10000; i++ {
		id := ids.New("test")
		if seen[id] {
			t.Fatalf("duplicate ID generated: %s", id)
		}
		seen[id] = true
	}
}

func TestNew_ThreadSafety(t *testing.T) {
	var wg sync.WaitGroup
	idChan := make(chan string, 1000)

	// Generate IDs from 10 goroutines
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				idChan <- ids.New("conc")
			}
		}()
	}

	go func() {
		wg.Wait()
		close(idChan)
	}()

	// Collect and check uniqueness
	seen := make(map[string]bool)
	for id := range idChan {
		if seen[id] {
			t.Fatalf("duplicate ID in concurrent test: %s", id)
		}
		seen[id] = true
	}

	if len(seen) != 1000 {
		t.Errorf("expected 1000 unique IDs, got %d", len(seen))
	}
}

func TestNew_Sortability(t *testing.T) {
	// Generate IDs in sequence with small delay to ensure different timestamps
	var idList []string
	for i := 0; i < 10; i++ {
		idList = append(idList, ids.New("sort"))
		time.Sleep(time.Millisecond) // Small delay to ensure different timestamps
	}

	// Verify they're monotonically increasing (ULIDs are time-sortable)
	for i := 1; i < len(idList); i++ {
		if idList[i] <= idList[i-1] {
			t.Errorf("IDs not sorted: %s <= %s", idList[i], idList[i-1])
		}
	}
}

func TestParse_Valid(t *testing.T) {
	id := ids.New("wo")
	prefix, idPart, err := ids.Parse(id)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if prefix != "wo" {
		t.Errorf("expected prefix 'wo', got %s", prefix)
	}
	if len(idPart) != 26 {
		t.Errorf("expected ULID length 26, got %d", len(idPart))
	}
}

func TestParse_Invalid(t *testing.T) {
	_, _, err := ids.Parse("invalidid")
	if err == nil {
		t.Error("expected error for invalid ID without underscore")
	}
}

func TestExtractTimestamp_Valid(t *testing.T) {
	before := time.Now()
	id := ids.New("wo")
	after := time.Now()

	ts := ids.ExtractTimestamp(id)

	if ts.IsZero() {
		t.Error("expected non-zero timestamp")
	}

	// Timestamp should be between before and after
	if ts.Before(before.Add(-time.Second)) || ts.After(after.Add(time.Second)) {
		t.Errorf("timestamp %v outside expected range [%v, %v]", ts, before, after)
	}
}

func TestExtractTimestamp_OldFormat(t *testing.T) {
	// Old hex format should return zero time
	ts := ids.ExtractTimestamp("wo_a1b2c3d4e5f6g7h8")
	if !ts.IsZero() {
		t.Errorf("expected zero time for old format, got %v", ts)
	}
}

func TestExtractTimestamp_Invalid(t *testing.T) {
	ts := ids.ExtractTimestamp("invalid")
	if !ts.IsZero() {
		t.Errorf("expected zero time for invalid ID, got %v", ts)
	}
}

func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ids.New("bench")
	}
}

func BenchmarkNew_Parallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ids.New("bench")
		}
	})
}
