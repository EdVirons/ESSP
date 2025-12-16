import * as React from 'react';
import type { WorkOrderFilters } from '@/types';

export function useWorkOrderFilters(initialFilters: WorkOrderFilters = { limit: 50 }) {
  const [filters, setFilters] = React.useState<WorkOrderFilters>(initialFilters);
  const [searchQuery, setSearchQuery] = React.useState('');

  const handleFilterChange = React.useCallback(
    (key: keyof WorkOrderFilters, value: string) => {
      setFilters((prev) => ({
        ...prev,
        [key]: value || undefined,
      }));
    },
    []
  );

  const resetFilters = React.useCallback(() => {
    setFilters(initialFilters);
    setSearchQuery('');
  }, [initialFilters]);

  return {
    filters,
    searchQuery,
    setSearchQuery,
    handleFilterChange,
    resetFilters,
  };
}
