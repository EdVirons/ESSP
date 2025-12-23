package store

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/edvirons/ssp/ims/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SurveysRepo struct{ pool *pgxpool.Pool }
type SurveyRoomsRepo struct{ pool *pgxpool.Pool }
type SurveyPhotosRepo struct{ pool *pgxpool.Pool }

func (r *SurveysRepo) Create(ctx context.Context, s models.SiteSurvey) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO site_surveys (
			id, tenant_id, project_id, status, conducted_by_user_id, conducted_at, summary, risks, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
	`, s.ID, s.TenantID, s.ProjectID, s.Status, s.ConductedByUserID, s.ConductedAt, s.Summary, s.Risks, s.CreatedAt, s.UpdatedAt)
	return err
}

func (r *SurveysRepo) GetByID(ctx context.Context, tenantID, id string) (models.SiteSurvey, error) {
	var s models.SiteSurvey
	row := r.pool.QueryRow(ctx, `
		SELECT id, tenant_id, project_id, status, conducted_by_user_id, conducted_at, summary, risks, created_at, updated_at
		FROM site_surveys
		WHERE tenant_id=$1 AND id=$2
	`, tenantID, id)
	if err := row.Scan(&s.ID, &s.TenantID, &s.ProjectID, &s.Status, &s.ConductedByUserID, &s.ConductedAt, &s.Summary, &s.Risks, &s.CreatedAt, &s.UpdatedAt); err != nil {
		return models.SiteSurvey{}, errors.New("not found")
	}
	return s, nil
}

type SurveyListParams struct {
	TenantID        string
	ProjectID       string
	Status          string
	Limit           int
	HasCursor       bool
	CursorCreatedAt time.Time
	CursorID        string
}

func (r *SurveysRepo) List(ctx context.Context, p SurveyListParams) ([]models.SiteSurvey, string, error) {
	conds := []string{"tenant_id=$1", "project_id=$2"}
	args := []any{p.TenantID, p.ProjectID}
	argN := 3
	if p.Status != "" {
		conds = append(conds, "status=$"+itoa(argN))
		args = append(args, p.Status)
		argN++
	}
	if p.HasCursor {
		conds = append(conds, "(created_at, id) < ($"+itoa(argN)+", $"+itoa(argN+1)+")")
		args = append(args, p.CursorCreatedAt, p.CursorID)
		argN += 2
	}
	limitPlus := p.Limit + 1
	args = append(args, limitPlus)
	sql := `
		SELECT id, tenant_id, project_id, status, conducted_by_user_id, conducted_at, summary, risks, created_at, updated_at
		FROM site_surveys
		WHERE ` + strings.Join(conds, " AND ") + `
		ORDER BY created_at DESC, id DESC
		LIMIT $` + itoa(argN)
	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()
	out := []models.SiteSurvey{}
	for rows.Next() {
		var x models.SiteSurvey
		if err := rows.Scan(&x.ID, &x.TenantID, &x.ProjectID, &x.Status, &x.ConductedByUserID, &x.ConductedAt, &x.Summary, &x.Risks, &x.CreatedAt, &x.UpdatedAt); err != nil {
			return nil, "", err
		}
		out = append(out, x)
	}
	next := ""
	if len(out) > p.Limit {
		last := out[p.Limit-1]
		next = EncodeCursor(last.CreatedAt, last.ID)
		out = out[:p.Limit]
	}
	return out, next, nil
}

func (r *SurveyRoomsRepo) Create(ctx context.Context, room models.SurveyRoom) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO survey_rooms (
			id, tenant_id, survey_id, name, room_type, floor, power_notes, network_notes, created_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
	`, room.ID, room.TenantID, room.SurveyID, room.Name, room.RoomType, room.Floor, room.PowerNotes, room.NetworkNotes, room.CreatedAt)
	return err
}

func (r *SurveyRoomsRepo) List(ctx context.Context, tenantID, surveyID string) ([]models.SurveyRoom, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, tenant_id, survey_id, name, room_type, floor, power_notes, network_notes, created_at
		FROM survey_rooms
		WHERE tenant_id=$1 AND survey_id=$2
		ORDER BY created_at ASC
	`, tenantID, surveyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []models.SurveyRoom{}
	for rows.Next() {
		var x models.SurveyRoom
		if err := rows.Scan(&x.ID, &x.TenantID, &x.SurveyID, &x.Name, &x.RoomType, &x.Floor, &x.PowerNotes, &x.NetworkNotes, &x.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, x)
	}
	return out, nil
}

func (r *SurveyPhotosRepo) Create(ctx context.Context, p models.SurveyPhoto) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO survey_photos (id, tenant_id, survey_id, room_id, attachment_id, caption, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
	`, p.ID, p.TenantID, p.SurveyID, p.RoomID, p.AttachmentID, p.Caption, p.CreatedAt)
	return err
}

func (r *SurveyPhotosRepo) List(ctx context.Context, tenantID, surveyID string) ([]models.SurveyPhoto, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, tenant_id, survey_id, room_id, attachment_id, caption, created_at
		FROM survey_photos
		WHERE tenant_id=$1 AND survey_id=$2
		ORDER BY created_at DESC
	`, tenantID, surveyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []models.SurveyPhoto{}
	for rows.Next() {
		var x models.SurveyPhoto
		if err := rows.Scan(&x.ID, &x.TenantID, &x.SurveyID, &x.RoomID, &x.AttachmentID, &x.Caption, &x.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, x)
	}
	return out, nil
}
