package store

import (
	"path"
	"strings"
	"time"

	"github.com/edvirons/ssp/ims/internal/models"
)

func ObjectKeyForAttachment(tenantID, schoolID string, entityType models.AttachmentEntityType, entityID string, now time.Time, fileName string) string {
	safeName := strings.ReplaceAll(strings.TrimSpace(fileName), " ", "_")
	datePrefix := now.UTC().Format("2006/01/02")
	return path.Join(
		"tenants", tenantID,
		"schools", schoolID,
		"attachments",
		string(entityType),
		entityID,
		datePrefix,
		now.UTC().Format("150405")+"_"+safeName,
	)
}
