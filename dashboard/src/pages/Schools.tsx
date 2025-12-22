import * as React from 'react';
import { type ColumnDef } from '@tanstack/react-table';
import {
  School,
  Search,
  MapPin,
  RefreshCw,
  Building2,
  GraduationCap,
  Info,
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardContent } from '@/components/ui/card';
import { Select } from '@/components/ui/select';
import { DataTable, SortableHeader } from '@/components/ui/data-table';
import { Badge } from '@/components/ui/badge';
import { SchoolDetailModal } from '@/components/schools';
import { useSchools, useSyncSchools, useCounties, useSubCounties } from '@/api/ssot';
import type { SchoolSnapshot, SchoolFilters } from '@/api/ssot';

// School level options (Kenya CBC system)
const LEVEL_OPTIONS = [
  { value: '', label: 'All Levels' },
  { value: 'JSS', label: 'JSS (Junior Secondary)' },
  { value: 'SSS', label: 'SSS (Senior Secondary)' },
  { value: 'Primary', label: 'Primary' },
  { value: 'TVET', label: 'TVET' },
  { value: 'University', label: 'University' },
];

// School type options
const TYPE_OPTIONS = [
  { value: '', label: 'All Types' },
  { value: 'public', label: 'Public' },
  { value: 'private', label: 'Private' },
];

// Helper function to get level badge color
function getLevelBadgeVariant(
  level?: string
): 'default' | 'secondary' | 'outline' | 'destructive' {
  switch (level) {
    case 'JSS':
      return 'default';
    case 'SSS':
      return 'secondary';
    case 'Primary':
      return 'outline';
    case 'TVET':
    case 'University':
      return 'outline';
    default:
      return 'secondary';
  }
}

export function Schools() {
  const [filters, setFilters] = React.useState<SchoolFilters>({
    limit: 50,
    offset: 0,
  });
  const [searchQuery, setSearchQuery] = React.useState('');
  const [selectedSchool, setSelectedSchool] = React.useState<SchoolSnapshot | null>(null);

  const { data, isLoading, refetch } = useSchools(filters);
  const syncSchools = useSyncSchools();
  const { data: countiesData } = useCounties();
  const { data: subCountiesData } = useSubCounties(filters.countyCode);

  const handleSearch = () => {
    setFilters((prev) => ({ ...prev, q: searchQuery, offset: 0 }));
  };

  const handleSync = () => {
    syncSchools.mutate(undefined, {
      onSuccess: () => {
        refetch();
      },
    });
  };

  // Counties for filter dropdown (from API)
  const counties = React.useMemo(() => {
    const items = countiesData?.items;
    if (!items) return [];
    return items.map((c) => ({
      value: c.code,
      label: c.name,
    }));
  }, [countiesData?.items]);

  const columns: ColumnDef<SchoolSnapshot>[] = [
    {
      accessorKey: 'name',
      header: ({ column }) => <SortableHeader column={column}>School Name</SortableHeader>,
      cell: ({ row }) => (
        <div className="flex items-center gap-3">
          <div className="flex h-8 w-8 items-center justify-center rounded-full bg-blue-50">
            <School className="h-4 w-4 text-blue-600" />
          </div>
          <div className="min-w-0">
            <div className="font-medium text-gray-900 truncate max-w-[250px]">
              {row.original.name}
            </div>
            <div className="text-sm text-gray-500">
              {row.original.knecCode || row.original.schoolId}
            </div>
          </div>
        </div>
      ),
    },
    {
      accessorKey: 'level',
      header: ({ column }) => <SortableHeader column={column}>Level</SortableHeader>,
      cell: ({ row }) => (
        <div className="flex items-center gap-2">
          <GraduationCap className="h-4 w-4 text-gray-400" />
          {row.original.level ? (
            <Badge variant={getLevelBadgeVariant(row.original.level)}>
              {row.original.level}
            </Badge>
          ) : (
            <span className="text-gray-400">-</span>
          )}
        </div>
      ),
    },
    {
      accessorKey: 'type',
      header: 'Type',
      cell: ({ row }) => (
        <span
          className={`capitalize ${row.original.type === 'public' ? 'text-green-600' : 'text-blue-600'}`}
        >
          {row.original.type || '-'}
        </span>
      ),
    },
    {
      accessorKey: 'countyName',
      header: ({ column }) => <SortableHeader column={column}>County</SortableHeader>,
      cell: ({ row }) => (
        <div className="flex items-center gap-2">
          <MapPin className="h-4 w-4 text-gray-400" />
          <span>{row.original.countyName || '-'}</span>
        </div>
      ),
    },
    {
      accessorKey: 'subCountyName',
      header: 'Sub-County',
      cell: ({ row }) => (
        <div className="flex items-center gap-2">
          <Building2 className="h-4 w-4 text-gray-400" />
          <span>{row.original.subCountyName || '-'}</span>
        </div>
      ),
    },
    {
      id: 'actions',
      header: '',
      cell: ({ row }) => (
        <Button variant="ghost" size="sm" onClick={() => setSelectedSchool(row.original)}>
          <Info className="h-4 w-4" />
        </Button>
      ),
    },
  ];

  const items = data?.items || [];
  const total = data?.total || 0;
  const currentPage = Math.floor((filters.offset || 0) / (filters.limit || 50)) + 1;
  const totalPages = Math.ceil(total / (filters.limit || 50));

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Schools</h1>
          <p className="text-gray-500">Browse and search schools from the SSOT registry</p>
        </div>
        <Button onClick={handleSync} disabled={syncSchools.isPending}>
          <RefreshCw
            className={`mr-2 h-4 w-4 ${syncSchools.isPending ? 'animate-spin' : ''}`}
          />
          {syncSchools.isPending ? 'Syncing...' : 'Sync Schools'}
        </Button>
      </div>

      {/* Stats */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <Card>
          <CardContent className="pt-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-gray-500">Total Schools</p>
                <p className="text-2xl font-bold text-gray-900">{total.toLocaleString()}</p>
              </div>
              <div className="h-12 w-12 rounded-full bg-blue-50 flex items-center justify-center">
                <School className="h-6 w-6 text-blue-600" />
              </div>
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="pt-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-gray-500">Counties</p>
                <p className="text-2xl font-bold text-gray-900">{countiesData?.total || 0}</p>
              </div>
              <div className="h-12 w-12 rounded-full bg-green-50 flex items-center justify-center">
                <MapPin className="h-6 w-6 text-green-600" />
              </div>
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="pt-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-gray-500">Sub-Counties</p>
                <p className="text-2xl font-bold text-gray-900">
                  {subCountiesData?.total || 0}
                </p>
              </div>
              <div className="h-12 w-12 rounded-full bg-orange-50 flex items-center justify-center">
                <Building2 className="h-6 w-6 text-orange-600" />
              </div>
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="pt-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-gray-500">Showing</p>
                <p className="text-2xl font-bold text-gray-900">{items.length}</p>
              </div>
              <div className="h-12 w-12 rounded-full bg-purple-50 flex items-center justify-center">
                <GraduationCap className="h-6 w-6 text-purple-600" />
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Filters */}
      <div className="flex flex-col sm:flex-row gap-4 flex-wrap">
        <div className="flex-1 flex gap-2 min-w-[300px]">
          <div className="relative flex-1 max-w-md">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-gray-400" />
            <Input
              placeholder="Search by name, ID, or KNEC code..."
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
          value={filters.countyCode || ''}
          onChange={(value) =>
            setFilters((prev) => ({ ...prev, countyCode: value || undefined, offset: 0 }))
          }
          options={[{ value: '', label: 'All Counties' }, ...counties]}
          className="w-48"
        />
        <Select
          value={filters.level || ''}
          onChange={(value) =>
            setFilters((prev) => ({ ...prev, level: value || undefined, offset: 0 }))
          }
          options={LEVEL_OPTIONS}
          className="w-40"
        />
        <Select
          value={filters.type || ''}
          onChange={(value) =>
            setFilters((prev) => ({ ...prev, type: value || undefined, offset: 0 }))
          }
          options={TYPE_OPTIONS}
          className="w-36"
        />
      </div>

      {/* Data Table */}
      <DataTable columns={columns} data={items} isLoading={isLoading} emptyMessage="No schools found" />

      {/* Pagination */}
      {totalPages > 1 && (
        <div className="flex flex-col items-center gap-3 sm:flex-row sm:justify-between">
          <p className="text-sm text-gray-500 text-center sm:text-left">
            Showing {(filters.offset || 0) + 1} to{' '}
            {Math.min((filters.offset || 0) + (filters.limit || 50), total)} of{' '}
            {total.toLocaleString()} schools
          </p>
          <div className="flex items-center gap-2">
            <Button
              variant="outline"
              size="sm"
              disabled={currentPage === 1}
              onClick={() =>
                setFilters((prev) => ({
                  ...prev,
                  offset: Math.max(0, (prev.offset || 0) - (prev.limit || 50)),
                }))
              }
            >
              Previous
            </Button>
            <span className="flex items-center px-3 text-sm text-gray-600 min-w-[80px] justify-center">
              <span className="hidden sm:inline">Page </span>
              {currentPage}
              <span className="sm:hidden">/{totalPages}</span>
              <span className="hidden sm:inline"> of {totalPages}</span>
            </span>
            <Button
              variant="outline"
              size="sm"
              disabled={currentPage === totalPages}
              onClick={() =>
                setFilters((prev) => ({
                  ...prev,
                  offset: (prev.offset || 0) + (prev.limit || 50),
                }))
              }
            >
              Next
            </Button>
          </div>
        </div>
      )}

      {/* School Detail Modal */}
      {selectedSchool && (
        <SchoolDetailModal school={selectedSchool} onClose={() => setSelectedSchool(null)} />
      )}
    </div>
  );
}
