import * as React from 'react';
import { Button } from '@/components/ui/button';
import { Modal, ModalHeader, ModalBody, ModalFooter } from '@/components/ui/modal';
import { Select } from '@/components/ui/select';
import { cn } from '@/lib/utils';
import type { ExportOptions, DeviceFilters } from '@/types/device';
import { exportFields, type ExportScope } from './constants';

interface ExportDevicesModalProps {
  open: boolean;
  onClose: () => void;
  onExport: (options: ExportOptions) => Promise<void>;
  isLoading: boolean;
  selectedCount?: number;
  totalCount?: number;
  filteredCount?: number;
  currentFilters?: DeviceFilters;
  selectedIds?: string[];
}

export function ExportDevicesModal({
  open,
  onClose,
  onExport,
  isLoading,
  selectedCount = 0,
  totalCount = 0,
  filteredCount = 0,
  currentFilters,
  selectedIds = [],
}: ExportDevicesModalProps) {
  const [exportScope, setExportScope] = React.useState<ExportScope>(
    selectedCount > 0 ? 'selected' : 'all'
  );
  const [exportFormat, setExportFormat] = React.useState<'csv' | 'xlsx'>('csv');
  const [selectedFields, setSelectedFields] = React.useState<string[]>(
    exportFields.filter((f) => f.default).map((f) => f.key)
  );

  // Reset state when modal opens
  React.useEffect(() => {
    if (open) {
      setExportScope(selectedCount > 0 ? 'selected' : 'all');
    }
  }, [open, selectedCount]);

  // Handle export
  const handleExport = async () => {
    const options: ExportOptions = {
      format: exportFormat,
      fields: selectedFields,
    };

    if (exportScope === 'selected') {
      options.ids = selectedIds;
    } else if (exportScope === 'filtered') {
      options.filters = currentFilters;
    }

    await onExport(options);
    onClose();
  };

  // Toggle field selection
  const toggleField = (key: string) => {
    setSelectedFields((prev) =>
      prev.includes(key) ? prev.filter((f) => f !== key) : [...prev, key]
    );
  };

  // Get scope count
  const getScopeCount = (scope: ExportScope): number => {
    switch (scope) {
      case 'all':
        return totalCount;
      case 'filtered':
        return filteredCount;
      case 'selected':
        return selectedCount;
      default:
        return 0;
    }
  };

  return (
    <Modal open={open} onClose={onClose} className="max-w-lg">
      <ModalHeader onClose={onClose}>Export Devices</ModalHeader>
      <ModalBody>
        <div className="space-y-6">
          {/* Export Scope */}
          <section>
            <h3 className="text-sm font-medium text-gray-900 mb-3">Export Scope</h3>
            <div className="space-y-2">
              <label className="flex items-center gap-3 p-3 border rounded-lg cursor-pointer hover:bg-gray-50">
                <input
                  type="radio"
                  name="scope"
                  checked={exportScope === 'all'}
                  onChange={() => setExportScope('all')}
                  className="text-blue-600"
                />
                <div className="flex-1">
                  <div className="font-medium text-gray-900">All devices</div>
                  <div className="text-sm text-gray-500">
                    Export all {totalCount.toLocaleString()} devices
                  </div>
                </div>
              </label>

              <label
                className={cn(
                  'flex items-center gap-3 p-3 border rounded-lg cursor-pointer',
                  filteredCount === totalCount
                    ? 'opacity-50 cursor-not-allowed'
                    : 'hover:bg-gray-50'
                )}
              >
                <input
                  type="radio"
                  name="scope"
                  checked={exportScope === 'filtered'}
                  onChange={() => setExportScope('filtered')}
                  disabled={filteredCount === totalCount}
                  className="text-blue-600"
                />
                <div className="flex-1">
                  <div className="font-medium text-gray-900">Filtered results</div>
                  <div className="text-sm text-gray-500">
                    Export {filteredCount.toLocaleString()} filtered devices
                  </div>
                </div>
              </label>

              <label
                className={cn(
                  'flex items-center gap-3 p-3 border rounded-lg cursor-pointer',
                  selectedCount === 0 ? 'opacity-50 cursor-not-allowed' : 'hover:bg-gray-50'
                )}
              >
                <input
                  type="radio"
                  name="scope"
                  checked={exportScope === 'selected'}
                  onChange={() => setExportScope('selected')}
                  disabled={selectedCount === 0}
                  className="text-blue-600"
                />
                <div className="flex-1">
                  <div className="font-medium text-gray-900">Selected devices</div>
                  <div className="text-sm text-gray-500">
                    Export {selectedCount.toLocaleString()} selected devices
                  </div>
                </div>
              </label>
            </div>
          </section>

          {/* Export Format */}
          <section>
            <h3 className="text-sm font-medium text-gray-900 mb-2">Format</h3>
            <Select
              value={exportFormat}
              onChange={(value) => setExportFormat(value as 'csv' | 'xlsx')}
              options={[
                { value: 'csv', label: 'CSV (Comma Separated Values)' },
                { value: 'xlsx', label: 'Excel (XLSX)' },
              ]}
            />
          </section>

          {/* Field Selection */}
          <section>
            <h3 className="text-sm font-medium text-gray-900 mb-2">Include Fields</h3>
            <div className="grid grid-cols-2 gap-2 max-h-48 overflow-auto">
              {exportFields.map((field) => (
                <label
                  key={field.key}
                  className="flex items-center gap-2 p-2 hover:bg-gray-50 rounded cursor-pointer"
                >
                  <input
                    type="checkbox"
                    checked={selectedFields.includes(field.key)}
                    onChange={() => toggleField(field.key)}
                    className="text-blue-600 rounded"
                  />
                  <span className="text-sm text-gray-700">{field.label}</span>
                </label>
              ))}
            </div>
          </section>
        </div>
      </ModalBody>
      <ModalFooter>
        <Button variant="outline" onClick={onClose} disabled={isLoading}>
          Cancel
        </Button>
        <Button
          onClick={handleExport}
          disabled={
            isLoading || selectedFields.length === 0 || getScopeCount(exportScope) === 0
          }
        >
          {isLoading
            ? 'Exporting...'
            : `Export ${getScopeCount(exportScope).toLocaleString()} Devices`}
        </Button>
      </ModalFooter>
    </Modal>
  );
}
