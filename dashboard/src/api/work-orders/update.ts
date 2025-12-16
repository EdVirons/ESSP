import { useMutation, useQueryClient } from '@tanstack/react-query';
import api from '../client';
import { toast } from '@/lib/toast';
import type {
  UpdateWorkOrderRequest,
  UpdateWorkOrderResponse,
} from '@/types/work-order';
import { WORK_ORDERS_KEY } from './keys';

/**
 * Update a work order's fields (PATCH update)
 */
export function useUpdateWorkOrder() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateWorkOrderRequest }) =>
      api.patch<UpdateWorkOrderResponse>(`/work-orders/${id}`, data),
    onMutate: async ({ id, data }) => {
      // Optimistic update
      await queryClient.cancelQueries({ queryKey: [WORK_ORDERS_KEY, id] });
      const previousData = queryClient.getQueryData([WORK_ORDERS_KEY, id]);
      queryClient.setQueryData([WORK_ORDERS_KEY, id], (old: unknown) => {
        if (!old) return old;
        return { ...old, ...data };
      });
      return { previousData };
    },
    onSuccess: (response, { id }) => {
      queryClient.invalidateQueries({ queryKey: [WORK_ORDERS_KEY] });
      queryClient.setQueryData([WORK_ORDERS_KEY, id], response.workOrder);
      toast.success('Work Order Updated', `Updated: ${response.updatedFields.join(', ')}`);
    },
    onError: (error: Error, { id }, context) => {
      // Rollback optimistic update
      if (context?.previousData) {
        queryClient.setQueryData([WORK_ORDERS_KEY, id], context.previousData);
      }
      toast.error('Update Failed', error.message || 'Failed to update work order');
    },
  });
}
