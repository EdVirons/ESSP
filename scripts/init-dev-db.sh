#!/bin/bash
# Initialize development databases for ESSP platform
# This script creates all required databases and runs migrations

set -e  # Exit on error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
POSTGRES_HOST=${POSTGRES_HOST:-localhost}
POSTGRES_PORT=${POSTGRES_PORT:-5432}
POSTGRES_USER=${POSTGRES_USER:-ssp}
POSTGRES_PASSWORD=${POSTGRES_PASSWORD:-ssp}

# Database names
DATABASES=("ssp_ims" "ssp_school" "ssp_devices" "ssp_parts" "ssp_hr")

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}ESSP Development Database Initialization${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""

# Check if PostgreSQL is running
echo -e "${YELLOW}Checking PostgreSQL connection...${NC}"
if ! PGPASSWORD=$POSTGRES_PASSWORD psql -h $POSTGRES_HOST -p $POSTGRES_PORT -U $POSTGRES_USER -c '\q' 2>/dev/null; then
    echo -e "${RED}Error: Cannot connect to PostgreSQL at $POSTGRES_HOST:$POSTGRES_PORT${NC}"
    echo -e "${YELLOW}Make sure PostgreSQL is running:${NC}"
    echo -e "  docker compose -f docker-compose.dev.yml up -d postgres"
    exit 1
fi
echo -e "${GREEN}✓ PostgreSQL is running${NC}"
echo ""

# Create databases
echo -e "${YELLOW}Creating databases...${NC}"
for db in "${DATABASES[@]}"; do
    echo -n "  Creating database '$db'... "

    # Check if database exists
    if PGPASSWORD=$POSTGRES_PASSWORD psql -h $POSTGRES_HOST -p $POSTGRES_PORT -U $POSTGRES_USER -lqt | cut -d \| -f 1 | grep -qw $db; then
        echo -e "${YELLOW}already exists${NC}"
    else
        PGPASSWORD=$POSTGRES_PASSWORD psql -h $POSTGRES_HOST -p $POSTGRES_PORT -U $POSTGRES_USER -c "CREATE DATABASE $db;" >/dev/null
        echo -e "${GREEN}created${NC}"
    fi
done
echo ""

# Run migrations
echo -e "${YELLOW}Running migrations...${NC}"

# Get the script directory
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_ROOT="$( cd "$SCRIPT_DIR/.." && pwd )"

# IMS API migrations
echo -n "  IMS API... "
cd "$PROJECT_ROOT/services/ims-api"
DB_URL="postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@$POSTGRES_HOST:$POSTGRES_PORT/ssp_ims?sslmode=disable" \
    go run ./cmd/migrate up 2>&1 | grep -v "applied" || true
echo -e "${GREEN}done${NC}"

# SSOT School migrations
echo -n "  SSOT School... "
cd "$PROJECT_ROOT/services/ssot-school"
if [ -d "migrations" ] && [ "$(ls -A migrations/*.sql 2>/dev/null)" ]; then
    DB_URL="postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@$POSTGRES_HOST:$POSTGRES_PORT/ssp_school?sslmode=disable" \
        go run ./cmd/migrate up 2>&1 | grep -v "applied" || true
    echo -e "${GREEN}done${NC}"
else
    echo -e "${YELLOW}no migrations${NC}"
fi

# SSOT Devices migrations
echo -n "  SSOT Devices... "
cd "$PROJECT_ROOT/services/ssot-devices"
if [ -d "migrations" ] && [ "$(ls -A migrations/*.sql 2>/dev/null)" ]; then
    DB_URL="postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@$POSTGRES_HOST:$POSTGRES_PORT/ssp_devices?sslmode=disable" \
        go run ./cmd/migrate up 2>&1 | grep -v "applied" || true
    echo -e "${GREEN}done${NC}"
else
    echo -e "${YELLOW}no migrations${NC}"
fi

# SSOT Parts migrations
echo -n "  SSOT Parts... "
cd "$PROJECT_ROOT/services/ssot-parts"
if [ -d "migrations" ] && [ "$(ls -A migrations/*.sql 2>/dev/null)" ]; then
    DB_URL="postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@$POSTGRES_HOST:$POSTGRES_PORT/ssp_parts?sslmode=disable" \
        go run ./cmd/migrate up 2>&1 | grep -v "applied" || true
    echo -e "${GREEN}done${NC}"
else
    echo -e "${YELLOW}no migrations${NC}"
fi

# SSOT HR migrations
echo -n "  SSOT HR... "
cd "$PROJECT_ROOT/services/ssot-hr"
if [ -d "migrations" ] && [ "$(ls -A migrations/*.sql 2>/dev/null)" ]; then
    PGPASSWORD=$POSTGRES_PASSWORD psql -h $POSTGRES_HOST -p $POSTGRES_PORT -U $POSTGRES_USER -d ssp_hr -f migrations/001_init.sql 2>&1 | grep -v "NOTICE" || true
    echo -e "${GREEN}done${NC}"
else
    echo -e "${YELLOW}no migrations${NC}"
fi

echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}Database initialization complete!${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo -e "${YELLOW}Database URLs:${NC}"
echo -e "  IMS:     postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@$POSTGRES_HOST:$POSTGRES_PORT/ssp_ims?sslmode=disable"
echo -e "  School:  postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@$POSTGRES_HOST:$POSTGRES_PORT/ssp_school?sslmode=disable"
echo -e "  Devices: postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@$POSTGRES_HOST:$POSTGRES_PORT/ssp_devices?sslmode=disable"
echo -e "  Parts:   postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@$POSTGRES_HOST:$POSTGRES_PORT/ssp_parts?sslmode=disable"
echo -e "  HR:      postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@$POSTGRES_HOST:$POSTGRES_PORT/ssp_hr?sslmode=disable"
echo ""
echo -e "${YELLOW}Next steps:${NC}"
echo -e "  1. Start services:"
echo -e "     cd services/ims-api && go run ./cmd/api"
echo -e "     cd services/ssot-school && go run ./cmd/ssot_school"
echo -e "     cd services/ssot-devices && go run ./cmd/ssot_devices"
echo -e "     cd services/ssot-parts && go run ./cmd/ssot_parts"
echo -e "     cd services/ssot-hr && go run ./cmd/ssot_hr"
echo -e "     cd services/sync-worker && go run ./cmd/worker"
echo ""
echo -e "  2. Access services:"
echo -e "     IMS API:      http://localhost:8100"
echo -e "     SSOT School:  http://localhost:8081"
echo -e "     SSOT Devices: http://localhost:8082"
echo -e "     SSOT Parts:   http://localhost:8083"
echo -e "     SSOT HR:      http://localhost:8300"
echo ""
