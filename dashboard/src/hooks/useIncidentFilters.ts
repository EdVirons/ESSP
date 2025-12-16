import * as React from 'react';
import { useSearchParams } from 'react-router-dom';
import type { IncidentFilters, IncidentStatus } from '@/types';

export function useIncidentFilters(initialFilters: IncidentFilters = { limit: 50 }) {
  const [searchParams, setSearchParams] = useSearchParams();

  // Initialize filters from URL params
  const getInitialFilters = React.useCallback((): IncidentFilters => {
    const statusParam = searchParams.get('status');
    return {
      ...initialFilters,
      status: statusParam as IncidentStatus | undefined,
    };
  }, [searchParams, initialFilters]);

  const [filters, setFilters] = React.useState<IncidentFilters>(getInitialFilters);
  const [searchQuery, setSearchQuery] = React.useState('');

  // Sync filters when URL changes (e.g., clicking nav links)
  React.useEffect(() => {
    const statusParam = searchParams.get('status');
    setFilters(prev => ({
      ...prev,
      status: statusParam as IncidentStatus | undefined,
    }));
  }, [searchParams]);

  const handleFilterChange = React.useCallback(
    (key: keyof IncidentFilters, value: string) => {
      setFilters((prev) => ({
        ...prev,
        [key]: value || undefined,
      }));
      // Update URL when status filter changes
      if (key === 'status') {
        if (value) {
          searchParams.set('status', value);
        } else {
          searchParams.delete('status');
        }
        setSearchParams(searchParams, { replace: true });
      }
    },
    [searchParams, setSearchParams]
  );

  const resetFilters = React.useCallback(() => {
    setFilters(initialFilters);
    setSearchQuery('');
    // Clear URL params
    searchParams.delete('status');
    setSearchParams(searchParams, { replace: true });
  }, [initialFilters, searchParams, setSearchParams]);

  return {
    filters,
    searchQuery,
    setSearchQuery,
    handleFilterChange,
    resetFilters,
  };
}
