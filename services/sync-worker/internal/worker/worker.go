// Package worker provides the SSOT sync worker implementation.
package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"time"

	"github.com/edvirons/ssp/sync_worker/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

// SyncWorker handles SSOT synchronization events.
type SyncWorker struct {
	log    *zap.Logger
	db     *pgxpool.Pool
	config config.Config
}

// New creates a new SyncWorker.
func New(log *zap.Logger, db *pgxpool.Pool, cfg config.Config) *SyncWorker {
	return &SyncWorker{
		log:    log,
		db:     db,
		config: cfg,
	}
}

// HandleEvent returns a NATS message handler for the given kind.
func (sw *SyncWorker) HandleEvent(kind string) nats.MsgHandler {
	return func(m *nats.Msg) {
		tenant := extractTenantID(m.Data)
		if tenant == "" {
			sw.log.Warn("event missing tenantId", zap.String("subject", m.Subject))
			return
		}

		sw.log.Info("event received",
			zap.String("subject", m.Subject),
			zap.String("kind", kind),
			zap.String("tenantId", tenant))

		// Process in background with retry logic
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), sw.config.FetchTimeout)
			defer cancel()

			if err := sw.fetchAndUpsertWithRetry(ctx, kind, tenant); err != nil {
				sw.log.Error("failed to sync after retries",
					zap.String("kind", kind),
					zap.String("tenantId", tenant),
					zap.Error(err))
			}
		}()
	}
}

// fetchAndUpsertWithRetry attempts to sync with exponential backoff retry.
func (sw *SyncWorker) fetchAndUpsertWithRetry(ctx context.Context, kind, tenant string) error {
	var lastErr error

	for attempt := 0; attempt <= sw.config.MaxRetries; attempt++ {
		if attempt > 0 {
			// Exponential backoff: 1s, 2s, 4s, 8s...
			backoff := time.Duration(math.Pow(2, float64(attempt-1))) * sw.config.InitialBackoff
			sw.log.Info("retrying after backoff",
				zap.String("kind", kind),
				zap.String("tenantId", tenant),
				zap.Int("attempt", attempt),
				zap.Duration("backoff", backoff))

			select {
			case <-time.After(backoff):
			case <-ctx.Done():
				return fmt.Errorf("context cancelled during backoff: %w", ctx.Err())
			}
		}

		err := sw.fetchAndUpsert(ctx, kind, tenant)
		if err == nil {
			if attempt > 0 {
				sw.log.Info("sync succeeded after retry",
					zap.String("kind", kind),
					zap.String("tenantId", tenant),
					zap.Int("attempt", attempt))
			}
			return nil
		}

		lastErr = err
		sw.log.Warn("sync attempt failed",
			zap.String("kind", kind),
			zap.String("tenantId", tenant),
			zap.Int("attempt", attempt),
			zap.Error(err))
	}

	return fmt.Errorf("max retries exceeded: %w", lastErr)
}

// fetchAndUpsert fetches data from SSOT and upserts into the database.
func (sw *SyncWorker) fetchAndUpsert(ctx context.Context, kind, tenant string) error {
	url := sw.getURLForKind(kind)
	if url == "" {
		return fmt.Errorf("unknown kind: %s", kind)
	}

	// Create request to SSOT export endpoint
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url+"/v1/export", nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("X-Tenant-Id", tenant)

	// Fetch export from SSOT service
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return fmt.Errorf("ssot export failed: status=%d body=%s", resp.StatusCode, string(body))
	}

	// Read response body
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}

	// Validate JSON
	var js any
	if err := json.Unmarshal(data, &js); err != nil {
		return fmt.Errorf("invalid json response: %w", err)
	}

	// Upsert into ims_ssot_snapshots table
	_, err = sw.db.Exec(ctx, `
		INSERT INTO ims_ssot_snapshots (tenant_id, kind, version, payload, updated_at)
		VALUES ($1, $2, '1', $3::jsonb, $4)
		ON CONFLICT (tenant_id, kind)
		DO UPDATE SET payload = EXCLUDED.payload, updated_at = EXCLUDED.updated_at
	`, tenant, kind, string(data), time.Now().UTC())

	if err != nil {
		return fmt.Errorf("database upsert failed: %w", err)
	}

	sw.log.Info("ssot snapshot synced",
		zap.String("kind", kind),
		zap.String("tenantId", tenant),
		zap.Int("bytes", len(data)))

	return nil
}

// getURLForKind returns the SSOT service URL for the given kind.
func (sw *SyncWorker) getURLForKind(kind string) string {
	switch kind {
	case "school":
		return sw.config.SchoolURL
	case "devices":
		return sw.config.DevicesURL
	case "parts":
		return sw.config.PartsURL
	default:
		return ""
	}
}

// extractTenantID extracts the tenantId from a JSON message.
func extractTenantID(data []byte) string {
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return ""
	}
	if v, ok := m["tenantId"].(string); ok {
		return v
	}
	return ""
}
