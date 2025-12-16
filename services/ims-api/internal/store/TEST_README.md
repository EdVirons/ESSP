# Store Repository Tests

This directory contains comprehensive unit tests for the repository layer of the IMS API service.

## Test Files

1. **testhelpers_test.go** - Common test utilities and helper functions
2. **incidents_repo_test.go** - Tests for IncidentRepo
3. **workorders_repo_test.go** - Tests for WorkOrderRepo
4. **service_shops_repo_test.go** - Tests for ServiceShopRepo
5. **workorder_parts_repo_test.go** - Tests for WorkOrderPartRepo

## Prerequisites

### Database Setup

The tests require a PostgreSQL database for integration testing. You have two options:

#### Option 1: Use Environment Variable
Set the `TEST_PG_DSN` environment variable to point to your test database:

```bash
export TEST_PG_DSN="postgres://username:password@localhost:5432/ims_test?sslmode=disable"
```

#### Option 2: Use Default (localhost)
If `TEST_PG_DSN` is not set, tests will use the default connection:
```
postgres://postgres:postgres@localhost:5432/ims_test?sslmode=disable
```

### Database Schema

Before running tests, ensure your test database has the required tables. You can:
1. Run migrations on the test database
2. Use the schema from your development/production database

Required tables:
- `incidents`
- `work_orders`
- `service_shops`
- `work_order_parts`

## Running Tests

### Run All Store Tests
```bash
cd /home/pato/opt/ESSP/services/ims-api
go test ./internal/store -v
```

### Run Specific Test File
```bash
go test ./internal/store -v -run TestIncidentRepo
go test ./internal/store -v -run TestWorkOrderRepo
go test ./internal/store -v -run TestServiceShopRepo
go test ./internal/store -v -run TestWorkOrderPartRepo
```

### Run Specific Test Function
```bash
go test ./internal/store -v -run TestIncidentRepo_Create
go test ./internal/store -v -run TestWorkOrderRepo_List
```

### Run with Coverage
```bash
go test ./internal/store -v -cover
go test ./internal/store -v -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Run with Race Detection
```bash
go test ./internal/store -v -race
```

## Test Coverage

### IncidentRepo Tests
- ✅ TestIncidentRepo_Create - Tests creating incidents
- ✅ TestIncidentRepo_GetByID - Tests retrieving incidents by ID with tenant/school isolation
- ✅ TestIncidentRepo_List - Tests listing with filters (status, device, query) and pagination
- ✅ TestIncidentRepo_UpdateStatus - Tests status updates
- ✅ TestIncidentRepo_MarkSLABreaches - Tests SLA breach detection

### WorkOrderRepo Tests
- ✅ TestWorkOrderRepo_Create - Tests creating work orders
- ✅ TestWorkOrderRepo_GetByID - Tests retrieving work orders by ID
- ✅ TestWorkOrderRepo_List - Tests listing with filters (status, device, incident) and pagination
- ✅ TestWorkOrderRepo_UpdateStatus - Tests status updates
- ✅ TestWorkOrderRepo_SetApprovalStatus - Tests approval status updates
- ✅ TestWorkOrderRepo_ListByPhase - Tests listing work orders by phase

### ServiceShopRepo Tests
- ✅ TestServiceShopRepo_Create - Tests creating service shops
- ✅ TestServiceShopRepo_GetByID - Tests retrieving service shops by ID
- ✅ TestServiceShopRepo_GetByCounty - Tests finding active shops by county
- ✅ TestServiceShopRepo_GetBySubCounty - Tests finding shops by sub-county with coverage level filtering
- ✅ TestServiceShopRepo_List - Tests listing with filters (county, active status) and pagination

### WorkOrderPartRepo Tests
- ✅ TestWorkOrderPartRepo_Create - Tests creating work order parts
- ✅ TestWorkOrderPartRepo_GetByID - Tests retrieving work order parts by ID
- ✅ TestWorkOrderPartRepo_List - Tests listing parts for a work order with pagination
- ✅ TestCreateWorkOrderPartTx - Tests transactional part creation
- ✅ TestUpdateWorkOrderPartUsedTx - Tests transactional quantity used updates
- ✅ TestUpdateWorkOrderPartPlannedTx - Tests transactional quantity planned updates

## Testing Patterns Used

### Table-Driven Tests
All tests use table-driven patterns for better maintainability:
```go
tests := []struct {
    name    string
    input   models.Incident
    wantErr bool
}{
    {name: "valid incident", input: validIncident(), wantErr: false},
    {name: "missing tenant", input: incidentNoTenant(), wantErr: true},
}
```

### Subtests with t.Run()
Each test case runs as a subtest for better isolation:
```go
for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        // Test implementation
    })
}
```

### Test Data Cleanup
Tests clean up their data before and after execution:
```go
defer cleanupIncidents(t, pool, tenantID, schoolID)
cleanupIncidents(t, pool, tenantID, schoolID)
```

### Factory Functions
Helper functions create valid test data:
```go
incident := validIncident()
workOrder := validWorkOrder()
shop := validServiceShop()
part := validWorkOrderPart()
```

## Test Data Isolation

Tests use unique identifiers to avoid conflicts:
- Tenant IDs: `tenant-test-*`
- School IDs: `school-test-*`
- Entity IDs: `inc-test-*`, `wo-test-*`, `shop-test-*`, `wop-test-*`

## Error Cases Tested

1. **Not Found Errors**: Attempting to retrieve non-existent records
2. **Tenant Isolation**: Verifying records cannot be accessed with wrong tenant ID
3. **School Isolation**: Verifying records cannot be accessed with wrong school ID
4. **Pagination**: Testing cursor-based pagination logic
5. **Filtering**: Testing various filter combinations
6. **Transactional Updates**: Testing quantity constraints in work order parts

## Continuous Integration

To run these tests in CI/CD:

```yaml
# Example GitHub Actions workflow
- name: Run Store Tests
  env:
    TEST_PG_DSN: postgres://postgres:postgres@localhost:5432/ims_test
  run: |
    go test ./internal/store -v -race -cover
```

## Troubleshooting

### Database Connection Failed
- Verify PostgreSQL is running
- Check the connection string (DSN)
- Ensure the database exists
- Verify user permissions

### Tests Fail with "relation does not exist"
- Run database migrations on the test database
- Verify schema is up to date

### Cleanup Warnings
- Cleanup warnings are logged but don't fail tests
- They indicate test data from previous runs

### Transaction Tests Fail
- Ensure your PostgreSQL supports transactions
- Verify foreign key constraints don't prevent cleanup

## Future Enhancements

Potential improvements to the test suite:
- [ ] Add benchmark tests for performance profiling
- [ ] Add concurrent access tests
- [ ] Add database migration tests
- [ ] Mock database for faster unit tests (optional)
- [ ] Add test data seeding scripts
- [ ] Add integration tests with other layers
