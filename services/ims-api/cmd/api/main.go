package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/edvirons/ssp/ims/internal/api"
	"github.com/edvirons/ssp/ims/internal/config"
	"github.com/edvirons/ssp/ims/internal/jobs"
	"github.com/edvirons/ssp/ims/internal/logging"
	"github.com/edvirons/ssp/ims/internal/store"
	"github.com/edvirons/ssp/ims/internal/ssotcache"
)

func main() {
	cfg := config.MustLoad()

	logger := logging.New(cfg.LogLevel)
	defer func() { _ = logger.Sync() }()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pg, err := store.NewPostgres(ctx, cfg.PGDSN)
	if err != nil {
		logger.Fatal("postgres connection failed", logging.Err(err))
	}
	defer pg.Close()

	// SSOT consumer: listen to ssot.*.changed events and cache latest SSOT exports in IMS DB
	if err := ssotcache.Start(ctx, logger, pg.RawPool(), ssotcache.ConfigFromEnv()); err != nil {
		logger.Error("ssotcache failed to start", logging.Err(err))
	}

	rdb := store.NewValkey(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
	defer func() { _ = rdb.Close() }()

	// Background jobs
	j := jobs.NewScheduler(logger, pg)
	j.Start(ctx)
	defer j.Stop()

	srv := api.NewServer(cfg, logger, pg, rdb)

	httpServer := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           srv.Router(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		logger.Info("http server started", logging.Str("addr", cfg.HTTPAddr))
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("http server crashed", logging.Err(err))
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	logger.Info("shutdown signal received")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		logger.Error("http shutdown failed", logging.Err(err))
	}
	logger.Info("shutdown complete")
}
