package lookups

import (
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Kind represents the type of lookup snapshot
type Kind string

const (
	KindSchool  Kind = "school"
	KindDevices Kind = "devices"
	KindParts   Kind = "parts"
)

// Store provides access to lookup data from SSOT snapshots
type Store struct {
	db *pgxpool.Pool
}

// New creates a new lookup store
func New(db *pgxpool.Pool) *Store { return &Store{db: db} }

// Snapshot represents a stored snapshot of lookup data
type Snapshot struct {
	TenantID  string
	Kind      Kind
	Version   string
	Payload   []byte
	UpdatedAt time.Time
}

// Common errors
var (
	ErrNotFound        = errors.New("not found")
	ErrSnapshotMissing = errors.New("ssot snapshot missing")
)
