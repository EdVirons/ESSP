package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/edvirons/ssp/ssot_parts/internal/api"
	"go.uber.org/zap"
)

func main() {
	log, _ := zap.NewProduction()
	defer log.Sync()

	addr := env("HTTP_ADDR", ":8083")
	if os.Getenv("DB_URL") == "" {
		_ = os.Setenv("DEFAULT_DB_URL", "postgres://ssp:ssp@postgres:5432/ssp_parts?sslmode=disable")
	}
	srv := &http.Server{Addr: addr, Handler: api.NewServer(log), ReadHeaderTimeout: 5 * time.Second}

	go func() {
		log.Info("ssot-parts starting", zap.String("addr", addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("server error", zap.Error(err))
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
	log.Info("shutdown complete")
	_ = ctx
}

func env(k, d string) string {
	if v := os.Getenv(k); v != "" { return v }
	return d
}
