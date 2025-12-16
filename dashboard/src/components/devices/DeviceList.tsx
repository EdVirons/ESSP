import * as React from 'react';
import { type ColumnDef } from '@tanstack/react-table';
import {
  Laptop,
  Tag,
  School,
  ChevronRight,
  MoreHorizontal,
  Edit,
  Trash2,
  ArrowRightLeft,
  ExternalLink,
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { DataTable, SortableHeader, createSelectColumn } from '@/components/ui/data-table';
import { formatRelativeTime, formatStatus, cn } from '@/lib/utils';
import type { Device, DeviceLifecycleStatus } from '@/types/device';
import {
  LIFECYCLE_STATUS_COLORS,
  ENROLLMENT_STATUS_COLORS,
  LIFECYCLE_TRANSITIONS,
} from '@/types/device';

interface DeviceListProps {
  devices: Device[];
  isLoading: boolean;
  selectedIds: string[];
  onSelectionChange: (ids: string[]) => void;
  onDeviceClick: (device: Device) => void;
  onEditDevice: (device: Device) => void;
  onDeleteDevice: (device: Device) => void;
  onStatusChange: (device: Device, status: DeviceLifecycleStatus) => void;
}

export function DeviceList({
  devices,
  isLoading,
  selectedIds,
  onSelectionChange,
  onDeviceClick,
  onEditDevice,
  onDeleteDevice,
  onStatusChange,
}: DeviceListProps) {
  const [openMenuId, setOpenMenuId] = React.useState<string | null>(null);

  // Convert selectedIds to row selection state
  const rowSelection = React.useMemo(() => {
    const selection: Record<string, boolean> = {};
    devices.forEach((device, index) => {
      if (selectedIds.includes(device.id)) {
        selection[index.toString()] = true;
      }
    });
    return selection;
  }, [devices, selectedIds]);

  // Handle row selection change
  const handleRowSelectionChange = React.useCallback(
    (newSelection: Record<string, boolean>) => {
      const newSelectedIds: string[] = [];
      Object.entries(newSelection).forEach(([index, isSelected]) => {
        if (isSelected) {
          const device = devices[parseInt(index)];
          if (device) {
            newSelectedIds.push(device.id);
          }
        }
      });
      onSelectionChange(newSelectedIds);
    },
    [devices, onSelectionChange]
  );

  // Close menu when clicking outside
  React.useEffect(() => {
    const handleClick = () => setOpenMenuId(null);
    if (openMenuId) {
      document.addEventListener('click', handleClick);
      return () => document.removeEventListener('click', handleClick);
    }
  }, [openMenuId]);

  const columns: ColumnDef<Device>[] = React.useMemo(
    () => [
      createSelectColumn<Device>(),
      {
        accessorKey: 'serial',
        header: ({ column }) => (
          <SortableHeader column={column}>Device</SortableHeader>
        ),
        cell: ({ row }) => (
          <div className="flex items-center gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-indigo-50">
              <Laptop className="h-5 w-5 text-indigo-600" />
            </div>
            <div className="min-w-0">
              <div className="font-mono font-medium text-gray-900">
                {row.original.serial}
              </div>
              <div className="flex items-center gap-1 text-sm text-gray-500">
                <Tag className="h-3 w-3" />
                {row.original.assetTag || '-'}
              </div>
            </div>
          </div>
        ),
      },
      {
        accessorKey: 'model',
        header: 'Model',
        cell: ({ row }) => {
          const model = row.original.model;
          return (
            <div className="min-w-0">
              <div className="font-medium text-gray-900">
                {model ? `${model.make} ${model.model}` : '-'}
              </div>
              <div className="text-sm text-gray-500 capitalize">
                {model?.category || '-'}
              </div>
            </div>
          );
        },
      },
      {
        accessorKey: 'schoolName',
        header: ({ column }) => (
          <SortableHeader column={column}>School</SortableHeader>
        ),
        cell: ({ row }) => (
          <div className="flex items-center gap-2">
            <School className="h-4 w-4 text-gray-400" />
            <span className="text-gray-900 truncate max-w-[150px]">
              {row.original.schoolName || row.original.schoolId || '-'}
            </span>
          </div>
        ),
      },
      {
        accessorKey: 'lifecycle',
        header: 'Status',
        cell: ({ row }) => {
          const status = row.original.lifecycle;
          return (
            <Badge className={cn('capitalize', LIFECYCLE_STATUS_COLORS[status])}>
              {formatStatus(status)}
            </Badge>
          );
        },
      },
      {
        accessorKey: 'enrolled',
        header: 'Enrollment',
        cell: ({ row }) => {
          const enrolled = row.original.enrolled;
          return (
            <Badge className={cn('capitalize', ENROLLMENT_STATUS_COLORS[enrolled])}>
              {formatStatus(enrolled)}
            </Badge>
          );
        },
      },
      {
        accessorKey: 'lastSeen',
        header: ({ column }) => (
          <SortableHeader column={column}>Last Seen</SortableHeader>
        ),
        cell: ({ row }) => (
          <span className="text-sm text-gray-600">
            {row.original.lastSeen
              ? formatRelativeTime(row.original.lastSeen)
              : 'Never'}
          </span>
        ),
      },
      {
        id: 'actions',
        cell: ({ row }) => {
          const device = row.original;
          const availableTransitions = LIFECYCLE_TRANSITIONS[device.lifecycle] || [];
          const isMenuOpen = openMenuId === device.id;

          return (
            <div className="flex items-center gap-1">
              <Button
                variant="ghost"
                size="sm"
                onClick={(e) => {
                  e.stopPropagation();
                  onDeviceClick(device);
                }}
                className="h-8 w-8 p-0"
                title="View details"
              >
                <ChevronRight className="h-4 w-4" />
              </Button>

              <div className="relative">
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={(e) => {
                    e.stopPropagation();
                    setOpenMenuId(isMenuOpen ? null : device.id);
                  }}
                  className="h-8 w-8 p-0"
                  title="More actions"
                >
                  <MoreHorizontal className="h-4 w-4" />
                </Button>

                {isMenuOpen && (
                  <div className="absolute right-0 top-full mt-1 z-50 w-48 rounded-md border border-gray-200 bg-white shadow-lg">
                    <div className="py-1">
                      <button
                        type="button"
                        onClick={(e) => {
                          e.stopPropagation();
                          onEditDevice(device);
                          setOpenMenuId(null);
                        }}
                        className="flex w-full items-center gap-2 px-3 py-2 text-sm text-gray-700 hover:bg-gray-50"
                      >
                        <Edit className="h-4 w-4" />
                        Edit Device
                      </button>

                      <button
                        type="button"
                        onClick={(e) => {
                          e.stopPropagation();
                          onDeviceClick(device);
                          setOpenMenuId(null);
                        }}
                        className="flex w-full items-center gap-2 px-3 py-2 text-sm text-gray-700 hover:bg-gray-50"
                      >
                        <ExternalLink className="h-4 w-4" />
                        View Details
                      </button>

                      {availableTransitions.length > 0 && (
                        <>
                          <div className="border-t border-gray-100 my-1" />
                          <div className="px-3 py-1.5">
                            <span className="text-xs font-medium text-gray-500 uppercase">
                              Change Status
                            </span>
                          </div>
                          {availableTransitions.map((newStatus) => (
                            <button
                              key={newStatus}
                              type="button"
                              onClick={(e) => {
                                e.stopPropagation();
                                onStatusChange(device, newStatus);
                                setOpenMenuId(null);
                              }}
                              className="flex w-full items-center gap-2 px-3 py-2 text-sm text-gray-700 hover:bg-gray-50"
                            >
                              <ArrowRightLeft className="h-4 w-4" />
                              {formatStatus(newStatus)}
                            </button>
                          ))}
                        </>
                      )}

                      <div className="border-t border-gray-100 my-1" />
                      <button
                        type="button"
                        onClick={(e) => {
                          e.stopPropagation();
                          onDeleteDevice(device);
                          setOpenMenuId(null);
                        }}
                        className="flex w-full items-center gap-2 px-3 py-2 text-sm text-red-600 hover:bg-red-50"
                      >
                        <Trash2 className="h-4 w-4" />
                        Delete Device
                      </button>
                    </div>
                  </div>
                )}
              </div>
            </div>
          );
        },
        size: 80,
      },
    ],
    [openMenuId, onDeviceClick, onEditDevice, onDeleteDevice, onStatusChange]
  );

  return (
    <DataTable
      columns={columns}
      data={devices}
      isLoading={isLoading}
      showRowSelection
      showColumnVisibility
      selectedRows={rowSelection}
      onRowSelectionChange={handleRowSelectionChange}
      onRowClick={onDeviceClick}
      emptyMessage="No devices found. Try adjusting your filters or add new devices."
      searchPlaceholder="Filter devices..."
    />
  );
}
