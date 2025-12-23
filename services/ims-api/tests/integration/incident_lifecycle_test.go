//go:build integration
// +build integration

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/edvirons/ssp/ims/internal/models"
	"github.com/edvirons/ssp/ims/internal/service"
	"github.com/edvirons/ssp/ims/internal/store"
	"github.com/edvirons/ssp/ims/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIncidentLifecycle_HappyPath(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	db, cfg, cleanup := setupTestWithFixtures(t)
	defer cleanup()

	ctx := context.Background()

	// Step 1: Create a device snapshot
	deviceID := "test-device-001"
	testutil.CreateDeviceSnapshot(t, db.RawPool(), cfg, deviceID)

	// Step 2: Create an incident
	now := time.Now().UTC()
	incident := models.Incident{
		ID:          store.NewID("inc"),
		TenantID:    cfg.TenantID,
		SchoolID:    cfg.SchoolID,
		DeviceID:    deviceID,
		SchoolName:  "Test School",
		Category:    "hardware",
		Severity:    models.SeverityMedium,
		Status:      models.IncidentNew,
		Title:       "Screen malfunction",
		Description: "Device screen is flickering",
		ReportedBy:  "test-user",
		SLADueAt:    service.SLADue(models.SeverityMedium, now),
		SLABreached: false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	err := db.Incidents().Create(ctx, incident)
	require.NoError(t, err, "should create incident successfully")

	// Step 3: Verify incident was created
	retrieved, err := db.Incidents().GetByID(ctx, cfg.TenantID, cfg.SchoolID, incident.ID)
	require.NoError(t, err, "should retrieve incident")
	assert.Equal(t, incident.ID, retrieved.ID)
	assert.Equal(t, models.IncidentNew, retrieved.Status)
	assert.Equal(t, "Screen malfunction", retrieved.Title)

	// Step 4: Transition through valid status changes
	validTransitions := []models.IncidentStatus{
		models.IncidentAcknowledged,
		models.IncidentInProgress,
		models.IncidentResolved,
		models.IncidentClosed,
	}

	for _, newStatus := range validTransitions {
		// Verify transition is allowed
		assert.True(t, service.CanTransitionIncident(retrieved.Status, newStatus),
			"should allow transition from %s to %s", retrieved.Status, newStatus)

		// Update status
		updated, err := db.Incidents().UpdateStatus(ctx, cfg.TenantID, cfg.SchoolID, incident.ID, newStatus, time.Now().UTC())
		require.NoError(t, err, "should update status to %s", newStatus)
		assert.Equal(t, newStatus, updated.Status)

		retrieved = updated
	}

	// Step 5: Verify final state
	final, err := db.Incidents().GetByID(ctx, cfg.TenantID, cfg.SchoolID, incident.ID)
	require.NoError(t, err)
	assert.Equal(t, models.IncidentClosed, final.Status)
}

func TestIncidentLifecycle_InvalidTransition(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	db, cfg, cleanup := setupTestWithFixtures(t)
	defer cleanup()

	ctx := context.Background()

	// Create incident
	deviceID := "test-device-002"
	testutil.CreateDeviceSnapshot(t, db.RawPool(), cfg, deviceID)

	incident := testutil.CreateIncident(t, db.RawPool(), cfg, deviceID)

	// Attempt invalid transition (new -> resolved, skipping intermediate steps)
	assert.False(t, service.CanTransitionIncident(models.IncidentNew, models.IncidentResolved),
		"should not allow direct transition from new to resolved")

	// Attempt invalid transition (new -> closed)
	assert.False(t, service.CanTransitionIncident(models.IncidentNew, models.IncidentClosed),
		"should not allow direct transition from new to closed")

	// Attempt transition from closed (terminal state)
	_, err := db.Incidents().UpdateStatus(ctx, cfg.TenantID, cfg.SchoolID, incident.ID, models.IncidentAcknowledged, time.Now().UTC())
	require.NoError(t, err) // First transition to acknowledged

	_, err = db.Incidents().UpdateStatus(ctx, cfg.TenantID, cfg.SchoolID, incident.ID, models.IncidentInProgress, time.Now().UTC())
	require.NoError(t, err)

	_, err = db.Incidents().UpdateStatus(ctx, cfg.TenantID, cfg.SchoolID, incident.ID, models.IncidentResolved, time.Now().UTC())
	require.NoError(t, err)

	_, err = db.Incidents().UpdateStatus(ctx, cfg.TenantID, cfg.SchoolID, incident.ID, models.IncidentClosed, time.Now().UTC())
	require.NoError(t, err)

	// Now try to transition from closed
	assert.False(t, service.CanTransitionIncident(models.IncidentClosed, models.IncidentNew),
		"should not allow transition from closed state")
}

func TestIncidentLifecycle_SLABreachDetection(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	db, cfg, cleanup := setupTestWithFixtures(t)
	defer cleanup()

	ctx := context.Background()

	deviceID := "test-device-003"
	testutil.CreateDeviceSnapshot(t, db.RawPool(), cfg, deviceID)

	// Test different severity levels and their SLA times
	testCases := []struct {
		severity    models.Severity
		expectedDue time.Duration
		description string
	}{
		{models.SeverityCritical, 4 * time.Hour, "Critical should have 4 hour SLA"},
		{models.SeverityHigh, 24 * time.Hour, "High should have 24 hour SLA"},
		{models.SeverityMedium, 48 * time.Hour, "Medium should have 48 hour SLA"},
		{models.SeverityLow, 72 * time.Hour, "Low should have 72 hour SLA"},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			now := time.Now().UTC()
			expectedDueAt := service.SLADue(tc.severity, now)

			incident := models.Incident{
				ID:          store.NewID("inc"),
				TenantID:    cfg.TenantID,
				SchoolID:    cfg.SchoolID,
				DeviceID:    deviceID,
				SchoolName:  "Test School",
				Category:    "hardware",
				Severity:    tc.severity,
				Status:      models.IncidentNew,
				Title:       "Test incident for SLA",
				Description: "Testing SLA calculation",
				ReportedBy:  "test-user",
				SLADueAt:    expectedDueAt,
				SLABreached: false,
				CreatedAt:   now,
				UpdatedAt:   now,
			}

			err := db.Incidents().Create(ctx, incident)
			require.NoError(t, err)

			// Verify SLA due time is calculated correctly
			retrieved, err := db.Incidents().GetByID(ctx, cfg.TenantID, cfg.SchoolID, incident.ID)
			require.NoError(t, err)

			// Allow 1 second tolerance for time comparison
			actualDuration := retrieved.SLADueAt.Sub(retrieved.CreatedAt)
			assert.InDelta(t, tc.expectedDue.Seconds(), actualDuration.Seconds(), 1.0,
				"SLA due time should be approximately %v after creation", tc.expectedDue)
		})
	}
}

func TestIncidentLifecycle_Escalation(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	db, cfg, cleanup := setupTestWithFixtures(t)
	defer cleanup()

	ctx := context.Background()

	deviceID := "test-device-004"
	testutil.CreateDeviceSnapshot(t, db.RawPool(), cfg, deviceID)

	incident := testutil.CreateIncident(t, db.RawPool(), cfg, deviceID)

	// Test escalation from different states
	escalationPaths := []struct {
		fromStatus models.IncidentStatus
		toStatus   models.IncidentStatus
		shouldWork bool
	}{
		{models.IncidentNew, models.IncidentEscalated, true},
		{models.IncidentAcknowledged, models.IncidentEscalated, true},
		{models.IncidentInProgress, models.IncidentEscalated, true},
	}

	for _, path := range escalationPaths {
		t.Run(string(path.fromStatus)+"_to_"+string(path.toStatus), func(t *testing.T) {
			// Create fresh incident
			newIncident := testutil.CreateIncident(t, db.RawPool(), cfg, deviceID)

			// Set to initial state (if not already new)
			if path.fromStatus != models.IncidentNew {
				// Transition to the starting state
				_, err := db.Incidents().UpdateStatus(ctx, cfg.TenantID, cfg.SchoolID, newIncident.ID, path.fromStatus, time.Now().UTC())
				require.NoError(t, err)
			}

			// Verify transition is allowed
			assert.Equal(t, path.shouldWork, service.CanTransitionIncident(path.fromStatus, path.toStatus))

			if path.shouldWork {
				// Perform escalation
				updated, err := db.Incidents().UpdateStatus(ctx, cfg.TenantID, cfg.SchoolID, newIncident.ID, path.toStatus, time.Now().UTC())
				require.NoError(t, err)
				assert.Equal(t, path.toStatus, updated.Status)
			}
		})
	}
}

func TestIncidentLifecycle_WithAutoWorkOrder(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	db, cfg, cleanup := setupTestWithFixtures(t)
	defer cleanup()

	ctx := context.Background()

	// Create device snapshot
	deviceID := "test-device-005"
	testutil.CreateDeviceSnapshot(t, db.RawPool(), cfg, deviceID)

	// Create service shop and staff
	shop := testutil.CreateServiceShop(t, db.RawPool(), cfg)
	staff := testutil.CreateServiceStaff(t, db.RawPool(), cfg, shop.ID)

	// Create incident
	incident := testutil.CreateIncident(t, db.RawPool(), cfg, deviceID)

	// Manually create auto-generated work order (simulating auto-route behavior)
	workOrder := models.WorkOrder{
		ID:              store.NewID("wo"),
		IncidentID:      incident.ID,
		TenantID:        cfg.TenantID,
		SchoolID:        cfg.SchoolID,
		DeviceID:        deviceID,
		SchoolName:      "Test School",
		Status:          models.WorkOrderDraft,
		ServiceShopID:   shop.ID,
		AssignedStaffID: staff.ID,
		RepairLocation:  models.RepairLocationServiceShop,
		AssignedTo:      staff.UserID,
		TaskType:        "triage",
		Notes:           "Auto-created from incident " + incident.ID,
		CreatedAt:       time.Now().UTC(),
		UpdatedAt:       time.Now().UTC(),
	}

	err := db.WorkOrders().Create(ctx, workOrder)
	require.NoError(t, err)

	// Verify work order was created and linked to incident
	wo, err := db.WorkOrders().GetByID(ctx, cfg.TenantID, cfg.SchoolID, workOrder.ID)
	require.NoError(t, err)
	assert.Equal(t, incident.ID, wo.IncidentID)
	assert.Equal(t, models.WorkOrderDraft, wo.Status)
	assert.Equal(t, shop.ID, wo.ServiceShopID)

	// Verify we can query work orders by incident ID
	workOrders, _, err := db.WorkOrders().List(ctx, store.WorkOrderListParams{
		TenantID:   cfg.TenantID,
		SchoolID:   cfg.SchoolID,
		IncidentID: incident.ID,
		Limit:      10,
	})
	require.NoError(t, err)
	assert.Len(t, workOrders, 1)
	assert.Equal(t, workOrder.ID, workOrders[0].ID)
}

func TestIncidentLifecycle_ConcurrentUpdates(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	db, cfg, cleanup := setupTestWithFixtures(t)
	defer cleanup()

	ctx := context.Background()

	deviceID := "test-device-006"
	testutil.CreateDeviceSnapshot(t, db.RawPool(), cfg, deviceID)

	incident := testutil.CreateIncident(t, db.RawPool(), cfg, deviceID)

	// Simulate concurrent status updates
	// Both goroutines try to update from new -> acknowledged
	// Only one should succeed (or both if database allows it)

	done := make(chan error, 2)

	go func() {
		_, err := db.Incidents().UpdateStatus(ctx, cfg.TenantID, cfg.SchoolID, incident.ID, models.IncidentAcknowledged, time.Now().UTC())
		done <- err
	}()

	go func() {
		_, err := db.Incidents().UpdateStatus(ctx, cfg.TenantID, cfg.SchoolID, incident.ID, models.IncidentAcknowledged, time.Now().UTC())
		done <- err
	}()

	// Wait for both to complete
	err1 := <-done
	err2 := <-done

	// At least one should succeed (both may succeed as they're doing the same transition)
	assert.True(t, err1 == nil || err2 == nil, "at least one concurrent update should succeed")

	// Verify final state
	final, err := db.Incidents().GetByID(ctx, cfg.TenantID, cfg.SchoolID, incident.ID)
	require.NoError(t, err)
	assert.Equal(t, models.IncidentAcknowledged, final.Status)
}
