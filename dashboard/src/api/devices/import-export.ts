import { useMutation, useQueryClient } from '@tanstack/react-query';
import type { ImportResult, ExportOptions } from '@/types/device';
import { DEVICES_KEY, DEVICE_MODELS_KEY, DEVICE_STATS_KEY } from './keys';

// ============================================================================
// Import/Export Operations
// ============================================================================

/**
 * Import devices from CSV file
 */
export function useImportDevices() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (file: File) => {
      const formData = new FormData();
      formData.append('file', file);
      // Use apiClient directly for multipart/form-data
      const response = await fetch('/api/v1/ssot/devices/import', {
        method: 'POST',
        body: formData,
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('auth_token')}`,
          'X-Tenant-ID': localStorage.getItem('tenant_id') || 'demo-tenant',
        },
      });
      if (!response.ok) {
        throw new Error('Import failed');
      }
      return response.json() as Promise<ImportResult>;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [DEVICES_KEY] });
      queryClient.invalidateQueries({ queryKey: [DEVICE_STATS_KEY] });
      queryClient.invalidateQueries({ queryKey: [DEVICE_MODELS_KEY] });
    },
  });
}

/**
 * Export devices to file
 */
export function useExportDevices() {
  return useMutation({
    mutationFn: async (options: ExportOptions) => {
      const response = await fetch('/api/v1/ssot/devices/export', {
        method: 'POST',
        body: JSON.stringify(options),
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('auth_token')}`,
          'X-Tenant-ID': localStorage.getItem('tenant_id') || 'demo-tenant',
        },
      });
      if (!response.ok) {
        throw new Error('Export failed');
      }
      const blob = await response.blob();
      const filename = options.format === 'csv' ? 'devices.csv' : 'devices.xlsx';

      // Trigger download
      const url = window.URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      link.download = filename;
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      window.URL.revokeObjectURL(url);
    },
  });
}

/**
 * Download import template
 */
export function downloadImportTemplate() {
  const headers = [
    'serial',
    'assetTag',
    'make',
    'model',
    'schoolId',
    'lifecycle',
    'enrolled',
    'notes',
    'warrantyExpiry',
    'purchaseDate',
  ];
  const csvContent = headers.join(',') + '\n';
  const blob = new Blob([csvContent], { type: 'text/csv' });
  const url = window.URL.createObjectURL(blob);
  const link = document.createElement('a');
  link.href = url;
  link.download = 'device_import_template.csv';
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
  window.URL.revokeObjectURL(url);
}
