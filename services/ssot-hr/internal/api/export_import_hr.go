package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/edvirons/ssp/shared/pkg/httpx"
	"github.com/edvirons/ssp/ssot_hr/internal/models"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

func (s *Server) Export(w http.ResponseWriter, r *http.Request) {
	tenant := httpx.TenantID(r)
	if tenant == "" {
		httpx.Error(w, 400, "X-Tenant-Id required")
		return
	}

	payload, err := s.exportAll(r.Context(), tenant)
	if err != nil {
		s.log.Error("export failed", zap.Error(err))
		httpx.Error(w, 500, "export failed")
		return
	}
	httpx.WriteJSON(w, 200, payload)
}

func (s *Server) Import(w http.ResponseWriter, r *http.Request) {
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

	res, err := s.importAll(r.Context(), tenant, body)
	if err != nil {
		s.log.Error("import failed", zap.Error(err))
		httpx.Error(w, 400, err.Error())
		return
	}

	// Publish NATS events
	_ = s.pub.PublishJSON("ssot.hr.snapshot", map[string]any{
		"tenantId":   tenant,
		"importedAt": time.Now().UTC(),
		"counts":     res,
	})
	_ = s.pub.PublishJSON("ssot.hr.changed", map[string]any{
		"tenantId": tenant,
		"at":       time.Now().UTC(),
	})

	httpx.WriteJSON(w, 200, map[string]any{"ok": true, "imported": res})
}

func (s *Server) exportAll(ctx context.Context, tenant string) (models.ExportPayload, error) {
	p := models.ExportPayload{
		Version:     "1",
		GeneratedAt: time.Now().UTC(),
	}

	// Export org_units
	rows, err := s.db.Query(ctx, `
		SELECT id, tenant_id, parent_id, code, name, kind, spec_json, created_at, updated_at
		FROM org_units WHERE tenant_id=$1 ORDER BY created_at
	`, tenant)
	if err != nil {
		return p, err
	}
	defer rows.Close()
	for rows.Next() {
		var o models.OrgUnit
		var specJSON string
		if err := rows.Scan(&o.ID, &o.TenantID, &o.ParentID, &o.Code, &o.Name, &o.Kind, &specJSON, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return p, err
		}
		o.SpecJSON = parseJSON(specJSON)
		p.OrgUnits = append(p.OrgUnits, o)
	}

	// Export people
	rows, err = s.db.Query(ctx, `
		SELECT id, tenant_id, org_unit_id, status, given_name, family_name, email, phone, title, avatar_url, spec_json, created_at, updated_at
		FROM people WHERE tenant_id=$1 ORDER BY created_at
	`, tenant)
	if err != nil {
		return p, err
	}
	defer rows.Close()
	for rows.Next() {
		var person models.Person
		var specJSON string
		if err := rows.Scan(&person.ID, &person.TenantID, &person.OrgUnitID, &person.Status, &person.GivenName, &person.FamilyName, &person.Email, &person.Phone, &person.Title, &person.AvatarURL, &specJSON, &person.CreatedAt, &person.UpdatedAt); err != nil {
			return p, err
		}
		person.SpecJSON = parseJSON(specJSON)
		p.People = append(p.People, person)
	}

	// Export teams
	rows, err = s.db.Query(ctx, `
		SELECT id, tenant_id, org_unit_id, key, name, description, spec_json, created_at, updated_at
		FROM teams WHERE tenant_id=$1 ORDER BY created_at
	`, tenant)
	if err != nil {
		return p, err
	}
	defer rows.Close()
	for rows.Next() {
		var t models.Team
		var specJSON string
		if err := rows.Scan(&t.ID, &t.TenantID, &t.OrgUnitID, &t.Key, &t.Name, &t.Description, &specJSON, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return p, err
		}
		t.SpecJSON = parseJSON(specJSON)
		p.Teams = append(p.Teams, t)
	}

	// Export team_memberships
	rows, err = s.db.Query(ctx, `
		SELECT id, tenant_id, team_id, person_id, role, status, started_at, ended_at, spec_json, created_at, updated_at
		FROM team_memberships WHERE tenant_id=$1 ORDER BY created_at
	`, tenant)
	if err != nil {
		return p, err
	}
	defer rows.Close()
	for rows.Next() {
		var m models.TeamMembership
		var specJSON string
		if err := rows.Scan(&m.ID, &m.TenantID, &m.TeamID, &m.PersonID, &m.Role, &m.Status, &m.StartedAt, &m.EndedAt, &specJSON, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return p, err
		}
		m.SpecJSON = parseJSON(specJSON)
		p.TeamMemberships = append(p.TeamMemberships, m)
	}

	return p, nil
}

func (s *Server) importAll(ctx context.Context, tenant string, body map[string]any) (map[string]int, error) {
	res := map[string]int{"orgUnits": 0, "people": 0, "teams": 0, "teamMemberships": 0}

	err := withTx(ctx, s.db, func(tx pgx.Tx) error {
		now := time.Now().UTC()

		// Import org_units
		if raw, ok := body["orgUnits"].([]any); ok {
			for _, item := range raw {
				m, ok := item.(map[string]any)
				if !ok {
					continue
				}
				id := trim(m["id"])
				if id == "" {
					id = newID(models.PrefixOrgUnit)
				}
				specJSON := "{}"
				if spec, ok := m["specJson"].(map[string]any); ok {
					if b, err := json.Marshal(spec); err == nil {
						specJSON = string(b)
					}
				}
				_, err := tx.Exec(ctx, `
					INSERT INTO org_units (id, tenant_id, parent_id, code, name, kind, spec_json, created_at, updated_at)
					VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $8)
					ON CONFLICT (id) DO UPDATE SET
						parent_id=EXCLUDED.parent_id, code=EXCLUDED.code, name=EXCLUDED.name,
						kind=EXCLUDED.kind, spec_json=EXCLUDED.spec_json, updated_at=EXCLUDED.updated_at
				`, id, tenant, trim(m["parentId"]), trim(m["code"]), trim(m["name"]), trim(m["kind"]), specJSON, now)
				if err != nil {
					return err
				}
				res["orgUnits"]++
			}
		}

		// Import people
		if raw, ok := body["people"].([]any); ok {
			for _, item := range raw {
				m, ok := item.(map[string]any)
				if !ok {
					continue
				}
				id := trim(m["id"])
				if id == "" {
					id = newID(models.PrefixPerson)
				}
				status := trim(m["status"])
				if status == "" {
					status = "active"
				}
				specJSON := "{}"
				if spec, ok := m["specJson"].(map[string]any); ok {
					if b, err := json.Marshal(spec); err == nil {
						specJSON = string(b)
					}
				}
				_, err := tx.Exec(ctx, `
					INSERT INTO people (id, tenant_id, org_unit_id, status, given_name, family_name, email, phone, title, avatar_url, spec_json, created_at, updated_at)
					VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $12)
					ON CONFLICT (id) DO UPDATE SET
						org_unit_id=EXCLUDED.org_unit_id, status=EXCLUDED.status, given_name=EXCLUDED.given_name,
						family_name=EXCLUDED.family_name, email=EXCLUDED.email, phone=EXCLUDED.phone,
						title=EXCLUDED.title, avatar_url=EXCLUDED.avatar_url, spec_json=EXCLUDED.spec_json, updated_at=EXCLUDED.updated_at
				`, id, tenant, trim(m["orgUnitId"]), status, trim(m["givenName"]), trim(m["familyName"]),
					trim(m["email"]), trim(m["phone"]), trim(m["title"]), trim(m["avatarUrl"]), specJSON, now)
				if err != nil {
					return err
				}
				res["people"]++
			}
		}

		// Import teams
		if raw, ok := body["teams"].([]any); ok {
			for _, item := range raw {
				m, ok := item.(map[string]any)
				if !ok {
					continue
				}
				id := trim(m["id"])
				if id == "" {
					id = newID(models.PrefixTeam)
				}
				specJSON := "{}"
				if spec, ok := m["specJson"].(map[string]any); ok {
					if b, err := json.Marshal(spec); err == nil {
						specJSON = string(b)
					}
				}
				_, err := tx.Exec(ctx, `
					INSERT INTO teams (id, tenant_id, org_unit_id, key, name, description, spec_json, created_at, updated_at)
					VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $8)
					ON CONFLICT (id) DO UPDATE SET
						org_unit_id=EXCLUDED.org_unit_id, key=EXCLUDED.key, name=EXCLUDED.name,
						description=EXCLUDED.description, spec_json=EXCLUDED.spec_json, updated_at=EXCLUDED.updated_at
				`, id, tenant, trim(m["orgUnitId"]), trim(m["key"]), trim(m["name"]), trim(m["description"]), specJSON, now)
				if err != nil {
					return err
				}
				res["teams"]++
			}
		}

		// Import team_memberships
		if raw, ok := body["teamMemberships"].([]any); ok {
			for _, item := range raw {
				m, ok := item.(map[string]any)
				if !ok {
					continue
				}
				id := trim(m["id"])
				if id == "" {
					id = newID(models.PrefixMembership)
				}
				role := trim(m["role"])
				if role == "" {
					role = "member"
				}
				status := trim(m["status"])
				if status == "" {
					status = "active"
				}
				specJSON := "{}"
				if spec, ok := m["specJson"].(map[string]any); ok {
					if b, err := json.Marshal(spec); err == nil {
						specJSON = string(b)
					}
				}
				_, err := tx.Exec(ctx, `
					INSERT INTO team_memberships (id, tenant_id, team_id, person_id, role, status, started_at, ended_at, spec_json, created_at, updated_at)
					VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $10)
					ON CONFLICT (id) DO UPDATE SET
						team_id=EXCLUDED.team_id, person_id=EXCLUDED.person_id, role=EXCLUDED.role,
						status=EXCLUDED.status, started_at=EXCLUDED.started_at, ended_at=EXCLUDED.ended_at,
						spec_json=EXCLUDED.spec_json, updated_at=EXCLUDED.updated_at
				`, id, tenant, trim(m["teamId"]), trim(m["personId"]), role, status, nil, nil, specJSON, now)
				if err != nil {
					return err
				}
				res["teamMemberships"]++
			}
		}

		return nil
	})

	return res, err
}
