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

// HR SSOT Sync handlers

func (h *SSOTSyncHandler) SyncPeople(w http.ResponseWriter, r *http.Request) {
	h.syncHR(w, r, store.SSOTPeople)
}

func (h *SSOTSyncHandler) SyncTeams(w http.ResponseWriter, r *http.Request) {
	h.syncHR(w, r, store.SSOTTeams)
}

func (h *SSOTSyncHandler) SyncOrgUnits(w http.ResponseWriter, r *http.Request) {
	h.syncHR(w, r, store.SSOTOrgUnits)
}

func (h *SSOTSyncHandler) SyncTeamMemberships(w http.ResponseWriter, r *http.Request) {
	h.syncHR(w, r, store.SSOTTeamMemberships)
}

// syncHR fetches the HR export and upserts the relevant data into snapshot tables
func (h *SSOTSyncHandler) syncHR(w http.ResponseWriter, r *http.Request, res store.SSOTResource) {
	tenant := middleware.TenantID(r.Context())

	if h.cfg.HRSSOTBaseURL == "" {
		http.Error(w, "hr ssot base url not configured", http.StatusBadRequest)
		return
	}

	client := ssot.NewClient(h.cfg.HRSSOTBaseURL).WithTenant(tenant)
	export, err := client.FetchHRExport()
	if err != nil {
		h.log.Error("hr ssot export failed", zap.Error(err))
		http.Error(w, "hr ssot sync failed: "+err.Error(), http.StatusBadGateway)
		return
	}

	var total int
	var maxSeen time.Time

	switch res {
	case store.SSOTPeople:
		for _, p := range export.People {
			fullName := p.GivenName
			if p.FamilyName != "" {
				if fullName != "" {
					fullName += " "
				}
				fullName += p.FamilyName
			}
			snap := store.PersonSnapshot{
				TenantID:   tenant,
				PersonID:   p.ID,
				OrgUnitID:  p.OrgUnitID,
				Status:     p.Status,
				GivenName:  p.GivenName,
				FamilyName: p.FamilyName,
				FullName:   fullName,
				Email:      p.Email,
				Phone:      p.Phone,
				Title:      p.Title,
				UpdatedAt:  p.UpdatedAt,
			}
			if err := h.pg.PeopleSnapshot().Upsert(r.Context(), snap); err != nil {
				h.log.Error("failed to upsert person", zap.Error(err), zap.String("personId", p.ID))
			}
			if p.UpdatedAt.After(maxSeen) {
				maxSeen = p.UpdatedAt
			}
			total++
		}

	case store.SSOTTeams:
		for _, t := range export.Teams {
			snap := store.TeamSnapshot{
				TenantID:    tenant,
				TeamID:      t.ID,
				OrgUnitID:   t.OrgUnitID,
				Key:         t.Key,
				Name:        t.Name,
				Description: t.Description,
				UpdatedAt:   t.UpdatedAt,
			}
			if err := h.pg.TeamsSnapshot().Upsert(r.Context(), snap); err != nil {
				h.log.Error("failed to upsert team", zap.Error(err), zap.String("teamId", t.ID))
			}
			if t.UpdatedAt.After(maxSeen) {
				maxSeen = t.UpdatedAt
			}
			total++
		}

	case store.SSOTOrgUnits:
		for _, o := range export.OrgUnits {
			snap := store.OrgUnitSnapshot{
				TenantID:  tenant,
				OrgUnitID: o.ID,
				ParentID:  o.ParentID,
				Code:      o.Code,
				Name:      o.Name,
				Kind:      o.Kind,
				UpdatedAt: o.UpdatedAt,
			}
			if err := h.pg.OrgUnitsSnapshot().Upsert(r.Context(), snap); err != nil {
				h.log.Error("failed to upsert org unit", zap.Error(err), zap.String("orgUnitId", o.ID))
			}
			if o.UpdatedAt.After(maxSeen) {
				maxSeen = o.UpdatedAt
			}
			total++
		}

	case store.SSOTTeamMemberships:
		for _, m := range export.TeamMemberships {
			snap := store.TeamMembershipSnapshot{
				TenantID:     tenant,
				MembershipID: m.ID,
				TeamID:       m.TeamID,
				PersonID:     m.PersonID,
				Role:         m.Role,
				Status:       m.Status,
				StartedAt:    m.StartedAt,
				EndedAt:      m.EndedAt,
				UpdatedAt:    m.UpdatedAt,
			}
			if err := h.pg.TeamMembershipsSnapshot().Upsert(r.Context(), snap); err != nil {
				h.log.Error("failed to upsert team membership", zap.Error(err), zap.String("membershipId", m.ID))
			}
			if m.UpdatedAt.After(maxSeen) {
				maxSeen = m.UpdatedAt
			}
			total++
		}

	default:
		http.Error(w, "unknown hr resource", http.StatusBadRequest)
		return
	}

	// Update sync state
	state := store.NewSSOTSyncState(tenant, res)
	state.LastUpdatedSince = maxSeen
	state.UpdatedAt = time.Now().UTC()
	_ = h.pg.SSOTState().Upsert(r.Context(), state)

	writeJSON(w, http.StatusOK, map[string]any{
		"ok":              true,
		"resource":        string(res),
		"synced":          total,
		"newUpdatedSince": maxSeen,
	})
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
			if err != nil {
				http.Error(w, "school ssot sync failed: "+err.Error(), http.StatusBadGateway)
				return
			}
			for _, it := range page.Items {
				it.TenantID = tenant
				_ = h.pg.SchoolsSnapshot().Upsert(r.Context(), it)
				if it.UpdatedAt.After(maxSeen) {
					maxSeen = it.UpdatedAt
				}
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
			if err != nil {
				http.Error(w, "device ssot sync failed: "+err.Error(), http.StatusBadGateway)
				return
			}
			for _, it := range page.Items {
				it.TenantID = tenant
				_ = h.pg.DevicesSnapshot().Upsert(r.Context(), it)
				if it.UpdatedAt.After(maxSeen) {
					maxSeen = it.UpdatedAt
				}
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
			if err != nil {
				http.Error(w, "parts ssot sync failed: "+err.Error(), http.StatusBadGateway)
				return
			}
			for _, it := range page.Items {
				it.TenantID = tenant
				_ = h.pg.PartsSnapshot().Upsert(r.Context(), it)
				if it.UpdatedAt.After(maxSeen) {
					maxSeen = it.UpdatedAt
				}
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
		"ok":              true,
		"resource":        string(res),
		"synced":          total,
		"newUpdatedSince": state.LastUpdatedSince,
	})
}
