package api

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/edvirons/ssp/shared/pkg/eventbus"
	"github.com/edvirons/ssp/shared/pkg/httpx"
	"github.com/edvirons/ssp/shared/pkg/ids"
	"github.com/edvirons/ssp/ssot_parts/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type Server struct {
	log *zap.Logger
	db  *pgxpool.Pool
	pub *eventbus.Publisher
}

func NewServer(log *zap.Logger) http.Handler {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		dbURL = os.Getenv("DEFAULT_DB_URL")
	}
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatal("db connect failed", zap.Error(err))
	}

	nc, err := nats.Connect(env("NATS_URL", "nats://localhost:4222"))
	if err != nil {
		log.Fatal("nats connect failed", zap.Error(err))
	}
	pub := eventbus.NewPublisher(nc)

	s := &Server{log: log, db: pool, pub: pub}

	r := chi.NewRouter()
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		httpx.WriteJSON(w, 200, map[string]any{"ok": true, "service": "ssot-parts"})
	})

	r.Route("/v1", func(r chi.Router) {
		r.Get("/export", s.Export)
		r.Post("/import", s.Import)
	})

	log.Info("routes ready")
	return r
}

func (s *Server) Export(w http.ResponseWriter, r *http.Request) {
	tenant := httpx.TenantID(r)
	if tenant == "" {
		httpx.Error(w, 400, "X-Tenant-Id required")
		return
	}

	payload, err := exportAll(r.Context(), s.db, tenant)
	if err != nil {
		httpx.Error(w, 500, "export failed")
		return
	}
	httpx.WriteJSON(w, 200, payload)
}

func (s *Server) Import(w http.ResponseWriter, r *http.Request) {
	tenant := httpx.TenantID(r)
	if tenant == "" {
		httpx.Error(w, 400, "X-Tenant-Id required")
		return
	}

	var body map[string]any
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		httpx.Error(w, 400, "invalid json")
		return
	}

	res, err := importAll(r.Context(), s.db, tenant, body)
	if err != nil {
		httpx.Error(w, 400, err.Error())
		return
	}

	_ = s.pub.PublishJSON("ssot.parts.snapshot", map[string]any{
		"tenantId":   tenant,
		"importedAt": time.Now().UTC(),
		"counts":     res,
	})

	// Also publish a lightweight "changed" signal
	_ = s.pub.PublishJSON("ssot.parts.changed", map[string]any{"tenantId": tenant, "at": time.Now().UTC()})

	httpx.WriteJSON(w, 200, map[string]any{"ok": true, "imported": res})
}

// ---------- helpers ----------
func env(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}

func withTx(ctx context.Context, db *pgxpool.Pool, fn func(pgx.Tx) error) error {
	tx, err := db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	if err := fn(tx); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func newID(prefix string) string { return ids.New(prefix) }

func trim(v any) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return strings.TrimSpace(s)
	}
	return ""
}

// exportAll/importAll are implemented in service-specific file:
// export_import_parts.go

// keep models referenced (avoid unused import in some IDEs)
var _ = models.ExportPayload{}
