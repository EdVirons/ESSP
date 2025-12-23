// Package health provides a health check HTTP server.
package health

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// StartServer starts the health check HTTP server on the given port.
func StartServer(log *zap.Logger, db *pgxpool.Pool, port string) {
	r := chi.NewRouter()

	r.Get("/health", healthHandler(log, db))

	addr := ":" + port
	log.Info("health check server starting", zap.String("addr", addr))
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Error("health server failed", zap.Error(err))
	}
}

// healthHandler returns an HTTP handler that checks database connectivity.
func healthHandler(log *zap.Logger, db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		// Check database connectivity
		if err := db.Ping(ctx); err != nil {
			log.Error("health check failed: database ping", zap.Error(err))
			w.WriteHeader(http.StatusServiceUnavailable)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"status": "unhealthy",
				"error":  "database unreachable",
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"status": "healthy",
		})
	}
}
