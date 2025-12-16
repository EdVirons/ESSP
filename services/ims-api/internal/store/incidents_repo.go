package store

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/edvirons/ssp/ims/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IncidentRepo struct {
	pool *pgxpool.Pool
}

func (r *IncidentRepo) Create(ctx context.Context, inc models.Incident) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO incidents (
			id, tenant_id, school_id, device_id,
			school_name, county_id, county_name, sub_county_id, sub_county_name,
			contact_name, contact_phone, contact_email,
			device_serial, device_asset_tag, device_model_id, device_make, device_model, device_category,
			category, severity, status,
			title, description, reported_by, sla_due_at, sla_breached, created_at, updated_at
		) VALUES (
			$1,$2,$3,$4,
			$5,$6,$7,$8,$9,
			$10,$11,$12,
			$13,$14,$15,$16,$17,$18,
			$19,$20,$21,
			$22,$23,$24,$25,$26,$27,$28
		)
	`, inc.ID, inc.TenantID, inc.SchoolID, inc.DeviceID,
		inc.SchoolName, inc.CountyID, inc.CountyName, inc.SubCountyID, inc.SubCountyName,
		inc.ContactName, inc.ContactPhone, inc.ContactEmail,
		inc.DeviceSerial, inc.DeviceAssetTag, inc.DeviceModelID, inc.DeviceMake, inc.DeviceModel, inc.DeviceCategory,
		inc.Category, inc.Severity, inc.Status,
		inc.Title, inc.Description, inc.ReportedBy, inc.SLADueAt, inc.SLABreached, inc.CreatedAt, inc.UpdatedAt)
	return err
}

func (r *IncidentRepo) GetByID(ctx context.Context, tenantID, schoolID, id string) (models.Incident, error) {
	var inc models.Incident
	row := r.pool.QueryRow(ctx, `
		SELECT id, tenant_id, school_id, device_id, category, severity, status,
		       title, description, reported_by, sla_due_at, sla_breached, created_at, updated_at
		FROM incidents
		WHERE tenant_id=$1 AND school_id=$2 AND id=$3
	`, tenantID, schoolID, id)

	err := row.Scan(&inc.ID, &inc.TenantID, &inc.SchoolID, &inc.DeviceID, &inc.Category, &inc.Severity, &inc.Status,
		&inc.Title, &inc.Description, &inc.ReportedBy, &inc.SLADueAt, &inc.SLABreached, &inc.CreatedAt, &inc.UpdatedAt)
	if err != nil {
		return models.Incident{}, errors.New("not found")
	}
	return inc, nil
}

type IncidentListParams struct {
	TenantID string
	SchoolID string
	Status   string
	DeviceID string
	Query    string
	Limit    int

	HasCursor      bool
	CursorCreatedAt time.Time
	CursorID       string
}

func (r *IncidentRepo) List(ctx context.Context, p IncidentListParams) ([]models.Incident, string, error) {
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
	if p.Query != "" {
		conds = append(conds, "(title ILIKE $"+itoa(argN)+" OR description ILIKE $"+itoa(argN)+")")
		args = append(args, "%"+p.Query+"%")
		argN++
	}
	if p.HasCursor {
		conds = append(conds, "(created_at, id) < ($"+itoa(argN)+", $"+itoa(argN+1)+")")
		args = append(args, p.CursorCreatedAt, p.CursorID)
		argN += 2
	}

	// +1 fetch for nextCursor
	limitPlus := p.Limit + 1
	args = append(args, limitPlus)

	sql := `
		SELECT id, tenant_id, school_id, device_id, category, severity, status,
		       title, description, reported_by, sla_due_at, sla_breached, created_at, updated_at
		FROM incidents
		WHERE ` + strings.Join(conds, " AND ") + `
		ORDER BY created_at DESC, id DESC
		LIMIT $` + itoa(argN)

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	out := []models.Incident{}
	for rows.Next() {
		var inc models.Incident
		if err := rows.Scan(&inc.ID, &inc.TenantID, &inc.SchoolID, &inc.DeviceID, &inc.Category, &inc.Severity, &inc.Status,
			&inc.Title, &inc.Description, &inc.ReportedBy, &inc.SLADueAt, &inc.SLABreached, &inc.CreatedAt, &inc.UpdatedAt); err != nil {
			return nil, "", err
		}
		out = append(out, inc)
	}

	next := ""
	if len(out) > p.Limit {
		last := out[p.Limit-1]
		next = EncodeCursor(last.CreatedAt, last.ID)
		out = out[:p.Limit]
	}
	return out, next, nil
}

func (r *IncidentRepo) UpdateStatus(ctx context.Context, tenantID, schoolID, id string, status models.IncidentStatus, now time.Time) (models.Incident, error) {
	_, err := r.pool.Exec(ctx, `
		UPDATE incidents
		SET status=$4, updated_at=$5
		WHERE tenant_id=$1 AND school_id=$2 AND id=$3
	`, tenantID, schoolID, id, status, now)
	if err != nil {
		return models.Incident{}, err
	}
	return r.GetByID(ctx, tenantID, schoolID, id)
}

func (r *IncidentRepo) MarkSLABreaches(ctx context.Context, now time.Time) (int, error) {
	tag, err := r.pool.Exec(ctx, `
		UPDATE incidents
		SET sla_breached=true, updated_at=$1
		WHERE sla_breached=false AND sla_due_at < $1 AND status NOT IN ('resolved','closed')
	`, now)
	if err != nil {
		return 0, err
	}
	return int(tag.RowsAffected()), nil
}
