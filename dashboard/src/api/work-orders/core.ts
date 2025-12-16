import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import api from '../client';
import { toast } from '@/lib/toast';
import type {
  WorkOrder,
  CreateWorkOrderRequest,
  WorkOrderFilters,
  WorkOrderStatus,
  PaginatedResponse,
} from '@/types';
import { WORK_ORDERS_KEY } from './keys';

// ============================================================================
// Core Work Order Operations
// ============================================================================

/**
 * List work orders with optional filters
 */
export function useWorkOrders(filters: WorkOrderFilters = {}) {
  return useQuery({
    queryKey: [WORK_ORDERS_KEY, filters],
    queryFn: () => api.get<PaginatedResponse<WorkOrder>>('/work-orders', filters),
    staleTime: 30_000,
  });
}

/**
 * Get a single work order by ID
 */
export function useWorkOrder(id: string) {
  return useQuery({
    queryKey: [WORK_ORDERS_KEY, id],
    queryFn: () => api.get<WorkOrder>(`/work-orders/${id}`),
    enabled: !!id,
  });
}

/**
 * Create a new work order
 */
export function useCreateWorkOrder() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateWorkOrderRequest) =>
      api.post<WorkOrder>('/work-orders', data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [WORK_ORDERS_KEY] });
      toast.success('Work Order Created', 'The work order has been created successfully');
    },
  });
}

/**
 * Update work order status
 */
export function useUpdateWorkOrderStatus() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, status }: { id: string; status: WorkOrderStatus }) =>
      api.patch<WorkOrder>(`/work-orders/${id}/status`, { status }),
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: [WORK_ORDERS_KEY] });
      queryClient.setQueryData([WORK_ORDERS_KEY, data.id], data);
      toast.success('Status Updated', `Work order status changed to ${data.status.replace('_', ' ')}`);
    },
  });
}
