package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/edvirons/ssp/ims/internal/blob"
	"github.com/edvirons/ssp/ims/internal/middleware"
	"github.com/edvirons/ssp/ims/internal/models"
	"github.com/edvirons/ssp/ims/internal/store"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type PresentationsHandler struct {
	log   *zap.Logger
	pg    *store.Postgres
	minio *blob.MinIO
}

func NewPresentationsHandler(log *zap.Logger, pg *store.Postgres, minio *blob.MinIO) *PresentationsHandler {
	return &PresentationsHandler{log: log, pg: pg, minio: minio}
}

// List returns presentations with optional filtering.
func (h *PresentationsHandler) List(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())

	filters := models.PresentationFilters{
		Limit:  50,
		Offset: 0,
	}

	if v := r.URL.Query().Get("type"); v != "" {
		t := models.PresentationType(v)
		filters.Type = &t
	}
	if v := r.URL.Query().Get("category"); v != "" {
		c := models.PresentationCategory(v)
		filters.Category = &c
	}
	if v := r.URL.Query().Get("featured"); v != "" {
		featured := v == "true"
		filters.IsFeatured = &featured
	}
	if v := r.URL.Query().Get("active"); v != "" {
		active := v == "true"
		filters.IsActive = &active
	}
	if v := r.URL.Query().Get("search"); v != "" {
		filters.Search = &v
	}
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			filters.Limit = n
		}
	}
	if v := r.URL.Query().Get("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			filters.Offset = n
		}
	}

	presentations, total, err := h.pg.Presentations().List(r.Context(), tenant, filters)
	if err != nil {
		h.log.Error("failed to list presentations", zap.Error(err))
		http.Error(w, "failed to list presentations", http.StatusInternalServerError)
		return
	}

	// Add presigned URLs for each presentation
	for i := range presentations {
		if presentations[i].FileKey != "" {
			url, _ := h.minio.PresignGet(r.Context(), presentations[i].FileKey)
			presentations[i].DownloadURL = url
		}
		if presentations[i].ThumbnailKey != "" {
			url, _ := h.minio.PresignGet(r.Context(), presentations[i].ThumbnailKey)
			presentations[i].PreviewURL = url
		}
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"presentations": presentations,
		"total":         total,
	})
}

// GetByID returns a single presentation.
func (h *PresentationsHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())
	id := chi.URLParam(r, "id")

	p, err := h.pg.Presentations().GetByID(r.Context(), tenant, id)
	if err != nil {
		if err.Error() == "not found" {
			http.Error(w, "presentation not found", http.StatusNotFound)
			return
		}
		h.log.Error("failed to get presentation", zap.Error(err))
		http.Error(w, "failed to get presentation", http.StatusInternalServerError)
		return
	}

	// Add presigned URLs
	if p.FileKey != "" {
		url, _ := h.minio.PresignGet(r.Context(), p.FileKey)
		p.DownloadURL = url
	}
	if p.ThumbnailKey != "" {
		url, _ := h.minio.PresignGet(r.Context(), p.ThumbnailKey)
		p.PreviewURL = url
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(p)
}

type createPresentationReq struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Type        string   `json:"type"`
	Category    string   `json:"category"`
	Tags        []string `json:"tags"`
	IsFeatured  bool     `json:"isFeatured"`
	FileName    string   `json:"fileName"`
	FileSize    int64    `json:"fileSize"`
	FileType    string   `json:"fileType"`
}

// Create creates a new presentation record and returns upload URLs.
func (h *PresentationsHandler) Create(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())
	userID := middleware.UserID(r.Context())

	var req createPresentationReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Title) == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(req.Type) == "" {
		http.Error(w, "type is required", http.StatusBadRequest)
		return
	}

	// Validate MIME type
	if req.FileType != "" && !isAllowedMIMEType(req.FileType) {
		http.Error(w, "file type not allowed", http.StatusBadRequest)
		return
	}

	now := time.Now().UTC()
	id := store.NewID("pres")
	fileKey := fmt.Sprintf("presentations/%s/%s/%s", tenant, id, req.FileName)
	thumbnailKey := fmt.Sprintf("presentations/%s/%s/thumbnail.jpg", tenant, id)

	p := models.Presentation{
		ID:            id,
		TenantID:      tenant,
		Title:         strings.TrimSpace(req.Title),
		Description:   strings.TrimSpace(req.Description),
		Type:          models.PresentationType(req.Type),
		Category:      models.PresentationCategory(req.Category),
		FileKey:       fileKey,
		FileName:      req.FileName,
		FileSize:      req.FileSize,
		FileType:      req.FileType,
		ThumbnailKey:  thumbnailKey,
		Tags:          req.Tags,
		Version:       1,
		IsActive:      true,
		IsFeatured:    req.IsFeatured,
		ViewCount:     0,
		DownloadCount: 0,
		CreatedBy:     userID,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if p.Tags == nil {
		p.Tags = []string{}
	}
	if p.Category == "" {
		p.Category = models.CategoryGeneral
	}

	if err := h.pg.Presentations().Create(r.Context(), p); err != nil {
		h.log.Error("failed to create presentation", zap.Error(err))
		http.Error(w, "failed to create presentation", http.StatusInternalServerError)
		return
	}

	// Generate presigned upload URLs
	var uploadURL, thumbnailUploadURL string
	if req.FileName != "" {
		url, err := h.minio.PresignPut(r.Context(), fileKey, req.FileType)
		if err != nil {
			h.log.Error("failed to generate upload URL", zap.Error(err))
		} else {
			uploadURL = url
		}

		tURL, err := h.minio.PresignPut(r.Context(), thumbnailKey, "image/jpeg")
		if err != nil {
			h.log.Error("failed to generate thumbnail upload URL", zap.Error(err))
		} else {
			thumbnailUploadURL = tURL
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"presentation":       p,
		"uploadUrl":          uploadURL,
		"thumbnailUploadUrl": thumbnailUploadURL,
	})
}

// Update updates presentation metadata.
func (h *PresentationsHandler) Update(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())
	userID := middleware.UserID(r.Context())
	id := chi.URLParam(r, "id")

	var req models.UpdatePresentationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if err := h.pg.Presentations().Update(r.Context(), tenant, id, req, userID); err != nil {
		h.log.Error("failed to update presentation", zap.Error(err))
		http.Error(w, "failed to update presentation", http.StatusInternalServerError)
		return
	}

	p, err := h.pg.Presentations().GetByID(r.Context(), tenant, id)
	if err != nil {
		http.Error(w, "failed to get updated presentation", http.StatusInternalServerError)
		return
	}

	// Add presigned URLs
	if p.FileKey != "" {
		url, _ := h.minio.PresignGet(r.Context(), p.FileKey)
		p.DownloadURL = url
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(p)
}

// Delete soft deletes a presentation.
func (h *PresentationsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())
	id := chi.URLParam(r, "id")

	if err := h.pg.Presentations().Delete(r.Context(), tenant, id); err != nil {
		h.log.Error("failed to delete presentation", zap.Error(err))
		http.Error(w, "failed to delete presentation", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DownloadURL returns a presigned download URL.
func (h *PresentationsHandler) DownloadURL(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())
	userID := middleware.UserID(r.Context())
	id := chi.URLParam(r, "id")

	p, err := h.pg.Presentations().GetByID(r.Context(), tenant, id)
	if err != nil {
		if err.Error() == "not found" {
			http.Error(w, "presentation not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to get presentation", http.StatusInternalServerError)
		return
	}

	if p.FileKey == "" {
		http.Error(w, "no file associated with this presentation", http.StatusNotFound)
		return
	}

	url, err := h.minio.PresignGet(r.Context(), p.FileKey)
	if err != nil {
		h.log.Error("failed to generate download URL", zap.Error(err))
		http.Error(w, "failed to generate download URL", http.StatusInternalServerError)
		return
	}

	// Increment download count (best-effort)
	_ = h.pg.Presentations().IncrementDownloadCount(r.Context(), tenant, id)

	// Record view event (best-effort)
	view := models.PresentationView{
		ID:             store.NewID("view"),
		TenantID:       tenant,
		PresentationID: id,
		ViewedBy:       userID,
		ViewedAt:       time.Now().UTC(),
		Context:        "download",
	}
	_ = h.pg.PresentationViews().Create(r.Context(), view)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"url":        url,
		"fileName":   p.FileName,
		"fileType":   p.FileType,
		"expiresInS": int(h.minio.Expiry.Seconds()),
	})
}

// RecordView records a view event without downloading.
func (h *PresentationsHandler) RecordView(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())
	userID := middleware.UserID(r.Context())
	id := chi.URLParam(r, "id")

	var req struct {
		Context         string `json:"context"`
		DurationSeconds *int   `json:"durationSeconds"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)

	// Increment view count
	_ = h.pg.Presentations().IncrementViewCount(r.Context(), tenant, id)

	// Record view event
	view := models.PresentationView{
		ID:              store.NewID("view"),
		TenantID:        tenant,
		PresentationID:  id,
		ViewedBy:        userID,
		ViewedAt:        time.Now().UTC(),
		Context:         req.Context,
		DurationSeconds: req.DurationSeconds,
	}
	_ = h.pg.PresentationViews().Create(r.Context(), view)

	w.WriteHeader(http.StatusNoContent)
}

// GetTypes returns available presentation types.
func (h *PresentationsHandler) GetTypes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"types":      models.PresentationTypes,
		"categories": models.PresentationCategories,
	})
}

func isAllowedMIMEType(mimeType string) bool {
	for _, allowed := range models.AllowedPresentationMIMETypes {
		if mimeType == allowed {
			return true
		}
	}
	return false
}
