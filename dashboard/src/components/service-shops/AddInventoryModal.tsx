import * as React from 'react';
import { Search, Package } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Modal, ModalHeader, ModalBody, ModalFooter } from '@/components/ui/modal';
import { useParts } from '@/api/service-shops';

interface AddInventoryModalProps {
  serviceShopId: string;
  open: boolean;
  onClose: () => void;
  onSubmit: (data: { serviceShopId: string; partId: string; qtyOnHand: number; reorderLevel: number }) => void;
  isLoading: boolean;
}

export function AddInventoryModal({
  serviceShopId,
  open,
  onClose,
  onSubmit,
  isLoading,
}: AddInventoryModalProps) {
  const [searchQuery, setSearchQuery] = React.useState('');
  const [selectedPart, setSelectedPart] = React.useState<{ id: string; name: string; sku: string } | null>(null);
  const [qtyOnHand, setQtyOnHand] = React.useState(0);
  const [reorderLevel, setReorderLevel] = React.useState(5);

  const { data: partsData, isLoading: partsLoading } = useParts({ limit: 20, active: true });

  const filteredParts = React.useMemo(() => {
    if (!partsData?.items) return [];
    if (!searchQuery) return partsData.items.slice(0, 10);
    const query = searchQuery.toLowerCase();
    return partsData.items.filter(
      (p) =>
        p.name.toLowerCase().includes(query) ||
        p.puk.toLowerCase().includes(query)
    );
  }, [partsData?.items, searchQuery]);

  const handleSubmit = () => {
    if (!selectedPart) return;
    onSubmit({
      serviceShopId,
      partId: selectedPart.id,
      qtyOnHand,
      reorderLevel,
    });
  };

  const handleClose = () => {
    setSearchQuery('');
    setSelectedPart(null);
    setQtyOnHand(0);
    setReorderLevel(5);
    onClose();
  };

  return (
    <Modal open={open} onClose={handleClose} className="max-w-lg">
      <ModalHeader onClose={handleClose}>Add Inventory Item</ModalHeader>
      <ModalBody>
        <div className="space-y-4">
          {!selectedPart ? (
            <>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Search Parts
                </label>
                <div className="relative">
                  <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-400" />
                  <Input
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                    placeholder="Search by name or SKU..."
                    className="pl-9"
                  />
                </div>
              </div>
              <div className="border rounded-lg max-h-60 overflow-auto">
                {partsLoading ? (
                  <div className="p-4 text-center text-gray-500">Loading parts...</div>
                ) : filteredParts.length === 0 ? (
                  <div className="p-4 text-center text-gray-500">
                    <Package className="h-8 w-8 mx-auto mb-2 text-gray-300" />
                    No parts found
                  </div>
                ) : (
                  <div className="divide-y">
                    {filteredParts.map((part) => (
                      <button
                        key={part.id}
                        type="button"
                        onClick={() => setSelectedPart({ id: part.id, name: part.name, sku: part.puk })}
                        className="w-full p-3 text-left hover:bg-gray-50 flex items-center gap-3"
                      >
                        <div className="flex h-8 w-8 items-center justify-center rounded-full bg-amber-50">
                          <Package className="h-4 w-4 text-amber-600" />
                        </div>
                        <div>
                          <div className="font-medium text-gray-900">{part.name}</div>
                          <div className="text-sm text-gray-500">{part.puk}</div>
                        </div>
                      </button>
                    ))}
                  </div>
                )}
              </div>
            </>
          ) : (
            <>
              <div className="p-3 bg-gray-50 rounded-lg flex items-center justify-between">
                <div className="flex items-center gap-3">
                  <div className="flex h-10 w-10 items-center justify-center rounded-full bg-amber-50">
                    <Package className="h-5 w-5 text-amber-600" />
                  </div>
                  <div>
                    <div className="font-medium text-gray-900">{selectedPart.name}</div>
                    <div className="text-sm text-gray-500">{selectedPart.sku}</div>
                  </div>
                </div>
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => setSelectedPart(null)}
                >
                  Change
                </Button>
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Quantity on Hand *
                  </label>
                  <Input
                    type="number"
                    min="0"
                    value={qtyOnHand}
                    onChange={(e) => setQtyOnHand(parseInt(e.target.value) || 0)}
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Reorder Level
                  </label>
                  <Input
                    type="number"
                    min="0"
                    value={reorderLevel}
                    onChange={(e) => setReorderLevel(parseInt(e.target.value) || 0)}
                  />
                </div>
              </div>
              <p className="text-sm text-gray-500">
                When quantity falls below the reorder level, a low stock alert will be shown.
              </p>
            </>
          )}
        </div>
      </ModalBody>
      <ModalFooter>
        <Button variant="outline" onClick={handleClose}>
          Cancel
        </Button>
        <Button
          onClick={handleSubmit}
          disabled={!selectedPart || isLoading}
        >
          {isLoading ? 'Adding...' : 'Add to Inventory'}
        </Button>
      </ModalFooter>
    </Modal>
  );
}
