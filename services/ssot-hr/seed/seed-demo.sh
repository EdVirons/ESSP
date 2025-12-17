#!/bin/bash
# seed-demo.sh - Seed HR Directory with demo data
set -e

SSOT_HR_URL="${SSOT_HR_URL:-http://localhost:8300}"
TENANT_ID="${TENANT_ID:-demo}"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "Seeding HR Directory at $SSOT_HR_URL for tenant $TENANT_ID..."

# Check if service is available
if ! curl -sf "$SSOT_HR_URL/healthz" > /dev/null 2>&1; then
    echo "Error: ssot-hr service is not available at $SSOT_HR_URL"
    echo "Please ensure the service is running."
    exit 1
fi

# Import the demo data
echo "Importing demo data..."
response=$(curl -sf -X POST "$SSOT_HR_URL/v1/import" \
    -H "Content-Type: application/json" \
    -H "X-Tenant-Id: $TENANT_ID" \
    -d @"$SCRIPT_DIR/demo-data.json")

echo "Import response: $response"
echo ""
echo "HR Directory seeded successfully!"
echo ""
echo "Summary:"
echo "  - 8 Org Units created"
echo "  - 8 People created"
echo "  - 7 Teams created"
echo "  - 9 Team Memberships created"
echo ""
echo "To sync to IMS-API, run:"
echo "  curl -X POST http://localhost:8100/v1/ssot/sync/all -H 'X-Tenant-Id: $TENANT_ID'"
