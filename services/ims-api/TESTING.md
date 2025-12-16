# Testing Infrastructure Documentation

This document describes the testing infrastructure for the IMS API service.

## Overview

The testing infrastructure provides:
- **Test utilities** - Common helpers for database setup, HTTP testing, and fixtures
- **Mock implementations** - In-memory mocks for Redis, MinIO, and NATS
- **Docker test environment** - Isolated services for integration testing
- **Makefile targets** - Convenient commands for running tests

## Quick Start

### Running Unit Tests

Unit tests run quickly without external dependencies:

```bash
cd services/ims-api
make test-unit
```

### Running Integration Tests

Integration tests require the test environment to be running:

```bash
# Start test environment
make test-env-up

# Set up test database schema
make test-db-setup

# Run integration tests
make test-integration

# Clean up
make test-env-down
```

### Generate Coverage Report

```bash
make test-coverage          # Text report
make test-coverage-html     # HTML report (opens in browser)
```

## Directory Structure

```
services/ims-api/
├── internal/
│   ├── testutil/          # Test utilities and helpers
│   │   ├── testutil.go    # Common test helpers
│   │   ├── db.go          # Database setup/teardown
│   │   ├── http.go        # HTTP test client
│   │   └── fixtures.go    # Test data fixtures
│   └── mocks/             # Mock implementations
│       ├── redis.go       # Mock Redis client
│       ├── minio.go       # Mock MinIO client
│       └── nats.go        # Mock NATS publisher
├── Makefile               # Test targets
└── TESTING.md            # This file
```

## Test Utilities

### Database Testing (`testutil.db.go`)

```go
import "github.com/edvirons/ssp/ims/internal/testutil"

func TestMyRepository(t *testing.T) {
    // Set up test database
    db := testutil.SetupTestDB(t)

    // Clean up tables
    testutil.CleanupAllTables(t, db.RawPool())

    // Test your repository...
}
```

**Available Functions:**
- `SetupTestDB(t)` - Create test database connection
- `SetupTestDBPool(t)` - Create raw pgxpool connection
- `CleanupAllTables(t, pool)` - Truncate all tables
- `TruncateTables(t, pool, tables...)` - Truncate specific tables
- `WithTransaction(t, pool, fn)` - Run test in transaction (auto-rollback)
- `ExecSQL(t, pool, sql, args...)` - Execute SQL statement
- `WaitForDB(t, dsn, timeout)` - Wait for database to be ready

### HTTP Testing (`testutil.http.go`)

```go
import "github.com/edvirons/ssp/ims/internal/testutil"

func TestMyHandler(t *testing.T) {
    // Create HTTP test client
    router := setupRouter() // Your router
    client := testutil.NewHTTPTestClient(t, router)

    // Make requests with fluent API
    client.
        WithTenant("test-tenant").
        WithSchool("test-school").
        Post("/api/incidents", map[string]string{
            "title": "Test incident",
        }).
        AssertStatus(201).
        AssertJSONField("title", "Test incident")
}
```

**Available Methods:**
- `NewHTTPTestClient(t, handler)` - Create test client
- `WithHeader(key, value)` - Add header to requests
- `WithTenant(tenantID)` - Set tenant header
- `WithSchool(schoolID)` - Set school header
- `WithAuth(token)` - Set auth header
- `Get(path)` - Make GET request
- `Post(path, body)` - Make POST request
- `Put(path, body)` - Make PUT request
- `Patch(path, body)` - Make PATCH request
- `Delete(path)` - Make DELETE request

**Response Assertions:**
- `AssertStatus(code)` - Assert HTTP status code
- `AssertJSON(v)` - Unmarshal response to struct
- `AssertJSONField(field, value)` - Assert JSON field value
- `AssertHeader(key, value)` - Assert header value
- `AssertContentType(contentType)` - Assert content type
- `AssertBodyContains(substr)` - Assert body contains substring

### Test Fixtures (`testutil.fixtures.go`)

```go
import "github.com/edvirons/ssp/ims/internal/testutil"

func TestIncidentWorkflow(t *testing.T) {
    pool := testutil.SetupTestDBPool(t)
    cfg := testutil.DefaultFixtureConfig()

    // Create test data
    school := testutil.CreateSchoolSnapshot(t, pool, cfg)
    device := testutil.CreateDeviceSnapshot(t, pool, cfg, "device-1")
    incident := testutil.CreateIncident(t, pool, cfg, "device-1")
    workOrder := testutil.CreateWorkOrder(t, pool, cfg, incident.ID, "device-1")

    // Test your workflow...
}
```

**Available Fixtures:**
- `CreateSchoolSnapshot(t, pool, cfg)` - Create test school
- `CreateDeviceSnapshot(t, pool, cfg, deviceID)` - Create test device
- `CreatePartSnapshot(t, pool, cfg, partID)` - Create test part
- `CreateIncident(t, pool, cfg, deviceID)` - Create test incident
- `CreateWorkOrder(t, pool, cfg, incidentID, deviceID)` - Create test work order
- `CreateServiceShop(t, pool, cfg)` - Create test service shop
- `CreateServiceStaff(t, pool, cfg, shopID)` - Create test staff
- `CreatePart(t, pool, cfg)` - Create test part
- `CreateInventoryItem(t, pool, cfg, shopID, partID, qty)` - Create inventory
- `CreateProgram(t, pool, cfg)` - Create test program
- `CreateAttachment(t, pool, cfg, entityType, entityID)` - Create test attachment

### Common Helpers (`testutil.testutil.go`)

```go
import "github.com/edvirons/ssp/ims/internal/testutil"

func TestWithHelpers(t *testing.T) {
    // Skip if running in short mode
    testutil.SkipIfShort(t)

    // Skip if INTEGRATION_TEST not set
    testutil.SkipIfNoIntegration(t)

    // Generate unique test ID
    id := testutil.GenerateTestID(t)

    // Wait for a condition
    testutil.WaitForCondition(t, func() bool {
        return someCondition()
    }, 5*time.Second, "condition description")
}
```

## Mock Implementations

### MockRedisClient (`mocks/redis.go`)

In-memory Redis client for unit testing:

```go
import "github.com/edvirons/ssp/ims/internal/mocks"

func TestWithRedis(t *testing.T) {
    redis := mocks.NewMockRedisClient()

    // Use like normal Redis client
    redis.Set(ctx, "key", "value", 10*time.Second)
    val := redis.Get(ctx, "key")

    // Verify operations
    assert.Equal(t, 1, redis.GetCallCount("Set"))
    assert.Equal(t, 1, redis.GetCallCount("Get"))

    // Reset for next test
    redis.Reset()
}
```

### MockMinIOClient (`mocks/minio.go`)

In-memory MinIO/S3 client for unit testing:

```go
import "github.com/edvirons/ssp/ims/internal/mocks"

func TestWithMinIO(t *testing.T) {
    minio := mocks.NewMockMinIOClient()

    // Create bucket
    minio.MakeBucket(ctx, "test-bucket", nil)

    // Upload object
    minio.PutObject(ctx, "test-bucket", "test.txt", []byte("content"), 7, "text/plain")

    // Generate presigned URL
    url, _ := minio.PresignedGetObject(ctx, "test-bucket", "test.txt", 15*time.Minute, nil)

    // Verify
    assert.True(t, minio.ObjectExists("test-bucket", "test.txt"))
    assert.Equal(t, 1, minio.GetObjectCount("test-bucket"))
}
```

### MockNATSPublisher (`mocks/nats.go`)

In-memory NATS publisher for unit testing:

```go
import "github.com/edvirons/ssp/ims/internal/mocks"

func TestWithNATS(t *testing.T) {
    nats := mocks.NewMockNATSPublisher()

    // Publish events
    nats.PublishJSON("incident.created", map[string]string{
        "id": "123",
        "type": "hardware",
    })

    // Verify
    assert.True(t, nats.WasPublished("incident.created"))
    assert.Equal(t, 1, nats.GetMessageCount("incident.created"))

    // Get message
    msg, ok := nats.GetLastMessage("incident.created")
    assert.True(t, ok)
}
```

## Test Environment

The test environment uses Docker Compose to run isolated instances of:
- **PostgreSQL** (port 5433) - Test database
- **Valkey/Redis** (port 6380) - Test cache
- **NATS** (port 4223) - Test message broker
- **MinIO** (port 9001) - Test object storage

All data is ephemeral (stored in tmpfs) for fast test execution.

### Managing Test Environment

```bash
# Start all services
make test-env-up

# Stop and remove all services (including data)
make test-env-down

# Set up database schema
make test-db-setup
```

### Environment Variables

Integration tests use these environment variables:

```bash
TEST_DB_DSN=postgres://ssp:ssp@localhost:5433/ssp_ims_test?sslmode=disable
INTEGRATION_TEST=1  # Set to run integration tests
```

## Writing Tests

### Unit Test Example

```go
package service_test

import (
    "testing"
    "github.com/edvirons/ssp/ims/internal/mocks"
    "github.com/stretchr/testify/assert"
)

func TestIncidentService_Create(t *testing.T) {
    // Arrange
    redis := mocks.NewMockRedisClient()
    nats := mocks.NewMockNATSPublisher()
    service := NewIncidentService(redis, nats)

    // Act
    incident, err := service.Create(ctx, request)

    // Assert
    assert.NoError(t, err)
    assert.NotEmpty(t, incident.ID)
    assert.True(t, nats.WasPublished("incident.created"))
}
```

### Integration Test Example

```go
package service_test

import (
    "testing"
    "github.com/edvirons/ssp/ims/internal/testutil"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestIncidentRepository_Integration(t *testing.T) {
    // Skip if not running integration tests
    testutil.SkipIfNoIntegration(t)

    // Arrange
    db := testutil.SetupTestDB(t)
    pool := db.RawPool()
    cfg := testutil.DefaultFixtureConfig()

    testutil.CleanupAllTables(t, pool)
    testutil.CreateSchoolSnapshot(t, pool, cfg)
    testutil.CreateDeviceSnapshot(t, pool, cfg, "device-1")

    // Act
    incident := testutil.CreateIncident(t, pool, cfg, "device-1")

    // Assert
    assert.NotEmpty(t, incident.ID)

    // Test repository methods...
}
```

### Table-Driven Test Example

```go
func TestValidateIncident(t *testing.T) {
    tests := []struct {
        name    string
        input   Incident
        wantErr bool
    }{
        {
            name: "valid incident",
            input: Incident{
                Title:    "Test",
                Severity: "medium",
            },
            wantErr: false,
        },
        {
            name: "missing title",
            input: Incident{
                Severity: "medium",
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateIncident(tt.input)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

## Best Practices

1. **Use table-driven tests** for testing multiple scenarios
2. **Skip integration tests** in short mode with `testutil.SkipIfNoIntegration(t)`
3. **Clean up test data** using `testutil.CleanupAllTables()` or transactions
4. **Use fixtures** for consistent test data setup
5. **Test in parallel** when safe using `t.Parallel()`
6. **Mock external dependencies** in unit tests
7. **Use integration tests** for database and critical paths
8. **Keep tests focused** - one assertion per test when possible
9. **Use descriptive test names** that explain what is being tested
10. **Clean up resources** using `t.Cleanup()` or defer

## Makefile Targets Reference

| Target | Description |
|--------|-------------|
| `make test` | Run all unit tests (short mode) |
| `make test-unit` | Run unit tests only |
| `make test-integration` | Run integration tests (requires test env) |
| `make test-coverage` | Generate coverage report |
| `make test-coverage-html` | Generate HTML coverage report |
| `make test-env-up` | Start test environment |
| `make test-env-down` | Stop test environment |
| `make test-db-setup` | Run migrations on test database |
| `make test-verbose` | Run tests with verbose output |
| `make test-report` | Generate JSON test report |
| `make test-bench` | Run benchmark tests |

## Continuous Integration

For CI/CD pipelines:

```bash
# Start test environment
make test-env-up

# Set up database
make test-db-setup

# Run all tests with coverage
INTEGRATION_TEST=1 make test-coverage

# Clean up
make test-env-down
```

## Troubleshooting

### Tests hang or timeout
- Check if test environment services are running: `docker-compose -f deployments/docker/docker-compose.test.yml ps`
- Check service logs: `docker-compose -f deployments/docker/docker-compose.test.yml logs`

### Database connection errors
- Ensure test database is running on port 5433
- Verify `TEST_DB_DSN` environment variable
- Run `make test-db-setup` to ensure schema is created

### Port conflicts
- The test environment uses different ports (5433, 6380, 4223, 9001) to avoid conflicts
- Check if these ports are available

### Stale test data
- Always use `testutil.CleanupAllTables()` or transactions
- Use `make test-env-down` to completely reset the environment

## Additional Resources

- [Go Testing Package](https://pkg.go.dev/testing)
- [Testify Documentation](https://github.com/stretchr/testify)
- [Table-Driven Tests in Go](https://dave.cheney.net/2019/05/07/prefer-table-driven-tests)
