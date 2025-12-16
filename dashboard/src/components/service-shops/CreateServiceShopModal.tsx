import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Modal, ModalHeader, ModalBody, ModalFooter } from '@/components/ui/modal';
import { Select } from '@/components/ui/select';
import type { CreateServiceShopRequest } from '@/types';

interface CreateServiceShopModalProps {
  open: boolean;
  onClose: () => void;
  formData: CreateServiceShopRequest;
  onFormChange: (data: CreateServiceShopRequest) => void;
  onSubmit: () => void;
  isLoading: boolean;
}

export function CreateServiceShopModal({
  open,
  onClose,
  formData,
  onFormChange,
  onSubmit,
  isLoading,
}: CreateServiceShopModalProps) {
  return (
    <Modal open={open} onClose={onClose} className="max-w-lg">
      <ModalHeader onClose={onClose}>Create New Service Shop</ModalHeader>
      <ModalBody>
        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Shop Name *
            </label>
            <Input
              value={formData.name}
              onChange={(e) =>
                onFormChange({ ...formData, name: e.target.value })
              }
              placeholder="Enter shop name"
            />
          </div>
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                County Code *
              </label>
              <Input
                value={formData.countyCode}
                onChange={(e) =>
                  onFormChange({ ...formData, countyCode: e.target.value })
                }
                placeholder="e.g., 047"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                County Name
              </label>
              <Input
                value={formData.countyName || ''}
                onChange={(e) =>
                  onFormChange({ ...formData, countyName: e.target.value })
                }
                placeholder="e.g., Nairobi"
              />
            </div>
          </div>
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Sub-County Code
              </label>
              <Input
                value={formData.subCountyCode || ''}
                onChange={(e) =>
                  onFormChange({ ...formData, subCountyCode: e.target.value })
                }
                placeholder="Optional"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Sub-County Name
              </label>
              <Input
                value={formData.subCountyName || ''}
                onChange={(e) =>
                  onFormChange({ ...formData, subCountyName: e.target.value })
                }
                placeholder="Optional"
              />
            </div>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Coverage Level
            </label>
            <Select
              value={formData.coverageLevel || 'county'}
              onChange={(value) =>
                onFormChange({ ...formData, coverageLevel: value })
              }
              options={[
                { value: 'county', label: 'County' },
                { value: 'sub_county', label: 'Sub-County' },
                { value: 'region', label: 'Region' },
              ]}
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Location / Address
            </label>
            <Input
              value={formData.location || ''}
              onChange={(e) =>
                onFormChange({ ...formData, location: e.target.value })
              }
              placeholder="Physical address or location description"
            />
          </div>
        </div>
      </ModalBody>
      <ModalFooter>
        <Button variant="outline" onClick={onClose}>
          Cancel
        </Button>
        <Button
          onClick={onSubmit}
          disabled={!formData.name || !formData.countyCode || isLoading}
        >
          {isLoading ? 'Creating...' : 'Create Service Shop'}
        </Button>
      </ModalFooter>
    </Modal>
  );
}
