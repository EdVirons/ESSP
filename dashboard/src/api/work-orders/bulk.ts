import { useMutation, useQueryClient } from '@tanstack/react-query';
import api from '../client';
import { toast } from '@/lib/toast';
import type {
  BulkStatusUpdateRequest,
  BulkAssignmentRequest,
  BulkApprovalRequest,
  BulkOperationResult,
} from '@/types/work-order';
import { WORK_ORDERS_KEY } from './keys';

/**
 * Bulk update status for multiple work orders
 */
export function useBulkStatusUpdate() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: BulkStatusUpdateRequest) =>
      api.post<BulkOperationResult>('/work-orders/bulk/status', data),
    onSuccess: (result) => {
      queryClient.invalidateQueries({ queryKey: [WORK_ORDERS_KEY] });
      if (result.failureCount === 0) {
        toast.success('Bulk Update Complete', `Updated ${result.successCount} work orders`);
      } else {
        toast.warning(
          'Bulk Update Partial',
          `${result.successCount} succeeded, ${result.failureCount} failed`
        );
      }
    },
    onError: (error: Error) => {
      toast.error('Bulk Update Failed', error.message || 'Failed to update work orders');
    },
  });
}

/**
 * Bulk update assignment for multiple work orders
 */
export function useBulkAssignment() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: BulkAssignmentRequest) =>
      api.post<BulkOperationResult>('/work-orders/bulk/assignment', data),
    onSuccess: (result) => {
      queryClient.invalidateQueries({ queryKey: [WORK_ORDERS_KEY] });
      if (result.failureCount === 0) {
        toast.success('Bulk Assignment Complete', `Assigned ${result.successCount} work orders`);
      } else {
        toast.warning(
          'Bulk Assignment Partial',
          `${result.successCount} succeeded, ${result.failureCount} failed`
        );
      }
    },
    onError: (error: Error) => {
      toast.error('Bulk Assignment Failed', error.message || 'Failed to assign work orders');
    },
  });
}

/**
 * Bulk approve/reject multiple work orders
 */
export function useBulkApproval() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: BulkApprovalRequest) =>
      api.post<BulkOperationResult>('/work-orders/bulk/approval', data),
    onSuccess: (result, variables) => {
      queryClient.invalidateQueries({ queryKey: [WORK_ORDERS_KEY] });
      const action = variables.decision === 'approved' ? 'approved' : 'rejected';
      if (result.failureCount === 0) {
        toast.success('Bulk Approval Complete', `${action} ${result.successCount} work orders`);
      } else {
        toast.warning(
          'Bulk Approval Partial',
          `${result.successCount} ${action}, ${result.failureCount} failed`
        );
      }
    },
    onError: (error: Error) => {
      toast.error('Bulk Approval Failed', error.message || 'Failed to process approvals');
    },
  });
}
