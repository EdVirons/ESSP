import * as React from 'react';
import {
  Search,
  Filter,
  X,
  RotateCcw,
  ChevronDown,
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardContent } from '@/components/ui/card';
import { Select } from '@/components/ui/select';
import { Badge } from '@/components/ui/badge';
import { cn } from '@/lib/utils';
import type { DeviceFilters as DeviceFiltersType } from '@/types/device';
import {
  LIFECYCLE_STATUS_OPTIONS,
  ENROLLMENT_STATUS_OPTIONS,
  DEVICE_CATEGORY_OPTIONS,
} from '@/types/device';

interface DeviceFiltersProps {
  filters: DeviceFiltersType;
  searchQuery: string;
  onSearchChange: (value: string) => void;
  onSearch: () => void;
  onFilterChange: <K extends keyof DeviceFiltersType>(key: K, value: DeviceFiltersType[K]) => void;
  onReset: () => void;
  hasActiveFilters: boolean;
  activeFilterCount: number;
  schools?: Array<{ value: string; label: string }>;
  models?: Array<{ value: string; label: string }>;
  makes?: string[];
}

export function DeviceFilters({
  filters,
  searchQuery,
  onSearchChange,
  onSearch,
  onFilterChange,
  onReset,
  hasActiveFilters,
  activeFilterCount,
  schools = [],
  models = [],
  makes = [],
}: DeviceFiltersProps) {
  const [showAdvanced, setShowAdvanced] = React.useState(false);

  // Build school options
  const schoolOptions = React.useMemo(() => [
    { value: '', label: 'All Schools' },
    ...schools,
  ], [schools]);

  // Build model options
  const modelOptions = React.useMemo(() => [
    { value: '', label: 'All Models' },
    ...models,
  ], [models]);

  // Build make options
  const makeOptions = React.useMemo(() => [
    { value: '', label: 'All Makes' },
    ...makes.map((make) => ({ value: make, label: make })),
  ], [makes]);

  // Status options with "All" prefix
  const lifecycleOptions = [
    { value: '', label: 'All Status' },
    ...LIFECYCLE_STATUS_OPTIONS,
  ];

  const enrollmentOptions = [
    { value: '', label: 'All Enrollment' },
    ...ENROLLMENT_STATUS_OPTIONS,
  ];

  const categoryOptions = [
    { value: '', label: 'All Categories' },
    ...DEVICE_CATEGORY_OPTIONS,
  ];

  // Handle keyboard search
  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') {
      onSearch();
    }
  };

  // Get active filter labels for chips
  const activeFilters = React.useMemo(() => {
    const active: Array<{ key: keyof DeviceFiltersType; label: string; value: string }> = [];

    if (filters.lifecycle) {
      const option = LIFECYCLE_STATUS_OPTIONS.find((o) => o.value === filters.lifecycle);
      if (option) active.push({ key: 'lifecycle', label: 'Status', value: option.label });
    }
    if (filters.enrolled) {
      const option = ENROLLMENT_STATUS_OPTIONS.find((o) => o.value === filters.enrolled);
      if (option) active.push({ key: 'enrolled', label: 'Enrollment', value: option.label });
    }
    if (filters.schoolId) {
      const option = schools.find((o) => o.value === filters.schoolId);
      if (option) active.push({ key: 'schoolId', label: 'School', value: option.label });
    }
    if (filters.modelId) {
      const option = models.find((o) => o.value === filters.modelId);
      if (option) active.push({ key: 'modelId', label: 'Model', value: option.label });
    }
    if (filters.category) {
      const option = DEVICE_CATEGORY_OPTIONS.find((o) => o.value === filters.category);
      if (option) active.push({ key: 'category', label: 'Category', value: option.label });
    }
    if (filters.make) {
      active.push({ key: 'make', label: 'Make', value: filters.make });
    }

    return active;
  }, [filters, schools, models]);

  return (
    <Card>
      <CardContent className="p-4">
        <div className="space-y-4">
          {/* Primary filters row */}
          <div className="flex flex-wrap items-center gap-4">
            {/* Search input */}
            <div className="relative flex-1 min-w-[200px] max-w-md">
              <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-400" />
              <Input
                placeholder="Search by serial, asset tag, or model..."
                value={searchQuery}
                onChange={(e) => onSearchChange(e.target.value)}
                onKeyDown={handleKeyDown}
                className="pl-9 pr-9"
              />
              {searchQuery && (
                <button
                  type="button"
                  onClick={() => {
                    onSearchChange('');
                    onSearch();
                  }}
                  className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600"
                >
                  <X className="h-4 w-4" />
                </button>
              )}
            </div>
            <Button variant="outline" onClick={onSearch}>
              Search
            </Button>

            {/* Quick filters */}
            <Select
              value={filters.lifecycle || ''}
              onChange={(value) => onFilterChange('lifecycle', value as DeviceFiltersType['lifecycle'])}
              options={lifecycleOptions}
              className="w-36"
            />

            <Select
              value={filters.schoolId || ''}
              onChange={(value) => onFilterChange('schoolId', value || undefined)}
              options={schoolOptions}
              className="w-44"
            />

            {/* Advanced filters toggle */}
            <Button
              variant="outline"
              onClick={() => setShowAdvanced(!showAdvanced)}
              className={cn(
                'gap-2',
                showAdvanced && 'bg-gray-100'
              )}
            >
              <Filter className="h-4 w-4" />
              More Filters
              {activeFilterCount > 0 && (
                <Badge variant="default" className="ml-1 h-5 px-1.5">
                  {activeFilterCount}
                </Badge>
              )}
              <ChevronDown
                className={cn(
                  'h-4 w-4 transition-transform',
                  showAdvanced && 'rotate-180'
                )}
              />
            </Button>

            {/* Reset button */}
            {hasActiveFilters && (
              <Button variant="ghost" onClick={onReset} className="gap-2">
                <RotateCcw className="h-4 w-4" />
                Reset
              </Button>
            )}
          </div>

          {/* Advanced filters */}
          {showAdvanced && (
            <div className="flex flex-wrap items-center gap-4 pt-2 border-t border-gray-200">
              <Select
                value={filters.enrolled || ''}
                onChange={(value) => onFilterChange('enrolled', value as DeviceFiltersType['enrolled'])}
                options={enrollmentOptions}
                className="w-40"
              />

              <Select
                value={filters.category || ''}
                onChange={(value) => onFilterChange('category', value as DeviceFiltersType['category'])}
                options={categoryOptions}
                className="w-40"
              />

              <Select
                value={filters.make || ''}
                onChange={(value) => onFilterChange('make', value || undefined)}
                options={makeOptions}
                className="w-40"
              />

              <Select
                value={filters.modelId || ''}
                onChange={(value) => onFilterChange('modelId', value || undefined)}
                options={modelOptions}
                className="w-52"
              />
            </div>
          )}

          {/* Active filter chips */}
          {activeFilters.length > 0 && (
            <div className="flex flex-wrap items-center gap-2 pt-2 border-t border-gray-200">
              <span className="text-sm text-gray-500">Active filters:</span>
              {activeFilters.map((filter) => (
                <Badge
                  key={filter.key}
                  variant="secondary"
                  className="gap-1 pl-2 pr-1 py-1"
                >
                  <span className="text-gray-500">{filter.label}:</span>
                  <span>{filter.value}</span>
                  <button
                    type="button"
                    onClick={() => onFilterChange(filter.key, undefined as unknown as DeviceFiltersType[typeof filter.key])}
                    className="ml-1 hover:bg-gray-300 rounded p-0.5"
                  >
                    <X className="h-3 w-3" />
                  </button>
                </Badge>
              ))}
              <button
                type="button"
                onClick={onReset}
                className="text-sm text-blue-600 hover:text-blue-800"
              >
                Clear all
              </button>
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  );
}
