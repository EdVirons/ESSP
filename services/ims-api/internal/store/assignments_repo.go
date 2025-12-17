package store

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/edvirons/ssp/ims/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AssignmentsRepo struct{ pool *pgxpool.Pool }

func (r *AssignmentsRepo) Create(ctx context.Context, a models.DeviceAssignment) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO device_assignments (id, tenant_id, device_id, location_id, assigned_to_user, assignment_type, effective_from, effective_to, notes, created_by, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
	`, a.ID, a.TenantID, a.DeviceID, a.LocationID, a.AssignedToUser, a.AssignmentType, a.EffectiveFrom, a.EffectiveTo, a.Notes, a.CreatedBy, a.CreatedAt)
	return err
}

// AssignDevice creates a new assignment and ends any current assignment
func (r *AssignmentsRepo) AssignDevice(ctx context.Context, a models.DeviceAssignment) error {
	now := time.Now().UTC()

	// End current assignment if exists
	_, err := r.pool.Exec(ctx, `
		UPDATE device_assignments SET effective_to=$3
		WHERE tenant_id=$1 AND device_id=$2 AND effective_to IS NULL
	`, a.TenantID, a.DeviceID, now)
	if err != nil {
		return err
	}

	// Create new assignment
	a.EffectiveFrom = now
	a.CreatedAt = now
	return r.Create(ctx, a)
}

// UnassignDevice ends the current assignment without creating a new one
func (r *AssignmentsRepo) UnassignDevice(ctx context.Context, tenantID, deviceID string) error {
	now := time.Now().UTC()
	result, err := r.pool.Exec(ctx, `
		UPDATE device_assignments SET effective_to=$3
		WHERE tenant_id=$1 AND device_id=$2 AND effective_to IS NULL
	`, tenantID, deviceID, now)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errors.New("no current assignment")
	}
	return nil
}

func (r *AssignmentsRepo) Get(ctx context.Context, tenantID, id string) (models.DeviceAssignment, error) {
	var a models.DeviceAssignment
	row := r.pool.QueryRow(ctx, `
		SELECT id, tenant_id, device_id, location_id, assigned_to_user, assignment_type, effective_from, effective_to, notes, created_by, created_at
		FROM device_assignments WHERE tenant_id=$1 AND id=$2
	`, tenantID, id)
	if err := row.Scan(&a.ID, &a.TenantID, &a.DeviceID, &a.LocationID, &a.AssignedToUser, &a.AssignmentType, &a.EffectiveFrom, &a.EffectiveTo, &a.Notes, &a.CreatedBy, &a.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.DeviceAssignment{}, errors.New("not found")
		}
		return models.DeviceAssignment{}, err
	}
	return a, nil
}

// GetCurrent returns the current (active) assignment for a device
func (r *AssignmentsRepo) GetCurrent(ctx context.Context, tenantID, deviceID string) (models.DeviceAssignment, error) {
	var a models.DeviceAssignment
	row := r.pool.QueryRow(ctx, `
		SELECT id, tenant_id, device_id, location_id, assigned_to_user, assignment_type, effective_from, effective_to, notes, created_by, created_at
		FROM device_assignments
		WHERE tenant_id=$1 AND device_id=$2 AND effective_to IS NULL
		ORDER BY effective_from DESC
		LIMIT 1
	`, tenantID, deviceID)
	if err := row.Scan(&a.ID, &a.TenantID, &a.DeviceID, &a.LocationID, &a.AssignedToUser, &a.AssignmentType, &a.EffectiveFrom, &a.EffectiveTo, &a.Notes, &a.CreatedBy, &a.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.DeviceAssignment{}, errors.New("not found")
		}
		return models.DeviceAssignment{}, err
	}
	return a, nil
}

type AssignmentListParams struct {
	TenantID       string
	DeviceID       string
	LocationID     string
	AssignmentType string
	CurrentOnly    bool
	Limit          int
	HasCursor      bool
	CursorAt       time.Time
	CursorID       string
}

func (r *AssignmentsRepo) List(ctx context.Context, p AssignmentListParams) ([]models.DeviceAssignment, string, error) {
	conds := []string{"tenant_id=$1"}
	args := []any{p.TenantID}
	argN := 2

	if p.DeviceID != "" {
		conds = append(conds, "device_id=$"+itoa(argN))
		args = append(args, p.DeviceID)
		argN++
	}
	if p.LocationID != "" {
		conds = append(conds, "location_id=$"+itoa(argN))
		args = append(args, p.LocationID)
		argN++
	}
	if p.AssignmentType != "" {
		conds = append(conds, "assignment_type=$"+itoa(argN))
		args = append(args, p.AssignmentType)
		argN++
	}
	if p.CurrentOnly {
		conds = append(conds, "effective_to IS NULL")
	}
	if p.HasCursor {
		conds = append(conds, "(effective_from, id) < ($"+itoa(argN)+", $"+itoa(argN+1)+")")
		args = append(args, p.CursorAt, p.CursorID)
		argN += 2
	}

	limit := p.Limit
	if limit <= 0 {
		limit = 50
	}
	limitPlus := limit + 1
	args = append(args, limitPlus)

	sql := `
		SELECT id, tenant_id, device_id, location_id, assigned_to_user, assignment_type, effective_from, effective_to, notes, created_by, created_at
		FROM device_assignments
		WHERE ` + strings.Join(conds, " AND ") + `
		ORDER BY effective_from DESC, id DESC
		LIMIT $` + itoa(argN)

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	out := []models.DeviceAssignment{}
	for rows.Next() {
		var a models.DeviceAssignment
		if err := rows.Scan(&a.ID, &a.TenantID, &a.DeviceID, &a.LocationID, &a.AssignedToUser, &a.AssignmentType, &a.EffectiveFrom, &a.EffectiveTo, &a.Notes, &a.CreatedBy, &a.CreatedAt); err != nil {
			return nil, "", err
		}
		out = append(out, a)
	}

	next := ""
	if len(out) > limit {
		last := out[limit-1]
		next = EncodeCursor(last.EffectiveFrom, last.ID)
		out = out[:limit]
	}
	return out, next, nil
}

// ListByLocation returns all current device assignments at a location
func (r *AssignmentsRepo) ListByLocation(ctx context.Context, tenantID, locationID string) ([]models.DeviceAssignment, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, tenant_id, device_id, location_id, assigned_to_user, assignment_type, effective_from, effective_to, notes, created_by, created_at
		FROM device_assignments
		WHERE tenant_id=$1 AND location_id=$2 AND effective_to IS NULL
		ORDER BY effective_from DESC
	`, tenantID, locationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []models.DeviceAssignment{}
	for rows.Next() {
		var a models.DeviceAssignment
		if err := rows.Scan(&a.ID, &a.TenantID, &a.DeviceID, &a.LocationID, &a.AssignedToUser, &a.AssignmentType, &a.EffectiveFrom, &a.EffectiveTo, &a.Notes, &a.CreatedBy, &a.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, nil
}

// GetHistory returns assignment history for a device
func (r *AssignmentsRepo) GetHistory(ctx context.Context, tenantID, deviceID string, limit int) ([]models.DeviceAssignment, error) {
	if limit <= 0 {
		limit = 20
	}
	rows, err := r.pool.Query(ctx, `
		SELECT id, tenant_id, device_id, location_id, assigned_to_user, assignment_type, effective_from, effective_to, notes, created_by, created_at
		FROM device_assignments
		WHERE tenant_id=$1 AND device_id=$2
		ORDER BY effective_from DESC
		LIMIT $3
	`, tenantID, deviceID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []models.DeviceAssignment{}
	for rows.Next() {
		var a models.DeviceAssignment
		if err := rows.Scan(&a.ID, &a.TenantID, &a.DeviceID, &a.LocationID, &a.AssignedToUser, &a.AssignmentType, &a.EffectiveFrom, &a.EffectiveTo, &a.Notes, &a.CreatedBy, &a.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, nil
}

// BulkAssign assigns multiple devices to a location
func (r *AssignmentsRepo) BulkAssign(ctx context.Context, tenantID string, deviceIDs []string, locationID string, assignmentType models.AssignmentType, createdBy string) (int, error) {
	if len(deviceIDs) == 0 {
		return 0, nil
	}

	now := time.Now().UTC()

	// End current assignments
	_, err := r.pool.Exec(ctx, `
		UPDATE device_assignments SET effective_to=$2
		WHERE tenant_id=$1 AND device_id = ANY($3) AND effective_to IS NULL
	`, tenantID, now, deviceIDs)
	if err != nil {
		return 0, err
	}

	// Create new assignments
	count := 0
	for _, deviceID := range deviceIDs {
		id := NewID("asn")
		_, err := r.pool.Exec(ctx, `
			INSERT INTO device_assignments (id, tenant_id, device_id, location_id, assignment_type, effective_from, created_by, created_at)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		`, id, tenantID, deviceID, locationID, assignmentType, now, createdBy, now)
		if err != nil {
			return count, err
		}
		count++
	}
	return count, nil
}

// GetCurrentAssignmentMap returns a map of deviceID -> current assignment for quick lookups
func (r *AssignmentsRepo) GetCurrentAssignmentMap(ctx context.Context, tenantID string, deviceIDs []string) (map[string]models.DeviceAssignment, error) {
	if len(deviceIDs) == 0 {
		return make(map[string]models.DeviceAssignment), nil
	}

	rows, err := r.pool.Query(ctx, `
		SELECT id, tenant_id, device_id, location_id, assigned_to_user, assignment_type, effective_from, effective_to, notes, created_by, created_at
		FROM device_assignments
		WHERE tenant_id=$1 AND device_id = ANY($2) AND effective_to IS NULL
	`, tenantID, deviceIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]models.DeviceAssignment)
	for rows.Next() {
		var a models.DeviceAssignment
		if err := rows.Scan(&a.ID, &a.TenantID, &a.DeviceID, &a.LocationID, &a.AssignedToUser, &a.AssignmentType, &a.EffectiveFrom, &a.EffectiveTo, &a.Notes, &a.CreatedBy, &a.CreatedAt); err != nil {
			return nil, err
		}
		result[a.DeviceID] = a
	}
	return result, nil
}
