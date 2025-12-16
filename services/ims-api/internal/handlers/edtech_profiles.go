package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/edvirons/ssp/ims/internal/claude"
	"github.com/edvirons/ssp/ims/internal/config"
	"github.com/edvirons/ssp/ims/internal/middleware"
	"github.com/edvirons/ssp/ims/internal/models"
	"github.com/edvirons/ssp/ims/internal/store"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// EdTechProfilesHandler handles EdTech profile endpoints
type EdTechProfilesHandler struct {
	log    *zap.Logger
	pg     *store.Postgres
	claude *claude.Client
	cfg    config.Config
}

// NewEdTechProfilesHandler creates a new handler
func NewEdTechProfilesHandler(log *zap.Logger, pg *store.Postgres, cfg config.Config) *EdTechProfilesHandler {
	var claudeClient *claude.Client
	if cfg.ClaudeAPIKey != "" {
		claudeClient = claude.NewClient(claude.ClientConfig{
			APIKey:         cfg.ClaudeAPIKey,
			Model:          cfg.ClaudeModel,
			MaxTokens:      cfg.ClaudeMaxTokens,
			TimeoutSeconds: cfg.ClaudeTimeoutSeconds,
		}, log)
	}

	return &EdTechProfilesHandler{
		log:    log,
		pg:     pg,
		claude: claudeClient,
		cfg:    cfg,
	}
}

// GetOptions returns form options for the assessment
func (h *EdTechProfilesHandler) GetOptions(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"networkQuality":     models.NetworkQualityOptions,
		"internetSpeed":      models.InternetSpeedOptions,
		"deviceAge":          models.DeviceAgeOptions,
		"supportFrequency":   models.SupportFrequencyOptions,
		"resolutionTime":     models.ResolutionTimeOptions,
		"budgetRange":        models.BudgetRangeOptions,
		"timeline":           models.TimelineOptions,
		"lmsPlatforms":       models.LMSPlatformOptions,
		"existingSoftware":   models.ExistingSoftwareOptions,
		"painPoints":         models.PainPointOptions,
		"strategicGoals":     models.StrategicGoalOptions,
	})
}

// GetBySchoolID retrieves a profile for the current user's school
func (h *EdTechProfilesHandler) GetBySchoolID(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.TenantID(r.Context())
	schoolID := chi.URLParam(r, "schoolId")

	if schoolID == "" {
		http.Error(w, "school_id is required", http.StatusBadRequest)
		return
	}

	profile, err := h.pg.EdTechProfiles().GetBySchoolID(r.Context(), tenantID, schoolID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			writeJSON(w, http.StatusOK, map[string]any{"profile": nil})
			return
		}
		h.log.Error("failed to get edtech profile", zap.Error(err))
		http.Error(w, "failed to get profile", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"profile": profile})
}

// SaveProfile creates or updates an EdTech profile
func (h *EdTechProfilesHandler) SaveProfile(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.TenantID(r.Context())
	userID := middleware.UserID(r.Context())

	var req models.EdTechProfile
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if req.SchoolID == "" {
		http.Error(w, "school_id is required", http.StatusBadRequest)
		return
	}

	now := time.Now().UTC()

	// Check if profile exists
	existing, err := h.pg.EdTechProfiles().GetBySchoolID(r.Context(), tenantID, req.SchoolID)
	isNew := err != nil

	var profile models.EdTechProfile
	if isNew {
		profile = req
		profile.ID = store.NewID("edtech")
		profile.TenantID = tenantID
		profile.Status = models.EdTechProfileDraft
		profile.Version = 1
		profile.CreatedAt = now
		profile.UpdatedAt = now
	} else {
		// Update existing profile
		profile = existing
		profile.TotalDevices = req.TotalDevices
		profile.DeviceTypes = req.DeviceTypes
		profile.NetworkQuality = req.NetworkQuality
		profile.InternetSpeed = req.InternetSpeed
		profile.LMSPlatform = req.LMSPlatform
		profile.ExistingSoftware = req.ExistingSoftware
		profile.ITStaffCount = req.ITStaffCount
		profile.DeviceAge = req.DeviceAge
		profile.PainPoints = req.PainPoints
		profile.SupportSatisfaction = req.SupportSatisfaction
		profile.BiggestChallenges = req.BiggestChallenges
		profile.SupportFrequency = req.SupportFrequency
		profile.AvgResolutionTime = req.AvgResolutionTime
		profile.BiggestFrustration = req.BiggestFrustration
		profile.WishList = req.WishList
		profile.StrategicGoals = req.StrategicGoals
		profile.BudgetRange = req.BudgetRange
		profile.Timeline = req.Timeline
		profile.ExpansionPlans = req.ExpansionPlans
		profile.PriorityRanking = req.PriorityRanking
		profile.DecisionMakers = req.DecisionMakers
		profile.UpdatedAt = now

		// Create history entry before update
		snapshot, _ := json.Marshal(existing)
		historyEntry := models.EdTechProfileHistory{
			ID:           store.NewID("edtech-hist"),
			ProfileID:    existing.ID,
			Snapshot:     snapshot,
			ChangedBy:    userID,
			ChangeReason: "Profile updated",
			ChangedAt:    now,
		}
		if err := h.pg.EdTechProfileHistory().Create(r.Context(), historyEntry); err != nil {
			h.log.Warn("failed to create history entry", zap.Error(err))
		}
	}

	if err := h.pg.EdTechProfiles().Upsert(r.Context(), profile); err != nil {
		h.log.Error("failed to save edtech profile", zap.Error(err))
		http.Error(w, "failed to save profile", http.StatusInternalServerError)
		return
	}

	// Refetch to get the updated version
	profile, _ = h.pg.EdTechProfiles().GetBySchoolID(r.Context(), tenantID, req.SchoolID)

	writeJSON(w, http.StatusOK, map[string]any{"profile": profile})
}

// GenerateAI generates AI summary and follow-up questions
func (h *EdTechProfilesHandler) GenerateAI(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.TenantID(r.Context())
	profileID := chi.URLParam(r, "id")

	if h.claude == nil || !h.claude.IsEnabled() {
		http.Error(w, "AI is not available", http.StatusServiceUnavailable)
		return
	}

	profile, err := h.pg.EdTechProfiles().GetByID(r.Context(), tenantID, profileID)
	if err != nil {
		http.Error(w, "profile not found", http.StatusNotFound)
		return
	}

	// Build AI prompt
	systemPrompt := buildEdTechAnalysisPrompt()
	userMessage := buildEdTechProfileMessage(profile)

	messages := []claude.Message{
		{Role: "user", Content: userMessage},
	}

	resp, err := h.claude.Chat(r.Context(), systemPrompt, messages)
	if err != nil {
		h.log.Error("failed to get AI response", zap.Error(err))
		http.Error(w, "failed to generate AI analysis", http.StatusInternalServerError)
		return
	}

	// Parse AI response
	analysis := parseAIAnalysisResponse(resp.Content)

	// Update profile with AI data
	if err := h.pg.EdTechProfiles().UpdateAI(
		r.Context(),
		tenantID,
		profileID,
		analysis.Summary,
		analysis.Recommendations,
		analysis.FollowUpQuestions,
	); err != nil {
		h.log.Error("failed to update AI data", zap.Error(err))
		http.Error(w, "failed to save AI analysis", http.StatusInternalServerError)
		return
	}

	// Refetch profile
	profile, _ = h.pg.EdTechProfiles().GetByID(r.Context(), tenantID, profileID)

	writeJSON(w, http.StatusOK, map[string]any{"profile": profile})
}

// SubmitFollowUp handles follow-up question responses
func (h *EdTechProfilesHandler) SubmitFollowUp(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.TenantID(r.Context())
	profileID := chi.URLParam(r, "id")

	var req struct {
		Responses map[string]string `json:"responses"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if err := h.pg.EdTechProfiles().UpdateFollowUpResponses(r.Context(), tenantID, profileID, req.Responses); err != nil {
		h.log.Error("failed to update follow-up responses", zap.Error(err))
		http.Error(w, "failed to save responses", http.StatusInternalServerError)
		return
	}

	// Refetch and return
	profile, _ := h.pg.EdTechProfiles().GetByID(r.Context(), tenantID, profileID)
	writeJSON(w, http.StatusOK, map[string]any{"profile": profile})
}

// Complete marks a profile as completed
func (h *EdTechProfilesHandler) Complete(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.TenantID(r.Context())
	userID := middleware.UserID(r.Context())
	profileID := chi.URLParam(r, "id")

	if err := h.pg.EdTechProfiles().Complete(r.Context(), tenantID, profileID, userID); err != nil {
		h.log.Error("failed to complete profile", zap.Error(err))
		http.Error(w, "failed to complete profile", http.StatusInternalServerError)
		return
	}

	profile, _ := h.pg.EdTechProfiles().GetByID(r.Context(), tenantID, profileID)
	writeJSON(w, http.StatusOK, map[string]any{"profile": profile})
}

// GetHistory returns version history for a profile
func (h *EdTechProfilesHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.TenantID(r.Context())
	schoolID := chi.URLParam(r, "schoolId")

	profile, err := h.pg.EdTechProfiles().GetBySchoolID(r.Context(), tenantID, schoolID)
	if err != nil {
		http.Error(w, "profile not found", http.StatusNotFound)
		return
	}

	limit := parseLimit(r.URL.Query().Get("limit"), 20, 100)
	history, err := h.pg.EdTechProfileHistory().List(r.Context(), profile.ID, limit)
	if err != nil {
		h.log.Error("failed to get history", zap.Error(err))
		http.Error(w, "failed to get history", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"history": history})
}

// AI prompt helpers

func buildEdTechAnalysisPrompt() string {
	return `You are an EdTech consultant analyzing a school's technology profile. Based on the assessment data provided, generate:

1. A 2-3 sentence summary of their current EdTech landscape
2. 3-4 actionable recommendations with priority levels (high/medium/low)
3. 2-3 follow-up questions to gather additional important information

Respond in valid JSON format:
{
  "summary": "2-3 sentence summary",
  "recommendations": [
    {"category": "infrastructure|training|software|security|support", "title": "short title", "description": "detailed recommendation", "priority": "high|medium|low"}
  ],
  "followUpQuestions": [
    {"id": "q1", "question": "the question", "context": "why this matters"}
  ]
}

Focus on practical, actionable advice appropriate for schools in developing markets. Consider resource constraints and prioritize high-impact, cost-effective solutions.`
}

func buildEdTechProfileMessage(p models.EdTechProfile) string {
	profileJSON, _ := json.MarshalIndent(map[string]any{
		"infrastructure": map[string]any{
			"totalDevices":     p.TotalDevices,
			"deviceTypes":      p.DeviceTypes,
			"networkQuality":   p.NetworkQuality,
			"internetSpeed":    p.InternetSpeed,
			"lmsPlatform":      p.LMSPlatform,
			"existingSoftware": p.ExistingSoftware,
			"itStaffCount":     p.ITStaffCount,
			"deviceAge":        p.DeviceAge,
		},
		"painPoints": map[string]any{
			"painPoints":          p.PainPoints,
			"supportSatisfaction": p.SupportSatisfaction,
			"biggestChallenges":   p.BiggestChallenges,
			"supportFrequency":    p.SupportFrequency,
			"avgResolutionTime":   p.AvgResolutionTime,
			"biggestFrustration":  p.BiggestFrustration,
			"wishList":            p.WishList,
		},
		"goals": map[string]any{
			"strategicGoals":  p.StrategicGoals,
			"budgetRange":     p.BudgetRange,
			"timeline":        p.Timeline,
			"expansionPlans":  p.ExpansionPlans,
			"priorityRanking": p.PriorityRanking,
			"decisionMakers":  p.DecisionMakers,
		},
	}, "", "  ")

	return "Please analyze this school's EdTech profile:\n\n" + string(profileJSON)
}

type aiAnalysisResult struct {
	Summary           string                     `json:"summary"`
	Recommendations   []models.AIRecommendation  `json:"recommendations"`
	FollowUpQuestions []models.FollowUpQuestion  `json:"followUpQuestions"`
}

func parseAIAnalysisResponse(content string) aiAnalysisResult {
	var result aiAnalysisResult

	// Try to parse as JSON
	if err := json.Unmarshal([]byte(content), &result); err == nil {
		return result
	}

	// If parsing fails, extract JSON from markdown code blocks
	start := strings.Index(content, "{")
	end := strings.LastIndex(content, "}")
	if start >= 0 && end > start {
		jsonStr := content[start : end+1]
		if err := json.Unmarshal([]byte(jsonStr), &result); err == nil {
			return result
		}
	}

	// Fallback: return a basic summary
	result.Summary = "Unable to parse AI response. Please try again."
	result.Recommendations = []models.AIRecommendation{}
	result.FollowUpQuestions = []models.FollowUpQuestion{}
	return result
}
