package models

import (
	"encoding/json"
	"time"
)

// FeatureKey represents a feature flag key.
type FeatureKey string

const (
	FeatureWorkOrderNotifications  FeatureKey = "work_order_notifications"
	FeatureWorkOrderRework         FeatureKey = "work_order_rework"
	FeatureWorkOrderBulkOperations FeatureKey = "work_order_bulk_operations"
	FeatureWorkOrderUpdate         FeatureKey = "work_order_update"
)

// FeatureConfig represents a feature flag configuration.
type FeatureConfig struct {
	ID          string          `json:"id"`
	TenantID    string          `json:"tenantId"`
	FeatureKey  FeatureKey      `json:"featureKey"`
	Enabled     bool            `json:"enabled"`
	ConfigValue json.RawMessage `json:"configValue"`
	Description string          `json:"description"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
}

// WorkOrderReworkConfig holds configuration for rework feature.
type WorkOrderReworkConfig struct {
	MaxReworkCount int  `json:"max_rework_count"`
	RequireReason  bool `json:"require_reason"`
}

// DefaultReworkConfig returns default rework configuration.
func DefaultReworkConfig() WorkOrderReworkConfig {
	return WorkOrderReworkConfig{
		MaxReworkCount: 5,
		RequireReason:  true,
	}
}

// ParseReworkConfig parses rework configuration from JSON.
func ParseReworkConfig(data json.RawMessage) WorkOrderReworkConfig {
	cfg := DefaultReworkConfig()
	if len(data) > 0 {
		_ = json.Unmarshal(data, &cfg)
	}
	return cfg
}

// WorkOrderNotificationConfig holds configuration for notifications.
type WorkOrderNotificationConfig struct {
	BatchDelaySeconds int `json:"batch_delay_seconds"`
}

// DefaultNotificationConfig returns default notification configuration.
func DefaultNotificationConfig() WorkOrderNotificationConfig {
	return WorkOrderNotificationConfig{
		BatchDelaySeconds: 5,
	}
}

// ParseNotificationConfig parses notification configuration from JSON.
func ParseNotificationConfig(data json.RawMessage) WorkOrderNotificationConfig {
	cfg := DefaultNotificationConfig()
	if len(data) > 0 {
		_ = json.Unmarshal(data, &cfg)
	}
	return cfg
}

// ParseBulkConfig parses bulk operation configuration from JSON.
func ParseBulkConfig(data json.RawMessage) BulkOperationConfig {
	cfg := DefaultBulkConfig()
	if len(data) > 0 {
		_ = json.Unmarshal(data, &cfg)
	}
	return cfg
}
