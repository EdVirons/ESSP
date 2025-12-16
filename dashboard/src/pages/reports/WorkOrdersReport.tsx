import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { ArrowLeft, Wrench, CheckCircle2, Clock, DollarSign, RefreshCw } from 'lucide-react';
import { useWorkOrdersReport } from '@/api/reports';
import { ReportFilters, DataTable, ReportSummaryCards, type Column } from '@/components/reports';
import type { ReportFilters as Filters, WorkOrderReportItem } from '@/types/reports';

const statusOptions = [
  { value: 'draft', label: 'Draft' },
  { value: 'assigned', label: 'Assigned' },
  { value: 'in_repair', label: 'In Repair' },
  { value: 'qa', label: 'QA' },
  { value: 'completed', label: 'Completed' },
  { value: 'approved', label: 'Approved' },
];

const statusColors: Record<string, string> = {
  draft: 'bg-gray-100 text-gray-700',
  assigned: 'bg-blue-100 text-blue-700',
  in_repair: 'bg-amber-100 text-amber-700',
  qa: 'bg-purple-100 text-purple-700',
  completed: 'bg-green-100 text-green-700',
  approved: 'bg-teal-100 text-teal-700',
};

export function WorkOrdersReport() {
  const navigate = useNavigate();
  const [filters, setFilters] = useState<Filters>({
    limit: 25,
    offset: 0,
    sortBy: 'createdAt',
    sortDir: 'desc',
  });

  const { data, isLoading } = useWorkOrdersReport(filters);

  const columns: Column<WorkOrderReportItem>[] = [
    {
      key: 'id',
      header: 'ID',
      render: (item) => (
        <span className="font-mono text-xs">{item.id.slice(0, 8)}</span>
      ),
    },
    {
      key: 'status',
      header: 'Status',
      sortable: true,
      render: (item) => (
        <Badge className={statusColors[item.status] || 'bg-gray-100'}>
          {item.status.replace('_', ' ')}
        </Badge>
      ),
    },
    { key: 'taskType', header: 'Task Type', sortable: true },
    { key: 'schoolName', header: 'School', sortable: true },
    { key: 'deviceCategory', header: 'Device' },
    {
      key: 'costCents',
      header: 'Cost',
      sortable: true,
      render: (item) => `KES ${(item.costCents / 100).toLocaleString()}`,
    },
    {
      key: 'reworkCount',
      header: 'Reworks',
      render: (item) => (
        <span className={item.reworkCount > 0 ? 'text-amber-600 font-medium' : ''}>
          {item.reworkCount}
        </span>
      ),
    },
    {
      key: 'createdAt',
      header: 'Created',
      sortable: true,
      render: (item) => new Date(item.createdAt).toLocaleDateString(),
    },
    {
      key: 'durationHours',
      header: 'Duration',
      render: (item) =>
        item.durationHours ? `${item.durationHours.toFixed(1)}h` : '-',
    },
  ];

  const summaryCards = data
    ? [
        {
          title: 'Total Work Orders',
          value: data.summary.total,
          icon: Wrench,
          color: 'info' as const,
        },
        {
          title: 'Completed',
          value: (data.summary.byStatus.completed || 0) + (data.summary.byStatus.approved || 0),
          icon: CheckCircle2,
          color: 'success' as const,
        },
        {
          title: 'Avg Completion',
          value: `${data.summary.avgCompletionHours.toFixed(1)}h`,
          icon: Clock,
          color: 'default' as const,
        },
        {
          title: 'Total Cost',
          value: `KES ${(data.summary.totalCostCents / 100).toLocaleString()}`,
          icon: DollarSign,
          color: 'default' as const,
        },
        {
          title: 'Rework Rate',
          value: `${data.summary.reworkRate.toFixed(1)}%`,
          icon: RefreshCw,
          color: data.summary.reworkRate > 10 ? 'warning' as const : 'success' as const,
        },
      ]
    : [];

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-4">
        <Button variant="ghost" size="icon" onClick={() => navigate('/reports')}>
          <ArrowLeft className="h-5 w-5" />
        </Button>
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Work Orders Report</h1>
          <p className="text-sm text-gray-500 mt-1">
            Track work order status, completion rates, and repair timelines
          </p>
        </div>
      </div>

      {data && <ReportSummaryCards cards={summaryCards} />}

      <ReportFilters
        filters={filters}
        onFiltersChange={setFilters}
        showStatusFilter
        statusOptions={statusOptions}
      />

      <Card>
        <CardHeader>
          <CardTitle>Work Orders</CardTitle>
        </CardHeader>
        <CardContent>
          <DataTable
            data={data?.items || []}
            columns={columns}
            pagination={data?.pagination}
            onPageChange={(offset) => setFilters({ ...filters, offset })}
            onLimitChange={(limit) => setFilters({ ...filters, limit, offset: 0 })}
            onSort={(sortBy, sortDir) => setFilters({ ...filters, sortBy, sortDir })}
            sortBy={filters.sortBy}
            sortDir={filters.sortDir}
            isLoading={isLoading}
            onRowClick={(item) => navigate(`/work-orders/${item.id}`)}
            emptyMessage="No work orders found"
          />
        </CardContent>
      </Card>

      {/* Status Distribution */}
      {data && (
        <Card>
          <CardHeader>
            <CardTitle>Status Distribution</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-3 md:grid-cols-6 gap-4">
              {Object.entries(data.summary.byStatus).map(([status, count]) => (
                <div key={status} className="text-center">
                  <Badge className={`${statusColors[status]} mb-2`}>
                    {status.replace('_', ' ')}
                  </Badge>
                  <p className="text-2xl font-bold">{count}</p>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  );
}
