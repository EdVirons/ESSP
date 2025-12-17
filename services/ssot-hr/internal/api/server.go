package api

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/edvirons/ssp/shared/pkg/eventbus"
	"github.com/edvirons/ssp/shared/pkg/httpx"
	"github.com/edvirons/ssp/shared/pkg/ids"
	"github.com/edvirons/ssp/ssot_hr/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Server struct {
	log   *zap.Logger
	db    *pgxpool.Pool
	pub   *eventbus.Publisher
	cache *redis.Client
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

	// Valkey/Redis cache (optional)
	var cache *redis.Client
	if addr := os.Getenv("VALKEY_ADDR"); addr != "" {
		cache = redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: os.Getenv("VALKEY_PASSWORD"),
			DB:       0,
		})
	}

	s := &Server{log: log, db: pool, pub: pub, cache: cache}

	r := chi.NewRouter()
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		httpx.WriteJSON(w, 200, map[string]any{"ok": true, "service": "ssot-hr"})
	})

	r.Route("/v1", func(r chi.Router) {
		// Export/Import (SSOT standard)
		r.Get("/export", s.Export)
		r.Post("/import", s.Import)

		// People CRUD
		r.Route("/people", func(r chi.Router) {
			r.Post("/", s.CreatePerson)
			r.Get("/", s.ListPeople)
			r.Get("/{id}", s.GetPerson)
			r.Patch("/{id}", s.PatchPerson)
			r.Delete("/{id}", s.DeletePerson)
		})

		// Teams CRUD
		r.Route("/teams", func(r chi.Router) {
			r.Post("/", s.CreateTeam)
			r.Get("/", s.ListTeams)
			r.Get("/{id}", s.GetTeam)
			r.Patch("/{id}", s.PatchTeam)
			r.Delete("/{id}", s.DeleteTeam)
		})

		// Team Memberships
		r.Route("/team-memberships", func(r chi.Router) {
			r.Post("/", s.CreateMembership)
			r.Get("/", s.ListMemberships)
			r.Delete("/{id}", s.DeleteMembership)
		})

		// Org Units
		r.Route("/org-units", func(r chi.Router) {
			r.Post("/", s.CreateOrgUnit)
			r.Get("/", s.ListOrgUnits)
			r.Get("/tree", s.GetOrgTree)
			r.Get("/{id}", s.GetOrgUnit)
			r.Patch("/{id}", s.PatchOrgUnit)
			r.Delete("/{id}", s.DeleteOrgUnit)
		})
	})

	log.Info("routes ready")
	return r
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
	defer tx.Rollback(ctx)
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

func parseJSON(data string) map[string]any {
	if data == "" || data == "{}" {
		return nil
	}
	var m map[string]any
	if err := json.Unmarshal([]byte(data), &m); err != nil {
		return nil
	}
	return m
}

// keep models referenced
var _ = models.ExportPayload{}
