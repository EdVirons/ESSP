package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/edvirons/ssp/ims/internal/config"
	"github.com/edvirons/ssp/ims/internal/middleware"
	"github.com/edvirons/ssp/ims/internal/models"
	"github.com/edvirons/ssp/ims/internal/store"
	"go.uber.org/zap"
)

type SSOTWebhookHandler struct {
	cfg config.Config
	log *zap.Logger
	pg  *store.Postgres
}

func NewSSOTWebhookHandler(cfg config.Config, log *zap.Logger, pg *store.Postgres) *SSOTWebhookHandler {
	return &SSOTWebhookHandler{cfg: cfg, log: log, pg: pg}
}

// Generic envelope supports either a single item or a batch
type envelope[T any] struct {
	Item  *T   `json:"item,omitempty"`
	Items []T  `json:"items,omitempty"`
}

func (h *SSOTWebhookHandler) Schools(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())
	var env envelope[models.SchoolSnapshot]
	if err := json.NewDecoder(r.Body).Decode(&env); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	items := env.Items
	if env.Item != nil {
		items = append(items, *env.Item)
	}
	maxSeen := time.Unix(0,0).UTC()
	for _, it := range items {
		it.TenantID = tenant
		_ = h.pg.SchoolsSnapshot().Upsert(r.Context(), it)
		if it.UpdatedAt.After(maxSeen) { maxSeen = it.UpdatedAt }
	}
	if len(items) > 0 {
		st, err := h.pg.SSOTState().Get(r.Context(), tenant, store.SSOTSchools)
		if err != nil { st = store.NewSSOTSyncState(tenant, store.SSOTSchools) }
		if maxSeen.After(st.LastUpdatedSince) { st.LastUpdatedSince = maxSeen }
		st.LastCursor = "" // push updates invalidate cursor
		st.UpdatedAt = time.Now().UTC()
		_ = h.pg.SSOTState().Upsert(r.Context(), st)
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "upserts": len(items)})
}

func (h *SSOTWebhookHandler) Devices(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())
	var env envelope[models.DeviceSnapshot]
	if err := json.NewDecoder(r.Body).Decode(&env); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	items := env.Items
	if env.Item != nil { items = append(items, *env.Item) }
	maxSeen := time.Unix(0,0).UTC()
	for _, it := range items {
		it.TenantID = tenant
		_ = h.pg.DevicesSnapshot().Upsert(r.Context(), it)
		if it.UpdatedAt.After(maxSeen) { maxSeen = it.UpdatedAt }
	}
	if len(items) > 0 {
		st, err := h.pg.SSOTState().Get(r.Context(), tenant, store.SSOTDevices)
		if err != nil { st = store.NewSSOTSyncState(tenant, store.SSOTDevices) }
		if maxSeen.After(st.LastUpdatedSince) { st.LastUpdatedSince = maxSeen }
		st.LastCursor = ""
		st.UpdatedAt = time.Now().UTC()
		_ = h.pg.SSOTState().Upsert(r.Context(), st)
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "upserts": len(items)})
}

func (h *SSOTWebhookHandler) Parts(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())
	var env envelope[models.PartSnapshot]
	if err := json.NewDecoder(r.Body).Decode(&env); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	items := env.Items
	if env.Item != nil { items = append(items, *env.Item) }
	maxSeen := time.Unix(0,0).UTC()
	for _, it := range items {
		it.TenantID = tenant
		_ = h.pg.PartsSnapshot().Upsert(r.Context(), it)
		if it.UpdatedAt.After(maxSeen) { maxSeen = it.UpdatedAt }
	}
	if len(items) > 0 {
		st, err := h.pg.SSOTState().Get(r.Context(), tenant, store.SSOTParts)
		if err != nil { st = store.NewSSOTSyncState(tenant, store.SSOTParts) }
		if maxSeen.After(st.LastUpdatedSince) { st.LastUpdatedSince = maxSeen }
		st.LastCursor = ""
		st.UpdatedAt = time.Now().UTC()
		_ = h.pg.SSOTState().Upsert(r.Context(), st)
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "upserts": len(items)})
}
