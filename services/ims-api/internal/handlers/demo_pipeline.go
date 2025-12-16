package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/edvirons/ssp/ims/internal/middleware"
	"github.com/edvirons/ssp/ims/internal/models"
	"github.com/edvirons/ssp/ims/internal/store"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type DemoPipelineHandler struct {
	log *zap.Logger
	pg  *store.Postgres
}

func NewDemoPipelineHandler(log *zap.Logger, pg *store.Postgres) *DemoPipelineHandler {
	return &DemoPipelineHandler{log: log, pg: pg}
}

// ListLeads returns all leads with optional filtering.
func (h *DemoPipelineHandler) ListLeads(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())

	filters := models.DemoLeadFilters{
		Limit:  50,
		Offset: 0,
	}

	if v := r.URL.Query().Get("stage"); v != "" {
		stage := models.DemoLeadStage(v)
		filters.Stage = &stage
	}
	if v := r.URL.Query().Get("assignedTo"); v != "" {
		filters.AssignedTo = &v
	}
	if v := r.URL.Query().Get("source"); v != "" {
		source := models.DemoLeadSource(v)
		filters.LeadSource = &source
	}
	if v := r.URL.Query().Get("search"); v != "" {
		filters.Search = &v
	}
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			filters.Limit = n
		}
	}
	if v := r.URL.Query().Get("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			filters.Offset = n
		}
	}

	leads, total, err := h.pg.DemoLeads().List(r.Context(), tenant, filters)
	if err != nil {
		h.log.Error("failed to list leads", zap.Error(err))
		http.Error(w, "failed to list leads", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"leads": leads,
		"total": total,
	})
}

// GetLead returns a single lead by ID.
func (h *DemoPipelineHandler) GetLead(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())
	id := chi.URLParam(r, "id")

	lead, err := h.pg.DemoLeads().GetByID(r.Context(), tenant, id)
	if err != nil {
		if err.Error() == "not found" {
			http.Error(w, "lead not found", http.StatusNotFound)
			return
		}
		h.log.Error("failed to get lead", zap.Error(err))
		http.Error(w, "failed to get lead", http.StatusInternalServerError)
		return
	}

	// Get recent activities
	activities, _ := h.pg.DemoLeadActivities().ListByLead(r.Context(), tenant, id, 10)

	// Get next scheduled demo
	nextDemo, _ := h.pg.DemoSchedules().GetNextByLead(r.Context(), tenant, id)

	result := models.DemoLeadWithActivities{
		DemoLead:         lead,
		RecentActivities: activities,
		NextDemo:         nextDemo,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// CreateLead creates a new lead.
func (h *DemoPipelineHandler) CreateLead(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())
	userID := middleware.UserID(r.Context())

	var req models.CreateDemoLeadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.SchoolName) == "" {
		http.Error(w, "schoolName is required", http.StatusBadRequest)
		return
	}

	now := time.Now().UTC()
	lead := models.DemoLead{
		ID:               store.NewID("lead"),
		TenantID:         tenant,
		SchoolID:         req.SchoolID,
		SchoolName:       strings.TrimSpace(req.SchoolName),
		ContactName:      strings.TrimSpace(req.ContactName),
		ContactEmail:     strings.TrimSpace(req.ContactEmail),
		ContactPhone:     strings.TrimSpace(req.ContactPhone),
		ContactRole:      strings.TrimSpace(req.ContactRole),
		CountyCode:       strings.TrimSpace(req.CountyCode),
		CountyName:       strings.TrimSpace(req.CountyName),
		SubCountyCode:    strings.TrimSpace(req.SubCountyCode),
		SubCountyName:    strings.TrimSpace(req.SubCountyName),
		Stage:            models.StageNewLead,
		StageChangedAt:   now,
		EstimatedValue:   req.EstimatedValue,
		EstimatedDevices: req.EstimatedDevices,
		Probability:      models.StageProbability[models.StageNewLead],
		LeadSource:       models.DemoLeadSource(req.LeadSource),
		AssignedTo:       userID,
		Notes:            strings.TrimSpace(req.Notes),
		Tags:             req.Tags,
		CreatedBy:        userID,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	if lead.Tags == nil {
		lead.Tags = []string{}
	}

	if err := h.pg.DemoLeads().Create(r.Context(), lead); err != nil {
		h.log.Error("failed to create lead", zap.Error(err))
		http.Error(w, "failed to create lead", http.StatusInternalServerError)
		return
	}

	// Create activity for lead creation
	activity := models.DemoLeadActivity{
		ID:           store.NewID("act"),
		TenantID:     tenant,
		LeadID:       lead.ID,
		ActivityType: models.DemoActivityCreated,
		Description:  "Lead created",
		CreatedBy:    userID,
		CreatedAt:    now,
	}
	h.pg.DemoLeadActivities().Create(r.Context(), activity)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(lead)
}

// UpdateLead updates an existing lead.
func (h *DemoPipelineHandler) UpdateLead(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())
	id := chi.URLParam(r, "id")

	var req models.UpdateDemoLeadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if err := h.pg.DemoLeads().Update(r.Context(), tenant, id, req); err != nil {
		h.log.Error("failed to update lead", zap.Error(err))
		http.Error(w, "failed to update lead", http.StatusInternalServerError)
		return
	}

	// Return updated lead
	lead, err := h.pg.DemoLeads().GetByID(r.Context(), tenant, id)
	if err != nil {
		http.Error(w, "failed to get updated lead", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lead)
}

// UpdateLeadStage changes the stage of a lead.
func (h *DemoPipelineHandler) UpdateLeadStage(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())
	userID := middleware.UserID(r.Context())
	id := chi.URLParam(r, "id")

	var req models.UpdateLeadStageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	// Get current lead to record stage change
	currentLead, err := h.pg.DemoLeads().GetByID(r.Context(), tenant, id)
	if err != nil {
		if err.Error() == "not found" {
			http.Error(w, "lead not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to get lead", http.StatusInternalServerError)
		return
	}

	if err := h.pg.DemoLeads().UpdateStage(r.Context(), tenant, id, req.Stage, req.LostReason, req.LostNotes); err != nil {
		h.log.Error("failed to update stage", zap.Error(err))
		http.Error(w, "failed to update stage", http.StatusInternalServerError)
		return
	}

	// Create activity for stage change
	now := time.Now().UTC()
	fromStage := currentLead.Stage
	activity := models.DemoLeadActivity{
		ID:           store.NewID("act"),
		TenantID:     tenant,
		LeadID:       id,
		ActivityType: models.DemoActivityStageChange,
		Description:  "Stage changed from " + string(fromStage) + " to " + string(req.Stage),
		FromStage:    &fromStage,
		ToStage:      &req.Stage,
		CreatedBy:    userID,
		CreatedAt:    now,
	}
	h.pg.DemoLeadActivities().Create(r.Context(), activity)

	// Update daily metrics based on stage change
	h.updateMetricsForStageChange(r, tenant, fromStage, req.Stage, currentLead.EstimatedValue)

	// Return updated lead
	lead, _ := h.pg.DemoLeads().GetByID(r.Context(), tenant, id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lead)
}

func (h *DemoPipelineHandler) updateMetricsForStageChange(r *http.Request, tenant string, from, to models.DemoLeadStage, value *float64) {
	val := float64(0)
	if value != nil {
		val = *value
	}

	switch to {
	case models.StageDemoScheduled:
		h.pg.SalesMetricsDaily().IncrementMetric(r.Context(), tenant, "demos_scheduled", 1)
	case models.StageDemoCompleted:
		h.pg.SalesMetricsDaily().IncrementMetric(r.Context(), tenant, "demos_completed", 1)
	case models.StageProposalSent:
		h.pg.SalesMetricsDaily().IncrementMetric(r.Context(), tenant, "proposals_sent", 1)
	case models.StageWon:
		h.pg.SalesMetricsDaily().IncrementMetric(r.Context(), tenant, "deals_won", 1)
		h.pg.SalesMetricsDaily().IncrementMetric(r.Context(), tenant, "won_value", val)
	case models.StageLost:
		h.pg.SalesMetricsDaily().IncrementMetric(r.Context(), tenant, "deals_lost", 1)
		h.pg.SalesMetricsDaily().IncrementMetric(r.Context(), tenant, "lost_value", val)
	}
}

// DeleteLead soft deletes a lead.
func (h *DemoPipelineHandler) DeleteLead(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())
	id := chi.URLParam(r, "id")

	if err := h.pg.DemoLeads().Delete(r.Context(), tenant, id); err != nil {
		h.log.Error("failed to delete lead", zap.Error(err))
		http.Error(w, "failed to delete lead", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// AddNote adds a note to a lead.
func (h *DemoPipelineHandler) AddNote(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())
	userID := middleware.UserID(r.Context())
	id := chi.URLParam(r, "id")

	var req models.AddLeadNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Note) == "" {
		http.Error(w, "note is required", http.StatusBadRequest)
		return
	}

	now := time.Now().UTC()
	activity := models.DemoLeadActivity{
		ID:           store.NewID("act"),
		TenantID:     tenant,
		LeadID:       id,
		ActivityType: models.DemoActivityNote,
		Description:  strings.TrimSpace(req.Note),
		CreatedBy:    userID,
		CreatedAt:    now,
	}

	if err := h.pg.DemoLeadActivities().Create(r.Context(), activity); err != nil {
		h.log.Error("failed to add note", zap.Error(err))
		http.Error(w, "failed to add note", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(activity)
}

// ListActivities returns activities for a lead.
func (h *DemoPipelineHandler) ListActivities(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())
	id := chi.URLParam(r, "id")

	limit := 50
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			limit = n
		}
	}

	activities, err := h.pg.DemoLeadActivities().ListByLead(r.Context(), tenant, id, limit)
	if err != nil {
		h.log.Error("failed to list activities", zap.Error(err))
		http.Error(w, "failed to list activities", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"activities": activities,
	})
}

// ScheduleDemo schedules a demo for a lead.
func (h *DemoPipelineHandler) ScheduleDemo(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())
	userID := middleware.UserID(r.Context())
	id := chi.URLParam(r, "id")

	var req models.CreateDemoScheduleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	scheduledDate, err := time.Parse("2006-01-02", req.ScheduledDate)
	if err != nil {
		http.Error(w, "invalid scheduledDate format (expected YYYY-MM-DD)", http.StatusBadRequest)
		return
	}

	now := time.Now().UTC()
	schedule := models.DemoSchedule{
		ID:              store.NewID("demo"),
		TenantID:        tenant,
		LeadID:          id,
		ScheduledDate:   scheduledDate,
		ScheduledTime:   req.ScheduledTime,
		DurationMinutes: req.DurationMinutes,
		Location:        req.Location,
		MeetingLink:     req.MeetingLink,
		Attendees:       req.Attendees,
		Status:          models.ScheduleStatusScheduled,
		ReminderSent:    false,
		CreatedBy:       userID,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if schedule.DurationMinutes == 0 {
		schedule.DurationMinutes = 60
	}
	if schedule.Attendees == nil {
		schedule.Attendees = []models.DemoAttendee{}
	}

	if err := h.pg.DemoSchedules().Create(r.Context(), schedule); err != nil {
		h.log.Error("failed to schedule demo", zap.Error(err))
		http.Error(w, "failed to schedule demo", http.StatusInternalServerError)
		return
	}

	// Update lead stage to demo_scheduled if it's before that stage
	lead, _ := h.pg.DemoLeads().GetByID(r.Context(), tenant, id)
	if lead.Stage == models.StageNewLead || lead.Stage == models.StageContacted {
		h.pg.DemoLeads().UpdateStage(r.Context(), tenant, id, models.StageDemoScheduled, "", "")
	}

	// Create activity
	activity := models.DemoLeadActivity{
		ID:           store.NewID("act"),
		TenantID:     tenant,
		LeadID:       id,
		ActivityType: models.DemoActivityDemo,
		Description:  "Demo scheduled for " + req.ScheduledDate + " " + req.ScheduledTime,
		ScheduledAt:  &scheduledDate,
		CreatedBy:    userID,
		CreatedAt:    now,
	}
	h.pg.DemoLeadActivities().Create(r.Context(), activity)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(schedule)
}

// GetPipelineSummary returns a summary of the pipeline by stage.
func (h *DemoPipelineHandler) GetPipelineSummary(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())

	summary, err := h.pg.DemoLeads().GetPipelineSummary(r.Context(), tenant)
	if err != nil {
		h.log.Error("failed to get pipeline summary", zap.Error(err))
		http.Error(w, "failed to get summary", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}

// GetRecentActivities returns recent activities across all leads.
func (h *DemoPipelineHandler) GetRecentActivities(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())

	limit := 20
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			limit = n
		}
	}

	activities, err := h.pg.DemoLeadActivities().ListRecent(r.Context(), tenant, limit)
	if err != nil {
		h.log.Error("failed to list recent activities", zap.Error(err))
		http.Error(w, "failed to list activities", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"activities": activities,
	})
}
