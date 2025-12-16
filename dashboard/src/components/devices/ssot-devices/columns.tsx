import { type ColumnDef } from '@tanstack/react-table';
import {
  Laptop,
  Package,
  Wrench,
  Archive,
} from 'lucide-react';
import { Badge } from '@/components/ui/badge';
import { SortableHeader } from '@/components/ui/data-table';
import { cn } from '@/lib/utils';
import type { SSOTDevice } from '@/types/device';

// Status colors for badges
export const STATUS_COLORS: Record<string, string> = {
  in_stock: 'bg-green-100 text-green-800',
  deployed: 'bg-blue-100 text-blue-800',
  repair: 'bg-yellow-100 text-yellow-800',
  retired: 'bg-gray-100 text-gray-800',
};

// Status icons
const STATUS_ICONS: Record<string, React.ReactNode> = {
  in_stock: <Package className="h-4 w-4" />,
  deployed: <Laptop className="h-4 w-4" />,
  repair: <Wrench className="h-4 w-4" />,
  retired: <Archive className="h-4 w-4" />,
};

export function formatStatus(status: string): string {
  return status.replace(/_/g, ' ').replace(/\b\w/g, l => l.toUpperCase());
}

export const ssotDeviceColumns: ColumnDef<SSOTDevice>[] = [
  {
    accessorKey: 'serial',
    header: ({ column }) => <SortableHeader column={column}>Serial / Asset Tag</SortableHeader>,
    cell: ({ row }) => (
      <div>
        <div className="font-mono text-sm font-medium text-gray-900">
          {row.original.serial}
        </div>
        {row.original.assetTag && (
          <div className="text-xs text-gray-500">
            {row.original.assetTag}
          </div>
        )}
      </div>
    ),
  },
  {
    accessorKey: 'model',
    header: ({ column }) => <SortableHeader column={column}>Model</SortableHeader>,
    cell: ({ row }) => (
      <span className="text-sm text-gray-900">{row.original.model || '-'}</span>
    ),
  },
  {
    accessorKey: 'schoolId',
    header: 'School',
    cell: ({ row }) => (
      <span className="text-sm text-gray-600">{row.original.schoolId}</span>
    ),
  },
  {
    accessorKey: 'status',
    header: ({ column }) => <SortableHeader column={column}>Status</SortableHeader>,
    cell: ({ row }) => (
      <Badge className={cn('capitalize gap-1', STATUS_COLORS[row.original.status] || 'bg-gray-100 text-gray-800')}>
        {STATUS_ICONS[row.original.status]}
        {formatStatus(row.original.status)}
      </Badge>
    ),
  },
  {
    accessorKey: 'updatedAt',
    header: ({ column }) => <SortableHeader column={column}>Updated</SortableHeader>,
    cell: ({ row }) => (
      <span className="text-sm text-gray-500">
        {new Date(row.original.updatedAt).toLocaleDateString()}
      </span>
    ),
  },
];
