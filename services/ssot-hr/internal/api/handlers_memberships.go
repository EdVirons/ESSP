package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/edvirons/ssp/shared/pkg/httpx"
	"github.com/edvirons/ssp/ssot_hr/internal/models"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func (s *Server) CreateMembership(w http.ResponseWriter, r *http.Request) {
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

	teamID := trim(body["teamId"])
	personID := trim(body["personId"])
	if teamID == "" || personID == "" {
		httpx.Error(w, 400, "teamId and personId required")
		return
	}

	id := trim(body["id"])
	if id == "" {
		id = newID(models.PrefixMembership)
	}
	role := trim(body["role"])
	if role == "" {
		role = "member"
	}
	status := trim(body["status"])
	if status == "" {
		status = "active"
	}
	specJSON := "{}"
	if spec, ok := body["specJson"].(map[string]any); ok {
		if b, err := json.Marshal(spec); err == nil {
			specJSON = string(b)
		}
	}

	now := time.Now().UTC()
	_, err := s.db.Exec(r.Context(), `
		INSERT INTO team_memberships (id, tenant_id, team_id, person_id, role, status, started_at, ended_at, spec_json, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NULL, $8, $9, $9)
	`, id, tenant, teamID, personID, role, status, now, specJSON, now)
	if err != nil {
		s.log.Error("create membership failed", zap.Error(err))
		httpx.Error(w, 500, "create failed")
		return
	}

	// Publish event
	_ = s.pub.PublishJSON("ssot.hr.membership.added", map[string]any{
		"tenantId": tenant, "membershipId": id, "teamId": teamID, "personId": personID, "role": role, "at": now,
	})

	httpx.WriteJSON(w, 201, map[string]any{"ok": true, "id": id})
}

func (s *Server) ListMemberships(w http.ResponseWriter, r *http.Request) {
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

	query := `SELECT id, tenant_id, team_id, person_id, role, status, started_at, ended_at, spec_json, created_at, updated_at FROM team_memberships WHERE tenant_id=$1`
	args := []any{tenant}
	argIdx := 2

	if teamID := q.Get("teamId"); teamID != "" {
		query += " AND team_id=$" + strconv.Itoa(argIdx)
		args = append(args, teamID)
		argIdx++
	}
	if personID := q.Get("personId"); personID != "" {
		query += " AND person_id=$" + strconv.Itoa(argIdx)
		args = append(args, personID)
		argIdx++
	}
	if status := q.Get("status"); status != "" {
		query += " AND status=$" + strconv.Itoa(argIdx)
		args = append(args, status)
		argIdx++
	}

	query += " ORDER BY created_at DESC LIMIT $" + strconv.Itoa(argIdx) + " OFFSET $" + strconv.Itoa(argIdx+1)
	args = append(args, limit, offset)

	rows, err := s.db.Query(r.Context(), query, args...)
	if err != nil {
		s.log.Error("list memberships failed", zap.Error(err))
		httpx.Error(w, 500, "query failed")
		return
	}
	defer rows.Close()

	var items []models.TeamMembership
	for rows.Next() {
		var m models.TeamMembership
		var specJSON string
		if err := rows.Scan(&m.ID, &m.TenantID, &m.TeamID, &m.PersonID, &m.Role, &m.Status, &m.StartedAt, &m.EndedAt, &specJSON, &m.CreatedAt, &m.UpdatedAt); err != nil {
			s.log.Error("scan membership failed", zap.Error(err))
			httpx.Error(w, 500, "scan failed")
			return
		}
		m.SpecJSON = parseJSON(specJSON)
		items = append(items, m)
	}

	httpx.WriteJSON(w, 200, map[string]any{"items": items, "limit": limit, "offset": offset})
}

func (s *Server) DeleteMembership(w http.ResponseWriter, r *http.Request) {
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

	// Get membership info for event
	var teamID, personID string
	_ = s.db.QueryRow(r.Context(), `SELECT team_id, person_id FROM team_memberships WHERE tenant_id=$1 AND id=$2`, tenant, id).Scan(&teamID, &personID)

	result, err := s.db.Exec(r.Context(), `DELETE FROM team_memberships WHERE tenant_id=$1 AND id=$2`, tenant, id)
	if err != nil {
		s.log.Error("delete membership failed", zap.Error(err))
		httpx.Error(w, 500, "delete failed")
		return
	}
	if result.RowsAffected() == 0 {
		httpx.Error(w, 404, "membership not found")
		return
	}

	// Publish event
	if teamID != "" && personID != "" {
		_ = s.pub.PublishJSON("ssot.hr.membership.removed", map[string]any{
			"tenantId": tenant, "membershipId": id, "teamId": teamID, "personId": personID, "at": time.Now().UTC(),
		})
	}

	httpx.WriteJSON(w, 200, map[string]any{"ok": true})
}
