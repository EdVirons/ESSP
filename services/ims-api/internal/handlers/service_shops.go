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

type ServiceShopHandler struct {
	log *zap.Logger
	pg  *store.Postgres
}

func NewServiceShopHandler(log *zap.Logger, pg *store.Postgres) *ServiceShopHandler {
	return &ServiceShopHandler{log: log, pg: pg}
}

type createShopReq struct {
	CountyCode    string `json:"countyCode"`
	CountyName    string `json:"countyName"`
	SubCountyCode string `json:"subCountyCode"`
	SubCountyName string `json:"subCountyName"`
	CoverageLevel string `json:"coverageLevel"`
	Name          string `json:"name"`
	Location      string `json:"location"`
	Active        bool   `json:"active"`
}

func (h *ServiceShopHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createShopReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(req.CountyCode) == "" || strings.TrimSpace(req.Name) == "" {
		http.Error(w, "countyCode and name are required", http.StatusBadRequest)
		return
	}

	tenant := middleware.TenantID(r.Context())
	now := time.Now().UTC()
	cov := strings.TrimSpace(req.CoverageLevel)
	if cov == "" {
		cov = "county"
	}
	shop := models.ServiceShop{
		ID:            store.NewID("shop"),
		TenantID:      tenant,
		CountyCode:    strings.TrimSpace(req.CountyCode),
		CountyName:    strings.TrimSpace(req.CountyName),
		SubCountyCode: strings.TrimSpace(req.SubCountyCode),
		SubCountyName: strings.TrimSpace(req.SubCountyName),
		CoverageLevel: cov,
		Name:          strings.TrimSpace(req.Name),
		Location:      strings.TrimSpace(req.Location),
		Active:        req.Active,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if err := h.pg.ServiceShops().Create(r.Context(), shop); err != nil {
		http.Error(w, "failed to create shop (maybe county already has one?)", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, shop)
}

func (h *ServiceShopHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	tenant := middleware.TenantID(r.Context())
	shop, err := h.pg.ServiceShops().GetByID(r.Context(), tenant, id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, shop)
}

func (h *ServiceShopHandler) List(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())
	county := strings.TrimSpace(r.URL.Query().Get("countyCode"))
	activeOnly := strings.TrimSpace(r.URL.Query().Get("active")) == "true"
	limit := parseLimit(r.URL.Query().Get("limit"), 50, 200)
	curT, curID, hasCur := decodeCursor(strings.TrimSpace(r.URL.Query().Get("cursor")))

	items, next, err := h.pg.ServiceShops().List(r.Context(), store.ShopListParams{
		TenantID: tenant, CountyCode: county, ActiveOnly: activeOnly,
		Limit: limit, HasCursor: hasCur, CursorCreatedAt: curT, CursorID: curID,
	})
	if err != nil {
		http.Error(w, "failed to list shops", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items, "nextCursor": next})
}
