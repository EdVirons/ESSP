package store

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/edvirons/ssp/ims/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lib/pq"
)

type ProjectTeamRepo struct{ pool *pgxpool.Pool }

// AddMember adds a team member to a project.
func (r *ProjectTeamRepo) AddMember(ctx context.Context, m models.ProjectTeamMember) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO project_team_members (
			id, tenant_id, project_id, user_id, user_email, user_name,
			role, assigned_phases, responsibility, assigned_by_user_id,
			assigned_at, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
	`, m.ID, m.TenantID, m.ProjectID, m.UserID, m.UserEmail, m.UserName,
		m.Role, pq.Array(m.AssignedPhases), m.Responsibility, m.AssignedByUserID,
		m.AssignedAt, m.CreatedAt, m.UpdatedAt)
	return err
}

// UpdateMember updates a team member's role, phases, or responsibility.
func (r *ProjectTeamRepo) UpdateMember(ctx context.Context, tenantID, memberID string, role string, phases []string, responsibility string) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE project_team_members
		SET role = $3, assigned_phases = $4, responsibility = $5, updated_at = $6
		WHERE tenant_id = $1 AND id = $2 AND removed_at IS NULL
	`, tenantID, memberID, role, pq.Array(phases), responsibility, time.Now())
	return err
}

// RemoveMember soft-deletes a team member.
func (r *ProjectTeamRepo) RemoveMember(ctx context.Context, tenantID, memberID string) error {
	now := time.Now()
	_, err := r.pool.Exec(ctx, `
		UPDATE project_team_members
		SET removed_at = $3, updated_at = $3
		WHERE tenant_id = $1 AND id = $2 AND removed_at IS NULL
	`, tenantID, memberID, now)
	return err
}

// GetMember retrieves a specific team member.
func (r *ProjectTeamRepo) GetMember(ctx context.Context, tenantID, memberID string) (models.ProjectTeamMember, error) {
	var m models.ProjectTeamMember
	var phases []string
	row := r.pool.QueryRow(ctx, `
		SELECT id, tenant_id, project_id, user_id, user_email, user_name,
			role, assigned_phases, responsibility, assigned_by_user_id,
			assigned_at, removed_at, created_at, updated_at
		FROM project_team_members
		WHERE tenant_id = $1 AND id = $2
	`, tenantID, memberID)
	if err := row.Scan(&m.ID, &m.TenantID, &m.ProjectID, &m.UserID, &m.UserEmail, &m.UserName,
		&m.Role, pq.Array(&phases), &m.Responsibility, &m.AssignedByUserID,
		&m.AssignedAt, &m.RemovedAt, &m.CreatedAt, &m.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.ProjectTeamMember{}, errors.New("not found")
		}
		return models.ProjectTeamMember{}, err
	}
	m.AssignedPhases = make([]models.PhaseType, len(phases))
	for i, p := range phases {
		m.AssignedPhases[i] = models.PhaseType(p)
	}
	return m, nil
}

// ListByProject lists all active team members for a project.
func (r *ProjectTeamRepo) ListByProject(ctx context.Context, tenantID, projectID string) ([]models.ProjectTeamMember, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, tenant_id, project_id, user_id, user_email, user_name,
			role, assigned_phases, responsibility, assigned_by_user_id,
			assigned_at, removed_at, created_at, updated_at
		FROM project_team_members
		WHERE tenant_id = $1 AND project_id = $2 AND removed_at IS NULL
		ORDER BY assigned_at ASC
	`, tenantID, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []models.ProjectTeamMember
	for rows.Next() {
		var m models.ProjectTeamMember
		var phases []string
		if err := rows.Scan(&m.ID, &m.TenantID, &m.ProjectID, &m.UserID, &m.UserEmail, &m.UserName,
			&m.Role, pq.Array(&phases), &m.Responsibility, &m.AssignedByUserID,
			&m.AssignedAt, &m.RemovedAt, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		m.AssignedPhases = make([]models.PhaseType, len(phases))
		for i, p := range phases {
			m.AssignedPhases[i] = models.PhaseType(p)
		}
		members = append(members, m)
	}
	return members, nil
}

// ListByUser lists all projects a user is assigned to.
func (r *ProjectTeamRepo) ListByUser(ctx context.Context, tenantID, userID string) ([]models.ProjectTeamMember, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, tenant_id, project_id, user_id, user_email, user_name,
			role, assigned_phases, responsibility, assigned_by_user_id,
			assigned_at, removed_at, created_at, updated_at
		FROM project_team_members
		WHERE tenant_id = $1 AND user_id = $2 AND removed_at IS NULL
		ORDER BY assigned_at DESC
	`, tenantID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []models.ProjectTeamMember
	for rows.Next() {
		var m models.ProjectTeamMember
		var phases []string
		if err := rows.Scan(&m.ID, &m.TenantID, &m.ProjectID, &m.UserID, &m.UserEmail, &m.UserName,
			&m.Role, pq.Array(&phases), &m.Responsibility, &m.AssignedByUserID,
			&m.AssignedAt, &m.RemovedAt, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		m.AssignedPhases = make([]models.PhaseType, len(phases))
		for i, p := range phases {
			m.AssignedPhases[i] = models.PhaseType(p)
		}
		members = append(members, m)
	}
	return members, nil
}

// IsMember checks if a user is an active member of a project.
func (r *ProjectTeamRepo) IsMember(ctx context.Context, tenantID, projectID, userID string) (bool, error) {
	var count int
	err := r.pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM project_team_members
		WHERE tenant_id = $1 AND project_id = $2 AND user_id = $3 AND removed_at IS NULL
	`, tenantID, projectID, userID).Scan(&count)
	return count > 0, err
}

// PhaseUserAssignments

// AddPhaseAssignment adds a user assignment to a phase.
func (r *ProjectTeamRepo) AddPhaseAssignment(ctx context.Context, a models.PhaseUserAssignment) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO phase_user_assignments (
			id, tenant_id, phase_id, project_id, user_id, user_email, user_name,
			assignment_type, assigned_by_user_id, assigned_at, created_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
	`, a.ID, a.TenantID, a.PhaseID, a.ProjectID, a.UserID, a.UserEmail, a.UserName,
		a.AssignmentType, a.AssignedByUserID, a.AssignedAt, a.CreatedAt)
	return err
}

// RemovePhaseAssignment soft-deletes a phase assignment.
func (r *ProjectTeamRepo) RemovePhaseAssignment(ctx context.Context, tenantID, assignmentID string) error {
	now := time.Now()
	_, err := r.pool.Exec(ctx, `
		UPDATE phase_user_assignments
		SET removed_at = $3
		WHERE tenant_id = $1 AND id = $2 AND removed_at IS NULL
	`, tenantID, assignmentID, now)
	return err
}

// ListPhaseAssignments lists all active assignments for a phase.
func (r *ProjectTeamRepo) ListPhaseAssignments(ctx context.Context, tenantID, phaseID string) ([]models.PhaseUserAssignment, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, tenant_id, phase_id, project_id, user_id, user_email, user_name,
			assignment_type, assigned_by_user_id, assigned_at, completed_at, removed_at, created_at
		FROM phase_user_assignments
		WHERE tenant_id = $1 AND phase_id = $2 AND removed_at IS NULL
		ORDER BY created_at ASC
	`, tenantID, phaseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assignments []models.PhaseUserAssignment
	for rows.Next() {
		var a models.PhaseUserAssignment
		if err := rows.Scan(&a.ID, &a.TenantID, &a.PhaseID, &a.ProjectID, &a.UserID, &a.UserEmail, &a.UserName,
			&a.AssignmentType, &a.AssignedByUserID, &a.AssignedAt, &a.CompletedAt, &a.RemovedAt, &a.CreatedAt); err != nil {
			return nil, err
		}
		assignments = append(assignments, a)
	}
	return assignments, nil
}

// ListPhaseAssignmentsByUser lists all phase assignments for a user in a project.
func (r *ProjectTeamRepo) ListPhaseAssignmentsByUser(ctx context.Context, tenantID, projectID, userID string) ([]models.PhaseUserAssignment, error) {
	conds := []string{"tenant_id=$1", "user_id=$2", "removed_at IS NULL"}
	args := []any{tenantID, userID}
	if projectID != "" {
		conds = append(conds, "project_id=$3")
		args = append(args, projectID)
	}

	rows, err := r.pool.Query(ctx, `
		SELECT id, tenant_id, phase_id, project_id, user_id, user_email, user_name,
			assignment_type, assigned_by_user_id, assigned_at, completed_at, removed_at, created_at
		FROM phase_user_assignments
		WHERE `+strings.Join(conds, " AND ")+`
		ORDER BY created_at DESC
	`, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assignments []models.PhaseUserAssignment
	for rows.Next() {
		var a models.PhaseUserAssignment
		if err := rows.Scan(&a.ID, &a.TenantID, &a.PhaseID, &a.ProjectID, &a.UserID, &a.UserEmail, &a.UserName,
			&a.AssignmentType, &a.AssignedByUserID, &a.AssignedAt, &a.CompletedAt, &a.RemovedAt, &a.CreatedAt); err != nil {
			return nil, err
		}
		assignments = append(assignments, a)
	}
	return assignments, nil
}
