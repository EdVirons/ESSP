import * as React from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Modal, ModalHeader, ModalBody, ModalFooter } from '@/components/ui/modal';
import { OrgUnitSelect } from './OrgUnitSelect';
import { OrgUnitKindSelect } from './OrgUnitKindSelect';
import type { CreateOrgUnitInput, OrgUnitSnapshot } from '@/types/hr';

interface CreateOrgUnitModalProps {
  open: boolean;
  onClose: () => void;
  onSubmit: (data: CreateOrgUnitInput) => void;
  isLoading: boolean;
  orgUnits?: OrgUnitSnapshot[];
}

const initialFormData: CreateOrgUnitInput = {
  code: '',
  name: '',
  kind: 'department',
  parentId: '',
};

export function CreateOrgUnitModal({
  open,
  onClose,
  onSubmit,
  isLoading,
  orgUnits = [],
}: CreateOrgUnitModalProps) {
  const [formData, setFormData] = React.useState<CreateOrgUnitInput>(initialFormData);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit(formData);
  };

  const handleClose = () => {
    setFormData(initialFormData);
    onClose();
  };

  const updateField = <K extends keyof CreateOrgUnitInput>(key: K, value: CreateOrgUnitInput[K]) => {
    setFormData((prev) => ({ ...prev, [key]: value }));
  };

  // Auto-generate code from name
  const handleNameChange = (name: string) => {
    const code = name
      .toUpperCase()
      .replace(/[^A-Z0-9]+/g, '-')
      .replace(/^-|-$/g, '');
    setFormData((prev) => ({ ...prev, name, code }));
  };

  return (
    <Modal open={open} onClose={handleClose} className="max-w-lg">
      <ModalHeader onClose={handleClose}>Create Organization Unit</ModalHeader>
      <form onSubmit={handleSubmit}>
        <ModalBody>
          <div className="space-y-4">
            <div>
              <label htmlFor="orgName" className="block text-sm font-medium text-gray-700 mb-1">
                Name *
              </label>
              <Input
                id="orgName"
                value={formData.name}
                onChange={(e) => handleNameChange(e.target.value)}
                placeholder="Engineering Department"
                required
              />
            </div>

            <div>
              <label htmlFor="orgCode" className="block text-sm font-medium text-gray-700 mb-1">
                Code *
              </label>
              <Input
                id="orgCode"
                value={formData.code}
                onChange={(e) => updateField('code', e.target.value)}
                placeholder="ENG-DEPT"
                required
              />
              <p className="text-xs text-gray-500 mt-1">
                Unique code for the org unit (auto-generated from name)
              </p>
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div>
                <label htmlFor="orgKind" className="block text-sm font-medium text-gray-700 mb-1">
                  Type
                </label>
                <OrgUnitKindSelect
                  id="orgKind"
                  value={formData.kind || 'department'}
                  onChange={(value) => updateField('kind', value)}
                />
              </div>
              <div>
                <label htmlFor="orgParent" className="block text-sm font-medium text-gray-700 mb-1">
                  Parent Unit
                </label>
                <OrgUnitSelect
                  id="orgParent"
                  value={formData.parentId || ''}
                  onChange={(value) => updateField('parentId', value)}
                  orgUnits={orgUnits}
                  noneLabel="None (Top Level)"
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
            disabled={!formData.code || !formData.name || isLoading}
          >
            {isLoading ? 'Creating...' : 'Create Org Unit'}
          </Button>
        </ModalFooter>
      </form>
    </Modal>
  );
}
