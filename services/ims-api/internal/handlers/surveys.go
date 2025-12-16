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

type SurveysHandler struct {
	log *zap.Logger
	pg  *store.Postgres
}

func NewSurveysHandler(log *zap.Logger, pg *store.Postgres) *SurveysHandler {
	return &SurveysHandler{log: log, pg: pg}
}

type createSurveyReq struct {
	ConductedByUserID string `json:"conductedByUserId"`
	ConductedAt string `json:"conductedAt"` // RFC3339
	Summary string `json:"summary"`
	Risks string `json:"risks"`
}

func (h *SurveysHandler) Create(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "id")
	var req createSurveyReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil { http.Error(w, "invalid json", http.StatusBadRequest); return }
	tenant := middleware.TenantID(r.Context())
	now := time.Now().UTC()

	var conductedAt *time.Time
	if strings.TrimSpace(req.ConductedAt) != "" {
		if t, err := time.Parse(time.RFC3339, strings.TrimSpace(req.ConductedAt)); err == nil {
			conductedAt = &t
		}
	}

	s := models.SiteSurvey{
		ID: store.NewID("survey"),
		TenantID: tenant,
		ProjectID: projectID,
		Status: models.SurveyDraft,
		ConductedByUserID: strings.TrimSpace(req.ConductedByUserID),
		ConductedAt: conductedAt,
		Summary: strings.TrimSpace(req.Summary),
		Risks: strings.TrimSpace(req.Risks),
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := h.pg.Surveys().Create(r.Context(), s); err != nil {
		http.Error(w, "failed to create survey", http.StatusInternalServerError); return
	}
	writeJSON(w, http.StatusCreated, s)
}

func (h *SurveysHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	surveyID := chi.URLParam(r, "surveyId")
	tenant := middleware.TenantID(r.Context())
	s, err := h.pg.Surveys().GetByID(r.Context(), tenant, surveyID)
	if err != nil { http.Error(w, "not found", http.StatusNotFound); return }
	rooms, _ := h.pg.SurveyRooms().List(r.Context(), tenant, surveyID)
	photos, _ := h.pg.SurveyPhotos().List(r.Context(), tenant, surveyID)
	writeJSON(w, http.StatusOK, map[string]any{"survey": s, "rooms": rooms, "photos": photos})
}

func (h *SurveysHandler) List(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "id")
	tenant := middleware.TenantID(r.Context())
	status := strings.TrimSpace(r.URL.Query().Get("status"))
	limit := parseLimit(r.URL.Query().Get("limit"), 50, 200)
	curT, curID, hasCur := decodeCursor(strings.TrimSpace(r.URL.Query().Get("cursor")))

	items, next, err := h.pg.Surveys().List(r.Context(), store.SurveyListParams{
		TenantID: tenant, ProjectID: projectID, Status: status,
		Limit: limit, HasCursor: hasCur, CursorCreatedAt: curT, CursorID: curID,
	})
	if err != nil { http.Error(w, "failed to list surveys", http.StatusInternalServerError); return }
	writeJSON(w, http.StatusOK, map[string]any{"items": items, "nextCursor": next})
}

type addRoomReq struct {
	Name string `json:"name"`
	RoomType string `json:"roomType"`
	Floor string `json:"floor"`
	PowerNotes string `json:"powerNotes"`
	NetworkNotes string `json:"networkNotes"`
}

func (h *SurveysHandler) AddRoom(w http.ResponseWriter, r *http.Request) {
	surveyID := chi.URLParam(r, "surveyId")
	var req addRoomReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil { http.Error(w, "invalid json", http.StatusBadRequest); return }
	if strings.TrimSpace(req.Name) == "" { http.Error(w, "name required", http.StatusBadRequest); return }
	tenant := middleware.TenantID(r.Context())
	now := time.Now().UTC()
	room := models.SurveyRoom{
		ID: store.NewID("room"),
		TenantID: tenant,
		SurveyID: surveyID,
		Name: strings.TrimSpace(req.Name),
		RoomType: strings.TrimSpace(req.RoomType),
		Floor: strings.TrimSpace(req.Floor),
		PowerNotes: strings.TrimSpace(req.PowerNotes),
		NetworkNotes: strings.TrimSpace(req.NetworkNotes),
		CreatedAt: now,
	}
	if err := h.pg.SurveyRooms().Create(r.Context(), room); err != nil {
		http.Error(w, "failed to add room", http.StatusInternalServerError); return
	}
	writeJSON(w, http.StatusCreated, room)
}

type addPhotoReq struct {
	RoomID string `json:"roomId"`
	AttachmentID string `json:"attachmentId"`
	Caption string `json:"caption"`
}

func (h *SurveysHandler) AddPhoto(w http.ResponseWriter, r *http.Request) {
	surveyID := chi.URLParam(r, "surveyId")
	var req addPhotoReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil { http.Error(w, "invalid json", http.StatusBadRequest); return }
	if strings.TrimSpace(req.AttachmentID) == "" { http.Error(w, "attachmentId required", http.StatusBadRequest); return }
	tenant := middleware.TenantID(r.Context())
	now := time.Now().UTC()
	p := models.SurveyPhoto{
		ID: store.NewID("photo"),
		TenantID: tenant,
		SurveyID: surveyID,
		RoomID: strings.TrimSpace(req.RoomID),
		AttachmentID: strings.TrimSpace(req.AttachmentID),
		Caption: strings.TrimSpace(req.Caption),
		CreatedAt: now,
	}
	if err := h.pg.SurveyPhotos().Create(r.Context(), p); err != nil {
		http.Error(w, "failed to add photo", http.StatusInternalServerError); return
	}
	writeJSON(w, http.StatusCreated, p)
}
