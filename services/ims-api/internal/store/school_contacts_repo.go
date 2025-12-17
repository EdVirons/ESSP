package store

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/edvirons/ssp/ims/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SchoolContactsRepo struct{ pool *pgxpool.Pool }

func (r *SchoolContactsRepo) Create(ctx context.Context, c models.SchoolContact) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO school_contacts (
			id, tenant_id, school_id, user_id, name, phone, email, role, is_primary, active, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
	`, c.ID,c.TenantID,c.SchoolID,c.UserID,c.Name,c.Phone,c.Email,c.Role,c.IsPrimary,c.Active,c.CreatedAt,c.UpdatedAt)
	return err
}

func (r *SchoolContactsRepo) List(ctx context.Context, tenantID, schoolID string) ([]models.SchoolContact, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, tenant_id, school_id, user_id, name, phone, email, role, is_primary, active, created_at, updated_at
		FROM school_contacts
		WHERE tenant_id=$1 AND school_id=$2
		ORDER BY is_primary DESC, active DESC, created_at DESC
	`, tenantID, schoolID)
	if err != nil { return nil, err }
	defer rows.Close()
	out := []models.SchoolContact{}
	for rows.Next() {
		var x models.SchoolContact
		if err := rows.Scan(&x.ID,&x.TenantID,&x.SchoolID,&x.UserID,&x.Name,&x.Phone,&x.Email,&x.Role,&x.IsPrimary,&x.Active,&x.CreatedAt,&x.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, x)
	}
	return out, nil
}

func (r *SchoolContactsRepo) GetPrimary(ctx context.Context, tenantID, schoolID string) (models.SchoolContact, error) {
	var x models.SchoolContact
	row := r.pool.QueryRow(ctx, `
		SELECT id, tenant_id, school_id, user_id, name, phone, email, role, is_primary, active, created_at, updated_at
		FROM school_contacts
		WHERE tenant_id=$1 AND school_id=$2 AND is_primary=true AND active=true
		ORDER BY updated_at DESC
		LIMIT 1
	`, tenantID, schoolID)
	if err := row.Scan(&x.ID,&x.TenantID,&x.SchoolID,&x.UserID,&x.Name,&x.Phone,&x.Email,&x.Role,&x.IsPrimary,&x.Active,&x.CreatedAt,&x.UpdatedAt); err != nil {
		return models.SchoolContact{}, errors.New("not found")
	}
	return x, nil
}

func (r *SchoolContactsRepo) SetPrimary(ctx context.Context, tenantID, schoolID, contactID string) error {
	now := time.Now().UTC()
	_, err := r.pool.Exec(ctx, `
		UPDATE school_contacts SET is_primary=false, updated_at=$3
		WHERE tenant_id=$1 AND school_id=$2
	`, tenantID, schoolID, now)
	if err != nil { return err }
	_, err = r.pool.Exec(ctx, `
		UPDATE school_contacts SET is_primary=true, updated_at=$4
		WHERE tenant_id=$1 AND school_id=$2 AND id=$3
	`, tenantID, schoolID, contactID, now)
	return err
}

func (r *SchoolContactsRepo) NormalizeRole(role string) string {
	role = strings.TrimSpace(role)
	if role == "" { role = "point_of_contact" }
	return role
}

// GetByUserID retrieves a school contact by user ID
func (r *SchoolContactsRepo) GetByUserID(ctx context.Context, tenantID, userID string) (models.SchoolContact, error) {
	var x models.SchoolContact
	row := r.pool.QueryRow(ctx, `
		SELECT id, tenant_id, school_id, user_id, name, phone, email, role, is_primary, active, created_at, updated_at
		FROM school_contacts
		WHERE tenant_id=$1 AND user_id=$2 AND active=true
		ORDER BY updated_at DESC
		LIMIT 1
	`, tenantID, userID)
	if err := row.Scan(&x.ID,&x.TenantID,&x.SchoolID,&x.UserID,&x.Name,&x.Phone,&x.Email,&x.Role,&x.IsPrimary,&x.Active,&x.CreatedAt,&x.UpdatedAt); err != nil {
		return models.SchoolContact{}, errors.New("not found")
	}
	return x, nil
}

// ListSchoolsByUserID returns all school IDs associated with a user
func (r *SchoolContactsRepo) ListSchoolsByUserID(ctx context.Context, tenantID, userID string) ([]string, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT DISTINCT school_id
		FROM school_contacts
		WHERE tenant_id=$1 AND user_id=$2 AND active=true
	`, tenantID, userID)
	if err != nil { return nil, err }
	defer rows.Close()
	out := []string{}
	for rows.Next() {
		var schoolID string
		if err := rows.Scan(&schoolID); err != nil {
			return nil, err
		}
		out = append(out, schoolID)
	}
	return out, nil
}

// ImpersonatableUser represents a user that can be impersonated
type ImpersonatableUser struct {
	UserID   string   `json:"userId"`
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Schools  []string `json:"schools"`
}

// ListImpersonatableUsers returns all school contacts that can be impersonated
func (r *SchoolContactsRepo) ListImpersonatableUsers(ctx context.Context, tenantID string, limit int) ([]ImpersonatableUser, error) {
	if limit <= 0 {
		limit = 50
	}
	rows, err := r.pool.Query(ctx, `
		SELECT user_id, name, email, array_agg(DISTINCT school_id) as schools
		FROM school_contacts
		WHERE tenant_id=$1 AND user_id IS NOT NULL AND user_id != '' AND active=true
		GROUP BY user_id, name, email
		ORDER BY name
		LIMIT $2
	`, tenantID, limit)
	if err != nil { return nil, err }
	defer rows.Close()
	out := []ImpersonatableUser{}
	for rows.Next() {
		var u ImpersonatableUser
		if err := rows.Scan(&u.UserID, &u.Name, &u.Email, &u.Schools); err != nil {
			return nil, err
		}
		out = append(out, u)
	}
	return out, nil
}
