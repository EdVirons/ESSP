import * as React from 'react';
import { toast } from 'sonner';
import {
  useCreateServiceStaff,
  useUpdateServiceStaff,
  useDeleteServiceStaff,
} from '@/api/service-shops';
import type { ServiceStaff, CreateServiceStaffRequest, UpdateServiceStaffRequest } from '@/types';

interface UseStaffActionsOptions {
  onRefetch: () => void;
}

export function useStaffActions({ onRefetch }: UseStaffActionsOptions) {
  const [selectedStaff, setSelectedStaff] = React.useState<ServiceStaff | null>(null);
  const [showDetail, setShowDetail] = React.useState(false);
  const [showCreateModal, setShowCreateModal] = React.useState(false);
  const [showEditModal, setShowEditModal] = React.useState(false);
  const [showDeleteModal, setShowDeleteModal] = React.useState(false);

  const createMutation = useCreateServiceStaff();
  const updateMutation = useUpdateServiceStaff();
  const deleteMutation = useDeleteServiceStaff();

  const handleCreateStaff = React.useCallback(
    async (data: CreateServiceStaffRequest) => {
      try {
        await createMutation.mutateAsync(data);
        toast.success('Staff member added successfully');
        setShowCreateModal(false);
        onRefetch();
      } catch (error) {
        toast.error('Failed to add staff member');
      }
    },
    [createMutation, onRefetch]
  );

  const handleUpdateStaff = React.useCallback(
    async (id: string, data: UpdateServiceStaffRequest) => {
      try {
        await updateMutation.mutateAsync({ id, data });
        toast.success('Staff member updated successfully');
        setShowEditModal(false);
        setSelectedStaff(null);
        onRefetch();
      } catch (error) {
        toast.error('Failed to update staff member');
      }
    },
    [updateMutation, onRefetch]
  );

  const handleDeleteStaff = React.useCallback(async () => {
    if (!selectedStaff) return;
    try {
      await deleteMutation.mutateAsync(selectedStaff.id);
      toast.success('Staff member removed successfully');
      setShowDeleteModal(false);
      setSelectedStaff(null);
      onRefetch();
    } catch (error) {
      toast.error('Failed to remove staff member');
    }
  }, [selectedStaff, deleteMutation, onRefetch]);

  const handleEditClick = React.useCallback((staff: ServiceStaff) => {
    setSelectedStaff(staff);
    setShowEditModal(true);
  }, []);

  const handleDeleteClick = React.useCallback((staff: ServiceStaff) => {
    setSelectedStaff(staff);
    setShowDeleteModal(true);
  }, []);

  const handleViewDetail = React.useCallback((staff: ServiceStaff) => {
    setSelectedStaff(staff);
    setShowDetail(true);
  }, []);

  const handleRowClick = React.useCallback((staff: ServiceStaff) => {
    setSelectedStaff(staff);
    setShowDetail(true);
  }, []);

  const closeEditModal = React.useCallback(() => {
    setShowEditModal(false);
  }, []);

  const closeDeleteModal = React.useCallback(() => {
    setShowDeleteModal(false);
  }, []);

  const closeDetail = React.useCallback(() => {
    setShowDetail(false);
  }, []);

  return {
    selectedStaff,
    showDetail,
    showCreateModal,
    setShowCreateModal,
    showEditModal,
    showDeleteModal,
    isCreating: createMutation.isPending,
    isUpdating: updateMutation.isPending,
    isDeleting: deleteMutation.isPending,
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
  };
}
