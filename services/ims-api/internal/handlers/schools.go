package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/edvirons/ssp/ims/internal/middleware"
	"github.com/edvirons/ssp/ims/internal/store"
	"go.uber.org/zap"
)

type SchoolHandler struct {
	log *zap.Logger
	pg  *store.Postgres
}

func NewSchoolHandler(log *zap.Logger, pg *store.Postgres) *SchoolHandler {
	return &SchoolHandler{log: log, pg: pg}
}

type upsertSchoolReq struct {
	SchoolID   string `json:"schoolId"`
	CountyCode string `json:"countyCode"`
	CountyName string `json:"countyName"`
}

func (h *SchoolHandler) Upsert(w http.ResponseWriter, r *http.Request) {
	var req upsertSchoolReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(req.SchoolID) == "" || strings.TrimSpace(req.CountyCode) == "" {
		http.Error(w, "schoolId and countyCode are required", http.StatusBadRequest)
		return
	}
	tenant := middleware.TenantID(r.Context())

	if err := h.pg.Schools().Upsert(r.Context(), store.NewSchoolProfile(tenant, strings.TrimSpace(req.SchoolID), strings.TrimSpace(req.CountyCode), strings.TrimSpace(req.CountyName))); err != nil {
		http.Error(w, "failed to upsert school", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}
