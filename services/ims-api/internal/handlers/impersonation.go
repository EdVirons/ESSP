package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/edvirons/ssp/ims/internal/middleware"
	"github.com/edvirons/ssp/ims/internal/store"
	"go.uber.org/zap"
)

// ImpersonationHandler handles impersonation-related requests
type ImpersonationHandler struct {
	log *zap.Logger
	pg  *store.Postgres
}

// NewImpersonationHandler creates a new impersonation handler
func NewImpersonationHandler(log *zap.Logger, pg *store.Postgres) *ImpersonationHandler {
	return &ImpersonationHandler{log: log, pg: pg}
}

// ListImpersonatableUsers returns users that can be impersonated
// GET /v1/impersonate/users
func (h *ImpersonationHandler) ListImpersonatableUsers(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())

	users, err := h.pg.SchoolContacts().ListImpersonatableUsers(r.Context(), tenant, 100)
	if err != nil {
		h.log.Error("failed to list impersonatable users", zap.Error(err))
		http.Error(w, "failed to list users", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"items": users})
}

// validateImpersonationReq is the request body for validation
type validateImpersonationReq struct {
	TargetUserID string `json:"targetUserId"`
}

// validateImpersonationResp is the response for validation
type validateImpersonationResp struct {
	Valid    bool     `json:"valid"`
	UserID   string   `json:"userId"`
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Schools  []string `json:"schools"`
	Error    string   `json:"error,omitempty"`
}

// ValidateImpersonation validates that a user can be impersonated
// POST /v1/impersonate/validate
func (h *ImpersonationHandler) ValidateImpersonation(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())

	var req validateImpersonationReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, validateImpersonationResp{
			Valid: false,
			Error: "invalid request body",
		})
		return
	}

	if req.TargetUserID == "" {
		writeJSON(w, http.StatusBadRequest, validateImpersonationResp{
			Valid: false,
			Error: "targetUserId is required",
		})
		return
	}

	// Look up the school contact by user ID
	contact, err := h.pg.SchoolContacts().GetByUserID(r.Context(), tenant, req.TargetUserID)
	if err != nil {
		h.log.Info("impersonation target not found",
			zap.String("targetUserId", req.TargetUserID),
			zap.Error(err))
		writeJSON(w, http.StatusOK, validateImpersonationResp{
			Valid: false,
			Error: "user not found or not a school contact",
		})
		return
	}

	// Get all schools for this user
	schools, err := h.pg.SchoolContacts().ListSchoolsByUserID(r.Context(), tenant, req.TargetUserID)
	if err != nil {
		h.log.Error("failed to get user schools", zap.Error(err))
		writeJSON(w, http.StatusOK, validateImpersonationResp{
			Valid: false,
			Error: "failed to get user schools",
		})
		return
	}

	writeJSON(w, http.StatusOK, validateImpersonationResp{
		Valid:   true,
		UserID:  contact.UserID,
		Name:    contact.Name,
		Email:   contact.Email,
		Schools: schools,
	})
}

// LoadImpersonationTarget is an ImpersonationLoader function that loads target user info
func (h *ImpersonationHandler) LoadImpersonationTarget(ctx context.Context, tenantID, targetUserID string) (*middleware.ImpersonationTarget, error) {
	// Look up the school contact by user ID
	contact, err := h.pg.SchoolContacts().GetByUserID(ctx, tenantID, targetUserID)
	if err != nil {
		return nil, err
	}

	// Get all schools for this user
	schools, err := h.pg.SchoolContacts().ListSchoolsByUserID(ctx, tenantID, targetUserID)
	if err != nil {
		return nil, err
	}

	return &middleware.ImpersonationTarget{
		UserID:   contact.UserID,
		Email:    contact.Email,
		Roles:    []string{"ssp_school_contact"}, // School contacts always have this role
		Schools:  schools,
		TenantID: tenantID,
	}, nil
}
