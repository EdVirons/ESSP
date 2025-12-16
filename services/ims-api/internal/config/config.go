package config

import (
	"log"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	AppEnv   string
	HTTPAddr string
	LogLevel string

	PGDSN string

	NATSURL string
	SSOTSchoolURL string
	SSOTDevicesURL string
	SSOTPartsURL string

	RedisAddr     string
	RedisPassword string
	RedisDB       int

	AuthEnabled  bool
	AuthIssuer   string
	AuthJWKSURL  string
	AuthAudience string

	TenantHeader string
	SchoolHeader string
	DevTenantID  string
	DevSchoolID  string

	CORSAllowedOrigins string

	AttachmentsPublicBaseURL string
	AttachmentsBucket        string

	MinIOEndpoint string
	MinIOAccessKey string
	MinIOSecretKey string
	MinIOUseSSL bool
	MinIORegion string
	MinIOPresignExpirySeconds int

	AutoRouteWorkOrders bool
	DefaultRepairLocation string

	SchoolSSOTBaseURL string
	DeviceSSOTBaseURL string
	PartsSSOTBaseURL string
	SSOTSyncPageSize int

	RateLimitEnabled  bool
	RateLimitReadRPM  int
	RateLimitWriteRPM int
	RateLimitBurst    int

	// Admin dashboard authentication
	AdminUsername   string
	AdminPassword   string
	AdminJWTExpiry  int // in hours
	AdminCookieSecure bool

	// Claude AI Support
	ClaudeAPIKey          string
	ClaudeModel           string
	ClaudeMaxTokens       int
	ClaudeTimeoutSeconds  int
	AIEnabled             bool
	AIMaxTurns            int
	AIFrustrationThreshold float64
}

func MustLoad() Config {
	c := Config{
		AppEnv:   getenv("APP_ENV", "dev"),
		HTTPAddr: getenv("HTTP_ADDR", ":8080"),
		LogLevel: getenv("LOG_LEVEL", "info"),

		PGDSN: getenvWithFallbacks("PG_DSN", "DATABASE_URL", "postgres://edvirons:edvirons@localhost:5432/ims?sslmode=disable"),

		RedisAddr:     getenvWithFallbacksForRedis("REDIS_ADDR", "REDIS_URL", "localhost:6379"),
		RedisPassword: getenv("REDIS_PASSWORD", ""),
		RedisDB:       mustAtoi(getenv("REDIS_DB", "0")),

		AuthEnabled:  mustAtob(getenv("AUTH_ENABLED", "false")),
		AuthIssuer:   getenv("AUTH_ISSUER", ""),
		AuthJWKSURL:  getenv("AUTH_JWKS_URL", ""),
		AuthAudience: getenv("AUTH_AUDIENCE", "ims-service"),

		TenantHeader: getenv("TENANT_HEADER", "X-Tenant-Id"),
		SchoolHeader: getenv("SCHOOL_HEADER", "X-School-Id"),
		DevTenantID:  getenv("DEV_TENANT_ID", "demo-tenant"),
		DevSchoolID:  getenv("DEV_SCHOOL_ID", "demo-school"),

		CORSAllowedOrigins: getenv("CORS_ALLOWED_ORIGINS", "*"),

		AttachmentsPublicBaseURL: getenv("ATTACHMENTS_PUBLIC_BASE_URL", "http://localhost:9000"),
		AttachmentsBucket:        getenv("ATTACHMENTS_BUCKET", "edvirons-ims"),

		MinIOEndpoint: getenv("MINIO_ENDPOINT", "localhost:9000"),
		MinIOAccessKey: getenv("MINIO_ACCESS_KEY", "minioadmin"),
		MinIOSecretKey: getenv("MINIO_SECRET_KEY", "minioadmin"),
		MinIOUseSSL: mustAtob(getenv("MINIO_USE_SSL", "false")),
		MinIORegion: getenv("MINIO_REGION", "us-east-1"),
		MinIOPresignExpirySeconds: mustAtoi(getenv("MINIO_PRESIGN_EXPIRY_SECONDS", "900")),

		AutoRouteWorkOrders: mustAtob(getenv("AUTO_ROUTE_WORK_ORDERS", "true")),
		DefaultRepairLocation: getenv("DEFAULT_REPAIR_LOCATION", "service_shop"),

		SchoolSSOTBaseURL: getenv("SCHOOL_SSOT_BASE_URL", ""),
		DeviceSSOTBaseURL: getenv("DEVICE_SSOT_BASE_URL", ""),
		PartsSSOTBaseURL: getenv("PARTS_SSOT_BASE_URL", ""),
		SSOTSyncPageSize: mustAtoi(getenv("SSOT_SYNC_PAGE_SIZE", "500")),

		RateLimitEnabled:  mustAtob(getenv("RATE_LIMIT_ENABLED", "true")),
		RateLimitReadRPM:  mustAtoi(getenv("RATE_LIMIT_READ_RPM", "300")),
		RateLimitWriteRPM: mustAtoi(getenv("RATE_LIMIT_WRITE_RPM", "100")),
		RateLimitBurst:    mustAtoi(getenv("RATE_LIMIT_BURST", "50")),

		AdminUsername:     getenv("ADMIN_USERNAME", "admin"),
		AdminPassword:     getenv("ADMIN_PASSWORD", "admin"),
		AdminJWTExpiry:    mustAtoi(getenv("ADMIN_JWT_EXPIRY_HOURS", "24")),
		AdminCookieSecure: mustAtob(getenv("ADMIN_COOKIE_SECURE", "false")),

		ClaudeAPIKey:          getenv("CLAUDE_API_KEY", ""),
		ClaudeModel:           getenv("CLAUDE_MODEL", "claude-sonnet-4-20250514"),
		ClaudeMaxTokens:       mustAtoi(getenv("CLAUDE_MAX_TOKENS", "1024")),
		ClaudeTimeoutSeconds:  mustAtoi(getenv("CLAUDE_TIMEOUT_SECONDS", "30")),
		AIEnabled:             mustAtob(getenv("AI_SUPPORT_ENABLED", "true")),
		AIMaxTurns:            mustAtoi(getenv("AI_MAX_TURNS", "10")),
		AIFrustrationThreshold: mustAtof(getenv("AI_FRUSTRATION_THRESHOLD", "0.7")),
	}
	if c.PGDSN == "" {
		log.Fatal("PG_DSN is required")
	}
	return c
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func mustAtoi(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf("invalid int: %q", s)
	}
	return n
}

func mustAtob(s string) bool {
	b, err := strconv.ParseBool(s)
	if err != nil {
		log.Fatalf("invalid bool: %q", s)
	}
	return b
}

func mustAtof(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		log.Fatalf("invalid float: %q", s)
	}
	return f
}

// getenvWithFallbacks checks multiple env var names, returning the first non-empty value
func getenvWithFallbacks(primary, fallback, def string) string {
	if v := os.Getenv(primary); v != "" {
		return v
	}
	if v := os.Getenv(fallback); v != "" {
		return v
	}
	return def
}

// getenvWithFallbacksForRedis handles REDIS_URL format (redis://host:port/db) and extracts host:port
func getenvWithFallbacksForRedis(primary, fallback, def string) string {
	if v := os.Getenv(primary); v != "" {
		return v
	}
	if v := os.Getenv(fallback); v != "" {
		// Parse redis://host:port/db format
		v = strings.TrimPrefix(v, "redis://")
		v = strings.TrimPrefix(v, "rediss://")
		// Remove /db suffix if present
		if idx := strings.Index(v, "/"); idx != -1 {
			v = v[:idx]
		}
		return v
	}
	return def
}
