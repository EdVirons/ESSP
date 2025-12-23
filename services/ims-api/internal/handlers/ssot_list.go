package handlers

import (
	"fmt"
	"net/http"

	"github.com/edvirons/ssp/ims/internal/config"
	"github.com/edvirons/ssp/ims/internal/middleware"
	"github.com/edvirons/ssp/ims/internal/store"
	"github.com/go-chi/chi/v5"
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
		"items":  items,
		"total":  total,
		"limit":  params.Limit,
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
		"items":  items,
		"total":  total,
		"limit":  params.Limit,
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
		"items":  items,
		"total":  total,
		"limit":  params.Limit,
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
			"count":      schoolCount,
			"lastSyncAt": schoolState.UpdatedAt,
			"lastCursor": schoolState.LastCursor,
		},
		"devices": map[string]any{
			"count":      deviceCount,
			"lastSyncAt": deviceState.UpdatedAt,
			"lastCursor": deviceState.LastCursor,
		},
		"parts": map[string]any{
			"count":      partCount,
			"lastSyncAt": partState.UpdatedAt,
			"lastCursor": partState.LastCursor,
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

// ==================== HR SSOT Handlers ====================

func (h *SSOTListHandler) ListPeople(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())
	q := r.URL.Query()

	params := store.PersonSnapshotListParams{
		TenantID:  tenant,
		Query:     q.Get("q"),
		Status:    q.Get("status"),
		OrgUnitID: q.Get("orgUnitId"),
		Limit:     parseLimit(q.Get("limit"), 50, 200),
		Offset:    parseOffset(q.Get("offset")),
	}

	items, total, err := h.pg.PeopleSnapshot().List(r.Context(), params)
	if err != nil {
		h.log.Error("failed to list people", zap.Error(err))
		http.Error(w, "failed to list people", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"items":  items,
		"total":  total,
		"limit":  params.Limit,
		"offset": params.Offset,
	})
}

func (h *SSOTListHandler) GetPerson(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())
	personID := chi.URLParam(r, "personId")

	person, err := h.pg.PeopleSnapshot().Get(r.Context(), tenant, personID)
	if err != nil {
		h.log.Error("failed to get person", zap.Error(err), zap.String("personId", personID))
		http.Error(w, "person not found", http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, person)
}

func (h *SSOTListHandler) ListTeams(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())
	q := r.URL.Query()

	params := store.TeamSnapshotListParams{
		TenantID:  tenant,
		Query:     q.Get("q"),
		Key:       q.Get("key"),
		OrgUnitID: q.Get("orgUnitId"),
		Limit:     parseLimit(q.Get("limit"), 50, 200),
		Offset:    parseOffset(q.Get("offset")),
	}

	items, total, err := h.pg.TeamsSnapshot().List(r.Context(), params)
	if err != nil {
		h.log.Error("failed to list teams", zap.Error(err))
		http.Error(w, "failed to list teams", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"items":  items,
		"total":  total,
		"limit":  params.Limit,
		"offset": params.Offset,
	})
}

func (h *SSOTListHandler) GetTeam(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())
	teamID := chi.URLParam(r, "teamId")

	team, err := h.pg.TeamsSnapshot().Get(r.Context(), tenant, teamID)
	if err != nil {
		h.log.Error("failed to get team", zap.Error(err), zap.String("teamId", teamID))
		http.Error(w, "team not found", http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, team)
}

func (h *SSOTListHandler) ListOrgUnits(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())
	q := r.URL.Query()

	params := store.OrgUnitSnapshotListParams{
		TenantID: tenant,
		Query:    q.Get("q"),
		Kind:     q.Get("kind"),
		ParentID: q.Get("parentId"),
		Limit:    parseLimit(q.Get("limit"), 50, 200),
		Offset:   parseOffset(q.Get("offset")),
	}

	items, total, err := h.pg.OrgUnitsSnapshot().List(r.Context(), params)
	if err != nil {
		h.log.Error("failed to list org units", zap.Error(err))
		http.Error(w, "failed to list org units", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"items":  items,
		"total":  total,
		"limit":  params.Limit,
		"offset": params.Offset,
	})
}

func (h *SSOTListHandler) GetOrgUnit(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())
	orgUnitID := chi.URLParam(r, "orgUnitId")

	unit, err := h.pg.OrgUnitsSnapshot().Get(r.Context(), tenant, orgUnitID)
	if err != nil {
		h.log.Error("failed to get org unit", zap.Error(err), zap.String("orgUnitId", orgUnitID))
		http.Error(w, "org unit not found", http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, unit)
}

func (h *SSOTListHandler) GetOrgTree(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())

	tree, err := h.pg.OrgUnitsSnapshot().GetTree(r.Context(), tenant)
	if err != nil {
		h.log.Error("failed to get org tree", zap.Error(err))
		http.Error(w, "failed to get org tree", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, tree)
}

// ==================== Team Memberships Handlers ====================

func (h *SSOTListHandler) ListTeamMemberships(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())
	q := r.URL.Query()

	params := store.TeamMembershipSnapshotListParams{
		TenantID: tenant,
		TeamID:   q.Get("teamId"),
		PersonID: q.Get("personId"),
		Role:     q.Get("role"),
		Status:   q.Get("status"),
		Limit:    parseLimit(q.Get("limit"), 50, 200),
		Offset:   parseOffset(q.Get("offset")),
	}

	items, total, err := h.pg.TeamMembershipsSnapshot().List(r.Context(), params)
	if err != nil {
		h.log.Error("failed to list team memberships", zap.Error(err))
		http.Error(w, "failed to list team memberships", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"items":  items,
		"total":  total,
		"limit":  params.Limit,
		"offset": params.Offset,
	})
}

func (h *SSOTListHandler) GetTeamMembership(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())
	membershipID := chi.URLParam(r, "membershipId")

	membership, err := h.pg.TeamMembershipsSnapshot().Get(r.Context(), tenant, membershipID)
	if err != nil {
		h.log.Error("failed to get team membership", zap.Error(err), zap.String("membershipId", membershipID))
		http.Error(w, "team membership not found", http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, membership)
}
