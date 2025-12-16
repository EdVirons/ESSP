package store

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/edvirons/ssp/ims/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ServiceStaffRepo struct{ pool *pgxpool.Pool }

func (r *ServiceStaffRepo) Create(ctx context.Context, s models.ServiceStaff) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO service_staff (id, tenant_id, service_shop_id, user_id, role, phone, active, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
	`, s.ID, s.TenantID, s.ServiceShopID, s.UserID, s.Role, s.Phone, s.Active, s.CreatedAt, s.UpdatedAt)
	return err
}

func (r *ServiceStaffRepo) GetByID(ctx context.Context, tenantID, id string) (models.ServiceStaff, error) {
	var s models.ServiceStaff
	row := r.pool.QueryRow(ctx, `
		SELECT id, tenant_id, service_shop_id, user_id, role, phone, active, created_at, updated_at
		FROM service_staff
		WHERE tenant_id=$1 AND id=$2
	`, tenantID, id)
	if err := row.Scan(&s.ID,&s.TenantID,&s.ServiceShopID,&s.UserID,&s.Role,&s.Phone,&s.Active,&s.CreatedAt,&s.UpdatedAt); err != nil {
		return models.ServiceStaff{}, errors.New("not found")
	}
	return s, nil
}

func (r *ServiceStaffRepo) GetLeadByShop(ctx context.Context, tenantID, shopID string) (models.ServiceStaff, error) {
	var s models.ServiceStaff
	row := r.pool.QueryRow(ctx, `
		SELECT id, tenant_id, service_shop_id, user_id, role, phone, active, created_at, updated_at
		FROM service_staff
		WHERE tenant_id=$1 AND service_shop_id=$2 AND role='lead_technician' AND active=true
		ORDER BY created_at ASC
		LIMIT 1
	`, tenantID, shopID)
	if err := row.Scan(&s.ID,&s.TenantID,&s.ServiceShopID,&s.UserID,&s.Role,&s.Phone,&s.Active,&s.CreatedAt,&s.UpdatedAt); err != nil {
		return models.ServiceStaff{}, errors.New("not found")
	}
	return s, nil
}

type StaffListParams struct {
	TenantID string
	ShopID string
	Role string
	ActiveOnly bool
	Limit int
	HasCursor bool
	CursorCreatedAt time.Time
	CursorID string
}

func (r *ServiceStaffRepo) List(ctx context.Context, p StaffListParams) ([]models.ServiceStaff, string, error) {
	conds := []string{"tenant_id=$1"}
	args := []any{p.TenantID}
	argN := 2
	if p.ShopID != "" {
		conds = append(conds, "service_shop_id=$"+itoa(argN)); args = append(args, p.ShopID); argN++
	}
	if p.Role != "" {
		conds = append(conds, "role=$"+itoa(argN)); args = append(args, p.Role); argN++
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
		SELECT id, tenant_id, service_shop_id, user_id, role, phone, active, created_at, updated_at
		FROM service_staff
		WHERE ` + strings.Join(conds, " AND ") + `
		ORDER BY created_at DESC, id DESC
		LIMIT $` + itoa(argN)

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil { return nil,"",err }
	defer rows.Close()
	out := []models.ServiceStaff{}
	for rows.Next() {
		var s models.ServiceStaff
		if err := rows.Scan(&s.ID,&s.TenantID,&s.ServiceShopID,&s.UserID,&s.Role,&s.Phone,&s.Active,&s.CreatedAt,&s.UpdatedAt); err != nil {
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
	return out,next,nil
}
