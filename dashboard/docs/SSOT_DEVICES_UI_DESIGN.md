# SSOT Devices Management UI Design Document

## Executive Summary

This document outlines the comprehensive UI design for the SSOT (Single Source of Truth) Devices management system within the ESSP Dashboard. The design follows existing dashboard patterns while introducing enhanced functionality for device lifecycle management, bulk operations, and model catalog management.

---

## 1. Architecture Overview

### 1.1 Component Hierarchy

```
DevicesPage (Main Container)
|
+-- DevicesStats (Summary cards)
|   +-- TotalDevicesCard
|   +-- StatusBreakdownCards (in_stock, deployed, repair, retired)
|   +-- ModelsCountCard
|
+-- DevicesToolbar
|   +-- SearchInput
|   +-- DeviceFilters (dropdown/inline filters)
|   +-- ViewToggle (table/grid view)
|   +-- ActionButtons (Import, Export, Add Device)
|
+-- DevicesContent
|   +-- DeviceList (table view)
|   |   +-- DataTable with selection
|   |   +-- DeviceRow (individual rows)
|   |   +-- BulkActionsBar (appears when items selected)
|   |
|   +-- DeviceGrid (card view - alternative)
|       +-- DeviceCard (individual cards)
|
+-- Pagination
|
+-- Modals/Sheets
    +-- DeviceDetailSheet (right-side slide panel)
    |   +-- DeviceInfoTab
    |   +-- DeviceHistoryTab
    |   +-- DeviceActionsTab
    |
    +-- CreateDeviceModal
    +-- EditDeviceModal
    +-- DeviceModelManagerModal
    +-- ImportExportPanel (modal/sheet)
    +-- BulkActionConfirmDialog
```

### 1.2 Data Flow

```
                                    +------------------+
                                    |   TanStack Query |
                                    |    (Cache)       |
                                    +--------+---------+
                                             |
     +---------------------------------------+---------------------------------------+
     |                                       |                                       |
     v                                       v                                       v
+----+----+                           +------+------+                         +------+------+
| useDevices |                        | useDeviceModels |                     | useSchools  |
| (list)     |                        | (catalog)       |                     | (lookup)    |
+----+----+                           +------+------+                         +------+------+
     |                                       |                                       |
     +---------------------------------------+---------------------------------------+
                                             |
                                             v
                                    +--------+---------+
                                    |   DevicesPage    |
                                    |   (State Mgmt)   |
                                    +--------+---------+
                                             |
              +------------------------------+------------------------------+
              |                              |                              |
              v                              v                              v
     +--------+--------+            +--------+--------+            +--------+--------+
     |   DeviceList    |            |  DeviceDetail   |            |   Modals        |
     |   Component     |            |  Sheet          |            |   (CRUD)        |
     +--------+--------+            +--------+--------+            +--------+--------+
```

---

## 2. Data Types

### 2.1 Core Types

```typescript
// Device lifecycle status
type DeviceLifecycleStatus = 'in_stock' | 'deployed' | 'repair' | 'retired';

// Device enrollment status
type DeviceEnrollmentStatus = 'enrolled' | 'pending' | 'unenrolled';

// Device model (from device_models table)
interface DeviceModel {
  id: string;
  make: string;           // e.g., "Dell", "HP", "Lenovo"
  model: string;          // e.g., "Latitude 5520"
  category: string;       // e.g., "laptop", "tablet", "desktop"
  specs: {
    processor?: string;
    ram?: string;
    storage?: string;
    display?: string;
    [key: string]: string | undefined;
  };
  imageUrl?: string;
  createdAt: string;
  updatedAt: string;
}

// Device (from devices table)
interface Device {
  id: string;
  tenantId: string;
  serial: string;
  assetTag: string;
  modelId: string;
  model?: DeviceModel;     // Joined model data
  schoolId: string;
  schoolName?: string;     // Joined school name
  lifecycle: DeviceLifecycleStatus;
  enrolled: DeviceEnrollmentStatus;
  assignedTo?: string;     // User ID if assigned
  assignedAt?: string;
  notes?: string;
  warrantyExpiry?: string;
  purchaseDate?: string;
  lastSeen?: string;       // Last check-in timestamp
  createdAt: string;
  updatedAt: string;
}

// Device filters
interface DeviceFilters {
  q?: string;              // Search query
  schoolId?: string;
  modelId?: string;
  lifecycle?: DeviceLifecycleStatus;
  enrolled?: DeviceEnrollmentStatus;
  category?: string;       // From model
  make?: string;           // From model
  limit?: number;
  offset?: number;
  sortBy?: string;
  sortOrder?: 'asc' | 'desc';
}

// Paginated response
interface PaginatedDevicesResponse {
  items: Device[];
  total: number;
  limit: number;
  offset: number;
}

// Import/Export types
interface DeviceImportRow {
  serial: string;
  assetTag: string;
  make: string;
  model: string;
  schoolId: string;
  lifecycle?: DeviceLifecycleStatus;
}

interface ImportResult {
  success: number;
  failed: number;
  errors: Array<{ row: number; error: string }>;
}
```

---

## 3. Component Specifications

### 3.1 DevicesPage (Main Container)

**File:** `/src/pages/DevicesPage.tsx`

**Responsibilities:**
- Orchestrate all device management functionality
- Manage global page state (selected devices, active modal, view mode)
- Coordinate API calls and data refresh
- Handle URL parameters for deep linking

**State Management:**
```typescript
interface DevicesPageState {
  // View state
  viewMode: 'table' | 'grid';

  // Filter state
  filters: DeviceFilters;
  searchQuery: string;

  // Selection state
  selectedDevices: string[];

  // Modal state
  activeModal: 'create' | 'edit' | 'import' | 'export' | 'models' | null;

  // Detail sheet state
  selectedDevice: Device | null;
  showDetail: boolean;
  detailTab: 'info' | 'history' | 'actions';
}
```

**Layout:**
```
+------------------------------------------------------------------+
|  [Header: "Devices" + description]              [Import] [Export] |
|                                                 [+ Add Device]    |
+------------------------------------------------------------------+
|  [Stats Cards Row - 5 cards across]                               |
+------------------------------------------------------------------+
|  [Search] [Status Filter] [School Filter] [Model Filter] [More]  |
|                                            [View: Table | Grid]   |
+------------------------------------------------------------------+
|  [Bulk Actions Bar - appears when items selected]                 |
+------------------------------------------------------------------+
|  +-------------------------------------------------------------+ |
|  |  Device Table/Grid                                           | |
|  |  - Checkbox column                                           | |
|  |  - Device info (serial, asset tag, model)                    | |
|  |  - School assignment                                         | |
|  |  - Status badges                                             | |
|  |  - Actions                                                   | |
|  +-------------------------------------------------------------+ |
+------------------------------------------------------------------+
|  [Pagination: Showing X-Y of Z | Page controls]                   |
+------------------------------------------------------------------+
```

### 3.2 DevicesStats Component

**File:** `/src/components/devices/DevicesStats.tsx`

**Props:**
```typescript
interface DevicesStatsProps {
  stats: {
    total: number;
    byStatus: Record<DeviceLifecycleStatus, number>;
    byEnrollment: Record<DeviceEnrollmentStatus, number>;
    modelsCount: number;
  };
  isLoading?: boolean;
}
```

**Visual Design:**
```
+----------+  +----------+  +----------+  +----------+  +----------+
| Total    |  | In Stock |  | Deployed |  | In Repair|  | Retired  |
| [icon]   |  | [icon]   |  | [icon]   |  | [icon]   |  | [icon]   |
| 1,234    |  | 456      |  | 678      |  | 45       |  | 55       |
| Devices  |  | Available|  | Active   |  | Pending  |  | Archived |
+----------+  +----------+  +----------+  +----------+  +----------+
   gray        green         blue          yellow        red
```

### 3.3 DeviceFilters Component

**File:** `/src/components/devices/DeviceFilters.tsx`

**Props:**
```typescript
interface DeviceFiltersProps {
  filters: DeviceFilters;
  searchQuery: string;
  onSearchChange: (value: string) => void;
  onFilterChange: (key: keyof DeviceFilters, value: string) => void;
  onReset: () => void;
  schools: Array<{ value: string; label: string }>;
  models: Array<{ value: string; label: string }>;
}
```

**Filter Options:**
- **Search:** Free text (serial, asset tag, model name)
- **Lifecycle Status:** All | In Stock | Deployed | Repair | Retired
- **Enrollment:** All | Enrolled | Pending | Unenrolled
- **School:** Dropdown with school names
- **Category:** Laptop | Tablet | Desktop | Other
- **Make:** Dynamic based on available models
- **Advanced Filters (expandable):**
  - Warranty expiry range
  - Purchase date range
  - Last seen range

### 3.4 DeviceList Component

**File:** `/src/components/devices/DeviceList.tsx`

**Props:**
```typescript
interface DeviceListProps {
  devices: Device[];
  isLoading: boolean;
  selectedIds: string[];
  onSelectionChange: (ids: string[]) => void;
  onDeviceClick: (device: Device) => void;
  onStatusChange: (deviceId: string, status: DeviceLifecycleStatus) => void;
}
```

**Table Columns:**
1. **Checkbox** - Row selection
2. **Device** - Icon + Serial + Asset Tag (stacked)
3. **Model** - Make + Model name
4. **School** - School name with icon
5. **Status** - Lifecycle status badge
6. **Enrollment** - Enrollment status badge
7. **Last Seen** - Relative timestamp
8. **Actions** - Quick action buttons

**Row Interactions:**
- Click row: Open detail sheet
- Checkbox: Add to selection
- Status badge: Quick status change dropdown
- Actions menu: Edit, Assign, Change Status, Delete

### 3.5 DeviceCard Component (Grid View)

**File:** `/src/components/devices/DeviceCard.tsx`

**Props:**
```typescript
interface DeviceCardProps {
  device: Device;
  isSelected: boolean;
  onSelect: (selected: boolean) => void;
  onClick: () => void;
}
```

**Visual Design:**
```
+----------------------------------+
| [Checkbox]           [Status]    |
| +----------------------------+   |
| |                            |   |
| |     [Device Icon/Image]    |   |
| |                            |   |
| +----------------------------+   |
| Serial: ABC123456789             |
| Asset: ESSP-2024-001             |
|                                  |
| Dell Latitude 5520               |
| +----------------------------+   |
| | [School Icon] School Name  |   |
| +----------------------------+   |
| [Enrollment Badge] [Last Seen]   |
+----------------------------------+
```

### 3.6 DeviceDetailSheet Component

**File:** `/src/components/devices/DeviceDetailSheet.tsx`

**Props:**
```typescript
interface DeviceDetailSheetProps {
  device: Device | null;
  open: boolean;
  onClose: () => void;
  activeTab: 'info' | 'history' | 'actions';
  onTabChange: (tab: string) => void;
  onUpdate: (device: Partial<Device>) => void;
  onDelete: () => void;
}
```

**Layout:**
```
+------------------------------------------+
| [X] Device Details                        |
+------------------------------------------+
| [Status Badge] [Enrollment Badge]         |
| Serial: ABC123456789                      |
| Asset Tag: ESSP-2024-001                  |
+------------------------------------------+
| [Info] [History] [Actions]   <- Tabs      |
+------------------------------------------+
| Tab Content Area                          |
|                                           |
| INFO TAB:                                 |
| - Device Information                      |
|   - Make/Model                            |
|   - Category                              |
|   - Specs (collapsible)                   |
| - Assignment                              |
|   - School                                |
|   - Assigned To                           |
|   - Assigned Date                         |
| - Dates                                   |
|   - Purchase Date                         |
|   - Warranty Expiry                       |
|   - Last Seen                             |
| - Notes                                   |
|                                           |
| HISTORY TAB:                              |
| - Status change timeline                  |
| - Assignment history                      |
| - Audit events                            |
|                                           |
| ACTIONS TAB:                              |
| - Change Status buttons                   |
| - Assign to School                        |
| - Assign to User                          |
| - Create Work Order                       |
| - Export Device Info                      |
+------------------------------------------+
| [Delete Device]              [Edit]       |
+------------------------------------------+
```

### 3.7 DeviceModelManager Component

**File:** `/src/components/devices/DeviceModelManager.tsx`

**Purpose:** Manage the device model catalog (makes, models, specs)

**Props:**
```typescript
interface DeviceModelManagerProps {
  open: boolean;
  onClose: () => void;
}
```

**Layout:**
```
+------------------------------------------+
| Device Model Catalog              [X]     |
+------------------------------------------+
| [Search models...] [+ Add Model]          |
+------------------------------------------+
| Filter: [All] [Laptop] [Tablet] [Desktop] |
+------------------------------------------+
| +--------------------------------------+ |
| | Make: Dell                           | |
| | - Latitude 5520          [Edit][Del] | |
| | - Latitude 7420          [Edit][Del] | |
| | - Precision 5560         [Edit][Del] | |
| +--------------------------------------+ |
| | Make: HP                             | |
| | - EliteBook 840 G8       [Edit][Del] | |
| | - ProBook 450 G8         [Edit][Del] | |
| +--------------------------------------+ |
| | Make: Lenovo                         | |
| | - ThinkPad X1 Carbon     [Edit][Del] | |
| | - ThinkPad T14           [Edit][Del] | |
| +--------------------------------------+ |
+------------------------------------------+
```

**Add/Edit Model Modal:**
```
+----------------------------------+
| Add Device Model          [X]    |
+----------------------------------+
| Make *         [____________]    |
| Model *        [____________]    |
| Category *     [Laptop     v]    |
+----------------------------------+
| Specifications                   |
| Processor      [____________]    |
| RAM            [____________]    |
| Storage        [____________]    |
| Display        [____________]    |
| + Add Spec                       |
+----------------------------------+
| [Cancel]            [Save Model] |
+----------------------------------+
```

### 3.8 ImportExportPanel Component

**File:** `/src/components/devices/ImportExportPanel.tsx`

**Props:**
```typescript
interface ImportExportPanelProps {
  open: boolean;
  mode: 'import' | 'export';
  onClose: () => void;
  onImportComplete: (result: ImportResult) => void;
  onExportComplete: () => void;
  selectedDeviceIds?: string[]; // For bulk export
}
```

**Import Mode Layout:**
```
+------------------------------------------+
| Import Devices                    [X]     |
+------------------------------------------+
| Step 1: Download Template                 |
| [Download CSV Template]                   |
|                                           |
| Step 2: Upload File                       |
| +--------------------------------------+ |
| |  [Drag & drop CSV file here]         | |
| |  or [Browse Files]                    | |
| +--------------------------------------+ |
|                                           |
| Step 3: Review & Import                   |
| +--------------------------------------+ |
| | Preview (first 5 rows)               | |
| | Serial    | Asset   | Make  | Model  | |
| | ABC123    | E-001   | Dell  | Lat... | |
| | ...       | ...     | ...   | ...    | |
| +--------------------------------------+ |
| Total rows: 150                           |
| Valid: 148 | Errors: 2                    |
|                                           |
| [View Errors]                             |
+------------------------------------------+
| [Cancel]                  [Import 148]    |
+------------------------------------------+
```

**Export Mode Layout:**
```
+------------------------------------------+
| Export Devices                    [X]     |
+------------------------------------------+
| Export Options                            |
|                                           |
| Scope:                                    |
| ( ) All devices (1,234)                   |
| ( ) Filtered results (456)                |
| (x) Selected devices (25)                 |
|                                           |
| Format:                                   |
| [CSV v]                                   |
|                                           |
| Include Fields:                           |
| [x] Serial Number                         |
| [x] Asset Tag                             |
| [x] Make/Model                            |
| [x] School                                |
| [x] Status                                |
| [x] Enrollment                            |
| [ ] Purchase Date                         |
| [ ] Warranty Expiry                       |
| [ ] Notes                                 |
+------------------------------------------+
| [Cancel]                      [Export]    |
+------------------------------------------+
```

### 3.9 BulkActionsBar Component

**File:** `/src/components/devices/BulkActionsBar.tsx`

**Props:**
```typescript
interface BulkActionsBarProps {
  selectedCount: number;
  onClearSelection: () => void;
  onBulkStatusChange: (status: DeviceLifecycleStatus) => void;
  onBulkAssign: () => void;
  onBulkExport: () => void;
  onBulkDelete: () => void;
}
```

**Visual Design:**
```
+------------------------------------------------------------------+
| [x] 25 devices selected                                           |
|                                                                   |
| [Change Status v] [Assign to School] [Export] [Delete]  [Clear]  |
+------------------------------------------------------------------+
```

---

## 4. API Hooks Specification

### 4.1 Device Hooks

**File:** `/src/api/devices.ts`

```typescript
// List devices with filters
export function useDevices(filters: DeviceFilters) {
  return useQuery({
    queryKey: ['devices', filters],
    queryFn: () => api.get<PaginatedDevicesResponse>('/v1/ssot/devices', filters),
    staleTime: 60_000,
  });
}

// Get single device
export function useDevice(id: string) {
  return useQuery({
    queryKey: ['device', id],
    queryFn: () => api.get<Device>(`/v1/ssot/devices/${id}`),
    enabled: !!id,
  });
}

// Create device
export function useCreateDevice() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: CreateDeviceInput) =>
      api.post<Device>('/v1/ssot/devices', data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['devices'] });
    },
  });
}

// Update device
export function useUpdateDevice() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateDeviceInput }) =>
      api.patch<Device>(`/v1/ssot/devices/${id}`, data),
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: ['devices'] });
      queryClient.invalidateQueries({ queryKey: ['device', id] });
    },
  });
}

// Delete device
export function useDeleteDevice() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => api.delete(`/v1/ssot/devices/${id}`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['devices'] });
    },
  });
}

// Bulk update devices
export function useBulkUpdateDevices() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: { ids: string[]; updates: Partial<Device> }) =>
      api.post('/v1/ssot/devices/bulk-update', data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['devices'] });
    },
  });
}

// Import devices
export function useImportDevices() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (file: File) => {
      const formData = new FormData();
      formData.append('file', file);
      return api.post<ImportResult>('/v1/ssot/devices/import', formData);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['devices'] });
    },
  });
}

// Export devices
export function useExportDevices() {
  return useMutation({
    mutationFn: (params: { ids?: string[]; filters?: DeviceFilters }) =>
      api.post<Blob>('/v1/ssot/devices/export', params, { responseType: 'blob' }),
  });
}

// Get device stats
export function useDeviceStats() {
  return useQuery({
    queryKey: ['device-stats'],
    queryFn: () => api.get<DeviceStats>('/v1/ssot/devices/stats'),
    staleTime: 30_000,
  });
}
```

### 4.2 Device Model Hooks

**File:** `/src/api/device-models.ts`

```typescript
// List device models
export function useDeviceModels(filters?: { category?: string; make?: string }) {
  return useQuery({
    queryKey: ['device-models', filters],
    queryFn: () => api.get<DeviceModel[]>('/v1/ssot/device-models', filters),
    staleTime: 300_000, // 5 minutes - models change infrequently
  });
}

// Create device model
export function useCreateDeviceModel() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: CreateDeviceModelInput) =>
      api.post<DeviceModel>('/v1/ssot/device-models', data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['device-models'] });
    },
  });
}

// Update device model
export function useUpdateDeviceModel() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateDeviceModelInput }) =>
      api.patch<DeviceModel>(`/v1/ssot/device-models/${id}`, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['device-models'] });
    },
  });
}

// Delete device model
export function useDeleteDeviceModel() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => api.delete(`/v1/ssot/device-models/${id}`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['device-models'] });
    },
  });
}
```

---

## 5. User Interactions & Workflows

### 5.1 Device Lifecycle Management

```
                    +------------+
                    |  In Stock  |
                    +-----+------+
                          |
                          | Deploy to School
                          v
+------------+      +-----+------+      +------------+
|   Repair   | <--> |  Deployed  | ---> |  Retired   |
+------------+      +------------+      +------------+
     ^                    |
     |                    |
     +--- Send to Repair -+
```

**Status Transition Rules:**
- `in_stock` -> `deployed` (requires school assignment)
- `deployed` -> `repair` (creates work order option)
- `deployed` -> `retired` (requires confirmation)
- `repair` -> `deployed` (after work order completion)
- `repair` -> `retired` (device unrepairable)
- Any status -> `retired` (with confirmation)

### 5.2 Bulk Operations Workflow

1. User selects multiple devices via checkboxes
2. Bulk action bar appears with available actions
3. User clicks action (e.g., "Change Status")
4. Confirmation modal shows:
   - List of affected devices
   - New status selection
   - Warning if any devices can't transition
5. User confirms action
6. Progress indicator during operation
7. Success/failure toast with summary

### 5.3 Import Workflow

1. User clicks "Import" button
2. Import panel opens
3. User downloads template (optional)
4. User uploads CSV file
5. System validates and shows preview
6. User reviews errors (if any)
7. User clicks "Import"
8. Progress indicator during import
9. Results summary shown
10. Device list refreshes

### 5.4 Search & Filter Workflow

1. User types in search box (debounced 300ms)
2. Results update as user types
3. User can add filters:
   - Status dropdown
   - School dropdown
   - Category dropdown
4. Active filters shown as chips
5. "Clear all" resets to default view
6. URL updates with filter state (shareable)

---

## 6. Responsive Design Considerations

### 6.1 Breakpoints

- **Mobile (< 640px):** Single column, cards only, simplified filters
- **Tablet (640px - 1024px):** 2-column grid, collapsible filters
- **Desktop (> 1024px):** Full table view, all filters visible

### 6.2 Mobile Adaptations

- Stats cards: Horizontal scroll or 2x2 grid
- Filters: Collapsed into "Filter" button with sheet
- Table: Switches to card view automatically
- Detail sheet: Full-screen modal instead of side sheet
- Bulk actions: Bottom sticky bar

### 6.3 Touch Interactions

- Swipe left on card: Quick actions
- Long press: Select for bulk operations
- Pull to refresh: Refresh device list

---

## 7. Accessibility Requirements

### 7.1 WCAG 2.1 AA Compliance

- All interactive elements keyboard accessible
- Focus indicators on all focusable elements
- Screen reader announcements for:
  - Selection changes
  - Filter updates
  - Status changes
  - Modal open/close
- Color contrast ratios meet 4.5:1 minimum
- Status indicators not solely color-dependent (icons + text)

### 7.2 Keyboard Navigation

- Tab: Move between interactive elements
- Enter/Space: Activate buttons, checkboxes
- Arrow keys: Navigate table rows
- Escape: Close modals/sheets
- Ctrl+A: Select all (in table context)

---

## 8. File Structure

```
/src
  /api
    devices.ts           # Device API hooks
    device-models.ts     # Device model API hooks
  /components
    /devices
      index.ts           # Barrel export
      DevicesStats.tsx
      DeviceFilters.tsx
      DeviceList.tsx
      DeviceCard.tsx
      DeviceDetailSheet.tsx
      DeviceModelManager.tsx
      ImportExportPanel.tsx
      BulkActionsBar.tsx
      CreateDeviceModal.tsx
      EditDeviceModal.tsx
  /hooks
    useDeviceFilters.ts
    useDeviceSelection.ts
  /pages
    DevicesPage.tsx      # Main page component
  /types
    device.ts            # Device type definitions
```

---

## 9. Implementation Priority

### Phase 1: Core Functionality
1. Device types and API hooks
2. DevicesPage basic structure
3. DeviceList with DataTable
4. DeviceFilters (basic)
5. Pagination

### Phase 2: Detail View & CRUD
1. DeviceDetailSheet
2. CreateDeviceModal
3. EditDeviceModal
4. Status change functionality

### Phase 3: Bulk Operations & Import/Export
1. Device selection
2. BulkActionsBar
3. ImportExportPanel
4. Bulk status changes

### Phase 4: Model Management & Polish
1. DeviceModelManager
2. Grid view (DeviceCard)
3. Mobile responsiveness
4. Accessibility audit
5. Performance optimization

---

## 10. Performance Considerations

- Virtual scrolling for large device lists (1000+ items)
- Debounced search (300ms)
- Optimistic UI updates for status changes
- Prefetch adjacent pages
- Memoized filter/sort operations
- Lazy load device specs in detail view
