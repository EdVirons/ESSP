package store

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/edvirons/ssp/ims/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MarketingKBRepo struct{ pool *pgxpool.Pool }

// ==================== MKB Articles ====================

func (r *MarketingKBRepo) CreateArticle(ctx context.Context, a models.MKBArticle) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO marketing_kb_articles (
			id, tenant_id, title, slug, summary, content, content_type,
			personas, context_tags, tags, version, status,
			usage_count, last_used_at,
			created_by_id, created_by_name, updated_by_id, updated_by_name,
			approved_at, approved_by_id, approved_by_name,
			created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23)
	`, a.ID, a.TenantID, a.Title, a.Slug, a.Summary, a.Content, a.ContentType,
		a.Personas, a.ContextTags, a.Tags, a.Version, a.Status,
		a.UsageCount, a.LastUsedAt,
		a.CreatedByID, a.CreatedByName, a.UpdatedByID, a.UpdatedByName,
		a.ApprovedAt, a.ApprovedByID, a.ApprovedByName,
		a.CreatedAt, a.UpdatedAt)
	return err
}

func (r *MarketingKBRepo) GetArticleByID(ctx context.Context, tenantID, id string) (models.MKBArticle, error) {
	var a models.MKBArticle
	row := r.pool.QueryRow(ctx, `
		SELECT id, tenant_id, title, slug, summary, content, content_type,
			personas, context_tags, tags, version, status,
			usage_count, last_used_at,
			created_by_id, created_by_name, updated_by_id, updated_by_name,
			approved_at, approved_by_id, approved_by_name,
			created_at, updated_at
		FROM marketing_kb_articles WHERE tenant_id=$1 AND id=$2
	`, tenantID, id)
	if err := row.Scan(
		&a.ID, &a.TenantID, &a.Title, &a.Slug, &a.Summary, &a.Content, &a.ContentType,
		&a.Personas, &a.ContextTags, &a.Tags, &a.Version, &a.Status,
		&a.UsageCount, &a.LastUsedAt,
		&a.CreatedByID, &a.CreatedByName, &a.UpdatedByID, &a.UpdatedByName,
		&a.ApprovedAt, &a.ApprovedByID, &a.ApprovedByName,
		&a.CreatedAt, &a.UpdatedAt,
	); err != nil {
		return models.MKBArticle{}, errors.New("not found")
	}
	return a, nil
}

func (r *MarketingKBRepo) GetArticleBySlug(ctx context.Context, tenantID, slug string) (models.MKBArticle, error) {
	var a models.MKBArticle
	row := r.pool.QueryRow(ctx, `
		SELECT id, tenant_id, title, slug, summary, content, content_type,
			personas, context_tags, tags, version, status,
			usage_count, last_used_at,
			created_by_id, created_by_name, updated_by_id, updated_by_name,
			approved_at, approved_by_id, approved_by_name,
			created_at, updated_at
		FROM marketing_kb_articles WHERE tenant_id=$1 AND slug=$2
	`, tenantID, slug)
	if err := row.Scan(
		&a.ID, &a.TenantID, &a.Title, &a.Slug, &a.Summary, &a.Content, &a.ContentType,
		&a.Personas, &a.ContextTags, &a.Tags, &a.Version, &a.Status,
		&a.UsageCount, &a.LastUsedAt,
		&a.CreatedByID, &a.CreatedByName, &a.UpdatedByID, &a.UpdatedByName,
		&a.ApprovedAt, &a.ApprovedByID, &a.ApprovedByName,
		&a.CreatedAt, &a.UpdatedAt,
	); err != nil {
		return models.MKBArticle{}, errors.New("not found")
	}
	return a, nil
}

func (r *MarketingKBRepo) UpdateArticle(ctx context.Context, a models.MKBArticle) error {
	result, err := r.pool.Exec(ctx, `
		UPDATE marketing_kb_articles SET
			title = $3,
			slug = $4,
			summary = $5,
			content = $6,
			content_type = $7,
			personas = $8,
			context_tags = $9,
			tags = $10,
			version = $11,
			status = $12,
			updated_by_id = $13,
			updated_by_name = $14,
			updated_at = $15
		WHERE tenant_id = $1 AND id = $2
	`, a.TenantID, a.ID, a.Title, a.Slug, a.Summary, a.Content, a.ContentType,
		a.Personas, a.ContextTags, a.Tags, a.Version, a.Status,
		a.UpdatedByID, a.UpdatedByName, a.UpdatedAt)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errors.New("not found")
	}
	return nil
}

// DeleteArticle archives the article (soft delete)
func (r *MarketingKBRepo) DeleteArticle(ctx context.Context, tenantID, id string) error {
	result, err := r.pool.Exec(ctx, `
		UPDATE marketing_kb_articles SET status = 'archived', updated_at = $3
		WHERE tenant_id = $1 AND id = $2
	`, tenantID, id, time.Now().UTC())
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errors.New("not found")
	}
	return nil
}

func (r *MarketingKBRepo) ListArticles(ctx context.Context, p models.MKBArticleListParams) ([]models.MKBArticle, string, error) {
	conds := []string{"tenant_id=$1"}
	args := []any{p.TenantID}
	argN := 2

	if p.ContentType != "" {
		conds = append(conds, "content_type = $"+itoa(argN))
		args = append(args, p.ContentType)
		argN++
	}
	if p.Persona != "" {
		conds = append(conds, "$"+itoa(argN)+" = ANY(personas)")
		args = append(args, p.Persona)
		argN++
	}
	if p.ContextTag != "" {
		conds = append(conds, "$"+itoa(argN)+" = ANY(context_tags)")
		args = append(args, p.ContextTag)
		argN++
	}
	if p.Status != "" {
		conds = append(conds, "status = $"+itoa(argN))
		args = append(args, p.Status)
		argN++
	}
	if p.Query != "" {
		conds = append(conds, "to_tsvector('english', title || ' ' || summary || ' ' || content) @@ plainto_tsquery('english', $"+itoa(argN)+")")
		args = append(args, p.Query)
		argN++
	}
	if p.HasCursor {
		conds = append(conds, "(updated_at, id) < ($"+itoa(argN)+", $"+itoa(argN+1)+")")
		args = append(args, p.CursorTime, p.CursorID)
		argN += 2
	}

	limitPlus := p.Limit + 1
	args = append(args, limitPlus)

	sql := `
		SELECT id, tenant_id, title, slug, summary, content, content_type,
			personas, context_tags, tags, version, status,
			usage_count, last_used_at,
			created_by_id, created_by_name, updated_by_id, updated_by_name,
			approved_at, approved_by_id, approved_by_name,
			created_at, updated_at
		FROM marketing_kb_articles
		WHERE ` + strings.Join(conds, " AND ") + `
		ORDER BY updated_at DESC, id DESC
		LIMIT $` + itoa(argN)

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	out := []models.MKBArticle{}
	for rows.Next() {
		var a models.MKBArticle
		if err := rows.Scan(
			&a.ID, &a.TenantID, &a.Title, &a.Slug, &a.Summary, &a.Content, &a.ContentType,
			&a.Personas, &a.ContextTags, &a.Tags, &a.Version, &a.Status,
			&a.UsageCount, &a.LastUsedAt,
			&a.CreatedByID, &a.CreatedByName, &a.UpdatedByID, &a.UpdatedByName,
			&a.ApprovedAt, &a.ApprovedByID, &a.ApprovedByName,
			&a.CreatedAt, &a.UpdatedAt,
		); err != nil {
			return nil, "", err
		}
		out = append(out, a)
	}

	next := ""
	if len(out) > p.Limit {
		last := out[p.Limit-1]
		next = EncodeCursor(last.UpdatedAt, last.ID)
		out = out[:p.Limit]
	}
	return out, next, nil
}

// SearchArticles performs full-text search with ranking
func (r *MarketingKBRepo) SearchArticles(ctx context.Context, tenantID, query string, limit int) ([]models.MKBArticle, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, tenant_id, title, slug, summary, content, content_type,
			personas, context_tags, tags, version, status,
			usage_count, last_used_at,
			created_by_id, created_by_name, updated_by_id, updated_by_name,
			approved_at, approved_by_id, approved_by_name,
			created_at, updated_at,
			ts_rank(to_tsvector('english', title || ' ' || summary || ' ' || content), plainto_tsquery('english', $2)) as rank
		FROM marketing_kb_articles
		WHERE tenant_id = $1
			AND status = 'approved'
			AND to_tsvector('english', title || ' ' || summary || ' ' || content) @@ plainto_tsquery('english', $2)
		ORDER BY rank DESC, updated_at DESC
		LIMIT $3
	`, tenantID, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.MKBArticle
	for rows.Next() {
		var a models.MKBArticle
		var rank float64
		if err := rows.Scan(
			&a.ID, &a.TenantID, &a.Title, &a.Slug, &a.Summary, &a.Content, &a.ContentType,
			&a.Personas, &a.ContextTags, &a.Tags, &a.Version, &a.Status,
			&a.UsageCount, &a.LastUsedAt,
			&a.CreatedByID, &a.CreatedByName, &a.UpdatedByID, &a.UpdatedByName,
			&a.ApprovedAt, &a.ApprovedByID, &a.ApprovedByName,
			&a.CreatedAt, &a.UpdatedAt, &rank,
		); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, nil
}

// GetArticleStats returns aggregate statistics for marketing KB articles
func (r *MarketingKBRepo) GetArticleStats(ctx context.Context, tenantID string) (models.MKBStats, error) {
	var stats models.MKBStats
	stats.ByContentType = make(map[string]int)
	stats.ByPersona = make(map[string]int)
	stats.ByContextTag = make(map[string]int)

	// Get total count
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM marketing_kb_articles WHERE tenant_id = $1 AND status != 'archived'`, tenantID).Scan(&stats.Total)
	if err != nil {
		return stats, err
	}

	// Get approved count
	err = r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM marketing_kb_articles WHERE tenant_id = $1 AND status = 'approved'`, tenantID).Scan(&stats.Approved)
	if err != nil {
		return stats, err
	}

	// Get in review count
	err = r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM marketing_kb_articles WHERE tenant_id = $1 AND status = 'review'`, tenantID).Scan(&stats.InReview)
	if err != nil {
		return stats, err
	}

	// Get draft count
	err = r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM marketing_kb_articles WHERE tenant_id = $1 AND status = 'draft'`, tenantID).Scan(&stats.Draft)
	if err != nil {
		return stats, err
	}

	// Count by content type
	rows, err := r.pool.Query(ctx, `
		SELECT content_type, COUNT(*) FROM marketing_kb_articles
		WHERE tenant_id = $1 AND status != 'archived'
		GROUP BY content_type
	`, tenantID)
	if err != nil {
		return stats, err
	}
	for rows.Next() {
		var ct string
		var count int
		if err := rows.Scan(&ct, &count); err != nil {
			rows.Close()
			return stats, err
		}
		stats.ByContentType[ct] = count
	}
	rows.Close()

	// Count by persona (unnest the array)
	rows, err = r.pool.Query(ctx, `
		SELECT unnest(personas) as persona, COUNT(*) FROM marketing_kb_articles
		WHERE tenant_id = $1 AND status != 'archived'
		GROUP BY persona
	`, tenantID)
	if err != nil {
		return stats, err
	}
	for rows.Next() {
		var p string
		var count int
		if err := rows.Scan(&p, &count); err != nil {
			rows.Close()
			return stats, err
		}
		stats.ByPersona[p] = count
	}
	rows.Close()

	// Count by context tag (unnest the array)
	rows, err = r.pool.Query(ctx, `
		SELECT unnest(context_tags) as tag, COUNT(*) FROM marketing_kb_articles
		WHERE tenant_id = $1 AND status != 'archived'
		GROUP BY tag
	`, tenantID)
	if err != nil {
		return stats, err
	}
	for rows.Next() {
		var t string
		var count int
		if err := rows.Scan(&t, &count); err != nil {
			rows.Close()
			return stats, err
		}
		stats.ByContextTag[t] = count
	}
	rows.Close()

	return stats, nil
}

// ApproveArticle changes article status to approved
func (r *MarketingKBRepo) ApproveArticle(ctx context.Context, tenantID, id, userID, userName string) error {
	now := time.Now().UTC()
	result, err := r.pool.Exec(ctx, `
		UPDATE marketing_kb_articles SET
			status = 'approved',
			approved_at = $3,
			approved_by_id = $4,
			approved_by_name = $5,
			updated_by_id = $4,
			updated_by_name = $5,
			updated_at = $3
		WHERE tenant_id = $1 AND id = $2 AND status IN ('draft', 'review')
	`, tenantID, id, now, userID, userName)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errors.New("not found or already approved")
	}
	return nil
}

// SubmitForReview changes article status to review
func (r *MarketingKBRepo) SubmitForReview(ctx context.Context, tenantID, id, userID, userName string) error {
	now := time.Now().UTC()
	result, err := r.pool.Exec(ctx, `
		UPDATE marketing_kb_articles SET
			status = 'review',
			updated_by_id = $3,
			updated_by_name = $4,
			updated_at = $5
		WHERE tenant_id = $1 AND id = $2 AND status = 'draft'
	`, tenantID, id, userID, userName, now)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errors.New("not found or not a draft")
	}
	return nil
}

// RecordUsage increments the usage count for an article
func (r *MarketingKBRepo) RecordUsage(ctx context.Context, tenantID, id string) error {
	now := time.Now().UTC()
	result, err := r.pool.Exec(ctx, `
		UPDATE marketing_kb_articles SET
			usage_count = usage_count + 1,
			last_used_at = $3
		WHERE tenant_id = $1 AND id = $2
	`, tenantID, id, now)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errors.New("not found")
	}
	return nil
}

// ArticleSlugExists checks if a slug already exists for a tenant (optionally excluding an article ID)
func (r *MarketingKBRepo) ArticleSlugExists(ctx context.Context, tenantID, slug, excludeID string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM marketing_kb_articles WHERE tenant_id = $1 AND slug = $2`
	args := []any{tenantID, slug}
	if excludeID != "" {
		query += ` AND id != $3`
		args = append(args, excludeID)
	}
	query += `)`
	err := r.pool.QueryRow(ctx, query, args...).Scan(&exists)
	return exists, err
}

// GetArticlesByIDs fetches multiple articles by their IDs
func (r *MarketingKBRepo) GetArticlesByIDs(ctx context.Context, tenantID string, ids []string) ([]models.MKBArticle, error) {
	if len(ids) == 0 {
		return []models.MKBArticle{}, nil
	}

	rows, err := r.pool.Query(ctx, `
		SELECT id, tenant_id, title, slug, summary, content, content_type,
			personas, context_tags, tags, version, status,
			usage_count, last_used_at,
			created_by_id, created_by_name, updated_by_id, updated_by_name,
			approved_at, approved_by_id, approved_by_name,
			created_at, updated_at
		FROM marketing_kb_articles
		WHERE tenant_id = $1 AND id = ANY($2)
	`, tenantID, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.MKBArticle
	for rows.Next() {
		var a models.MKBArticle
		if err := rows.Scan(
			&a.ID, &a.TenantID, &a.Title, &a.Slug, &a.Summary, &a.Content, &a.ContentType,
			&a.Personas, &a.ContextTags, &a.Tags, &a.Version, &a.Status,
			&a.UsageCount, &a.LastUsedAt,
			&a.CreatedByID, &a.CreatedByName, &a.UpdatedByID, &a.UpdatedByName,
			&a.ApprovedAt, &a.ApprovedByID, &a.ApprovedByName,
			&a.CreatedAt, &a.UpdatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, nil
}

// ==================== Pitch Kits ====================

func (r *MarketingKBRepo) CreatePitchKit(ctx context.Context, pk models.PitchKit) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO marketing_pitch_kits (
			id, tenant_id, name, description, target_persona,
			context_tags, article_ids, is_template,
			created_by_id, created_by_name, updated_by_id, updated_by_name,
			created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)
	`, pk.ID, pk.TenantID, pk.Name, pk.Description, pk.TargetPersona,
		pk.ContextTags, pk.ArticleIDs, pk.IsTemplate,
		pk.CreatedByID, pk.CreatedByName, pk.UpdatedByID, pk.UpdatedByName,
		pk.CreatedAt, pk.UpdatedAt)
	return err
}

func (r *MarketingKBRepo) GetPitchKitByID(ctx context.Context, tenantID, id string) (models.PitchKit, error) {
	var pk models.PitchKit
	row := r.pool.QueryRow(ctx, `
		SELECT id, tenant_id, name, description, target_persona,
			context_tags, article_ids, is_template,
			created_by_id, created_by_name, updated_by_id, updated_by_name,
			created_at, updated_at
		FROM marketing_pitch_kits WHERE tenant_id=$1 AND id=$2
	`, tenantID, id)
	if err := row.Scan(
		&pk.ID, &pk.TenantID, &pk.Name, &pk.Description, &pk.TargetPersona,
		&pk.ContextTags, &pk.ArticleIDs, &pk.IsTemplate,
		&pk.CreatedByID, &pk.CreatedByName, &pk.UpdatedByID, &pk.UpdatedByName,
		&pk.CreatedAt, &pk.UpdatedAt,
	); err != nil {
		return models.PitchKit{}, errors.New("not found")
	}
	return pk, nil
}

func (r *MarketingKBRepo) UpdatePitchKit(ctx context.Context, pk models.PitchKit) error {
	result, err := r.pool.Exec(ctx, `
		UPDATE marketing_pitch_kits SET
			name = $3,
			description = $4,
			target_persona = $5,
			context_tags = $6,
			article_ids = $7,
			is_template = $8,
			updated_by_id = $9,
			updated_by_name = $10,
			updated_at = $11
		WHERE tenant_id = $1 AND id = $2
	`, pk.TenantID, pk.ID, pk.Name, pk.Description, pk.TargetPersona,
		pk.ContextTags, pk.ArticleIDs, pk.IsTemplate,
		pk.UpdatedByID, pk.UpdatedByName, pk.UpdatedAt)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errors.New("not found")
	}
	return nil
}

func (r *MarketingKBRepo) DeletePitchKit(ctx context.Context, tenantID, id string) error {
	result, err := r.pool.Exec(ctx, `
		DELETE FROM marketing_pitch_kits WHERE tenant_id = $1 AND id = $2
	`, tenantID, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errors.New("not found")
	}
	return nil
}

func (r *MarketingKBRepo) ListPitchKits(ctx context.Context, p models.PitchKitListParams) ([]models.PitchKit, string, error) {
	conds := []string{"tenant_id=$1"}
	args := []any{p.TenantID}
	argN := 2

	if p.TargetPersona != "" {
		conds = append(conds, "target_persona = $"+itoa(argN))
		args = append(args, p.TargetPersona)
		argN++
	}
	if p.IsTemplate != nil {
		conds = append(conds, "is_template = $"+itoa(argN))
		args = append(args, *p.IsTemplate)
		argN++
	}
	if p.HasCursor {
		conds = append(conds, "(updated_at, id) < ($"+itoa(argN)+", $"+itoa(argN+1)+")")
		args = append(args, p.CursorTime, p.CursorID)
		argN += 2
	}

	limitPlus := p.Limit + 1
	args = append(args, limitPlus)

	sql := `
		SELECT id, tenant_id, name, description, target_persona,
			context_tags, article_ids, is_template,
			created_by_id, created_by_name, updated_by_id, updated_by_name,
			created_at, updated_at
		FROM marketing_pitch_kits
		WHERE ` + strings.Join(conds, " AND ") + `
		ORDER BY updated_at DESC, id DESC
		LIMIT $` + itoa(argN)

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	out := []models.PitchKit{}
	for rows.Next() {
		var pk models.PitchKit
		if err := rows.Scan(
			&pk.ID, &pk.TenantID, &pk.Name, &pk.Description, &pk.TargetPersona,
			&pk.ContextTags, &pk.ArticleIDs, &pk.IsTemplate,
			&pk.CreatedByID, &pk.CreatedByName, &pk.UpdatedByID, &pk.UpdatedByName,
			&pk.CreatedAt, &pk.UpdatedAt,
		); err != nil {
			return nil, "", err
		}
		out = append(out, pk)
	}

	next := ""
	if len(out) > p.Limit {
		last := out[p.Limit-1]
		next = EncodeCursor(last.UpdatedAt, last.ID)
		out = out[:p.Limit]
	}
	return out, next, nil
}

// GetPitchKitWithArticles fetches a pitch kit and populates its articles
func (r *MarketingKBRepo) GetPitchKitWithArticles(ctx context.Context, tenantID, id string) (models.PitchKit, error) {
	pk, err := r.GetPitchKitByID(ctx, tenantID, id)
	if err != nil {
		return models.PitchKit{}, err
	}

	if len(pk.ArticleIDs) > 0 {
		articles, err := r.GetArticlesByIDs(ctx, tenantID, pk.ArticleIDs)
		if err != nil {
			return models.PitchKit{}, err
		}
		// Preserve order from ArticleIDs
		articleMap := make(map[string]models.MKBArticle)
		for _, a := range articles {
			articleMap[a.ID] = a
		}
		pk.Articles = make([]models.MKBArticle, 0, len(pk.ArticleIDs))
		for _, id := range pk.ArticleIDs {
			if a, ok := articleMap[id]; ok {
				pk.Articles = append(pk.Articles, a)
			}
		}
	}

	return pk, nil
}
