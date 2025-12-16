import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import api from '../client';
import type { WorkOrderSchedule, PaginatedResponse } from '@/types';
import { WORK_ORDERS_KEY } from './keys';

// ============================================================================
// Scheduling Operations
// ============================================================================

/**
 * Get schedules for a work order
 */
export function useWorkOrderSchedules(workOrderId: string) {
  return useQuery({
    queryKey: [WORK_ORDERS_KEY, workOrderId, 'schedules'],
    queryFn: () =>
      api.get<PaginatedResponse<WorkOrderSchedule>>(`/work-orders/${workOrderId}/schedules`),
    enabled: !!workOrderId,
  });
}

/**
 * Create a schedule for a work order
 */
export function useCreateSchedule() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      workOrderId,
      scheduledStart,
      scheduledEnd,
      timezone = 'Africa/Nairobi',
      notes,
    }: {
      workOrderId: string;
      scheduledStart?: string;
      scheduledEnd?: string;
      timezone?: string;
      notes?: string;
    }) =>
      api.post<WorkOrderSchedule>(`/work-orders/${workOrderId}/schedule`, {
        scheduledStart,
        scheduledEnd,
        timezone,
        notes,
      }),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: [WORK_ORDERS_KEY, variables.workOrderId, 'schedules'],
      });
    },
  });
}
