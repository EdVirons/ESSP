module github.com/edvirons/ssp/ssot_devices

go 1.22

require (
	github.com/go-chi/chi/v5 v5.0.12
	github.com/jackc/pgx/v5 v5.6.0
	go.uber.org/zap v1.27.0
	github.com/nats-io/nats.go v1.35.0
)

replace github.com/edvirons/ssp/shared => ../../shared
