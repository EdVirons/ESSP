package handlers_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/edvirons/ssp/ims/internal/config"
	"github.com/edvirons/ssp/ims/internal/handlers"
	"github.com/edvirons/ssp/ims/internal/middleware"
	"github.com/edvirons/ssp/ims/internal/mocks"
	"github.com/edvirons/ssp/ims/internal/models"
	"github.com/edvirons/ssp/ims/internal/store"
	"github.com/edvirons/ssp/ims/internal/testutil"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestIncidentHandler_Create(t *testing.T) {
	// Setup
	cfg := config.Config{
		AutoRouteWorkOrders:   false,
		DefaultRepairLocation: "service_shop",
	}
	logger := zap.NewNop()
	mockRedis := mocks.NewMockRedisClient()

	// Create test database connection
	pg := testutil.SetupTestDB(t)

	// Setup test fixtures
	fixtureConfig := testutil.DefaultFixtureConfig()
	testutil.CreateSchoolSnapshot(t, pg.RawPool(), fixtureConfig)

	handler := handlers.NewIncidentHandler(cfg, logger, pg, mockRedis)

	tests := []struct {
		name       string
		body       string
		tenant     string
		school     string
		wantStatus int
		wantFields map[string]interface{}
		validate   func(t *testing.T, body string)
	}{
		{
			name: "valid incident creation",
			body: `{
				"deviceId": "dev-001",
				"title": "Screen broken",
				"description": "Device screen is cracked",
				"category": "hardware",
				"severity": "medium",
				"reportedBy": "teacher@school.com"
			}`,
			tenant:     "test-tenant",
			school:     "test-school",
			wantStatus: http.StatusOK,
			validate: func(t *testing.T, body string) {
				var result models.Incident
				err := json.Unmarshal([]byte(body), &result)
				require.NoError(t, err)
				assert.Equal(t, "dev-001", result.DeviceID)
				assert.Equal(t, "Screen broken", result.Title)
				assert.Equal(t, models.SeverityMedium, result.Severity)
				assert.Equal(t, models.IncidentNew, result.Status)
				assert.NotEmpty(t, result.ID)
				assert.True(t, strings.HasPrefix(result.ID, "inc_"))
			},
		},
		{
			name: "missing deviceId",
			body: `{
				"title": "Test incident",
				"severity": "low"
			}`,
			tenant:     "test-tenant",
			school:     "test-school",
			wantStatus: http.StatusBadRequest,
			validate: func(t *testing.T, body string) {
				assert.Contains(t, body, "deviceId")
			},
		},
		{
			name: "missing title",
			body: `{
				"deviceId": "dev-001",
				"severity": "low"
			}`,
			tenant:     "test-tenant",
			school:     "test-school",
			wantStatus: http.StatusBadRequest,
			validate: func(t *testing.T, body string) {
				assert.Contains(t, body, "title")
			},
		},
		{
			name:       "invalid json",
			body:       `{invalid json}`,
			tenant:     "test-tenant",
			school:     "test-school",
			wantStatus: http.StatusBadRequest,
			validate: func(t *testing.T, body string) {
				assert.Contains(t, body, "invalid json")
			},
		},
		{
			name: "empty deviceId after trim",
			body: `{
				"deviceId": "   ",
				"title": "Test incident"
			}`,
			tenant:     "test-tenant",
			school:     "test-school",
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "minimal valid incident",
			body: `{
				"deviceId": "dev-002",
				"title": "Minimal incident"
			}`,
			tenant:     "test-tenant",
			school:     "test-school",
			wantStatus: http.StatusOK,
			validate: func(t *testing.T, body string) {
				var result models.Incident
				err := json.Unmarshal([]byte(body), &result)
				require.NoError(t, err)
				assert.Equal(t, "dev-002", result.DeviceID)
				assert.Equal(t, "Minimal incident", result.Title)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/v1/incidents", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")

			// Add tenant and school context
			ctx := context.Background()
			ctx = middleware.WithTenantID(ctx, tt.tenant)
			ctx = middleware.WithSchoolID(ctx, tt.school)
			req = req.WithContext(ctx)

			rec := httptest.NewRecorder()
			handler.Create(rec, req)

			assert.Equal(t, tt.wantStatus, rec.Code, "unexpected status code")

			if tt.validate != nil {
				tt.validate(t, rec.Body.String())
			}
		})
	}
}

func TestIncidentHandler_GetByID(t *testing.T) {
	logger := zap.NewNop()
	cfg := config.Config{}
	mockRedis := mocks.NewMockRedisClient()

	pg := testutil.SetupTestDB(t)

	handler := handlers.NewIncidentHandler(cfg, logger, pg, mockRedis)

	// Create test incident
	fixtureConfig := testutil.DefaultFixtureConfig()
	incident := testutil.CreateIncident(t, pg.RawPool(), fixtureConfig, "dev-001")

	tests := []struct {
		name       string
		incidentID string
		tenant     string
		school     string
		wantStatus int
		validate   func(t *testing.T, body string)
	}{
		{
			name:       "found incident",
			incidentID: incident.ID,
			tenant:     "test-tenant",
			school:     "test-school",
			wantStatus: http.StatusOK,
			validate: func(t *testing.T, body string) {
				var result models.Incident
				err := json.Unmarshal([]byte(body), &result)
				require.NoError(t, err)
				assert.Equal(t, incident.ID, result.ID)
				assert.Equal(t, incident.Title, result.Title)
			},
		},
		{
			name:       "not found incident",
			incidentID: "inc_nonexistent",
			tenant:     "test-tenant",
			school:     "test-school",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "wrong tenant",
			incidentID: incident.ID,
			tenant:     "wrong-tenant",
			school:     "test-school",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "wrong school",
			incidentID: incident.ID,
			tenant:     "test-tenant",
			school:     "wrong-school",
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup router with chi for URL params
			r := chi.NewRouter()
			r.Get("/incidents/{id}", handler.GetByID)

			req := httptest.NewRequest(http.MethodGet, "/incidents/"+tt.incidentID, nil)
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

func TestIncidentHandler_List(t *testing.T) {
	logger := zap.NewNop()
	cfg := config.Config{}
	mockRedis := mocks.NewMockRedisClient()

	pg := testutil.SetupTestDB(t)

	handler := handlers.NewIncidentHandler(cfg, logger, pg, mockRedis)
	ctx := context.Background()

	// Create test incidents
	fixtureConfig := testutil.DefaultFixtureConfig()
	inc1 := testutil.CreateIncident(t, pg.RawPool(), fixtureConfig, "dev-001")
	inc2 := testutil.CreateIncident(t, pg.RawPool(), fixtureConfig, "dev-002")

	// Update inc2 status for filtering test
	_, _ = pg.Incidents().UpdateStatus(ctx, inc2.TenantID, inc2.SchoolID, inc2.ID, models.IncidentAcknowledged, time.Now())

	tests := []struct {
		name       string
		queryParams string
		tenant     string
		school     string
		wantStatus int
		validate   func(t *testing.T, body string)
	}{
		{
			name:       "list all incidents",
			queryParams: "",
			tenant:     "test-tenant",
			school:     "test-school",
			wantStatus: http.StatusOK,
			validate: func(t *testing.T, body string) {
				var result map[string]interface{}
				err := json.Unmarshal([]byte(body), &result)
				require.NoError(t, err)
				items := result["items"].([]interface{})
				assert.GreaterOrEqual(t, len(items), 2)
			},
		},
		{
			name:       "filter by status",
			queryParams: "?status=acknowledged",
			tenant:     "test-tenant",
			school:     "test-school",
			wantStatus: http.StatusOK,
			validate: func(t *testing.T, body string) {
				var result map[string]interface{}
				err := json.Unmarshal([]byte(body), &result)
				require.NoError(t, err)
				items := result["items"].([]interface{})
				assert.GreaterOrEqual(t, len(items), 1)
			},
		},
		{
			name:       "filter by deviceId",
			queryParams: "?deviceId=dev-001",
			tenant:     "test-tenant",
			school:     "test-school",
			wantStatus: http.StatusOK,
			validate: func(t *testing.T, body string) {
				var result map[string]interface{}
				err := json.Unmarshal([]byte(body), &result)
				require.NoError(t, err)
				items := result["items"].([]interface{})
				assert.GreaterOrEqual(t, len(items), 1)
			},
		},
		{
			name:       "with limit",
			queryParams: "?limit=1",
			tenant:     "test-tenant",
			school:     "test-school",
			wantStatus: http.StatusOK,
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
			req := httptest.NewRequest(http.MethodGet, "/v1/incidents"+tt.queryParams, nil)
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

func TestIncidentHandler_UpdateStatus(t *testing.T) {
	logger := zap.NewNop()
	cfg := config.Config{}
	mockRedis := mocks.NewMockRedisClient()

	pg := testutil.SetupTestDB(t)
	ctx := context.Background()

	handler := handlers.NewIncidentHandler(cfg, logger, pg, mockRedis)

	// Create test incident
	fixtureConfig := testutil.DefaultFixtureConfig()
	incident := testutil.CreateIncident(t, pg.RawPool(), fixtureConfig, "dev-001")

	tests := []struct {
		name       string
		incidentID string
		body       string
		tenant     string
		school     string
		wantStatus int
		validate   func(t *testing.T, body string)
	}{
		{
			name:       "valid transition new to acknowledged",
			incidentID: incident.ID,
			body:       `{"status": "acknowledged"}`,
			tenant:     "test-tenant",
			school:     "test-school",
			wantStatus: http.StatusOK,
			validate: func(t *testing.T, body string) {
				var result models.Incident
				err := json.Unmarshal([]byte(body), &result)
				require.NoError(t, err)
				assert.Equal(t, models.IncidentAcknowledged, result.Status)
			},
		},
		{
			name:       "invalid transition new to resolved",
			incidentID: incident.ID,
			body:       `{"status": "resolved"}`,
			tenant:     "test-tenant",
			school:     "test-school",
			wantStatus: http.StatusBadRequest,
			validate: func(t *testing.T, body string) {
				assert.Contains(t, body, "invalid status transition")
			},
		},
		{
			name:       "invalid json",
			incidentID: incident.ID,
			body:       `{invalid}`,
			tenant:     "test-tenant",
			school:     "test-school",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "incident not found",
			incidentID: "inc_nonexistent",
			body:       `{"status": "acknowledged"}`,
			tenant:     "test-tenant",
			school:     "test-school",
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup router with chi for URL params
			r := chi.NewRouter()
			r.Patch("/incidents/{id}/status", handler.UpdateStatus)

			req := httptest.NewRequest(http.MethodPatch, "/incidents/"+tt.incidentID+"/status", strings.NewReader(tt.body))
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

func TestIncidentHandler_StatusTransitions(t *testing.T) {
	logger := zap.NewNop()
	cfg := config.Config{}
	mockRedis := mocks.NewMockRedisClient()

	pg := testutil.SetupTestDB(t)

	handler := handlers.NewIncidentHandler(cfg, logger, pg, mockRedis)
	fixtureConfig := testutil.DefaultFixtureConfig()

	// Test valid transition sequence: new -> acknowledged -> in_progress -> resolved -> closed
	t.Run("valid transition sequence", func(t *testing.T) {
		incident := testutil.CreateIncident(t, pg.RawPool(), fixtureConfig, "dev-seq-001")

		transitions := []struct {
			from models.IncidentStatus
			to   models.IncidentStatus
		}{
			{models.IncidentNew, models.IncidentAcknowledged},
			{models.IncidentAcknowledged, models.IncidentInProgress},
			{models.IncidentInProgress, models.IncidentResolved},
			{models.IncidentResolved, models.IncidentClosed},
		}

		r := chi.NewRouter()
		r.Patch("/incidents/{id}/status", handler.UpdateStatus)

		for _, tr := range transitions {
			body := `{"status": "` + string(tr.to) + `"}`
			req := httptest.NewRequest(http.MethodPatch, "/incidents/"+incident.ID+"/status", strings.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			ctx := middleware.WithTenantID(context.Background(), "test-tenant")
			ctx = middleware.WithSchoolID(ctx, "test-school")
			req = req.WithContext(ctx)

			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusOK, rec.Code, "transition from %s to %s failed", tr.from, tr.to)
		}
	})

	// Test escalation path: new -> escalated
	t.Run("escalation from new", func(t *testing.T) {
		incident := testutil.CreateIncident(t, pg.RawPool(), fixtureConfig, "dev-esc-001")

		r := chi.NewRouter()
		r.Patch("/incidents/{id}/status", handler.UpdateStatus)

		body := `{"status": "escalated"}`
		req := httptest.NewRequest(http.MethodPatch, "/incidents/"+incident.ID+"/status", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		ctx := middleware.WithTenantID(context.Background(), "test-tenant")
		ctx = middleware.WithSchoolID(ctx, "test-school")
		req = req.WithContext(ctx)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})
}
