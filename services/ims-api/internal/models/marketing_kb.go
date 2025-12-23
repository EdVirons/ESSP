package models

import "time"

// MKBContentType represents the type of marketing KB article
type MKBContentType string

const (
	MKBContentTypeMessaging MKBContentType = "messaging"
	MKBContentTypeCaseStudy MKBContentType = "case_study"
	MKBContentTypeDeck      MKBContentType = "deck"
	MKBContentTypeObjection MKBContentType = "objection"
	MKBContentTypeROI       MKBContentType = "roi"
)

// MKBPersona represents the target persona for the content
type MKBPersona string

const (
	MKBPersonaDirector       MKBPersona = "director"
	MKBPersonaPrincipal      MKBPersona = "principal"
	MKBPersonaTeacher        MKBPersona = "teacher"
	MKBPersonaParent         MKBPersona = "parent"
	MKBPersonaITAdmin        MKBPersona = "it_admin"
	MKBPersonaCountyOfficial MKBPersona = "county_official"
)

// MKBContextTag represents school context tags
type MKBContextTag string

const (
	MKBContextRural           MKBContextTag = "rural"
	MKBContextUrban           MKBContextTag = "urban"
	MKBContextLowConnectivity MKBContextTag = "low_connectivity"
	MKBContextNoISP           MKBContextTag = "no_isp"
	MKBContextConnected       MKBContextTag = "connected"
	MKBContextPrivate         MKBContextTag = "private"
	MKBContextPublic          MKBContextTag = "public"
	MKBContextCBC             MKBContextTag = "cbc"
	MKBContextIGCSE           MKBContextTag = "igcse"
	MKBContext844             MKBContextTag = "8-4-4"
)

// MKBArticleStatus represents the publication status
type MKBArticleStatus string

const (
	MKBStatusDraft    MKBArticleStatus = "draft"
	MKBStatusReview   MKBArticleStatus = "review"
	MKBStatusApproved MKBArticleStatus = "approved"
	MKBStatusArchived MKBArticleStatus = "archived"
)

// MKBArticle represents a marketing KB article
type MKBArticle struct {
	ID             string           `json:"id"`
	TenantID       string           `json:"tenantId"`
	Title          string           `json:"title"`
	Slug           string           `json:"slug"`
	Summary        string           `json:"summary"`
	Content        string           `json:"content"`
	ContentType    MKBContentType   `json:"contentType"`
	Personas       []string         `json:"personas"`
	ContextTags    []string         `json:"contextTags"`
	Tags           []string         `json:"tags"`
	Version        int              `json:"version"`
	Status         MKBArticleStatus `json:"status"`
	UsageCount     int              `json:"usageCount"`
	LastUsedAt     *time.Time       `json:"lastUsedAt,omitempty"`
	CreatedByID    string           `json:"createdById"`
	CreatedByName  string           `json:"createdByName"`
	UpdatedByID    string           `json:"updatedById"`
	UpdatedByName  string           `json:"updatedByName"`
	ApprovedAt     *time.Time       `json:"approvedAt,omitempty"`
	ApprovedByID   string           `json:"approvedById,omitempty"`
	ApprovedByName string           `json:"approvedByName,omitempty"`
	CreatedAt      time.Time        `json:"createdAt"`
	UpdatedAt      time.Time        `json:"updatedAt"`
}

// PitchKit represents a saved collection of marketing articles
type PitchKit struct {
	ID            string       `json:"id"`
	TenantID      string       `json:"tenantId"`
	Name          string       `json:"name"`
	Description   string       `json:"description"`
	TargetPersona string       `json:"targetPersona"`
	ContextTags   []string     `json:"contextTags"`
	ArticleIDs    []string     `json:"articleIds"`
	Articles      []MKBArticle `json:"articles,omitempty"`
	IsTemplate    bool         `json:"isTemplate"`
	CreatedByID   string       `json:"createdById"`
	CreatedByName string       `json:"createdByName"`
	UpdatedByID   string       `json:"updatedById"`
	UpdatedByName string       `json:"updatedByName"`
	CreatedAt     time.Time    `json:"createdAt"`
	UpdatedAt     time.Time    `json:"updatedAt"`
}

// MKBArticleListParams holds parameters for listing marketing KB articles
type MKBArticleListParams struct {
	TenantID    string
	ContentType string
	Persona     string
	ContextTag  string
	Status      string
	Query       string
	Limit       int
	HasCursor   bool
	CursorTime  time.Time
	CursorID    string
}

// PitchKitListParams holds parameters for listing pitch kits
type PitchKitListParams struct {
	TenantID      string
	TargetPersona string
	IsTemplate    *bool
	Limit         int
	HasCursor     bool
	CursorTime    time.Time
	CursorID      string
}

// MKBStats holds aggregate statistics for marketing KB articles
type MKBStats struct {
	Total         int            `json:"total"`
	Approved      int            `json:"approved"`
	InReview      int            `json:"inReview"`
	Draft         int            `json:"draft"`
	ByContentType map[string]int `json:"byContentType"`
	ByPersona     map[string]int `json:"byPersona"`
	ByContextTag  map[string]int `json:"byContextTag"`
}

// ValidMKBContentTypes returns all valid content types
func ValidMKBContentTypes() []MKBContentType {
	return []MKBContentType{
		MKBContentTypeMessaging,
		MKBContentTypeCaseStudy,
		MKBContentTypeDeck,
		MKBContentTypeObjection,
		MKBContentTypeROI,
	}
}

// ValidMKBPersonas returns all valid personas
func ValidMKBPersonas() []MKBPersona {
	return []MKBPersona{
		MKBPersonaDirector,
		MKBPersonaPrincipal,
		MKBPersonaTeacher,
		MKBPersonaParent,
		MKBPersonaITAdmin,
		MKBPersonaCountyOfficial,
	}
}

// ValidMKBContextTags returns all valid context tags
func ValidMKBContextTags() []MKBContextTag {
	return []MKBContextTag{
		MKBContextRural,
		MKBContextUrban,
		MKBContextLowConnectivity,
		MKBContextNoISP,
		MKBContextConnected,
		MKBContextPrivate,
		MKBContextPublic,
		MKBContextCBC,
		MKBContextIGCSE,
		MKBContext844,
	}
}

// ValidMKBStatuses returns all valid statuses
func ValidMKBStatuses() []MKBArticleStatus {
	return []MKBArticleStatus{
		MKBStatusDraft,
		MKBStatusReview,
		MKBStatusApproved,
		MKBStatusArchived,
	}
}
