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
import { Textarea } from '@/components/ui/textarea';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Loader2, AlertTriangle } from 'lucide-react';
import { useRejectWorkOrder } from '@/api/work-orders';
import type { WorkOrder, WorkOrderStatus, RejectionCategory } from '@/types/work-order';

interface RejectWorkOrderModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  workOrder: WorkOrder;
}

// Valid backward transitions
const validReworkTransitions: Record<WorkOrderStatus, WorkOrderStatus[]> = {
  draft: [],
  assigned: ['draft'],
  in_repair: ['assigned', 'draft'],
  qa: ['in_repair', 'assigned'],
  completed: ['qa', 'in_repair'],
  approved: ['completed'],
};

const rejectionCategories: { value: RejectionCategory; label: string }[] = [
  { value: 'quality', label: 'Quality Issues' },
  { value: 'incomplete', label: 'Incomplete Work' },
  { value: 'wrong_parts', label: 'Wrong Parts Used' },
  { value: 'safety', label: 'Safety Concern' },
  { value: 'other', label: 'Other' },
];

const statusLabels: Record<WorkOrderStatus, string> = {
  draft: 'Draft',
  assigned: 'Assigned',
  in_repair: 'In Repair',
  qa: 'QA',
  completed: 'Completed',
  approved: 'Approved',
};

export function RejectWorkOrderModal({
  open,
  onOpenChange,
  workOrder,
}: RejectWorkOrderModalProps) {
  const rejectMutation = useRejectWorkOrder();

  const availableTargets = validReworkTransitions[workOrder.status] || [];
  const [targetStatus, setTargetStatus] = useState<WorkOrderStatus | ''>('');
  const [category, setCategory] = useState<RejectionCategory>('quality');
  const [reason, setReason] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!targetStatus || !reason.trim()) return;

    try {
      await rejectMutation.mutateAsync({
        id: workOrder.id,
        data: {
          targetStatus: targetStatus as WorkOrderStatus,
          reason: reason.trim(),
          category,
        },
      });
      onOpenChange(false);
      setTargetStatus('');
      setCategory('quality');
      setReason('');
    } catch {
      // Error handled by mutation
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[500px]">
        <form onSubmit={handleSubmit}>
          <DialogHeader>
            <DialogTitle className="flex items-center gap-2">
              <AlertTriangle className="h-5 w-5 text-orange-500" />
              Reject Work Order
            </DialogTitle>
            <DialogDescription>
              Send this work order back to a previous stage for rework.
              {workOrder.reworkCount > 0 && (
                <span className="block mt-1 text-orange-600">
                  This work order has been rejected {workOrder.reworkCount} time(s) already.
                </span>
              )}
            </DialogDescription>
          </DialogHeader>

          <div className="grid gap-4 py-4">
            <div className="grid gap-2">
              <Label htmlFor="targetStatus">Send Back To</Label>
              <Select
                value={targetStatus}
                onValueChange={(value) => setTargetStatus(value as WorkOrderStatus)}
              >
                <SelectTrigger id="targetStatus">
                  <SelectValue placeholder="Select target status" />
                </SelectTrigger>
                <SelectContent>
                  {availableTargets.map((status) => (
                    <SelectItem key={status} value={status}>
                      {statusLabels[status]}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>

            <div className="grid gap-2">
              <Label htmlFor="category">Rejection Category</Label>
              <Select
                value={category}
                onValueChange={(value) => setCategory(value as RejectionCategory)}
              >
                <SelectTrigger id="category">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  {rejectionCategories.map((cat) => (
                    <SelectItem key={cat.value} value={cat.value}>
                      {cat.label}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>

            <div className="grid gap-2">
              <Label htmlFor="reason">Rejection Reason</Label>
              <Textarea
                id="reason"
                placeholder="Explain why this work order is being rejected..."
                value={reason}
                onChange={(e) => setReason(e.target.value)}
                rows={4}
                required
              />
            </div>
          </div>

          <DialogFooter>
            <Button
              type="button"
              variant="outline"
              onClick={() => onOpenChange(false)}
              disabled={rejectMutation.isPending}
            >
              Cancel
            </Button>
            <Button
              type="submit"
              variant="destructive"
              disabled={!targetStatus || !reason.trim() || rejectMutation.isPending}
            >
              {rejectMutation.isPending ? (
                <>
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  Rejecting...
                </>
              ) : (
                'Reject Work Order'
              )}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
