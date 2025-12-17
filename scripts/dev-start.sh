#!/bin/bash
# Quick start script for ESSP local development
# This script starts all infrastructure services and optionally runs database initialization

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Get script directory
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_ROOT="$( cd "$SCRIPT_DIR/.." && pwd )"

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}ESSP Development Environment${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Parse command line arguments
INIT_DB=false
WITH_TOOLS=false
PROFILE_ARG=""

while [[ $# -gt 0 ]]; do
    case $1 in
        --init-db)
            INIT_DB=true
            shift
            ;;
        --with-tools)
            WITH_TOOLS=true
            PROFILE_ARG="--profile tools"
            shift
            ;;
        --help|-h)
            echo "Usage: $0 [OPTIONS]"
            echo ""
            echo "Options:"
            echo "  --init-db      Initialize databases and run migrations"
            echo "  --with-tools   Start optional tools (Adminer, Redis Commander)"
            echo "  --help, -h     Show this help message"
            echo ""
            echo "Examples:"
            echo "  $0                        # Start infrastructure only"
            echo "  $0 --init-db              # Start infrastructure and initialize databases"
            echo "  $0 --with-tools           # Start with optional management tools"
            echo "  $0 --init-db --with-tools # Start everything with database initialization"
            exit 0
            ;;
        *)
            echo -e "${RED}Unknown option: $1${NC}"
            echo "Use --help for usage information"
            exit 1
            ;;
    esac
done

cd "$PROJECT_ROOT"

# Check if docker compose is available
if ! command -v docker &> /dev/null; then
    echo -e "${RED}Error: Docker is not installed${NC}"
    echo "Please install Docker: https://docs.docker.com/get-docker/"
    exit 1
fi

if ! docker compose version &> /dev/null; then
    echo -e "${RED}Error: Docker Compose is not available${NC}"
    echo "Please install Docker Compose: https://docs.docker.com/compose/install/"
    exit 1
fi

# Start infrastructure services
echo -e "${YELLOW}Starting infrastructure services...${NC}"
docker compose -f docker-compose.dev.yml $PROFILE_ARG up -d

echo ""
echo -e "${YELLOW}Waiting for services to be ready...${NC}"

# Wait for PostgreSQL
echo -n "  PostgreSQL... "
timeout 30 bash -c 'until docker exec essp-postgres-dev pg_isready -U ssp -q 2>/dev/null; do sleep 1; done' && echo -e "${GREEN}ready${NC}" || (echo -e "${RED}failed${NC}" && exit 1)

# Wait for Valkey
echo -n "  Valkey... "
timeout 30 bash -c 'until docker exec essp-valkey-dev valkey-cli ping 2>/dev/null | grep -q PONG; do sleep 1; done' && echo -e "${GREEN}ready${NC}" || (echo -e "${RED}failed${NC}" && exit 1)

# Wait for MinIO
echo -n "  MinIO... "
timeout 30 bash -c 'until curl -sf http://localhost:9000/minio/health/live >/dev/null 2>&1; do sleep 1; done' && echo -e "${GREEN}ready${NC}" || (echo -e "${RED}failed${NC}" && exit 1)

# Wait for NATS
echo -n "  NATS... "
timeout 30 bash -c 'until curl -sf http://localhost:8222/healthz >/dev/null 2>&1; do sleep 1; done' && echo -e "${GREEN}ready${NC}" || (echo -e "${RED}failed${NC}" && exit 1)

echo ""

# Initialize databases if requested
if [ "$INIT_DB" = true ]; then
    echo -e "${YELLOW}Initializing databases...${NC}"
    bash "$SCRIPT_DIR/init-dev-db.sh"
else
    echo -e "${YELLOW}Skipping database initialization (use --init-db to initialize)${NC}"
    echo ""
fi

# Show service status
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}Infrastructure services are running!${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""

echo -e "${BLUE}Service URLs:${NC}"
echo -e "  ${YELLOW}PostgreSQL:${NC}     localhost:5432 (user: ssp, password: ssp)"
echo -e "  ${YELLOW}Valkey/Redis:${NC}   localhost:6379"
echo -e "  ${YELLOW}MinIO API:${NC}      http://localhost:9000"
echo -e "  ${YELLOW}MinIO Console:${NC}  http://localhost:9001 (user: minio, password: minio12345)"
echo -e "  ${YELLOW}NATS:${NC}           nats://localhost:4222"
echo -e "  ${YELLOW}NATS Monitoring:${NC} http://localhost:8222"

if [ "$WITH_TOOLS" = true ]; then
    echo ""
    echo -e "${BLUE}Management Tools:${NC}"
    echo -e "  ${YELLOW}Adminer:${NC}          http://localhost:8090"
    echo -e "  ${YELLOW}Redis Commander:${NC}  http://localhost:8091"
fi

echo ""
echo -e "${BLUE}Docker Compose Commands:${NC}"
echo -e "  ${YELLOW}View logs:${NC}        docker compose -f docker-compose.dev.yml logs -f"
echo -e "  ${YELLOW}Stop services:${NC}    docker compose -f docker-compose.dev.yml down"
echo -e "  ${YELLOW}Restart:${NC}          docker compose -f docker-compose.dev.yml restart"
echo -e "  ${YELLOW}Clean all data:${NC}   docker compose -f docker-compose.dev.yml down -v"
echo ""

if [ "$INIT_DB" = false ]; then
    echo -e "${YELLOW}To initialize databases, run:${NC}"
    echo -e "  bash scripts/init-dev-db.sh"
    echo -e "  ${BLUE}or${NC}"
    echo -e "  bash scripts/dev-start.sh --init-db"
    echo ""
fi

echo -e "${BLUE}Start Application Services:${NC}"
echo ""
echo -e "  ${YELLOW}Terminal 1 - IMS API (port 8100):${NC}"
echo -e "    cd services/ims-api && go run ./cmd/api"
echo -e "    ${BLUE}(Uses .env file for configuration)${NC}"
echo ""
echo -e "  ${YELLOW}Terminal 2 - Dashboard (port 5173):${NC}"
echo -e "    cd dashboard && npm run dev"
echo -e "    ${BLUE}(Proxies /v1 and /api to localhost:8100)${NC}"
echo ""
echo -e "  ${YELLOW}Terminal 3 - SSOT School:${NC}"
echo -e "    cd services/ssot-school && go run ./cmd/ssot_school"
echo ""
echo -e "  ${YELLOW}Terminal 4 - SSOT Devices:${NC}"
echo -e "    cd services/ssot-devices && go run ./cmd/ssot_devices"
echo ""
echo -e "  ${YELLOW}Terminal 5 - SSOT Parts:${NC}"
echo -e "    cd services/ssot-parts && go run ./cmd/ssot_parts"
echo ""
echo -e "  ${YELLOW}Terminal 6 - SSOT HR:${NC}"
echo -e "    cd services/ssot-hr && go run ./cmd/ssot_hr"
echo ""
echo -e "  ${YELLOW}Terminal 7 - Sync Worker:${NC}"
echo -e "    cd services/sync-worker && go run ./cmd/worker"
echo ""

echo -e "${BLUE}Application Service URLs (when running):${NC}"
echo -e "  ${YELLOW}Dashboard:${NC}        http://localhost:5173"
echo -e "  ${YELLOW}IMS API:${NC}          http://localhost:8100"
echo -e "  ${YELLOW}SSOT School:${NC}      http://localhost:8081"
echo -e "  ${YELLOW}SSOT Devices:${NC}     http://localhost:8082"
echo -e "  ${YELLOW}SSOT Parts:${NC}       http://localhost:8083"
echo -e "  ${YELLOW}SSOT HR:${NC}          http://localhost:8300"
echo ""

echo -e "${GREEN}Ready for development!${NC}"
