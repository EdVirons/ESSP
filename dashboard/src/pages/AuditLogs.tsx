import * as React from 'react';
import { Download, FileText, Loader2, AlertCircle, RefreshCw } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import { DataTable } from '@/components/ui/data-table';
import {
  createAuditLogColumns,
  AuditLogStats,
  AuditLogDetail,
  AuditLogFilters,
} from '@/components/audit-logs';
import { useAuditLogs, useAuditLogEntityTypes, exportAuditLogs } from '@/api/audit-logs';
import type { AuditLog, AuditLogFilters as AuditLogFiltersType } from '@/types';

export function AuditLogs() {
  // Filters state
  const [filters, setFilters] = React.useState<AuditLogFiltersType>({
    limit: 50,
  });
  const [searchQuery, setSearchQuery] = React.useState('');
  const [startDate, setStartDate] = React.useState<Date | null>(null);
  const [endDate, setEndDate] = React.useState<Date | null>(null);
  const [isExporting, setIsExporting] = React.useState(false);

  // Selected log for detail view
  const [selectedLog, setSelectedLog] = React.useState<AuditLog | null>(null);
  const [showDetail, setShowDetail] = React.useState(false);

  // API hooks
  const { data, isLoading, error, refetch } = useAuditLogs({
    ...filters,
    startDate: startDate?.toISOString(),
    endDate: endDate?.toISOString(),
  });
  const { data: entityTypes } = useAuditLogEntityTypes();

  // Table columns
  const columns = React.useMemo(
    () =>
      createAuditLogColumns({
        onViewDetail: (log) => {
          setSelectedLog(log);
          setShowDetail(true);
        },
      }),
    []
  );

  // Handle filter changes
  const handleFilterChange = (key: keyof AuditLogFiltersType, value: string) => {
    setFilters((prev) => ({
      ...prev,
      [key]: value || undefined,
    }));
  };

  // Handle export
  const handleExport = async () => {
    setIsExporting(true);
    try {
      const blob = await exportAuditLogs({
        ...filters,
        startDate: startDate?.toISOString(),
        endDate: endDate?.toISOString(),
      });

      // Download the file
      const url = URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      link.download = `audit-logs-${new Date().toISOString().split('T')[0]}.csv`;
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      URL.revokeObjectURL(url);
    } catch (err) {
      console.error('Failed to export audit logs:', err);
    } finally {
      setIsExporting(false);
    }
  };

  // Handle clear filters
  const handleClearFilters = () => {
    setFilters({ limit: 50 });
    setSearchQuery('');
    setStartDate(null);
    setEndDate(null);
  };

  const logs = data?.items || [];
  const entityTypeOptions = [
    { value: '', label: 'All Entity Types' },
    ...(entityTypes || []),
  ];

  // Error state
  if (error) {
    return (
      <div className="space-y-6">
        {/* Page Header */}
        <div className="rounded-xl bg-gradient-to-r from-slate-700 via-slate-800 to-slate-900 p-6 text-white shadow-lg">
          <div className="flex items-center gap-4">
            <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-white/10 backdrop-blur">
              <FileText className="h-6 w-6" />
            </div>
            <div>
              <h1 className="text-2xl font-bold">Audit Logs</h1>
              <p className="text-slate-300">View system activity and change history</p>
            </div>
          </div>
        </div>
        <div className="flex flex-col items-center justify-center py-12">
          <AlertCircle className="h-12 w-12 text-red-500 mb-4" />
          <p className="text-red-600 text-center mb-4">
            {(error as Error).message || 'Failed to load audit logs'}
          </p>
          <Button onClick={() => refetch()}>
            <RefreshCw className="h-4 w-4 mr-2" />
            Try Again
          </Button>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Page Header */}
      <div className="rounded-xl bg-gradient-to-r from-slate-700 via-slate-800 to-slate-900 p-6 text-white shadow-lg">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-4">
            <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-white/10 backdrop-blur">
              <FileText className="h-6 w-6" />
            </div>
            <div>
              <h1 className="text-2xl font-bold">Audit Logs</h1>
              <p className="text-slate-300">View system activity and change history</p>
            </div>
          </div>
          <Button
            onClick={handleExport}
            disabled={isExporting}
            variant="outline"
            className="border-white/20 bg-white/10 text-white hover:bg-white/20 hover:text-white"
          >
            {isExporting ? (
              <Loader2 className="h-4 w-4 animate-spin" />
            ) : (
              <Download className="h-4 w-4" />
            )}
            {isExporting ? 'Exporting...' : 'Export CSV'}
          </Button>
        </div>
      </div>

      {/* Stats Cards */}
      <AuditLogStats logs={logs} />

      {/* Filters */}
      <AuditLogFilters
        filters={filters}
        onFilterChange={handleFilterChange}
        searchQuery={searchQuery}
        onSearchChange={setSearchQuery}
        entityTypeOptions={entityTypeOptions}
        startDate={startDate}
        endDate={endDate}
        onDateChange={(start, end) => {
          setStartDate(start);
          setEndDate(end);
        }}
        onClearFilters={handleClearFilters}
      />

      {/* Audit Logs Table */}
      <Card className="border-0 shadow-md overflow-hidden">
        <CardContent className="p-0">
          <DataTable
            columns={columns}
            data={logs}
            isLoading={isLoading}
            showColumnVisibility
            onRowClick={(row) => {
              setSelectedLog(row);
              setShowDetail(true);
            }}
            emptyMessage="No audit logs found"
          />
        </CardContent>
      </Card>

      {/* Audit Log Detail Sheet */}
      <AuditLogDetail
        log={selectedLog}
        open={showDetail}
        onClose={() => setShowDetail(false)}
      />
    </div>
  );
}
