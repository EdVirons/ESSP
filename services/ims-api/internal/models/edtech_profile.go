package models

import (
	"encoding/json"
	"time"
)

// EdTechProfileStatus represents the status of an EdTech profile.
type EdTechProfileStatus string

const (
	EdTechProfileDraft     EdTechProfileStatus = "draft"
	EdTechProfileCompleted EdTechProfileStatus = "completed"
)

// DeviceTypes represents the breakdown of device types.
type DeviceTypes struct {
	Laptops     int `json:"laptops"`
	Chromebooks int `json:"chromebooks"`
	Tablets     int `json:"tablets"`
	Desktops    int `json:"desktops"`
	Other       int `json:"other"`
}

// AIRecommendation represents a single AI recommendation.
type AIRecommendation struct {
	Category    string `json:"category"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
}

// FollowUpQuestion represents an AI-generated follow-up question.
type FollowUpQuestion struct {
	ID       string `json:"id"`
	Question string `json:"question"`
	Context  string `json:"context"`
}

// EdTechProfile represents a school's EdTech profile assessment.
type EdTechProfile struct {
	ID       string `json:"id"`
	TenantID string `json:"tenantId"`
	SchoolID string `json:"schoolId"`

	// Infrastructure Section
	TotalDevices     int              `json:"totalDevices"`
	DeviceTypes      DeviceTypes      `json:"deviceTypes"`
	NetworkQuality   string           `json:"networkQuality"`
	InternetSpeed    string           `json:"internetSpeed"`
	LMSPlatform      string           `json:"lmsPlatform"`
	ExistingSoftware []string         `json:"existingSoftware"`
	ITStaffCount     int              `json:"itStaffCount"`
	DeviceAge        string           `json:"deviceAge"`

	// Pain Points Section
	PainPoints          []string `json:"painPoints"`
	SupportSatisfaction int      `json:"supportSatisfaction"`
	BiggestChallenges   []string `json:"biggestChallenges"`
	SupportFrequency    string   `json:"supportFrequency"`
	AvgResolutionTime   string   `json:"avgResolutionTime"`
	BiggestFrustration  string   `json:"biggestFrustration"`
	WishList            string   `json:"wishList"`

	// Goals Section
	StrategicGoals  []string `json:"strategicGoals"`
	BudgetRange     string   `json:"budgetRange"`
	Timeline        string   `json:"timeline"`
	ExpansionPlans  string   `json:"expansionPlans"`
	PriorityRanking []string `json:"priorityRanking"`
	DecisionMakers  []string `json:"decisionMakers"`

	// AI Section
	AISummary          string              `json:"aiSummary"`
	AIRecommendations  []AIRecommendation  `json:"aiRecommendations"`
	FollowUpQuestions  []FollowUpQuestion  `json:"followUpQuestions"`
	FollowUpResponses  map[string]string   `json:"followUpResponses"`

	// Metadata
	Status      EdTechProfileStatus `json:"status"`
	CompletedAt *time.Time          `json:"completedAt"`
	CompletedBy string              `json:"completedBy"`
	Version     int                 `json:"version"`
	CreatedAt   time.Time           `json:"createdAt"`
	UpdatedAt   time.Time           `json:"updatedAt"`
}

// EdTechProfileHistory represents a historical snapshot of a profile.
type EdTechProfileHistory struct {
	ID           string          `json:"id"`
	ProfileID    string          `json:"profileId"`
	Snapshot     json.RawMessage `json:"snapshot"`
	ChangedBy    string          `json:"changedBy"`
	ChangeReason string          `json:"changeReason"`
	ChangedAt    time.Time       `json:"changedAt"`
}

// Predefined options for the assessment form.
var (
	NetworkQualityOptions = []string{"excellent", "good", "fair", "poor"}
	InternetSpeedOptions  = []string{"fiber", "broadband", "dsl", "limited", "unreliable"}
	DeviceAgeOptions      = []string{"under_2_years", "2_4_years", "4_6_years", "over_6_years", "mixed"}
	SupportFrequencyOptions = []string{"daily", "weekly", "monthly", "rarely"}
	ResolutionTimeOptions = []string{"same_day", "1_3_days", "week", "longer"}
	BudgetRangeOptions    = []string{"none", "limited", "moderate", "substantial"}
	TimelineOptions       = []string{"immediate", "this_year", "next_year", "long_term"}

	LMSPlatformOptions = []string{
		"Google Classroom",
		"Microsoft Teams",
		"Canvas",
		"Moodle",
		"Blackboard",
		"Schoology",
		"None",
		"Other",
	}

	ExistingSoftwareOptions = []string{
		"Google Workspace",
		"Microsoft 365",
		"Zoom",
		"Kahoot",
		"Quizlet",
		"Khan Academy",
		"Duolingo",
		"Code.org",
		"Scratch",
		"Adobe Creative Suite",
		"AutoCAD",
	}

	PainPointOptions = []string{
		"Device maintenance and repairs",
		"Slow or unreliable internet",
		"Software licensing costs",
		"Teacher training on technology",
		"Student device access equity",
		"Network security concerns",
		"Outdated hardware",
		"Technical support availability",
		"Integration between systems",
		"Data backup and recovery",
	}

	StrategicGoalOptions = []string{
		"1:1 device program",
		"STEM curriculum expansion",
		"Digital literacy training",
		"Network infrastructure upgrade",
		"Learning management system adoption",
		"Device refresh program",
		"Cybersecurity enhancement",
		"Remote learning capability",
		"AI/adaptive learning tools",
		"Coding and robotics programs",
	}
)
