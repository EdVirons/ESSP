package store

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/edvirons/ssp/ims/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProjectsRepo struct{ pool *pgxpool.Pool }

func (r *ProjectsRepo) Create(ctx context.Context, p models.SchoolServiceProject) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO school_service_projects (
			id, tenant_id, school_id, project_type, status, current_phase, start_date, go_live_date,
			account_manager_user_id, notes, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,NULLIF($7,''),NULLIF($8,''),$9,$10,$11,$12)
	`, p.ID, p.TenantID, p.SchoolID, p.ProjectType, p.Status, p.CurrentPhase, p.StartDate, p.GoLiveDate, p.AccountManagerUserID, p.Notes, p.CreatedAt, p.UpdatedAt)
	return err
}

func (r *ProjectsRepo) GetByID(ctx context.Context, tenantID, id string) (models.SchoolServiceProject, error) {
	var p models.SchoolServiceProject
	var sd, gd *time.Time
	row := r.pool.QueryRow(ctx, `
		SELECT id, tenant_id, school_id, project_type, status, current_phase, start_date, go_live_date,
			account_manager_user_id, notes, created_at, updated_at
		FROM school_service_projects
		WHERE tenant_id=$1 AND id=$2
	`, tenantID, id)
	if err := row.Scan(&p.ID, &p.TenantID, &p.SchoolID, &p.ProjectType, &p.Status, &p.CurrentPhase, &sd, &gd, &p.AccountManagerUserID, &p.Notes, &p.CreatedAt, &p.UpdatedAt); err != nil {
		return models.SchoolServiceProject{}, errors.New("not found")
	}
	if sd != nil {
		p.StartDate = sd.Format("2006-01-02")
	}
	if gd != nil {
		p.GoLiveDate = gd.Format("2006-01-02")
	}
	return p, nil
}

type ProjectListParams struct {
	TenantID    string
	SchoolID    string
	ProjectType string
	Status      string
	Limit       int
	HasCursor   bool
	CursorCreatedAt time.Time
	CursorID    string
}

func (r *ProjectsRepo) List(ctx context.Context, p ProjectListParams) ([]models.SchoolServiceProject, string, error) {
	conds := []string{"tenant_id=$1"}
	args := []any{p.TenantID}
	argN := 2
	if p.SchoolID != "" {
		conds = append(conds, "school_id=$"+itoa(argN))
		args = append(args, p.SchoolID)
		argN++
	}
	if p.ProjectType != "" {
		conds = append(conds, "project_type=$"+itoa(argN))
		args = append(args, p.ProjectType)
		argN++
	}
	if p.Status != "" {
		conds = append(conds, "status=$"+itoa(argN))
		args = append(args, p.Status)
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
		SELECT id, tenant_id, school_id, project_type, status, current_phase, start_date, go_live_date,
			account_manager_user_id, notes, created_at, updated_at
		FROM school_service_projects
		WHERE ` + strings.Join(conds, " AND ") + `
		ORDER BY created_at DESC, id DESC
		LIMIT $` + itoa(argN)

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	out := []models.SchoolServiceProject{}
	for rows.Next() {
		var x models.SchoolServiceProject
		var sd, gd *time.Time
		if err := rows.Scan(&x.ID, &x.TenantID, &x.SchoolID, &x.ProjectType, &x.Status, &x.CurrentPhase, &sd, &gd, &x.AccountManagerUserID, &x.Notes, &x.CreatedAt, &x.UpdatedAt); err != nil {
			return nil, "", err
		}
		if sd != nil {
			x.StartDate = sd.Format("2006-01-02")
		}
		if gd != nil {
			x.GoLiveDate = gd.Format("2006-01-02")
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

// CountByType returns the count of projects grouped by project type.
func (r *ProjectsRepo) CountByType(ctx context.Context, tenantID string) (map[string]int, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT project_type, COUNT(*) as count
		FROM school_service_projects
		WHERE tenant_id = $1
		GROUP BY project_type
	`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	counts := make(map[string]int)
	for rows.Next() {
		var pt string
		var count int
		if err := rows.Scan(&pt, &count); err != nil {
			return nil, err
		}
		counts[pt] = count
	}
	return counts, nil
}
