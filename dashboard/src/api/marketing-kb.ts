import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import api from './client';
import type {
  MKBArticle,
  PitchKit,
  CreateMKBArticleRequest,
  UpdateMKBArticleRequest,
  CreatePitchKitRequest,
  UpdatePitchKitRequest,
  MKBArticleFilters,
  PitchKitFilters,
  MKBStats,
  PaginatedResponse,
} from '@/types';

const MKB_KEY = 'marketing-kb-articles';
const MKB_STATS_KEY = 'marketing-kb-stats';
const PITCH_KIT_KEY = 'pitch-kits';

// ============================================================================
// Marketing KB Articles CRUD
// ============================================================================

export function useMKBArticles(filters: MKBArticleFilters = {}) {
  return useQuery({
    queryKey: [MKB_KEY, filters],
    queryFn: () => api.get<PaginatedResponse<MKBArticle>>('/marketing-kb/articles', filters),
    staleTime: 30_000,
  });
}

export function useMKBArticle(id: string) {
  return useQuery({
    queryKey: [MKB_KEY, 'id', id],
    queryFn: () => api.get<MKBArticle>(`/marketing-kb/articles/${id}`),
    enabled: !!id,
  });
}

export function useMKBArticleBySlug(slug: string) {
  return useQuery({
    queryKey: [MKB_KEY, 'slug', slug],
    queryFn: () => api.get<MKBArticle>(`/marketing-kb/articles/slug/${slug}`),
    enabled: !!slug,
  });
}

export function useCreateMKBArticle() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateMKBArticleRequest) =>
      api.post<MKBArticle>('/marketing-kb/articles', data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [MKB_KEY] });
      queryClient.invalidateQueries({ queryKey: [MKB_STATS_KEY] });
    },
  });
}

export function useUpdateMKBArticle() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateMKBArticleRequest }) =>
      api.patch<MKBArticle>(`/marketing-kb/articles/${id}`, data),
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: [MKB_KEY] });
      queryClient.invalidateQueries({ queryKey: [MKB_KEY, 'id', id] });
      queryClient.invalidateQueries({ queryKey: [MKB_STATS_KEY] });
    },
  });
}

export function useDeleteMKBArticle() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => api.delete(`/marketing-kb/articles/${id}`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [MKB_KEY] });
      queryClient.invalidateQueries({ queryKey: [MKB_STATS_KEY] });
    },
  });
}

export function useApproveMKBArticle() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => api.post<MKBArticle>(`/marketing-kb/articles/${id}/approve`, {}),
    onSuccess: (_, id) => {
      queryClient.invalidateQueries({ queryKey: [MKB_KEY] });
      queryClient.invalidateQueries({ queryKey: [MKB_KEY, 'id', id] });
      queryClient.invalidateQueries({ queryKey: [MKB_STATS_KEY] });
    },
  });
}

export function useSubmitForReview() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => api.post<MKBArticle>(`/marketing-kb/articles/${id}/submit-review`, {}),
    onSuccess: (_, id) => {
      queryClient.invalidateQueries({ queryKey: [MKB_KEY] });
      queryClient.invalidateQueries({ queryKey: [MKB_KEY, 'id', id] });
      queryClient.invalidateQueries({ queryKey: [MKB_STATS_KEY] });
    },
  });
}

export function useRecordMKBUsage() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => api.post(`/marketing-kb/articles/${id}/usage`, {}),
    onSuccess: (_, id) => {
      queryClient.invalidateQueries({ queryKey: [MKB_KEY, 'id', id] });
    },
  });
}

// ============================================================================
// Marketing KB Search & Stats
// ============================================================================

export function useMKBSearch(query: string, limit = 20) {
  return useQuery({
    queryKey: [MKB_KEY, 'search', query, limit],
    queryFn: () => api.get<{ items: MKBArticle[] }>('/marketing-kb/search', { q: query, limit }),
    enabled: query.length >= 2,
    staleTime: 30_000,
  });
}

export function useMKBStats() {
  return useQuery({
    queryKey: [MKB_STATS_KEY],
    queryFn: () => api.get<MKBStats>('/marketing-kb/stats'),
    staleTime: 60_000,
  });
}

// ============================================================================
// Pitch Kits CRUD
// ============================================================================

export function usePitchKits(filters: PitchKitFilters = {}) {
  return useQuery({
    queryKey: [PITCH_KIT_KEY, filters],
    queryFn: () => api.get<PaginatedResponse<PitchKit>>('/marketing-kb/pitch-kits', filters),
    staleTime: 30_000,
  });
}

export function usePitchKit(id: string) {
  return useQuery({
    queryKey: [PITCH_KIT_KEY, 'id', id],
    queryFn: () => api.get<PitchKit>(`/marketing-kb/pitch-kits/${id}`),
    enabled: !!id,
  });
}

export function useCreatePitchKit() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreatePitchKitRequest) =>
      api.post<PitchKit>('/marketing-kb/pitch-kits', data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [PITCH_KIT_KEY] });
    },
  });
}

export function useUpdatePitchKit() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdatePitchKitRequest }) =>
      api.patch<PitchKit>(`/marketing-kb/pitch-kits/${id}`, data),
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: [PITCH_KIT_KEY] });
      queryClient.invalidateQueries({ queryKey: [PITCH_KIT_KEY, 'id', id] });
    },
  });
}

export function useDeletePitchKit() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => api.delete(`/marketing-kb/pitch-kits/${id}`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [PITCH_KIT_KEY] });
    },
  });
}
