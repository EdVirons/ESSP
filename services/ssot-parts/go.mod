module github.com/edvirons/ssp/ssot_parts

go 1.22

require (
	github.com/edvirons/ssp/shared v0.0.0-00010101000000-000000000000
	github.com/go-chi/chi/v5 v5.0.12
	github.com/jackc/pgx/v5 v5.6.0
	github.com/nats-io/nats.go v1.35.0
	go.uber.org/zap v1.27.0
)

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/klauspost/compress v1.17.2 // indirect
	github.com/nats-io/nkeys v0.4.7 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/oklog/ulid/v2 v2.1.0 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	golang.org/x/crypto v0.18.0 // indirect
	golang.org/x/sync v0.1.0 // indirect
	golang.org/x/sys v0.16.0 // indirect
	golang.org/x/text v0.14.0 // indirect
)

replace github.com/edvirons/ssp/shared => ../../shared
