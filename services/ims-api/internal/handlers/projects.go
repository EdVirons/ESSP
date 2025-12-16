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

type ProjectsHandler struct {
	log *zap.Logger
	pg  *store.Postgres
}

func NewProjectsHandler(log *zap.Logger, pg *store.Postgres) *ProjectsHandler {
	return &ProjectsHandler{log: log, pg: pg}
}

type createProjectReq struct {
	SchoolID             string `json:"schoolId"`
	ProjectType          string `json:"projectType"`
	StartDate            string `json:"startDate"`
	GoLiveDate           string `json:"goLiveDate"`
	AccountManagerUserID string `json:"accountManagerUserId"`
	Notes                string `json:"notes"`
}

func (h *ProjectsHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createProjectReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(req.SchoolID) == "" {
		http.Error(w, "schoolId required", http.StatusBadRequest)
		return
	}

	// Determine project type (default to full_installation)
	projectType := models.ProjectType(strings.TrimSpace(req.ProjectType))
	if projectType == "" {
		projectType = models.ProjectTypeFullInstallation
	}

	// Validate project type and get config
	config, ok := models.ProjectTypeConfigs[projectType]
	if !ok {
		http.Error(w, "invalid projectType", http.StatusBadRequest)
		return
	}

	tenant := middleware.TenantID(r.Context())
	now := time.Now().UTC()
	p := models.SchoolServiceProject{
		ID:                   store.NewID("proj"),
		TenantID:             tenant,
		SchoolID:             strings.TrimSpace(req.SchoolID),
		ProjectType:          projectType,
		Status:               models.ProjectActive,
		CurrentPhase:         config.DefaultPhase,
		StartDate:            strings.TrimSpace(req.StartDate),
		GoLiveDate:           strings.TrimSpace(req.GoLiveDate),
		AccountManagerUserID: strings.TrimSpace(req.AccountManagerUserID),
		Notes:                strings.TrimSpace(req.Notes),
		CreatedAt:            now,
		UpdatedAt:            now,
	}
	if err := h.pg.Projects().Create(r.Context(), p); err != nil {
		http.Error(w, "failed to create project", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, p)
}

func (h *ProjectsHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	tenant := middleware.TenantID(r.Context())
	p, err := h.pg.Projects().GetByID(r.Context(), tenant, id)
	if err != nil { http.Error(w, "not found", http.StatusNotFound); return }
	writeJSON(w, http.StatusOK, p)
}

func (h *ProjectsHandler) List(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())
	schoolID := strings.TrimSpace(r.URL.Query().Get("schoolId"))
	projectType := strings.TrimSpace(r.URL.Query().Get("projectType"))
	status := strings.TrimSpace(r.URL.Query().Get("status"))
	limit := parseLimit(r.URL.Query().Get("limit"), 50, 200)
	curT, curID, hasCur := decodeCursor(strings.TrimSpace(r.URL.Query().Get("cursor")))

	items, next, err := h.pg.Projects().List(r.Context(), store.ProjectListParams{
		TenantID:        tenant,
		SchoolID:        schoolID,
		ProjectType:     projectType,
		Status:          status,
		Limit:           limit,
		HasCursor:       hasCur,
		CursorCreatedAt: curT,
		CursorID:        curID,
	})
	if err != nil {
		http.Error(w, "failed to list projects", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items, "nextCursor": next})
}

// GetProjectTypes returns all available project type configurations.
func (h *ProjectsHandler) GetProjectTypes(w http.ResponseWriter, r *http.Request) {
	configs := make([]models.ProjectTypeConfig, 0, len(models.ProjectTypeConfigs))
	for _, pt := range models.ValidProjectTypes() {
		configs = append(configs, models.ProjectTypeConfigs[pt])
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": configs})
}

// GetProjectTypeCounts returns the count of projects for each project type.
func (h *ProjectsHandler) GetProjectTypeCounts(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())
	counts, err := h.pg.Projects().CountByType(r.Context(), tenant)
	if err != nil {
		h.log.Error("failed to get project type counts", zap.Error(err))
		http.Error(w, "failed to get counts", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, counts)
}
