import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import api from '../client';
import type { WorkOrderPart, PaginatedResponse } from '@/types';
import { WORK_ORDERS_KEY } from './keys';

// ============================================================================
// Bill of Materials (BOM) Operations
// ============================================================================

/**
 * Get BOM for a work order
 */
export function useWorkOrderBOM(workOrderId: string) {
  return useQuery({
    queryKey: [WORK_ORDERS_KEY, workOrderId, 'bom'],
    queryFn: () => api.get<PaginatedResponse<WorkOrderPart>>(`/work-orders/${workOrderId}/bom`),
    enabled: !!workOrderId,
  });
}

/**
 * Add an item to the BOM
 */
export function useAddBOMItem() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      workOrderId,
      partId,
      qtyPlanned,
      allowIncompatible = false,
    }: {
      workOrderId: string;
      partId: string;
      qtyPlanned: number;
      allowIncompatible?: boolean;
    }) =>
      api.post<WorkOrderPart>(
        `/work-orders/${workOrderId}/bom/items?allowIncompatible=${allowIncompatible}`,
        { partId, qtyPlanned }
      ),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: [WORK_ORDERS_KEY, variables.workOrderId, 'bom'],
      });
    },
  });
}

/**
 * Consume (use) a BOM item
 */
export function useConsumeBOMItem() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      workOrderId,
      itemId,
      qtyUsed,
    }: {
      workOrderId: string;
      itemId: string;
      qtyUsed: number;
    }) =>
      api.patch<WorkOrderPart>(
        `/work-orders/${workOrderId}/bom/items/${itemId}/consume`,
        { qtyUsed }
      ),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: [WORK_ORDERS_KEY, variables.workOrderId, 'bom'],
      });
    },
  });
}

/**
 * Release (return) a BOM item
 */
export function useReleaseBOMItem() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      workOrderId,
      itemId,
      qty,
    }: {
      workOrderId: string;
      itemId: string;
      qty: number;
    }) =>
      api.patch<WorkOrderPart>(
        `/work-orders/${workOrderId}/bom/items/${itemId}/release`,
        { qty }
      ),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: [WORK_ORDERS_KEY, variables.workOrderId, 'bom'],
      });
    },
  });
}
