# ESSP Development Scripts

Helper scripts for local development and testing.

## Available Scripts

### `dev-start.sh`

Quick start script for the entire development environment.

**Usage:**
```bash
# Start infrastructure only
./scripts/dev-start.sh

# Start infrastructure and initialize databases
./scripts/dev-start.sh --init-db

# Start with optional management tools (Adminer, Redis Commander)
./scripts/dev-start.sh --with-tools

# Start everything
./scripts/dev-start.sh --init-db --with-tools

# Show help
./scripts/dev-start.sh --help
```

**What it does:**
- Starts all infrastructure services (PostgreSQL, Valkey, MinIO, NATS)
- Waits for services to be ready
- Optionally initializes databases and runs migrations
- Optionally starts management tools
- Displays helpful information about running services

### `init-dev-db.sh`

Initialize development databases and run migrations.

**Usage:**
```bash
./scripts/init-dev-db.sh
```

**What it does:**
- Creates all required databases (ssp_ims, ssp_school, ssp_devices, ssp_parts)
- Runs migrations for all services
- Displays connection URLs and next steps

**Prerequisites:**
- PostgreSQL must be running (start with `docker-compose.dev.yml`)

## Quick Start

For first-time setup:

```bash
# 1. Start infrastructure and initialize everything
./scripts/dev-start.sh --init-db --with-tools

# 2. In separate terminals, start each service
cd services/ims-api && go run ./cmd/api
cd services/ssot-school && go run ./cmd/ssot_school
cd services/ssot-devices && go run ./cmd/ssot_devices
cd services/ssot-parts && go run ./cmd/ssot_parts
cd services/sync-worker && go run ./cmd/worker
```

## Subsequent Development Sessions

```bash
# Just start infrastructure (databases persist in Docker volumes)
./scripts/dev-start.sh

# Then start the services you're working on
cd services/ims-api && go run ./cmd/api
```

## Stopping Services

```bash
# Stop infrastructure services
docker compose -f docker-compose.dev.yml down

# Stop and remove all data (clean slate)
docker compose -f docker-compose.dev.yml down -v
```

## Troubleshooting

### Database Connection Issues

```bash
# Check if PostgreSQL is running
docker ps | grep postgres

# View PostgreSQL logs
docker compose -f docker-compose.dev.yml logs postgres

# Restart PostgreSQL
docker compose -f docker-compose.dev.yml restart postgres
```

### Reinitialize Databases

```bash
# Stop services and remove data
docker compose -f docker-compose.dev.yml down -v

# Start fresh
./scripts/dev-start.sh --init-db
```

### Port Conflicts

If you see "port already in use" errors:

```bash
# Check what's using the port
lsof -i :5432  # PostgreSQL
lsof -i :6379  # Valkey
lsof -i :9000  # MinIO
lsof -i :4222  # NATS

# Stop conflicting services or change ports in docker-compose.dev.yml
```

## Service Ports

| Service | Port | Purpose |
|---------|------|---------|
| PostgreSQL | 5432 | Database |
| Valkey | 6379 | Cache/Redis |
| MinIO API | 9000 | Object storage |
| MinIO Console | 9001 | Web UI |
| NATS | 4222 | Messaging |
| NATS Monitoring | 8222 | NATS Web UI |
| Adminer | 8090 | Database UI (optional) |
| Redis Commander | 8091 | Redis UI (optional) |
| IMS API | 8080 | IMS service |
| SSOT School | 8081 | School SSOT |
| SSOT Devices | 8082 | Devices SSOT |
| SSOT Parts | 8083 | Parts SSOT |

## Environment Variables

Scripts support these environment variables for customization:

```bash
# PostgreSQL connection
export POSTGRES_HOST=localhost
export POSTGRES_PORT=5432
export POSTGRES_USER=ssp
export POSTGRES_PASSWORD=ssp

# Then run scripts
./scripts/init-dev-db.sh
```
