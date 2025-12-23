package claude

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// ContextBuilder builds SSOT context for AI conversations
type ContextBuilder struct {
	pool *pgxpool.Pool
	log  *zap.Logger
}

// NewContextBuilder creates a new context builder
func NewContextBuilder(pool *pgxpool.Pool, log *zap.Logger) *ContextBuilder {
	return &ContextBuilder{
		pool: pool,
		log:  log,
	}
}

// BuildContext fetches relevant context for a chat session
func (cb *ContextBuilder) BuildContext(ctx context.Context, tenantID, schoolID string, deviceSerial *string) (*SSOTContext, error) {
	result := &SSOTContext{}

	// Fetch school context
	school, err := cb.fetchSchoolContext(ctx, tenantID, schoolID)
	if err != nil {
		cb.log.Warn("Failed to fetch school context", zap.Error(err), zap.String("school_id", schoolID))
	} else {
		result.School = school
	}

	// Fetch device context if serial provided
	if deviceSerial != nil && *deviceSerial != "" {
		device, err := cb.fetchDeviceContext(ctx, tenantID, *deviceSerial)
		if err != nil {
			cb.log.Warn("Failed to fetch device context", zap.Error(err), zap.String("serial", *deviceSerial))
		} else {
			result.Device = device
		}
	}

	// Fetch history context
	history, err := cb.fetchHistoryContext(ctx, tenantID, schoolID)
	if err != nil {
		cb.log.Warn("Failed to fetch history context", zap.Error(err))
	} else {
		result.History = history
	}

	return result, nil
}

// fetchSchoolContext retrieves school information from the database
func (cb *ContextBuilder) fetchSchoolContext(ctx context.Context, tenantID, schoolID string) (*SchoolContext, error) {
	query := `
		SELECT school_id, county_name, lea_name, school_type
		FROM schools
		WHERE tenant_id = $1 AND school_id = $2
		LIMIT 1
	`

	var school SchoolContext
	var schoolType *string
	err := cb.pool.QueryRow(ctx, query, tenantID, schoolID).Scan(
		&school.ID,
		&school.CountyName,
		&school.Name,
		&schoolType,
	)
	if err != nil {
		return nil, err
	}

	school.District = school.CountyName
	if schoolType != nil {
		school.Type = *schoolType
	}

	return &school, nil
}

// fetchDeviceContext retrieves device information from the database
func (cb *ContextBuilder) fetchDeviceContext(ctx context.Context, tenantID, serial string) (*DeviceContext, error) {
	query := `
		SELECT
			device_id,
			serial_number,
			COALESCE(make, ''),
			COALESCE(model, ''),
			COALESCE(device_type, ''),
			COALESCE(warranty_status, 'unknown'),
			warranty_expiry_date,
			assigned_to
		FROM devices
		WHERE tenant_id = $1 AND serial_number = $2
		LIMIT 1
	`

	var device DeviceContext
	var warrantyExpiry, assignedTo *string
	err := cb.pool.QueryRow(ctx, query, tenantID, serial).Scan(
		&device.ID,
		&device.SerialNumber,
		&device.Make,
		&device.Model,
		&device.DeviceType,
		&device.WarrantyStatus,
		&warrantyExpiry,
		&assignedTo,
	)
	if err != nil {
		return nil, err
	}

	if warrantyExpiry != nil {
		device.WarrantyExpiry = *warrantyExpiry
	}
	if assignedTo != nil {
		device.AssignedTo = *assignedTo
	}

	// Fetch last repair date
	repairQuery := `
		SELECT MAX(completed_at)
		FROM work_orders
		WHERE tenant_id = $1 AND device_serial = $2 AND status = 'completed'
	`
	var lastRepair *time.Time
	if err := cb.pool.QueryRow(ctx, repairQuery, tenantID, serial).Scan(&lastRepair); err == nil && lastRepair != nil {
		device.LastRepairDate = lastRepair.Format("2006-01-02")
	}

	return &device, nil
}

// fetchHistoryContext retrieves support history for a school
func (cb *ContextBuilder) fetchHistoryContext(ctx context.Context, tenantID, schoolID string) (*HistoryContext, error) {
	history := &HistoryContext{}

	// Count recent incidents (last 30 days)
	incidentQuery := `
		SELECT COUNT(*), MAX(created_at)
		FROM incidents
		WHERE tenant_id = $1 AND school_id = $2 AND created_at > $3
	`
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	var lastIncident *time.Time
	err := cb.pool.QueryRow(ctx, incidentQuery, tenantID, schoolID, thirtyDaysAgo).Scan(
		&history.RecentIncidents,
		&lastIncident,
	)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}
	if lastIncident != nil {
		history.LastIncidentDate = lastIncident.Format("2006-01-02")
	}

	// Count total repairs
	repairQuery := `
		SELECT COUNT(*)
		FROM work_orders
		WHERE tenant_id = $1 AND school_id = $2 AND status = 'completed'
	`
	_ = cb.pool.QueryRow(ctx, repairQuery, tenantID, schoolID).Scan(&history.TotalRepairs)

	// Get common issue categories
	categoryQuery := `
		SELECT category, COUNT(*) as cnt
		FROM incidents
		WHERE tenant_id = $1 AND school_id = $2 AND created_at > $3
		GROUP BY category
		ORDER BY cnt DESC
		LIMIT 3
	`
	rows, err := cb.pool.Query(ctx, categoryQuery, tenantID, schoolID, thirtyDaysAgo)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var category string
			var count int
			if rows.Scan(&category, &count) == nil && category != "" {
				history.CommonIssues = append(history.CommonIssues, category)
			}
		}
	}

	return history, nil
}

// UpdateDeviceFromConversation updates the context with a device mentioned in conversation
func (cb *ContextBuilder) UpdateDeviceFromConversation(ctx context.Context, ssotCtx *SSOTContext, tenantID, serial string) error {
	device, err := cb.fetchDeviceContext(ctx, tenantID, serial)
	if err != nil {
		return err
	}
	ssotCtx.Device = device
	return nil
}
