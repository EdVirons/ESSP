// bom_consume.go provides Consume and Release handlers for BOM inventory operations.
package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/edvirons/ssp/ims/internal/middleware"
	"github.com/edvirons/ssp/ims/internal/store"
	"github.com/go-chi/chi/v5"
)

type consumeReq struct {
	QtyUsed int64 `json:"qtyUsed"`
}

// Consume marks qty as used and decrements inventory (reserved + on-hand).
func (h *BOMHandler) Consume(w http.ResponseWriter, r *http.Request) {
	woID := chi.URLParam(r, "id")
	itemID := chi.URLParam(r, "itemId")

	var req consumeReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if req.QtyUsed <= 0 {
		http.Error(w, "qtyUsed>0 required", http.StatusBadRequest)
		return
	}

	tenant := middleware.TenantID(r.Context())
	school := middleware.SchoolID(r.Context())
	_ = woID // validated through item

	item, err := h.pg.WorkOrderParts().GetByID(r.Context(), tenant, school, itemID)
	if err != nil {
		http.Error(w, "bom item not found", http.StatusNotFound)
		return
	}

	now := time.Now().UTC()
	err = h.pgTx(r.Context(), func(ctx context.Context, tx store.Tx) error {
		// ensure belongs to wo
		if item.WorkOrderID != woID {
			return context.Canceled
		}
		// update bom used with guard (planned-used >= add)
		if err := store.UpdateWorkOrderPartUsedTx(ctx, tx, tenant, school, itemID, req.QtyUsed, now); err != nil {
			return err
		}
		// consume inventory: reduce reserved and available
		_, err := tx.Exec(ctx, `
			UPDATE inventory
			SET qty_reserved = GREATEST(qty_reserved - $4, 0),
			    qty_available = GREATEST(qty_available - $4, 0),
			    updated_at=$5
			WHERE tenant_id=$1 AND service_shop_id=$2 AND part_id=$3
		`, tenant, item.ServiceShopID, item.PartID, req.QtyUsed, now)
		return err
	})
	if err != nil {
		http.Error(w, "failed to consume", http.StatusConflict)
		return
	}
	updated, _ := h.pg.WorkOrderParts().GetByID(r.Context(), tenant, school, itemID)
	writeJSON(w, http.StatusOK, updated)
}

type releaseReq struct {
	Qty int64 `json:"qty"`
}

// Release returns reserved parts back to availability (de-reserve).
func (h *BOMHandler) Release(w http.ResponseWriter, r *http.Request) {
	woID := chi.URLParam(r, "id")
	itemID := chi.URLParam(r, "itemId")

	var req releaseReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if req.Qty <= 0 {
		http.Error(w, "qty>0 required", http.StatusBadRequest)
		return
	}

	tenant := middleware.TenantID(r.Context())
	school := middleware.SchoolID(r.Context())
	item, err := h.pg.WorkOrderParts().GetByID(r.Context(), tenant, school, itemID)
	if err != nil {
		http.Error(w, "bom item not found", http.StatusNotFound)
		return
	}
	if item.WorkOrderID != woID {
		http.Error(w, "bom item does not belong to work order", http.StatusBadRequest)
		return
	}

	now := time.Now().UTC()
	err = h.pgTx(r.Context(), func(ctx context.Context, tx store.Tx) error {
		// reduce planned by qty (but not below used)
		newPlanned := item.QtyPlanned - req.Qty
		if newPlanned < item.QtyUsed {
			return context.Canceled
		}
		if err := store.UpdateWorkOrderPartPlannedTx(ctx, tx, tenant, school, itemID, newPlanned, now); err != nil {
			return err
		}
		// release inventory reserved by qty
		_, err := tx.Exec(ctx, `
			UPDATE inventory
			SET qty_reserved = GREATEST(qty_reserved - $4, 0), updated_at=$5
			WHERE tenant_id=$1 AND service_shop_id=$2 AND part_id=$3
		`, tenant, item.ServiceShopID, item.PartID, req.Qty, now)
		return err
	})
	if err != nil {
		http.Error(w, "failed to release", http.StatusConflict)
		return
	}
	updated, _ := h.pg.WorkOrderParts().GetByID(r.Context(), tenant, school, itemID)
	writeJSON(w, http.StatusOK, updated)
}
