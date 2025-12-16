import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { ArrowLeft, AlertTriangle, CheckCircle2, Clock, ShieldCheck, XCircle } from 'lucide-react';
import { useIncidentsReport } from '@/api/reports';
import { ReportFilters, DataTable, ReportSummaryCards, type Column } from '@/components/reports';
import type { ReportFilters as Filters, IncidentReportItem } from '@/types/reports';

const statusOptions = [
  { value: 'new', label: 'New' },
  { value: 'acknowledged', label: 'Acknowledged' },
  { value: 'in_progress', label: 'In Progress' },
  { value: 'escalated', label: 'Escalated' },
  { value: 'resolved', label: 'Resolved' },
  { value: 'closed', label: 'Closed' },
];

const statusColors: Record<string, string> = {
  new: 'bg-blue-100 text-blue-700',
  acknowledged: 'bg-amber-100 text-amber-700',
  in_progress: 'bg-purple-100 text-purple-700',
  escalated: 'bg-red-100 text-red-700',
  resolved: 'bg-green-100 text-green-700',
  closed: 'bg-gray-100 text-gray-700',
};

const severityColors: Record<string, string> = {
  low: 'bg-green-100 text-green-700',
  medium: 'bg-amber-100 text-amber-700',
  high: 'bg-orange-100 text-orange-700',
  critical: 'bg-red-100 text-red-700',
};

export function IncidentsReport() {
  const navigate = useNavigate();
  const [filters, setFilters] = useState<Filters>({
    limit: 25,
    offset: 0,
    sortBy: 'createdAt',
    sortDir: 'desc',
  });

  const { data, isLoading } = useIncidentsReport(filters);

  const columns: Column<IncidentReportItem>[] = [
    {
      key: 'id',
      header: 'ID',
      render: (item) => (
        <span className="font-mono text-xs">{item.id.slice(0, 8)}</span>
      ),
    },
    {
      key: 'title',
      header: 'Title',
      render: (item) => (
        <span className="truncate max-w-[200px] block">{item.title}</span>
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
    {
      key: 'severity',
      header: 'Severity',
      sortable: true,
      render: (item) => (
        <Badge className={severityColors[item.severity] || 'bg-gray-100'}>
          {item.severity}
        </Badge>
      ),
    },
    { key: 'schoolName', header: 'School' },
    {
      key: 'slaBreached',
      header: 'SLA',
      render: (item) => (
        item.slaBreached ? (
          <Badge variant="destructive">Breached</Badge>
        ) : (
          <Badge className="bg-green-100 text-green-700">On Track</Badge>
        )
      ),
    },
    {
      key: 'createdAt',
      header: 'Created',
      sortable: true,
      render: (item) => new Date(item.createdAt).toLocaleDateString(),
    },
    {
      key: 'resolutionHours',
      header: 'Resolution Time',
      render: (item) =>
        item.resolutionHours ? `${item.resolutionHours.toFixed(1)}h` : '-',
    },
  ];

  const summaryCards = data
    ? [
        {
          title: 'Total Incidents',
          value: data.summary.total,
          icon: AlertTriangle,
          color: 'info' as const,
        },
        {
          title: 'Resolved',
          value: (data.summary.byStatus.resolved || 0) + (data.summary.byStatus.closed || 0),
          icon: CheckCircle2,
          color: 'success' as const,
        },
        {
          title: 'SLA Compliance',
          value: `${data.summary.slaComplianceRate.toFixed(1)}%`,
          icon: ShieldCheck,
          color: data.summary.slaComplianceRate >= 90 ? 'success' as const : 'warning' as const,
        },
        {
          title: 'Avg Resolution',
          value: `${data.summary.avgResolutionHours.toFixed(1)}h`,
          icon: Clock,
          color: 'default' as const,
        },
        {
          title: 'SLA Breached',
          value: data.summary.slaBreachedCount,
          icon: XCircle,
          color: data.summary.slaBreachedCount > 0 ? 'danger' as const : 'success' as const,
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
          <h1 className="text-2xl font-bold text-gray-900">Incidents Report</h1>
          <p className="text-sm text-gray-500 mt-1">
            Analyze incidents, SLA compliance, and resolution times
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
          <CardTitle>Incidents</CardTitle>
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
            onRowClick={(item) => navigate(`/incidents/${item.id}`)}
            emptyMessage="No incidents found"
          />
        </CardContent>
      </Card>

      {/* Severity Distribution */}
      {data && (
        <div className="grid gap-6 md:grid-cols-2">
          <Card>
            <CardHeader>
              <CardTitle>Status Distribution</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-3 gap-4">
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

          <Card>
            <CardHeader>
              <CardTitle>Severity Distribution</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-2 gap-4">
                {Object.entries(data.summary.bySeverity).map(([severity, count]) => (
                  <div key={severity} className="text-center">
                    <Badge className={`${severityColors[severity]} mb-2`}>
                      {severity}
                    </Badge>
                    <p className="text-2xl font-bold">{count}</p>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
        </div>
      )}
    </div>
  );
}
