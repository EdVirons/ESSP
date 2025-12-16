import { useMutation, useQueryClient } from '@tanstack/react-query';
import api from '../client';
import type { WorkOrderApproval } from '@/types';
import { WORK_ORDERS_KEY } from './keys';

// ============================================================================
// Approval Operations
// ============================================================================

/**
 * Request approval for a work order
 */
export function useRequestApproval() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      workOrderId,
      approvalType = 'school_signoff',
    }: {
      workOrderId: string;
      approvalType?: string;
    }) =>
      api.post<WorkOrderApproval>(`/work-orders/${workOrderId}/approvals`, {
        approvalType,
      }),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: [WORK_ORDERS_KEY, variables.workOrderId],
      });
    },
  });
}

/**
 * Make a decision on an approval request
 */
export function useDecideApproval() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      workOrderId,
      approvalId,
      status,
      notes,
    }: {
      workOrderId: string;
      approvalId: string;
      status: 'approved' | 'rejected';
      notes?: string;
    }) =>
      api.patch<WorkOrderApproval>(
        `/work-orders/${workOrderId}/approvals/${approvalId}/decide`,
        { status, notes }
      ),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: [WORK_ORDERS_KEY, variables.workOrderId],
      });
    },
  });
}
