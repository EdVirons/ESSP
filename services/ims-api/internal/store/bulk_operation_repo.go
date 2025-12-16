package store

import (
	"context"
	"encoding/json"
	"time"

	"github.com/edvirons/ssp/ims/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

// BulkOperationRepo handles bulk operation log operations.
type BulkOperationRepo struct {
	pool *pgxpool.Pool
}

// Create inserts a new bulk operation log entry.
func (r *BulkOperationRepo) Create(ctx context.Context, log models.BulkOperationLog) error {
	errorsJSON, _ := json.Marshal([]models.BulkOperationError{})
	if log.Errors != nil {
		errorsJSON = log.Errors
	}

	_, err := r.pool.Exec(ctx, `
		INSERT INTO bulk_operation_log (
			id, tenant_id, user_id, operation_type, entity_type, requested_ids,
			successful_ids, failed_ids, errors, started_at, completed_at,
			total_count, success_count, failure_count, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	`, log.ID, log.TenantID, log.UserID, log.OperationType, log.EntityType, log.RequestedIDs,
		log.SuccessfulIDs, log.FailedIDs, errorsJSON, log.StartedAt, log.CompletedAt,
		log.TotalCount, log.SuccessCount, log.FailureCount, log.CreatedAt)
	return err
}

// Update updates a bulk operation log with results.
func (r *BulkOperationRepo) Update(ctx context.Context, id string, successfulIDs, failedIDs []string, errors []models.BulkOperationError, completedAt time.Time) error {
	errorsJSON, _ := json.Marshal(errors)
	_, err := r.pool.Exec(ctx, `
		UPDATE bulk_operation_log
		SET successful_ids = $2,
			failed_ids = $3,
			errors = $4,
			completed_at = $5,
			success_count = array_length($2, 1),
			failure_count = array_length($3, 1)
		WHERE id = $1
	`, id, successfulIDs, failedIDs, errorsJSON, completedAt)
	return err
}

// GetByID retrieves a bulk operation log by ID.
func (r *BulkOperationRepo) GetByID(ctx context.Context, tenantID, id string) (models.BulkOperationLog, error) {
	var log models.BulkOperationLog
	var completedAt *time.Time

	err := r.pool.QueryRow(ctx, `
		SELECT id, tenant_id, user_id, operation_type, entity_type, requested_ids,
			successful_ids, failed_ids, errors, started_at, completed_at,
			total_count, success_count, failure_count, created_at
		FROM bulk_operation_log
		WHERE tenant_id = $1 AND id = $2
	`, tenantID, id).Scan(&log.ID, &log.TenantID, &log.UserID, &log.OperationType, &log.EntityType, &log.RequestedIDs,
		&log.SuccessfulIDs, &log.FailedIDs, &log.Errors, &log.StartedAt, &completedAt,
		&log.TotalCount, &log.SuccessCount, &log.FailureCount, &log.CreatedAt)
	if err != nil {
		return models.BulkOperationLog{}, err
	}
	log.CompletedAt = completedAt
	return log, nil
}

// BulkOperationListParams holds parameters for listing bulk operations.
type BulkOperationListParams struct {
	TenantID      string
	UserID        string
	OperationType string
	Limit         int
	HasCursor     bool
	CursorTime    time.Time
	CursorID      string
}

// List retrieves bulk operation logs with optional filters.
func (r *BulkOperationRepo) List(ctx context.Context, p BulkOperationListParams) ([]models.BulkOperationLog, string, error) {
	conds := []string{"tenant_id = $1"}
	args := []any{p.TenantID}
	argN := 2

	if p.UserID != "" {
		conds = append(conds, "user_id = $"+itoa(argN))
		args = append(args, p.UserID)
		argN++
	}
	if p.OperationType != "" {
		conds = append(conds, "operation_type = $"+itoa(argN))
		args = append(args, p.OperationType)
		argN++
	}
	if p.HasCursor {
		conds = append(conds, "(created_at, id) < ($"+itoa(argN)+", $"+itoa(argN+1)+")")
		args = append(args, p.CursorTime, p.CursorID)
		argN += 2
	}

	limitPlus := p.Limit + 1
	args = append(args, limitPlus)

	query := `
		SELECT id, tenant_id, user_id, operation_type, entity_type, requested_ids,
			successful_ids, failed_ids, errors, started_at, completed_at,
			total_count, success_count, failure_count, created_at
		FROM bulk_operation_log
		WHERE ` + joinConds(conds) + `
		ORDER BY created_at DESC, id DESC
		LIMIT $` + itoa(argN)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	out := []models.BulkOperationLog{}
	for rows.Next() {
		var log models.BulkOperationLog
		var completedAt *time.Time
		if err := rows.Scan(&log.ID, &log.TenantID, &log.UserID, &log.OperationType, &log.EntityType, &log.RequestedIDs,
			&log.SuccessfulIDs, &log.FailedIDs, &log.Errors, &log.StartedAt, &completedAt,
			&log.TotalCount, &log.SuccessCount, &log.FailureCount, &log.CreatedAt); err != nil {
			return nil, "", err
		}
		log.CompletedAt = completedAt
		out = append(out, log)
	}

	next := ""
	if len(out) > p.Limit {
		last := out[p.Limit-1]
		next = EncodeCursor(last.CreatedAt, last.ID)
		out = out[:p.Limit]
	}
	return out, next, nil
}

func joinConds(conds []string) string {
	if len(conds) == 0 {
		return "TRUE"
	}
	result := conds[0]
	for i := 1; i < len(conds); i++ {
		result += " AND " + conds[i]
	}
	return result
}
