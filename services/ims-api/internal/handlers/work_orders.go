package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/edvirons/ssp/ims/internal/audit"
	"github.com/edvirons/ssp/ims/internal/middleware"
	"github.com/edvirons/ssp/ims/internal/models"
	"github.com/edvirons/ssp/ims/internal/service"
	"github.com/edvirons/ssp/ims/internal/store"
	"github.com/edvirons/ssp/ims/internal/lookups"
	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type WorkOrderHandler struct {
	log   *zap.Logger
	pg    *store.Postgres
	rdb   *redis.Client
	audit audit.AuditLogger
}

func NewWorkOrderHandler(log *zap.Logger, pg *store.Postgres, rdb *redis.Client, auditLogger audit.AuditLogger) *WorkOrderHandler {
	return &WorkOrderHandler{log: log, pg: pg, rdb: rdb, audit: auditLogger}
}

type createWOReq struct {
	IncidentID        string `json:"incidentId"`
	DeviceID          string `json:"deviceId"`
	TaskType          string `json:"taskType"`
	ServiceShopID     string `json:"serviceShopId"`
	AssignedStaffID   string `json:"assignedStaffId"`
	RepairLocation    models.RepairLocation `json:"repairLocation"`
	AssignedTo        string `json:"assignedTo"`
	CostEstimateCents int64  `json:"costEstimateCents"`
	Notes             string `json:"notes"`
}

func (h *WorkOrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createWOReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(req.DeviceID) == "" || strings.TrimSpace(req.TaskType) == "" {
		http.Error(w, "deviceId and taskType are required", http.StatusBadRequest)
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

	wo := models.WorkOrder{
		ID:                store.NewID("wo"),
		IncidentID:        strings.TrimSpace(req.IncidentID),
		TenantID:          tenant,
		SchoolID:          school,
		DeviceID:          strings.TrimSpace(req.DeviceID),
		SchoolName:        func() string { if sc!=nil { return sc.Name }; return "" }(),
		ContactName:       func() string { if pc!=nil { return pc.Name }; return "" }(),
		ContactPhone:      func() string { if pc!=nil { return pc.Phone }; return "" }(),
		ContactEmail:      func() string { if pc!=nil { return pc.Email }; return "" }(),
		DeviceSerial:      func() string { if dv!=nil { return dv.Serial }; return "" }(),
		DeviceAssetTag:    func() string { if dv!=nil { return dv.AssetTag }; return "" }(),
		DeviceModelID:     func() string { if dv!=nil { return dv.ModelID }; return "" }(),
		DeviceMake:        func() string { if dv!=nil { return dv.Make }; return "" }(),
		DeviceModel:       func() string { if dv!=nil { return dv.Model }; return "" }(),
		DeviceCategory:    func() string { if dv!=nil { return dv.Category }; return "" }(),
		Status:            models.WorkOrderDraft,
		ServiceShopID:     strings.TrimSpace(req.ServiceShopID),
		AssignedStaffID:   strings.TrimSpace(req.AssignedStaffID),
		RepairLocation:    req.RepairLocation,
		AssignedTo:        strings.TrimSpace(req.AssignedTo),
		TaskType:          strings.TrimSpace(req.TaskType),
		CostEstimateCents: req.CostEstimateCents,
		Notes:             strings.TrimSpace(req.Notes),
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	if err := h.pg.WorkOrders().Create(r.Context(), wo); err != nil {
		http.Error(w, "failed to create work order", http.StatusInternalServerError)
		return
	}

	// Log the creation in audit trail
	if err := h.audit.LogCreate(r.Context(), "work_order", wo.ID, wo); err != nil {
		h.log.Error("failed to log work order creation audit", zap.Error(err))
		// Don't fail the request if audit logging fails
	}

	writeJSON(w, http.StatusCreated, wo)
}

func (h *WorkOrderHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	tenant := middleware.TenantID(r.Context())
	school := middleware.SchoolID(r.Context())

	wo, err := h.pg.WorkOrders().GetByID(r.Context(), tenant, school, id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, wo)
}

func (h *WorkOrderHandler) List(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())
	school := middleware.SchoolID(r.Context())

	status := strings.TrimSpace(r.URL.Query().Get("status"))
	deviceID := strings.TrimSpace(r.URL.Query().Get("deviceId"))
	incidentID := strings.TrimSpace(r.URL.Query().Get("incidentId"))
	limit := parseLimit(r.URL.Query().Get("limit"), 50, 200)

	curT, curID, hasCur := decodeCursor(strings.TrimSpace(r.URL.Query().Get("cursor")))

	items, next, err := h.pg.WorkOrders().List(r.Context(), store.WorkOrderListParams{
		TenantID: tenant,
		SchoolID: school,
		Status: status,
		DeviceID: deviceID,
		IncidentID: incidentID,
		Limit: limit,
		CursorCreatedAt: curT,
		CursorID: curID,
		HasCursor: hasCur,
	})
	if err != nil {
		http.Error(w, "failed to list work orders", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items, "nextCursor": next})
}

type updateWOStatusReq struct {
	Status models.WorkOrderStatus `json:"status"`
}

func (h *WorkOrderHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req updateWOStatusReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	tenant := middleware.TenantID(r.Context())
	school := middleware.SchoolID(r.Context())

	cur, err := h.pg.WorkOrders().GetByID(r.Context(), tenant, school, id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if !service.CanTransitionWorkOrder(cur.Status, req.Status) {
		http.Error(w, "invalid status transition", http.StatusBadRequest)
		return
	}

	updated, err := h.pg.WorkOrders().UpdateStatus(r.Context(), tenant, school, id, req.Status, time.Now().UTC())
	if err != nil {
		http.Error(w, "failed to update status", http.StatusInternalServerError)
		return
	}

	// Log the update in audit trail
	if err := h.audit.LogUpdate(r.Context(), "work_order", id, cur, updated); err != nil {
		h.log.Error("failed to log work order update audit", zap.Error(err))
		// Don't fail the request if audit logging fails
	}

	writeJSON(w, http.StatusOK, updated)
}
