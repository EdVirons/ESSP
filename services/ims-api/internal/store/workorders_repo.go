package store

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/edvirons/ssp/ims/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WorkOrderRepo struct {
	pool *pgxpool.Pool
}

func (r *WorkOrderRepo) Create(ctx context.Context, wo models.WorkOrder) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO work_orders (
			id, incident_id, tenant_id, school_id, device_id, status, service_shop_id, assigned_staff_id, repair_location,
			assigned_to, task_type, project_id, phase_id, cost_estimate_cents, notes,
			created_from_project, created_by_user_id, created_by_user_name, created_at, updated_at
		) VALUES (
			$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,NULLIF($12,''),NULLIF($13,''),$14,$15,$16,$17,$18,$19,$20
		)
	`, wo.ID, wo.IncidentID, wo.TenantID, wo.SchoolID, wo.DeviceID, wo.Status,
		wo.ServiceShopID, wo.AssignedStaffID, wo.RepairLocation,
		wo.AssignedTo, wo.TaskType, wo.ProjectID, wo.PhaseID, wo.CostEstimateCents, wo.Notes,
		wo.CreatedFromProject, wo.CreatedByUserID, wo.CreatedByUserName, wo.CreatedAt, wo.UpdatedAt)
	return err
}

func (r *WorkOrderRepo) GetByID(ctx context.Context, tenantID, schoolID, id string) (models.WorkOrder, error) {
	var wo models.WorkOrder
	row := r.pool.QueryRow(ctx, `
		SELECT id, incident_id, tenant_id, school_id, device_id, status, service_shop_id, assigned_staff_id, repair_location, service_shop_id, assigned_staff_id, repair_location,
		       assigned_to, task_type, cost_estimate_cents, notes, created_at, updated_at
		FROM work_orders
		WHERE tenant_id=$1 AND school_id=$2 AND id=$3
	`, tenantID, schoolID, id)

	err := row.Scan(&wo.ID, &wo.IncidentID, &wo.TenantID, &wo.SchoolID, &wo.DeviceID, &wo.Status, &wo.ServiceShopID, &wo.AssignedStaffID, &wo.RepairLocation,
		&wo.AssignedTo, &wo.TaskType, &wo.CostEstimateCents, &wo.Notes, &wo.CreatedAt, &wo.UpdatedAt)
	if err != nil {
		return models.WorkOrder{}, errors.New("not found")
	}
	return wo, nil
}

type WorkOrderListParams struct {
	TenantID string
	SchoolID string
	Status string
	DeviceID string
	IncidentID string
	Limit int

	HasCursor bool
	CursorCreatedAt time.Time
	CursorID string
}

func (r *WorkOrderRepo) List(ctx context.Context, p WorkOrderListParams) ([]models.WorkOrder, string, error) {
	conds := []string{"tenant_id=$1", "school_id=$2"}
	args := []any{p.TenantID, p.SchoolID}
	argN := 3

	if p.Status != "" {
		conds = append(conds, "status=$"+itoa(argN))
		args = append(args, p.Status)
		argN++
	}
	if p.DeviceID != "" {
		conds = append(conds, "device_id=$"+itoa(argN))
		args = append(args, p.DeviceID)
		argN++
	}
	if p.IncidentID != "" {
		conds = append(conds, "incident_id=$"+itoa(argN))
		args = append(args, p.IncidentID)
		argN++
	}
	if p.HasCursor {
		conds = append(conds, "(created_at, id) < ($"+itoa(argN)+", $"+itoa(argN+1)+")")
		args = append(args, p.CursorCreatedAt, p.CursorID)
		argN += 2
	}

	limitPlus := p.Limit + 1
	args = append(args, limitPlus)

	sql := `
		SELECT id, incident_id, tenant_id, school_id, device_id, status, service_shop_id, assigned_staff_id, repair_location, service_shop_id, assigned_staff_id, repair_location,
		       assigned_to, task_type, cost_estimate_cents, notes, created_at, updated_at
		FROM work_orders
		WHERE ` + strings.Join(conds, " AND ") + `
		ORDER BY created_at DESC, id DESC
		LIMIT $` + itoa(argN)

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	out := []models.WorkOrder{}
	for rows.Next() {
		var wo models.WorkOrder
		if err := rows.Scan(&wo.ID, &wo.IncidentID, &wo.TenantID, &wo.SchoolID, &wo.DeviceID, &wo.Status, &wo.ServiceShopID, &wo.AssignedStaffID, &wo.RepairLocation,
			&wo.AssignedTo, &wo.TaskType, &wo.CostEstimateCents, &wo.Notes, &wo.CreatedAt, &wo.UpdatedAt); err != nil {
			return nil, "", err
		}
		out = append(out, wo)
	}

	next := ""
	if len(out) > p.Limit {
		last := out[p.Limit-1]
		next = EncodeCursor(last.CreatedAt, last.ID)
		out = out[:p.Limit]
	}
	return out, next, nil
}

func (r *WorkOrderRepo) UpdateStatus(ctx context.Context, tenantID, schoolID, id string, status models.WorkOrderStatus, now time.Time) (models.WorkOrder, error) {
	_, err := r.pool.Exec(ctx, `
		UPDATE work_orders
		SET status=$4, updated_at=$5
		WHERE tenant_id=$1 AND school_id=$2 AND id=$3
	`, tenantID, schoolID, id, status, now)
	if err != nil {
		return models.WorkOrder{}, err
	}
	return r.GetByID(ctx, tenantID, schoolID, id)
}


func (r *WorkOrderRepo) SetApprovalStatus(ctx context.Context, tenantID, schoolID, workOrderID, status string) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE work_orders SET approval_status=$4, updated_at=$5
		WHERE tenant_id=$1 AND school_id=$2 AND id=$3
	`, tenantID, schoolID, workOrderID, status, time.Now().UTC())
	return err
}


func (r *WorkOrderRepo) ListByPhase(ctx context.Context, tenantID, phaseID string) ([]models.WorkOrder, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, incident_id, tenant_id, school_id, device_id, status, service_shop_id, assigned_staff_id, repair_location, assigned_to,
			task_type, project_id, phase_id, onsite_contact_id, approval_status, cost_estimate_cents, notes,
			COALESCE(created_from_project, FALSE), COALESCE(created_by_user_id, ''), COALESCE(created_by_user_name, ''),
			created_at, updated_at
		FROM work_orders
		WHERE tenant_id=$1 AND phase_id=$2
		ORDER BY created_at DESC, id DESC
	`, tenantID, phaseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []models.WorkOrder{}
	for rows.Next() {
		var x models.WorkOrder
		if err := rows.Scan(&x.ID, &x.IncidentID, &x.TenantID, &x.SchoolID, &x.DeviceID, &x.Status, &x.ServiceShopID, &x.AssignedStaffID, &x.RepairLocation, &x.AssignedTo,
			&x.TaskType, &x.ProjectID, &x.PhaseID, &x.OnsiteContactID, &x.ApprovalStatus, &x.CostEstimateCents, &x.Notes,
			&x.CreatedFromProject, &x.CreatedByUserID, &x.CreatedByUserName,
			&x.CreatedAt, &x.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, x)
	}
	return out, nil
}

// ProjectWOListParams holds parameters for listing work orders by project.
type ProjectWOListParams struct {
	TenantID   string
	ProjectID  string
	PhaseID    string
	Status     string
	Limit      int
	HasCursor  bool
	CursorTime time.Time
	CursorID   string
}

// ListByProject lists work orders for a specific project with optional filters.
func (r *WorkOrderRepo) ListByProject(ctx context.Context, p ProjectWOListParams) ([]models.WorkOrder, string, error) {
	conds := []string{"tenant_id=$1", "project_id=$2"}
	args := []any{p.TenantID, p.ProjectID}
	argN := 3

	if p.PhaseID != "" {
		conds = append(conds, "phase_id=$"+itoa(argN))
		args = append(args, p.PhaseID)
		argN++
	}
	if p.Status != "" {
		conds = append(conds, "status=$"+itoa(argN))
		args = append(args, p.Status)
		argN++
	}
	if p.HasCursor {
		conds = append(conds, "(created_at, id) < ($"+itoa(argN)+", $"+itoa(argN+1)+")")
		args = append(args, p.CursorTime, p.CursorID)
		argN += 2
	}

	limitPlus := p.Limit + 1
	args = append(args, limitPlus)

	sql := `
		SELECT id, incident_id, tenant_id, school_id, device_id, status, service_shop_id, assigned_staff_id, repair_location, assigned_to,
			task_type, project_id, phase_id, onsite_contact_id, approval_status, cost_estimate_cents, notes,
			COALESCE(created_from_project, FALSE), COALESCE(created_by_user_id, ''), COALESCE(created_by_user_name, ''),
			created_at, updated_at
		FROM work_orders
		WHERE ` + strings.Join(conds, " AND ") + `
		ORDER BY created_at DESC, id DESC
		LIMIT $` + itoa(argN)

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	out := []models.WorkOrder{}
	for rows.Next() {
		var x models.WorkOrder
		if err := rows.Scan(&x.ID, &x.IncidentID, &x.TenantID, &x.SchoolID, &x.DeviceID, &x.Status, &x.ServiceShopID, &x.AssignedStaffID, &x.RepairLocation, &x.AssignedTo,
			&x.TaskType, &x.ProjectID, &x.PhaseID, &x.OnsiteContactID, &x.ApprovalStatus, &x.CostEstimateCents, &x.Notes,
			&x.CreatedFromProject, &x.CreatedByUserID, &x.CreatedByUserName,
			&x.CreatedAt, &x.UpdatedAt); err != nil {
			return nil, "", err
		}
		out = append(out, x)
	}

	next := ""
	if len(out) > p.Limit {
		last := out[p.Limit-1]
		next = EncodeCursor(last.CreatedAt, last.ID)
		out = out[:p.Limit]
	}
	return out, next, nil
}

// UpdateFields performs a partial update on a work order using dynamic field updates.
// Only fields present in the updates map will be modified.
func (r *WorkOrderRepo) UpdateFields(ctx context.Context, tenantID, schoolID, id string, updates map[string]any, now time.Time) (models.WorkOrder, error) {
	if len(updates) == 0 {
		return r.GetByID(ctx, tenantID, schoolID, id)
	}

	// Build dynamic UPDATE query
	setClauses := []string{}
	args := []any{tenantID, schoolID, id}
	argN := 4

	for field, value := range updates {
		setClauses = append(setClauses, field+"=$"+itoa(argN))
		args = append(args, value)
		argN++
	}

	// Always update updated_at
	setClauses = append(setClauses, "updated_at=$"+itoa(argN))
	args = append(args, now)

	query := `
		UPDATE work_orders
		SET ` + strings.Join(setClauses, ", ") + `
		WHERE tenant_id=$1 AND school_id=$2 AND id=$3
	`

	result, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		return models.WorkOrder{}, err
	}
	if result.RowsAffected() == 0 {
		return models.WorkOrder{}, errors.New("not found")
	}

	return r.GetByID(ctx, tenantID, schoolID, id)
}

// UpdateStatusWithRework updates status and increments rework count.
func (r *WorkOrderRepo) UpdateStatusWithRework(ctx context.Context, tenantID, schoolID, id string, status models.WorkOrderStatus, reason string, now time.Time) (models.WorkOrder, error) {
	_, err := r.pool.Exec(ctx, `
		UPDATE work_orders
		SET status = $4,
			rework_count = COALESCE(rework_count, 0) + 1,
			last_rework_at = $5,
			last_rework_reason = $6,
			updated_at = $5
		WHERE tenant_id = $1 AND school_id = $2 AND id = $3
	`, tenantID, schoolID, id, status, now, reason)
	if err != nil {
		return models.WorkOrder{}, err
	}
	return r.GetByID(ctx, tenantID, schoolID, id)
}

// GetByIDs fetches multiple work orders by their IDs for bulk operations.
func (r *WorkOrderRepo) GetByIDs(ctx context.Context, tenantID, schoolID string, ids []string) ([]models.WorkOrder, error) {
	if len(ids) == 0 {
		return []models.WorkOrder{}, nil
	}

	// Build IN clause
	placeholders := make([]string, len(ids))
	args := []any{tenantID, schoolID}
	for i, id := range ids {
		placeholders[i] = "$" + itoa(i+3)
		args = append(args, id)
	}

	query := `
		SELECT id, incident_id, tenant_id, school_id, device_id, status, service_shop_id, assigned_staff_id, repair_location, assigned_to,
			task_type, COALESCE(project_id, ''), COALESCE(phase_id, ''), COALESCE(onsite_contact_id, ''), COALESCE(approval_status, 'not_required'), cost_estimate_cents, notes,
			COALESCE(created_from_project, FALSE), COALESCE(created_by_user_id, ''), COALESCE(created_by_user_name, ''),
			COALESCE(rework_count, 0), last_rework_at, COALESCE(last_rework_reason, ''),
			created_at, updated_at
		FROM work_orders
		WHERE tenant_id = $1 AND school_id = $2 AND id IN (` + strings.Join(placeholders, ", ") + `)
		ORDER BY created_at DESC
	`

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []models.WorkOrder{}
	for rows.Next() {
		var x models.WorkOrder
		if err := rows.Scan(&x.ID, &x.IncidentID, &x.TenantID, &x.SchoolID, &x.DeviceID, &x.Status, &x.ServiceShopID, &x.AssignedStaffID, &x.RepairLocation, &x.AssignedTo,
			&x.TaskType, &x.ProjectID, &x.PhaseID, &x.OnsiteContactID, &x.ApprovalStatus, &x.CostEstimateCents, &x.Notes,
			&x.CreatedFromProject, &x.CreatedByUserID, &x.CreatedByUserName,
			&x.ReworkCount, &x.LastReworkAt, &x.LastReworkReason,
			&x.CreatedAt, &x.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, x)
	}
	return out, nil
}

// BulkUpdateStatus updates status for multiple work orders in a single transaction.
func (r *WorkOrderRepo) BulkUpdateStatus(ctx context.Context, tenantID, schoolID string, ids []string, status models.WorkOrderStatus, now time.Time) error {
	if len(ids) == 0 {
		return nil
	}

	// Build IN clause
	placeholders := make([]string, len(ids))
	args := []any{tenantID, schoolID, status, now}
	for i, id := range ids {
		placeholders[i] = "$" + itoa(i+5)
		args = append(args, id)
	}

	query := `
		UPDATE work_orders
		SET status = $3, updated_at = $4
		WHERE tenant_id = $1 AND school_id = $2 AND id IN (` + strings.Join(placeholders, ", ") + `)
	`

	_, err := r.pool.Exec(ctx, query, args...)
	return err
}

// BulkUpdateAssignment updates assignment for multiple work orders.
func (r *WorkOrderRepo) BulkUpdateAssignment(ctx context.Context, tenantID, schoolID string, ids []string, staffID, shopID *string, now time.Time) error {
	if len(ids) == 0 {
		return nil
	}

	// Build dynamic SET clause
	setClauses := []string{"updated_at = $3"}
	args := []any{tenantID, schoolID, now}
	argN := 4

	if staffID != nil {
		setClauses = append(setClauses, "assigned_staff_id = $"+itoa(argN))
		args = append(args, *staffID)
		argN++
	}
	if shopID != nil {
		setClauses = append(setClauses, "service_shop_id = $"+itoa(argN))
		args = append(args, *shopID)
		argN++
	}

	// Build IN clause
	placeholders := make([]string, len(ids))
	for i, id := range ids {
		placeholders[i] = "$" + itoa(argN+i)
		args = append(args, id)
	}

	query := `
		UPDATE work_orders
		SET ` + strings.Join(setClauses, ", ") + `
		WHERE tenant_id = $1 AND school_id = $2 AND id IN (` + strings.Join(placeholders, ", ") + `)
	`

	_, err := r.pool.Exec(ctx, query, args...)
	return err
}

// GetReworkCount returns the rework count for a work order.
func (r *WorkOrderRepo) GetReworkCount(ctx context.Context, tenantID, schoolID, id string) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx, `
		SELECT COALESCE(rework_count, 0)
		FROM work_orders
		WHERE tenant_id = $1 AND school_id = $2 AND id = $3
	`, tenantID, schoolID, id).Scan(&count)
	return count, err
}
