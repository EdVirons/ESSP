import * as React from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Modal, ModalHeader, ModalBody, ModalFooter } from '@/components/ui/modal';
import { OrgUnitSelect } from './OrgUnitSelect';
import { PersonStatusSelect } from './PersonStatusSelect';
import type { CreatePersonInput, OrgUnitSnapshot } from '@/types/hr';

interface CreatePersonModalProps {
  open: boolean;
  onClose: () => void;
  onSubmit: (data: CreatePersonInput) => void;
  isLoading: boolean;
  orgUnits?: OrgUnitSnapshot[];
}

const initialFormData: CreatePersonInput = {
  givenName: '',
  familyName: '',
  email: '',
  phone: '',
  title: '',
  status: 'active',
  orgUnitId: '',
};

export function CreatePersonModal({
  open,
  onClose,
  onSubmit,
  isLoading,
  orgUnits = [],
}: CreatePersonModalProps) {
  const [formData, setFormData] = React.useState<CreatePersonInput>(initialFormData);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit(formData);
  };

  const handleClose = () => {
    setFormData(initialFormData);
    onClose();
  };

  const updateField = <K extends keyof CreatePersonInput>(key: K, value: CreatePersonInput[K]) => {
    setFormData((prev) => ({ ...prev, [key]: value }));
  };

  return (
    <Modal open={open} onClose={handleClose} className="max-w-lg">
      <ModalHeader onClose={handleClose}>Add New Person</ModalHeader>
      <form onSubmit={handleSubmit}>
        <ModalBody>
          <div className="space-y-4">
            <div className="grid grid-cols-2 gap-4">
              <div>
                <label htmlFor="givenName" className="block text-sm font-medium text-gray-700 mb-1">
                  First Name *
                </label>
                <Input
                  id="givenName"
                  value={formData.givenName}
                  onChange={(e) => updateField('givenName', e.target.value)}
                  placeholder="John"
                  required
                />
              </div>
              <div>
                <label htmlFor="familyName" className="block text-sm font-medium text-gray-700 mb-1">
                  Last Name *
                </label>
                <Input
                  id="familyName"
                  value={formData.familyName}
                  onChange={(e) => updateField('familyName', e.target.value)}
                  placeholder="Doe"
                  required
                />
              </div>
            </div>

            <div>
              <label htmlFor="email" className="block text-sm font-medium text-gray-700 mb-1">
                Email *
              </label>
              <Input
                id="email"
                type="email"
                value={formData.email}
                onChange={(e) => updateField('email', e.target.value)}
                placeholder="john.doe@example.com"
                required
              />
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div>
                <label htmlFor="phone" className="block text-sm font-medium text-gray-700 mb-1">
                  Phone
                </label>
                <Input
                  id="phone"
                  type="tel"
                  value={formData.phone || ''}
                  onChange={(e) => updateField('phone', e.target.value)}
                  placeholder="+1 555-123-4567"
                />
              </div>
              <div>
                <label htmlFor="title" className="block text-sm font-medium text-gray-700 mb-1">
                  Title
                </label>
                <Input
                  id="title"
                  value={formData.title || ''}
                  onChange={(e) => updateField('title', e.target.value)}
                  placeholder="Software Engineer"
                />
              </div>
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div>
                <label htmlFor="status" className="block text-sm font-medium text-gray-700 mb-1">
                  Status
                </label>
                <PersonStatusSelect
                  id="status"
                  value={formData.status || 'active'}
                  onChange={(value) => updateField('status', value)}
                />
              </div>
              <div>
                <label htmlFor="orgUnit" className="block text-sm font-medium text-gray-700 mb-1">
                  Organization Unit
                </label>
                <OrgUnitSelect
                  id="orgUnit"
                  value={formData.orgUnitId || ''}
                  onChange={(value) => updateField('orgUnitId', value)}
                  orgUnits={orgUnits}
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
            disabled={!formData.givenName || !formData.familyName || !formData.email || isLoading}
          >
            {isLoading ? 'Creating...' : 'Add Person'}
          </Button>
        </ModalFooter>
      </form>
    </Modal>
  );
}
