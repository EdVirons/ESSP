package api

import (
	"context"
	"errors"
	"time"

	"github.com/edvirons/ssp/ssot_parts/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func exportAll(ctx context.Context, db *pgxpool.Pool, tenant string) (models.ExportPayload, error) {
	p := models.ExportPayload{Version: "1", GeneratedAt: time.Now().UTC()}
	rows, err := db.Query(ctx, `SELECT id, tenant_id, name, category, puk, spec_json, created_at, updated_at FROM parts WHERE tenant_id=$1 ORDER BY name`, tenant)
	if err != nil { return p, err }
	for rows.Next() {
		var x models.Part
		if err := rows.Scan(&x.ID,&x.TenantID,&x.Name,&x.Category,&x.Puk,&x.SpecJSON,&x.CreatedAt,&x.UpdatedAt); err != nil { rows.Close(); return p, err }
		p.Parts = append(p.Parts, x)
	}
	rows.Close()

	rows, err = db.Query(ctx, `SELECT id, tenant_id, part_id, device_model_id, created_at FROM part_compatibility WHERE tenant_id=$1 ORDER BY created_at DESC`, tenant)
	if err != nil { return p, err }
	for rows.Next() {
		var x models.PartCompatibility
		if err := rows.Scan(&x.ID,&x.TenantID,&x.PartID,&x.DeviceModelID,&x.CreatedAt); err != nil { rows.Close(); return p, err }
		p.Compatibility = append(p.Compatibility, x)
	}
	rows.Close()

	rows, err = db.Query(ctx, `SELECT id, tenant_id, part_id, vendor_id, sku, unit_price_cents, currency, lead_time_days, created_at, updated_at FROM vendor_skus WHERE tenant_id=$1 ORDER BY created_at DESC`, tenant)
	if err != nil { return p, err }
	for rows.Next() {
		var x models.VendorSKU
		if err := rows.Scan(&x.ID,&x.TenantID,&x.PartID,&x.VendorID,&x.SKU,&x.UnitPriceCents,&x.Currency,&x.LeadTimeDays,&x.CreatedAt,&x.UpdatedAt); err != nil { rows.Close(); return p, err }
		p.VendorSKUs = append(p.VendorSKUs, x)
	}
	rows.Close()
	return p, nil
}

func importAll(ctx context.Context, db *pgxpool.Pool, tenant string, body map[string]any) (map[string]any, error) {
	parts, _ := body["parts"].([]any)
	compat, _ := body["compatibility"].([]any)
	vskus, _ := body["vendorSkus"].([]any)
	if len(parts)==0 && len(compat)==0 && len(vskus)==0 { return nil, errors.New("no ssot data provided") }
	res := map[string]any{"parts": 0, "compatibility": 0, "vendorSkus": 0}

	err := withTx(ctx, db, func(tx pgx.Tx) error {
	now := time.Now().UTC()

	for _, it := range parts {
		m, ok := it.(map[string]any); if !ok { continue }
		id := trim(m["id"]); if id=="" { id=newID("part") }
		name := trim(m["name"]); if name=="" { continue }
		cat := trim(m["category"]); if cat=="" { cat="misc" }
		puk := trim(m["puk"])
		spec := trim(m["specJson"]); if spec=="" { spec="{}" }
		_, err := tx.Exec(ctx, `
			INSERT INTO parts (id, tenant_id, name, category, puk, spec_json, created_at, updated_at)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$7)
			ON CONFLICT (id) DO UPDATE SET name=EXCLUDED.name, category=EXCLUDED.category, puk=EXCLUDED.puk, spec_json=EXCLUDED.spec_json, updated_at=$7
		`, id, tenant, name, cat, puk, spec, now)
		if err != nil { return err }
		res["parts"] = res["parts"].(int) + 1
	}

	for _, it := range compat {
		m, ok := it.(map[string]any); if !ok { continue }
		id := trim(m["id"]); if id=="" { id=newID("pc") }
		partID := trim(m["partId"]); dmid := trim(m["deviceModelId"])
		if partID=="" || dmid=="" { continue }
		_, err := tx.Exec(ctx, `
			INSERT INTO part_compatibility (id, tenant_id, part_id, device_model_id, created_at)
			VALUES ($1,$2,$3,$4,$5)
			ON CONFLICT (id) DO UPDATE SET part_id=EXCLUDED.part_id, device_model_id=EXCLUDED.device_model_id
		`, id, tenant, partID, dmid, now)
		if err != nil { return err }
		res["compatibility"] = res["compatibility"].(int) + 1
	}

	for _, it := range vskus {
		m, ok := it.(map[string]any); if !ok { continue }
		id := trim(m["id"]); if id=="" { id=newID("vsku") }
		partID := trim(m["partId"]); if partID=="" { continue }
		vendor := trim(m["vendorId"]); sku := trim(m["sku"])
		if vendor=="" || sku=="" { continue }
		unit := int64(0)
		if v, ok := m["unitPriceCents"].(float64); ok { unit=int64(v) }
		cur := trim(m["currency"]); if cur=="" { cur="KES" }
		lt := 7
		if v, ok := m["leadTimeDays"].(float64); ok { lt=int(v) }
		_, err := tx.Exec(ctx, `
			INSERT INTO vendor_skus (id, tenant_id, part_id, vendor_id, sku, unit_price_cents, currency, lead_time_days, created_at, updated_at)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$9)
			ON CONFLICT (id) DO UPDATE SET vendor_id=EXCLUDED.vendor_id, sku=EXCLUDED.sku, unit_price_cents=EXCLUDED.unit_price_cents,
				currency=EXCLUDED.currency, lead_time_days=EXCLUDED.lead_time_days, updated_at=$9
		`, id, tenant, partID, vendor, sku, unit, cur, lt, now)
		if err != nil { return err }
		res["vendorSkus"] = res["vendorSkus"].(int) + 1
	}
	return nil
	})
	return res, err
}
