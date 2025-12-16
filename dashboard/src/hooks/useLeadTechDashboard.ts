import { useQuery } from '@tanstack/react-query';
import { api } from '@/api/client';

// Types for lead tech dashboard data
export interface PendingApproval {
  id: string;
  workOrderId: string;
  title: string;
  schoolName: string;
  priority: string;
  requestedAt: string;
  requestedByName: string;
}

export interface ScheduledWorkOrder {
  id: string;
  title: string;
  schoolName: string;
  scheduledStart: string | null;
  assignedTo: string;
  status: string;
  priority: string;
}

export interface TeamWorkOrderMetrics {
  inProgress: number;
  completed: number;
  pending: number;
  scheduled: number;
}

export interface BOMReadinessItem {
  workOrderId: string;
  workOrderTitle: string;
  totalParts: number;
  availableParts: number;
  missingParts: number;
  isReady: boolean;
}

export interface TeamActivityItem {
  id: string;
  actorName: string;
  action: string;
  entityType: string;
  description: string;
  createdAt: string;
}

export interface LeadTechDashboardSummary {
  pendingApprovalsCount: number;
  todaysScheduledCount: number;
  teamMetrics: TeamWorkOrderMetrics;
  pendingApprovals: PendingApproval[];
  todaysSchedule: ScheduledWorkOrder[];
  bomReadiness: BOMReadinessItem[];
  recentTeamActivity: TeamActivityItem[];
}

// Query keys
export const leadTechKeys = {
  all: ['leadtech'] as const,
  dashboard: () => [...leadTechKeys.all, 'dashboard'] as const,
  approvals: () => [...leadTechKeys.all, 'approvals'] as const,
  schedule: () => [...leadTechKeys.all, 'schedule'] as const,
  metrics: () => [...leadTechKeys.all, 'metrics'] as const,
};

// Main dashboard summary hook
export function useLeadTechDashboard() {
  return useQuery({
    queryKey: leadTechKeys.dashboard(),
    queryFn: () => api.get<LeadTechDashboardSummary>('/leadtech/dashboard'),
    staleTime: 30_000, // Consider data stale after 30 seconds
    refetchInterval: 60_000, // Auto refresh every minute
  });
}

// Pending approvals hook
export function usePendingApprovals(limit = 10) {
  return useQuery({
    queryKey: [...leadTechKeys.approvals(), limit],
    queryFn: () => api.get<{ items: PendingApproval[] }>('/leadtech/approvals', { limit }),
  });
}

// Today's schedule hook
export function useTodaysSchedule(limit = 20) {
  return useQuery({
    queryKey: [...leadTechKeys.schedule(), limit],
    queryFn: () => api.get<{ items: ScheduledWorkOrder[] }>('/leadtech/schedule', { limit }),
  });
}

// Team metrics hook
export function useTeamMetrics() {
  return useQuery({
    queryKey: leadTechKeys.metrics(),
    queryFn: () => api.get<TeamWorkOrderMetrics>('/leadtech/team-metrics'),
  });
}
