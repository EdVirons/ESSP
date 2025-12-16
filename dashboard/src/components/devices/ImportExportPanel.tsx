import type { ImportResult, ExportOptions, DeviceFilters } from '@/types/device';
import { ImportDevicesModal } from './import-export/ImportDevicesModal';
import { ExportDevicesModal } from './import-export/ExportDevicesModal';

// Re-export the split components for direct usage
export { ImportDevicesModal, ExportDevicesModal } from './import-export';

interface ImportExportPanelProps {
  open: boolean;
  mode: 'import' | 'export';
  onClose: () => void;
  onImport: (file: File) => Promise<ImportResult>;
  onExport: (options: ExportOptions) => Promise<void>;
  isLoading: boolean;
  selectedCount?: number;
  totalCount?: number;
  filteredCount?: number;
  currentFilters?: DeviceFilters;
  selectedIds?: string[];
}

/**
 * ImportExportPanel - Wrapper component for backwards compatibility
 * Renders either ImportDevicesModal or ExportDevicesModal based on mode
 */
export function ImportExportPanel({
  open,
  mode,
  onClose,
  onImport,
  onExport,
  isLoading,
  selectedCount,
  totalCount,
  filteredCount,
  currentFilters,
  selectedIds,
}: ImportExportPanelProps) {
  if (mode === 'import') {
    return (
      <ImportDevicesModal
        open={open}
        onClose={onClose}
        onImport={onImport}
        isLoading={isLoading}
      />
    );
  }

  return (
    <ExportDevicesModal
      open={open}
      onClose={onClose}
      onExport={onExport}
      isLoading={isLoading}
      selectedCount={selectedCount}
      totalCount={totalCount}
      filteredCount={filteredCount}
      currentFilters={currentFilters}
      selectedIds={selectedIds}
    />
  );
}
