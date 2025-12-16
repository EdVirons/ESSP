import * as React from 'react';
import { type ColumnDef } from '@tanstack/react-table';
import {
  Laptop,
  Search,
  RefreshCw,
  Hash,
  Tag,
  CheckCircle2,
  XCircle,
  AlertCircle,
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Select } from '@/components/ui/select';
import { DataTable, SortableHeader } from '@/components/ui/data-table';
import { useDevices, useSyncDevices } from '@/api/ssot';
import type { DeviceSnapshot } from '@/api/ssot';
import { formatDate } from '@/lib/utils';

const statusOptions = [
  { value: '', label: 'All Status' },
  { value: 'active', label: 'Active' },
  { value: 'inactive', label: 'Inactive' },
  { value: 'repair', label: 'In Repair' },
  { value: 'disposed', label: 'Disposed' },
];

const statusColors: Record<string, string> = {
  active: 'bg-green-100 text-green-800',
  inactive: 'bg-gray-100 text-gray-800',
  repair: 'bg-yellow-100 text-yellow-800',
  disposed: 'bg-red-100 text-red-800',
};

const statusIcons: Record<string, React.ReactNode> = {
  active: <CheckCircle2 className="h-3 w-3" />,
  inactive: <XCircle className="h-3 w-3" />,
  repair: <AlertCircle className="h-3 w-3" />,
  disposed: <XCircle className="h-3 w-3" />,
};

export function Devices() {
  const [filters, setFilters] = React.useState<{
    q?: string;
    schoolId?: string;
    status?: string;
    limit?: number;
    offset?: number;
  }>({
    limit: 50,
    offset: 0,
  });
  const [searchQuery, setSearchQuery] = React.useState('');

  const { data, isLoading, refetch } = useDevices(filters);
  const syncDevices = useSyncDevices();

  const handleSearch = () => {
    setFilters((prev) => ({ ...prev, q: searchQuery, offset: 0 }));
  };

  const handleSync = () => {
    syncDevices.mutate(undefined, {
      onSuccess: () => {
        refetch();
      },
    });
  };

  const columns: ColumnDef<DeviceSnapshot>[] = [
    {
      accessorKey: 'serial',
      header: ({ column }) => <SortableHeader column={column}>Serial Number</SortableHeader>,
      cell: ({ row }) => (
        <div className="flex items-center gap-3">
          <div className="flex h-8 w-8 items-center justify-center rounded-full bg-indigo-50">
            <Laptop className="h-4 w-4 text-indigo-600" />
          </div>
          <div className="min-w-0">
            <div className="font-mono font-medium text-gray-900">
              {row.original.serial || '-'}
            </div>
            <div className="text-sm text-gray-500">{row.original.deviceId}</div>
          </div>
        </div>
      ),
    },
    {
      accessorKey: 'assetTag',
      header: ({ column }) => <SortableHeader column={column}>Asset Tag</SortableHeader>,
      cell: ({ row }) => (
        <div className="flex items-center gap-2">
          <Tag className="h-4 w-4 text-gray-400" />
          <span className="font-mono">{row.original.assetTag || '-'}</span>
        </div>
      ),
    },
    {
      accessorKey: 'model',
      header: 'Model',
      cell: ({ row }) => (
        <span className="text-gray-900">{row.original.model || '-'}</span>
      ),
    },
    {
      accessorKey: 'schoolId',
      header: 'School ID',
      cell: ({ row }) => (
        <div className="flex items-center gap-2">
          <Hash className="h-4 w-4 text-gray-400" />
          <span className="font-mono text-sm">{row.original.schoolId || '-'}</span>
        </div>
      ),
    },
    {
      accessorKey: 'status',
      header: 'Status',
      cell: ({ row }) => {
        const status = row.original.status || 'unknown';
        return (
          <Badge className={statusColors[status] || 'bg-gray-100 text-gray-800'}>
            <span className="flex items-center gap-1">
              {statusIcons[status]}
              {status.charAt(0).toUpperCase() + status.slice(1)}
            </span>
          </Badge>
        );
      },
    },
    {
      accessorKey: 'updatedAt',
      header: ({ column }) => <SortableHeader column={column}>Last Updated</SortableHeader>,
      cell: ({ row }) => (
        <span className="text-gray-600 text-sm">
          {formatDate(row.original.updatedAt)}
        </span>
      ),
    },
  ];

  const items = data?.items || [];
  const total = data?.total || 0;
  const currentPage = Math.floor((filters.offset || 0) / (filters.limit || 50)) + 1;
  const totalPages = Math.ceil(total / (filters.limit || 50));

  // Calculate status counts
  const statusCounts = React.useMemo(() => {
    const counts: Record<string, number> = { active: 0, inactive: 0, repair: 0, disposed: 0 };
    items.forEach((d) => {
      const status = d.status || 'unknown';
      if (status in counts) {
        counts[status]++;
      }
    });
    return counts;
  }, [items]);

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Devices</h1>
          <p className="text-gray-500">Browse and search devices from the SSOT inventory</p>
        </div>
        <Button onClick={handleSync} disabled={syncDevices.isPending}>
          <RefreshCw className={`mr-2 h-4 w-4 ${syncDevices.isPending ? 'animate-spin' : ''}`} />
          {syncDevices.isPending ? 'Syncing...' : 'Sync Devices'}
        </Button>
      </div>

      {/* Stats */}
      <div className="grid grid-cols-2 md:grid-cols-5 gap-4">
        <Card>
          <CardContent className="pt-6">
            <div className="text-center">
              <p className="text-sm text-gray-500">Total</p>
              <p className="text-2xl font-bold text-gray-900">{total}</p>
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="pt-6">
            <div className="text-center">
              <p className="text-sm text-green-600">Active</p>
              <p className="text-2xl font-bold text-green-600">{statusCounts.active}</p>
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="pt-6">
            <div className="text-center">
              <p className="text-sm text-gray-500">Inactive</p>
              <p className="text-2xl font-bold text-gray-600">{statusCounts.inactive}</p>
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="pt-6">
            <div className="text-center">
              <p className="text-sm text-yellow-600">In Repair</p>
              <p className="text-2xl font-bold text-yellow-600">{statusCounts.repair}</p>
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="pt-6">
            <div className="text-center">
              <p className="text-sm text-red-600">Disposed</p>
              <p className="text-2xl font-bold text-red-600">{statusCounts.disposed}</p>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Filters */}
      <div className="flex flex-col sm:flex-row gap-4">
        <div className="flex-1 flex gap-2">
          <div className="relative flex-1 max-w-md">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-gray-400" />
            <Input
              placeholder="Search by serial, asset tag, or model..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              onKeyDown={(e) => e.key === 'Enter' && handleSearch()}
              className="pl-9"
            />
          </div>
          <Button variant="outline" onClick={handleSearch}>
            Search
          </Button>
        </div>
        <Select
          value={filters.status || ''}
          onChange={(value) => setFilters((prev) => ({ ...prev, status: value || undefined, offset: 0 }))}
          options={statusOptions}
          className="w-40"
        />
      </div>

      {/* Data Table */}
      <DataTable
        columns={columns}
        data={items}
        isLoading={isLoading}
        emptyMessage="No devices found"
      />

      {/* Pagination */}
      {totalPages > 1 && (
        <div className="flex items-center justify-between">
          <p className="text-sm text-gray-500">
            Showing {(filters.offset || 0) + 1} to {Math.min((filters.offset || 0) + (filters.limit || 50), total)} of {total} devices
          </p>
          <div className="flex gap-2">
            <Button
              variant="outline"
              size="sm"
              disabled={currentPage === 1}
              onClick={() => setFilters((prev) => ({ ...prev, offset: Math.max(0, (prev.offset || 0) - (prev.limit || 50)) }))}
            >
              Previous
            </Button>
            <span className="flex items-center px-3 text-sm text-gray-600">
              Page {currentPage} of {totalPages}
            </span>
            <Button
              variant="outline"
              size="sm"
              disabled={currentPage === totalPages}
              onClick={() => setFilters((prev) => ({ ...prev, offset: (prev.offset || 0) + (prev.limit || 50) }))}
            >
              Next
            </Button>
          </div>
        </div>
      )}
    </div>
  );
}
