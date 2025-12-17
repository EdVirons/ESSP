import * as React from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { Modal, ModalHeader, ModalBody, ModalFooter } from '@/components/ui/modal';
import type { RegisterDeviceRequest, Location } from '@/types';

interface AddDeviceModalProps {
  open: boolean;
  onClose: () => void;
  onSubmit: (data: RegisterDeviceRequest) => void;
  isLoading: boolean;
  locations?: Location[];
}

export function AddDeviceModal({
  open,
  onClose,
  onSubmit,
  isLoading,
  locations = [],
}: AddDeviceModalProps) {
  const [formData, setFormData] = React.useState<RegisterDeviceRequest>({
    serial: '',
    assetTag: '',
    model: '',
    make: '',
    notes: '',
    locationId: '',
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit({
      ...formData,
      locationId: formData.locationId || undefined,
    });
  };

  const handleClose = () => {
    setFormData({
      serial: '',
      assetTag: '',
      model: '',
      make: '',
      notes: '',
      locationId: '',
    });
    onClose();
  };

  return (
    <Modal open={open} onClose={handleClose} className="max-w-lg">
      <ModalHeader onClose={handleClose}>Register New Device</ModalHeader>
      <form onSubmit={handleSubmit}>
        <ModalBody>
          <div className="space-y-4">
            <div className="grid grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Serial Number *
                </label>
                <Input
                  value={formData.serial}
                  onChange={(e) => setFormData({ ...formData, serial: e.target.value })}
                  placeholder="e.g., SN123456789"
                  required
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Asset Tag
                </label>
                <Input
                  value={formData.assetTag}
                  onChange={(e) => setFormData({ ...formData, assetTag: e.target.value })}
                  placeholder="e.g., TAG-001"
                />
              </div>
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Make
                </label>
                <Input
                  value={formData.make}
                  onChange={(e) => setFormData({ ...formData, make: e.target.value })}
                  placeholder="e.g., HP, Dell, Lenovo"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Model *
                </label>
                <Input
                  value={formData.model}
                  onChange={(e) => setFormData({ ...formData, model: e.target.value })}
                  placeholder="e.g., Chromebook 11 G8"
                  required
                />
              </div>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Initial Location
              </label>
              <select
                value={formData.locationId}
                onChange={(e) => setFormData({ ...formData, locationId: e.target.value })}
                className="w-full rounded-md border border-gray-300 bg-white px-3 py-2 text-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
              >
                <option value="">No location assigned</option>
                {locations.filter(l => l.active).map((loc) => (
                  <option key={loc.id} value={loc.id}>
                    {loc.name} ({loc.locationType})
                  </option>
                ))}
              </select>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Notes
              </label>
              <Textarea
                value={formData.notes}
                onChange={(e) => setFormData({ ...formData, notes: e.target.value })}
                placeholder="Any additional notes about this device"
                rows={3}
              />
            </div>
          </div>
        </ModalBody>
        <ModalFooter>
          <Button type="button" variant="outline" onClick={handleClose}>
            Cancel
          </Button>
          <Button
            type="submit"
            disabled={!formData.serial || !formData.model || isLoading}
          >
            {isLoading ? 'Registering...' : 'Register Device'}
          </Button>
        </ModalFooter>
      </form>
    </Modal>
  );
}
