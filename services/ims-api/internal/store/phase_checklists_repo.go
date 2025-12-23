package store

import (
	"context"
	"strings"
	"time"

	"github.com/edvirons/ssp/ims/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PhaseChecklistsRepo struct{ pool *pgxpool.Pool }

func (r *PhaseChecklistsRepo) Create(ctx context.Context, t models.PhaseChecklistTemplate) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO phase_checklist_templates (
			id, tenant_id, phase_type, title, description, required, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	`, t.ID, t.TenantID, t.PhaseType, t.Title, t.Description, t.Required, t.CreatedAt, t.UpdatedAt)
	return err
}

func (r *PhaseChecklistsRepo) List(ctx context.Context, tenantID string, phaseType string) ([]models.PhaseChecklistTemplate, error) {
	conds := "tenant_id=$1"
	args := []any{tenantID}
	argN := 2
	if strings.TrimSpace(phaseType) != "" {
		conds += " AND phase_type=$" + itoa(argN)
		args = append(args, strings.TrimSpace(phaseType))
	}
	rows, err := r.pool.Query(ctx, `
		SELECT id, tenant_id, phase_type, title, description, required, created_at, updated_at
		FROM phase_checklist_templates
		WHERE `+conds+`
		ORDER BY phase_type ASC, required DESC, created_at ASC
	`, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []models.PhaseChecklistTemplate{}
	for rows.Next() {
		var x models.PhaseChecklistTemplate
		if err := rows.Scan(&x.ID, &x.TenantID, &x.PhaseType, &x.Title, &x.Description, &x.Required, &x.CreatedAt, &x.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, x)
	}
	return out, nil
}

func (r *PhaseChecklistsRepo) SeedDefaults(ctx context.Context, tenantID string) error {
	now := time.Now().UTC()
	defaults := []models.PhaseChecklistTemplate{
		{ID: NewID("pct"), TenantID: tenantID, PhaseType: models.PhaseSurvey, Title: "Power points verified", Description: "Confirm sockets, load, UPS needs", Required: true, CreatedAt: now, UpdatedAt: now},
		{ID: NewID("pct"), TenantID: tenantID, PhaseType: models.PhaseSurvey, Title: "Network path surveyed", Description: "Cable paths, switch locations, rack", Required: true, CreatedAt: now, UpdatedAt: now},
		{ID: NewID("pct"), TenantID: tenantID, PhaseType: models.PhaseInstall, Title: "Low voltage cabling completed", Description: "Drops installed and tested", Required: true, CreatedAt: now, UpdatedAt: now},
		{ID: NewID("pct"), TenantID: tenantID, PhaseType: models.PhaseIntegrate, Title: "Devices enrolled to MDM", Description: "Group policies applied", Required: true, CreatedAt: now, UpdatedAt: now},
		{ID: NewID("pct"), TenantID: tenantID, PhaseType: models.PhaseCommission, Title: "School sign-off", Description: "Primary contact acceptance", Required: true, CreatedAt: now, UpdatedAt: now},
	}
	for _, d := range defaults {
		_ = r.Create(ctx, d)
	}
	return nil
}
