package seed

import (
	"time"

	"github.com/edvirons/ssp/ims/internal/models"
	"github.com/edvirons/ssp/ims/internal/store"
)

// SeedData holds all the seed data for development testing.
type SeedData struct {
	Schools      []models.SchoolSnapshot
	ServiceShops []models.ServiceShop
	Staff        []models.ServiceStaff
	Parts        []models.PartSnapshot
	Devices      []models.DeviceSnapshot
	Inventory    []models.InventoryItem
	Incidents    []models.Incident
	WorkOrders   []models.WorkOrder
	Schedules    []models.WorkOrderSchedule
	Deliverables []models.WorkOrderDeliverable
	Contacts     []models.SchoolContact
	Projects     []models.SchoolServiceProject
	Phases       []models.ServicePhase
}

// KenyanCounty represents a Kenyan county with code.
type KenyanCounty struct {
	Code      string
	Name      string
	SubCounty string
	SubCode   string
}

// Counties in Kenya for seed data.
var kenyanCounties = []KenyanCounty{
	{Code: "047", Name: "Nairobi", SubCounty: "Westlands", SubCode: "047-01"},
	{Code: "001", Name: "Mombasa", SubCounty: "Mvita", SubCode: "001-01"},
	{Code: "042", Name: "Kisumu", SubCounty: "Kisumu Central", SubCode: "042-01"},
	{Code: "032", Name: "Nakuru", SubCounty: "Nakuru Town East", SubCode: "032-01"},
	{Code: "027", Name: "Uasin Gishu", SubCounty: "Eldoret East", SubCode: "027-01"},
}

// GenerateSeedData generates all seed data for development.
func GenerateSeedData(tenantID string) *SeedData {
	now := time.Now().UTC()
	data := &SeedData{}

	// Generate IDs
	// Use 'demo-school' for the first school to match frontend default
	schoolIDs := make([]string, 5)
	schoolIDs[0] = "demo-school" // Must match frontend default in client.ts
	for i := 1; i < len(schoolIDs); i++ {
		schoolIDs[i] = store.NewID("sch")
	}

	shopIDs := make([]string, 3)
	for i := range shopIDs {
		shopIDs[i] = store.NewID("shop")
	}

	staffIDs := make([]string, 6)
	for i := range staffIDs {
		staffIDs[i] = store.NewID("staff")
	}

	partIDs := make([]string, 8)
	for i := range partIDs {
		partIDs[i] = store.NewID("part")
	}

	deviceIDs := make([]string, 10)
	for i := range deviceIDs {
		deviceIDs[i] = store.NewID("dev")
	}

	incidentIDs := make([]string, 8)
	for i := range incidentIDs {
		incidentIDs[i] = store.NewID("inc")
	}

	workOrderIDs := make([]string, 6)
	for i := range workOrderIDs {
		workOrderIDs[i] = store.NewID("wo")
	}

	projectIDs := make([]string, 2)
	for i := range projectIDs {
		projectIDs[i] = store.NewID("proj")
	}

	// 1. Schools (5 schools across Kenya)
	data.Schools = []models.SchoolSnapshot{
		{
			TenantID:      tenantID,
			SchoolID:      schoolIDs[0],
			Name:          "Nairobi Primary School",
			CountyCode:    kenyanCounties[0].Code,
			CountyName:    kenyanCounties[0].Name,
			SubCountyCode: kenyanCounties[0].SubCode,
			SubCountyName: kenyanCounties[0].SubCounty,
			Level:         "primary",
			Type:          "public",
			UpdatedAt:     now,
		},
		{
			TenantID:      tenantID,
			SchoolID:      schoolIDs[1],
			Name:          "Mombasa Girls Secondary",
			CountyCode:    kenyanCounties[1].Code,
			CountyName:    kenyanCounties[1].Name,
			SubCountyCode: kenyanCounties[1].SubCode,
			SubCountyName: kenyanCounties[1].SubCounty,
			Level:         "secondary",
			Type:          "public",
			Sex:           "girls",
			UpdatedAt:     now,
		},
		{
			TenantID:      tenantID,
			SchoolID:      schoolIDs[2],
			Name:          "Kisumu Boys High School",
			CountyCode:    kenyanCounties[2].Code,
			CountyName:    kenyanCounties[2].Name,
			SubCountyCode: kenyanCounties[2].SubCode,
			SubCountyName: kenyanCounties[2].SubCounty,
			Level:         "secondary",
			Type:          "public",
			Sex:           "boys",
			UpdatedAt:     now,
		},
		{
			TenantID:      tenantID,
			SchoolID:      schoolIDs[3],
			Name:          "Nakuru Academy",
			CountyCode:    kenyanCounties[3].Code,
			CountyName:    kenyanCounties[3].Name,
			SubCountyCode: kenyanCounties[3].SubCode,
			SubCountyName: kenyanCounties[3].SubCounty,
			Level:         "secondary",
			Type:          "private",
			Sex:           "mixed",
			UpdatedAt:     now,
		},
		{
			TenantID:      tenantID,
			SchoolID:      schoolIDs[4],
			Name:          "Eldoret Technical Institute",
			CountyCode:    kenyanCounties[4].Code,
			CountyName:    kenyanCounties[4].Name,
			SubCountyCode: kenyanCounties[4].SubCode,
			SubCountyName: kenyanCounties[4].SubCounty,
			Level:         "tertiary",
			Type:          "public",
			UpdatedAt:     now,
		},
	}

	// 2. Service Shops (3 shops in different regions)
	data.ServiceShops = []models.ServiceShop{
		{
			ID:            shopIDs[0],
			TenantID:      tenantID,
			CountyCode:    kenyanCounties[0].Code,
			CountyName:    kenyanCounties[0].Name,
			SubCountyCode: kenyanCounties[0].SubCode,
			SubCountyName: kenyanCounties[0].SubCounty,
			CoverageLevel: "county",
			Name:          "Nairobi Service Center",
			Location:      "Westlands, Nairobi",
			Active:        true,
			CreatedAt:     now,
			UpdatedAt:     now,
		},
		{
			ID:            shopIDs[1],
			TenantID:      tenantID,
			CountyCode:    kenyanCounties[1].Code,
			CountyName:    kenyanCounties[1].Name,
			SubCountyCode: kenyanCounties[1].SubCode,
			SubCountyName: kenyanCounties[1].SubCounty,
			CoverageLevel: "regional",
			Name:          "Coast Repair Hub",
			Location:      "Moi Avenue, Mombasa",
			Active:        true,
			CreatedAt:     now,
			UpdatedAt:     now,
		},
		{
			ID:            shopIDs[2],
			TenantID:      tenantID,
			CountyCode:    kenyanCounties[3].Code,
			CountyName:    kenyanCounties[3].Name,
			SubCountyCode: kenyanCounties[3].SubCode,
			SubCountyName: kenyanCounties[3].SubCounty,
			CoverageLevel: "regional",
			Name:          "Rift Valley Tech Shop",
			Location:      "Kenyatta Avenue, Nakuru",
			Active:        true,
			CreatedAt:     now,
			UpdatedAt:     now,
		},
	}

	// 3. Service Staff (6 staff, 2 per shop)
	staffNames := []string{"John Kamau", "Mary Wanjiku", "Peter Ochieng", "Jane Akinyi", "David Kiprop", "Grace Chebet"}
	staffPhones := []string{"+254722000001", "+254722000002", "+254722000003", "+254722000004", "+254722000005", "+254722000006"}
	staffRoles := []models.StaffRole{
		models.StaffRoleLeadTechnician, models.StaffRoleAssistantTechnician,
		models.StaffRoleLeadTechnician, models.StaffRoleStorekeeper,
		models.StaffRoleLeadTechnician, models.StaffRoleAssistantTechnician,
	}

	for i := 0; i < 6; i++ {
		data.Staff = append(data.Staff, models.ServiceStaff{
			ID:            staffIDs[i],
			TenantID:      tenantID,
			ServiceShopID: shopIDs[i/2],
			UserID:        store.NewID("user"),
			Role:          staffRoles[i],
			Phone:         staffPhones[i],
			Active:        true,
			CreatedAt:     now,
			UpdatedAt:     now,
		})
	}

	// Store staff names for later use
	_ = staffNames

	// 4. Parts (8 common spare parts)
	partsData := []struct {
		SKU      string
		Name     string
		Category string
		Unit     string
	}{
		{"SCR-LCD-15", "15.6\" LCD Screen", "screens", "piece"},
		{"SCR-LCD-14", "14\" LCD Screen", "screens", "piece"},
		{"BAT-DELL-65", "Dell 65Wh Battery", "batteries", "piece"},
		{"BAT-HP-45", "HP 45Wh Battery", "batteries", "piece"},
		{"KBD-US-STD", "US Standard Keyboard", "keyboards", "piece"},
		{"CHG-UNIV-65", "Universal 65W Charger", "chargers", "piece"},
		{"RAM-DDR4-8", "8GB DDR4 RAM", "memory", "piece"},
		{"SSD-256-SATA", "256GB SATA SSD", "storage", "piece"},
	}

	for i, p := range partsData {
		data.Parts = append(data.Parts, models.PartSnapshot{
			TenantID:  tenantID,
			PartID:    partIDs[i],
			PUK:       p.SKU,
			Name:      p.Name,
			Category:  p.Category,
			Unit:      p.Unit,
			UpdatedAt: now,
		})
	}

	// 5. Devices (10 devices across schools)
	deviceModels := []struct {
		Make     string
		Model    string
		Category string
	}{
		{"Dell", "Latitude 5520", "laptop"},
		{"Dell", "Latitude 3520", "laptop"},
		{"HP", "ProBook 450 G8", "laptop"},
		{"HP", "ProBook 440 G8", "laptop"},
		{"Lenovo", "ThinkPad E15", "laptop"},
		{"Lenovo", "ThinkPad L15", "laptop"},
		{"Dell", "OptiPlex 3090", "desktop"},
		{"HP", "ProDesk 400 G7", "desktop"},
		{"Lenovo", "ThinkCentre M70q", "desktop"},
		{"Apple", "MacBook Air M1", "laptop"},
	}

	// Put all devices in the first school so they're visible in dev mode
	for i := 0; i < 10; i++ {
		dm := deviceModels[i]
		data.Devices = append(data.Devices, models.DeviceSnapshot{
			TenantID:  tenantID,
			DeviceID:  deviceIDs[i],
			SchoolID:  schoolIDs[0], // All devices go to first school for dev visibility
			Model:     dm.Make + " " + dm.Model,
			Serial:    "SN" + string(rune('A'+i)) + "2024" + padInt(i+1, 4),
			AssetTag:  "ASSET-" + padInt(i+1, 5),
			Status:    "active",
			UpdatedAt: now,
		})
	}

	// 6. Inventory (parts in stock at shops)
	inventoryIdx := 0
	for shopIdx := 0; shopIdx < 3; shopIdx++ {
		for partIdx := 0; partIdx < 8; partIdx++ {
			qty := int64((shopIdx+1)*5 + partIdx*2)
			data.Inventory = append(data.Inventory, models.InventoryItem{
				ID:               store.NewID("inv"),
				TenantID:         tenantID,
				ServiceShopID:    shopIDs[shopIdx],
				PartID:           partIDs[partIdx],
				QtyAvailable:     qty,
				QtyReserved:      0,
				ReorderThreshold: 5,
				UpdatedAt:        now,
			})
			inventoryIdx++
		}
	}

	// 7. Incidents (8 incidents with varying severity and status)
	incidentData := []struct {
		DeviceIdx   int
		Category    string
		Severity    models.Severity
		Status      models.IncidentStatus
		Title       string
		Description string
		DaysAgo     int
	}{
		{0, "hardware", models.SeverityHigh, models.IncidentNew, "Screen cracked", "Laptop screen has visible cracks after being dropped", 1},
		{1, "hardware", models.SeverityMedium, models.IncidentAcknowledged, "Battery not charging", "Device not charging when plugged in", 3},
		{2, "hardware", models.SeverityLow, models.IncidentInProgress, "Keyboard malfunction", "Several keys not responding", 5},
		{3, "software", models.SeverityMedium, models.IncidentInProgress, "OS boot failure", "Device fails to boot past BIOS", 7},
		{4, "hardware", models.SeverityCritical, models.IncidentEscalated, "Device not turning on", "No response when power button pressed", 2},
		{5, "hardware", models.SeverityMedium, models.IncidentResolved, "Overheating issue", "Device shuts down due to overheating", 14},
		{6, "hardware", models.SeverityLow, models.IncidentClosed, "USB port not working", "Front USB ports not recognizing devices", 21},
		{7, "network", models.SeverityHigh, models.IncidentNew, "WiFi connectivity issues", "Intermittent WiFi disconnections", 0},
	}

	for i, inc := range incidentData {
		school := data.Schools[0] // All incidents go to first school for dev visibility
		device := data.Devices[inc.DeviceIdx]
		createdAt := now.Add(-time.Duration(inc.DaysAgo) * 24 * time.Hour)
		slaDue := createdAt.Add(48 * time.Hour) // 48 hour SLA

		data.Incidents = append(data.Incidents, models.Incident{
			ID:             incidentIDs[i],
			TenantID:       tenantID,
			SchoolID:       school.SchoolID,
			DeviceID:       device.DeviceID,
			SchoolName:     school.Name,
			CountyID:       school.CountyCode,
			CountyName:     school.CountyName,
			SubCountyID:    school.SubCountyCode,
			SubCountyName:  school.SubCountyName,
			DeviceSerial:   device.Serial,
			DeviceAssetTag: device.AssetTag,
			DeviceMake:     deviceModels[inc.DeviceIdx].Make,
			DeviceModel:    deviceModels[inc.DeviceIdx].Model,
			DeviceCategory: deviceModels[inc.DeviceIdx].Category,
			Category:       inc.Category,
			Severity:       inc.Severity,
			Status:         inc.Status,
			Title:          inc.Title,
			Description:    inc.Description,
			ReportedBy:     store.NewID("user"),
			SLADueAt:       slaDue,
			SLABreached:    slaDue.Before(now),
			CreatedAt:      createdAt,
			UpdatedAt:      now,
		})
	}

	// 8. Work Orders (6 work orders in different statuses)
	workOrderData := []struct {
		IncidentIdx    int
		Status         models.WorkOrderStatus
		ShopIdx        int
		StaffIdx       int
		TaskType       string
		CostCents      int64
		Notes          string
		RepairLocation models.RepairLocation
	}{
		{0, models.WorkOrderDraft, 0, 0, "repair", 15000, "Screen replacement needed", models.RepairLocationServiceShop},
		{1, models.WorkOrderAssigned, 0, 1, "repair", 8500, "Battery replacement", models.RepairLocationServiceShop},
		{2, models.WorkOrderInRepair, 1, 2, "repair", 5000, "Keyboard replacement in progress", models.RepairLocationServiceShop},
		{3, models.WorkOrderQA, 1, 3, "diagnosis", 3500, "OS reinstallation completed, testing", models.RepairLocationOnSite},
		{4, models.WorkOrderCompleted, 2, 4, "repair", 25000, "Motherboard repair completed", models.RepairLocationServiceShop},
		{5, models.WorkOrderApproved, 2, 5, "inspection", 2000, "Thermal paste replaced, signed off", models.RepairLocationOnSite},
	}

	for i, wo := range workOrderData {
		incident := data.Incidents[wo.IncidentIdx]
		shop := data.ServiceShops[wo.ShopIdx]
		staff := data.Staff[wo.StaffIdx]
		createdAt := incident.CreatedAt.Add(2 * time.Hour)

		data.WorkOrders = append(data.WorkOrders, models.WorkOrder{
			ID:                workOrderIDs[i],
			IncidentID:        incident.ID,
			TenantID:          tenantID,
			SchoolID:          incident.SchoolID,
			DeviceID:          incident.DeviceID,
			SchoolName:        incident.SchoolName,
			ContactName:       "School Admin",
			ContactPhone:      "+254700000000",
			DeviceSerial:      incident.DeviceSerial,
			DeviceAssetTag:    incident.DeviceAssetTag,
			DeviceMake:        incident.DeviceMake,
			DeviceModel:       incident.DeviceModel,
			DeviceCategory:    incident.DeviceCategory,
			Status:            wo.Status,
			ServiceShopID:     shop.ID,
			AssignedStaffID:   staff.ID,
			RepairLocation:    wo.RepairLocation,
			AssignedTo:        staff.UserID,
			TaskType:          wo.TaskType,
			CostEstimateCents: wo.CostCents,
			Notes:             wo.Notes,
			CreatedAt:         createdAt,
			UpdatedAt:         now,
		})
	}

	// 9. Schedules (3 scheduled work orders)
	scheduleData := []struct {
		WorkOrderIdx int
		DaysFromNow  int
		Notes        string
	}{
		{1, 1, "Scheduled for battery replacement"},
		{2, 2, "Keyboard parts arriving tomorrow"},
		{3, 0, "On-site visit scheduled for today"},
	}

	for _, sched := range scheduleData {
		wo := data.WorkOrders[sched.WorkOrderIdx]
		staff := data.Staff[sched.WorkOrderIdx]
		schedStart := now.Add(time.Duration(sched.DaysFromNow) * 24 * time.Hour).Truncate(time.Hour)
		schedEnd := schedStart.Add(2 * time.Hour)

		data.Schedules = append(data.Schedules, models.WorkOrderSchedule{
			ID:              store.NewID("sched"),
			TenantID:        tenantID,
			SchoolID:        wo.SchoolID,
			WorkOrderID:     wo.ID,
			ScheduledStart:  &schedStart,
			ScheduledEnd:    &schedEnd,
			Timezone:        "Africa/Nairobi",
			Notes:           sched.Notes,
			CreatedByUserID: staff.UserID,
			CreatedAt:       now,
		})
	}

	// 10. Deliverables (4 deliverables for completed/QA work orders)
	deliverableData := []struct {
		WorkOrderIdx int
		Title        string
		Description  string
		Status       models.DeliverableStatus
	}{
		{3, "Diagnostic report", "Initial diagnosis findings", models.DeliverableSubmitted},
		{4, "Repair photos", "Before and after photos of repair", models.DeliverableApproved},
		{4, "Test results", "Functional test results", models.DeliverableApproved},
		{5, "Inspection checklist", "Completed inspection checklist", models.DeliverableApproved},
	}

	for _, del := range deliverableData {
		wo := data.WorkOrders[del.WorkOrderIdx]
		data.Deliverables = append(data.Deliverables, models.WorkOrderDeliverable{
			ID:          store.NewID("deliv"),
			TenantID:    tenantID,
			SchoolID:    wo.SchoolID,
			WorkOrderID: wo.ID,
			Title:       del.Title,
			Description: del.Description,
			Status:      del.Status,
			CreatedAt:   now,
			UpdatedAt:   now,
		})
	}

	// 11. School Contacts (1-2 per school)
	contactData := []struct {
		SchoolIdx int
		Name      string
		Phone     string
		Email     string
		Role      string
		IsPrimary bool
	}{
		{0, "James Mwangi", "+254711111111", "jmwangi@nairobi-primary.ac.ke", "ICT Teacher", true},
		{0, "Sarah Njeri", "+254711111112", "snjeri@nairobi-primary.ac.ke", "Deputy Principal", false},
		{1, "Fatuma Hassan", "+254711111113", "fhassan@mombasa-girls.ac.ke", "Principal", true},
		{2, "Brian Otieno", "+254711111114", "botieno@kisumu-boys.ac.ke", "Lab Technician", true},
		{3, "Anne Wambui", "+254711111115", "awambui@nakuru-academy.ac.ke", "IT Administrator", true},
		{4, "Robert Kibet", "+254711111116", "rkibet@eldoret-tech.ac.ke", "HOD ICT", true},
	}

	for _, c := range contactData {
		school := data.Schools[c.SchoolIdx]
		data.Contacts = append(data.Contacts, models.SchoolContact{
			ID:        store.NewID("contact"),
			TenantID:  tenantID,
			SchoolID:  school.SchoolID,
			UserID:    store.NewID("user"),
			Name:      c.Name,
			Phone:     c.Phone,
			Email:     c.Email,
			Role:      c.Role,
			IsPrimary: c.IsPrimary,
			Active:    true,
			CreatedAt: now,
			UpdatedAt: now,
		})
	}

	// 12. Projects (2 projects)
	data.Projects = []models.SchoolServiceProject{
		{
			ID:                   projectIDs[0],
			TenantID:             tenantID,
			SchoolID:             schoolIDs[0],
			ProjectType:          models.ProjectTypeFullInstallation,
			Status:               models.ProjectActive,
			CurrentPhase:         models.PhaseInstall,
			StartDate:            now.Add(-30 * 24 * time.Hour).Format("2006-01-02"),
			GoLiveDate:           now.Add(60 * 24 * time.Hour).Format("2006-01-02"),
			AccountManagerUserID: store.NewID("user"),
			Notes:                "Full ICT lab installation project",
			CreatedAt:            now.Add(-30 * 24 * time.Hour),
			UpdatedAt:            now,
		},
		{
			ID:                   projectIDs[1],
			TenantID:             tenantID,
			SchoolID:             schoolIDs[2],
			ProjectType:          models.ProjectTypeDeviceRefresh,
			Status:               models.ProjectActive,
			CurrentPhase:         models.PhaseAssessment,
			StartDate:            now.Add(-7 * 24 * time.Hour).Format("2006-01-02"),
			GoLiveDate:           now.Add(30 * 24 * time.Hour).Format("2006-01-02"),
			AccountManagerUserID: store.NewID("user"),
			Notes:                "Upgrade aging laptops to new models",
			CreatedAt:            now.Add(-7 * 24 * time.Hour),
			UpdatedAt:            now,
		},
	}

	// 13. Service Phases (phases for projects)
	phaseData := []struct {
		ProjectIdx int
		PhaseType  models.PhaseType
		Status     models.PhaseStatus
		OwnerRole  string
	}{
		{0, models.PhaseDemo, models.PhaseDone, "sales"},
		{0, models.PhaseSurvey, models.PhaseDone, "field_tech"},
		{0, models.PhaseProcurement, models.PhaseDone, "operations"},
		{0, models.PhaseInstall, models.PhaseInProgress, "field_tech"},
		{1, models.PhaseAssessment, models.PhaseInProgress, "field_tech"},
	}

	for _, p := range phaseData {
		project := data.Projects[p.ProjectIdx]
		data.Phases = append(data.Phases, models.ServicePhase{
			ID:          store.NewID("phase"),
			TenantID:    tenantID,
			ProjectID:   project.ID,
			PhaseType:   p.PhaseType,
			Status:      p.Status,
			OwnerRole:   p.OwnerRole,
			OwnerUserID: store.NewID("user"),
			StartDate:   project.StartDate,
			CreatedAt:   now,
			UpdatedAt:   now,
		})
	}

	return data
}

// padInt pads an integer with leading zeros.
func padInt(n, width int) string {
	s := ""
	for i := 0; i < width; i++ {
		s = string(rune('0'+n%10)) + s
		n /= 10
	}
	return s
}
