import { useQuery } from '@tanstack/react-query';
import { api } from '@/api/client';

// Types for warehouse dashboard data
export interface LowStockItem {
  id: string;
  partId: string;
  partName: string;
  shopId: string;
  shopName: string;
  qtyAvailable: number;
  reorderThreshold: number;
}

export interface PendingPartIssue {
  id: string;
  title: string;
  schoolName: string;
  priority: string;
  partsNeeded: number;
  createdAt: string;
}

export interface InventoryActivity {
  id: string;
  type: 'receipt' | 'issue' | 'adjustment' | string;
  description: string;
  partName: string;
  actorName: string;
  qtyChange: number;
  createdAt: string;
}

export interface WarehouseDashboardSummary {
  lowStockCount: number;
  pendingWorkOrders: number;
  todayMovements: number;
  totalParts: number;
  partsCategories: Record<string, number>;
  recentActivity: InventoryActivity[];
  lowStockItems: LowStockItem[];
  pendingPartIssues: PendingPartIssue[];
}

// Query keys
export const warehouseKeys = {
  all: ['warehouse'] as const,
  dashboard: () => [...warehouseKeys.all, 'dashboard'] as const,
  lowStock: () => [...warehouseKeys.all, 'low-stock'] as const,
  pendingIssues: () => [...warehouseKeys.all, 'pending-issues'] as const,
  movements: () => [...warehouseKeys.all, 'movements'] as const,
};

// Main dashboard summary hook
export function useWarehouseDashboard() {
  return useQuery({
    queryKey: warehouseKeys.dashboard(),
    queryFn: () => api.get<WarehouseDashboardSummary>('/warehouse/dashboard'),
    staleTime: 30_000, // Consider data stale after 30 seconds
    refetchInterval: 60_000, // Auto refresh every minute
  });
}

// Low stock items hook
export function useLowStockItems(limit = 20) {
  return useQuery({
    queryKey: [...warehouseKeys.lowStock(), limit],
    queryFn: () => api.get<{ items: LowStockItem[] }>('/warehouse/low-stock', { limit }),
  });
}

// Pending part issues hook
export function usePendingPartIssues(limit = 20) {
  return useQuery({
    queryKey: [...warehouseKeys.pendingIssues(), limit],
    queryFn: () => api.get<{ items: PendingPartIssue[] }>('/warehouse/pending-issues', { limit }),
  });
}

// Stock movements hook
export function useStockMovements(limit = 20) {
  return useQuery({
    queryKey: [...warehouseKeys.movements(), limit],
    queryFn: () => api.get<{ items: InventoryActivity[] }>('/warehouse/movements', { limit }),
  });
}
