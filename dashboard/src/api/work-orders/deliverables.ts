import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import api from '../client';
import type { WorkOrderDeliverable, PaginatedResponse } from '@/types';
import { WORK_ORDERS_KEY } from './keys';

// ============================================================================
// Deliverables Operations
// ============================================================================

/**
 * Get deliverables for a work order
 */
export function useWorkOrderDeliverables(workOrderId: string) {
  return useQuery({
    queryKey: [WORK_ORDERS_KEY, workOrderId, 'deliverables'],
    queryFn: () =>
      api.get<PaginatedResponse<WorkOrderDeliverable>>(`/work-orders/${workOrderId}/deliverables`),
    enabled: !!workOrderId,
  });
}

/**
 * Add a deliverable to a work order
 */
export function useAddDeliverable() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      workOrderId,
      title,
      description,
    }: {
      workOrderId: string;
      title: string;
      description?: string;
    }) =>
      api.post<WorkOrderDeliverable>(`/work-orders/${workOrderId}/deliverables`, {
        title,
        description,
      }),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: [WORK_ORDERS_KEY, variables.workOrderId, 'deliverables'],
      });
    },
  });
}

/**
 * Submit a deliverable with evidence
 */
export function useSubmitDeliverable() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      workOrderId,
      deliverableId,
      evidenceAttachmentId,
      notes,
    }: {
      workOrderId: string;
      deliverableId: string;
      evidenceAttachmentId: string;
      notes?: string;
    }) =>
      api.patch<WorkOrderDeliverable>(
        `/work-orders/${workOrderId}/deliverables/${deliverableId}/submit`,
        { evidenceAttachmentId, notes }
      ),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: [WORK_ORDERS_KEY, variables.workOrderId, 'deliverables'],
      });
    },
  });
}

/**
 * Review (approve/reject) a deliverable
 */
export function useReviewDeliverable() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      workOrderId,
      deliverableId,
      status,
      notes,
    }: {
      workOrderId: string;
      deliverableId: string;
      status: 'approved' | 'rejected';
      notes?: string;
    }) =>
      api.patch<WorkOrderDeliverable>(
        `/work-orders/${workOrderId}/deliverables/${deliverableId}/review`,
        { status, notes }
      ),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: [WORK_ORDERS_KEY, variables.workOrderId, 'deliverables'],
      });
    },
  });
}
