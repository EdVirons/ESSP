package store

import (
	"context"
	"strings"
	"time"

	"github.com/edvirons/ssp/ims/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PhasesRepo struct{ pool *pgxpool.Pool }

func (r *PhasesRepo) Create(ctx context.Context, p models.ServicePhase) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO service_phases (
			id, tenant_id, project_id, phase_type, status, owner_role, owner_user_id, owner_user_name,
			start_date, end_date, notes, status_changed_at, status_changed_by_user_id, status_changed_by_user_name,
			created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,NULLIF($9,''),NULLIF($10,''),$11,$12,$13,$14,$15,$16)
	`, p.ID, p.TenantID, p.ProjectID, p.PhaseType, p.Status, p.OwnerRole, p.OwnerUserID, p.OwnerUserName,
		p.StartDate, p.EndDate, p.Notes, p.StatusChangedAt, p.StatusChangedByUserID, p.StatusChangedByUserName,
		p.CreatedAt, p.UpdatedAt)
	return err
}

// UpdateStatus updates the phase status with tracking info.
func (r *PhasesRepo) UpdateStatus(ctx context.Context, tenantID, phaseID string, status models.PhaseStatus, changedByUserID, changedByUserName string) error {
	now := time.Now()
	_, err := r.pool.Exec(ctx, `
		UPDATE service_phases
		SET status = $3, status_changed_at = $4, status_changed_by_user_id = $5, status_changed_by_user_name = $6, updated_at = $4
		WHERE tenant_id = $1 AND id = $2
	`, tenantID, phaseID, status, now, changedByUserID, changedByUserName)
	return err
}

// UpdateOwner updates the phase owner user.
func (r *PhasesRepo) UpdateOwner(ctx context.Context, tenantID, phaseID, ownerUserID, ownerUserName string) error {
	now := time.Now()
	_, err := r.pool.Exec(ctx, `
		UPDATE service_phases
		SET owner_user_id = $3, owner_user_name = $4, updated_at = $5
		WHERE tenant_id = $1 AND id = $2
	`, tenantID, phaseID, ownerUserID, ownerUserName, now)
	return err
}

// GetByID retrieves a phase by ID.
func (r *PhasesRepo) GetByID(ctx context.Context, tenantID, phaseID string) (models.ServicePhase, error) {
	var p models.ServicePhase
	var sd, ed *time.Time
	row := r.pool.QueryRow(ctx, `
		SELECT id, tenant_id, project_id, phase_type, status, owner_role, owner_user_id, owner_user_name,
			start_date, end_date, notes, status_changed_at, status_changed_by_user_id, status_changed_by_user_name,
			created_at, updated_at
		FROM service_phases
		WHERE tenant_id = $1 AND id = $2
	`, tenantID, phaseID)
	if err := row.Scan(&p.ID, &p.TenantID, &p.ProjectID, &p.PhaseType, &p.Status, &p.OwnerRole, &p.OwnerUserID, &p.OwnerUserName,
		&sd, &ed, &p.Notes, &p.StatusChangedAt, &p.StatusChangedByUserID, &p.StatusChangedByUserName,
		&p.CreatedAt, &p.UpdatedAt); err != nil {
		return models.ServicePhase{}, err
	}
	if sd != nil {
		p.StartDate = sd.Format("2006-01-02")
	}
	if ed != nil {
		p.EndDate = ed.Format("2006-01-02")
	}
	return p, nil
}

type PhaseListParams struct {
	TenantID string
	ProjectID string
	PhaseType string
	Status string
	Limit int
	HasCursor bool
	CursorCreatedAt time.Time
	CursorID string
}

func (r *PhasesRepo) List(ctx context.Context, p PhaseListParams) ([]models.ServicePhase, string, error) {
	conds := []string{"tenant_id=$1", "project_id=$2"}
	args := []any{p.TenantID, p.ProjectID}
	argN := 3
	if p.PhaseType != "" { conds = append(conds, "phase_type=$"+itoa(argN)); args=append(args,p.PhaseType); argN++ }
	if p.Status != "" { conds = append(conds, "status=$"+itoa(argN)); args=append(args,p.Status); argN++ }
	if p.HasCursor {
		conds = append(conds, "(created_at, id) < ($"+itoa(argN)+", $"+itoa(argN+1)+")")
		args = append(args, p.CursorCreatedAt, p.CursorID); argN += 2
	}
	limitPlus := p.Limit+1
	args = append(args, limitPlus)

	sql := `
		SELECT id, tenant_id, project_id, phase_type, status, owner_role, owner_user_id, owner_user_name,
			start_date, end_date, notes, status_changed_at, status_changed_by_user_id, status_changed_by_user_name,
			created_at, updated_at
		FROM service_phases
		WHERE ` + strings.Join(conds, " AND ") + `
		ORDER BY created_at DESC, id DESC
		LIMIT $` + itoa(argN)

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil { return nil,"",err }
	defer rows.Close()

	out := []models.ServicePhase{}
	for rows.Next() {
		var x models.ServicePhase
		var sd, ed *time.Time
		if err := rows.Scan(&x.ID, &x.TenantID, &x.ProjectID, &x.PhaseType, &x.Status, &x.OwnerRole, &x.OwnerUserID, &x.OwnerUserName,
			&sd, &ed, &x.Notes, &x.StatusChangedAt, &x.StatusChangedByUserID, &x.StatusChangedByUserName,
			&x.CreatedAt, &x.UpdatedAt); err != nil { return nil,"",err }
		if sd != nil { x.StartDate = sd.Format("2006-01-02") }
		if ed != nil { x.EndDate = ed.Format("2006-01-02") }
		out = append(out, x)
	}
	next := ""
	if len(out) > p.Limit {
		last := out[p.Limit-1]
		next = EncodeCursor(last.CreatedAt, last.ID)
		out = out[:p.Limit]
	}
	return out,next,nil
}
