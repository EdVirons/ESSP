package models

import (
	"time"
)

// DemoLeadStage represents the pipeline stage of a lead.
type DemoLeadStage string

const (
	StageNewLead       DemoLeadStage = "new_lead"
	StageContacted     DemoLeadStage = "contacted"
	StageDemoScheduled DemoLeadStage = "demo_scheduled"
	StageDemoCompleted DemoLeadStage = "demo_completed"
	StageProposalSent  DemoLeadStage = "proposal_sent"
	StageNegotiation   DemoLeadStage = "negotiation"
	StageWon           DemoLeadStage = "won"
	StageLost          DemoLeadStage = "lost"
)

// DemoLeadSource represents how the lead was acquired.
type DemoLeadSource string

const (
	SourceWebsite      DemoLeadSource = "website"
	SourceReferral     DemoLeadSource = "referral"
	SourceEvent        DemoLeadSource = "event"
	SourceColdOutreach DemoLeadSource = "cold_outreach"
	SourceInbound      DemoLeadSource = "inbound"
)

// DemoActivityType represents the type of activity on a lead.
type DemoActivityType string

const (
	DemoActivityNote        DemoActivityType = "note"
	DemoActivityCall        DemoActivityType = "call"
	DemoActivityEmail       DemoActivityType = "email"
	DemoActivityMeeting     DemoActivityType = "meeting"
	DemoActivityDemo        DemoActivityType = "demo"
	DemoActivityStageChange DemoActivityType = "stage_change"
	DemoActivityCreated     DemoActivityType = "created"
	DemoActivityUpdated     DemoActivityType = "updated"
)

// DemoScheduleStatus represents the status of a scheduled demo.
type DemoScheduleStatus string

const (
	ScheduleStatusScheduled   DemoScheduleStatus = "scheduled"
	ScheduleStatusCompleted   DemoScheduleStatus = "completed"
	ScheduleStatusCancelled   DemoScheduleStatus = "cancelled"
	ScheduleStatusRescheduled DemoScheduleStatus = "rescheduled"
)

// DemoLead represents a sales lead in the pipeline.
type DemoLead struct {
	ID       string `json:"id"`
	TenantID string `json:"tenantId"`

	// School/Lead Info
	SchoolID     *string `json:"schoolId"`
	SchoolName   string  `json:"schoolName"`
	ContactName  string  `json:"contactName"`
	ContactEmail string  `json:"contactEmail"`
	ContactPhone string  `json:"contactPhone"`
	ContactRole  string  `json:"contactRole"`

	// Location Info
	CountyCode    string `json:"countyCode"`
	CountyName    string `json:"countyName"`
	SubCountyCode string `json:"subCountyCode"`
	SubCountyName string `json:"subCountyName"`

	// Pipeline Info
	Stage          DemoLeadStage `json:"stage"`
	StageChangedAt time.Time     `json:"stageChangedAt"`

	// Deal Info
	EstimatedValue    *float64   `json:"estimatedValue"`
	EstimatedDevices  *int       `json:"estimatedDevices"`
	Probability       int        `json:"probability"`
	ExpectedCloseDate *time.Time `json:"expectedCloseDate"`

	// Source & Attribution
	LeadSource DemoLeadSource `json:"leadSource"`
	AssignedTo string         `json:"assignedTo"`

	// Notes & Details
	Notes string   `json:"notes"`
	Tags  []string `json:"tags"`

	// Lost reason
	LostReason string `json:"lostReason"`
	LostNotes  string `json:"lostNotes"`

	// Metadata
	CreatedBy string    `json:"createdBy"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// DemoLeadActivity represents an activity or interaction with a lead.
type DemoLeadActivity struct {
	ID       string `json:"id"`
	TenantID string `json:"tenantId"`
	LeadID   string `json:"leadId"`

	ActivityType DemoActivityType `json:"activityType"`
	Description  string           `json:"description"`

	// For stage changes
	FromStage *DemoLeadStage `json:"fromStage"`
	ToStage   *DemoLeadStage `json:"toStage"`

	// For scheduled activities
	ScheduledAt *time.Time `json:"scheduledAt"`
	CompletedAt *time.Time `json:"completedAt"`

	CreatedBy string    `json:"createdBy"`
	CreatedAt time.Time `json:"createdAt"`
}

// DemoAttendee represents an attendee for a scheduled demo.
type DemoAttendee struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

// DemoSchedule represents a scheduled demo meeting.
type DemoSchedule struct {
	ID       string `json:"id"`
	TenantID string `json:"tenantId"`
	LeadID   string `json:"leadId"`

	ScheduledDate   time.Time `json:"scheduledDate"`
	ScheduledTime   string    `json:"scheduledTime"`
	DurationMinutes int       `json:"durationMinutes"`

	Location    string         `json:"location"`
	MeetingLink string         `json:"meetingLink"`
	Attendees   []DemoAttendee `json:"attendees"`

	Status       DemoScheduleStatus `json:"status"`
	Outcome      string             `json:"outcome"`
	OutcomeNotes string             `json:"outcomeNotes"`

	ReminderSent bool `json:"reminderSent"`

	CreatedBy string    `json:"createdBy"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// DemoLeadWithActivities includes the lead with recent activities.
type DemoLeadWithActivities struct {
	DemoLead
	RecentActivities []DemoLeadActivity `json:"recentActivities"`
	NextDemo         *DemoSchedule      `json:"nextDemo"`
}

// CreateDemoLeadRequest is the request payload for creating a lead.
type CreateDemoLeadRequest struct {
	SchoolID         *string  `json:"schoolId"`
	SchoolName       string   `json:"schoolName" validate:"required"`
	ContactName      string   `json:"contactName"`
	ContactEmail     string   `json:"contactEmail"`
	ContactPhone     string   `json:"contactPhone"`
	ContactRole      string   `json:"contactRole"`
	CountyCode       string   `json:"countyCode"`
	CountyName       string   `json:"countyName"`
	SubCountyCode    string   `json:"subCountyCode"`
	SubCountyName    string   `json:"subCountyName"`
	EstimatedValue   *float64 `json:"estimatedValue"`
	EstimatedDevices *int     `json:"estimatedDevices"`
	LeadSource       string   `json:"leadSource"`
	Notes            string   `json:"notes"`
	Tags             []string `json:"tags"`
}

// UpdateDemoLeadRequest is the request payload for updating a lead.
type UpdateDemoLeadRequest struct {
	SchoolName        *string   `json:"schoolName"`
	ContactName       *string   `json:"contactName"`
	ContactEmail      *string   `json:"contactEmail"`
	ContactPhone      *string   `json:"contactPhone"`
	ContactRole       *string   `json:"contactRole"`
	CountyCode        *string   `json:"countyCode"`
	CountyName        *string   `json:"countyName"`
	SubCountyCode     *string   `json:"subCountyCode"`
	SubCountyName     *string   `json:"subCountyName"`
	EstimatedValue    *float64  `json:"estimatedValue"`
	EstimatedDevices  *int      `json:"estimatedDevices"`
	Probability       *int      `json:"probability"`
	ExpectedCloseDate *string   `json:"expectedCloseDate"`
	AssignedTo        *string   `json:"assignedTo"`
	Notes             *string   `json:"notes"`
	Tags              *[]string `json:"tags"`
}

// UpdateLeadStageRequest is the request for changing a lead's stage.
type UpdateLeadStageRequest struct {
	Stage      DemoLeadStage `json:"stage" validate:"required"`
	LostReason string        `json:"lostReason"`
	LostNotes  string        `json:"lostNotes"`
}

// AddLeadNoteRequest is the request for adding a note to a lead.
type AddLeadNoteRequest struct {
	Note string `json:"note" validate:"required"`
}

// CreateDemoScheduleRequest is the request for scheduling a demo.
type CreateDemoScheduleRequest struct {
	ScheduledDate   string         `json:"scheduledDate" validate:"required"`
	ScheduledTime   string         `json:"scheduledTime"`
	DurationMinutes int            `json:"durationMinutes"`
	Location        string         `json:"location"`
	MeetingLink     string         `json:"meetingLink"`
	Attendees       []DemoAttendee `json:"attendees"`
}

// PipelineStageCount represents the count and value of leads in a stage.
type PipelineStageCount struct {
	Stage      DemoLeadStage `json:"stage"`
	Count      int           `json:"count"`
	TotalValue float64       `json:"totalValue"`
}

// PipelineSummary is the summary of the pipeline by stage.
type PipelineSummary struct {
	Stages         []PipelineStageCount `json:"stages"`
	TotalLeads     int                  `json:"totalLeads"`
	TotalValue     float64              `json:"totalValue"`
	AverageValue   float64              `json:"averageValue"`
	ConversionRate float64              `json:"conversionRate"`
}

// DemoLeadFilters represents filters for querying leads.
type DemoLeadFilters struct {
	Stage      *DemoLeadStage  `json:"stage"`
	AssignedTo *string         `json:"assignedTo"`
	LeadSource *DemoLeadSource `json:"leadSource"`
	Search     *string         `json:"search"`
	Limit      int             `json:"limit"`
	Offset     int             `json:"offset"`
}

// Predefined stage configurations.
var (
	DemoLeadStages = []DemoLeadStage{
		StageNewLead,
		StageContacted,
		StageDemoScheduled,
		StageDemoCompleted,
		StageProposalSent,
		StageNegotiation,
		StageWon,
		StageLost,
	}

	DemoLeadSources = []DemoLeadSource{
		SourceWebsite,
		SourceReferral,
		SourceEvent,
		SourceColdOutreach,
		SourceInbound,
	}

	// Stage to probability mapping (default values)
	StageProbability = map[DemoLeadStage]int{
		StageNewLead:       10,
		StageContacted:     20,
		StageDemoScheduled: 40,
		StageDemoCompleted: 60,
		StageProposalSent:  70,
		StageNegotiation:   80,
		StageWon:           100,
		StageLost:          0,
	}
)
