package mocks

import (
	"context"
)

// MockAuditLogger is a no-op implementation of audit.AuditLogger for testing.
type MockAuditLogger struct {
	CreateCalls []AuditCall
	UpdateCalls []AuditCall
	DeleteCalls []AuditCall
}

// AuditCall records an audit method call
type AuditCall struct {
	EntityType string
	EntityID   string
	Before     any
	After      any
}

// NewMockAuditLogger creates a new mock audit logger.
func NewMockAuditLogger() *MockAuditLogger {
	return &MockAuditLogger{
		CreateCalls: make([]AuditCall, 0),
		UpdateCalls: make([]AuditCall, 0),
		DeleteCalls: make([]AuditCall, 0),
	}
}

// LogCreate implements audit.AuditLogger
func (m *MockAuditLogger) LogCreate(ctx context.Context, entityType, entityID string, after any) error {
	m.CreateCalls = append(m.CreateCalls, AuditCall{
		EntityType: entityType,
		EntityID:   entityID,
		After:      after,
	})
	return nil
}

// LogUpdate implements audit.AuditLogger
func (m *MockAuditLogger) LogUpdate(ctx context.Context, entityType, entityID string, before, after any) error {
	m.UpdateCalls = append(m.UpdateCalls, AuditCall{
		EntityType: entityType,
		EntityID:   entityID,
		Before:     before,
		After:      after,
	})
	return nil
}

// LogDelete implements audit.AuditLogger
func (m *MockAuditLogger) LogDelete(ctx context.Context, entityType, entityID string, before any) error {
	m.DeleteCalls = append(m.DeleteCalls, AuditCall{
		EntityType: entityType,
		EntityID:   entityID,
		Before:     before,
	})
	return nil
}
