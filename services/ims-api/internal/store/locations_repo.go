package store

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/edvirons/ssp/ims/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LocationsRepo struct{ pool *pgxpool.Pool }

func (r *LocationsRepo) Create(ctx context.Context, loc models.Location) error {
	metadata := loc.Metadata
	if len(metadata) == 0 {
		metadata = []byte("{}")
	}
	_, err := r.pool.Exec(ctx, `
		INSERT INTO locations (id, tenant_id, school_id, parent_id, name, location_type, code, capacity, metadata, active, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
	`, loc.ID, loc.TenantID, loc.SchoolID, loc.ParentID, loc.Name, loc.LocationType, loc.Code, loc.Capacity, metadata, loc.Active, loc.CreatedAt, loc.UpdatedAt)
	return err
}

func (r *LocationsRepo) Update(ctx context.Context, loc models.Location) error {
	metadata := loc.Metadata
	if len(metadata) == 0 {
		metadata = []byte("{}")
	}
	result, err := r.pool.Exec(ctx, `
		UPDATE locations SET
			parent_id=$4, name=$5, location_type=$6, code=$7, capacity=$8, metadata=$9, active=$10, updated_at=$11
		WHERE tenant_id=$1 AND school_id=$2 AND id=$3
	`, loc.TenantID, loc.SchoolID, loc.ID, loc.ParentID, loc.Name, loc.LocationType, loc.Code, loc.Capacity, metadata, loc.Active, loc.UpdatedAt)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errors.New("not found")
	}
	return nil
}

func (r *LocationsRepo) Get(ctx context.Context, tenantID, id string) (models.Location, error) {
	var loc models.Location
	var metadata []byte
	row := r.pool.QueryRow(ctx, `
		SELECT id, tenant_id, school_id, parent_id, name, location_type, code, capacity, metadata, active, created_at, updated_at
		FROM locations WHERE tenant_id=$1 AND id=$2
	`, tenantID, id)
	if err := row.Scan(&loc.ID, &loc.TenantID, &loc.SchoolID, &loc.ParentID, &loc.Name, &loc.LocationType, &loc.Code, &loc.Capacity, &metadata, &loc.Active, &loc.CreatedAt, &loc.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Location{}, errors.New("not found")
		}
		return models.Location{}, err
	}
	loc.Metadata = metadata
	return loc, nil
}

func (r *LocationsRepo) Delete(ctx context.Context, tenantID, id string) error {
	result, err := r.pool.Exec(ctx, `DELETE FROM locations WHERE tenant_id=$1 AND id=$2`, tenantID, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errors.New("not found")
	}
	return nil
}

type LocationListParams struct {
	TenantID  string
	SchoolID  string
	ParentID  *string // nil = root locations, empty string = all
	Active    *bool
	Type      string
	Limit     int
	HasCursor bool
	CursorAt  time.Time
	CursorID  string
}

func (r *LocationsRepo) List(ctx context.Context, p LocationListParams) ([]models.Location, string, error) {
	conds := []string{"tenant_id=$1"}
	args := []any{p.TenantID}
	argN := 2

	if p.SchoolID != "" {
		conds = append(conds, "school_id=$"+itoa(argN))
		args = append(args, p.SchoolID)
		argN++
	}
	if p.ParentID != nil {
		if *p.ParentID == "" {
			conds = append(conds, "parent_id IS NULL")
		} else {
			conds = append(conds, "parent_id=$"+itoa(argN))
			args = append(args, *p.ParentID)
			argN++
		}
	}
	if p.Active != nil {
		conds = append(conds, "active=$"+itoa(argN))
		args = append(args, *p.Active)
		argN++
	}
	if p.Type != "" {
		conds = append(conds, "location_type=$"+itoa(argN))
		args = append(args, p.Type)
		argN++
	}
	if p.HasCursor {
		conds = append(conds, "(created_at, id) < ($"+itoa(argN)+", $"+itoa(argN+1)+")")
		args = append(args, p.CursorAt, p.CursorID)
		argN += 2
	}

	limit := p.Limit
	if limit <= 0 {
		limit = 50
	}
	limitPlus := limit + 1
	args = append(args, limitPlus)

	sql := `
		SELECT id, tenant_id, school_id, parent_id, name, location_type, code, capacity, metadata, active, created_at, updated_at
		FROM locations
		WHERE ` + strings.Join(conds, " AND ") + `
		ORDER BY created_at DESC, id DESC
		LIMIT $` + itoa(argN)

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	out := []models.Location{}
	for rows.Next() {
		var x models.Location
		var metadata []byte
		if err := rows.Scan(&x.ID, &x.TenantID, &x.SchoolID, &x.ParentID, &x.Name, &x.LocationType, &x.Code, &x.Capacity, &metadata, &x.Active, &x.CreatedAt, &x.UpdatedAt); err != nil {
			return nil, "", err
		}
		x.Metadata = metadata
		out = append(out, x)
	}

	next := ""
	if len(out) > limit {
		last := out[limit-1]
		next = EncodeCursor(last.CreatedAt, last.ID)
		out = out[:limit]
	}
	return out, next, nil
}

// ListBySchool returns all active locations for a school (no pagination, for inventory view)
func (r *LocationsRepo) ListBySchool(ctx context.Context, tenantID, schoolID string) ([]models.Location, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, tenant_id, school_id, parent_id, name, location_type, code, capacity, metadata, active, created_at, updated_at
		FROM locations
		WHERE tenant_id=$1 AND school_id=$2 AND active=true
		ORDER BY name
	`, tenantID, schoolID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []models.Location{}
	for rows.Next() {
		var x models.Location
		var metadata []byte
		if err := rows.Scan(&x.ID, &x.TenantID, &x.SchoolID, &x.ParentID, &x.Name, &x.LocationType, &x.Code, &x.Capacity, &metadata, &x.Active, &x.CreatedAt, &x.UpdatedAt); err != nil {
			return nil, err
		}
		x.Metadata = metadata
		out = append(out, x)
	}
	return out, nil
}

// GetPath returns the hierarchical path for a location (e.g., "Block A > Floor 1 > Lab 101")
func (r *LocationsRepo) GetPath(ctx context.Context, tenantID, locationID string) (string, error) {
	var path string
	err := r.pool.QueryRow(ctx, `SELECT get_location_path($1)`, locationID).Scan(&path)
	if err != nil {
		return "", err
	}
	return path, nil
}

// GetTree returns locations as a tree structure starting from root
func (r *LocationsRepo) GetTree(ctx context.Context, tenantID, schoolID string) ([]models.LocationTreeNode, error) {
	locs, err := r.ListBySchool(ctx, tenantID, schoolID)
	if err != nil {
		return nil, err
	}

	// Build tree in memory
	byID := make(map[string]*models.LocationTreeNode)
	for i := range locs {
		byID[locs[i].ID] = &models.LocationTreeNode{Location: locs[i]}
	}

	var roots []models.LocationTreeNode
	for _, node := range byID {
		if node.ParentID == nil || *node.ParentID == "" {
			roots = append(roots, *node)
		} else if parent, ok := byID[*node.ParentID]; ok {
			parent.Children = append(parent.Children, *node)
		}
	}
	return roots, nil
}

// CountDevicesByLocation returns device counts per location for a school
func (r *LocationsRepo) CountDevicesByLocation(ctx context.Context, tenantID, schoolID string) (map[string]int, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT l.id, COUNT(da.device_id)::int
		FROM locations l
		LEFT JOIN device_assignments da ON da.location_id = l.id AND da.effective_to IS NULL
		WHERE l.tenant_id=$1 AND l.school_id=$2 AND l.active=true
		GROUP BY l.id
	`, tenantID, schoolID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	counts := make(map[string]int)
	for rows.Next() {
		var id string
		var count int
		if err := rows.Scan(&id, &count); err != nil {
			return nil, err
		}
		counts[id] = count
	}
	return counts, nil
}

// GetMetadata parses the metadata JSON for a location
func (r *LocationsRepo) GetMetadata(loc models.Location) (map[string]any, error) {
	if len(loc.Metadata) == 0 {
		return make(map[string]any), nil
	}
	var m map[string]any
	if err := json.Unmarshal(loc.Metadata, &m); err != nil {
		return nil, err
	}
	return m, nil
}
