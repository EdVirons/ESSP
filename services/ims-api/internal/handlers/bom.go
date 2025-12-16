// Package handlers provides HTTP handlers for BOM (Bill of Materials) operations.
// This file contains AddItem and List handlers. Additional handlers are in:
// - bom_consume.go: Consume and Release operations
// - bom_suggest.go: Suggest compatible parts
package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/edvirons/ssp/ims/internal/lookups"
	"github.com/edvirons/ssp/ims/internal/middleware"
	"github.com/edvirons/ssp/ims/internal/models"
	"github.com/edvirons/ssp/ims/internal/store"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// BOMHandler handles Bill of Materials operations for work orders.
type BOMHandler struct {
	log *zap.Logger
	pg  *store.Postgres
}

// NewBOMHandler creates a new BOM handler.
func NewBOMHandler(log *zap.Logger, pg *store.Postgres) *BOMHandler {
	return &BOMHandler{log: log, pg: pg}
}

type addBOMItemReq struct {
	PartID     string `json:"partId"`
	QtyPlanned int64  `json:"qtyPlanned"`
}

// AddItem creates a BOM item and reserves inventory in the work order's service shop.
func (h *BOMHandler) AddItem(w http.ResponseWriter, r *http.Request) {
	woID := chi.URLParam(r, "id")
	var req addBOMItemReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(req.PartID) == "" || req.QtyPlanned <= 0 {
		http.Error(w, "partId and qtyPlanned>0 required", http.StatusBadRequest)
		return
	}

	tenant := middleware.TenantID(r.Context())
	school := middleware.SchoolID(r.Context())

	wo, err := h.pg.WorkOrders().GetByID(r.Context(), tenant, school, woID)
	if err != nil {
		http.Error(w, "work order not found", http.StatusNotFound)
		return
	}
	if strings.TrimSpace(wo.ServiceShopID) == "" {
		http.Error(w, "work order has no serviceShopId", http.StatusBadRequest)
		return
	}

	lk := lookups.New(h.pg.RawPool())
	dev, _ := lk.DeviceByID(r.Context(), tenant, wo.DeviceID)
	part, _ := lk.PartByID(r.Context(), tenant, strings.TrimSpace(req.PartID))
	deviceModelID := ""
	if dev != nil {
		deviceModelID = dev.ModelID
	}
	isCompat := true
	if deviceModelID != "" {
		ok, _ := lk.IsPartCompatibleWithDeviceModel(r.Context(), tenant, strings.TrimSpace(req.PartID), deviceModelID)
		isCompat = ok
	}
	allowIncompat := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("allowIncompatible")))
	if (allowIncompat == "" || allowIncompat == "0" || allowIncompat == "false") && enforceCompatibility() && !isCompat {
		http.Error(w, "part is not compatible with device model", http.StatusConflict)
		return
	}

	now := time.Now().UTC()
	item := models.WorkOrderPart{
		ID:            store.NewID("bom"),
		TenantID:      tenant,
		SchoolID:      school,
		WorkOrderID:   woID,
		ServiceShopID: wo.ServiceShopID,
		PartID:        strings.TrimSpace(req.PartID),
		PartName:      getPartName(part),
		PartPUK:       getPartPUK(part),
		PartCategory:  getPartCategory(part),
		DeviceModelID: deviceModelID,
		IsCompatible:  isCompat,
		QtyPlanned:    req.QtyPlanned,
		QtyUsed:       0,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	// Transaction: reserve inventory + create BOM item
	err = h.pgTx(r.Context(), func(ctx context.Context, tx store.Tx) error {
		if _, err := tx.Exec(ctx, `
			UPDATE inventory
			SET qty_reserved = qty_reserved + $4, updated_at = $5
			WHERE tenant_id=$1 AND service_shop_id=$2 AND part_id=$3
			  AND (qty_available - qty_reserved) >= $4
		`, tenant, wo.ServiceShopID, item.PartID, item.QtyPlanned, now); err != nil {
			return err
		}
		return store.CreateWorkOrderPartTx(ctx, tx, item)
	})
	if err != nil {
		http.Error(w, "failed to reserve inventory or create bom item", http.StatusConflict)
		return
	}

	writeJSON(w, http.StatusCreated, item)
}

// List returns BOM items for a work order with pagination.
func (h *BOMHandler) List(w http.ResponseWriter, r *http.Request) {
	woID := chi.URLParam(r, "id")
	tenant := middleware.TenantID(r.Context())
	school := middleware.SchoolID(r.Context())

	limit := parseLimit(r.URL.Query().Get("limit"), 50, 200)
	curT, curID, hasCur := decodeCursor(strings.TrimSpace(r.URL.Query().Get("cursor")))

	items, next, err := h.pg.WorkOrderParts().List(r.Context(), store.WorkOrderPartListParams{
		TenantID: tenant, SchoolID: school, WorkOrderID: woID,
		Limit: limit, HasCursor: hasCur, CursorCreatedAt: curT, CursorID: curID,
	})
	if err != nil {
		http.Error(w, "failed to list bom", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items, "nextCursor": next})
}

// pgTx runs a pg transaction using the underlying pool.
func (h *BOMHandler) pgTx(ctx context.Context, fn func(context.Context, store.Tx) error) error {
	tx, err := h.pg.RawPool().Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := fn(ctx, tx); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// enforceCompatibility returns whether BOM compatibility enforcement is enabled.
func enforceCompatibility() bool {
	v := strings.ToLower(strings.TrimSpace(os.Getenv("BOM_ENFORCE_COMPATIBILITY")))
	if v == "" {
		return true // default ON
	}
	return !(v == "0" || v == "false" || v == "no")
}

// Part attribute helpers
func getPartName(part *lookups.PartSummary) string {
	if part != nil {
		return part.Name
	}
	return ""
}

func getPartPUK(part *lookups.PartSummary) string {
	if part != nil {
		return part.PUK
	}
	return ""
}

func getPartCategory(part *lookups.PartSummary) string {
	if part != nil {
		return part.Category
	}
	return ""
}
