package handlers

import (
	"net/http"
	"time"

	"github.com/edvirons/ssp/ims/internal/config"
	"github.com/edvirons/ssp/ims/internal/middleware"
	"github.com/edvirons/ssp/ims/internal/ssot"
	"github.com/edvirons/ssp/ims/internal/store"
	"go.uber.org/zap"
)

type SSOTSyncHandler struct {
	cfg config.Config
	log *zap.Logger
	pg  *store.Postgres
}

func NewSSOTSyncHandler(cfg config.Config, log *zap.Logger, pg *store.Postgres) *SSOTSyncHandler {
	return &SSOTSyncHandler{cfg: cfg, log: log, pg: pg}
}

func (h *SSOTSyncHandler) SyncSchools(w http.ResponseWriter, r *http.Request) {
	h.sync(w, r, store.SSOTSchools)
}

func (h *SSOTSyncHandler) SyncDevices(w http.ResponseWriter, r *http.Request) {
	h.sync(w, r, store.SSOTDevices)
}

func (h *SSOTSyncHandler) SyncParts(w http.ResponseWriter, r *http.Request) {
	h.sync(w, r, store.SSOTParts)
}

func (h *SSOTSyncHandler) sync(w http.ResponseWriter, r *http.Request, res store.SSOTResource) {
	tenant := middleware.TenantID(r.Context())
	limit := h.cfg.SSOTSyncPageSize

	// Load checkpoint
	state, err := h.pg.SSOTState().Get(r.Context(), tenant, res)
	if err != nil {
		state = store.NewSSOTSyncState(tenant, res)
		_ = h.pg.SSOTState().Upsert(r.Context(), state)
	}

	cursor := state.LastCursor
	since := state.LastUpdatedSince
	total := 0
	maxSeen := since

	var client *ssot.Client
	switch res {
	case store.SSOTSchools:
		client = ssot.NewClient(h.cfg.SchoolSSOTBaseURL).WithTenant(tenant)
	case store.SSOTDevices:
		client = ssot.NewClient(h.cfg.DeviceSSOTBaseURL).WithTenant(tenant)
	case store.SSOTParts:
		client = ssot.NewClient(h.cfg.PartsSSOTBaseURL).WithTenant(tenant)
	default:
		http.Error(w, "unknown resource", http.StatusBadRequest)
		return
	}

	for {
		if client.BaseURL == "" {
			http.Error(w, "ssot base url not configured", http.StatusBadRequest)
			return
		}

		switch res {
		case store.SSOTSchools:
			page, err := client.ListSchools(since, cursor, limit)
			if err != nil { http.Error(w, "school ssot sync failed: "+err.Error(), http.StatusBadGateway); return }
			for _, it := range page.Items {
				it.TenantID = tenant
				_ = h.pg.SchoolsSnapshot().Upsert(r.Context(), it)
				if it.UpdatedAt.After(maxSeen) { maxSeen = it.UpdatedAt }
				total++
			}
			if page.Next == "" {
				cursor = ""
				break
			}
			cursor = page.Next
			continue

		case store.SSOTDevices:
			page, err := client.ListDevices(since, cursor, limit)
			if err != nil { http.Error(w, "device ssot sync failed: "+err.Error(), http.StatusBadGateway); return }
			for _, it := range page.Items {
				it.TenantID = tenant
				_ = h.pg.DevicesSnapshot().Upsert(r.Context(), it)
				if it.UpdatedAt.After(maxSeen) { maxSeen = it.UpdatedAt }
				total++
			}
			if page.Next == "" {
				cursor = ""
				break
			}
			cursor = page.Next
			continue

		case store.SSOTParts:
			page, err := client.ListParts(since, cursor, limit)
			if err != nil { http.Error(w, "parts ssot sync failed: "+err.Error(), http.StatusBadGateway); return }
			for _, it := range page.Items {
				it.TenantID = tenant
				_ = h.pg.PartsSnapshot().Upsert(r.Context(), it)
				if it.UpdatedAt.After(maxSeen) { maxSeen = it.UpdatedAt }
				total++
			}
			if page.Next == "" {
				cursor = ""
				break
			}
			cursor = page.Next
			continue
		}

		break
	}

	// Advance checkpoint: set updatedSince to maxSeen, clear cursor
	state.LastUpdatedSince = maxSeen
	state.LastCursor = cursor
	state.UpdatedAt = time.Now().UTC()
	_ = h.pg.SSOTState().Upsert(r.Context(), state)

	writeJSON(w, http.StatusOK, map[string]any{
		"ok": true,
		"resource": string(res),
		"synced": total,
		"newUpdatedSince": state.LastUpdatedSince,
	})
}
