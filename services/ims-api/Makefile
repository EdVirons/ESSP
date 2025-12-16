APP?=app
DB_URL?=postgres://ssp:ssp@localhost:5432/ssp_ims?sslmode=disable
TEST_DB_URL?=postgres://ssp:ssp@localhost:5433/ssp_ims_test?sslmode=disable

.PHONY: run test test-unit test-integration test-coverage test-coverage-html tidy migrate-up migrate-down test-env-up test-env-down test-db-setup seed seed-clean

run:
	go run ./cmd/api

# Run all tests (unit tests by default, use short mode to skip integration tests)
test:
	go test -v -race -short ./...

# Run only unit tests (fast, no external dependencies)
test-unit:
	go test -v -race -short ./...

# Run integration tests (requires test environment)
test-integration:
	INTEGRATION_TEST=1 TEST_DB_DSN=$(TEST_DB_URL) go test -v -race ./...

# Generate coverage report
test-coverage:
	go test -v -race -short -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -func=coverage.out

# Generate and open HTML coverage report
test-coverage-html:
	go test -v -race -short -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Start test environment (PostgreSQL, Redis, NATS, MinIO)
test-env-up:
	docker-compose -f ../../deployments/docker/docker-compose.test.yml up -d
	@echo "Waiting for services to be ready..."
	@sleep 5
	@echo "Test environment is ready!"
	@echo "PostgreSQL: localhost:5433"
	@echo "Valkey/Redis: localhost:6380"
	@echo "NATS: localhost:4223"
	@echo "MinIO: localhost:9001"

# Stop and remove test environment
test-env-down:
	docker-compose -f ../../deployments/docker/docker-compose.test.yml down -v

# Set up test database schema (run migrations)
test-db-setup:
	DB_URL=$(TEST_DB_URL) go run ./cmd/migrate up

# Run tests with verbose output and show individual test results
test-verbose:
	go test -v -race ./...

# Run tests and generate a detailed report
test-report:
	go test -v -race -json ./... > test-report.json
	@echo "Test report generated: test-report.json"

# Benchmark tests
test-bench:
	go test -v -race -run=^$$ -bench=. -benchmem ./...

tidy:
	go mod tidy

migrate-up:
	go run ./cmd/migrate up

migrate-down:
	go run ./cmd/migrate down

# Seed database with realistic Kenyan test data
seed:
	DB_URL=$(DB_URL) go run ./cmd/seed

# Clean and re-seed database
seed-clean:
	DB_URL=$(DB_URL) go run ./cmd/seed --clean
