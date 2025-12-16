package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/edvirons/ssp/ims/internal/audit"
	"github.com/edvirons/ssp/ims/internal/middleware"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// AuditLogsHandler handles audit log API requests
type AuditLogsHandler struct {
	log   *zap.Logger
	store *audit.Store
}

// NewAuditLogsHandler creates a new audit logs handler
func NewAuditLogsHandler(log *zap.Logger, auditStore *audit.Store) *AuditLogsHandler {
	return &AuditLogsHandler{
		log:   log,
		store: auditStore,
	}
}

// List retrieves audit logs with filters and pagination
func (h *AuditLogsHandler) List(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())

	// Parse query parameters
	entityType := strings.TrimSpace(r.URL.Query().Get("entityType"))
	entityID := strings.TrimSpace(r.URL.Query().Get("entityId"))
	userID := strings.TrimSpace(r.URL.Query().Get("userId"))
	action := strings.TrimSpace(r.URL.Query().Get("action"))
	startDateStr := strings.TrimSpace(r.URL.Query().Get("startDate"))
	endDateStr := strings.TrimSpace(r.URL.Query().Get("endDate"))
	limit := parseLimit(r.URL.Query().Get("limit"), 50, 200)

	curT, curID, hasCur := decodeCursor(strings.TrimSpace(r.URL.Query().Get("cursor")))

	// Parse dates
	var startDate, endDate time.Time
	var err error

	if startDateStr != "" {
		startDate, err = time.Parse(time.RFC3339, startDateStr)
		if err != nil {
			// Try date-only format
			startDate, err = time.Parse("2006-01-02", startDateStr)
			if err != nil {
				http.Error(w, "invalid startDate format, use RFC3339 or YYYY-MM-DD", http.StatusBadRequest)
				return
			}
		}
	}

	if endDateStr != "" {
		endDate, err = time.Parse(time.RFC3339, endDateStr)
		if err != nil {
			// Try date-only format
			endDate, err = time.Parse("2006-01-02", endDateStr)
			if err != nil {
				http.Error(w, "invalid endDate format, use RFC3339 or YYYY-MM-DD", http.StatusBadRequest)
				return
			}
			// If using date-only format, set to end of day
			endDate = endDate.Add(24*time.Hour - time.Second)
		}
	}

	items, next, err := h.store.List(r.Context(), audit.ListParams{
		TenantID:        tenant,
		EntityType:      entityType,
		EntityID:        entityID,
		UserID:          userID,
		Action:          action,
		StartDate:       startDate,
		EndDate:         endDate,
		Limit:           limit,
		CursorCreatedAt: curT,
		CursorID:        curID,
		HasCursor:       hasCur,
	})
	if err != nil {
		h.log.Error("failed to list audit logs", zap.Error(err))
		http.Error(w, "failed to list audit logs", http.StatusInternalServerError)
		return
	}

	resp := map[string]any{
		"items":      items,
		"nextCursor": next,
	}
	writeJSON(w, http.StatusOK, resp)
}

// GetByID retrieves a single audit log by ID
func (h *AuditLogsHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	tenant := middleware.TenantID(r.Context())

	log, err := h.store.GetByID(r.Context(), tenant, id)
	if err != nil {
		http.Error(w, "audit log not found", http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, log)
}
