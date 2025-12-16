package handlers

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/edvirons/ssp/ims/internal/middleware"
	"github.com/edvirons/ssp/ims/internal/models"
	"github.com/edvirons/ssp/ims/internal/store"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type KBArticlesHandler struct {
	log *zap.Logger
	pg  *store.Postgres
}

func NewKBArticlesHandler(log *zap.Logger, pg *store.Postgres) *KBArticlesHandler {
	return &KBArticlesHandler{log: log, pg: pg}
}

type createKBArticleReq struct {
	Title          string   `json:"title"`
	Slug           string   `json:"slug"`
	Summary        string   `json:"summary"`
	Content        string   `json:"content"`
	ContentType    string   `json:"contentType"`
	Module         string   `json:"module"`
	LifecycleStage string   `json:"lifecycleStage"`
	Tags           []string `json:"tags"`
}

func (h *KBArticlesHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createKBArticleReq
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
	exists, err := h.pg.KBArticles().SlugExists(r.Context(), tenant, slug, "")
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
	contentType := models.KBContentType(req.ContentType)
	if contentType == "" {
		contentType = models.KBContentTypeRunbook
	}
	if !isValidContentType(contentType) {
		http.Error(w, "invalid content type", http.StatusBadRequest)
		return
	}

	// Validate module
	module := models.KBModule(req.Module)
	if module == "" {
		module = models.KBModuleGeneral
	}
	if !isValidModule(module) {
		http.Error(w, "invalid module", http.StatusBadRequest)
		return
	}

	// Validate lifecycle stage
	lifecycleStage := models.KBLifecycleStage(req.LifecycleStage)
	if lifecycleStage == "" {
		lifecycleStage = models.KBLifecycleSupport
	}
	if !isValidLifecycleStage(lifecycleStage) {
		http.Error(w, "invalid lifecycle stage", http.StatusBadRequest)
		return
	}

	now := time.Now().UTC()
	tags := req.Tags
	if tags == nil {
		tags = []string{}
	}

	article := models.KBArticle{
		ID:             store.NewID("kb"),
		TenantID:       tenant,
		Title:          strings.TrimSpace(req.Title),
		Slug:           slug,
		Summary:        strings.TrimSpace(req.Summary),
		Content:        req.Content,
		ContentType:    contentType,
		Module:         module,
		LifecycleStage: lifecycleStage,
		Tags:           tags,
		Version:        1,
		Status:         models.KBArticleStatusDraft,
		CreatedByID:    userID,
		CreatedByName:  userName,
		UpdatedByID:    userID,
		UpdatedByName:  userName,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := h.pg.KBArticles().Create(r.Context(), article); err != nil {
		h.log.Error("failed to create kb article", zap.Error(err))
		http.Error(w, "failed to create article", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, article)
}

func (h *KBArticlesHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	tenant := middleware.TenantID(r.Context())

	article, err := h.pg.KBArticles().GetByID(r.Context(), tenant, id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, article)
}

func (h *KBArticlesHandler) GetBySlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	tenant := middleware.TenantID(r.Context())

	article, err := h.pg.KBArticles().GetBySlug(r.Context(), tenant, slug)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, article)
}

type updateKBArticleReq struct {
	Title          *string   `json:"title"`
	Slug           *string   `json:"slug"`
	Summary        *string   `json:"summary"`
	Content        *string   `json:"content"`
	ContentType    *string   `json:"contentType"`
	Module         *string   `json:"module"`
	LifecycleStage *string   `json:"lifecycleStage"`
	Tags           *[]string `json:"tags"`
}

func (h *KBArticlesHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	tenant := middleware.TenantID(r.Context())
	userID := middleware.UserID(r.Context())
	userName := middleware.UserName(r.Context())

	// Get existing article
	article, err := h.pg.KBArticles().GetByID(r.Context(), tenant, id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	// Parse update request
	var req updateKBArticleReq
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
			// Check if new slug exists
			exists, err := h.pg.KBArticles().SlugExists(r.Context(), tenant, newSlug, article.ID)
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
		ct := models.KBContentType(*req.ContentType)
		if !isValidContentType(ct) {
			http.Error(w, "invalid content type", http.StatusBadRequest)
			return
		}
		article.ContentType = ct
	}
	if req.Module != nil {
		m := models.KBModule(*req.Module)
		if !isValidModule(m) {
			http.Error(w, "invalid module", http.StatusBadRequest)
			return
		}
		article.Module = m
	}
	if req.LifecycleStage != nil {
		ls := models.KBLifecycleStage(*req.LifecycleStage)
		if !isValidLifecycleStage(ls) {
			http.Error(w, "invalid lifecycle stage", http.StatusBadRequest)
			return
		}
		article.LifecycleStage = ls
	}
	if req.Tags != nil {
		article.Tags = *req.Tags
	}

	article.Version++
	article.UpdatedByID = userID
	article.UpdatedByName = userName
	article.UpdatedAt = time.Now().UTC()

	if err := h.pg.KBArticles().Update(r.Context(), article); err != nil {
		h.log.Error("failed to update kb article", zap.Error(err))
		http.Error(w, "failed to update article", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, article)
}

func (h *KBArticlesHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	tenant := middleware.TenantID(r.Context())

	if err := h.pg.KBArticles().Delete(r.Context(), tenant, id); err != nil {
		if err.Error() == "not found" {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		h.log.Error("failed to delete kb article", zap.Error(err))
		http.Error(w, "failed to delete article", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *KBArticlesHandler) List(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	contentType := strings.TrimSpace(r.URL.Query().Get("contentType"))
	module := strings.TrimSpace(r.URL.Query().Get("module"))
	lifecycleStage := strings.TrimSpace(r.URL.Query().Get("lifecycleStage"))
	status := strings.TrimSpace(r.URL.Query().Get("status"))
	limit := parseLimit(r.URL.Query().Get("limit"), 50, 200)
	curT, curID, hasCur := decodeCursor(strings.TrimSpace(r.URL.Query().Get("cursor")))

	items, next, err := h.pg.KBArticles().List(r.Context(), models.KBArticleListParams{
		TenantID:       tenant,
		ContentType:    contentType,
		Module:         module,
		LifecycleStage: lifecycleStage,
		Status:         status,
		Query:          q,
		Limit:          limit,
		HasCursor:      hasCur,
		CursorTime:     curT,
		CursorID:       curID,
	})
	if err != nil {
		h.log.Error("failed to list kb articles", zap.Error(err))
		http.Error(w, "failed to list articles", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items, "nextCursor": next})
}

func (h *KBArticlesHandler) Search(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	if q == "" {
		http.Error(w, "query parameter q is required", http.StatusBadRequest)
		return
	}

	limit := parseLimit(r.URL.Query().Get("limit"), 20, 100)
	items, err := h.pg.KBArticles().Search(r.Context(), tenant, q, limit)
	if err != nil {
		h.log.Error("failed to search kb articles", zap.Error(err))
		http.Error(w, "failed to search articles", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

func (h *KBArticlesHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())
	stats, err := h.pg.KBArticles().GetStats(r.Context(), tenant)
	if err != nil {
		h.log.Error("failed to get kb stats", zap.Error(err))
		http.Error(w, "failed to get stats", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, stats)
}

func (h *KBArticlesHandler) Publish(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	tenant := middleware.TenantID(r.Context())
	userID := middleware.UserID(r.Context())
	userName := middleware.UserName(r.Context())

	if err := h.pg.KBArticles().Publish(r.Context(), tenant, id, userID, userName); err != nil {
		if err.Error() == "not found or not a draft" {
			http.Error(w, "article not found or not a draft", http.StatusNotFound)
			return
		}
		h.log.Error("failed to publish kb article", zap.Error(err))
		http.Error(w, "failed to publish article", http.StatusInternalServerError)
		return
	}

	// Return the updated article
	article, err := h.pg.KBArticles().GetByID(r.Context(), tenant, id)
	if err != nil {
		h.log.Error("failed to get published article", zap.Error(err))
		http.Error(w, "article published but failed to retrieve", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, article)
}

// Helper functions

func generateSlug(title string) string {
	slug := strings.ToLower(strings.TrimSpace(title))
	// Replace spaces with dashes
	slug = strings.ReplaceAll(slug, " ", "-")
	// Remove non-alphanumeric characters except dashes
	reg := regexp.MustCompile(`[^a-z0-9-]`)
	slug = reg.ReplaceAllString(slug, "")
	// Remove multiple consecutive dashes
	reg = regexp.MustCompile(`-+`)
	slug = reg.ReplaceAllString(slug, "-")
	// Trim leading/trailing dashes
	slug = strings.Trim(slug, "-")
	return slug
}

func isValidContentType(ct models.KBContentType) bool {
	for _, valid := range models.ValidKBContentTypes() {
		if ct == valid {
			return true
		}
	}
	return false
}

func isValidModule(m models.KBModule) bool {
	for _, valid := range models.ValidKBModules() {
		if m == valid {
			return true
		}
	}
	return false
}

func isValidLifecycleStage(ls models.KBLifecycleStage) bool {
	for _, valid := range models.ValidKBLifecycleStages() {
		if ls == valid {
			return true
		}
	}
	return false
}
