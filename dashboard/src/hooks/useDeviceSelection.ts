import * as React from 'react';
import type { Device } from '@/types/device';

/**
 * Hook for managing device selection state
 */
export function useDeviceSelection(devices: Device[] = []) {
  const [selectedIds, setSelectedIds] = React.useState<Set<string>>(new Set());

  // Get selected devices
  const selectedDevices = React.useMemo(() => {
    return devices.filter((device) => selectedIds.has(device.id));
  }, [devices, selectedIds]);

  // Check if a device is selected
  const isSelected = React.useCallback((deviceId: string) => {
    return selectedIds.has(deviceId);
  }, [selectedIds]);

  // Check if all devices are selected
  const isAllSelected = React.useMemo(() => {
    if (devices.length === 0) return false;
    return devices.every((device) => selectedIds.has(device.id));
  }, [devices, selectedIds]);

  // Check if some (but not all) devices are selected
  const isSomeSelected = React.useMemo(() => {
    if (devices.length === 0) return false;
    const someSelected = devices.some((device) => selectedIds.has(device.id));
    return someSelected && !isAllSelected;
  }, [devices, selectedIds, isAllSelected]);

  // Toggle selection for a single device
  const toggleSelection = React.useCallback((deviceId: string) => {
    setSelectedIds((prev) => {
      const next = new Set(prev);
      if (next.has(deviceId)) {
        next.delete(deviceId);
      } else {
        next.add(deviceId);
      }
      return next;
    });
  }, []);

  // Select a single device
  const selectDevice = React.useCallback((deviceId: string) => {
    setSelectedIds((prev) => {
      const next = new Set(prev);
      next.add(deviceId);
      return next;
    });
  }, []);

  // Deselect a single device
  const deselectDevice = React.useCallback((deviceId: string) => {
    setSelectedIds((prev) => {
      const next = new Set(prev);
      next.delete(deviceId);
      return next;
    });
  }, []);

  // Select multiple devices
  const selectDevices = React.useCallback((deviceIds: string[]) => {
    setSelectedIds((prev) => {
      const next = new Set(prev);
      deviceIds.forEach((id) => next.add(id));
      return next;
    });
  }, []);

  // Deselect multiple devices
  const deselectDevices = React.useCallback((deviceIds: string[]) => {
    setSelectedIds((prev) => {
      const next = new Set(prev);
      deviceIds.forEach((id) => next.delete(id));
      return next;
    });
  }, []);

  // Select all visible devices
  const selectAll = React.useCallback(() => {
    setSelectedIds(new Set(devices.map((d) => d.id)));
  }, [devices]);

  // Deselect all devices
  const deselectAll = React.useCallback(() => {
    setSelectedIds(new Set());
  }, []);

  // Toggle all visible devices
  const toggleAll = React.useCallback(() => {
    if (isAllSelected) {
      deselectAll();
    } else {
      selectAll();
    }
  }, [isAllSelected, selectAll, deselectAll]);

  // Clear selection when devices change (e.g., pagination)
  React.useEffect(() => {
    // Keep only selected IDs that still exist in current devices
    setSelectedIds((prev) => {
      const currentDeviceIds = new Set(devices.map((d) => d.id));
      const next = new Set<string>();
      prev.forEach((id) => {
        if (currentDeviceIds.has(id)) {
          next.add(id);
        }
      });
      return next;
    });
  }, [devices]);

  return {
    selectedIds: Array.from(selectedIds),
    selectedDevices,
    selectedCount: selectedIds.size,
    isSelected,
    isAllSelected,
    isSomeSelected,
    hasSelection: selectedIds.size > 0,
    toggleSelection,
    selectDevice,
    deselectDevice,
    selectDevices,
    deselectDevices,
    selectAll,
    deselectAll,
    toggleAll,
  };
}
