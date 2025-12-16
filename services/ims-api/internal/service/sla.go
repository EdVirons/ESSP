package service

import (
	"time"

	"github.com/edvirons/ssp/ims/internal/models"
)

func SLADue(sev models.Severity, now time.Time) time.Time {
	switch sev {
	case models.SeverityCritical:
		return now.Add(4 * time.Hour)
	case models.SeverityHigh:
		return now.Add(24 * time.Hour)
	case models.SeverityMedium:
		return now.Add(48 * time.Hour)
	default:
		return now.Add(72 * time.Hour)
	}
}
