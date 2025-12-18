// Role-to-permissions mapping
// This defines what permissions each role has
// Mirrors backend permissions from services/ims-api/internal/auth/permissions.go

export const rolePermissions: Record<string, string[]> = {
  // Admin has all permissions (wildcard)
  ssp_admin: [
    '*', // Wildcard - has all permissions
    'incident:read', 'incident:create', 'incident:update', 'incident:delete',
    'workorder:read', 'workorder:create', 'workorder:update', 'workorder:delete',
    'workorder:schedule', 'workorder:deliverable', 'workorder:approval', 'workorder:review',
    'project:read', 'project:create', 'project:update',
    'phase:read', 'phase:create', 'phase:update',
    'survey:read', 'survey:create', 'survey:update',
    'boq:read', 'boq:create', 'boq:update',
    'serviceshop:read', 'serviceshop:create', 'serviceshop:update',
    'servicestaff:read', 'servicestaff:create', 'servicestaff:update',
    'school:read', 'school:create', 'school:update',
    'school:contact:read', 'school:contact:create', 'school:contact:update',
    'device:read', 'device:create', 'device:update',
    'parts:read', 'parts:create', 'parts:update', 'parts:delete',
    'inventory:read', 'inventory:create', 'inventory:update',
    'bom:read', 'bom:create', 'bom:update', 'bom:consume',
    'attachment:read', 'attachment:create', 'attachment:delete',
    'project:team:read', 'project:team:update',
    'activity:read', 'activity:create', 'activity:update', 'activity:delete',
    'notification:read', 'notification:update',
    'ssot:read', 'ssot:sync', 'ssot:webhook',
    'telemetry:ingest',
    'messages:read', 'messages:create', 'messages:manage',
    'chat:accept', 'chat:transfer', 'chat:manage',
    'kb:read', 'kb:create', 'kb:update', 'kb:delete',
    'mkb:read', 'mkb:create', 'mkb:update', 'mkb:delete', 'mkb:approve',
    'audit:read',
    'settings:read', 'settings:update',
    // HR SSOT permissions
    'hr:read', 'hr:create', 'hr:update', 'hr:delete',
    'people:read', 'teams:read', 'org:read',
  ],

  // Operations Manager - global field operations lead (between Admin and Lead Tech)
  ssp_ops_manager: [
    // SSOT read access for school directory
    'ssot:read',
    // All Lead Tech permissions
    'workorder:read', 'workorder:create', 'workorder:update',
    'workorder:schedule', 'workorder:deliverable', 'workorder:approval',
    'bom:read', 'bom:update', 'bom:consume',
    'attachment:read', 'attachment:create',
    'school:read', 'school:contact:read',
    'servicestaff:read', 'servicestaff:create', 'servicestaff:update',
    'serviceshop:read', 'serviceshop:create', 'serviceshop:update',
    'device:read',
    'parts:read',
    'inventory:read', 'inventory:update',
    'telemetry:ingest',
    'project:team:read', 'project:team:update',
    'activity:read', 'activity:create', 'activity:update',
    'notification:read', 'notification:update',
    'messages:read', 'messages:create', 'messages:manage',
    'kb:read',

    // Operations Manager specific - global/cross-shop capabilities
    'ops:manage_shops',      // Create/update service shops globally
    'ops:global_inventory',  // View/manage inventory across all shops
    'ops:reassign_work',     // Reassign work orders between shops
    'ops:global_reports',    // Access global operations reports
    'ops:manage_staff',      // Manage staff across all shops
    'ops:dashboard',         // Access operations dashboard

    // Additional management permissions
    'reports:read',
    'dashboard:read',
    // HR permissions (read-only for ops manager)
    'hr:read', 'people:read', 'teams:read', 'org:read',
    // Impersonation - can act as school contacts
    'impersonate:user',
  ],

  // Support agent - tickets/dispatch
  ssp_support_agent: [
    'incident:read', 'incident:create', 'incident:update',
    'workorder:read', 'workorder:create', 'workorder:update',
    'workorder:schedule', 'workorder:review',
    'attachment:read', 'attachment:create',
    'school:read', 'school:contact:read',
    'serviceshop:read', 'servicestaff:read',
    'parts:read', 'inventory:read',
    'device:read',
    'project:team:read',
    'activity:read', 'activity:create', 'activity:update',
    'notification:read', 'notification:update',
    'messages:read', 'messages:create', 'messages:manage',
    'chat:accept', 'chat:transfer',
    'kb:read',
  ],

  // Field tech - work orders + deliverables (RESTRICTED - no project/activity access)
  ssp_field_tech: [
    'workorder:read', 'workorder:update', 'workorder:deliverable',
    'bom:read', 'bom:consume',
    'attachment:read', 'attachment:create',
    'school:read', 'school:contact:read',
    'telemetry:ingest',
    'notification:read', 'notification:update',
    'messages:read', 'messages:create',
    'kb:read',
  ],

  // Lead tech - scheduling + approval requests + team management
  ssp_lead_tech: [
    'workorder:read', 'workorder:update',
    'workorder:schedule', 'workorder:deliverable', 'workorder:approval',
    'bom:read', 'bom:update', 'bom:consume',
    'attachment:read', 'attachment:create',
    'school:read', 'school:contact:read',
    'servicestaff:read',
    'serviceshop:read', // NEW: View service shop locations for dispatch
    'device:read',      // NEW: View device details in work orders
    'parts:read', 'inventory:read',
    'telemetry:ingest',
    'project:team:read', 'project:team:update',
    'activity:read', 'activity:create', 'activity:update',
    'notification:read', 'notification:update',
    'messages:read', 'messages:create', 'messages:manage',
    'kb:read',
  ],

  // Demo team - demos, surveys, pipeline
  ssp_demo_team: [
    'project:read', 'project:create', 'project:update',
    'phase:read', 'phase:create', 'phase:update',
    'survey:read', 'survey:create', 'survey:update',
    'boq:read', 'boq:create', 'boq:update',
    'school:read',
    'attachment:read', 'attachment:create',
    'telemetry:ingest',
    'project:team:read', 'project:team:update',
    'activity:read', 'activity:create', 'activity:update', 'activity:delete',
    'notification:read', 'notification:update',
    'messages:read', 'messages:create',
  ],

  // Sales/Marketing - demos, pipeline, customer data, reports, presentations
  ssp_sales_marketing: [
    // School/Customer data (read-only across all schools)
    'school:read',
    'school:read_all',
    'school:contact:read',
    'device:read',
    'parts:read',

    // Demo/Pipeline management
    'demo:manage',
    'demo:pipeline',
    'project:read',
    'phase:read',
    'survey:read',
    'boq:read',

    // Work orders (view only for support context)
    'workorder:read',

    // Reports and analytics
    'reporting:sales',
    'reports:read',
    'dashboard:read',

    // Content/Presentations
    'content:manage',
    'presentations:view',

    // Marketing Knowledge Base (can create/read/update, but not delete/approve)
    'mkb:create', 'mkb:read', 'mkb:update',

    // Activity and team
    'project:team:read',
    'activity:read',

    // Communication
    'notification:read', 'notification:update',
    'messages:read', 'messages:create',
  ],

  // School contact - create incidents, approve sign-offs
  ssp_school_contact: [
    'incident:read', 'incident:create',
    'workorder:read', 'workorder:approval',
    'attachment:read', 'attachment:create',
    'school:read', 'school:contact:read',
    'project:team:read',
    'activity:read', 'activity:create',
    'notification:read', 'notification:update',
  ],

  // Supplier - parts catalog + fulfillment visibility
  ssp_supplier: [
    'parts:read', 'inventory:read',
    'workorder:read',
    'bom:read',
    'notification:read', 'notification:update',
    'messages:read', 'messages:create',
  ],

  // Contractor - work packages, deliverables submission
  ssp_contractor: [
    'workorder:read', 'workorder:deliverable',
    'bom:read',
    'attachment:read', 'attachment:create',
    'school:read',
    'project:team:read',
    'activity:read', 'activity:create',
    'notification:read', 'notification:update',
    'messages:read', 'messages:create',
  ],

  // Warehouse manager - inventory, BOM pick/issue, deliverables confirmation
  ssp_warehouse_manager: [
    // Inventory - full CRUD + advanced operations
    'inventory:read', 'inventory:create', 'inventory:update',
    'inventory:adjust', 'inventory:transfer', 'inventory:audit',

    // Parts - full catalog management
    'parts:read', 'parts:create', 'parts:update', 'parts:delete',

    // BOM operations
    'bom:read', 'bom:update', 'bom:consume',

    // Work order - read + deliverables
    'workorder:read', 'workorder:deliverable',

    // Activity logging
    'activity:create', 'activity:read',

    // Supporting access
    'serviceshop:read', 'project:team:read', 'device:read',

    // Reporting
    'report:inventory',

    // Notifications and messages
    'notification:read', 'notification:update',
    'messages:read', 'messages:create',
  ],
};

// Get all permissions for a set of roles
export function getPermissionsForRoles(roles: string[]): string[] {
  const permissions = new Set<string>();

  for (const role of roles) {
    const rolePerms = rolePermissions[role];
    if (rolePerms) {
      for (const perm of rolePerms) {
        permissions.add(perm);
      }
    }
  }

  return Array.from(permissions);
}

// Check if a role has a specific permission
export function roleHasPermission(role: string, permission: string): boolean {
  const perms = rolePermissions[role];
  return perms ? perms.includes(permission) : false;
}
