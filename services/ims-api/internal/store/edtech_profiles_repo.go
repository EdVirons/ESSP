package store

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/edvirons/ssp/ims/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EdTechProfilesRepo struct{ pool *pgxpool.Pool }

// Upsert creates or updates an EdTech profile.
func (r *EdTechProfilesRepo) Upsert(ctx context.Context, p models.EdTechProfile) error {
	deviceTypesJSON, _ := json.Marshal(p.DeviceTypes)
	existingSoftwareJSON, _ := json.Marshal(p.ExistingSoftware)
	painPointsJSON, _ := json.Marshal(p.PainPoints)
	aiRecsJSON, _ := json.Marshal(p.AIRecommendations)
	followUpQsJSON, _ := json.Marshal(p.FollowUpQuestions)
	followUpRsJSON, _ := json.Marshal(p.FollowUpResponses)
	priorityRankingJSON, _ := json.Marshal(p.PriorityRanking)

	_, err := r.pool.Exec(ctx, `
		INSERT INTO edtech_profiles (
			id, tenant_id, school_id,
			total_devices, device_types, network_quality, internet_speed, lms_platform,
			existing_software, it_staff_count, device_age,
			pain_points, support_satisfaction, biggest_challenges, support_frequency,
			avg_resolution_time, biggest_frustration, wish_list,
			strategic_goals, budget_range, timeline, expansion_plans, priority_ranking, decision_makers,
			ai_summary, ai_recommendations, follow_up_questions, follow_up_responses,
			status, completed_at, completed_by, version, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18,
			$19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34
		)
		ON CONFLICT (tenant_id, school_id) DO UPDATE SET
			total_devices = EXCLUDED.total_devices,
			device_types = EXCLUDED.device_types,
			network_quality = EXCLUDED.network_quality,
			internet_speed = EXCLUDED.internet_speed,
			lms_platform = EXCLUDED.lms_platform,
			existing_software = EXCLUDED.existing_software,
			it_staff_count = EXCLUDED.it_staff_count,
			device_age = EXCLUDED.device_age,
			pain_points = EXCLUDED.pain_points,
			support_satisfaction = EXCLUDED.support_satisfaction,
			biggest_challenges = EXCLUDED.biggest_challenges,
			support_frequency = EXCLUDED.support_frequency,
			avg_resolution_time = EXCLUDED.avg_resolution_time,
			biggest_frustration = EXCLUDED.biggest_frustration,
			wish_list = EXCLUDED.wish_list,
			strategic_goals = EXCLUDED.strategic_goals,
			budget_range = EXCLUDED.budget_range,
			timeline = EXCLUDED.timeline,
			expansion_plans = EXCLUDED.expansion_plans,
			priority_ranking = EXCLUDED.priority_ranking,
			decision_makers = EXCLUDED.decision_makers,
			ai_summary = EXCLUDED.ai_summary,
			ai_recommendations = EXCLUDED.ai_recommendations,
			follow_up_questions = EXCLUDED.follow_up_questions,
			follow_up_responses = EXCLUDED.follow_up_responses,
			status = EXCLUDED.status,
			completed_at = EXCLUDED.completed_at,
			completed_by = EXCLUDED.completed_by,
			version = edtech_profiles.version + 1,
			updated_at = EXCLUDED.updated_at
	`,
		p.ID, p.TenantID, p.SchoolID,
		p.TotalDevices, deviceTypesJSON, p.NetworkQuality, p.InternetSpeed, p.LMSPlatform,
		existingSoftwareJSON, p.ITStaffCount, p.DeviceAge,
		painPointsJSON, p.SupportSatisfaction, p.BiggestChallenges, p.SupportFrequency,
		p.AvgResolutionTime, p.BiggestFrustration, p.WishList,
		p.StrategicGoals, p.BudgetRange, p.Timeline, p.ExpansionPlans, priorityRankingJSON, p.DecisionMakers,
		p.AISummary, aiRecsJSON, followUpQsJSON, followUpRsJSON,
		p.Status, p.CompletedAt, p.CompletedBy, p.Version, p.CreatedAt, p.UpdatedAt,
	)
	return err
}

// GetBySchoolID retrieves the profile for a school.
func (r *EdTechProfilesRepo) GetBySchoolID(ctx context.Context, tenantID, schoolID string) (models.EdTechProfile, error) {
	var p models.EdTechProfile
	var deviceTypesJSON, existingSoftwareJSON, painPointsJSON []byte
	var aiRecsJSON, followUpQsJSON, followUpRsJSON, priorityRankingJSON []byte

	row := r.pool.QueryRow(ctx, `
		SELECT id, tenant_id, school_id,
			total_devices, device_types, network_quality, internet_speed, lms_platform,
			existing_software, it_staff_count, device_age,
			pain_points, support_satisfaction, biggest_challenges, support_frequency,
			avg_resolution_time, biggest_frustration, wish_list,
			strategic_goals, budget_range, timeline, expansion_plans, priority_ranking, decision_makers,
			ai_summary, ai_recommendations, follow_up_questions, follow_up_responses,
			status, completed_at, completed_by, version, created_at, updated_at
		FROM edtech_profiles
		WHERE tenant_id = $1 AND school_id = $2
	`, tenantID, schoolID)

	err := row.Scan(
		&p.ID, &p.TenantID, &p.SchoolID,
		&p.TotalDevices, &deviceTypesJSON, &p.NetworkQuality, &p.InternetSpeed, &p.LMSPlatform,
		&existingSoftwareJSON, &p.ITStaffCount, &p.DeviceAge,
		&painPointsJSON, &p.SupportSatisfaction, &p.BiggestChallenges, &p.SupportFrequency,
		&p.AvgResolutionTime, &p.BiggestFrustration, &p.WishList,
		&p.StrategicGoals, &p.BudgetRange, &p.Timeline, &p.ExpansionPlans, &priorityRankingJSON, &p.DecisionMakers,
		&p.AISummary, &aiRecsJSON, &followUpQsJSON, &followUpRsJSON,
		&p.Status, &p.CompletedAt, &p.CompletedBy, &p.Version, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.EdTechProfile{}, errors.New("not found")
		}
		return models.EdTechProfile{}, err
	}

	// Unmarshal JSON fields
	json.Unmarshal(deviceTypesJSON, &p.DeviceTypes)
	json.Unmarshal(existingSoftwareJSON, &p.ExistingSoftware)
	json.Unmarshal(painPointsJSON, &p.PainPoints)
	json.Unmarshal(aiRecsJSON, &p.AIRecommendations)
	json.Unmarshal(followUpQsJSON, &p.FollowUpQuestions)
	json.Unmarshal(followUpRsJSON, &p.FollowUpResponses)
	json.Unmarshal(priorityRankingJSON, &p.PriorityRanking)

	// Initialize nil slices/maps
	if p.ExistingSoftware == nil {
		p.ExistingSoftware = []string{}
	}
	if p.PainPoints == nil {
		p.PainPoints = []string{}
	}
	if p.BiggestChallenges == nil {
		p.BiggestChallenges = []string{}
	}
	if p.StrategicGoals == nil {
		p.StrategicGoals = []string{}
	}
	if p.PriorityRanking == nil {
		p.PriorityRanking = []string{}
	}
	if p.DecisionMakers == nil {
		p.DecisionMakers = []string{}
	}
	if p.AIRecommendations == nil {
		p.AIRecommendations = []models.AIRecommendation{}
	}
	if p.FollowUpQuestions == nil {
		p.FollowUpQuestions = []models.FollowUpQuestion{}
	}
	if p.FollowUpResponses == nil {
		p.FollowUpResponses = map[string]string{}
	}

	return p, nil
}

// GetByID retrieves a profile by its ID.
func (r *EdTechProfilesRepo) GetByID(ctx context.Context, tenantID, id string) (models.EdTechProfile, error) {
	var p models.EdTechProfile
	var deviceTypesJSON, existingSoftwareJSON, painPointsJSON []byte
	var aiRecsJSON, followUpQsJSON, followUpRsJSON, priorityRankingJSON []byte

	row := r.pool.QueryRow(ctx, `
		SELECT id, tenant_id, school_id,
			total_devices, device_types, network_quality, internet_speed, lms_platform,
			existing_software, it_staff_count, device_age,
			pain_points, support_satisfaction, biggest_challenges, support_frequency,
			avg_resolution_time, biggest_frustration, wish_list,
			strategic_goals, budget_range, timeline, expansion_plans, priority_ranking, decision_makers,
			ai_summary, ai_recommendations, follow_up_questions, follow_up_responses,
			status, completed_at, completed_by, version, created_at, updated_at
		FROM edtech_profiles
		WHERE tenant_id = $1 AND id = $2
	`, tenantID, id)

	err := row.Scan(
		&p.ID, &p.TenantID, &p.SchoolID,
		&p.TotalDevices, &deviceTypesJSON, &p.NetworkQuality, &p.InternetSpeed, &p.LMSPlatform,
		&existingSoftwareJSON, &p.ITStaffCount, &p.DeviceAge,
		&painPointsJSON, &p.SupportSatisfaction, &p.BiggestChallenges, &p.SupportFrequency,
		&p.AvgResolutionTime, &p.BiggestFrustration, &p.WishList,
		&p.StrategicGoals, &p.BudgetRange, &p.Timeline, &p.ExpansionPlans, &priorityRankingJSON, &p.DecisionMakers,
		&p.AISummary, &aiRecsJSON, &followUpQsJSON, &followUpRsJSON,
		&p.Status, &p.CompletedAt, &p.CompletedBy, &p.Version, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.EdTechProfile{}, errors.New("not found")
		}
		return models.EdTechProfile{}, err
	}

	// Unmarshal JSON fields
	json.Unmarshal(deviceTypesJSON, &p.DeviceTypes)
	json.Unmarshal(existingSoftwareJSON, &p.ExistingSoftware)
	json.Unmarshal(painPointsJSON, &p.PainPoints)
	json.Unmarshal(aiRecsJSON, &p.AIRecommendations)
	json.Unmarshal(followUpQsJSON, &p.FollowUpQuestions)
	json.Unmarshal(followUpRsJSON, &p.FollowUpResponses)
	json.Unmarshal(priorityRankingJSON, &p.PriorityRanking)

	return p, nil
}

// UpdateAI updates only the AI-related fields of a profile.
func (r *EdTechProfilesRepo) UpdateAI(ctx context.Context, tenantID, id, summary string, recommendations []models.AIRecommendation, questions []models.FollowUpQuestion) error {
	aiRecsJSON, _ := json.Marshal(recommendations)
	followUpQsJSON, _ := json.Marshal(questions)

	_, err := r.pool.Exec(ctx, `
		UPDATE edtech_profiles
		SET ai_summary = $3, ai_recommendations = $4, follow_up_questions = $5, updated_at = $6
		WHERE tenant_id = $1 AND id = $2
	`, tenantID, id, summary, aiRecsJSON, followUpQsJSON, time.Now().UTC())
	return err
}

// UpdateFollowUpResponses updates the follow-up responses.
func (r *EdTechProfilesRepo) UpdateFollowUpResponses(ctx context.Context, tenantID, id string, responses map[string]string) error {
	responsesJSON, _ := json.Marshal(responses)
	_, err := r.pool.Exec(ctx, `
		UPDATE edtech_profiles
		SET follow_up_responses = $3, updated_at = $4
		WHERE tenant_id = $1 AND id = $2
	`, tenantID, id, responsesJSON, time.Now().UTC())
	return err
}

// Complete marks a profile as completed.
func (r *EdTechProfilesRepo) Complete(ctx context.Context, tenantID, id, completedBy string) error {
	now := time.Now().UTC()
	_, err := r.pool.Exec(ctx, `
		UPDATE edtech_profiles
		SET status = $3, completed_at = $4, completed_by = $5, updated_at = $4
		WHERE tenant_id = $1 AND id = $2
	`, tenantID, id, models.EdTechProfileCompleted, now, completedBy)
	return err
}

// History repo
type EdTechProfileHistoryRepo struct{ pool *pgxpool.Pool }

// Create adds a history entry.
func (r *EdTechProfileHistoryRepo) Create(ctx context.Context, h models.EdTechProfileHistory) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO edtech_profile_history (id, profile_id, snapshot, changed_by, change_reason, changed_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, h.ID, h.ProfileID, h.Snapshot, h.ChangedBy, h.ChangeReason, h.ChangedAt)
	return err
}

// List returns history for a profile.
func (r *EdTechProfileHistoryRepo) List(ctx context.Context, profileID string, limit int) ([]models.EdTechProfileHistory, error) {
	if limit <= 0 {
		limit = 20
	}
	rows, err := r.pool.Query(ctx, `
		SELECT id, profile_id, snapshot, changed_by, change_reason, changed_at
		FROM edtech_profile_history
		WHERE profile_id = $1
		ORDER BY changed_at DESC
		LIMIT $2
	`, profileID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.EdTechProfileHistory
	for rows.Next() {
		var h models.EdTechProfileHistory
		if err := rows.Scan(&h.ID, &h.ProfileID, &h.Snapshot, &h.ChangedBy, &h.ChangeReason, &h.ChangedAt); err != nil {
			return nil, err
		}
		out = append(out, h)
	}
	return out, nil
}
