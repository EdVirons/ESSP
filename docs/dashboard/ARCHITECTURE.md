# ESSP Management Dashboard - Architecture Design Document

## Document Information

- **Version**: 1.0.0
- **Date**: 2025-12-12
- **Status**: Draft
- **Author**: System Architect

---

## 1. Executive Summary

This document describes the architecture for the ESSP (EdVirons School Services Platform) Management Dashboard - a comprehensive web-based administration interface for managing the microservices platform. The dashboard provides operations teams with visibility into system health, data management capabilities, and operational controls.

### 1.1 Goals

1. **Unified Management Interface**: Single pane of glass for all ESSP services
2. **Operational Visibility**: Real-time health monitoring and alerting
3. **Data Management**: Full CRUD operations for core business entities
4. **Audit & Compliance**: Complete audit trail visibility
5. **Configuration Management**: Runtime configuration and feature flags

### 1.2 Scope

The dashboard covers:
- 5 microservices (ims-api, ssot-school, ssot-devices, ssot-parts, sync-worker)
- Infrastructure components (PostgreSQL, Valkey, NATS, MinIO)
- Business entities (Incidents, Work Orders, Programs, Service Shops, etc.)
- Operational functions (logs, jobs, events, audits)

---

## 2. Current System Analysis

### 2.1 Existing Architecture

```
+---------------------+    +------------------+    +------------------+
|   Frontend Clients  |    |   Mobile Apps    |    |   Admin Tools    |
+----------+----------+    +--------+---------+    +--------+---------+
           |                        |                       |
           +------------------------+-----------------------+
                                    |
                          +---------v---------+
                          |   Load Balancer   |
                          |   (NGINX/Traefik) |
                          +---------+---------+
                                    |
        +---------------------------+---------------------------+
        |                           |                           |
+-------v-------+           +-------v-------+           +-------v-------+
|   IMS-API     |           |  SSOT-School  |           |  SSOT-Devices |
|   :8080       |           |   :8081       |           |   :8082       |
+-------+-------+           +-------+-------+           +-------+-------+
        |                           |                           |
        +---------------------------+---------------------------+
                                    |
+-----------------------------------+-----------------------------------+
|                           Shared Infrastructure                       |
+-------+---------------+---------------+---------------+---------------+
        |               |               |               |
+-------v-------+ +-----v-----+ +-------v-------+ +-----v-----+
|  PostgreSQL   | |   Valkey  | |     NATS      | |   MinIO   |
|     :5432     | |   :6379   | |    :4222      | |   :9000   |
+---------------+ +-----------+ +---------------+ +-----------+
```

### 2.2 API Summary

Based on OpenAPI spec analysis (`/home/pato/opt/ESSP/docs/openapi/ims-api.yaml`):

| Endpoint Group | Operations | Authentication | Notes |
|----------------|------------|----------------|-------|
| Health | 2 | None | /healthz, /readyz |
| Incidents | 4 | JWT + RBAC | CRUD + status updates |
| Work Orders | 12 | JWT + RBAC | Full lifecycle |
| BOM | 5 | JWT + RBAC | Parts management |
| Service Shops | 3 | JWT + RBAC | Shop registry |
| Service Staff | 3 | JWT + RBAC | Staff registry |
| Programs | 3 | JWT + RBAC | Program lifecycle |
| Phases | 3 | JWT + RBAC | Phase management |
| Surveys | 5 | JWT + RBAC | Site surveys |
| Attachments | 5 | JWT + RBAC | File management |
| Audit Logs | 2 | JWT + Admin | Admin only |

### 2.3 Authentication & Authorization

- **Method**: JWT Bearer tokens (Keycloak compatible)
- **Headers**: `X-Tenant-ID`, `X-School-ID` for multi-tenancy
- **RBAC Roles**: 10 defined roles with granular permissions
- **Key Roles for Dashboard**:
  - `ssp_admin`: Full access (dashboard primary user)
  - `ssp_support_agent`: Tickets/dispatch
  - `ssp_warehouse_manager`: Inventory/BOM

---

## 3. Dashboard Architecture

### 3.1 High-Level Architecture

```
+------------------------------------------------------------------+
|                        ESSP Dashboard                             |
+------------------------------------------------------------------+
|  +------------------+  +------------------+  +------------------+ |
|  |   Overview       |  |   Data Mgmt      |  |   Operations     | |
|  +------------------+  +------------------+  +------------------+ |
|  | - Service Health |  | - Incidents      |  | - Audit Logs     | |
|  | - Key Metrics    |  | - Work Orders    |  | - Event Stream   | |
|  | - Activity Feed  |  | - Programs       |  | - Job Monitor    | |
|  | - Alerts         |  | - Service Shops  |  | - Config Mgmt    | |
|  +------------------+  +------------------+  +------------------+ |
+------------------------------------------------------------------+
                              |
                              | API Calls (REST/JSON)
                              |
+------------------------------------------------------------------+
|                      Dashboard Backend                            |
|  +------------------+  +------------------+  +------------------+ |
|  | API Gateway      |  | Auth Middleware  |  | WebSocket Hub    | |
|  +------------------+  +------------------+  +------------------+ |
|  +------------------+  +------------------+  +------------------+ |
|  | Service Proxy    |  | Metrics Agg      |  | Log Aggregator   | |
|  +------------------+  +------------------+  +------------------+ |
+------------------------------------------------------------------+
                              |
        +---------------------+---------------------+
        |                     |                     |
+-------v-------+     +-------v-------+     +-------v-------+
|   IMS-API     |     |  SSOT APIs    |     |  Prometheus   |
|   (primary)   |     |  (school,     |     |   Grafana     |
|               |     |   devices,    |     |               |
|               |     |   parts)      |     |               |
+---------------+     +---------------+     +---------------+
```

### 3.2 Deployment Options

#### Option A: Embedded Dashboard (Recommended)

Dashboard served directly from `ims-api` as static files with a dedicated `/admin` route group.

**Pros**:
- Single deployment unit
- Shared authentication infrastructure
- No CORS complexity
- Simplified operations

**Cons**:
- Larger binary size
- Coupled release cycle
- Build complexity

#### Option B: Separate SPA with BFF

Standalone frontend application with a dedicated Backend-for-Frontend (BFF) service.

**Pros**:
- Independent deployment
- Specialized team ownership
- Technology flexibility

**Cons**:
- Additional infrastructure
- Authentication complexity
- CORS configuration

#### Option C: Server-Side Rendered (SSR)

Go templates with HTMX for interactivity.

**Pros**:
- Minimal JavaScript
- Fast initial load
- SEO friendly (not relevant for admin)

**Cons**:
- Limited interactivity
- More server-side logic
- Non-standard for SPAs

### 3.3 Recommended Approach

**Hybrid Embedded SPA (Option A+)**:
- React SPA embedded in `ims-api` binary using `embed.FS`
- Served at `/admin/*` routes
- Shares JWT authentication with existing API
- Uses existing middleware stack

---

## 4. Technology Stack

### 4.1 Frontend

| Category | Choice | Justification |
|----------|--------|---------------|
| **Framework** | React 18+ | Industry standard, large ecosystem, team familiarity |
| **Build Tool** | Vite | Fast HMR, modern bundling, TypeScript support |
| **Type System** | TypeScript | Type safety, better DX, refactoring support |
| **UI Library** | shadcn/ui | Tailwind-based, accessible, copy-paste components |
| **CSS** | Tailwind CSS | Utility-first, consistent design, small bundle |
| **State Management** | TanStack Query | Server state caching, mutations, real-time |
| **Routing** | React Router v6 | Standard routing, nested layouts |
| **Forms** | React Hook Form + Zod | Performance, validation, type inference |
| **Tables** | TanStack Table | Headless, virtualization, sorting/filtering |
| **Charts** | Recharts | React-native, responsive, composable |

### 4.2 Backend Extensions

| Category | Choice | Justification |
|----------|--------|---------------|
| **Admin Routes** | Chi router (existing) | Consistent with current codebase |
| **Static Serving** | `embed.FS` | Single binary deployment |
| **WebSocket** | nhooyr.io/websocket | Modern, RFC compliant, context support |
| **Metrics Aggregation** | Existing Prometheus | Leverage monitoring stack |

### 4.3 Development Tools

- **Package Manager**: pnpm (faster, disk efficient)
- **Linting**: ESLint + Prettier
- **Testing**: Vitest + React Testing Library
- **E2E Testing**: Playwright

---

## 5. Component Architecture

### 5.1 Dashboard Layout Structure

```
+------------------------------------------------------------------+
|  [Logo]  ESSP Dashboard                    [User] [Settings] [?] |
+------------------------------------------------------------------+
|         |                                                         |
| +-----+ |  +------------------------------------------------------+
| | Nav | |  | Breadcrumb: Overview > Incidents > INC-2024-001      |
| |     | |  +------------------------------------------------------+
| | [x] | |  |                                                      |
| | Ovr | |  |                    MAIN CONTENT AREA                 |
| | [_] | |  |                                                      |
| | Inc | |  |  +-------------------+  +-------------------+        |
| | [_] | |  |  |   Summary Card    |  |   Summary Card    |        |
| | WOs | |  |  +-------------------+  +-------------------+        |
| | [_] | |  |                                                      |
| | Prg | |  |  +-----------------------------------------------+  |
| | [_] | |  |  |                                               |  |
| | Shp | |  |  |           Data Table / Detail View            |  |
| | [_] | |  |  |                                               |  |
| | Ops | |  |  +-----------------------------------------------+  |
| |     | |  |                                                      |
| +-----+ |  +------------------------------------------------------+
+------------------------------------------------------------------+
|  Status: Connected | Last sync: 2s ago | v1.0.0                   |
+------------------------------------------------------------------+
```

### 5.2 Navigation Structure

```
Dashboard
+-- Overview
|   +-- Service Health
|   +-- Key Metrics
|   +-- Activity Feed
|   +-- System Alerts
|
+-- Data Management
|   +-- Incidents
|   |   +-- List
|   |   +-- Create
|   |   +-- [id] Detail
|   |
|   +-- Work Orders
|   |   +-- List
|   |   +-- Create
|   |   +-- [id] Detail
|   |   +-- [id] BOM
|   |   +-- [id] Deliverables
|   |
|   +-- Programs
|   |   +-- List
|   |   +-- Create
|   |   +-- [id] Detail
|   |   +-- [id] Phases
|   |   +-- [id] Surveys
|   |
|   +-- Service Shops
|   |   +-- List
|   |   +-- Create
|   |   +-- [id] Detail
|   |   +-- [id] Staff
|   |
|   +-- SSOT Registry
|       +-- Schools
|       +-- Devices
|       +-- Parts
|
+-- Operations
|   +-- Audit Logs
|   +-- Event Stream
|   +-- Background Jobs
|   +-- NATS Monitor
|
+-- Configuration
    +-- Environment
    +-- Feature Flags
    +-- Rate Limits
    +-- RBAC Roles
```

### 5.3 Component Hierarchy

```
App
+-- ThemeProvider
+-- AuthProvider
+-- QueryClientProvider
+-- RouterProvider
    +-- Layout
        +-- Header
        |   +-- Logo
        |   +-- GlobalSearch
        |   +-- NotificationBell
        |   +-- UserMenu
        |
        +-- Sidebar
        |   +-- NavGroup
        |   |   +-- NavItem
        |   +-- CollapseToggle
        |
        +-- MainContent
        |   +-- Breadcrumb
        |   +-- PageHeader
        |   +-- <Outlet /> (Page Content)
        |
        +-- StatusBar
            +-- ConnectionStatus
            +-- SyncIndicator
            +-- VersionInfo
```

---

## 6. Data Flow Architecture

### 6.1 API Communication

```typescript
// API Client Pattern
const apiClient = createApiClient({
  baseURL: '/api/v1',
  headers: {
    'X-Tenant-ID': tenantId,
    'X-School-ID': schoolId,
  },
});

// Query Hook Pattern
function useIncidents(filters: IncidentFilters) {
  return useQuery({
    queryKey: ['incidents', filters],
    queryFn: () => apiClient.get('/incidents', { params: filters }),
    staleTime: 30_000, // 30 seconds
  });
}

// Mutation Hook Pattern
function useCreateIncident() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: CreateIncidentRequest) =>
      apiClient.post('/incidents', data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['incidents'] });
    },
  });
}
```

### 6.2 Real-Time Updates

```
Browser                    Dashboard Backend              NATS
   |                              |                        |
   | WebSocket Connect            |                        |
   |----------------------------->|                        |
   |                              |                        |
   |                              | Subscribe to events    |
   |                              |----------------------->|
   |                              |                        |
   |                              |     Event: incident.created
   |                              |<-----------------------|
   |    JSON Event                |                        |
   |<-----------------------------|                        |
   |                              |                        |
   | Update UI (TanStack Query)   |                        |
   |                              |                        |
```

### 6.3 State Management Strategy

| State Type | Solution | Example |
|------------|----------|---------|
| Server State | TanStack Query | Incidents list, Work Order details |
| Form State | React Hook Form | Create/Edit forms |
| UI State | React useState/useReducer | Modal open, sidebar collapsed |
| Auth State | Context + localStorage | JWT token, user info |
| URL State | React Router | Filters, pagination, sorting |

---

## 7. Security Architecture

### 7.1 Authentication Flow

```
+----------+          +----------+          +----------+
|  Browser |          | Dashboard|          | Keycloak |
+----+-----+          +----+-----+          +----+-----+
     |                     |                     |
     | 1. Access /admin    |                     |
     |-------------------->|                     |
     |                     |                     |
     | 2. Redirect to login|                     |
     |<--------------------|                     |
     |                     |                     |
     | 3. Login with creds |                     |
     |-------------------------------------------->
     |                     |                     |
     | 4. JWT tokens       |                     |
     |<--------------------------------------------
     |                     |                     |
     | 5. Store tokens     |                     |
     |                     |                     |
     | 6. API request + JWT|                     |
     |-------------------->|                     |
     |                     | 7. Verify JWT       |
     |                     |-------------------->|
     |                     |                     |
     |                     | 8. Token valid      |
     |                     |<--------------------|
     |                     |                     |
     | 9. API response     |                     |
     |<--------------------|                     |
```

### 7.2 Authorization Model

```typescript
// Permission check in frontend
function canAccessPage(user: User, requiredPermissions: string[]): boolean {
  if (user.roles.includes('ssp_admin')) return true;
  return requiredPermissions.every(perm =>
    user.permissions.includes(perm)
  );
}

// Route protection
<Route
  path="/incidents"
  element={
    <RequirePermission permissions={['incident:read']}>
      <IncidentsPage />
    </RequirePermission>
  }
/>
```

### 7.3 Security Headers

The existing middleware provides:
- Content-Security-Policy
- X-Frame-Options: DENY
- X-Content-Type-Options: nosniff
- X-XSS-Protection: 1; mode=block
- Referrer-Policy: strict-origin-when-cross-origin

---

## 8. API Extensions Required

### 8.1 Dashboard-Specific Endpoints

New endpoints to add to `ims-api`:

```yaml
# Admin dashboard endpoints
/admin/v1/health/services:
  get:
    summary: Aggregated health status of all services
    responses:
      200:
        content:
          application/json:
            schema:
              type: object
              properties:
                services:
                  type: array
                  items:
                    type: object
                    properties:
                      name: { type: string }
                      status: { type: string, enum: [healthy, degraded, unhealthy] }
                      latencyMs: { type: integer }
                      lastCheck: { type: string, format: date-time }

/admin/v1/metrics/summary:
  get:
    summary: Dashboard summary metrics
    responses:
      200:
        content:
          application/json:
            schema:
              type: object
              properties:
                incidents:
                  type: object
                  properties:
                    total: { type: integer }
                    open: { type: integer }
                    slaBreached: { type: integer }
                workOrders:
                  type: object
                  properties:
                    total: { type: integer }
                    inProgress: { type: integer }
                    completedToday: { type: integer }
                programs:
                  type: object
                  properties:
                    active: { type: integer }
                    pending: { type: integer }

/admin/v1/activity:
  get:
    summary: Recent activity feed
    parameters:
      - name: limit
        in: query
        schema: { type: integer, default: 50 }
    responses:
      200:
        content:
          application/json:
            schema:
              type: array
              items:
                type: object
                properties:
                  id: { type: string }
                  type: { type: string }
                  action: { type: string }
                  actor: { type: string }
                  target: { type: string }
                  timestamp: { type: string, format: date-time }
                  metadata: { type: object }

/admin/v1/ws:
  get:
    summary: WebSocket endpoint for real-time updates
    description: |
      Establishes WebSocket connection for:
      - Health status changes
      - New incidents/work orders
      - Audit log events
      - Job completions
```

### 8.2 SSOT API Extensions

Lightweight proxy endpoints to aggregate SSOT data:

```yaml
/admin/v1/ssot/schools:
  get:
    summary: List schools from SSOT cache

/admin/v1/ssot/devices:
  get:
    summary: List devices from SSOT cache

/admin/v1/ssot/parts:
  get:
    summary: List parts from SSOT cache
```

---

## 9. File Structure

### 9.1 Frontend Structure

```
/home/pato/opt/ESSP/dashboard/
+-- package.json
+-- pnpm-lock.yaml
+-- tsconfig.json
+-- vite.config.ts
+-- tailwind.config.ts
+-- postcss.config.js
+-- .eslintrc.cjs
+-- .prettierrc
+-- index.html
+-- public/
|   +-- favicon.ico
|   +-- logo.svg
|
+-- src/
    +-- main.tsx
    +-- App.tsx
    +-- index.css
    |
    +-- api/
    |   +-- client.ts           # API client setup
    |   +-- incidents.ts        # Incident API hooks
    |   +-- work-orders.ts      # Work order API hooks
    |   +-- programs.ts         # Program API hooks
    |   +-- service-shops.ts    # Service shop API hooks
    |   +-- audit-logs.ts       # Audit log API hooks
    |   +-- health.ts           # Health check API hooks
    |   +-- types.ts            # API type definitions
    |
    +-- components/
    |   +-- ui/                 # shadcn/ui components
    |   |   +-- button.tsx
    |   |   +-- card.tsx
    |   |   +-- dialog.tsx
    |   |   +-- table.tsx
    |   |   +-- ...
    |   |
    |   +-- layout/
    |   |   +-- Layout.tsx
    |   |   +-- Header.tsx
    |   |   +-- Sidebar.tsx
    |   |   +-- StatusBar.tsx
    |   |   +-- Breadcrumb.tsx
    |   |
    |   +-- shared/
    |   |   +-- DataTable.tsx
    |   |   +-- StatusBadge.tsx
    |   |   +-- LoadingSpinner.tsx
    |   |   +-- ErrorBoundary.tsx
    |   |   +-- ConfirmDialog.tsx
    |   |
    |   +-- incidents/
    |   |   +-- IncidentList.tsx
    |   |   +-- IncidentForm.tsx
    |   |   +-- IncidentDetail.tsx
    |   |   +-- IncidentStatusSelect.tsx
    |   |
    |   +-- work-orders/
    |   |   +-- WorkOrderList.tsx
    |   |   +-- WorkOrderForm.tsx
    |   |   +-- WorkOrderDetail.tsx
    |   |   +-- WorkOrderBOM.tsx
    |   |   +-- WorkOrderDeliverables.tsx
    |   |
    |   +-- programs/
    |   |   +-- ProgramList.tsx
    |   |   +-- ProgramForm.tsx
    |   |   +-- ProgramDetail.tsx
    |   |   +-- PhaseList.tsx
    |   |   +-- SurveyList.tsx
    |   |
    |   +-- service-shops/
    |   |   +-- ServiceShopList.tsx
    |   |   +-- ServiceShopForm.tsx
    |   |   +-- ServiceShopDetail.tsx
    |   |   +-- StaffList.tsx
    |   |
    |   +-- operations/
    |   |   +-- AuditLogViewer.tsx
    |   |   +-- EventStream.tsx
    |   |   +-- JobMonitor.tsx
    |   |
    |   +-- overview/
    |       +-- ServiceHealthCard.tsx
    |       +-- MetricsSummary.tsx
    |       +-- ActivityFeed.tsx
    |       +-- AlertsPanel.tsx
    |
    +-- pages/
    |   +-- Overview.tsx
    |   +-- Incidents.tsx
    |   +-- IncidentDetail.tsx
    |   +-- WorkOrders.tsx
    |   +-- WorkOrderDetail.tsx
    |   +-- Programs.tsx
    |   +-- ProgramDetail.tsx
    |   +-- ServiceShops.tsx
    |   +-- ServiceShopDetail.tsx
    |   +-- AuditLogs.tsx
    |   +-- Settings.tsx
    |   +-- Login.tsx
    |   +-- NotFound.tsx
    |
    +-- hooks/
    |   +-- useAuth.ts
    |   +-- useWebSocket.ts
    |   +-- useDebounce.ts
    |   +-- usePagination.ts
    |
    +-- lib/
    |   +-- utils.ts
    |   +-- constants.ts
    |   +-- permissions.ts
    |
    +-- types/
    |   +-- index.ts
    |   +-- incidents.ts
    |   +-- work-orders.ts
    |   +-- programs.ts
    |
    +-- context/
        +-- AuthContext.tsx
        +-- ThemeContext.tsx
```

### 9.2 Backend Extensions

```
/home/pato/opt/ESSP/services/ims-api/
+-- cmd/
|   +-- api/
|       +-- main.go             # Updated to serve dashboard
|
+-- internal/
    +-- admin/                  # NEW: Admin dashboard handlers
    |   +-- routes.go           # Admin route registration
    |   +-- health.go           # Service health aggregation
    |   +-- metrics.go          # Dashboard metrics
    |   +-- activity.go         # Activity feed
    |   +-- websocket.go        # WebSocket hub
    |
    +-- dashboard/              # NEW: Embedded dashboard
        +-- embed.go            # embed.FS for static files
```

---

## 10. Implementation Phases

### Phase 1: Foundation (Week 1-2)

**Deliverables**:
1. Dashboard project scaffolding with Vite + React + TypeScript
2. UI component library setup (shadcn/ui + Tailwind)
3. Layout components (Header, Sidebar, Main)
4. API client with authentication
5. Backend: Admin routes group and static file serving

**Tasks**:
- [ ] Initialize dashboard project with Vite
- [ ] Configure TypeScript and ESLint
- [ ] Set up Tailwind CSS
- [ ] Install and configure shadcn/ui
- [ ] Create base Layout component
- [ ] Implement authentication flow
- [ ] Add admin routes to ims-api
- [ ] Embed static files in Go binary

### Phase 2: Overview Dashboard (Week 3)

**Deliverables**:
1. Service health monitoring panel
2. Key metrics summary cards
3. Activity feed component
4. System alerts panel

**Tasks**:
- [ ] Create health check aggregation endpoint
- [ ] Build ServiceHealthCard component
- [ ] Implement MetricsSummary with Recharts
- [ ] Create ActivityFeed with virtual scrolling
- [ ] Add real-time updates via WebSocket

### Phase 3: Incidents Management (Week 4)

**Deliverables**:
1. Incidents list with filtering/sorting
2. Incident creation form
3. Incident detail view
4. Status transition controls

**Tasks**:
- [ ] Build IncidentList with DataTable
- [ ] Create IncidentForm with validation
- [ ] Implement IncidentDetail view
- [ ] Add status workflow UI
- [ ] Integrate with audit logging

### Phase 4: Work Orders Management (Week 5-6)

**Deliverables**:
1. Work orders list and detail views
2. BOM management interface
3. Deliverables tracking
4. Approval workflow UI
5. Scheduling interface

**Tasks**:
- [ ] Build WorkOrderList component
- [ ] Create WorkOrderForm
- [ ] Implement WorkOrderDetail with tabs
- [ ] Add BOM management (add/consume/release parts)
- [ ] Build deliverable submission/review UI
- [ ] Create scheduling calendar component

### Phase 5: Programs & Shops (Week 7)

**Deliverables**:
1. Programs list and detail views
2. Phase management UI
3. Survey viewer
4. Service shops and staff management

**Tasks**:
- [ ] Build ProgramList and detail views
- [ ] Create phase timeline component
- [ ] Implement survey viewer
- [ ] Build ServiceShopList
- [ ] Add staff management interface

### Phase 6: Operations (Week 8)

**Deliverables**:
1. Audit log viewer with advanced filters
2. NATS event stream viewer
3. Background job monitor
4. Configuration management UI

**Tasks**:
- [ ] Build AuditLogViewer with filter panel
- [ ] Create EventStream real-time viewer
- [ ] Implement JobMonitor component
- [ ] Add configuration display/edit UI

### Phase 7: Polish & Testing (Week 9-10)

**Deliverables**:
1. End-to-end testing suite
2. Performance optimization
3. Documentation
4. Deployment automation

**Tasks**:
- [ ] Write E2E tests with Playwright
- [ ] Optimize bundle size and lazy loading
- [ ] Add loading states and error handling
- [ ] Create user documentation
- [ ] Update Helm charts for dashboard
- [ ] Add Dockerfile for dashboard build

---

## 11. Performance Considerations

### 11.1 Frontend Optimization

| Technique | Implementation |
|-----------|----------------|
| Code Splitting | React.lazy() for route-level splitting |
| List Virtualization | TanStack Virtual for long lists |
| Image Optimization | Sharp for thumbnails, lazy loading |
| Bundle Size | Tree shaking, no moment.js |
| Caching | TanStack Query with staleTime |

### 11.2 API Optimization

| Technique | Implementation |
|-----------|----------------|
| Pagination | Cursor-based (existing) |
| Filtering | Server-side with indexes |
| Caching | Redis caching for aggregations |
| Compression | gzip response compression |

### 11.3 Target Metrics

| Metric | Target |
|--------|--------|
| Initial Load | < 3s |
| Time to Interactive | < 5s |
| API Response (p95) | < 200ms |
| Bundle Size (gzipped) | < 500KB |

---

## 12. Monitoring & Observability

### 12.1 Dashboard Metrics

Add to Prometheus:

```yaml
# Dashboard-specific metrics
dashboard_page_views_total:
  type: counter
  labels: [page, user_role]

dashboard_api_calls_total:
  type: counter
  labels: [endpoint, status]

dashboard_websocket_connections:
  type: gauge
  labels: [tenant_id]
```

### 12.2 Error Tracking

- Frontend: Integration with Sentry or similar
- Backend: Existing structured logging with Zap

### 12.3 User Analytics

- Page views and navigation patterns
- Feature usage tracking
- Error frequency by feature

---

## 13. Risks and Mitigations

| Risk | Probability | Impact | Mitigation |
|------|-------------|--------|------------|
| Bundle size bloat | Medium | Medium | Strict bundle analysis, lazy loading |
| API performance | Low | High | Caching, pagination, indexes |
| Auth complexity | Medium | High | Reuse existing JWT infrastructure |
| Browser compatibility | Low | Medium | Target modern browsers only |
| Feature creep | High | High | Strict scope control, MVP focus |

---

## 14. Open Questions

1. **Multi-tenancy in dashboard**: Should admin users see cross-tenant data?
2. **Offline support**: Is offline capability required for field use?
3. **Mobile responsiveness**: Full mobile support or desktop-only initially?
4. **Internationalization**: i18n support in Phase 1 or later?
5. **Dark mode**: Required for initial release?

---

## 15. Appendix

### A. Technology Comparison Matrix

| Criteria | React | Vue 3 | Svelte | HTMX |
|----------|-------|-------|--------|------|
| Ecosystem | 5/5 | 4/5 | 3/5 | 2/5 |
| Learning Curve | 4/5 | 4/5 | 5/5 | 4/5 |
| TypeScript | 5/5 | 4/5 | 4/5 | 2/5 |
| Component Libraries | 5/5 | 4/5 | 3/5 | 2/5 |
| Team Familiarity | 4/5 | 3/5 | 2/5 | 3/5 |
| **Total** | **23** | 19 | 17 | 13 |

### B. UI Library Comparison

| Criteria | shadcn/ui | Ant Design | Material UI |
|----------|-----------|------------|-------------|
| Bundle Size | 5/5 | 2/5 | 3/5 |
| Customization | 5/5 | 3/5 | 4/5 |
| Accessibility | 4/5 | 4/5 | 5/5 |
| Design Quality | 5/5 | 4/5 | 4/5 |
| Learning Curve | 4/5 | 3/5 | 3/5 |
| **Total** | **23** | 16 | 19 |

### C. Related Documents

- OpenAPI Spec: `/home/pato/opt/ESSP/docs/openapi/ims-api.yaml`
- Kubernetes Deployment: `/home/pato/opt/ESSP/deployments/k8s/README.md`
- Monitoring Setup: `/home/pato/opt/ESSP/deployments/monitoring/README.md`
- Helm Chart: `/home/pato/opt/ESSP/charts/essp/values.yaml`

---

## Document History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0.0 | 2025-12-12 | System Architect | Initial draft |
