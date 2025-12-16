# ESSP Development Guide

This guide will help you set up a complete local development environment for the EdVirons School Services Platform (ESSP).

## Table of Contents

- [Prerequisites](#prerequisites)
- [Environment Setup](#environment-setup)
- [Local Infrastructure](#local-infrastructure)
- [Configuration](#configuration)
- [Running Services](#running-services)
- [Database](#database)
- [Testing](#testing)
- [Common Tasks](#common-tasks)
- [Debugging](#debugging)
- [Code Style](#code-style)

## Prerequisites

### Required Software

1. **Go 1.22+**
   ```bash
   # Check your Go version
   go version

   # Download from https://go.dev/dl/ if needed
   ```

2. **Docker & Docker Compose**
   ```bash
   # Check Docker installation
   docker --version
   docker compose version

   # Minimum versions:
   # - Docker: 20.10+
   # - Docker Compose: 2.0+
   ```

3. **Git**
   ```bash
   git --version
   ```

4. **Make** (optional but recommended)
   ```bash
   make --version
   ```

### Optional Tools

- **golangci-lint** - For code linting
  ```bash
  # Install from https://golangci-lint.run/usage/install/
  # macOS
  brew install golangci-lint

  # Linux
  curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
  ```

- **Air** - For hot reload during development
  ```bash
  go install github.com/cosmtrek/air@latest
  ```

- **Delve** - Go debugger
  ```bash
  go install github.com/go-delve/delve/cmd/dlv@latest
  ```

## Environment Setup

### 1. Clone the Repository

```bash
git clone <repository-url> /home/pato/opt/ESSP
cd /home/pato/opt/ESSP
```

### 2. Install Go Dependencies

The project uses Go workspaces. Install dependencies for all modules:

```bash
# Sync workspace and tidy all modules
make tidy
```

This will run `go mod tidy` for:
- `shared/` - Shared packages
- `services/ims-api/` - IMS API service
- `services/ssot-school/` - School SSOT service
- `services/ssot-devices/` - Devices SSOT service
- `services/ssot-parts/` - Parts SSOT service
- `services/sync-worker/` - Event sync worker

### 3. IDE Setup

#### VS Code

Install recommended extensions:
- Go (golang.go)
- Go Test Explorer
- Docker
- YAML

Create `.vscode/settings.json`:
```json
{
  "go.useLanguageServer": true,
  "go.lintTool": "golangci-lint",
  "go.lintOnSave": "workspace",
  "go.formatTool": "goimports",
  "go.testFlags": ["-v", "-race"],
  "go.coverOnSave": true,
  "editor.formatOnSave": true,
  "[go]": {
    "editor.codeActionsOnSave": {
      "source.organizeImports": true
    }
  }
}
```

#### GoLand

1. Open the project root directory
2. GoLand should automatically detect the Go workspace
3. Enable Go Modules support: **Settings > Go > Go Modules**
4. Configure file watchers for goimports and gofmt

## Local Infrastructure

### Starting Infrastructure Services

The platform requires several backing services. Use Docker Compose for local development:

```bash
# Start all infrastructure services
docker compose -f docker-compose.dev.yml up -d

# Check service status
docker compose -f docker-compose.dev.yml ps

# View logs
docker compose -f docker-compose.dev.yml logs -f

# Stop all services
docker compose -f docker-compose.dev.yml down

# Stop and remove all data (clean slate)
docker compose -f docker-compose.dev.yml down -v
```

The development compose file (`docker-compose.dev.yml`) includes:

- **PostgreSQL 16** - Main database (port 5432)
- **Valkey/Redis 7** - Caching and rate limiting (port 6379)
- **MinIO** - S3-compatible object storage for attachments (ports 9000, 9001)
- **NATS** - Event messaging (port 4222)

### Accessing Services

| Service | URL | Credentials |
|---------|-----|-------------|
| PostgreSQL | `localhost:5432` | user: `ssp`, password: `ssp` |
| Valkey/Redis | `localhost:6379` | (no password) |
| MinIO Console | `http://localhost:9001` | user: `minio`, password: `minio12345` |
| MinIO API | `http://localhost:9000` | user: `minio`, password: `minio12345` |
| NATS | `nats://localhost:4222` | (no auth) |

## Configuration

### Environment Variables

Each service reads configuration from environment variables. For local development, you can either:

1. **Set environment variables directly:**
   ```bash
   export PG_DSN="postgres://ssp:ssp@localhost:5432/ssp_ims?sslmode=disable"
   export REDIS_ADDR="localhost:6379"
   export NATS_URL="nats://localhost:4222"
   ```

2. **Create a `.env` file** (not tracked in git):
   ```bash
   # Create .env file in service directory
   cd services/ims-api
   cp .env.ratelimit.example .env
   # Edit .env with your settings
   ```

### IMS API Configuration

Key environment variables for `services/ims-api`:

```bash
# Application
APP_ENV=dev
HTTP_ADDR=:8080
LOG_LEVEL=debug

# Database
PG_DSN=postgres://ssp:ssp@localhost:5432/ssp_ims?sslmode=disable

# Redis/Valkey
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0

# NATS
NATS_URL=nats://localhost:4222

# MinIO
MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY=minio
MINIO_SECRET_KEY=minio12345
MINIO_USE_SSL=false
MINIO_REGION=us-east-1
ATTACHMENTS_BUCKET=edvirons-ims
ATTACHMENTS_PUBLIC_BASE_URL=http://localhost:9000

# Authentication (disable for local dev)
AUTH_ENABLED=false
DEV_TENANT_ID=demo-tenant
DEV_SCHOOL_ID=demo-school

# Rate Limiting (disable or use high limits for local dev)
RATE_LIMIT_ENABLED=false
RATE_LIMIT_READ_RPM=1000
RATE_LIMIT_WRITE_RPM=500
RATE_LIMIT_BURST=100

# CORS
CORS_ALLOWED_ORIGINS=*

# Work Order Settings
AUTO_ROUTE_WORK_ORDERS=true
DEFAULT_REPAIR_LOCATION=service_shop

# SSOT Sync
SSOT_SYNC_PAGE_SIZE=500
SCHOOL_SSOT_BASE_URL=http://localhost:8081
DEVICE_SSOT_BASE_URL=http://localhost:8082
PARTS_SSOT_BASE_URL=http://localhost:8083
```

### SSOT Services Configuration

Similar configuration for SSOT services (School, Devices, Parts):

```bash
# services/ssot-school
HTTP_ADDR=:8081
DB_URL=postgres://ssp:ssp@localhost:5432/ssp_school?sslmode=disable
NATS_URL=nats://localhost:4222
LOG_LEVEL=debug

# services/ssot-devices
HTTP_ADDR=:8082
DB_URL=postgres://ssp:ssp@localhost:5432/ssp_devices?sslmode=disable
NATS_URL=nats://localhost:4222
LOG_LEVEL=debug

# services/ssot-parts
HTTP_ADDR=:8083
DB_URL=postgres://ssp:ssp@localhost:5432/ssp_parts?sslmode=disable
NATS_URL=nats://localhost:4222
LOG_LEVEL=debug
```

### Sync Worker Configuration

```bash
NATS_URL=nats://localhost:4222
LOG_LEVEL=debug
```

## Running Services

### Running Individual Services

Each service can be run independently:

```bash
# IMS API
cd services/ims-api
go run ./cmd/api

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
```

### Running with Make

```bash
# From service directory
cd services/ims-api
make run

# Build specific service
make build-ims-api

# Build all services
make build-all
```

### Running All Services

You have two options:

**Option 1: Using Docker Compose** (includes all services)
```bash
cd deployments/docker
docker compose up --build
```

**Option 2: Run infrastructure in Docker, services locally** (recommended for development)
```bash
# Terminal 1: Start infrastructure
docker compose -f docker-compose.dev.yml up -d

# Terminal 2: IMS API
cd services/ims-api && go run ./cmd/api

# Terminal 3: SSOT School
cd services/ssot-school && go run ./cmd/ssot_school

# Terminal 4: SSOT Devices
cd services/ssot-devices && go run ./cmd/ssot_devices

# Terminal 5: SSOT Parts
cd services/ssot-parts && go run ./cmd/ssot_parts

# Terminal 6: Sync Worker
cd services/sync-worker && go run ./cmd/worker
```

### Hot Reload with Air

For automatic reloading during development:

1. Install Air:
   ```bash
   go install github.com/cosmtrek/air@latest
   ```

2. Create `.air.toml` in service directory:
   ```toml
   root = "."
   testdata_dir = "testdata"
   tmp_dir = "tmp"

   [build]
     args_bin = []
     bin = "./tmp/main"
     cmd = "go build -o ./tmp/main ./cmd/api"
     delay = 1000
     exclude_dir = ["assets", "tmp", "vendor", "testdata"]
     exclude_file = []
     exclude_regex = ["_test.go"]
     exclude_unchanged = false
     follow_symlink = false
     full_bin = ""
     include_dir = []
     include_ext = ["go", "tpl", "tmpl", "html"]
     include_file = []
     kill_delay = "0s"
     log = "build-errors.log"
     poll = false
     poll_interval = 0
     rerun = false
     rerun_delay = 500
     send_interrupt = false
     stop_on_error = false

   [color]
     app = ""
     build = "yellow"
     main = "magenta"
     runner = "green"
     watcher = "cyan"

   [log]
     main_only = false
     time = false

   [misc]
     clean_on_exit = false

   [screen]
     clear_on_rebuild = false
     keep_scroll = true
   ```

3. Run with Air:
   ```bash
   cd services/ims-api
   air
   ```

## Database

### Creating Databases

Create the required databases:

```bash
# Connect to PostgreSQL
docker exec -it essp-postgres-dev psql -U ssp

# Create databases
CREATE DATABASE ssp_ims;
CREATE DATABASE ssp_school;
CREATE DATABASE ssp_devices;
CREATE DATABASE ssp_parts;
\q
```

Or use psql directly:
```bash
PGPASSWORD=ssp psql -h localhost -U ssp -c "CREATE DATABASE ssp_ims;"
PGPASSWORD=ssp psql -h localhost -U ssp -c "CREATE DATABASE ssp_school;"
PGPASSWORD=ssp psql -h localhost -U ssp -c "CREATE DATABASE ssp_devices;"
PGPASSWORD=ssp psql -h localhost -U ssp -c "CREATE DATABASE ssp_parts;"
```

### Running Migrations

Each service has its own migrations in `migrations/` directory.

**IMS API:**
```bash
cd services/ims-api

# Run migrations
make migrate-up

# Or manually
DB_URL=postgres://ssp:ssp@localhost:5432/ssp_ims?sslmode=disable go run ./cmd/migrate up
```

**SSOT Services:**
```bash
# School SSOT
cd services/ssot-school
DB_URL=postgres://ssp:ssp@localhost:5432/ssp_school?sslmode=disable go run ./cmd/migrate up

# Devices SSOT
cd services/ssot-devices
DB_URL=postgres://ssp:ssp@localhost:5432/ssp_devices?sslmode=disable go run ./cmd/migrate up

# Parts SSOT
cd services/ssot-parts
DB_URL=postgres://ssp:ssp@localhost:5432/ssp_parts?sslmode=disable go run ./cmd/migrate up
```

### Migration Files

Migrations follow a simple pattern:
- Named: `NNN_description.sql` (e.g., `001_init.sql`)
- Executed in lexical order
- Support Goose-style markers for up/down migrations

Example migration:
```sql
-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id TEXT NOT NULL,
    email TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_tenant ON users(tenant_id);

-- +goose Down
DROP TABLE IF EXISTS users;
```

### Seed Data

For local development, you may want to add seed data:

```bash
# Connect to database
PGPASSWORD=ssp psql -h localhost -U ssp -d ssp_ims

# Run SQL commands or load from file
\i /path/to/seed.sql
```

### Database Connection Troubleshooting

**Connection refused:**
```bash
# Check PostgreSQL is running
docker compose -f docker-compose.dev.yml ps postgres

# Check logs
docker compose -f docker-compose.dev.yml logs postgres

# Restart PostgreSQL
docker compose -f docker-compose.dev.yml restart postgres
```

**Password authentication failed:**
- Verify credentials in connection string
- Check environment variables
- Ensure database exists

**Too many connections:**
- Adjust `max_connections` in PostgreSQL config
- Check for connection leaks in code
- Use connection pooling properly

## Testing

### Unit Tests

Run unit tests (no external dependencies):

```bash
# From project root - all services
make test

# From specific service
cd services/ims-api
make test-unit

# With coverage
make test-coverage

# HTML coverage report
make test-coverage-html
```

### Integration Tests

Integration tests require a test environment:

```bash
cd services/ims-api

# Start test infrastructure
make test-env-up

# Setup test database
make test-db-setup

# Run integration tests
make test-integration

# Cleanup test environment
make test-env-down
```

Test environment uses different ports:
- PostgreSQL: 5433
- Redis: 6380
- NATS: 4223
- MinIO: 9001

### Running Specific Tests

```bash
# Run specific test
go test -v -run TestIncidentCreate ./internal/handlers

# Run tests in specific package
go test -v ./internal/store

# Run with race detector
go test -race ./...

# Verbose output
go test -v ./...
```

### Test Structure

Tests are organized as:
- `*_test.go` - Unit tests alongside code
- `tests/` - Integration tests
- `internal/testutil/` - Test utilities
- `internal/mocks/` - Mock implementations

Example test:
```go
func TestCreateIncident(t *testing.T) {
    // Setup
    db := testutil.NewTestDB(t)
    defer db.Close()

    repo := store.NewIncidentRepo(db)

    // Execute
    incident, err := repo.Create(ctx, &store.Incident{
        TenantID: "test-tenant",
        Title: "Test incident",
    })

    // Assert
    require.NoError(t, err)
    assert.NotEmpty(t, incident.ID)
    assert.Equal(t, "Test incident", incident.Title)
}
```

## Common Tasks

### Adding a New Endpoint

1. **Define the model** in `internal/models/`:
   ```go
   type Resource struct {
       ID        string    `json:"id"`
       TenantID  string    `json:"tenant_id"`
       Name      string    `json:"name"`
       CreatedAt time.Time `json:"created_at"`
   }
   ```

2. **Create repository** in `internal/store/`:
   ```go
   type ResourceRepo struct {
       db *pgxpool.Pool
   }

   func (r *ResourceRepo) Create(ctx context.Context, resource *Resource) error {
       // Implementation
   }
   ```

3. **Create handler** in `internal/handlers/`:
   ```go
   func (h *Handler) CreateResource(w http.ResponseWriter, r *http.Request) {
       // Parse request
       // Validate
       // Call repository
       // Return response
   }
   ```

4. **Register route** in `internal/api/router.go`:
   ```go
   r.Post("/resources", h.CreateResource)
   ```

5. **Add tests** in `internal/handlers/*_test.go`

### Creating Database Migrations

1. Create new file in `services/{service}/migrations/`:
   ```bash
   cd services/ims-api/migrations
   touch 009_add_resources.sql
   ```

2. Write migration:
   ```sql
   -- +goose Up
   CREATE TABLE resources (
       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
       tenant_id TEXT NOT NULL,
       name TEXT NOT NULL,
       created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
   );

   -- +goose Down
   DROP TABLE resources;
   ```

3. Run migration:
   ```bash
   make migrate-up
   ```

### Adding New Models

1. Define struct in `internal/models/`:
   ```go
   type Device struct {
       ID          string     `json:"id" db:"id"`
       TenantID    string     `json:"tenant_id" db:"tenant_id"`
       SerialNo    string     `json:"serial_no" db:"serial_no"`
       Model       string     `json:"model" db:"model"`
       Status      string     `json:"status" db:"status"`
       CreatedAt   time.Time  `json:"created_at" db:"created_at"`
       UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
   }
   ```

2. Add validation:
   ```go
   func (d *Device) Validate() error {
       if d.TenantID == "" {
           return errors.New("tenant_id required")
       }
       if d.SerialNo == "" {
           return errors.New("serial_no required")
       }
       return nil
   }
   ```

3. Create repository methods in `internal/store/`

### Working with Shared Packages

The `shared/` module contains code used across services:

```bash
# Make changes in shared/
cd shared
# Edit files

# Update in all services
cd ..
make tidy
```

Common shared packages:
- `shared/events/` - Event definitions
- `shared/middleware/` - Shared middleware
- `shared/logging/` - Logging utilities

## Debugging

### Using Delve Debugger

**Install Delve:**
```bash
go install github.com/go-delve/delve/cmd/dlv@latest
```

**Debug a service:**
```bash
cd services/ims-api
dlv debug ./cmd/api

# In Delve console
(dlv) break main.main
(dlv) continue
(dlv) next
(dlv) print cfg
```

**Debug with arguments:**
```bash
dlv debug ./cmd/migrate -- up
```

**Attach to running process:**
```bash
# Find process
ps aux | grep ims-api

# Attach
dlv attach <PID>
```

### VS Code Debugging

Create `.vscode/launch.json`:
```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Debug IMS API",
            "type": "go",
            "request": "launch",
            "mode": "debug",
            "program": "${workspaceFolder}/services/ims-api/cmd/api",
            "env": {
                "PG_DSN": "postgres://ssp:ssp@localhost:5432/ssp_ims?sslmode=disable",
                "REDIS_ADDR": "localhost:6379",
                "LOG_LEVEL": "debug"
            }
        },
        {
            "name": "Debug Current Test",
            "type": "go",
            "request": "launch",
            "mode": "test",
            "program": "${fileDirname}",
            "args": [
                "-test.run",
                "^${selectedText}$"
            ]
        }
    ]
}
```

### Logging Configuration

Adjust log levels:
```bash
# In code (internal/config/config.go)
LOG_LEVEL=debug   # debug, info, warn, error

# Or per service
cd services/ims-api
LOG_LEVEL=debug go run ./cmd/api
```

View structured logs:
```go
import "go.uber.org/zap"

logger.Info("processing request",
    zap.String("tenant_id", tenantID),
    zap.String("request_id", requestID),
    zap.Duration("elapsed", elapsed),
)
```

### Common Issues

**Service won't start:**
- Check port conflicts: `lsof -i :8080`
- Verify database is running
- Check environment variables
- Review logs

**Database connection fails:**
- Ensure PostgreSQL is running
- Verify connection string
- Check firewall/network

**Tests fail:**
- Ensure test environment is up: `make test-env-up`
- Run migrations: `make test-db-setup`
- Check test database connection string

**Import cycle:**
- Reorganize packages
- Move shared code to `shared/`
- Use interfaces to break dependencies

## Code Style

### Go Conventions

Follow standard Go conventions:
- Use `gofmt` and `goimports`
- Follow [Effective Go](https://go.dev/doc/effective_go)
- Use meaningful variable names
- Keep functions small and focused
- Document exported functions

### Project Structure

Standard service structure:
```
services/ims-api/
├── cmd/
│   ├── api/          # Main service entry point
│   └── migrate/      # Migration runner
├── internal/
│   ├── api/          # HTTP server setup
│   ├── config/       # Configuration
│   ├── handlers/     # HTTP handlers
│   ├── middleware/   # HTTP middleware
│   ├── models/       # Data models
│   ├── store/        # Database repositories
│   └── logging/      # Logging setup
├── migrations/       # SQL migrations
├── tests/           # Integration tests
├── go.mod
├── go.sum
├── Makefile
└── Dockerfile
```

### Naming Conventions

- **Files:** `snake_case.go`
- **Packages:** Short, lowercase, no underscores
- **Types:** `PascalCase`
- **Functions/Methods:** `PascalCase` (exported), `camelCase` (unexported)
- **Variables:** `camelCase`
- **Constants:** `PascalCase` or `SCREAMING_SNAKE_CASE` for enums

### Error Handling

```go
// Good: Return errors, don't panic
func DoSomething() error {
    if err := validate(); err != nil {
        return fmt.Errorf("validation failed: %w", err)
    }
    return nil
}

// Good: Handle errors explicitly
if err := DoSomething(); err != nil {
    logger.Error("operation failed", logging.Err(err))
    return err
}

// Bad: Ignoring errors
_ = DoSomething()

// Bad: Panic in library code
panic("something went wrong")
```

### Linting

Run linters before committing:

```bash
# Lint all code
make lint

# From project root
golangci-lint run ./...

# From specific service
cd services/ims-api
golangci-lint run ./...

# Auto-fix issues
golangci-lint run --fix ./...
```

Linter configuration in `.golangci.yml` at project root.

### Commit Message Format

Follow conventional commits:

```
type(scope): subject

body

footer
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `refactor`: Code refactoring
- `test`: Test changes
- `chore`: Build/tooling changes

Examples:
```
feat(ims-api): add incident attachment endpoint

Implements file upload for incident attachments using MinIO.
Includes validation, virus scanning, and size limits.

Closes #123

---

fix(ssot-school): correct school import validation

School codes must be unique per county, not globally.

---

docs: update development setup guide

Add section on debugging with Delve.
```

### Code Review Checklist

Before submitting PRs:
- [ ] Code compiles without errors
- [ ] All tests pass
- [ ] Linters pass
- [ ] Added tests for new functionality
- [ ] Updated documentation
- [ ] No commented-out code
- [ ] Error handling is correct
- [ ] Logging is appropriate
- [ ] Security considerations addressed
- [ ] Database migrations included if needed

## Additional Resources

- [Go Documentation](https://go.dev/doc/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Docker Documentation](https://docs.docker.com/)
- [NATS Documentation](https://docs.nats.io/)
- [MinIO Documentation](https://min.io/docs/)

## Getting Help

- Check existing documentation in `docs/`
- Review test files for usage examples
- Ask in team chat
- Open an issue for bugs or feature requests

---

Happy coding!
