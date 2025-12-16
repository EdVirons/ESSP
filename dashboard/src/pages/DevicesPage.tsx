import * as React from 'react';
import { RefreshCw, Laptop, Filter } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { DataTable } from '@/components/ui/data-table';
import {
  useDevices,
  useDeviceStats,
  useDeviceModels,
  useDeviceMakes,
  useSyncDevicesFromSSO,
} from '@/api/devices';
import { useSchools } from '@/api/ssot';
import { toast } from '@/lib/toast';
import { cn } from '@/lib/utils';
import {
  ssotDeviceColumns,
  SSOTDeviceStats,
  SSOTDeviceFilters,
  ModelsSummary,
} from '@/components/devices/ssot-devices';

export function DevicesPage() {
  // Filters
  const [searchQuery, setSearchQuery] = React.useState('');
  const [statusFilter, setStatusFilter] = React.useState<string>('');
  const [schoolFilter, setSchoolFilter] = React.useState<string>('');
  const [page, setPage] = React.useState(0);
  const limit = 50;

  // Build filters object
  const filters = React.useMemo(() => ({
    q: searchQuery || undefined,
    status: statusFilter || undefined,
    schoolId: schoolFilter || undefined,
    limit,
    offset: page * limit,
  }), [searchQuery, statusFilter, schoolFilter, page]);

  // API queries
  const { data: devicesData, isLoading: devicesLoading, refetch } = useDevices(filters);
  const { data: statsData, isLoading: statsLoading } = useDeviceStats();
  const { data: modelsData } = useDeviceModels();
  const { data: makesData } = useDeviceMakes();
  const { data: schoolsData } = useSchools({ limit: 1000 });

  // Mutations
  const syncDevices = useSyncDevicesFromSSO();

  // Derived data
  const devices = devicesData?.items || [];
  const totalDevices = devicesData?.total || 0;
  const models = modelsData?.items || [];
  const makes = makesData?.makes || [];
  const schools = schoolsData?.items || [];

  // School options for filter
  const schoolOptions = React.useMemo(() => {
    return [
      { value: '', label: 'All Schools' },
      ...schools.map((s) => ({
        value: s.schoolId,
        label: s.name,
      })),
    ];
  }, [schools]);

  // Handle sync
  const handleSync = async () => {
    try {
      const result = await syncDevices.mutateAsync();
      toast.success(`Synced ${result.synced} devices from SSOT`);
      refetch();
    } catch {
      toast.error('Failed to sync devices');
    }
  };

  // Clear all filters
  const clearFilters = () => {
    setSearchQuery('');
    setStatusFilter('');
    setSchoolFilter('');
    setPage(0);
  };

  const hasFilters = searchQuery || statusFilter || schoolFilter;

  // Pagination
  const totalPages = Math.ceil(totalDevices / limit);

  // Handle pagination changes
  const handlePaginationChange = (pageIndex: number) => {
    setPage(pageIndex);
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-semibold text-gray-900">SSOT Devices</h1>
          <p className="text-sm text-gray-500 mt-1">
            Device inventory from Single Source of Truth
          </p>
        </div>
        <div className="flex items-center gap-3">
          <Button
            variant="outline"
            onClick={handleSync}
            disabled={syncDevices.isPending}
            className="gap-2"
          >
            <RefreshCw className={cn('h-4 w-4', syncDevices.isPending && 'animate-spin')} />
            Sync from SSOT
          </Button>
        </div>
      </div>

      {/* Stats Cards */}
      {!statsLoading && statsData && (
        <SSOTDeviceStats stats={statsData} />
      )}

      {/* Filters */}
      <SSOTDeviceFilters
        searchQuery={searchQuery}
        onSearchChange={(value) => {
          setSearchQuery(value);
          setPage(0);
        }}
        statusFilter={statusFilter}
        onStatusChange={(value) => {
          setStatusFilter(value);
          setPage(0);
        }}
        schoolFilter={schoolFilter}
        onSchoolChange={(value) => {
          setSchoolFilter(value);
          setPage(0);
        }}
        schoolOptions={schoolOptions}
        onClearFilters={clearFilters}
        hasFilters={!!hasFilters}
      />

      {/* Devices Table */}
      <Card>
        <CardHeader className="flex flex-row items-center justify-between">
          <CardTitle>
            Devices {totalDevices > 0 && `(${totalDevices})`}
          </CardTitle>
          {hasFilters && (
            <div className="flex items-center gap-2 text-sm text-muted-foreground">
              <Filter className="h-4 w-4" />
              Filtered
            </div>
          )}
        </CardHeader>
        <CardContent>
          {devices.length === 0 && !devicesLoading ? (
            <div className="text-center py-12">
              <Laptop className="h-12 w-12 text-gray-300 mx-auto mb-3" />
              <p className="text-gray-500">
                {hasFilters ? 'No devices match your filters' : 'No devices found'}
              </p>
              {hasFilters && (
                <Button variant="link" onClick={clearFilters} className="mt-2">
                  Clear filters
                </Button>
              )}
            </div>
          ) : (
            <>
              <DataTable
                columns={ssotDeviceColumns}
                data={devices}
                isLoading={devicesLoading}
                emptyMessage="No devices found"
                pageSize={limit}
                pageCount={totalPages}
                manualPagination
                onPaginationChange={(pageIndex) => handlePaginationChange(pageIndex)}
              />

              {/* Custom pagination info */}
              {totalPages > 1 && (
                <div className="flex items-center justify-between pt-4 border-t mt-4">
                  <div className="text-sm text-gray-500">
                    Showing {page * limit + 1} to {Math.min((page + 1) * limit, totalDevices)} of {totalDevices}
                  </div>
                  <div className="flex gap-2">
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => setPage(p => Math.max(0, p - 1))}
                      disabled={page === 0}
                    >
                      Previous
                    </Button>
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => setPage(p => Math.min(totalPages - 1, p + 1))}
                      disabled={page >= totalPages - 1}
                    >
                      Next
                    </Button>
                  </div>
                </div>
              )}
            </>
          )}
        </CardContent>
      </Card>

      {/* Models Summary */}
      <ModelsSummary models={models} makes={makes} />
    </div>
  );
}
