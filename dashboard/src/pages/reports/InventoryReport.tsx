import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { ArrowLeft, Package, AlertTriangle, CheckCircle2, Boxes } from 'lucide-react';
import { useInventoryReport } from '@/api/reports';
import { ReportFilters, DataTable, ReportSummaryCards, type Column } from '@/components/reports';
import type { ReportFilters as Filters, InventoryReportItem } from '@/types/reports';

export function InventoryReport() {
  const navigate = useNavigate();
  const [filters, setFilters] = useState<Filters>({
    limit: 25,
    offset: 0,
    sortBy: 'partName',
    sortDir: 'asc',
  });

  const { data, isLoading } = useInventoryReport(filters);

  // Extract unique categories for filter
  const categoryOptions = data
    ? Object.keys(data.summary.byCategory).map((cat) => ({
        value: cat,
        label: cat,
      }))
    : [];

  const columns: Column<InventoryReportItem>[] = [
    {
      key: 'partSku',
      header: 'SKU',
      render: (item) => (
        <span className="font-mono text-xs">{item.partSku}</span>
      ),
    },
    { key: 'partName', header: 'Part Name', sortable: true },
    { key: 'category', header: 'Category', sortable: true },
    { key: 'serviceShopName', header: 'Service Shop' },
    {
      key: 'qtyAvailable',
      header: 'Available',
      sortable: true,
      render: (item) => (
        <span className={item.isLowStock ? 'text-amber-600 font-medium' : ''}>
          {item.qtyAvailable}
        </span>
      ),
    },
    { key: 'qtyReserved', header: 'Reserved' },
    { key: 'reorderThreshold', header: 'Reorder At' },
    {
      key: 'isLowStock',
      header: 'Stock Status',
      render: (item) => (
        item.isLowStock ? (
          <Badge variant="destructive" className="bg-amber-100 text-amber-700">
            Low Stock
          </Badge>
        ) : (
          <Badge className="bg-green-100 text-green-700">OK</Badge>
        )
      ),
    },
  ];

  const summaryCards = data
    ? [
        {
          title: 'Total Parts',
          value: data.summary.totalParts,
          icon: Package,
          color: 'info' as const,
        },
        {
          title: 'Total Quantity',
          value: data.summary.totalQtyAvailable,
          icon: Boxes,
          color: 'default' as const,
        },
        {
          title: 'Low Stock Items',
          value: data.summary.lowStockCount,
          icon: AlertTriangle,
          color: data.summary.lowStockCount > 0 ? 'warning' as const : 'success' as const,
        },
        {
          title: 'Stock Health',
          value:
            data.summary.totalParts > 0
              ? `${(((data.summary.totalParts - data.summary.lowStockCount) / data.summary.totalParts) * 100).toFixed(0)}%`
              : '100%',
          icon: CheckCircle2,
          color: 'success' as const,
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
          <h1 className="text-2xl font-bold text-gray-900">Inventory Report</h1>
          <p className="text-sm text-gray-500 mt-1">
            Monitor stock levels, low stock alerts, and parts distribution
          </p>
        </div>
      </div>

      {data && <ReportSummaryCards cards={summaryCards} />}

      <ReportFilters
        filters={filters}
        onFiltersChange={setFilters}
        showCategoryFilter
        categoryOptions={categoryOptions}
      />

      <Card>
        <CardHeader>
          <CardTitle>Inventory Items</CardTitle>
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
            keyField="partId"
            emptyMessage="No inventory items found"
          />
        </CardContent>
      </Card>

      {/* Category Distribution */}
      {data && Object.keys(data.summary.byCategory).length > 0 && (
        <Card>
          <CardHeader>
            <CardTitle>Parts by Category</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-4">
              {Object.entries(data.summary.byCategory).map(([category, count]) => (
                <div key={category} className="text-center p-4 bg-gray-50 rounded-lg">
                  <p className="text-sm text-gray-500 mb-1">{category}</p>
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
