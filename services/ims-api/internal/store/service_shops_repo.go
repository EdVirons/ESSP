package store

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/edvirons/ssp/ims/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ServiceShopRepo struct{ pool *pgxpool.Pool }

func (r *ServiceShopRepo) Create(ctx context.Context, s models.ServiceShop) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO service_shops (id, tenant_id, county_code, county_name, sub_county_code, sub_county_name, coverage_level, name, location, active, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
	`, s.ID, s.TenantID, s.CountyCode, s.CountyName, s.SubCountyCode, s.SubCountyName, s.CoverageLevel, s.Name, s.Location, s.Active, s.CreatedAt, s.UpdatedAt)
	return err
}

func (r *ServiceShopRepo) GetByID(ctx context.Context, tenantID, id string) (models.ServiceShop, error) {
	var s models.ServiceShop
	row := r.pool.QueryRow(ctx, `
		SELECT id, tenant_id, county_code, county_name, sub_county_code, sub_county_name, coverage_level, name, location, active, created_at, updated_at
		FROM service_shops
		WHERE tenant_id=$1 AND id=$2
	`, tenantID, id)
	if err := row.Scan(&s.ID, &s.TenantID, &s.CountyCode, &s.CountyName, &s.SubCountyCode, &s.SubCountyName, &s.CoverageLevel, &s.Name, &s.Location, &s.Active, &s.CreatedAt, &s.UpdatedAt); err != nil {
		return models.ServiceShop{}, errors.New("not found")
	}
	return s, nil
}

func (r *ServiceShopRepo) GetByCounty(ctx context.Context, tenantID, countyCode string) (models.ServiceShop, error) {
	var s models.ServiceShop
	row := r.pool.QueryRow(ctx, `
		SELECT id, tenant_id, county_code, county_name, sub_county_code, sub_county_name, coverage_level, name, location, active, created_at, updated_at
		FROM service_shops
		WHERE tenant_id=$1 AND county_code=$2 AND active=true
	`, tenantID, countyCode)
	if err := row.Scan(&s.ID, &s.TenantID, &s.CountyCode, &s.CountyName, &s.SubCountyCode, &s.SubCountyName, &s.CoverageLevel, &s.Name, &s.Location, &s.Active, &s.CreatedAt, &s.UpdatedAt); err != nil {
		return models.ServiceShop{}, errors.New("not found")
	}
	return s, nil
}

type ShopListParams struct {
	TenantID string
	CountyCode string
	ActiveOnly bool
	Limit int
	HasCursor bool
	CursorCreatedAt time.Time
	CursorID string
}

func (r *ServiceShopRepo) List(ctx context.Context, p ShopListParams) ([]models.ServiceShop, string, error) {
	conds := []string{"tenant_id=$1"}
	args := []any{p.TenantID}
	argN := 2
	if p.CountyCode != "" {
		conds = append(conds, "county_code=$"+itoa(argN))
		args = append(args, p.CountyCode); argN++
	}
	if p.ActiveOnly {
		conds = append(conds, "active=true")
	}
	if p.HasCursor {
		conds = append(conds, "(created_at, id) < ($"+itoa(argN)+", $"+itoa(argN+1)+")")
		args = append(args, p.CursorCreatedAt, p.CursorID); argN += 2
	}
	limitPlus := p.Limit+1
	args = append(args, limitPlus)

	sql := `
		SELECT id, tenant_id, county_code, county_name, sub_county_code, sub_county_name, coverage_level, name, location, active, created_at, updated_at
		FROM service_shops
		WHERE ` + strings.Join(conds, " AND ") + `
		ORDER BY created_at DESC, id DESC
		LIMIT $` + itoa(argN)

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil { return nil,"",err }
	defer rows.Close()
	out := []models.ServiceShop{}
	for rows.Next() {
		var s models.ServiceShop
		if err := rows.Scan(&s.ID,&s.TenantID,&s.CountyCode,&s.CountyName,&s.SubCountyCode,&s.SubCountyName,&s.CoverageLevel,&s.Name,&s.Location,&s.Active,&s.CreatedAt,&s.UpdatedAt); err != nil {
			return nil,"",err
		}
		out = append(out, s)
	}
	next := ""
	if len(out) > p.Limit {
		last := out[p.Limit-1]
		next = EncodeCursor(last.CreatedAt, last.ID)
		out = out[:p.Limit]
	}
	return out, next, nil
}


func (r *ServiceShopRepo) GetBySubCounty(ctx context.Context, tenantID, countyCode, subCountyCode string) (models.ServiceShop, error) {
	var s models.ServiceShop
	row := r.pool.QueryRow(ctx, `
		SELECT id, tenant_id, county_code, county_name, sub_county_code, sub_county_name, coverage_level,
		       name, location, active, created_at, updated_at
		FROM service_shops
		WHERE tenant_id=$1 AND county_code=$2 AND sub_county_code=$3 AND active=true
		  AND coverage_level='sub_county'
		LIMIT 1
	`, tenantID, countyCode, subCountyCode)
	if err := row.Scan(&s.ID, &s.TenantID, &s.CountyCode, &s.CountyName, &s.SubCountyCode, &s.SubCountyName, &s.CoverageLevel,
		&s.Name, &s.Location, &s.Active, &s.CreatedAt, &s.UpdatedAt); err != nil {
		return models.ServiceShop{}, errors.New("not found")
	}
	return s, nil
}
