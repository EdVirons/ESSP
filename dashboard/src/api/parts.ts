import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import api from './client';
import type {
  Part,
  CreatePartRequest,
  UpdatePartRequest,
  PartFilters,
  PartStats,
  ImportPartsResult,
  PaginatedResponse,
} from '@/types';

const PARTS_KEY = 'parts';
const PARTS_STATS_KEY = 'parts-stats';
const PARTS_CATEGORIES_KEY = 'parts-categories';

// ============================================================================
// Parts CRUD
// ============================================================================

export function useParts(filters: PartFilters = {}) {
  return useQuery({
    queryKey: [PARTS_KEY, filters],
    queryFn: () => api.get<PaginatedResponse<Part>>('/parts', filters),
    staleTime: 30_000,
  });
}

export function usePart(id: string) {
  return useQuery({
    queryKey: [PARTS_KEY, id],
    queryFn: () => api.get<Part>(`/parts/${id}`),
    enabled: !!id,
  });
}

export function useCreatePart() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreatePartRequest) =>
      api.post<Part>('/parts', data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [PARTS_KEY] });
      queryClient.invalidateQueries({ queryKey: [PARTS_STATS_KEY] });
      queryClient.invalidateQueries({ queryKey: [PARTS_CATEGORIES_KEY] });
    },
  });
}

export function useUpdatePart() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdatePartRequest }) =>
      api.patch<Part>(`/parts/${id}`, data),
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: [PARTS_KEY] });
      queryClient.invalidateQueries({ queryKey: [PARTS_KEY, id] });
      queryClient.invalidateQueries({ queryKey: [PARTS_STATS_KEY] });
      queryClient.invalidateQueries({ queryKey: [PARTS_CATEGORIES_KEY] });
    },
  });
}

export function useDeletePart() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => api.delete(`/parts/${id}`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [PARTS_KEY] });
      queryClient.invalidateQueries({ queryKey: [PARTS_STATS_KEY] });
      queryClient.invalidateQueries({ queryKey: [PARTS_CATEGORIES_KEY] });
    },
  });
}

// ============================================================================
// Parts Stats & Categories
// ============================================================================

export function usePartsStats() {
  return useQuery({
    queryKey: [PARTS_STATS_KEY],
    queryFn: () => api.get<PartStats>('/parts/stats'),
    staleTime: 60_000,
  });
}

export function usePartsCategories() {
  return useQuery({
    queryKey: [PARTS_CATEGORIES_KEY],
    queryFn: () => api.get<{ items: string[] }>('/parts/categories'),
    staleTime: 60_000,
  });
}

// ============================================================================
// Import/Export
// ============================================================================

export function useImportParts() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (file: File) => {
      const formData = new FormData();
      formData.append('file', file);
      const response = await fetch('/v1/parts/import', {
        method: 'POST',
        body: formData,
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('auth_token')}`,
          'X-Tenant-ID': localStorage.getItem('tenant_id') || 'demo-tenant',
        },
      });
      if (!response.ok) {
        throw new Error('Import failed');
      }
      return response.json() as Promise<ImportPartsResult>;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [PARTS_KEY] });
      queryClient.invalidateQueries({ queryKey: [PARTS_STATS_KEY] });
      queryClient.invalidateQueries({ queryKey: [PARTS_CATEGORIES_KEY] });
    },
  });
}

export function useExportParts() {
  return useMutation({
    mutationFn: async () => {
      const response = await fetch('/v1/parts/export', {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('auth_token')}`,
          'X-Tenant-ID': localStorage.getItem('tenant_id') || 'demo-tenant',
        },
      });
      if (!response.ok) {
        throw new Error('Export failed');
      }
      const blob = await response.blob();

      // Trigger download
      const url = window.URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      link.download = 'parts.csv';
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      window.URL.revokeObjectURL(url);
    },
  });
}

/**
 * Download import template CSV
 */
export function downloadPartsTemplate() {
  const headers = [
    'sku',
    'name',
    'category',
    'description',
    'unitCostCents',
    'supplier',
    'supplierSku',
  ];

  const csvContent = headers.join(',') + '\n' +
    'EXAMPLE-001,Example Part,Electronics,An example part description,1500,Acme Corp,ACM-001\n' +
    'EXAMPLE-002,Another Part,Cables,Another description,250,Tech Supply,TS-002';

  const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
  const url = window.URL.createObjectURL(blob);
  const link = document.createElement('a');
  link.href = url;
  link.download = 'parts-template.csv';
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
  window.URL.revokeObjectURL(url);
}
