package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/edvirons/ssp/ims/internal/middleware"
	"github.com/edvirons/ssp/ims/internal/models"
	"github.com/edvirons/ssp/ims/internal/store"
	"go.uber.org/zap"
)

type InventoryHandler struct {
	log *zap.Logger
	pg  *store.Postgres
}

func NewInventoryHandler(log *zap.Logger, pg *store.Postgres) *InventoryHandler {
	return &InventoryHandler{log: log, pg: pg}
}

type upsertInventoryReq struct {
	ServiceShopID    string `json:"serviceShopId"`
	PartID           string `json:"partId"`
	QtyAvailable     int64  `json:"qtyAvailable"`
	ReorderThreshold int64  `json:"reorderThreshold"`
}

func (h *InventoryHandler) Upsert(w http.ResponseWriter, r *http.Request) {
	var req upsertInventoryReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(req.ServiceShopID) == "" || strings.TrimSpace(req.PartID) == "" {
		http.Error(w, "serviceShopId and partId are required", http.StatusBadRequest)
		return
	}
	tenant := middleware.TenantID(r.Context())
	item := models.InventoryItem{
		ID:               store.NewID("inv"),
		TenantID:         tenant,
		ServiceShopID:    strings.TrimSpace(req.ServiceShopID),
		PartID:           strings.TrimSpace(req.PartID),
		QtyAvailable:     req.QtyAvailable,
		QtyReserved:      0,
		ReorderThreshold: req.ReorderThreshold,
		UpdatedAt:        time.Now().UTC(),
	}
	if err := h.pg.Inventory().Upsert(r.Context(), item); err != nil {
		http.Error(w, "failed to upsert inventory", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *InventoryHandler) List(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())
	shopID := strings.TrimSpace(r.URL.Query().Get("serviceShopId"))
	partID := strings.TrimSpace(r.URL.Query().Get("partId"))
	limit := parseLimit(r.URL.Query().Get("limit"), 50, 200)

	curT, curID, hasCur := decodeCursor(strings.TrimSpace(r.URL.Query().Get("cursor")))

	items, next, err := h.pg.Inventory().List(r.Context(), store.InventoryListParams{
		TenantID: tenant, ShopID: shopID, PartID: partID, Limit: limit,
		HasCursor: hasCur, CursorUpdatedAt: curT, CursorID: curID,
	})
	if err != nil {
		http.Error(w, "failed to list inventory", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items, "nextCursor": next})
}
