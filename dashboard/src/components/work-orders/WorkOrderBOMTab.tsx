import { Plus, Package } from 'lucide-react';
import { Button } from '@/components/ui/button';
import type { WorkOrderPart } from '@/types';

interface WorkOrderBOMTabProps {
  bomItems: WorkOrderPart[];
}

export function WorkOrderBOMTab({ bomItems }: WorkOrderBOMTabProps) {
  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h3 className="font-medium text-gray-900">Bill of Materials</h3>
        <Button size="sm">
          <Plus className="h-4 w-4" />
          Add Part
        </Button>
      </div>
      {bomItems.length === 0 ? (
        <div className="text-center py-8 text-gray-500">
          <Package className="h-12 w-12 mx-auto mb-2 text-gray-300" />
          <p>No parts added yet</p>
        </div>
      ) : (
        <div className="space-y-2">
          {bomItems.map((item) => (
            <div
              key={item.id}
              className="flex items-center justify-between p-3 bg-gray-50 rounded-lg"
            >
              <div>
                <div className="font-medium text-gray-900">{item.partName}</div>
                <div className="text-sm text-gray-500">{item.partPuk}</div>
              </div>
              <div className="text-right">
                <div className="font-medium">
                  {item.qtyUsed} / {item.qtyPlanned}
                </div>
                <div className="text-sm text-gray-500">used / planned</div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
