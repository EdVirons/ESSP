package store

import (
	"context"
	"testing"
	"time"

	"github.com/edvirons/ssp/ims/internal/models"
)

func TestIncidentRepo_Create(t *testing.T) {
	pool := setupTestDB(t)
	repo := &IncidentRepo{pool: pool}

	tests := []struct {
		name    string
		input   models.Incident
		wantErr bool
	}{
		{
			name:    "valid incident",
			input:   validIncident(),
			wantErr: false,
		},
		{
			name: "incident with different ID",
			input: func() models.Incident {
				inc := validIncident()
				inc.ID = "inc-test-create-002"
				return inc
			}(),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Cleanup before and after test
			defer cleanupIncidents(t, pool, tt.input.TenantID, tt.input.SchoolID)
			cleanupIncidents(t, pool, tt.input.TenantID, tt.input.SchoolID)

			err := repo.Create(context.Background(), tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Verify creation if no error expected
			if !tt.wantErr && err == nil {
				got, err := repo.GetByID(context.Background(), tt.input.TenantID, tt.input.SchoolID, tt.input.ID)
				if err != nil {
					t.Errorf("Failed to retrieve created incident: %v", err)
					return
				}
				if got.ID != tt.input.ID {
					t.Errorf("Created incident ID = %v, want %v", got.ID, tt.input.ID)
				}
				if got.Title != tt.input.Title {
					t.Errorf("Created incident Title = %v, want %v", got.Title, tt.input.Title)
				}
			}
		})
	}
}

func TestIncidentRepo_GetByID(t *testing.T) {
	pool := setupTestDB(t)
	repo := &IncidentRepo{pool: pool}

	// Setup: Create a test incident
	incident := validIncident()
	incident.ID = "inc-test-getbyid-001"
	defer cleanupIncidents(t, pool, incident.TenantID, incident.SchoolID)
	cleanupIncidents(t, pool, incident.TenantID, incident.SchoolID)

	err := repo.Create(context.Background(), incident)
	if err != nil {
		t.Fatalf("Failed to create test incident: %v", err)
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
			name:     "existing incident",
			tenantID: incident.TenantID,
			schoolID: incident.SchoolID,
			id:       incident.ID,
			wantErr:  false,
			wantID:   incident.ID,
		},
		{
			name:     "non-existent incident",
			tenantID: incident.TenantID,
			schoolID: incident.SchoolID,
			id:       "non-existent-id",
			wantErr:  true,
		},
		{
			name:     "wrong tenant",
			tenantID: "wrong-tenant",
			schoolID: incident.SchoolID,
			id:       incident.ID,
			wantErr:  true,
		},
		{
			name:     "wrong school",
			tenantID: incident.TenantID,
			schoolID: "wrong-school",
			id:       incident.ID,
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

func TestIncidentRepo_List(t *testing.T) {
	pool := setupTestDB(t)
	repo := &IncidentRepo{pool: pool}

	tenantID := "tenant-test-list"
	schoolID := "school-test-list"

	// Setup: Create multiple test incidents
	defer cleanupIncidents(t, pool, tenantID, schoolID)
	cleanupIncidents(t, pool, tenantID, schoolID)

	now := time.Now().UTC()
	incidents := []models.Incident{
		func() models.Incident {
			inc := validIncident()
			inc.ID = "inc-list-001"
			inc.TenantID = tenantID
			inc.SchoolID = schoolID
			inc.Status = models.IncidentNew
			inc.DeviceID = "device-001"
			inc.Title = "First incident"
			inc.CreatedAt = now.Add(-3 * time.Hour)
			inc.UpdatedAt = now.Add(-3 * time.Hour)
			return inc
		}(),
		func() models.Incident {
			inc := validIncident()
			inc.ID = "inc-list-002"
			inc.TenantID = tenantID
			inc.SchoolID = schoolID
			inc.Status = models.IncidentInProgress
			inc.DeviceID = "device-002"
			inc.Title = "Second incident"
			inc.CreatedAt = now.Add(-2 * time.Hour)
			inc.UpdatedAt = now.Add(-2 * time.Hour)
			return inc
		}(),
		func() models.Incident {
			inc := validIncident()
			inc.ID = "inc-list-003"
			inc.TenantID = tenantID
			inc.SchoolID = schoolID
			inc.Status = models.IncidentNew
			inc.DeviceID = "device-001"
			inc.Title = "Third incident"
			inc.Description = "Contains search term"
			inc.CreatedAt = now.Add(-1 * time.Hour)
			inc.UpdatedAt = now.Add(-1 * time.Hour)
			return inc
		}(),
	}

	for _, inc := range incidents {
		if err := repo.Create(context.Background(), inc); err != nil {
			t.Fatalf("Failed to create test incident %s: %v", inc.ID, err)
		}
	}

	tests := []struct {
		name      string
		params    IncidentListParams
		wantCount int
		wantNext  bool // whether to expect a next cursor
	}{
		{
			name: "list all incidents",
			params: IncidentListParams{
				TenantID: tenantID,
				SchoolID: schoolID,
				Limit:    10,
			},
			wantCount: 3,
			wantNext:  false,
		},
		{
			name: "filter by status",
			params: IncidentListParams{
				TenantID: tenantID,
				SchoolID: schoolID,
				Status:   string(models.IncidentNew),
				Limit:    10,
			},
			wantCount: 2,
			wantNext:  false,
		},
		{
			name: "filter by device",
			params: IncidentListParams{
				TenantID: tenantID,
				SchoolID: schoolID,
				DeviceID: "device-001",
				Limit:    10,
			},
			wantCount: 2,
			wantNext:  false,
		},
		{
			name: "search by query",
			params: IncidentListParams{
				TenantID: tenantID,
				SchoolID: schoolID,
				Query:    "search term",
				Limit:    10,
			},
			wantCount: 1,
			wantNext:  false,
		},
		{
			name: "pagination with limit",
			params: IncidentListParams{
				TenantID: tenantID,
				SchoolID: schoolID,
				Limit:    2,
			},
			wantCount: 2,
			wantNext:  true,
		},
		{
			name: "pagination with cursor",
			params: func() IncidentListParams {
				// Get first page to obtain cursor
				firstPage, nextCursor, err := repo.List(context.Background(), IncidentListParams{
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

				return IncidentListParams{
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

func TestIncidentRepo_UpdateStatus(t *testing.T) {
	pool := setupTestDB(t)
	repo := &IncidentRepo{pool: pool}

	// Setup: Create a test incident
	incident := validIncident()
	incident.ID = "inc-test-updatestatus-001"
	incident.Status = models.IncidentNew
	defer cleanupIncidents(t, pool, incident.TenantID, incident.SchoolID)
	cleanupIncidents(t, pool, incident.TenantID, incident.SchoolID)

	err := repo.Create(context.Background(), incident)
	if err != nil {
		t.Fatalf("Failed to create test incident: %v", err)
	}

	tests := []struct {
		name       string
		tenantID   string
		schoolID   string
		id         string
		newStatus  models.IncidentStatus
		wantErr    bool
		wantStatus models.IncidentStatus
	}{
		{
			name:       "update to in_progress",
			tenantID:   incident.TenantID,
			schoolID:   incident.SchoolID,
			id:         incident.ID,
			newStatus:  models.IncidentInProgress,
			wantErr:    false,
			wantStatus: models.IncidentInProgress,
		},
		{
			name:      "update non-existent incident",
			tenantID:  incident.TenantID,
			schoolID:  incident.SchoolID,
			id:        "non-existent-id",
			newStatus: models.IncidentResolved,
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
			if !tt.wantErr && !got.UpdatedAt.After(incident.CreatedAt) {
				t.Errorf("UpdateStatus() UpdatedAt was not updated")
			}
		})
	}
}

func TestIncidentRepo_MarkSLABreaches(t *testing.T) {
	pool := setupTestDB(t)
	repo := &IncidentRepo{pool: pool}

	tenantID := "tenant-test-sla"
	schoolID := "school-test-sla"

	// Setup: Create incidents with different SLA states
	defer cleanupIncidents(t, pool, tenantID, schoolID)
	cleanupIncidents(t, pool, tenantID, schoolID)

	now := time.Now().UTC()
	incidents := []models.Incident{
		func() models.Incident {
			inc := validIncident()
			inc.ID = "inc-sla-001"
			inc.TenantID = tenantID
			inc.SchoolID = schoolID
			inc.Status = models.IncidentNew
			inc.SLADueAt = now.Add(-1 * time.Hour) // Breached
			inc.SLABreached = false
			return inc
		}(),
		func() models.Incident {
			inc := validIncident()
			inc.ID = "inc-sla-002"
			inc.TenantID = tenantID
			inc.SchoolID = schoolID
			inc.Status = models.IncidentInProgress
			inc.SLADueAt = now.Add(-2 * time.Hour) // Breached
			inc.SLABreached = false
			return inc
		}(),
		func() models.Incident {
			inc := validIncident()
			inc.ID = "inc-sla-003"
			inc.TenantID = tenantID
			inc.SchoolID = schoolID
			inc.Status = models.IncidentResolved // Should not be marked
			inc.SLADueAt = now.Add(-1 * time.Hour)
			inc.SLABreached = false
			return inc
		}(),
		func() models.Incident {
			inc := validIncident()
			inc.ID = "inc-sla-004"
			inc.TenantID = tenantID
			inc.SchoolID = schoolID
			inc.Status = models.IncidentNew
			inc.SLADueAt = now.Add(1 * time.Hour) // Not breached yet
			inc.SLABreached = false
			return inc
		}(),
	}

	for _, inc := range incidents {
		if err := repo.Create(context.Background(), inc); err != nil {
			t.Fatalf("Failed to create test incident %s: %v", inc.ID, err)
		}
	}

	// Run MarkSLABreaches
	count, err := repo.MarkSLABreaches(context.Background(), now)
	if err != nil {
		t.Fatalf("MarkSLABreaches() error = %v", err)
	}

	// We expect 2 incidents to be marked (inc-sla-001 and inc-sla-002)
	wantCount := 2
	if count != wantCount {
		t.Errorf("MarkSLABreaches() count = %v, want %v", count, wantCount)
	}

	// Verify the breached incidents are marked
	for _, incID := range []string{"inc-sla-001", "inc-sla-002"} {
		inc, err := repo.GetByID(context.Background(), tenantID, schoolID, incID)
		if err != nil {
			t.Errorf("Failed to get incident %s: %v", incID, err)
			continue
		}
		if !inc.SLABreached {
			t.Errorf("Incident %s should be marked as SLA breached", incID)
		}
	}

	// Verify the non-breached incidents are not marked
	for _, incID := range []string{"inc-sla-003", "inc-sla-004"} {
		inc, err := repo.GetByID(context.Background(), tenantID, schoolID, incID)
		if err != nil {
			t.Errorf("Failed to get incident %s: %v", incID, err)
			continue
		}
		if inc.SLABreached {
			t.Errorf("Incident %s should not be marked as SLA breached", incID)
		}
	}
}
