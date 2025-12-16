import { Plus, FileCheck } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import type { WorkOrderDeliverable } from '@/types';

interface WorkOrderDeliverablesTabProps {
  deliverables: WorkOrderDeliverable[];
}

export function WorkOrderDeliverablesTab({ deliverables }: WorkOrderDeliverablesTabProps) {
  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h3 className="font-medium text-gray-900">Deliverables</h3>
        <Button size="sm">
          <Plus className="h-4 w-4" />
          Add Deliverable
        </Button>
      </div>
      {deliverables.length === 0 ? (
        <div className="text-center py-8 text-gray-500">
          <FileCheck className="h-12 w-12 mx-auto mb-2 text-gray-300" />
          <p>No deliverables defined</p>
        </div>
      ) : (
        <div className="space-y-2">
          {deliverables.map((deliverable) => (
            <div key={deliverable.id} className="p-3 bg-gray-50 rounded-lg">
              <div className="flex items-center justify-between mb-1">
                <span className="font-medium text-gray-900">
                  {deliverable.title}
                </span>
                <Badge
                  variant={
                    deliverable.status === 'approved'
                      ? 'success'
                      : deliverable.status === 'rejected'
                      ? 'destructive'
                      : 'outline'
                  }
                >
                  {deliverable.status}
                </Badge>
              </div>
              {deliverable.description && (
                <p className="text-sm text-gray-500">{deliverable.description}</p>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
