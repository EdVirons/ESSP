import * as React from 'react';
import {
  useCreateServiceShop,
  useUpdateServiceShop,
  useDeleteServiceShop,
  useCreateServiceStaff,
  useUpsertInventory,
} from '@/api/service-shops';
import type { ServiceShop, CreateServiceShopRequest, CreateServiceStaffRequest } from '@/types';
import { toast } from 'sonner';

interface UseServiceShopActionsOptions {
  onRefetch: () => void;
  onStaffRefetch?: () => void;
  onInventoryRefetch?: () => void;
}

const initialCreateForm: CreateServiceShopRequest = {
  name: '',
  countyCode: '',
  countyName: '',
  subCountyCode: '',
  subCountyName: '',
  coverageLevel: 'county',
  location: '',
  active: true,
};

export function useServiceShopActions({
  onRefetch,
  onStaffRefetch,
  onInventoryRefetch,
}: UseServiceShopActionsOptions) {
  // Selected shop state
  const [selectedShop, setSelectedShop] = React.useState<ServiceShop | null>(null);
  const [detailTab, setDetailTab] = React.useState('staff');

  // Modal states
  const [showDetail, setShowDetail] = React.useState(false);
  const [showCreateModal, setShowCreateModal] = React.useState(false);
  const [showEditModal, setShowEditModal] = React.useState(false);
  const [showDeleteModal, setShowDeleteModal] = React.useState(false);
  const [showAddStaffModal, setShowAddStaffModal] = React.useState(false);
  const [showAddInventoryModal, setShowAddInventoryModal] = React.useState(false);

  // Create form state
  const [createForm, setCreateForm] = React.useState<CreateServiceShopRequest>(initialCreateForm);

  // Mutations
  const createShop = useCreateServiceShop();
  const updateShop = useUpdateServiceShop();
  const deleteShop = useDeleteServiceShop();
  const createStaff = useCreateServiceStaff();
  const upsertInventory = useUpsertInventory();

  // Handle create shop
  const handleCreateShop = React.useCallback(async () => {
    try {
      await createShop.mutateAsync(createForm);
      toast.success('Service shop created successfully');
      setShowCreateModal(false);
      setCreateForm(initialCreateForm);
      onRefetch();
    } catch (err) {
      toast.error('Failed to create service shop');
    }
  }, [createShop, createForm, onRefetch]);

  // Handle update shop
  const handleUpdateShop = React.useCallback(
    async (id: string, data: Partial<CreateServiceShopRequest>) => {
      try {
        await updateShop.mutateAsync({ id, data });
        toast.success('Service shop updated successfully');
        setShowEditModal(false);
        setSelectedShop(null);
        onRefetch();
      } catch (err) {
        toast.error('Failed to update service shop');
      }
    },
    [updateShop, onRefetch]
  );

  // Handle delete shop
  const handleDeleteShop = React.useCallback(async () => {
    if (!selectedShop) return;
    try {
      await deleteShop.mutateAsync(selectedShop.id);
      toast.success('Service shop deleted successfully');
      setShowDeleteModal(false);
      setSelectedShop(null);
      onRefetch();
    } catch (err) {
      toast.error('Failed to delete service shop');
    }
  }, [deleteShop, selectedShop, onRefetch]);

  // Handle add staff
  const handleAddStaff = React.useCallback(
    async (data: CreateServiceStaffRequest) => {
      try {
        await createStaff.mutateAsync(data);
        toast.success('Staff member added successfully');
        setShowAddStaffModal(false);
        onStaffRefetch?.();
      } catch (err) {
        toast.error('Failed to add staff member');
      }
    },
    [createStaff, onStaffRefetch]
  );

  // Handle add inventory
  const handleAddInventory = React.useCallback(
    async (data: {
      serviceShopId: string;
      partId: string;
      qtyOnHand: number;
      reorderLevel: number;
    }) => {
      try {
        await upsertInventory.mutateAsync(data);
        toast.success('Inventory item added successfully');
        setShowAddInventoryModal(false);
        onInventoryRefetch?.();
      } catch (err) {
        toast.error('Failed to add inventory item');
      }
    },
    [upsertInventory, onInventoryRefetch]
  );

  // Handle edit click
  const handleEditClick = React.useCallback((shop: ServiceShop) => {
    setSelectedShop(shop);
    setShowEditModal(true);
  }, []);

  // Handle delete click
  const handleDeleteClick = React.useCallback((shop: ServiceShop) => {
    setSelectedShop(shop);
    setShowDeleteModal(true);
  }, []);

  // Handle view detail
  const handleViewDetail = React.useCallback((shop: ServiceShop) => {
    setSelectedShop(shop);
    setShowDetail(true);
    setDetailTab('staff');
  }, []);

  // Handle row click
  const handleRowClick = React.useCallback((shop: ServiceShop) => {
    setSelectedShop(shop);
    setShowDetail(true);
    setDetailTab('staff');
  }, []);

  // Close modals
  const closeEditModal = React.useCallback(() => {
    setShowEditModal(false);
    setSelectedShop(null);
  }, []);

  const closeDeleteModal = React.useCallback(() => {
    setShowDeleteModal(false);
    setSelectedShop(null);
  }, []);

  const closeDetail = React.useCallback(() => {
    setShowDetail(false);
  }, []);

  // Open edit from detail
  const openEditFromDetail = React.useCallback(() => {
    setShowDetail(false);
    setShowEditModal(true);
  }, []);

  return {
    // Selected shop
    selectedShop,
    setSelectedShop,
    detailTab,
    setDetailTab,

    // Modal states
    showDetail,
    showCreateModal,
    setShowCreateModal,
    showEditModal,
    showDeleteModal,
    showAddStaffModal,
    setShowAddStaffModal,
    showAddInventoryModal,
    setShowAddInventoryModal,

    // Create form
    createForm,
    setCreateForm,

    // Loading states
    isCreating: createShop.isPending,
    isUpdating: updateShop.isPending,
    isDeleting: deleteShop.isPending,
    isAddingStaff: createStaff.isPending,
    isAddingInventory: upsertInventory.isPending,

    // Handlers
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
  };
}
