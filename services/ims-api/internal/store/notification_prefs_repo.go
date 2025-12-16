package store

import (
	"context"
	"errors"
	"time"

	"github.com/edvirons/ssp/ims/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lib/pq"
)

// NotificationPrefsRepo handles user notification preferences operations.
type NotificationPrefsRepo struct {
	pool *pgxpool.Pool
}

// GetPreferences retrieves notification preferences for a user.
func (r *NotificationPrefsRepo) GetPreferences(ctx context.Context, tenantID, userID string) (models.UserNotificationPreferences, error) {
	var prefs models.UserNotificationPreferences
	var enabledTypes []string

	err := r.pool.QueryRow(ctx, `
		SELECT id, tenant_id, user_id, enabled_types, in_app_enabled, email_enabled,
			quiet_hours_start, quiet_hours_end, quiet_hours_timezone, created_at, updated_at
		FROM user_notification_preferences
		WHERE tenant_id = $1 AND user_id = $2
	`, tenantID, userID).Scan(&prefs.ID, &prefs.TenantID, &prefs.UserID, pq.Array(&enabledTypes),
		&prefs.InAppEnabled, &prefs.EmailEnabled, &prefs.QuietHoursStart, &prefs.QuietHoursEnd,
		&prefs.QuietHoursTimezone, &prefs.CreatedAt, &prefs.UpdatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		// Return default preferences
		return models.UserNotificationPreferences{
			TenantID:           tenantID,
			UserID:             userID,
			EnabledTypes:       models.AllNotificationTypes,
			InAppEnabled:       true,
			EmailEnabled:       false,
			QuietHoursTimezone: "Africa/Nairobi",
		}, nil
	}
	if err != nil {
		return models.UserNotificationPreferences{}, err
	}

	// Convert string slice to NotificationType slice
	prefs.EnabledTypes = make([]models.NotificationType, len(enabledTypes))
	for i, t := range enabledTypes {
		prefs.EnabledTypes[i] = models.NotificationType(t)
	}
	return prefs, nil
}

// UpsertPreferences creates or updates notification preferences for a user.
func (r *NotificationPrefsRepo) UpsertPreferences(ctx context.Context, prefs models.UserNotificationPreferences) error {
	// Convert NotificationType slice to string slice for pq.Array
	enabledTypes := make([]string, len(prefs.EnabledTypes))
	for i, t := range prefs.EnabledTypes {
		enabledTypes[i] = string(t)
	}

	now := time.Now().UTC()
	_, err := r.pool.Exec(ctx, `
		INSERT INTO user_notification_preferences (
			id, tenant_id, user_id, enabled_types, in_app_enabled, email_enabled,
			quiet_hours_start, quiet_hours_end, quiet_hours_timezone, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (tenant_id, user_id) DO UPDATE SET
			enabled_types = $4,
			in_app_enabled = $5,
			email_enabled = $6,
			quiet_hours_start = $7,
			quiet_hours_end = $8,
			quiet_hours_timezone = $9,
			updated_at = $11
	`, prefs.ID, prefs.TenantID, prefs.UserID, pq.Array(enabledTypes),
		prefs.InAppEnabled, prefs.EmailEnabled, prefs.QuietHoursStart, prefs.QuietHoursEnd,
		prefs.QuietHoursTimezone, now, now)
	return err
}

// DeletePreferences removes notification preferences for a user.
func (r *NotificationPrefsRepo) DeletePreferences(ctx context.Context, tenantID, userID string) error {
	_, err := r.pool.Exec(ctx, `
		DELETE FROM user_notification_preferences
		WHERE tenant_id = $1 AND user_id = $2
	`, tenantID, userID)
	return err
}
