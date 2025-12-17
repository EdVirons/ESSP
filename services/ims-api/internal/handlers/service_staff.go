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

type ServiceStaffHandler struct {
	log *zap.Logger
	pg  *store.Postgres
}

func NewServiceStaffHandler(log *zap.Logger, pg *store.Postgres) *ServiceStaffHandler {
	return &ServiceStaffHandler{log: log, pg: pg}
}

type createStaffReq struct {
	ServiceShopID string          `json:"serviceShopId"`
	UserID        string          `json:"userId"`
	Role          models.StaffRole `json:"role"`
	Phone         string          `json:"phone"`
	Active        bool            `json:"active"`
}

func (h *ServiceStaffHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createStaffReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(req.ServiceShopID) == "" || strings.TrimSpace(req.UserID) == "" || strings.TrimSpace(string(req.Role)) == "" {
		http.Error(w, "serviceShopId, userId, role are required", http.StatusBadRequest)
		return
	}

	tenant := middleware.TenantID(r.Context())
	now := time.Now().UTC()
	st := models.ServiceStaff{
		ID:            store.NewID("staff"),
		TenantID:      tenant,
		ServiceShopID: strings.TrimSpace(req.ServiceShopID),
		UserID:        strings.TrimSpace(req.UserID),
		Role:          req.Role,
		Phone:         strings.TrimSpace(req.Phone),
		Active:        req.Active,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	if err := h.pg.ServiceStaff().Create(r.Context(), st); err != nil {
		http.Error(w, "failed to create staff", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, st)
}

func (h *ServiceStaffHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	tenant := middleware.TenantID(r.Context())
	st, err := h.pg.ServiceStaff().GetByID(r.Context(), tenant, id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, st)
}

func (h *ServiceStaffHandler) List(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())
	shopID := strings.TrimSpace(r.URL.Query().Get("serviceShopId"))
	role := strings.TrimSpace(r.URL.Query().Get("role"))
	activeOnly := strings.TrimSpace(r.URL.Query().Get("active")) == "true"
	limit := parseLimit(r.URL.Query().Get("limit"), 50, 200)
	curT, curID, hasCur := decodeCursor(strings.TrimSpace(r.URL.Query().Get("cursor")))

	items, next, err := h.pg.ServiceStaff().List(r.Context(), store.StaffListParams{
		TenantID: tenant, ShopID: shopID, Role: role, ActiveOnly: activeOnly,
		Limit: limit, HasCursor: hasCur, CursorCreatedAt: curT, CursorID: curID,
	})
	if err != nil {
		http.Error(w, "failed to list staff", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items, "nextCursor": next})
}

type updateStaffReq struct {
	ServiceShopID *string           `json:"serviceShopId"`
	Role          *models.StaffRole `json:"role"`
	Phone         *string           `json:"phone"`
	Active        *bool             `json:"active"`
}

func (h *ServiceStaffHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	tenant := middleware.TenantID(r.Context())

	// Get existing staff
	st, err := h.pg.ServiceStaff().GetByID(r.Context(), tenant, id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	// Parse update request
	var req updateStaffReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	// Apply updates
	if req.ServiceShopID != nil {
		st.ServiceShopID = strings.TrimSpace(*req.ServiceShopID)
	}
	if req.Role != nil {
		st.Role = *req.Role
	}
	if req.Phone != nil {
		st.Phone = strings.TrimSpace(*req.Phone)
	}
	if req.Active != nil {
		st.Active = *req.Active
	}
	st.UpdatedAt = time.Now().UTC()

	if err := h.pg.ServiceStaff().Update(r.Context(), st); err != nil {
		h.log.Error("failed to update staff", zap.Error(err))
		http.Error(w, "failed to update staff", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, st)
}

func (h *ServiceStaffHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	tenant := middleware.TenantID(r.Context())

	if err := h.pg.ServiceStaff().Delete(r.Context(), tenant, id); err != nil {
		if err.Error() == "not found" {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		h.log.Error("failed to delete staff", zap.Error(err))
		http.Error(w, "failed to delete staff", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *ServiceStaffHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())
	stats, err := h.pg.ServiceStaff().GetStats(r.Context(), tenant)
	if err != nil {
		h.log.Error("failed to get staff stats", zap.Error(err))
		http.Error(w, "failed to get stats", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, stats)
}
