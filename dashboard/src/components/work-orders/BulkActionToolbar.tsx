import { useState } from 'react';
import { Button } from '@/components/ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { X, ChevronDown, Play, UserCheck, CheckCircle2, Loader2 } from 'lucide-react';
import { useBulkApproval } from '@/api/work-orders';
import type { WorkOrderStatus } from '@/types/work-order';
import { BulkStatusUpdateModal } from './BulkStatusUpdateModal';
import { BulkAssignModal } from './BulkAssignModal';

interface BulkActionToolbarProps {
  selectedIds: string[];
  onClearSelection: () => void;
  currentStatuses: WorkOrderStatus[];
}

export function BulkActionToolbar({
  selectedIds,
  onClearSelection,
  currentStatuses,
}: BulkActionToolbarProps) {
  const [statusModalOpen, setStatusModalOpen] = useState(false);
  const [assignModalOpen, setAssignModalOpen] = useState(false);

  const bulkApproval = useBulkApproval();

  const count = selectedIds.length;
  if (count === 0) return null;

  // Check if all selected are in completed status for approval
  const allCompleted = currentStatuses.every((s) => s === 'completed');

  const handleBulkApprove = async () => {
    await bulkApproval.mutateAsync({
      workOrderIds: selectedIds,
      decision: 'approved',
    });
    onClearSelection();
  };

  return (
    <>
      <div className="fixed bottom-6 left-1/2 -translate-x-1/2 z-50">
        <div className="bg-gray-900 text-white rounded-lg shadow-2xl px-4 py-3 flex items-center gap-4">
          <span className="text-sm font-medium">
            {count} work order{count > 1 ? 's' : ''} selected
          </span>

          <div className="h-6 w-px bg-gray-700" />

          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="secondary" size="sm">
                <Play className="h-4 w-4 mr-2" />
                Change Status
                <ChevronDown className="h-4 w-4 ml-2" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent>
              <DropdownMenuItem onClick={() => setStatusModalOpen(true)}>
                Update Status...
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>

          <Button
            variant="secondary"
            size="sm"
            onClick={() => setAssignModalOpen(true)}
          >
            <UserCheck className="h-4 w-4 mr-2" />
            Assign
          </Button>

          {allCompleted && (
            <Button
              variant="default"
              size="sm"
              className="bg-green-600 hover:bg-green-700"
              onClick={handleBulkApprove}
              disabled={bulkApproval.isPending}
            >
              {bulkApproval.isPending ? (
                <Loader2 className="h-4 w-4 mr-2 animate-spin" />
              ) : (
                <CheckCircle2 className="h-4 w-4 mr-2" />
              )}
              Approve All
            </Button>
          )}

          <div className="h-6 w-px bg-gray-700" />

          <Button
            variant="ghost"
            size="sm"
            className="text-gray-400 hover:text-white"
            onClick={onClearSelection}
          >
            <X className="h-4 w-4" />
          </Button>
        </div>
      </div>

      <BulkStatusUpdateModal
        open={statusModalOpen}
        onOpenChange={setStatusModalOpen}
        selectedIds={selectedIds}
        onSuccess={onClearSelection}
      />

      <BulkAssignModal
        open={assignModalOpen}
        onOpenChange={setAssignModalOpen}
        selectedIds={selectedIds}
        onSuccess={onClearSelection}
      />
    </>
  );
}
