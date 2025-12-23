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
	"github.com/edvirons/ssp/ssot_school/internal/models"
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
		httpx.WriteJSON(w, 200, map[string]any{"ok": true, "service": "ssot-school"})
	})

	r.Route("/v1", func(r chi.Router) {
		r.Get("/export", s.Export)
		r.Post("/import", s.Import)
		r.Get("/schools", s.ListSchools)
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

	_ = s.pub.PublishJSON("ssot.school.snapshot", map[string]any{
		"tenantId":   tenant,
		"importedAt": time.Now().UTC(),
		"counts":     res,
	})

	// Also publish a lightweight "changed" signal
	_ = s.pub.PublishJSON("ssot.school.changed", map[string]any{"tenantId": tenant, "at": time.Now().UTC()})

	httpx.WriteJSON(w, 200, map[string]any{"ok": true, "imported": res})
}

// SchoolItem is the response format for schools list (matches IMS SchoolSnapshot)
type SchoolItem struct {
	TenantID      string    `json:"tenantId"`
	SchoolID      string    `json:"schoolId"`
	Name          string    `json:"name"`
	CountyCode    string    `json:"countyCode"`
	CountyName    string    `json:"countyName"`
	SubCountyCode string    `json:"subCountyCode"`
	SubCountyName string    `json:"subCountyName"`
	Level         string    `json:"level"`
	Type          string    `json:"type"`
	KnecCode      string    `json:"knecCode"`
	Uic           string    `json:"uic"`
	Sex           string    `json:"sex"`
	Cluster       string    `json:"cluster"`
	Accommodation string    `json:"accommodation"`
	Latitude      float64   `json:"latitude"`
	Longitude     float64   `json:"longitude"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type SchoolsPage struct {
	Items      []SchoolItem `json:"items"`
	NextCursor string       `json:"nextCursor"`
}

func (s *Server) ListSchools(w http.ResponseWriter, r *http.Request) {
	tenant := httpx.TenantID(r)
	if tenant == "" {
		httpx.Error(w, 400, "X-Tenant-Id required")
		return
	}

	// Query schools with county/subcounty names
	rows, err := s.db.Query(r.Context(), `
		SELECT
			s.id, s.tenant_id, s.name,
			COALESCE(c.code, ''), COALESCE(c.name, ''),
			COALESCE(sc.code, ''), COALESCE(sc.name, ''),
			s.level, s.type, s.knec_code, s.uic, s.sex, s.cluster, s.accommodation,
			s.latitude, s.longitude, s.updated_at
		FROM schools s
		LEFT JOIN counties c ON s.county_id = c.id AND c.tenant_id = s.tenant_id
		LEFT JOIN sub_counties sc ON s.sub_county_id = sc.id AND sc.tenant_id = s.tenant_id
		WHERE s.tenant_id = $1
		ORDER BY s.name
	`, tenant)
	if err != nil {
		s.log.Error("failed to query schools", zap.Error(err))
		httpx.Error(w, 500, "query failed")
		return
	}
	defer rows.Close()

	items := []SchoolItem{}
	for rows.Next() {
		var it SchoolItem
		if err := rows.Scan(&it.SchoolID, &it.TenantID, &it.Name,
			&it.CountyCode, &it.CountyName, &it.SubCountyCode, &it.SubCountyName,
			&it.Level, &it.Type, &it.KnecCode, &it.Uic, &it.Sex, &it.Cluster, &it.Accommodation,
			&it.Latitude, &it.Longitude, &it.UpdatedAt); err != nil {
			s.log.Error("failed to scan school", zap.Error(err))
			continue
		}
		items = append(items, it)
	}

	httpx.WriteJSON(w, 200, SchoolsPage{Items: items, NextCursor: ""})
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
// export_import_school.go

// keep models referenced (avoid unused import in some IDEs)
var _ = models.ExportPayload{}
