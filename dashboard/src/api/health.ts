import { useQuery } from '@tanstack/react-query';
import api from './client';
import type { ServiceHealth, DashboardMetrics, ActivityEvent } from '@/types';
import { HEALTH_CHECK_INTERVAL, METRICS_REFRESH_INTERVAL, ACTIVITY_REFRESH_INTERVAL } from '@/lib/constants';

// Service health
export function useServiceHealth() {
  return useQuery({
    queryKey: ['admin', 'health', 'services'],
    queryFn: () =>
      api.get<{ services: ServiceHealth[] }>('/v1/health/services'),
    refetchInterval: HEALTH_CHECK_INTERVAL,
    staleTime: 10_000, // 10 seconds
  });
}

// Dashboard metrics summary
export function useDashboardMetrics() {
  return useQuery({
    queryKey: ['admin', 'metrics', 'summary'],
    queryFn: () => api.get<DashboardMetrics>('/v1/metrics/summary'),
    refetchInterval: METRICS_REFRESH_INTERVAL,
    staleTime: 30_000, // 30 seconds
  });
}

// Activity feed
export function useActivityFeed(limit: number = 50) {
  return useQuery({
    queryKey: ['admin', 'activity', limit],
    queryFn: () =>
      api.get<ActivityEvent[]>('/v1/activity', { limit }),
    refetchInterval: ACTIVITY_REFRESH_INTERVAL,
    staleTime: 5_000, // 5 seconds
  });
}

// Individual service health check
export function useHealthCheck(service: string) {
  const endpoint = service === 'ims-api' ? '/healthz' : `/ssot/${service}/healthz`;

  return useQuery({
    queryKey: ['health', service],
    queryFn: async () => {
      const response = await fetch(endpoint);
      return response.ok;
    },
    refetchInterval: HEALTH_CHECK_INTERVAL,
  });
}

// Readiness check
export function useReadinessCheck() {
  return useQuery({
    queryKey: ['readyz'],
    queryFn: async () => {
      const response = await fetch('/readyz');
      return response.ok;
    },
    refetchInterval: HEALTH_CHECK_INTERVAL,
  });
}
