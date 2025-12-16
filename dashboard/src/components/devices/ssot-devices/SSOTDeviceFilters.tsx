import { Search, X } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Select } from '@/components/ui/select';

interface FilterOption {
  value: string;
  label: string;
}

interface SSOTDeviceFiltersProps {
  searchQuery: string;
  onSearchChange: (value: string) => void;
  statusFilter: string;
  onStatusChange: (value: string) => void;
  schoolFilter: string;
  onSchoolChange: (value: string) => void;
  schoolOptions: FilterOption[];
  onClearFilters: () => void;
  hasFilters: boolean;
}

// Status options for filter
const STATUS_OPTIONS: FilterOption[] = [
  { value: '', label: 'All Statuses' },
  { value: 'in_stock', label: 'In Stock' },
  { value: 'deployed', label: 'Deployed' },
  { value: 'repair', label: 'In Repair' },
  { value: 'retired', label: 'Retired' },
];

export function SSOTDeviceFilters({
  searchQuery,
  onSearchChange,
  statusFilter,
  onStatusChange,
  schoolFilter,
  onSchoolChange,
  schoolOptions,
  onClearFilters,
  hasFilters,
}: SSOTDeviceFiltersProps) {
  return (
    <Card>
      <CardContent className="pt-6">
        <div className="flex flex-wrap gap-4 items-end">
          {/* Search */}
          <div className="flex-1 min-w-[200px]">
            <label className="text-sm font-medium text-gray-700 mb-1.5 block">
              Search
            </label>
            <div className="relative">
              <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-400" />
              <Input
                placeholder="Search by serial, asset tag, model..."
                value={searchQuery}
                onChange={(e) => onSearchChange(e.target.value)}
                className="pl-9"
              />
            </div>
          </div>

          {/* Status Filter */}
          <div className="w-[180px]">
            <label className="text-sm font-medium text-gray-700 mb-1.5 block">
              Status
            </label>
            <Select
              value={statusFilter}
              onChange={onStatusChange}
              options={STATUS_OPTIONS}
            />
          </div>

          {/* School Filter */}
          <div className="w-[200px]">
            <label className="text-sm font-medium text-gray-700 mb-1.5 block">
              School
            </label>
            <Select
              value={schoolFilter}
              onChange={onSchoolChange}
              options={schoolOptions}
            />
          </div>

          {/* Clear Filters */}
          {hasFilters && (
            <Button variant="ghost" onClick={onClearFilters} className="gap-2">
              <X className="h-4 w-4" />
              Clear
            </Button>
          )}
        </div>
      </CardContent>
    </Card>
  );
}
