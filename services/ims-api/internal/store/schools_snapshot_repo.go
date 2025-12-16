package store

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/edvirons/ssp/ims/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SchoolsSnapshotRepo struct{ pool *pgxpool.Pool }

func (r *SchoolsSnapshotRepo) Upsert(ctx context.Context, s models.SchoolSnapshot) error {
	// Auto-lookup county/sub-county codes if missing but names are provided
	s = r.enrichWithCodes(ctx, s)

	_, err := r.pool.Exec(ctx, `
		INSERT INTO schools_snapshot (
			tenant_id, school_id, name, county_code, county_name, sub_county_code, sub_county_name,
			level, type, knec_code, uic, sex, cluster, accommodation, latitude, longitude, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)
		ON CONFLICT (tenant_id, school_id)
		DO UPDATE SET
		  name=EXCLUDED.name,
		  county_code=EXCLUDED.county_code,
		  county_name=EXCLUDED.county_name,
		  sub_county_code=EXCLUDED.sub_county_code,
		  sub_county_name=EXCLUDED.sub_county_name,
		  level=EXCLUDED.level,
		  type=EXCLUDED.type,
		  knec_code=EXCLUDED.knec_code,
		  uic=EXCLUDED.uic,
		  sex=EXCLUDED.sex,
		  cluster=EXCLUDED.cluster,
		  accommodation=EXCLUDED.accommodation,
		  latitude=EXCLUDED.latitude,
		  longitude=EXCLUDED.longitude,
		  updated_at=EXCLUDED.updated_at
	`, s.TenantID, s.SchoolID, s.Name, s.CountyCode, s.CountyName, s.SubCountyCode, s.SubCountyName,
		s.Level, s.Type, s.KnecCode, s.Uic, s.Sex, s.Cluster, s.Accommodation, s.Latitude, s.Longitude, s.UpdatedAt)
	return err
}

// enrichWithCodes looks up county/sub-county codes from SSOT tables if codes are missing but names exist.
func (r *SchoolsSnapshotRepo) enrichWithCodes(ctx context.Context, s models.SchoolSnapshot) models.SchoolSnapshot {
	// If county code is missing but name exists, look it up
	if s.CountyCode == "" && s.CountyName != "" {
		var code string
		_ = r.pool.QueryRow(ctx, `SELECT code FROM ssot_counties WHERE LOWER(name) = LOWER($1)`, s.CountyName).Scan(&code)
		if code != "" {
			s.CountyCode = code
		}
	}

	// If sub-county code is missing but name and county code exist, look it up
	if s.SubCountyCode == "" && s.SubCountyName != "" && s.CountyCode != "" {
		var code string
		_ = r.pool.QueryRow(ctx, `SELECT code FROM ssot_sub_counties WHERE LOWER(name) = LOWER($1) AND county_code = $2`, s.SubCountyName, s.CountyCode).Scan(&code)
		if code != "" {
			s.SubCountyCode = code
		}
	}

	return s
}

func (r *SchoolsSnapshotRepo) Get(ctx context.Context, tenantID, schoolID string) (models.SchoolSnapshot, error) {
	var s models.SchoolSnapshot
	row := r.pool.QueryRow(ctx, `
		SELECT tenant_id, school_id, name, county_code, county_name, sub_county_code, sub_county_name,
			level, type, knec_code, uic, sex, cluster, accommodation, latitude, longitude, updated_at
		FROM schools_snapshot
		WHERE tenant_id=$1 AND school_id=$2
	`, tenantID, schoolID)
	if err := row.Scan(&s.TenantID, &s.SchoolID, &s.Name, &s.CountyCode, &s.CountyName, &s.SubCountyCode, &s.SubCountyName,
		&s.Level, &s.Type, &s.KnecCode, &s.Uic, &s.Sex, &s.Cluster, &s.Accommodation, &s.Latitude, &s.Longitude, &s.UpdatedAt); err != nil {
		return models.SchoolSnapshot{}, errors.New("not found")
	}
	return s, nil
}

func NewSchoolSnapshot(tenantID, schoolID string) models.SchoolSnapshot {
	return models.SchoolSnapshot{TenantID: tenantID, SchoolID: schoolID, UpdatedAt: time.Now().UTC()}
}

type SchoolSnapshotListParams struct {
	TenantID   string
	Query      string
	CountyCode string
	Level      string
	Type       string
	Limit      int
	Offset     int
}

func (r *SchoolsSnapshotRepo) List(ctx context.Context, p SchoolSnapshotListParams) ([]models.SchoolSnapshot, int, error) {
	conds := []string{"tenant_id=$1"}
	args := []any{p.TenantID}
	argN := 2

	if p.Query != "" {
		conds = append(conds, "(name ILIKE $"+itoa(argN)+" OR school_id ILIKE $"+itoa(argN)+" OR knec_code ILIKE $"+itoa(argN)+")")
		args = append(args, "%"+p.Query+"%")
		argN++
	}
	if p.CountyCode != "" {
		conds = append(conds, "county_code=$"+itoa(argN))
		args = append(args, p.CountyCode)
		argN++
	}
	if p.Level != "" {
		conds = append(conds, "level=$"+itoa(argN))
		args = append(args, p.Level)
		argN++
	}
	if p.Type != "" {
		conds = append(conds, "type=$"+itoa(argN))
		args = append(args, p.Type)
		argN++
	}

	// Count total
	countSQL := "SELECT COUNT(*) FROM schools_snapshot WHERE " + strings.Join(conds, " AND ")
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
		SELECT tenant_id, school_id, name, county_code, county_name, sub_county_code, sub_county_name,
			level, type, knec_code, uic, sex, cluster, accommodation, latitude, longitude, updated_at
		FROM schools_snapshot
		WHERE ` + strings.Join(conds, " AND ") + `
		ORDER BY name ASC
		LIMIT $` + itoa(argN) + ` OFFSET $` + itoa(argN+1)

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var items []models.SchoolSnapshot
	for rows.Next() {
		var s models.SchoolSnapshot
		if err := rows.Scan(&s.TenantID, &s.SchoolID, &s.Name, &s.CountyCode, &s.CountyName, &s.SubCountyCode, &s.SubCountyName,
			&s.Level, &s.Type, &s.KnecCode, &s.Uic, &s.Sex, &s.Cluster, &s.Accommodation, &s.Latitude, &s.Longitude, &s.UpdatedAt); err != nil {
			return nil, 0, err
		}
		items = append(items, s)
	}

	return items, total, nil
}

func (r *SchoolsSnapshotRepo) Count(ctx context.Context, tenantID string) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM schools_snapshot WHERE tenant_id=$1", tenantID).Scan(&count)
	return count, err
}

// County represents a unique county from schools
type County struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// SubCounty represents a unique sub-county from schools
type SubCounty struct {
	Code       string `json:"code"`
	Name       string `json:"name"`
	CountyCode string `json:"countyCode"`
}

// ListCounties returns all distinct counties for a tenant
func (r *SchoolsSnapshotRepo) ListCounties(ctx context.Context, tenantID string) ([]County, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT DISTINCT county_code, county_name
		FROM schools_snapshot
		WHERE tenant_id=$1 AND county_code != ''
		ORDER BY county_name
	`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var counties []County
	for rows.Next() {
		var c County
		if err := rows.Scan(&c.Code, &c.Name); err != nil {
			return nil, err
		}
		counties = append(counties, c)
	}
	return counties, nil
}

// ListSubCounties returns all distinct sub-counties for a tenant, optionally filtered by county
func (r *SchoolsSnapshotRepo) ListSubCounties(ctx context.Context, tenantID, countyCode string) ([]SubCounty, error) {
	var query string
	var args []any

	if countyCode != "" {
		query = `
			SELECT DISTINCT sub_county_code, sub_county_name, county_code
			FROM schools_snapshot
			WHERE tenant_id=$1 AND county_code=$2 AND sub_county_code != ''
			ORDER BY sub_county_name
		`
		args = []any{tenantID, countyCode}
	} else {
		query = `
			SELECT DISTINCT sub_county_code, sub_county_name, county_code
			FROM schools_snapshot
			WHERE tenant_id=$1 AND sub_county_code != ''
			ORDER BY sub_county_name
		`
		args = []any{tenantID}
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subCounties []SubCounty
	for rows.Next() {
		var sc SubCounty
		if err := rows.Scan(&sc.Code, &sc.Name, &sc.CountyCode); err != nil {
			return nil, err
		}
		subCounties = append(subCounties, sc)
	}
	return subCounties, nil
}
