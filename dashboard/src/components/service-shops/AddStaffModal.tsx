import * as React from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Modal, ModalHeader, ModalBody, ModalFooter } from '@/components/ui/modal';
import { Select } from '@/components/ui/select';
import type { CreateServiceStaffRequest, StaffRole } from '@/types';

const roleOptions = [
  { value: 'lead_technician', label: 'Lead Technician' },
  { value: 'assistant_technician', label: 'Assistant Technician' },
  { value: 'storekeeper', label: 'Storekeeper' },
];

interface AddStaffModalProps {
  serviceShopId: string;
  open: boolean;
  onClose: () => void;
  onSubmit: (data: CreateServiceStaffRequest) => void;
  isLoading: boolean;
}

export function AddStaffModal({
  serviceShopId,
  open,
  onClose,
  onSubmit,
  isLoading,
}: AddStaffModalProps) {
  const [formData, setFormData] = React.useState<Omit<CreateServiceStaffRequest, 'serviceShopId'>>({
    userId: '',
    role: 'assistant_technician',
    phone: '',
    active: true,
  });

  const handleSubmit = () => {
    onSubmit({
      serviceShopId,
      ...formData,
    });
  };

  const handleClose = () => {
    setFormData({
      userId: '',
      role: 'assistant_technician',
      phone: '',
      active: true,
    });
    onClose();
  };

  return (
    <Modal open={open} onClose={handleClose} className="max-w-md">
      <ModalHeader onClose={handleClose}>Add Staff Member</ModalHeader>
      <ModalBody>
        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Staff Name / ID *
            </label>
            <Input
              value={formData.userId}
              onChange={(e) =>
                setFormData({ ...formData, userId: e.target.value })
              }
              placeholder="Enter staff name or ID"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Role *
            </label>
            <Select
              value={formData.role}
              onChange={(value) =>
                setFormData({ ...formData, role: value as StaffRole })
              }
              options={roleOptions}
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Phone Number
            </label>
            <Input
              value={formData.phone || ''}
              onChange={(e) =>
                setFormData({ ...formData, phone: e.target.value })
              }
              placeholder="+254..."
            />
          </div>
        </div>
      </ModalBody>
      <ModalFooter>
        <Button variant="outline" onClick={handleClose}>
          Cancel
        </Button>
        <Button
          onClick={handleSubmit}
          disabled={!formData.userId || !formData.role || isLoading}
        >
          {isLoading ? 'Adding...' : 'Add Staff'}
        </Button>
      </ModalFooter>
    </Modal>
  );
}
