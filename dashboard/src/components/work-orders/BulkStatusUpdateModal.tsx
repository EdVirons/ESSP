import { useState } from 'react';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { Label } from '@/components/ui/label';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Loader2, AlertCircle } from 'lucide-react';
import { useBulkStatusUpdate } from '@/api/work-orders';
import type { WorkOrderStatus } from '@/types/work-order';

interface BulkStatusUpdateModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  selectedIds: string[];
  onSuccess: () => void;
}

const statusOptions: { value: WorkOrderStatus; label: string }[] = [
  { value: 'assigned', label: 'Assigned' },
  { value: 'in_repair', label: 'In Repair' },
  { value: 'qa', label: 'QA' },
  { value: 'completed', label: 'Completed' },
];

export function BulkStatusUpdateModal({
  open,
  onOpenChange,
  selectedIds,
  onSuccess,
}: BulkStatusUpdateModalProps) {
  const bulkUpdate = useBulkStatusUpdate();
  const [status, setStatus] = useState<WorkOrderStatus | ''>('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!status) return;

    try {
      await bulkUpdate.mutateAsync({
        workOrderIds: selectedIds,
        status: status as WorkOrderStatus,
      });
      onOpenChange(false);
      setStatus('');
      onSuccess();
    } catch {
      // Error handled by mutation
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[425px]">
        <form onSubmit={handleSubmit}>
          <DialogHeader>
            <DialogTitle>Bulk Status Update</DialogTitle>
            <DialogDescription>
              Update the status of {selectedIds.length} selected work order
              {selectedIds.length > 1 ? 's' : ''}.
            </DialogDescription>
          </DialogHeader>

          <div className="py-4">
            <div className="flex items-start gap-3 p-3 bg-amber-50 border border-amber-200 rounded-lg mb-4">
              <AlertCircle className="h-5 w-5 text-amber-600 shrink-0 mt-0.5" />
              <p className="text-sm text-amber-800">
                Only valid forward status transitions will be applied. Work orders that cannot
                transition to the selected status will be skipped.
              </p>
            </div>

            <div className="grid gap-2">
              <Label htmlFor="status">New Status</Label>
              <Select
                value={status}
                onValueChange={(value) => setStatus(value as WorkOrderStatus)}
              >
                <SelectTrigger id="status">
                  <SelectValue placeholder="Select new status" />
                </SelectTrigger>
                <SelectContent>
                  {statusOptions.map((option) => (
                    <SelectItem key={option.value} value={option.value}>
                      {option.label}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
          </div>

          <DialogFooter>
            <Button
              type="button"
              variant="outline"
              onClick={() => onOpenChange(false)}
              disabled={bulkUpdate.isPending}
            >
              Cancel
            </Button>
            <Button type="submit" disabled={!status || bulkUpdate.isPending}>
              {bulkUpdate.isPending ? (
                <>
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  Updating...
                </>
              ) : (
                `Update ${selectedIds.length} Work Order${selectedIds.length > 1 ? 's' : ''}`
              )}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
