package lookups

import (
	"context"
	"time"
)

// GetSnapshot retrieves a snapshot from the database
func (s *Store) GetSnapshot(ctx context.Context, tenant string, kind Kind) (*Snapshot, error) {
	var version string
	var payload []byte
	var updated time.Time

	err := s.db.QueryRow(ctx, `
		SELECT version, payload::text, updated_at
		FROM ims_ssot_snapshots
		WHERE tenant_id=$1 AND kind=$2
	`, tenant, string(kind)).Scan(&version, &payload, &updated)
	if err != nil {
		return nil, err
	}

	return &Snapshot{
		TenantID:  tenant,
		Kind:      kind,
		Version:   version,
		Payload:   payload,
		UpdatedAt: updated,
	}, nil
}
