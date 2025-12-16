package store

import "github.com/edvirons/ssp/shared/pkg/ids"

// NewID generates a new prefixed ULID.
// Delegates to shared ids package for consistent ID generation across all services.
// Format: {prefix}_{26-char-ulid}
// Example: wo_01ARZ3NDEKTSV4RRFFQ69G5FAV
func NewID(prefix string) string {
	return ids.New(prefix)
}
