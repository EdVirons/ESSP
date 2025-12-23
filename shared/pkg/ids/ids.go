package ids

import (
	"crypto/rand"
	"fmt"
	"sync"
	"time"

	"github.com/oklog/ulid/v2"
)

// entropyPool provides thread-safe access to entropy sources.
// Using a pool avoids lock contention in high-throughput scenarios.
var entropyPool = sync.Pool{
	New: func() interface{} {
		return ulid.Monotonic(rand.Reader, 0)
	},
}

// New generates a new prefixed ULID.
// Format: {prefix}_{26-char-ulid}
// Example: wo_01ARZ3NDEKTSV4RRFFQ69G5FAV
//
// ULIDs are:
// - Lexicographically sortable
// - Time-ordered (first 48 bits are timestamp)
// - Cryptographically random (last 80 bits)
func New(prefix string) string {
	entropy, _ := entropyPool.Get().(ulid.MonotonicReader)
	defer entropyPool.Put(entropy)

	id := ulid.MustNew(ulid.Timestamp(time.Now()), entropy)
	return prefix + "_" + id.String()
}

// Parse extracts prefix and ID from a prefixed ID string.
// Returns error if the ID format is invalid.
func Parse(prefixedID string) (prefix string, id string, err error) {
	for i := 0; i < len(prefixedID); i++ {
		if prefixedID[i] == '_' {
			return prefixedID[:i], prefixedID[i+1:], nil
		}
	}
	return "", "", fmt.Errorf("invalid ID format: missing underscore separator")
}

// ExtractTimestamp extracts the timestamp from a ULID-formatted ID.
// Returns zero time for old hex-format IDs or invalid IDs.
func ExtractTimestamp(prefixedID string) time.Time {
	_, idPart, err := Parse(prefixedID)
	if err != nil || len(idPart) != 26 {
		return time.Time{} // Old format or invalid
	}

	parsed, err := ulid.Parse(idPart)
	if err != nil {
		return time.Time{}
	}

	return ulid.Time(parsed.Time())
}
