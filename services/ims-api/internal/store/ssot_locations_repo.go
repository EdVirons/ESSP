package store

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// SSOTLocationsRepo provides lookup methods for SSOT county/sub-county data.
type SSOTLocationsRepo struct{ pool *pgxpool.Pool }

// County represents a Kenya county with code and name.
type SSOTCounty struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// SubCounty represents a Kenya sub-county with code, name, and parent county code.
type SSOTSubCounty struct {
	Code       string `json:"code"`
	Name       string `json:"name"`
	CountyCode string `json:"countyCode"`
}

// LookupCountyCode returns the county code for a given county name.
// Returns empty string if not found.
func (r *SSOTLocationsRepo) LookupCountyCode(ctx context.Context, name string) string {
	if name == "" {
		return ""
	}
	var code string
	err := r.pool.QueryRow(ctx, `
		SELECT code FROM ssot_counties WHERE LOWER(name) = LOWER($1)
	`, name).Scan(&code)
	if err != nil {
		return ""
	}
	return code
}

// LookupSubCountyCode returns the sub-county code for a given sub-county name and county code.
// Returns empty string if not found.
func (r *SSOTLocationsRepo) LookupSubCountyCode(ctx context.Context, name, countyCode string) string {
	if name == "" || countyCode == "" {
		return ""
	}
	var code string
	err := r.pool.QueryRow(ctx, `
		SELECT code FROM ssot_sub_counties
		WHERE LOWER(name) = LOWER($1) AND county_code = $2
	`, name, countyCode).Scan(&code)
	if err != nil {
		return ""
	}
	return code
}

// LookupSubCountyCodeByName returns the sub-county code for a given sub-county name only.
// This is less precise but useful when county code is not known.
// Returns empty string if not found or multiple matches exist.
func (r *SSOTLocationsRepo) LookupSubCountyCodeByName(ctx context.Context, name string) (code, countyCode string) {
	if name == "" {
		return "", ""
	}
	err := r.pool.QueryRow(ctx, `
		SELECT code, county_code FROM ssot_sub_counties
		WHERE LOWER(name) = LOWER($1)
		LIMIT 1
	`, name).Scan(&code, &countyCode)
	if err != nil {
		return "", ""
	}
	return code, countyCode
}

// ListCounties returns all counties.
func (r *SSOTLocationsRepo) ListCounties(ctx context.Context) ([]SSOTCounty, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT code, name FROM ssot_counties ORDER BY name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var counties []SSOTCounty
	for rows.Next() {
		var c SSOTCounty
		if err := rows.Scan(&c.Code, &c.Name); err != nil {
			return nil, err
		}
		counties = append(counties, c)
	}
	return counties, nil
}

// ListSubCounties returns sub-counties, optionally filtered by county code.
func (r *SSOTLocationsRepo) ListSubCounties(ctx context.Context, countyCode string) ([]SSOTSubCounty, error) {
	var query string
	var args []any

	if countyCode != "" {
		query = `SELECT code, name, county_code FROM ssot_sub_counties WHERE county_code = $1 ORDER BY name`
		args = []any{countyCode}
	} else {
		query = `SELECT code, name, county_code FROM ssot_sub_counties ORDER BY name`
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subCounties []SSOTSubCounty
	for rows.Next() {
		var sc SSOTSubCounty
		if err := rows.Scan(&sc.Code, &sc.Name, &sc.CountyCode); err != nil {
			return nil, err
		}
		subCounties = append(subCounties, sc)
	}
	return subCounties, nil
}
