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

type GroupsRepo struct{ pool *pgxpool.Pool }

func (r *GroupsRepo) Create(ctx context.Context, g models.DeviceGroup) error {
	policies := g.Policies
	if len(policies) == 0 {
		policies = []byte("{}")
	}
	var selectorJSON []byte
	if g.Selector != nil {
		var err error
		selectorJSON, err = json.Marshal(g.Selector)
		if err != nil {
			return err
		}
	}
	_, err := r.pool.Exec(ctx, `
		INSERT INTO device_groups (id, tenant_id, school_id, name, description, group_type, location_id, selector, policies, active, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
	`, g.ID, g.TenantID, g.SchoolID, g.Name, g.Description, g.GroupType, g.LocationID, selectorJSON, policies, g.Active, g.CreatedAt, g.UpdatedAt)
	return err
}

func (r *GroupsRepo) Update(ctx context.Context, g models.DeviceGroup) error {
	policies := g.Policies
	if len(policies) == 0 {
		policies = []byte("{}")
	}
	var selectorJSON []byte
	if g.Selector != nil {
		var err error
		selectorJSON, err = json.Marshal(g.Selector)
		if err != nil {
			return err
		}
	}
	result, err := r.pool.Exec(ctx, `
		UPDATE device_groups SET
			school_id=$3, name=$4, description=$5, group_type=$6, location_id=$7, selector=$8, policies=$9, active=$10, updated_at=$11
		WHERE tenant_id=$1 AND id=$2
	`, g.TenantID, g.ID, g.SchoolID, g.Name, g.Description, g.GroupType, g.LocationID, selectorJSON, policies, g.Active, g.UpdatedAt)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errors.New("not found")
	}
	return nil
}

func (r *GroupsRepo) Get(ctx context.Context, tenantID, id string) (models.DeviceGroup, error) {
	var g models.DeviceGroup
	var policies, selectorJSON []byte
	row := r.pool.QueryRow(ctx, `
		SELECT id, tenant_id, school_id, name, description, group_type, location_id, selector, policies, active, created_at, updated_at
		FROM device_groups WHERE tenant_id=$1 AND id=$2
	`, tenantID, id)
	if err := row.Scan(&g.ID, &g.TenantID, &g.SchoolID, &g.Name, &g.Description, &g.GroupType, &g.LocationID, &selectorJSON, &policies, &g.Active, &g.CreatedAt, &g.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.DeviceGroup{}, errors.New("not found")
		}
		return models.DeviceGroup{}, err
	}
	g.Policies = policies
	if len(selectorJSON) > 0 {
		var sel models.GroupSelector
		if err := json.Unmarshal(selectorJSON, &sel); err == nil {
			g.Selector = &sel
		}
	}
	return g, nil
}

func (r *GroupsRepo) Delete(ctx context.Context, tenantID, id string) error {
	result, err := r.pool.Exec(ctx, `DELETE FROM device_groups WHERE tenant_id=$1 AND id=$2`, tenantID, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errors.New("not found")
	}
	return nil
}

type GroupListParams struct {
	TenantID  string
	SchoolID  string
	Type      string
	Active    *bool
	Limit     int
	HasCursor bool
	CursorAt  time.Time
	CursorID  string
}

func (r *GroupsRepo) List(ctx context.Context, p GroupListParams) ([]models.DeviceGroup, string, error) {
	conds := []string{"tenant_id=$1"}
	args := []any{p.TenantID}
	argN := 2

	if p.SchoolID != "" {
		conds = append(conds, "(school_id=$"+itoa(argN)+" OR school_id IS NULL)")
		args = append(args, p.SchoolID)
		argN++
	}
	if p.Type != "" {
		conds = append(conds, "group_type=$"+itoa(argN))
		args = append(args, p.Type)
		argN++
	}
	if p.Active != nil {
		conds = append(conds, "active=$"+itoa(argN))
		args = append(args, *p.Active)
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
		SELECT id, tenant_id, school_id, name, description, group_type, location_id, selector, policies, active, created_at, updated_at
		FROM device_groups
		WHERE ` + strings.Join(conds, " AND ") + `
		ORDER BY created_at DESC, id DESC
		LIMIT $` + itoa(argN)

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	out := []models.DeviceGroup{}
	for rows.Next() {
		var g models.DeviceGroup
		var policies, selectorJSON []byte
		if err := rows.Scan(&g.ID, &g.TenantID, &g.SchoolID, &g.Name, &g.Description, &g.GroupType, &g.LocationID, &selectorJSON, &policies, &g.Active, &g.CreatedAt, &g.UpdatedAt); err != nil {
			return nil, "", err
		}
		g.Policies = policies
		if len(selectorJSON) > 0 {
			var sel models.GroupSelector
			if err := json.Unmarshal(selectorJSON, &sel); err == nil {
				g.Selector = &sel
			}
		}
		out = append(out, g)
	}

	next := ""
	if len(out) > limit {
		last := out[limit-1]
		next = EncodeCursor(last.CreatedAt, last.ID)
		out = out[:limit]
	}
	return out, next, nil
}

// ListBySchool returns all active groups for a school (including tenant-wide)
func (r *GroupsRepo) ListBySchool(ctx context.Context, tenantID, schoolID string) ([]models.DeviceGroup, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, tenant_id, school_id, name, description, group_type, location_id, selector, policies, active, created_at, updated_at
		FROM device_groups
		WHERE tenant_id=$1 AND (school_id=$2 OR school_id IS NULL) AND active=true
		ORDER BY name
	`, tenantID, schoolID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []models.DeviceGroup{}
	for rows.Next() {
		var g models.DeviceGroup
		var policies, selectorJSON []byte
		if err := rows.Scan(&g.ID, &g.TenantID, &g.SchoolID, &g.Name, &g.Description, &g.GroupType, &g.LocationID, &selectorJSON, &policies, &g.Active, &g.CreatedAt, &g.UpdatedAt); err != nil {
			return nil, err
		}
		g.Policies = policies
		if len(selectorJSON) > 0 {
			var sel models.GroupSelector
			if err := json.Unmarshal(selectorJSON, &sel); err == nil {
				g.Selector = &sel
			}
		}
		out = append(out, g)
	}
	return out, nil
}

// ---------- Group Members ----------

func (r *GroupsRepo) AddMember(ctx context.Context, m models.GroupMember) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO group_members (id, tenant_id, group_id, device_id, added_at, added_by)
		VALUES ($1,$2,$3,$4,$5,$6)
		ON CONFLICT (group_id, device_id) DO NOTHING
	`, m.ID, m.TenantID, m.GroupID, m.DeviceID, m.AddedAt, m.AddedBy)
	return err
}

func (r *GroupsRepo) RemoveMember(ctx context.Context, tenantID, groupID, deviceID string) error {
	result, err := r.pool.Exec(ctx, `DELETE FROM group_members WHERE tenant_id=$1 AND group_id=$2 AND device_id=$3`, tenantID, groupID, deviceID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errors.New("not found")
	}
	return nil
}

func (r *GroupsRepo) ListMembers(ctx context.Context, tenantID, groupID string) ([]models.GroupMember, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, tenant_id, group_id, device_id, added_at, added_by
		FROM group_members
		WHERE tenant_id=$1 AND group_id=$2
		ORDER BY added_at DESC
	`, tenantID, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []models.GroupMember{}
	for rows.Next() {
		var m models.GroupMember
		if err := rows.Scan(&m.ID, &m.TenantID, &m.GroupID, &m.DeviceID, &m.AddedAt, &m.AddedBy); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, nil
}

func (r *GroupsRepo) GetMemberCount(ctx context.Context, tenantID, groupID string) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM group_members WHERE tenant_id=$1 AND group_id=$2`, tenantID, groupID).Scan(&count)
	return count, err
}

// GetGroupsForDevice returns all groups a device belongs to
func (r *GroupsRepo) GetGroupsForDevice(ctx context.Context, tenantID, deviceID string) ([]models.DeviceGroup, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT g.id, g.tenant_id, g.school_id, g.name, g.description, g.group_type, g.location_id, g.selector, g.policies, g.active, g.created_at, g.updated_at
		FROM device_groups g
		JOIN group_members gm ON gm.group_id = g.id
		WHERE g.tenant_id=$1 AND gm.device_id=$2 AND g.active=true
		ORDER BY g.name
	`, tenantID, deviceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []models.DeviceGroup{}
	for rows.Next() {
		var g models.DeviceGroup
		var policies, selectorJSON []byte
		if err := rows.Scan(&g.ID, &g.TenantID, &g.SchoolID, &g.Name, &g.Description, &g.GroupType, &g.LocationID, &selectorJSON, &policies, &g.Active, &g.CreatedAt, &g.UpdatedAt); err != nil {
			return nil, err
		}
		g.Policies = policies
		if len(selectorJSON) > 0 {
			var sel models.GroupSelector
			if err := json.Unmarshal(selectorJSON, &sel); err == nil {
				g.Selector = &sel
			}
		}
		out = append(out, g)
	}
	return out, nil
}

// BulkAddMembers adds multiple devices to a group
func (r *GroupsRepo) BulkAddMembers(ctx context.Context, tenantID, groupID string, deviceIDs []string, addedBy string) (int, error) {
	if len(deviceIDs) == 0 {
		return 0, nil
	}

	now := time.Now().UTC()
	count := 0
	for _, deviceID := range deviceIDs {
		id := NewID("gm")
		result, err := r.pool.Exec(ctx, `
			INSERT INTO group_members (id, tenant_id, group_id, device_id, added_at, added_by)
			VALUES ($1,$2,$3,$4,$5,$6)
			ON CONFLICT (group_id, device_id) DO NOTHING
		`, id, tenantID, groupID, deviceID, now, addedBy)
		if err != nil {
			return count, err
		}
		if result.RowsAffected() > 0 {
			count++
		}
	}
	return count, nil
}

// BulkRemoveMembers removes multiple devices from a group
func (r *GroupsRepo) BulkRemoveMembers(ctx context.Context, tenantID, groupID string, deviceIDs []string) (int, error) {
	if len(deviceIDs) == 0 {
		return 0, nil
	}

	result, err := r.pool.Exec(ctx, `
		DELETE FROM group_members WHERE tenant_id=$1 AND group_id=$2 AND device_id = ANY($3)
	`, tenantID, groupID, deviceIDs)
	if err != nil {
		return 0, err
	}
	return int(result.RowsAffected()), nil
}
