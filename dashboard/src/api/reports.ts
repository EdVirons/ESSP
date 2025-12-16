import { useQuery } from '@tanstack/react-query';
import api from './client';
import type {
  ReportFilters,
  WorkOrderReportResponse,
  IncidentReportResponse,
  InventoryReportResponse,
  SchoolReportResponse,
  ExecutiveDashboard,
} from '@/types/reports';

const REPORTS_KEY = 'reports';

// Work Orders Report
export function useWorkOrdersReport(filters: ReportFilters = {}) {
  return useQuery({
    queryKey: [REPORTS_KEY, 'work-orders', filters],
    queryFn: () => api.get<WorkOrderReportResponse>('/reports/work-orders', filters),
    staleTime: 60_000, // 1 minute
  });
}

// Incidents Report
export function useIncidentsReport(filters: ReportFilters = {}) {
  return useQuery({
    queryKey: [REPORTS_KEY, 'incidents', filters],
    queryFn: () => api.get<IncidentReportResponse>('/reports/incidents', filters),
    staleTime: 60_000,
  });
}

// Inventory Report
export function useInventoryReport(filters: ReportFilters = {}) {
  return useQuery({
    queryKey: [REPORTS_KEY, 'inventory', filters],
    queryFn: () => api.get<InventoryReportResponse>('/reports/inventory', filters),
    staleTime: 60_000,
  });
}

// Schools Report
export function useSchoolsReport(filters: ReportFilters = {}) {
  return useQuery({
    queryKey: [REPORTS_KEY, 'schools', filters],
    queryFn: () => api.get<SchoolReportResponse>('/reports/schools', filters),
    staleTime: 60_000,
  });
}

// Executive Dashboard
export function useExecutiveDashboard() {
  return useQuery({
    queryKey: [REPORTS_KEY, 'executive'],
    queryFn: () => api.get<ExecutiveDashboard>('/reports/executive'),
    staleTime: 60_000,
  });
}
