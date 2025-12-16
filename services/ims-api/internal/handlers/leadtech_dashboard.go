package handlers

import (
	"net/http"
	"time"

	"github.com/edvirons/ssp/ims/internal/middleware"
	"github.com/edvirons/ssp/ims/internal/store"
	"go.uber.org/zap"
)

// LeadTechDashboardHandler provides aggregated data for the lead tech dashboard
type LeadTechDashboardHandler struct {
	log *zap.Logger
	pg  *store.Postgres
}

// NewLeadTechDashboardHandler creates a new lead tech dashboard handler
func NewLeadTechDashboardHandler(log *zap.Logger, pg *store.Postgres) *LeadTechDashboardHandler {
	return &LeadTechDashboardHandler{log: log, pg: pg}
}

// PendingApproval represents a work order awaiting approval
type PendingApproval struct {
	ID              string `json:"id"`
	WorkOrderID     string `json:"workOrderId"`
	Title           string `json:"title"`
	SchoolName      string `json:"schoolName"`
	Priority        string `json:"priority"`
	RequestedAt     string `json:"requestedAt"`
	RequestedByName string `json:"requestedByName"`
}

// ScheduledWorkOrder represents a work order scheduled for today/this week
type ScheduledWorkOrder struct {
	ID             string  `json:"id"`
	Title          string  `json:"title"`
	SchoolName     string  `json:"schoolName"`
	ScheduledStart *string `json:"scheduledStart"`
	AssignedTo     string  `json:"assignedTo"`
	Status         string  `json:"status"`
	Priority       string  `json:"priority"`
}

// TeamWorkOrderMetrics represents work order stats by status
type TeamWorkOrderMetrics struct {
	InProgress int `json:"inProgress"`
	Completed  int `json:"completed"`
	Pending    int `json:"pending"`
	Scheduled  int `json:"scheduled"`
}

// BOMReadinessItem represents parts availability for a scheduled job
type BOMReadinessItem struct {
	WorkOrderID    string `json:"workOrderId"`
	WorkOrderTitle string `json:"workOrderTitle"`
	TotalParts     int    `json:"totalParts"`
	AvailableParts int    `json:"availableParts"`
	MissingParts   int    `json:"missingParts"`
	IsReady        bool   `json:"isReady"`
}

// TeamActivityItem represents recent team activity
type TeamActivityItem struct {
	ID          string `json:"id"`
	ActorName   string `json:"actorName"`
	Action      string `json:"action"`
	EntityType  string `json:"entityType"`
	Description string `json:"description"`
	CreatedAt   string `json:"createdAt"`
}

// LeadTechDashboardSummary contains all dashboard data
type LeadTechDashboardSummary struct {
	PendingApprovalsCount int                  `json:"pendingApprovalsCount"`
	TodaysScheduledCount  int                  `json:"todaysScheduledCount"`
	TeamMetrics           TeamWorkOrderMetrics `json:"teamMetrics"`
	PendingApprovals      []PendingApproval    `json:"pendingApprovals"`
	TodaysSchedule        []ScheduledWorkOrder `json:"todaysSchedule"`
	BOMReadiness          []BOMReadinessItem   `json:"bomReadiness"`
	RecentTeamActivity    []TeamActivityItem   `json:"recentTeamActivity"`
}

// GetDashboardSummary returns aggregated dashboard data for lead tech
func (h *LeadTechDashboardHandler) GetDashboardSummary(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())
	ctx := r.Context()
	pool := h.pg.RawPool()

	summary := LeadTechDashboardSummary{
		TeamMetrics:        TeamWorkOrderMetrics{},
		PendingApprovals:   []PendingApproval{},
		TodaysSchedule:     []ScheduledWorkOrder{},
		BOMReadiness:       []BOMReadinessItem{},
		RecentTeamActivity: []TeamActivityItem{},
	}

	// Get pending approvals
	approvalRows, err := pool.Query(ctx, `
		SELECT
			w.id,
			COALESCE(w.notes, 'Work Order ' || LEFT(w.id, 8)) as title,
			COALESCE(s.county_name, 'Unknown School') as school_name,
			COALESCE(w.priority, 'medium') as priority,
			w.created_at
		FROM work_orders w
		LEFT JOIN schools s ON s.tenant_id = w.tenant_id AND s.school_id = w.school_id
		WHERE w.tenant_id = $1
			AND w.status = 'pending_approval'
		ORDER BY
			CASE COALESCE(w.priority, 'medium')
				WHEN 'critical' THEN 1
				WHEN 'high' THEN 2
				WHEN 'medium' THEN 3
				ELSE 4
			END,
			w.created_at ASC
		LIMIT 10
	`, tenant)
	if err != nil {
		h.log.Error("failed to query pending approvals", zap.Error(err))
	} else {
		defer approvalRows.Close()
		for approvalRows.Next() {
			var approval PendingApproval
			var createdAt time.Time
			if err := approvalRows.Scan(
				&approval.WorkOrderID, &approval.Title, &approval.SchoolName,
				&approval.Priority, &createdAt,
			); err != nil {
				continue
			}
			approval.ID = approval.WorkOrderID
			approval.RequestedAt = createdAt.UTC().Format(time.RFC3339)
			approval.RequestedByName = "System"
			summary.PendingApprovals = append(summary.PendingApprovals, approval)
		}
		summary.PendingApprovalsCount = len(summary.PendingApprovals)
	}

	// Get today's scheduled work orders
	today := time.Now().UTC().Truncate(24 * time.Hour)
	tomorrow := today.Add(24 * time.Hour)

	scheduleRows, err := pool.Query(ctx, `
		SELECT
			w.id,
			COALESCE(w.notes, 'Work Order ' || LEFT(w.id, 8)) as title,
			COALESCE(sch.county_name, 'Unknown School') as school_name,
			ws.scheduled_start,
			COALESCE(w.assigned_to, '') as assigned_to,
			COALESCE(w.status, 'open') as status,
			COALESCE(w.priority, 'medium') as priority
		FROM work_order_schedules ws
		JOIN work_orders w ON w.id = ws.work_order_id AND w.tenant_id = ws.tenant_id
		LEFT JOIN schools sch ON sch.tenant_id = w.tenant_id AND sch.school_id = w.school_id
		WHERE ws.tenant_id = $1
			AND ws.scheduled_start IS NOT NULL
			AND ws.scheduled_start >= $2
			AND ws.scheduled_start < $3
		ORDER BY ws.scheduled_start ASC
		LIMIT 20
	`, tenant, today, tomorrow)
	if err != nil {
		h.log.Error("failed to query scheduled work orders", zap.Error(err))
	} else {
		defer scheduleRows.Close()
		for scheduleRows.Next() {
			var wo ScheduledWorkOrder
			var scheduledStart *time.Time
			if err := scheduleRows.Scan(
				&wo.ID, &wo.Title, &wo.SchoolName,
				&scheduledStart, &wo.AssignedTo, &wo.Status, &wo.Priority,
			); err != nil {
				continue
			}
			if scheduledStart != nil {
				t := scheduledStart.UTC().Format(time.RFC3339)
				wo.ScheduledStart = &t
			}
			summary.TodaysSchedule = append(summary.TodaysSchedule, wo)
		}
		summary.TodaysScheduledCount = len(summary.TodaysSchedule)
	}

	// Get team work order metrics
	metricsRows, err := pool.Query(ctx, `
		SELECT
			COALESCE(status, 'open') as status,
			COUNT(*) as count
		FROM work_orders
		WHERE tenant_id = $1
		GROUP BY status
	`, tenant)
	if err != nil {
		h.log.Error("failed to query team metrics", zap.Error(err))
	} else {
		defer metricsRows.Close()
		for metricsRows.Next() {
			var status string
			var count int
			if err := metricsRows.Scan(&status, &count); err != nil {
				continue
			}
			switch status {
			case "in_repair", "in_progress":
				summary.TeamMetrics.InProgress += count
			case "completed", "approved":
				summary.TeamMetrics.Completed += count
			case "open", "assigned", "pending_approval":
				summary.TeamMetrics.Pending += count
			case "scheduled":
				summary.TeamMetrics.Scheduled += count
			}
		}
	}

	// Get BOM readiness for scheduled jobs
	bomRows, err := pool.Query(ctx, `
		SELECT
			w.id,
			COALESCE(w.notes, 'Work Order ' || LEFT(w.id, 8)) as title,
			COUNT(DISTINCT wop.part_id) as total_parts,
			COUNT(DISTINCT CASE WHEN i.qty_available >= wop.qty_planned THEN wop.part_id END) as available_parts
		FROM work_order_schedules ws
		JOIN work_orders w ON w.id = ws.work_order_id AND w.tenant_id = ws.tenant_id
		LEFT JOIN work_order_parts wop ON wop.work_order_id = w.id AND wop.tenant_id = w.tenant_id
		LEFT JOIN inventory i ON i.part_id = wop.part_id AND i.tenant_id = w.tenant_id
		WHERE ws.tenant_id = $1
			AND ws.scheduled_start IS NOT NULL
			AND ws.scheduled_start >= $2
			AND ws.scheduled_start < $3
			AND w.status NOT IN ('completed', 'approved', 'cancelled')
		GROUP BY w.id, w.notes, ws.scheduled_start
		HAVING COUNT(DISTINCT wop.part_id) > 0
		ORDER BY ws.scheduled_start ASC
		LIMIT 10
	`, tenant, today, tomorrow.Add(7*24*time.Hour)) // Look ahead 7 days
	if err != nil {
		h.log.Error("failed to query BOM readiness", zap.Error(err))
	} else {
		defer bomRows.Close()
		for bomRows.Next() {
			var item BOMReadinessItem
			if err := bomRows.Scan(
				&item.WorkOrderID, &item.WorkOrderTitle,
				&item.TotalParts, &item.AvailableParts,
			); err != nil {
				continue
			}
			item.MissingParts = item.TotalParts - item.AvailableParts
			item.IsReady = item.MissingParts == 0
			summary.BOMReadiness = append(summary.BOMReadiness, item)
		}
	}

	// Get recent team activity from audit logs
	activityRows, err := pool.Query(ctx, `
		SELECT
			a.id,
			COALESCE(a.actor_name, 'System') as actor_name,
			a.action,
			a.entity_type,
			a.created_at
		FROM audit_logs a
		WHERE a.tenant_id = $1
			AND a.entity_type IN ('work_order', 'bom_item', 'work_order_approval')
		ORDER BY a.created_at DESC
		LIMIT 10
	`, tenant)
	if err != nil {
		h.log.Error("failed to query recent activity", zap.Error(err))
	} else {
		defer activityRows.Close()
		for activityRows.Next() {
			var activity TeamActivityItem
			var createdAt time.Time
			if err := activityRows.Scan(
				&activity.ID, &activity.ActorName, &activity.Action,
				&activity.EntityType, &createdAt,
			); err != nil {
				continue
			}
			// Generate description based on action and entity type
			switch activity.EntityType {
			case "work_order":
				switch activity.Action {
				case "create":
					activity.Description = "Created work order"
				case "update":
					activity.Description = "Updated work order"
				case "delete":
					activity.Description = "Deleted work order"
				default:
					activity.Description = "Modified work order"
				}
			case "bom_item":
				switch activity.Action {
				case "create":
					activity.Description = "Added BOM item"
				case "consume":
					activity.Description = "Consumed BOM parts"
				default:
					activity.Description = "Modified BOM"
				}
			case "work_order_approval":
				switch activity.Action {
				case "approve":
					activity.Description = "Approved work order"
				case "reject":
					activity.Description = "Rejected work order"
				default:
					activity.Description = "Processed approval"
				}
			default:
				activity.Description = activity.Action + " " + activity.EntityType
			}
			activity.CreatedAt = createdAt.UTC().Format(time.RFC3339)
			summary.RecentTeamActivity = append(summary.RecentTeamActivity, activity)
		}
	}

	writeJSON(w, http.StatusOK, summary)
}

// GetPendingApprovals returns paginated list of pending approvals
func (h *LeadTechDashboardHandler) GetPendingApprovals(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())
	ctx := r.Context()
	pool := h.pg.RawPool()
	limit := parseLimit(r.URL.Query().Get("limit"), 20, 100)

	rows, err := pool.Query(ctx, `
		SELECT
			w.id,
			COALESCE(w.notes, 'Work Order ' || LEFT(w.id, 8)) as title,
			COALESCE(s.county_name, 'Unknown School') as school_name,
			COALESCE(w.priority, 'medium') as priority,
			w.created_at
		FROM work_orders w
		LEFT JOIN schools s ON s.tenant_id = w.tenant_id AND s.school_id = w.school_id
		WHERE w.tenant_id = $1
			AND w.status = 'pending_approval'
		ORDER BY
			CASE COALESCE(w.priority, 'medium')
				WHEN 'critical' THEN 1
				WHEN 'high' THEN 2
				WHEN 'medium' THEN 3
				ELSE 4
			END,
			w.created_at ASC
		LIMIT $2
	`, tenant, limit)
	if err != nil {
		h.log.Error("failed to query pending approvals", zap.Error(err))
		http.Error(w, "failed to query pending approvals", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	items := []PendingApproval{}
	for rows.Next() {
		var approval PendingApproval
		var createdAt time.Time
		if err := rows.Scan(
			&approval.WorkOrderID, &approval.Title, &approval.SchoolName,
			&approval.Priority, &createdAt,
		); err != nil {
			continue
		}
		approval.ID = approval.WorkOrderID
		approval.RequestedAt = createdAt.UTC().Format(time.RFC3339)
		approval.RequestedByName = "System"
		items = append(items, approval)
	}

	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// GetTodaysSchedule returns today's scheduled work orders
func (h *LeadTechDashboardHandler) GetTodaysSchedule(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())
	ctx := r.Context()
	pool := h.pg.RawPool()
	limit := parseLimit(r.URL.Query().Get("limit"), 20, 100)

	today := time.Now().UTC().Truncate(24 * time.Hour)
	tomorrow := today.Add(24 * time.Hour)

	rows, err := pool.Query(ctx, `
		SELECT
			w.id,
			COALESCE(w.notes, 'Work Order ' || LEFT(w.id, 8)) as title,
			COALESCE(sch.county_name, 'Unknown School') as school_name,
			ws.scheduled_start,
			COALESCE(w.assigned_to, '') as assigned_to,
			COALESCE(w.status, 'open') as status,
			COALESCE(w.priority, 'medium') as priority
		FROM work_order_schedules ws
		JOIN work_orders w ON w.id = ws.work_order_id AND w.tenant_id = ws.tenant_id
		LEFT JOIN schools sch ON sch.tenant_id = w.tenant_id AND sch.school_id = w.school_id
		WHERE ws.tenant_id = $1
			AND ws.scheduled_start IS NOT NULL
			AND ws.scheduled_start >= $2
			AND ws.scheduled_start < $3
		ORDER BY ws.scheduled_start ASC
		LIMIT $4
	`, tenant, today, tomorrow, limit)
	if err != nil {
		h.log.Error("failed to query today's schedule", zap.Error(err))
		http.Error(w, "failed to query schedule", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	items := []ScheduledWorkOrder{}
	for rows.Next() {
		var wo ScheduledWorkOrder
		var scheduledStart *time.Time
		if err := rows.Scan(
			&wo.ID, &wo.Title, &wo.SchoolName,
			&scheduledStart, &wo.AssignedTo, &wo.Status, &wo.Priority,
		); err != nil {
			continue
		}
		if scheduledStart != nil {
			t := scheduledStart.UTC().Format(time.RFC3339)
			wo.ScheduledStart = &t
		}
		items = append(items, wo)
	}

	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// GetTeamMetrics returns work order metrics for the team
func (h *LeadTechDashboardHandler) GetTeamMetrics(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())
	ctx := r.Context()
	pool := h.pg.RawPool()

	metrics := TeamWorkOrderMetrics{}

	rows, err := pool.Query(ctx, `
		SELECT
			COALESCE(status, 'open') as status,
			COUNT(*) as count
		FROM work_orders
		WHERE tenant_id = $1
		GROUP BY status
	`, tenant)
	if err != nil {
		h.log.Error("failed to query team metrics", zap.Error(err))
		http.Error(w, "failed to query metrics", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			continue
		}
		switch status {
		case "in_repair", "in_progress":
			metrics.InProgress += count
		case "completed", "approved":
			metrics.Completed += count
		case "open", "assigned", "pending_approval":
			metrics.Pending += count
		case "scheduled":
			metrics.Scheduled += count
		}
	}

	writeJSON(w, http.StatusOK, metrics)
}
