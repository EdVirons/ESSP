# ESSP Management Dashboard - Implementation Plan

## Overview

This document provides a detailed task breakdown for implementing the ESSP Management Dashboard. Each task is scoped for execution by a specialized development agent with clear inputs, outputs, and acceptance criteria.

---

## Phase 1: Foundation (Week 1-2)

### 1.1 Project Scaffolding

**Task ID**: DASH-001
**Estimated Time**: 2 hours
**Agent Type**: Frontend Developer

**Description**:
Initialize the dashboard frontend project with Vite, React 18, TypeScript, and essential development tooling.

**Steps**:
1. Create project at `/home/pato/opt/ESSP/dashboard/`
2. Initialize with Vite + React + TypeScript template
3. Configure pnpm as package manager
4. Set up path aliases (@/ for src/)
5. Configure ESLint with React and TypeScript rules
6. Configure Prettier with consistent formatting
7. Add .gitignore for Node.js projects

**Acceptance Criteria**:
- [ ] `pnpm dev` starts development server on port 5173
- [ ] `pnpm build` produces production build in `dist/`
- [ ] `pnpm lint` runs ESLint without errors
- [ ] TypeScript strict mode enabled

**Output Files**:
- `/home/pato/opt/ESSP/dashboard/package.json`
- `/home/pato/opt/ESSP/dashboard/tsconfig.json`
- `/home/pato/opt/ESSP/dashboard/vite.config.ts`
- `/home/pato/opt/ESSP/dashboard/.eslintrc.cjs`
- `/home/pato/opt/ESSP/dashboard/.prettierrc`

---

### 1.2 Tailwind CSS Setup

**Task ID**: DASH-002
**Estimated Time**: 1 hour
**Agent Type**: Frontend Developer
**Depends On**: DASH-001

**Description**:
Install and configure Tailwind CSS with custom theme extending the ESSP brand colors.

**Steps**:
1. Install tailwindcss, postcss, autoprefixer
2. Initialize Tailwind configuration
3. Configure content paths for purging
4. Add custom color palette (primary, secondary, success, warning, error)
5. Configure dark mode support (class-based)
6. Add custom font family (Inter)
7. Set up index.css with Tailwind directives

**Acceptance Criteria**:
- [ ] Tailwind classes compile correctly
- [ ] Custom theme colors available
- [ ] Dark mode toggle works
- [ ] Production build has purged unused CSS

**Output Files**:
- `/home/pato/opt/ESSP/dashboard/tailwind.config.ts`
- `/home/pato/opt/ESSP/dashboard/postcss.config.js`
- `/home/pato/opt/ESSP/dashboard/src/index.css`

---

### 1.3 shadcn/ui Installation

**Task ID**: DASH-003
**Estimated Time**: 2 hours
**Agent Type**: Frontend Developer
**Depends On**: DASH-002

**Description**:
Install shadcn/ui CLI and add core UI components needed across the dashboard.

**Steps**:
1. Initialize shadcn/ui with `npx shadcn-ui@latest init`
2. Configure components.json for Tailwind and React
3. Add core components:
   - Button, Card, Dialog, Input, Label
   - Table, Tabs, Badge, Avatar
   - Dropdown Menu, Command, Popover
   - Sheet (for mobile nav), Toast
   - Form (with react-hook-form integration)
4. Set up utils.ts with cn() helper

**Acceptance Criteria**:
- [ ] All components render correctly
- [ ] Components are accessible (keyboard navigation, ARIA)
- [ ] Dark mode variants work
- [ ] Components are customizable via className

**Output Files**:
- `/home/pato/opt/ESSP/dashboard/components.json`
- `/home/pato/opt/ESSP/dashboard/src/components/ui/*.tsx`
- `/home/pato/opt/ESSP/dashboard/src/lib/utils.ts`

---

### 1.4 Layout Components

**Task ID**: DASH-004
**Estimated Time**: 4 hours
**Agent Type**: Frontend Developer
**Depends On**: DASH-003

**Description**:
Create the main layout components for the dashboard shell.

**Components to Build**:

1. **Layout.tsx**: Main layout wrapper with sidebar and header
2. **Header.tsx**: Top navigation bar with logo, search, user menu
3. **Sidebar.tsx**: Collapsible navigation sidebar
4. **StatusBar.tsx**: Bottom status bar with connection status
5. **Breadcrumb.tsx**: Dynamic breadcrumb based on route

**Requirements**:
- Sidebar collapses to icons on small screens
- Header shows user avatar and dropdown menu
- Breadcrumb auto-generates from route path
- Status bar shows WebSocket connection status
- Responsive design for tablet/desktop

**Acceptance Criteria**:
- [ ] Layout renders with all components
- [ ] Sidebar toggles open/closed
- [ ] Header displays logo and user menu
- [ ] Navigation links are keyboard accessible
- [ ] Mobile responsive (sidebar becomes sheet)

**Output Files**:
- `/home/pato/opt/ESSP/dashboard/src/components/layout/Layout.tsx`
- `/home/pato/opt/ESSP/dashboard/src/components/layout/Header.tsx`
- `/home/pato/opt/ESSP/dashboard/src/components/layout/Sidebar.tsx`
- `/home/pato/opt/ESSP/dashboard/src/components/layout/StatusBar.tsx`
- `/home/pato/opt/ESSP/dashboard/src/components/layout/Breadcrumb.tsx`

---

### 1.5 Routing Setup

**Task ID**: DASH-005
**Estimated Time**: 2 hours
**Agent Type**: Frontend Developer
**Depends On**: DASH-004

**Description**:
Configure React Router with all dashboard routes and lazy loading.

**Routes to Configure**:
```
/                    -> Redirect to /overview
/overview            -> Overview dashboard
/incidents           -> Incidents list
/incidents/:id       -> Incident detail
/work-orders         -> Work orders list
/work-orders/:id     -> Work order detail
/programs            -> Programs list
/programs/:id        -> Program detail
/service-shops       -> Service shops list
/service-shops/:id   -> Service shop detail
/audit-logs          -> Audit logs viewer
/settings            -> Settings page
/login               -> Login page (outside layout)
/*                   -> 404 Not Found
```

**Requirements**:
- Code splitting with React.lazy()
- Loading states during chunk loading
- Route guards for authentication
- Preserve query params on navigation

**Acceptance Criteria**:
- [ ] All routes accessible
- [ ] Lazy loading works with suspense
- [ ] 404 page shows for unknown routes
- [ ] Protected routes redirect to login

**Output Files**:
- `/home/pato/opt/ESSP/dashboard/src/App.tsx`
- `/home/pato/opt/ESSP/dashboard/src/routes.tsx`
- `/home/pato/opt/ESSP/dashboard/src/components/shared/RouteGuard.tsx`

---

### 1.6 API Client Setup

**Task ID**: DASH-006
**Estimated Time**: 3 hours
**Agent Type**: Frontend Developer
**Depends On**: DASH-005

**Description**:
Create the API client with authentication, error handling, and type-safe request/response handling.

**Features**:
1. Axios-based HTTP client with interceptors
2. JWT token management (refresh flow)
3. Automatic tenant/school headers
4. Request/response type definitions
5. Error response handling
6. TanStack Query provider setup

**Type Definitions** (from OpenAPI spec):
- Incident, CreateIncidentRequest
- WorkOrder, CreateWorkOrderRequest
- Program, ServiceShop, ServiceStaff
- AuditLog, Attachment
- Paginated response wrapper

**Acceptance Criteria**:
- [ ] API calls include JWT token
- [ ] 401 responses trigger token refresh
- [ ] Errors are properly typed
- [ ] Types match OpenAPI spec
- [ ] TanStack Query devtools available

**Output Files**:
- `/home/pato/opt/ESSP/dashboard/src/api/client.ts`
- `/home/pato/opt/ESSP/dashboard/src/api/types.ts`
- `/home/pato/opt/ESSP/dashboard/src/providers/QueryProvider.tsx`

---

### 1.7 Authentication Implementation

**Task ID**: DASH-007
**Estimated Time**: 4 hours
**Agent Type**: Frontend Developer
**Depends On**: DASH-006

**Description**:
Implement authentication flow with JWT tokens and protected routes.

**Features**:
1. Login page with form validation
2. JWT token storage (httpOnly cookie preferred, localStorage fallback)
3. Auth context with user state
4. Token refresh mechanism
5. Logout functionality
6. Role-based UI visibility

**Context API**:
```typescript
interface AuthContext {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  login: (credentials: LoginCredentials) => Promise<void>;
  logout: () => void;
  hasPermission: (permission: string) => boolean;
}
```

**Acceptance Criteria**:
- [ ] Login form validates input
- [ ] Successful login stores token and user
- [ ] Failed login shows error message
- [ ] Token refresh works before expiry
- [ ] Logout clears all auth state
- [ ] Permission checks work correctly

**Output Files**:
- `/home/pato/opt/ESSP/dashboard/src/context/AuthContext.tsx`
- `/home/pato/opt/ESSP/dashboard/src/hooks/useAuth.ts`
- `/home/pato/opt/ESSP/dashboard/src/pages/Login.tsx`
- `/home/pato/opt/ESSP/dashboard/src/lib/permissions.ts`

---

### 1.8 Backend Admin Routes

**Task ID**: DASH-008
**Estimated Time**: 3 hours
**Agent Type**: Go Developer
**Depends On**: None

**Description**:
Add admin-specific routes to ims-api for dashboard functionality.

**New Endpoints**:
1. `GET /admin/v1/health/services` - Aggregated service health
2. `GET /admin/v1/metrics/summary` - Dashboard summary metrics
3. `GET /admin/v1/activity` - Recent activity feed
4. `GET /admin/ws` - WebSocket endpoint for real-time updates

**Implementation**:
1. Create `/internal/admin/` package
2. Add health aggregation (calls /healthz on all services)
3. Add metrics summary (queries DB for counts)
4. Add activity feed (recent audit logs + events)
5. Register routes under `/admin` group with ssp_admin role requirement

**Acceptance Criteria**:
- [ ] All endpoints return valid JSON
- [ ] Endpoints require ssp_admin role
- [ ] Health check aggregates all services
- [ ] Metrics summary returns correct counts
- [ ] WebSocket connection works

**Output Files**:
- `/home/pato/opt/ESSP/services/ims-api/internal/admin/routes.go`
- `/home/pato/opt/ESSP/services/ims-api/internal/admin/health.go`
- `/home/pato/opt/ESSP/services/ims-api/internal/admin/metrics.go`
- `/home/pato/opt/ESSP/services/ims-api/internal/admin/activity.go`
- `/home/pato/opt/ESSP/services/ims-api/internal/admin/websocket.go`

---

### 1.9 Static File Embedding

**Task ID**: DASH-009
**Estimated Time**: 2 hours
**Agent Type**: Go Developer
**Depends On**: DASH-008

**Description**:
Embed the dashboard frontend build into the ims-api binary and serve it.

**Implementation**:
1. Create `/internal/dashboard/` package
2. Use `embed.FS` to embed `dashboard/dist/`
3. Serve at `/admin/*` with SPA fallback
4. Handle index.html for client-side routing
5. Set correct MIME types and caching headers

**Build Integration**:
1. Add Makefile target to build dashboard first
2. Update Docker build to include dashboard
3. Ensure binary size is reasonable

**Acceptance Criteria**:
- [ ] Dashboard accessible at `/admin/`
- [ ] All static assets load correctly
- [ ] Client-side routing works (refresh on /admin/incidents)
- [ ] Assets have proper caching headers
- [ ] Binary size increase < 5MB

**Output Files**:
- `/home/pato/opt/ESSP/services/ims-api/internal/dashboard/embed.go`
- `/home/pato/opt/ESSP/services/ims-api/internal/dashboard/handler.go`
- `/home/pato/opt/ESSP/Makefile` (updated)

---

## Phase 2: Overview Dashboard (Week 3)

### 2.1 Service Health Card

**Task ID**: DASH-010
**Estimated Time**: 3 hours
**Agent Type**: Frontend Developer
**Depends On**: DASH-006, DASH-008

**Description**:
Build the service health monitoring card showing status of all microservices.

**Features**:
- Service name, status (healthy/degraded/unhealthy)
- Response latency indicator
- Last check timestamp
- Auto-refresh every 30 seconds
- Visual status indicators (green/yellow/red)

**API Hook**:
```typescript
function useServiceHealth() {
  return useQuery({
    queryKey: ['admin', 'health'],
    queryFn: () => api.get('/admin/v1/health/services'),
    refetchInterval: 30000,
  });
}
```

**Acceptance Criteria**:
- [ ] Shows all 5 services
- [ ] Status color coding works
- [ ] Auto-refresh updates data
- [ ] Loading state displayed
- [ ] Error state handled

**Output Files**:
- `/home/pato/opt/ESSP/dashboard/src/components/overview/ServiceHealthCard.tsx`
- `/home/pato/opt/ESSP/dashboard/src/api/health.ts`

---

### 2.2 Metrics Summary Cards

**Task ID**: DASH-011
**Estimated Time**: 3 hours
**Agent Type**: Frontend Developer
**Depends On**: DASH-006, DASH-008

**Description**:
Build summary metric cards showing key operational statistics.

**Cards**:
1. **Incidents**: Total, Open, SLA Breached
2. **Work Orders**: Total, In Progress, Completed Today
3. **Programs**: Active, Pending Completion
4. **Inventory**: Low Stock Items

**Features**:
- Animated count transitions
- Trend indicators (up/down arrows)
- Click to navigate to filtered list
- Compact design for grid layout

**Acceptance Criteria**:
- [ ] All cards display correct data
- [ ] Numbers animate on change
- [ ] Click navigation works
- [ ] Responsive grid layout
- [ ] Loading skeletons shown

**Output Files**:
- `/home/pato/opt/ESSP/dashboard/src/components/overview/MetricsSummary.tsx`
- `/home/pato/opt/ESSP/dashboard/src/components/overview/MetricCard.tsx`
- `/home/pato/opt/ESSP/dashboard/src/api/metrics.ts`

---

### 2.3 Activity Feed

**Task ID**: DASH-012
**Estimated Time**: 4 hours
**Agent Type**: Frontend Developer
**Depends On**: DASH-006, DASH-008

**Description**:
Build real-time activity feed showing recent system events.

**Features**:
- Virtual scrolling for performance
- Event type icons and colors
- Relative timestamps (2 minutes ago)
- Actor information
- Click to view entity detail
- Real-time updates via WebSocket

**Event Types**:
- incident.created, incident.updated
- workorder.created, workorder.status_changed
- program.created, phase.completed
- user.login, user.logout

**Acceptance Criteria**:
- [ ] Shows 50+ items with virtual scroll
- [ ] Events render with correct icons
- [ ] Timestamps update dynamically
- [ ] WebSocket updates add new items
- [ ] Click navigation works

**Output Files**:
- `/home/pato/opt/ESSP/dashboard/src/components/overview/ActivityFeed.tsx`
- `/home/pato/opt/ESSP/dashboard/src/components/overview/ActivityItem.tsx`
- `/home/pato/opt/ESSP/dashboard/src/hooks/useWebSocket.ts`

---

### 2.4 Alerts Panel

**Task ID**: DASH-013
**Estimated Time**: 2 hours
**Agent Type**: Frontend Developer
**Depends On**: DASH-010, DASH-011

**Description**:
Build alerts panel showing critical system notifications.

**Alert Types**:
- Service degraded/down
- SLA breach warnings
- High error rate
- Low inventory warnings

**Features**:
- Severity levels (critical, warning, info)
- Dismissible alerts
- Persistence in localStorage
- Sound notification for critical

**Acceptance Criteria**:
- [ ] Alerts sorted by severity
- [ ] Dismiss works and persists
- [ ] New alerts animate in
- [ ] Critical alerts highlighted
- [ ] Empty state displayed

**Output Files**:
- `/home/pato/opt/ESSP/dashboard/src/components/overview/AlertsPanel.tsx`
- `/home/pato/opt/ESSP/dashboard/src/components/overview/AlertItem.tsx`

---

### 2.5 Overview Page Assembly

**Task ID**: DASH-014
**Estimated Time**: 2 hours
**Agent Type**: Frontend Developer
**Depends On**: DASH-010, DASH-011, DASH-012, DASH-013

**Description**:
Assemble all overview components into the main dashboard page.

**Layout**:
```
+---------------------------+---------------------------+
|     Service Health        |     Metrics Summary       |
+---------------------------+---------------------------+
|                           |                           |
|     Activity Feed         |     Alerts Panel          |
|                           |                           |
+---------------------------+---------------------------+
```

**Features**:
- Responsive grid layout
- Component loading coordination
- Error boundaries per section
- Refresh all button

**Acceptance Criteria**:
- [ ] All components render correctly
- [ ] Responsive layout works
- [ ] Individual errors don't break page
- [ ] Refresh button works

**Output Files**:
- `/home/pato/opt/ESSP/dashboard/src/pages/Overview.tsx`

---

## Phase 3: Incidents Management (Week 4)

### 3.1 DataTable Component

**Task ID**: DASH-015
**Estimated Time**: 4 hours
**Agent Type**: Frontend Developer
**Depends On**: DASH-003

**Description**:
Build a reusable data table component using TanStack Table.

**Features**:
- Column sorting (client and server-side)
- Column filtering
- Pagination (cursor-based)
- Row selection
- Column visibility toggle
- Responsive (horizontal scroll on mobile)
- Loading and empty states

**API**:
```typescript
interface DataTableProps<TData> {
  columns: ColumnDef<TData>[];
  data: TData[];
  pagination?: PaginationState;
  sorting?: SortingState;
  onPaginationChange?: (pagination: PaginationState) => void;
  onSortingChange?: (sorting: SortingState) => void;
  isLoading?: boolean;
  emptyMessage?: string;
}
```

**Acceptance Criteria**:
- [ ] Sorting works on all column types
- [ ] Pagination shows correct page info
- [ ] Loading state displays skeleton
- [ ] Empty state shows message
- [ ] Row selection tracks correctly

**Output Files**:
- `/home/pato/opt/ESSP/dashboard/src/components/shared/DataTable.tsx`
- `/home/pato/opt/ESSP/dashboard/src/components/shared/DataTablePagination.tsx`
- `/home/pato/opt/ESSP/dashboard/src/components/shared/DataTableToolbar.tsx`

---

### 3.2 Incidents List

**Task ID**: DASH-016
**Estimated Time**: 4 hours
**Agent Type**: Frontend Developer
**Depends On**: DASH-015, DASH-006

**Description**:
Build the incidents list page with filtering and actions.

**Columns**:
- ID (link to detail)
- Title
- Status (badge)
- Severity (badge with color)
- School Name
- Device
- Reported By
- Created At
- Actions (view, update status)

**Filters**:
- Status (multi-select)
- Severity (multi-select)
- Date range
- Search (title, description)

**Acceptance Criteria**:
- [ ] List loads with pagination
- [ ] Filters update URL params
- [ ] Sorting persists in URL
- [ ] Quick status update works
- [ ] Create button navigates to form

**Output Files**:
- `/home/pato/opt/ESSP/dashboard/src/pages/Incidents.tsx`
- `/home/pato/opt/ESSP/dashboard/src/components/incidents/IncidentList.tsx`
- `/home/pato/opt/ESSP/dashboard/src/components/incidents/IncidentFilters.tsx`
- `/home/pato/opt/ESSP/dashboard/src/api/incidents.ts`

---

### 3.3 Incident Form

**Task ID**: DASH-017
**Estimated Time**: 3 hours
**Agent Type**: Frontend Developer
**Depends On**: DASH-003, DASH-006

**Description**:
Build incident creation/edit form with validation.

**Fields**:
- Device ID (searchable select from SSOT)
- Category (select)
- Severity (select)
- Title (required, max 200 chars)
- Description (textarea, optional)
- Reported By (text)

**Validation** (using Zod):
- Device ID required
- Title required, 1-200 characters
- Severity must be valid enum

**Features**:
- Real-time validation
- Async device lookup
- Submit loading state
- Success/error toasts
- Redirect on success

**Acceptance Criteria**:
- [ ] Form validates on submit
- [ ] Device search works
- [ ] Submit creates incident
- [ ] Error messages display
- [ ] Success redirects to list

**Output Files**:
- `/home/pato/opt/ESSP/dashboard/src/components/incidents/IncidentForm.tsx`
- `/home/pato/opt/ESSP/dashboard/src/components/incidents/IncidentFormSchema.ts`

---

### 3.4 Incident Detail

**Task ID**: DASH-018
**Estimated Time**: 4 hours
**Agent Type**: Frontend Developer
**Depends On**: DASH-016, DASH-017

**Description**:
Build incident detail view with status workflow and related data.

**Sections**:
1. Header: Title, status badge, severity, actions
2. Details: Device info, school, contact, dates
3. Description: Full description text
4. Status History: Timeline of status changes
5. Related Work Orders: Linked work orders
6. Attachments: Associated files
7. Audit Trail: Change history

**Actions**:
- Update Status (dropdown with allowed transitions)
- Create Work Order (prefilled from incident)
- Add Attachment

**Acceptance Criteria**:
- [ ] All sections render correctly
- [ ] Status update works with confirmation
- [ ] Create WO navigates with prefill
- [ ] Attachments display/download
- [ ] Back navigation works

**Output Files**:
- `/home/pato/opt/ESSP/dashboard/src/pages/IncidentDetail.tsx`
- `/home/pato/opt/ESSP/dashboard/src/components/incidents/IncidentDetail.tsx`
- `/home/pato/opt/ESSP/dashboard/src/components/incidents/IncidentStatusSelect.tsx`
- `/home/pato/opt/ESSP/dashboard/src/components/incidents/IncidentTimeline.tsx`

---

## Phase 4: Work Orders Management (Week 5-6)

### 4.1 Work Orders List

**Task ID**: DASH-019
**Estimated Time**: 4 hours
**Agent Type**: Frontend Developer
**Depends On**: DASH-015

**Description**:
Build work orders list with advanced filtering.

**Columns**:
- ID, Incident ID
- Status, Repair Location
- Device, School
- Assigned Staff
- Service Shop
- Cost Estimate
- Created, Updated

**Filters**:
- Status
- Repair Location
- Service Shop
- Date range

**Acceptance Criteria**:
- [ ] List with all columns
- [ ] Filters work correctly
- [ ] Sorting on key columns
- [ ] Bulk status update

**Output Files**:
- `/home/pato/opt/ESSP/dashboard/src/pages/WorkOrders.tsx`
- `/home/pato/opt/ESSP/dashboard/src/components/work-orders/WorkOrderList.tsx`
- `/home/pato/opt/ESSP/dashboard/src/api/work-orders.ts`

---

### 4.2 Work Order Form

**Task ID**: DASH-020
**Estimated Time**: 3 hours
**Agent Type**: Frontend Developer
**Depends On**: DASH-006

**Description**:
Build work order creation form.

**Fields**:
- Incident ID (optional, searchable)
- Device ID (required, searchable)
- Task Type
- Repair Location (on_site/service_shop)
- Service Shop (if service_shop)
- Assigned Staff
- Cost Estimate
- Notes

**Acceptance Criteria**:
- [ ] Conditional service shop field
- [ ] Staff filtered by shop
- [ ] Incident prefill works
- [ ] Validation complete

**Output Files**:
- `/home/pato/opt/ESSP/dashboard/src/components/work-orders/WorkOrderForm.tsx`

---

### 4.3 Work Order Detail (Tabs)

**Task ID**: DASH-021
**Estimated Time**: 6 hours
**Agent Type**: Frontend Developer
**Depends On**: DASH-019

**Description**:
Build work order detail with tabbed interface.

**Tabs**:
1. **Overview**: Basic info, status, device details
2. **BOM**: Bill of materials management
3. **Deliverables**: Deliverable tracking
4. **Schedule**: Scheduling entries
5. **Approvals**: Approval requests
6. **Audit**: Change history

**Features**:
- Tab state in URL
- Quick actions in header
- Status workflow

**Acceptance Criteria**:
- [ ] All tabs render correctly
- [ ] Tab state persists in URL
- [ ] Actions work from detail
- [ ] Loading states per tab

**Output Files**:
- `/home/pato/opt/ESSP/dashboard/src/pages/WorkOrderDetail.tsx`
- `/home/pato/opt/ESSP/dashboard/src/components/work-orders/WorkOrderDetail.tsx`
- `/home/pato/opt/ESSP/dashboard/src/components/work-orders/WorkOrderTabs.tsx`

---

### 4.4 BOM Management

**Task ID**: DASH-022
**Estimated Time**: 5 hours
**Agent Type**: Frontend Developer
**Depends On**: DASH-021

**Description**:
Build Bill of Materials management interface for work orders.

**Features**:
- List current BOM items
- Add parts (with suggestion endpoint)
- Consume parts (mark as used)
- Release reserved parts
- Show compatibility warnings

**Add Part Flow**:
1. Search/select part from suggestions
2. Enter quantity
3. Confirm reservation
4. Handle insufficient inventory error

**Acceptance Criteria**:
- [ ] BOM list shows all items
- [ ] Part search with suggestions
- [ ] Add/consume/release work
- [ ] Inventory errors handled
- [ ] Quantities update correctly

**Output Files**:
- `/home/pato/opt/ESSP/dashboard/src/components/work-orders/WorkOrderBOM.tsx`
- `/home/pato/opt/ESSP/dashboard/src/components/work-orders/AddPartDialog.tsx`
- `/home/pato/opt/ESSP/dashboard/src/api/bom.ts`

---

### 4.5 Deliverables Management

**Task ID**: DASH-023
**Estimated Time**: 4 hours
**Agent Type**: Frontend Developer
**Depends On**: DASH-021

**Description**:
Build deliverables tracking and review interface.

**Features**:
- List deliverables with status
- Add new deliverable
- Submit deliverable (with evidence upload)
- Review deliverable (approve/reject)
- Show submission/review history

**Status Flow**:
pending -> submitted -> approved/rejected

**Acceptance Criteria**:
- [ ] Deliverable list correct
- [ ] Add deliverable works
- [ ] Submit with file upload
- [ ] Review approve/reject
- [ ] Status updates reflect

**Output Files**:
- `/home/pato/opt/ESSP/dashboard/src/components/work-orders/WorkOrderDeliverables.tsx`
- `/home/pato/opt/ESSP/dashboard/src/components/work-orders/DeliverableCard.tsx`
- `/home/pato/opt/ESSP/dashboard/src/components/work-orders/SubmitDeliverableDialog.tsx`
- `/home/pato/opt/ESSP/dashboard/src/components/work-orders/ReviewDeliverableDialog.tsx`

---

## Phase 5: Programs & Service Shops (Week 7)

### 5.1 Programs List

**Task ID**: DASH-024
**Estimated Time**: 3 hours
**Agent Type**: Frontend Developer
**Depends On**: DASH-015

**Description**:
Build programs list with status and phase indicators.

**Output Files**:
- `/home/pato/opt/ESSP/dashboard/src/pages/Programs.tsx`
- `/home/pato/opt/ESSP/dashboard/src/components/programs/ProgramList.tsx`
- `/home/pato/opt/ESSP/dashboard/src/api/programs.ts`

---

### 5.2 Program Detail

**Task ID**: DASH-025
**Estimated Time**: 5 hours
**Agent Type**: Frontend Developer
**Depends On**: DASH-024

**Description**:
Build program detail with phases timeline and surveys.

**Sections**:
1. Overview: Status, dates, account manager
2. Phases: Timeline view with status
3. Surveys: List with status
4. BOQ: Bill of quantities

**Output Files**:
- `/home/pato/opt/ESSP/dashboard/src/pages/ProgramDetail.tsx`
- `/home/pato/opt/ESSP/dashboard/src/components/programs/PhaseTimeline.tsx`
- `/home/pato/opt/ESSP/dashboard/src/components/programs/SurveyList.tsx`

---

### 5.3 Service Shops List

**Task ID**: DASH-026
**Estimated Time**: 3 hours
**Agent Type**: Frontend Developer
**Depends On**: DASH-015

**Output Files**:
- `/home/pato/opt/ESSP/dashboard/src/pages/ServiceShops.tsx`
- `/home/pato/opt/ESSP/dashboard/src/components/service-shops/ServiceShopList.tsx`
- `/home/pato/opt/ESSP/dashboard/src/api/service-shops.ts`

---

### 5.4 Service Shop Detail

**Task ID**: DASH-027
**Estimated Time**: 4 hours
**Agent Type**: Frontend Developer
**Depends On**: DASH-026

**Description**:
Build service shop detail with staff and inventory.

**Sections**:
1. Overview: Location, coverage, status
2. Staff: Staff list with roles
3. Inventory: Parts inventory at shop

**Output Files**:
- `/home/pato/opt/ESSP/dashboard/src/pages/ServiceShopDetail.tsx`
- `/home/pato/opt/ESSP/dashboard/src/components/service-shops/StaffList.tsx`
- `/home/pato/opt/ESSP/dashboard/src/components/service-shops/InventoryList.tsx`

---

## Phase 6: Operations (Week 8)

### 6.1 Audit Log Viewer

**Task ID**: DASH-028
**Estimated Time**: 5 hours
**Agent Type**: Frontend Developer
**Depends On**: DASH-015

**Description**:
Build comprehensive audit log viewer with advanced filters.

**Features**:
- Filterable by entity type, action, user, date range
- JSON diff viewer for before/after states
- Export to CSV
- Infinite scroll

**Acceptance Criteria**:
- [ ] All filters work
- [ ] JSON diff renders correctly
- [ ] Export downloads CSV
- [ ] Performance with 1000+ entries

**Output Files**:
- `/home/pato/opt/ESSP/dashboard/src/pages/AuditLogs.tsx`
- `/home/pato/opt/ESSP/dashboard/src/components/operations/AuditLogViewer.tsx`
- `/home/pato/opt/ESSP/dashboard/src/components/operations/AuditLogFilters.tsx`
- `/home/pato/opt/ESSP/dashboard/src/components/operations/JsonDiffViewer.tsx`
- `/home/pato/opt/ESSP/dashboard/src/api/audit-logs.ts`

---

### 6.2 Event Stream Viewer

**Task ID**: DASH-029
**Estimated Time**: 4 hours
**Agent Type**: Frontend Developer
**Depends On**: DASH-012

**Description**:
Build real-time NATS event stream viewer.

**Features**:
- Real-time event display
- Filter by subject pattern
- Pause/resume stream
- Event detail expansion
- Clear/export functionality

**Output Files**:
- `/home/pato/opt/ESSP/dashboard/src/components/operations/EventStream.tsx`
- `/home/pato/opt/ESSP/dashboard/src/components/operations/EventCard.tsx`

---

### 6.3 Settings Page

**Task ID**: DASH-030
**Estimated Time**: 3 hours
**Agent Type**: Frontend Developer
**Depends On**: DASH-007

**Description**:
Build settings page with configuration display.

**Sections**:
1. User Profile
2. Theme Settings
3. Notification Preferences
4. Environment Info (read-only)

**Output Files**:
- `/home/pato/opt/ESSP/dashboard/src/pages/Settings.tsx`
- `/home/pato/opt/ESSP/dashboard/src/components/settings/ProfileSection.tsx`
- `/home/pato/opt/ESSP/dashboard/src/components/settings/ThemeSection.tsx`

---

## Phase 7: Testing & Polish (Week 9-10)

### 7.1 Unit Tests

**Task ID**: DASH-031
**Estimated Time**: 8 hours
**Agent Type**: Frontend Developer

**Description**:
Write unit tests for all components and hooks.

**Coverage Targets**:
- Components: 80%
- Hooks: 90%
- Utils: 100%

**Output Files**:
- `/home/pato/opt/ESSP/dashboard/src/**/*.test.tsx`

---

### 7.2 E2E Tests

**Task ID**: DASH-032
**Estimated Time**: 6 hours
**Agent Type**: Frontend Developer

**Description**:
Write Playwright E2E tests for critical flows.

**Flows to Test**:
1. Login flow
2. Create incident
3. Create and complete work order
4. View and filter audit logs
5. Dashboard overview interaction

**Output Files**:
- `/home/pato/opt/ESSP/dashboard/e2e/*.spec.ts`

---

### 7.3 Documentation

**Task ID**: DASH-033
**Estimated Time**: 4 hours
**Agent Type**: Technical Writer

**Description**:
Write user documentation for dashboard.

**Sections**:
1. Getting Started
2. Dashboard Overview
3. Managing Incidents
4. Managing Work Orders
5. Administration
6. Troubleshooting

**Output Files**:
- `/home/pato/opt/ESSP/docs/dashboard/USER_GUIDE.md`

---

### 7.4 Deployment Updates

**Task ID**: DASH-034
**Estimated Time**: 3 hours
**Agent Type**: DevOps Engineer

**Description**:
Update deployment configuration for dashboard.

**Updates**:
1. Helm chart for dashboard toggle
2. Docker build integration
3. CI/CD pipeline for frontend build
4. Monitoring for dashboard metrics

**Output Files**:
- `/home/pato/opt/ESSP/charts/essp/values.yaml` (updated)
- `/home/pato/opt/ESSP/.github/workflows/dashboard.yml`
- `/home/pato/opt/ESSP/services/ims-api/Dockerfile` (updated)

---

## Summary

| Phase | Tasks | Estimated Hours |
|-------|-------|-----------------|
| 1. Foundation | 9 | 23 |
| 2. Overview | 5 | 14 |
| 3. Incidents | 4 | 15 |
| 4. Work Orders | 5 | 22 |
| 5. Programs/Shops | 4 | 15 |
| 6. Operations | 3 | 12 |
| 7. Testing/Polish | 4 | 21 |
| **Total** | **34** | **122** |

---

## Execution Notes

### Task Dependencies Graph

```
DASH-001 -> DASH-002 -> DASH-003 -> DASH-004 -> DASH-005
                                          |
                                          v
                                    DASH-006 -> DASH-007
                                          |
         DASH-008 -----------+------------+
              |              |
              v              v
         DASH-009     DASH-010, DASH-011, DASH-012
                            |
                            v
                      DASH-013, DASH-014

         DASH-015 -> DASH-016 -> DASH-017 -> DASH-018
              |
              +-> DASH-019 -> DASH-020 -> DASH-021 -> DASH-022, DASH-023
              |
              +-> DASH-024 -> DASH-025
              |
              +-> DASH-026 -> DASH-027
              |
              +-> DASH-028

DASH-012 -> DASH-029
DASH-007 -> DASH-030

All -> DASH-031, DASH-032, DASH-033, DASH-034
```

### Parallel Execution Opportunities

1. **Week 1**: DASH-001 through DASH-007 (sequential), DASH-008 (parallel)
2. **Week 2**: DASH-009 (depends on 008), DASH-010 through DASH-014 (parallel after 006)
3. **Week 3-4**: DASH-015 first, then DASH-016 through DASH-018 (sequential)
4. **Week 5-6**: DASH-019 through DASH-023 (some parallel)
5. **Week 7**: DASH-024 through DASH-027 (some parallel)
6. **Week 8**: DASH-028 through DASH-030 (parallel)
7. **Week 9-10**: DASH-031 through DASH-034 (parallel)

### Risk Mitigation

1. **Scope Creep**: Strict adherence to task acceptance criteria
2. **Integration Issues**: Early integration testing after Phase 1
3. **Performance**: Performance testing after Phase 2
4. **Authentication Complexity**: Fallback to simple token auth if Keycloak issues

---

## Appendix: Quick Reference

### Key Files Created

| Type | Path |
|------|------|
| Frontend Root | `/home/pato/opt/ESSP/dashboard/` |
| React App | `/home/pato/opt/ESSP/dashboard/src/` |
| UI Components | `/home/pato/opt/ESSP/dashboard/src/components/ui/` |
| Feature Components | `/home/pato/opt/ESSP/dashboard/src/components/{feature}/` |
| Pages | `/home/pato/opt/ESSP/dashboard/src/pages/` |
| API Hooks | `/home/pato/opt/ESSP/dashboard/src/api/` |
| Backend Admin | `/home/pato/opt/ESSP/services/ims-api/internal/admin/` |
| Embedded Dashboard | `/home/pato/opt/ESSP/services/ims-api/internal/dashboard/` |

### Technology Stack Quick Reference

| Tool | Purpose | Version |
|------|---------|---------|
| React | UI Framework | 18.x |
| TypeScript | Type Safety | 5.x |
| Vite | Build Tool | 5.x |
| Tailwind CSS | Styling | 3.x |
| shadcn/ui | Component Library | latest |
| TanStack Query | Data Fetching | 5.x |
| TanStack Table | Tables | 8.x |
| React Router | Routing | 6.x |
| React Hook Form | Forms | 7.x |
| Zod | Validation | 3.x |
| Recharts | Charts | 2.x |
