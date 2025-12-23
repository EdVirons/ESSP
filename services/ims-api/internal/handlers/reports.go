package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/edvirons/ssp/ims/internal/middleware"
	"github.com/edvirons/ssp/ims/internal/store"
	"go.uber.org/zap"
)

// ReportsHandler provides report data endpoints
type ReportsHandler struct {
	log *zap.Logger
	pg  *store.Postgres
}

// NewReportsHandler creates a new reports handler
func NewReportsHandler(log *zap.Logger, pg *store.Postgres) *ReportsHandler {
	return &ReportsHandler{log: log, pg: pg}
}

// ReportFilters contains common filter parameters for reports
type ReportFilters struct {
	DateFrom   *time.Time
	DateTo     *time.Time
	Status     []string
	SchoolID   string
	CountyCode string
	Category   string
	SortBy     string
	SortDir    string
	Limit      int
	Offset     int
}

func parseReportFilters(r *http.Request) ReportFilters {
	q := r.URL.Query()
	filters := ReportFilters{
		SchoolID:   q.Get("schoolId"),
		CountyCode: q.Get("countyCode"),
		Category:   q.Get("category"),
		SortBy:     q.Get("sortBy"),
		SortDir:    q.Get("sortDir"),
		Limit:      50,
		Offset:     0,
	}

	if v := q.Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 100 {
			filters.Limit = n
		}
	}
	if v := q.Get("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			filters.Offset = n
		}
	}
	if v := q.Get("dateFrom"); v != "" {
		if t, err := time.Parse("2006-01-02", v); err == nil {
			filters.DateFrom = &t
		}
	}
	if v := q.Get("dateTo"); v != "" {
		if t, err := time.Parse("2006-01-02", v); err == nil {
			endOfDay := t.Add(24*time.Hour - time.Second)
			filters.DateTo = &endOfDay
		}
	}
	if v := q["status"]; len(v) > 0 {
		filters.Status = v
	}

	if filters.SortDir != "asc" && filters.SortDir != "desc" {
		filters.SortDir = "desc"
	}

	return filters
}

// --- Work Orders Report ---

// WorkOrderReportItem represents a single work order in the report
type WorkOrderReportItem struct {
	ID             string     `json:"id"`
	IncidentID     string     `json:"incidentId"`
	Status         string     `json:"status"`
	TaskType       string     `json:"taskType"`
	SchoolName     string     `json:"schoolName"`
	DeviceCategory string     `json:"deviceCategory"`
	AssignedTo     string     `json:"assignedTo"`
	CostCents      int64      `json:"costCents"`
	ReworkCount    int        `json:"reworkCount"`
	CreatedAt      time.Time  `json:"createdAt"`
	CompletedAt    *time.Time `json:"completedAt,omitempty"`
	DurationHours  *float64   `json:"durationHours,omitempty"`
}

// WorkOrderReportSummary contains aggregated work order metrics
type WorkOrderReportSummary struct {
	Total              int            `json:"total"`
	ByStatus           map[string]int `json:"byStatus"`
	AvgCompletionHours float64        `json:"avgCompletionHours"`
	TotalCostCents     int64          `json:"totalCostCents"`
	ReworkRate         float64        `json:"reworkRate"`
}

// WorkOrderReportResponse is the full work order report response
type WorkOrderReportResponse struct {
	Items      []WorkOrderReportItem  `json:"items"`
	Summary    WorkOrderReportSummary `json:"summary"`
	Pagination struct {
		Total  int `json:"total"`
		Offset int `json:"offset"`
		Limit  int `json:"limit"`
	} `json:"pagination"`
}

// WorkOrdersReport returns work order report data
func (h *ReportsHandler) WorkOrdersReport(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())
	ctx := r.Context()
	pool := h.pg.RawPool()
	filters := parseReportFilters(r)

	response := WorkOrderReportResponse{
		Items: []WorkOrderReportItem{},
		Summary: WorkOrderReportSummary{
			ByStatus: make(map[string]int),
		},
	}

	// Build WHERE clause
	whereClause := "WHERE w.tenant_id = $1"
	args := []any{tenant}
	argIdx := 2

	if filters.DateFrom != nil {
		whereClause += " AND w.created_at >= $" + strconv.Itoa(argIdx)
		args = append(args, *filters.DateFrom)
		argIdx++
	}
	if filters.DateTo != nil {
		whereClause += " AND w.created_at <= $" + strconv.Itoa(argIdx)
		args = append(args, *filters.DateTo)
		argIdx++
	}
	if len(filters.Status) > 0 {
		whereClause += " AND w.status = ANY($" + strconv.Itoa(argIdx) + ")"
		args = append(args, filters.Status)
		argIdx++
	}
	if filters.SchoolID != "" {
		whereClause += " AND w.school_id = $" + strconv.Itoa(argIdx)
		args = append(args, filters.SchoolID)
		argIdx++
	}

	// Get summary
	summaryQuery := `
		SELECT
			COUNT(*) as total,
			COALESCE(SUM(CASE WHEN status = 'draft' THEN 1 ELSE 0 END), 0) as draft,
			COALESCE(SUM(CASE WHEN status = 'assigned' THEN 1 ELSE 0 END), 0) as assigned,
			COALESCE(SUM(CASE WHEN status = 'in_repair' THEN 1 ELSE 0 END), 0) as in_repair,
			COALESCE(SUM(CASE WHEN status = 'qa' THEN 1 ELSE 0 END), 0) as qa,
			COALESCE(SUM(CASE WHEN status = 'completed' THEN 1 ELSE 0 END), 0) as completed,
			COALESCE(SUM(CASE WHEN status = 'approved' THEN 1 ELSE 0 END), 0) as approved,
			COALESCE(SUM(cost_estimate_cents), 0) as total_cost,
			COALESCE(SUM(CASE WHEN rework_count > 0 THEN 1 ELSE 0 END), 0) as reworked_count,
			COALESCE(AVG(EXTRACT(EPOCH FROM (updated_at - created_at))/3600) FILTER (WHERE status IN ('completed', 'approved')), 0) as avg_completion_hours
		FROM work_orders w
		` + whereClause

	var draft, assigned, inRepair, qa, completed, approved, reworkedCount int
	err := pool.QueryRow(ctx, summaryQuery, args...).Scan(
		&response.Summary.Total,
		&draft, &assigned, &inRepair, &qa, &completed, &approved,
		&response.Summary.TotalCostCents,
		&reworkedCount,
		&response.Summary.AvgCompletionHours,
	)
	if err != nil {
		h.log.Error("failed to query work order summary", zap.Error(err))
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to generate report"})
		return
	}

	response.Summary.ByStatus = map[string]int{
		"draft":     draft,
		"assigned":  assigned,
		"in_repair": inRepair,
		"qa":        qa,
		"completed": completed,
		"approved":  approved,
	}

	if response.Summary.Total > 0 {
		response.Summary.ReworkRate = float64(reworkedCount) / float64(response.Summary.Total) * 100
	}

	// Get items with pagination
	orderBy := "w.created_at"
	if filters.SortBy == "status" || filters.SortBy == "schoolName" || filters.SortBy == "costCents" {
		orderBy = "w." + filters.SortBy
	}

	itemsQuery := `
		SELECT
			w.id, COALESCE(w.incident_id, ''), w.status, w.task_type,
			COALESCE(w.school_name, ''), COALESCE(w.device_category, ''),
			COALESCE(w.assigned_staff_id, ''),
			COALESCE(w.cost_estimate_cents, 0),
			COALESCE(w.rework_count, 0),
			w.created_at,
			CASE WHEN w.status IN ('completed', 'approved') THEN w.updated_at ELSE NULL END as completed_at
		FROM work_orders w
		` + whereClause + `
		ORDER BY ` + orderBy + ` ` + filters.SortDir + `
		LIMIT $` + strconv.Itoa(argIdx) + ` OFFSET $` + strconv.Itoa(argIdx+1)

	args = append(args, filters.Limit, filters.Offset)

	rows, err := pool.Query(ctx, itemsQuery, args...)
	if err != nil {
		h.log.Error("failed to query work order items", zap.Error(err))
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to generate report"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var item WorkOrderReportItem
		var completedAt *time.Time
		err := rows.Scan(
			&item.ID, &item.IncidentID, &item.Status, &item.TaskType,
			&item.SchoolName, &item.DeviceCategory, &item.AssignedTo,
			&item.CostCents, &item.ReworkCount, &item.CreatedAt, &completedAt,
		)
		if err != nil {
			h.log.Error("failed to scan work order item", zap.Error(err))
			continue
		}
		item.CompletedAt = completedAt
		if completedAt != nil {
			hours := completedAt.Sub(item.CreatedAt).Hours()
			item.DurationHours = &hours
		}
		response.Items = append(response.Items, item)
	}

	response.Pagination.Total = response.Summary.Total
	response.Pagination.Offset = filters.Offset
	response.Pagination.Limit = filters.Limit

	writeJSON(w, http.StatusOK, response)
}

// --- Incidents Report ---

// IncidentReportItem represents a single incident in the report
type IncidentReportItem struct {
	ID              string     `json:"id"`
	Title           string     `json:"title"`
	Status          string     `json:"status"`
	Severity        string     `json:"severity"`
	Category        string     `json:"category"`
	SchoolName      string     `json:"schoolName"`
	SLABreached     bool       `json:"slaBreached"`
	CreatedAt       time.Time  `json:"createdAt"`
	ResolvedAt      *time.Time `json:"resolvedAt,omitempty"`
	ResolutionHours *float64   `json:"resolutionHours,omitempty"`
}

// IncidentReportSummary contains aggregated incident metrics
type IncidentReportSummary struct {
	Total              int            `json:"total"`
	ByStatus           map[string]int `json:"byStatus"`
	BySeverity         map[string]int `json:"bySeverity"`
	SLABreachedCount   int            `json:"slaBreachedCount"`
	SLAComplianceRate  float64        `json:"slaComplianceRate"`
	AvgResolutionHours float64        `json:"avgResolutionHours"`
}

// IncidentReportResponse is the full incident report response
type IncidentReportResponse struct {
	Items      []IncidentReportItem  `json:"items"`
	Summary    IncidentReportSummary `json:"summary"`
	Pagination struct {
		Total  int `json:"total"`
		Offset int `json:"offset"`
		Limit  int `json:"limit"`
	} `json:"pagination"`
}

// IncidentsReport returns incident report data
func (h *ReportsHandler) IncidentsReport(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())
	ctx := r.Context()
	pool := h.pg.RawPool()
	filters := parseReportFilters(r)

	response := IncidentReportResponse{
		Items: []IncidentReportItem{},
		Summary: IncidentReportSummary{
			ByStatus:   make(map[string]int),
			BySeverity: make(map[string]int),
		},
	}

	// Build WHERE clause
	whereClause := "WHERE i.tenant_id = $1"
	args := []any{tenant}
	argIdx := 2

	if filters.DateFrom != nil {
		whereClause += " AND i.created_at >= $" + strconv.Itoa(argIdx)
		args = append(args, *filters.DateFrom)
		argIdx++
	}
	if filters.DateTo != nil {
		whereClause += " AND i.created_at <= $" + strconv.Itoa(argIdx)
		args = append(args, *filters.DateTo)
		argIdx++
	}
	if len(filters.Status) > 0 {
		whereClause += " AND i.status = ANY($" + strconv.Itoa(argIdx) + ")"
		args = append(args, filters.Status)
		argIdx++
	}
	if filters.SchoolID != "" {
		whereClause += " AND i.school_id = $" + strconv.Itoa(argIdx)
		args = append(args, filters.SchoolID)
		argIdx++
	}

	// Get summary
	summaryQuery := `
		SELECT
			COUNT(*) as total,
			COALESCE(SUM(CASE WHEN status = 'new' THEN 1 ELSE 0 END), 0) as new_count,
			COALESCE(SUM(CASE WHEN status = 'acknowledged' THEN 1 ELSE 0 END), 0) as acknowledged,
			COALESCE(SUM(CASE WHEN status = 'in_progress' THEN 1 ELSE 0 END), 0) as in_progress,
			COALESCE(SUM(CASE WHEN status = 'escalated' THEN 1 ELSE 0 END), 0) as escalated,
			COALESCE(SUM(CASE WHEN status = 'resolved' THEN 1 ELSE 0 END), 0) as resolved,
			COALESCE(SUM(CASE WHEN status = 'closed' THEN 1 ELSE 0 END), 0) as closed,
			COALESCE(SUM(CASE WHEN severity = 'low' THEN 1 ELSE 0 END), 0) as sev_low,
			COALESCE(SUM(CASE WHEN severity = 'medium' THEN 1 ELSE 0 END), 0) as sev_medium,
			COALESCE(SUM(CASE WHEN severity = 'high' THEN 1 ELSE 0 END), 0) as sev_high,
			COALESCE(SUM(CASE WHEN severity = 'critical' THEN 1 ELSE 0 END), 0) as sev_critical,
			COALESCE(SUM(CASE WHEN sla_breached THEN 1 ELSE 0 END), 0) as sla_breached,
			COALESCE(AVG(EXTRACT(EPOCH FROM (updated_at - created_at))/3600) FILTER (WHERE status IN ('resolved', 'closed')), 0) as avg_resolution_hours
		FROM incidents i
		` + whereClause

	var newCount, acknowledged, inProgress, escalated, resolved, closed int
	var sevLow, sevMedium, sevHigh, sevCritical int
	err := pool.QueryRow(ctx, summaryQuery, args...).Scan(
		&response.Summary.Total,
		&newCount, &acknowledged, &inProgress, &escalated, &resolved, &closed,
		&sevLow, &sevMedium, &sevHigh, &sevCritical,
		&response.Summary.SLABreachedCount,
		&response.Summary.AvgResolutionHours,
	)
	if err != nil {
		h.log.Error("failed to query incident summary", zap.Error(err))
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to generate report"})
		return
	}

	response.Summary.ByStatus = map[string]int{
		"new":          newCount,
		"acknowledged": acknowledged,
		"in_progress":  inProgress,
		"escalated":    escalated,
		"resolved":     resolved,
		"closed":       closed,
	}
	response.Summary.BySeverity = map[string]int{
		"low":      sevLow,
		"medium":   sevMedium,
		"high":     sevHigh,
		"critical": sevCritical,
	}

	if response.Summary.Total > 0 {
		response.Summary.SLAComplianceRate = float64(response.Summary.Total-response.Summary.SLABreachedCount) / float64(response.Summary.Total) * 100
	}

	// Get items
	orderBy := "i.created_at"
	itemsQuery := `
		SELECT
			i.id, i.title, i.status, i.severity,
			COALESCE(i.category, ''), COALESCE(i.school_name, ''),
			COALESCE(i.sla_breached, false),
			i.created_at,
			CASE WHEN i.status IN ('resolved', 'closed') THEN i.updated_at ELSE NULL END as resolved_at
		FROM incidents i
		` + whereClause + `
		ORDER BY ` + orderBy + ` ` + filters.SortDir + `
		LIMIT $` + strconv.Itoa(argIdx) + ` OFFSET $` + strconv.Itoa(argIdx+1)

	args = append(args, filters.Limit, filters.Offset)

	rows, err := pool.Query(ctx, itemsQuery, args...)
	if err != nil {
		h.log.Error("failed to query incident items", zap.Error(err))
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to generate report"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var item IncidentReportItem
		var resolvedAt *time.Time
		err := rows.Scan(
			&item.ID, &item.Title, &item.Status, &item.Severity,
			&item.Category, &item.SchoolName, &item.SLABreached,
			&item.CreatedAt, &resolvedAt,
		)
		if err != nil {
			h.log.Error("failed to scan incident item", zap.Error(err))
			continue
		}
		item.ResolvedAt = resolvedAt
		if resolvedAt != nil {
			hours := resolvedAt.Sub(item.CreatedAt).Hours()
			item.ResolutionHours = &hours
		}
		response.Items = append(response.Items, item)
	}

	response.Pagination.Total = response.Summary.Total
	response.Pagination.Offset = filters.Offset
	response.Pagination.Limit = filters.Limit

	writeJSON(w, http.StatusOK, response)
}

// --- Inventory Report ---

// InventoryReportItem represents a single inventory item in the report
type InventoryReportItem struct {
	PartID           string `json:"partId"`
	PartSKU          string `json:"partSku"`
	PartName         string `json:"partName"`
	Category         string `json:"category"`
	ServiceShopName  string `json:"serviceShopName"`
	QtyAvailable     int64  `json:"qtyAvailable"`
	QtyReserved      int64  `json:"qtyReserved"`
	ReorderThreshold int64  `json:"reorderThreshold"`
	IsLowStock       bool   `json:"isLowStock"`
}

// InventoryReportSummary contains aggregated inventory metrics
type InventoryReportSummary struct {
	TotalParts        int            `json:"totalParts"`
	LowStockCount     int            `json:"lowStockCount"`
	TotalQtyAvailable int64          `json:"totalQtyAvailable"`
	ByCategory        map[string]int `json:"byCategory"`
}

// InventoryReportResponse is the full inventory report response
type InventoryReportResponse struct {
	Items      []InventoryReportItem  `json:"items"`
	Summary    InventoryReportSummary `json:"summary"`
	Pagination struct {
		Total  int `json:"total"`
		Offset int `json:"offset"`
		Limit  int `json:"limit"`
	} `json:"pagination"`
}

// InventoryReport returns inventory report data
func (h *ReportsHandler) InventoryReport(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())
	ctx := r.Context()
	pool := h.pg.RawPool()
	filters := parseReportFilters(r)

	response := InventoryReportResponse{
		Items: []InventoryReportItem{},
		Summary: InventoryReportSummary{
			ByCategory: make(map[string]int),
		},
	}

	// Build WHERE clause
	whereClause := "WHERE i.tenant_id = $1"
	args := []any{tenant}
	argIdx := 2

	if filters.Category != "" {
		whereClause += " AND p.category = $" + strconv.Itoa(argIdx)
		args = append(args, filters.Category)
		argIdx++
	}

	// Get summary
	summaryQuery := `
		SELECT
			COUNT(*) as total,
			COALESCE(SUM(CASE WHEN i.qty_available <= i.reorder_threshold AND i.reorder_threshold > 0 THEN 1 ELSE 0 END), 0) as low_stock,
			COALESCE(SUM(i.qty_available), 0) as total_qty
		FROM inventory i
		JOIN parts p ON p.id = i.part_id
		` + whereClause

	err := pool.QueryRow(ctx, summaryQuery, args...).Scan(
		&response.Summary.TotalParts,
		&response.Summary.LowStockCount,
		&response.Summary.TotalQtyAvailable,
	)
	if err != nil {
		h.log.Error("failed to query inventory summary", zap.Error(err))
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to generate report"})
		return
	}

	// Get category breakdown
	categoryQuery := `
		SELECT p.category, COUNT(*) as cnt
		FROM inventory i
		JOIN parts p ON p.id = i.part_id
		` + whereClause + `
		GROUP BY p.category
	`
	catRows, err := pool.Query(ctx, categoryQuery, args[:argIdx-1]...)
	if err == nil {
		defer catRows.Close()
		for catRows.Next() {
			var cat string
			var cnt int
			if err := catRows.Scan(&cat, &cnt); err == nil {
				response.Summary.ByCategory[cat] = cnt
			}
		}
	}

	// Get items
	itemsQuery := `
		SELECT
			p.id, p.sku, p.name, p.category,
			COALESCE(s.name, 'Unknown') as shop_name,
			i.qty_available, i.qty_reserved, i.reorder_threshold,
			(i.qty_available <= i.reorder_threshold AND i.reorder_threshold > 0) as is_low_stock
		FROM inventory i
		JOIN parts p ON p.id = i.part_id
		LEFT JOIN service_shops s ON s.id = i.service_shop_id
		` + whereClause + `
		ORDER BY p.name
		LIMIT $` + strconv.Itoa(argIdx) + ` OFFSET $` + strconv.Itoa(argIdx+1)

	args = append(args, filters.Limit, filters.Offset)

	rows, err := pool.Query(ctx, itemsQuery, args...)
	if err != nil {
		h.log.Error("failed to query inventory items", zap.Error(err))
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to generate report"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var item InventoryReportItem
		err := rows.Scan(
			&item.PartID, &item.PartSKU, &item.PartName, &item.Category,
			&item.ServiceShopName,
			&item.QtyAvailable, &item.QtyReserved, &item.ReorderThreshold,
			&item.IsLowStock,
		)
		if err != nil {
			h.log.Error("failed to scan inventory item", zap.Error(err))
			continue
		}
		response.Items = append(response.Items, item)
	}

	response.Pagination.Total = response.Summary.TotalParts
	response.Pagination.Offset = filters.Offset
	response.Pagination.Limit = filters.Limit

	writeJSON(w, http.StatusOK, response)
}

// --- Schools Report ---

// SchoolReportItem represents a single school in the report
type SchoolReportItem struct {
	SchoolID       string `json:"schoolId"`
	SchoolName     string `json:"schoolName"`
	CountyName     string `json:"countyName"`
	DeviceCount    int    `json:"deviceCount"`
	IncidentCount  int    `json:"incidentCount"`
	WorkOrderCount int    `json:"workOrderCount"`
}

// SchoolReportSummary contains aggregated school metrics
type SchoolReportSummary struct {
	TotalSchools int            `json:"totalSchools"`
	TotalDevices int            `json:"totalDevices"`
	ByCounty     map[string]int `json:"byCounty"`
}

// SchoolReportResponse is the full schools report response
type SchoolReportResponse struct {
	Items      []SchoolReportItem  `json:"items"`
	Summary    SchoolReportSummary `json:"summary"`
	Pagination struct {
		Total  int `json:"total"`
		Offset int `json:"offset"`
		Limit  int `json:"limit"`
	} `json:"pagination"`
}

// SchoolsReport returns schools report data
func (h *ReportsHandler) SchoolsReport(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())
	ctx := r.Context()
	pool := h.pg.RawPool()
	filters := parseReportFilters(r)

	response := SchoolReportResponse{
		Items: []SchoolReportItem{},
		Summary: SchoolReportSummary{
			ByCounty: make(map[string]int),
		},
	}

	// Build WHERE clause
	whereClause := "WHERE s.tenant_id = $1"
	args := []any{tenant}
	argIdx := 2

	if filters.CountyCode != "" {
		whereClause += " AND s.county_code = $" + strconv.Itoa(argIdx)
		args = append(args, filters.CountyCode)
		argIdx++
	}

	// Get summary
	summaryQuery := `
		SELECT COUNT(DISTINCT s.school_id) as total_schools
		FROM schools_snapshot s
		` + whereClause

	err := pool.QueryRow(ctx, summaryQuery, args...).Scan(&response.Summary.TotalSchools)
	if err != nil {
		h.log.Error("failed to query schools summary", zap.Error(err))
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to generate report"})
		return
	}

	// Get device count
	deviceQuery := `SELECT COUNT(*) FROM devices_snapshot WHERE tenant_id = $1`
	_ = pool.QueryRow(ctx, deviceQuery, tenant).Scan(&response.Summary.TotalDevices)

	// Get county breakdown
	countyQuery := `
		SELECT COALESCE(s.county_name, 'Unknown'), COUNT(*) as cnt
		FROM schools_snapshot s
		` + whereClause + `
		GROUP BY s.county_name
	`
	countyRows, err := pool.Query(ctx, countyQuery, args[:argIdx-1]...)
	if err == nil {
		defer countyRows.Close()
		for countyRows.Next() {
			var county string
			var cnt int
			if err := countyRows.Scan(&county, &cnt); err == nil {
				response.Summary.ByCounty[county] = cnt
			}
		}
	}

	// Get items with incident/work order counts
	itemsQuery := `
		SELECT
			s.school_id,
			s.school_name,
			COALESCE(s.county_name, 'Unknown') as county_name,
			(SELECT COUNT(*) FROM devices_snapshot d WHERE d.school_id = s.school_id AND d.tenant_id = s.tenant_id) as device_count,
			(SELECT COUNT(*) FROM incidents i WHERE i.school_id = s.school_id AND i.tenant_id = s.tenant_id) as incident_count,
			(SELECT COUNT(*) FROM work_orders w WHERE w.school_id = s.school_id AND w.tenant_id = s.tenant_id) as work_order_count
		FROM schools_snapshot s
		` + whereClause + `
		ORDER BY s.school_name
		LIMIT $` + strconv.Itoa(argIdx) + ` OFFSET $` + strconv.Itoa(argIdx+1)

	args = append(args, filters.Limit, filters.Offset)

	rows, err := pool.Query(ctx, itemsQuery, args...)
	if err != nil {
		h.log.Error("failed to query school items", zap.Error(err))
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to generate report"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var item SchoolReportItem
		err := rows.Scan(
			&item.SchoolID, &item.SchoolName, &item.CountyName,
			&item.DeviceCount, &item.IncidentCount, &item.WorkOrderCount,
		)
		if err != nil {
			h.log.Error("failed to scan school item", zap.Error(err))
			continue
		}
		response.Items = append(response.Items, item)
	}

	response.Pagination.Total = response.Summary.TotalSchools
	response.Pagination.Offset = filters.Offset
	response.Pagination.Limit = filters.Limit

	writeJSON(w, http.StatusOK, response)
}

// --- Executive Dashboard ---

// ExecutiveDashboardResponse contains high-level KPIs
type ExecutiveDashboardResponse struct {
	WorkOrders struct {
		Total             int     `json:"total"`
		Completed         int     `json:"completed"`
		InProgress        int     `json:"inProgress"`
		CompletionRate    float64 `json:"completionRate"`
		AvgCompletionDays float64 `json:"avgCompletionDays"`
	} `json:"workOrders"`

	Incidents struct {
		Total         int     `json:"total"`
		Open          int     `json:"open"`
		Resolved      int     `json:"resolved"`
		SLACompliance float64 `json:"slaCompliance"`
		Critical      int     `json:"critical"`
	} `json:"incidents"`

	Inventory struct {
		TotalParts int `json:"totalParts"`
		LowStock   int `json:"lowStock"`
		OutOfStock int `json:"outOfStock"`
	} `json:"inventory"`

	Schools struct {
		TotalSchools   int `json:"totalSchools"`
		TotalDevices   int `json:"totalDevices"`
		ActiveProjects int `json:"activeProjects"`
	} `json:"schools"`
}

// ExecutiveDashboard returns executive-level KPIs
func (h *ReportsHandler) ExecutiveDashboard(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())
	ctx := r.Context()
	pool := h.pg.RawPool()

	response := ExecutiveDashboardResponse{}

	// Work Orders KPIs
	woQuery := `
		SELECT
			COUNT(*) as total,
			COALESCE(SUM(CASE WHEN status IN ('completed', 'approved') THEN 1 ELSE 0 END), 0) as completed,
			COALESCE(SUM(CASE WHEN status = 'in_repair' THEN 1 ELSE 0 END), 0) as in_progress,
			COALESCE(AVG(EXTRACT(EPOCH FROM (updated_at - created_at))/86400) FILTER (WHERE status IN ('completed', 'approved')), 0) as avg_days
		FROM work_orders
		WHERE tenant_id = $1
	`
	_ = pool.QueryRow(ctx, woQuery, tenant).Scan(
		&response.WorkOrders.Total,
		&response.WorkOrders.Completed,
		&response.WorkOrders.InProgress,
		&response.WorkOrders.AvgCompletionDays,
	)
	if response.WorkOrders.Total > 0 {
		response.WorkOrders.CompletionRate = float64(response.WorkOrders.Completed) / float64(response.WorkOrders.Total) * 100
	}

	// Incidents KPIs
	incQuery := `
		SELECT
			COUNT(*) as total,
			COALESCE(SUM(CASE WHEN status NOT IN ('resolved', 'closed') THEN 1 ELSE 0 END), 0) as open,
			COALESCE(SUM(CASE WHEN status IN ('resolved', 'closed') THEN 1 ELSE 0 END), 0) as resolved,
			COALESCE(SUM(CASE WHEN severity = 'critical' THEN 1 ELSE 0 END), 0) as critical,
			COALESCE(SUM(CASE WHEN sla_breached THEN 1 ELSE 0 END), 0) as sla_breached
		FROM incidents
		WHERE tenant_id = $1
	`
	var slaBreached int
	_ = pool.QueryRow(ctx, incQuery, tenant).Scan(
		&response.Incidents.Total,
		&response.Incidents.Open,
		&response.Incidents.Resolved,
		&response.Incidents.Critical,
		&slaBreached,
	)
	if response.Incidents.Total > 0 {
		response.Incidents.SLACompliance = float64(response.Incidents.Total-slaBreached) / float64(response.Incidents.Total) * 100
	}

	// Inventory KPIs
	invQuery := `
		SELECT
			COUNT(*) as total,
			COALESCE(SUM(CASE WHEN qty_available <= reorder_threshold AND reorder_threshold > 0 THEN 1 ELSE 0 END), 0) as low_stock,
			COALESCE(SUM(CASE WHEN qty_available = 0 THEN 1 ELSE 0 END), 0) as out_of_stock
		FROM inventory
		WHERE tenant_id = $1
	`
	_ = pool.QueryRow(ctx, invQuery, tenant).Scan(
		&response.Inventory.TotalParts,
		&response.Inventory.LowStock,
		&response.Inventory.OutOfStock,
	)

	// Schools KPIs
	schoolQuery := `SELECT COUNT(DISTINCT school_id) FROM schools_snapshot WHERE tenant_id = $1`
	_ = pool.QueryRow(ctx, schoolQuery, tenant).Scan(&response.Schools.TotalSchools)

	deviceQuery := `SELECT COUNT(*) FROM devices_snapshot WHERE tenant_id = $1`
	_ = pool.QueryRow(ctx, deviceQuery, tenant).Scan(&response.Schools.TotalDevices)

	projectQuery := `SELECT COUNT(*) FROM school_service_projects WHERE tenant_id = $1 AND status NOT IN ('completed', 'cancelled')`
	_ = pool.QueryRow(ctx, projectQuery, tenant).Scan(&response.Schools.ActiveProjects)

	writeJSON(w, http.StatusOK, response)
}
