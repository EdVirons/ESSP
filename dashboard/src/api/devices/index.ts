// Query keys (for external use if needed)
export {
  DEVICES_KEY,
  DEVICE_KEY,
  DEVICE_MODELS_KEY,
  DEVICE_STATS_KEY,
  DEVICE_HISTORY_KEY,
} from './keys';

// Device & Device Model Queries
export {
  useDevices,
  useDevice,
  useDeviceHistory,
  useDeviceStats,
  useDeviceModels,
  useDeviceModel,
  useDeviceMakes,
} from './queries';

// Device & Device Model Mutations
export {
  useCreateDevice,
  useUpdateDevice,
  useDeleteDevice,
  useBulkUpdateDevices,
  useBulkDeleteDevices,
  useCreateDeviceModel,
  useUpdateDeviceModel,
  useDeleteDeviceModel,
  useSyncDevicesFromSSO,
} from './mutations';

// Import/Export Operations
export {
  useImportDevices,
  useExportDevices,
  downloadImportTemplate,
} from './import-export';
