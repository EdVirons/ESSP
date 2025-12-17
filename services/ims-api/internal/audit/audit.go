package audit

import (
	"context"
	"encoding/json"
	"time"

	"github.com/edvirons/ssp/ims/internal/store"
)

// Action represents the type of audit action
type Action string

const (
	ActionCreate Action = "create"
	ActionUpdate Action = "update"
	ActionDelete Action = "delete"
)

// AuditContext contains contextual information for audit logging
type AuditContext struct {
	TenantID  string
	UserID    string
	UserEmail string
	IPAddress string
	UserAgent string
	RequestID string
	// Impersonation fields
	ImpersonatedUserID    string
	ImpersonatedUserEmail string
	ImpersonationReason   string
}

// AuditLogger provides methods for logging entity changes
type AuditLogger interface {
	LogCreate(ctx context.Context, entityType, entityID string, after any) error
	LogUpdate(ctx context.Context, entityType, entityID string, before, after any) error
	LogDelete(ctx context.Context, entityType, entityID string, before any) error
}

// Logger is the concrete implementation of AuditLogger
type Logger struct {
	store *Store
}

// NewLogger creates a new audit logger
func NewLogger(store *Store) *Logger {
	return &Logger{store: store}
}

// LogCreate logs a create operation
func (l *Logger) LogCreate(ctx context.Context, entityType, entityID string, after any) error {
	auditCtx := GetAuditContext(ctx)

	afterJSON, err := marshalEntity(after)
	if err != nil {
		return err
	}

	log := AuditLog{
		ID:                    store.NewID("audit"),
		TenantID:              auditCtx.TenantID,
		UserID:                auditCtx.UserID,
		UserEmail:             auditCtx.UserEmail,
		Action:                string(ActionCreate),
		EntityType:            entityType,
		EntityID:              entityID,
		BeforeState:           nil,
		AfterState:            afterJSON,
		IPAddress:             auditCtx.IPAddress,
		UserAgent:             auditCtx.UserAgent,
		RequestID:             auditCtx.RequestID,
		CreatedAt:             time.Now().UTC(),
		ImpersonatedUserID:    auditCtx.ImpersonatedUserID,
		ImpersonatedUserEmail: auditCtx.ImpersonatedUserEmail,
		ImpersonationReason:   auditCtx.ImpersonationReason,
	}

	return l.store.Create(ctx, log)
}

// LogUpdate logs an update operation
func (l *Logger) LogUpdate(ctx context.Context, entityType, entityID string, before, after any) error {
	auditCtx := GetAuditContext(ctx)

	beforeJSON, err := marshalEntity(before)
	if err != nil {
		return err
	}

	afterJSON, err := marshalEntity(after)
	if err != nil {
		return err
	}

	log := AuditLog{
		ID:                    store.NewID("audit"),
		TenantID:              auditCtx.TenantID,
		UserID:                auditCtx.UserID,
		UserEmail:             auditCtx.UserEmail,
		Action:                string(ActionUpdate),
		EntityType:            entityType,
		EntityID:              entityID,
		BeforeState:           beforeJSON,
		AfterState:            afterJSON,
		IPAddress:             auditCtx.IPAddress,
		UserAgent:             auditCtx.UserAgent,
		RequestID:             auditCtx.RequestID,
		CreatedAt:             time.Now().UTC(),
		ImpersonatedUserID:    auditCtx.ImpersonatedUserID,
		ImpersonatedUserEmail: auditCtx.ImpersonatedUserEmail,
		ImpersonationReason:   auditCtx.ImpersonationReason,
	}

	return l.store.Create(ctx, log)
}

// LogDelete logs a delete operation
func (l *Logger) LogDelete(ctx context.Context, entityType, entityID string, before any) error {
	auditCtx := GetAuditContext(ctx)

	beforeJSON, err := marshalEntity(before)
	if err != nil {
		return err
	}

	log := AuditLog{
		ID:                    store.NewID("audit"),
		TenantID:              auditCtx.TenantID,
		UserID:                auditCtx.UserID,
		UserEmail:             auditCtx.UserEmail,
		Action:                string(ActionDelete),
		EntityType:            entityType,
		EntityID:              entityID,
		BeforeState:           beforeJSON,
		AfterState:            nil,
		IPAddress:             auditCtx.IPAddress,
		UserAgent:             auditCtx.UserAgent,
		RequestID:             auditCtx.RequestID,
		CreatedAt:             time.Now().UTC(),
		ImpersonatedUserID:    auditCtx.ImpersonatedUserID,
		ImpersonatedUserEmail: auditCtx.ImpersonatedUserEmail,
		ImpersonationReason:   auditCtx.ImpersonationReason,
	}

	return l.store.Create(ctx, log)
}

// marshalEntity converts an entity to JSON bytes
func marshalEntity(entity any) ([]byte, error) {
	if entity == nil {
		return nil, nil
	}
	return json.Marshal(entity)
}

// AuditLog represents a single audit log entry
type AuditLog struct {
	ID          string
	TenantID    string
	UserID      string
	UserEmail   string
	Action      string
	EntityType  string
	EntityID    string
	BeforeState []byte
	AfterState  []byte
	IPAddress   string
	UserAgent   string
	RequestID   string
	CreatedAt   time.Time
	// Impersonation fields
	ImpersonatedUserID    string
	ImpersonatedUserEmail string
	ImpersonationReason   string
}
