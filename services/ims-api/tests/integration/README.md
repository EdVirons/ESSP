# IMS-API Integration Tests

This directory contains comprehensive integration tests for the IMS (Incident Management System) API service.

## Overview

Integration tests verify the complete workflow and interaction between different components of the system, including:

- **Incident Lifecycle** - Full workflow from incident creation to closure
- **Work Order Lifecycle** - Complete work order management including BOM, scheduling, deliverables, and approvals
- **SSOT Synchronization** - Testing of snapshot caching and lookup enrichment
- **BOM Operations** - Part reservation, consumption, and release operations

## Test Files

### `incident_lifecycle_test.go`
Tests the complete incident management workflow:
- Creating incidents with SLA calculation
- Status transitions and validation
- Invalid transition prevention
- SLA breach detection for different severity levels
- Incident escalation paths
- Auto-creation of work orders from incidents
- Concurrent incident updates

### `work_order_lifecycle_test.go`
Tests the complete work order management workflow:
- Creating work orders from incidents
- Adding parts to BOM with inventory reservation
- Scheduling work orders
- Adding and managing deliverables
- Submitting and reviewing deliverables
- Requesting and approving work orders
- Status transition validation
- Multiple work orders per incident

### `ssot_sync_test.go`
Tests SSOT (Single Source of Truth) synchronization:
- School snapshot caching and retrieval
- Device snapshot caching and retrieval
- Part snapshot caching and retrieval
- Lookup enrichment from cached snapshots
- Sync state management for incremental updates
- Bulk snapshot upserts
- Concurrent snapshot updates
- Handling missing snapshots

### `bom_operations_test.go`
Tests Bill of Materials (BOM) operations:
- Adding parts to work order BOM
- Consuming parts from inventory
- Releasing unused parts back to inventory
- Handling insufficient inventory scenarios
- Managing multiple parts in a BOM
- Concurrent inventory reservation
- Complete inventory tracking through reservation, consumption, and release

## Prerequisites

1. **PostgreSQL Database**: A test database must be available
2. **Environment Variables**: Configure the test database connection
3. **Database Migrations**: Ensure all migrations are applied to the test database

## Setup

### 1. Create Test Database

```bash
# Create test database
createdb ssp_ims_test

# Apply migrations
cd /home/pato/opt/ESSP/services/ims-api
make migrate-test
```

### 2. Configure Environment

Set the test database connection string:

```bash
export TEST_DB_DSN="postgres://ssp:ssp@localhost:5432/ssp_ims_test?sslmode=disable"
export INTEGRATION_TEST=1
```

Alternatively, create a `.env.test` file in the service root:

```env
TEST_DB_DSN=postgres://ssp:ssp@localhost:5432/ssp_ims_test?sslmode=disable
INTEGRATION_TEST=1
```

## Running Tests

### Run All Integration Tests

```bash
# From the service root
cd /home/pato/opt/ESSP/services/ims-api

# Run all integration tests
INTEGRATION_TEST=1 go test -v ./tests/integration/...

# Or using make
make test-integration
```

### Run Specific Test File

```bash
# Run only incident lifecycle tests
INTEGRATION_TEST=1 go test -v ./tests/integration/ -run TestIncidentLifecycle

# Run only work order lifecycle tests
INTEGRATION_TEST=1 go test -v ./tests/integration/ -run TestWorkOrderLifecycle

# Run only SSOT sync tests
INTEGRATION_TEST=1 go test -v ./tests/integration/ -run TestSSOTSync

# Run only BOM operations tests
INTEGRATION_TEST=1 go test -v ./tests/integration/ -run TestBOMOperations
```

### Run Specific Test Case

```bash
# Run a single test case
INTEGRATION_TEST=1 go test -v ./tests/integration/ -run TestIncidentLifecycle_HappyPath

# Run with timeout
INTEGRATION_TEST=1 go test -v -timeout 30m ./tests/integration/...
```

### Run with Build Tags

Integration tests use build tags to prevent accidental execution:

```bash
# Explicitly specify build tag
go test -v -tags=integration ./tests/integration/...
```

## Test Isolation

Each test:
- Uses a dedicated test database connection
- Cleans all tables before execution using `testutil.CleanupAllTables()`
- Creates necessary fixtures using `testutil` helpers
- Runs in isolation from other tests

## Common Test Patterns

### Basic Test Setup

```go
func TestMyFeature(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test in short mode")
    }

    db, cfg, cleanup := setupTestWithFixtures(t)
    defer cleanup()

    ctx := context.Background()

    // Your test code here
}
```

### Using Test Fixtures

```go
// Create device snapshot
deviceID := "test-device-001"
testutil.CreateDeviceSnapshot(t, db.RawPool(), cfg, deviceID)

// Create incident
incident := testutil.CreateIncident(t, db.RawPool(), cfg, deviceID)

// Create work order
workOrder := testutil.CreateWorkOrder(t, db.RawPool(), cfg, incident.ID, deviceID)

// Create service shop and staff
shop := testutil.CreateServiceShop(t, db.RawPool(), cfg)
staff := testutil.CreateServiceStaff(t, db.RawPool(), cfg, shop.ID)

// Create parts and inventory
part := testutil.CreatePart(t, db.RawPool(), cfg)
inventory := testutil.CreateInventoryItem(t, db.RawPool(), cfg, shop.ID, part.ID, 100)
```

## Test Data Management

### Fixture Configuration

Tests use `testutil.FixtureConfig` to manage tenant and school IDs:

```go
cfg := testutil.DefaultFixtureConfig()
// cfg.TenantID = "test-tenant"
// cfg.SchoolID = "test-school"
```

### Database Cleanup

Tests automatically clean up using:
- `testutil.CleanupAllTables()` - Truncates all tables
- `t.Cleanup()` - Registers cleanup functions
- Deferred cleanup functions

## Debugging Tests

### Verbose Output

```bash
# Run with verbose output
INTEGRATION_TEST=1 go test -v ./tests/integration/...

# Run with race detection
INTEGRATION_TEST=1 go test -v -race ./tests/integration/...
```

### Database Inspection

During test development, you can inspect the database state:

```bash
# Connect to test database
psql ssp_ims_test

# Check data
SELECT * FROM incidents;
SELECT * FROM work_orders;
SELECT * FROM work_order_parts;
SELECT * FROM inventory;
```

### Test Logging

Tests use the testing.T logger:

```go
t.Logf("Creating incident with ID: %s", incident.ID)
```

## CI/CD Integration

### GitHub Actions Example

```yaml
- name: Run Integration Tests
  env:
    TEST_DB_DSN: postgres://postgres:postgres@localhost:5432/test_db?sslmode=disable
    INTEGRATION_TEST: 1
  run: |
    make migrate-test
    make test-integration
```

### Docker Compose

Use the test docker-compose configuration:

```bash
# Start test database
docker-compose -f deployments/docker/docker-compose.test.yml up -d

# Run migrations
make migrate-test

# Run tests
INTEGRATION_TEST=1 go test -v ./tests/integration/...

# Cleanup
docker-compose -f deployments/docker/docker-compose.test.yml down -v
```

## Performance Considerations

- Integration tests are slower than unit tests (database I/O)
- Each test cleans the database, adding overhead
- Consider parallel test execution for large test suites:

```bash
# Run tests in parallel (use with caution - may cause race conditions)
INTEGRATION_TEST=1 go test -v -parallel 4 ./tests/integration/...
```

## Troubleshooting

### Test Database Connection Fails

```bash
# Verify database exists
psql -l | grep ssp_ims_test

# Check connection string
echo $TEST_DB_DSN

# Test connection
psql "$TEST_DB_DSN" -c "SELECT 1;"
```

### Migrations Not Applied

```bash
# Check migration status
make migrate-status-test

# Apply migrations
make migrate-test
```

### Tests Skip Automatically

If tests skip with "Skipping integration test", ensure:

```bash
export INTEGRATION_TEST=1
```

### Foreign Key Violations

Ensure test data is created in the correct order:
1. Snapshots (schools, devices, parts)
2. Service shops and staff
3. Incidents
4. Work orders
5. BOM items, schedules, deliverables

## Best Practices

1. **Always use fixtures** - Use `testutil` helpers for test data creation
2. **Clean state** - Each test should start with a clean database
3. **Use assertions** - Use `testify/assert` and `testify/require`
4. **Test realistic scenarios** - Test complete workflows, not just happy paths
5. **Handle errors** - Always check and handle errors appropriately
6. **Add context** - Use descriptive test names and assertions with messages
7. **Test concurrency** - Include concurrent operation tests where relevant
8. **Verify side effects** - Check that operations have the expected database effects

## Contributing

When adding new integration tests:

1. Follow the existing test structure and patterns
2. Use the `// +build integration` build tag
3. Add `SkipIfShort()` check at the beginning
4. Use `setupTestWithFixtures()` for consistent setup
5. Document complex test scenarios with comments
6. Ensure tests are idempotent and isolated
7. Add the test to this README if it covers new functionality

## Related Documentation

- [Testing Strategy](../../TESTING.md) - Overall testing approach
- [Test Setup Summary](../../TEST_SETUP_SUMMARY.md) - Quick setup guide
- [API Documentation](../../README.md) - API service overview
