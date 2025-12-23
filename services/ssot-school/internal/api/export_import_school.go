package api

import (
	"context"
	"errors"
	"time"

	"github.com/edvirons/ssp/ssot_school/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func exportAll(ctx context.Context, db *pgxpool.Pool, tenant string) (models.ExportPayload, error) {
	p := models.ExportPayload{Version: "1", GeneratedAt: time.Now().UTC()}
	rows, err := db.Query(ctx, `SELECT id, tenant_id, name, code, created_at, updated_at FROM counties WHERE tenant_id=$1 ORDER BY name`, tenant)
	if err != nil {
		return p, err
	}
	for rows.Next() {
		var x models.County
		if err := rows.Scan(&x.ID, &x.TenantID, &x.Name, &x.Code, &x.CreatedAt, &x.UpdatedAt); err != nil {
			rows.Close()
			return p, err
		}
		p.Counties = append(p.Counties, x)
	}
	rows.Close()

	rows, err = db.Query(ctx, `SELECT id, tenant_id, county_id, name, code, created_at, updated_at FROM sub_counties WHERE tenant_id=$1 ORDER BY name`, tenant)
	if err != nil {
		return p, err
	}
	for rows.Next() {
		var x models.SubCounty
		if err := rows.Scan(&x.ID, &x.TenantID, &x.CountyID, &x.Name, &x.Code, &x.CreatedAt, &x.UpdatedAt); err != nil {
			rows.Close()
			return p, err
		}
		p.SubCounties = append(p.SubCounties, x)
	}
	rows.Close()

	rows, err = db.Query(ctx, `SELECT id, tenant_id, name, code, county_id, sub_county_id, level, type, active,
		knec_code, uic, sex, cluster, accommodation, latitude, longitude, created_at, updated_at
		FROM schools WHERE tenant_id=$1 ORDER BY name`, tenant)
	if err != nil {
		return p, err
	}
	for rows.Next() {
		var x models.School
		if err := rows.Scan(&x.ID, &x.TenantID, &x.Name, &x.Code, &x.CountyID, &x.SubCountyID, &x.Level, &x.Type, &x.Active,
			&x.KnecCode, &x.Uic, &x.Sex, &x.Cluster, &x.Accommodation, &x.Latitude, &x.Longitude,
			&x.CreatedAt, &x.UpdatedAt); err != nil {
			rows.Close()
			return p, err
		}
		p.Schools = append(p.Schools, x)
	}
	rows.Close()

	rows, err = db.Query(ctx, `SELECT id, tenant_id, school_id, name, phone, email, role, is_primary, active, created_at, updated_at FROM school_contacts WHERE tenant_id=$1 ORDER BY created_at DESC`, tenant)
	if err != nil {
		return p, err
	}
	for rows.Next() {
		var x models.Contact
		if err := rows.Scan(&x.ID, &x.TenantID, &x.SchoolID, &x.Name, &x.Phone, &x.Email, &x.Role, &x.IsPrimary, &x.Active, &x.CreatedAt, &x.UpdatedAt); err != nil {
			rows.Close()
			return p, err
		}
		p.Contacts = append(p.Contacts, x)
	}
	rows.Close()
	return p, nil
}

func importAll(ctx context.Context, db *pgxpool.Pool, tenant string, body map[string]any) (map[string]any, error) {
	counties, _ := body["counties"].([]any)
	sub, _ := body["subCounties"].([]any)
	schools, _ := body["schools"].([]any)
	contacts, _ := body["contacts"].([]any)

	if len(counties) == 0 && len(sub) == 0 && len(schools) == 0 && len(contacts) == 0 {
		return nil, errors.New("no ssot data provided")
	}

	res := map[string]any{"counties": 0, "subCounties": 0, "schools": 0, "contacts": 0}

	err := withTx(ctx, db, func(tx pgx.Tx) error {
		now := time.Now().UTC()

		for _, it := range counties {
			m, ok := it.(map[string]any)
			if !ok {
				continue
			}
			id := trim(m["id"])
			if id == "" {
				id = newID("county")
			}
			name := trim(m["name"])
			if name == "" {
				continue
			}
			code := trim(m["code"])
			_, err := tx.Exec(ctx, `
			INSERT INTO counties (id, tenant_id, name, code, created_at, updated_at)
			VALUES ($1,$2,$3,$4,$5,$5)
			ON CONFLICT (id) DO UPDATE SET name=EXCLUDED.name, code=EXCLUDED.code, updated_at=$5
		`, id, tenant, name, code, now)
			if err != nil {
				return err
			}
			if c, ok := res["counties"].(int); ok {
				res["counties"] = c + 1
			}
		}

		for _, it := range sub {
			m, ok := it.(map[string]any)
			if !ok {
				continue
			}
			id := trim(m["id"])
			if id == "" {
				id = newID("subcounty")
			}
			name := trim(m["name"])
			if name == "" {
				continue
			}
			countyID := trim(m["countyId"])
			code := trim(m["code"])
			_, err := tx.Exec(ctx, `
			INSERT INTO sub_counties (id, tenant_id, county_id, name, code, created_at, updated_at)
			VALUES ($1,$2,$3,$4,$5,$6,$6)
			ON CONFLICT (id) DO UPDATE SET county_id=EXCLUDED.county_id, name=EXCLUDED.name, code=EXCLUDED.code, updated_at=$6
		`, id, tenant, countyID, name, code, now)
			if err != nil {
				return err
			}
			if c, ok := res["subCounties"].(int); ok {
				res["subCounties"] = c + 1
			}
		}

		for _, it := range schools {
			m, ok := it.(map[string]any)
			if !ok {
				continue
			}
			id := trim(m["id"])
			if id == "" {
				id = newID("school")
			}
			name := trim(m["name"])
			if name == "" {
				continue
			}
			code := trim(m["code"])
			countyID := trim(m["countyId"])
			subID := trim(m["subCountyId"])
			level := trim(m["level"])
			if level == "" {
				level = "Other"
			}
			type_ := trim(m["type"])
			if type_ == "" {
				type_ = "public"
			}
			active := true
			if v, ok := m["active"].(bool); ok {
				active = v
			}
			knecCode := trim(m["knecCode"])
			uic := trim(m["uic"])
			sex := trim(m["sex"])
			cluster := trim(m["cluster"])
			accommodation := trim(m["accommodation"])
			lat := 0.0
			if v, ok := m["latitude"].(float64); ok {
				lat = v
			}
			lng := 0.0
			if v, ok := m["longitude"].(float64); ok {
				lng = v
			}
			_, err := tx.Exec(ctx, `
			INSERT INTO schools (id, tenant_id, name, code, county_id, sub_county_id, level, type, active,
				knec_code, uic, sex, cluster, accommodation, latitude, longitude, created_at, updated_at)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$17)
			ON CONFLICT (id) DO UPDATE SET name=EXCLUDED.name, code=EXCLUDED.code, county_id=EXCLUDED.county_id, sub_county_id=EXCLUDED.sub_county_id,
				level=EXCLUDED.level, type=EXCLUDED.type, active=EXCLUDED.active, knec_code=EXCLUDED.knec_code, uic=EXCLUDED.uic,
				sex=EXCLUDED.sex, cluster=EXCLUDED.cluster, accommodation=EXCLUDED.accommodation, latitude=EXCLUDED.latitude,
				longitude=EXCLUDED.longitude, updated_at=$17
		`, id, tenant, name, code, countyID, subID, level, type_, active, knecCode, uic, sex, cluster, accommodation, lat, lng, now)
			if err != nil {
				return err
			}
			if c, ok := res["schools"].(int); ok {
				res["schools"] = c + 1
			}
		}

		for _, it := range contacts {
			m, ok := it.(map[string]any)
			if !ok {
				continue
			}
			id := trim(m["id"])
			if id == "" {
				id = newID("contact")
			}
			schoolID := trim(m["schoolId"])
			if schoolID == "" {
				continue
			}
			name := trim(m["name"])
			if name == "" {
				continue
			}
			phone := trim(m["phone"])
			email := trim(m["email"])
			role := trim(m["role"])
			if role == "" {
				role = "point_of_contact"
			}
			isPrimary := false
			if v, ok := m["isPrimary"].(bool); ok {
				isPrimary = v
			}
			active := true
			if v, ok := m["active"].(bool); ok {
				active = v
			}

			_, err := tx.Exec(ctx, `
			INSERT INTO school_contacts (id, tenant_id, school_id, name, phone, email, role, is_primary, active, created_at, updated_at)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$10)
			ON CONFLICT (id) DO UPDATE SET school_id=EXCLUDED.school_id, name=EXCLUDED.name, phone=EXCLUDED.phone, email=EXCLUDED.email,
				role=EXCLUDED.role, is_primary=EXCLUDED.is_primary, active=EXCLUDED.active, updated_at=$10
		`, id, tenant, schoolID, name, phone, email, role, isPrimary, active, now)
			if err != nil {
				return err
			}
			if c, ok := res["contacts"].(int); ok {
				res["contacts"] = c + 1
			}
		}
		return nil
	})
	return res, err
}
