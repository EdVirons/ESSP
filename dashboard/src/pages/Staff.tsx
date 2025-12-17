import * as React from 'react';
import { Plus, Search, X } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardContent } from '@/components/ui/card';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { DataTable } from '@/components/ui/data-table';
import { ConfirmDialog } from '@/components/ui/modal';
import {
  StaffStats,
  StaffDetail,
  CreateStaffModal,
  EditStaffModal,
  createStaffColumns,
  roleOptions,
  statusOptions,
  useStaffActions,
} from '@/components/staff';
import {
  useServiceStaff,
  useServiceStaffStats,
  useServiceShops,
} from '@/api/service-shops';

export function Staff() {
  // Filters state
  const [filters, setFilters] = React.useState<{
    serviceShopId?: string;
    role?: string;
    active?: boolean;
    limit?: number;
  }>({
    limit: 100,
  });
  const [searchQuery, setSearchQuery] = React.useState('');

  // API hooks
  const { data, isLoading, refetch } = useServiceStaff(filters);
  const { data: statsData, isLoading: statsLoading } = useServiceStaffStats();
  const { data: shopsData } = useServiceShops({ limit: 200 });

  // Actions hook
  const {
    selectedStaff,
    showDetail,
    showCreateModal,
    setShowCreateModal,
    showEditModal,
    showDeleteModal,
    isCreating,
    isUpdating,
    isDeleting,
    handleCreateStaff,
    handleUpdateStaff,
    handleDeleteStaff,
    handleEditClick,
    handleDeleteClick,
    handleViewDetail,
    handleRowClick,
    closeEditModal,
    closeDeleteModal,
    closeDetail,
  } = useStaffActions({ onRefetch: refetch });

  // Shop lookup for display
  const shopLookup = React.useMemo(() => {
    const map = new Map<string, string>();
    shopsData?.items?.forEach((shop) => map.set(shop.id, shop.name));
    return map;
  }, [shopsData?.items]);

  // Filter staff by search query (client-side)
  const filteredStaff = React.useMemo(() => {
    let staff = data?.items || [];
    if (searchQuery) {
      const query = searchQuery.toLowerCase();
      staff = staff.filter(
        (s) =>
          s.userId.toLowerCase().includes(query) ||
          s.phone?.toLowerCase().includes(query) ||
          shopLookup.get(s.serviceShopId)?.toLowerCase().includes(query)
      );
    }
    return staff;
  }, [data?.items, searchQuery, shopLookup]);

  // Shop options for filter dropdown
  const shopOptions = React.useMemo(
    () => [
      { value: '', label: 'All Shops' },
      ...(shopsData?.items || []).map((s) => ({ value: s.id, label: s.name })),
    ],
    [shopsData?.items]
  );

  // Table columns
  const columns = React.useMemo(
    () =>
      createStaffColumns({
        shopLookup,
        onEdit: handleEditClick,
        onDelete: handleDeleteClick,
        onViewDetail: handleViewDetail,
      }),
    [shopLookup, handleEditClick, handleDeleteClick, handleViewDetail]
  );

  // Check if any filters are active
  const hasActiveFilters =
    filters.serviceShopId || filters.role || filters.active !== undefined;

  // Clear all filters
  const clearFilters = () => {
    setFilters({ limit: 100 });
    setSearchQuery('');
  };

  return (
    <div className="space-y-6">
      {/* Page Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Staff Management</h1>
          <p className="text-sm text-gray-500">
            Manage technicians and staff across all service shops
          </p>
        </div>
        <Button onClick={() => setShowCreateModal(true)}>
          <Plus className="h-4 w-4 mr-2" />
          Add Staff Member
        </Button>
      </div>

      {/* Stats Cards */}
      <StaffStats stats={statsData} isLoading={statsLoading} />

      {/* Filters */}
      <Card>
        <CardContent className="p-4">
          <div className="flex flex-wrap items-center gap-4">
            <div className="relative flex-1 min-w-[200px] max-w-md">
              <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-400" />
              <Input
                placeholder="Search by name, phone, or shop..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="pl-9"
              />
            </div>

            <Select
              value={filters.serviceShopId || ''}
              onValueChange={(value) =>
                setFilters((prev) => ({
                  ...prev,
                  serviceShopId: value || undefined,
                }))
              }
            >
              <SelectTrigger className="w-48">
                <SelectValue placeholder="All Shops" />
              </SelectTrigger>
              <SelectContent>
                {shopOptions.map((opt) => (
                  <SelectItem key={opt.value} value={opt.value}>
                    {opt.label}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>

            <Select
              value={filters.role || ''}
              onValueChange={(value) =>
                setFilters((prev) => ({
                  ...prev,
                  role: value || undefined,
                }))
              }
            >
              <SelectTrigger className="w-44">
                <SelectValue placeholder="All Roles" />
              </SelectTrigger>
              <SelectContent>
                {roleOptions.map((opt) => (
                  <SelectItem key={opt.value} value={opt.value}>
                    {opt.label}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>

            <Select
              value={filters.active === undefined ? '' : String(filters.active)}
              onValueChange={(value) =>
                setFilters((prev) => ({
                  ...prev,
                  active: value === '' ? undefined : value === 'true',
                }))
              }
            >
              <SelectTrigger className="w-32">
                <SelectValue placeholder="All Status" />
              </SelectTrigger>
              <SelectContent>
                {statusOptions.map((opt) => (
                  <SelectItem key={opt.value} value={opt.value}>
                    {opt.label}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>

            {hasActiveFilters && (
              <Button
                variant="ghost"
                size="sm"
                onClick={clearFilters}
                className="text-gray-500"
              >
                <X className="h-4 w-4 mr-1" />
                Clear Filters
              </Button>
            )}
          </div>
        </CardContent>
      </Card>

      {/* Staff Table */}
      <Card>
        <CardContent className="p-0">
          <DataTable
            columns={columns}
            data={filteredStaff}
            isLoading={isLoading}
            onRowClick={handleRowClick}
          />
        </CardContent>
      </Card>

      {/* Detail Sheet */}
      <StaffDetail
        staff={selectedStaff}
        shopName={
          selectedStaff ? shopLookup.get(selectedStaff.serviceShopId) : undefined
        }
        open={showDetail}
        onClose={closeDetail}
        onEditClick={() => {
          closeDetail();
          if (selectedStaff) handleEditClick(selectedStaff);
        }}
      />

      {/* Create Modal */}
      <CreateStaffModal
        open={showCreateModal}
        onClose={() => setShowCreateModal(false)}
        onSubmit={handleCreateStaff}
        isLoading={isCreating}
        shops={shopsData?.items || []}
      />

      {/* Edit Modal */}
      <EditStaffModal
        staff={selectedStaff}
        open={showEditModal}
        onClose={closeEditModal}
        onSubmit={handleUpdateStaff}
        isLoading={isUpdating}
        shops={shopsData?.items || []}
      />

      {/* Delete Confirm Modal */}
      <ConfirmDialog
        open={showDeleteModal}
        onClose={closeDeleteModal}
        onConfirm={handleDeleteStaff}
        title="Remove Staff Member"
        description={`Are you sure you want to remove "${selectedStaff?.userId}"? This action cannot be undone.`}
        confirmText="Remove"
        isLoading={isDeleting}
        variant="destructive"
      />
    </div>
  );
}
