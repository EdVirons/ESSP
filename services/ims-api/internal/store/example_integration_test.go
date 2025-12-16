package store_test

import (
	"context"
	"testing"
	"time"

	"github.com/edvirons/ssp/ims/internal/models"
	"github.com/edvirons/ssp/ims/internal/store"
	"github.com/edvirons/ssp/ims/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Example integration test demonstrating the full testing infrastructure
// This test shows how to:
// - Set up a test database
// - Create test fixtures
// - Test repository methods
// - Clean up test data

func TestIncidentRepository_Integration(t *testing.T) {
	// Skip if not running integration tests
	testutil.SkipIfNoIntegration(t)

	// Set up test database
	db := testutil.SetupTestDB(t)
	pool := db.RawPool()
	cfg := testutil.DefaultFixtureConfig()

	// Clean up all tables to ensure clean state
	testutil.CleanupAllTables(t, pool)

	// Create required snapshots (foreign key dependencies)
	testutil.CreateSchoolSnapshot(t, pool, cfg)
	testutil.CreateDeviceSnapshot(t, pool, cfg, "device-001")

	// Get the repository
	repo := db.Incidents()

	t.Run("Create and Get Incident", func(t *testing.T) {
		ctx := context.Background()

		// Create an incident
		incident := models.Incident{
			ID:          testutil.GenerateTestID(t),
			TenantID:    cfg.TenantID,
			SchoolID:    cfg.SchoolID,
			DeviceID:    "device-001",
			SchoolName:  "Test School",
			Category:    "hardware",
			Severity:    models.SeverityMedium,
			Status:      models.IncidentNew,
			Title:       "Screen is cracked",
			Description: "The device screen is cracked and needs replacement",
			ReportedBy:  "teacher@school.com",
			SLADueAt:    time.Now().Add(24 * time.Hour),
			SLABreached: false,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		err := repo.Create(ctx, incident)
		require.NoError(t, err, "failed to create incident")

		// Retrieve the incident
		retrieved, err := repo.GetByID(ctx, cfg.TenantID, cfg.SchoolID, incident.ID)
		require.NoError(t, err, "failed to get incident")

		// Verify the data
		assert.Equal(t, incident.ID, retrieved.ID)
		assert.Equal(t, incident.TenantID, retrieved.TenantID)
		assert.Equal(t, incident.SchoolID, retrieved.SchoolID)
		assert.Equal(t, incident.DeviceID, retrieved.DeviceID)
		assert.Equal(t, incident.Title, retrieved.Title)
		assert.Equal(t, incident.Category, retrieved.Category)
		assert.Equal(t, incident.Severity, retrieved.Severity)
		assert.Equal(t, incident.Status, retrieved.Status)
	})

	t.Run("List Incidents with Pagination", func(t *testing.T) {
		ctx := context.Background()

		// Clean up from previous test
		testutil.TruncateTables(t, pool, "incidents")

		// Create multiple incidents
		for i := 1; i <= 5; i++ {
			incident := models.Incident{
				ID:          testutil.GenerateTestID(t),
				TenantID:    cfg.TenantID,
				SchoolID:    cfg.SchoolID,
				DeviceID:    "device-001",
				Category:    "hardware",
				Severity:    models.SeverityMedium,
				Status:      models.IncidentNew,
				Title:       "Test Incident",
				Description: "Description",
				ReportedBy:  "test@example.com",
				SLADueAt:    time.Now().Add(24 * time.Hour),
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}
			err := repo.Create(ctx, incident)
			require.NoError(t, err)

			// Small delay to ensure different timestamps
			time.Sleep(10 * time.Millisecond)
		}

		// List incidents with limit
		params := store.IncidentListParams{
			TenantID: cfg.TenantID,
			SchoolID: cfg.SchoolID,
			Limit:    3,
		}

		incidents, nextCursor, err := repo.List(ctx, params)
		require.NoError(t, err, "failed to list incidents")

		// Should return 3 incidents
		assert.Len(t, incidents, 3)
		assert.NotEmpty(t, nextCursor, "should have next cursor for pagination")

		// All incidents should be ordered by created_at DESC
		for i := 1; i < len(incidents); i++ {
			assert.True(t, incidents[i-1].CreatedAt.After(incidents[i].CreatedAt) ||
				incidents[i-1].CreatedAt.Equal(incidents[i].CreatedAt))
		}
	})

	t.Run("Update Incident Status", func(t *testing.T) {
		ctx := context.Background()

		// Clean up from previous test
		testutil.TruncateTables(t, pool, "incidents")

		// Create an incident
		incident := testutil.CreateIncident(t, pool, cfg, "device-001")

		// Update status
		updated, err := repo.UpdateStatus(ctx, cfg.TenantID, cfg.SchoolID, incident.ID, models.IncidentAcknowledged, time.Now())
		require.NoError(t, err, "failed to update incident status")

		// Verify status changed
		assert.Equal(t, models.IncidentAcknowledged, updated.Status)

		// Verify updated_at changed
		assert.True(t, updated.UpdatedAt.After(incident.UpdatedAt))
	})

	t.Run("Filter Incidents by Status", func(t *testing.T) {
		ctx := context.Background()

		// Clean up from previous test
		testutil.TruncateTables(t, pool, "incidents")

		// Create incidents with different statuses
		statuses := []models.IncidentStatus{
			models.IncidentNew,
			models.IncidentAcknowledged,
			models.IncidentInProgress,
			models.IncidentNew,
		}

		for _, status := range statuses {
			incident := models.Incident{
				ID:          testutil.GenerateTestID(t),
				TenantID:    cfg.TenantID,
				SchoolID:    cfg.SchoolID,
				DeviceID:    "device-001",
				Category:    "hardware",
				Severity:    models.SeverityMedium,
				Status:      status,
				Title:       "Test Incident",
				Description: "Description",
				ReportedBy:  "test@example.com",
				SLADueAt:    time.Now().Add(24 * time.Hour),
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}
			err := repo.Create(ctx, incident)
			require.NoError(t, err)
		}

		// Filter by status
		params := store.IncidentListParams{
			TenantID: cfg.TenantID,
			SchoolID: cfg.SchoolID,
			Status:   string(models.IncidentNew),
			Limit:    10,
		}

		incidents, _, err := repo.List(ctx, params)
		require.NoError(t, err)

		// Should return only incidents with "new" status
		assert.Len(t, incidents, 2)
		for _, inc := range incidents {
			assert.Equal(t, models.IncidentNew, inc.Status)
		}
	})

	t.Run("Search Incidents by Query", func(t *testing.T) {
		ctx := context.Background()

		// Clean up from previous test
		testutil.TruncateTables(t, pool, "incidents")

		// Create incidents with different titles
		titles := []string{
			"Screen is broken",
			"Battery not charging",
			"Keyboard issue",
		}

		for _, title := range titles {
			incident := models.Incident{
				ID:          testutil.GenerateTestID(t),
				TenantID:    cfg.TenantID,
				SchoolID:    cfg.SchoolID,
				DeviceID:    "device-001",
				Category:    "hardware",
				Severity:    models.SeverityMedium,
				Status:      models.IncidentNew,
				Title:       title,
				Description: "Test description",
				ReportedBy:  "test@example.com",
				SLADueAt:    time.Now().Add(24 * time.Hour),
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}
			err := repo.Create(ctx, incident)
			require.NoError(t, err)
		}

		// Search for "screen"
		params := store.IncidentListParams{
			TenantID: cfg.TenantID,
			SchoolID: cfg.SchoolID,
			Query:    "screen",
			Limit:    10,
		}

		incidents, _, err := repo.List(ctx, params)
		require.NoError(t, err)

		// Should find only the screen incident
		assert.Len(t, incidents, 1)
		assert.Contains(t, incidents[0].Title, "Screen")
	})

	t.Run("Mark SLA Breaches", func(t *testing.T) {
		ctx := context.Background()

		// Clean up from previous test
		testutil.TruncateTables(t, pool, "incidents")

		// Create incidents with different SLA due dates
		now := time.Now()

		// Breached incident (due date in the past)
		breachedIncident := models.Incident{
			ID:          testutil.GenerateTestID(t),
			TenantID:    cfg.TenantID,
			SchoolID:    cfg.SchoolID,
			DeviceID:    "device-001",
			Category:    "hardware",
			Severity:    models.SeverityCritical,
			Status:      models.IncidentInProgress,
			Title:       "Breached incident",
			Description: "This SLA has been breached",
			ReportedBy:  "test@example.com",
			SLADueAt:    now.Add(-1 * time.Hour),
			SLABreached: false,
			CreatedAt:   now.Add(-2 * time.Hour),
			UpdatedAt:   now.Add(-2 * time.Hour),
		}
		err := repo.Create(ctx, breachedIncident)
		require.NoError(t, err)

		// Non-breached incident (due date in the future)
		okIncident := models.Incident{
			ID:          testutil.GenerateTestID(t),
			TenantID:    cfg.TenantID,
			SchoolID:    cfg.SchoolID,
			DeviceID:    "device-001",
			Category:    "hardware",
			Severity:    models.SeverityMedium,
			Status:      models.IncidentNew,
			Title:       "OK incident",
			Description: "This SLA is fine",
			ReportedBy:  "test@example.com",
			SLADueAt:    now.Add(24 * time.Hour),
			SLABreached: false,
			CreatedAt:   now,
			UpdatedAt:   now,
		}
		err = repo.Create(ctx, okIncident)
		require.NoError(t, err)

		// Mark SLA breaches
		count, err := repo.MarkSLABreaches(ctx, now)
		require.NoError(t, err)

		// Should have marked 1 incident
		assert.Equal(t, 1, count)

		// Verify the breach was marked
		retrieved, err := repo.GetByID(ctx, cfg.TenantID, cfg.SchoolID, breachedIncident.ID)
		require.NoError(t, err)
		assert.True(t, retrieved.SLABreached)

		// Verify the other incident is still not breached
		retrieved, err = repo.GetByID(ctx, cfg.TenantID, cfg.SchoolID, okIncident.ID)
		require.NoError(t, err)
		assert.False(t, retrieved.SLABreached)
	})
}

// Example test using transactions for isolation
func TestIncidentRepository_WithTransaction(t *testing.T) {
	testutil.SkipIfNoIntegration(t)

	db := testutil.SetupTestDB(t)
	pool := db.RawPool()
	cfg := testutil.DefaultFixtureConfig()

	testutil.CleanupAllTables(t, pool)
	testutil.CreateSchoolSnapshot(t, pool, cfg)
	testutil.CreateDeviceSnapshot(t, pool, cfg, "device-001")

	t.Run("Transaction rollback doesn't affect database", func(t *testing.T) {
		// This test demonstrates using transactions for test isolation
		testutil.WithTransaction(t, pool, func(tx store.Tx) {
			// Create incident within transaction
			incident := models.Incident{
				ID:          testutil.GenerateTestID(t),
				TenantID:    cfg.TenantID,
				SchoolID:    cfg.SchoolID,
				DeviceID:    "device-001",
				Category:    "hardware",
				Severity:    models.SeverityMedium,
				Status:      models.IncidentNew,
				Title:       "Transaction test",
				Description: "This should be rolled back",
				ReportedBy:  "test@example.com",
				SLADueAt:    time.Now().Add(24 * time.Hour),
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}

			_, err := tx.Exec(context.Background(), `
				INSERT INTO incidents (
					id, tenant_id, school_id, device_id,
					category, severity, status, title, description,
					reported_by, sla_due_at, sla_breached, created_at, updated_at
				) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)
			`, incident.ID, incident.TenantID, incident.SchoolID, incident.DeviceID,
				incident.Category, incident.Severity, incident.Status,
				incident.Title, incident.Description, incident.ReportedBy,
				incident.SLADueAt, incident.SLABreached,
				incident.CreatedAt, incident.UpdatedAt)

			require.NoError(t, err)
			// Transaction will be rolled back automatically
		})

		// Verify nothing was persisted
		repo := db.Incidents()
		params := store.IncidentListParams{
			TenantID: cfg.TenantID,
			SchoolID: cfg.SchoolID,
			Limit:    10,
		}
		incidents, _, err := repo.List(context.Background(), params)
		require.NoError(t, err)
		assert.Len(t, incidents, 0, "transaction should have been rolled back")
	})
}
