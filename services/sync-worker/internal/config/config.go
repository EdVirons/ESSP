// Package config provides configuration loading for the sync worker.
package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds all configuration for the sync worker.
type Config struct {
	NATSURL        string
	PGDSN          string
	SchoolURL      string
	DevicesURL     string
	PartsURL       string
	FetchTimeout   time.Duration
	HealthPort     string
	MaxRetries     int
	InitialBackoff time.Duration
}

// Load reads configuration from environment variables with defaults.
func Load() Config {
	return Config{
		NATSURL:        env("NATS_URL", "nats://localhost:4222"),
		PGDSN:          env("PG_DSN", "postgres://postgres:postgres@localhost:5432/ims?sslmode=disable"),
		SchoolURL:      env("SSOT_SCHOOL_URL", "http://ssot-school:8081"),
		DevicesURL:     env("SSOT_DEVICES_URL", "http://ssot-devices:8082"),
		PartsURL:       env("SSOT_PARTS_URL", "http://ssot-parts:8083"),
		FetchTimeout:   envDuration("SSOT_FETCH_TIMEOUT_SECONDS", 30*time.Second),
		HealthPort:     env("HEALTH_PORT", "8084"),
		MaxRetries:     envInt("MAX_RETRIES", 3),
		InitialBackoff: envDuration("INITIAL_BACKOFF_SECONDS", 1*time.Second),
	}
}

// env returns the value of an environment variable or a default.
func env(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}

// envInt returns an integer environment variable or a default.
func envInt(k string, d int) int {
	if v := os.Getenv(k); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return d
}

// envDuration returns a duration environment variable (in seconds) or a default.
func envDuration(k string, d time.Duration) time.Duration {
	if v := os.Getenv(k); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return time.Duration(i) * time.Second
		}
	}
	return d
}
