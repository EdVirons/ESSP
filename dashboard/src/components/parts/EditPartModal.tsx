import * as React from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { Modal, ModalHeader, ModalBody, ModalFooter } from '@/components/ui/modal';
import type { Part, UpdatePartRequest } from '@/types';

interface EditPartModalProps {
  part: Part | null;
  open: boolean;
  onClose: () => void;
  onSubmit: (id: string, data: UpdatePartRequest) => void;
  isLoading: boolean;
  categories?: string[];
}

export function EditPartModal({
  part,
  open,
  onClose,
  onSubmit,
  isLoading,
  categories = [],
}: EditPartModalProps) {
  const [formData, setFormData] = React.useState<UpdatePartRequest>({});

  // Initialize form when part changes
  React.useEffect(() => {
    if (part) {
      setFormData({
        name: part.name,
        category: part.category,
        description: part.description,
        unitCostCents: part.unitCostCents,
        supplier: part.supplier,
        supplierSku: part.supplierSku,
        active: part.active,
      });
    }
  }, [part]);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (part) {
      onSubmit(part.id, formData);
    }
  };

  // Format cents to dollars for display
  const displayPrice = (cents: number | undefined) =>
    cents !== undefined ? (cents / 100).toFixed(2) : '0.00';

  // Parse dollars to cents
  const parsePriceToCents = (value: string) => {
    const num = parseFloat(value);
    return isNaN(num) ? 0 : Math.round(num * 100);
  };

  if (!part) return null;

  return (
    <Modal open={open} onClose={onClose} className="max-w-lg">
      <ModalHeader onClose={onClose}>Edit Part</ModalHeader>
      <form onSubmit={handleSubmit}>
        <ModalBody>
          <div className="space-y-4">
            <div className="grid grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  SKU
                </label>
                <Input
                  value={part.sku}
                  disabled
                  className="bg-gray-50"
                />
                <p className="text-xs text-gray-500 mt-1">SKU cannot be changed</p>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Name *
                </label>
                <Input
                  value={formData.name || ''}
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
                  value={formData.category || ''}
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
                  value={displayPrice(formData.unitCostCents)}
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
                value={formData.description || ''}
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
                  value={formData.supplier || ''}
                  onChange={(e) => setFormData({ ...formData, supplier: e.target.value })}
                  placeholder="Supplier name"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Supplier SKU
                </label>
                <Input
                  value={formData.supplierSku || ''}
                  onChange={(e) => setFormData({ ...formData, supplierSku: e.target.value })}
                  placeholder="Supplier's part number"
                />
              </div>
            </div>

            <div className="flex items-center gap-2">
              <input
                type="checkbox"
                id="active"
                checked={formData.active ?? true}
                onChange={(e) => setFormData({ ...formData, active: e.target.checked })}
                className="h-4 w-4 rounded border-gray-300 text-blue-600 focus:ring-blue-500"
              />
              <label htmlFor="active" className="text-sm text-gray-700">
                Active (visible in catalog)
              </label>
            </div>
          </div>
        </ModalBody>
        <ModalFooter>
          <Button type="button" variant="outline" onClick={onClose}>
            Cancel
          </Button>
          <Button type="submit" disabled={!formData.name || isLoading}>
            {isLoading ? 'Saving...' : 'Save Changes'}
          </Button>
        </ModalFooter>
      </form>
    </Modal>
  );
}
