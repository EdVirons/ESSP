import * as React from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { Modal, ModalHeader, ModalBody, ModalFooter } from '@/components/ui/modal';
import type { CreatePartRequest } from '@/types';

interface CreatePartModalProps {
  open: boolean;
  onClose: () => void;
  onSubmit: (data: CreatePartRequest) => void;
  isLoading: boolean;
  categories?: string[];
}

export function CreatePartModal({
  open,
  onClose,
  onSubmit,
  isLoading,
  categories = [],
}: CreatePartModalProps) {
  const [formData, setFormData] = React.useState<CreatePartRequest>({
    sku: '',
    name: '',
    category: '',
    description: '',
    unitCostCents: 0,
    supplier: '',
    supplierSku: '',
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit(formData);
  };

  const handleClose = () => {
    setFormData({
      sku: '',
      name: '',
      category: '',
      description: '',
      unitCostCents: 0,
      supplier: '',
      supplierSku: '',
    });
    onClose();
  };

  // Format cents to dollars for display
  const displayPrice = (cents: number) => (cents / 100).toFixed(2);

  // Parse dollars to cents
  const parsePriceToCents = (value: string) => {
    const num = parseFloat(value);
    return isNaN(num) ? 0 : Math.round(num * 100);
  };

  return (
    <Modal open={open} onClose={handleClose} className="max-w-lg">
      <ModalHeader onClose={handleClose}>Create New Part</ModalHeader>
      <form onSubmit={handleSubmit}>
        <ModalBody>
          <div className="space-y-4">
            <div className="grid grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  SKU *
                </label>
                <Input
                  value={formData.sku}
                  onChange={(e) => setFormData({ ...formData, sku: e.target.value })}
                  placeholder="e.g., WDG-001"
                  required
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Name *
                </label>
                <Input
                  value={formData.name}
                  onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                  placeholder="Part name"
                  required
                />
              </div>
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Category
                </label>
                <Input
                  value={formData.category}
                  onChange={(e) => setFormData({ ...formData, category: e.target.value })}
                  placeholder="e.g., Electronics"
                  list="category-list"
                />
                <datalist id="category-list">
                  {(categories ?? []).map((cat) => (
                    <option key={cat} value={cat} />
                  ))}
                </datalist>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Unit Price ($)
                </label>
                <Input
                  type="number"
                  step="0.01"
                  min="0"
                  value={displayPrice(formData.unitCostCents || 0)}
                  onChange={(e) => setFormData({ ...formData, unitCostCents: parsePriceToCents(e.target.value) })}
                  placeholder="0.00"
                />
              </div>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Description
              </label>
              <Textarea
                value={formData.description}
                onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                placeholder="Part description"
                rows={3}
              />
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Supplier
                </label>
                <Input
                  value={formData.supplier}
                  onChange={(e) => setFormData({ ...formData, supplier: e.target.value })}
                  placeholder="Supplier name"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Supplier SKU
                </label>
                <Input
                  value={formData.supplierSku}
                  onChange={(e) => setFormData({ ...formData, supplierSku: e.target.value })}
                  placeholder="Supplier's part number"
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
            disabled={!formData.sku || !formData.name || isLoading}
          >
            {isLoading ? 'Creating...' : 'Create Part'}
          </Button>
        </ModalFooter>
      </form>
    </Modal>
  );
}
