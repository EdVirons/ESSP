import { Search, Filter } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardContent } from '@/components/ui/card';
import { Select } from '@/components/ui/select';
import type { WorkOrderFilters } from '@/types';

const statusOptions = [
  { value: '', label: 'All Statuses' },
  { value: 'draft', label: 'Draft' },
  { value: 'assigned', label: 'Assigned' },
  { value: 'in_repair', label: 'In Repair' },
  { value: 'qa', label: 'QA' },
  { value: 'completed', label: 'Completed' },
  { value: 'approved', label: 'Approved' },
];

interface WorkOrdersFiltersProps {
  filters: WorkOrderFilters;
  searchQuery: string;
  onSearchChange: (value: string) => void;
  onFilterChange: (key: keyof WorkOrderFilters, value: string) => void;
}

export function WorkOrdersFilters({
  filters,
  searchQuery,
  onSearchChange,
  onFilterChange,
}: WorkOrdersFiltersProps) {
  return (
    <Card className="border-0 shadow-md">
      <CardContent className="p-4">
        <div className="flex flex-wrap items-center gap-4">
          <div className="relative flex-1 min-w-[200px] max-w-md">
            <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-400" />
            <Input
              placeholder="Search work orders..."
              value={searchQuery}
              onChange={(e) => onSearchChange(e.target.value)}
              className="pl-9 border-gray-200 focus:border-blue-400 focus:ring-blue-400"
            />
          </div>
          <Select
            value={filters.status || ''}
            onChange={(value) => onFilterChange('status', value)}
            options={statusOptions}
            placeholder="Status"
            className="w-40"
          />
          <Button variant="outline" className="border-gray-200 hover:bg-blue-50 hover:text-blue-700 hover:border-blue-200">
            <Filter className="h-4 w-4" />
            More Filters
          </Button>
        </div>
      </CardContent>
    </Card>
  );
}
