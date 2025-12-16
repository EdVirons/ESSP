// Device lifecycle status
export type DeviceLifecycleStatus = 'in_stock' | 'deployed' | 'repair' | 'retired';

// SSOT Device Snapshot (simplified view from SSOT service)
export interface SSOTDevice {
  tenantId: string;
  deviceId: string;
  schoolId: string;
  model: string;
  serial: string;
  assetTag: string;
  status: DeviceLifecycleStatus;
  updatedAt: string;
}

// SSOT Device Stats
export interface SSOTDeviceStats {
  total: number;
  byStatus: Record<string, number>;
  bySchool: Record<string, number>;
  uniqueModels: number;
}

// SSOT Device Model (derived from snapshot data)
export interface SSOTDeviceModel {
  model: string;
  count: number;
}

// Paginated SSOT responses
export interface PaginatedSSOTDevicesResponse {
  items: SSOTDevice[];
  total: number;
  limit: number;
  offset: number;
}

export interface SSOTDeviceModelListResponse {
  items: SSOTDeviceModel[];
  total: number;
}

export interface SSOTDeviceMakesResponse {
  makes: string[];
}

// Device enrollment status
export type DeviceEnrollmentStatus = 'enrolled' | 'pending' | 'unenrolled';

// Device category
export type DeviceCategory = 'laptop' | 'tablet' | 'desktop' | 'chromebook' | 'other';

// Device model specifications
export interface DeviceModelSpecs {
  processor?: string;
  ram?: string;
  storage?: string;
  display?: string;
  battery?: string;
  weight?: string;
  os?: string;
  [key: string]: string | undefined;
}

// Device model (from device_models table)
export interface DeviceModel {
  id: string;
  tenantId: string;
  make: string;
  model: string;
  category: DeviceCategory;
  specs: DeviceModelSpecs;
  imageUrl?: string;
  createdAt: string;
  updatedAt: string;
}

// Device (from devices table)
export interface Device {
  id: string;
  tenantId: string;
  serial: string;
  assetTag: string;
  modelId: string;
  model?: DeviceModel;
  schoolId: string;
  schoolName?: string;
  lifecycle: DeviceLifecycleStatus;
  enrolled: DeviceEnrollmentStatus;
  assignedTo?: string;
  assignedToName?: string;
  assignedAt?: string;
  notes?: string;
  warrantyExpiry?: string;
  purchaseDate?: string;
  lastSeen?: string;
  createdAt: string;
  updatedAt: string;
}

// Device history/audit entry
export interface DeviceHistoryEntry {
  id: string;
  deviceId: string;
  action: string;
  field?: string;
  oldValue?: string;
  newValue?: string;
  actor: string;
  actorName?: string;
  timestamp: string;
  metadata?: Record<string, unknown>;
}

// Device filters for API queries
export interface DeviceFilters {
  q?: string;
  schoolId?: string;
  modelId?: string;
  lifecycle?: DeviceLifecycleStatus;
  enrolled?: DeviceEnrollmentStatus;
  category?: DeviceCategory;
  make?: string;
  warrantyBefore?: string;
  warrantyAfter?: string;
  purchaseBefore?: string;
  purchaseAfter?: string;
  limit?: number;
  offset?: number;
  sortBy?: 'serial' | 'assetTag' | 'model' | 'school' | 'lifecycle' | 'lastSeen' | 'createdAt';
  sortOrder?: 'asc' | 'desc';
}

// Device model filters
export interface DeviceModelFilters {
  q?: string;
  category?: DeviceCategory;
  make?: string;
}

// Paginated responses
export interface PaginatedDevicesResponse {
  items: Device[];
  total: number;
  limit: number;
  offset: number;
}

export interface DeviceModelListResponse {
  items: DeviceModel[];
  total: number;
}

// Device stats
export interface DeviceStats {
  total: number;
  byLifecycle: Record<DeviceLifecycleStatus, number>;
  byEnrollment: Record<DeviceEnrollmentStatus, number>;
  byCategory: Record<DeviceCategory, number>;
  modelsCount: number;
  schoolsWithDevices: number;
}

// Create/Update inputs
export interface CreateDeviceInput {
  serial: string;
  assetTag: string;
  modelId: string;
  schoolId: string;
  lifecycle?: DeviceLifecycleStatus;
  enrolled?: DeviceEnrollmentStatus;
  assignedTo?: string;
  notes?: string;
  warrantyExpiry?: string;
  purchaseDate?: string;
}

export interface UpdateDeviceInput {
  serial?: string;
  assetTag?: string;
  modelId?: string;
  schoolId?: string;
  lifecycle?: DeviceLifecycleStatus;
  enrolled?: DeviceEnrollmentStatus;
  assignedTo?: string;
  notes?: string;
  warrantyExpiry?: string;
  purchaseDate?: string;
}

export interface CreateDeviceModelInput {
  make: string;
  model: string;
  category: DeviceCategory;
  specs?: DeviceModelSpecs;
  imageUrl?: string;
}

export interface UpdateDeviceModelInput {
  make?: string;
  model?: string;
  category?: DeviceCategory;
  specs?: DeviceModelSpecs;
  imageUrl?: string;
}

// Bulk operations
export interface BulkUpdateDevicesInput {
  ids: string[];
  updates: {
    lifecycle?: DeviceLifecycleStatus;
    enrolled?: DeviceEnrollmentStatus;
    schoolId?: string;
  };
}

export interface BulkDeleteDevicesInput {
  ids: string[];
}

// Import/Export
export interface DeviceImportRow {
  serial: string;
  assetTag: string;
  make: string;
  model: string;
  schoolId: string;
  lifecycle?: DeviceLifecycleStatus;
  enrolled?: DeviceEnrollmentStatus;
  notes?: string;
  warrantyExpiry?: string;
  purchaseDate?: string;
}

export interface ImportResult {
  success: number;
  failed: number;
  created: number;
  updated: number;
  errors: Array<{
    row: number;
    serial?: string;
    error: string;
  }>;
}

export interface ExportOptions {
  ids?: string[];
  filters?: DeviceFilters;
  format: 'csv' | 'xlsx';
  fields: string[];
}

// Status transition helpers
export const LIFECYCLE_STATUS_OPTIONS: Array<{ value: DeviceLifecycleStatus; label: string }> = [
  { value: 'in_stock', label: 'In Stock' },
  { value: 'deployed', label: 'Deployed' },
  { value: 'repair', label: 'In Repair' },
  { value: 'retired', label: 'Retired' },
];

export const ENROLLMENT_STATUS_OPTIONS: Array<{ value: DeviceEnrollmentStatus; label: string }> = [
  { value: 'enrolled', label: 'Enrolled' },
  { value: 'pending', label: 'Pending' },
  { value: 'unenrolled', label: 'Unenrolled' },
];

export const DEVICE_CATEGORY_OPTIONS: Array<{ value: DeviceCategory; label: string }> = [
  { value: 'laptop', label: 'Laptop' },
  { value: 'tablet', label: 'Tablet' },
  { value: 'desktop', label: 'Desktop' },
  { value: 'chromebook', label: 'Chromebook' },
  { value: 'other', label: 'Other' },
];

// Status colors for badges
export const LIFECYCLE_STATUS_COLORS: Record<DeviceLifecycleStatus, string> = {
  in_stock: 'bg-green-100 text-green-800',
  deployed: 'bg-blue-100 text-blue-800',
  repair: 'bg-yellow-100 text-yellow-800',
  retired: 'bg-gray-100 text-gray-800',
};

export const ENROLLMENT_STATUS_COLORS: Record<DeviceEnrollmentStatus, string> = {
  enrolled: 'bg-emerald-100 text-emerald-800',
  pending: 'bg-orange-100 text-orange-800',
  unenrolled: 'bg-red-100 text-red-800',
};

// Status transition rules
export const LIFECYCLE_TRANSITIONS: Record<DeviceLifecycleStatus, DeviceLifecycleStatus[]> = {
  in_stock: ['deployed', 'retired'],
  deployed: ['repair', 'retired', 'in_stock'],
  repair: ['deployed', 'retired', 'in_stock'],
  retired: [], // Terminal state
};
