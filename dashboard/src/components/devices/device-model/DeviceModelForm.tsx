import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Modal, ModalHeader, ModalBody, ModalFooter } from '@/components/ui/modal';
import { Select } from '@/components/ui/select';
import { DEVICE_CATEGORY_OPTIONS } from '@/types/device';
import type { DeviceCategory } from '@/types/device';
import type { ModelFormData, DeviceModel } from './types';
import { specLabels } from './types';

interface DeviceModelFormProps {
  open: boolean;
  onClose: () => void;
  editingModel: DeviceModel | null;
  formData: ModelFormData;
  formErrors: Record<string, string>;
  isSubmitting: boolean;
  onSubmit: () => void;
  onUpdateField: <K extends keyof ModelFormData>(field: K, value: ModelFormData[K]) => void;
  onUpdateSpec: (key: string, value: string) => void;
}

export function DeviceModelForm({
  open,
  onClose,
  editingModel,
  formData,
  formErrors,
  isSubmitting,
  onSubmit,
  onUpdateField,
  onUpdateSpec,
}: DeviceModelFormProps) {
  return (
    <Modal open={open} onClose={onClose} className="max-w-md">
      <ModalHeader onClose={onClose}>
        {editingModel ? 'Edit Model' : 'Add Device Model'}
      </ModalHeader>
      <ModalBody>
        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Make *
            </label>
            <Input
              value={formData.make}
              onChange={(e) => onUpdateField('make', e.target.value)}
              placeholder="e.g., Dell, HP, Lenovo"
              error={!!formErrors.make}
            />
            {formErrors.make && (
              <p className="text-sm text-red-600 mt-1">{formErrors.make}</p>
            )}
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Model *
            </label>
            <Input
              value={formData.model}
              onChange={(e) => onUpdateField('model', e.target.value)}
              placeholder="e.g., Latitude 5520"
              error={!!formErrors.model}
            />
            {formErrors.model && (
              <p className="text-sm text-red-600 mt-1">{formErrors.model}</p>
            )}
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Category *
            </label>
            <Select
              value={formData.category}
              onChange={(value) => onUpdateField('category', value as DeviceCategory)}
              options={DEVICE_CATEGORY_OPTIONS}
            />
          </div>

          {/* Specifications */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Specifications
            </label>
            <div className="space-y-2 bg-gray-50 rounded-lg p-3">
              {specLabels.map(({ key, label }) => (
                <div key={key} className="flex items-center gap-2">
                  <label className="text-sm text-gray-600 w-24 flex-shrink-0">
                    {label}
                  </label>
                  <Input
                    value={formData.specs[key] || ''}
                    onChange={(e) => onUpdateSpec(key, e.target.value)}
                    placeholder={`Enter ${label.toLowerCase()}`}
                    className="flex-1"
                  />
                </div>
              ))}
            </div>
          </div>
        </div>
      </ModalBody>
      <ModalFooter>
        <Button variant="outline" onClick={onClose} disabled={isSubmitting}>
          Cancel
        </Button>
        <Button onClick={onSubmit} disabled={isSubmitting}>
          {isSubmitting
            ? 'Saving...'
            : editingModel
            ? 'Save Changes'
            : 'Create Model'}
        </Button>
      </ModalFooter>
    </Modal>
  );
}
