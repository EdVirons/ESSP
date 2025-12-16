package handlers

import (
	"net/http"
	"time"

	"github.com/edvirons/ssp/ims/internal/middleware"
	"github.com/edvirons/ssp/ims/internal/store"
	"go.uber.org/zap"
)

// WarehouseDashboardHandler provides aggregated data for the warehouse manager dashboard
type WarehouseDashboardHandler struct {
	log *zap.Logger
	pg  *store.Postgres
}

// NewWarehouseDashboardHandler creates a new warehouse dashboard handler
func NewWarehouseDashboardHandler(log *zap.Logger, pg *store.Postgres) *WarehouseDashboardHandler {
	return &WarehouseDashboardHandler{log: log, pg: pg}
}

// LowStockItem represents an inventory item below reorder threshold
type LowStockItem struct {
	ID               string `json:"id"`
	PartID           string `json:"partId"`
	PartName         string `json:"partName"`
	ShopID           string `json:"shopId"`
	ShopName         string `json:"shopName"`
	QtyAvailable     int64  `json:"qtyAvailable"`
	ReorderThreshold int64  `json:"reorderThreshold"`
}

// PendingPartIssue represents a work order needing parts
type PendingPartIssue struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	SchoolName  string `json:"schoolName"`
	Priority    string `json:"priority"`
	PartsNeeded int    `json:"partsNeeded"`
	CreatedAt   string `json:"createdAt"`
}

// InventoryActivity represents a stock movement event
type InventoryActivity struct {
	ID          string `json:"id"`
	Type        string `json:"type"` // receipt, issue, adjustment, transfer
	Description string `json:"description"`
	PartName    string `json:"partName"`
	ActorName   string `json:"actorName"`
	QtyChange   int64  `json:"qtyChange"`
	CreatedAt   string `json:"createdAt"`
}

// WarehouseDashboardSummary contains all dashboard data
type WarehouseDashboardSummary struct {
	LowStockCount     int                `json:"lowStockCount"`
	PendingWorkOrders int                `json:"pendingWorkOrders"`
	TodayMovements    int                `json:"todayMovements"`
	TotalParts        int                `json:"totalParts"`
	PartsCategories   map[string]int     `json:"partsCategories"`
	RecentActivity    []InventoryActivity `json:"recentActivity"`
	LowStockItems     []LowStockItem      `json:"lowStockItems"`
	PendingPartIssues []PendingPartIssue  `json:"pendingPartIssues"`
}

// GetDashboardSummary returns aggregated dashboard data for warehouse manager
func (h *WarehouseDashboardHandler) GetDashboardSummary(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())
	ctx := r.Context()
	pool := h.pg.RawPool()

	summary := WarehouseDashboardSummary{
		PartsCategories:   make(map[string]int),
		RecentActivity:    []InventoryActivity{},
		LowStockItems:     []LowStockItem{},
		PendingPartIssues: []PendingPartIssue{},
	}

	// Get low stock items count and list
	lowStockRows, err := pool.Query(ctx, `
		SELECT
			i.id, i.part_id, p.name as part_name,
			i.service_shop_id, s.name as shop_name,
			i.qty_available, i.reorder_threshold
		FROM inventory i
		JOIN parts p ON p.id = i.part_id
		JOIN service_shops s ON s.id = i.service_shop_id
		WHERE i.tenant_id = $1
			AND i.qty_available <= i.reorder_threshold
			AND i.reorder_threshold > 0
		ORDER BY (i.qty_available::float / NULLIF(i.reorder_threshold, 0)) ASC
		LIMIT 10
	`, tenant)
	if err != nil {
		h.log.Error("failed to query low stock items", zap.Error(err))
	} else {
		defer lowStockRows.Close()
		for lowStockRows.Next() {
			var item LowStockItem
			if err := lowStockRows.Scan(
				&item.ID, &item.PartID, &item.PartName,
				&item.ShopID, &item.ShopName,
				&item.QtyAvailable, &item.ReorderThreshold,
			); err != nil {
				continue
			}
			summary.LowStockItems = append(summary.LowStockItems, item)
		}
		summary.LowStockCount = len(summary.LowStockItems)
	}

	// Get pending work orders needing parts (status in_repair or assigned with BOM items)
	pendingRows, err := pool.Query(ctx, `
		SELECT
			w.id,
			COALESCE(w.notes, 'Work Order ' || LEFT(w.id, 8)) as title,
			COALESCE(sc.name, 'Unknown School') as school_name,
			COALESCE(w.priority, 'medium') as priority,
			COUNT(DISTINCT b.part_id) as parts_needed,
			w.created_at
		FROM work_orders w
		LEFT JOIN schools sc ON sc.id = w.school_id
		LEFT JOIN bom_items b ON b.work_order_id = w.id AND b.consumed_qty < b.qty
		WHERE w.tenant_id = $1
			AND w.status IN ('assigned', 'in_repair')
		GROUP BY w.id, w.notes, sc.name, w.priority, w.created_at
		HAVING COUNT(DISTINCT b.part_id) > 0
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
		h.log.Error("failed to query pending work orders", zap.Error(err))
	} else {
		defer pendingRows.Close()
		for pendingRows.Next() {
			var wo PendingPartIssue
			var createdAt time.Time
			if err := pendingRows.Scan(
				&wo.ID, &wo.Title, &wo.SchoolName,
				&wo.Priority, &wo.PartsNeeded, &createdAt,
			); err != nil {
				continue
			}
			wo.CreatedAt = createdAt.UTC().Format(time.RFC3339)
			summary.PendingPartIssues = append(summary.PendingPartIssues, wo)
		}
		summary.PendingWorkOrders = len(summary.PendingPartIssues)
	}

	// Get today's inventory movements from audit logs
	today := time.Now().UTC().Truncate(24 * time.Hour)
	var todayMovements int
	err = pool.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM audit_logs
		WHERE tenant_id = $1
			AND entity_type = 'inventory'
			AND created_at >= $2
	`, tenant, today).Scan(&todayMovements)
	if err != nil {
		h.log.Error("failed to count today movements", zap.Error(err))
	}
	summary.TodayMovements = todayMovements

	// Get parts by category
	categoryRows, err := pool.Query(ctx, `
		SELECT COALESCE(category, 'uncategorized'), COUNT(*)
		FROM parts
		WHERE tenant_id = $1
		GROUP BY category
	`, tenant)
	if err != nil {
		h.log.Error("failed to query parts categories", zap.Error(err))
	} else {
		defer categoryRows.Close()
		for categoryRows.Next() {
			var category string
			var count int
			if err := categoryRows.Scan(&category, &count); err != nil {
				continue
			}
			summary.PartsCategories[category] = count
			summary.TotalParts += count
		}
	}

	// Get recent activity from audit logs
	activityRows, err := pool.Query(ctx, `
		SELECT
			a.id,
			a.action,
			COALESCE(a.actor_name, 'System') as actor_name,
			a.entity_id,
			a.created_at
		FROM audit_logs a
		WHERE a.tenant_id = $1
			AND a.entity_type IN ('inventory', 'parts', 'bom_item')
		ORDER BY a.created_at DESC
		LIMIT 10
	`, tenant)
	if err != nil {
		h.log.Error("failed to query recent activity", zap.Error(err))
	} else {
		defer activityRows.Close()
		for activityRows.Next() {
			var activity InventoryActivity
			var createdAt time.Time
			var entityID string
			if err := activityRows.Scan(
				&activity.ID, &activity.Type, &activity.ActorName,
				&entityID, &createdAt,
			); err != nil {
				continue
			}
			// Map action to type
			switch activity.Type {
			case "create":
				activity.Type = "receipt"
				activity.Description = "Added to inventory"
			case "update":
				activity.Type = "adjustment"
				activity.Description = "Stock adjusted"
			case "delete":
				activity.Type = "issue"
				activity.Description = "Removed from inventory"
			default:
				activity.Description = "Inventory updated"
			}
			activity.PartName = entityID[:12] + "..."
			activity.CreatedAt = createdAt.UTC().Format(time.RFC3339)
			summary.RecentActivity = append(summary.RecentActivity, activity)
		}
	}

	writeJSON(w, http.StatusOK, summary)
}

// GetLowStockItems returns paginated list of low stock items
func (h *WarehouseDashboardHandler) GetLowStockItems(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())
	ctx := r.Context()
	pool := h.pg.RawPool()
	limit := parseLimit(r.URL.Query().Get("limit"), 20, 100)

	rows, err := pool.Query(ctx, `
		SELECT
			i.id, i.part_id, p.name as part_name,
			i.service_shop_id, s.name as shop_name,
			i.qty_available, i.reorder_threshold
		FROM inventory i
		JOIN parts p ON p.id = i.part_id
		JOIN service_shops s ON s.id = i.service_shop_id
		WHERE i.tenant_id = $1
			AND i.qty_available <= i.reorder_threshold
			AND i.reorder_threshold > 0
		ORDER BY (i.qty_available::float / NULLIF(i.reorder_threshold, 0)) ASC
		LIMIT $2
	`, tenant, limit)
	if err != nil {
		h.log.Error("failed to query low stock items", zap.Error(err))
		http.Error(w, "failed to query low stock items", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	items := []LowStockItem{}
	for rows.Next() {
		var item LowStockItem
		if err := rows.Scan(
			&item.ID, &item.PartID, &item.PartName,
			&item.ShopID, &item.ShopName,
			&item.QtyAvailable, &item.ReorderThreshold,
		); err != nil {
			continue
		}
		items = append(items, item)
	}

	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// GetPendingPartIssues returns work orders that need parts issued
func (h *WarehouseDashboardHandler) GetPendingPartIssues(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())
	ctx := r.Context()
	pool := h.pg.RawPool()
	limit := parseLimit(r.URL.Query().Get("limit"), 20, 100)

	rows, err := pool.Query(ctx, `
		SELECT
			w.id,
			COALESCE(w.notes, 'Work Order ' || LEFT(w.id, 8)) as title,
			COALESCE(sc.name, 'Unknown School') as school_name,
			COALESCE(w.priority, 'medium') as priority,
			COUNT(DISTINCT b.part_id) as parts_needed,
			w.created_at
		FROM work_orders w
		LEFT JOIN schools sc ON sc.id = w.school_id
		LEFT JOIN bom_items b ON b.work_order_id = w.id AND b.consumed_qty < b.qty
		WHERE w.tenant_id = $1
			AND w.status IN ('assigned', 'in_repair')
		GROUP BY w.id, w.notes, sc.name, w.priority, w.created_at
		HAVING COUNT(DISTINCT b.part_id) > 0
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
		h.log.Error("failed to query pending work orders", zap.Error(err))
		http.Error(w, "failed to query pending work orders", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	items := []PendingPartIssue{}
	for rows.Next() {
		var wo PendingPartIssue
		var createdAt time.Time
		if err := rows.Scan(
			&wo.ID, &wo.Title, &wo.SchoolName,
			&wo.Priority, &wo.PartsNeeded, &createdAt,
		); err != nil {
			continue
		}
		wo.CreatedAt = createdAt.UTC().Format(time.RFC3339)
		items = append(items, wo)
	}

	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// GetStockMovements returns recent stock movement history
func (h *WarehouseDashboardHandler) GetStockMovements(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())
	ctx := r.Context()
	pool := h.pg.RawPool()
	limit := parseLimit(r.URL.Query().Get("limit"), 20, 100)

	rows, err := pool.Query(ctx, `
		SELECT
			a.id,
			a.action,
			COALESCE(a.actor_name, 'System') as actor_name,
			a.entity_id,
			a.created_at
		FROM audit_logs a
		WHERE a.tenant_id = $1
			AND a.entity_type IN ('inventory', 'parts', 'bom_item')
		ORDER BY a.created_at DESC
		LIMIT $2
	`, tenant, limit)
	if err != nil {
		h.log.Error("failed to query stock movements", zap.Error(err))
		http.Error(w, "failed to query stock movements", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	items := []InventoryActivity{}
	for rows.Next() {
		var activity InventoryActivity
		var createdAt time.Time
		var entityID string
		if err := rows.Scan(
			&activity.ID, &activity.Type, &activity.ActorName,
			&entityID, &createdAt,
		); err != nil {
			continue
		}
		switch activity.Type {
		case "create":
			activity.Type = "receipt"
			activity.Description = "Added to inventory"
		case "update":
			activity.Type = "adjustment"
			activity.Description = "Stock adjusted"
		case "delete":
			activity.Type = "issue"
			activity.Description = "Removed from inventory"
		default:
			activity.Description = "Inventory updated"
		}
		activity.PartName = entityID[:12] + "..."
		activity.CreatedAt = createdAt.UTC().Format(time.RFC3339)
		items = append(items, activity)
	}

	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}
