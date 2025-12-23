package store

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/edvirons/ssp/ims/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DevicesSnapshotRepo struct{ pool *pgxpool.Pool }

func (r *DevicesSnapshotRepo) Upsert(ctx context.Context, d models.DeviceSnapshot) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO devices_snapshot (
			tenant_id, device_id, school_id, model, serial, asset_tag, status, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		ON CONFLICT (tenant_id, device_id)
		DO UPDATE SET
		  school_id=EXCLUDED.school_id,
		  model=EXCLUDED.model,
		  serial=EXCLUDED.serial,
		  asset_tag=EXCLUDED.asset_tag,
		  status=EXCLUDED.status,
		  updated_at=EXCLUDED.updated_at
	`, d.TenantID, d.DeviceID, d.SchoolID, d.Model, d.Serial, d.AssetTag, d.Status, d.UpdatedAt)
	return err
}

func (r *DevicesSnapshotRepo) Get(ctx context.Context, tenantID, deviceID string) (models.DeviceSnapshot, error) {
	var d models.DeviceSnapshot
	row := r.pool.QueryRow(ctx, `
		SELECT tenant_id, device_id, school_id, model, serial, asset_tag, status, updated_at
		FROM devices_snapshot
		WHERE tenant_id=$1 AND device_id=$2
	`, tenantID, deviceID)
	if err := row.Scan(&d.TenantID, &d.DeviceID, &d.SchoolID, &d.Model, &d.Serial, &d.AssetTag, &d.Status, &d.UpdatedAt); err != nil {
		return models.DeviceSnapshot{}, errors.New("not found")
	}
	return d, nil
}

func NewDeviceSnapshot(tenantID, deviceID string) models.DeviceSnapshot {
	return models.DeviceSnapshot{TenantID: tenantID, DeviceID: deviceID, UpdatedAt: time.Now().UTC()}
}

type DeviceSnapshotListParams struct {
	TenantID string
	SchoolID string
	Query    string
	Status   string
	Limit    int
	Offset   int
}

func (r *DevicesSnapshotRepo) List(ctx context.Context, p DeviceSnapshotListParams) ([]models.DeviceSnapshot, int, error) {
	conds := []string{"tenant_id=$1"}
	args := []any{p.TenantID}
	argN := 2

	if p.SchoolID != "" {
		conds = append(conds, "school_id=$"+itoa(argN))
		args = append(args, p.SchoolID)
		argN++
	}
	if p.Query != "" {
		conds = append(conds, "(serial ILIKE $"+itoa(argN)+" OR asset_tag ILIKE $"+itoa(argN)+" OR model ILIKE $"+itoa(argN)+")")
		args = append(args, "%"+p.Query+"%")
		argN++
	}
	if p.Status != "" {
		conds = append(conds, "status=$"+itoa(argN))
		args = append(args, p.Status)
		argN++
	}

	// Count total
	countSQL := "SELECT COUNT(*) FROM devices_snapshot WHERE " + strings.Join(conds, " AND ")
	var total int
	if err := r.pool.QueryRow(ctx, countSQL, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Fetch items
	limit := p.Limit
	if limit <= 0 {
		limit = 50
	}
	offset := p.Offset
	if offset < 0 {
		offset = 0
	}

	args = append(args, limit, offset)
	sql := `
		SELECT tenant_id, device_id, school_id, model, serial, asset_tag, status, updated_at
		FROM devices_snapshot
		WHERE ` + strings.Join(conds, " AND ") + `
		ORDER BY updated_at DESC
		LIMIT $` + itoa(argN) + ` OFFSET $` + itoa(argN+1)

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	items := []models.DeviceSnapshot{} // Initialize as empty slice, not nil
	for rows.Next() {
		var d models.DeviceSnapshot
		if err := rows.Scan(&d.TenantID, &d.DeviceID, &d.SchoolID, &d.Model, &d.Serial, &d.AssetTag, &d.Status, &d.UpdatedAt); err != nil {
			return nil, 0, err
		}
		items = append(items, d)
	}

	return items, total, nil
}

func (r *DevicesSnapshotRepo) Count(ctx context.Context, tenantID string) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM devices_snapshot WHERE tenant_id=$1", tenantID).Scan(&count)
	return count, err
}

// DeviceStats represents aggregate statistics for devices
type DeviceStats struct {
	Total        int            `json:"total"`
	ByStatus     map[string]int `json:"byStatus"`
	BySchool     map[string]int `json:"bySchool"`
	UniqueModels int            `json:"uniqueModels"`
}

func (r *DevicesSnapshotRepo) Stats(ctx context.Context, tenantID string) (DeviceStats, error) {
	stats := DeviceStats{
		ByStatus: make(map[string]int),
		BySchool: make(map[string]int),
	}

	// Get total count
	if err := r.pool.QueryRow(ctx,
		"SELECT COUNT(*) FROM devices_snapshot WHERE tenant_id=$1", tenantID).Scan(&stats.Total); err != nil {
		return stats, err
	}

	// Get counts by status
	rows, err := r.pool.Query(ctx, `
		SELECT COALESCE(status, 'unknown'), COUNT(*)
		FROM devices_snapshot
		WHERE tenant_id=$1
		GROUP BY status
	`, tenantID)
	if err != nil {
		return stats, err
	}
	defer rows.Close()
	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			continue
		}
		stats.ByStatus[status] = count
	}

	// Get unique models count
	if err := r.pool.QueryRow(ctx,
		"SELECT COUNT(DISTINCT model) FROM devices_snapshot WHERE tenant_id=$1", tenantID).Scan(&stats.UniqueModels); err != nil {
		return stats, err
	}

	return stats, nil
}

// DeviceModel represents a unique device model
type DeviceModel struct {
	Model string `json:"model"`
	Count int    `json:"count"`
}

func (r *DevicesSnapshotRepo) ListModels(ctx context.Context, tenantID string) ([]DeviceModel, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT COALESCE(model, 'Unknown'), COUNT(*) as count
		FROM devices_snapshot
		WHERE tenant_id=$1
		GROUP BY model
		ORDER BY count DESC
	`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	models := []DeviceModel{} // Initialize as empty slice, not nil
	for rows.Next() {
		var m DeviceModel
		if err := rows.Scan(&m.Model, &m.Count); err != nil {
			continue
		}
		models = append(models, m)
	}
	return models, nil
}

// ListBySchool returns all devices for a school (no pagination, for inventory view)
func (r *DevicesSnapshotRepo) ListBySchool(ctx context.Context, tenantID, schoolID string) ([]models.DeviceSnapshot, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT tenant_id, device_id, school_id, model, serial, asset_tag, status, updated_at
		FROM devices_snapshot
		WHERE tenant_id=$1 AND school_id=$2
		ORDER BY updated_at DESC
	`, tenantID, schoolID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []models.DeviceSnapshot{}
	for rows.Next() {
		var d models.DeviceSnapshot
		if err := rows.Scan(&d.TenantID, &d.DeviceID, &d.SchoolID, &d.Model, &d.Serial, &d.AssetTag, &d.Status, &d.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, d)
	}
	return items, nil
}

func (r *DevicesSnapshotRepo) ListMakes(ctx context.Context, tenantID string) ([]string, error) {
	// Extract make from model (first word before space or the whole string)
	rows, err := r.pool.Query(ctx, `
		SELECT DISTINCT
			CASE
				WHEN position(' ' in model) > 0 THEN split_part(model, ' ', 1)
				ELSE model
			END as make
		FROM devices_snapshot
		WHERE tenant_id=$1 AND model IS NOT NULL AND model != ''
		ORDER BY make
	`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	makes := []string{} // Initialize as empty slice, not nil
	for rows.Next() {
		var make string
		if err := rows.Scan(&make); err != nil {
			continue
		}
		makes = append(makes, make)
	}
	return makes, nil
}
