import {
  X,
  CheckSquare,
  ArrowRightLeft,
  School,
  Download,
  Trash2,
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Select } from '@/components/ui/select';
import type { DeviceLifecycleStatus } from '@/types/device';
import { LIFECYCLE_STATUS_OPTIONS } from '@/types/device';

interface BulkActionsBarProps {
  selectedCount: number;
  onClearSelection: () => void;
  onBulkStatusChange: (status: DeviceLifecycleStatus) => void;
  onBulkAssign: () => void;
  onBulkExport: () => void;
  onBulkDelete: () => void;
  isUpdating?: boolean;
  isDeleting?: boolean;
}

export function BulkActionsBar({
  selectedCount,
  onClearSelection,
  onBulkStatusChange,
  onBulkAssign,
  onBulkExport,
  onBulkDelete,
  isUpdating,
  isDeleting,
}: BulkActionsBarProps) {
  if (selectedCount === 0) return null;

  const statusOptions = LIFECYCLE_STATUS_OPTIONS.map((opt) => ({
    value: opt.value,
    label: `Set ${opt.label}`,
  }));

  return (
    <div className="sticky top-0 z-20 bg-blue-50 border border-blue-200 rounded-lg p-4 mb-4 shadow-sm">
      <div className="flex flex-wrap items-center justify-between gap-4">
        {/* Selection info */}
        <div className="flex items-center gap-3">
          <div className="flex items-center gap-2 text-blue-700">
            <CheckSquare className="h-5 w-5" />
            <span className="font-medium">
              {selectedCount} device{selectedCount !== 1 ? 's' : ''} selected
            </span>
          </div>
          <button
            type="button"
            onClick={onClearSelection}
            className="text-blue-600 hover:text-blue-800 flex items-center gap-1 text-sm"
          >
            <X className="h-4 w-4" />
            Clear
          </button>
        </div>

        {/* Actions */}
        <div className="flex flex-wrap items-center gap-2">
          {/* Status change dropdown */}
          <div className="flex items-center gap-2">
            <ArrowRightLeft className="h-4 w-4 text-gray-500" />
            <Select
              value=""
              onChange={(value) => {
                if (value) {
                  onBulkStatusChange(value as DeviceLifecycleStatus);
                }
              }}
              options={[{ value: '', label: 'Change Status' }, ...statusOptions]}
              className="w-40"
              disabled={isUpdating}
            />
          </div>

          {/* Assign to school */}
          <Button
            variant="outline"
            size="sm"
            onClick={onBulkAssign}
            disabled={isUpdating}
            className="gap-2"
          >
            <School className="h-4 w-4" />
            Assign School
          </Button>

          {/* Export selected */}
          <Button
            variant="outline"
            size="sm"
            onClick={onBulkExport}
            className="gap-2"
          >
            <Download className="h-4 w-4" />
            Export
          </Button>

          {/* Delete selected */}
          <Button
            variant="outline"
            size="sm"
            onClick={onBulkDelete}
            disabled={isDeleting}
            className="gap-2 text-red-600 hover:text-red-700 hover:bg-red-50 border-red-200"
          >
            <Trash2 className="h-4 w-4" />
            Delete
          </Button>
        </div>
      </div>
    </div>
  );
}
