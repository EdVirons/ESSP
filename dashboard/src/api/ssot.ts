import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import api from './client';

// Types
export interface SchoolSnapshot {
  tenantId: string;
  schoolId: string;
  name: string;
  countyCode: string;
  countyName: string;
  subCountyCode: string;
  subCountyName: string;
  level?: string;
  type?: string;
  knecCode?: string;
  uic?: string;
  sex?: string;
  cluster?: string;
  accommodation?: string;
  latitude?: number;
  longitude?: number;
  updatedAt: string;
}

export interface DeviceSnapshot {
  tenantId: string;
  deviceId: string;
  schoolId: string;
  model: string;
  serial: string;
  assetTag: string;
  status: string;
  updatedAt: string;
}

export interface PartSnapshot {
  tenantId: string;
  partId: string;
  puk: string;
  name: string;
  category: string;
  unit: string;
  updatedAt: string;
}

export interface SSOTSyncStatus {
  schools: {
    count: number;
    lastSyncAt: string;
    lastCursor: string;
  };
  devices: {
    count: number;
    lastSyncAt: string;
    lastCursor: string;
  };
  parts: {
    count: number;
    lastSyncAt: string;
    lastCursor: string;
  };
}

export interface PaginatedSSOTResponse<T> {
  items: T[];
  total: number;
  limit: number;
  offset: number;
}

export interface SchoolFilters {
  q?: string;
  countyCode?: string;
  level?: string;
  type?: string;
  limit?: number;
  offset?: number;
}

export interface DeviceFilters {
  q?: string;
  schoolId?: string;
  status?: string;
  limit?: number;
  offset?: number;
}

export interface PartFilters {
  q?: string;
  category?: string;
  limit?: number;
  offset?: number;
}

// County and Sub-County types
export interface County {
  code: string;
  name: string;
}

export interface SubCounty {
  code: string;
  name: string;
  countyCode: string;
}

const SCHOOLS_KEY = 'ssot-schools';
const COUNTIES_KEY = 'ssot-counties';
const SUBCOUNTIES_KEY = 'ssot-subcounties';
const DEVICES_KEY = 'ssot-devices';
const PARTS_KEY = 'ssot-parts';
const SYNC_STATUS_KEY = 'ssot-sync-status';

// List schools from snapshot
export function useSchools(filters: SchoolFilters = {}) {
  return useQuery({
    queryKey: [SCHOOLS_KEY, filters],
    queryFn: () => api.get<PaginatedSSOTResponse<SchoolSnapshot>>('/ssot/schools', filters),
    staleTime: 60_000, // 1 minute
  });
}

// List all counties
export function useCounties() {
  return useQuery({
    queryKey: [COUNTIES_KEY],
    queryFn: () => api.get<{ items: County[]; total: number }>('/ssot/schools/counties'),
    staleTime: 300_000, // 5 minutes - counties don't change often
  });
}

// List sub-counties, optionally filtered by county
export function useSubCounties(countyCode?: string) {
  return useQuery({
    queryKey: [SUBCOUNTIES_KEY, countyCode],
    queryFn: () => api.get<{ items: SubCounty[]; total: number }>('/ssot/schools/sub-counties',
      countyCode ? { countyCode } : {}
    ),
    staleTime: 300_000,
  });
}

// List devices from snapshot
export function useDevices(filters: DeviceFilters = {}) {
  return useQuery({
    queryKey: [DEVICES_KEY, filters],
    queryFn: () => api.get<PaginatedSSOTResponse<DeviceSnapshot>>('/ssot/devices', filters),
    staleTime: 60_000,
  });
}

// List parts from snapshot
export function usePartsSnapshot(filters: PartFilters = {}) {
  return useQuery({
    queryKey: [PARTS_KEY, filters],
    queryFn: () => api.get<PaginatedSSOTResponse<PartSnapshot>>('/ssot/parts', filters),
    staleTime: 60_000,
  });
}

// Get sync status
export function useSyncStatus() {
  return useQuery({
    queryKey: [SYNC_STATUS_KEY],
    queryFn: () => api.get<SSOTSyncStatus>('/ssot/status'),
    staleTime: 30_000,
    refetchInterval: 30_000, // Auto refresh every 30 seconds
  });
}

// Trigger sync for schools
export function useSyncSchools() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: () => api.post<{ ok: boolean; synced: number }>('/ssot/sync/schools', {}),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [SCHOOLS_KEY] });
      queryClient.invalidateQueries({ queryKey: [SYNC_STATUS_KEY] });
    },
  });
}

// Trigger sync for devices
export function useSyncDevices() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: () => api.post<{ ok: boolean; synced: number }>('/ssot/sync/devices', {}),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [DEVICES_KEY] });
      queryClient.invalidateQueries({ queryKey: [SYNC_STATUS_KEY] });
    },
  });
}

// Trigger sync for parts
export function useSyncParts() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: () => api.post<{ ok: boolean; synced: number }>('/ssot/sync/parts', {}),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [PARTS_KEY] });
      queryClient.invalidateQueries({ queryKey: [SYNC_STATUS_KEY] });
    },
  });
}
