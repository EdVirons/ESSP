package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/edvirons/ssp/ims/internal/metrics"
	chimw "github.com/go-chi/chi/v5/middleware"
)

// MetricsMiddleware returns a middleware that records HTTP metrics using Prometheus
func MetricsMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := chimw.NewWrapResponseWriter(w, r.ProtoMajor)

			next.ServeHTTP(ww, r)

			duration := time.Since(start)
			status := strconv.Itoa(ww.Status())

			// Record metrics
			metrics.HTTPRequestsTotal.WithLabelValues(r.Method, r.URL.Path, status).Inc()
			metrics.HTTPRequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration.Seconds())
		})
	}
}
