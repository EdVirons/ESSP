package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/edvirons/ssp/ims/internal/blob"
	"github.com/edvirons/ssp/ims/internal/config"
	"github.com/edvirons/ssp/ims/internal/middleware"
	"github.com/edvirons/ssp/ims/internal/models"
	"github.com/edvirons/ssp/ims/internal/store"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type AttachmentHandler struct {
	cfg   config.Config
	log   *zap.Logger
	pg    *store.Postgres
	minio *blob.MinIO
}

func NewAttachmentHandler(cfg config.Config, log *zap.Logger, pg *store.Postgres, minioClient *blob.MinIO) *AttachmentHandler {
	return &AttachmentHandler{cfg: cfg, log: log, pg: pg, minio: minioClient}
}

type createAttachmentReq struct {
	EntityType  models.AttachmentEntityType `json:"entityType"`
	EntityID    string                      `json:"entityId"`
	FileName    string                      `json:"fileName"`
	ContentType string                      `json:"contentType"`
	SizeBytes   int64                       `json:"sizeBytes"`
}

func (h *AttachmentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createAttachmentReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(req.EntityID) == "" || strings.TrimSpace(req.FileName) == "" {
		http.Error(w, "entityId and fileName are required", http.StatusBadRequest)
		return
	}

	tenant := middleware.TenantID(r.Context())
	school := middleware.SchoolID(r.Context())

	now := time.Now().UTC()
	att := models.Attachment{
		ID:          store.NewID("att"),
		TenantID:    tenant,
		SchoolID:    school,
		EntityType:  req.EntityType,
		EntityID:    strings.TrimSpace(req.EntityID),
		FileName:    strings.TrimSpace(req.FileName),
		ContentType: strings.TrimSpace(req.ContentType),
		SizeBytes:   req.SizeBytes,
		ObjectKey:   store.ObjectKeyForAttachment(tenant, school, req.EntityType, req.EntityID, now, req.FileName),
		CreatedAt:   now,
	}

	if err := h.pg.Attachments().Create(r.Context(), att); err != nil {
		http.Error(w, "failed to create attachment", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, att)
}

func (h *AttachmentHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	tenant := middleware.TenantID(r.Context())
	school := middleware.SchoolID(r.Context())

	att, err := h.pg.Attachments().GetByID(r.Context(), tenant, school, id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, att)
}

func (h *AttachmentHandler) List(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())
	school := middleware.SchoolID(r.Context())

	entityType := strings.TrimSpace(r.URL.Query().Get("entityType"))
	entityID := strings.TrimSpace(r.URL.Query().Get("entityId"))
	limit := parseLimit(r.URL.Query().Get("limit"), 50, 200)

	curT, curID, hasCur := decodeCursor(strings.TrimSpace(r.URL.Query().Get("cursor")))

	items, next, err := h.pg.Attachments().List(r.Context(), store.AttachmentListParams{
		TenantID:        tenant,
		SchoolID:        school,
		EntityType:      entityType,
		EntityID:        entityID,
		Limit:           limit,
		CursorCreatedAt: curT,
		CursorID:        curID,
		HasCursor:       hasCur,
	})
	if err != nil {
		http.Error(w, "failed to list attachments", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items, "nextCursor": next})
}

// UploadURL returns a stub URL. In production, replace with real presigned PUT from MinIO/S3.
// UploadURL returns a **presigned PUT** URL from MinIO.
func (h *AttachmentHandler) UploadURL(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	tenant := middleware.TenantID(r.Context())
	school := middleware.SchoolID(r.Context())

	att, err := h.pg.Attachments().GetByID(r.Context(), tenant, school, id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	url, err := h.minio.PresignPut(r.Context(), att.ObjectKey, att.ContentType)
	if err != nil {
		http.Error(w, "failed to presign upload url", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"method":      "PUT",
		"url":         url,
		"bucket":      h.minio.Bucket,
		"objectKey":   att.ObjectKey,
		"expiresInS":  int(h.minio.Expiry.Seconds()),
		"contentType": att.ContentType,
	})
}

// DownloadURL returns a **presigned GET** URL from MinIO.
func (h *AttachmentHandler) DownloadURL(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	tenant := middleware.TenantID(r.Context())
	school := middleware.SchoolID(r.Context())

	att, err := h.pg.Attachments().GetByID(r.Context(), tenant, school, id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	url, err := h.minio.PresignGet(r.Context(), att.ObjectKey)
	if err != nil {
		http.Error(w, "failed to presign download url", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"url":        url,
		"bucket":     h.minio.Bucket,
		"objectKey":  att.ObjectKey,
		"expiresInS": int(h.minio.Expiry.Seconds()),
	})
}
