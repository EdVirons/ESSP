package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/edvirons/ssp/ims/internal/middleware"
	"github.com/edvirons/ssp/ims/internal/models"
	"github.com/edvirons/ssp/ims/internal/store"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type MarketingKBHandler struct {
	log *zap.Logger
	pg  *store.Postgres
}

func NewMarketingKBHandler(log *zap.Logger, pg *store.Postgres) *MarketingKBHandler {
	return &MarketingKBHandler{log: log, pg: pg}
}

// ==================== Article Handlers ====================

type createMKBArticleReq struct {
	Title       string   `json:"title"`
	Slug        string   `json:"slug"`
	Summary     string   `json:"summary"`
	Content     string   `json:"content"`
	ContentType string   `json:"contentType"`
	Personas    []string `json:"personas"`
	ContextTags []string `json:"contextTags"`
	Tags        []string `json:"tags"`
}

func (h *MarketingKBHandler) CreateArticle(w http.ResponseWriter, r *http.Request) {
	var req createMKBArticleReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(req.Title) == "" || strings.TrimSpace(req.Content) == "" {
		http.Error(w, "title and content are required", http.StatusBadRequest)
		return
	}

	tenant := middleware.TenantID(r.Context())
	userID := middleware.UserID(r.Context())
	userName := middleware.UserName(r.Context())

	// Generate slug if not provided
	slug := strings.TrimSpace(req.Slug)
	if slug == "" {
		slug = generateSlug(req.Title)
	}

	// Check if slug exists
	exists, err := h.pg.MarketingKB().ArticleSlugExists(r.Context(), tenant, slug, "")
	if err != nil {
		h.log.Error("failed to check slug", zap.Error(err))
		http.Error(w, "failed to create article", http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "slug already exists", http.StatusConflict)
		return
	}

	// Validate content type
	contentType := models.MKBContentType(req.ContentType)
	if contentType == "" {
		contentType = models.MKBContentTypeMessaging
	}
	if !isValidMKBContentType(contentType) {
		http.Error(w, "invalid content type", http.StatusBadRequest)
		return
	}

	now := time.Now().UTC()
	personas := req.Personas
	if personas == nil {
		personas = []string{}
	}
	contextTags := req.ContextTags
	if contextTags == nil {
		contextTags = []string{}
	}
	tags := req.Tags
	if tags == nil {
		tags = []string{}
	}

	article := models.MKBArticle{
		ID:            store.NewID("mkb"),
		TenantID:      tenant,
		Title:         strings.TrimSpace(req.Title),
		Slug:          slug,
		Summary:       strings.TrimSpace(req.Summary),
		Content:       req.Content,
		ContentType:   contentType,
		Personas:      personas,
		ContextTags:   contextTags,
		Tags:          tags,
		Version:       1,
		Status:        models.MKBStatusDraft,
		UsageCount:    0,
		CreatedByID:   userID,
		CreatedByName: userName,
		UpdatedByID:   userID,
		UpdatedByName: userName,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if err := h.pg.MarketingKB().CreateArticle(r.Context(), article); err != nil {
		h.log.Error("failed to create mkb article", zap.Error(err))
		http.Error(w, "failed to create article", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, article)
}

func (h *MarketingKBHandler) GetArticleByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	tenant := middleware.TenantID(r.Context())

	article, err := h.pg.MarketingKB().GetArticleByID(r.Context(), tenant, id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, article)
}

func (h *MarketingKBHandler) GetArticleBySlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	tenant := middleware.TenantID(r.Context())

	article, err := h.pg.MarketingKB().GetArticleBySlug(r.Context(), tenant, slug)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, article)
}

type updateMKBArticleReq struct {
	Title       *string   `json:"title"`
	Slug        *string   `json:"slug"`
	Summary     *string   `json:"summary"`
	Content     *string   `json:"content"`
	ContentType *string   `json:"contentType"`
	Personas    *[]string `json:"personas"`
	ContextTags *[]string `json:"contextTags"`
	Tags        *[]string `json:"tags"`
	Status      *string   `json:"status"`
}

func (h *MarketingKBHandler) UpdateArticle(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	tenant := middleware.TenantID(r.Context())
	userID := middleware.UserID(r.Context())
	userName := middleware.UserName(r.Context())

	// Get existing article
	article, err := h.pg.MarketingKB().GetArticleByID(r.Context(), tenant, id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	// Parse update request
	var req updateMKBArticleReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	// Apply updates
	if req.Title != nil {
		article.Title = strings.TrimSpace(*req.Title)
	}
	if req.Slug != nil {
		newSlug := strings.TrimSpace(*req.Slug)
		if newSlug != article.Slug {
			exists, err := h.pg.MarketingKB().ArticleSlugExists(r.Context(), tenant, newSlug, article.ID)
			if err != nil {
				h.log.Error("failed to check slug", zap.Error(err))
				http.Error(w, "failed to update article", http.StatusInternalServerError)
				return
			}
			if exists {
				http.Error(w, "slug already exists", http.StatusConflict)
				return
			}
			article.Slug = newSlug
		}
	}
	if req.Summary != nil {
		article.Summary = strings.TrimSpace(*req.Summary)
	}
	if req.Content != nil {
		article.Content = *req.Content
	}
	if req.ContentType != nil {
		ct := models.MKBContentType(*req.ContentType)
		if !isValidMKBContentType(ct) {
			http.Error(w, "invalid content type", http.StatusBadRequest)
			return
		}
		article.ContentType = ct
	}
	if req.Personas != nil {
		article.Personas = *req.Personas
	}
	if req.ContextTags != nil {
		article.ContextTags = *req.ContextTags
	}
	if req.Tags != nil {
		article.Tags = *req.Tags
	}
	if req.Status != nil {
		status := models.MKBArticleStatus(*req.Status)
		// Only allow certain status transitions
		if status == models.MKBStatusReview && article.Status == models.MKBStatusDraft {
			article.Status = status
		} else if status == models.MKBStatusDraft && article.Status == models.MKBStatusReview {
			// Allow going back to draft from review
			article.Status = status
		}
		// Note: approved status change should use the Approve endpoint
	}

	article.Version++
	article.UpdatedByID = userID
	article.UpdatedByName = userName
	article.UpdatedAt = time.Now().UTC()

	if err := h.pg.MarketingKB().UpdateArticle(r.Context(), article); err != nil {
		h.log.Error("failed to update mkb article", zap.Error(err))
		http.Error(w, "failed to update article", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, article)
}

func (h *MarketingKBHandler) DeleteArticle(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	tenant := middleware.TenantID(r.Context())

	if err := h.pg.MarketingKB().DeleteArticle(r.Context(), tenant, id); err != nil {
		if err.Error() == "not found" {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		h.log.Error("failed to delete mkb article", zap.Error(err))
		http.Error(w, "failed to delete article", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *MarketingKBHandler) ListArticles(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	contentType := strings.TrimSpace(r.URL.Query().Get("contentType"))
	persona := strings.TrimSpace(r.URL.Query().Get("persona"))
	contextTag := strings.TrimSpace(r.URL.Query().Get("contextTag"))
	status := strings.TrimSpace(r.URL.Query().Get("status"))
	limit := parseLimit(r.URL.Query().Get("limit"), 50, 200)
	curT, curID, hasCur := decodeCursor(strings.TrimSpace(r.URL.Query().Get("cursor")))

	items, next, err := h.pg.MarketingKB().ListArticles(r.Context(), models.MKBArticleListParams{
		TenantID:    tenant,
		ContentType: contentType,
		Persona:     persona,
		ContextTag:  contextTag,
		Status:      status,
		Query:       q,
		Limit:       limit,
		HasCursor:   hasCur,
		CursorTime:  curT,
		CursorID:    curID,
	})
	if err != nil {
		h.log.Error("failed to list mkb articles", zap.Error(err))
		http.Error(w, "failed to list articles", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items, "nextCursor": next})
}

func (h *MarketingKBHandler) SearchArticles(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	if q == "" {
		http.Error(w, "query parameter q is required", http.StatusBadRequest)
		return
	}

	limit := parseLimit(r.URL.Query().Get("limit"), 20, 100)
	items, err := h.pg.MarketingKB().SearchArticles(r.Context(), tenant, q, limit)
	if err != nil {
		h.log.Error("failed to search mkb articles", zap.Error(err))
		http.Error(w, "failed to search articles", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

func (h *MarketingKBHandler) GetArticleStats(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())
	stats, err := h.pg.MarketingKB().GetArticleStats(r.Context(), tenant)
	if err != nil {
		h.log.Error("failed to get mkb stats", zap.Error(err))
		http.Error(w, "failed to get stats", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, stats)
}

func (h *MarketingKBHandler) ApproveArticle(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	tenant := middleware.TenantID(r.Context())
	userID := middleware.UserID(r.Context())
	userName := middleware.UserName(r.Context())

	if err := h.pg.MarketingKB().ApproveArticle(r.Context(), tenant, id, userID, userName); err != nil {
		if err.Error() == "not found or already approved" {
			http.Error(w, "article not found or already approved", http.StatusNotFound)
			return
		}
		h.log.Error("failed to approve mkb article", zap.Error(err))
		http.Error(w, "failed to approve article", http.StatusInternalServerError)
		return
	}

	// Return the updated article
	article, err := h.pg.MarketingKB().GetArticleByID(r.Context(), tenant, id)
	if err != nil {
		h.log.Error("failed to get approved article", zap.Error(err))
		http.Error(w, "article approved but failed to retrieve", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, article)
}

func (h *MarketingKBHandler) SubmitForReview(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	tenant := middleware.TenantID(r.Context())
	userID := middleware.UserID(r.Context())
	userName := middleware.UserName(r.Context())

	if err := h.pg.MarketingKB().SubmitForReview(r.Context(), tenant, id, userID, userName); err != nil {
		if err.Error() == "not found or not a draft" {
			http.Error(w, "article not found or not a draft", http.StatusNotFound)
			return
		}
		h.log.Error("failed to submit mkb article for review", zap.Error(err))
		http.Error(w, "failed to submit article for review", http.StatusInternalServerError)
		return
	}

	// Return the updated article
	article, err := h.pg.MarketingKB().GetArticleByID(r.Context(), tenant, id)
	if err != nil {
		h.log.Error("failed to get article after submit", zap.Error(err))
		http.Error(w, "article submitted but failed to retrieve", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, article)
}

func (h *MarketingKBHandler) RecordUsage(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	tenant := middleware.TenantID(r.Context())

	if err := h.pg.MarketingKB().RecordUsage(r.Context(), tenant, id); err != nil {
		if err.Error() == "not found" {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		h.log.Error("failed to record mkb usage", zap.Error(err))
		http.Error(w, "failed to record usage", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ==================== Pitch Kit Handlers ====================

type createPitchKitReq struct {
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	TargetPersona string   `json:"targetPersona"`
	ContextTags   []string `json:"contextTags"`
	ArticleIDs    []string `json:"articleIds"`
	IsTemplate    bool     `json:"isTemplate"`
}

func (h *MarketingKBHandler) CreatePitchKit(w http.ResponseWriter, r *http.Request) {
	var req createPitchKitReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(req.Name) == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	tenant := middleware.TenantID(r.Context())
	userID := middleware.UserID(r.Context())
	userName := middleware.UserName(r.Context())

	now := time.Now().UTC()
	contextTags := req.ContextTags
	if contextTags == nil {
		contextTags = []string{}
	}
	articleIDs := req.ArticleIDs
	if articleIDs == nil {
		articleIDs = []string{}
	}

	targetPersona := req.TargetPersona
	if targetPersona == "" {
		targetPersona = string(models.MKBPersonaDirector)
	}

	pk := models.PitchKit{
		ID:            store.NewID("pk"),
		TenantID:      tenant,
		Name:          strings.TrimSpace(req.Name),
		Description:   strings.TrimSpace(req.Description),
		TargetPersona: targetPersona,
		ContextTags:   contextTags,
		ArticleIDs:    articleIDs,
		IsTemplate:    req.IsTemplate,
		CreatedByID:   userID,
		CreatedByName: userName,
		UpdatedByID:   userID,
		UpdatedByName: userName,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if err := h.pg.MarketingKB().CreatePitchKit(r.Context(), pk); err != nil {
		h.log.Error("failed to create pitch kit", zap.Error(err))
		http.Error(w, "failed to create pitch kit", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, pk)
}

func (h *MarketingKBHandler) GetPitchKitByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	tenant := middleware.TenantID(r.Context())

	// Get pitch kit with articles populated
	pk, err := h.pg.MarketingKB().GetPitchKitWithArticles(r.Context(), tenant, id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, pk)
}

type updatePitchKitReq struct {
	Name          *string   `json:"name"`
	Description   *string   `json:"description"`
	TargetPersona *string   `json:"targetPersona"`
	ContextTags   *[]string `json:"contextTags"`
	ArticleIDs    *[]string `json:"articleIds"`
	IsTemplate    *bool     `json:"isTemplate"`
}

func (h *MarketingKBHandler) UpdatePitchKit(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	tenant := middleware.TenantID(r.Context())
	userID := middleware.UserID(r.Context())
	userName := middleware.UserName(r.Context())

	// Get existing pitch kit
	pk, err := h.pg.MarketingKB().GetPitchKitByID(r.Context(), tenant, id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	// Parse update request
	var req updatePitchKitReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	// Apply updates
	if req.Name != nil {
		pk.Name = strings.TrimSpace(*req.Name)
	}
	if req.Description != nil {
		pk.Description = strings.TrimSpace(*req.Description)
	}
	if req.TargetPersona != nil {
		pk.TargetPersona = *req.TargetPersona
	}
	if req.ContextTags != nil {
		pk.ContextTags = *req.ContextTags
	}
	if req.ArticleIDs != nil {
		pk.ArticleIDs = *req.ArticleIDs
	}
	if req.IsTemplate != nil {
		pk.IsTemplate = *req.IsTemplate
	}

	pk.UpdatedByID = userID
	pk.UpdatedByName = userName
	pk.UpdatedAt = time.Now().UTC()

	if err := h.pg.MarketingKB().UpdatePitchKit(r.Context(), pk); err != nil {
		h.log.Error("failed to update pitch kit", zap.Error(err))
		http.Error(w, "failed to update pitch kit", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, pk)
}

func (h *MarketingKBHandler) DeletePitchKit(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	tenant := middleware.TenantID(r.Context())

	if err := h.pg.MarketingKB().DeletePitchKit(r.Context(), tenant, id); err != nil {
		if err.Error() == "not found" {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		h.log.Error("failed to delete pitch kit", zap.Error(err))
		http.Error(w, "failed to delete pitch kit", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *MarketingKBHandler) ListPitchKits(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())
	targetPersona := strings.TrimSpace(r.URL.Query().Get("targetPersona"))
	isTemplateStr := strings.TrimSpace(r.URL.Query().Get("isTemplate"))
	limit := parseLimit(r.URL.Query().Get("limit"), 50, 200)
	curT, curID, hasCur := decodeCursor(strings.TrimSpace(r.URL.Query().Get("cursor")))

	var isTemplate *bool
	if isTemplateStr == "true" {
		v := true
		isTemplate = &v
	} else if isTemplateStr == "false" {
		v := false
		isTemplate = &v
	}

	items, next, err := h.pg.MarketingKB().ListPitchKits(r.Context(), models.PitchKitListParams{
		TenantID:      tenant,
		TargetPersona: targetPersona,
		IsTemplate:    isTemplate,
		Limit:         limit,
		HasCursor:     hasCur,
		CursorTime:    curT,
		CursorID:      curID,
	})
	if err != nil {
		h.log.Error("failed to list pitch kits", zap.Error(err))
		http.Error(w, "failed to list pitch kits", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items, "nextCursor": next})
}

// ==================== Helper Functions ====================

func isValidMKBContentType(ct models.MKBContentType) bool {
	for _, valid := range models.ValidMKBContentTypes() {
		if ct == valid {
			return true
		}
	}
	return false
}
