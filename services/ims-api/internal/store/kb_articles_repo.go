package store

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/edvirons/ssp/ims/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type KBArticleRepo struct{ pool *pgxpool.Pool }

func (r *KBArticleRepo) Create(ctx context.Context, a models.KBArticle) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO kb_articles (
			id, tenant_id, title, slug, summary, content, content_type,
			module, lifecycle_stage, tags, version, status,
			created_by_id, created_by_name, updated_by_id, updated_by_name,
			published_at, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19)
	`, a.ID, a.TenantID, a.Title, a.Slug, a.Summary, a.Content, a.ContentType,
		a.Module, a.LifecycleStage, a.Tags, a.Version, a.Status,
		a.CreatedByID, a.CreatedByName, a.UpdatedByID, a.UpdatedByName,
		a.PublishedAt, a.CreatedAt, a.UpdatedAt)
	return err
}

func (r *KBArticleRepo) GetByID(ctx context.Context, tenantID, id string) (models.KBArticle, error) {
	var a models.KBArticle
	row := r.pool.QueryRow(ctx, `
		SELECT id, tenant_id, title, slug, summary, content, content_type,
			module, lifecycle_stage, tags, version, status,
			created_by_id, created_by_name, updated_by_id, updated_by_name,
			published_at, created_at, updated_at
		FROM kb_articles WHERE tenant_id=$1 AND id=$2
	`, tenantID, id)
	if err := row.Scan(
		&a.ID, &a.TenantID, &a.Title, &a.Slug, &a.Summary, &a.Content, &a.ContentType,
		&a.Module, &a.LifecycleStage, &a.Tags, &a.Version, &a.Status,
		&a.CreatedByID, &a.CreatedByName, &a.UpdatedByID, &a.UpdatedByName,
		&a.PublishedAt, &a.CreatedAt, &a.UpdatedAt,
	); err != nil {
		return models.KBArticle{}, errors.New("not found")
	}
	return a, nil
}

func (r *KBArticleRepo) GetBySlug(ctx context.Context, tenantID, slug string) (models.KBArticle, error) {
	var a models.KBArticle
	row := r.pool.QueryRow(ctx, `
		SELECT id, tenant_id, title, slug, summary, content, content_type,
			module, lifecycle_stage, tags, version, status,
			created_by_id, created_by_name, updated_by_id, updated_by_name,
			published_at, created_at, updated_at
		FROM kb_articles WHERE tenant_id=$1 AND slug=$2
	`, tenantID, slug)
	if err := row.Scan(
		&a.ID, &a.TenantID, &a.Title, &a.Slug, &a.Summary, &a.Content, &a.ContentType,
		&a.Module, &a.LifecycleStage, &a.Tags, &a.Version, &a.Status,
		&a.CreatedByID, &a.CreatedByName, &a.UpdatedByID, &a.UpdatedByName,
		&a.PublishedAt, &a.CreatedAt, &a.UpdatedAt,
	); err != nil {
		return models.KBArticle{}, errors.New("not found")
	}
	return a, nil
}

func (r *KBArticleRepo) Update(ctx context.Context, a models.KBArticle) error {
	result, err := r.pool.Exec(ctx, `
		UPDATE kb_articles SET
			title = $3,
			slug = $4,
			summary = $5,
			content = $6,
			content_type = $7,
			module = $8,
			lifecycle_stage = $9,
			tags = $10,
			version = $11,
			status = $12,
			updated_by_id = $13,
			updated_by_name = $14,
			published_at = $15,
			updated_at = $16
		WHERE tenant_id = $1 AND id = $2
	`, a.TenantID, a.ID, a.Title, a.Slug, a.Summary, a.Content, a.ContentType,
		a.Module, a.LifecycleStage, a.Tags, a.Version, a.Status,
		a.UpdatedByID, a.UpdatedByName, a.PublishedAt, a.UpdatedAt)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errors.New("not found")
	}
	return nil
}

// Delete archives the article (soft delete)
func (r *KBArticleRepo) Delete(ctx context.Context, tenantID, id string) error {
	result, err := r.pool.Exec(ctx, `
		UPDATE kb_articles SET status = 'archived', updated_at = $3
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

func (r *KBArticleRepo) List(ctx context.Context, p models.KBArticleListParams) ([]models.KBArticle, string, error) {
	conds := []string{"tenant_id=$1"}
	args := []any{p.TenantID}
	argN := 2

	if p.ContentType != "" {
		conds = append(conds, "content_type = $"+itoa(argN))
		args = append(args, p.ContentType)
		argN++
	}
	if p.Module != "" {
		conds = append(conds, "module = $"+itoa(argN))
		args = append(args, p.Module)
		argN++
	}
	if p.LifecycleStage != "" {
		conds = append(conds, "lifecycle_stage = $"+itoa(argN))
		args = append(args, p.LifecycleStage)
		argN++
	}
	if p.Status != "" {
		conds = append(conds, "status = $"+itoa(argN))
		args = append(args, p.Status)
		argN++
	}
	if p.Query != "" {
		// Use full-text search
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
			module, lifecycle_stage, tags, version, status,
			created_by_id, created_by_name, updated_by_id, updated_by_name,
			published_at, created_at, updated_at
		FROM kb_articles
		WHERE ` + strings.Join(conds, " AND ") + `
		ORDER BY updated_at DESC, id DESC
		LIMIT $` + itoa(argN)

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	out := []models.KBArticle{}
	for rows.Next() {
		var a models.KBArticle
		if err := rows.Scan(
			&a.ID, &a.TenantID, &a.Title, &a.Slug, &a.Summary, &a.Content, &a.ContentType,
			&a.Module, &a.LifecycleStage, &a.Tags, &a.Version, &a.Status,
			&a.CreatedByID, &a.CreatedByName, &a.UpdatedByID, &a.UpdatedByName,
			&a.PublishedAt, &a.CreatedAt, &a.UpdatedAt,
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

// Search performs full-text search with ranking
func (r *KBArticleRepo) Search(ctx context.Context, tenantID, query string, limit int) ([]models.KBArticle, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, tenant_id, title, slug, summary, content, content_type,
			module, lifecycle_stage, tags, version, status,
			created_by_id, created_by_name, updated_by_id, updated_by_name,
			published_at, created_at, updated_at,
			ts_rank(to_tsvector('english', title || ' ' || summary || ' ' || content), plainto_tsquery('english', $2)) as rank
		FROM kb_articles
		WHERE tenant_id = $1
			AND status = 'published'
			AND to_tsvector('english', title || ' ' || summary || ' ' || content) @@ plainto_tsquery('english', $2)
		ORDER BY rank DESC, updated_at DESC
		LIMIT $3
	`, tenantID, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.KBArticle
	for rows.Next() {
		var a models.KBArticle
		var rank float64
		if err := rows.Scan(
			&a.ID, &a.TenantID, &a.Title, &a.Slug, &a.Summary, &a.Content, &a.ContentType,
			&a.Module, &a.LifecycleStage, &a.Tags, &a.Version, &a.Status,
			&a.CreatedByID, &a.CreatedByName, &a.UpdatedByID, &a.UpdatedByName,
			&a.PublishedAt, &a.CreatedAt, &a.UpdatedAt, &rank,
		); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, nil
}

// GetStats returns aggregate statistics for KB articles
func (r *KBArticleRepo) GetStats(ctx context.Context, tenantID string) (models.KBStats, error) {
	var stats models.KBStats
	stats.ByContentType = make(map[string]int)
	stats.ByModule = make(map[string]int)
	stats.ByLifecycleStage = make(map[string]int)

	// Get total count
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM kb_articles WHERE tenant_id = $1`, tenantID).Scan(&stats.Total)
	if err != nil {
		return stats, err
	}

	// Get published count
	err = r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM kb_articles WHERE tenant_id = $1 AND status = 'published'`, tenantID).Scan(&stats.Published)
	if err != nil {
		return stats, err
	}

	// Get draft count
	err = r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM kb_articles WHERE tenant_id = $1 AND status = 'draft'`, tenantID).Scan(&stats.Draft)
	if err != nil {
		return stats, err
	}

	// Count by content type
	rows, err := r.pool.Query(ctx, `
		SELECT content_type, COUNT(*) FROM kb_articles
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

	// Count by module
	rows, err = r.pool.Query(ctx, `
		SELECT module, COUNT(*) FROM kb_articles
		WHERE tenant_id = $1 AND status != 'archived'
		GROUP BY module
	`, tenantID)
	if err != nil {
		return stats, err
	}
	for rows.Next() {
		var m string
		var count int
		if err := rows.Scan(&m, &count); err != nil {
			rows.Close()
			return stats, err
		}
		stats.ByModule[m] = count
	}
	rows.Close()

	// Count by lifecycle stage
	rows, err = r.pool.Query(ctx, `
		SELECT lifecycle_stage, COUNT(*) FROM kb_articles
		WHERE tenant_id = $1 AND status != 'archived'
		GROUP BY lifecycle_stage
	`, tenantID)
	if err != nil {
		return stats, err
	}
	for rows.Next() {
		var ls string
		var count int
		if err := rows.Scan(&ls, &count); err != nil {
			rows.Close()
			return stats, err
		}
		stats.ByLifecycleStage[ls] = count
	}
	rows.Close()

	return stats, nil
}

// Publish changes article status to published
func (r *KBArticleRepo) Publish(ctx context.Context, tenantID, id, userID, userName string) error {
	now := time.Now().UTC()
	result, err := r.pool.Exec(ctx, `
		UPDATE kb_articles SET
			status = 'published',
			published_at = $3,
			updated_by_id = $4,
			updated_by_name = $5,
			updated_at = $3
		WHERE tenant_id = $1 AND id = $2 AND status = 'draft'
	`, tenantID, id, now, userID, userName)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errors.New("not found or not a draft")
	}
	return nil
}

// SlugExists checks if a slug already exists for a tenant (optionally excluding an article ID)
func (r *KBArticleRepo) SlugExists(ctx context.Context, tenantID, slug, excludeID string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM kb_articles WHERE tenant_id = $1 AND slug = $2`
	args := []any{tenantID, slug}
	if excludeID != "" {
		query += ` AND id != $3`
		args = append(args, excludeID)
	}
	query += `)`
	err := r.pool.QueryRow(ctx, query, args...).Scan(&exists)
	return exists, err
}
