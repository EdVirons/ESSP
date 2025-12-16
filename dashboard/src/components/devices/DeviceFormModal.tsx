import * as React from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Modal, ModalHeader, ModalBody, ModalFooter } from '@/components/ui/modal';
import { Select } from '@/components/ui/select';
import { Textarea } from '@/components/ui/textarea';
import type {
  Device,
  DeviceModel,
  CreateDeviceInput,
  UpdateDeviceInput,
} from '@/types/device';
import {
  LIFECYCLE_STATUS_OPTIONS,
  ENROLLMENT_STATUS_OPTIONS,
} from '@/types/device';

interface DeviceFormModalProps {
  open: boolean;
  onClose: () => void;
  onSubmit: (data: CreateDeviceInput | UpdateDeviceInput) => void;
  isLoading: boolean;
  device?: Device | null; // If provided, it's edit mode
  models: DeviceModel[];
  schools: Array<{ value: string; label: string }>;
}

interface FormData {
  serial: string;
  assetTag: string;
  modelId: string;
  schoolId: string;
  lifecycle: string;
  enrolled: string;
  assignedTo: string;
  notes: string;
  warrantyExpiry: string;
  purchaseDate: string;
}

const initialFormData: FormData = {
  serial: '',
  assetTag: '',
  modelId: '',
  schoolId: '',
  lifecycle: 'in_stock',
  enrolled: 'unenrolled',
  assignedTo: '',
  notes: '',
  warrantyExpiry: '',
  purchaseDate: '',
};

export function DeviceFormModal({
  open,
  onClose,
  onSubmit,
  isLoading,
  device,
  models,
  schools,
}: DeviceFormModalProps) {
  const isEditMode = !!device;
  const [formData, setFormData] = React.useState<FormData>(initialFormData);
  const [errors, setErrors] = React.useState<Partial<Record<keyof FormData, string>>>({});

  // Populate form when editing
  React.useEffect(() => {
    if (device) {
      setFormData({
        serial: device.serial || '',
        assetTag: device.assetTag || '',
        modelId: device.modelId || '',
        schoolId: device.schoolId || '',
        lifecycle: device.lifecycle || 'in_stock',
        enrolled: device.enrolled || 'unenrolled',
        assignedTo: device.assignedTo || '',
        notes: device.notes || '',
        warrantyExpiry: device.warrantyExpiry
          ? device.warrantyExpiry.split('T')[0]
          : '',
        purchaseDate: device.purchaseDate
          ? device.purchaseDate.split('T')[0]
          : '',
      });
    } else {
      setFormData(initialFormData);
    }
    setErrors({});
  }, [device, open]);

  // Build model options
  const modelOptions = React.useMemo(() => {
    const options = models.map((m) => ({
      value: m.id,
      label: `${m.make} ${m.model}`,
    }));
    return [{ value: '', label: 'Select a model...' }, ...options];
  }, [models]);

  // Build school options
  const schoolOptions = React.useMemo(() => {
    return [{ value: '', label: 'Select a school...' }, ...schools];
  }, [schools]);

  // Handle input change
  const handleChange = (
    field: keyof FormData,
    value: string
  ) => {
    setFormData((prev) => ({ ...prev, [field]: value }));
    // Clear error when user starts typing
    if (errors[field]) {
      setErrors((prev) => ({ ...prev, [field]: undefined }));
    }
  };

  // Validate form
  const validate = (): boolean => {
    const newErrors: Partial<Record<keyof FormData, string>> = {};

    if (!formData.serial.trim()) {
      newErrors.serial = 'Serial number is required';
    }
    if (!formData.modelId) {
      newErrors.modelId = 'Device model is required';
    }
    if (!formData.schoolId) {
      newErrors.schoolId = 'School is required';
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  // Handle form submission
  const handleSubmit = () => {
    if (!validate()) return;

    const data: CreateDeviceInput | UpdateDeviceInput = {
      serial: formData.serial.trim(),
      assetTag: formData.assetTag.trim() || undefined,
      modelId: formData.modelId,
      schoolId: formData.schoolId,
      lifecycle: formData.lifecycle as CreateDeviceInput['lifecycle'],
      enrolled: formData.enrolled as CreateDeviceInput['enrolled'],
      assignedTo: formData.assignedTo.trim() || undefined,
      notes: formData.notes.trim() || undefined,
      warrantyExpiry: formData.warrantyExpiry || undefined,
      purchaseDate: formData.purchaseDate || undefined,
    };

    onSubmit(data);
  };

  return (
    <Modal open={open} onClose={onClose} className="max-w-lg">
      <ModalHeader onClose={onClose}>
        {isEditMode ? 'Edit Device' : 'Add New Device'}
      </ModalHeader>
      <ModalBody>
        <div className="space-y-4">
          {/* Serial Number */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Serial Number *
            </label>
            <Input
              value={formData.serial}
              onChange={(e) => handleChange('serial', e.target.value)}
              placeholder="Enter serial number"
              error={!!errors.serial}
            />
            {errors.serial && (
              <p className="text-sm text-red-600 mt-1">{errors.serial}</p>
            )}
          </div>

          {/* Asset Tag */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Asset Tag
            </label>
            <Input
              value={formData.assetTag}
              onChange={(e) => handleChange('assetTag', e.target.value)}
              placeholder="Enter asset tag (optional)"
            />
          </div>

          {/* Device Model */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Device Model *
            </label>
            <Select
              value={formData.modelId}
              onChange={(value) => handleChange('modelId', value)}
              options={modelOptions}
              error={!!errors.modelId}
            />
            {errors.modelId && (
              <p className="text-sm text-red-600 mt-1">{errors.modelId}</p>
            )}
          </div>

          {/* School */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              School *
            </label>
            <Select
              value={formData.schoolId}
              onChange={(value) => handleChange('schoolId', value)}
              options={schoolOptions}
              error={!!errors.schoolId}
            />
            {errors.schoolId && (
              <p className="text-sm text-red-600 mt-1">{errors.schoolId}</p>
            )}
          </div>

          {/* Status and Enrollment */}
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Lifecycle Status
              </label>
              <Select
                value={formData.lifecycle}
                onChange={(value) => handleChange('lifecycle', value)}
                options={LIFECYCLE_STATUS_OPTIONS}
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Enrollment
              </label>
              <Select
                value={formData.enrolled}
                onChange={(value) => handleChange('enrolled', value)}
                options={ENROLLMENT_STATUS_OPTIONS}
              />
            </div>
          </div>

          {/* Assigned To */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Assigned To
            </label>
            <Input
              value={formData.assignedTo}
              onChange={(e) => handleChange('assignedTo', e.target.value)}
              placeholder="User ID or name (optional)"
            />
          </div>

          {/* Dates */}
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Purchase Date
              </label>
              <Input
                type="date"
                value={formData.purchaseDate}
                onChange={(e) => handleChange('purchaseDate', e.target.value)}
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Warranty Expiry
              </label>
              <Input
                type="date"
                value={formData.warrantyExpiry}
                onChange={(e) => handleChange('warrantyExpiry', e.target.value)}
              />
            </div>
          </div>

          {/* Notes */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Notes
            </label>
            <Textarea
              value={formData.notes}
              onChange={(e) => handleChange('notes', e.target.value)}
              placeholder="Additional notes (optional)"
              rows={3}
            />
          </div>
        </div>
      </ModalBody>
      <ModalFooter>
        <Button variant="outline" onClick={onClose} disabled={isLoading}>
          Cancel
        </Button>
        <Button onClick={handleSubmit} disabled={isLoading}>
          {isLoading
            ? isEditMode
              ? 'Saving...'
              : 'Creating...'
            : isEditMode
            ? 'Save Changes'
            : 'Create Device'}
        </Button>
      </ModalFooter>
    </Modal>
  );
}
