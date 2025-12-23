package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/edvirons/ssp/ims/internal/config"
	"github.com/edvirons/ssp/ims/internal/middleware"
	"github.com/edvirons/ssp/ims/internal/models"
	"github.com/edvirons/ssp/ims/internal/service"
	"github.com/edvirons/ssp/ims/internal/store"
	"go.uber.org/zap"
)

type TelemetryHandler struct {
	cfg config.Config
	log *zap.Logger
	pg  *store.Postgres
}

func NewTelemetryHandler(cfg config.Config, log *zap.Logger, pg *store.Postgres) *TelemetryHandler {
	return &TelemetryHandler{cfg: cfg, log: log, pg: pg}
}

type telemetryEvent struct {
	DeviceID   string `json:"deviceId"`
	Type       string `json:"type"`       // e.g. "policy_breach", "crash", "offline"
	Message    string `json:"message"`    // details
	Severity   string `json:"severity"`   // low|medium|high|criticalgritical (typos tolerated)
	ReportedBy string `json:"reportedBy"` // e.g. "nexus-mdm"
}

func (h *TelemetryHandler) Ingest(w http.ResponseWriter, r *http.Request) {
	var e telemetryEvent
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(e.DeviceID) == "" || strings.TrimSpace(e.Type) == "" {
		http.Error(w, "deviceId and type are required", http.StatusBadRequest)
		return
	}

	tenant := middleware.TenantID(r.Context())
	school := middleware.SchoolID(r.Context())

	sev := parseSeverity(e.Severity)
	now := time.Now().UTC()
	inc := models.Incident{
		ID:          store.NewID("inc"),
		TenantID:    tenant,
		SchoolID:    school,
		DeviceID:    strings.TrimSpace(e.DeviceID),
		Category:    "telemetry:" + strings.TrimSpace(e.Type),
		Severity:    sev,
		Status:      models.IncidentNew,
		Title:       "Telemetry: " + strings.TrimSpace(e.Type),
		Description: strings.TrimSpace(e.Message),
		ReportedBy:  firstNonEmpty(strings.TrimSpace(e.ReportedBy), "nexus-mdm"),
		SLADueAt:    service.SLADue(sev, now),
		SLABreached: false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := h.pg.Incidents().Create(r.Context(), inc); err != nil {
		http.Error(w, "failed to create incident", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, inc)
}

func parseSeverity(s string) models.Severity {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "critical", "critial", "crit":
		return models.SeverityCritical
	case "high":
		return models.SeverityHigh
	case "medium", "med":
		return models.SeverityMedium
	default:
		return models.SeverityLow
	}
}

func firstNonEmpty(a, b string) string {
	if a != "" {
		return a
	}
	return b
}
