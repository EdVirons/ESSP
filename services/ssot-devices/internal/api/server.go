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
	"github.com/edvirons/ssp/ssot_devices/internal/models"
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
		httpx.WriteJSON(w, 200, map[string]any{"ok": true, "service": "ssot-devices"})
	})

	r.Route("/v1", func(r chi.Router) {
		r.Get("/export", s.Export)
		r.Post("/import", s.Import)

		// Network identity (MAC address) endpoints
		r.Get("/devices/{deviceId}/network-identities", s.ListNetworkIdentities)
		r.Post("/devices/{deviceId}/network-identities", s.UpsertNetworkIdentity)
		r.Delete("/devices/{deviceId}/network-identities/{id}", s.DeleteNetworkIdentity)

		// MAC lookup
		r.Get("/lookup/mac/{mac}", s.LookupByMAC)
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

	_ = s.pub.PublishJSON("ssot.devices.snapshot", map[string]any{
		"tenantId":   tenant,
		"importedAt": time.Now().UTC(),
		"counts":     res,
	})

	// Also publish a lightweight "changed" signal
	_ = s.pub.PublishJSON("ssot.devices.changed", map[string]any{"tenantId": tenant, "at": time.Now().UTC()})

	httpx.WriteJSON(w, 200, map[string]any{"ok": true, "imported": res})
}

// ---------- Network Identity Handlers ----------

func (s *Server) ListNetworkIdentities(w http.ResponseWriter, r *http.Request) {
	tenant := httpx.TenantID(r)
	if tenant == "" {
		httpx.Error(w, 400, "X-Tenant-Id required")
		return
	}
	deviceID := chi.URLParam(r, "deviceId")
	if deviceID == "" {
		httpx.Error(w, 400, "deviceId required")
		return
	}

	rows, err := s.db.Query(r.Context(), `
		SELECT id, tenant_id, device_id, mac_address, interface_name, interface_type, is_primary, first_seen_at, last_seen_at, created_at, updated_at
		FROM device_network_identities
		WHERE tenant_id=$1 AND device_id=$2
		ORDER BY is_primary DESC, created_at ASC
	`, tenant, deviceID)
	if err != nil {
		httpx.Error(w, 500, "query failed")
		return
	}
	defer rows.Close()

	var items []models.DeviceNetworkIdentity
	for rows.Next() {
		var x models.DeviceNetworkIdentity
		var ifType string
		if err := rows.Scan(&x.ID, &x.TenantID, &x.DeviceID, &x.MACAddress, &x.InterfaceName, &ifType, &x.IsPrimary, &x.FirstSeenAt, &x.LastSeenAt, &x.CreatedAt, &x.UpdatedAt); err != nil {
			httpx.Error(w, 500, "scan failed")
			return
		}
		x.InterfaceType = models.InterfaceType(ifType)
		items = append(items, x)
	}
	httpx.WriteJSON(w, 200, map[string]any{"items": items})
}

func (s *Server) UpsertNetworkIdentity(w http.ResponseWriter, r *http.Request) {
	tenant := httpx.TenantID(r)
	if tenant == "" {
		httpx.Error(w, 400, "X-Tenant-Id required")
		return
	}
	deviceID := chi.URLParam(r, "deviceId")
	if deviceID == "" {
		httpx.Error(w, 400, "deviceId required")
		return
	}

	var body map[string]any
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		httpx.Error(w, 400, "invalid json")
		return
	}

	mac := trim(body["macAddress"])
	if mac == "" {
		httpx.Error(w, 400, "macAddress required")
		return
	}
	mac = models.NormalizeMACAddress(mac)

	id := trim(body["id"])
	if id == "" {
		id = newID("netid")
	}
	ifName := trim(body["interfaceName"])
	ifType := trim(body["interfaceType"])
	if ifType == "" {
		ifType = "unknown"
	}
	isPrimary := false
	if v, ok := body["isPrimary"].(bool); ok {
		isPrimary = v
	}

	now := time.Now().UTC()
	_, err := s.db.Exec(r.Context(), `
		INSERT INTO device_network_identities (id, tenant_id, device_id, mac_address, interface_name, interface_type, is_primary, first_seen_at, last_seen_at, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$8,$8,$8)
		ON CONFLICT (tenant_id, mac_address) DO UPDATE SET
			device_id=EXCLUDED.device_id, interface_name=EXCLUDED.interface_name, interface_type=EXCLUDED.interface_type,
			is_primary=EXCLUDED.is_primary, last_seen_at=$8, updated_at=$8
	`, id, tenant, deviceID, mac, ifName, ifType, isPrimary, now)
	if err != nil {
		s.log.Error("upsert network identity failed", zap.Error(err))
		httpx.Error(w, 500, "upsert failed")
		return
	}

	// Publish event for sync
	_ = s.pub.PublishJSON("ssot.devices.network.changed", map[string]any{
		"tenantId": tenant, "deviceId": deviceID, "macAddress": mac, "action": "upsert", "at": now,
	})

	httpx.WriteJSON(w, 200, map[string]any{"ok": true, "id": id, "macAddress": mac})
}

func (s *Server) DeleteNetworkIdentity(w http.ResponseWriter, r *http.Request) {
	tenant := httpx.TenantID(r)
	if tenant == "" {
		httpx.Error(w, 400, "X-Tenant-Id required")
		return
	}
	deviceID := chi.URLParam(r, "deviceId")
	id := chi.URLParam(r, "id")
	if deviceID == "" || id == "" {
		httpx.Error(w, 400, "deviceId and id required")
		return
	}

	// Get MAC for event before delete
	var mac string
	_ = s.db.QueryRow(r.Context(), `SELECT mac_address FROM device_network_identities WHERE tenant_id=$1 AND id=$2`, tenant, id).Scan(&mac)

	result, err := s.db.Exec(r.Context(), `DELETE FROM device_network_identities WHERE tenant_id=$1 AND id=$2 AND device_id=$3`, tenant, id, deviceID)
	if err != nil {
		httpx.Error(w, 500, "delete failed")
		return
	}
	if result.RowsAffected() == 0 {
		httpx.Error(w, 404, "not found")
		return
	}

	// Publish event
	if mac != "" {
		_ = s.pub.PublishJSON("ssot.devices.network.changed", map[string]any{
			"tenantId": tenant, "deviceId": deviceID, "macAddress": mac, "action": "deleted", "at": time.Now().UTC(),
		})
	}

	httpx.WriteJSON(w, 200, map[string]any{"ok": true})
}

func (s *Server) LookupByMAC(w http.ResponseWriter, r *http.Request) {
	tenant := httpx.TenantID(r)
	if tenant == "" {
		httpx.Error(w, 400, "X-Tenant-Id required")
		return
	}
	mac := chi.URLParam(r, "mac")
	if mac == "" {
		httpx.Error(w, 400, "mac required")
		return
	}
	mac = models.NormalizeMACAddress(mac)

	var netId models.DeviceNetworkIdentity
	var ifType string
	err := s.db.QueryRow(r.Context(), `
		SELECT id, tenant_id, device_id, mac_address, interface_name, interface_type, is_primary, first_seen_at, last_seen_at, created_at, updated_at
		FROM device_network_identities
		WHERE tenant_id=$1 AND mac_address=$2
	`, tenant, mac).Scan(&netId.ID, &netId.TenantID, &netId.DeviceID, &netId.MACAddress, &netId.InterfaceName, &ifType, &netId.IsPrimary, &netId.FirstSeenAt, &netId.LastSeenAt, &netId.CreatedAt, &netId.UpdatedAt)
	if err == pgx.ErrNoRows {
		httpx.Error(w, 404, "MAC not found")
		return
	}
	if err != nil {
		httpx.Error(w, 500, "query failed")
		return
	}
	netId.InterfaceType = models.InterfaceType(ifType)

	// Also fetch the device
	var device models.Device
	err = s.db.QueryRow(r.Context(), `
		SELECT id, tenant_id, serial, asset_tag, device_model_id, school_id, assigned_to, lifecycle, enrolled, created_at, updated_at
		FROM devices WHERE tenant_id=$1 AND id=$2
	`, tenant, netId.DeviceID).Scan(&device.ID, &device.TenantID, &device.Serial, &device.AssetTag, &device.DeviceModelID, &device.SchoolID, &device.AssignedTo, &device.Lifecycle, &device.Enrolled, &device.CreatedAt, &device.UpdatedAt)
	if err != nil && err != pgx.ErrNoRows {
		httpx.Error(w, 500, "device query failed")
		return
	}

	httpx.WriteJSON(w, 200, map[string]any{"networkIdentity": netId, "device": device})
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
// export_import_devices.go

// keep models referenced (avoid unused import in some IDEs)
var _ = models.ExportPayload{}
