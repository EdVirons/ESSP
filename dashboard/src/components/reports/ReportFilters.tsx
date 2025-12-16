import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Calendar, X } from 'lucide-react';
import type { ReportFilters as Filters } from '@/types/reports';

interface ReportFiltersProps {
  filters: Filters;
  onFiltersChange: (filters: Filters) => void;
  showStatusFilter?: boolean;
  statusOptions?: { value: string; label: string }[];
  showCategoryFilter?: boolean;
  categoryOptions?: { value: string; label: string }[];
  showCountyFilter?: boolean;
  countyOptions?: { value: string; label: string }[];
}

export function ReportFilters({
  filters,
  onFiltersChange,
  showStatusFilter = false,
  statusOptions = [],
  showCategoryFilter = false,
  categoryOptions = [],
  showCountyFilter = false,
  countyOptions = [],
}: ReportFiltersProps) {
  const hasActiveFilters =
    filters.dateFrom ||
    filters.dateTo ||
    (filters.status && filters.status.length > 0) ||
    filters.category ||
    filters.countyCode;

  const handleClearFilters = () => {
    onFiltersChange({
      ...filters,
      dateFrom: undefined,
      dateTo: undefined,
      status: undefined,
      category: undefined,
      countyCode: undefined,
    });
  };

  return (
    <div className="bg-white p-4 rounded-lg border space-y-4">
      <div className="flex items-center justify-between">
        <h3 className="font-medium text-gray-900">Filters</h3>
        {hasActiveFilters && (
          <Button variant="ghost" size="sm" onClick={handleClearFilters}>
            <X className="h-4 w-4 mr-1" />
            Clear Filters
          </Button>
        )}
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 xl:grid-cols-5 gap-4">
        {/* Date From */}
        <div className="space-y-1.5">
          <Label htmlFor="dateFrom" className="text-xs text-gray-500">
            From Date
          </Label>
          <div className="relative">
            <Calendar className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-gray-400" />
            <Input
              id="dateFrom"
              type="date"
              className="pl-10"
              value={filters.dateFrom || ''}
              onChange={(e) =>
                onFiltersChange({ ...filters, dateFrom: e.target.value || undefined })
              }
            />
          </div>
        </div>

        {/* Date To */}
        <div className="space-y-1.5">
          <Label htmlFor="dateTo" className="text-xs text-gray-500">
            To Date
          </Label>
          <div className="relative">
            <Calendar className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-gray-400" />
            <Input
              id="dateTo"
              type="date"
              className="pl-10"
              value={filters.dateTo || ''}
              onChange={(e) =>
                onFiltersChange({ ...filters, dateTo: e.target.value || undefined })
              }
            />
          </div>
        </div>

        {/* Status Filter */}
        {showStatusFilter && statusOptions.length > 0 && (
          <div className="space-y-1.5">
            <Label className="text-xs text-gray-500">Status</Label>
            <Select
              value={filters.status?.[0] || 'all'}
              onValueChange={(value) =>
                onFiltersChange({
                  ...filters,
                  status: value === 'all' ? undefined : [value],
                })
              }
            >
              <SelectTrigger>
                <SelectValue placeholder="All statuses" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All Statuses</SelectItem>
                {statusOptions.map((opt) => (
                  <SelectItem key={opt.value} value={opt.value}>
                    {opt.label}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
        )}

        {/* Category Filter */}
        {showCategoryFilter && categoryOptions.length > 0 && (
          <div className="space-y-1.5">
            <Label className="text-xs text-gray-500">Category</Label>
            <Select
              value={filters.category || 'all'}
              onValueChange={(value) =>
                onFiltersChange({
                  ...filters,
                  category: value === 'all' ? undefined : value,
                })
              }
            >
              <SelectTrigger>
                <SelectValue placeholder="All categories" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All Categories</SelectItem>
                {categoryOptions.map((opt) => (
                  <SelectItem key={opt.value} value={opt.value}>
                    {opt.label}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
        )}

        {/* County Filter */}
        {showCountyFilter && countyOptions.length > 0 && (
          <div className="space-y-1.5">
            <Label className="text-xs text-gray-500">County</Label>
            <Select
              value={filters.countyCode || 'all'}
              onValueChange={(value) =>
                onFiltersChange({
                  ...filters,
                  countyCode: value === 'all' ? undefined : value,
                })
              }
            >
              <SelectTrigger>
                <SelectValue placeholder="All counties" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All Counties</SelectItem>
                {countyOptions.map((opt) => (
                  <SelectItem key={opt.value} value={opt.value}>
                    {opt.label}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
        )}
      </div>
    </div>
  );
}
