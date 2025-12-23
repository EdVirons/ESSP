package store

import (
	"context"
	"testing"
	"time"

	"github.com/edvirons/ssp/ims/internal/models"
)

func TestWorkOrderRepo_Create(t *testing.T) {
	pool := setupTestDB(t)
	repo := &WorkOrderRepo{pool: pool}

	tests := []struct {
		name    string
		input   models.WorkOrder
		wantErr bool
	}{
		{
			name:    "valid work order",
			input:   validWorkOrder(),
			wantErr: false,
		},
		{
			name: "work order with different ID",
			input: func() models.WorkOrder {
				wo := validWorkOrder()
				wo.ID = "wo-test-create-002"
				return wo
			}(),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Cleanup before and after test
			defer cleanupWorkOrders(t, pool, tt.input.TenantID, tt.input.SchoolID)
			cleanupWorkOrders(t, pool, tt.input.TenantID, tt.input.SchoolID)

			err := repo.Create(context.Background(), tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Verify creation if no error expected
			if !tt.wantErr && err == nil {
				got, err := repo.GetByID(context.Background(), tt.input.TenantID, tt.input.SchoolID, tt.input.ID)
				if err != nil {
					t.Errorf("Failed to retrieve created work order: %v", err)
					return
				}
				if got.ID != tt.input.ID {
					t.Errorf("Created work order ID = %v, want %v", got.ID, tt.input.ID)
				}
				if got.TaskType != tt.input.TaskType {
					t.Errorf("Created work order TaskType = %v, want %v", got.TaskType, tt.input.TaskType)
				}
			}
		})
	}
}

func TestWorkOrderRepo_GetByID(t *testing.T) {
	pool := setupTestDB(t)
	repo := &WorkOrderRepo{pool: pool}

	// Setup: Create a test work order
	workOrder := validWorkOrder()
	workOrder.ID = "wo-test-getbyid-001"
	defer cleanupWorkOrders(t, pool, workOrder.TenantID, workOrder.SchoolID)
	cleanupWorkOrders(t, pool, workOrder.TenantID, workOrder.SchoolID)

	err := repo.Create(context.Background(), workOrder)
	if err != nil {
		t.Fatalf("Failed to create test work order: %v", err)
	}

	tests := []struct {
		name     string
		tenantID string
		schoolID string
		id       string
		wantErr  bool
		wantID   string
	}{
		{
			name:     "existing work order",
			tenantID: workOrder.TenantID,
			schoolID: workOrder.SchoolID,
			id:       workOrder.ID,
			wantErr:  false,
			wantID:   workOrder.ID,
		},
		{
			name:     "non-existent work order",
			tenantID: workOrder.TenantID,
			schoolID: workOrder.SchoolID,
			id:       "non-existent-id",
			wantErr:  true,
		},
		{
			name:     "wrong tenant",
			tenantID: "wrong-tenant",
			schoolID: workOrder.SchoolID,
			id:       workOrder.ID,
			wantErr:  true,
		},
		{
			name:     "wrong school",
			tenantID: workOrder.TenantID,
			schoolID: "wrong-school",
			id:       workOrder.ID,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.GetByID(context.Background(), tt.tenantID, tt.schoolID, tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.ID != tt.wantID {
				t.Errorf("GetByID() ID = %v, want %v", got.ID, tt.wantID)
			}
		})
	}
}

func TestWorkOrderRepo_List(t *testing.T) {
	pool := setupTestDB(t)
	repo := &WorkOrderRepo{pool: pool}

	tenantID := "tenant-test-wolist"
	schoolID := "school-test-wolist"

	// Setup: Create multiple test work orders
	defer cleanupWorkOrders(t, pool, tenantID, schoolID)
	cleanupWorkOrders(t, pool, tenantID, schoolID)

	now := time.Now().UTC()
	workOrders := []models.WorkOrder{
		func() models.WorkOrder {
			wo := validWorkOrder()
			wo.ID = "wo-list-001"
			wo.TenantID = tenantID
			wo.SchoolID = schoolID
			wo.Status = models.WorkOrderDraft
			wo.DeviceID = "device-001"
			wo.IncidentID = "inc-001"
			wo.CreatedAt = now.Add(-3 * time.Hour)
			wo.UpdatedAt = now.Add(-3 * time.Hour)
			return wo
		}(),
		func() models.WorkOrder {
			wo := validWorkOrder()
			wo.ID = "wo-list-002"
			wo.TenantID = tenantID
			wo.SchoolID = schoolID
			wo.Status = models.WorkOrderAssigned
			wo.DeviceID = "device-002"
			wo.IncidentID = "inc-002"
			wo.CreatedAt = now.Add(-2 * time.Hour)
			wo.UpdatedAt = now.Add(-2 * time.Hour)
			return wo
		}(),
		func() models.WorkOrder {
			wo := validWorkOrder()
			wo.ID = "wo-list-003"
			wo.TenantID = tenantID
			wo.SchoolID = schoolID
			wo.Status = models.WorkOrderDraft
			wo.DeviceID = "device-001"
			wo.IncidentID = "inc-001"
			wo.CreatedAt = now.Add(-1 * time.Hour)
			wo.UpdatedAt = now.Add(-1 * time.Hour)
			return wo
		}(),
	}

	for _, wo := range workOrders {
		if err := repo.Create(context.Background(), wo); err != nil {
			t.Fatalf("Failed to create test work order %s: %v", wo.ID, err)
		}
	}

	tests := []struct {
		name      string
		params    WorkOrderListParams
		wantCount int
		wantNext  bool
	}{
		{
			name: "list all work orders",
			params: WorkOrderListParams{
				TenantID: tenantID,
				SchoolID: schoolID,
				Limit:    10,
			},
			wantCount: 3,
			wantNext:  false,
		},
		{
			name: "filter by status",
			params: WorkOrderListParams{
				TenantID: tenantID,
				SchoolID: schoolID,
				Status:   string(models.WorkOrderDraft),
				Limit:    10,
			},
			wantCount: 2,
			wantNext:  false,
		},
		{
			name: "filter by device",
			params: WorkOrderListParams{
				TenantID: tenantID,
				SchoolID: schoolID,
				DeviceID: "device-001",
				Limit:    10,
			},
			wantCount: 2,
			wantNext:  false,
		},
		{
			name: "filter by incident",
			params: WorkOrderListParams{
				TenantID:   tenantID,
				SchoolID:   schoolID,
				IncidentID: "inc-001",
				Limit:      10,
			},
			wantCount: 2,
			wantNext:  false,
		},
		{
			name: "pagination with limit",
			params: WorkOrderListParams{
				TenantID: tenantID,
				SchoolID: schoolID,
				Limit:    2,
			},
			wantCount: 2,
			wantNext:  true,
		},
		{
			name: "pagination with cursor",
			params: func() WorkOrderListParams {
				// Get first page to obtain cursor
				firstPage, nextCursor, err := repo.List(context.Background(), WorkOrderListParams{
					TenantID: tenantID,
					SchoolID: schoolID,
					Limit:    1,
				})
				if err != nil || len(firstPage) == 0 {
					t.Fatalf("Failed to get first page for cursor test: %v", err)
				}

				cursorTime, cursorID, ok := DecodeCursor(nextCursor)
				if !ok {
					t.Fatalf("Failed to decode cursor: %s", nextCursor)
				}

				return WorkOrderListParams{
					TenantID:        tenantID,
					SchoolID:        schoolID,
					Limit:           10,
					HasCursor:       true,
					CursorCreatedAt: cursorTime,
					CursorID:        cursorID,
				}
			}(),
			wantCount: 2,
			wantNext:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, nextCursor, err := repo.List(context.Background(), tt.params)
			if err != nil {
				t.Errorf("List() error = %v", err)
				return
			}
			if len(got) != tt.wantCount {
				t.Errorf("List() count = %v, want %v", len(got), tt.wantCount)
			}
			if tt.wantNext && nextCursor == "" {
				t.Errorf("List() expected next cursor, got empty")
			}
			if !tt.wantNext && nextCursor != "" {
				t.Errorf("List() expected no next cursor, got %s", nextCursor)
			}
		})
	}
}

func TestWorkOrderRepo_UpdateStatus(t *testing.T) {
	pool := setupTestDB(t)
	repo := &WorkOrderRepo{pool: pool}

	// Setup: Create a test work order
	workOrder := validWorkOrder()
	workOrder.ID = "wo-test-updatestatus-001"
	workOrder.Status = models.WorkOrderDraft
	defer cleanupWorkOrders(t, pool, workOrder.TenantID, workOrder.SchoolID)
	cleanupWorkOrders(t, pool, workOrder.TenantID, workOrder.SchoolID)

	err := repo.Create(context.Background(), workOrder)
	if err != nil {
		t.Fatalf("Failed to create test work order: %v", err)
	}

	tests := []struct {
		name       string
		tenantID   string
		schoolID   string
		id         string
		newStatus  models.WorkOrderStatus
		wantErr    bool
		wantStatus models.WorkOrderStatus
	}{
		{
			name:       "update to assigned",
			tenantID:   workOrder.TenantID,
			schoolID:   workOrder.SchoolID,
			id:         workOrder.ID,
			newStatus:  models.WorkOrderAssigned,
			wantErr:    false,
			wantStatus: models.WorkOrderAssigned,
		},
		{
			name:      "update non-existent work order",
			tenantID:  workOrder.TenantID,
			schoolID:  workOrder.SchoolID,
			id:        "non-existent-id",
			newStatus: models.WorkOrderCompleted,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			now := time.Now().UTC()
			got, err := repo.UpdateStatus(context.Background(), tt.tenantID, tt.schoolID, tt.id, tt.newStatus, now)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.Status != tt.wantStatus {
				t.Errorf("UpdateStatus() status = %v, want %v", got.Status, tt.wantStatus)
			}
			// Verify UpdatedAt was updated
			if !tt.wantErr && !got.UpdatedAt.After(workOrder.CreatedAt) {
				t.Errorf("UpdateStatus() UpdatedAt was not updated")
			}
		})
	}
}

func TestWorkOrderRepo_SetApprovalStatus(t *testing.T) {
	pool := setupTestDB(t)
	repo := &WorkOrderRepo{pool: pool}

	// Setup: Create a test work order
	workOrder := validWorkOrder()
	workOrder.ID = "wo-test-approval-001"
	workOrder.ApprovalStatus = ""
	defer cleanupWorkOrders(t, pool, workOrder.TenantID, workOrder.SchoolID)
	cleanupWorkOrders(t, pool, workOrder.TenantID, workOrder.SchoolID)

	err := repo.Create(context.Background(), workOrder)
	if err != nil {
		t.Fatalf("Failed to create test work order: %v", err)
	}

	tests := []struct {
		name         string
		tenantID     string
		schoolID     string
		workOrderID  string
		status       string
		wantErr      bool
		wantApproval string
	}{
		{
			name:         "set approval to approved",
			tenantID:     workOrder.TenantID,
			schoolID:     workOrder.SchoolID,
			workOrderID:  workOrder.ID,
			status:       "approved",
			wantErr:      false,
			wantApproval: "approved",
		},
		{
			name:        "set approval on non-existent work order",
			tenantID:    workOrder.TenantID,
			schoolID:    workOrder.SchoolID,
			workOrderID: "non-existent-id",
			status:      "approved",
			wantErr:     false, // UPDATE doesn't fail on no rows affected
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.SetApprovalStatus(context.Background(), tt.tenantID, tt.schoolID, tt.workOrderID, tt.status)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetApprovalStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Verify the approval status was set (only if work order exists)
			if !tt.wantErr && tt.workOrderID == workOrder.ID {
				got, err := repo.GetByID(context.Background(), tt.tenantID, tt.schoolID, tt.workOrderID)
				if err != nil {
					t.Errorf("Failed to retrieve work order after setting approval: %v", err)
					return
				}
				if got.ApprovalStatus != tt.wantApproval {
					t.Errorf("SetApprovalStatus() approval = %v, want %v", got.ApprovalStatus, tt.wantApproval)
				}
			}
		})
	}
}

func TestWorkOrderRepo_ListByPhase(t *testing.T) {
	pool := setupTestDB(t)
	repo := &WorkOrderRepo{pool: pool}

	tenantID := "tenant-test-phase"
	schoolID := "school-test-phase"

	// Setup: Create work orders with different phases
	defer cleanupWorkOrders(t, pool, tenantID, schoolID)
	cleanupWorkOrders(t, pool, tenantID, schoolID)

	now := time.Now().UTC()
	workOrders := []models.WorkOrder{
		func() models.WorkOrder {
			wo := validWorkOrder()
			wo.ID = "wo-phase-001"
			wo.TenantID = tenantID
			wo.SchoolID = schoolID
			wo.ProjectID = "project-001"
			wo.PhaseID = "phase-001"
			wo.CreatedAt = now.Add(-3 * time.Hour)
			wo.UpdatedAt = now.Add(-3 * time.Hour)
			return wo
		}(),
		func() models.WorkOrder {
			wo := validWorkOrder()
			wo.ID = "wo-phase-002"
			wo.TenantID = tenantID
			wo.SchoolID = schoolID
			wo.ProjectID = "project-001"
			wo.PhaseID = "phase-001"
			wo.CreatedAt = now.Add(-2 * time.Hour)
			wo.UpdatedAt = now.Add(-2 * time.Hour)
			return wo
		}(),
		func() models.WorkOrder {
			wo := validWorkOrder()
			wo.ID = "wo-phase-003"
			wo.TenantID = tenantID
			wo.SchoolID = schoolID
			wo.ProjectID = "project-001"
			wo.PhaseID = "phase-002" // Different phase
			wo.CreatedAt = now.Add(-1 * time.Hour)
			wo.UpdatedAt = now.Add(-1 * time.Hour)
			return wo
		}(),
	}

	for _, wo := range workOrders {
		if err := repo.Create(context.Background(), wo); err != nil {
			t.Fatalf("Failed to create test work order %s: %v", wo.ID, err)
		}
	}

	tests := []struct {
		name      string
		tenantID  string
		phaseID   string
		wantCount int
		wantErr   bool
	}{
		{
			name:      "list work orders for phase-001",
			tenantID:  tenantID,
			phaseID:   "phase-001",
			wantCount: 2,
			wantErr:   false,
		},
		{
			name:      "list work orders for phase-002",
			tenantID:  tenantID,
			phaseID:   "phase-002",
			wantCount: 1,
			wantErr:   false,
		},
		{
			name:      "list work orders for non-existent phase",
			tenantID:  tenantID,
			phaseID:   "phase-999",
			wantCount: 0,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.ListByPhase(context.Background(), tt.tenantID, tt.phaseID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListByPhase() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.wantCount {
				t.Errorf("ListByPhase() count = %v, want %v", len(got), tt.wantCount)
			}

			// Verify all returned work orders belong to the requested phase
			for _, wo := range got {
				if wo.PhaseID != tt.phaseID {
					t.Errorf("ListByPhase() returned work order with phase %v, want %v", wo.PhaseID, tt.phaseID)
				}
			}
		})
	}
}
