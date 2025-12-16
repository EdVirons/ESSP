package models

import "time"

// TeamMemberRole represents the role of a team member in a project.
type TeamMemberRole string

const (
	TeamRoleOwner        TeamMemberRole = "owner"
	TeamRoleCollaborator TeamMemberRole = "collaborator"
	TeamRoleViewer       TeamMemberRole = "viewer"
)

// ValidTeamMemberRoles returns all valid team member roles.
func ValidTeamMemberRoles() []TeamMemberRole {
	return []TeamMemberRole{TeamRoleOwner, TeamRoleCollaborator, TeamRoleViewer}
}

// IsValidTeamMemberRole checks if a role is valid.
func IsValidTeamMemberRole(role string) bool {
	switch TeamMemberRole(role) {
	case TeamRoleOwner, TeamRoleCollaborator, TeamRoleViewer:
		return true
	default:
		return false
	}
}

// ProjectTeamMember represents a team member assigned to a project.
type ProjectTeamMember struct {
	ID               string         `json:"id"`
	TenantID         string         `json:"tenantId"`
	ProjectID        string         `json:"projectId"`
	UserID           string         `json:"userId"`
	UserEmail        string         `json:"userEmail"`
	UserName         string         `json:"userName"`
	Role             TeamMemberRole `json:"role"`
	AssignedPhases   []PhaseType    `json:"assignedPhases"`
	Responsibility   string         `json:"responsibility"`
	AssignedByUserID string         `json:"assignedByUserId"`
	AssignedAt       time.Time      `json:"assignedAt"`
	RemovedAt        *time.Time     `json:"removedAt,omitempty"`
	CreatedAt        time.Time      `json:"createdAt"`
	UpdatedAt        time.Time      `json:"updatedAt"`
}

// PhaseAssignmentType represents the type of assignment to a phase.
type PhaseAssignmentType string

const (
	PhaseAssignmentOwner        PhaseAssignmentType = "owner"
	PhaseAssignmentCollaborator PhaseAssignmentType = "collaborator"
	PhaseAssignmentReviewer     PhaseAssignmentType = "reviewer"
)

// PhaseUserAssignment represents a user assigned to a specific phase.
type PhaseUserAssignment struct {
	ID               string              `json:"id"`
	TenantID         string              `json:"tenantId"`
	PhaseID          string              `json:"phaseId"`
	ProjectID        string              `json:"projectId"`
	UserID           string              `json:"userId"`
	UserEmail        string              `json:"userEmail"`
	UserName         string              `json:"userName"`
	AssignmentType   PhaseAssignmentType `json:"assignmentType"`
	AssignedByUserID string              `json:"assignedByUserId"`
	AssignedAt       time.Time           `json:"assignedAt"`
	CompletedAt      *time.Time          `json:"completedAt,omitempty"`
	RemovedAt        *time.Time          `json:"removedAt,omitempty"`
	CreatedAt        time.Time           `json:"createdAt"`
}

// DefaultPhaseOwnerRole returns the default owner role for a phase type.
func DefaultPhaseOwnerRole(phaseType PhaseType) string {
	switch phaseType {
	case PhaseDemo, PhaseSurvey:
		return "ssp_demo_team"
	case PhaseProcurement:
		return "ssp_admin"
	case PhaseInstall, PhaseIntegrate, PhaseCommission:
		return "ssp_lead_tech"
	case PhaseOps:
		return "ssp_support_agent"
	case PhaseAssessment, PhaseDeployment, PhaseVerification:
		return "ssp_lead_tech"
	case PhaseOnboarding, PhaseActive, PhaseRenewal:
		return "ssp_support_agent"
	case PhaseIntake, PhaseDiagnosis, PhaseRepair, PhaseTesting, PhaseHandover:
		return "ssp_lead_tech"
	case PhasePlanning, PhaseDelivery, PhaseCertification:
		return "ssp_demo_team"
	default:
		return "ssp_admin"
	}
}
