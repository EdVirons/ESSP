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

type BOQHandler struct {
	log *zap.Logger
	pg  *store.Postgres
}

func NewBOQHandler(log *zap.Logger, pg *store.Postgres) *BOQHandler {
	return &BOQHandler{log: log, pg: pg}
}

type addBOQReq struct {
	Category string `json:"category"`
	Description string `json:"description"`
	PartID string `json:"partId"`
	Qty int64 `json:"qty"`
	Unit string `json:"unit"`
	EstimatedCostCents int64 `json:"estimatedCostCents"`
	Approved bool `json:"approved"`
}

func (h *BOQHandler) AddItem(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "id")
	var req addBOQReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil { http.Error(w, "invalid json", http.StatusBadRequest); return }
	if strings.TrimSpace(req.Category) == "" || strings.TrimSpace(req.Description) == "" {
		http.Error(w, "category and description required", http.StatusBadRequest); return
	}
	tenant := middleware.TenantID(r.Context())
	now := time.Now().UTC()
	b := models.BOQItem{
		ID: store.NewID("boq"),
		TenantID: tenant,
		ProjectID: projectID,
		Category: strings.TrimSpace(req.Category),
		Description: strings.TrimSpace(req.Description),
		PartID: strings.TrimSpace(req.PartID),
		Qty: req.Qty,
		Unit: strings.TrimSpace(req.Unit),
		EstimatedCostCents: req.EstimatedCostCents,
		Approved: req.Approved,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := h.pg.BOQ().Create(r.Context(), b); err != nil {
		http.Error(w, "failed to add boq item", http.StatusInternalServerError); return
	}
	writeJSON(w, http.StatusCreated, b)
}

func (h *BOQHandler) List(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "id")
	tenant := middleware.TenantID(r.Context())
	limit := parseLimit(r.URL.Query().Get("limit"), 50, 200)
	curT, curID, hasCur := decodeCursor(strings.TrimSpace(r.URL.Query().Get("cursor")))
	var approvedPtr *bool
	if strings.TrimSpace(r.URL.Query().Get("approved")) == "true" {
		v := true; approvedPtr = &v
	} else if strings.TrimSpace(r.URL.Query().Get("approved")) == "false" {
		v := false; approvedPtr = &v
	}

	items, next, err := h.pg.BOQ().List(r.Context(), store.BOQListParams{
		TenantID: tenant, ProjectID: projectID, Approved: approvedPtr,
		Limit: limit, HasCursor: hasCur, CursorCreatedAt: curT, CursorID: curID,
	})
	if err != nil { http.Error(w, "failed to list boq", http.StatusInternalServerError); return }
	writeJSON(w, http.StatusOK, map[string]any{"items": items, "nextCursor": next})
}
