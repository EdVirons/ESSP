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

func (s *Server) CreateTeam(w http.ResponseWriter, r *http.Request) {
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
		id = newID(models.PrefixTeam)
	}
	specJSON := "{}"
	if spec, ok := body["specJson"].(map[string]any); ok {
		if b, err := json.Marshal(spec); err == nil {
			specJSON = string(b)
		}
	}

	now := time.Now().UTC()
	_, err := s.db.Exec(r.Context(), `
		INSERT INTO teams (id, tenant_id, org_unit_id, key, name, description, spec_json, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $8)
	`, id, tenant, trim(body["orgUnitId"]), trim(body["key"]), trim(body["name"]), trim(body["description"]), specJSON, now)
	if err != nil {
		s.log.Error("create team failed", zap.Error(err))
		httpx.Error(w, 500, "create failed")
		return
	}

	// Publish event
	_ = s.pub.PublishJSON("ssot.hr.team.created", map[string]any{
		"tenantId": tenant, "teamId": id, "key": trim(body["key"]), "at": now,
	})

	httpx.WriteJSON(w, 201, map[string]any{"ok": true, "id": id})
}

func (s *Server) GetTeam(w http.ResponseWriter, r *http.Request) {
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

	var t models.Team
	var specJSON string
	err := s.db.QueryRow(r.Context(), `
		SELECT id, tenant_id, org_unit_id, key, name, description, spec_json, created_at, updated_at
		FROM teams WHERE tenant_id=$1 AND id=$2
	`, tenant, id).Scan(&t.ID, &t.TenantID, &t.OrgUnitID, &t.Key, &t.Name, &t.Description, &specJSON, &t.CreatedAt, &t.UpdatedAt)
	if err == pgx.ErrNoRows {
		httpx.Error(w, 404, "team not found")
		return
	}
	if err != nil {
		s.log.Error("get team failed", zap.Error(err))
		httpx.Error(w, 500, "query failed")
		return
	}
	t.SpecJSON = parseJSON(specJSON)

	httpx.WriteJSON(w, 200, t)
}

func (s *Server) ListTeams(w http.ResponseWriter, r *http.Request) {
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

	query := `SELECT id, tenant_id, org_unit_id, key, name, description, spec_json, created_at, updated_at FROM teams WHERE tenant_id=$1`
	args := []any{tenant}
	argIdx := 2

	if orgUnitID := q.Get("orgUnitId"); orgUnitID != "" {
		query += " AND org_unit_id=$" + strconv.Itoa(argIdx)
		args = append(args, orgUnitID)
		argIdx++
	}
	if key := q.Get("key"); key != "" {
		query += " AND key=$" + strconv.Itoa(argIdx)
		args = append(args, key)
		argIdx++
	}

	query += " ORDER BY name LIMIT $" + strconv.Itoa(argIdx) + " OFFSET $" + strconv.Itoa(argIdx+1)
	args = append(args, limit, offset)

	rows, err := s.db.Query(r.Context(), query, args...)
	if err != nil {
		s.log.Error("list teams failed", zap.Error(err))
		httpx.Error(w, 500, "query failed")
		return
	}
	defer rows.Close()

	var items []models.Team
	for rows.Next() {
		var t models.Team
		var specJSON string
		if err := rows.Scan(&t.ID, &t.TenantID, &t.OrgUnitID, &t.Key, &t.Name, &t.Description, &specJSON, &t.CreatedAt, &t.UpdatedAt); err != nil {
			s.log.Error("scan team failed", zap.Error(err))
			httpx.Error(w, 500, "scan failed")
			return
		}
		t.SpecJSON = parseJSON(specJSON)
		items = append(items, t)
	}

	httpx.WriteJSON(w, 200, map[string]any{"items": items, "limit": limit, "offset": offset})
}

func (s *Server) PatchTeam(w http.ResponseWriter, r *http.Request) {
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

	if v, ok := body["orgUnitId"]; ok {
		setClauses = append(setClauses, "org_unit_id=$"+strconv.Itoa(argIdx))
		args = append(args, trim(v))
		argIdx++
	}
	if v, ok := body["key"]; ok {
		setClauses = append(setClauses, "key=$"+strconv.Itoa(argIdx))
		args = append(args, trim(v))
		argIdx++
	}
	if v, ok := body["name"]; ok {
		setClauses = append(setClauses, "name=$"+strconv.Itoa(argIdx))
		args = append(args, trim(v))
		argIdx++
	}
	if v, ok := body["description"]; ok {
		setClauses = append(setClauses, "description=$"+strconv.Itoa(argIdx))
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

	query := "UPDATE teams SET "
	for i, c := range setClauses {
		if i > 0 {
			query += ", "
		}
		query += c
	}
	query += " WHERE tenant_id=$1 AND id=$2"

	result, err := s.db.Exec(r.Context(), query, args...)
	if err != nil {
		s.log.Error("patch team failed", zap.Error(err))
		httpx.Error(w, 500, "update failed")
		return
	}
	if result.RowsAffected() == 0 {
		httpx.Error(w, 404, "team not found")
		return
	}

	httpx.WriteJSON(w, 200, map[string]any{"ok": true})
}

func (s *Server) DeleteTeam(w http.ResponseWriter, r *http.Request) {
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

	// Get team info for event
	var key, name string
	_ = s.db.QueryRow(r.Context(), `SELECT key, name FROM teams WHERE tenant_id=$1 AND id=$2`, tenant, id).Scan(&key, &name)

	// Delete associated team memberships first
	_, err := s.db.Exec(r.Context(), `DELETE FROM team_memberships WHERE tenant_id=$1 AND team_id=$2`, tenant, id)
	if err != nil {
		s.log.Error("delete team memberships failed", zap.Error(err))
		httpx.Error(w, 500, "delete memberships failed")
		return
	}

	// Delete the team
	result, err := s.db.Exec(r.Context(), `DELETE FROM teams WHERE tenant_id=$1 AND id=$2`, tenant, id)
	if err != nil {
		s.log.Error("delete team failed", zap.Error(err))
		httpx.Error(w, 500, "delete failed")
		return
	}
	if result.RowsAffected() == 0 {
		httpx.Error(w, 404, "team not found")
		return
	}

	// Publish event
	_ = s.pub.PublishJSON("ssot.hr.team.deleted", map[string]any{
		"tenantId": tenant, "teamId": id, "key": key, "name": name, "at": time.Now().UTC(),
	})

	httpx.WriteJSON(w, 200, map[string]any{"ok": true})
}
