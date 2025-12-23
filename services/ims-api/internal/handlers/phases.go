package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/edvirons/ssp/ims/internal/middleware"
	"github.com/edvirons/ssp/ims/internal/models"
	"github.com/edvirons/ssp/ims/internal/store"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type PhasesHandler struct {
	log *zap.Logger
	pg  *store.Postgres
}

func NewPhasesHandler(log *zap.Logger, pg *store.Postgres) *PhasesHandler {
	return &PhasesHandler{log: log, pg: pg}
}

type createPhaseReq struct {
	PhaseType models.PhaseType `json:"phaseType"`
	OwnerRole string           `json:"ownerRole"`
	StartDate string           `json:"startDate"`
	EndDate   string           `json:"endDate"`
	Notes     string           `json:"notes"`
}

func (h *PhasesHandler) Create(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "id")
	var req createPhaseReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(string(req.PhaseType)) == "" {
		http.Error(w, "phaseType required", http.StatusBadRequest)
		return
	}
	tenant := middleware.TenantID(r.Context())
	now := time.Now().UTC()
	p := models.ServicePhase{
		ID:        store.NewID("phase"),
		TenantID:  tenant,
		ProjectID: projectID,
		PhaseType: req.PhaseType,
		Status:    models.PhasePending,
		OwnerRole: strings.TrimSpace(req.OwnerRole),
		StartDate: strings.TrimSpace(req.StartDate),
		EndDate:   strings.TrimSpace(req.EndDate),
		Notes:     strings.TrimSpace(req.Notes),
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := h.pg.Phases().Create(r.Context(), p); err != nil {
		http.Error(w, "failed to create phase", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, p)
}

func (h *PhasesHandler) List(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "id")
	tenant := middleware.TenantID(r.Context())
	phaseType := strings.TrimSpace(r.URL.Query().Get("phaseType"))
	status := strings.TrimSpace(r.URL.Query().Get("status"))
	limit := parseLimit(r.URL.Query().Get("limit"), 50, 200)
	curT, curID, hasCur := decodeCursor(strings.TrimSpace(r.URL.Query().Get("cursor")))

	items, next, err := h.pg.Phases().List(r.Context(), store.PhaseListParams{
		TenantID: tenant, ProjectID: projectID, PhaseType: phaseType, Status: status,
		Limit: limit, HasCursor: hasCur, CursorCreatedAt: curT, CursorID: curID,
	})
	if err != nil {
		http.Error(w, "failed to list phases", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items, "nextCursor": next})
}

type updatePhaseStatusReq struct {
	Status models.PhaseStatus `json:"status"` // pending|in_progress|blocked|done
}

func (h *PhasesHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	phaseID := chi.URLParam(r, "phaseId")
	tenant := middleware.TenantID(r.Context())

	var req updatePhaseStatusReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	st := strings.TrimSpace(string(req.Status))
	if st == "" {
		http.Error(w, "status required", http.StatusBadRequest)
		return
	}

	// Gate: if moving to done, ensure all WOs under this phase have:
	// - all deliverables approved
	// - approval_status is approved or not_required
	// - work order status is complete/approved
	if req.Status == models.PhaseDone {
		wos, err := h.pg.WorkOrders().ListByPhase(r.Context(), tenant, phaseID)
		if err != nil {
			http.Error(w, "failed to load work orders", http.StatusInternalServerError)
			return
		}

		for _, wo := range wos {
			cnt, err := h.pg.WorkOrderDeliverables().CountNotApprovedByWorkOrder(r.Context(), tenant, wo.SchoolID, wo.ID)
			if err != nil {
				http.Error(w, "failed to validate deliverables", http.StatusInternalServerError)
				return
			}
			if cnt > 0 {
				http.Error(w, "phase blocked: unapproved deliverables exist", http.StatusConflict)
				return
			}
			if wo.ApprovalStatus != "approved" && wo.ApprovalStatus != "not_required" {
				http.Error(w, "phase blocked: work order approvals pending", http.StatusConflict)
				return
			}
			if string(wo.Status) != "complete" && string(wo.Status) != "approved" {
				http.Error(w, "phase blocked: work orders not complete", http.StatusConflict)
				return
			}
		}
	}

	// Update phase status
	_, err := h.pg.RawPool().Exec(r.Context(), `
		UPDATE service_phases SET status=$3, updated_at=$4
		WHERE tenant_id=$1 AND id=$2
	`, tenant, phaseID, req.Status, time.Now().UTC())
	if err != nil {
		http.Error(w, "failed to update phase", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}
