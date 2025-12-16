package handlers

import (
	"net/http"
	"time"

	"github.com/edvirons/ssp/ims/internal/middleware"
	"github.com/edvirons/ssp/ims/internal/store"
	"go.uber.org/zap"
)

// SupportAgentDashboardHandler provides aggregated data for the support agent dashboard
type SupportAgentDashboardHandler struct {
	log *zap.Logger
	pg  *store.Postgres
}

// NewSupportAgentDashboardHandler creates a new support agent dashboard handler
func NewSupportAgentDashboardHandler(log *zap.Logger, pg *store.Postgres) *SupportAgentDashboardHandler {
	return &SupportAgentDashboardHandler{log: log, pg: pg}
}

// IncidentQueueItem represents an incident in the queue
type IncidentQueueItem struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	SchoolName  string  `json:"schoolName"`
	Category    string  `json:"category"`
	Severity    string  `json:"severity"`
	Status      string  `json:"status"`
	ReportedBy  string  `json:"reportedBy"`
	SLADueAt    *string `json:"slaDueAt"`
	SLABreached bool    `json:"slaBreached"`
	CreatedAt   string  `json:"createdAt"`
}

// ChatQueueItem represents a chat session in the queue
type ChatQueueItem struct {
	ID                string  `json:"id"`
	SchoolName        string  `json:"schoolName"`
	ContactName       string  `json:"contactName"`
	Status            string  `json:"status"`
	QueuePosition     *int    `json:"queuePosition"`
	AssignedAgentName *string `json:"assignedAgentName"`
	StartedAt         string  `json:"startedAt"`
	WaitTimeSeconds   int     `json:"waitTimeSeconds"`
}

// WorkOrderQueueItem represents a work order for support agent view
type WorkOrderQueueItem struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	SchoolName string `json:"schoolName"`
	Status     string `json:"status"`
	TaskType   string `json:"taskType"`
	AssignedTo string `json:"assignedTo"`
	CreatedAt  string `json:"createdAt"`
}

// SupportAgentActivity represents recent activity for support agent
type SupportAgentActivity struct {
	ID          string `json:"id"`
	ActorName   string `json:"actorName"`
	Action      string `json:"action"`
	EntityType  string `json:"entityType"`
	Description string `json:"description"`
	CreatedAt   string `json:"createdAt"`
}

// IncidentMetrics represents incident statistics
type IncidentMetrics struct {
	Open        int `json:"open"`
	InProgress  int `json:"inProgress"`
	Resolved    int `json:"resolved"`
	SLABreached int `json:"slaBreached"`
}

// SupportAgentDashboardSummary contains all dashboard data for support agent
type SupportAgentDashboardSummary struct {
	// Counts
	OpenIncidentsCount  int `json:"openIncidentsCount"`
	WaitingChatsCount   int `json:"waitingChatsCount"`
	ActiveChatsCount    int `json:"activeChatsCount"`
	ActiveWorkOrders    int `json:"activeWorkOrders"`
	UnreadMessagesCount int `json:"unreadMessagesCount"`

	// Metrics
	IncidentMetrics IncidentMetrics `json:"incidentMetrics"`

	// Lists
	IncidentQueue    []IncidentQueueItem    `json:"incidentQueue"`
	ChatQueue        []ChatQueueItem        `json:"chatQueue"`
	WorkOrderQueue   []WorkOrderQueueItem   `json:"workOrderQueue"`
	RecentActivity   []SupportAgentActivity `json:"recentActivity"`
}

// GetDashboardSummary returns aggregated dashboard data for support agent
func (h *SupportAgentDashboardHandler) GetDashboardSummary(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())
	ctx := r.Context()
	pool := h.pg.RawPool()

	summary := SupportAgentDashboardSummary{
		IncidentMetrics: IncidentMetrics{},
		IncidentQueue:   []IncidentQueueItem{},
		ChatQueue:       []ChatQueueItem{},
		WorkOrderQueue:  []WorkOrderQueueItem{},
		RecentActivity:  []SupportAgentActivity{},
	}

	// Get incident metrics and queue
	incidentRows, err := pool.Query(ctx, `
		SELECT
			i.id,
			i.title,
			COALESCE(s.county_name, 'Unknown School') as school_name,
			i.category,
			i.severity,
			i.status,
			COALESCE(i.reported_by, '') as reported_by,
			i.sla_due_at,
			i.sla_breached,
			i.created_at
		FROM incidents i
		LEFT JOIN schools s ON s.tenant_id = i.tenant_id AND s.school_id = i.school_id
		WHERE i.tenant_id = $1
			AND i.status NOT IN ('closed', 'resolved')
		ORDER BY
			i.sla_breached DESC,
			CASE i.severity
				WHEN 'critical' THEN 1
				WHEN 'high' THEN 2
				WHEN 'medium' THEN 3
				ELSE 4
			END,
			i.sla_due_at ASC,
			i.created_at ASC
		LIMIT 15
	`, tenant)
	if err != nil {
		h.log.Error("failed to query incidents", zap.Error(err))
	} else {
		defer incidentRows.Close()
		for incidentRows.Next() {
			var item IncidentQueueItem
			var createdAt time.Time
			var slaDueAt *time.Time
			if err := incidentRows.Scan(
				&item.ID, &item.Title, &item.SchoolName,
				&item.Category, &item.Severity, &item.Status,
				&item.ReportedBy, &slaDueAt, &item.SLABreached, &createdAt,
			); err != nil {
				continue
			}
			if slaDueAt != nil {
				t := slaDueAt.UTC().Format(time.RFC3339)
				item.SLADueAt = &t
			}
			item.CreatedAt = createdAt.UTC().Format(time.RFC3339)
			summary.IncidentQueue = append(summary.IncidentQueue, item)
		}
	}

	// Get incident metrics by status
	metricsRows, err := pool.Query(ctx, `
		SELECT
			status,
			sla_breached,
			COUNT(*) as count
		FROM incidents
		WHERE tenant_id = $1
		GROUP BY status, sla_breached
	`, tenant)
	if err != nil {
		h.log.Error("failed to query incident metrics", zap.Error(err))
	} else {
		defer metricsRows.Close()
		for metricsRows.Next() {
			var status string
			var slaBreached bool
			var count int
			if err := metricsRows.Scan(&status, &slaBreached, &count); err != nil {
				continue
			}
			if slaBreached {
				summary.IncidentMetrics.SLABreached += count
			}
			switch status {
			case "open", "new", "assigned":
				summary.IncidentMetrics.Open += count
			case "in_progress", "investigating":
				summary.IncidentMetrics.InProgress += count
			case "resolved", "closed":
				summary.IncidentMetrics.Resolved += count
			}
		}
	}
	summary.OpenIncidentsCount = summary.IncidentMetrics.Open + summary.IncidentMetrics.InProgress

	// Get chat queue (waiting and active sessions)
	chatRows, err := pool.Query(ctx, `
		SELECT
			cs.id,
			COALESCE(s.county_name, 'Unknown School') as school_name,
			COALESCE(cs.school_contact_name, '') as contact_name,
			cs.status,
			cs.queue_position,
			cs.assigned_agent_name,
			cs.started_at
		FROM chat_sessions cs
		LEFT JOIN schools s ON s.tenant_id = cs.tenant_id AND s.school_id = cs.school_id
		WHERE cs.tenant_id = $1
			AND cs.status IN ('waiting', 'active')
		ORDER BY
			CASE cs.status
				WHEN 'waiting' THEN 1
				WHEN 'active' THEN 2
			END,
			cs.started_at ASC
		LIMIT 15
	`, tenant)
	if err != nil {
		h.log.Error("failed to query chat sessions", zap.Error(err))
	} else {
		defer chatRows.Close()
		now := time.Now()
		for chatRows.Next() {
			var item ChatQueueItem
			var startedAt time.Time
			if err := chatRows.Scan(
				&item.ID, &item.SchoolName, &item.ContactName,
				&item.Status, &item.QueuePosition, &item.AssignedAgentName,
				&startedAt,
			); err != nil {
				continue
			}
			item.StartedAt = startedAt.UTC().Format(time.RFC3339)
			item.WaitTimeSeconds = int(now.Sub(startedAt).Seconds())
			summary.ChatQueue = append(summary.ChatQueue, item)

			// Count by status
			if item.Status == "waiting" {
				summary.WaitingChatsCount++
			} else if item.Status == "active" {
				summary.ActiveChatsCount++
			}
		}
	}

	// Get active work orders
	woRows, err := pool.Query(ctx, `
		SELECT
			w.id,
			COALESCE(w.notes, 'Work Order ' || LEFT(w.id, 8)) as title,
			COALESCE(s.county_name, 'Unknown School') as school_name,
			w.status,
			COALESCE(w.task_type, '') as task_type,
			COALESCE(w.assigned_to, '') as assigned_to,
			w.created_at
		FROM work_orders w
		LEFT JOIN schools s ON s.tenant_id = w.tenant_id AND s.school_id = w.school_id
		WHERE w.tenant_id = $1
			AND w.status NOT IN ('completed', 'approved', 'cancelled', 'closed')
		ORDER BY
			w.created_at DESC
		LIMIT 10
	`, tenant)
	if err != nil {
		h.log.Error("failed to query work orders", zap.Error(err))
	} else {
		defer woRows.Close()
		for woRows.Next() {
			var item WorkOrderQueueItem
			var createdAt time.Time
			if err := woRows.Scan(
				&item.ID, &item.Title, &item.SchoolName,
				&item.Status, &item.TaskType, &item.AssignedTo, &createdAt,
			); err != nil {
				continue
			}
			item.CreatedAt = createdAt.UTC().Format(time.RFC3339)
			summary.WorkOrderQueue = append(summary.WorkOrderQueue, item)
		}
		summary.ActiveWorkOrders = len(summary.WorkOrderQueue)
	}

	// Get unread messages count for support team
	var unreadCount int
	err = pool.QueryRow(ctx, `
		SELECT COALESCE(SUM(unread_count_support), 0)
		FROM message_threads
		WHERE tenant_id = $1
			AND status = 'open'
	`, tenant).Scan(&unreadCount)
	if err != nil {
		h.log.Error("failed to count unread messages", zap.Error(err))
	}
	summary.UnreadMessagesCount = unreadCount

	// Get recent activity from audit logs
	activityRows, err := pool.Query(ctx, `
		SELECT
			a.id,
			COALESCE(a.actor_name, 'System') as actor_name,
			a.action,
			a.entity_type,
			a.created_at
		FROM audit_logs a
		WHERE a.tenant_id = $1
			AND a.entity_type IN ('incident', 'work_order', 'chat_session', 'message')
		ORDER BY a.created_at DESC
		LIMIT 10
	`, tenant)
	if err != nil {
		h.log.Error("failed to query recent activity", zap.Error(err))
	} else {
		defer activityRows.Close()
		for activityRows.Next() {
			var activity SupportAgentActivity
			var createdAt time.Time
			if err := activityRows.Scan(
				&activity.ID, &activity.ActorName, &activity.Action,
				&activity.EntityType, &createdAt,
			); err != nil {
				continue
			}
			// Generate description based on action and entity type
			switch activity.EntityType {
			case "incident":
				switch activity.Action {
				case "create":
					activity.Description = "Created new incident"
				case "update":
					activity.Description = "Updated incident"
				case "resolve":
					activity.Description = "Resolved incident"
				default:
					activity.Description = "Modified incident"
				}
			case "work_order":
				switch activity.Action {
				case "create":
					activity.Description = "Created work order"
				case "update":
					activity.Description = "Updated work order"
				default:
					activity.Description = "Modified work order"
				}
			case "chat_session":
				switch activity.Action {
				case "accept":
					activity.Description = "Accepted chat session"
				case "end":
					activity.Description = "Ended chat session"
				case "transfer":
					activity.Description = "Transferred chat"
				default:
					activity.Description = "Modified chat session"
				}
			case "message":
				activity.Description = "Sent message"
			default:
				activity.Description = activity.Action + " " + activity.EntityType
			}
			activity.CreatedAt = createdAt.UTC().Format(time.RFC3339)
			summary.RecentActivity = append(summary.RecentActivity, activity)
		}
	}

	writeJSON(w, http.StatusOK, summary)
}

// GetIncidentQueue returns paginated list of incidents for support agent
func (h *SupportAgentDashboardHandler) GetIncidentQueue(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())
	ctx := r.Context()
	pool := h.pg.RawPool()
	limit := parseLimit(r.URL.Query().Get("limit"), 20, 100)

	rows, err := pool.Query(ctx, `
		SELECT
			i.id,
			i.title,
			COALESCE(s.county_name, 'Unknown School') as school_name,
			i.category,
			i.severity,
			i.status,
			COALESCE(i.reported_by, '') as reported_by,
			i.sla_due_at,
			i.sla_breached,
			i.created_at
		FROM incidents i
		LEFT JOIN schools s ON s.tenant_id = i.tenant_id AND s.school_id = i.school_id
		WHERE i.tenant_id = $1
			AND i.status NOT IN ('closed', 'resolved')
		ORDER BY
			i.sla_breached DESC,
			CASE i.severity
				WHEN 'critical' THEN 1
				WHEN 'high' THEN 2
				WHEN 'medium' THEN 3
				ELSE 4
			END,
			i.sla_due_at ASC,
			i.created_at ASC
		LIMIT $2
	`, tenant, limit)
	if err != nil {
		h.log.Error("failed to query incidents", zap.Error(err))
		http.Error(w, "failed to query incidents", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	items := []IncidentQueueItem{}
	for rows.Next() {
		var item IncidentQueueItem
		var createdAt time.Time
		var slaDueAt *time.Time
		if err := rows.Scan(
			&item.ID, &item.Title, &item.SchoolName,
			&item.Category, &item.Severity, &item.Status,
			&item.ReportedBy, &slaDueAt, &item.SLABreached, &createdAt,
		); err != nil {
			continue
		}
		if slaDueAt != nil {
			t := slaDueAt.UTC().Format(time.RFC3339)
			item.SLADueAt = &t
		}
		item.CreatedAt = createdAt.UTC().Format(time.RFC3339)
		items = append(items, item)
	}

	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// GetChatQueue returns paginated list of chat sessions
func (h *SupportAgentDashboardHandler) GetChatQueue(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())
	ctx := r.Context()
	pool := h.pg.RawPool()
	limit := parseLimit(r.URL.Query().Get("limit"), 20, 100)

	rows, err := pool.Query(ctx, `
		SELECT
			cs.id,
			COALESCE(s.county_name, 'Unknown School') as school_name,
			COALESCE(cs.school_contact_name, '') as contact_name,
			cs.status,
			cs.queue_position,
			cs.assigned_agent_name,
			cs.started_at
		FROM chat_sessions cs
		LEFT JOIN schools s ON s.tenant_id = cs.tenant_id AND s.school_id = cs.school_id
		WHERE cs.tenant_id = $1
			AND cs.status IN ('waiting', 'active')
		ORDER BY
			CASE cs.status
				WHEN 'waiting' THEN 1
				WHEN 'active' THEN 2
			END,
			cs.started_at ASC
		LIMIT $2
	`, tenant, limit)
	if err != nil {
		h.log.Error("failed to query chat sessions", zap.Error(err))
		http.Error(w, "failed to query chat sessions", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	items := []ChatQueueItem{}
	now := time.Now()
	for rows.Next() {
		var item ChatQueueItem
		var startedAt time.Time
		if err := rows.Scan(
			&item.ID, &item.SchoolName, &item.ContactName,
			&item.Status, &item.QueuePosition, &item.AssignedAgentName,
			&startedAt,
		); err != nil {
			continue
		}
		item.StartedAt = startedAt.UTC().Format(time.RFC3339)
		item.WaitTimeSeconds = int(now.Sub(startedAt).Seconds())
		items = append(items, item)
	}

	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// GetIncidentMetrics returns incident statistics
func (h *SupportAgentDashboardHandler) GetIncidentMetrics(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())
	ctx := r.Context()
	pool := h.pg.RawPool()

	metrics := IncidentMetrics{}

	rows, err := pool.Query(ctx, `
		SELECT
			status,
			sla_breached,
			COUNT(*) as count
		FROM incidents
		WHERE tenant_id = $1
		GROUP BY status, sla_breached
	`, tenant)
	if err != nil {
		h.log.Error("failed to query incident metrics", zap.Error(err))
		http.Error(w, "failed to query metrics", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var status string
		var slaBreached bool
		var count int
		if err := rows.Scan(&status, &slaBreached, &count); err != nil {
			continue
		}
		if slaBreached {
			metrics.SLABreached += count
		}
		switch status {
		case "open", "new", "assigned":
			metrics.Open += count
		case "in_progress", "investigating":
			metrics.InProgress += count
		case "resolved", "closed":
			metrics.Resolved += count
		}
	}

	writeJSON(w, http.StatusOK, metrics)
}

// GetWorkOrderQueue returns active work orders for support agent
func (h *SupportAgentDashboardHandler) GetWorkOrderQueue(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())
	ctx := r.Context()
	pool := h.pg.RawPool()
	limit := parseLimit(r.URL.Query().Get("limit"), 20, 100)

	rows, err := pool.Query(ctx, `
		SELECT
			w.id,
			COALESCE(w.notes, 'Work Order ' || LEFT(w.id, 8)) as title,
			COALESCE(s.county_name, 'Unknown School') as school_name,
			w.status,
			COALESCE(w.task_type, '') as task_type,
			COALESCE(w.assigned_to, '') as assigned_to,
			w.created_at
		FROM work_orders w
		LEFT JOIN schools s ON s.tenant_id = w.tenant_id AND s.school_id = w.school_id
		WHERE w.tenant_id = $1
			AND w.status NOT IN ('completed', 'approved', 'cancelled', 'closed')
		ORDER BY
			w.created_at DESC
		LIMIT $2
	`, tenant, limit)
	if err != nil {
		h.log.Error("failed to query work orders", zap.Error(err))
		http.Error(w, "failed to query work orders", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	items := []WorkOrderQueueItem{}
	for rows.Next() {
		var item WorkOrderQueueItem
		var createdAt time.Time
		if err := rows.Scan(
			&item.ID, &item.Title, &item.SchoolName,
			&item.Status, &item.TaskType, &item.AssignedTo, &createdAt,
		); err != nil {
			continue
		}
		item.CreatedAt = createdAt.UTC().Format(time.RFC3339)
		items = append(items, item)
	}

	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}
