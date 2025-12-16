# Getting Started with ESSP Development

Welcome to the EdVirons School Services Platform (ESSP) development environment!

## Quick Start (5 minutes)

### 1. Prerequisites Check

Ensure you have the following installed:
- Go 1.22+ (`go version`)
- Docker & Docker Compose (`docker --version` and `docker compose version`)
- Git (`git --version`)

### 2. Install Dependencies

```bash
cd /home/pato/opt/ESSP
make tidy
```

### 3. Start Development Environment

```bash
# Start all infrastructure services and initialize databases
./scripts/dev-start.sh --init-db --with-tools
```

This will:
- Start PostgreSQL, Valkey (Redis), MinIO, and NATS
- Create all databases
- Run migrations
- Start optional management tools (Adminer, Redis Commander)

### 4. Start Application Services

Open 5 terminals and run:

```bash
# Terminal 1 - IMS API (main service)
cd services/ims-api
go run ./cmd/api

# Terminal 2 - SSOT School
cd services/ssot-school
go run ./cmd/ssot_school

# Terminal 3 - SSOT Devices
cd services/ssot-devices
go run ./cmd/ssot_devices

# Terminal 4 - SSOT Parts
cd services/ssot-parts
go run ./cmd/ssot_parts

# Terminal 5 - Sync Worker
cd services/sync-worker
go run ./cmd/worker
```

### 5. Verify Everything is Running

```bash
# Check services
curl http://localhost:8080/health  # IMS API
curl http://localhost:8081/health  # SSOT School
curl http://localhost:8082/health  # SSOT Devices
curl http://localhost:8083/health  # SSOT Parts
```

## What's Next?

### Learn More

- **Full Development Guide**: [docs/DEVELOPMENT.md](docs/DEVELOPMENT.md)
  - Detailed setup instructions
  - Configuration options
  - Testing guide
  - Debugging tips

- **Quick Reference**: [docs/QUICK_REFERENCE.md](docs/QUICK_REFERENCE.md)
  - Common commands
  - API endpoints
  - Database connections
  - Troubleshooting

- **Architecture**: [docs/architecture.md](docs/architecture.md)
  - System overview
  - Service boundaries
  - Data ownership

- **API Documentation**: [docs/api.md](docs/api.md)
  - Routing conventions
  - Headers and authentication
  - Error handling

### Access Management Tools

With `--with-tools` flag, you get:

- **MinIO Console**: http://localhost:9001
  - User: `minio` / Password: `minio12345`
  - Manage file attachments

- **Adminer**: http://localhost:8090
  - Database: PostgreSQL
  - Server: `postgres`
  - User: `ssp` / Password: `ssp`
  - Visual database management

- **Redis Commander**: http://localhost:8091
  - Explore cache data
  - Debug rate limiting

- **NATS Monitoring**: http://localhost:8222
  - Monitor message queues
  - View event streams

### Development Workflow

1. **Make Code Changes**
   - Edit files in your IDE
   - Services auto-reload (if using Air) or restart manually

2. **Run Tests**
   ```bash
   cd services/ims-api
   make test
   ```

3. **Commit Changes**
   ```bash
   git add .
   git commit -m "feat(ims-api): add new feature"
   git push
   ```

### Common Tasks

**Add a new database table:**
```bash
cd services/ims-api/migrations
touch 009_add_my_table.sql
# Edit the migration file
make migrate-up
```

**Reset databases:**
```bash
docker compose -f docker-compose.dev.yml down -v
./scripts/dev-start.sh --init-db
```

**View logs:**
```bash
# Infrastructure logs
docker compose -f docker-compose.dev.yml logs -f

# Specific service
docker compose -f docker-compose.dev.yml logs -f postgres
```

**Debug with VS Code:**
- Press F5
- Select "Debug IMS API" (or other service)
- Set breakpoints and debug

## Project Structure

```
ESSP/
├── docs/                      # Documentation
│   ├── DEVELOPMENT.md        # Full development guide
│   ├── QUICK_REFERENCE.md    # Command reference
│   └── ...
├── services/                  # Microservices
│   ├── ims-api/              # Main API service
│   ├── ssot-school/          # School SSOT
│   ├── ssot-devices/         # Devices SSOT
│   ├── ssot-parts/           # Parts SSOT
│   └── sync-worker/          # Event sync worker
├── shared/                    # Shared libraries
├── scripts/                   # Development scripts
│   ├── dev-start.sh          # Start development environment
│   └── init-dev-db.sh        # Initialize databases
├── docker-compose.dev.yml    # Local infrastructure
├── Makefile                   # Build commands
└── go.work                    # Go workspace
```

## Troubleshooting

### Services won't start

Check if ports are already in use:
```bash
lsof -i :8080  # IMS API
lsof -i :5432  # PostgreSQL
```

### Database connection fails

Ensure PostgreSQL is running:
```bash
docker compose -f docker-compose.dev.yml ps postgres
docker compose -f docker-compose.dev.yml logs postgres
```

### Tests fail

Start test environment:
```bash
cd services/ims-api
make test-env-up
make test-db-setup
make test-integration
```

### Need a clean slate

```bash
# Stop everything and remove all data
docker compose -f docker-compose.dev.yml down -v

# Start fresh
./scripts/dev-start.sh --init-db
```

## Getting Help

- Review documentation in `docs/`
- Check test files for usage examples
- Ask in team chat
- Open an issue for bugs

## IDE Setup

### VS Code (Recommended)

The `.vscode/` directory contains pre-configured settings:

1. Open the project: `code /home/pato/opt/ESSP`
2. Install recommended extensions (VS Code will prompt)
3. Press F5 to debug any service

### GoLand

1. Open the project directory
2. Enable Go Modules support
3. Use run configurations for each service

## Helpful Aliases

Add to your `~/.bashrc` or `~/.zshrc`:

```bash
alias essp-start='./scripts/dev-start.sh'
alias essp-stop='cd /home/pato/opt/ESSP && docker compose -f docker-compose.dev.yml down'
alias essp-ims='cd /home/pato/opt/ESSP/services/ims-api && go run ./cmd/api'
```

Then use:
```bash
essp-start --init-db
essp-ims
```

## Daily Development Flow

```bash
# Morning: Start infrastructure
essp-start

# Start the service you're working on
cd services/ims-api
go run ./cmd/api

# Make changes, test, commit
# ...

# Evening: Stop infrastructure
essp-stop
```

---

**Happy coding!** For more details, see [docs/DEVELOPMENT.md](docs/DEVELOPMENT.md)
