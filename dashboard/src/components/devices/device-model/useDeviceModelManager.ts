import * as React from 'react';
import type {
  DeviceModel,
  ModelFormData,
  CreateDeviceModelInput,
  UpdateDeviceModelInput,
} from './types';
import { initialFormData } from './types';

interface UseDeviceModelManagerProps {
  models: DeviceModel[];
  onCreate: (data: CreateDeviceModelInput) => Promise<void>;
  onUpdate: (id: string, data: UpdateDeviceModelInput) => Promise<void>;
  onDelete: (id: string) => Promise<void>;
}

export function useDeviceModelManager({
  models,
  onCreate,
  onUpdate,
  onDelete,
}: UseDeviceModelManagerProps) {
  const [searchQuery, setSearchQuery] = React.useState('');
  const [categoryFilter, setCategoryFilter] = React.useState<string>('');
  const [expandedMakes, setExpandedMakes] = React.useState<Set<string>>(new Set());

  // Form modal state
  const [showFormModal, setShowFormModal] = React.useState(false);
  const [editingModel, setEditingModel] = React.useState<DeviceModel | null>(null);
  const [formData, setFormData] = React.useState<ModelFormData>(initialFormData);
  const [formErrors, setFormErrors] = React.useState<Record<string, string>>({});
  const [isSubmitting, setIsSubmitting] = React.useState(false);

  // Delete confirmation state
  const [deleteModel, setDeleteModel] = React.useState<DeviceModel | null>(null);
  const [isDeleting, setIsDeleting] = React.useState(false);

  // Filter models
  const filteredModels = React.useMemo(() => {
    return models.filter((m) => {
      const matchesSearch =
        !searchQuery ||
        m.make.toLowerCase().includes(searchQuery.toLowerCase()) ||
        m.model.toLowerCase().includes(searchQuery.toLowerCase());
      const matchesCategory = !categoryFilter || m.category === categoryFilter;
      return matchesSearch && matchesCategory;
    });
  }, [models, searchQuery, categoryFilter]);

  // Group models by make
  const groupedModels = React.useMemo(() => {
    const groups: Record<string, DeviceModel[]> = {};
    filteredModels.forEach((m) => {
      if (!groups[m.make]) {
        groups[m.make] = [];
      }
      groups[m.make].push(m);
    });
    return groups;
  }, [filteredModels]);

  // Sort makes alphabetically
  const sortedMakes = Object.keys(groupedModels).sort();

  // Toggle make expansion
  const toggleMake = (make: string) => {
    setExpandedMakes((prev) => {
      const next = new Set(prev);
      if (next.has(make)) {
        next.delete(make);
      } else {
        next.add(make);
      }
      return next;
    });
  };

  // Open form modal for create
  const handleCreate = () => {
    setEditingModel(null);
    setFormData(initialFormData);
    setFormErrors({});
    setShowFormModal(true);
  };

  // Open form modal for edit
  const handleEdit = (model: DeviceModel) => {
    setEditingModel(model);
    setFormData({
      make: model.make,
      model: model.model,
      category: model.category,
      specs: { ...model.specs },
    });
    setFormErrors({});
    setShowFormModal(true);
  };

  // Validate form
  const validateForm = (): boolean => {
    const errors: Record<string, string> = {};
    if (!formData.make.trim()) {
      errors.make = 'Make is required';
    }
    if (!formData.model.trim()) {
      errors.model = 'Model is required';
    }
    setFormErrors(errors);
    return Object.keys(errors).length === 0;
  };

  // Submit form
  const handleSubmit = async () => {
    if (!validateForm()) return;

    setIsSubmitting(true);
    try {
      const data: CreateDeviceModelInput | UpdateDeviceModelInput = {
        make: formData.make.trim(),
        model: formData.model.trim(),
        category: formData.category,
        specs: formData.specs,
      };

      if (editingModel) {
        await onUpdate(editingModel.id, data);
      } else {
        await onCreate(data as CreateDeviceModelInput);
      }
      setShowFormModal(false);
    } finally {
      setIsSubmitting(false);
    }
  };

  // Handle delete
  const handleDelete = async () => {
    if (!deleteModel) return;

    setIsDeleting(true);
    try {
      await onDelete(deleteModel.id);
      setDeleteModel(null);
    } finally {
      setIsDeleting(false);
    }
  };

  // Update spec value
  const updateSpec = (key: string, value: string) => {
    setFormData((prev) => ({
      ...prev,
      specs: {
        ...prev.specs,
        [key]: value || undefined,
      },
    }));
  };

  // Update form field
  const updateFormField = <K extends keyof ModelFormData>(
    field: K,
    value: ModelFormData[K]
  ) => {
    setFormData((prev) => ({ ...prev, [field]: value }));
  };

  return {
    // Search and filter state
    searchQuery,
    setSearchQuery,
    categoryFilter,
    setCategoryFilter,
    expandedMakes,
    toggleMake,

    // Computed data
    filteredModels,
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
  };
}
