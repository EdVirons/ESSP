package handlers

import (
	"encoding/csv"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/edvirons/ssp/ims/internal/middleware"
	"github.com/edvirons/ssp/ims/internal/models"
	"github.com/edvirons/ssp/ims/internal/store"
	"go.uber.org/zap"
)

// Import handles CSV import of parts
func (h *PartsHandler) Import(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())

	// Parse multipart form (max 10MB)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "failed to parse form", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "file is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Read header
	header, err := reader.Read()
	if err != nil {
		http.Error(w, "failed to read CSV header", http.StatusBadRequest)
		return
	}

	// Build column index map
	colIndex := make(map[string]int)
	for i, col := range header {
		colIndex[strings.ToLower(strings.TrimSpace(col))] = i
	}

	// Validate required columns
	if _, ok := colIndex["sku"]; !ok {
		http.Error(w, "missing required column: sku", http.StatusBadRequest)
		return
	}
	if _, ok := colIndex["name"]; !ok {
		http.Error(w, "missing required column: name", http.StatusBadRequest)
		return
	}

	var created, failed int
	var errors []string
	now := time.Now().UTC()

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			failed++
			errors = append(errors, "failed to read row")
			continue
		}

		sku := getCSVValue(record, colIndex, "sku")
		name := getCSVValue(record, colIndex, "name")

		if sku == "" || name == "" {
			failed++
			errors = append(errors, "row missing sku or name")
			continue
		}

		p := models.Part{
			ID:            store.NewID("part"),
			TenantID:      tenant,
			SKU:           sku,
			Name:          name,
			Category:      getCSVValue(record, colIndex, "category"),
			Description:   getCSVValue(record, colIndex, "description"),
			UnitCostCents: parseIntCSV(getCSVValue(record, colIndex, "unitcostcents")),
			Supplier:      getCSVValue(record, colIndex, "supplier"),
			SupplierSku:   getCSVValue(record, colIndex, "suppliersku"),
			Active:        true,
			CreatedAt:     now,
			UpdatedAt:     now,
		}

		if err := h.pg.Parts().Create(r.Context(), p); err != nil {
			failed++
			errors = append(errors, "failed to create: "+sku)
			continue
		}
		created++
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"created": created,
		"failed":  failed,
		"errors":  errors,
	})
}

// Export handles CSV export of parts
func (h *PartsHandler) Export(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())

	// Get all parts (no pagination for export)
	items, _, err := h.pg.Parts().List(r.Context(), store.PartListParams{
		TenantID: tenant,
		Limit:    10000, // Max export limit
	})
	if err != nil {
		h.log.Error("failed to list parts for export", zap.Error(err))
		http.Error(w, "failed to export", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=parts.csv")

	writer := csv.NewWriter(w)
	defer writer.Flush()

	// Write header
	writer.Write([]string{"sku", "name", "category", "description", "unitCostCents", "supplier", "supplierSku", "active", "createdAt", "updatedAt"})

	// Write rows
	for _, p := range items {
		writer.Write([]string{
			p.SKU,
			p.Name,
			p.Category,
			p.Description,
			strconv.Itoa(p.UnitCostCents),
			p.Supplier,
			p.SupplierSku,
			boolToStr(p.Active),
			p.CreatedAt.Format(time.RFC3339),
			p.UpdatedAt.Format(time.RFC3339),
		})
	}
}

// CSV helper functions

func getCSVValue(record []string, colIndex map[string]int, col string) string {
	if idx, ok := colIndex[col]; ok && idx < len(record) {
		return strings.TrimSpace(record[idx])
	}
	return ""
}

func parseIntCSV(s string) int {
	var v int
	for _, c := range s {
		if c >= '0' && c <= '9' {
			v = v*10 + int(c-'0')
		}
	}
	return v
}

func boolToStr(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
