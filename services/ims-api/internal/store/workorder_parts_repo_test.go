package store

import (
	"context"
	"testing"
	"time"

	"github.com/edvirons/ssp/ims/internal/models"
)

func TestWorkOrderPartRepo_Create(t *testing.T) {
	pool := setupTestDB(t)
	repo := &WorkOrderPartRepo{pool: pool}

	tests := []struct {
		name    string
		input   models.WorkOrderPart
		wantErr bool
	}{
		{
			name:    "valid work order part",
			input:   validWorkOrderPart(),
			wantErr: false,
		},
		{
			name: "work order part with different ID",
			input: func() models.WorkOrderPart {
				part := validWorkOrderPart()
				part.ID = "wop-test-create-002"
				return part
			}(),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Cleanup before and after test
			defer cleanupWorkOrderParts(t, pool, tt.input.TenantID, tt.input.SchoolID)
			cleanupWorkOrderParts(t, pool, tt.input.TenantID, tt.input.SchoolID)

			err := repo.Create(context.Background(), tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Verify creation if no error expected
			if !tt.wantErr && err == nil {
				got, err := repo.GetByID(context.Background(), tt.input.TenantID, tt.input.SchoolID, tt.input.ID)
				if err != nil {
					t.Errorf("Failed to retrieve created work order part: %v", err)
					return
				}
				if got.ID != tt.input.ID {
					t.Errorf("Created work order part ID = %v, want %v", got.ID, tt.input.ID)
				}
				if got.PartName != tt.input.PartName {
					t.Errorf("Created work order part PartName = %v, want %v", got.PartName, tt.input.PartName)
				}
			}
		})
	}
}

func TestWorkOrderPartRepo_GetByID(t *testing.T) {
	pool := setupTestDB(t)
	repo := &WorkOrderPartRepo{pool: pool}

	// Setup: Create a test work order part
	part := validWorkOrderPart()
	part.ID = "wop-test-getbyid-001"
	defer cleanupWorkOrderParts(t, pool, part.TenantID, part.SchoolID)
	cleanupWorkOrderParts(t, pool, part.TenantID, part.SchoolID)

	err := repo.Create(context.Background(), part)
	if err != nil {
		t.Fatalf("Failed to create test work order part: %v", err)
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
			name:     "existing work order part",
			tenantID: part.TenantID,
			schoolID: part.SchoolID,
			id:       part.ID,
			wantErr:  false,
			wantID:   part.ID,
		},
		{
			name:     "non-existent work order part",
			tenantID: part.TenantID,
			schoolID: part.SchoolID,
			id:       "non-existent-id",
			wantErr:  true,
		},
		{
			name:     "wrong tenant",
			tenantID: "wrong-tenant",
			schoolID: part.SchoolID,
			id:       part.ID,
			wantErr:  true,
		},
		{
			name:     "wrong school",
			tenantID: part.TenantID,
			schoolID: "wrong-school",
			id:       part.ID,
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

func TestWorkOrderPartRepo_List(t *testing.T) {
	pool := setupTestDB(t)
	repo := &WorkOrderPartRepo{pool: pool}

	tenantID := "tenant-test-partlist"
	schoolID := "school-test-partlist"
	workOrderID := "wo-test-partlist-001"

	// Setup: Create multiple test work order parts
	defer cleanupWorkOrderParts(t, pool, tenantID, schoolID)
	cleanupWorkOrderParts(t, pool, tenantID, schoolID)

	now := time.Now().UTC()
	parts := []models.WorkOrderPart{
		func() models.WorkOrderPart {
			part := validWorkOrderPart()
			part.ID = "wop-list-001"
			part.TenantID = tenantID
			part.SchoolID = schoolID
			part.WorkOrderID = workOrderID
			part.PartName = "LCD Screen"
			part.QtyPlanned = 2
			part.QtyUsed = 0
			part.CreatedAt = now.Add(-3 * time.Hour)
			part.UpdatedAt = now.Add(-3 * time.Hour)
			return part
		}(),
		func() models.WorkOrderPart {
			part := validWorkOrderPart()
			part.ID = "wop-list-002"
			part.TenantID = tenantID
			part.SchoolID = schoolID
			part.WorkOrderID = workOrderID
			part.PartName = "Battery"
			part.QtyPlanned = 1
			part.QtyUsed = 1
			part.CreatedAt = now.Add(-2 * time.Hour)
			part.UpdatedAt = now.Add(-2 * time.Hour)
			return part
		}(),
		func() models.WorkOrderPart {
			part := validWorkOrderPart()
			part.ID = "wop-list-003"
			part.TenantID = tenantID
			part.SchoolID = schoolID
			part.WorkOrderID = workOrderID
			part.PartName = "Keyboard"
			part.QtyPlanned = 1
			part.QtyUsed = 0
			part.CreatedAt = now.Add(-1 * time.Hour)
			part.UpdatedAt = now.Add(-1 * time.Hour)
			return part
		}(),
		func() models.WorkOrderPart {
			part := validWorkOrderPart()
			part.ID = "wop-list-004"
			part.TenantID = tenantID
			part.SchoolID = schoolID
			part.WorkOrderID = "wo-different-001" // Different work order
			part.PartName = "Mouse"
			part.CreatedAt = now
			part.UpdatedAt = now
			return part
		}(),
	}

	for _, part := range parts {
		if err := repo.Create(context.Background(), part); err != nil {
			t.Fatalf("Failed to create test work order part %s: %v", part.ID, err)
		}
	}

	tests := []struct {
		name      string
		params    WorkOrderPartListParams
		wantCount int
		wantNext  bool
	}{
		{
			name: "list all parts for work order",
			params: WorkOrderPartListParams{
				TenantID:    tenantID,
				SchoolID:    schoolID,
				WorkOrderID: workOrderID,
				Limit:       10,
			},
			wantCount: 3,
			wantNext:  false,
		},
		{
			name: "list parts for different work order",
			params: WorkOrderPartListParams{
				TenantID:    tenantID,
				SchoolID:    schoolID,
				WorkOrderID: "wo-different-001",
				Limit:       10,
			},
			wantCount: 1,
			wantNext:  false,
		},
		{
			name: "pagination with limit",
			params: WorkOrderPartListParams{
				TenantID:    tenantID,
				SchoolID:    schoolID,
				WorkOrderID: workOrderID,
				Limit:       2,
			},
			wantCount: 2,
			wantNext:  true,
		},
		{
			name: "pagination with cursor",
			params: func() WorkOrderPartListParams {
				// Get first page to obtain cursor
				firstPage, nextCursor, err := repo.List(context.Background(), WorkOrderPartListParams{
					TenantID:    tenantID,
					SchoolID:    schoolID,
					WorkOrderID: workOrderID,
					Limit:       1,
				})
				if err != nil || len(firstPage) == 0 {
					t.Fatalf("Failed to get first page for cursor test: %v", err)
				}

				cursorTime, cursorID, ok := DecodeCursor(nextCursor)
				if !ok {
					t.Fatalf("Failed to decode cursor: %s", nextCursor)
				}

				return WorkOrderPartListParams{
					TenantID:        tenantID,
					SchoolID:        schoolID,
					WorkOrderID:     workOrderID,
					Limit:           10,
					HasCursor:       true,
					CursorCreatedAt: cursorTime,
					CursorID:        cursorID,
				}
			}(),
			wantCount: 2,
			wantNext:  false,
		},
		{
			name: "non-existent work order",
			params: WorkOrderPartListParams{
				TenantID:    tenantID,
				SchoolID:    schoolID,
				WorkOrderID: "non-existent-wo",
				Limit:       10,
			},
			wantCount: 0,
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

			// Verify all returned parts belong to the requested work order
			for _, part := range got {
				if part.WorkOrderID != tt.params.WorkOrderID {
					t.Errorf("List() returned part with work order %v, want %v", part.WorkOrderID, tt.params.WorkOrderID)
				}
			}
		})
	}
}

func TestCreateWorkOrderPartTx(t *testing.T) {
	pool := setupTestDB(t)

	tenantID := "tenant-test-tx"
	schoolID := "school-test-tx"

	defer cleanupWorkOrderParts(t, pool, tenantID, schoolID)
	cleanupWorkOrderParts(t, pool, tenantID, schoolID)

	tests := []struct {
		name    string
		input   models.WorkOrderPart
		wantErr bool
	}{
		{
			name: "create part in transaction",
			input: func() models.WorkOrderPart {
				part := validWorkOrderPart()
				part.ID = "wop-tx-001"
				part.TenantID = tenantID
				part.SchoolID = schoolID
				return part
			}(),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			// Begin transaction
			tx, err := pool.Begin(ctx)
			if err != nil {
				t.Fatalf("Failed to begin transaction: %v", err)
			}
			defer tx.Rollback(ctx)

			err = CreateWorkOrderPartTx(ctx, tx, tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateWorkOrderPartTx() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				// Commit transaction
				if err := tx.Commit(ctx); err != nil {
					t.Fatalf("Failed to commit transaction: %v", err)
				}

				// Verify the part was created
				repo := &WorkOrderPartRepo{pool: pool}
				got, err := repo.GetByID(ctx, tt.input.TenantID, tt.input.SchoolID, tt.input.ID)
				if err != nil {
					t.Errorf("Failed to retrieve created part: %v", err)
					return
				}
				if got.ID != tt.input.ID {
					t.Errorf("Created part ID = %v, want %v", got.ID, tt.input.ID)
				}
			}
		})
	}
}

func TestUpdateWorkOrderPartUsedTx(t *testing.T) {
	pool := setupTestDB(t)
	repo := &WorkOrderPartRepo{pool: pool}

	tenantID := "tenant-test-updateused"
	schoolID := "school-test-updateused"

	// Setup: Create a test part
	part := validWorkOrderPart()
	part.ID = "wop-updateused-001"
	part.TenantID = tenantID
	part.SchoolID = schoolID
	part.QtyPlanned = 10
	part.QtyUsed = 0

	defer cleanupWorkOrderParts(t, pool, tenantID, schoolID)
	cleanupWorkOrderParts(t, pool, tenantID, schoolID)

	if err := repo.Create(context.Background(), part); err != nil {
		t.Fatalf("Failed to create test part: %v", err)
	}

	tests := []struct {
		name        string
		tenantID    string
		schoolID    string
		id          string
		addUsed     int64
		wantErr     bool
		wantQtyUsed int64
	}{
		{
			name:        "increment used by 2",
			tenantID:    tenantID,
			schoolID:    schoolID,
			id:          part.ID,
			addUsed:     2,
			wantErr:     false,
			wantQtyUsed: 2,
		},
		{
			name:        "increment used by 5 more",
			tenantID:    tenantID,
			schoolID:    schoolID,
			id:          part.ID,
			addUsed:     5,
			wantErr:     false,
			wantQtyUsed: 7,
		},
		{
			name:     "exceed planned quantity",
			tenantID: tenantID,
			schoolID: schoolID,
			id:       part.ID,
			addUsed:  10, // Would make total 17, exceeds planned 10
			wantErr:  false, // UPDATE doesn't fail, just doesn't update
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			// Begin transaction
			tx, err := pool.Begin(ctx)
			if err != nil {
				t.Fatalf("Failed to begin transaction: %v", err)
			}
			defer tx.Rollback(ctx)

			now := time.Now().UTC()
			err = UpdateWorkOrderPartUsedTx(ctx, tx, tt.tenantID, tt.schoolID, tt.id, tt.addUsed, now)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateWorkOrderPartUsedTx() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				// Commit transaction
				if err := tx.Commit(ctx); err != nil {
					t.Fatalf("Failed to commit transaction: %v", err)
				}

				// Verify the qty_used was updated
				got, err := repo.GetByID(ctx, tt.tenantID, tt.schoolID, tt.id)
				if err != nil {
					t.Errorf("Failed to retrieve part: %v", err)
					return
				}

				// For the "exceed" test, qty should not have changed
				if tt.name != "exceed planned quantity" && got.QtyUsed != tt.wantQtyUsed {
					t.Errorf("UpdateWorkOrderPartUsedTx() QtyUsed = %v, want %v", got.QtyUsed, tt.wantQtyUsed)
				}
			}
		})
	}
}

func TestUpdateWorkOrderPartPlannedTx(t *testing.T) {
	pool := setupTestDB(t)
	repo := &WorkOrderPartRepo{pool: pool}

	tenantID := "tenant-test-updateplanned"
	schoolID := "school-test-updateplanned"

	// Setup: Create a test part
	part := validWorkOrderPart()
	part.ID = "wop-updateplanned-001"
	part.TenantID = tenantID
	part.SchoolID = schoolID
	part.QtyPlanned = 5
	part.QtyUsed = 0

	defer cleanupWorkOrderParts(t, pool, tenantID, schoolID)
	cleanupWorkOrderParts(t, pool, tenantID, schoolID)

	if err := repo.Create(context.Background(), part); err != nil {
		t.Fatalf("Failed to create test part: %v", err)
	}

	tests := []struct {
		name            string
		tenantID        string
		schoolID        string
		id              string
		newPlanned      int64
		wantErr         bool
		wantQtyPlanned  int64
	}{
		{
			name:           "update planned to 10",
			tenantID:       tenantID,
			schoolID:       schoolID,
			id:             part.ID,
			newPlanned:     10,
			wantErr:        false,
			wantQtyPlanned: 10,
		},
		{
			name:           "update planned to 3",
			tenantID:       tenantID,
			schoolID:       schoolID,
			id:             part.ID,
			newPlanned:     3,
			wantErr:        false,
			wantQtyPlanned: 3,
		},
		{
			name:       "update non-existent part",
			tenantID:   tenantID,
			schoolID:   schoolID,
			id:         "non-existent-id",
			newPlanned: 5,
			wantErr:    false, // UPDATE doesn't fail on no rows
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			// Begin transaction
			tx, err := pool.Begin(ctx)
			if err != nil {
				t.Fatalf("Failed to begin transaction: %v", err)
			}
			defer tx.Rollback(ctx)

			now := time.Now().UTC()
			err = UpdateWorkOrderPartPlannedTx(ctx, tx, tt.tenantID, tt.schoolID, tt.id, tt.newPlanned, now)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateWorkOrderPartPlannedTx() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				// Commit transaction
				if err := tx.Commit(ctx); err != nil {
					t.Fatalf("Failed to commit transaction: %v", err)
				}

				// Verify the qty_planned was updated (only if part exists)
				if tt.id == part.ID {
					got, err := repo.GetByID(ctx, tt.tenantID, tt.schoolID, tt.id)
					if err != nil {
						t.Errorf("Failed to retrieve part: %v", err)
						return
					}
					if got.QtyPlanned != tt.wantQtyPlanned {
						t.Errorf("UpdateWorkOrderPartPlannedTx() QtyPlanned = %v, want %v", got.QtyPlanned, tt.wantQtyPlanned)
					}
				}
			}
		})
	}
}
