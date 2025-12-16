package seed

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Seeder handles seeding the database with test data.
type Seeder struct {
	pool     *pgxpool.Pool
	tenantID string
	verbose  bool
}

// NewSeeder creates a new seeder instance.
func NewSeeder(pool *pgxpool.Pool, tenantID string, verbose bool) *Seeder {
	return &Seeder{
		pool:     pool,
		tenantID: tenantID,
		verbose:  verbose,
	}
}

// Run seeds the database with test data.
func (s *Seeder) Run(ctx context.Context) error {
	s.log("Generating seed data for tenant: %s", s.tenantID)

	data := GenerateSeedData(s.tenantID)

	// Seed in dependency order
	if err := s.seedSchools(ctx, data); err != nil {
		return err
	}

	if err := s.seedServiceShops(ctx, data); err != nil {
		return err
	}

	if err := s.seedServiceStaff(ctx, data); err != nil {
		return err
	}

	if err := s.seedParts(ctx, data); err != nil {
		return err
	}

	if err := s.seedDevices(ctx, data); err != nil {
		return err
	}

	if err := s.seedInventory(ctx, data); err != nil {
		return err
	}

	if err := s.seedIncidents(ctx, data); err != nil {
		return err
	}

	if err := s.seedWorkOrders(ctx, data); err != nil {
		return err
	}

	if err := s.seedSchedules(ctx, data); err != nil {
		return err
	}

	if err := s.seedDeliverables(ctx, data); err != nil {
		return err
	}

	if err := s.seedContacts(ctx, data); err != nil {
		return err
	}

	if err := s.seedProjects(ctx, data); err != nil {
		return err
	}

	if err := s.seedPhases(ctx, data); err != nil {
		return err
	}

	s.log("Seeding completed successfully!")
	return nil
}

// Clean removes all seed data from the database.
func (s *Seeder) Clean(ctx context.Context) error {
	s.log("Cleaning database for tenant: %s", s.tenantID)

	// Tables to clean in reverse dependency order
	tables := []string{
		"service_phases",
		"school_service_projects",
		"work_order_deliverables",
		"work_order_schedules",
		"work_order_parts",
		"work_orders",
		"incidents",
		"school_contacts",
		"inventory",
		"service_staff",
		"service_shops",
		"devices_snapshot",
		"parts_snapshot",
		"schools_snapshot",
	}

	for _, table := range tables {
		query := fmt.Sprintf("DELETE FROM %s WHERE tenant_id = $1", table)
		result, err := s.pool.Exec(ctx, query, s.tenantID)
		if err != nil {
			return fmt.Errorf("failed to clean table %s: %w", table, err)
		}
		s.log("  Cleaned %s: %d rows deleted", table, result.RowsAffected())
	}

	s.log("Database cleaned successfully!")
	return nil
}

func (s *Seeder) log(format string, args ...interface{}) {
	if s.verbose {
		log.Printf(format, args...)
	}
}

func (s *Seeder) seedSchools(ctx context.Context, data *SeedData) error {
	s.log("  Seeding %d schools...", len(data.Schools))

	query := `
		INSERT INTO schools_snapshot
		(tenant_id, school_id, name, county_code, county_name, sub_county_code, sub_county_name, level, type, sex, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (tenant_id, school_id) DO UPDATE
		SET name = EXCLUDED.name, updated_at = EXCLUDED.updated_at
	`

	for _, school := range data.Schools {
		_, err := s.pool.Exec(ctx, query,
			school.TenantID, school.SchoolID, school.Name,
			school.CountyCode, school.CountyName,
			school.SubCountyCode, school.SubCountyName,
			school.Level, school.Type, school.Sex,
			school.UpdatedAt,
		)
		if err != nil {
			return fmt.Errorf("failed to seed school %s: %w", school.Name, err)
		}
	}
	return nil
}

func (s *Seeder) seedServiceShops(ctx context.Context, data *SeedData) error {
	s.log("  Seeding %d service shops...", len(data.ServiceShops))

	query := `
		INSERT INTO service_shops
		(id, tenant_id, county_code, county_name, sub_county_code, sub_county_name,
		 coverage_level, name, location, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		ON CONFLICT (id) DO NOTHING
	`

	for _, shop := range data.ServiceShops {
		_, err := s.pool.Exec(ctx, query,
			shop.ID, shop.TenantID, shop.CountyCode, shop.CountyName,
			shop.SubCountyCode, shop.SubCountyName, shop.CoverageLevel,
			shop.Name, shop.Location, shop.Active,
			shop.CreatedAt, shop.UpdatedAt,
		)
		if err != nil {
			return fmt.Errorf("failed to seed service shop %s: %w", shop.Name, err)
		}
	}
	return nil
}

func (s *Seeder) seedServiceStaff(ctx context.Context, data *SeedData) error {
	s.log("  Seeding %d service staff...", len(data.Staff))

	query := `
		INSERT INTO service_staff
		(id, tenant_id, service_shop_id, user_id, role, phone, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (id) DO NOTHING
	`

	for _, st := range data.Staff {
		_, err := s.pool.Exec(ctx, query,
			st.ID, st.TenantID, st.ServiceShopID,
			st.UserID, st.Role, st.Phone, st.Active,
			st.CreatedAt, st.UpdatedAt,
		)
		if err != nil {
			return fmt.Errorf("failed to seed staff: %w", err)
		}
	}
	return nil
}

func (s *Seeder) seedParts(ctx context.Context, data *SeedData) error {
	s.log("  Seeding %d parts...", len(data.Parts))

	query := `
		INSERT INTO parts_snapshot
		(tenant_id, part_id, puk, name, category, unit, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (tenant_id, part_id) DO UPDATE
		SET name = EXCLUDED.name, updated_at = EXCLUDED.updated_at
	`

	for _, part := range data.Parts {
		_, err := s.pool.Exec(ctx, query,
			part.TenantID, part.PartID, part.PUK,
			part.Name, part.Category, part.Unit,
			part.UpdatedAt,
		)
		if err != nil {
			return fmt.Errorf("failed to seed part %s: %w", part.Name, err)
		}
	}
	return nil
}

func (s *Seeder) seedDevices(ctx context.Context, data *SeedData) error {
	s.log("  Seeding %d devices...", len(data.Devices))

	query := `
		INSERT INTO devices_snapshot
		(tenant_id, device_id, school_id, model, serial, asset_tag, status, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (tenant_id, device_id) DO UPDATE
		SET model = EXCLUDED.model, updated_at = EXCLUDED.updated_at
	`

	for _, device := range data.Devices {
		_, err := s.pool.Exec(ctx, query,
			device.TenantID, device.DeviceID, device.SchoolID,
			device.Model, device.Serial, device.AssetTag,
			device.Status, device.UpdatedAt,
		)
		if err != nil {
			return fmt.Errorf("failed to seed device %s: %w", device.Serial, err)
		}
	}
	return nil
}

func (s *Seeder) seedInventory(ctx context.Context, data *SeedData) error {
	s.log("  Seeding %d inventory items...", len(data.Inventory))

	query := `
		INSERT INTO inventory
		(id, tenant_id, service_shop_id, part_id, qty_available, qty_reserved, reorder_threshold, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (id) DO NOTHING
	`

	for _, item := range data.Inventory {
		_, err := s.pool.Exec(ctx, query,
			item.ID, item.TenantID, item.ServiceShopID, item.PartID,
			item.QtyAvailable, item.QtyReserved, item.ReorderThreshold,
			item.UpdatedAt,
		)
		if err != nil {
			return fmt.Errorf("failed to seed inventory: %w", err)
		}
	}
	return nil
}

func (s *Seeder) seedIncidents(ctx context.Context, data *SeedData) error {
	s.log("  Seeding %d incidents...", len(data.Incidents))

	query := `
		INSERT INTO incidents
		(id, tenant_id, school_id, device_id, school_name, county_id, county_name,
		 sub_county_id, sub_county_name, device_serial, device_asset_tag,
		 device_make, device_model, device_category,
		 category, severity, status, title, description, reported_by,
		 sla_due_at, sla_breached, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24)
		ON CONFLICT (id) DO NOTHING
	`

	for _, inc := range data.Incidents {
		_, err := s.pool.Exec(ctx, query,
			inc.ID, inc.TenantID, inc.SchoolID, inc.DeviceID,
			inc.SchoolName, inc.CountyID, inc.CountyName,
			inc.SubCountyID, inc.SubCountyName,
			inc.DeviceSerial, inc.DeviceAssetTag,
			inc.DeviceMake, inc.DeviceModel, inc.DeviceCategory,
			inc.Category, inc.Severity, inc.Status,
			inc.Title, inc.Description, inc.ReportedBy,
			inc.SLADueAt, inc.SLABreached,
			inc.CreatedAt, inc.UpdatedAt,
		)
		if err != nil {
			return fmt.Errorf("failed to seed incident %s: %w", inc.Title, err)
		}
	}
	return nil
}

func (s *Seeder) seedWorkOrders(ctx context.Context, data *SeedData) error {
	s.log("  Seeding %d work orders...", len(data.WorkOrders))

	query := `
		INSERT INTO work_orders
		(id, incident_id, tenant_id, school_id, device_id, school_name,
		 contact_name, contact_phone, device_serial, device_asset_tag,
		 device_make, device_model, device_category,
		 status, service_shop_id, assigned_staff_id, repair_location, assigned_to,
		 task_type, cost_estimate_cents, notes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23)
		ON CONFLICT (id) DO NOTHING
	`

	for _, wo := range data.WorkOrders {
		_, err := s.pool.Exec(ctx, query,
			wo.ID, wo.IncidentID, wo.TenantID, wo.SchoolID, wo.DeviceID,
			wo.SchoolName, wo.ContactName, wo.ContactPhone,
			wo.DeviceSerial, wo.DeviceAssetTag,
			wo.DeviceMake, wo.DeviceModel, wo.DeviceCategory,
			wo.Status, wo.ServiceShopID, wo.AssignedStaffID,
			wo.RepairLocation, wo.AssignedTo,
			wo.TaskType, wo.CostEstimateCents, wo.Notes,
			wo.CreatedAt, wo.UpdatedAt,
		)
		if err != nil {
			return fmt.Errorf("failed to seed work order: %w", err)
		}
	}
	return nil
}

func (s *Seeder) seedSchedules(ctx context.Context, data *SeedData) error {
	s.log("  Seeding %d schedules...", len(data.Schedules))

	query := `
		INSERT INTO work_order_schedules
		(id, tenant_id, school_id, work_order_id, scheduled_start, scheduled_end,
		 timezone, notes, created_by_user_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (id) DO NOTHING
	`

	for _, sched := range data.Schedules {
		_, err := s.pool.Exec(ctx, query,
			sched.ID, sched.TenantID, sched.SchoolID, sched.WorkOrderID,
			sched.ScheduledStart, sched.ScheduledEnd,
			sched.Timezone, sched.Notes, sched.CreatedByUserID, sched.CreatedAt,
		)
		if err != nil {
			return fmt.Errorf("failed to seed schedule: %w", err)
		}
	}
	return nil
}

func (s *Seeder) seedDeliverables(ctx context.Context, data *SeedData) error {
	s.log("  Seeding %d deliverables...", len(data.Deliverables))

	query := `
		INSERT INTO work_order_deliverables
		(id, tenant_id, school_id, work_order_id, title, description, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (id) DO NOTHING
	`

	for _, del := range data.Deliverables {
		_, err := s.pool.Exec(ctx, query,
			del.ID, del.TenantID, del.SchoolID, del.WorkOrderID,
			del.Title, del.Description, del.Status,
			del.CreatedAt, del.UpdatedAt,
		)
		if err != nil {
			return fmt.Errorf("failed to seed deliverable: %w", err)
		}
	}
	return nil
}

func (s *Seeder) seedContacts(ctx context.Context, data *SeedData) error {
	s.log("  Seeding %d school contacts...", len(data.Contacts))

	query := `
		INSERT INTO school_contacts
		(id, tenant_id, school_id, user_id, name, phone, email, role, is_primary, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		ON CONFLICT (id) DO NOTHING
	`

	for _, contact := range data.Contacts {
		_, err := s.pool.Exec(ctx, query,
			contact.ID, contact.TenantID, contact.SchoolID, contact.UserID,
			contact.Name, contact.Phone, contact.Email, contact.Role,
			contact.IsPrimary, contact.Active,
			contact.CreatedAt, contact.UpdatedAt,
		)
		if err != nil {
			return fmt.Errorf("failed to seed contact %s: %w", contact.Name, err)
		}
	}
	return nil
}

func (s *Seeder) seedProjects(ctx context.Context, data *SeedData) error {
	s.log("  Seeding %d projects...", len(data.Projects))

	query := `
		INSERT INTO school_service_projects
		(id, tenant_id, school_id, project_type, status, current_phase,
		 start_date, go_live_date, account_manager_user_id, notes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		ON CONFLICT (id) DO NOTHING
	`

	for _, proj := range data.Projects {
		_, err := s.pool.Exec(ctx, query,
			proj.ID, proj.TenantID, proj.SchoolID,
			proj.ProjectType, proj.Status, proj.CurrentPhase,
			proj.StartDate, proj.GoLiveDate, proj.AccountManagerUserID,
			proj.Notes, proj.CreatedAt, proj.UpdatedAt,
		)
		if err != nil {
			return fmt.Errorf("failed to seed project: %w", err)
		}
	}
	return nil
}

func (s *Seeder) seedPhases(ctx context.Context, data *SeedData) error {
	s.log("  Seeding %d phases...", len(data.Phases))

	query := `
		INSERT INTO service_phases
		(id, tenant_id, project_id, phase_type, status, owner_role, owner_user_id, start_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (id) DO NOTHING
	`

	for _, phase := range data.Phases {
		_, err := s.pool.Exec(ctx, query,
			phase.ID, phase.TenantID, phase.ProjectID,
			phase.PhaseType, phase.Status, phase.OwnerRole, phase.OwnerUserID,
			phase.StartDate, phase.CreatedAt, phase.UpdatedAt,
		)
		if err != nil {
			return fmt.Errorf("failed to seed phase: %w", err)
		}
	}
	return nil
}
