import { Plus, Package } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { cn } from '@/lib/utils';

interface InventoryItem {
  id: string;
  partName: string;
  partPuk: string;
  qtyOnHand: number;
  qtyAvailable: number;
  reorderLevel: number;
}

interface InventoryListProps {
  inventory: InventoryItem[];
  onAddClick?: () => void;
}

export function InventoryList({ inventory, onAddClick }: InventoryListProps) {
  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h3 className="font-medium text-gray-900">Inventory</h3>
        <Button size="sm" onClick={onAddClick}>
          <Plus className="h-4 w-4" />
          Add Stock
        </Button>
      </div>
      {inventory.length === 0 ? (
        <div className="text-center py-8 text-gray-500">
          <Package className="h-12 w-12 mx-auto mb-2 text-gray-300" />
          <p>No inventory items</p>
        </div>
      ) : (
        <div className="space-y-2">
          {inventory.map((item) => (
            <div
              key={item.id}
              className="flex items-center justify-between p-3 bg-gray-50 rounded-lg"
            >
              <div>
                <div className="font-medium text-gray-900">{item.partName}</div>
                <div className="text-sm text-gray-500">{item.partPuk}</div>
              </div>
              <div className="text-right">
                <div
                  className={cn(
                    'font-medium',
                    item.qtyAvailable <= item.reorderLevel
                      ? 'text-red-600'
                      : 'text-gray-900'
                  )}
                >
                  {item.qtyAvailable} / {item.qtyOnHand}
                </div>
                <div className="text-sm text-gray-500">available / on hand</div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
