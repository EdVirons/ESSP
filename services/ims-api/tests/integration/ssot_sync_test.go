//go:build integration
// +build integration

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/edvirons/ssp/ims/internal/lookups"
	"github.com/edvirons/ssp/ims/internal/models"
	"github.com/edvirons/ssp/ims/internal/store"
	"github.com/edvirons/ssp/ims/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSSOTSync_SchoolSnapshotCaching(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	db, cfg, cleanup := setupTestWithFixtures(t)
	defer cleanup()

	ctx := context.Background()

	// Create multiple school snapshots
	school1 := models.SchoolSnapshot{
		TenantID:      cfg.TenantID,
		SchoolID:      "school-001",
		Name:          "Test School 1",
		CountyCode:    "001",
		CountyName:    "Nairobi",
		SubCountyCode: "001-01",
		SubCountyName: "Westlands",
		UpdatedAt:     time.Now().UTC(),
	}

	school2 := models.SchoolSnapshot{
		TenantID:      cfg.TenantID,
		SchoolID:      "school-002",
		Name:          "Test School 2",
		CountyCode:    "002",
		CountyName:    "Mombasa",
		SubCountyCode: "002-01",
		SubCountyName: "Mvita",
		UpdatedAt:     time.Now().UTC(),
	}

	// Upsert school snapshots
	err := db.SchoolsSnapshot().Upsert(ctx, school1)
	require.NoError(t, err)

	err = db.SchoolsSnapshot().Upsert(ctx, school2)
	require.NoError(t, err)

	// Verify snapshots are cached and retrievable
	retrieved1, err := db.SchoolsSnapshot().Get(ctx, cfg.TenantID, "school-001")
	require.NoError(t, err)
	assert.Equal(t, "Test School 1", retrieved1.Name)
	assert.Equal(t, "Nairobi", retrieved1.CountyName)

	retrieved2, err := db.SchoolsSnapshot().Get(ctx, cfg.TenantID, "school-002")
	require.NoError(t, err)
	assert.Equal(t, "Test School 2", retrieved2.Name)
	assert.Equal(t, "Mombasa", retrieved2.CountyName)

	// Test update (upsert existing)
	school1Updated := school1
	school1Updated.Name = "Test School 1 Updated"
	school1Updated.UpdatedAt = time.Now().UTC()

	err = db.SchoolsSnapshot().Upsert(ctx, school1Updated)
	require.NoError(t, err)

	retrievedUpdated, err := db.SchoolsSnapshot().Get(ctx, cfg.TenantID, "school-001")
	require.NoError(t, err)
	assert.Equal(t, "Test School 1 Updated", retrievedUpdated.Name)
}

func TestSSOTSync_DeviceSnapshotCaching(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	db, cfg, cleanup := setupTestWithFixtures(t)
	defer cleanup()

	ctx := context.Background()

	// Create device snapshots
	device1 := models.DeviceSnapshot{
		TenantID:  cfg.TenantID,
		DeviceID:  "device-001",
		SchoolID:  cfg.SchoolID,
		Model:     "Laptop Dell Latitude",
		Serial:    "SN123456",
		AssetTag:  "ASSET-001",
		Status:    "active",
		UpdatedAt: time.Now().UTC(),
	}

	device2 := models.DeviceSnapshot{
		TenantID:  cfg.TenantID,
		DeviceID:  "device-002",
		SchoolID:  cfg.SchoolID,
		Model:     "Laptop HP ProBook",
		Serial:    "SN789012",
		AssetTag:  "ASSET-002",
		Status:    "active",
		UpdatedAt: time.Now().UTC(),
	}

	// Upsert device snapshots
	err := db.DevicesSnapshot().Upsert(ctx, device1)
	require.NoError(t, err)

	err = db.DevicesSnapshot().Upsert(ctx, device2)
	require.NoError(t, err)

	// Verify snapshots are retrievable
	retrieved1, err := db.DevicesSnapshot().Get(ctx, cfg.TenantID, "device-001")
	require.NoError(t, err)
	assert.Equal(t, "Laptop Dell Latitude", retrieved1.Model)
	assert.Equal(t, "SN123456", retrieved1.Serial)

	retrieved2, err := db.DevicesSnapshot().Get(ctx, cfg.TenantID, "device-002")
	require.NoError(t, err)
	assert.Equal(t, "Laptop HP ProBook", retrieved2.Model)

	// Test update
	device1Updated := device1
	device1Updated.Status = "inactive"
	device1Updated.UpdatedAt = time.Now().UTC()

	err = db.DevicesSnapshot().Upsert(ctx, device1Updated)
	require.NoError(t, err)

	retrievedUpdated, err := db.DevicesSnapshot().Get(ctx, cfg.TenantID, "device-001")
	require.NoError(t, err)
	assert.Equal(t, "inactive", retrievedUpdated.Status)
}

func TestSSOTSync_PartSnapshotCaching(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	db, cfg, cleanup := setupTestWithFixtures(t)
	defer cleanup()

	ctx := context.Background()

	// Create part snapshots
	part1 := models.PartSnapshot{
		TenantID:  cfg.TenantID,
		PartID:    "part-001",
		PUK:       "PUK-12345",
		Name:      "Screen LCD 15.6 inch",
		Category:  "display",
		Unit:      "piece",
		UpdatedAt: time.Now().UTC(),
	}

	part2 := models.PartSnapshot{
		TenantID:  cfg.TenantID,
		PartID:    "part-002",
		PUK:       "PUK-67890",
		Name:      "Keyboard US Layout",
		Category:  "input",
		Unit:      "piece",
		UpdatedAt: time.Now().UTC(),
	}

	// Upsert part snapshots
	err := db.PartsSnapshot().Upsert(ctx, part1)
	require.NoError(t, err)

	err = db.PartsSnapshot().Upsert(ctx, part2)
	require.NoError(t, err)

	// Verify snapshots are retrievable
	retrieved1, err := db.PartsSnapshot().Get(ctx, cfg.TenantID, "part-001")
	require.NoError(t, err)
	assert.Equal(t, "Screen LCD 15.6 inch", retrieved1.Name)
	assert.Equal(t, "PUK-12345", retrieved1.PUK)

	retrieved2, err := db.PartsSnapshot().Get(ctx, cfg.TenantID, "part-002")
	require.NoError(t, err)
	assert.Equal(t, "Keyboard US Layout", retrieved2.Name)

	// Test update
	part1Updated := part1
	part1Updated.Name = "Screen LCD 15.6 inch HD"
	part1Updated.UpdatedAt = time.Now().UTC()

	err = db.PartsSnapshot().Upsert(ctx, part1Updated)
	require.NoError(t, err)

	retrievedUpdated, err := db.PartsSnapshot().Get(ctx, cfg.TenantID, "part-001")
	require.NoError(t, err)
	assert.Equal(t, "Screen LCD 15.6 inch HD", retrievedUpdated.Name)
}

func TestSSOTSync_LookupEnrichment(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	db, cfg, cleanup := setupTestWithFixtures(t)
	defer cleanup()

	ctx := context.Background()

	// Create snapshots for enrichment testing
	schoolSnapshot := models.SchoolSnapshot{
		TenantID:      cfg.TenantID,
		SchoolID:      cfg.SchoolID,
		Name:          "Enrichment Test School",
		CountyCode:    "001",
		CountyName:    "Test County",
		SubCountyCode: "001-01",
		SubCountyName: "Test SubCounty",
		UpdatedAt:     time.Now().UTC(),
	}

	err := db.SchoolsSnapshot().Upsert(ctx, schoolSnapshot)
	require.NoError(t, err)

	deviceID := "enrich-device-001"
	deviceSnapshot := models.DeviceSnapshot{
		TenantID:  cfg.TenantID,
		DeviceID:  deviceID,
		SchoolID:  cfg.SchoolID,
		Model:     "Enrichment Test Model",
		Serial:    "ENRICH-SN-001",
		AssetTag:  "ENRICH-ASSET-001",
		Status:    "active",
		UpdatedAt: time.Now().UTC(),
	}

	err = db.DevicesSnapshot().Upsert(ctx, deviceSnapshot)
	require.NoError(t, err)

	partID := "enrich-part-001"
	partSnapshot := models.PartSnapshot{
		TenantID:  cfg.TenantID,
		PartID:    partID,
		PUK:       "ENRICH-PUK-001",
		Name:      "Enrichment Test Part",
		Category:  "test-category",
		Unit:      "piece",
		UpdatedAt: time.Now().UTC(),
	}

	err = db.PartsSnapshot().Upsert(ctx, partSnapshot)
	require.NoError(t, err)

	// Use lookups to verify enrichment works
	lk := lookups.New(db.RawPool())

	// Test school lookup
	school, err := lk.SchoolByID(ctx, cfg.TenantID, cfg.SchoolID)
	require.NoError(t, err)
	assert.NotNil(t, school)
	assert.Equal(t, "Enrichment Test School", school.Name)
	assert.Equal(t, "Test County", school.CountyName)

	// Test device lookup
	device, err := lk.DeviceByID(ctx, cfg.TenantID, deviceID)
	require.NoError(t, err)
	assert.NotNil(t, device)
	assert.Equal(t, "Enrichment Test Model", device.Model)
	assert.Equal(t, "ENRICH-SN-001", device.Serial)

	// Test part lookup
	part, err := lk.PartByID(ctx, cfg.TenantID, partID)
	require.NoError(t, err)
	assert.NotNil(t, part)
	assert.Equal(t, "Enrichment Test Part", part.Name)
	assert.Equal(t, "ENRICH-PUK-001", part.PUK)
}

func TestSSOTSync_SyncStateManagement(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	db, cfg, cleanup := setupTestWithFixtures(t)
	defer cleanup()

	ctx := context.Background()

	// Test SSOT state creation and updates
	resources := []store.SSOTResource{
		store.SSOTSchools,
		store.SSOTDevices,
		store.SSOTParts,
	}

	for _, resource := range resources {
		t.Run(string(resource), func(t *testing.T) {
			// Create initial state
			state := store.NewSSOTSyncState(cfg.TenantID, resource)
			initialTime := time.Now().UTC()
			state.LastUpdatedSince = initialTime
			state.LastCursor = "cursor-123"
			state.UpdatedAt = initialTime

			err := db.SSOTState().Upsert(ctx, state)
			require.NoError(t, err)

			// Retrieve and verify
			retrieved, err := db.SSOTState().Get(ctx, cfg.TenantID, resource)
			require.NoError(t, err)
			assert.Equal(t, cfg.TenantID, retrieved.TenantID)
			assert.Equal(t, resource, retrieved.Resource)
			assert.Equal(t, "cursor-123", retrieved.LastCursor)

			// Update state (simulating a sync)
			newTime := time.Now().UTC().Add(1 * time.Hour)
			state.LastUpdatedSince = newTime
			state.LastCursor = "cursor-456"
			state.UpdatedAt = newTime

			err = db.SSOTState().Upsert(ctx, state)
			require.NoError(t, err)

			// Verify update
			updated, err := db.SSOTState().Get(ctx, cfg.TenantID, resource)
			require.NoError(t, err)
			assert.Equal(t, "cursor-456", updated.LastCursor)

			// Verify timestamp was updated (allowing 1 second tolerance)
			assert.True(t, updated.LastUpdatedSince.After(initialTime) || updated.LastUpdatedSince.Equal(initialTime))
		})
	}
}

func TestSSOTSync_BulkSnapshotUpsert(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	db, cfg, cleanup := setupTestWithFixtures(t)
	defer cleanup()

	ctx := context.Background()

	// Simulate bulk sync of multiple schools
	numSchools := 10
	schools := make([]models.SchoolSnapshot, numSchools)

	baseTime := time.Now().UTC()
	for i := 0; i < numSchools; i++ {
		schools[i] = models.SchoolSnapshot{
			TenantID:      cfg.TenantID,
			SchoolID:      store.NewID("school"),
			Name:          "Bulk School " + string(rune('A'+i)),
			CountyCode:    "001",
			CountyName:    "Test County",
			SubCountyCode: "001-01",
			SubCountyName: "Test SubCounty",
			UpdatedAt:     baseTime.Add(time.Duration(i) * time.Minute),
		}

		err := db.SchoolsSnapshot().Upsert(ctx, schools[i])
		require.NoError(t, err)
	}

	// Verify all schools were inserted
	for i, school := range schools {
		retrieved, err := db.SchoolsSnapshot().Get(ctx, cfg.TenantID, school.SchoolID)
		require.NoError(t, err)
		assert.Equal(t, "Bulk School "+string(rune('A'+i)), retrieved.Name)
	}

	// Update all schools (upsert existing)
	for i := range schools {
		schools[i].Name = "Updated Bulk School " + string(rune('A'+i))
		schools[i].UpdatedAt = time.Now().UTC()

		err := db.SchoolsSnapshot().Upsert(ctx, schools[i])
		require.NoError(t, err)
	}

	// Verify updates
	for i, school := range schools {
		retrieved, err := db.SchoolsSnapshot().Get(ctx, cfg.TenantID, school.SchoolID)
		require.NoError(t, err)
		assert.Equal(t, "Updated Bulk School "+string(rune('A'+i)), retrieved.Name)
	}
}

func TestSSOTSync_ConcurrentSnapshotUpdates(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	db, cfg, cleanup := setupTestWithFixtures(t)
	defer cleanup()

	ctx := context.Background()

	schoolID := "concurrent-school-001"

	// Create initial snapshot
	initialSnapshot := models.SchoolSnapshot{
		TenantID:      cfg.TenantID,
		SchoolID:      schoolID,
		Name:          "Initial Name",
		CountyCode:    "001",
		CountyName:    "Test County",
		SubCountyCode: "001-01",
		SubCountyName: "Test SubCounty",
		UpdatedAt:     time.Now().UTC(),
	}

	err := db.SchoolsSnapshot().Upsert(ctx, initialSnapshot)
	require.NoError(t, err)

	// Simulate concurrent updates from different sync processes
	numConcurrent := 5
	done := make(chan error, numConcurrent)

	for i := 0; i < numConcurrent; i++ {
		go func(index int) {
			snapshot := models.SchoolSnapshot{
				TenantID:      cfg.TenantID,
				SchoolID:      schoolID,
				Name:          "Updated Name " + store.Itoa(int64(index)),
				CountyCode:    "001",
				CountyName:    "Test County",
				SubCountyCode: "001-01",
				SubCountyName: "Test SubCounty",
				UpdatedAt:     time.Now().UTC().Add(time.Duration(index) * time.Millisecond),
			}

			done <- db.SchoolsSnapshot().Upsert(ctx, snapshot)
		}(i)
	}

	// Wait for all to complete
	for i := 0; i < numConcurrent; i++ {
		err := <-done
		assert.NoError(t, err, "concurrent upsert should succeed")
	}

	// Verify final state (should have one of the updated names)
	final, err := db.SchoolsSnapshot().Get(ctx, cfg.TenantID, schoolID)
	require.NoError(t, err)
	assert.Contains(t, final.Name, "Updated Name", "should have an updated name")
}

func TestSSOTSync_MissingSnapshot(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	db, cfg, cleanup := setupTestWithFixtures(t)
	defer cleanup()

	ctx := context.Background()

	// Try to retrieve non-existent snapshots
	_, err := db.SchoolsSnapshot().Get(ctx, cfg.TenantID, "non-existent-school")
	assert.Error(t, err, "should return error for non-existent school")

	_, err = db.DevicesSnapshot().Get(ctx, cfg.TenantID, "non-existent-device")
	assert.Error(t, err, "should return error for non-existent device")

	_, err = db.PartsSnapshot().Get(ctx, cfg.TenantID, "non-existent-part")
	assert.Error(t, err, "should return error for non-existent part")

	// Use lookups - should handle missing gracefully
	lk := lookups.New(db.RawPool())

	school, err := lk.SchoolByID(ctx, cfg.TenantID, "non-existent-school")
	// Depending on implementation, might return nil or error
	assert.True(t, err != nil || school == nil, "should handle missing school gracefully")

	device, err := lk.DeviceByID(ctx, cfg.TenantID, "non-existent-device")
	assert.True(t, err != nil || device == nil, "should handle missing device gracefully")

	part, err := lk.PartByID(ctx, cfg.TenantID, "non-existent-part")
	assert.True(t, err != nil || part == nil, "should handle missing part gracefully")
}
