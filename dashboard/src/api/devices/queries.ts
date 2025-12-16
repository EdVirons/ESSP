import { useQuery } from '@tanstack/react-query';
import api from '../client';
import type {
  SSOTDevice,
  SSOTDeviceModel,
  SSOTDeviceStats,
  DeviceFilters,
  DeviceModelFilters,
  DeviceHistoryEntry,
  PaginatedSSOTDevicesResponse,
  SSOTDeviceModelListResponse,
  SSOTDeviceMakesResponse,
} from '@/types/device';
import {
  DEVICES_KEY,
  DEVICE_KEY,
  DEVICE_MODELS_KEY,
  DEVICE_STATS_KEY,
  DEVICE_HISTORY_KEY,
} from './keys';

// ============================================================================
// Device Queries
// ============================================================================

/**
 * Fetch paginated list of devices with filters
 */
export function useDevices(filters: DeviceFilters = {}) {
  return useQuery({
    queryKey: [DEVICES_KEY, filters],
    queryFn: () => api.get<PaginatedSSOTDevicesResponse>('/ssot/devices', filters),
    staleTime: 60_000, // 1 minute
  });
}

/**
 * Fetch a single device by ID
 */
export function useDevice(id: string) {
  return useQuery({
    queryKey: [DEVICE_KEY, id],
    queryFn: () => api.get<SSOTDevice>(`/ssot/devices/${id}`),
    enabled: !!id,
    staleTime: 60_000,
  });
}

/**
 * Fetch device history/audit log
 */
export function useDeviceHistory(deviceId: string) {
  return useQuery({
    queryKey: [DEVICE_HISTORY_KEY, deviceId],
    queryFn: () => api.get<{ items: DeviceHistoryEntry[] }>(`/ssot/devices/${deviceId}/history`),
    enabled: !!deviceId,
    staleTime: 30_000,
  });
}

/**
 * Fetch device statistics
 */
export function useDeviceStats() {
  return useQuery({
    queryKey: [DEVICE_STATS_KEY],
    queryFn: () => api.get<SSOTDeviceStats>('/ssot/devices/stats'),
    staleTime: 30_000,
    refetchInterval: 60_000, // Refresh every minute
  });
}

// ============================================================================
// Device Model Queries
// ============================================================================

/**
 * Fetch list of device models
 */
export function useDeviceModels(filters: DeviceModelFilters = {}) {
  return useQuery({
    queryKey: [DEVICE_MODELS_KEY, filters],
    queryFn: () => api.get<SSOTDeviceModelListResponse>('/ssot/device-models', filters),
    staleTime: 300_000, // 5 minutes - models change infrequently
  });
}

/**
 * Fetch a single device model by ID
 */
export function useDeviceModel(id: string) {
  return useQuery({
    queryKey: [DEVICE_MODELS_KEY, id],
    queryFn: () => api.get<SSOTDeviceModel>(`/ssot/device-models/${id}`),
    enabled: !!id,
    staleTime: 300_000,
  });
}

/**
 * Get unique makes from device models
 */
export function useDeviceMakes() {
  return useQuery({
    queryKey: [DEVICE_MODELS_KEY, 'makes'],
    queryFn: () => api.get<SSOTDeviceMakesResponse>('/ssot/device-models/makes'),
    staleTime: 300_000,
  });
}
