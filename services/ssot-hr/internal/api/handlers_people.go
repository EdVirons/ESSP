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

func (s *Server) CreatePerson(w http.ResponseWriter, r *http.Request) {
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
		id = newID(models.PrefixPerson)
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
		INSERT INTO people (id, tenant_id, org_unit_id, status, given_name, family_name, email, phone, title, avatar_url, spec_json, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $12)
	`, id, tenant, trim(body["orgUnitId"]), status, trim(body["givenName"]), trim(body["familyName"]),
		trim(body["email"]), trim(body["phone"]), trim(body["title"]), trim(body["avatarUrl"]), specJSON, now)
	if err != nil {
		s.log.Error("create person failed", zap.Error(err))
		httpx.Error(w, 500, "create failed")
		return
	}

	// Publish event
	_ = s.pub.PublishJSON("ssot.hr.person.created", map[string]any{
		"tenantId": tenant, "personId": id, "email": trim(body["email"]), "at": now,
	})

	httpx.WriteJSON(w, 201, map[string]any{"ok": true, "id": id})
}

func (s *Server) GetPerson(w http.ResponseWriter, r *http.Request) {
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

	var p models.Person
	var specJSON string
	err := s.db.QueryRow(r.Context(), `
		SELECT id, tenant_id, org_unit_id, status, given_name, family_name, email, phone, title, avatar_url, spec_json, created_at, updated_at
		FROM people WHERE tenant_id=$1 AND id=$2
	`, tenant, id).Scan(&p.ID, &p.TenantID, &p.OrgUnitID, &p.Status, &p.GivenName, &p.FamilyName, &p.Email, &p.Phone, &p.Title, &p.AvatarURL, &specJSON, &p.CreatedAt, &p.UpdatedAt)
	if err == pgx.ErrNoRows {
		httpx.Error(w, 404, "person not found")
		return
	}
	if err != nil {
		s.log.Error("get person failed", zap.Error(err))
		httpx.Error(w, 500, "query failed")
		return
	}
	p.SpecJSON = parseJSON(specJSON)

	httpx.WriteJSON(w, 200, p)
}

func (s *Server) ListPeople(w http.ResponseWriter, r *http.Request) {
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

	// Build query with filters
	query := `SELECT id, tenant_id, org_unit_id, status, given_name, family_name, email, phone, title, avatar_url, spec_json, created_at, updated_at FROM people WHERE tenant_id=$1`
	args := []any{tenant}
	argIdx := 2

	if status := q.Get("status"); status != "" {
		query += " AND status=$" + strconv.Itoa(argIdx)
		args = append(args, status)
		argIdx++
	}
	if orgUnitID := q.Get("orgUnitId"); orgUnitID != "" {
		query += " AND org_unit_id=$" + strconv.Itoa(argIdx)
		args = append(args, orgUnitID)
		argIdx++
	}
	if email := q.Get("email"); email != "" {
		query += " AND email=$" + strconv.Itoa(argIdx)
		args = append(args, email)
		argIdx++
	}

	query += " ORDER BY family_name, given_name LIMIT $" + strconv.Itoa(argIdx) + " OFFSET $" + strconv.Itoa(argIdx+1)
	args = append(args, limit, offset)

	rows, err := s.db.Query(r.Context(), query, args...)
	if err != nil {
		s.log.Error("list people failed", zap.Error(err))
		httpx.Error(w, 500, "query failed")
		return
	}
	defer rows.Close()

	var items []models.Person
	for rows.Next() {
		var p models.Person
		var specJSON string
		if err := rows.Scan(&p.ID, &p.TenantID, &p.OrgUnitID, &p.Status, &p.GivenName, &p.FamilyName, &p.Email, &p.Phone, &p.Title, &p.AvatarURL, &specJSON, &p.CreatedAt, &p.UpdatedAt); err != nil {
			s.log.Error("scan person failed", zap.Error(err))
			httpx.Error(w, 500, "scan failed")
			return
		}
		p.SpecJSON = parseJSON(specJSON)
		items = append(items, p)
	}

	httpx.WriteJSON(w, 200, map[string]any{"items": items, "limit": limit, "offset": offset})
}

func (s *Server) PatchPerson(w http.ResponseWriter, r *http.Request) {
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

	// Build dynamic UPDATE
	setClauses := []string{}
	args := []any{tenant, id}
	argIdx := 3

	if v, ok := body["orgUnitId"]; ok {
		setClauses = append(setClauses, "org_unit_id=$"+strconv.Itoa(argIdx))
		args = append(args, trim(v))
		argIdx++
	}
	if v, ok := body["status"]; ok {
		setClauses = append(setClauses, "status=$"+strconv.Itoa(argIdx))
		args = append(args, trim(v))
		argIdx++
	}
	if v, ok := body["givenName"]; ok {
		setClauses = append(setClauses, "given_name=$"+strconv.Itoa(argIdx))
		args = append(args, trim(v))
		argIdx++
	}
	if v, ok := body["familyName"]; ok {
		setClauses = append(setClauses, "family_name=$"+strconv.Itoa(argIdx))
		args = append(args, trim(v))
		argIdx++
	}
	if v, ok := body["phone"]; ok {
		setClauses = append(setClauses, "phone=$"+strconv.Itoa(argIdx))
		args = append(args, trim(v))
		argIdx++
	}
	if v, ok := body["title"]; ok {
		setClauses = append(setClauses, "title=$"+strconv.Itoa(argIdx))
		args = append(args, trim(v))
		argIdx++
	}
	if v, ok := body["avatarUrl"]; ok {
		setClauses = append(setClauses, "avatar_url=$"+strconv.Itoa(argIdx))
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

	query := "UPDATE people SET "
	for i, c := range setClauses {
		if i > 0 {
			query += ", "
		}
		query += c
	}
	query += " WHERE tenant_id=$1 AND id=$2"

	result, err := s.db.Exec(r.Context(), query, args...)
	if err != nil {
		s.log.Error("patch person failed", zap.Error(err))
		httpx.Error(w, 500, "update failed")
		return
	}
	if result.RowsAffected() == 0 {
		httpx.Error(w, 404, "person not found")
		return
	}

	// Publish event
	_ = s.pub.PublishJSON("ssot.hr.person.updated", map[string]any{
		"tenantId": tenant, "personId": id, "at": time.Now().UTC(),
	})

	httpx.WriteJSON(w, 200, map[string]any{"ok": true})
}

func (s *Server) DeletePerson(w http.ResponseWriter, r *http.Request) {
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

	result, err := s.db.Exec(r.Context(), `DELETE FROM people WHERE tenant_id=$1 AND id=$2`, tenant, id)
	if err != nil {
		s.log.Error("delete person failed", zap.Error(err))
		httpx.Error(w, 500, "delete failed")
		return
	}
	if result.RowsAffected() == 0 {
		httpx.Error(w, 404, "person not found")
		return
	}

	httpx.WriteJSON(w, 200, map[string]any{"ok": true})
}
