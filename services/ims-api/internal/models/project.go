package models

import "time"

// ProjectType represents the type of service project.
type ProjectType string

const (
	ProjectTypeFullInstallation ProjectType = "full_installation"
	ProjectTypeDeviceRefresh    ProjectType = "device_refresh"
	ProjectTypeSupport          ProjectType = "support"
	ProjectTypeRepair           ProjectType = "repair"
	ProjectTypeTraining         ProjectType = "training"
)

// ProjectStatus represents the status of a project.
type ProjectStatus string

const (
	ProjectActive    ProjectStatus = "active"
	ProjectPaused    ProjectStatus = "paused"
	ProjectCompleted ProjectStatus = "completed"
)

// PhaseType represents the type of phase.
type PhaseType string

const (
	// Full Installation phases
	PhaseDemo        PhaseType = "demo"
	PhaseSurvey      PhaseType = "survey"
	PhaseProcurement PhaseType = "procurement"
	PhaseInstall     PhaseType = "install"
	PhaseIntegrate   PhaseType = "integrate"
	PhaseCommission  PhaseType = "commission"
	PhaseOps         PhaseType = "ops"

	// Device Refresh phases
	PhaseAssessment   PhaseType = "assessment"
	PhaseDeployment   PhaseType = "deployment"
	PhaseVerification PhaseType = "verification"

	// Support phases
	PhaseOnboarding PhaseType = "onboarding"
	PhaseActive     PhaseType = "active"
	PhaseRenewal    PhaseType = "renewal"

	// Repair phases
	PhaseIntake    PhaseType = "intake"
	PhaseDiagnosis PhaseType = "diagnosis"
	PhaseRepair    PhaseType = "repair"
	PhaseTesting   PhaseType = "testing"
	PhaseHandover  PhaseType = "handover"

	// Training phases
	PhasePlanning      PhaseType = "planning"
	PhaseDelivery      PhaseType = "delivery"
	PhaseCertification PhaseType = "certification"
)

// ProjectTypeConfig defines the configuration for a project type.
type ProjectTypeConfig struct {
	Type         ProjectType `json:"type"`
	Label        string      `json:"label"`
	Description  string      `json:"description"`
	Phases       []PhaseType `json:"phases"`
	DefaultPhase PhaseType   `json:"defaultPhase"`
}

// ProjectTypeConfigs contains the configuration for all project types.
var ProjectTypeConfigs = map[ProjectType]ProjectTypeConfig{
	ProjectTypeFullInstallation: {
		Type:         ProjectTypeFullInstallation,
		Label:        "Full Installation",
		Description:  "Complete school technology installation",
		Phases:       []PhaseType{PhaseDemo, PhaseSurvey, PhaseProcurement, PhaseInstall, PhaseIntegrate, PhaseCommission, PhaseOps},
		DefaultPhase: PhaseDemo,
	},
	ProjectTypeDeviceRefresh: {
		Type:         ProjectTypeDeviceRefresh,
		Label:        "Device Refresh",
		Description:  "Upgrade or replace existing devices",
		Phases:       []PhaseType{PhaseAssessment, PhaseProcurement, PhaseDeployment, PhaseVerification},
		DefaultPhase: PhaseAssessment,
	},
	ProjectTypeSupport: {
		Type:         ProjectTypeSupport,
		Label:        "Support",
		Description:  "Ongoing technical support contract",
		Phases:       []PhaseType{PhaseOnboarding, PhaseActive, PhaseRenewal},
		DefaultPhase: PhaseOnboarding,
	},
	ProjectTypeRepair: {
		Type:         ProjectTypeRepair,
		Label:        "Repair",
		Description:  "Device repair workflow",
		Phases:       []PhaseType{PhaseIntake, PhaseDiagnosis, PhaseRepair, PhaseTesting, PhaseHandover},
		DefaultPhase: PhaseIntake,
	},
	ProjectTypeTraining: {
		Type:         ProjectTypeTraining,
		Label:        "Training",
		Description:  "Staff training project",
		Phases:       []PhaseType{PhasePlanning, PhaseDelivery, PhaseAssessment, PhaseCertification},
		DefaultPhase: PhasePlanning,
	},
}

// ValidProjectTypes returns all valid project types.
func ValidProjectTypes() []ProjectType {
	return []ProjectType{
		ProjectTypeFullInstallation,
		ProjectTypeDeviceRefresh,
		ProjectTypeSupport,
		ProjectTypeRepair,
		ProjectTypeTraining,
	}
}

// IsValidPhaseForType checks if a phase is valid for a given project type.
func IsValidPhaseForType(projectType ProjectType, phase PhaseType) bool {
	config, ok := ProjectTypeConfigs[projectType]
	if !ok {
		return false
	}
	for _, p := range config.Phases {
		if p == phase {
			return true
		}
	}
	return false
}

// PhaseStatus represents the status of a phase.
type PhaseStatus string

const (
	PhasePending    PhaseStatus = "pending"
	PhaseInProgress PhaseStatus = "in_progress"
	PhaseBlocked    PhaseStatus = "blocked"
	PhaseDone       PhaseStatus = "done"
)

// SchoolServiceProject represents a school service project.
type SchoolServiceProject struct {
	ID                   string        `json:"id"`
	TenantID             string        `json:"tenantId"`
	SchoolID             string        `json:"schoolId"`
	ProjectType          ProjectType   `json:"projectType"`
	Status               ProjectStatus `json:"status"`
	CurrentPhase         PhaseType     `json:"currentPhase"`
	StartDate            string        `json:"startDate"`  // YYYY-MM-DD (kept as string for simplicity)
	GoLiveDate           string        `json:"goLiveDate"` // YYYY-MM-DD
	AccountManagerUserID string        `json:"accountManagerUserId"`
	Notes                string        `json:"notes"`
	CreatedAt            time.Time     `json:"createdAt"`
	UpdatedAt            time.Time     `json:"updatedAt"`
}

// ServicePhase represents a phase in a service project.
type ServicePhase struct {
	ID                      string      `json:"id"`
	TenantID                string      `json:"tenantId"`
	ProjectID               string      `json:"projectId"`
	PhaseType               PhaseType   `json:"phaseType"`
	Status                  PhaseStatus `json:"status"`
	OwnerRole               string      `json:"ownerRole"`
	OwnerUserID             string      `json:"ownerUserId"`
	OwnerUserName           string      `json:"ownerUserName"`
	StartDate               string      `json:"startDate"`
	EndDate                 string      `json:"endDate"`
	Notes                   string      `json:"notes"`
	StatusChangedAt         *time.Time  `json:"statusChangedAt,omitempty"`
	StatusChangedByUserID   string      `json:"statusChangedByUserId"`
	StatusChangedByUserName string      `json:"statusChangedByUserName"`
	CreatedAt               time.Time   `json:"createdAt"`
	UpdatedAt               time.Time   `json:"updatedAt"`
}

// PhaseChecklistTemplate represents a checklist template for a phase.
type PhaseChecklistTemplate struct {
	ID          string      `json:"id"`
	TenantID    string      `json:"tenantId"`
	ProjectType ProjectType `json:"projectType"`
	PhaseType   PhaseType   `json:"phaseType"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Required    bool        `json:"required"`
	CreatedAt   time.Time   `json:"createdAt"`
	UpdatedAt   time.Time   `json:"updatedAt"`
}
