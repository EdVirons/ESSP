APP?=ssot-hr
HTTP_ADDR?=:8300
DB_URL?=postgres://ssp:ssp@localhost:5432/ssp_hr?sslmode=disable

.PHONY: tidy test run build migrate-up migrate-down seed

tidy:
	go mod tidy

test:
	go test ./...

run:
	APP_ENV=dev HTTP_ADDR=$(HTTP_ADDR) DB_URL=$(DB_URL) go run ./cmd/ssot_hr

build:
	go build -o bin/$(APP) ./cmd/ssot_hr

migrate-up:
	psql $(DB_URL) -f migrations/001_init.sql

migrate-down:
	psql $(DB_URL) -c "DROP TABLE IF EXISTS hr_outbox_events, hr_audit_log, team_memberships, teams, people, org_units CASCADE;"

seed: ## Seed demo data into HR directory
	@./seed/seed-demo.sh
