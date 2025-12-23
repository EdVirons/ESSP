package store

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/edvirons/ssp/ims/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

// setupTestDB creates a connection pool to the test database.
// It uses the TEST_PG_DSN environment variable or falls back to a default localhost connection.
// The caller is responsible for cleaning up test data.
func setupTestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()

	dsn := os.Getenv("TEST_PG_DSN")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/ims_test?sslmode=disable"
	}

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Verify connection
	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		t.Fatalf("Failed to ping test database: %v", err)
	}

	t.Cleanup(func() {
		pool.Close()
	})

	return pool
}

// cleanupIncidents removes all test incidents for a given tenant/school
func cleanupIncidents(t *testing.T, pool *pgxpool.Pool, tenantID, schoolID string) {
	t.Helper()
	_, err := pool.Exec(context.Background(),
		"DELETE FROM incidents WHERE tenant_id=$1 AND school_id=$2",
		tenantID, schoolID)
	if err != nil {
		t.Logf("Warning: Failed to cleanup incidents: %v", err)
	}
}

// cleanupWorkOrders removes all test work orders for a given tenant/school
func cleanupWorkOrders(t *testing.T, pool *pgxpool.Pool, tenantID, schoolID string) {
	t.Helper()
	_, err := pool.Exec(context.Background(),
		"DELETE FROM work_orders WHERE tenant_id=$1 AND school_id=$2",
		tenantID, schoolID)
	if err != nil {
		t.Logf("Warning: Failed to cleanup work orders: %v", err)
	}
}

// cleanupServiceShops removes all test service shops for a given tenant
func cleanupServiceShops(t *testing.T, pool *pgxpool.Pool, tenantID string) {
	t.Helper()
	_, err := pool.Exec(context.Background(),
		"DELETE FROM service_shops WHERE tenant_id=$1",
		tenantID)
	if err != nil {
		t.Logf("Warning: Failed to cleanup service shops: %v", err)
	}
}

// cleanupWorkOrderParts removes all test work order parts for a given tenant/school
func cleanupWorkOrderParts(t *testing.T, pool *pgxpool.Pool, tenantID, schoolID string) {
	t.Helper()
	_, err := pool.Exec(context.Background(),
		"DELETE FROM work_order_parts WHERE tenant_id=$1 AND school_id=$2",
		tenantID, schoolID)
	if err != nil {
		t.Logf("Warning: Failed to cleanup work order parts: %v", err)
	}
}

// Test data factories

func validIncident() models.Incident {
	now := time.Now().UTC()
	return models.Incident{
		ID:          "inc-test-001",
		TenantID:    "tenant-test",
		SchoolID:    "school-test",
		DeviceID:    "device-001",
		SchoolName:  "Test School",
		CountyID:    "county-001",
		CountyName:  "Test County",
		SubCountyID: "subcounty-001",
		SubCountyName: "Test SubCounty",
		ContactName: "John Doe",
		ContactPhone: "+254700000000",
		ContactEmail: "john@example.com",
		DeviceSerial: "SN12345",
		DeviceAssetTag: "AT12345",
		DeviceModelID: "model-001",
		DeviceMake: "Dell",
		DeviceModel: "Latitude 3400",
		DeviceCategory: "Laptop",
		Category:    "hardware",
		Severity:    models.SeverityMedium,
		Status:      models.IncidentNew,
		Title:       "Test Incident",
		Description: "Test incident description",
		ReportedBy:  "user-test",
		SLADueAt:    now.Add(24 * time.Hour),
		SLABreached: false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func incidentNoTenant() models.Incident {
	inc := validIncident()
	inc.ID = "inc-test-002"
	inc.TenantID = ""
	return inc
}

func validWorkOrder() models.WorkOrder {
	now := time.Now().UTC()
	return models.WorkOrder{
		ID:                "wo-test-001",
		IncidentID:        "inc-test-001",
		TenantID:          "tenant-test",
		SchoolID:          "school-test",
		DeviceID:          "device-001",
		SchoolName:        "Test School",
		ContactName:       "John Doe",
		ContactPhone:      "+254700000000",
		ContactEmail:      "john@example.com",
		DeviceSerial:      "SN12345",
		DeviceAssetTag:    "AT12345",
		DeviceModelID:     "model-001",
		DeviceMake:        "Dell",
		DeviceModel:       "Latitude 3400",
		DeviceCategory:    "Laptop",
		Status:            models.WorkOrderDraft,
		ServiceShopID:     "shop-001",
		AssignedStaffID:   "staff-001",
		RepairLocation:    models.RepairLocationServiceShop,
		AssignedTo:        "Technician Name",
		TaskType:          "repair",
		ProjectID:         "",
		PhaseID:           "",
		OnsiteContactID:   "",
		ApprovalStatus:    "",
		CostEstimateCents: 10000,
		Notes:             "Test work order notes",
		CreatedAt:         now,
		UpdatedAt:         now,
	}
}

func workOrderNoTenant() models.WorkOrder {
	wo := validWorkOrder()
	wo.ID = "wo-test-002"
	wo.TenantID = ""
	return wo
}

func validServiceShop() models.ServiceShop {
	now := time.Now().UTC()
	return models.ServiceShop{
		ID:            "shop-test-001",
		TenantID:      "tenant-test",
		CountyCode:    "001",
		CountyName:    "Nairobi",
		SubCountyCode: "001-001",
		SubCountyName: "Westlands",
		CoverageLevel: "sub_county",
		Name:          "Test Service Shop",
		Location:      "Nairobi, Kenya",
		Active:        true,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

func serviceShopNoTenant() models.ServiceShop {
	shop := validServiceShop()
	shop.ID = "shop-test-002"
	shop.TenantID = ""
	return shop
}

func validWorkOrderPart() models.WorkOrderPart {
	now := time.Now().UTC()
	return models.WorkOrderPart{
		ID:            "wop-test-001",
		TenantID:      "tenant-test",
		SchoolID:      "school-test",
		WorkOrderID:   "wo-test-001",
		ServiceShopID: "shop-001",
		PartID:        "part-001",
		PartName:      "LCD Screen",
		PartPUK:       "PUK-LCD-001",
		PartCategory:  "Display",
		DeviceModelID: "model-001",
		IsCompatible:  true,
		QtyPlanned:    2,
		QtyUsed:       0,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

func workOrderPartNoTenant() models.WorkOrderPart {
	part := validWorkOrderPart()
	part.ID = "wop-test-002"
	part.TenantID = ""
	return part
}
