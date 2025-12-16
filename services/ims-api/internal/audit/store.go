package audit

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/edvirons/ssp/ims/internal/store"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Store handles database operations for audit logs
type Store struct {
	pool *pgxpool.Pool
}

// NewStore creates a new audit store
func NewStore(pool *pgxpool.Pool) *Store {
	return &Store{pool: pool}
}

// Create inserts a new audit log entry
func (s *Store) Create(ctx context.Context, log AuditLog) error {
	_, err := s.pool.Exec(ctx, `
		INSERT INTO audit_logs (
			id, tenant_id, user_id, user_email, action,
			entity_type, entity_id, before_state, after_state,
			ip_address, user_agent, request_id, created_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
		)
	`, log.ID, log.TenantID, log.UserID, log.UserEmail, log.Action,
		log.EntityType, log.EntityID, log.BeforeState, log.AfterState,
		log.IPAddress, log.UserAgent, log.RequestID, log.CreatedAt)
	return err
}

// ListParams defines parameters for listing audit logs
type ListParams struct {
	TenantID   string
	EntityType string
	EntityID   string
	UserID     string
	Action     string
	StartDate  time.Time
	EndDate    time.Time
	Limit      int

	HasCursor       bool
	CursorCreatedAt time.Time
	CursorID        string
}

// ListResult represents a single audit log result with decoded JSON
type ListResult struct {
	ID          string    `json:"id"`
	TenantID    string    `json:"tenantId"`
	UserID      string    `json:"userId"`
	UserEmail   string    `json:"userEmail"`
	Action      string    `json:"action"`
	EntityType  string    `json:"entityType"`
	EntityID    string    `json:"entityId"`
	BeforeState any       `json:"beforeState,omitempty"`
	AfterState  any       `json:"afterState,omitempty"`
	IPAddress   string    `json:"ipAddress"`
	UserAgent   string    `json:"userAgent"`
	RequestID   string    `json:"requestId"`
	CreatedAt   time.Time `json:"createdAt"`
}

// List retrieves audit logs with filtering and pagination
func (s *Store) List(ctx context.Context, p ListParams) ([]ListResult, string, error) {
	conds := []string{"tenant_id=$1"}
	args := []any{p.TenantID}
	argN := 2

	if p.EntityType != "" {
		conds = append(conds, "entity_type=$"+itoa(argN))
		args = append(args, p.EntityType)
		argN++
	}

	if p.EntityID != "" {
		conds = append(conds, "entity_id=$"+itoa(argN))
		args = append(args, p.EntityID)
		argN++
	}

	if p.UserID != "" {
		conds = append(conds, "user_id=$"+itoa(argN))
		args = append(args, p.UserID)
		argN++
	}

	if p.Action != "" {
		conds = append(conds, "action=$"+itoa(argN))
		args = append(args, p.Action)
		argN++
	}

	if !p.StartDate.IsZero() {
		conds = append(conds, "created_at >= $"+itoa(argN))
		args = append(args, p.StartDate)
		argN++
	}

	if !p.EndDate.IsZero() {
		conds = append(conds, "created_at <= $"+itoa(argN))
		args = append(args, p.EndDate)
		argN++
	}

	if p.HasCursor {
		conds = append(conds, "(created_at, id) < ($"+itoa(argN)+", $"+itoa(argN+1)+")")
		args = append(args, p.CursorCreatedAt, p.CursorID)
		argN += 2
	}

	// +1 fetch for nextCursor
	limitPlus := p.Limit + 1
	args = append(args, limitPlus)

	sql := `
		SELECT id, tenant_id, user_id, user_email, action,
		       entity_type, entity_id, before_state, after_state,
		       ip_address, user_agent, request_id, created_at
		FROM audit_logs
		WHERE ` + strings.Join(conds, " AND ") + `
		ORDER BY created_at DESC, id DESC
		LIMIT $` + itoa(argN)

	rows, err := s.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	out := []ListResult{}
	for rows.Next() {
		var log ListResult
		var beforeState, afterState []byte

		if err := rows.Scan(
			&log.ID, &log.TenantID, &log.UserID, &log.UserEmail, &log.Action,
			&log.EntityType, &log.EntityID, &beforeState, &afterState,
			&log.IPAddress, &log.UserAgent, &log.RequestID, &log.CreatedAt,
		); err != nil {
			return nil, "", err
		}

		// Decode JSONB fields
		if beforeState != nil {
			var before any
			if err := json.Unmarshal(beforeState, &before); err == nil {
				log.BeforeState = before
			}
		}

		if afterState != nil {
			var after any
			if err := json.Unmarshal(afterState, &after); err == nil {
				log.AfterState = after
			}
		}

		out = append(out, log)
	}

	next := ""
	if len(out) > p.Limit {
		last := out[p.Limit-1]
		next = store.EncodeCursor(last.CreatedAt, last.ID)
		out = out[:p.Limit]
	}
	return out, next, nil
}

// GetByID retrieves a single audit log by ID
func (s *Store) GetByID(ctx context.Context, tenantID, id string) (ListResult, error) {
	var log ListResult
	var beforeState, afterState []byte

	err := s.pool.QueryRow(ctx, `
		SELECT id, tenant_id, user_id, user_email, action,
		       entity_type, entity_id, before_state, after_state,
		       ip_address, user_agent, request_id, created_at
		FROM audit_logs
		WHERE tenant_id=$1 AND id=$2
	`, tenantID, id).Scan(
		&log.ID, &log.TenantID, &log.UserID, &log.UserEmail, &log.Action,
		&log.EntityType, &log.EntityID, &beforeState, &afterState,
		&log.IPAddress, &log.UserAgent, &log.RequestID, &log.CreatedAt,
	)

	if err != nil {
		return ListResult{}, err
	}

	// Decode JSONB fields
	if beforeState != nil {
		var before any
		if err := json.Unmarshal(beforeState, &before); err == nil {
			log.BeforeState = before
		}
	}

	if afterState != nil {
		var after any
		if err := json.Unmarshal(afterState, &after); err == nil {
			log.AfterState = after
		}
	}

	return log, nil
}

// itoa converts int to string for building SQL queries
func itoa(i int) string {
	if i < 0 {
		return "0"
	}
	if i == 0 {
		return "0"
	}
	neg := false
	if i < 0 {
		neg = true
		i = -i
	}
	buf := [20]byte{}
	pos := len(buf) - 1
	for i >= 10 {
		buf[pos] = byte('0' + i%10)
		pos--
		i /= 10
	}
	buf[pos] = byte('0' + i)
	if neg {
		pos--
		buf[pos] = '-'
	}
	return string(buf[pos:])
}
