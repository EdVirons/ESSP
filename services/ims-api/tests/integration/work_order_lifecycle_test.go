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

func TestWorkOrderLifecycle_CompleteWorkflow(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	db, cfg, cleanup := setupTestWithFixtures(t)
	defer cleanup()

	ctx := context.Background()

	// Setup: Create necessary entities
	deviceID := "test-device-wo-001"
	testutil.CreateDeviceSnapshot(t, db.RawPool(), cfg, deviceID)

	incident := testutil.CreateIncident(t, db.RawPool(), cfg, deviceID)
	shop := testutil.CreateServiceShop(t, db.RawPool(), cfg)
	staff := testutil.CreateServiceStaff(t, db.RawPool(), cfg, shop.ID)

	// Step 1: Create work order from incident
	now := time.Now().UTC()
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
		TaskType:        "repair",
		Notes:           "Repair screen malfunction",
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	err := db.WorkOrders().Create(ctx, workOrder)
	require.NoError(t, err, "should create work order")

	// Step 2: Verify work order was created
	wo, err := db.WorkOrders().GetByID(ctx, cfg.TenantID, cfg.SchoolID, workOrder.ID)
	require.NoError(t, err)
	assert.Equal(t, models.WorkOrderDraft, wo.Status)

	// Step 3: Add parts to BOM
	part := testutil.CreatePart(t, db.RawPool(), cfg)
	inventory := testutil.CreateInventoryItem(t, db.RawPool(), cfg, shop.ID, part.ID, 100)

	bomItem := models.WorkOrderPart{
		ID:            store.NewID("bom"),
		TenantID:      cfg.TenantID,
		SchoolID:      cfg.SchoolID,
		WorkOrderID:   workOrder.ID,
		ServiceShopID: shop.ID,
		PartID:        part.ID,
		PartName:      part.Name,
		PartCategory:  part.Category,
		QtyPlanned:    5,
		QtyUsed:       0,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	// Reserve inventory
	execSQL(t, db.RawPool(), `
		UPDATE inventory
		SET qty_reserved = qty_reserved + $3, updated_at = $4
		WHERE tenant_id=$1 AND service_shop_id=$2 AND part_id=$5
	`, cfg.TenantID, shop.ID, bomItem.QtyPlanned, now, part.ID)

	err = store.CreateWorkOrderPartTx(ctx, db.RawPool(), bomItem)
	require.NoError(t, err, "should add part to BOM")

	// Step 4: Schedule work order
	schedStart := time.Now().UTC().Add(24 * time.Hour)
	schedEnd := schedStart.Add(2 * time.Hour)
	schedule := models.WorkOrderSchedule{
		ID:              store.NewID("sched"),
		TenantID:        cfg.TenantID,
		SchoolID:        cfg.SchoolID,
		WorkOrderID:     workOrder.ID,
		ScheduledStart:  &schedStart,
		ScheduledEnd:    &schedEnd,
		Timezone:        "Africa/Nairobi",
		Notes:           "Scheduled for tomorrow",
		CreatedByUserID: staff.UserID,
		CreatedAt:       now,
	}

	err = db.WorkOrderSchedules().Create(ctx, schedule)
	require.NoError(t, err, "should create schedule")

	// Step 5: Transition work order to Assigned
	updated, err := db.WorkOrders().UpdateStatus(ctx, cfg.TenantID, cfg.SchoolID, workOrder.ID, models.WorkOrderAssigned, now)
	require.NoError(t, err)
	assert.Equal(t, models.WorkOrderAssigned, updated.Status)

	// Step 6: Transition to InRepair and consume parts
	updated, err = db.WorkOrders().UpdateStatus(ctx, cfg.TenantID, cfg.SchoolID, workOrder.ID, models.WorkOrderInRepair, now)
	require.NoError(t, err)
	assert.Equal(t, models.WorkOrderInRepair, updated.Status)

	// Consume parts (using 3 out of 5 planned)
	qtyUsed := int64(3)
	err = store.UpdateWorkOrderPartUsedTx(ctx, db.RawPool(), cfg.TenantID, cfg.SchoolID, bomItem.ID, qtyUsed, now)
	require.NoError(t, err)

	// Update inventory
	execSQL(t, db.RawPool(), `
		UPDATE inventory
		SET qty_reserved = GREATEST(qty_reserved - $3, 0),
			qty_available = GREATEST(qty_available - $3, 0),
			updated_at=$4
		WHERE tenant_id=$1 AND service_shop_id=$2 AND part_id=$5
	`, cfg.TenantID, shop.ID, qtyUsed, now, part.ID)

	// Verify parts were consumed
	updatedBOM, err := db.WorkOrderParts().GetByID(ctx, cfg.TenantID, cfg.SchoolID, bomItem.ID)
	require.NoError(t, err)
	assert.Equal(t, qtyUsed, updatedBOM.QtyUsed)

	// Verify inventory was updated
	updatedInv, err := db.Inventory().GetByPartID(ctx, cfg.TenantID, shop.ID, part.ID)
	require.NoError(t, err)
	assert.Equal(t, inventory.QtyAvailable-qtyUsed, updatedInv.QtyAvailable)
	assert.Equal(t, inventory.QtyReserved-qtyUsed, updatedInv.QtyReserved)

	// Step 7: Add deliverable
	deliverable := models.WorkOrderDeliverable{
		ID:          store.NewID("deliv"),
		TenantID:    cfg.TenantID,
		SchoolID:    cfg.SchoolID,
		WorkOrderID: workOrder.ID,
		Title:       "Screen replacement evidence",
		Description: "Photo of replaced screen",
		Status:      models.DeliverablePending,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	err = db.WorkOrderDeliverables().Create(ctx, deliverable)
	require.NoError(t, err, "should create deliverable")

	// Step 8: Submit deliverable
	evidenceID := "attachment-001"
	err = db.WorkOrderDeliverables().MarkSubmitted(ctx, cfg.TenantID, cfg.SchoolID, deliverable.ID, staff.UserID, evidenceID, "Screen replaced successfully")
	require.NoError(t, err, "should submit deliverable")

	submittedDeliv, err := db.WorkOrderDeliverables().GetByID(ctx, cfg.TenantID, cfg.SchoolID, deliverable.ID)
	require.NoError(t, err)
	assert.Equal(t, models.DeliverableSubmitted, submittedDeliv.Status)
	assert.Equal(t, evidenceID, submittedDeliv.EvidenceAttachmentID)

	// Step 9: Review and approve deliverable
	err = db.WorkOrderDeliverables().Review(ctx, cfg.TenantID, cfg.SchoolID, deliverable.ID, "reviewer-user", "approved", "Looks good")
	require.NoError(t, err, "should approve deliverable")

	approvedDeliv, err := db.WorkOrderDeliverables().GetByID(ctx, cfg.TenantID, cfg.SchoolID, deliverable.ID)
	require.NoError(t, err)
	assert.Equal(t, models.DeliverableApproved, approvedDeliv.Status)

	// Step 10: Transition to QA
	updated, err = db.WorkOrders().UpdateStatus(ctx, cfg.TenantID, cfg.SchoolID, workOrder.ID, models.WorkOrderQA, now)
	require.NoError(t, err)
	assert.Equal(t, models.WorkOrderQA, updated.Status)

	// Step 11: Transition to Completed
	updated, err = db.WorkOrders().UpdateStatus(ctx, cfg.TenantID, cfg.SchoolID, workOrder.ID, models.WorkOrderCompleted, now)
	require.NoError(t, err)
	assert.Equal(t, models.WorkOrderCompleted, updated.Status)

	// Step 12: Request approval
	approval := models.WorkOrderApproval{
		ID:                store.NewID("appr"),
		TenantID:          cfg.TenantID,
		SchoolID:          cfg.SchoolID,
		WorkOrderID:       workOrder.ID,
		ApprovalType:      "school_signoff",
		RequestedByUserID: staff.UserID,
		RequestedAt:       now,
		Status:            models.ApprovalPending,
	}

	err = db.WorkOrderApprovals().Request(ctx, approval)
	require.NoError(t, err, "should request approval")

	// Step 13: Approve the work order
	err = db.WorkOrderApprovals().Decide(ctx, cfg.TenantID, cfg.SchoolID, approval.ID, "approver-user", "approved", "Work completed satisfactorily")
	require.NoError(t, err, "should approve work order")

	approvedAppr, err := db.WorkOrderApprovals().GetByID(ctx, cfg.TenantID, cfg.SchoolID, approval.ID)
	require.NoError(t, err)
	assert.Equal(t, models.ApprovalApproved, approvedAppr.Status)

	// Step 14: Transition to Approved (final state)
	updated, err = db.WorkOrders().UpdateStatus(ctx, cfg.TenantID, cfg.SchoolID, workOrder.ID, models.WorkOrderApproved, now)
	require.NoError(t, err)
	assert.Equal(t, models.WorkOrderApproved, updated.Status)

	// Final verification: Check complete workflow
	finalWO, err := db.WorkOrders().GetByID(ctx, cfg.TenantID, cfg.SchoolID, workOrder.ID)
	require.NoError(t, err)
	assert.Equal(t, models.WorkOrderApproved, finalWO.Status)
}

func TestWorkOrderLifecycle_InvalidTransitions(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	db, cfg, cleanup := setupTestWithFixtures(t)
	defer cleanup()

	ctx := context.Background()

	deviceID := "test-device-wo-002"
	testutil.CreateDeviceSnapshot(t, db.RawPool(), cfg, deviceID)
	incident := testutil.CreateIncident(t, db.RawPool(), cfg, deviceID)
	shop := testutil.CreateServiceShop(t, db.RawPool(), cfg)

	workOrder := testutil.CreateWorkOrder(t, db.RawPool(), cfg, incident.ID, deviceID)

	// Test invalid transitions
	invalidTransitions := []struct {
		from models.WorkOrderStatus
		to   models.WorkOrderStatus
	}{
		{models.WorkOrderDraft, models.WorkOrderInRepair},     // Can't skip assigned
		{models.WorkOrderDraft, models.WorkOrderCompleted},    // Can't skip to completed
		{models.WorkOrderAssigned, models.WorkOrderCompleted}, // Can't skip in_repair
		{models.WorkOrderCompleted, models.WorkOrderInRepair}, // Can't go backwards
		{models.WorkOrderApproved, models.WorkOrderDraft},     // Can't go backwards from approved
	}

	for _, tc := range invalidTransitions {
		assert.False(t, service.CanTransitionWorkOrder(tc.from, tc.to),
			"should not allow transition from %s to %s", tc.from, tc.to)
	}

	// Verify valid transitions are allowed
	validTransitions := []struct {
		from models.WorkOrderStatus
		to   models.WorkOrderStatus
	}{
		{models.WorkOrderDraft, models.WorkOrderAssigned},
		{models.WorkOrderAssigned, models.WorkOrderInRepair},
		{models.WorkOrderInRepair, models.WorkOrderQA},
		{models.WorkOrderInRepair, models.WorkOrderCompleted}, // Can skip QA
		{models.WorkOrderQA, models.WorkOrderCompleted},
		{models.WorkOrderCompleted, models.WorkOrderApproved},
	}

	for _, tc := range validTransitions {
		assert.True(t, service.CanTransitionWorkOrder(tc.from, tc.to),
			"should allow transition from %s to %s", tc.from, tc.to)
	}
}

func TestWorkOrderLifecycle_ScheduleManagement(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	db, cfg, cleanup := setupTestWithFixtures(t)
	defer cleanup()

	ctx := context.Background()

	deviceID := "test-device-wo-003"
	testutil.CreateDeviceSnapshot(t, db.RawPool(), cfg, deviceID)
	incident := testutil.CreateIncident(t, db.RawPool(), cfg, deviceID)
	workOrder := testutil.CreateWorkOrder(t, db.RawPool(), cfg, incident.ID, deviceID)

	now := time.Now().UTC()

	// Create multiple schedules for the same work order
	schedule1Start := now.Add(24 * time.Hour)
	schedule1End := schedule1Start.Add(2 * time.Hour)

	schedule1 := models.WorkOrderSchedule{
		ID:              store.NewID("sched"),
		TenantID:        cfg.TenantID,
		SchoolID:        cfg.SchoolID,
		WorkOrderID:     workOrder.ID,
		ScheduledStart:  &schedule1Start,
		ScheduledEnd:    &schedule1End,
		Timezone:        "Africa/Nairobi",
		Notes:           "First schedule",
		CreatedByUserID: "user-1",
		CreatedAt:       now,
	}

	err := db.WorkOrderSchedules().Create(ctx, schedule1)
	require.NoError(t, err)

	// Create a rescheduled entry
	schedule2Start := now.Add(48 * time.Hour)
	schedule2End := schedule2Start.Add(2 * time.Hour)

	schedule2 := models.WorkOrderSchedule{
		ID:              store.NewID("sched"),
		TenantID:        cfg.TenantID,
		SchoolID:        cfg.SchoolID,
		WorkOrderID:     workOrder.ID,
		ScheduledStart:  &schedule2Start,
		ScheduledEnd:    &schedule2End,
		Timezone:        "Africa/Nairobi",
		Notes:           "Rescheduled",
		CreatedByUserID: "user-1",
		CreatedAt:       now.Add(1 * time.Hour),
	}

	err = db.WorkOrderSchedules().Create(ctx, schedule2)
	require.NoError(t, err)

	// List schedules
	schedules, _, err := db.WorkOrderSchedules().List(ctx, store.ScheduleListParams{
		TenantID:    cfg.TenantID,
		SchoolID:    cfg.SchoolID,
		WorkOrderID: workOrder.ID,
		Limit:       10,
	})
	require.NoError(t, err)
	assert.Len(t, schedules, 2, "should have two schedule entries")
}

func TestWorkOrderLifecycle_DeliverablesWorkflow(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	db, cfg, cleanup := setupTestWithFixtures(t)
	defer cleanup()

	ctx := context.Background()

	deviceID := "test-device-wo-004"
	testutil.CreateDeviceSnapshot(t, db.RawPool(), cfg, deviceID)
	incident := testutil.CreateIncident(t, db.RawPool(), cfg, deviceID)
	workOrder := testutil.CreateWorkOrder(t, db.RawPool(), cfg, incident.ID, deviceID)

	now := time.Now().UTC()

	// Create deliverable
	deliverable := models.WorkOrderDeliverable{
		ID:          store.NewID("deliv"),
		TenantID:    cfg.TenantID,
		SchoolID:    cfg.SchoolID,
		WorkOrderID: workOrder.ID,
		Title:       "Test deliverable",
		Description: "Test description",
		Status:      models.DeliverablePending,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	err := db.WorkOrderDeliverables().Create(ctx, deliverable)
	require.NoError(t, err)

	// Test submission
	err = db.WorkOrderDeliverables().MarkSubmitted(ctx, cfg.TenantID, cfg.SchoolID, deliverable.ID, "submitter-user", "evidence-id", "Submitted")
	require.NoError(t, err)

	submitted, err := db.WorkOrderDeliverables().GetByID(ctx, cfg.TenantID, cfg.SchoolID, deliverable.ID)
	require.NoError(t, err)
	assert.Equal(t, models.DeliverableSubmitted, submitted.Status)
	assert.NotNil(t, submitted.SubmittedAt)

	// Test rejection
	err = db.WorkOrderDeliverables().Review(ctx, cfg.TenantID, cfg.SchoolID, deliverable.ID, "reviewer-user", "rejected", "Please redo")
	require.NoError(t, err)

	rejected, err := db.WorkOrderDeliverables().GetByID(ctx, cfg.TenantID, cfg.SchoolID, deliverable.ID)
	require.NoError(t, err)
	assert.Equal(t, models.DeliverableRejected, rejected.Status)
	assert.NotNil(t, rejected.ReviewedAt)
	assert.Equal(t, "Please redo", rejected.ReviewNotes)

	// Resubmit after rejection
	err = db.WorkOrderDeliverables().MarkSubmitted(ctx, cfg.TenantID, cfg.SchoolID, deliverable.ID, "submitter-user", "evidence-id-2", "Resubmitted with fixes")
	require.NoError(t, err)

	// Approve on second review
	err = db.WorkOrderDeliverables().Review(ctx, cfg.TenantID, cfg.SchoolID, deliverable.ID, "reviewer-user", "approved", "Good now")
	require.NoError(t, err)

	approved, err := db.WorkOrderDeliverables().GetByID(ctx, cfg.TenantID, cfg.SchoolID, deliverable.ID)
	require.NoError(t, err)
	assert.Equal(t, models.DeliverableApproved, approved.Status)
}

func TestWorkOrderLifecycle_ApprovalWorkflow(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	db, cfg, cleanup := setupTestWithFixtures(t)
	defer cleanup()

	ctx := context.Background()

	deviceID := "test-device-wo-005"
	testutil.CreateDeviceSnapshot(t, db.RawPool(), cfg, deviceID)
	incident := testutil.CreateIncident(t, db.RawPool(), cfg, deviceID)
	workOrder := testutil.CreateWorkOrder(t, db.RawPool(), cfg, incident.ID, deviceID)

	now := time.Now().UTC()

	// Request approval
	approval := models.WorkOrderApproval{
		ID:                store.NewID("appr"),
		TenantID:          cfg.TenantID,
		SchoolID:          cfg.SchoolID,
		WorkOrderID:       workOrder.ID,
		ApprovalType:      "school_signoff",
		RequestedByUserID: "requester-user",
		RequestedAt:       now,
		Status:            models.ApprovalPending,
	}

	err := db.WorkOrderApprovals().Request(ctx, approval)
	require.NoError(t, err)

	retrieved, err := db.WorkOrderApprovals().GetByID(ctx, cfg.TenantID, cfg.SchoolID, approval.ID)
	require.NoError(t, err)
	assert.Equal(t, models.ApprovalPending, retrieved.Status)

	// Test rejection
	err = db.WorkOrderApprovals().Decide(ctx, cfg.TenantID, cfg.SchoolID, approval.ID, "approver-user", "rejected", "Needs more work")
	require.NoError(t, err)

	rejected, err := db.WorkOrderApprovals().GetByID(ctx, cfg.TenantID, cfg.SchoolID, approval.ID)
	require.NoError(t, err)
	assert.Equal(t, models.ApprovalRejected, rejected.Status)
	assert.NotNil(t, rejected.DecidedAt)
	assert.Equal(t, "Needs more work", rejected.DecisionNotes)

	// Request new approval after fixes
	approval2 := models.WorkOrderApproval{
		ID:                store.NewID("appr"),
		TenantID:          cfg.TenantID,
		SchoolID:          cfg.SchoolID,
		WorkOrderID:       workOrder.ID,
		ApprovalType:      "school_signoff",
		RequestedByUserID: "requester-user",
		RequestedAt:       now.Add(1 * time.Hour),
		Status:            models.ApprovalPending,
	}

	err = db.WorkOrderApprovals().Request(ctx, approval2)
	require.NoError(t, err)

	// Approve
	err = db.WorkOrderApprovals().Decide(ctx, cfg.TenantID, cfg.SchoolID, approval2.ID, "approver-user", "approved", "Approved")
	require.NoError(t, err)

	approved, err := db.WorkOrderApprovals().GetByID(ctx, cfg.TenantID, cfg.SchoolID, approval2.ID)
	require.NoError(t, err)
	assert.Equal(t, models.ApprovalApproved, approved.Status)
}

func TestWorkOrderLifecycle_MultipleWorkOrders(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	db, cfg, cleanup := setupTestWithFixtures(t)
	defer cleanup()

	ctx := context.Background()

	deviceID := "test-device-wo-006"
	testutil.CreateDeviceSnapshot(t, db.RawPool(), cfg, deviceID)
	incident := testutil.CreateIncident(t, db.RawPool(), cfg, deviceID)

	// Create multiple work orders for the same incident
	numWorkOrders := 3
	workOrderIDs := make([]string, numWorkOrders)

	for i := 0; i < numWorkOrders; i++ {
		wo := testutil.CreateWorkOrder(t, db.RawPool(), cfg, incident.ID, deviceID)
		workOrderIDs[i] = wo.ID
	}

	// Query work orders by incident
	workOrders, _, err := db.WorkOrders().List(ctx, store.WorkOrderListParams{
		TenantID:   cfg.TenantID,
		SchoolID:   cfg.SchoolID,
		IncidentID: incident.ID,
		Limit:      10,
	})
	require.NoError(t, err)
	assert.Len(t, workOrders, numWorkOrders)

	// Verify all work orders are for the same incident
	for _, wo := range workOrders {
		assert.Equal(t, incident.ID, wo.IncidentID)
	}
}
