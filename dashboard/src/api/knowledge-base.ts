import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import api from './client';
import type {
  KBArticle,
  CreateKBArticleRequest,
  UpdateKBArticleRequest,
  KBArticleFilters,
  KBStats,
  PaginatedResponse,
} from '@/types';

const KB_KEY = 'kb-articles';
const KB_STATS_KEY = 'kb-stats';

// ============================================================================
// KB Articles CRUD
// ============================================================================

export function useKBArticles(filters: KBArticleFilters = {}) {
  return useQuery({
    queryKey: [KB_KEY, filters],
    queryFn: () => api.get<PaginatedResponse<KBArticle>>('/kb/articles', filters),
    staleTime: 30_000,
  });
}

export function useKBArticle(id: string) {
  return useQuery({
    queryKey: [KB_KEY, 'id', id],
    queryFn: () => api.get<KBArticle>(`/kb/articles/${id}`),
    enabled: !!id,
  });
}

export function useKBArticleBySlug(slug: string) {
  return useQuery({
    queryKey: [KB_KEY, 'slug', slug],
    queryFn: () => api.get<KBArticle>(`/kb/articles/slug/${slug}`),
    enabled: !!slug,
  });
}

export function useCreateKBArticle() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateKBArticleRequest) =>
      api.post<KBArticle>('/kb/articles', data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [KB_KEY] });
      queryClient.invalidateQueries({ queryKey: [KB_STATS_KEY] });
    },
  });
}

export function useUpdateKBArticle() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateKBArticleRequest }) =>
      api.patch<KBArticle>(`/kb/articles/${id}`, data),
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: [KB_KEY] });
      queryClient.invalidateQueries({ queryKey: [KB_KEY, 'id', id] });
      queryClient.invalidateQueries({ queryKey: [KB_STATS_KEY] });
    },
  });
}

export function useDeleteKBArticle() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => api.delete(`/kb/articles/${id}`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [KB_KEY] });
      queryClient.invalidateQueries({ queryKey: [KB_STATS_KEY] });
    },
  });
}

export function usePublishKBArticle() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => api.post<KBArticle>(`/kb/articles/${id}/publish`, {}),
    onSuccess: (_, id) => {
      queryClient.invalidateQueries({ queryKey: [KB_KEY] });
      queryClient.invalidateQueries({ queryKey: [KB_KEY, 'id', id] });
      queryClient.invalidateQueries({ queryKey: [KB_STATS_KEY] });
    },
  });
}

// ============================================================================
// KB Search & Stats
// ============================================================================

export function useKBSearch(query: string, limit = 20) {
  return useQuery({
    queryKey: [KB_KEY, 'search', query, limit],
    queryFn: () => api.get<{ items: KBArticle[] }>('/kb/search', { q: query, limit }),
    enabled: query.length >= 2,
    staleTime: 30_000,
  });
}

export function useKBStats() {
  return useQuery({
    queryKey: [KB_STATS_KEY],
    queryFn: () => api.get<KBStats>('/kb/stats'),
    staleTime: 60_000,
  });
}
