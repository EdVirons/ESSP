// bom_suggest.go provides the Suggest handler for BOM compatible parts lookup.
package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/edvirons/ssp/ims/internal/lookups"
	"github.com/edvirons/ssp/ims/internal/middleware"
	"github.com/go-chi/chi/v5"
)

type suggestPartsResp struct {
	DeviceModelID string                `json:"deviceModelId"`
	Count         int                   `json:"count"`
	Items         []lookups.PartSummary `json:"items"`
}

// Suggest returns compatible parts for the device model on this work order.
// Optional query params:
// - q: substring filter over part name/category/puk
// - limit: max number (default 25, max 200)
func (h *BOMHandler) Suggest(w http.ResponseWriter, r *http.Request) {
	woID := chi.URLParam(r, "id")
	tenant := middleware.TenantID(r.Context())
	school := middleware.SchoolID(r.Context())

	wo, err := h.pg.WorkOrders().GetByID(r.Context(), tenant, school, woID)
	if err != nil {
		http.Error(w, "work order not found", http.StatusNotFound)
		return
	}

	lk := lookups.New(h.pg.RawPool())
	dv, err := lk.DeviceByID(r.Context(), tenant, wo.DeviceID)
	if err != nil || dv == nil || strings.TrimSpace(dv.ModelID) == "" {
		http.Error(w, "device model not resolved (ssot snapshot missing?)", http.StatusFailedDependency)
		return
	}

	px, err := lk.LoadPartsExport(r.Context(), tenant)
	if err != nil {
		http.Error(w, "parts snapshot missing", http.StatusFailedDependency)
		return
	}

	q := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("q")))
	limit := int64(25)
	if v := strings.TrimSpace(r.URL.Query().Get("limit")); v != "" {
		if n, err := parseSuggestLimit(v); err == nil {
			limit = n
		}
	}
	if limit <= 0 {
		limit = 25
	}
	if limit > 200 {
		limit = 200
	}

	partByID := map[string]lookups.PartSummary{}
	for _, p := range px.Parts {
		partByID[p.ID] = lookups.PartSummary{ID: p.ID, Name: p.Name, Category: p.Category, PUK: p.PUK}
	}

	seen := map[string]bool{}
	items := []lookups.PartSummary{}
	for _, c := range px.Compatibility {
		if c.DeviceModelID != dv.ModelID {
			continue
		}
		if seen[c.PartID] {
			continue
		}
		seen[c.PartID] = true
		ps, ok := partByID[c.PartID]
		if !ok {
			continue
		}
		if q != "" {
			hay := strings.ToLower(ps.Name + " " + ps.Category + " " + ps.PUK)
			if !strings.Contains(hay, q) {
				continue
			}
		}
		items = append(items, ps)
		if int64(len(items)) >= limit {
			break
		}
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(suggestPartsResp{
		DeviceModelID: dv.ModelID,
		Count:         len(items),
		Items:         items,
	})
}

// parseSuggestLimit parses a string to int64 for limit parameter.
func parseSuggestLimit(s string) (int64, error) {
	if s == "" {
		return 0, errors.New("empty")
	}
	var n int64
	for _, ch := range s {
		if ch < '0' || ch > '9' {
			return 0, errors.New("invalid int")
		}
		n = n*10 + int64(ch-'0')
	}
	return n, nil
}
