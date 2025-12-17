// Location types
export type LocationType = 'block' | 'floor' | 'room' | 'lab' | 'storage' | 'office';

export interface Location {
  id: string;
  tenantId: string;
  schoolId: string;
  parentId: string | null;
  name: string;
  locationType: LocationType;
  code: string;
  capacity: number;
  metadata?: Record<string, unknown>;
  active: boolean;
  createdAt: string;
  updatedAt: string;
  // Computed
  path?: string;
  deviceCount?: number;
}

export interface LocationTreeNode extends Location {
  children?: LocationTreeNode[];
}

export interface CreateLocationRequest {
  parentId?: string;
  name: string;
  locationType?: LocationType;
  code?: string;
  capacity?: number;
}

export interface UpdateLocationRequest {
  parentId?: string;
  name?: string;
  locationType?: LocationType;
  code?: string;
  capacity?: number;
  active?: boolean;
}

// Assignment types
export type AssignmentType = 'permanent' | 'temporary' | 'loan' | 'repair' | 'storage';

export interface DeviceAssignment {
  id: string;
  tenantId: string;
  deviceId: string;
  locationId: string | null;
  assignedToUser: string;
  assignmentType: AssignmentType;
  effectiveFrom: string;
  effectiveTo: string | null;
  notes: string;
  createdBy: string;
  createdAt: string;
  // Joined
  location?: Location;
}

export interface AssignDeviceRequest {
  locationId?: string;
  assignedToUser?: string;
  assignmentType?: AssignmentType;
  notes?: string;
}

// Group types
export type GroupType = 'manual' | 'location' | 'dynamic';

export interface GroupSelector {
  model?: string;
  lifecycle?: string[];
  locationType?: string;
  locationId?: string;
  schoolId?: string;
}

export interface DeviceGroup {
  id: string;
  tenantId: string;
  schoolId: string | null;
  name: string;
  description: string;
  groupType: GroupType;
  locationId: string | null;
  selector?: GroupSelector;
  policies?: Record<string, unknown>;
  active: boolean;
  createdAt: string;
  updatedAt: string;
  // Computed
  memberCount?: number;
}

export interface GroupMember {
  id: string;
  tenantId: string;
  groupId: string;
  deviceId: string;
  addedAt: string;
  addedBy: string;
}

export interface CreateGroupRequest {
  schoolId?: string;
  name: string;
  description?: string;
  groupType?: GroupType;
  locationId?: string;
  selector?: GroupSelector;
}

export interface AddGroupMembersRequest {
  deviceIds: string[];
}

// Device registration
export interface RegisterDeviceRequest {
  serial: string;
  assetTag?: string;
  model: string;
  make?: string;
  notes?: string;
  locationId?: string;
}

// Inventory device (enriched device view)
export interface InventoryDevice {
  id: string;
  tenantId: string;
  serial: string;
  assetTag: string;
  model: string;
  make: string;
  schoolId: string;
  lifecycle: string;
  enrolled: boolean;
  lastSeenAt?: string;
  createdAt: string;
  updatedAt: string;
  // Joined
  location?: Location;
  locationPath?: string;
  macAddresses?: string[];
  groups?: string[];
}

// Inventory summary
export interface InventorySummary {
  totalDevices: number;
  byStatus: Record<string, number>;
  byLocation: {
    assigned: number;
    unassigned: number;
  };
  byModel?: Record<string, number>;
}

// School inventory response
export interface SchoolInventoryResponse {
  school: {
    id: string;
    name: string;
  };
  summary: InventorySummary;
  devices: InventoryDevice[];
  locations: Location[];
}

// List responses
export interface LocationsListResponse {
  items: Location[];
}

export interface LocationTreeResponse {
  tree: LocationTreeNode[];
}

export interface AssignmentsListResponse {
  items: DeviceAssignment[];
}

export interface GroupsListResponse {
  items: DeviceGroup[];
}

export interface GroupMembersListResponse {
  items: GroupMember[];
}

// Constants for UI
export const LOCATION_TYPE_OPTIONS: Array<{ value: LocationType; label: string }> = [
  { value: 'block', label: 'Block' },
  { value: 'floor', label: 'Floor' },
  { value: 'room', label: 'Room' },
  { value: 'lab', label: 'Computer Lab' },
  { value: 'storage', label: 'Storage' },
  { value: 'office', label: 'Office' },
];

export const ASSIGNMENT_TYPE_OPTIONS: Array<{ value: AssignmentType; label: string }> = [
  { value: 'permanent', label: 'Permanent' },
  { value: 'temporary', label: 'Temporary' },
  { value: 'loan', label: 'On Loan' },
  { value: 'repair', label: 'In Repair' },
  { value: 'storage', label: 'In Storage' },
];

export const GROUP_TYPE_OPTIONS: Array<{ value: GroupType; label: string }> = [
  { value: 'manual', label: 'Manual' },
  { value: 'location', label: 'Location-based' },
  { value: 'dynamic', label: 'Dynamic' },
];

export const LOCATION_TYPE_ICONS: Record<LocationType, string> = {
  block: 'Building2',
  floor: 'Layers',
  room: 'DoorOpen',
  lab: 'Monitor',
  storage: 'Warehouse',
  office: 'Briefcase',
};
