import api from '@/api/client';
import type {
  SalesDashboardResponse,
  SalesMetricsSummary,
  PipelineSummary,
} from '@/types/sales';

export const salesApi = {
  // Get full sales dashboard data
  getDashboard: (days?: number): Promise<SalesDashboardResponse> => {
    const params = days ? `?days=${days}` : '';
    return api.get<SalesDashboardResponse>(`/sales/dashboard${params}`);
  },

  // Get metrics summary for a date range
  getMetrics: (startDate?: string, endDate?: string): Promise<SalesMetricsSummary> => {
    const params = new URLSearchParams();
    if (startDate) params.set('startDate', startDate);
    if (endDate) params.set('endDate', endDate);
    const queryString = params.toString();
    return api.get<SalesMetricsSummary>(`/sales/metrics${queryString ? `?${queryString}` : ''}`);
  },

  // Get pipeline stage counts
  getPipelineStages: (): Promise<PipelineSummary> => {
    return api.get<PipelineSummary>('/sales/pipeline-stages');
  },

  // Increment a metric (for admin/testing)
  incrementMetric: (metric: string, value: number): Promise<void> => {
    return api.post('/sales/metrics/increment', { metric, value });
  },
};

export default salesApi;
