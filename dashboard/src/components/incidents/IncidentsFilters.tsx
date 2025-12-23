import * as React from 'react';
import { Search, Filter } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardContent } from '@/components/ui/card';
import { Select } from '@/components/ui/select';
import { DateRangePicker } from '@/components/ui/date-picker';
import type { IncidentFilters } from '@/types';

const severityOptions = [
  { value: '', label: 'All Severities' },
  { value: 'low', label: 'Low' },
  { value: 'medium', label: 'Medium' },
  { value: 'high', label: 'High' },
  { value: 'critical', label: 'Critical' },
];

const statusOptions = [
  { value: '', label: 'All Statuses' },
  { value: 'new', label: 'New' },
  { value: 'acknowledged', label: 'Acknowledged' },
  { value: 'in_progress', label: 'In Progress' },
  { value: 'escalated', label: 'Escalated' },
  { value: 'resolved', label: 'Resolved' },
  { value: 'closed', label: 'Closed' },
];

interface IncidentsFiltersProps {
  filters: IncidentFilters;
  searchQuery: string;
  onSearchChange: (value: string) => void;
  onFilterChange: (key: keyof IncidentFilters, value: string) => void;
  onClearFilters: () => void;
}

export function IncidentsFilters({
  filters,
  searchQuery,
  onSearchChange,
  onFilterChange,
  onClearFilters,
}: IncidentsFiltersProps) {
  const [showFilters, setShowFilters] = React.useState(false);
  const [startDate, setStartDate] = React.useState<Date | null>(null);
  const [endDate, setEndDate] = React.useState<Date | null>(null);

  return (
    <Card>
      <CardContent className="p-4">
        <div className="flex flex-col gap-3 sm:flex-row sm:flex-wrap sm:items-center sm:gap-4">
          <div className="relative w-full sm:flex-1 sm:min-w-[200px] sm:max-w-md">
            <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-400" />
            <Input
              placeholder="Search incidents..."
              value={searchQuery}
              onChange={(e) => onSearchChange(e.target.value)}
              className="pl-9"
            />
          </div>
          <Select
            value={filters.status || ''}
            onChange={(value) => onFilterChange('status', value)}
            options={statusOptions}
            placeholder="Status"
            className="w-40"
          />
          <Select
            value={filters.severity || ''}
            onChange={(value) => onFilterChange('severity', value)}
            options={severityOptions}
            placeholder="Severity"
            className="w-40"
          />
          <Button variant="outline" onClick={() => setShowFilters(!showFilters)}>
            <Filter className="h-4 w-4" />
            More Filters
          </Button>
        </div>

        {showFilters && (
          <div className="mt-4 pt-4 border-t border-gray-200">
            <div className="flex flex-col gap-4 sm:flex-row sm:flex-wrap sm:items-center">
              <div className="w-full sm:flex-1 sm:min-w-[300px]">
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Date Range
                </label>
                <DateRangePicker
                  startDate={startDate}
                  endDate={endDate}
                  onChange={(start, end) => {
                    setStartDate(start);
                    setEndDate(end);
                  }}
                />
              </div>
              <div className="flex items-end gap-2">
                <Button
                  variant="outline"
                  onClick={() => {
                    onClearFilters();
                    setStartDate(null);
                    setEndDate(null);
                  }}
                >
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
