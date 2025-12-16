package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/edvirons/ssp/ims/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DemoLeadsRepo struct{ pool *pgxpool.Pool }

// Create creates a new demo lead.
func (r *DemoLeadsRepo) Create(ctx context.Context, lead models.DemoLead) error {
	tagsJSON, _ := json.Marshal(lead.Tags)

	_, err := r.pool.Exec(ctx, `
		INSERT INTO demo_leads (
			id, tenant_id, school_id, school_name, contact_name, contact_email, contact_phone, contact_role,
			county_code, county_name, sub_county_code, sub_county_name,
			stage, stage_changed_at, estimated_value, estimated_devices, probability, expected_close_date,
			lead_source, assigned_to, notes, tags, lost_reason, lost_notes, created_by, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27
		)
	`,
		lead.ID, lead.TenantID, lead.SchoolID, lead.SchoolName, lead.ContactName, lead.ContactEmail,
		lead.ContactPhone, lead.ContactRole, lead.CountyCode, lead.CountyName, lead.SubCountyCode, lead.SubCountyName,
		lead.Stage, lead.StageChangedAt, lead.EstimatedValue, lead.EstimatedDevices, lead.Probability, lead.ExpectedCloseDate,
		lead.LeadSource, lead.AssignedTo, lead.Notes, tagsJSON, lead.LostReason, lead.LostNotes, lead.CreatedBy, lead.CreatedAt, lead.UpdatedAt,
	)
	return err
}

// GetByID retrieves a lead by ID.
func (r *DemoLeadsRepo) GetByID(ctx context.Context, tenantID, id string) (models.DemoLead, error) {
	var lead models.DemoLead
	var tagsJSON []byte

	row := r.pool.QueryRow(ctx, `
		SELECT id, tenant_id, school_id, school_name, contact_name, contact_email, contact_phone, contact_role,
			county_code, county_name, sub_county_code, sub_county_name,
			stage, stage_changed_at, estimated_value, estimated_devices, probability, expected_close_date,
			lead_source, assigned_to, notes, tags, lost_reason, lost_notes, created_by, created_at, updated_at
		FROM demo_leads
		WHERE tenant_id = $1 AND id = $2
	`, tenantID, id)

	err := row.Scan(
		&lead.ID, &lead.TenantID, &lead.SchoolID, &lead.SchoolName, &lead.ContactName, &lead.ContactEmail,
		&lead.ContactPhone, &lead.ContactRole, &lead.CountyCode, &lead.CountyName, &lead.SubCountyCode, &lead.SubCountyName,
		&lead.Stage, &lead.StageChangedAt, &lead.EstimatedValue, &lead.EstimatedDevices, &lead.Probability, &lead.ExpectedCloseDate,
		&lead.LeadSource, &lead.AssignedTo, &lead.Notes, &tagsJSON, &lead.LostReason, &lead.LostNotes, &lead.CreatedBy, &lead.CreatedAt, &lead.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.DemoLead{}, errors.New("not found")
		}
		return models.DemoLead{}, err
	}

	json.Unmarshal(tagsJSON, &lead.Tags)
	if lead.Tags == nil {
		lead.Tags = []string{}
	}

	return lead, nil
}

// List retrieves leads with optional filtering.
func (r *DemoLeadsRepo) List(ctx context.Context, tenantID string, filters models.DemoLeadFilters) ([]models.DemoLead, int, error) {
	var conditions []string
	var args []interface{}
	args = append(args, tenantID)
	conditions = append(conditions, "tenant_id = $1")

	argIdx := 2
	if filters.Stage != nil {
		conditions = append(conditions, fmt.Sprintf("stage = $%d", argIdx))
		args = append(args, *filters.Stage)
		argIdx++
	}
	if filters.AssignedTo != nil {
		conditions = append(conditions, fmt.Sprintf("assigned_to = $%d", argIdx))
		args = append(args, *filters.AssignedTo)
		argIdx++
	}
	if filters.LeadSource != nil {
		conditions = append(conditions, fmt.Sprintf("lead_source = $%d", argIdx))
		args = append(args, *filters.LeadSource)
		argIdx++
	}
	if filters.Search != nil && *filters.Search != "" {
		conditions = append(conditions, fmt.Sprintf("(school_name ILIKE $%d OR contact_name ILIKE $%d OR contact_email ILIKE $%d)", argIdx, argIdx, argIdx))
		args = append(args, "%"+*filters.Search+"%")
		argIdx++
	}

	whereClause := strings.Join(conditions, " AND ")

	// Count total
	var total int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM demo_leads WHERE %s", whereClause)
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Fetch rows
	limit := filters.Limit
	if limit <= 0 {
		limit = 50
	}
	offset := filters.Offset
	if offset < 0 {
		offset = 0
	}

	query := fmt.Sprintf(`
		SELECT id, tenant_id, school_id, school_name, contact_name, contact_email, contact_phone, contact_role,
			county_code, county_name, sub_county_code, sub_county_name,
			stage, stage_changed_at, estimated_value, estimated_devices, probability, expected_close_date,
			lead_source, assigned_to, notes, tags, lost_reason, lost_notes, created_by, created_at, updated_at
		FROM demo_leads
		WHERE %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var leads []models.DemoLead
	for rows.Next() {
		var lead models.DemoLead
		var tagsJSON []byte
		if err := rows.Scan(
			&lead.ID, &lead.TenantID, &lead.SchoolID, &lead.SchoolName, &lead.ContactName, &lead.ContactEmail,
			&lead.ContactPhone, &lead.ContactRole, &lead.CountyCode, &lead.CountyName, &lead.SubCountyCode, &lead.SubCountyName,
			&lead.Stage, &lead.StageChangedAt, &lead.EstimatedValue, &lead.EstimatedDevices, &lead.Probability, &lead.ExpectedCloseDate,
			&lead.LeadSource, &lead.AssignedTo, &lead.Notes, &tagsJSON, &lead.LostReason, &lead.LostNotes, &lead.CreatedBy, &lead.CreatedAt, &lead.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		json.Unmarshal(tagsJSON, &lead.Tags)
		if lead.Tags == nil {
			lead.Tags = []string{}
		}
		leads = append(leads, lead)
	}

	if leads == nil {
		leads = []models.DemoLead{}
	}

	return leads, total, nil
}

// Update updates a lead.
func (r *DemoLeadsRepo) Update(ctx context.Context, tenantID, id string, req models.UpdateDemoLeadRequest) error {
	var sets []string
	var args []interface{}
	args = append(args, tenantID, id)
	argIdx := 3

	if req.SchoolName != nil {
		sets = append(sets, fmt.Sprintf("school_name = $%d", argIdx))
		args = append(args, *req.SchoolName)
		argIdx++
	}
	if req.ContactName != nil {
		sets = append(sets, fmt.Sprintf("contact_name = $%d", argIdx))
		args = append(args, *req.ContactName)
		argIdx++
	}
	if req.ContactEmail != nil {
		sets = append(sets, fmt.Sprintf("contact_email = $%d", argIdx))
		args = append(args, *req.ContactEmail)
		argIdx++
	}
	if req.ContactPhone != nil {
		sets = append(sets, fmt.Sprintf("contact_phone = $%d", argIdx))
		args = append(args, *req.ContactPhone)
		argIdx++
	}
	if req.ContactRole != nil {
		sets = append(sets, fmt.Sprintf("contact_role = $%d", argIdx))
		args = append(args, *req.ContactRole)
		argIdx++
	}
	if req.CountyCode != nil {
		sets = append(sets, fmt.Sprintf("county_code = $%d", argIdx))
		args = append(args, *req.CountyCode)
		argIdx++
	}
	if req.CountyName != nil {
		sets = append(sets, fmt.Sprintf("county_name = $%d", argIdx))
		args = append(args, *req.CountyName)
		argIdx++
	}
	if req.SubCountyCode != nil {
		sets = append(sets, fmt.Sprintf("sub_county_code = $%d", argIdx))
		args = append(args, *req.SubCountyCode)
		argIdx++
	}
	if req.SubCountyName != nil {
		sets = append(sets, fmt.Sprintf("sub_county_name = $%d", argIdx))
		args = append(args, *req.SubCountyName)
		argIdx++
	}
	if req.EstimatedValue != nil {
		sets = append(sets, fmt.Sprintf("estimated_value = $%d", argIdx))
		args = append(args, *req.EstimatedValue)
		argIdx++
	}
	if req.EstimatedDevices != nil {
		sets = append(sets, fmt.Sprintf("estimated_devices = $%d", argIdx))
		args = append(args, *req.EstimatedDevices)
		argIdx++
	}
	if req.Probability != nil {
		sets = append(sets, fmt.Sprintf("probability = $%d", argIdx))
		args = append(args, *req.Probability)
		argIdx++
	}
	if req.ExpectedCloseDate != nil {
		sets = append(sets, fmt.Sprintf("expected_close_date = $%d", argIdx))
		args = append(args, *req.ExpectedCloseDate)
		argIdx++
	}
	if req.AssignedTo != nil {
		sets = append(sets, fmt.Sprintf("assigned_to = $%d", argIdx))
		args = append(args, *req.AssignedTo)
		argIdx++
	}
	if req.Notes != nil {
		sets = append(sets, fmt.Sprintf("notes = $%d", argIdx))
		args = append(args, *req.Notes)
		argIdx++
	}
	if req.Tags != nil {
		tagsJSON, _ := json.Marshal(*req.Tags)
		sets = append(sets, fmt.Sprintf("tags = $%d", argIdx))
		args = append(args, tagsJSON)
		argIdx++
	}

	if len(sets) == 0 {
		return nil
	}

	sets = append(sets, fmt.Sprintf("updated_at = $%d", argIdx))
	args = append(args, time.Now().UTC())

	query := fmt.Sprintf("UPDATE demo_leads SET %s WHERE tenant_id = $1 AND id = $2", strings.Join(sets, ", "))
	_, err := r.pool.Exec(ctx, query, args...)
	return err
}

// UpdateStage updates the stage of a lead.
func (r *DemoLeadsRepo) UpdateStage(ctx context.Context, tenantID, id string, stage models.DemoLeadStage, lostReason, lostNotes string) error {
	now := time.Now().UTC()
	probability := models.StageProbability[stage]

	_, err := r.pool.Exec(ctx, `
		UPDATE demo_leads
		SET stage = $3, stage_changed_at = $4, probability = $5, lost_reason = $6, lost_notes = $7, updated_at = $4
		WHERE tenant_id = $1 AND id = $2
	`, tenantID, id, stage, now, probability, lostReason, lostNotes)
	return err
}

// Delete deletes a lead.
func (r *DemoLeadsRepo) Delete(ctx context.Context, tenantID, id string) error {
	_, err := r.pool.Exec(ctx, "DELETE FROM demo_leads WHERE tenant_id = $1 AND id = $2", tenantID, id)
	return err
}

// GetPipelineSummary returns counts and values by stage.
func (r *DemoLeadsRepo) GetPipelineSummary(ctx context.Context, tenantID string) (models.PipelineSummary, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT stage, COUNT(*) as count, COALESCE(SUM(estimated_value), 0) as total_value
		FROM demo_leads
		WHERE tenant_id = $1 AND stage NOT IN ('won', 'lost')
		GROUP BY stage
	`, tenantID)
	if err != nil {
		return models.PipelineSummary{}, err
	}
	defer rows.Close()

	stageMap := make(map[models.DemoLeadStage]models.PipelineStageCount)
	var totalLeads int
	var totalValue float64

	for rows.Next() {
		var stage models.DemoLeadStage
		var count int
		var value float64
		if err := rows.Scan(&stage, &count, &value); err != nil {
			return models.PipelineSummary{}, err
		}
		stageMap[stage] = models.PipelineStageCount{Stage: stage, Count: count, TotalValue: value}
		totalLeads += count
		totalValue += value
	}

	// Build ordered list
	var stages []models.PipelineStageCount
	for _, s := range models.DemoLeadStages {
		if s == models.StageWon || s == models.StageLost {
			continue
		}
		if sc, ok := stageMap[s]; ok {
			stages = append(stages, sc)
		} else {
			stages = append(stages, models.PipelineStageCount{Stage: s, Count: 0, TotalValue: 0})
		}
	}

	var avgValue float64
	if totalLeads > 0 {
		avgValue = totalValue / float64(totalLeads)
	}

	// Get conversion rate (won / total closed)
	var won, lost int
	r.pool.QueryRow(ctx, `
		SELECT
			COUNT(*) FILTER (WHERE stage = 'won'),
			COUNT(*) FILTER (WHERE stage = 'lost')
		FROM demo_leads WHERE tenant_id = $1
	`, tenantID).Scan(&won, &lost)

	var conversionRate float64
	if won+lost > 0 {
		conversionRate = float64(won) / float64(won+lost) * 100
	}

	return models.PipelineSummary{
		Stages:         stages,
		TotalLeads:     totalLeads,
		TotalValue:     totalValue,
		AverageValue:   avgValue,
		ConversionRate: conversionRate,
	}, nil
}

// DemoLeadActivitiesRepo handles lead activities.
type DemoLeadActivitiesRepo struct{ pool *pgxpool.Pool }

// Create creates a new activity.
func (r *DemoLeadActivitiesRepo) Create(ctx context.Context, activity models.DemoLeadActivity) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO demo_lead_activities (
			id, tenant_id, lead_id, activity_type, description, from_stage, to_stage,
			scheduled_at, completed_at, created_by, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`,
		activity.ID, activity.TenantID, activity.LeadID, activity.ActivityType, activity.Description,
		activity.FromStage, activity.ToStage, activity.ScheduledAt, activity.CompletedAt,
		activity.CreatedBy, activity.CreatedAt,
	)
	return err
}

// ListByLead retrieves activities for a lead.
func (r *DemoLeadActivitiesRepo) ListByLead(ctx context.Context, tenantID, leadID string, limit int) ([]models.DemoLeadActivity, error) {
	if limit <= 0 {
		limit = 50
	}

	rows, err := r.pool.Query(ctx, `
		SELECT id, tenant_id, lead_id, activity_type, description, from_stage, to_stage,
			scheduled_at, completed_at, created_by, created_at
		FROM demo_lead_activities
		WHERE tenant_id = $1 AND lead_id = $2
		ORDER BY created_at DESC
		LIMIT $3
	`, tenantID, leadID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []models.DemoLeadActivity
	for rows.Next() {
		var a models.DemoLeadActivity
		if err := rows.Scan(
			&a.ID, &a.TenantID, &a.LeadID, &a.ActivityType, &a.Description, &a.FromStage, &a.ToStage,
			&a.ScheduledAt, &a.CompletedAt, &a.CreatedBy, &a.CreatedAt,
		); err != nil {
			return nil, err
		}
		activities = append(activities, a)
	}

	if activities == nil {
		activities = []models.DemoLeadActivity{}
	}

	return activities, nil
}

// ListRecent retrieves recent activities across all leads.
func (r *DemoLeadActivitiesRepo) ListRecent(ctx context.Context, tenantID string, limit int) ([]models.DemoLeadActivity, error) {
	if limit <= 0 {
		limit = 20
	}

	rows, err := r.pool.Query(ctx, `
		SELECT id, tenant_id, lead_id, activity_type, description, from_stage, to_stage,
			scheduled_at, completed_at, created_by, created_at
		FROM demo_lead_activities
		WHERE tenant_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`, tenantID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []models.DemoLeadActivity
	for rows.Next() {
		var a models.DemoLeadActivity
		if err := rows.Scan(
			&a.ID, &a.TenantID, &a.LeadID, &a.ActivityType, &a.Description, &a.FromStage, &a.ToStage,
			&a.ScheduledAt, &a.CompletedAt, &a.CreatedBy, &a.CreatedAt,
		); err != nil {
			return nil, err
		}
		activities = append(activities, a)
	}

	if activities == nil {
		activities = []models.DemoLeadActivity{}
	}

	return activities, nil
}

// DemoSchedulesRepo handles demo schedules.
type DemoSchedulesRepo struct{ pool *pgxpool.Pool }

// Create creates a new schedule.
func (r *DemoSchedulesRepo) Create(ctx context.Context, schedule models.DemoSchedule) error {
	attendeesJSON, _ := json.Marshal(schedule.Attendees)

	_, err := r.pool.Exec(ctx, `
		INSERT INTO demo_schedules (
			id, tenant_id, lead_id, scheduled_date, scheduled_time, duration_minutes,
			location, meeting_link, attendees, status, outcome, outcome_notes,
			reminder_sent, created_by, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
	`,
		schedule.ID, schedule.TenantID, schedule.LeadID, schedule.ScheduledDate, schedule.ScheduledTime,
		schedule.DurationMinutes, schedule.Location, schedule.MeetingLink, attendeesJSON, schedule.Status,
		schedule.Outcome, schedule.OutcomeNotes, schedule.ReminderSent, schedule.CreatedBy,
		schedule.CreatedAt, schedule.UpdatedAt,
	)
	return err
}

// GetByLead retrieves the next scheduled demo for a lead.
func (r *DemoSchedulesRepo) GetNextByLead(ctx context.Context, tenantID, leadID string) (*models.DemoSchedule, error) {
	var schedule models.DemoSchedule
	var attendeesJSON []byte

	row := r.pool.QueryRow(ctx, `
		SELECT id, tenant_id, lead_id, scheduled_date, scheduled_time, duration_minutes,
			location, meeting_link, attendees, status, outcome, outcome_notes,
			reminder_sent, created_by, created_at, updated_at
		FROM demo_schedules
		WHERE tenant_id = $1 AND lead_id = $2 AND status = 'scheduled' AND scheduled_date >= CURRENT_DATE
		ORDER BY scheduled_date ASC, scheduled_time ASC
		LIMIT 1
	`, tenantID, leadID)

	err := row.Scan(
		&schedule.ID, &schedule.TenantID, &schedule.LeadID, &schedule.ScheduledDate, &schedule.ScheduledTime,
		&schedule.DurationMinutes, &schedule.Location, &schedule.MeetingLink, &attendeesJSON, &schedule.Status,
		&schedule.Outcome, &schedule.OutcomeNotes, &schedule.ReminderSent, &schedule.CreatedBy,
		&schedule.CreatedAt, &schedule.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	json.Unmarshal(attendeesJSON, &schedule.Attendees)
	if schedule.Attendees == nil {
		schedule.Attendees = []models.DemoAttendee{}
	}

	return &schedule, nil
}

// UpdateStatus updates the status and outcome of a scheduled demo.
func (r *DemoSchedulesRepo) UpdateStatus(ctx context.Context, tenantID, id string, status models.DemoScheduleStatus, outcome, outcomeNotes string) error {
	now := time.Now().UTC()
	_, err := r.pool.Exec(ctx, `
		UPDATE demo_schedules
		SET status = $3, outcome = $4, outcome_notes = $5, updated_at = $6
		WHERE tenant_id = $1 AND id = $2
	`, tenantID, id, status, outcome, outcomeNotes, now)
	return err
}
