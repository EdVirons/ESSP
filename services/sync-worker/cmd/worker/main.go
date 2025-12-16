package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/edvirons/ssp/sync_worker/internal/config"
	"github.com/edvirons/ssp/sync_worker/internal/health"
	"github.com/edvirons/ssp/sync_worker/internal/worker"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

func main() {
	log, _ := zap.NewProduction()
	defer log.Sync()

	cfg := config.Load()

	// Connect to PostgreSQL
	db, err := pgxpool.New(context.Background(), cfg.PGDSN)
	if err != nil {
		log.Fatal("database connect failed", zap.Error(err))
	}
	defer db.Close()

	// Verify database connection
	if err := db.Ping(context.Background()); err != nil {
		log.Fatal("database ping failed", zap.Error(err))
	}
	log.Info("database connected")

	// Connect to NATS
	nc, err := nats.Connect(cfg.NATSURL)
	if err != nil {
		log.Fatal("nats connect failed", zap.Error(err))
	}
	defer nc.Drain()
	log.Info("nats connected", zap.String("url", cfg.NATSURL))

	// Create sync worker
	sw := worker.New(log, db, cfg)

	// Subscribe to SSOT change events
	for _, kind := range []string{"school", "devices", "parts"} {
		subject := fmt.Sprintf("ssot.%s.changed", kind)
		_, err := nc.Subscribe(subject, sw.HandleEvent(kind))
		if err != nil {
			log.Fatal("subscribe failed", zap.String("subject", subject), zap.Error(err))
		}
		log.Info("subscribed to NATS subject", zap.String("subject", subject))
	}

	// Start health check server
	go health.StartServer(log, db, cfg.HealthPort)

	log.Info("sync-worker running", zap.String("nats", cfg.NATSURL), zap.String("health_port", cfg.HealthPort))

	// Wait for shutdown signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Info("shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = ctx
	log.Info("shutdown complete")
}
