package testutil

import (
	"context"
	"testing"
	"time"

	"github.com/edvirons/ssp/ims/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

// FixtureConfig holds configuration for creating test fixtures.
type FixtureConfig struct {
	TenantID string
	SchoolID string
}

// DefaultFixtureConfig returns a default fixture configuration.
func DefaultFixtureConfig() FixtureConfig {
	return FixtureConfig{
		TenantID: "test-tenant",
		SchoolID: "test-school",
	}
}

// CreateSchoolSnapshot creates a test school snapshot in the database.
func CreateSchoolSnapshot(t *testing.T, pool *pgxpool.Pool, cfg FixtureConfig) models.SchoolSnapshot {
	t.Helper()

	snapshot := models.SchoolSnapshot{
		TenantID:      cfg.TenantID,
		SchoolID:      cfg.SchoolID,
		Name:          "Test School",
		CountyCode:    "001",
		CountyName:    "Test County",
		SubCountyCode: "001-01",
		SubCountyName: "Test SubCounty",
		UpdatedAt:     time.Now(),
	}

	query := `
		INSERT INTO schools_snapshots
		(tenant_id, school_id, name, county_code, county_name, sub_county_code, sub_county_name, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (tenant_id, school_id) DO UPDATE
		SET name = EXCLUDED.name, updated_at = EXCLUDED.updated_at
	`

	_, err := pool.Exec(context.Background(), query,
		snapshot.TenantID, snapshot.SchoolID, snapshot.Name,
		snapshot.CountyCode, snapshot.CountyName,
		snapshot.SubCountyCode, snapshot.SubCountyName,
		snapshot.UpdatedAt,
	)
	require.NoError(t, err, "failed to create school snapshot")

	return snapshot
}

// CreateDeviceSnapshot creates a test device snapshot in the database.
func CreateDeviceSnapshot(t *testing.T, pool *pgxpool.Pool, cfg FixtureConfig, deviceID string) models.DeviceSnapshot {
	t.Helper()

	snapshot := models.DeviceSnapshot{
		TenantID:  cfg.TenantID,
		DeviceID:  deviceID,
		SchoolID:  cfg.SchoolID,
		Model:     "Test Model",
		Serial:    "TEST-SERIAL-001",
		AssetTag:  "ASSET-001",
		Status:    "active",
		UpdatedAt: time.Now(),
	}

	query := `
		INSERT INTO devices_snapshots
		(tenant_id, device_id, school_id, model, serial, asset_tag, status, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (tenant_id, device_id) DO UPDATE
		SET model = EXCLUDED.model, updated_at = EXCLUDED.updated_at
	`

	_, err := pool.Exec(context.Background(), query,
		snapshot.TenantID, snapshot.DeviceID, snapshot.SchoolID,
		snapshot.Model, snapshot.Serial, snapshot.AssetTag,
		snapshot.Status, snapshot.UpdatedAt,
	)
	require.NoError(t, err, "failed to create device snapshot")

	return snapshot
}

// CreatePartSnapshot creates a test part snapshot in the database.
func CreatePartSnapshot(t *testing.T, pool *pgxpool.Pool, cfg FixtureConfig, partID string) models.PartSnapshot {
	t.Helper()

	snapshot := models.PartSnapshot{
		TenantID:  cfg.TenantID,
		PartID:    partID,
		PUK:       "TEST-PUK-001",
		Name:      "Test Part",
		Category:  "electronics",
		Unit:      "piece",
		UpdatedAt: time.Now(),
	}

	query := `
		INSERT INTO parts_snapshots
		(tenant_id, part_id, puk, name, category, unit, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (tenant_id, part_id) DO UPDATE
		SET name = EXCLUDED.name, updated_at = EXCLUDED.updated_at
	`

	_, err := pool.Exec(context.Background(), query,
		snapshot.TenantID, snapshot.PartID, snapshot.PUK,
		snapshot.Name, snapshot.Category, snapshot.Unit,
		snapshot.UpdatedAt,
	)
	require.NoError(t, err, "failed to create part snapshot")

	return snapshot
}

// CreateIncident creates a test incident in the database.
func CreateIncident(t *testing.T, pool *pgxpool.Pool, cfg FixtureConfig, deviceID string) models.Incident {
	t.Helper()

	incident := models.Incident{
		ID:          GenerateTestID(t),
		TenantID:    cfg.TenantID,
		SchoolID:    cfg.SchoolID,
		DeviceID:    deviceID,
		SchoolName:  "Test School",
		Category:    "hardware",
		Severity:    models.SeverityMedium,
		Status:      models.IncidentNew,
		Title:       "Test Incident",
		Description: "This is a test incident",
		ReportedBy:  "test-user",
		SLADueAt:    time.Now().Add(24 * time.Hour),
		SLABreached: false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	query := `
		INSERT INTO incidents
		(id, tenant_id, school_id, device_id, school_name, category, severity, status,
		 title, description, reported_by, sla_due_at, sla_breached, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	`

	_, err := pool.Exec(context.Background(), query,
		incident.ID, incident.TenantID, incident.SchoolID, incident.DeviceID,
		incident.SchoolName, incident.Category, incident.Severity, incident.Status,
		incident.Title, incident.Description, incident.ReportedBy,
		incident.SLADueAt, incident.SLABreached,
		incident.CreatedAt, incident.UpdatedAt,
	)
	require.NoError(t, err, "failed to create incident")

	return incident
}

// CreateWorkOrder creates a test work order in the database.
func CreateWorkOrder(t *testing.T, pool *pgxpool.Pool, cfg FixtureConfig, incidentID, deviceID string) models.WorkOrder {
	t.Helper()

	workOrder := models.WorkOrder{
		ID:             GenerateTestID(t),
		IncidentID:     incidentID,
		TenantID:       cfg.TenantID,
		SchoolID:       cfg.SchoolID,
		DeviceID:       deviceID,
		SchoolName:     "Test School",
		Status:         models.WorkOrderDraft,
		RepairLocation: models.RepairLocationServiceShop,
		TaskType:       "repair",
		Notes:          "Test work order",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	query := `
		INSERT INTO work_orders
		(id, incident_id, tenant_id, school_id, device_id, school_name, status,
		 repair_location, task_type, notes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	_, err := pool.Exec(context.Background(), query,
		workOrder.ID, workOrder.IncidentID, workOrder.TenantID,
		workOrder.SchoolID, workOrder.DeviceID, workOrder.SchoolName,
		workOrder.Status, workOrder.RepairLocation, workOrder.TaskType,
		workOrder.Notes, workOrder.CreatedAt, workOrder.UpdatedAt,
	)
	require.NoError(t, err, "failed to create work order")

	return workOrder
}

// CreateServiceShop creates a test service shop in the database.
func CreateServiceShop(t *testing.T, pool *pgxpool.Pool, cfg FixtureConfig) models.ServiceShop {
	t.Helper()

	shop := models.ServiceShop{
		ID:            GenerateTestID(t),
		TenantID:      cfg.TenantID,
		CountyCode:    "001",
		CountyName:    "Test County",
		SubCountyCode: "001-01",
		SubCountyName: "Test SubCounty",
		CoverageLevel: "county",
		Name:          "Test Service Shop",
		Location:      "Test Location",
		Active:        true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	query := `
		INSERT INTO service_shops
		(id, tenant_id, county_code, county_name, sub_county_code, sub_county_name,
		 coverage_level, name, location, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	_, err := pool.Exec(context.Background(), query,
		shop.ID, shop.TenantID, shop.CountyCode, shop.CountyName,
		shop.SubCountyCode, shop.SubCountyName, shop.CoverageLevel,
		shop.Name, shop.Location, shop.Active,
		shop.CreatedAt, shop.UpdatedAt,
	)
	require.NoError(t, err, "failed to create service shop")

	return shop
}

// CreateServiceStaff creates a test service staff member in the database.
func CreateServiceStaff(t *testing.T, pool *pgxpool.Pool, cfg FixtureConfig, shopID string) models.ServiceStaff {
	t.Helper()

	staff := models.ServiceStaff{
		ID:            GenerateTestID(t),
		TenantID:      cfg.TenantID,
		ServiceShopID: shopID,
		UserID:        "test-user-id",
		Role:          models.StaffRoleLeadTechnician,
		Phone:         "+254700000000",
		Active:        true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	query := `
		INSERT INTO service_staff
		(id, tenant_id, service_shop_id, user_id, role, phone, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := pool.Exec(context.Background(), query,
		staff.ID, staff.TenantID, staff.ServiceShopID,
		staff.UserID, staff.Role, staff.Phone, staff.Active,
		staff.CreatedAt, staff.UpdatedAt,
	)
	require.NoError(t, err, "failed to create service staff")

	return staff
}

// CreatePart creates a test part in the database.
func CreatePart(t *testing.T, pool *pgxpool.Pool, cfg FixtureConfig) models.Part {
	t.Helper()

	part := models.Part{
		ID:        GenerateTestID(t),
		TenantID:  cfg.TenantID,
		SKU:       "TEST-SKU-001",
		Name:      "Test Part",
		Category:  "electronics",
		CreatedAt: time.Now(),
	}

	query := `
		INSERT INTO parts
		(id, tenant_id, sku, name, category, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := pool.Exec(context.Background(), query,
		part.ID, part.TenantID, part.SKU,
		part.Name, part.Category, part.CreatedAt,
	)
	require.NoError(t, err, "failed to create part")

	return part
}

// CreateInventoryItem creates a test inventory item in the database.
func CreateInventoryItem(t *testing.T, pool *pgxpool.Pool, cfg FixtureConfig, shopID, partID string, qty int64) models.InventoryItem {
	t.Helper()

	item := models.InventoryItem{
		ID:               GenerateTestID(t),
		TenantID:         cfg.TenantID,
		ServiceShopID:    shopID,
		PartID:           partID,
		QtyAvailable:     qty,
		QtyReserved:      0,
		ReorderThreshold: 10,
		UpdatedAt:        time.Now(),
	}

	query := `
		INSERT INTO inventory
		(id, tenant_id, service_shop_id, part_id, qty_available, qty_reserved, reorder_threshold, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := pool.Exec(context.Background(), query,
		item.ID, item.TenantID, item.ServiceShopID, item.PartID,
		item.QtyAvailable, item.QtyReserved, item.ReorderThreshold,
		item.UpdatedAt,
	)
	require.NoError(t, err, "failed to create inventory item")

	return item
}

// CreateAttachment creates a test attachment in the database.
func CreateAttachment(t *testing.T, pool *pgxpool.Pool, cfg FixtureConfig, entityType models.AttachmentEntityType, entityID string) models.Attachment {
	t.Helper()

	attachment := models.Attachment{
		ID:          GenerateTestID(t),
		TenantID:    cfg.TenantID,
		SchoolID:    cfg.SchoolID,
		EntityType:  entityType,
		EntityID:    entityID,
		FileName:    "test-file.jpg",
		ContentType: "image/jpeg",
		SizeBytes:   1024,
		ObjectKey:   "test/test-file.jpg",
		CreatedAt:   time.Now(),
	}

	query := `
		INSERT INTO attachments
		(id, tenant_id, school_id, entity_type, entity_id, file_name, content_type, size_bytes, object_key, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := pool.Exec(context.Background(), query,
		attachment.ID, attachment.TenantID, attachment.SchoolID,
		attachment.EntityType, attachment.EntityID,
		attachment.FileName, attachment.ContentType, attachment.SizeBytes,
		attachment.ObjectKey, attachment.CreatedAt,
	)
	require.NoError(t, err, "failed to create attachment")

	return attachment
}
