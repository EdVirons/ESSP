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

type PresentationsRepo struct{ pool *pgxpool.Pool }

// Create creates a new presentation.
func (r *PresentationsRepo) Create(ctx context.Context, p models.Presentation) error {
	tagsJSON, _ := json.Marshal(p.Tags)

	_, err := r.pool.Exec(ctx, `
		INSERT INTO presentations (
			id, tenant_id, title, description, type, category,
			file_key, file_name, file_size, file_type, thumbnail_key, preview_type,
			tags, version, is_active, is_featured, view_count, download_count,
			created_by, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21
		)
	`,
		p.ID, p.TenantID, p.Title, p.Description, p.Type, p.Category,
		p.FileKey, p.FileName, p.FileSize, p.FileType, p.ThumbnailKey, p.PreviewType,
		tagsJSON, p.Version, p.IsActive, p.IsFeatured, p.ViewCount, p.DownloadCount,
		p.CreatedBy, p.CreatedAt, p.UpdatedAt,
	)
	return err
}

// GetByID retrieves a presentation by ID.
func (r *PresentationsRepo) GetByID(ctx context.Context, tenantID, id string) (models.Presentation, error) {
	var p models.Presentation
	var tagsJSON []byte

	row := r.pool.QueryRow(ctx, `
		SELECT id, tenant_id, title, description, type, category,
			file_key, file_name, file_size, file_type, thumbnail_key, preview_type,
			tags, version, is_active, is_featured, view_count, download_count, last_viewed_at,
			created_by, updated_by, created_at, updated_at
		FROM presentations
		WHERE tenant_id = $1 AND id = $2
	`, tenantID, id)

	err := row.Scan(
		&p.ID, &p.TenantID, &p.Title, &p.Description, &p.Type, &p.Category,
		&p.FileKey, &p.FileName, &p.FileSize, &p.FileType, &p.ThumbnailKey, &p.PreviewType,
		&tagsJSON, &p.Version, &p.IsActive, &p.IsFeatured, &p.ViewCount, &p.DownloadCount, &p.LastViewedAt,
		&p.CreatedBy, &p.UpdatedBy, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Presentation{}, errors.New("not found")
		}
		return models.Presentation{}, err
	}

	_ = json.Unmarshal(tagsJSON, &p.Tags)
	if p.Tags == nil {
		p.Tags = []string{}
	}

	return p, nil
}

// List retrieves presentations with optional filtering.
func (r *PresentationsRepo) List(ctx context.Context, tenantID string, filters models.PresentationFilters) ([]models.Presentation, int, error) {
	var conditions []string
	var args []interface{}
	args = append(args, tenantID)
	conditions = append(conditions, "tenant_id = $1")

	argIdx := 2
	if filters.Type != nil {
		conditions = append(conditions, fmt.Sprintf("type = $%d", argIdx))
		args = append(args, *filters.Type)
		argIdx++
	}
	if filters.Category != nil {
		conditions = append(conditions, fmt.Sprintf("category = $%d", argIdx))
		args = append(args, *filters.Category)
		argIdx++
	}
	if filters.IsFeatured != nil {
		conditions = append(conditions, fmt.Sprintf("is_featured = $%d", argIdx))
		args = append(args, *filters.IsFeatured)
		argIdx++
	}
	if filters.IsActive != nil {
		conditions = append(conditions, fmt.Sprintf("is_active = $%d", argIdx))
		args = append(args, *filters.IsActive)
		argIdx++
	}
	if filters.Search != nil && *filters.Search != "" {
		conditions = append(conditions, fmt.Sprintf("(title ILIKE $%d OR description ILIKE $%d)", argIdx, argIdx))
		args = append(args, "%"+*filters.Search+"%")
		argIdx++
	}

	whereClause := strings.Join(conditions, " AND ")

	// Count total
	var total int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM presentations WHERE %s", whereClause)
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
		SELECT id, tenant_id, title, description, type, category,
			file_key, file_name, file_size, file_type, thumbnail_key, preview_type,
			tags, version, is_active, is_featured, view_count, download_count, last_viewed_at,
			created_by, updated_by, created_at, updated_at
		FROM presentations
		WHERE %s
		ORDER BY is_featured DESC, created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var presentations []models.Presentation
	for rows.Next() {
		var p models.Presentation
		var tagsJSON []byte
		if err := rows.Scan(
			&p.ID, &p.TenantID, &p.Title, &p.Description, &p.Type, &p.Category,
			&p.FileKey, &p.FileName, &p.FileSize, &p.FileType, &p.ThumbnailKey, &p.PreviewType,
			&tagsJSON, &p.Version, &p.IsActive, &p.IsFeatured, &p.ViewCount, &p.DownloadCount, &p.LastViewedAt,
			&p.CreatedBy, &p.UpdatedBy, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		_ = json.Unmarshal(tagsJSON, &p.Tags)
		if p.Tags == nil {
			p.Tags = []string{}
		}
		presentations = append(presentations, p)
	}

	if presentations == nil {
		presentations = []models.Presentation{}
	}

	return presentations, total, nil
}

// Update updates a presentation.
func (r *PresentationsRepo) Update(ctx context.Context, tenantID, id string, req models.UpdatePresentationRequest, updatedBy string) error {
	var sets []string
	var args []interface{}
	args = append(args, tenantID, id)
	argIdx := 3

	if req.Title != nil {
		sets = append(sets, fmt.Sprintf("title = $%d", argIdx))
		args = append(args, *req.Title)
		argIdx++
	}
	if req.Description != nil {
		sets = append(sets, fmt.Sprintf("description = $%d", argIdx))
		args = append(args, *req.Description)
		argIdx++
	}
	if req.Type != nil {
		sets = append(sets, fmt.Sprintf("type = $%d", argIdx))
		args = append(args, *req.Type)
		argIdx++
	}
	if req.Category != nil {
		sets = append(sets, fmt.Sprintf("category = $%d", argIdx))
		args = append(args, *req.Category)
		argIdx++
	}
	if req.Tags != nil {
		tagsJSON, _ := json.Marshal(*req.Tags)
		sets = append(sets, fmt.Sprintf("tags = $%d", argIdx))
		args = append(args, tagsJSON)
		argIdx++
	}
	if req.IsActive != nil {
		sets = append(sets, fmt.Sprintf("is_active = $%d", argIdx))
		args = append(args, *req.IsActive)
		argIdx++
	}
	if req.IsFeatured != nil {
		sets = append(sets, fmt.Sprintf("is_featured = $%d", argIdx))
		args = append(args, *req.IsFeatured)
		argIdx++
	}

	if len(sets) == 0 {
		return nil
	}

	sets = append(sets, fmt.Sprintf("updated_by = $%d", argIdx))
	args = append(args, updatedBy)
	argIdx++

	sets = append(sets, fmt.Sprintf("updated_at = $%d", argIdx))
	args = append(args, time.Now().UTC())

	query := fmt.Sprintf("UPDATE presentations SET %s WHERE tenant_id = $1 AND id = $2", strings.Join(sets, ", "))
	_, err := r.pool.Exec(ctx, query, args...)
	return err
}

// Delete soft deletes a presentation by setting is_active to false.
func (r *PresentationsRepo) Delete(ctx context.Context, tenantID, id string) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE presentations SET is_active = false, updated_at = $3
		WHERE tenant_id = $1 AND id = $2
	`, tenantID, id, time.Now().UTC())
	return err
}

// IncrementViewCount increments the view count.
func (r *PresentationsRepo) IncrementViewCount(ctx context.Context, tenantID, id string) error {
	now := time.Now().UTC()
	_, err := r.pool.Exec(ctx, `
		UPDATE presentations
		SET view_count = view_count + 1, last_viewed_at = $3
		WHERE tenant_id = $1 AND id = $2
	`, tenantID, id, now)
	return err
}

// IncrementDownloadCount increments the download count.
func (r *PresentationsRepo) IncrementDownloadCount(ctx context.Context, tenantID, id string) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE presentations
		SET download_count = download_count + 1
		WHERE tenant_id = $1 AND id = $2
	`, tenantID, id)
	return err
}

// PresentationViewsRepo handles view tracking.
type PresentationViewsRepo struct{ pool *pgxpool.Pool }

// Create records a view event.
func (r *PresentationViewsRepo) Create(ctx context.Context, view models.PresentationView) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO presentation_views (id, tenant_id, presentation_id, viewed_by, viewed_at, context, duration_seconds)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, view.ID, view.TenantID, view.PresentationID, view.ViewedBy, view.ViewedAt, view.Context, view.DurationSeconds)
	return err
}

// SalesMetricsDailyRepo handles daily sales metrics.
type SalesMetricsDailyRepo struct{ pool *pgxpool.Pool }

// GetOrCreate gets or creates the daily metrics record.
func (r *SalesMetricsDailyRepo) GetOrCreate(ctx context.Context, tenantID string, date time.Time) (models.SalesMetricsDaily, error) {
	var m models.SalesMetricsDaily
	metricDate := date.Truncate(24 * time.Hour)

	row := r.pool.QueryRow(ctx, `
		SELECT id, tenant_id, metric_date, new_leads, leads_contacted, demos_scheduled, demos_completed,
			proposals_sent, deals_won, deals_lost, pipeline_value, won_value, lost_value,
			calls_made, emails_sent, meetings_held, created_at, updated_at
		FROM sales_metrics_daily
		WHERE tenant_id = $1 AND metric_date = $2
	`, tenantID, metricDate)

	err := row.Scan(
		&m.ID, &m.TenantID, &m.MetricDate, &m.NewLeads, &m.LeadsContacted, &m.DemosScheduled, &m.DemosCompleted,
		&m.ProposalsSent, &m.DealsWon, &m.DealsLost, &m.PipelineValue, &m.WonValue, &m.LostValue,
		&m.CallsMade, &m.EmailsSent, &m.MeetingsHeld, &m.CreatedAt, &m.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Create new record
			now := time.Now().UTC()
			m = models.SalesMetricsDaily{
				ID:         generateID(),
				TenantID:   tenantID,
				MetricDate: metricDate,
				CreatedAt:  now,
				UpdatedAt:  now,
			}
			_, err = r.pool.Exec(ctx, `
				INSERT INTO sales_metrics_daily (id, tenant_id, metric_date, created_at, updated_at)
				VALUES ($1, $2, $3, $4, $5)
			`, m.ID, m.TenantID, m.MetricDate, m.CreatedAt, m.UpdatedAt)
			if err != nil {
				return models.SalesMetricsDaily{}, err
			}
			return m, nil
		}
		return models.SalesMetricsDaily{}, err
	}

	return m, nil
}

// GetSummary retrieves aggregated metrics for a date range.
func (r *SalesMetricsDailyRepo) GetSummary(ctx context.Context, tenantID string, startDate, endDate time.Time) (models.SalesMetricsSummary, error) {
	var s models.SalesMetricsSummary

	err := r.pool.QueryRow(ctx, `
		SELECT
			COALESCE(SUM(new_leads), 0),
			COALESCE(SUM(demos_scheduled), 0),
			COALESCE(SUM(demos_completed), 0),
			COALESCE(SUM(proposals_sent), 0),
			COALESCE(SUM(deals_won), 0),
			COALESCE(SUM(deals_lost), 0),
			COALESCE(SUM(won_value), 0)
		FROM sales_metrics_daily
		WHERE tenant_id = $1 AND metric_date BETWEEN $2 AND $3
	`, tenantID, startDate, endDate).Scan(
		&s.NewLeadsThisPeriod, &s.DemosScheduled, &s.DemosCompleted,
		&s.ProposalsSent, &s.DealsWon, &s.DealsLost, &s.WonValueThisPeriod,
	)
	if err != nil {
		return models.SalesMetricsSummary{}, err
	}

	// Calculate rates
	if s.DealsWon+s.DealsLost > 0 {
		s.WinRate = float64(s.DealsWon) / float64(s.DealsWon+s.DealsLost) * 100
	}
	if s.DealsWon > 0 {
		s.AverageDealSize = s.WonValueThisPeriod / float64(s.DealsWon)
	}

	return s, nil
}

// IncrementMetric increments a specific metric for today.
func (r *SalesMetricsDailyRepo) IncrementMetric(ctx context.Context, tenantID, metric string, value float64) error {
	today := time.Now().UTC().Truncate(24 * time.Hour)

	// Ensure record exists
	_, err := r.GetOrCreate(ctx, tenantID, today)
	if err != nil {
		return err
	}

	query := fmt.Sprintf(`
		UPDATE sales_metrics_daily
		SET %s = %s + $3, updated_at = $4
		WHERE tenant_id = $1 AND metric_date = $2
	`, metric, metric)

	_, err = r.pool.Exec(ctx, query, tenantID, today, value, time.Now().UTC())
	return err
}

func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
