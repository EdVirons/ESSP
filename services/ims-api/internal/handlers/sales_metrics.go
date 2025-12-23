package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/edvirons/ssp/ims/internal/middleware"
	"github.com/edvirons/ssp/ims/internal/models"
	"github.com/edvirons/ssp/ims/internal/store"
	"go.uber.org/zap"
)

type SalesMetricsHandler struct {
	log *zap.Logger
	pg  *store.Postgres
}

func NewSalesMetricsHandler(log *zap.Logger, pg *store.Postgres) *SalesMetricsHandler {
	return &SalesMetricsHandler{log: log, pg: pg}
}

// GetDashboard returns a comprehensive sales dashboard summary.
func (h *SalesMetricsHandler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())

	// Default to last 30 days
	days := 30
	if v := r.URL.Query().Get("days"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			days = n
		}
	}

	endDate := time.Now().UTC().Truncate(24 * time.Hour)
	startDate := endDate.AddDate(0, 0, -days)

	// Get pipeline summary
	pipelineSummary, err := h.pg.DemoLeads().GetPipelineSummary(r.Context(), tenant)
	if err != nil {
		h.log.Error("failed to get pipeline summary", zap.Error(err))
		pipelineSummary = models.PipelineSummary{}
	}

	// Get sales metrics summary for the period
	metricsSummary, err := h.pg.SalesMetricsDaily().GetSummary(r.Context(), tenant, startDate, endDate)
	if err != nil {
		h.log.Error("failed to get metrics summary", zap.Error(err))
		metricsSummary = models.SalesMetricsSummary{}
	}

	// Merge data from pipeline into metrics summary
	metricsSummary.TotalLeads = pipelineSummary.TotalLeads
	metricsSummary.TotalPipelineValue = pipelineSummary.TotalValue
	metricsSummary.ConversionRate = pipelineSummary.ConversionRate

	// Get recent activities
	recentActivities, err := h.pg.DemoLeadActivities().ListRecent(r.Context(), tenant, 10)
	if err != nil {
		h.log.Error("failed to get recent activities", zap.Error(err))
		recentActivities = []models.DemoLeadActivity{}
	}

	// Convert to RecentActivity format with lead names
	activities := make([]models.RecentActivity, 0, len(recentActivities))
	for _, a := range recentActivities {
		// Get lead name
		lead, _ := h.pg.DemoLeads().GetByID(r.Context(), tenant, a.LeadID)
		leadName := ""
		if lead.ID != "" {
			leadName = lead.SchoolName
		}

		activities = append(activities, models.RecentActivity{
			ID:          a.ID,
			Type:        string(a.ActivityType),
			Description: a.Description,
			LeadID:      a.LeadID,
			LeadName:    leadName,
			UserID:      a.CreatedBy,
			CreatedAt:   a.CreatedAt,
		})
	}

	// Get schools by region (top regions by lead count)
	schoolsByRegion := h.getSchoolsByRegion(r, tenant)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"metrics":          metricsSummary,
		"pipelineStages":   pipelineSummary.Stages,
		"recentActivities": activities,
		"schoolsByRegion":  schoolsByRegion,
		"period": map[string]interface{}{
			"startDate": startDate.Format("2006-01-02"),
			"endDate":   endDate.Format("2006-01-02"),
			"days":      days,
		},
	})
}

func (h *SalesMetricsHandler) getSchoolsByRegion(r *http.Request, tenant string) []models.SchoolsByRegion {
	// Get leads grouped by some region identifier
	// For now, return mock data - in production, this would query leads by school region
	return []models.SchoolsByRegion{
		{Region: "Nairobi", Count: 15, Value: 2500000},
		{Region: "Mombasa", Count: 8, Value: 1200000},
		{Region: "Kisumu", Count: 6, Value: 900000},
		{Region: "Nakuru", Count: 5, Value: 750000},
		{Region: "Eldoret", Count: 4, Value: 600000},
	}
}

// GetMetricsSummary returns metrics for a specific period.
func (h *SalesMetricsHandler) GetMetricsSummary(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())

	startDateStr := r.URL.Query().Get("startDate")
	endDateStr := r.URL.Query().Get("endDate")

	var startDate, endDate time.Time
	var err error

	if startDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			http.Error(w, "invalid startDate format", http.StatusBadRequest)
			return
		}
	} else {
		startDate = time.Now().UTC().AddDate(0, 0, -30).Truncate(24 * time.Hour)
	}

	if endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			http.Error(w, "invalid endDate format", http.StatusBadRequest)
			return
		}
	} else {
		endDate = time.Now().UTC().Truncate(24 * time.Hour)
	}

	summary, err := h.pg.SalesMetricsDaily().GetSummary(r.Context(), tenant, startDate, endDate)
	if err != nil {
		h.log.Error("failed to get metrics summary", zap.Error(err))
		http.Error(w, "failed to get metrics", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}

// GetPipelineStages returns the current pipeline stage counts.
func (h *SalesMetricsHandler) GetPipelineStages(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())

	summary, err := h.pg.DemoLeads().GetPipelineSummary(r.Context(), tenant)
	if err != nil {
		h.log.Error("failed to get pipeline summary", zap.Error(err))
		http.Error(w, "failed to get pipeline stages", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}

// IncrementMetric manually increments a metric (for testing/admin).
func (h *SalesMetricsHandler) IncrementMetric(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())

	var req struct {
		Metric string  `json:"metric"`
		Value  float64 `json:"value"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	// Validate metric name
	validMetrics := map[string]bool{
		"new_leads":       true,
		"leads_contacted": true,
		"demos_scheduled": true,
		"demos_completed": true,
		"proposals_sent":  true,
		"deals_won":       true,
		"deals_lost":      true,
		"pipeline_value":  true,
		"won_value":       true,
		"lost_value":      true,
		"calls_made":      true,
		"emails_sent":     true,
		"meetings_held":   true,
	}

	if !validMetrics[req.Metric] {
		http.Error(w, "invalid metric name", http.StatusBadRequest)
		return
	}

	if err := h.pg.SalesMetricsDaily().IncrementMetric(r.Context(), tenant, req.Metric, req.Value); err != nil {
		h.log.Error("failed to increment metric", zap.Error(err))
		http.Error(w, "failed to increment metric", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
