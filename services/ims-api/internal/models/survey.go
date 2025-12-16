package models

import "time"

// SurveyStatus represents the status of a survey.
type SurveyStatus string

const (
	SurveyDraft     SurveyStatus = "draft"
	SurveySubmitted SurveyStatus = "submitted"
	SurveyApproved  SurveyStatus = "approved"
)

// SiteSurvey represents a site survey.
type SiteSurvey struct {
	ID                string       `json:"id"`
	TenantID          string       `json:"tenantId"`
	ProjectID         string       `json:"projectId"`
	Status            SurveyStatus `json:"status"`
	ConductedByUserID string       `json:"conductedByUserId"`
	ConductedAt       *time.Time   `json:"conductedAt"`
	Summary           string       `json:"summary"`
	Risks             string       `json:"risks"`
	CreatedAt         time.Time    `json:"createdAt"`
	UpdatedAt         time.Time    `json:"updatedAt"`
}

// SurveyRoom represents a room surveyed.
type SurveyRoom struct {
	ID           string    `json:"id"`
	TenantID     string    `json:"tenantId"`
	SurveyID     string    `json:"surveyId"`
	Name         string    `json:"name"`
	RoomType     string    `json:"roomType"`
	Floor        string    `json:"floor"`
	PowerNotes   string    `json:"powerNotes"`
	NetworkNotes string    `json:"networkNotes"`
	CreatedAt    time.Time `json:"createdAt"`
}

// SurveyPhoto represents a photo attached to a survey.
type SurveyPhoto struct {
	ID           string    `json:"id"`
	TenantID     string    `json:"tenantId"`
	SurveyID     string    `json:"surveyId"`
	RoomID       string    `json:"roomId"`
	AttachmentID string    `json:"attachmentId"`
	Caption      string    `json:"caption"`
	CreatedAt    time.Time `json:"createdAt"`
}
