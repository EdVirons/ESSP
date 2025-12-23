package handlers_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/edvirons/ssp/ims/internal/handlers"
	"github.com/edvirons/ssp/ims/internal/middleware"
	"github.com/edvirons/ssp/ims/internal/models"
	"github.com/edvirons/ssp/ims/internal/testutil"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestBOMHandler_AddItem(t *testing.T) {
	logger := zap.NewNop()
	pg := testutil.SetupTestDB(t)

	handler := handlers.NewBOMHandler(logger, pg)

	// Setup test fixtures
	fixtureConfig := testutil.DefaultFixtureConfig()
	shop := testutil.CreateServiceShop(t, pg.RawPool(), fixtureConfig)

	// Create a work order with service shop
	ctx := context.Background()
	wo := models.WorkOrder{
		ID:            "wo_test_001",
		TenantID:      "test-tenant",
		SchoolID:      "test-school",
		DeviceID:      "dev-001",
		Status:        models.WorkOrderDraft,
		ServiceShopID: shop.ID,
		TaskType:      "repair",
	}
	query := `
		INSERT INTO work_orders
		(id, tenant_id, school_id, device_id, status, service_shop_id, task_type, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
	`
	_, err := pg.RawPool().Exec(ctx, query, wo.ID, wo.TenantID, wo.SchoolID, wo.DeviceID, wo.Status, wo.ServiceShopID, wo.TaskType)
	require.NoError(t, err)

	// Create a part
	part := testutil.CreatePart(t, pg.RawPool(), fixtureConfig)

	// Create inventory for the part
	_ = testutil.CreateInventoryItem(t, pg.RawPool(), fixtureConfig, shop.ID, part.ID, 100)

	tests := []struct {
		name        string
		workOrderID string
		body        string
		tenant      string
		school      string
		wantStatus  int
		validate    func(t *testing.T, body string)
	}{
		{
			name:        "valid BOM item addition",
			workOrderID: wo.ID,
			body: `{
				"partId": "` + part.ID + `",
				"qtyPlanned": 5
			}`,
			tenant:     "test-tenant",
			school:     "test-school",
			wantStatus: http.StatusCreated,
			validate: func(t *testing.T, body string) {
				var result models.WorkOrderPart
				err := json.Unmarshal([]byte(body), &result)
				require.NoError(t, err)
				assert.Equal(t, part.ID, result.PartID)
				assert.Equal(t, int64(5), result.QtyPlanned)
				assert.Equal(t, int64(0), result.QtyUsed)
				assert.NotEmpty(t, result.ID)
			},
		},
		{
			name:        "missing partId",
			workOrderID: wo.ID,
			body:        `{"qtyPlanned": 5}`,
			tenant:      "test-tenant",
			school:      "test-school",
			wantStatus:  http.StatusBadRequest,
			validate: func(t *testing.T, body string) {
				assert.Contains(t, body, "partId")
			},
		},
		{
			name:        "invalid qtyPlanned zero",
			workOrderID: wo.ID,
			body:        `{"partId": "` + part.ID + `", "qtyPlanned": 0}`,
			tenant:      "test-tenant",
			school:      "test-school",
			wantStatus:  http.StatusBadRequest,
		},
		{
			name:        "invalid qtyPlanned negative",
			workOrderID: wo.ID,
			body:        `{"partId": "` + part.ID + `", "qtyPlanned": -5}`,
			tenant:      "test-tenant",
			school:      "test-school",
			wantStatus:  http.StatusBadRequest,
		},
		{
			name:        "invalid json",
			workOrderID: wo.ID,
			body:        `{invalid}`,
			tenant:      "test-tenant",
			school:      "test-school",
			wantStatus:  http.StatusBadRequest,
		},
		{
			name:        "work order not found",
			workOrderID: "wo_nonexistent",
			body:        `{"partId": "` + part.ID + `", "qtyPlanned": 5}`,
			tenant:      "test-tenant",
			school:      "test-school",
			wantStatus:  http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			r.Post("/work-orders/{id}/bom", handler.AddItem)

			req := httptest.NewRequest(http.MethodPost, "/work-orders/"+tt.workOrderID+"/bom", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			ctx := context.Background()
			ctx = middleware.WithTenantID(ctx, tt.tenant)
			ctx = middleware.WithSchoolID(ctx, tt.school)
			req = req.WithContext(ctx)

			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantStatus, rec.Code)

			if tt.validate != nil {
				tt.validate(t, rec.Body.String())
			}
		})
	}
}

func TestBOMHandler_List(t *testing.T) {
	logger := zap.NewNop()
	pg := testutil.SetupTestDB(t)

	handler := handlers.NewBOMHandler(logger, pg)

	// Setup test fixtures
	fixtureConfig := testutil.DefaultFixtureConfig()
	shop := testutil.CreateServiceShop(t, pg.RawPool(), fixtureConfig)

	// Create work order
	ctx := context.Background()
	wo := models.WorkOrder{
		ID:            "wo_list_001",
		TenantID:      "test-tenant",
		SchoolID:      "test-school",
		DeviceID:      "dev-001",
		Status:        models.WorkOrderDraft,
		ServiceShopID: shop.ID,
		TaskType:      "repair",
	}
	query := `
		INSERT INTO work_orders
		(id, tenant_id, school_id, device_id, status, service_shop_id, task_type, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
	`
	_, err := pg.RawPool().Exec(ctx, query, wo.ID, wo.TenantID, wo.SchoolID, wo.DeviceID, wo.Status, wo.ServiceShopID, wo.TaskType)
	require.NoError(t, err)

	// Create parts and BOM items
	part1 := testutil.CreatePart(t, pg.RawPool(), fixtureConfig)
	part2 := testutil.CreatePart(t, pg.RawPool(), fixtureConfig)

	// Add BOM items
	bomItem1 := models.WorkOrderPart{
		ID:            "bom_001",
		TenantID:      "test-tenant",
		SchoolID:      "test-school",
		WorkOrderID:   wo.ID,
		ServiceShopID: shop.ID,
		PartID:        part1.ID,
		QtyPlanned:    10,
		QtyUsed:       0,
	}
	bomItem2 := models.WorkOrderPart{
		ID:            "bom_002",
		TenantID:      "test-tenant",
		SchoolID:      "test-school",
		WorkOrderID:   wo.ID,
		ServiceShopID: shop.ID,
		PartID:        part2.ID,
		QtyPlanned:    5,
		QtyUsed:       0,
	}

	bomQuery := `
		INSERT INTO work_order_parts
		(id, tenant_id, school_id, work_order_id, service_shop_id, part_id, qty_planned, qty_used, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
	`
	_, err = pg.RawPool().Exec(ctx, bomQuery, bomItem1.ID, bomItem1.TenantID, bomItem1.SchoolID, bomItem1.WorkOrderID, bomItem1.ServiceShopID, bomItem1.PartID, bomItem1.QtyPlanned, bomItem1.QtyUsed)
	require.NoError(t, err)
	_, err = pg.RawPool().Exec(ctx, bomQuery, bomItem2.ID, bomItem2.TenantID, bomItem2.SchoolID, bomItem2.WorkOrderID, bomItem2.ServiceShopID, bomItem2.PartID, bomItem2.QtyPlanned, bomItem2.QtyUsed)
	require.NoError(t, err)

	tests := []struct {
		name        string
		workOrderID string
		queryParams string
		tenant      string
		school      string
		wantStatus  int
		validate    func(t *testing.T, body string)
	}{
		{
			name:        "list all BOM items",
			workOrderID: wo.ID,
			queryParams: "",
			tenant:      "test-tenant",
			school:      "test-school",
			wantStatus:  http.StatusOK,
			validate: func(t *testing.T, body string) {
				var result map[string]interface{}
				err := json.Unmarshal([]byte(body), &result)
				require.NoError(t, err)
				items := result["items"].([]interface{})
				assert.Equal(t, 2, len(items))
			},
		},
		{
			name:        "list with limit",
			workOrderID: wo.ID,
			queryParams: "?limit=1",
			tenant:      "test-tenant",
			school:      "test-school",
			wantStatus:  http.StatusOK,
			validate: func(t *testing.T, body string) {
				var result map[string]interface{}
				err := json.Unmarshal([]byte(body), &result)
				require.NoError(t, err)
				items := result["items"].([]interface{})
				assert.LessOrEqual(t, len(items), 1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			r.Get("/work-orders/{id}/bom", handler.List)

			req := httptest.NewRequest(http.MethodGet, "/work-orders/"+tt.workOrderID+"/bom"+tt.queryParams, nil)
			ctx := context.Background()
			ctx = middleware.WithTenantID(ctx, tt.tenant)
			ctx = middleware.WithSchoolID(ctx, tt.school)
			req = req.WithContext(ctx)

			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantStatus, rec.Code)

			if tt.validate != nil {
				tt.validate(t, rec.Body.String())
			}
		})
	}
}

func TestBOMHandler_Consume(t *testing.T) {
	logger := zap.NewNop()
	pg := testutil.SetupTestDB(t)

	handler := handlers.NewBOMHandler(logger, pg)

	// Setup test fixtures
	fixtureConfig := testutil.DefaultFixtureConfig()
	shop := testutil.CreateServiceShop(t, pg.RawPool(), fixtureConfig)
	part := testutil.CreatePart(t, pg.RawPool(), fixtureConfig)

	// Create inventory
	_ = testutil.CreateInventoryItem(t, pg.RawPool(), fixtureConfig, shop.ID, part.ID, 100)

	// Create work order
	ctx := context.Background()
	wo := models.WorkOrder{
		ID:            "wo_consume_001",
		TenantID:      "test-tenant",
		SchoolID:      "test-school",
		DeviceID:      "dev-001",
		Status:        models.WorkOrderDraft,
		ServiceShopID: shop.ID,
		TaskType:      "repair",
	}
	woQuery := `
		INSERT INTO work_orders
		(id, tenant_id, school_id, device_id, status, service_shop_id, task_type, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
	`
	_, err := pg.RawPool().Exec(ctx, woQuery, wo.ID, wo.TenantID, wo.SchoolID, wo.DeviceID, wo.Status, wo.ServiceShopID, wo.TaskType)
	require.NoError(t, err)

	// Create BOM item
	bomItem := models.WorkOrderPart{
		ID:            "bom_consume_001",
		TenantID:      "test-tenant",
		SchoolID:      "test-school",
		WorkOrderID:   wo.ID,
		ServiceShopID: shop.ID,
		PartID:        part.ID,
		QtyPlanned:    10,
		QtyUsed:       0,
	}
	bomQuery := `
		INSERT INTO work_order_parts
		(id, tenant_id, school_id, work_order_id, service_shop_id, part_id, qty_planned, qty_used, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
	`
	_, err = pg.RawPool().Exec(ctx, bomQuery, bomItem.ID, bomItem.TenantID, bomItem.SchoolID, bomItem.WorkOrderID, bomItem.ServiceShopID, bomItem.PartID, bomItem.QtyPlanned, bomItem.QtyUsed)
	require.NoError(t, err)

	// Reserve inventory
	_, err = pg.RawPool().Exec(ctx, `
		UPDATE inventory
		SET qty_reserved = qty_reserved + $4
		WHERE tenant_id=$1 AND service_shop_id=$2 AND part_id=$3
	`, "test-tenant", shop.ID, part.ID, int64(10))
	require.NoError(t, err)

	tests := []struct {
		name        string
		workOrderID string
		itemID      string
		body        string
		tenant      string
		school      string
		wantStatus  int
		validate    func(t *testing.T, body string)
	}{
		{
			name:        "valid consumption",
			workOrderID: wo.ID,
			itemID:      bomItem.ID,
			body:        `{"qtyUsed": 3}`,
			tenant:      "test-tenant",
			school:      "test-school",
			wantStatus:  http.StatusOK,
			validate: func(t *testing.T, body string) {
				var result models.WorkOrderPart
				err := json.Unmarshal([]byte(body), &result)
				require.NoError(t, err)
				assert.Equal(t, int64(3), result.QtyUsed)
			},
		},
		{
			name:        "invalid qtyUsed zero",
			workOrderID: wo.ID,
			itemID:      bomItem.ID,
			body:        `{"qtyUsed": 0}`,
			tenant:      "test-tenant",
			school:      "test-school",
			wantStatus:  http.StatusBadRequest,
		},
		{
			name:        "invalid qtyUsed negative",
			workOrderID: wo.ID,
			itemID:      bomItem.ID,
			body:        `{"qtyUsed": -5}`,
			tenant:      "test-tenant",
			school:      "test-school",
			wantStatus:  http.StatusBadRequest,
		},
		{
			name:        "invalid json",
			workOrderID: wo.ID,
			itemID:      bomItem.ID,
			body:        `{invalid}`,
			tenant:      "test-tenant",
			school:      "test-school",
			wantStatus:  http.StatusBadRequest,
		},
		{
			name:        "BOM item not found",
			workOrderID: wo.ID,
			itemID:      "bom_nonexistent",
			body:        `{"qtyUsed": 3}`,
			tenant:      "test-tenant",
			school:      "test-school",
			wantStatus:  http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			r.Post("/work-orders/{id}/bom/{itemId}/consume", handler.Consume)

			req := httptest.NewRequest(http.MethodPost, "/work-orders/"+tt.workOrderID+"/bom/"+tt.itemID+"/consume", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			ctx := context.Background()
			ctx = middleware.WithTenantID(ctx, tt.tenant)
			ctx = middleware.WithSchoolID(ctx, tt.school)
			req = req.WithContext(ctx)

			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantStatus, rec.Code)

			if tt.validate != nil {
				tt.validate(t, rec.Body.String())
			}
		})
	}
}

func TestBOMHandler_Suggest(t *testing.T) {
	logger := zap.NewNop()
	pg := testutil.SetupTestDB(t)

	handler := handlers.NewBOMHandler(logger, pg)

	// Setup test fixtures
	fixtureConfig := testutil.DefaultFixtureConfig()
	shop := testutil.CreateServiceShop(t, pg.RawPool(), fixtureConfig)

	// Create work order
	ctx := context.Background()
	wo := models.WorkOrder{
		ID:            "wo_suggest_001",
		TenantID:      "test-tenant",
		SchoolID:      "test-school",
		DeviceID:      "dev-suggest-001",
		Status:        models.WorkOrderDraft,
		ServiceShopID: shop.ID,
		TaskType:      "repair",
	}
	woQuery := `
		INSERT INTO work_orders
		(id, tenant_id, school_id, device_id, status, service_shop_id, task_type, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
	`
	_, err := pg.RawPool().Exec(ctx, woQuery, wo.ID, wo.TenantID, wo.SchoolID, wo.DeviceID, wo.Status, wo.ServiceShopID, wo.TaskType)
	require.NoError(t, err)

	// Note: Suggest endpoint requires SSOT snapshots for device/part data which may not be
	// fully available in unit tests. These tests verify the basic HTTP flow.
	tests := []struct {
		name        string
		workOrderID string
		queryParams string
		tenant      string
		school      string
		wantStatus  int
	}{
		{
			name:        "work order not found",
			workOrderID: "wo_nonexistent",
			queryParams: "",
			tenant:      "test-tenant",
			school:      "test-school",
			wantStatus:  http.StatusNotFound,
		},
		// Additional tests would require setting up device snapshots and parts compatibility data
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			r.Get("/work-orders/{id}/bom/suggest", handler.Suggest)

			req := httptest.NewRequest(http.MethodGet, "/work-orders/"+tt.workOrderID+"/bom/suggest"+tt.queryParams, nil)
			ctx := context.Background()
			ctx = middleware.WithTenantID(ctx, tt.tenant)
			ctx = middleware.WithSchoolID(ctx, tt.school)
			req = req.WithContext(ctx)

			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestBOMHandler_InventoryReservation(t *testing.T) {
	// Integration test to verify inventory reservation when adding BOM items
	logger := zap.NewNop()
	pg := testutil.SetupTestDB(t)

	handler := handlers.NewBOMHandler(logger, pg)

	// Setup test fixtures
	fixtureConfig := testutil.DefaultFixtureConfig()
	shop := testutil.CreateServiceShop(t, pg.RawPool(), fixtureConfig)
	part := testutil.CreatePart(t, pg.RawPool(), fixtureConfig)

	// Create inventory with specific quantity
	_ = testutil.CreateInventoryItem(t, pg.RawPool(), fixtureConfig, shop.ID, part.ID, 100)

	// Create work order
	ctx := context.Background()
	wo := models.WorkOrder{
		ID:            "wo_inv_001",
		TenantID:      "test-tenant",
		SchoolID:      "test-school",
		DeviceID:      "dev-001",
		Status:        models.WorkOrderDraft,
		ServiceShopID: shop.ID,
		TaskType:      "repair",
	}
	woQuery := `
		INSERT INTO work_orders
		(id, tenant_id, school_id, device_id, status, service_shop_id, task_type, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
	`
	_, err := pg.RawPool().Exec(ctx, woQuery, wo.ID, wo.TenantID, wo.SchoolID, wo.DeviceID, wo.Status, wo.ServiceShopID, wo.TaskType)
	require.NoError(t, err)

	t.Run("inventory reserved when BOM item added", func(t *testing.T) {
		r := chi.NewRouter()
		r.Post("/work-orders/{id}/bom", handler.AddItem)

		body := `{"partId": "` + part.ID + `", "qtyPlanned": 10}`
		req := httptest.NewRequest(http.MethodPost, "/work-orders/"+wo.ID+"/bom", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		ctx := middleware.WithTenantID(context.Background(), "test-tenant")
		ctx = middleware.WithSchoolID(ctx, "test-school")
		req = req.WithContext(ctx)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)

		// Verify inventory was reserved
		var qtyReserved int64
		err := pg.RawPool().QueryRow(context.Background(), `
			SELECT qty_reserved FROM inventory
			WHERE tenant_id=$1 AND service_shop_id=$2 AND part_id=$3
		`, "test-tenant", shop.ID, part.ID).Scan(&qtyReserved)
		require.NoError(t, err)
		assert.Equal(t, int64(10), qtyReserved)
	})
}
