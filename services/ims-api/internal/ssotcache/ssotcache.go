package ssotcache

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type Config struct {
	NATSURL      string
	SchoolURL    string
	DevicesURL   string
	PartsURL     string
	FetchTimeout time.Duration
}

func ConfigFromEnv() Config {
	c := Config{
		NATSURL:      env("NATS_URL", ""),
		SchoolURL:    env("SSOT_SCHOOL_URL", "http://ssot-school:8081"),
		DevicesURL:   env("SSOT_DEVICES_URL", "http://ssot-devices:8082"),
		PartsURL:     env("SSOT_PARTS_URL", "http://ssot-parts:8083"),
		FetchTimeout: 20 * time.Second,
	}
	if v := os.Getenv("SSOT_FETCH_TIMEOUT_SECONDS"); v != "" {
		if d, err := time.ParseDuration(v + "s"); err == nil {
			c.FetchTimeout = d
		}
	}
	return c
}

func Start(ctx context.Context, log *zap.Logger, db *pgxpool.Pool, c Config) error {
	if c.NATSURL == "" {
		log.Info("ssotcache disabled (NATS_URL not set)")
		return nil
	}

	nc, err := nats.Connect(c.NATSURL)
	if err != nil {
		return fmt.Errorf("nats connect: %w", err)
	}

	subscribe := func(kind string) error {
		subj := fmt.Sprintf("ssot.%s.changed", kind)
		_, err := nc.Subscribe(subj, func(m *nats.Msg) {
			tenant := tenantFromEvent(m.Data)
			if tenant == "" {
				log.Warn("ssot event missing tenantId", zap.String("subject", m.Subject))
				return
			}
			go func() {
				ctx2, cancel := context.WithTimeout(ctx, c.FetchTimeout)
				defer cancel()
				if err := FetchAndUpsert(ctx2, log, db, c, kind, tenant); err != nil {
					log.Error("ssot snapshot fetch failed", zap.String("kind", kind), zap.String("tenantId", tenant), zap.Error(err))
				}
			}()
		})
		return err
	}

	for _, k := range []string{"school", "devices", "parts"} {
		if err := subscribe(k); err != nil {
			return fmt.Errorf("subscribe %s: %w", k, err)
		}
	}

	go func() {
		<-ctx.Done()
		_ = nc.Drain()
	}()

	log.Info("ssotcache started", zap.String("nats", c.NATSURL))
	return nil
}

func FetchAndUpsert(ctx context.Context, log *zap.Logger, db *pgxpool.Pool, c Config, kind, tenant string) error {
	url := endpointForKind(c, kind)
	if url == "" {
		return fmt.Errorf("unknown kind: %s", kind)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url+"/v1/export", nil)
	if err != nil { return err }
	req.Header.Set("X-Tenant-Id", tenant)

	resp, err := http.DefaultClient.Do(req)
	if err != nil { return err }
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return fmt.Errorf("ssot export %s: status=%d body=%s", kind, resp.StatusCode, string(b))
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil { return err }

	var js any
	if err := json.Unmarshal(b, &js); err != nil {
		return fmt.Errorf("invalid ssot json: %w", err)
	}

	_, err = db.Exec(ctx, `
		INSERT INTO ims_ssot_snapshots (tenant_id, kind, version, payload, updated_at)
		VALUES ($1,$2,'1',$3::jsonb,$4)
		ON CONFLICT (tenant_id, kind) DO UPDATE SET payload=EXCLUDED.payload, updated_at=EXCLUDED.updated_at
	`, tenant, kind, string(b), time.Now().UTC())
	if err != nil { return err }

	log.Info("ssot snapshot updated", zap.String("kind", kind), zap.String("tenantId", tenant), zap.Int("bytes", len(b)))
	return nil
}

func endpointForKind(c Config, kind string) string {
	switch kind {
	case "school":
		return c.SchoolURL
	case "devices":
		return c.DevicesURL
	case "parts":
		return c.PartsURL
	default:
		return ""
	}
}

func tenantFromEvent(b []byte) string {
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil { return "" }
	if v, ok := m["tenantId"].(string); ok { return v }
	return ""
}

func env(k, d string) string {
	if v := os.Getenv(k); v != "" { return v }
	return d
}
