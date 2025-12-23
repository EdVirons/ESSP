package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/edvirons/ssp/shared/pkg/httpx"
	"github.com/edvirons/ssp/ssot_hr/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

func (s *Server) CreateOrgUnit(w http.ResponseWriter, r *http.Request) {
	tenant := httpx.TenantID(r)
	if tenant == "" {
		httpx.Error(w, 400, "X-Tenant-Id required")
		return
	}

	var body map[string]any
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		httpx.Error(w, 400, "invalid json")
		return
	}

	id := trim(body["id"])
	if id == "" {
		id = newID(models.PrefixOrgUnit)
	}
	kind := trim(body["kind"])
	if kind == "" {
		kind = "department"
	}
	specJSON := "{}"
	if spec, ok := body["specJson"].(map[string]any); ok {
		if b, err := json.Marshal(spec); err == nil {
			specJSON = string(b)
		}
	}

	now := time.Now().UTC()
	_, err := s.db.Exec(r.Context(), `
		INSERT INTO org_units (id, tenant_id, parent_id, code, name, kind, spec_json, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $8)
	`, id, tenant, trim(body["parentId"]), trim(body["code"]), trim(body["name"]), kind, specJSON, now)
	if err != nil {
		s.log.Error("create org unit failed", zap.Error(err))
		httpx.Error(w, 500, "create failed")
		return
	}

	// Invalidate org tree cache
	if s.cache != nil {
		_ = s.cache.Del(r.Context(), "orgtree:"+tenant)
	}

	httpx.WriteJSON(w, 201, map[string]any{"ok": true, "id": id})
}

func (s *Server) GetOrgUnit(w http.ResponseWriter, r *http.Request) {
	tenant := httpx.TenantID(r)
	if tenant == "" {
		httpx.Error(w, 400, "X-Tenant-Id required")
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		httpx.Error(w, 400, "id required")
		return
	}

	var o models.OrgUnit
	var specJSON string
	err := s.db.QueryRow(r.Context(), `
		SELECT id, tenant_id, parent_id, code, name, kind, spec_json, created_at, updated_at
		FROM org_units WHERE tenant_id=$1 AND id=$2
	`, tenant, id).Scan(&o.ID, &o.TenantID, &o.ParentID, &o.Code, &o.Name, &o.Kind, &specJSON, &o.CreatedAt, &o.UpdatedAt)
	if err == pgx.ErrNoRows {
		httpx.Error(w, 404, "org unit not found")
		return
	}
	if err != nil {
		s.log.Error("get org unit failed", zap.Error(err))
		httpx.Error(w, 500, "query failed")
		return
	}
	o.SpecJSON = parseJSON(specJSON)

	httpx.WriteJSON(w, 200, o)
}

func (s *Server) ListOrgUnits(w http.ResponseWriter, r *http.Request) {
	tenant := httpx.TenantID(r)
	if tenant == "" {
		httpx.Error(w, 400, "X-Tenant-Id required")
		return
	}

	q := r.URL.Query()
	limit := 50
	if l, err := strconv.Atoi(q.Get("limit")); err == nil && l > 0 && l <= 200 {
		limit = l
	}
	offset := 0
	if o, err := strconv.Atoi(q.Get("offset")); err == nil && o >= 0 {
		offset = o
	}

	query := `SELECT id, tenant_id, parent_id, code, name, kind, spec_json, created_at, updated_at FROM org_units WHERE tenant_id=$1`
	args := []any{tenant}
	argIdx := 2

	if parentID := q.Get("parentId"); parentID != "" {
		query += " AND parent_id=$" + strconv.Itoa(argIdx)
		args = append(args, parentID)
		argIdx++
	}
	if kind := q.Get("kind"); kind != "" {
		query += " AND kind=$" + strconv.Itoa(argIdx)
		args = append(args, kind)
		argIdx++
	}

	query += " ORDER BY name LIMIT $" + strconv.Itoa(argIdx) + " OFFSET $" + strconv.Itoa(argIdx+1)
	args = append(args, limit, offset)

	rows, err := s.db.Query(r.Context(), query, args...)
	if err != nil {
		s.log.Error("list org units failed", zap.Error(err))
		httpx.Error(w, 500, "query failed")
		return
	}
	defer rows.Close()

	var items []models.OrgUnit
	for rows.Next() {
		var o models.OrgUnit
		var specJSON string
		if err := rows.Scan(&o.ID, &o.TenantID, &o.ParentID, &o.Code, &o.Name, &o.Kind, &specJSON, &o.CreatedAt, &o.UpdatedAt); err != nil {
			s.log.Error("scan org unit failed", zap.Error(err))
			httpx.Error(w, 500, "scan failed")
			return
		}
		o.SpecJSON = parseJSON(specJSON)
		items = append(items, o)
	}

	httpx.WriteJSON(w, 200, map[string]any{"items": items, "limit": limit, "offset": offset})
}

// OrgTreeNode represents a node in the org tree
type OrgTreeNode struct {
	models.OrgUnit
	Children []OrgTreeNode `json:"children,omitempty"`
}

func (s *Server) GetOrgTree(w http.ResponseWriter, r *http.Request) {
	tenant := httpx.TenantID(r)
	if tenant == "" {
		httpx.Error(w, 400, "X-Tenant-Id required")
		return
	}

	// Try cache first
	if s.cache != nil {
		cacheKey := "orgtree:" + tenant
		cached, err := s.cache.Get(r.Context(), cacheKey).Result()
		if err == nil && cached != "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			_, _ = w.Write([]byte(cached))
			return
		}
	}

	// Query all org units
	rows, err := s.db.Query(r.Context(), `
		SELECT id, tenant_id, parent_id, code, name, kind, spec_json, created_at, updated_at
		FROM org_units WHERE tenant_id=$1 ORDER BY name
	`, tenant)
	if err != nil {
		s.log.Error("get org tree failed", zap.Error(err))
		httpx.Error(w, 500, "query failed")
		return
	}
	defer rows.Close()

	var units []models.OrgUnit
	for rows.Next() {
		var o models.OrgUnit
		var specJSON string
		if err := rows.Scan(&o.ID, &o.TenantID, &o.ParentID, &o.Code, &o.Name, &o.Kind, &specJSON, &o.CreatedAt, &o.UpdatedAt); err != nil {
			s.log.Error("scan org unit failed", zap.Error(err))
			httpx.Error(w, 500, "scan failed")
			return
		}
		o.SpecJSON = parseJSON(specJSON)
		units = append(units, o)
	}

	// Build tree
	tree := buildOrgTree(units)

	// Cache result
	if s.cache != nil {
		if b, err := json.Marshal(tree); err == nil {
			_ = s.cache.Set(r.Context(), "orgtree:"+tenant, string(b), 5*time.Minute)
		}
	}

	httpx.WriteJSON(w, 200, tree)
}

func buildOrgTree(units []models.OrgUnit) []OrgTreeNode {
	// Build map of id -> unit
	unitMap := make(map[string]models.OrgUnit)
	for _, u := range units {
		unitMap[u.ID] = u
	}

	// Build children map
	childrenMap := make(map[string][]models.OrgUnit)
	var roots []models.OrgUnit
	for _, u := range units {
		if u.ParentID == "" {
			roots = append(roots, u)
		} else {
			childrenMap[u.ParentID] = append(childrenMap[u.ParentID], u)
		}
	}

	// Recursive function to build tree nodes
	var buildNode func(u models.OrgUnit) OrgTreeNode
	buildNode = func(u models.OrgUnit) OrgTreeNode {
		node := OrgTreeNode{OrgUnit: u}
		for _, child := range childrenMap[u.ID] {
			node.Children = append(node.Children, buildNode(child))
		}
		return node
	}

	var tree []OrgTreeNode
	for _, root := range roots {
		tree = append(tree, buildNode(root))
	}
	return tree
}

func (s *Server) PatchOrgUnit(w http.ResponseWriter, r *http.Request) {
	tenant := httpx.TenantID(r)
	if tenant == "" {
		httpx.Error(w, 400, "X-Tenant-Id required")
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		httpx.Error(w, 400, "id required")
		return
	}

	var body map[string]any
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		httpx.Error(w, 400, "invalid json")
		return
	}

	setClauses := []string{}
	args := []any{tenant, id}
	argIdx := 3

	if v, ok := body["parentId"]; ok {
		setClauses = append(setClauses, "parent_id=$"+strconv.Itoa(argIdx))
		args = append(args, trim(v))
		argIdx++
	}
	if v, ok := body["code"]; ok {
		setClauses = append(setClauses, "code=$"+strconv.Itoa(argIdx))
		args = append(args, trim(v))
		argIdx++
	}
	if v, ok := body["name"]; ok {
		setClauses = append(setClauses, "name=$"+strconv.Itoa(argIdx))
		args = append(args, trim(v))
		argIdx++
	}
	if v, ok := body["kind"]; ok {
		setClauses = append(setClauses, "kind=$"+strconv.Itoa(argIdx))
		args = append(args, trim(v))
		argIdx++
	}
	if spec, ok := body["specJson"].(map[string]any); ok {
		if b, err := json.Marshal(spec); err == nil {
			setClauses = append(setClauses, "spec_json=$"+strconv.Itoa(argIdx))
			args = append(args, string(b))
			argIdx++
		}
	}

	if len(setClauses) == 0 {
		httpx.Error(w, 400, "no fields to update")
		return
	}

	setClauses = append(setClauses, "updated_at=$"+strconv.Itoa(argIdx))
	args = append(args, time.Now().UTC())

	query := "UPDATE org_units SET "
	for i, c := range setClauses {
		if i > 0 {
			query += ", "
		}
		query += c
	}
	query += " WHERE tenant_id=$1 AND id=$2"

	result, err := s.db.Exec(r.Context(), query, args...)
	if err != nil {
		s.log.Error("patch org unit failed", zap.Error(err))
		httpx.Error(w, 500, "update failed")
		return
	}
	if result.RowsAffected() == 0 {
		httpx.Error(w, 404, "org unit not found")
		return
	}

	// Invalidate org tree cache
	if s.cache != nil {
		_ = s.cache.Del(r.Context(), "orgtree:"+tenant)
	}

	httpx.WriteJSON(w, 200, map[string]any{"ok": true})
}

func (s *Server) DeleteOrgUnit(w http.ResponseWriter, r *http.Request) {
	tenant := httpx.TenantID(r)
	if tenant == "" {
		httpx.Error(w, 400, "X-Tenant-Id required")
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		httpx.Error(w, 400, "id required")
		return
	}

	// Check if org unit has children
	var childCount int
	err := s.db.QueryRow(r.Context(), `SELECT COUNT(*) FROM org_units WHERE tenant_id=$1 AND parent_id=$2`, tenant, id).Scan(&childCount)
	if err != nil {
		s.log.Error("check children failed", zap.Error(err))
		httpx.Error(w, 500, "check failed")
		return
	}
	if childCount > 0 {
		httpx.Error(w, 400, "cannot delete org unit with children")
		return
	}

	// Check if org unit has people assigned
	var peopleCount int
	err = s.db.QueryRow(r.Context(), `SELECT COUNT(*) FROM people WHERE tenant_id=$1 AND org_unit_id=$2`, tenant, id).Scan(&peopleCount)
	if err != nil {
		s.log.Error("check people failed", zap.Error(err))
		httpx.Error(w, 500, "check failed")
		return
	}
	if peopleCount > 0 {
		httpx.Error(w, 400, "cannot delete org unit with assigned people")
		return
	}

	// Get org unit info for event
	var code, name string
	_ = s.db.QueryRow(r.Context(), `SELECT code, name FROM org_units WHERE tenant_id=$1 AND id=$2`, tenant, id).Scan(&code, &name)

	// Delete the org unit
	result, err := s.db.Exec(r.Context(), `DELETE FROM org_units WHERE tenant_id=$1 AND id=$2`, tenant, id)
	if err != nil {
		s.log.Error("delete org unit failed", zap.Error(err))
		httpx.Error(w, 500, "delete failed")
		return
	}
	if result.RowsAffected() == 0 {
		httpx.Error(w, 404, "org unit not found")
		return
	}

	// Invalidate org tree cache
	if s.cache != nil {
		_ = s.cache.Del(r.Context(), "orgtree:"+tenant)
	}

	// Publish event
	_ = s.pub.PublishJSON("ssot.hr.org_unit.deleted", map[string]any{
		"tenantId": tenant, "orgUnitId": id, "code": code, "name": name, "at": time.Now().UTC(),
	})

	httpx.WriteJSON(w, 200, map[string]any{"ok": true})
}
