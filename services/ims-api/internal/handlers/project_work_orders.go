package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/edvirons/ssp/ims/internal/lookups"
	"github.com/edvirons/ssp/ims/internal/middleware"
	"github.com/edvirons/ssp/ims/internal/models"
	"github.com/edvirons/ssp/ims/internal/store"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type ProjectWorkOrdersHandler struct {
	log        *zap.Logger
	pg         *store.Postgres
	activities *ProjectActivitiesHandler
}

func NewProjectWorkOrdersHandler(log *zap.Logger, pg *store.Postgres, activities *ProjectActivitiesHandler) *ProjectWorkOrdersHandler {
	return &ProjectWorkOrdersHandler{log: log, pg: pg, activities: activities}
}

type createProjectWOReq struct {
	PhaseID           string                `json:"phaseId"`
	DeviceID          string                `json:"deviceId"`
	TaskType          string                `json:"taskType"`
	ServiceShopID     string                `json:"serviceShopId"`
	AssignedStaffID   string                `json:"assignedStaffId"`
	RepairLocation    models.RepairLocation `json:"repairLocation"`
	AssignedTo        string                `json:"assignedTo"`
	CostEstimateCents int64                 `json:"costEstimateCents"`
	Notes             string                `json:"notes"`
}

func (h *ProjectWorkOrdersHandler) CreateFromProject(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "id")
	tenant := middleware.TenantID(r.Context())
	school := middleware.SchoolID(r.Context())
	userID := middleware.UserID(r.Context())
	userName := middleware.UserName(r.Context())

	var req createProjectWOReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	// Validate required fields
	taskType := strings.TrimSpace(req.TaskType)
	if taskType == "" {
		http.Error(w, "taskType required", http.StatusBadRequest)
		return
	}

	// Verify project exists
	project, err := h.pg.Projects().GetByID(r.Context(), tenant, projectID)
	if err != nil {
		http.Error(w, "project not found", http.StatusNotFound)
		return
	}

	// Use project's school if not specified
	if school == "" {
		school = project.SchoolID
	}

	now := time.Now().UTC()

	// Enrich from SSOT snapshots (best-effort)
	lk := lookups.New(h.pg.RawPool())
	sc, _ := lk.SchoolByID(r.Context(), tenant, school)
	pc, _ := lk.PrimaryContactBySchoolID(r.Context(), tenant, school)

	var dv *lookups.DeviceSummary
	deviceID := strings.TrimSpace(req.DeviceID)
	if deviceID != "" {
		dv, _ = lk.DeviceByID(r.Context(), tenant, deviceID)
	}

	wo := models.WorkOrder{
		ID:       store.NewID("wo"),
		TenantID: tenant,
		SchoolID: school,
		DeviceID: deviceID,

		// Denormalized lookups
		SchoolName:     safeString(sc, func(s *lookups.SchoolSummary) string { return s.Name }),
		ContactName:    safeString(pc, func(c *lookups.ContactSummary) string { return c.Name }),
		ContactPhone:   safeString(pc, func(c *lookups.ContactSummary) string { return c.Phone }),
		ContactEmail:   safeString(pc, func(c *lookups.ContactSummary) string { return c.Email }),
		DeviceSerial:   safeString(dv, func(d *lookups.DeviceSummary) string { return d.Serial }),
		DeviceAssetTag: safeString(dv, func(d *lookups.DeviceSummary) string { return d.AssetTag }),
		DeviceModelID:  safeString(dv, func(d *lookups.DeviceSummary) string { return d.ModelID }),
		DeviceMake:     safeString(dv, func(d *lookups.DeviceSummary) string { return d.Make }),
		DeviceModel:    safeString(dv, func(d *lookups.DeviceSummary) string { return d.Model }),
		DeviceCategory: safeString(dv, func(d *lookups.DeviceSummary) string { return d.Category }),

		Status:            models.WorkOrderDraft,
		ServiceShopID:     strings.TrimSpace(req.ServiceShopID),
		AssignedStaffID:   strings.TrimSpace(req.AssignedStaffID),
		RepairLocation:    req.RepairLocation,
		AssignedTo:        strings.TrimSpace(req.AssignedTo),
		TaskType:          taskType,
		ProjectID:         projectID,
		PhaseID:           strings.TrimSpace(req.PhaseID),
		CostEstimateCents: req.CostEstimateCents,
		Notes:             strings.TrimSpace(req.Notes),

		// Project-originated fields
		CreatedFromProject: true,
		CreatedByUserID:    userID,
		CreatedByUserName:  userName,

		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := h.pg.WorkOrders().Create(r.Context(), wo); err != nil {
		h.log.Error("failed to create work order from project", zap.Error(err))
		http.Error(w, "failed to create work order", http.StatusInternalServerError)
		return
	}

	// Log activity for work order creation
	if h.activities != nil {
		h.activities.LogWorkOrderActivity(r.Context(), tenant, projectID, wo.PhaseID, wo.ID, userID, userName, "created", string(wo.Status))
	}

	writeJSON(w, http.StatusCreated, wo)
}

func (h *ProjectWorkOrdersHandler) ListByProject(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "id")
	tenant := middleware.TenantID(r.Context())

	phaseID := strings.TrimSpace(r.URL.Query().Get("phaseId"))
	status := strings.TrimSpace(r.URL.Query().Get("status"))
	limit := parseLimit(r.URL.Query().Get("limit"), 50, 200)
	curT, curID, hasCur := decodeCursor(strings.TrimSpace(r.URL.Query().Get("cursor")))

	workOrders, next, err := h.pg.WorkOrders().ListByProject(r.Context(), store.ProjectWOListParams{
		TenantID:   tenant,
		ProjectID:  projectID,
		PhaseID:    phaseID,
		Status:     status,
		Limit:      limit,
		HasCursor:  hasCur,
		CursorTime: curT,
		CursorID:   curID,
	})
	if err != nil {
		h.log.Error("failed to list project work orders", zap.Error(err))
		http.Error(w, "failed to list work orders", http.StatusInternalServerError)
		return
	}
	if workOrders == nil {
		workOrders = []models.WorkOrder{}
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"items":      workOrders,
		"nextCursor": next,
	})
}

// Helper for safe nil pointer access
func safeString[T any](ptr *T, fn func(*T) string) string {
	if ptr == nil {
		return ""
	}
	return fn(ptr)
}
