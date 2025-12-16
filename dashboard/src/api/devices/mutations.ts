import { useMutation, useQueryClient } from '@tanstack/react-query';
import api from '../client';
import type {
  SSOTDevice,
  SSOTDeviceModel,
  CreateDeviceInput,
  UpdateDeviceInput,
  CreateDeviceModelInput,
  UpdateDeviceModelInput,
  BulkUpdateDevicesInput,
  BulkDeleteDevicesInput,
} from '@/types/device';
import {
  DEVICES_KEY,
  DEVICE_KEY,
  DEVICE_MODELS_KEY,
  DEVICE_STATS_KEY,
} from './keys';

// ============================================================================
// Device Mutations
// ============================================================================

/**
 * Create a new device
 */
export function useCreateDevice() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateDeviceInput) =>
      api.post<SSOTDevice>('/ssot/devices', data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [DEVICES_KEY] });
      queryClient.invalidateQueries({ queryKey: [DEVICE_STATS_KEY] });
    },
  });
}

/**
 * Update an existing device
 */
export function useUpdateDevice() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateDeviceInput }) =>
      api.patch<SSOTDevice>(`/ssot/devices/${id}`, data),
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: [DEVICES_KEY] });
      queryClient.invalidateQueries({ queryKey: [DEVICE_KEY, id] });
      queryClient.invalidateQueries({ queryKey: [DEVICE_STATS_KEY] });
    },
  });
}

/**
 * Delete a device
 */
export function useDeleteDevice() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => api.delete(`/ssot/devices/${id}`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [DEVICES_KEY] });
      queryClient.invalidateQueries({ queryKey: [DEVICE_STATS_KEY] });
    },
  });
}

/**
 * Bulk update multiple devices
 */
export function useBulkUpdateDevices() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: BulkUpdateDevicesInput) =>
      api.post<{ updated: number }>('/ssot/devices/bulk-update', data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [DEVICES_KEY] });
      queryClient.invalidateQueries({ queryKey: [DEVICE_STATS_KEY] });
    },
  });
}

/**
 * Bulk delete multiple devices
 */
export function useBulkDeleteDevices() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: BulkDeleteDevicesInput) =>
      api.post<{ deleted: number }>('/ssot/devices/bulk-delete', data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [DEVICES_KEY] });
      queryClient.invalidateQueries({ queryKey: [DEVICE_STATS_KEY] });
    },
  });
}

// ============================================================================
// Device Model Mutations
// ============================================================================

/**
 * Create a new device model
 */
export function useCreateDeviceModel() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateDeviceModelInput) =>
      api.post<SSOTDeviceModel>('/ssot/device-models', data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [DEVICE_MODELS_KEY] });
    },
  });
}

/**
 * Update an existing device model
 */
export function useUpdateDeviceModel() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateDeviceModelInput }) =>
      api.patch<SSOTDeviceModel>(`/ssot/device-models/${id}`, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [DEVICE_MODELS_KEY] });
    },
  });
}

/**
 * Delete a device model
 */
export function useDeleteDeviceModel() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => api.delete(`/ssot/device-models/${id}`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [DEVICE_MODELS_KEY] });
    },
  });
}

// ============================================================================
// Sync Operations (from SSOT service)
// ============================================================================

/**
 * Trigger device sync from SSOT service
 */
export function useSyncDevicesFromSSO() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: () => api.post<{ ok: boolean; synced: number }>('/ssot/sync/devices', {}),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [DEVICES_KEY] });
      queryClient.invalidateQueries({ queryKey: [DEVICE_STATS_KEY] });
    },
  });
}
