import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import api from '../client';
import { toast } from '@/lib/toast';
import type {
  WorkOrderReworkHistory,
  RejectWorkOrderRequest,
  RejectWorkOrderResponse,
} from '@/types/work-order';
import { WORK_ORDERS_KEY } from './keys';

export const REWORK_HISTORY_KEY = 'work-order-rework-history';

/**
 * Reject a work order and send it back to a previous status
 */
export function useRejectWorkOrder() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: RejectWorkOrderRequest }) =>
      api.post<RejectWorkOrderResponse>(`/work-orders/${id}/reject`, data),
    onSuccess: (response, { id }) => {
      queryClient.invalidateQueries({ queryKey: [WORK_ORDERS_KEY] });
      queryClient.setQueryData([WORK_ORDERS_KEY, id], response.workOrder);
      queryClient.invalidateQueries({ queryKey: [REWORK_HISTORY_KEY, id] });
      toast.success('Work Order Rejected', `Sent back to ${response.workOrder.status.replace('_', ' ')} status`);
    },
    onError: (error: Error) => {
      toast.error('Rejection Failed', error.message || 'Failed to reject work order');
    },
  });
}

/**
 * Get rework history for a work order
 */
export function useWorkOrderReworkHistory(workOrderId: string, enabled = true) {
  return useQuery({
    queryKey: [REWORK_HISTORY_KEY, workOrderId],
    queryFn: () =>
      api.get<{ history: WorkOrderReworkHistory[]; nextCursor?: string }>(
        `/work-orders/${workOrderId}/rework-history`
      ),
    enabled: enabled && !!workOrderId,
    staleTime: 60_000,
  });
}
