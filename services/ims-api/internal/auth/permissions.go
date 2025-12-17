package auth

// Permission constants for all IMS operations
const (
	// Incident permissions
	PermIncidentCreate = "incident:create"
	PermIncidentRead   = "incident:read"
	PermIncidentUpdate = "incident:update"
	PermIncidentDelete = "incident:delete"

	// Work Order permissions
	PermWorkOrderCreate = "workorder:create"
	PermWorkOrderRead   = "workorder:read"
	PermWorkOrderUpdate = "workorder:update"
	PermWorkOrderDelete = "workorder:delete"

	// Work Order Operations (scheduling, deliverables, approvals)
	PermWorkOrderSchedule     = "workorder:schedule"
	PermWorkOrderDeliverable  = "workorder:deliverable"
	PermWorkOrderApproval     = "workorder:approval"
	PermWorkOrderReview       = "workorder:review"

	// BOM permissions
	PermBOMCreate = "bom:create"
	PermBOMRead   = "bom:read"
	PermBOMUpdate = "bom:update"
	PermBOMConsume = "bom:consume"

	// School permissions
	PermSchoolCreate = "school:create"
	PermSchoolRead   = "school:read"
	PermSchoolUpdate = "school:update"

	// School Contacts permissions
	PermSchoolContactCreate = "school:contact:create"
	PermSchoolContactRead   = "school:contact:read"
	PermSchoolContactUpdate = "school:contact:update"

	// Attachment permissions
	PermAttachmentCreate = "attachment:create"
	PermAttachmentRead   = "attachment:read"
	PermAttachmentDelete = "attachment:delete"

	// Service Shop permissions
	PermServiceShopCreate = "serviceshop:create"
	PermServiceShopRead   = "serviceshop:read"
	PermServiceShopUpdate = "serviceshop:update"

	// Service Staff permissions
	PermServiceStaffCreate = "servicestaff:create"
	PermServiceStaffRead   = "servicestaff:read"
	PermServiceStaffUpdate = "servicestaff:update"

	// Parts permissions
	PermPartsCreate = "parts:create"
	PermPartsRead   = "parts:read"
	PermPartsUpdate = "parts:update"
	PermPartsDelete = "parts:delete"

	// Inventory permissions
	PermInventoryCreate = "inventory:create"
	PermInventoryRead   = "inventory:read"
	PermInventoryUpdate = "inventory:update"

	// Advanced Inventory Operations (warehouse manager)
	PermInventoryAdjust   = "inventory:adjust"
	PermInventoryTransfer = "inventory:transfer"
	PermInventoryAudit    = "inventory:audit"

	// Device Inventory permissions (school device tracking)
	PermLocationRead     = "location:read"
	PermLocationWrite    = "location:write"
	PermAssignmentRead   = "assignment:read"
	PermAssignmentWrite  = "assignment:write"
	PermGroupRead        = "group:read"
	PermGroupWrite       = "group:write"
	PermDeviceInventory  = "device:inventory" // View school device inventory

	// Reporting permissions
	PermReportInventory = "report:inventory"

	// Project permissions
	PermProjectCreate = "project:create"
	PermProjectRead   = "project:read"
	PermProjectUpdate = "project:update"

	// Phase permissions
	PermPhaseCreate = "phase:create"
	PermPhaseRead   = "phase:read"
	PermPhaseUpdate = "phase:update"

	// Survey permissions
	PermSurveyCreate = "survey:create"
	PermSurveyRead   = "survey:read"
	PermSurveyUpdate = "survey:update"

	// BOQ permissions
	PermBOQCreate = "boq:create"
	PermBOQRead   = "boq:read"
	PermBOQUpdate = "boq:update"

	// Project Team permissions
	PermProjectTeamRead   = "project:team:read"
	PermProjectTeamUpdate = "project:team:update"

	// Program Activity permissions
	PermActivityCreate = "activity:create"
	PermActivityRead   = "activity:read"
	PermActivityUpdate = "activity:update"
	PermActivityDelete = "activity:delete"

	// User Notification permissions
	PermNotificationRead   = "notification:read"
	PermNotificationUpdate = "notification:update"

	// SSOT permissions
	PermSSOTRead    = "ssot:read"
	PermSSOTSync    = "ssot:sync"
	PermSSOTWebhook = "ssot:webhook"

	// Telemetry permissions
	PermTelemetryIngest = "telemetry:ingest"

	// Messaging permissions
	PermMessagesRead   = "messages:read"
	PermMessagesCreate = "messages:create"
	PermMessagesManage = "messages:manage"

	// Live Chat permissions
	PermChatAccept   = "chat:accept"
	PermChatTransfer = "chat:transfer"
	PermChatManage   = "chat:manage"

	// Device permissions
	PermDeviceRead   = "device:read"
	PermDeviceCreate = "device:create"
	PermDeviceUpdate = "device:update"

	// Knowledge Base permissions
	PermKBCreate = "kb:create"
	PermKBRead   = "kb:read"
	PermKBUpdate = "kb:update"
	PermKBDelete = "kb:delete"

	// Impersonation permission (ops managers can act on behalf of school contacts)
	PermImpersonate = "impersonate:user"

	// Marketing Knowledge Base permissions
	PermMKBCreate  = "mkb:create"
	PermMKBRead    = "mkb:read"
	PermMKBUpdate  = "mkb:update"
	PermMKBDelete  = "mkb:delete"
	PermMKBApprove = "mkb:approve"

	// Sales/Marketing permissions
	PermDemoManage     = "demo:manage"        // Manage demo pipeline
	PermDemoPipeline   = "demo:pipeline"      // View/update demo pipeline stages
	PermSchoolReadAll  = "school:read_all"    // Read all school data (not just own)
	PermReportingSales = "reporting:sales"    // Access sales reports
	PermContentManage  = "content:manage"     // Manage marketing content
	PermPresentations  = "presentations:view" // Access sales presentations
	PermReportsRead    = "reports:read"       // Read general reports
	PermDashboardRead  = "dashboard:read"     // Access dashboard

	// Operations Manager permissions (global field operations)
	PermOpsManageShops     = "ops:manage_shops"      // Create/update service shops globally
	PermOpsGlobalInventory = "ops:global_inventory"  // View/manage inventory across all shops
	PermOpsReassignWork    = "ops:reassign_work"     // Reassign work orders between shops
	PermOpsGlobalReports   = "ops:global_reports"    // Access global operations reports
	PermOpsManageStaff     = "ops:manage_staff"      // Manage staff across all shops
	PermOpsDashboard       = "ops:dashboard"         // Access operations manager dashboard

	// Wildcard for admin
	PermAll = "*"
)

// RolePermissions maps Keycloak roles to their permitted operations
// Based on roles defined in /home/pato/opt/ESSP/docs/rbac.md
var RolePermissions = map[string][]string{
	// Tenant super admin - all permissions
	"ssp_admin": {
		PermAll,
		// Explicitly list critical admin permissions for clarity
		PermSSOTRead,
		PermSSOTSync,
		PermSSOTWebhook,
	},

	// Operations Manager - global field operations lead (between Admin and Lead Tech)
	"ssp_ops_manager": {
		// All Lead Tech permissions
		PermWorkOrderRead,
		PermWorkOrderUpdate,
		PermWorkOrderSchedule,
		PermWorkOrderDeliverable,
		PermWorkOrderApproval,
		PermBOMRead,
		PermBOMUpdate,
		PermBOMConsume,
		PermAttachmentCreate,
		PermAttachmentRead,
		PermSchoolRead,
		PermSchoolContactRead,
		PermServiceStaffRead,
		PermServiceShopRead,
		PermDeviceRead,
		PermPartsRead,
		PermInventoryRead,
		PermTelemetryIngest,
		PermProjectTeamRead,
		PermProjectTeamUpdate,
		PermActivityCreate,
		PermActivityRead,
		PermActivityUpdate,
		PermNotificationRead,
		PermNotificationUpdate,
		PermMessagesRead,
		PermMessagesCreate,
		PermMessagesManage,
		PermKBRead,

		// Operations Manager specific - global/cross-shop capabilities
		PermOpsManageShops,     // Create/update service shops
		PermOpsGlobalInventory, // View/manage inventory across all shops
		PermOpsReassignWork,    // Reassign work orders between shops
		PermOpsGlobalReports,   // Access global operations reports
		PermOpsManageStaff,     // Manage staff across all shops
		PermOpsDashboard,       // Access operations dashboard

		// Additional management permissions
		PermServiceShopCreate, // Create new service shops
		PermServiceShopUpdate, // Update service shops
		PermServiceStaffCreate, // Create staff
		PermServiceStaffUpdate, // Update staff
		PermWorkOrderCreate,   // Create work orders
		PermInventoryUpdate,   // Update inventory
		PermReportsRead,       // View reports
		PermDashboardRead,     // View dashboard

		// Device inventory management
		PermDeviceInventory,
		PermLocationRead,
		PermLocationWrite,
		PermAssignmentRead,
		PermAssignmentWrite,
		PermGroupRead,
		PermGroupWrite,

		// Impersonation - can act on behalf of school contacts
		PermImpersonate,
	},

	// Support agent - tickets/dispatch
	"ssp_support_agent": {
		PermIncidentCreate,
		PermIncidentRead,
		PermIncidentUpdate,
		PermWorkOrderCreate,
		PermWorkOrderRead,
		PermWorkOrderUpdate,
		PermWorkOrderSchedule,
		PermWorkOrderReview,
		PermAttachmentCreate,
		PermAttachmentRead,
		PermSchoolRead,
		PermSchoolContactRead,
		PermServiceShopRead,
		PermServiceStaffRead,
		PermPartsRead,
		PermInventoryRead,
		PermDeviceRead,
		PermProjectTeamRead,
		PermActivityCreate,
		PermActivityRead,
		PermActivityUpdate,
		PermNotificationRead,
		PermNotificationUpdate,
		PermMessagesRead,
		PermMessagesCreate,
		PermMessagesManage,
		PermChatAccept,
		PermChatTransfer,
		PermKBRead,
	},

	// Field tech - work orders + deliverables (RESTRICTED - no project/activity access)
	"ssp_field_tech": {
		PermWorkOrderRead,
		PermWorkOrderUpdate,
		PermWorkOrderDeliverable,
		PermBOMRead,
		PermBOMConsume,
		PermAttachmentCreate,
		PermAttachmentRead,
		PermSchoolRead,
		PermSchoolContactRead,
		PermTelemetryIngest,
		PermNotificationRead,
		PermNotificationUpdate,
		PermMessagesRead,
		PermMessagesCreate,
		PermKBRead,
	},

	// Lead tech - scheduling + approval requests + team management
	"ssp_lead_tech": {
		PermWorkOrderRead,
		PermWorkOrderUpdate,
		PermWorkOrderSchedule,
		PermWorkOrderDeliverable,
		PermWorkOrderApproval,
		PermBOMRead,
		PermBOMUpdate,
		PermBOMConsume,
		PermAttachmentCreate,
		PermAttachmentRead,
		PermSchoolRead,
		PermSchoolContactRead,
		PermServiceStaffRead,
		PermServiceShopRead, // NEW: View service shop locations for dispatch
		PermDeviceRead,      // NEW: View device details in work orders
		PermPartsRead,
		PermInventoryRead,
		PermTelemetryIngest,
		PermProjectTeamRead,
		PermProjectTeamUpdate,
		PermActivityCreate,
		PermActivityRead,
		PermActivityUpdate,
		PermNotificationRead,
		PermNotificationUpdate,
		PermMessagesRead,
		PermMessagesCreate,
		PermMessagesManage,
		PermKBRead,
	},

	// Demo team - demos, surveys, pipeline
	"ssp_demo_team": {
		PermProjectCreate,
		PermProjectRead,
		PermProjectUpdate,
		PermPhaseCreate,
		PermPhaseRead,
		PermPhaseUpdate,
		PermSurveyCreate,
		PermSurveyRead,
		PermSurveyUpdate,
		PermBOQCreate,
		PermBOQRead,
		PermBOQUpdate,
		PermSchoolRead,
		PermAttachmentCreate,
		PermAttachmentRead,
		PermTelemetryIngest,
		PermProjectTeamRead,
		PermProjectTeamUpdate,
		PermActivityCreate,
		PermActivityRead,
		PermActivityUpdate,
		PermActivityDelete,
		PermNotificationRead,
		PermNotificationUpdate,
		PermMessagesRead,
		PermMessagesCreate,
	},

	// Sales/Marketing - demos, pipeline, customer data, reports, presentations
	"ssp_sales_marketing": {
		// School/Customer data (read-only across all schools)
		PermSchoolRead,
		PermSchoolReadAll,
		PermSchoolContactRead,
		PermDeviceRead,
		PermPartsRead,

		// Demo/Pipeline management
		PermDemoManage,
		PermDemoPipeline,
		PermProjectRead,
		PermPhaseRead,
		PermSurveyRead,
		PermBOQRead,

		// Work orders (view only for support context)
		PermWorkOrderRead,

		// Reports and analytics
		PermReportingSales,
		PermReportsRead,
		PermDashboardRead,

		// Content/Presentations
		PermContentManage,
		PermPresentations,

		// Marketing Knowledge Base (can create/read/update, but not delete/approve)
		PermMKBCreate,
		PermMKBRead,
		PermMKBUpdate,

		// Activity and team
		PermProjectTeamRead,
		PermActivityRead,

		// Communication
		PermNotificationRead,
		PermNotificationUpdate,
		PermMessagesRead,
		PermMessagesCreate,
	},

	// School contact - create incidents, approve sign-offs, manage device inventory
	"ssp_school_contact": {
		PermIncidentCreate,
		PermIncidentRead,
		PermWorkOrderRead,
		PermWorkOrderApproval,
		PermAttachmentCreate,
		PermAttachmentRead,
		PermSchoolRead,
		PermSchoolContactRead,
		PermProjectTeamRead,
		PermActivityCreate,
		PermActivityRead,
		PermNotificationRead,
		PermNotificationUpdate,
		// Device inventory management for their school
		PermDeviceInventory,
		PermLocationRead,
		PermLocationWrite,    // Create/update school locations
		PermAssignmentRead,
		PermAssignmentWrite,  // Assign/unassign devices
		PermGroupRead,
		PermGroupWrite,       // Create/manage device groups
		PermDeviceCreate,     // Register new devices
	},

	// Supplier - parts catalog + fulfillment visibility
	"ssp_supplier": {
		PermPartsRead,
		PermInventoryRead,
		PermWorkOrderRead,
		PermBOMRead,
		PermNotificationRead,
		PermNotificationUpdate,
		PermMessagesRead,
		PermMessagesCreate,
	},

	// Contractor - work packages, deliverables submission
	"ssp_contractor": {
		PermWorkOrderRead,
		PermWorkOrderDeliverable,
		PermBOMRead,
		PermAttachmentCreate,
		PermAttachmentRead,
		PermSchoolRead,
		PermProjectTeamRead,
		PermActivityCreate,
		PermActivityRead,
		PermNotificationRead,
		PermNotificationUpdate,
		PermMessagesRead,
		PermMessagesCreate,
	},

	// Warehouse manager - inventory, BOM pick/issue, deliverables confirmation
	"ssp_warehouse_manager": {
		// Inventory - full CRUD + advanced operations
		PermInventoryCreate,
		PermInventoryRead,
		PermInventoryUpdate,
		PermInventoryAdjust,
		PermInventoryTransfer,
		PermInventoryAudit,

		// Parts - full catalog management
		PermPartsCreate,
		PermPartsRead,
		PermPartsUpdate,
		PermPartsDelete,

		// BOM operations
		PermBOMRead,
		PermBOMUpdate,
		PermBOMConsume,

		// Work Order - read + deliverables
		PermWorkOrderRead,
		PermWorkOrderDeliverable,

		// Activity logging
		PermActivityCreate,
		PermActivityRead,

		// Supporting access
		PermServiceShopRead,
		PermProjectTeamRead,
		PermDeviceRead,

		// Reporting
		PermReportInventory,

		// Notifications and messages
		PermNotificationRead,
		PermNotificationUpdate,
		PermMessagesRead,
		PermMessagesCreate,
	},
}

// HasPermission checks if a given role has a specific permission
func HasPermission(role, permission string) bool {
	perms, ok := RolePermissions[role]
	if !ok {
		return false
	}

	for _, p := range perms {
		if p == PermAll || p == permission {
			return true
		}
	}
	return false
}

// HasAnyPermission checks if a given role has any of the specified permissions
func HasAnyPermission(role string, permissions ...string) bool {
	for _, perm := range permissions {
		if HasPermission(role, perm) {
			return true
		}
	}
	return false
}

// GetUserPermissions returns all permissions for a list of roles
func GetUserPermissions(roles []string) []string {
	permSet := make(map[string]bool)

	for _, role := range roles {
		if perms, ok := RolePermissions[role]; ok {
			for _, p := range perms {
				if p == PermAll {
					// If user has wildcard permission, return it immediately
					return []string{PermAll}
				}
				permSet[p] = true
			}
		}
	}

	permissions := make([]string, 0, len(permSet))
	for perm := range permSet {
		permissions = append(permissions, perm)
	}
	return permissions
}

// UserHasPermission checks if a user with given roles has a specific permission
func UserHasPermission(roles []string, permission string) bool {
	for _, role := range roles {
		if HasPermission(role, permission) {
			return true
		}
	}
	return false
}

// UserHasAnyPermission checks if a user with given roles has any of the specified permissions
func UserHasAnyPermission(roles []string, permissions ...string) bool {
	for _, perm := range permissions {
		if UserHasPermission(roles, perm) {
			return true
		}
	}
	return false
}
