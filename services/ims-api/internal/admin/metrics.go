package admin

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// MetricsResponse represents the dashboard metrics summary
type MetricsResponse struct {
	Incidents  IncidentMetrics  `json:"incidents"`
	WorkOrders WorkOrderMetrics `json:"workOrders"`
	Programs   ProgramMetrics   `json:"programs"`
}

// IncidentMetrics contains incident statistics
type IncidentMetrics struct {
	Total       int64 `json:"total"`
	Open        int64 `json:"open"`
	SLABreached int64 `json:"slaBreached"`
}

// WorkOrderMetrics contains work order statistics
type WorkOrderMetrics struct {
	Total          int64 `json:"total"`
	InProgress     int64 `json:"inProgress"`
	CompletedToday int64 `json:"completedToday"`
}

// ProgramMetrics contains program statistics
type ProgramMetrics struct {
	Active  int64 `json:"active"`
	Pending int64 `json:"pending"`
}

// GetMetricsSummary returns dashboard metrics summary
func (h *Handler) GetMetricsSummary(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	pool := h.pg.RawPool()

	var resp MetricsResponse

	// Get incident metrics
	err := pool.QueryRow(ctx, `
		SELECT
			COUNT(*) as total,
			COUNT(*) FILTER (WHERE status IN ('new', 'acknowledged', 'in_progress', 'escalated')) as open,
			COUNT(*) FILTER (WHERE sla_breached = true) as sla_breached
		FROM incidents
	`).Scan(&resp.Incidents.Total, &resp.Incidents.Open, &resp.Incidents.SLABreached)
	if err != nil {
		h.logger.Error("failed to query incident metrics", zap.Error(err))
	}

	// Get work order metrics
	todayStart := time.Now().Truncate(24 * time.Hour)
	err = pool.QueryRow(ctx, `
		SELECT
			COUNT(*) as total,
			COUNT(*) FILTER (WHERE status IN ('assigned', 'in_repair', 'qa')) as in_progress,
			COUNT(*) FILTER (WHERE status = 'completed' AND updated_at >= $1) as completed_today
		FROM work_orders
	`, todayStart).Scan(&resp.WorkOrders.Total, &resp.WorkOrders.InProgress, &resp.WorkOrders.CompletedToday)
	if err != nil {
		h.logger.Error("failed to query work order metrics", zap.Error(err))
	}

	// Get program metrics
	err = pool.QueryRow(ctx, `
		SELECT
			COUNT(*) FILTER (WHERE status = 'active') as active,
			COUNT(*) FILTER (WHERE status = 'paused') as pending
		FROM school_service_programs
	`).Scan(&resp.Programs.Active, &resp.Programs.Pending)
	if err != nil {
		h.logger.Error("failed to query program metrics", zap.Error(err))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
