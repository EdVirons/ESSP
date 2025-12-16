APP?=app
DB_URL?=postgres://ssp:ssp@localhost:5432/ssp_devices?sslmode=disable

.PHONY: run test tidy migrate-up migrate-down

run:
	go run ./cmd/ssot_devices

test:
	go test ./...

tidy:
	go mod tidy

migrate-up:
	go run ./cmd/migrate up

migrate-down:
	go run ./cmd/migrate down
