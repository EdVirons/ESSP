package store

import (
	"context"
	"errors"
	"time"

	"github.com/edvirons/ssp/ims/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SchoolRepo struct{ pool *pgxpool.Pool }

// Upsert stores county mapping for a school (from School SSOT sync).
func (r *SchoolRepo) Upsert(ctx context.Context, s models.SchoolProfile) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO schools (tenant_id, school_id, county_code, county_name, updated_at)
		VALUES ($1,$2,$3,$4,$5)
		ON CONFLICT (tenant_id, school_id)
		DO UPDATE SET county_code=EXCLUDED.county_code, county_name=EXCLUDED.county_name, updated_at=EXCLUDED.updated_at
	`, s.TenantID, s.SchoolID, s.CountyCode, s.CountyName, s.UpdatedAt)
	return err
}

func (r *SchoolRepo) Get(ctx context.Context, tenantID, schoolID string) (models.SchoolProfile, error) {
	var s models.SchoolProfile
	row := r.pool.QueryRow(ctx, `
		SELECT tenant_id, school_id, county_code, county_name, updated_at
		FROM schools
		WHERE tenant_id=$1 AND school_id=$2
	`, tenantID, schoolID)
	if err := row.Scan(&s.TenantID, &s.SchoolID, &s.CountyCode, &s.CountyName, &s.UpdatedAt); err != nil {
		return models.SchoolProfile{}, errors.New("not found")
	}
	return s, nil
}

func NewSchoolProfile(tenantID, schoolID, countyCode, countyName string) models.SchoolProfile {
	return models.SchoolProfile{
		TenantID: tenantID, SchoolID: schoolID,
		CountyCode: countyCode, CountyName: countyName,
		UpdatedAt: time.Now().UTC(),
	}
}
