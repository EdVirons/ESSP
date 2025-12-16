import { Button } from '@/components/ui/button';
import { Modal, ModalHeader, ModalBody, ModalFooter } from '@/components/ui/modal';
import { ConfirmDialog } from '@/components/ui/modal';
import type { DeviceModelManagerProps } from './types';
import { useDeviceModelManager } from './useDeviceModelManager';
import { DeviceModelForm } from './DeviceModelForm';
import { DeviceModelList } from './DeviceModelList';

export function DeviceModelManager({
  open,
  onClose,
  models,
  isLoading,
  onCreate,
  onUpdate,
  onDelete,
}: DeviceModelManagerProps) {
  const {
    // Search and filter state
    searchQuery,
    setSearchQuery,
    categoryFilter,
    setCategoryFilter,
    expandedMakes,
    toggleMake,

    // Computed data
    groupedModels,
    sortedMakes,

    // Form modal state
    showFormModal,
    setShowFormModal,
    editingModel,
    formData,
    formErrors,
    isSubmitting,

    // Delete state
    deleteModel,
    setDeleteModel,
    isDeleting,

    // Actions
    handleCreate,
    handleEdit,
    handleSubmit,
    handleDelete,
    updateSpec,
    updateFormField,
  } = useDeviceModelManager({
    models,
    onCreate,
    onUpdate,
    onDelete,
  });

  return (
    <>
      <Modal open={open} onClose={onClose} className="max-w-2xl">
        <ModalHeader onClose={onClose}>Device Model Catalog</ModalHeader>
        <ModalBody className="p-0">
          <DeviceModelList
            isLoading={isLoading}
            searchQuery={searchQuery}
            onSearchChange={setSearchQuery}
            categoryFilter={categoryFilter}
            onCategoryFilterChange={setCategoryFilter}
            sortedMakes={sortedMakes}
            groupedModels={groupedModels}
            expandedMakes={expandedMakes}
            onToggleMake={toggleMake}
            onCreateClick={handleCreate}
            onEditClick={handleEdit}
            onDeleteClick={setDeleteModel}
          />
        </ModalBody>
        <ModalFooter>
          <Button variant="outline" onClick={onClose}>
            Close
          </Button>
        </ModalFooter>
      </Modal>

      {/* Create/Edit Form Modal */}
      <DeviceModelForm
        open={showFormModal}
        onClose={() => setShowFormModal(false)}
        editingModel={editingModel}
        formData={formData}
        formErrors={formErrors}
        isSubmitting={isSubmitting}
        onSubmit={handleSubmit}
        onUpdateField={updateFormField}
        onUpdateSpec={updateSpec}
      />

      {/* Delete Confirmation */}
      <ConfirmDialog
        open={!!deleteModel}
        onClose={() => setDeleteModel(null)}
        onConfirm={handleDelete}
        title="Delete Device Model"
        description={`Are you sure you want to delete "${deleteModel?.make} ${deleteModel?.model}"? This action cannot be undone.`}
        confirmText="Delete"
        variant="destructive"
        isLoading={isDeleting}
      />
    </>
  );
}
