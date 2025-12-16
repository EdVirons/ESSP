import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { ArrowLeft, School, Laptop, AlertTriangle, Wrench, MapPin } from 'lucide-react';
import { useSchoolsReport } from '@/api/reports';
import { ReportFilters, DataTable, ReportSummaryCards, type Column } from '@/components/reports';
import type { ReportFilters as Filters, SchoolReportItem } from '@/types/reports';

export function SchoolsReport() {
  const navigate = useNavigate();
  const [filters, setFilters] = useState<Filters>({
    limit: 25,
    offset: 0,
    sortBy: 'schoolName',
    sortDir: 'asc',
  });

  const { data, isLoading } = useSchoolsReport(filters);

  // Extract unique counties for filter
  const countyOptions = data
    ? Object.keys(data.summary.byCounty).map((county) => ({
        value: county,
        label: county,
      }))
    : [];

  const columns: Column<SchoolReportItem>[] = [
    {
      key: 'schoolId',
      header: 'ID',
      render: (item) => (
        <span className="font-mono text-xs">{item.schoolId.slice(0, 8)}</span>
      ),
    },
    { key: 'schoolName', header: 'School Name', sortable: true },
    { key: 'countyName', header: 'County', sortable: true },
    {
      key: 'deviceCount',
      header: 'Devices',
      sortable: true,
      render: (item) => (
        <div className="flex items-center gap-1">
          <Laptop className="h-4 w-4 text-gray-400" />
          <span>{item.deviceCount}</span>
        </div>
      ),
    },
    {
      key: 'incidentCount',
      header: 'Incidents',
      sortable: true,
      render: (item) => (
        <div className="flex items-center gap-1">
          <AlertTriangle className="h-4 w-4 text-amber-500" />
          <span className={item.incidentCount > 5 ? 'text-amber-600 font-medium' : ''}>
            {item.incidentCount}
          </span>
        </div>
      ),
    },
    {
      key: 'workOrderCount',
      header: 'Work Orders',
      sortable: true,
      render: (item) => (
        <div className="flex items-center gap-1">
          <Wrench className="h-4 w-4 text-blue-500" />
          <span>{item.workOrderCount}</span>
        </div>
      ),
    },
  ];

  const summaryCards = data
    ? [
        {
          title: 'Total Schools',
          value: data.summary.totalSchools,
          icon: School,
          color: 'info' as const,
        },
        {
          title: 'Total Devices',
          value: data.summary.totalDevices,
          icon: Laptop,
          color: 'default' as const,
        },
        {
          title: 'Counties',
          value: Object.keys(data.summary.byCounty).length,
          icon: MapPin,
          color: 'default' as const,
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
          <h1 className="text-2xl font-bold text-gray-900">Schools Report</h1>
          <p className="text-sm text-gray-500 mt-1">
            View schools, device counts, and incident history by location
          </p>
        </div>
      </div>

      {data && <ReportSummaryCards cards={summaryCards} />}

      <ReportFilters
        filters={filters}
        onFiltersChange={setFilters}
        showCountyFilter
        countyOptions={countyOptions}
      />

      <Card>
        <CardHeader>
          <CardTitle>Schools</CardTitle>
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
            keyField="schoolId"
            emptyMessage="No schools found"
          />
        </CardContent>
      </Card>

      {/* County Distribution */}
      {data && Object.keys(data.summary.byCounty).length > 0 && (
        <Card>
          <CardHeader>
            <CardTitle>Schools by County</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-4">
              {Object.entries(data.summary.byCounty)
                .sort((a, b) => b[1] - a[1])
                .map(([county, count]) => (
                  <div key={county} className="text-center p-4 bg-gray-50 rounded-lg">
                    <p className="text-sm text-gray-500 mb-1 truncate">{county}</p>
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
