package handlers

import (
	"fmt"
	"net/http"

	"github.com/edvirons/ssp/ims/internal/config"
	"github.com/edvirons/ssp/ims/internal/middleware"
	"github.com/edvirons/ssp/ims/internal/store"
	"go.uber.org/zap"
)

type SSOTListHandler struct {
	cfg config.Config
	log *zap.Logger
	pg  *store.Postgres
}

func NewSSOTListHandler(cfg config.Config, log *zap.Logger, pg *store.Postgres) *SSOTListHandler {
	return &SSOTListHandler{cfg: cfg, log: log, pg: pg}
}

func (h *SSOTListHandler) ListSchools(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())
	q := r.URL.Query()

	params := store.SchoolSnapshotListParams{
		TenantID:   tenant,
		Query:      q.Get("q"),
		CountyCode: q.Get("countyCode"),
		Level:      q.Get("level"),
		Type:       q.Get("type"),
		Limit:      parseLimit(q.Get("limit"), 50, 200),
		Offset:     parseOffset(q.Get("offset")),
	}

	items, total, err := h.pg.SchoolsSnapshot().List(r.Context(), params)
	if err != nil {
		h.log.Error("failed to list schools", zap.Error(err))
		http.Error(w, "failed to list schools", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"items": items,
		"total": total,
		"limit": params.Limit,
		"offset": params.Offset,
	})
}

func (h *SSOTListHandler) ListDevices(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())
	q := r.URL.Query()

	params := store.DeviceSnapshotListParams{
		TenantID: tenant,
		SchoolID: q.Get("schoolId"),
		Query:    q.Get("q"),
		Status:   q.Get("status"),
		Limit:    parseLimit(q.Get("limit"), 50, 200),
		Offset:   parseOffset(q.Get("offset")),
	}

	items, total, err := h.pg.DevicesSnapshot().List(r.Context(), params)
	if err != nil {
		h.log.Error("failed to list devices", zap.Error(err))
		http.Error(w, "failed to list devices", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"items": items,
		"total": total,
		"limit": params.Limit,
		"offset": params.Offset,
	})
}

func (h *SSOTListHandler) GetDeviceStats(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())

	stats, err := h.pg.DevicesSnapshot().Stats(r.Context(), tenant)
	if err != nil {
		h.log.Error("failed to get device stats", zap.Error(err))
		http.Error(w, "failed to get device stats", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, stats)
}

func (h *SSOTListHandler) ListDeviceModels(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())

	models, err := h.pg.DevicesSnapshot().ListModels(r.Context(), tenant)
	if err != nil {
		h.log.Error("failed to list device models", zap.Error(err))
		http.Error(w, "failed to list device models", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"items": models,
		"total": len(models),
	})
}

func (h *SSOTListHandler) GetDeviceMakes(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())

	makes, err := h.pg.DevicesSnapshot().ListMakes(r.Context(), tenant)
	if err != nil {
		h.log.Error("failed to get device makes", zap.Error(err))
		http.Error(w, "failed to get device makes", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"makes": makes,
	})
}

func (h *SSOTListHandler) ListParts(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())
	q := r.URL.Query()

	params := store.PartSnapshotListParams{
		TenantID: tenant,
		Category: q.Get("category"),
		Query:    q.Get("q"),
		Limit:    parseLimit(q.Get("limit"), 50, 200),
		Offset:   parseOffset(q.Get("offset")),
	}

	items, total, err := h.pg.PartsSnapshot().List(r.Context(), params)
	if err != nil {
		h.log.Error("failed to list parts", zap.Error(err))
		http.Error(w, "failed to list parts", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"items": items,
		"total": total,
		"limit": params.Limit,
		"offset": params.Offset,
	})
}

func (h *SSOTListHandler) GetSyncStatus(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())

	// Get counts from each snapshot table
	schoolCount, _ := h.pg.SchoolsSnapshot().Count(r.Context(), tenant)
	deviceCount, _ := h.pg.DevicesSnapshot().Count(r.Context(), tenant)
	partCount, _ := h.pg.PartsSnapshot().Count(r.Context(), tenant)

	// Get sync states
	schoolState, _ := h.pg.SSOTState().Get(r.Context(), tenant, store.SSOTSchools)
	deviceState, _ := h.pg.SSOTState().Get(r.Context(), tenant, store.SSOTDevices)
	partState, _ := h.pg.SSOTState().Get(r.Context(), tenant, store.SSOTParts)

	writeJSON(w, http.StatusOK, map[string]any{
		"schools": map[string]any{
			"count":       schoolCount,
			"lastSyncAt":  schoolState.UpdatedAt,
			"lastCursor":  schoolState.LastCursor,
		},
		"devices": map[string]any{
			"count":       deviceCount,
			"lastSyncAt":  deviceState.UpdatedAt,
			"lastCursor":  deviceState.LastCursor,
		},
		"parts": map[string]any{
			"count":       partCount,
			"lastSyncAt":  partState.UpdatedAt,
			"lastCursor":  partState.LastCursor,
		},
	})
}

func (h *SSOTListHandler) ListCounties(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())

	counties, err := h.pg.SchoolsSnapshot().ListCounties(r.Context(), tenant)
	if err != nil {
		h.log.Error("failed to list counties", zap.Error(err))
		http.Error(w, "failed to list counties", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"items": counties,
		"total": len(counties),
	})
}

func (h *SSOTListHandler) ListSubCounties(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())
	countyCode := r.URL.Query().Get("countyCode")

	subCounties, err := h.pg.SchoolsSnapshot().ListSubCounties(r.Context(), tenant, countyCode)
	if err != nil {
		h.log.Error("failed to list sub-counties", zap.Error(err))
		http.Error(w, "failed to list sub-counties", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"items": subCounties,
		"total": len(subCounties),
	})
}

func parseOffset(s string) int {
	if s == "" {
		return 0
	}
	var v int
	if _, err := fmt.Sscanf(s, "%d", &v); err != nil {
		return 0
	}
	if v < 0 {
		return 0
	}
	return v
}
