import * as React from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { Modal, ModalHeader, ModalBody, ModalFooter } from '@/components/ui/modal';
import { OrgUnitSelect } from './OrgUnitSelect';
import type { CreateTeamInput, OrgUnitSnapshot } from '@/types/hr';

interface CreateTeamModalProps {
  open: boolean;
  onClose: () => void;
  onSubmit: (data: CreateTeamInput) => void;
  isLoading: boolean;
  orgUnits?: OrgUnitSnapshot[];
}

const initialFormData: CreateTeamInput = {
  key: '',
  name: '',
  description: '',
  orgUnitId: '',
};

export function CreateTeamModal({
  open,
  onClose,
  onSubmit,
  isLoading,
  orgUnits = [],
}: CreateTeamModalProps) {
  const [formData, setFormData] = React.useState<CreateTeamInput>(initialFormData);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit(formData);
  };

  const handleClose = () => {
    setFormData(initialFormData);
    onClose();
  };

  const updateField = <K extends keyof CreateTeamInput>(key: K, value: CreateTeamInput[K]) => {
    setFormData((prev) => ({ ...prev, [key]: value }));
  };

  // Auto-generate key from name
  const handleNameChange = (name: string) => {
    const key = name
      .toLowerCase()
      .replace(/[^a-z0-9]+/g, '-')
      .replace(/^-|-$/g, '');
    setFormData((prev) => ({ ...prev, name, key }));
  };

  return (
    <Modal open={open} onClose={handleClose} className="max-w-lg">
      <ModalHeader onClose={handleClose}>Create New Team</ModalHeader>
      <form onSubmit={handleSubmit}>
        <ModalBody>
          <div className="space-y-4">
            <div>
              <label htmlFor="teamName" className="block text-sm font-medium text-gray-700 mb-1">
                Team Name *
              </label>
              <Input
                id="teamName"
                value={formData.name}
                onChange={(e) => handleNameChange(e.target.value)}
                placeholder="Engineering Team"
                required
              />
            </div>

            <div>
              <label htmlFor="teamKey" className="block text-sm font-medium text-gray-700 mb-1">
                Team Key *
              </label>
              <Input
                id="teamKey"
                value={formData.key}
                onChange={(e) => updateField('key', e.target.value)}
                placeholder="engineering-team"
                required
              />
              <p className="text-xs text-gray-500 mt-1">
                Unique identifier for the team (auto-generated from name)
              </p>
            </div>

            <div>
              <label htmlFor="teamDescription" className="block text-sm font-medium text-gray-700 mb-1">
                Description
              </label>
              <Textarea
                id="teamDescription"
                value={formData.description || ''}
                onChange={(e) => updateField('description', e.target.value)}
                placeholder="Describe the team's purpose and responsibilities"
                rows={3}
              />
            </div>

            <div>
              <label htmlFor="teamOrgUnit" className="block text-sm font-medium text-gray-700 mb-1">
                Organization Unit
              </label>
              <OrgUnitSelect
                id="teamOrgUnit"
                value={formData.orgUnitId || ''}
                onChange={(value) => updateField('orgUnitId', value)}
                orgUnits={orgUnits}
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
            disabled={!formData.key || !formData.name || isLoading}
          >
            {isLoading ? 'Creating...' : 'Create Team'}
          </Button>
        </ModalFooter>
      </form>
    </Modal>
  );
}
