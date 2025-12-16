import * as React from 'react';
import { useSearchParams } from 'react-router-dom';
import type { DeviceFilters, DeviceLifecycleStatus, DeviceEnrollmentStatus, DeviceCategory } from '@/types/device';

const DEFAULT_FILTERS: DeviceFilters = {
  limit: 50,
  offset: 0,
  sortBy: 'createdAt',
  sortOrder: 'desc',
};

/**
 * Hook for managing device filters with URL synchronization
 */
export function useDeviceFilters(initialFilters: DeviceFilters = DEFAULT_FILTERS) {
  const [searchParams, setSearchParams] = useSearchParams();

  // Initialize filters from URL params or defaults
  const getInitialFilters = React.useCallback((): DeviceFilters => {
    const q = searchParams.get('q') || undefined;
    const schoolId = searchParams.get('schoolId') || undefined;
    const modelId = searchParams.get('modelId') || undefined;
    const lifecycle = searchParams.get('lifecycle') as DeviceLifecycleStatus | undefined;
    const enrolled = searchParams.get('enrolled') as DeviceEnrollmentStatus | undefined;
    const category = searchParams.get('category') as DeviceCategory | undefined;
    const make = searchParams.get('make') || undefined;
    const limit = searchParams.get('limit') ? parseInt(searchParams.get('limit')!) : initialFilters.limit;
    const offset = searchParams.get('offset') ? parseInt(searchParams.get('offset')!) : initialFilters.offset;
    const sortBy = searchParams.get('sortBy') as DeviceFilters['sortBy'] || initialFilters.sortBy;
    const sortOrder = searchParams.get('sortOrder') as DeviceFilters['sortOrder'] || initialFilters.sortOrder;

    return {
      q,
      schoolId,
      modelId,
      lifecycle,
      enrolled,
      category,
      make,
      limit,
      offset,
      sortBy,
      sortOrder,
    };
  }, [searchParams, initialFilters]);

  const [filters, setFiltersState] = React.useState<DeviceFilters>(getInitialFilters);
  const [searchQuery, setSearchQuery] = React.useState(filters.q || '');

  // Sync filters to URL
  const syncToUrl = React.useCallback((newFilters: DeviceFilters) => {
    const params = new URLSearchParams();

    if (newFilters.q) params.set('q', newFilters.q);
    if (newFilters.schoolId) params.set('schoolId', newFilters.schoolId);
    if (newFilters.modelId) params.set('modelId', newFilters.modelId);
    if (newFilters.lifecycle) params.set('lifecycle', newFilters.lifecycle);
    if (newFilters.enrolled) params.set('enrolled', newFilters.enrolled);
    if (newFilters.category) params.set('category', newFilters.category);
    if (newFilters.make) params.set('make', newFilters.make);
    if (newFilters.limit && newFilters.limit !== DEFAULT_FILTERS.limit) {
      params.set('limit', newFilters.limit.toString());
    }
    if (newFilters.offset && newFilters.offset !== 0) {
      params.set('offset', newFilters.offset.toString());
    }
    if (newFilters.sortBy && newFilters.sortBy !== DEFAULT_FILTERS.sortBy) {
      params.set('sortBy', newFilters.sortBy);
    }
    if (newFilters.sortOrder && newFilters.sortOrder !== DEFAULT_FILTERS.sortOrder) {
      params.set('sortOrder', newFilters.sortOrder);
    }

    setSearchParams(params, { replace: true });
  }, [setSearchParams]);

  // Update filters
  const setFilters = React.useCallback((
    updater: DeviceFilters | ((prev: DeviceFilters) => DeviceFilters)
  ) => {
    setFiltersState((prev) => {
      const newFilters = typeof updater === 'function' ? updater(prev) : updater;
      syncToUrl(newFilters);
      return newFilters;
    });
  }, [syncToUrl]);

  // Handle individual filter change
  const handleFilterChange = React.useCallback(
    <K extends keyof DeviceFilters>(key: K, value: DeviceFilters[K]) => {
      setFilters((prev) => ({
        ...prev,
        [key]: value || undefined,
        offset: 0, // Reset pagination when filter changes
      }));
    },
    [setFilters]
  );

  // Handle search (triggered explicitly)
  const handleSearch = React.useCallback(() => {
    setFilters((prev) => ({
      ...prev,
      q: searchQuery || undefined,
      offset: 0,
    }));
  }, [searchQuery, setFilters]);

  // Handle search query change with debounce option
  const handleSearchQueryChange = React.useCallback((value: string) => {
    setSearchQuery(value);
  }, []);

  // Handle pagination
  const handlePageChange = React.useCallback((newOffset: number) => {
    setFilters((prev) => ({
      ...prev,
      offset: newOffset,
    }));
  }, [setFilters]);

  // Handle page size change
  const handlePageSizeChange = React.useCallback((newLimit: number) => {
    setFilters((prev) => ({
      ...prev,
      limit: newLimit,
      offset: 0, // Reset to first page
    }));
  }, [setFilters]);

  // Handle sort change
  const handleSortChange = React.useCallback((
    sortBy: DeviceFilters['sortBy'],
    sortOrder: DeviceFilters['sortOrder']
  ) => {
    setFilters((prev) => ({
      ...prev,
      sortBy,
      sortOrder,
    }));
  }, [setFilters]);

  // Reset all filters
  const resetFilters = React.useCallback(() => {
    setSearchQuery('');
    setFilters(initialFilters);
  }, [initialFilters, setFilters]);

  // Check if any filters are active
  const hasActiveFilters = React.useMemo(() => {
    return !!(
      filters.q ||
      filters.schoolId ||
      filters.modelId ||
      filters.lifecycle ||
      filters.enrolled ||
      filters.category ||
      filters.make
    );
  }, [filters]);

  // Get active filter count
  const activeFilterCount = React.useMemo(() => {
    let count = 0;
    if (filters.q) count++;
    if (filters.schoolId) count++;
    if (filters.modelId) count++;
    if (filters.lifecycle) count++;
    if (filters.enrolled) count++;
    if (filters.category) count++;
    if (filters.make) count++;
    return count;
  }, [filters]);

  return {
    filters,
    searchQuery,
    setFilters,
    setSearchQuery: handleSearchQueryChange,
    handleFilterChange,
    handleSearch,
    handlePageChange,
    handlePageSizeChange,
    handleSortChange,
    resetFilters,
    hasActiveFilters,
    activeFilterCount,
  };
}
