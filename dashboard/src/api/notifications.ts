import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import axios from 'axios';
import type { NotificationsResponse, UnreadCountResponse } from '@/types/notification';

const NOTIFICATIONS_KEY = 'notifications';

// Admin API client - doesn't use /v1 base URL
const adminApi = axios.create({
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add tenant headers
adminApi.interceptors.request.use((config) => {
  const tenantId = localStorage.getItem('tenant_id') || 'demo-tenant';
  config.headers['X-Tenant-ID'] = tenantId;
  return config;
});

// Get notifications list
export function useNotifications(options?: { limit?: number }) {
  const limit = options?.limit ?? 50;

  return useQuery({
    queryKey: [NOTIFICATIONS_KEY, { limit }],
    queryFn: () =>
      adminApi
        .get<NotificationsResponse>('/v1/notifications', { params: { limit } })
        .then((res) => res.data),
    staleTime: 30_000,
    refetchInterval: 60_000, // Refresh every minute
  });
}

// Get unread count
export function useUnreadCount() {
  return useQuery({
    queryKey: [NOTIFICATIONS_KEY, 'unread-count'],
    queryFn: () =>
      adminApi
        .get<UnreadCountResponse>('/v1/notifications/unread-count')
        .then((res) => res.data),
    staleTime: 30_000,
    refetchInterval: 30_000, // Refresh every 30 seconds
  });
}

// Mark notifications as read
export function useMarkNotificationsRead() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (ids: string) =>
      adminApi
        .post<{ success: boolean }>('/v1/notifications/mark-read', { ids })
        .then((res) => res.data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [NOTIFICATIONS_KEY] });
    },
  });
}
