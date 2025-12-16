import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import api from './client';
import type {
  ServiceShop,
  CreateServiceShopRequest,
  ServiceStaff,
  CreateServiceStaffRequest,
  PaginatedResponse,
} from '@/types';

const SERVICE_SHOPS_KEY = 'service-shops';
const SERVICE_STAFF_KEY = 'service-staff';
const INVENTORY_KEY = 'inventory';

// Service Shops
export function useServiceShops(filters: { countyCode?: string; active?: boolean; limit?: number; cursor?: string } = {}) {
  return useQuery({
    queryKey: [SERVICE_SHOPS_KEY, filters],
    queryFn: () => api.get<PaginatedResponse<ServiceShop>>('/service-shops', filters),
    staleTime: 30_000,
  });
}

export function useServiceShop(id: string) {
  return useQuery({
    queryKey: [SERVICE_SHOPS_KEY, id],
    queryFn: () => api.get<ServiceShop>(`/service-shops/${id}`),
    enabled: !!id,
  });
}

export function useCreateServiceShop() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateServiceShopRequest) =>
      api.post<ServiceShop>('/service-shops', data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [SERVICE_SHOPS_KEY] });
    },
  });
}

export function useUpdateServiceShop() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<CreateServiceShopRequest> }) =>
      api.patch<ServiceShop>(`/service-shops/${id}`, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [SERVICE_SHOPS_KEY] });
    },
  });
}

export function useDeleteServiceShop() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => api.delete(`/service-shops/${id}`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [SERVICE_SHOPS_KEY] });
    },
  });
}

// Service Shops Stats
export function useServiceShopsStats() {
  return useQuery({
    queryKey: [SERVICE_SHOPS_KEY, 'stats'],
    queryFn: () => api.get<{
      total: number;
      active: number;
      totalStaff: number;
      countiesCovered: number;
      lowStockShops: number;
    }>('/service-shops/stats'),
    staleTime: 30_000,
  });
}

// Service Staff
export function useServiceStaff(filters: { serviceShopId?: string; role?: string; active?: boolean; limit?: number; cursor?: string } = {}) {
  return useQuery({
    queryKey: [SERVICE_STAFF_KEY, filters],
    queryFn: () => api.get<PaginatedResponse<ServiceStaff>>('/service-staff', filters),
    staleTime: 30_000,
  });
}

export function useServiceStaffMember(id: string) {
  return useQuery({
    queryKey: [SERVICE_STAFF_KEY, id],
    queryFn: () => api.get<ServiceStaff>(`/service-staff/${id}`),
    enabled: !!id,
  });
}

export function useCreateServiceStaff() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateServiceStaffRequest) =>
      api.post<ServiceStaff>('/service-staff', data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [SERVICE_STAFF_KEY] });
    },
  });
}

export function useUpdateServiceStaff() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<CreateServiceStaffRequest> }) =>
      api.patch<ServiceStaff>(`/service-staff/${id}`, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [SERVICE_STAFF_KEY] });
    },
  });
}

export function useDeleteServiceStaff() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => api.delete(`/service-staff/${id}`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [SERVICE_STAFF_KEY] });
    },
  });
}

// Inventory
interface InventoryItem {
  id: string;
  tenantId: string;
  serviceShopId: string;
  partId: string;
  partName: string;
  partPuk: string;
  partCategory: string;
  qtyOnHand: number;
  qtyReserved: number;
  qtyAvailable: number;
  reorderLevel: number;
  lastRestockedAt: string | null;
  createdAt: string;
  updatedAt: string;
}

export function useInventory(filters: { serviceShopId?: string; partId?: string; lowStock?: boolean; limit?: number; cursor?: string } = {}) {
  return useQuery({
    queryKey: [INVENTORY_KEY, filters],
    queryFn: () => api.get<PaginatedResponse<InventoryItem>>('/inventory', filters),
    staleTime: 30_000,
  });
}

export function useUpsertInventory() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      serviceShopId,
      partId,
      qtyOnHand,
      reorderLevel,
    }: {
      serviceShopId: string;
      partId: string;
      qtyOnHand: number;
      reorderLevel?: number;
    }) =>
      api.post<InventoryItem>('/inventory/upsert', {
        serviceShopId,
        partId,
        qtyOnHand,
        reorderLevel,
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [INVENTORY_KEY] });
    },
  });
}

// Parts (for selecting parts in inventory)
interface Part {
  id: string;
  tenantId: string;
  puk: string;
  name: string;
  category: string;
  deviceModelId: string;
  unitCostCents: number;
  active: boolean;
  createdAt: string;
  updatedAt: string;
}

export function useParts(filters: { category?: string; deviceModelId?: string; active?: boolean; limit?: number; cursor?: string } = {}) {
  return useQuery({
    queryKey: ['parts', filters],
    queryFn: () => api.get<PaginatedResponse<Part>>('/parts', filters),
    staleTime: 60_000,
  });
}

export function usePart(id: string) {
  return useQuery({
    queryKey: ['parts', id],
    queryFn: () => api.get<Part>(`/parts/${id}`),
    enabled: !!id,
  });
}
