package admin

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

// ServiceHealth represents the health status of a single service
type ServiceHealth struct {
	Name      string `json:"name"`
	Status    string `json:"status"` // healthy, degraded, unhealthy
	LatencyMs int64  `json:"latencyMs"`
	LastCheck string `json:"lastCheck"`
}

// HealthResponse is the response for /admin/v1/health/services
type HealthResponse struct {
	Services []ServiceHealth `json:"services"`
}

// GetServicesHealth returns aggregated health status of all ESSP services
func (h *Handler) GetServicesHealth(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	// Services to check
	services := []struct {
		name     string
		endpoint string
	}{
		{"ims-api", "http://localhost:8080/healthz"},
		{"ssot-school", "http://localhost:8081/healthz"},
		{"ssot-devices", "http://localhost:8082/healthz"},
		{"ssot-parts", "http://localhost:8083/healthz"},
		{"sync-worker", "http://localhost:8084/healthz"},
	}

	var wg sync.WaitGroup
	results := make([]ServiceHealth, len(services))

	for i, svc := range services {
		wg.Add(1)
		go func(idx int, name, endpoint string) {
			defer wg.Done()
			results[idx] = checkServiceHealth(ctx, name, endpoint)
		}(i, svc.name, svc.endpoint)
	}

	wg.Wait()

	resp := HealthResponse{Services: results}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func checkServiceHealth(ctx context.Context, name, endpoint string) ServiceHealth {
	start := time.Now()
	status := "healthy"

	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return ServiceHealth{
			Name:      name,
			Status:    "unhealthy",
			LatencyMs: 0,
			LastCheck: time.Now().UTC().Format(time.RFC3339),
		}
	}

	resp, err := client.Do(req)
	latency := time.Since(start).Milliseconds()

	if err != nil {
		status = "unhealthy"
	} else {
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			status = "degraded"
		}
		// Check latency threshold
		if latency > 500 {
			status = "degraded"
		}
	}

	return ServiceHealth{
		Name:      name,
		Status:    status,
		LatencyMs: latency,
		LastCheck: time.Now().UTC().Format(time.RFC3339),
	}
}
