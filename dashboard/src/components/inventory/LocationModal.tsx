import * as React from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Modal, ModalHeader, ModalBody, ModalFooter } from '@/components/ui/modal';
import type { Location, CreateLocationRequest, UpdateLocationRequest, LocationType } from '@/types';

const LOCATION_TYPES: Array<{ value: LocationType; label: string }> = [
  { value: 'block', label: 'Block' },
  { value: 'floor', label: 'Floor' },
  { value: 'room', label: 'Room' },
  { value: 'lab', label: 'Computer Lab' },
  { value: 'storage', label: 'Storage' },
  { value: 'office', label: 'Office' },
];

interface LocationModalProps {
  open: boolean;
  onClose: () => void;
  onSubmit: (data: CreateLocationRequest | UpdateLocationRequest) => void;
  isLoading: boolean;
  location?: Location | null; // If provided, edit mode
  locations?: Location[]; // For parent selection
}

export function LocationModal({
  open,
  onClose,
  onSubmit,
  isLoading,
  location,
  locations = [],
}: LocationModalProps) {
  const isEdit = !!location;

  const [formData, setFormData] = React.useState<CreateLocationRequest>({
    parentId: '',
    name: '',
    locationType: 'room',
    code: '',
    capacity: 0,
  });

  // Reset form when modal opens/closes or location changes
  React.useEffect(() => {
    if (open) {
      if (location) {
        setFormData({
          parentId: location.parentId || '',
          name: location.name,
          locationType: location.locationType,
          code: location.code,
          capacity: location.capacity,
        });
      } else {
        setFormData({
          parentId: '',
          name: '',
          locationType: 'room',
          code: '',
          capacity: 0,
        });
      }
    }
  }, [open, location]);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit({
      ...formData,
      parentId: formData.parentId || undefined,
    });
  };

  const handleClose = () => {
    setFormData({
      parentId: '',
      name: '',
      locationType: 'room',
      code: '',
      capacity: 0,
    });
    onClose();
  };

  // Filter out the current location from parent options (can't be parent of itself)
  const parentOptions = locations.filter(l => l.active && (!location || l.id !== location.id));

  return (
    <Modal open={open} onClose={handleClose} className="max-w-lg">
      <ModalHeader onClose={handleClose}>
        {isEdit ? 'Edit Location' : 'Create New Location'}
      </ModalHeader>
      <form onSubmit={handleSubmit}>
        <ModalBody>
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Name *
              </label>
              <Input
                value={formData.name}
                onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                placeholder="e.g., Computer Lab A"
                required
              />
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Type *
                </label>
                <select
                  value={formData.locationType}
                  onChange={(e) => setFormData({ ...formData, locationType: e.target.value as LocationType })}
                  className="w-full rounded-md border border-gray-300 bg-white px-3 py-2 text-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
                  required
                >
                  {LOCATION_TYPES.map((type) => (
                    <option key={type.value} value={type.value}>
                      {type.label}
                    </option>
                  ))}
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Code
                </label>
                <Input
                  value={formData.code}
                  onChange={(e) => setFormData({ ...formData, code: e.target.value })}
                  placeholder="e.g., B1-R101"
                />
              </div>
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Parent Location
                </label>
                <select
                  value={formData.parentId}
                  onChange={(e) => setFormData({ ...formData, parentId: e.target.value })}
                  className="w-full rounded-md border border-gray-300 bg-white px-3 py-2 text-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
                >
                  <option value="">No parent (top level)</option>
                  {parentOptions.map((loc) => (
                    <option key={loc.id} value={loc.id}>
                      {loc.name} ({loc.locationType})
                    </option>
                  ))}
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Capacity
                </label>
                <Input
                  type="number"
                  min="0"
                  value={formData.capacity}
                  onChange={(e) => setFormData({ ...formData, capacity: parseInt(e.target.value) || 0 })}
                  placeholder="0"
                />
              </div>
            </div>
          </div>
        </ModalBody>
        <ModalFooter>
          <Button type="button" variant="outline" onClick={handleClose}>
            Cancel
          </Button>
          <Button
            type="submit"
            disabled={!formData.name || isLoading}
          >
            {isLoading ? (isEdit ? 'Updating...' : 'Creating...') : (isEdit ? 'Update Location' : 'Create Location')}
          </Button>
        </ModalFooter>
      </form>
    </Modal>
  );
}
