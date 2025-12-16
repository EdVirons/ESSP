import * as React from 'react';
import {
  useCreatePart,
  useUpdatePart,
  useDeletePart,
  useImportParts,
  useExportParts,
} from '@/api/parts';
import type { Part, CreatePartRequest, UpdatePartRequest } from '@/types';
import { toast } from 'sonner';

interface UsePartsCRUDOptions {
  onRefetch: () => void;
}

export function usePartsCRUD({ onRefetch }: UsePartsCRUDOptions) {
  // Modal states
  const [createModalOpen, setCreateModalOpen] = React.useState(false);
  const [editModalOpen, setEditModalOpen] = React.useState(false);
  const [deleteModalOpen, setDeleteModalOpen] = React.useState(false);
  const [importModalOpen, setImportModalOpen] = React.useState(false);
  const [selectedPart, setSelectedPart] = React.useState<Part | null>(null);

  // API mutations
  const createPart = useCreatePart();
  const updatePart = useUpdatePart();
  const deletePart = useDeletePart();
  const importParts = useImportParts();
  const exportParts = useExportParts();

  const handleCreate = React.useCallback(
    (data: CreatePartRequest) => {
      createPart.mutate(data, {
        onSuccess: () => {
          toast.success('Part created successfully');
          setCreateModalOpen(false);
          onRefetch();
        },
        onError: (error) => {
          toast.error(error.message || 'Failed to create part');
        },
      });
    },
    [createPart, onRefetch]
  );

  const handleEdit = React.useCallback((part: Part) => {
    setSelectedPart(part);
    setEditModalOpen(true);
  }, []);

  const handleUpdate = React.useCallback(
    (id: string, data: UpdatePartRequest) => {
      updatePart.mutate(
        { id, data },
        {
          onSuccess: () => {
            toast.success('Part updated successfully');
            setEditModalOpen(false);
            setSelectedPart(null);
            onRefetch();
          },
          onError: (error) => {
            toast.error(error.message || 'Failed to update part');
          },
        }
      );
    },
    [updatePart, onRefetch]
  );

  const handleDeleteClick = React.useCallback((part: Part) => {
    setSelectedPart(part);
    setDeleteModalOpen(true);
  }, []);

  const handleDelete = React.useCallback(() => {
    if (!selectedPart) return;
    deletePart.mutate(selectedPart.id, {
      onSuccess: () => {
        toast.success('Part deleted successfully');
        setDeleteModalOpen(false);
        setSelectedPart(null);
        onRefetch();
      },
      onError: (error) => {
        toast.error(error.message || 'Failed to delete part');
      },
    });
  }, [selectedPart, deletePart, onRefetch]);

  const handleImport = React.useCallback(
    (file: File) => {
      importParts.mutate(file, {
        onSuccess: (result) => {
          const message = `Created ${result.created} parts${result.failed > 0 ? `. ${result.failed} failed.` : ''}`;
          toast.success(message);
          setImportModalOpen(false);
          onRefetch();
        },
        onError: (error) => {
          toast.error(error.message || 'Import failed');
        },
      });
    },
    [importParts, onRefetch]
  );

  const handleExport = React.useCallback(() => {
    exportParts.mutate(undefined, {
      onSuccess: () => {
        toast.success('CSV file downloaded');
      },
      onError: (error) => {
        toast.error(error.message || 'Export failed');
      },
    });
  }, [exportParts]);

  const closeEditModal = React.useCallback(() => {
    setEditModalOpen(false);
    setSelectedPart(null);
  }, []);

  const closeDeleteModal = React.useCallback(() => {
    setDeleteModalOpen(false);
    setSelectedPart(null);
  }, []);

  return {
    // Modal states
    createModalOpen,
    setCreateModalOpen,
    editModalOpen,
    deleteModalOpen,
    importModalOpen,
    setImportModalOpen,
    selectedPart,

    // Loading states
    isCreating: createPart.isPending,
    isUpdating: updatePart.isPending,
    isDeleting: deletePart.isPending,
    isImporting: importParts.isPending,
    isExporting: exportParts.isPending,

    // Handlers
    handleCreate,
    handleEdit,
    handleUpdate,
    handleDeleteClick,
    handleDelete,
    handleImport,
    handleExport,
    closeEditModal,
    closeDeleteModal,
  };
}
