// +build integration

package integration

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

func TestBOMOperations_AddPartsToBOM(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	db, cfg, cleanup := setupTestWithFixtures(t)
	defer cleanup()

	ctx := context.Background()

	// Setup: Create work order, part, and inventory
	deviceID := "test-device-bom-001"
	testutil.CreateDeviceSnapshot(t, db.RawPool(), cfg, deviceID)
	incident := testutil.CreateIncident(t, db.RawPool(), cfg, deviceID)
	shop := testutil.CreateServiceShop(t, db.RawPool(), cfg)
	workOrder := testutil.CreateWorkOrder(t, db.RawPool(), cfg, incident.ID, deviceID)

	// Update work order to have service shop
	execSQL(t, db.RawPool(), `
		UPDATE work_orders SET service_shop_id = $2 WHERE id = $1
	`, workOrder.ID, shop.ID)

	part := testutil.CreatePart(t, db.RawPool(), cfg)
	inventory := testutil.CreateInventoryItem(t, db.RawPool(), cfg, shop.ID, part.ID, 100)

	// Add part to BOM
	now := time.Now().UTC()
	qtyPlanned := int64(10)

	bomItem := models.WorkOrderPart{
		ID:            store.NewID("bom"),
		TenantID:      cfg.TenantID,
		SchoolID:      cfg.SchoolID,
		WorkOrderID:   workOrder.ID,
		ServiceShopID: shop.ID,
		PartID:        part.ID,
		PartName:      part.Name,
		PartCategory:  part.Category,
		QtyPlanned:    qtyPlanned,
		QtyUsed:       0,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	// Reserve inventory
	execSQL(t, db.RawPool(), `
		UPDATE inventory
		SET qty_reserved = qty_reserved + $3, updated_at = $4
		WHERE tenant_id=$1 AND service_shop_id=$2 AND part_id=$5
		  AND (qty_available - qty_reserved) >= $3
	`, cfg.TenantID, shop.ID, qtyPlanned, now, part.ID)

	err := store.CreateWorkOrderPartTx(ctx, db.RawPool(), bomItem)
	require.NoError(t, err, "should add part to BOM")

	// Verify BOM item was created
	retrieved, err := db.WorkOrderParts().GetByID(ctx, cfg.TenantID, cfg.SchoolID, bomItem.ID)
	require.NoError(t, err)
	assert.Equal(t, qtyPlanned, retrieved.QtyPlanned)
	assert.Equal(t, int64(0), retrieved.QtyUsed)

	// Verify inventory was reserved
	updatedInv, err := db.Inventory().GetByPartID(ctx, cfg.TenantID, shop.ID, part.ID)
	require.NoError(t, err)
	assert.Equal(t, inventory.QtyAvailable, updatedInv.QtyAvailable)
	assert.Equal(t, qtyPlanned, updatedInv.QtyReserved)
}

func TestBOMOperations_ConsumePartsFromInventory(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	db, cfg, cleanup := setupTestWithFixtures(t)
	defer cleanup()

	ctx := context.Background()

	// Setup
	deviceID := "test-device-bom-002"
	testutil.CreateDeviceSnapshot(t, db.RawPool(), cfg, deviceID)
	incident := testutil.CreateIncident(t, db.RawPool(), cfg, deviceID)
	shop := testutil.CreateServiceShop(t, db.RawPool(), cfg)
	workOrder := testutil.CreateWorkOrder(t, db.RawPool(), cfg, incident.ID, deviceID)

	execSQL(t, db.RawPool(), `UPDATE work_orders SET service_shop_id = $2 WHERE id = $1`, workOrder.ID, shop.ID)

	part := testutil.CreatePart(t, db.RawPool(), cfg)
	inventory := testutil.CreateInventoryItem(t, db.RawPool(), cfg, shop.ID, part.ID, 100)

	// Add part to BOM and reserve
	now := time.Now().UTC()
	qtyPlanned := int64(20)

	bomItem := models.WorkOrderPart{
		ID:            store.NewID("bom"),
		TenantID:      cfg.TenantID,
		SchoolID:      cfg.SchoolID,
		WorkOrderID:   workOrder.ID,
		ServiceShopID: shop.ID,
		PartID:        part.ID,
		PartName:      part.Name,
		QtyPlanned:    qtyPlanned,
		QtyUsed:       0,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	execSQL(t, db.RawPool(), `
		UPDATE inventory
		SET qty_reserved = qty_reserved + $3, updated_at = $4
		WHERE tenant_id=$1 AND service_shop_id=$2 AND part_id=$5
	`, cfg.TenantID, shop.ID, qtyPlanned, now, part.ID)

	err := store.CreateWorkOrderPartTx(ctx, db.RawPool(), bomItem)
	require.NoError(t, err)

	// Consume parts (use 15 out of 20 planned)
	qtyUsed := int64(15)

	err = store.UpdateWorkOrderPartUsedTx(ctx, db.RawPool(), cfg.TenantID, cfg.SchoolID, bomItem.ID, qtyUsed, now)
	require.NoError(t, err, "should consume parts")

	// Update inventory to reflect consumption
	execSQL(t, db.RawPool(), `
		UPDATE inventory
		SET qty_reserved = GREATEST(qty_reserved - $3, 0),
			qty_available = GREATEST(qty_available - $3, 0),
			updated_at=$4
		WHERE tenant_id=$1 AND service_shop_id=$2 AND part_id=$5
	`, cfg.TenantID, shop.ID, qtyUsed, now, part.ID)

	// Verify BOM item updated
	updatedBOM, err := db.WorkOrderParts().GetByID(ctx, cfg.TenantID, cfg.SchoolID, bomItem.ID)
	require.NoError(t, err)
	assert.Equal(t, qtyUsed, updatedBOM.QtyUsed)

	// Verify inventory consumed
	updatedInv, err := db.Inventory().GetByPartID(ctx, cfg.TenantID, shop.ID, part.ID)
	require.NoError(t, err)
	assert.Equal(t, inventory.QtyAvailable-qtyUsed, updatedInv.QtyAvailable, "available should be reduced")
	assert.Equal(t, qtyPlanned-qtyUsed, updatedInv.QtyReserved, "reserved should be reduced")

	// Try to consume more than planned (should fail or be guarded)
	extraQty := int64(10) // Would exceed planned
	err = store.UpdateWorkOrderPartUsedTx(ctx, db.RawPool(), cfg.TenantID, cfg.SchoolID, bomItem.ID, extraQty, now)
	assert.Error(t, err, "should not allow consuming more than planned")
}

func TestBOMOperations_ReleaseUnusedParts(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	db, cfg, cleanup := setupTestWithFixtures(t)
	defer cleanup()

	ctx := context.Background()

	// Setup
	deviceID := "test-device-bom-003"
	testutil.CreateDeviceSnapshot(t, db.RawPool(), cfg, deviceID)
	incident := testutil.CreateIncident(t, db.RawPool(), cfg, deviceID)
	shop := testutil.CreateServiceShop(t, db.RawPool(), cfg)
	workOrder := testutil.CreateWorkOrder(t, db.RawPool(), cfg, incident.ID, deviceID)

	execSQL(t, db.RawPool(), `UPDATE work_orders SET service_shop_id = $2 WHERE id = $1`, workOrder.ID, shop.ID)

	part := testutil.CreatePart(t, db.RawPool(), cfg)
	inventory := testutil.CreateInventoryItem(t, db.RawPool(), cfg, shop.ID, part.ID, 100)

	// Add part to BOM and reserve
	now := time.Now().UTC()
	qtyPlanned := int64(20)

	bomItem := models.WorkOrderPart{
		ID:            store.NewID("bom"),
		TenantID:      cfg.TenantID,
		SchoolID:      cfg.SchoolID,
		WorkOrderID:   workOrder.ID,
		ServiceShopID: shop.ID,
		PartID:        part.ID,
		PartName:      part.Name,
		QtyPlanned:    qtyPlanned,
		QtyUsed:       0,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	execSQL(t, db.RawPool(), `
		UPDATE inventory
		SET qty_reserved = qty_reserved + $3, updated_at = $4
		WHERE tenant_id=$1 AND service_shop_id=$2 AND part_id=$5
	`, cfg.TenantID, shop.ID, qtyPlanned, now, part.ID)

	err := store.CreateWorkOrderPartTx(ctx, db.RawPool(), bomItem)
	require.NoError(t, err)

	// Consume some parts
	qtyUsed := int64(8)
	err = store.UpdateWorkOrderPartUsedTx(ctx, db.RawPool(), cfg.TenantID, cfg.SchoolID, bomItem.ID, qtyUsed, now)
	require.NoError(t, err)

	execSQL(t, db.RawPool(), `
		UPDATE inventory
		SET qty_reserved = GREATEST(qty_reserved - $3, 0),
			qty_available = GREATEST(qty_available - $3, 0),
			updated_at=$4
		WHERE tenant_id=$1 AND service_shop_id=$2 AND part_id=$5
	`, cfg.TenantID, shop.ID, qtyUsed, now, part.ID)

	// Release unused parts (planned was 20, used 8, so release 12)
	qtyToRelease := qtyPlanned - qtyUsed
	newPlanned := qtyUsed // Reduce planned to match used

	err = store.UpdateWorkOrderPartPlannedTx(ctx, db.RawPool(), cfg.TenantID, cfg.SchoolID, bomItem.ID, newPlanned, now)
	require.NoError(t, err, "should release unused parts")

	// Update inventory
	execSQL(t, db.RawPool(), `
		UPDATE inventory
		SET qty_reserved = GREATEST(qty_reserved - $3, 0), updated_at=$4
		WHERE tenant_id=$1 AND service_shop_id=$2 AND part_id=$5
	`, cfg.TenantID, shop.ID, qtyToRelease, now, part.ID)

	// Verify BOM updated
	updatedBOM, err := db.WorkOrderParts().GetByID(ctx, cfg.TenantID, cfg.SchoolID, bomItem.ID)
	require.NoError(t, err)
	assert.Equal(t, newPlanned, updatedBOM.QtyPlanned)
	assert.Equal(t, qtyUsed, updatedBOM.QtyUsed)

	// Verify inventory released
	updatedInv, err := db.Inventory().GetByPartID(ctx, cfg.TenantID, shop.ID, part.ID)
	require.NoError(t, err)
	expectedAvailable := inventory.QtyAvailable - qtyUsed
	expectedReserved := int64(0) // All reserved should be released
	assert.Equal(t, expectedAvailable, updatedInv.QtyAvailable)
	assert.Equal(t, expectedReserved, updatedInv.QtyReserved)
}

func TestBOMOperations_InsufficientInventory(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	db, cfg, cleanup := setupTestWithFixtures(t)
	defer cleanup()

	ctx := context.Background()

	// Setup
	deviceID := "test-device-bom-004"
	testutil.CreateDeviceSnapshot(t, db.RawPool(), cfg, deviceID)
	incident := testutil.CreateIncident(t, db.RawPool(), cfg, deviceID)
	shop := testutil.CreateServiceShop(t, db.RawPool(), cfg)
	workOrder := testutil.CreateWorkOrder(t, db.RawPool(), cfg, incident.ID, deviceID)

	execSQL(t, db.RawPool(), `UPDATE work_orders SET service_shop_id = $2 WHERE id = $1`, workOrder.ID, shop.ID)

	part := testutil.CreatePart(t, db.RawPool(), cfg)
	// Create inventory with only 5 items available
	testutil.CreateInventoryItem(t, db.RawPool(), cfg, shop.ID, part.ID, 5)

	// Try to reserve 10 items (more than available)
	now := time.Now().UTC()
	qtyPlanned := int64(10)

	// Try to reserve - should fail with constraint
	result, err := db.RawPool().Exec(ctx, `
		UPDATE inventory
		SET qty_reserved = qty_reserved + $3, updated_at = $4
		WHERE tenant_id=$1 AND service_shop_id=$2 AND part_id=$5
		  AND (qty_available - qty_reserved) >= $3
	`, cfg.TenantID, shop.ID, qtyPlanned, now, part.ID)

	require.NoError(t, err)
	rowsAffected := result.RowsAffected()

	// Should not update any rows (constraint not met)
	assert.Equal(t, int64(0), rowsAffected, "should not reserve when insufficient inventory")

	// Verify inventory unchanged
	inv, err := db.Inventory().GetByPartID(ctx, cfg.TenantID, shop.ID, part.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(5), inv.QtyAvailable)
	assert.Equal(t, int64(0), inv.QtyReserved)
}

func TestBOMOperations_MultipleParts(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	db, cfg, cleanup := setupTestWithFixtures(t)
	defer cleanup()

	ctx := context.Background()

	// Setup
	deviceID := "test-device-bom-005"
	testutil.CreateDeviceSnapshot(t, db.RawPool(), cfg, deviceID)
	incident := testutil.CreateIncident(t, db.RawPool(), cfg, deviceID)
	shop := testutil.CreateServiceShop(t, db.RawPool(), cfg)
	workOrder := testutil.CreateWorkOrder(t, db.RawPool(), cfg, incident.ID, deviceID)

	execSQL(t, db.RawPool(), `UPDATE work_orders SET service_shop_id = $2 WHERE id = $1`, workOrder.ID, shop.ID)

	// Create multiple parts and inventory
	numParts := 3
	parts := make([]models.Part, numParts)
	bomItems := make([]models.WorkOrderPart, numParts)

	now := time.Now().UTC()

	for i := 0; i < numParts; i++ {
		parts[i] = testutil.CreatePart(t, db.RawPool(), cfg)
		testutil.CreateInventoryItem(t, db.RawPool(), cfg, shop.ID, parts[i].ID, 100)

		qtyPlanned := int64((i + 1) * 5) // 5, 10, 15

		bomItems[i] = models.WorkOrderPart{
			ID:            store.NewID("bom"),
			TenantID:      cfg.TenantID,
			SchoolID:      cfg.SchoolID,
			WorkOrderID:   workOrder.ID,
			ServiceShopID: shop.ID,
			PartID:        parts[i].ID,
			PartName:      parts[i].Name,
			QtyPlanned:    qtyPlanned,
			QtyUsed:       0,
			CreatedAt:     now,
			UpdatedAt:     now,
		}

		// Reserve inventory
		execSQL(t, db.RawPool(), `
			UPDATE inventory
			SET qty_reserved = qty_reserved + $3, updated_at = $4
			WHERE tenant_id=$1 AND service_shop_id=$2 AND part_id=$5
		`, cfg.TenantID, shop.ID, qtyPlanned, now, parts[i].ID)

		err := store.CreateWorkOrderPartTx(ctx, db.RawPool(), bomItems[i])
		require.NoError(t, err)
	}

	// List all BOM items for the work order
	bomList, _, err := db.WorkOrderParts().List(ctx, store.WorkOrderPartListParams{
		TenantID:    cfg.TenantID,
		SchoolID:    cfg.SchoolID,
		WorkOrderID: workOrder.ID,
		Limit:       10,
	})
	require.NoError(t, err)
	assert.Len(t, bomList, numParts)

	// Verify quantities
	totalPlanned := int64(0)
	for _, item := range bomList {
		totalPlanned += item.QtyPlanned
		assert.Equal(t, int64(0), item.QtyUsed)
	}
	assert.Equal(t, int64(30), totalPlanned) // 5 + 10 + 15

	// Consume different amounts from each part
	for i, item := range bomItems {
		qtyUsed := int64((i + 1) * 2) // 2, 4, 6

		err := store.UpdateWorkOrderPartUsedTx(ctx, db.RawPool(), cfg.TenantID, cfg.SchoolID, item.ID, qtyUsed, now)
		require.NoError(t, err)

		execSQL(t, db.RawPool(), `
			UPDATE inventory
			SET qty_reserved = GREATEST(qty_reserved - $3, 0),
				qty_available = GREATEST(qty_available - $3, 0),
				updated_at=$4
			WHERE tenant_id=$1 AND service_shop_id=$2 AND part_id=$5
		`, cfg.TenantID, shop.ID, qtyUsed, now, parts[i].ID)
	}

	// Verify consumption
	totalUsed := int64(0)
	for _, item := range bomItems {
		updated, err := db.WorkOrderParts().GetByID(ctx, cfg.TenantID, cfg.SchoolID, item.ID)
		require.NoError(t, err)
		totalUsed += updated.QtyUsed
	}
	assert.Equal(t, int64(12), totalUsed) // 2 + 4 + 6
}

func TestBOMOperations_ConcurrentReservation(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	db, cfg, cleanup := setupTestWithFixtures(t)
	defer cleanup()

	ctx := context.Background()

	// Setup
	deviceID := "test-device-bom-006"
	testutil.CreateDeviceSnapshot(t, db.RawPool(), cfg, deviceID)
	shop := testutil.CreateServiceShop(t, db.RawPool(), cfg)
	part := testutil.CreatePart(t, db.RawPool(), cfg)

	// Create inventory with exactly 20 items
	testutil.CreateInventoryItem(t, db.RawPool(), cfg, shop.ID, part.ID, 20)

	// Create two work orders trying to reserve 15 items each concurrently
	// Only one should fully succeed
	incident1 := testutil.CreateIncident(t, db.RawPool(), cfg, deviceID)
	incident2 := testutil.CreateIncident(t, db.RawPool(), cfg, deviceID)

	wo1 := testutil.CreateWorkOrder(t, db.RawPool(), cfg, incident1.ID, deviceID)
	wo2 := testutil.CreateWorkOrder(t, db.RawPool(), cfg, incident2.ID, deviceID)

	execSQL(t, db.RawPool(), `UPDATE work_orders SET service_shop_id = $2 WHERE id = $1`, wo1.ID, shop.ID)
	execSQL(t, db.RawPool(), `UPDATE work_orders SET service_shop_id = $2 WHERE id = $1`, wo2.ID, shop.ID)

	now := time.Now().UTC()
	qtyToReserve := int64(15)

	done := make(chan int64, 2)

	// Try concurrent reservations
	reserve := func(woID string) {
		result, err := db.RawPool().Exec(ctx, `
			UPDATE inventory
			SET qty_reserved = qty_reserved + $3, updated_at = $4
			WHERE tenant_id=$1 AND service_shop_id=$2 AND part_id=$5
			  AND (qty_available - qty_reserved) >= $3
		`, cfg.TenantID, shop.ID, qtyToReserve, now, part.ID)

		if err != nil {
			done <- 0
			return
		}

		done <- result.RowsAffected()
	}

	go reserve(wo1.ID)
	go reserve(wo2.ID)

	// Collect results
	result1 := <-done
	result2 := <-done

	// Exactly one should succeed (rows affected = 1), one should fail (rows affected = 0)
	totalSuccesses := result1 + result2
	assert.Equal(t, int64(1), totalSuccesses, "exactly one reservation should succeed")

	// Verify final inventory state
	inv, err := db.Inventory().GetByPartID(ctx, cfg.TenantID, shop.ID, part.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(20), inv.QtyAvailable)
	assert.Equal(t, int64(15), inv.QtyReserved, "should have 15 reserved from the one successful reservation")
}

func TestBOMOperations_InventoryTracking(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	db, cfg, cleanup := setupTestWithFixtures(t)
	defer cleanup()

	ctx := context.Background()

	// Setup
	deviceID := "test-device-bom-007"
	testutil.CreateDeviceSnapshot(t, db.RawPool(), cfg, deviceID)
	incident := testutil.CreateIncident(t, db.RawPool(), cfg, deviceID)
	shop := testutil.CreateServiceShop(t, db.RawPool(), cfg)
	workOrder := testutil.CreateWorkOrder(t, db.RawPool(), cfg, incident.ID, deviceID)

	execSQL(t, db.RawPool(), `UPDATE work_orders SET service_shop_id = $2 WHERE id = $1`, workOrder.ID, shop.ID)

	part := testutil.CreatePart(t, db.RawPool(), cfg)
	initialQty := int64(100)
	testutil.CreateInventoryItem(t, db.RawPool(), cfg, shop.ID, part.ID, initialQty)

	now := time.Now().UTC()

	// Step 1: Reserve 30 items
	qtyPlanned := int64(30)
	execSQL(t, db.RawPool(), `
		UPDATE inventory
		SET qty_reserved = qty_reserved + $3, updated_at = $4
		WHERE tenant_id=$1 AND service_shop_id=$2 AND part_id=$5
	`, cfg.TenantID, shop.ID, qtyPlanned, now, part.ID)

	inv1, err := db.Inventory().GetByPartID(ctx, cfg.TenantID, shop.ID, part.ID)
	require.NoError(t, err)
	assert.Equal(t, initialQty, inv1.QtyAvailable)
	assert.Equal(t, qtyPlanned, inv1.QtyReserved)

	// Step 2: Consume 20 items
	qtyUsed := int64(20)
	execSQL(t, db.RawPool(), `
		UPDATE inventory
		SET qty_reserved = GREATEST(qty_reserved - $3, 0),
			qty_available = GREATEST(qty_available - $3, 0),
			updated_at=$4
		WHERE tenant_id=$1 AND service_shop_id=$2 AND part_id=$5
	`, cfg.TenantID, shop.ID, qtyUsed, now, part.ID)

	inv2, err := db.Inventory().GetByPartID(ctx, cfg.TenantID, shop.ID, part.ID)
	require.NoError(t, err)
	assert.Equal(t, initialQty-qtyUsed, inv2.QtyAvailable, "available should decrease by used amount")
	assert.Equal(t, qtyPlanned-qtyUsed, inv2.QtyReserved, "reserved should decrease by used amount")

	// Step 3: Release remaining 10 reserved items
	qtyToRelease := qtyPlanned - qtyUsed
	execSQL(t, db.RawPool(), `
		UPDATE inventory
		SET qty_reserved = GREATEST(qty_reserved - $3, 0), updated_at=$4
		WHERE tenant_id=$1 AND service_shop_id=$2 AND part_id=$5
	`, cfg.TenantID, shop.ID, qtyToRelease, now, part.ID)

	inv3, err := db.Inventory().GetByPartID(ctx, cfg.TenantID, shop.ID, part.ID)
	require.NoError(t, err)
	assert.Equal(t, initialQty-qtyUsed, inv3.QtyAvailable, "available should stay the same after release")
	assert.Equal(t, int64(0), inv3.QtyReserved, "all reservations should be released")

	// Final state: 80 available, 0 reserved (20 were consumed)
	assert.Equal(t, int64(80), inv3.QtyAvailable)
	assert.Equal(t, int64(0), inv3.QtyReserved)
}
