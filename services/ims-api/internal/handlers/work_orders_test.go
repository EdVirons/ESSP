package handlers_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/edvirons/ssp/ims/internal/handlers"
	"github.com/edvirons/ssp/ims/internal/middleware"
	"github.com/edvirons/ssp/ims/internal/mocks"
	"github.com/edvirons/ssp/ims/internal/models"
	"github.com/edvirons/ssp/ims/internal/testutil"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestWorkOrderHandler_Create(t *testing.T) {
	logger := zap.NewNop()
	mockAudit := mocks.NewMockAuditLogger()

	pg := testutil.SetupTestDB(t)

	handler := handlers.NewWorkOrderHandler(logger, pg, nil, mockAudit)

	// Setup test fixtures
	fixtureConfig := testutil.DefaultFixtureConfig()
	testutil.CreateSchoolSnapshot(t, pg.RawPool(), fixtureConfig)

	tests := []struct {
		name       string
		body       string
		tenant     string
		school     string
		wantStatus int
		validate   func(t *testing.T, body string)
	}{
		{
			name: "valid work order creation",
			body: `{
				"deviceId": "dev-001",
				"taskType": "repair",
				"serviceShopId": "shop-001",
				"repairLocation": "service_shop",
				"costEstimateCents": 5000,
				"notes": "Screen replacement"
			}`,
			tenant:     "test-tenant",
			school:     "test-school",
			wantStatus: http.StatusCreated,
			validate: func(t *testing.T, body string) {
				var result models.WorkOrder
				err := json.Unmarshal([]byte(body), &result)
				require.NoError(t, err)
				assert.Equal(t, "dev-001", result.DeviceID)
				assert.Equal(t, "repair", result.TaskType)
				assert.Equal(t, models.WorkOrderDraft, result.Status)
				assert.NotEmpty(t, result.ID)
				assert.True(t, strings.HasPrefix(result.ID, "wo_"))
			},
		},
		{
			name: "missing deviceId",
			body: `{
				"taskType": "repair"
			}`,
			tenant:     "test-tenant",
			school:     "test-school",
			wantStatus: http.StatusBadRequest,
			validate: func(t *testing.T, body string) {
				assert.Contains(t, body, "deviceId")
			},
		},
		{
			name: "missing taskType",
			body: `{
				"deviceId": "dev-001"
			}`,
			tenant:     "test-tenant",
			school:     "test-school",
			wantStatus: http.StatusBadRequest,
			validate: func(t *testing.T, body string) {
				assert.Contains(t, body, "taskType")
			},
		},
		{
			name:       "invalid json",
			body:       `{invalid json}`,
			tenant:     "test-tenant",
			school:     "test-school",
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "with incident ID",
			body: `{
				"incidentId": "inc-001",
				"deviceId": "dev-002",
				"taskType": "triage"
			}`,
			tenant:     "test-tenant",
			school:     "test-school",
			wantStatus: http.StatusCreated,
			validate: func(t *testing.T, body string) {
				var result models.WorkOrder
				err := json.Unmarshal([]byte(body), &result)
				require.NoError(t, err)
				assert.Equal(t, "inc-001", result.IncidentID)
			},
		},
		{
			name: "minimal valid work order",
			body: `{
				"deviceId": "dev-003",
				"taskType": "inspection"
			}`,
			tenant:     "test-tenant",
			school:     "test-school",
			wantStatus: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/v1/work-orders", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")

			ctx := context.Background()
			ctx = middleware.WithTenantID(ctx, tt.tenant)
			ctx = middleware.WithSchoolID(ctx, tt.school)
			req = req.WithContext(ctx)

			rec := httptest.NewRecorder()
			handler.Create(rec, req)

			assert.Equal(t, tt.wantStatus, rec.Code)

			if tt.validate != nil {
				tt.validate(t, rec.Body.String())
			}
		})
	}
}

func TestWorkOrderHandler_GetByID(t *testing.T) {
	logger := zap.NewNop()
	mockAudit := mocks.NewMockAuditLogger()

	pg := testutil.SetupTestDB(t)

	handler := handlers.NewWorkOrderHandler(logger, pg, nil, mockAudit)

	// Create test work order
	fixtureConfig := testutil.DefaultFixtureConfig()
	wo := testutil.CreateWorkOrder(t, pg.RawPool(), fixtureConfig, "", "dev-001")

	tests := []struct {
		name        string
		workOrderID string
		tenant      string
		school      string
		wantStatus  int
		validate    func(t *testing.T, body string)
	}{
		{
			name:        "found work order",
			workOrderID: wo.ID,
			tenant:      "test-tenant",
			school:      "test-school",
			wantStatus:  http.StatusOK,
			validate: func(t *testing.T, body string) {
				var result models.WorkOrder
				err := json.Unmarshal([]byte(body), &result)
				require.NoError(t, err)
				assert.Equal(t, wo.ID, result.ID)
				assert.Equal(t, wo.DeviceID, result.DeviceID)
			},
		},
		{
			name:        "not found work order",
			workOrderID: "wo_nonexistent",
			tenant:      "test-tenant",
			school:      "test-school",
			wantStatus:  http.StatusNotFound,
		},
		{
			name:        "wrong tenant",
			workOrderID: wo.ID,
			tenant:      "wrong-tenant",
			school:      "test-school",
			wantStatus:  http.StatusNotFound,
		},
		{
			name:        "wrong school",
			workOrderID: wo.ID,
			tenant:      "test-tenant",
			school:      "wrong-school",
			wantStatus:  http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			r.Get("/work-orders/{id}", handler.GetByID)

			req := httptest.NewRequest(http.MethodGet, "/work-orders/"+tt.workOrderID, nil)
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

func TestWorkOrderHandler_List(t *testing.T) {
	logger := zap.NewNop()
	mockAudit := mocks.NewMockAuditLogger()

	pg := testutil.SetupTestDB(t)
	ctx := context.Background()

	handler := handlers.NewWorkOrderHandler(logger, pg, nil, mockAudit)

	// Create test work orders
	fixtureConfig := testutil.DefaultFixtureConfig()
	_ = testutil.CreateWorkOrder(t, pg.RawPool(), fixtureConfig, "", "dev-001")
	wo2 := testutil.CreateWorkOrder(t, pg.RawPool(), fixtureConfig, "", "dev-002")

	// Update wo2 status for filtering test
	_, _ = pg.WorkOrders().UpdateStatus(ctx, wo2.TenantID, wo2.SchoolID, wo2.ID, models.WorkOrderAssigned, time.Now())

	tests := []struct {
		name        string
		queryParams string
		tenant      string
		school      string
		wantStatus  int
		validate    func(t *testing.T, body string)
	}{
		{
			name:        "list all work orders",
			queryParams: "",
			tenant:      "test-tenant",
			school:      "test-school",
			wantStatus:  http.StatusOK,
			validate: func(t *testing.T, body string) {
				var result map[string]interface{}
				err := json.Unmarshal([]byte(body), &result)
				require.NoError(t, err)
				items := result["items"].([]interface{})
				assert.GreaterOrEqual(t, len(items), 2)
			},
		},
		{
			name:        "filter by status",
			queryParams: "?status=assigned",
			tenant:      "test-tenant",
			school:      "test-school",
			wantStatus:  http.StatusOK,
			validate: func(t *testing.T, body string) {
				var result map[string]interface{}
				err := json.Unmarshal([]byte(body), &result)
				require.NoError(t, err)
				items := result["items"].([]interface{})
				assert.GreaterOrEqual(t, len(items), 1)
			},
		},
		{
			name:        "filter by deviceId",
			queryParams: "?deviceId=dev-001",
			tenant:      "test-tenant",
			school:      "test-school",
			wantStatus:  http.StatusOK,
			validate: func(t *testing.T, body string) {
				var result map[string]interface{}
				err := json.Unmarshal([]byte(body), &result)
				require.NoError(t, err)
				items := result["items"].([]interface{})
				assert.GreaterOrEqual(t, len(items), 1)
			},
		},
		{
			name:        "with limit",
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
			req := httptest.NewRequest(http.MethodGet, "/v1/work-orders"+tt.queryParams, nil)
			ctx := context.Background()
			ctx = middleware.WithTenantID(ctx, tt.tenant)
			ctx = middleware.WithSchoolID(ctx, tt.school)
			req = req.WithContext(ctx)

			rec := httptest.NewRecorder()
			handler.List(rec, req)

			assert.Equal(t, tt.wantStatus, rec.Code)

			if tt.validate != nil {
				tt.validate(t, rec.Body.String())
			}
		})
	}
}

func TestWorkOrderHandler_UpdateStatus(t *testing.T) {
	logger := zap.NewNop()
	mockAudit := mocks.NewMockAuditLogger()

	pg := testutil.SetupTestDB(t)

	handler := handlers.NewWorkOrderHandler(logger, pg, nil, mockAudit)

	// Create test work order
	fixtureConfig := testutil.DefaultFixtureConfig()
	wo := testutil.CreateWorkOrder(t, pg.RawPool(), fixtureConfig, "", "dev-001")

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
			name:        "valid transition draft to assigned",
			workOrderID: wo.ID,
			body:        `{"status": "assigned"}`,
			tenant:      "test-tenant",
			school:      "test-school",
			wantStatus:  http.StatusOK,
			validate: func(t *testing.T, body string) {
				var result models.WorkOrder
				err := json.Unmarshal([]byte(body), &result)
				require.NoError(t, err)
				assert.Equal(t, models.WorkOrderAssigned, result.Status)
			},
		},
		{
			name:        "invalid transition draft to completed",
			workOrderID: wo.ID,
			body:        `{"status": "completed"}`,
			tenant:      "test-tenant",
			school:      "test-school",
			wantStatus:  http.StatusBadRequest,
			validate: func(t *testing.T, body string) {
				assert.Contains(t, body, "invalid status transition")
			},
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
			body:        `{"status": "assigned"}`,
			tenant:      "test-tenant",
			school:      "test-school",
			wantStatus:  http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			r.Patch("/work-orders/{id}/status", handler.UpdateStatus)

			req := httptest.NewRequest(http.MethodPatch, "/work-orders/"+tt.workOrderID+"/status", strings.NewReader(tt.body))
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

func TestWorkOrderHandler_StatusTransitions(t *testing.T) {
	logger := zap.NewNop()
	mockAudit := mocks.NewMockAuditLogger()

	pg := testutil.SetupTestDB(t)

	handler := handlers.NewWorkOrderHandler(logger, pg, nil, mockAudit)
	fixtureConfig := testutil.DefaultFixtureConfig()

	// Test valid transition sequence: draft -> assigned -> in_repair -> qa -> completed -> approved
	t.Run("valid transition sequence", func(t *testing.T) {
		wo := testutil.CreateWorkOrder(t, pg.RawPool(), fixtureConfig, "", "dev-seq-001")

		transitions := []struct {
			from models.WorkOrderStatus
			to   models.WorkOrderStatus
		}{
			{models.WorkOrderDraft, models.WorkOrderAssigned},
			{models.WorkOrderAssigned, models.WorkOrderInRepair},
			{models.WorkOrderInRepair, models.WorkOrderQA},
			{models.WorkOrderQA, models.WorkOrderCompleted},
			{models.WorkOrderCompleted, models.WorkOrderApproved},
		}

		r := chi.NewRouter()
		r.Patch("/work-orders/{id}/status", handler.UpdateStatus)

		for _, tr := range transitions {
			body := `{"status": "` + string(tr.to) + `"}`
			req := httptest.NewRequest(http.MethodPatch, "/work-orders/"+wo.ID+"/status", strings.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			ctx := middleware.WithTenantID(context.Background(), "test-tenant")
			ctx = middleware.WithSchoolID(ctx, "test-school")
			req = req.WithContext(ctx)

			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusOK, rec.Code, "transition from %s to %s failed", tr.from, tr.to)
		}
	})

	// Test direct completion path (skipping QA)
	t.Run("skip qa path in_repair to completed", func(t *testing.T) {
		wo := testutil.CreateWorkOrder(t, pg.RawPool(), fixtureConfig, "", "dev-skip-001")

		r := chi.NewRouter()
		r.Patch("/work-orders/{id}/status", handler.UpdateStatus)

		// First transition to assigned
		body := `{"status": "assigned"}`
		req := httptest.NewRequest(http.MethodPatch, "/work-orders/"+wo.ID+"/status", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		ctx := middleware.WithTenantID(context.Background(), "test-tenant")
		ctx = middleware.WithSchoolID(ctx, "test-school")
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)

		// Then to in_repair
		body = `{"status": "in_repair"}`
		req = httptest.NewRequest(http.MethodPatch, "/work-orders/"+wo.ID+"/status", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(ctx)
		rec = httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)

		// Finally directly to completed (skip QA)
		body = `{"status": "completed"}`
		req = httptest.NewRequest(http.MethodPatch, "/work-orders/"+wo.ID+"/status", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(ctx)
		rec = httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
	})
}
