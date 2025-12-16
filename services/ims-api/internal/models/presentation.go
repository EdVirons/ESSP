package models

import (
	"time"
)

// PresentationType represents the type of presentation material.
type PresentationType string

const (
	TypePresentation  PresentationType = "presentation"
	TypeBrochure      PresentationType = "brochure"
	TypeCaseStudy     PresentationType = "case_study"
	TypeVideo         PresentationType = "video"
	TypeROICalculator PresentationType = "roi_calculator"
	TypeTemplate      PresentationType = "template"
	TypeOther         PresentationType = "other"
)

// PresentationCategory represents the category of content.
type PresentationCategory string

const (
	CategoryGeneral         PresentationCategory = "general"
	CategoryProductOverview PresentationCategory = "product_overview"
	CategoryTechnical       PresentationCategory = "technical"
	CategoryPricing         PresentationCategory = "pricing"
	CategoryOnboarding      PresentationCategory = "onboarding"
	CategoryTraining        PresentationCategory = "training"
)

// Presentation represents a sales presentation or marketing material.
type Presentation struct {
	ID       string `json:"id"`
	TenantID string `json:"tenantId"`

	// Content Info
	Title       string               `json:"title"`
	Description string               `json:"description"`
	Type        PresentationType     `json:"type"`
	Category    PresentationCategory `json:"category"`

	// File Storage
	FileKey  string `json:"fileKey"`
	FileName string `json:"fileName"`
	FileSize int64  `json:"fileSize"`
	FileType string `json:"fileType"`

	// Preview
	ThumbnailKey string `json:"thumbnailKey"`
	PreviewType  string `json:"previewType"`

	// Metadata
	Tags       []string `json:"tags"`
	Version    int      `json:"version"`
	IsActive   bool     `json:"isActive"`
	IsFeatured bool     `json:"isFeatured"`

	// Usage Stats
	ViewCount    int        `json:"viewCount"`
	DownloadCount int       `json:"downloadCount"`
	LastViewedAt *time.Time `json:"lastViewedAt"`

	// Audit
	CreatedBy string    `json:"createdBy"`
	UpdatedBy string    `json:"updatedBy"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	// Computed/joined fields (not stored)
	DownloadURL string `json:"downloadUrl,omitempty"`
	PreviewURL  string `json:"previewUrl,omitempty"`
}

// PresentationVersion represents a historical version of a presentation.
type PresentationVersion struct {
	ID             string `json:"id"`
	PresentationID string `json:"presentationId"`

	Version     int    `json:"version"`
	FileKey     string `json:"fileKey"`
	FileName    string `json:"fileName"`
	FileSize    int64  `json:"fileSize"`
	ChangeNotes string `json:"changeNotes"`

	CreatedBy string    `json:"createdBy"`
	CreatedAt time.Time `json:"createdAt"`
}

// PresentationView represents a view event for analytics.
type PresentationView struct {
	ID             string `json:"id"`
	TenantID       string `json:"tenantId"`
	PresentationID string `json:"presentationId"`

	ViewedBy        string    `json:"viewedBy"`
	ViewedAt        time.Time `json:"viewedAt"`
	Context         string    `json:"context"`
	DurationSeconds *int      `json:"durationSeconds"`
}

// SalesMetricsDaily represents daily aggregated sales metrics.
type SalesMetricsDaily struct {
	ID         string    `json:"id"`
	TenantID   string    `json:"tenantId"`
	MetricDate time.Time `json:"metricDate"`

	// Pipeline metrics
	NewLeads       int `json:"newLeads"`
	LeadsContacted int `json:"leadsContacted"`
	DemosScheduled int `json:"demosScheduled"`
	DemosCompleted int `json:"demosCompleted"`
	ProposalsSent  int `json:"proposalsSent"`
	DealsWon       int `json:"dealsWon"`
	DealsLost      int `json:"dealsLost"`

	// Value metrics
	PipelineValue float64 `json:"pipelineValue"`
	WonValue      float64 `json:"wonValue"`
	LostValue     float64 `json:"lostValue"`

	// Activity metrics
	CallsMade    int `json:"callsMade"`
	EmailsSent   int `json:"emailsSent"`
	MeetingsHeld int `json:"meetingsHeld"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// CreatePresentationRequest is the request payload for uploading a presentation.
type CreatePresentationRequest struct {
	Title       string   `json:"title" validate:"required"`
	Description string   `json:"description"`
	Type        string   `json:"type" validate:"required"`
	Category    string   `json:"category"`
	Tags        []string `json:"tags"`
	IsFeatured  bool     `json:"isFeatured"`
}

// UpdatePresentationRequest is the request payload for updating a presentation.
type UpdatePresentationRequest struct {
	Title       *string   `json:"title"`
	Description *string   `json:"description"`
	Type        *string   `json:"type"`
	Category    *string   `json:"category"`
	Tags        *[]string `json:"tags"`
	IsActive    *bool     `json:"isActive"`
	IsFeatured  *bool     `json:"isFeatured"`
}

// PresentationFilters represents filters for querying presentations.
type PresentationFilters struct {
	Type       *PresentationType     `json:"type"`
	Category   *PresentationCategory `json:"category"`
	Search     *string               `json:"search"`
	IsFeatured *bool                 `json:"isFeatured"`
	IsActive   *bool                 `json:"isActive"`
	Limit      int                   `json:"limit"`
	Offset     int                   `json:"offset"`
}

// SalesMetricsSummary is the dashboard summary response.
type SalesMetricsSummary struct {
	// Current period metrics
	TotalLeads        int     `json:"totalLeads"`
	NewLeadsThisPeriod int    `json:"newLeadsThisPeriod"`
	DemosScheduled    int     `json:"demosScheduled"`
	DemosCompleted    int     `json:"demosCompleted"`
	ProposalsSent     int     `json:"proposalsSent"`
	DealsWon          int     `json:"dealsWon"`
	DealsLost         int     `json:"dealsLost"`

	// Value metrics
	TotalPipelineValue float64 `json:"totalPipelineValue"`
	WonValueThisPeriod float64 `json:"wonValueThisPeriod"`

	// Rates
	ConversionRate  float64 `json:"conversionRate"`
	WinRate         float64 `json:"winRate"`
	AverageDealSize float64 `json:"averageDealSize"`

	// Activity metrics
	TotalActivities int `json:"totalActivities"`
}

// RecentActivity represents a recent activity for the dashboard feed.
type RecentActivity struct {
	ID           string    `json:"id"`
	Type         string    `json:"type"`
	Description  string    `json:"description"`
	LeadID       string    `json:"leadId"`
	LeadName     string    `json:"leadName"`
	UserID       string    `json:"userId"`
	UserName     string    `json:"userName"`
	CreatedAt    time.Time `json:"createdAt"`
}

// SchoolsByRegion represents the count of schools/leads by region.
type SchoolsByRegion struct {
	Region string `json:"region"`
	Count  int    `json:"count"`
	Value  float64 `json:"value"`
}

// Predefined options.
var (
	PresentationTypes = []PresentationType{
		TypePresentation,
		TypeBrochure,
		TypeCaseStudy,
		TypeVideo,
		TypeROICalculator,
		TypeTemplate,
		TypeOther,
	}

	PresentationCategories = []PresentationCategory{
		CategoryGeneral,
		CategoryProductOverview,
		CategoryTechnical,
		CategoryPricing,
		CategoryOnboarding,
		CategoryTraining,
	}

	// Allowed MIME types for uploads
	AllowedPresentationMIMETypes = []string{
		"application/pdf",
		"application/vnd.ms-powerpoint",
		"application/vnd.openxmlformats-officedocument.presentationml.presentation",
		"video/mp4",
		"video/webm",
		"application/vnd.ms-excel",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		"image/png",
		"image/jpeg",
	}
)
