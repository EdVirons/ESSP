# ESSP Quick Reference Guide

Quick command reference for common development tasks.

## First Time Setup

```bash
# 1. Clone and enter repository
git clone <repo-url> /home/pato/opt/ESSP
cd /home/pato/opt/ESSP

# 2. Install dependencies
make tidy

# 3. Start infrastructure and initialize databases
./scripts/dev-start.sh --init-db --with-tools

# 4. Start services (in separate terminals)
cd services/ims-api && go run ./cmd/api
cd services/ssot-school && go run ./cmd/ssot_school
cd services/ssot-devices && go run ./cmd/ssot_devices
cd services/ssot-parts && go run ./cmd/ssot_parts
cd services/sync-worker && go run ./cmd/worker
```

## Daily Development

```bash
# Start infrastructure
./scripts/dev-start.sh

# Start the service you're working on
cd services/ims-api && go run ./cmd/api
```

## Docker Commands

```bash
# Start all infrastructure
docker compose -f docker-compose.dev.yml up -d

# Start with management tools
docker compose -f docker-compose.dev.yml --profile tools up -d

# View logs
docker compose -f docker-compose.dev.yml logs -f
docker compose -f docker-compose.dev.yml logs -f postgres

# Check service status
docker compose -f docker-compose.dev.yml ps

# Stop services
docker compose -f docker-compose.dev.yml down

# Stop and remove all data (fresh start)
docker compose -f docker-compose.dev.yml down -v

# Restart specific service
docker compose -f docker-compose.dev.yml restart postgres
```

## Database Commands

```bash
# Initialize databases
./scripts/init-dev-db.sh

# Create databases manually
docker exec -it essp-postgres-dev psql -U ssp -c "CREATE DATABASE ssp_ims;"

# Connect to database
docker exec -it essp-postgres-dev psql -U ssp -d ssp_ims

# Run migrations (from service directory)
cd services/ims-api
make migrate-up
# or
DB_URL=postgres://ssp:ssp@localhost:5432/ssp_ims?sslmode=disable go run ./cmd/migrate up

# Drop and recreate database
docker exec -it essp-postgres-dev psql -U ssp -c "DROP DATABASE IF EXISTS ssp_ims;"
docker exec -it essp-postgres-dev psql -U ssp -c "CREATE DATABASE ssp_ims;"
cd services/ims-api && make migrate-up
```

## Running Services

```bash
# IMS API
cd services/ims-api
go run ./cmd/api
# or with environment variables
PG_DSN=postgres://ssp:ssp@localhost:5432/ssp_ims?sslmode=disable go run ./cmd/api

# SSOT School
cd services/ssot-school
go run ./cmd/ssot_school

# SSOT Devices
cd services/ssot-devices
go run ./cmd/ssot_devices

# SSOT Parts
cd services/ssot-parts
go run ./cmd/ssot_parts

# Sync Worker
cd services/sync-worker
go run ./cmd/worker

# Run with hot reload using Air
cd services/ims-api
air
```

## Build Commands

```bash
# Build all services
make build-all

# Build specific service
make build-ims-api
make build-ssot-school

# Build from service directory
cd services/ims-api
go build -o bin/ims-api ./cmd/api

# Build with make
cd services/ims-api
make run
```

## Testing Commands

```bash
# Run all tests (from project root)
make test

# Run tests for specific service
cd services/ims-api
make test-unit

# Run integration tests
make test-env-up        # Start test environment
make test-db-setup      # Initialize test database
make test-integration   # Run integration tests
make test-env-down      # Stop test environment

# Run specific test
go test -v -run TestIncidentCreate ./internal/handlers

# Run tests with coverage
make test-coverage

# Generate HTML coverage report
make test-coverage-html

# Run tests with race detector
go test -race ./...

# Run benchmarks
make test-bench
```

## Code Quality

```bash
# Run linter
make lint
# or
golangci-lint run ./...

# Auto-fix linting issues
golangci-lint run --fix ./...

# Format code
gofmt -w .
goimports -w .

# Tidy dependencies
make tidy
# or per service
cd services/ims-api && go mod tidy

# Update dependencies
go get -u ./...
go mod tidy
```

## Debugging

```bash
# Run with Delve debugger
cd services/ims-api
dlv debug ./cmd/api

# Debug with breakpoint
dlv debug ./cmd/api
(dlv) break main.main
(dlv) continue

# Attach to running process
ps aux | grep ims-api
dlv attach <PID>

# Debug tests
dlv test ./internal/handlers -- -test.run TestIncidentCreate
```

## Git Workflow

```bash
# Create feature branch
git checkout -b feature/add-incident-attachments

# Commit changes
git add .
git commit -m "feat(ims-api): add incident attachment endpoint"

# Push branch
git push origin feature/add-incident-attachments

# Create PR (if gh CLI is installed)
gh pr create --title "Add incident attachment endpoint" --body "Description..."
```

## Environment Variables

```bash
# Load from .env file
cd services/ims-api
cp .env.example .env
# Edit .env
# Then run
source .env  # if using bash/zsh
go run ./cmd/api

# Or set inline
PG_DSN=postgres://... go run ./cmd/api

# Check current environment
env | grep PG_DSN
```

## Troubleshooting

```bash
# Check if port is in use
lsof -i :8080
lsof -i :5432

# Kill process on port
lsof -ti :8080 | xargs kill -9

# Check service health
curl http://localhost:8080/health
curl http://localhost:9000/minio/health/live

# View Docker logs
docker logs essp-postgres-dev
docker logs essp-postgres-dev --tail 100 -f

# Check PostgreSQL connection
docker exec -it essp-postgres-dev psql -U ssp -c "SELECT version();"

# Check Valkey/Redis connection
docker exec -it essp-valkey-dev valkey-cli ping

# Clear Valkey/Redis cache
docker exec -it essp-valkey-dev valkey-cli FLUSHALL

# Restart infrastructure from scratch
docker compose -f docker-compose.dev.yml down -v
./scripts/dev-start.sh --init-db
```

## Service URLs

### Infrastructure
- PostgreSQL: `localhost:5432` (user: ssp, password: ssp)
- Valkey/Redis: `localhost:6379`
- MinIO API: `http://localhost:9000`
- MinIO Console: `http://localhost:9001` (user: minio, password: minio12345)
- NATS: `nats://localhost:4222`
- NATS Monitoring: `http://localhost:8222`
- Adminer (optional): `http://localhost:8090`
- Redis Commander (optional): `http://localhost:8091`

### Application Services
- IMS API: `http://localhost:8080`
- SSOT School: `http://localhost:8081`
- SSOT Devices: `http://localhost:8082`
- SSOT Parts: `http://localhost:8083`

## Database Connection Strings

```bash
# IMS
postgres://ssp:ssp@localhost:5432/ssp_ims?sslmode=disable

# SSOT School
postgres://ssp:ssp@localhost:5432/ssp_school?sslmode=disable

# SSOT Devices
postgres://ssp:ssp@localhost:5432/ssp_devices?sslmode=disable

# SSOT Parts
postgres://ssp:ssp@localhost:5432/ssp_parts?sslmode=disable

# Test Database (different port)
postgres://ssp:ssp@localhost:5433/ssp_ims_test?sslmode=disable
```

## Common API Endpoints

```bash
# IMS API Examples
# Health check
curl http://localhost:8080/health

# List incidents (with tenant header)
curl -H "X-Tenant-Id: demo-tenant" http://localhost:8080/api/v1/incidents

# Create incident
curl -X POST http://localhost:8080/api/v1/incidents \
  -H "X-Tenant-Id: demo-tenant" \
  -H "Content-Type: application/json" \
  -d '{"title":"Test incident","description":"Test"}'

# SSOT School
curl http://localhost:8081/health
curl -H "X-Tenant-Id: demo-tenant" http://localhost:8081/api/v1/schools

# SSOT Devices
curl http://localhost:8082/health
curl -H "X-Tenant-Id: demo-tenant" http://localhost:8082/api/v1/devices

# SSOT Parts
curl http://localhost:8083/health
curl -H "X-Tenant-Id: demo-tenant" http://localhost:8083/api/v1/parts
```

## Performance Profiling

```bash
# Run with CPU profiling
go run -cpuprofile=cpu.prof ./cmd/api

# Run with memory profiling
go run -memprofile=mem.prof ./cmd/api

# Analyze profile
go tool pprof cpu.prof
go tool pprof -http=:8081 cpu.prof

# Benchmark with profiling
go test -bench=. -cpuprofile=cpu.prof ./...
```

## Docker Image Build

```bash
# Build specific service image
make docker-build-ims-api

# Build all service images
make docker-build

# Build and tag manually
docker build -t essp/ims-api:latest -f services/ims-api/Dockerfile .

# Push to registry
make docker-push
```

## Useful Aliases

Add to your `~/.bashrc` or `~/.zshrc`:

```bash
# ESSP Development Aliases
alias essp-start='cd /home/pato/opt/ESSP && ./scripts/dev-start.sh'
alias essp-init='cd /home/pato/opt/ESSP && ./scripts/dev-start.sh --init-db'
alias essp-stop='cd /home/pato/opt/ESSP && docker compose -f docker-compose.dev.yml down'
alias essp-logs='cd /home/pato/opt/ESSP && docker compose -f docker-compose.dev.yml logs -f'
alias essp-clean='cd /home/pato/opt/ESSP && docker compose -f docker-compose.dev.yml down -v'
alias essp-db='docker exec -it essp-postgres-dev psql -U ssp -d ssp_ims'
alias essp-redis='docker exec -it essp-valkey-dev valkey-cli'

# Service aliases
alias essp-ims='cd /home/pato/opt/ESSP/services/ims-api && go run ./cmd/api'
alias essp-school='cd /home/pato/opt/ESSP/services/ssot-school && go run ./cmd/ssot_school'
alias essp-devices='cd /home/pato/opt/ESSP/services/ssot-devices && go run ./cmd/ssot_devices'
alias essp-parts='cd /home/pato/opt/ESSP/services/ssot-parts && go run ./cmd/ssot_parts'
alias essp-worker='cd /home/pato/opt/ESSP/services/sync-worker && go run ./cmd/worker'
```

## Additional Resources

- Full documentation: `/home/pato/opt/ESSP/docs/DEVELOPMENT.md`
- Architecture: `/home/pato/opt/ESSP/docs/architecture.md`
- API documentation: `/home/pato/opt/ESSP/docs/api.md`
- Scripts README: `/home/pato/opt/ESSP/scripts/README.md`
