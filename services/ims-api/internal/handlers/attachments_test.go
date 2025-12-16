package handlers_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/edvirons/ssp/ims/internal/blob"
	"github.com/edvirons/ssp/ims/internal/config"
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

// MockMinIOForAttachments creates a mock MinIO client configured for attachment testing
func MockMinIOForAttachments() *blob.MinIO {
	return &blob.MinIO{
		Bucket: "test-bucket",
		Expiry: 15 * time.Minute,
	}
}

func TestAttachmentHandler_Create(t *testing.T) {
	cfg := config.Config{
		AttachmentsBucket: "test-bucket",
	}
	logger := zap.NewNop()
	pg := testutil.SetupTestDB(t)
	minioClient := MockMinIOForAttachments()

	handler := handlers.NewAttachmentHandler(cfg, logger, pg, minioClient)

	tests := []struct {
		name       string
		body       string
		tenant     string
		school     string
		wantStatus int
		validate   func(t *testing.T, body string)
	}{
		{
			name: "valid attachment creation",
			body: `{
				"entityType": "incident",
				"entityId": "inc-001",
				"fileName": "screenshot.png",
				"contentType": "image/png",
				"sizeBytes": 1024
			}`,
			tenant:     "test-tenant",
			school:     "test-school",
			wantStatus: http.StatusCreated,
			validate: func(t *testing.T, body string) {
				var result models.Attachment
				err := json.Unmarshal([]byte(body), &result)
				require.NoError(t, err)
				assert.Equal(t, models.AttachmentIncident, result.EntityType)
				assert.Equal(t, "inc-001", result.EntityID)
				assert.Equal(t, "screenshot.png", result.FileName)
				assert.Equal(t, "image/png", result.ContentType)
				assert.Equal(t, int64(1024), result.SizeBytes)
				assert.NotEmpty(t, result.ID)
				assert.True(t, strings.HasPrefix(result.ID, "att_"))
				assert.NotEmpty(t, result.ObjectKey)
			},
		},
		{
			name: "work order attachment",
			body: `{
				"entityType": "work_order",
				"entityId": "wo-001",
				"fileName": "repair_photo.jpg",
				"contentType": "image/jpeg",
				"sizeBytes": 2048
			}`,
			tenant:     "test-tenant",
			school:     "test-school",
			wantStatus: http.StatusCreated,
			validate: func(t *testing.T, body string) {
				var result models.Attachment
				err := json.Unmarshal([]byte(body), &result)
				require.NoError(t, err)
				assert.Equal(t, models.AttachmentWorkOrder, result.EntityType)
			},
		},
		{
			name: "missing entityId",
			body: `{
				"entityType": "incident",
				"fileName": "test.png",
				"contentType": "image/png",
				"sizeBytes": 1024
			}`,
			tenant:     "test-tenant",
			school:     "test-school",
			wantStatus: http.StatusBadRequest,
			validate: func(t *testing.T, body string) {
				assert.Contains(t, body, "entityId")
			},
		},
		{
			name: "missing fileName",
			body: `{
				"entityType": "incident",
				"entityId": "inc-001",
				"contentType": "image/png",
				"sizeBytes": 1024
			}`,
			tenant:     "test-tenant",
			school:     "test-school",
			wantStatus: http.StatusBadRequest,
			validate: func(t *testing.T, body string) {
				assert.Contains(t, body, "fileName")
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
			name: "empty entityId after trim",
			body: `{
				"entityType": "incident",
				"entityId": "   ",
				"fileName": "test.png"
			}`,
			tenant:     "test-tenant",
			school:     "test-school",
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "PDF document attachment",
			body: `{
				"entityType": "work_order",
				"entityId": "wo-002",
				"fileName": "invoice.pdf",
				"contentType": "application/pdf",
				"sizeBytes": 51200
			}`,
			tenant:     "test-tenant",
			school:     "test-school",
			wantStatus: http.StatusCreated,
			validate: func(t *testing.T, body string) {
				var result models.Attachment
				err := json.Unmarshal([]byte(body), &result)
				require.NoError(t, err)
				assert.Equal(t, "application/pdf", result.ContentType)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/v1/attachments", strings.NewReader(tt.body))
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

func TestAttachmentHandler_GetByID(t *testing.T) {
	cfg := config.Config{}
	logger := zap.NewNop()
	pg := testutil.SetupTestDB(t)
	minioClient := MockMinIOForAttachments()

	handler := handlers.NewAttachmentHandler(cfg, logger, pg, minioClient)

	// Create test attachment
	fixtureConfig := testutil.DefaultFixtureConfig()
	att := testutil.CreateAttachment(t, pg.RawPool(), fixtureConfig, models.AttachmentIncident, "inc-001")

	tests := []struct {
		name       string
		attachmentID string
		tenant     string
		school     string
		wantStatus int
		validate   func(t *testing.T, body string)
	}{
		{
			name:       "found attachment",
			attachmentID: att.ID,
			tenant:     "test-tenant",
			school:     "test-school",
			wantStatus: http.StatusOK,
			validate: func(t *testing.T, body string) {
				var result models.Attachment
				err := json.Unmarshal([]byte(body), &result)
				require.NoError(t, err)
				assert.Equal(t, att.ID, result.ID)
				assert.Equal(t, att.FileName, result.FileName)
			},
		},
		{
			name:       "not found attachment",
			attachmentID: "att_nonexistent",
			tenant:     "test-tenant",
			school:     "test-school",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "wrong tenant",
			attachmentID: att.ID,
			tenant:     "wrong-tenant",
			school:     "test-school",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "wrong school",
			attachmentID: att.ID,
			tenant:     "test-tenant",
			school:     "wrong-school",
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			r.Get("/attachments/{id}", handler.GetByID)

			req := httptest.NewRequest(http.MethodGet, "/attachments/"+tt.attachmentID, nil)
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

func TestAttachmentHandler_List(t *testing.T) {
	cfg := config.Config{}
	logger := zap.NewNop()
	pg := testutil.SetupTestDB(t)
	minioClient := MockMinIOForAttachments()

	handler := handlers.NewAttachmentHandler(cfg, logger, pg, minioClient)

	// Create test attachments
	fixtureConfig := testutil.DefaultFixtureConfig()
	att1 := testutil.CreateAttachment(t, pg.RawPool(), fixtureConfig, models.AttachmentIncident, "inc-001")
	att2 := testutil.CreateAttachment(t, pg.RawPool(), fixtureConfig, models.AttachmentWorkOrder, "wo-001")
	_ = att1
	_ = att2

	tests := []struct {
		name       string
		queryParams string
		tenant     string
		school     string
		wantStatus int
		validate   func(t *testing.T, body string)
	}{
		{
			name:       "list all attachments",
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
			name:       "filter by entityType incident",
			queryParams: "?entityType=incident",
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
			name:       "filter by entityId",
			queryParams: "?entityId=inc-001",
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
			req := httptest.NewRequest(http.MethodGet, "/v1/attachments"+tt.queryParams, nil)
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

func TestAttachmentHandler_UploadURL(t *testing.T) {
	cfg := config.Config{}
	logger := zap.NewNop()
	pg := testutil.SetupTestDB(t)

	// Create mock MinIO client
	mockMinIO := mocks.NewMockMinIOClient()
	mockMinIO.MakeBucket(context.Background(), "test-bucket", nil)

	// Wrap in blob.MinIO struct
	minioClient := &blob.MinIO{
		Client: nil, // We'll use the mock directly in handler, but this shows the structure
		Bucket: "test-bucket",
		Expiry: 15 * time.Minute,
	}

	handler := handlers.NewAttachmentHandler(cfg, logger, pg, minioClient)

	// Create test attachment
	fixtureConfig := testutil.DefaultFixtureConfig()
	att := testutil.CreateAttachment(t, pg.RawPool(), fixtureConfig, models.AttachmentIncident, "inc-001")

	tests := []struct {
		name       string
		attachmentID string
		tenant     string
		school     string
		wantStatus int
		validate   func(t *testing.T, body string)
	}{
		{
			name:       "get upload URL for existing attachment",
			attachmentID: att.ID,
			tenant:     "test-tenant",
			school:     "test-school",
			wantStatus: http.StatusOK,
			validate: func(t *testing.T, body string) {
				var result map[string]interface{}
				err := json.Unmarshal([]byte(body), &result)
				require.NoError(t, err)
				assert.Equal(t, "PUT", result["method"])
				assert.NotEmpty(t, result["url"])
				assert.NotEmpty(t, result["objectKey"])
			},
		},
		{
			name:       "attachment not found",
			attachmentID: "att_nonexistent",
			tenant:     "test-tenant",
			school:     "test-school",
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			r.Get("/attachments/{id}/upload-url", handler.UploadURL)

			req := httptest.NewRequest(http.MethodGet, "/attachments/"+tt.attachmentID+"/upload-url", nil)
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

func TestAttachmentHandler_DownloadURL(t *testing.T) {
	cfg := config.Config{}
	logger := zap.NewNop()
	pg := testutil.SetupTestDB(t)

	minioClient := MockMinIOForAttachments()

	handler := handlers.NewAttachmentHandler(cfg, logger, pg, minioClient)

	// Create test attachment
	fixtureConfig := testutil.DefaultFixtureConfig()
	att := testutil.CreateAttachment(t, pg.RawPool(), fixtureConfig, models.AttachmentIncident, "inc-001")

	tests := []struct {
		name       string
		attachmentID string
		tenant     string
		school     string
		wantStatus int
		validate   func(t *testing.T, body string)
	}{
		{
			name:       "get download URL for existing attachment",
			attachmentID: att.ID,
			tenant:     "test-tenant",
			school:     "test-school",
			wantStatus: http.StatusOK,
			validate: func(t *testing.T, body string) {
				var result map[string]interface{}
				err := json.Unmarshal([]byte(body), &result)
				require.NoError(t, err)
				assert.NotEmpty(t, result["url"])
				assert.NotEmpty(t, result["objectKey"])
			},
		},
		{
			name:       "attachment not found",
			attachmentID: "att_nonexistent",
			tenant:     "test-tenant",
			school:     "test-school",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "wrong tenant",
			attachmentID: att.ID,
			tenant:     "wrong-tenant",
			school:     "test-school",
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			r.Get("/attachments/{id}/download-url", handler.DownloadURL)

			req := httptest.NewRequest(http.MethodGet, "/attachments/"+tt.attachmentID+"/download-url", nil)
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

func TestAttachmentHandler_MultipleEntityTypes(t *testing.T) {
	// Integration test to verify attachments work with different entity types
	cfg := config.Config{}
	logger := zap.NewNop()
	pg := testutil.SetupTestDB(t)
	minioClient := MockMinIOForAttachments()

	handler := handlers.NewAttachmentHandler(cfg, logger, pg, minioClient)

	entityTypes := []struct {
		entityType models.AttachmentEntityType
		entityID   string
	}{
		{models.AttachmentIncident, "inc-multi-001"},
		{models.AttachmentWorkOrder, "wo-multi-001"},
	}

	for _, et := range entityTypes {
		t.Run(string(et.entityType), func(t *testing.T) {
			body := `{
				"entityType": "` + string(et.entityType) + `",
				"entityId": "` + et.entityID + `",
				"fileName": "test.png",
				"contentType": "image/png",
				"sizeBytes": 1024
			}`

			req := httptest.NewRequest(http.MethodPost, "/v1/attachments", strings.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			ctx := middleware.WithTenantID(context.Background(), "test-tenant")
			ctx = middleware.WithSchoolID(ctx, "test-school")
			req = req.WithContext(ctx)

			rec := httptest.NewRecorder()
			handler.Create(rec, req)

			assert.Equal(t, http.StatusCreated, rec.Code)

			var result models.Attachment
			err := json.Unmarshal(rec.Body.Bytes(), &result)
			require.NoError(t, err)
			assert.Equal(t, et.entityType, result.EntityType)
			assert.Equal(t, et.entityID, result.EntityID)
		})
	}
}
