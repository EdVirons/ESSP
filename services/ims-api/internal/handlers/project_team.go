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

type ProjectTeamHandler struct {
	log *zap.Logger
	pg  *store.Postgres
}

func NewProjectTeamHandler(log *zap.Logger, pg *store.Postgres) *ProjectTeamHandler {
	return &ProjectTeamHandler{log: log, pg: pg}
}

type addTeamMemberReq struct {
	UserID         string   `json:"userId"`
	UserEmail      string   `json:"userEmail"`
	UserName       string   `json:"userName"`
	Role           string   `json:"role"`
	AssignedPhases []string `json:"assignedPhases"`
	Responsibility string   `json:"responsibility"`
}

func (h *ProjectTeamHandler) AddMember(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "id")
	tenant := middleware.TenantID(r.Context())
	currentUserID := middleware.UserID(r.Context())

	var req addTeamMemberReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.UserID) == "" {
		http.Error(w, "userId required", http.StatusBadRequest)
		return
	}

	// Validate role
	if !models.IsValidTeamMemberRole(req.Role) {
		req.Role = string(models.TeamRoleCollaborator)
	}

	// Check if already a member
	isMember, err := h.pg.ProjectTeam().IsMember(r.Context(), tenant, projectID, req.UserID)
	if err != nil {
		h.log.Error("failed to check membership", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	if isMember {
		http.Error(w, "user is already a team member", http.StatusConflict)
		return
	}

	now := time.Now().UTC()
	phases := make([]models.PhaseType, len(req.AssignedPhases))
	for i, p := range req.AssignedPhases {
		phases[i] = models.PhaseType(p)
	}

	member := models.ProjectTeamMember{
		ID:               store.NewID("ptm"),
		TenantID:         tenant,
		ProjectID:        projectID,
		UserID:           strings.TrimSpace(req.UserID),
		UserEmail:        strings.TrimSpace(req.UserEmail),
		UserName:         strings.TrimSpace(req.UserName),
		Role:             models.TeamMemberRole(req.Role),
		AssignedPhases:   phases,
		Responsibility:   strings.TrimSpace(req.Responsibility),
		AssignedByUserID: currentUserID,
		AssignedAt:       now,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	if err := h.pg.ProjectTeam().AddMember(r.Context(), member); err != nil {
		h.log.Error("failed to add team member", zap.Error(err))
		http.Error(w, "failed to add team member", http.StatusInternalServerError)
		return
	}

	// Create activity for the assignment
	h.logTeamActivity(r, projectID, member.UserID, member.UserName, "added", member.Role)

	writeJSON(w, http.StatusCreated, member)
}

type updateTeamMemberReq struct {
	Role           string   `json:"role"`
	AssignedPhases []string `json:"assignedPhases"`
	Responsibility string   `json:"responsibility"`
}

func (h *ProjectTeamHandler) UpdateMember(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "id")
	memberID := chi.URLParam(r, "memberId")
	tenant := middleware.TenantID(r.Context())

	var req updateTeamMemberReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	// Validate role
	if req.Role != "" && !models.IsValidTeamMemberRole(req.Role) {
		http.Error(w, "invalid role", http.StatusBadRequest)
		return
	}

	// Get existing member
	member, err := h.pg.ProjectTeam().GetMember(r.Context(), tenant, memberID)
	if err != nil {
		http.Error(w, "member not found", http.StatusNotFound)
		return
	}
	if member.ProjectID != projectID {
		http.Error(w, "member not found in this project", http.StatusNotFound)
		return
	}

	role := req.Role
	if role == "" {
		role = string(member.Role)
	}

	if err := h.pg.ProjectTeam().UpdateMember(r.Context(), tenant, memberID, role, req.AssignedPhases, req.Responsibility); err != nil {
		h.log.Error("failed to update team member", zap.Error(err))
		http.Error(w, "failed to update team member", http.StatusInternalServerError)
		return
	}

	// Get updated member
	updated, _ := h.pg.ProjectTeam().GetMember(r.Context(), tenant, memberID)
	writeJSON(w, http.StatusOK, updated)
}

func (h *ProjectTeamHandler) RemoveMember(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "id")
	memberID := chi.URLParam(r, "memberId")
	tenant := middleware.TenantID(r.Context())

	// Get member info before removing
	member, err := h.pg.ProjectTeam().GetMember(r.Context(), tenant, memberID)
	if err != nil {
		http.Error(w, "member not found", http.StatusNotFound)
		return
	}
	if member.ProjectID != projectID {
		http.Error(w, "member not found in this project", http.StatusNotFound)
		return
	}

	if err := h.pg.ProjectTeam().RemoveMember(r.Context(), tenant, memberID); err != nil {
		h.log.Error("failed to remove team member", zap.Error(err))
		http.Error(w, "failed to remove team member", http.StatusInternalServerError)
		return
	}

	// Create activity for the removal
	h.logTeamActivity(r, projectID, member.UserID, member.UserName, "removed", member.Role)

	w.WriteHeader(http.StatusNoContent)
}

func (h *ProjectTeamHandler) ListMembers(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "id")
	tenant := middleware.TenantID(r.Context())

	members, err := h.pg.ProjectTeam().ListByProject(r.Context(), tenant, projectID)
	if err != nil {
		h.log.Error("failed to list team members", zap.Error(err))
		http.Error(w, "failed to list team members", http.StatusInternalServerError)
		return
	}
	if members == nil {
		members = []models.ProjectTeamMember{}
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"members": members,
		"total":   len(members),
	})
}

func (h *ProjectTeamHandler) ListMyProjects(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())
	userID := middleware.UserID(r.Context())

	members, err := h.pg.ProjectTeam().ListByUser(r.Context(), tenant, userID)
	if err != nil {
		h.log.Error("failed to list user projects", zap.Error(err))
		http.Error(w, "failed to list projects", http.StatusInternalServerError)
		return
	}
	if members == nil {
		members = []models.ProjectTeamMember{}
	}

	// Extract project IDs
	projectIDs := make([]string, len(members))
	for i, m := range members {
		projectIDs[i] = m.ProjectID
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"memberships": members,
		"projectIds":  projectIDs,
		"total":       len(members),
	})
}

// Phase assignments

type addPhaseAssignmentReq struct {
	UserID         string `json:"userId"`
	UserEmail      string `json:"userEmail"`
	UserName       string `json:"userName"`
	AssignmentType string `json:"assignmentType"`
}

func (h *ProjectTeamHandler) AddPhaseAssignment(w http.ResponseWriter, r *http.Request) {
	phaseID := chi.URLParam(r, "phaseId")
	tenant := middleware.TenantID(r.Context())
	currentUserID := middleware.UserID(r.Context())

	var req addPhaseAssignmentReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.UserID) == "" {
		http.Error(w, "userId required", http.StatusBadRequest)
		return
	}

	// Get phase to get project ID
	phase, err := h.pg.Phases().GetByID(r.Context(), tenant, phaseID)
	if err != nil {
		http.Error(w, "phase not found", http.StatusNotFound)
		return
	}

	assignmentType := models.PhaseAssignmentType(req.AssignmentType)
	if assignmentType == "" {
		assignmentType = models.PhaseAssignmentCollaborator
	}

	now := time.Now().UTC()
	assignment := models.PhaseUserAssignment{
		ID:               store.NewID("pha"),
		TenantID:         tenant,
		PhaseID:          phaseID,
		ProjectID:        phase.ProjectID,
		UserID:           strings.TrimSpace(req.UserID),
		UserEmail:        strings.TrimSpace(req.UserEmail),
		UserName:         strings.TrimSpace(req.UserName),
		AssignmentType:   assignmentType,
		AssignedByUserID: currentUserID,
		AssignedAt:       now,
		CreatedAt:        now,
	}

	if err := h.pg.ProjectTeam().AddPhaseAssignment(r.Context(), assignment); err != nil {
		h.log.Error("failed to add phase assignment", zap.Error(err))
		http.Error(w, "failed to add phase assignment", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, assignment)
}

func (h *ProjectTeamHandler) RemovePhaseAssignment(w http.ResponseWriter, r *http.Request) {
	assignmentID := chi.URLParam(r, "assignmentId")
	tenant := middleware.TenantID(r.Context())

	if err := h.pg.ProjectTeam().RemovePhaseAssignment(r.Context(), tenant, assignmentID); err != nil {
		h.log.Error("failed to remove phase assignment", zap.Error(err))
		http.Error(w, "failed to remove phase assignment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ProjectTeamHandler) ListPhaseAssignments(w http.ResponseWriter, r *http.Request) {
	phaseID := chi.URLParam(r, "phaseId")
	tenant := middleware.TenantID(r.Context())

	assignments, err := h.pg.ProjectTeam().ListPhaseAssignments(r.Context(), tenant, phaseID)
	if err != nil {
		h.log.Error("failed to list phase assignments", zap.Error(err))
		http.Error(w, "failed to list phase assignments", http.StatusInternalServerError)
		return
	}
	if assignments == nil {
		assignments = []models.PhaseUserAssignment{}
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"assignments": assignments,
		"total":       len(assignments),
	})
}

// Helper to log team activities
func (h *ProjectTeamHandler) logTeamActivity(r *http.Request, projectID, userID, userName, action string, role models.TeamMemberRole) {
	tenant := middleware.TenantID(r.Context())
	actorID := middleware.UserID(r.Context())
	actorName := middleware.UserName(r.Context())
	actorEmail := "" // Not always available

	now := time.Now().UTC()
	activity := models.ProjectActivity{
		ID:           store.NewID("act"),
		TenantID:     tenant,
		ProjectID:    projectID,
		ActivityType: models.ActivityAssignment,
		ActorUserID:  actorID,
		ActorEmail:   actorEmail,
		ActorName:    actorName,
		Content:      "",
		Metadata: map[string]any{
			"userId":   userID,
			"userName": userName,
			"action":   action,
			"role":     string(role),
		},
		AttachmentIDs: []string{},
		Visibility:    models.VisibilityTeam,
		IsPinned:      false,
		CreatedAt:     now,
	}

	if err := h.pg.ProjectActivities().CreateActivity(r.Context(), activity); err != nil {
		h.log.Warn("failed to log team activity", zap.Error(err))
	}
}
