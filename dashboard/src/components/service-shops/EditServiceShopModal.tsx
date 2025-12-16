import * as React from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Modal, ModalHeader, ModalBody, ModalFooter } from '@/components/ui/modal';
import { Select } from '@/components/ui/select';
import { Badge } from '@/components/ui/badge';
import type { ServiceShop, CreateServiceShopRequest } from '@/types';

interface EditServiceShopModalProps {
  shop: ServiceShop | null;
  open: boolean;
  onClose: () => void;
  onSubmit: (id: string, data: Partial<CreateServiceShopRequest>) => void;
  isLoading: boolean;
}

export function EditServiceShopModal({
  shop,
  open,
  onClose,
  onSubmit,
  isLoading,
}: EditServiceShopModalProps) {
  const [formData, setFormData] = React.useState<Partial<CreateServiceShopRequest>>({});

  React.useEffect(() => {
    if (shop) {
      setFormData({
        name: shop.name,
        countyCode: shop.countyCode,
        countyName: shop.countyName,
        subCountyCode: shop.subCountyCode,
        subCountyName: shop.subCountyName,
        coverageLevel: shop.coverageLevel,
        location: shop.location,
        active: shop.active,
      });
    }
  }, [shop]);

  const handleSubmit = () => {
    if (shop) {
      onSubmit(shop.id, formData);
    }
  };

  if (!shop) return null;

  return (
    <Modal open={open} onClose={onClose} className="max-w-lg">
      <ModalHeader onClose={onClose}>Edit Service Shop</ModalHeader>
      <ModalBody>
        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Shop Name *
            </label>
            <Input
              value={formData.name || ''}
              onChange={(e) =>
                setFormData({ ...formData, name: e.target.value })
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
                value={formData.countyCode || ''}
                onChange={(e) =>
                  setFormData({ ...formData, countyCode: e.target.value })
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
                  setFormData({ ...formData, countyName: e.target.value })
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
                  setFormData({ ...formData, subCountyCode: e.target.value })
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
                  setFormData({ ...formData, subCountyName: e.target.value })
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
                setFormData({ ...formData, coverageLevel: value })
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
                setFormData({ ...formData, location: e.target.value })
              }
              placeholder="Physical address or location description"
            />
          </div>
          <div className="flex items-center justify-between py-2">
            <div>
              <label className="block text-sm font-medium text-gray-700">
                Active Status
              </label>
              <p className="text-sm text-gray-500">
                Inactive shops won't receive new work orders
              </p>
            </div>
            <Button
              type="button"
              variant="outline"
              size="sm"
              onClick={() => setFormData({ ...formData, active: !formData.active })}
              className={formData.active ? 'bg-green-50 border-green-200' : 'bg-gray-50'}
            >
              <Badge className={formData.active ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-600'}>
                {formData.active ? 'Active' : 'Inactive'}
              </Badge>
            </Button>
          </div>
        </div>
      </ModalBody>
      <ModalFooter>
        <Button variant="outline" onClick={onClose}>
          Cancel
        </Button>
        <Button
          onClick={handleSubmit}
          disabled={!formData.name || !formData.countyCode || isLoading}
        >
          {isLoading ? 'Saving...' : 'Save Changes'}
        </Button>
      </ModalFooter>
    </Modal>
  );
}
