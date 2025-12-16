package store

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/edvirons/ssp/ims/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// FeatureConfigRepo handles feature configuration operations.
type FeatureConfigRepo struct {
	pool *pgxpool.Pool
}

// GetFeature retrieves a feature configuration by key.
// First checks for tenant-specific config, then falls back to global config (tenant_id = '*').
func (r *FeatureConfigRepo) GetFeature(ctx context.Context, tenantID string, key models.FeatureKey) (models.FeatureConfig, error) {
	var cfg models.FeatureConfig

	// First try tenant-specific config
	err := r.pool.QueryRow(ctx, `
		SELECT id, tenant_id, feature_key, enabled, config_value, description, created_at, updated_at
		FROM feature_config
		WHERE tenant_id = $1 AND feature_key = $2
	`, tenantID, key).Scan(&cfg.ID, &cfg.TenantID, &cfg.FeatureKey, &cfg.Enabled, &cfg.ConfigValue, &cfg.Description, &cfg.CreatedAt, &cfg.UpdatedAt)

	if err == nil {
		return cfg, nil
	}

	if !errors.Is(err, pgx.ErrNoRows) {
		return models.FeatureConfig{}, err
	}

	// Fall back to global config
	err = r.pool.QueryRow(ctx, `
		SELECT id, tenant_id, feature_key, enabled, config_value, description, created_at, updated_at
		FROM feature_config
		WHERE tenant_id = '*' AND feature_key = $1
	`, key).Scan(&cfg.ID, &cfg.TenantID, &cfg.FeatureKey, &cfg.Enabled, &cfg.ConfigValue, &cfg.Description, &cfg.CreatedAt, &cfg.UpdatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return models.FeatureConfig{}, models.ErrFeatureNotFound
	}
	return cfg, err
}

// IsFeatureEnabled checks if a feature is enabled for a tenant.
func (r *FeatureConfigRepo) IsFeatureEnabled(ctx context.Context, tenantID string, key models.FeatureKey) (bool, error) {
	cfg, err := r.GetFeature(ctx, tenantID, key)
	if err != nil {
		if errors.Is(err, models.ErrFeatureNotFound) {
			// Default to enabled if config not found
			return true, nil
		}
		return false, err
	}
	return cfg.Enabled, nil
}

// GetBulkConfig returns the bulk operation configuration for a tenant.
func (r *FeatureConfigRepo) GetBulkConfig(ctx context.Context, tenantID string) (models.BulkOperationConfig, error) {
	cfg, err := r.GetFeature(ctx, tenantID, models.FeatureWorkOrderBulkOperations)
	if err != nil {
		if errors.Is(err, models.ErrFeatureNotFound) {
			return models.DefaultBulkConfig(), nil
		}
		return models.BulkOperationConfig{}, err
	}
	return models.ParseBulkConfig(cfg.ConfigValue), nil
}

// GetReworkConfig returns the rework configuration for a tenant.
func (r *FeatureConfigRepo) GetReworkConfig(ctx context.Context, tenantID string) (models.WorkOrderReworkConfig, error) {
	cfg, err := r.GetFeature(ctx, tenantID, models.FeatureWorkOrderRework)
	if err != nil {
		if errors.Is(err, models.ErrFeatureNotFound) {
			return models.DefaultReworkConfig(), nil
		}
		return models.WorkOrderReworkConfig{}, err
	}
	return models.ParseReworkConfig(cfg.ConfigValue), nil
}

// UpsertFeature creates or updates a feature configuration.
func (r *FeatureConfigRepo) UpsertFeature(ctx context.Context, cfg models.FeatureConfig) error {
	configJSON, err := json.Marshal(cfg.ConfigValue)
	if err != nil {
		configJSON = []byte("{}")
	}

	now := time.Now().UTC()
	_, err = r.pool.Exec(ctx, `
		INSERT INTO feature_config (id, tenant_id, feature_key, enabled, config_value, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (tenant_id, feature_key) DO UPDATE SET
			enabled = $4,
			config_value = $5,
			description = $6,
			updated_at = $8
	`, cfg.ID, cfg.TenantID, cfg.FeatureKey, cfg.Enabled, configJSON, cfg.Description, now, now)
	return err
}

// ListFeatures retrieves all feature configurations for a tenant.
func (r *FeatureConfigRepo) ListFeatures(ctx context.Context, tenantID string) ([]models.FeatureConfig, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, tenant_id, feature_key, enabled, config_value, description, created_at, updated_at
		FROM feature_config
		WHERE tenant_id = $1 OR tenant_id = '*'
		ORDER BY feature_key
	`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Use a map to handle tenant-specific overriding global
	configMap := make(map[models.FeatureKey]models.FeatureConfig)
	for rows.Next() {
		var cfg models.FeatureConfig
		if err := rows.Scan(&cfg.ID, &cfg.TenantID, &cfg.FeatureKey, &cfg.Enabled, &cfg.ConfigValue, &cfg.Description, &cfg.CreatedAt, &cfg.UpdatedAt); err != nil {
			return nil, err
		}
		// Tenant-specific config overrides global
		existing, exists := configMap[cfg.FeatureKey]
		if !exists || (exists && existing.TenantID == "*" && cfg.TenantID != "*") {
			configMap[cfg.FeatureKey] = cfg
		}
	}

	out := make([]models.FeatureConfig, 0, len(configMap))
	for _, cfg := range configMap {
		out = append(out, cfg)
	}
	return out, nil
}
