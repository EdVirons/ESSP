package store

import (
	"context"
	"testing"
	"time"

	"github.com/edvirons/ssp/ims/internal/models"
)

func TestServiceShopRepo_Create(t *testing.T) {
	pool := setupTestDB(t)
	repo := &ServiceShopRepo{pool: pool}

	tests := []struct {
		name    string
		input   models.ServiceShop
		wantErr bool
	}{
		{
			name:    "valid service shop",
			input:   validServiceShop(),
			wantErr: false,
		},
		{
			name: "service shop with different ID",
			input: func() models.ServiceShop {
				shop := validServiceShop()
				shop.ID = "shop-test-create-002"
				return shop
			}(),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Cleanup before and after test
			defer cleanupServiceShops(t, pool, tt.input.TenantID)
			cleanupServiceShops(t, pool, tt.input.TenantID)

			err := repo.Create(context.Background(), tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Verify creation if no error expected
			if !tt.wantErr && err == nil {
				got, err := repo.GetByID(context.Background(), tt.input.TenantID, tt.input.ID)
				if err != nil {
					t.Errorf("Failed to retrieve created service shop: %v", err)
					return
				}
				if got.ID != tt.input.ID {
					t.Errorf("Created service shop ID = %v, want %v", got.ID, tt.input.ID)
				}
				if got.Name != tt.input.Name {
					t.Errorf("Created service shop Name = %v, want %v", got.Name, tt.input.Name)
				}
			}
		})
	}
}

func TestServiceShopRepo_GetByID(t *testing.T) {
	pool := setupTestDB(t)
	repo := &ServiceShopRepo{pool: pool}

	// Setup: Create a test service shop
	shop := validServiceShop()
	shop.ID = "shop-test-getbyid-001"
	defer cleanupServiceShops(t, pool, shop.TenantID)
	cleanupServiceShops(t, pool, shop.TenantID)

	err := repo.Create(context.Background(), shop)
	if err != nil {
		t.Fatalf("Failed to create test service shop: %v", err)
	}

	tests := []struct {
		name     string
		tenantID string
		id       string
		wantErr  bool
		wantID   string
	}{
		{
			name:     "existing service shop",
			tenantID: shop.TenantID,
			id:       shop.ID,
			wantErr:  false,
			wantID:   shop.ID,
		},
		{
			name:     "non-existent service shop",
			tenantID: shop.TenantID,
			id:       "non-existent-id",
			wantErr:  true,
		},
		{
			name:     "wrong tenant",
			tenantID: "wrong-tenant",
			id:       shop.ID,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.GetByID(context.Background(), tt.tenantID, tt.id)
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

func TestServiceShopRepo_GetByCounty(t *testing.T) {
	pool := setupTestDB(t)
	repo := &ServiceShopRepo{pool: pool}

	tenantID := "tenant-test-county"

	// Setup: Create service shops in different counties
	defer cleanupServiceShops(t, pool, tenantID)
	cleanupServiceShops(t, pool, tenantID)

	shops := []models.ServiceShop{
		func() models.ServiceShop {
			shop := validServiceShop()
			shop.ID = "shop-county-001"
			shop.TenantID = tenantID
			shop.CountyCode = "001"
			shop.CountyName = "Nairobi"
			shop.Active = true
			return shop
		}(),
		func() models.ServiceShop {
			shop := validServiceShop()
			shop.ID = "shop-county-002"
			shop.TenantID = tenantID
			shop.CountyCode = "002"
			shop.CountyName = "Mombasa"
			shop.Active = true
			return shop
		}(),
		func() models.ServiceShop {
			shop := validServiceShop()
			shop.ID = "shop-county-003"
			shop.TenantID = tenantID
			shop.CountyCode = "001"
			shop.CountyName = "Nairobi"
			shop.Active = false // Inactive
			return shop
		}(),
	}

	for _, shop := range shops {
		if err := repo.Create(context.Background(), shop); err != nil {
			t.Fatalf("Failed to create test service shop %s: %v", shop.ID, err)
		}
	}

	tests := []struct {
		name       string
		tenantID   string
		countyCode string
		wantErr    bool
		wantID     string
	}{
		{
			name:       "get active shop in Nairobi",
			tenantID:   tenantID,
			countyCode: "001",
			wantErr:    false,
			wantID:     "shop-county-001",
		},
		{
			name:       "get active shop in Mombasa",
			tenantID:   tenantID,
			countyCode: "002",
			wantErr:    false,
			wantID:     "shop-county-002",
		},
		{
			name:       "non-existent county",
			tenantID:   tenantID,
			countyCode: "999",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.GetByCounty(context.Background(), tt.tenantID, tt.countyCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByCounty() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.ID != tt.wantID {
				t.Errorf("GetByCounty() ID = %v, want %v", got.ID, tt.wantID)
			}
			// Verify the shop is active
			if !tt.wantErr && !got.Active {
				t.Errorf("GetByCounty() returned inactive shop")
			}
		})
	}
}

func TestServiceShopRepo_GetBySubCounty(t *testing.T) {
	pool := setupTestDB(t)
	repo := &ServiceShopRepo{pool: pool}

	tenantID := "tenant-test-subcounty"

	// Setup: Create service shops with different coverage levels
	defer cleanupServiceShops(t, pool, tenantID)
	cleanupServiceShops(t, pool, tenantID)

	shops := []models.ServiceShop{
		func() models.ServiceShop {
			shop := validServiceShop()
			shop.ID = "shop-subcounty-001"
			shop.TenantID = tenantID
			shop.CountyCode = "001"
			shop.SubCountyCode = "001-001"
			shop.SubCountyName = "Westlands"
			shop.CoverageLevel = "sub_county"
			shop.Active = true
			return shop
		}(),
		func() models.ServiceShop {
			shop := validServiceShop()
			shop.ID = "shop-subcounty-002"
			shop.TenantID = tenantID
			shop.CountyCode = "001"
			shop.SubCountyCode = "001-002"
			shop.SubCountyName = "Dagoretti"
			shop.CoverageLevel = "sub_county"
			shop.Active = true
			return shop
		}(),
		func() models.ServiceShop {
			shop := validServiceShop()
			shop.ID = "shop-subcounty-003"
			shop.TenantID = tenantID
			shop.CountyCode = "001"
			shop.SubCountyCode = "001-001"
			shop.SubCountyName = "Westlands"
			shop.CoverageLevel = "county" // Not sub_county level
			shop.Active = true
			return shop
		}(),
	}

	for _, shop := range shops {
		if err := repo.Create(context.Background(), shop); err != nil {
			t.Fatalf("Failed to create test service shop %s: %v", shop.ID, err)
		}
	}

	tests := []struct {
		name          string
		tenantID      string
		countyCode    string
		subCountyCode string
		wantErr       bool
		wantID        string
	}{
		{
			name:          "get shop in Westlands",
			tenantID:      tenantID,
			countyCode:    "001",
			subCountyCode: "001-001",
			wantErr:       false,
			wantID:        "shop-subcounty-001",
		},
		{
			name:          "get shop in Dagoretti",
			tenantID:      tenantID,
			countyCode:    "001",
			subCountyCode: "001-002",
			wantErr:       false,
			wantID:        "shop-subcounty-002",
		},
		{
			name:          "non-existent sub-county",
			tenantID:      tenantID,
			countyCode:    "001",
			subCountyCode: "001-999",
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.GetBySubCounty(context.Background(), tt.tenantID, tt.countyCode, tt.subCountyCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBySubCounty() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.ID != tt.wantID {
				t.Errorf("GetBySubCounty() ID = %v, want %v", got.ID, tt.wantID)
			}
			// Verify coverage level is sub_county
			if !tt.wantErr && got.CoverageLevel != "sub_county" {
				t.Errorf("GetBySubCounty() CoverageLevel = %v, want sub_county", got.CoverageLevel)
			}
			// Verify the shop is active
			if !tt.wantErr && !got.Active {
				t.Errorf("GetBySubCounty() returned inactive shop")
			}
		})
	}
}

func TestServiceShopRepo_List(t *testing.T) {
	pool := setupTestDB(t)
	repo := &ServiceShopRepo{pool: pool}

	tenantID := "tenant-test-shoplist"

	// Setup: Create multiple test service shops
	defer cleanupServiceShops(t, pool, tenantID)
	cleanupServiceShops(t, pool, tenantID)

	now := time.Now().UTC()
	shops := []models.ServiceShop{
		func() models.ServiceShop {
			shop := validServiceShop()
			shop.ID = "shop-list-001"
			shop.TenantID = tenantID
			shop.CountyCode = "001"
			shop.Active = true
			shop.CreatedAt = now.Add(-3 * time.Hour)
			shop.UpdatedAt = now.Add(-3 * time.Hour)
			return shop
		}(),
		func() models.ServiceShop {
			shop := validServiceShop()
			shop.ID = "shop-list-002"
			shop.TenantID = tenantID
			shop.CountyCode = "002"
			shop.Active = true
			shop.CreatedAt = now.Add(-2 * time.Hour)
			shop.UpdatedAt = now.Add(-2 * time.Hour)
			return shop
		}(),
		func() models.ServiceShop {
			shop := validServiceShop()
			shop.ID = "shop-list-003"
			shop.TenantID = tenantID
			shop.CountyCode = "001"
			shop.Active = false
			shop.CreatedAt = now.Add(-1 * time.Hour)
			shop.UpdatedAt = now.Add(-1 * time.Hour)
			return shop
		}(),
	}

	for _, shop := range shops {
		if err := repo.Create(context.Background(), shop); err != nil {
			t.Fatalf("Failed to create test service shop %s: %v", shop.ID, err)
		}
	}

	tests := []struct {
		name      string
		params    ShopListParams
		wantCount int
		wantNext  bool
	}{
		{
			name: "list all service shops",
			params: ShopListParams{
				TenantID:   tenantID,
				ActiveOnly: false,
				Limit:      10,
			},
			wantCount: 3,
			wantNext:  false,
		},
		{
			name: "list active shops only",
			params: ShopListParams{
				TenantID:   tenantID,
				ActiveOnly: true,
				Limit:      10,
			},
			wantCount: 2,
			wantNext:  false,
		},
		{
			name: "filter by county",
			params: ShopListParams{
				TenantID:   tenantID,
				CountyCode: "001",
				ActiveOnly: false,
				Limit:      10,
			},
			wantCount: 2,
			wantNext:  false,
		},
		{
			name: "filter by county - active only",
			params: ShopListParams{
				TenantID:   tenantID,
				CountyCode: "001",
				ActiveOnly: true,
				Limit:      10,
			},
			wantCount: 1,
			wantNext:  false,
		},
		{
			name: "pagination with limit",
			params: ShopListParams{
				TenantID:   tenantID,
				ActiveOnly: false,
				Limit:      2,
			},
			wantCount: 2,
			wantNext:  true,
		},
		{
			name: "pagination with cursor",
			params: func() ShopListParams {
				// Get first page to obtain cursor
				firstPage, nextCursor, err := repo.List(context.Background(), ShopListParams{
					TenantID: tenantID,
					Limit:    1,
				})
				if err != nil || len(firstPage) == 0 {
					t.Fatalf("Failed to get first page for cursor test: %v", err)
				}

				cursorTime, cursorID, ok := DecodeCursor(nextCursor)
				if !ok {
					t.Fatalf("Failed to decode cursor: %s", nextCursor)
				}

				return ShopListParams{
					TenantID:        tenantID,
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

			// Verify ActiveOnly filter
			if tt.params.ActiveOnly {
				for _, shop := range got {
					if !shop.Active {
						t.Errorf("List() with ActiveOnly returned inactive shop %s", shop.ID)
					}
				}
			}
		})
	}
}
