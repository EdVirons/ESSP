package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/edvirons/ssp/ims/internal/audit"
	"github.com/edvirons/ssp/ims/internal/config"
	"github.com/edvirons/ssp/ims/internal/lookups"
	"github.com/edvirons/ssp/ims/internal/middleware"
	"github.com/edvirons/ssp/ims/internal/models"
	"github.com/edvirons/ssp/ims/internal/service"
	"github.com/edvirons/ssp/ims/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type IncidentHandler struct {
	cfg   config.Config
	log   *zap.Logger
	pg    *store.Postgres
	rdb   *redis.Client
	audit audit.AuditLogger
}

func NewIncidentHandler(cfg config.Config, log *zap.Logger, pg *store.Postgres, rdb *redis.Client, auditLogger audit.AuditLogger) *IncidentHandler {
	return &IncidentHandler{cfg: cfg, log: log, pg: pg, rdb: rdb, audit: auditLogger}
}

type createIncidentReq struct {
	DeviceID    string          `json:"deviceId"`
	Category    string          `json:"category"`
	Severity    models.Severity `json:"severity"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	ReportedBy  string          `json:"reportedBy"`
}

func (h *IncidentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createIncidentReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(req.DeviceID) == "" || strings.TrimSpace(req.Title) == "" {
		http.Error(w, "deviceId and title are required", http.StatusBadRequest)
		return
	}

	tenant := middleware.TenantID(r.Context())
	school := middleware.SchoolID(r.Context())

	now := time.Now().UTC()

	// Enrich from SSOT snapshots (best-effort)
	lk := lookups.New(h.pg.RawPool())
	sc, _ := lk.SchoolByID(r.Context(), tenant, school)
	pc, _ := lk.PrimaryContactBySchoolID(r.Context(), tenant, school)
	dv, _ := lk.DeviceByID(r.Context(), tenant, strings.TrimSpace(req.DeviceID))

	inc := models.Incident{
		ID:       store.NewID("inc"),
		TenantID: tenant,
		SchoolID: school,
		DeviceID: strings.TrimSpace(req.DeviceID),
		SchoolName: func() string {
			if sc != nil {
				return sc.Name
			}
			return ""
		}(),
		CountyID: func() string {
			if sc != nil {
				return sc.CountyID
			}
			return ""
		}(),
		CountyName: func() string {
			if sc != nil {
				return sc.CountyName
			}
			return ""
		}(),
		SubCountyID: func() string {
			if sc != nil {
				return sc.SubCountyID
			}
			return ""
		}(),
		SubCountyName: func() string {
			if sc != nil {
				return sc.SubCountyName
			}
			return ""
		}(),
		ContactName: func() string {
			if pc != nil {
				return pc.Name
			}
			return ""
		}(),
		ContactPhone: func() string {
			if pc != nil {
				return pc.Phone
			}
			return ""
		}(),
		ContactEmail: func() string {
			if pc != nil {
				return pc.Email
			}
			return ""
		}(),
		DeviceSerial: func() string {
			if dv != nil {
				return dv.Serial
			}
			return ""
		}(),
		DeviceAssetTag: func() string {
			if dv != nil {
				return dv.AssetTag
			}
			return ""
		}(),
		DeviceModelID: func() string {
			if dv != nil {
				return dv.ModelID
			}
			return ""
		}(),
		DeviceMake: func() string {
			if dv != nil {
				return dv.Make
			}
			return ""
		}(),
		DeviceModel: func() string {
			if dv != nil {
				return dv.Model
			}
			return ""
		}(),
		DeviceCategory: func() string {
			if dv != nil {
				return dv.Category
			}
			return ""
		}(),
		Category:    strings.TrimSpace(req.Category),
		Severity:    req.Severity,
		Status:      models.IncidentNew,
		Title:       strings.TrimSpace(req.Title),
		Description: strings.TrimSpace(req.Description),
		ReportedBy:  strings.TrimSpace(req.ReportedBy),
		SLADueAt:    service.SLADue(req.Severity, now),
		SLABreached: false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := h.pg.Incidents().Create(r.Context(), inc); err != nil {
		http.Error(w, "failed to create incident", http.StatusInternalServerError)
		return
	}

	// Log the creation in audit trail
	if err := h.audit.LogCreate(r.Context(), "incident", inc.ID, inc); err != nil {
		h.log.Error("failed to log incident creation audit", zap.Error(err))
		// Don't fail the request if audit logging fails
	}

	// Optional: auto-route to service shop (sub-county first, fallback to county) -> create work order
	if h.cfg.AutoRouteWorkOrders {
		// Resolve school's geography from SSOT snapshot cache
		sp, err := h.pg.SchoolsSnapshot().Get(r.Context(), tenant, school)
		if err == nil {
			var shop models.ServiceShop
			// Try sub-county coverage first
			if strings.TrimSpace(sp.SubCountyCode) != "" && strings.TrimSpace(sp.CountyCode) != "" {
				shop, err = h.pg.ServiceShops().GetBySubCounty(r.Context(), tenant, strings.TrimSpace(sp.CountyCode), strings.TrimSpace(sp.SubCountyCode))
			}
			// Fallback to county coverage
			if err != nil || shop.ID == "" {
				if strings.TrimSpace(sp.CountyCode) != "" {
					shop, err = h.pg.ServiceShops().GetByCounty(r.Context(), tenant, strings.TrimSpace(sp.CountyCode))
				}
			}
			if err == nil && shop.ID != "" {
				// Assign to lead technician if available
				lead, _ := h.pg.ServiceStaff().GetLeadByShop(r.Context(), tenant, shop.ID)
				rl := models.RepairLocation(h.cfg.DefaultRepairLocation)
				if rl == "" {
					rl = models.RepairLocationServiceShop
				}
				wo := models.WorkOrder{
					ID:                store.NewID("wo"),
					IncidentID:        inc.ID,
					TenantID:          tenant,
					SchoolID:          school,
					DeviceID:          inc.DeviceID,
					Status:            models.WorkOrderDraft,
					ServiceShopID:     shop.ID,
					AssignedStaffID:   lead.ID,
					RepairLocation:    rl,
					AssignedTo:        lead.UserID,
					TaskType:          "triage",
					CostEstimateCents: 0,
					Notes:             "Auto-created from incident " + inc.ID,
					CreatedAt:         now,
					UpdatedAt:         now,
				}
				_ = h.pg.WorkOrders().Create(r.Context(), wo)
			}
		}
	}
	writeJSON(w, http.StatusOK, inc)
}

func (h *IncidentHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	tenant := middleware.TenantID(r.Context())
	school := middleware.SchoolID(r.Context())

	inc, err := h.pg.Incidents().GetByID(r.Context(), tenant, school, id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, inc)
}

func (h *IncidentHandler) List(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())
	school := middleware.SchoolID(r.Context())

	status := strings.TrimSpace(r.URL.Query().Get("status"))
	deviceID := strings.TrimSpace(r.URL.Query().Get("deviceId"))
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	limit := parseLimit(r.URL.Query().Get("limit"), 50, 200)

	curT, curID, hasCur := decodeCursor(strings.TrimSpace(r.URL.Query().Get("cursor")))

	items, next, err := h.pg.Incidents().List(r.Context(), store.IncidentListParams{
		TenantID:        tenant,
		SchoolID:        school,
		Status:          status,
		DeviceID:        deviceID,
		Query:           q,
		Limit:           limit,
		CursorCreatedAt: curT,
		CursorID:        curID,
		HasCursor:       hasCur,
	})
	if err != nil {
		http.Error(w, "failed to list incidents", http.StatusInternalServerError)
		return
	}

	resp := map[string]any{
		"items":      items,
		"nextCursor": next,
	}
	writeJSON(w, http.StatusOK, resp)
}

type updateIncidentStatusReq struct {
	Status models.IncidentStatus `json:"status"`
}

func (h *IncidentHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req updateIncidentStatusReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	tenant := middleware.TenantID(r.Context())
	school := middleware.SchoolID(r.Context())

	cur, err := h.pg.Incidents().GetByID(r.Context(), tenant, school, id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if !service.CanTransitionIncident(cur.Status, req.Status) {
		http.Error(w, "invalid status transition", http.StatusBadRequest)
		return
	}

	updated, err := h.pg.Incidents().UpdateStatus(r.Context(), tenant, school, id, req.Status, time.Now().UTC())
	if err != nil {
		http.Error(w, "failed to update status", http.StatusInternalServerError)
		return
	}

	// Log the update in audit trail
	if err := h.audit.LogUpdate(r.Context(), "incident", id, cur, updated); err != nil {
		h.log.Error("failed to log incident update audit", zap.Error(err))
		// Don't fail the request if audit logging fails
	}

	writeJSON(w, http.StatusOK, updated)
}
