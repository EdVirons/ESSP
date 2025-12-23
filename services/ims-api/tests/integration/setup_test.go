//go:build integration
// +build integration

package integration

import (
	"context"
	"testing"

	"github.com/edvirons/ssp/ims/internal/store"
	"github.com/edvirons/ssp/ims/internal/testutil"
	"github.com/jackc/pgx/v5/pgxpool"
)

// setupTest prepares a clean database state for integration tests.
// It returns a Postgres store instance and a cleanup function.
func setupTest(t *testing.T) (*store.Postgres, func()) {
	t.Helper()

	// Skip if not running integration tests
	testutil.SkipIfNoIntegration(t)

	// Setup database connection
	db := testutil.SetupTestDB(t)

	// Clean all tables to ensure isolation
	testutil.CleanupAllTables(t, db.RawPool())

	// Return cleanup function
	cleanup := func() {
		db.Close()
	}

	return db, cleanup
}

// setupTestWithFixtures prepares a test with common fixtures already created.
// Returns db, config, and cleanup function.
func setupTestWithFixtures(t *testing.T) (*store.Postgres, testutil.FixtureConfig, func()) {
	t.Helper()

	db, cleanup := setupTest(t)

	// Create standard fixtures
	cfg := testutil.DefaultFixtureConfig()
	testutil.CreateSchoolSnapshot(t, db.RawPool(), cfg)
	testutil.CreateServiceShop(t, db.RawPool(), cfg)

	return db, cfg, cleanup
}

// execSQL is a helper to execute SQL statements in tests.
func execSQL(t *testing.T, pool *pgxpool.Pool, sql string, args ...interface{}) {
	t.Helper()
	testutil.ExecSQL(t, pool, sql, args...)
}

// queryCount returns the count from a COUNT(*) query.
func queryCount(t *testing.T, pool *pgxpool.Pool, query string, args ...interface{}) int64 {
	t.Helper()

	var count int64
	err := pool.QueryRow(context.Background(), query, args...).Scan(&count)
	if err != nil {
		t.Fatalf("failed to query count: %v", err)
	}
	return count
}
