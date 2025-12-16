import * as React from 'react';
import { Search, Filter } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardContent } from '@/components/ui/card';
import { Select } from '@/components/ui/select';
import { DateRangePicker } from '@/components/ui/date-picker';
import type { AuditLogFilters as AuditLogFiltersType } from '@/types';
import { actionOptions } from './columns';

interface AuditLogFiltersProps {
  filters: AuditLogFiltersType;
  onFilterChange: (key: keyof AuditLogFiltersType, value: string) => void;
  searchQuery: string;
  onSearchChange: (query: string) => void;
  entityTypeOptions: Array<{ value: string; label: string }>;
  startDate: Date | null;
  endDate: Date | null;
  onDateChange: (start: Date | null, end: Date | null) => void;
  onClearFilters: () => void;
}

export function AuditLogFilters({
  filters,
  onFilterChange,
  searchQuery,
  onSearchChange,
  entityTypeOptions,
  startDate,
  endDate,
  onDateChange,
  onClearFilters,
}: AuditLogFiltersProps) {
  const [showFilters, setShowFilters] = React.useState(false);

  return (
    <Card className="border-0 shadow-md">
      <CardContent className="p-4">
        <div className="flex flex-wrap items-center gap-4">
          <div className="relative flex-1 min-w-[200px] max-w-md">
            <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-400" />
            <Input
              placeholder="Search by user or entity ID..."
              value={searchQuery}
              onChange={(e) => onSearchChange(e.target.value)}
              className="pl-9 border-gray-200 focus:border-slate-400 focus:ring-slate-400"
            />
          </div>
          <Select
            value={filters.action || ''}
            onChange={(value) => onFilterChange('action', value)}
            options={actionOptions}
            placeholder="Action"
            className="w-36"
          />
          <Select
            value={filters.entityType || ''}
            onChange={(value) => onFilterChange('entityType', value)}
            options={entityTypeOptions}
            placeholder="Entity Type"
            className="w-44"
          />
          <Button
            variant="outline"
            onClick={() => setShowFilters(!showFilters)}
            className={showFilters ? 'bg-slate-100 border-slate-300' : ''}
          >
            <Filter className="h-4 w-4" />
            More Filters
          </Button>
        </div>

        {showFilters && (
          <div className="mt-4 pt-4 border-t border-gray-100">
            <div className="flex flex-wrap items-center gap-4">
              <div className="flex-1 min-w-[300px]">
                <label className="block text-sm font-semibold text-gray-700 mb-1.5">
                  Date Range
                </label>
                <DateRangePicker
                  startDate={startDate}
                  endDate={endDate}
                  onChange={(start, end) => onDateChange(start, end)}
                />
              </div>
              <div>
                <label className="block text-sm font-semibold text-gray-700 mb-1.5">User ID</label>
                <Input
                  value={filters.userId || ''}
                  onChange={(e) => onFilterChange('userId', e.target.value)}
                  placeholder="Filter by user ID"
                  className="w-48 border-gray-200"
                />
              </div>
              <div className="flex items-end gap-2">
                <Button variant="outline" onClick={onClearFilters} className="text-slate-600">
                  Clear Filters
                </Button>
              </div>
            </div>
          </div>
        )}
      </CardContent>
    </Card>
  );
}
