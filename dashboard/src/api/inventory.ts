import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import api from './client';
import type {
  Location,
  CreateLocationRequest,
  UpdateLocationRequest,
  DeviceAssignment,
  AssignDeviceRequest,
  DeviceGroup,
  CreateGroupRequest,
  AddGroupMembersRequest,
  RegisterDeviceRequest,
  InventoryDevice,
  SchoolInventoryResponse,
  LocationsListResponse,
  LocationTreeResponse,
  AssignmentsListResponse,
  GroupsListResponse,
  GroupMembersListResponse,
} from '@/types';

// Query keys
export const INVENTORY_KEYS = {
  all: ['device-inventory'] as const,
  schoolInventory: (schoolId: string) => [...INVENTORY_KEYS.all, 'school', schoolId] as const,
  locations: (schoolId: string) => [...INVENTORY_KEYS.all, 'locations', schoolId] as const,
  locationTree: (schoolId: string) => [...INVENTORY_KEYS.all, 'locations', schoolId, 'tree'] as const,
  assignments: (deviceId: string) => [...INVENTORY_KEYS.all, 'assignments', deviceId] as const,
  locationDevices: (locationId: string) => [...INVENTORY_KEYS.all, 'location-devices', locationId] as const,
  groups: (schoolId: string) => [...INVENTORY_KEYS.all, 'groups', schoolId] as const,
  group: (id: string) => [...INVENTORY_KEYS.all, 'group', id] as const,
  groupDevices: (id: string) => [...INVENTORY_KEYS.all, 'group-devices', id] as const,
};

// ============================================
// School Inventory
// ============================================

export function useSchoolInventory(schoolId: string) {
  return useQuery({
    queryKey: INVENTORY_KEYS.schoolInventory(schoolId),
    queryFn: () => api.get<SchoolInventoryResponse>(`/schools/${schoolId}/inventory`),
    enabled: !!schoolId,
    staleTime: 30_000,
  });
}

// ============================================
// Locations
// ============================================

export function useLocations(schoolId: string) {
  return useQuery({
    queryKey: INVENTORY_KEYS.locations(schoolId),
    queryFn: () => api.get<LocationsListResponse>(`/schools/${schoolId}/locations`),
    enabled: !!schoolId,
    staleTime: 60_000,
  });
}

export function useLocationTree(schoolId: string) {
  return useQuery({
    queryKey: INVENTORY_KEYS.locationTree(schoolId),
    queryFn: () => api.get<LocationTreeResponse>(`/schools/${schoolId}/locations?tree=true`),
    enabled: !!schoolId,
    staleTime: 60_000,
  });
}

export function useCreateLocation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ schoolId, data }: { schoolId: string; data: CreateLocationRequest }) =>
      api.post<Location>(`/schools/${schoolId}/locations`, data),
    onSuccess: (_, { schoolId }) => {
      queryClient.invalidateQueries({ queryKey: INVENTORY_KEYS.locations(schoolId) });
      queryClient.invalidateQueries({ queryKey: INVENTORY_KEYS.locationTree(schoolId) });
      queryClient.invalidateQueries({ queryKey: INVENTORY_KEYS.schoolInventory(schoolId) });
    },
  });
}

export function useUpdateLocation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ schoolId, id, data }: { schoolId: string; id: string; data: UpdateLocationRequest }) =>
      api.put<Location>(`/schools/${schoolId}/locations/${id}`, data),
    onSuccess: (_, { schoolId }) => {
      queryClient.invalidateQueries({ queryKey: INVENTORY_KEYS.locations(schoolId) });
      queryClient.invalidateQueries({ queryKey: INVENTORY_KEYS.locationTree(schoolId) });
      queryClient.invalidateQueries({ queryKey: INVENTORY_KEYS.schoolInventory(schoolId) });
    },
  });
}

export function useDeleteLocation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ schoolId, id }: { schoolId: string; id: string }) =>
      api.delete(`/schools/${schoolId}/locations/${id}`),
    onSuccess: (_, { schoolId }) => {
      queryClient.invalidateQueries({ queryKey: INVENTORY_KEYS.locations(schoolId) });
      queryClient.invalidateQueries({ queryKey: INVENTORY_KEYS.locationTree(schoolId) });
      queryClient.invalidateQueries({ queryKey: INVENTORY_KEYS.schoolInventory(schoolId) });
    },
  });
}

// ============================================
// Device Assignments
// ============================================

export function useDeviceAssignments(deviceId: string) {
  return useQuery({
    queryKey: INVENTORY_KEYS.assignments(deviceId),
    queryFn: () => api.get<AssignmentsListResponse>(`/devices/${deviceId}/assignments`),
    enabled: !!deviceId,
  });
}

export function useLocationDevices(locationId: string) {
  return useQuery({
    queryKey: INVENTORY_KEYS.locationDevices(locationId),
    queryFn: () => api.get<AssignmentsListResponse>(`/locations/${locationId}/devices`),
    enabled: !!locationId,
  });
}

export function useAssignDevice() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ deviceId, data }: { deviceId: string; data: AssignDeviceRequest }) =>
      api.post<{ ok: boolean; assignment: DeviceAssignment }>(`/devices/${deviceId}/assign`, data),
    onSuccess: (_, { deviceId }) => {
      queryClient.invalidateQueries({ queryKey: INVENTORY_KEYS.assignments(deviceId) });
      queryClient.invalidateQueries({ queryKey: INVENTORY_KEYS.all });
    },
  });
}

export function useUnassignDevice() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (deviceId: string) =>
      api.delete<{ ok: boolean }>(`/devices/${deviceId}/assign`),
    onSuccess: (_, deviceId) => {
      queryClient.invalidateQueries({ queryKey: INVENTORY_KEYS.assignments(deviceId) });
      queryClient.invalidateQueries({ queryKey: INVENTORY_KEYS.all });
    },
  });
}

// ============================================
// Device Groups
// ============================================

export function useGroups(schoolId: string) {
  return useQuery({
    queryKey: INVENTORY_KEYS.groups(schoolId),
    queryFn: () => api.get<GroupsListResponse>(`/schools/${schoolId}/groups`),
    enabled: !!schoolId,
    staleTime: 60_000,
  });
}

export function useGroup(id: string) {
  return useQuery({
    queryKey: INVENTORY_KEYS.group(id),
    queryFn: () => api.get<DeviceGroup>(`/groups/${id}`),
    enabled: !!id,
  });
}

export function useGroupDevices(id: string) {
  return useQuery({
    queryKey: INVENTORY_KEYS.groupDevices(id),
    queryFn: () => api.get<GroupMembersListResponse>(`/groups/${id}/devices`),
    enabled: !!id,
  });
}

export function useCreateGroup() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateGroupRequest) =>
      api.post<DeviceGroup>('/groups', data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: INVENTORY_KEYS.all });
    },
  });
}

export function useAddGroupMembers() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ groupId, data }: { groupId: string; data: AddGroupMembersRequest }) =>
      api.post<{ ok: boolean; added: number }>(`/groups/${groupId}/members`, data),
    onSuccess: (_, { groupId }) => {
      queryClient.invalidateQueries({ queryKey: INVENTORY_KEYS.group(groupId) });
      queryClient.invalidateQueries({ queryKey: INVENTORY_KEYS.groupDevices(groupId) });
    },
  });
}

export function useRemoveGroupMember() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ groupId, deviceId }: { groupId: string; deviceId: string }) =>
      api.delete<{ ok: boolean }>(`/groups/${groupId}/members/${deviceId}`),
    onSuccess: (_, { groupId }) => {
      queryClient.invalidateQueries({ queryKey: INVENTORY_KEYS.group(groupId) });
      queryClient.invalidateQueries({ queryKey: INVENTORY_KEYS.groupDevices(groupId) });
    },
  });
}

// ============================================
// Device Registration
// ============================================

export function useRegisterDevice() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ schoolId, data }: { schoolId: string; data: RegisterDeviceRequest }) =>
      api.post<InventoryDevice>(`/schools/${schoolId}/devices`, data),
    onSuccess: (_, { schoolId }) => {
      queryClient.invalidateQueries({ queryKey: INVENTORY_KEYS.schoolInventory(schoolId) });
    },
  });
}

// ============================================
// Impersonation
// ============================================

export interface ImpersonatableUser {
  userId: string;
  name: string;
  email: string;
  schools: string[];
}

export interface ImpersonatableUsersResponse {
  items: ImpersonatableUser[];
}

export interface ValidateImpersonationRequest {
  targetUserId: string;
}

export interface ValidateImpersonationResponse {
  valid: boolean;
  userId: string;
  name: string;
  email: string;
  schools: string[];
  error?: string;
}

export const IMPERSONATION_KEYS = {
  all: ['impersonation'] as const,
  users: () => [...IMPERSONATION_KEYS.all, 'users'] as const,
};

export function useImpersonatableUsers() {
  return useQuery({
    queryKey: IMPERSONATION_KEYS.users(),
    queryFn: () => api.get<ImpersonatableUsersResponse>('/impersonate/users'),
    staleTime: 5 * 60_000, // 5 minutes
  });
}

export function useValidateImpersonation() {
  return useMutation({
    mutationFn: (data: ValidateImpersonationRequest) =>
      api.post<ValidateImpersonationResponse>('/impersonate/validate', data),
  });
}
