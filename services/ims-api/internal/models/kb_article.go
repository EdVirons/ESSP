package models

import "time"

// KBContentType represents the type of KB article
type KBContentType string

const (
	KBContentTypeRunbook         KBContentType = "runbook"
	KBContentTypeTroubleshooting KBContentType = "troubleshooting"
	KBContentTypeKEDB            KBContentType = "kedb"
	KBContentTypeChecklist       KBContentType = "checklist"
	KBContentTypeSOP             KBContentType = "sop"
)

// KBModule represents the ESSP module the article relates to
type KBModule string

const (
	KBModuleLearningPortal KBModule = "learning_portal"
	KBModuleMDM            KBModule = "mdm"
	KBModuleSSO            KBModule = "sso"
	KBModuleDevices        KBModule = "devices"
	KBModuleInventory      KBModule = "inventory"
	KBModuleGeneral        KBModule = "general"
)

// KBLifecycleStage represents when the article is most relevant
type KBLifecycleStage string

const (
	KBLifecycleDemo       KBLifecycleStage = "demo"
	KBLifecycleInstall    KBLifecycleStage = "install"
	KBLifecycleCommission KBLifecycleStage = "commission"
	KBLifecycleSupport    KBLifecycleStage = "support"
)

// KBArticleStatus represents the publication status
type KBArticleStatus string

const (
	KBArticleStatusDraft     KBArticleStatus = "draft"
	KBArticleStatusPublished KBArticleStatus = "published"
	KBArticleStatusArchived  KBArticleStatus = "archived"
)

// KBArticle represents a knowledge base article
type KBArticle struct {
	ID             string           `json:"id"`
	TenantID       string           `json:"tenantId"`
	Title          string           `json:"title"`
	Slug           string           `json:"slug"`
	Summary        string           `json:"summary"`
	Content        string           `json:"content"`
	ContentType    KBContentType    `json:"contentType"`
	Module         KBModule         `json:"module"`
	LifecycleStage KBLifecycleStage `json:"lifecycleStage"`
	Tags           []string         `json:"tags"`
	Version        int              `json:"version"`
	Status         KBArticleStatus  `json:"status"`
	CreatedByID    string           `json:"createdById"`
	CreatedByName  string           `json:"createdByName"`
	UpdatedByID    string           `json:"updatedById"`
	UpdatedByName  string           `json:"updatedByName"`
	PublishedAt    *time.Time       `json:"publishedAt,omitempty"`
	CreatedAt      time.Time        `json:"createdAt"`
	UpdatedAt      time.Time        `json:"updatedAt"`
}

// KBArticleListParams holds parameters for listing KB articles
type KBArticleListParams struct {
	TenantID       string
	ContentType    string
	Module         string
	LifecycleStage string
	Status         string
	Query          string
	Limit          int
	HasCursor      bool
	CursorTime     time.Time
	CursorID       string
}

// KBStats holds aggregate statistics for KB articles
type KBStats struct {
	Total            int            `json:"total"`
	Published        int            `json:"published"`
	Draft            int            `json:"draft"`
	ByContentType    map[string]int `json:"byContentType"`
	ByModule         map[string]int `json:"byModule"`
	ByLifecycleStage map[string]int `json:"byLifecycleStage"`
}

// ValidKBContentTypes returns all valid content types
func ValidKBContentTypes() []KBContentType {
	return []KBContentType{
		KBContentTypeRunbook,
		KBContentTypeTroubleshooting,
		KBContentTypeKEDB,
		KBContentTypeChecklist,
		KBContentTypeSOP,
	}
}

// ValidKBModules returns all valid modules
func ValidKBModules() []KBModule {
	return []KBModule{
		KBModuleLearningPortal,
		KBModuleMDM,
		KBModuleSSO,
		KBModuleDevices,
		KBModuleInventory,
		KBModuleGeneral,
	}
}

// ValidKBLifecycleStages returns all valid lifecycle stages
func ValidKBLifecycleStages() []KBLifecycleStage {
	return []KBLifecycleStage{
		KBLifecycleDemo,
		KBLifecycleInstall,
		KBLifecycleCommission,
		KBLifecycleSupport,
	}
}
