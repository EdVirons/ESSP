import * as React from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { Modal, ModalHeader, ModalBody, ModalFooter } from '@/components/ui/modal';
import type { InventoryDevice, Location, AssignDeviceRequest, AssignmentType } from '@/types';

const ASSIGNMENT_TYPES: Array<{ value: AssignmentType; label: string }> = [
  { value: 'permanent', label: 'Permanent' },
  { value: 'temporary', label: 'Temporary' },
  { value: 'loan', label: 'On Loan' },
  { value: 'repair', label: 'In Repair' },
  { value: 'storage', label: 'In Storage' },
];

interface AssignDeviceModalProps {
  open: boolean;
  onClose: () => void;
  onSubmit: (deviceId: string, data: AssignDeviceRequest) => void;
  isLoading: boolean;
  device: InventoryDevice | null;
  locations: Location[];
  mode?: 'location' | 'user'; // Which field to focus on
}

export function AssignDeviceModal({
  open,
  onClose,
  onSubmit,
  isLoading,
  device,
  locations,
  mode = 'location',
}: AssignDeviceModalProps) {
  const [formData, setFormData] = React.useState<AssignDeviceRequest>({
    locationId: '',
    assignedToUser: '',
    assignmentType: 'permanent',
    notes: '',
  });

  // Reset form when modal opens or device changes
  React.useEffect(() => {
    if (open && device) {
      setFormData({
        locationId: device.location?.id || '',
        assignedToUser: '',
        assignmentType: 'permanent',
        notes: '',
      });
    }
  }, [open, device]);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!device) return;

    onSubmit(device.id, {
      ...formData,
      locationId: formData.locationId || undefined,
      assignedToUser: formData.assignedToUser || undefined,
    });
  };

  const handleClose = () => {
    setFormData({
      locationId: '',
      assignedToUser: '',
      assignmentType: 'permanent',
      notes: '',
    });
    onClose();
  };

  if (!device) return null;

  return (
    <Modal open={open} onClose={handleClose} className="max-w-lg">
      <ModalHeader onClose={handleClose}>
        {mode === 'user' ? 'Assign Device to User' : 'Change Device Location'}
      </ModalHeader>
      <form onSubmit={handleSubmit}>
        <ModalBody>
          <div className="space-y-4">
            {/* Device info */}
            <div className="bg-gray-50 rounded-lg p-3">
              <p className="text-xs text-gray-500">Device</p>
              <p className="font-medium text-gray-900">
                {device.model} - {device.serial}
              </p>
            </div>

            {/* Location selection */}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Location
              </label>
              <select
                value={formData.locationId}
                onChange={(e) => setFormData({ ...formData, locationId: e.target.value })}
                className="w-full rounded-md border border-gray-300 bg-white px-3 py-2 text-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
                autoFocus={mode === 'location'}
              >
                <option value="">No location assigned</option>
                {locations.filter(l => l.active).map((loc) => (
                  <option key={loc.id} value={loc.id}>
                    {loc.name} ({loc.locationType})
                  </option>
                ))}
              </select>
            </div>

            {/* User assignment */}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Assigned User
              </label>
              <Input
                value={formData.assignedToUser}
                onChange={(e) => setFormData({ ...formData, assignedToUser: e.target.value })}
                placeholder="e.g., Student ID or name"
                autoFocus={mode === 'user'}
              />
              <p className="mt-1 text-xs text-gray-500">
                Enter student ID, staff ID, or name
              </p>
            </div>

            {/* Assignment type */}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Assignment Type
              </label>
              <select
                value={formData.assignmentType}
                onChange={(e) => setFormData({ ...formData, assignmentType: e.target.value as AssignmentType })}
                className="w-full rounded-md border border-gray-300 bg-white px-3 py-2 text-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
              >
                {ASSIGNMENT_TYPES.map((type) => (
                  <option key={type.value} value={type.value}>
                    {type.label}
                  </option>
                ))}
              </select>
            </div>

            {/* Notes */}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Notes
              </label>
              <Textarea
                value={formData.notes}
                onChange={(e) => setFormData({ ...formData, notes: e.target.value })}
                placeholder="Any additional notes about this assignment"
                rows={2}
              />
            </div>
          </div>
        </ModalBody>
        <ModalFooter>
          <Button type="button" variant="outline" onClick={handleClose}>
            Cancel
          </Button>
          <Button type="submit" disabled={isLoading}>
            {isLoading ? 'Saving...' : 'Save Assignment'}
          </Button>
        </ModalFooter>
      </form>
    </Modal>
  );
}
