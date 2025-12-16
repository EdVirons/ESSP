import { useWorkOrderReworkHistory } from '@/api/work-orders';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Loader2, History, ArrowRight, AlertTriangle } from 'lucide-react';
import { format } from 'date-fns';
import type { WorkOrderReworkHistory, RejectionCategory } from '@/types/work-order';

interface ReworkHistoryTabProps {
  workOrderId: string;
}

const statusLabels: Record<string, string> = {
  draft: 'Draft',
  assigned: 'Assigned',
  in_repair: 'In Repair',
  qa: 'QA',
  completed: 'Completed',
  approved: 'Approved',
};

const categoryLabels: Record<RejectionCategory, string> = {
  quality: 'Quality Issues',
  incomplete: 'Incomplete Work',
  wrong_parts: 'Wrong Parts',
  safety: 'Safety Concern',
  other: 'Other',
};

const categoryColors: Record<RejectionCategory, string> = {
  quality: 'bg-orange-100 text-orange-800',
  incomplete: 'bg-yellow-100 text-yellow-800',
  wrong_parts: 'bg-red-100 text-red-800',
  safety: 'bg-purple-100 text-purple-800',
  other: 'bg-gray-100 text-gray-800',
};

export function ReworkHistoryTab({ workOrderId }: ReworkHistoryTabProps) {
  const { data, isLoading, error } = useWorkOrderReworkHistory(workOrderId);

  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-12">
        <Loader2 className="h-8 w-8 animate-spin text-gray-400" />
      </div>
    );
  }

  if (error) {
    return (
      <div className="text-center py-12 text-red-600">
        Failed to load rework history
      </div>
    );
  }

  const history = data?.history || [];

  if (history.length === 0) {
    return (
      <div className="text-center py-12">
        <History className="h-12 w-12 mx-auto text-gray-300 mb-4" />
        <p className="text-gray-500">No rework history</p>
        <p className="text-sm text-gray-400">
          This work order has not been rejected or sent back for rework.
        </p>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      <div className="flex items-center gap-2 text-amber-600 bg-amber-50 border border-amber-200 rounded-lg px-4 py-3">
        <AlertTriangle className="h-5 w-5" />
        <span className="text-sm font-medium">
          This work order has been rejected {history.length} time{history.length > 1 ? 's' : ''}
        </span>
      </div>

      <div className="space-y-3">
        {history.map((entry: WorkOrderReworkHistory) => (
          <Card key={entry.id}>
            <CardHeader className="pb-2">
              <div className="flex items-center justify-between">
                <CardTitle className="text-sm font-medium flex items-center gap-2">
                  <span className="text-gray-500">#{entry.reworkSequence}</span>
                  <span className="text-gray-700">
                    {statusLabels[entry.fromStatus]}
                  </span>
                  <ArrowRight className="h-4 w-4 text-gray-400" />
                  <span className="text-gray-700">
                    {statusLabels[entry.toStatus]}
                  </span>
                </CardTitle>
                <Badge
                  variant="secondary"
                  className={categoryColors[entry.rejectionCategory]}
                >
                  {categoryLabels[entry.rejectionCategory]}
                </Badge>
              </div>
            </CardHeader>
            <CardContent>
              <p className="text-sm text-gray-600 mb-3">
                {entry.rejectionReason}
              </p>
              <div className="flex items-center justify-between text-xs text-gray-500">
                <span>By: {entry.rejectedByName || 'Unknown'}</span>
                <span>{format(new Date(entry.createdAt), 'MMM d, yyyy h:mm a')}</span>
              </div>
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  );
}
