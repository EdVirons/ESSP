package testutil

import (
	"context"
	"fmt"
	"testing"

	"github.com/edvirons/ssp/ims/internal/store"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

// SetupTestDB creates a test database connection and returns a Postgres store.
// It automatically cleans up the connection when the test completes.
//
// Usage:
//   db := testutil.SetupTestDB(t)
//   // Use db for testing
func SetupTestDB(t *testing.T) *store.Postgres {
	t.Helper()

	// Use environment variable or default test DSN
	dsn := GetEnvOrDefault("TEST_DB_DSN", "postgres://ssp:ssp@localhost:5432/ssp_ims_test?sslmode=disable")

	ctx := context.Background()
	db, err := store.NewPostgres(ctx, dsn)
	require.NoError(t, err, "failed to connect to test database")

	// Verify connection
	err = db.Ping(ctx)
	require.NoError(t, err, "failed to ping test database")

	// Register cleanup
	t.Cleanup(func() {
		db.Close()
	})

	return db
}

// SetupTestDBPool creates a raw pgxpool connection for tests that need direct pool access.
func SetupTestDBPool(t *testing.T) *pgxpool.Pool {
	t.Helper()

	dsn := GetEnvOrDefault("TEST_DB_DSN", "postgres://ssp:ssp@localhost:5432/ssp_ims_test?sslmode=disable")

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	require.NoError(t, err, "failed to connect to test database pool")

	err = pool.Ping(ctx)
	require.NoError(t, err, "failed to ping test database")

	t.Cleanup(func() {
		pool.Close()
	})

	return pool
}

// TruncateTables truncates specified tables to clean up test data.
// This is useful for resetting state between tests.
func TruncateTables(t *testing.T, pool *pgxpool.Pool, tables ...string) {
	t.Helper()

	ctx := context.Background()
	for _, table := range tables {
		query := fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table)
		_, err := pool.Exec(ctx, query)
		require.NoError(t, err, "failed to truncate table %s", table)
	}
}

// CleanupAllTables truncates all main tables in the test database.
// Use this to ensure a clean state before running tests.
func CleanupAllTables(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()

	tables := []string{
		"work_order_approvals",
		"work_order_deliverables",
		"work_order_schedules",
		"phase_checklist_templates",
		"survey_photos",
		"survey_rooms",
		"site_surveys",
		"boq_items",
		"service_phases",
		"school_service_programs",
		"school_contacts",
		"work_order_parts",
		"inventory",
		"parts",
		"service_staff",
		"service_shops",
		"attachments",
		"work_orders",
		"incidents",
		"parts_snapshots",
		"devices_snapshots",
		"schools_snapshots",
		"ssot_state",
	}

	TruncateTables(t, pool, tables...)
}

// WithTransaction runs a test function within a database transaction and rolls it back.
// This ensures test isolation without affecting the database state.
func WithTransaction(t *testing.T, pool *pgxpool.Pool, fn func(tx pgx.Tx)) {
	t.Helper()

	ctx := context.Background()
	tx, err := pool.Begin(ctx)
	require.NoError(t, err, "failed to begin transaction")

	defer func() {
		err := tx.Rollback(ctx)
		if err != nil {
			t.Logf("warning: failed to rollback transaction: %v", err)
		}
	}()

	fn(tx)
}

// ExecSQL executes a SQL statement and requires it to succeed.
func ExecSQL(t *testing.T, pool *pgxpool.Pool, sql string, args ...interface{}) {
	t.Helper()

	ctx := context.Background()
	_, err := pool.Exec(ctx, sql, args...)
	require.NoError(t, err, "failed to execute SQL: %s", sql)
}

// QueryRow executes a query that returns a single row.
func QueryRow(t *testing.T, pool *pgxpool.Pool, sql string, args ...interface{}) pgx.Row {
	t.Helper()

	ctx := context.Background()
	return pool.QueryRow(ctx, sql, args...)
}

// WaitForDB waits for the database to be ready, useful in integration tests.
func WaitForDB(t *testing.T, dsn string, timeout context.Context) *pgxpool.Pool {
	t.Helper()

	var pool *pgxpool.Pool
	var err error

	for {
		select {
		case <-timeout.Done():
			t.Fatal("timeout waiting for database to be ready")
		default:
			pool, err = pgxpool.New(context.Background(), dsn)
			if err == nil {
				if pool.Ping(context.Background()) == nil {
					t.Cleanup(func() { pool.Close() })
					return pool
				}
				pool.Close()
			}
		}
	}
}
