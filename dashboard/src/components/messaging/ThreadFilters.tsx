import { Search, Filter, X } from 'lucide-react';
import { Input } from '@/components/ui/input';
import { Select } from '@/components/ui/select';

interface ThreadFiltersProps {
  searchQuery: string;
  onSearchChange: (query: string) => void;
  statusFilter: string;
  onStatusChange: (status: string) => void;
}

const statusOptions = [
  { value: 'all', label: 'All Conversations' },
  { value: 'open', label: 'Open' },
  { value: 'closed', label: 'Closed' },
];

export function ThreadFilters({
  searchQuery,
  onSearchChange,
  statusFilter,
  onStatusChange,
}: ThreadFiltersProps) {
  return (
    <div className="p-3 border-b border-gray-200 space-y-3">
      {/* Search */}
      <div className="relative">
        <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-gray-400" />
        <Input
          type="text"
          placeholder="Search conversations..."
          value={searchQuery}
          onChange={(e) => onSearchChange(e.target.value)}
          className="pl-9 pr-8"
        />
        {searchQuery && (
          <button
            onClick={() => onSearchChange('')}
            className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600"
          >
            <X className="h-4 w-4" />
          </button>
        )}
      </div>

      {/* Status filter */}
      <div className="flex items-center gap-2">
        <Filter className="h-4 w-4 text-gray-400" />
        <Select
          value={statusFilter}
          onChange={onStatusChange}
          options={statusOptions}
          placeholder="Filter by status"
          className="flex-1"
        />
      </div>
    </div>
  );
}
