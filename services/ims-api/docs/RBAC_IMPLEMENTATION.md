# RBAC Implementation for IMS API

This document describes the Role-Based Access Control (RBAC) implementation for the ESSP IMS API service (FE-001).

## Overview

The RBAC system provides fine-grained access control for all API endpoints based on JWT token claims. It extracts user roles from JWT tokens and enforces permission checks before allowing access to protected resources.

## Implementation Components

### 1. Permission Definitions (`internal/auth/permissions.go`)

Defines all permission constants and role-to-permission mappings.

**Permission Categories:**
- Incident operations: `incident:create`, `incident:read`, `incident:update`, `incident:delete`
- Work Order operations: `workorder:create`, `workorder:read`, `workorder:update`, `workorder:delete`
- Work Order Operations: `workorder:schedule`, `workorder:deliverable`, `workorder:approval`, `workorder:review`
- BOM operations: `bom:create`, `bom:read`, `bom:update`, `bom:consume`
- School operations: `school:create`, `school:read`, `school:update`
- School Contacts: `school:contact:create`, `school:contact:read`, `school:contact:update`
- Attachments: `attachment:create`, `attachment:read`, `attachment:delete`
- Service Shops: `serviceshop:create`, `serviceshop:read`, `serviceshop:update`
- Service Staff: `servicestaff:create`, `servicestaff:read`, `servicestaff:update`
- Parts: `parts:create`, `parts:read`, `parts:update`
- Inventory: `inventory:create`, `inventory:read`, `inventory:update`
- Programs: `program:create`, `program:read`, `program:update`
- Phases: `phase:create`, `phase:read`, `phase:update`
- Surveys: `survey:create`, `survey:read`, `survey:update`
- BOQ: `boq:create`, `boq:read`, `boq:update`
- SSOT: `ssot:sync`, `ssot:webhook`
- Telemetry: `telemetry:ingest`

**Role Mappings:**
- `ssp_admin`: All permissions (wildcard `*`)
- `ssp_support_agent`: Incident and work order management
- `ssp_field_tech`: Work order execution and deliverables
- `ssp_lead_tech`: Scheduling, approvals, and BOM management
- `ssp_demo_team`: Programs, surveys, and BOQ management
- `ssp_sales`: Read-only access to programs and sales pipeline
- `ssp_school_contact`: Create incidents and approve sign-offs
- `ssp_supplier`: Read-only access to parts and inventory
- `ssp_contractor`: Work order deliverables submission
- `ssp_warehouse_manager`: Inventory and BOM operations

### 2. RBAC Middleware (`internal/middleware/rbac.go`)

Provides middleware functions for enforcing permission checks:

**Functions:**
- `RequirePermission(permission, logger)`: Checks if user has a specific permission
- `RequireAnyPermission(logger, permissions...)`: Checks if user has any of the specified permissions
- `RequireRole(role, logger)`: Checks if user has a specific role
- `RequireAnyRole(logger, roles...)`: Checks if user has any of the specified roles
- `RequireSchoolAccess(logger)`: Validates school-scoped access

**Features:**
- Returns HTTP 403 Forbidden for unauthorized access
- Logs permission denials with context (roles, permissions, path)
- Supports admin bypass for school-scoped endpoints
- Graceful handling of missing roles

### 3. Context Enhancements (`internal/middleware/context.go`)

Extended to store user authentication details:

**New Context Values:**
- `roles`: User's roles extracted from JWT
- `assignedSchools`: User's assigned school IDs
- `claims`: Full JWT claims map

**Functions:**
- `WithRoles(ctx, roles)` / `Roles(ctx)`
- `WithAssignedSchools(ctx, schools)` / `AssignedSchools(ctx)`
- `WithClaims(ctx, claims)` / `Claims(ctx)`

### 4. JWT Authentication Updates (`internal/middleware/auth.go`)

Enhanced to extract and store roles and school assignments:

**JWT Claim Extraction:**
- Supports multiple JWT formats:
  - Direct `roles` array
  - Keycloak `realm_access.roles`
  - Keycloak `resource_access` format
- Filters for SSP-specific roles (prefix: `ssp_`)
- Extracts school assignments from `schools` or `schoolId` claims
- Stores all extracted data in request context

**Logging:**
- Debug logs authentication details (roles, schools, path)

### 5. Server Route Protection (`internal/api/server.go`)

All API endpoints are now protected with appropriate RBAC middleware:

**Route Groups by Permission:**

1. **Incidents**
   - Read: `incident:read`
   - Create: `incident:create`
   - Update: `incident:update`

2. **Work Orders**
   - Read: `workorder:read`
   - Create: `workorder:create`
   - Update: `workorder:update`
   - Schedule: `workorder:schedule`
   - Deliverables: `workorder:deliverable`
   - Review: `workorder:review`
   - Approvals: `workorder:approval`

3. **BOM**
   - Read: `bom:read`
   - Create/Update: `bom:create` OR `bom:update`
   - Consume: `bom:consume`

4. **School Contacts**
   - Read: `school:contact:read`
   - Create/Update: `school:contact:create` OR `school:contact:update`

5. **Attachments**
   - Read: `attachment:read`
   - Create: `attachment:create`

6. **Telemetry**
   - Ingest: `telemetry:ingest`

7. **Schools**
   - Upsert: `school:create` OR `school:update`

8. **SSOT Operations** (Admin Only)
   - Sync: `ssot:sync`
   - Webhooks: `ssot:webhook`

9. **Programs, Phases, Surveys, BOQ**
   - Separate read/create/update permissions

10. **Service Shops, Staff, Parts, Inventory**
    - Separate read/create permissions

11. **Audit Logs** (Admin Only)
    - Requires `ssp_admin` role

**Public Endpoints:**
- `/healthz`: Health check (no auth required)
- `/readyz`: Readiness check (no auth required)

## JWT Token Format

The system expects JWT tokens with the following structure:

```json
{
  "iss": "https://your-keycloak.example.com/realms/ssp",
  "aud": "ims-api",
  "tenantId": "tenant-123",
  "schoolId": "school-456",  // Optional, for backward compatibility
  "schools": ["school-456", "school-789"],  // Optional, for multi-school access
  "realm_access": {
    "roles": [
      "ssp_admin",
      "ssp_field_tech"
    ]
  }
}
```

**Supported Claim Formats:**
1. `roles` - Direct array of role strings
2. `realm_access.roles` - Keycloak realm roles
3. `resource_access.{client}.roles` - Keycloak client-specific roles

**School Access:**
- `schools` - Array of school IDs the user can access
- `schoolId` - Single school ID (backward compatibility)

## Usage Examples

### Example 1: Protect a route with a specific permission

```go
r.Group(func(r chi.Router) {
    r.Use(middleware.RequirePermission(auth.PermIncidentRead, logger))
    r.Get("/incidents", handler.List)
})
```

### Example 2: Require any of multiple permissions

```go
r.Group(func(r chi.Router) {
    r.Use(middleware.RequireAnyPermission(logger,
        auth.PermBOMCreate,
        auth.PermBOMUpdate))
    r.Post("/work-orders/{id}/bom/items", handler.AddItem)
})
```

### Example 3: Require a specific role

```go
r.Group(func(r chi.Router) {
    r.Use(middleware.RequireRole("ssp_admin", logger))
    r.Get("/audit-logs", handler.List)
})
```

### Example 4: School-scoped access control

```go
r.Group(func(r chi.Router) {
    r.Use(middleware.RequirePermission(auth.PermSchoolContactRead, logger))
    r.Use(middleware.RequireSchoolAccess(logger))
    r.Get("/schools/{schoolId}/contacts", handler.List)
})
```

## Error Responses

**403 Forbidden** - Returned when:
- User has no roles assigned
- User lacks required permission(s)
- User lacks required role(s)
- User doesn't have access to requested school

Error messages:
- `"forbidden: no roles assigned"`
- `"forbidden: insufficient permissions"`
- `"forbidden: required role not assigned"`
- `"forbidden: no school access"`
- `"forbidden: school access denied"`

**401 Unauthorized** - Returned when:
- Authorization header is missing
- Authorization header is malformed
- JWT token is invalid or expired

## Testing

To test RBAC enforcement:

1. **Generate test JWT tokens** with different roles
2. **Call protected endpoints** with different tokens
3. **Verify 403 responses** for unauthorized access
4. **Check logs** for permission denial details

Example curl command:
```bash
curl -X GET http://localhost:8080/v1/incidents \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Security Considerations

1. **Role Hierarchy**: Admin role (`ssp_admin`) has wildcard permission (`*`)
2. **School Scoping**: Some roles can only access assigned schools
3. **Permission Logging**: All permission denials are logged for audit
4. **Token Validation**: JWT signature and claims are validated before role extraction
5. **Defense in Depth**: RBAC is applied at middleware layer, separate from business logic

## Future Enhancements

Potential improvements:
1. **Dynamic Permissions**: Load permissions from database
2. **Resource-Level Permissions**: Check ownership/assignment before allowing access
3. **Permission Caching**: Cache role-to-permission mappings in Redis
4. **Policy Engine**: Integrate with external policy engine (e.g., OPA)
5. **Fine-grained School Access**: Support school-level permission overrides

## Files Modified/Created

**Created:**
- `/home/pato/opt/ESSP/services/ims-api/internal/auth/permissions.go`
- `/home/pato/opt/ESSP/services/ims-api/internal/middleware/rbac.go`

**Modified:**
- `/home/pato/opt/ESSP/services/ims-api/internal/middleware/context.go`
- `/home/pato/opt/ESSP/services/ims-api/internal/middleware/auth.go`
- `/home/pato/opt/ESSP/services/ims-api/internal/api/server.go`

## Reference Documentation

- RBAC roles and permissions: `/home/pato/opt/ESSP/docs/rbac.md`
- JWT verification: `/home/pato/opt/ESSP/services/ims-api/internal/auth/verifier.go`
