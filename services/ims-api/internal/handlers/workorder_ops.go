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

type WorkOrderOpsHandler struct {
	log *zap.Logger
	pg  *store.Postgres
}

func NewWorkOrderOpsHandler(log *zap.Logger, pg *store.Postgres) *WorkOrderOpsHandler {
	return &WorkOrderOpsHandler{log: log, pg: pg}
}

type scheduleReq struct {
	ScheduledStart  string `json:"scheduledStart"` // RFC3339
	ScheduledEnd    string `json:"scheduledEnd"`   // RFC3339
	Timezone        string `json:"timezone"`
	Notes           string `json:"notes"`
	CreatedByUserID string `json:"createdByUserId"`
}

func (h *WorkOrderOpsHandler) Schedule(w http.ResponseWriter, r *http.Request) {
	woID := chi.URLParam(r, "id")
	tenant := middleware.TenantID(r.Context())
	school := middleware.SchoolID(r.Context())

	var req scheduleReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	var startPtr *time.Time
	var endPtr *time.Time
	if strings.TrimSpace(req.ScheduledStart) != "" {
		if t, err := time.Parse(time.RFC3339, strings.TrimSpace(req.ScheduledStart)); err == nil {
			startPtr = &t
		}
	}
	if strings.TrimSpace(req.ScheduledEnd) != "" {
		if t, err := time.Parse(time.RFC3339, strings.TrimSpace(req.ScheduledEnd)); err == nil {
			endPtr = &t
		}
	}

	now := time.Now().UTC()
	// inherit phase_id from work order
	wo, err := h.pg.WorkOrders().GetByID(r.Context(), tenant, school, woID)
	var phaseID string
	if err == nil {
		phaseID = wo.PhaseID
	}

	z := strings.TrimSpace(req.Timezone)
	if z == "" {
		z = "Africa/Nairobi"
	}

	s := models.WorkOrderSchedule{
		ID:              store.NewID("sched"),
		TenantID:        tenant,
		SchoolID:        school,
		WorkOrderID:     woID,
		PhaseID:         phaseID,
		ScheduledStart:  startPtr,
		ScheduledEnd:    endPtr,
		Timezone:        z,
		Notes:           strings.TrimSpace(req.Notes),
		CreatedByUserID: strings.TrimSpace(req.CreatedByUserID),
		CreatedAt:       now,
	}
	if err := h.pg.WorkOrderSchedules().Create(r.Context(), s); err != nil {
		http.Error(w, "failed to create schedule", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, s)
}

func (h *WorkOrderOpsHandler) Schedules(w http.ResponseWriter, r *http.Request) {
	woID := chi.URLParam(r, "id")
	tenant := middleware.TenantID(r.Context())
	school := middleware.SchoolID(r.Context())
	limit := parseLimit(r.URL.Query().Get("limit"), 50, 200)
	curT, curID, hasCur := decodeCursor(strings.TrimSpace(r.URL.Query().Get("cursor")))

	items, next, err := h.pg.WorkOrderSchedules().List(r.Context(), store.ScheduleListParams{
		TenantID: tenant, SchoolID: school, WorkOrderID: woID,
		Limit: limit, HasCursor: hasCur, CursorCreatedAt: curT, CursorID: curID,
	})
	if err != nil {
		http.Error(w, "failed to list schedules", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items, "nextCursor": next})
}

type addDeliverableReq struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (h *WorkOrderOpsHandler) AddDeliverable(w http.ResponseWriter, r *http.Request) {
	woID := chi.URLParam(r, "id")
	tenant := middleware.TenantID(r.Context())
	school := middleware.SchoolID(r.Context())
	var req addDeliverableReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(req.Title) == "" {
		http.Error(w, "title required", http.StatusBadRequest)
		return
	}
	now := time.Now().UTC()
	// inherit phase_id from work order
	wo, err := h.pg.WorkOrders().GetByID(r.Context(), tenant, school, woID)
	var phaseID string
	if err == nil {
		phaseID = wo.PhaseID
	}

	d := models.WorkOrderDeliverable{
		ID:          store.NewID("deliv"),
		TenantID:    tenant,
		SchoolID:    school,
		WorkOrderID: woID,
		PhaseID:     phaseID,
		Title:       strings.TrimSpace(req.Title),
		Description: strings.TrimSpace(req.Description),
		Status:      models.DeliverablePending,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := h.pg.WorkOrderDeliverables().Create(r.Context(), d); err != nil {
		http.Error(w, "failed to create deliverable", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, d)
}

func (h *WorkOrderOpsHandler) Deliverables(w http.ResponseWriter, r *http.Request) {
	woID := chi.URLParam(r, "id")
	tenant := middleware.TenantID(r.Context())
	school := middleware.SchoolID(r.Context())
	status := strings.TrimSpace(r.URL.Query().Get("status"))
	limit := parseLimit(r.URL.Query().Get("limit"), 50, 200)
	curT, curID, hasCur := decodeCursor(strings.TrimSpace(r.URL.Query().Get("cursor")))

	items, next, err := h.pg.WorkOrderDeliverables().List(r.Context(), store.DeliverableListParams{
		TenantID: tenant, SchoolID: school, WorkOrderID: woID, Status: status,
		Limit: limit, HasCursor: hasCur, CursorCreatedAt: curT, CursorID: curID,
	})
	if err != nil {
		http.Error(w, "failed to list deliverables", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items, "nextCursor": next})
}

type submitDeliverableReq struct {
	EvidenceAttachmentID string `json:"evidenceAttachmentId"`
	SubmittedByUserID    string `json:"submittedByUserId"`
	Notes                string `json:"notes"`
}

func (h *WorkOrderOpsHandler) SubmitDeliverable(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "deliverableId")
	tenant := middleware.TenantID(r.Context())
	school := middleware.SchoolID(r.Context())
	var req submitDeliverableReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(req.EvidenceAttachmentID) == "" {
		http.Error(w, "evidenceAttachmentId required", http.StatusBadRequest)
		return
	}

	if err := h.pg.WorkOrderDeliverables().MarkSubmitted(r.Context(), tenant, school, id, strings.TrimSpace(req.SubmittedByUserID), strings.TrimSpace(req.EvidenceAttachmentID), strings.TrimSpace(req.Notes)); err != nil {
		http.Error(w, "failed to submit deliverable", http.StatusInternalServerError)
		return
	}
	updated, _ := h.pg.WorkOrderDeliverables().GetByID(r.Context(), tenant, school, id)
	writeJSON(w, http.StatusOK, updated)
}

type reviewDeliverableReq struct {
	ReviewerUserID string `json:"reviewerUserId"`
	Status         string `json:"status"` // approved|rejected
	Notes          string `json:"notes"`
}

func (h *WorkOrderOpsHandler) ReviewDeliverable(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "deliverableId")
	tenant := middleware.TenantID(r.Context())
	school := middleware.SchoolID(r.Context())
	var req reviewDeliverableReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if err := h.pg.WorkOrderDeliverables().Review(r.Context(), tenant, school, id, strings.TrimSpace(req.ReviewerUserID), strings.TrimSpace(req.Status), strings.TrimSpace(req.Notes)); err != nil {
		http.Error(w, "failed to review", http.StatusBadRequest)
		return
	}
	updated, _ := h.pg.WorkOrderDeliverables().GetByID(r.Context(), tenant, school, id)
	writeJSON(w, http.StatusOK, updated)
}

type requestApprovalReq struct {
	ApprovalType      string `json:"approvalType"`
	RequestedByUserID string `json:"requestedByUserId"`
}

func (h *WorkOrderOpsHandler) RequestApproval(w http.ResponseWriter, r *http.Request) {
	woID := chi.URLParam(r, "id")
	tenant := middleware.TenantID(r.Context())
	school := middleware.SchoolID(r.Context())
	var req requestApprovalReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	at := strings.TrimSpace(req.ApprovalType)
	if at == "" {
		at = "school_signoff"
	}
	now := time.Now().UTC()
	// inherit phase_id from work order
	wo, err := h.pg.WorkOrders().GetByID(r.Context(), tenant, school, woID)
	var phaseID string
	if err == nil {
		phaseID = wo.PhaseID
	}

	a := models.WorkOrderApproval{
		ID:                store.NewID("appr"),
		TenantID:          tenant,
		SchoolID:          school,
		WorkOrderID:       woID,
		PhaseID:           phaseID,
		ApprovalType:      at,
		RequestedByUserID: strings.TrimSpace(req.RequestedByUserID),
		RequestedAt:       now,
		Status:            models.ApprovalPending,
	}
	if err := h.pg.WorkOrderApprovals().Request(r.Context(), a); err != nil {
		http.Error(w, "failed to request approval", http.StatusInternalServerError)
		return
	}
	_ = h.pg.WorkOrders().SetApprovalStatus(r.Context(), tenant, school, woID, "pending")
	writeJSON(w, http.StatusCreated, a)
}

type decideApprovalReq struct {
	DecidedByUserID string `json:"decidedByUserId"`
	Status          string `json:"status"` // approved|rejected
	Notes           string `json:"notes"`
}

func (h *WorkOrderOpsHandler) DecideApproval(w http.ResponseWriter, r *http.Request) {
	approvalID := chi.URLParam(r, "approvalId")
	woID := chi.URLParam(r, "id")
	tenant := middleware.TenantID(r.Context())
	school := middleware.SchoolID(r.Context())
	var req decideApprovalReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if err := h.pg.WorkOrderApprovals().Decide(r.Context(), tenant, school, approvalID, strings.TrimSpace(req.DecidedByUserID), strings.TrimSpace(req.Status), strings.TrimSpace(req.Notes)); err != nil {
		http.Error(w, "failed to decide", http.StatusBadRequest)
		return
	}
	st := strings.TrimSpace(req.Status)
	_ = h.pg.WorkOrders().SetApprovalStatus(r.Context(), tenant, school, woID, st)
	updated, _ := h.pg.WorkOrderApprovals().GetByID(r.Context(), tenant, school, approvalID)
	writeJSON(w, http.StatusOK, updated)
}
