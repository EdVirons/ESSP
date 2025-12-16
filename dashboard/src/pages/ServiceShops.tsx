import * as React from 'react';
import { Plus, Search } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardContent } from '@/components/ui/card';
import { Select } from '@/components/ui/select';
import { DataTable } from '@/components/ui/data-table';
import { ConfirmDialog } from '@/components/ui/modal';
import {
  ServiceShopsStats,
  ServiceShopDetail,
  CreateServiceShopModal,
  EditServiceShopModal,
  AddStaffModal,
  AddInventoryModal,
  createServiceShopColumns,
  statusOptions,
  coverageOptions,
  useServiceShopActions,
} from '@/components/service-shops';
import { useServiceShops, useServiceStaff, useInventory } from '@/api/service-shops';

export function ServiceShops() {
  // Filters state
  const [filters, setFilters] = React.useState<{
    active?: boolean;
    countyCode?: string;
    limit?: number;
  }>({
    limit: 50,
  });
  const [searchQuery, setSearchQuery] = React.useState('');
  const [coverageFilter, setCoverageFilter] = React.useState('');

  // API hooks
  const { data, isLoading, refetch } = useServiceShops(filters);

  // Staff and inventory for all shops (for stats)
  const { data: allStaffData } = useServiceStaff({ limit: 1000 });
  const { data: allInventoryData } = useInventory({ limit: 1000, lowStock: true });

  // Actions hook
  const {
    selectedShop,
    detailTab,
    setDetailTab,
    showDetail,
    showCreateModal,
    setShowCreateModal,
    showEditModal,
    showDeleteModal,
    showAddStaffModal,
    setShowAddStaffModal,
    showAddInventoryModal,
    setShowAddInventoryModal,
    createForm,
    setCreateForm,
    isCreating,
    isUpdating,
    isDeleting,
    isAddingStaff,
    isAddingInventory,
    handleCreateShop,
    handleUpdateShop,
    handleDeleteShop,
    handleAddStaff,
    handleAddInventory,
    handleEditClick,
    handleDeleteClick,
    handleViewDetail,
    handleRowClick,
    closeEditModal,
    closeDeleteModal,
    closeDetail,
    openEditFromDetail,
  } = useServiceShopActions({
    onRefetch: refetch,
    onStaffRefetch: () => refetchStaff(),
    onInventoryRefetch: () => refetchInventory(),
  });

  // Detail view data
  const { data: staffData, refetch: refetchStaff } = useServiceStaff({
    serviceShopId: selectedShop?.id,
  });
  const { data: inventoryData, refetch: refetchInventory } = useInventory({
    serviceShopId: selectedShop?.id,
  });

  // Filter shops by coverage level (client-side)
  const filteredShops = React.useMemo(() => {
    let shops = data?.items || [];
    if (coverageFilter) {
      shops = shops.filter((s) => s.coverageLevel === coverageFilter);
    }
    if (searchQuery) {
      const query = searchQuery.toLowerCase();
      shops = shops.filter(
        (s) =>
          s.name.toLowerCase().includes(query) ||
          s.countyName.toLowerCase().includes(query) ||
          s.location?.toLowerCase().includes(query)
      );
    }
    return shops;
  }, [data?.items, coverageFilter, searchQuery]);

  // Table columns
  const columns = React.useMemo(
    () =>
      createServiceShopColumns({
        allStaff: allStaffData?.items,
        allInventory: allInventoryData?.items,
        onEdit: handleEditClick,
        onDelete: handleDeleteClick,
        onViewDetail: handleViewDetail,
      }),
    [allStaffData?.items, allInventoryData?.items, handleEditClick, handleDeleteClick, handleViewDetail]
  );

  const shops = filteredShops;
  const staff = staffData?.items || [];
  const inventory = inventoryData?.items || [];
  const totalStaff = allStaffData?.items?.length || 0;
  const lowStockCount = allInventoryData?.items?.length || 0;

  return (
    <div className="space-y-6">
      {/* Page Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Service Shops</h1>
          <p className="text-sm text-gray-500">Manage service shops, staff, and inventory</p>
        </div>
        <Button onClick={() => setShowCreateModal(true)}>
          <Plus className="h-4 w-4" />
          Create Service Shop
        </Button>
      </div>

      {/* Stats Cards */}
      <ServiceShopsStats shops={shops} totalStaff={totalStaff} lowStockCount={lowStockCount} />

      {/* Filters */}
      <Card>
        <CardContent className="p-4">
          <div className="flex flex-wrap items-center gap-4">
            <div className="relative flex-1 min-w-[200px] max-w-md">
              <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-400" />
              <Input
                placeholder="Search shops by name, county, or location..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="pl-9"
              />
            </div>
            <Select
              value={filters.active === undefined ? '' : String(filters.active)}
              onChange={(value) =>
                setFilters((prev) => ({
                  ...prev,
                  active: value === '' ? undefined : value === 'true',
                }))
              }
              options={statusOptions}
              placeholder="Status"
              className="w-32"
            />
            <Select
              value={coverageFilter}
              onChange={(value) => setCoverageFilter(value)}
              options={coverageOptions}
              placeholder="Coverage"
              className="w-36"
            />
          </div>
        </CardContent>
      </Card>

      {/* Service Shops Table */}
      <Card>
        <CardContent className="p-0">
          <DataTable
            columns={columns}
            data={shops}
            isLoading={isLoading}
            searchKey="name"
            searchPlaceholder="Search by name..."
            showRowSelection
            showColumnVisibility
            onRowClick={handleRowClick}
            emptyMessage="No service shops found. Create your first service shop to get started."
          />
        </CardContent>
      </Card>

      {/* Service Shop Detail Sheet */}
      <ServiceShopDetail
        shop={selectedShop}
        open={showDetail}
        onClose={closeDetail}
        detailTab={detailTab}
        onDetailTabChange={setDetailTab}
        staff={staff}
        inventory={inventory}
        onEditClick={openEditFromDetail}
        onAddStaffClick={() => setShowAddStaffModal(true)}
        onAddInventoryClick={() => setShowAddInventoryModal(true)}
      />

      {/* Create Service Shop Modal */}
      <CreateServiceShopModal
        open={showCreateModal}
        onClose={() => setShowCreateModal(false)}
        formData={createForm}
        onFormChange={setCreateForm}
        onSubmit={handleCreateShop}
        isLoading={isCreating}
      />

      {/* Edit Service Shop Modal */}
      <EditServiceShopModal
        shop={selectedShop}
        open={showEditModal}
        onClose={closeEditModal}
        onSubmit={handleUpdateShop}
        isLoading={isUpdating}
      />

      {/* Delete Confirm Modal */}
      <ConfirmDialog
        open={showDeleteModal}
        onClose={closeDeleteModal}
        onConfirm={handleDeleteShop}
        title="Delete Service Shop"
        description={`Are you sure you want to delete "${selectedShop?.name}"? This will also remove all staff assignments and inventory. This action cannot be undone.`}
        confirmText="Delete"
        isLoading={isDeleting}
        variant="destructive"
      />

      {/* Add Staff Modal */}
      {selectedShop && (
        <AddStaffModal
          serviceShopId={selectedShop.id}
          open={showAddStaffModal}
          onClose={() => setShowAddStaffModal(false)}
          onSubmit={handleAddStaff}
          isLoading={isAddingStaff}
        />
      )}

      {/* Add Inventory Modal */}
      {selectedShop && (
        <AddInventoryModal
          serviceShopId={selectedShop.id}
          open={showAddInventoryModal}
          onClose={() => setShowAddInventoryModal(false)}
          onSubmit={handleAddInventory}
          isLoading={isAddingInventory}
        />
      )}
    </div>
  );
}
