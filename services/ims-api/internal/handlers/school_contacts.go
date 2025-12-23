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

type SchoolContactsHandler struct {
	log *zap.Logger
	pg  *store.Postgres
}

func NewSchoolContactsHandler(log *zap.Logger, pg *store.Postgres) *SchoolContactsHandler {
	return &SchoolContactsHandler{log: log, pg: pg}
}

type createContactReq struct {
	UserID    string `json:"userId"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	IsPrimary bool   `json:"isPrimary"`
}

func (h *SchoolContactsHandler) Create(w http.ResponseWriter, r *http.Request) {
	schoolID := chi.URLParam(r, "schoolId")
	var req createContactReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(req.Name) == "" {
		http.Error(w, "name required", http.StatusBadRequest)
		return
	}

	tenant := middleware.TenantID(r.Context())
	now := time.Now().UTC()
	c := models.SchoolContact{
		ID:        store.NewID("contact"),
		TenantID:  tenant,
		SchoolID:  schoolID,
		UserID:    strings.TrimSpace(req.UserID),
		Name:      strings.TrimSpace(req.Name),
		Phone:     strings.TrimSpace(req.Phone),
		Email:     strings.TrimSpace(req.Email),
		Role:      h.pg.SchoolContacts().NormalizeRole(req.Role),
		IsPrimary: req.IsPrimary,
		Active:    true,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := h.pg.SchoolContacts().Create(r.Context(), c); err != nil {
		http.Error(w, "failed to create contact", http.StatusInternalServerError)
		return
	}
	if c.IsPrimary {
		_ = h.pg.SchoolContacts().SetPrimary(r.Context(), tenant, schoolID, c.ID)
	}
	writeJSON(w, http.StatusCreated, c)
}

func (h *SchoolContactsHandler) List(w http.ResponseWriter, r *http.Request) {
	schoolID := chi.URLParam(r, "schoolId")
	tenant := middleware.TenantID(r.Context())
	items, err := h.pg.SchoolContacts().List(r.Context(), tenant, schoolID)
	if err != nil {
		http.Error(w, "failed to list contacts", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

type setPrimaryReq struct {
	ContactID string `json:"contactId"`
}

func (h *SchoolContactsHandler) SetPrimary(w http.ResponseWriter, r *http.Request) {
	schoolID := chi.URLParam(r, "schoolId")
	var req setPrimaryReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(req.ContactID) == "" {
		http.Error(w, "contactId required", http.StatusBadRequest)
		return
	}
	tenant := middleware.TenantID(r.Context())
	if err := h.pg.SchoolContacts().SetPrimary(r.Context(), tenant, schoolID, strings.TrimSpace(req.ContactID)); err != nil {
		http.Error(w, "failed to set primary", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}
