package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/edvirons/ssp/ims/internal/middleware"
	"github.com/edvirons/ssp/ims/internal/models"
	"github.com/edvirons/ssp/ims/internal/store"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type PartsHandler struct {
	log *zap.Logger
	pg  *store.Postgres
}

func NewPartsHandler(log *zap.Logger, pg *store.Postgres) *PartsHandler {
	return &PartsHandler{log: log, pg: pg}
}

type createPartReq struct {
	SKU           string `json:"sku"`
	Name          string `json:"name"`
	Category      string `json:"category"`
	Description   string `json:"description"`
	UnitCostCents int    `json:"unitCostCents"`
	Supplier      string `json:"supplier"`
	SupplierSku   string `json:"supplierSku"`
}

func (h *PartsHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createPartReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(req.SKU) == "" || strings.TrimSpace(req.Name) == "" {
		http.Error(w, "sku and name are required", http.StatusBadRequest)
		return
	}
	tenant := middleware.TenantID(r.Context())
	now := time.Now().UTC()
	p := models.Part{
		ID:            store.NewID("part"),
		TenantID:      tenant,
		SKU:           strings.TrimSpace(req.SKU),
		Name:          strings.TrimSpace(req.Name),
		Category:      strings.TrimSpace(req.Category),
		Description:   strings.TrimSpace(req.Description),
		UnitCostCents: req.UnitCostCents,
		Supplier:      strings.TrimSpace(req.Supplier),
		SupplierSku:   strings.TrimSpace(req.SupplierSku),
		Active:        true,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	if err := h.pg.Parts().Create(r.Context(), p); err != nil {
		h.log.Error("failed to create part", zap.Error(err))
		http.Error(w, "failed to create part", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, p)
}

func (h *PartsHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	tenant := middleware.TenantID(r.Context())
	p, err := h.pg.Parts().GetByID(r.Context(), tenant, id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, p)
}

type updatePartReq struct {
	Name          *string `json:"name"`
	Category      *string `json:"category"`
	Description   *string `json:"description"`
	UnitCostCents *int    `json:"unitCostCents"`
	Supplier      *string `json:"supplier"`
	SupplierSku   *string `json:"supplierSku"`
	Active        *bool   `json:"active"`
}

func (h *PartsHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	tenant := middleware.TenantID(r.Context())

	// Get existing part
	p, err := h.pg.Parts().GetByID(r.Context(), tenant, id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	// Parse update request
	var req updatePartReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	// Apply updates
	if req.Name != nil {
		p.Name = strings.TrimSpace(*req.Name)
	}
	if req.Category != nil {
		p.Category = strings.TrimSpace(*req.Category)
	}
	if req.Description != nil {
		p.Description = strings.TrimSpace(*req.Description)
	}
	if req.UnitCostCents != nil {
		p.UnitCostCents = *req.UnitCostCents
	}
	if req.Supplier != nil {
		p.Supplier = strings.TrimSpace(*req.Supplier)
	}
	if req.SupplierSku != nil {
		p.SupplierSku = strings.TrimSpace(*req.SupplierSku)
	}
	if req.Active != nil {
		p.Active = *req.Active
	}
	p.UpdatedAt = time.Now().UTC()

	if err := h.pg.Parts().Update(r.Context(), p); err != nil {
		h.log.Error("failed to update part", zap.Error(err))
		http.Error(w, "failed to update part", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, p)
}

func (h *PartsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	tenant := middleware.TenantID(r.Context())

	if err := h.pg.Parts().Delete(r.Context(), tenant, id); err != nil {
		if err.Error() == "not found" {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		h.log.Error("failed to delete part", zap.Error(err))
		http.Error(w, "failed to delete part", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *PartsHandler) List(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	category := strings.TrimSpace(r.URL.Query().Get("category"))
	limit := parseLimit(r.URL.Query().Get("limit"), 50, 200)
	curT, curID, hasCur := decodeCursor(strings.TrimSpace(r.URL.Query().Get("cursor")))

	// Parse active filter
	var active *bool
	if activeStr := r.URL.Query().Get("active"); activeStr != "" {
		val := activeStr == "true"
		active = &val
	}

	items, next, err := h.pg.Parts().List(r.Context(), store.PartListParams{
		TenantID:        tenant,
		Q:               q,
		Category:        category,
		Active:          active,
		Limit:           limit,
		HasCursor:       hasCur,
		CursorCreatedAt: curT,
		CursorID:        curID,
	})
	if err != nil {
		h.log.Error("failed to list parts", zap.Error(err))
		http.Error(w, "failed to list parts", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items, "nextCursor": next})
}

func (h *PartsHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())
	categories, err := h.pg.Parts().GetCategories(r.Context(), tenant)
	if err != nil {
		h.log.Error("failed to get categories", zap.Error(err))
		http.Error(w, "failed to get categories", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": categories})
}

func (h *PartsHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())
	total, err := h.pg.Parts().Count(r.Context(), tenant)
	if err != nil {
		h.log.Error("failed to count parts", zap.Error(err))
		http.Error(w, "failed to get stats", http.StatusInternalServerError)
		return
	}
	byCategory, err := h.pg.Parts().CountByCategory(r.Context(), tenant)
	if err != nil {
		h.log.Error("failed to count by category", zap.Error(err))
		http.Error(w, "failed to get stats", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"total":      total,
		"byCategory": byCategory,
	})
}
