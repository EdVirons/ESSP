import { useQuery } from '@tanstack/react-query';
import { api } from '@/api/client';

// Types for support agent dashboard data
export interface IncidentQueueItem {
  id: string;
  title: string;
  schoolName: string;
  category: string;
  severity: string;
  status: string;
  reportedBy: string;
  slaDueAt: string | null;
  slaBreached: boolean;
  createdAt: string;
}

export interface ChatQueueItem {
  id: string;
  schoolName: string;
  contactName: string;
  status: string;
  queuePosition: number | null;
  assignedAgentName: string | null;
  startedAt: string;
  waitTimeSeconds: number;
}

export interface WorkOrderQueueItem {
  id: string;
  title: string;
  schoolName: string;
  status: string;
  taskType: string;
  assignedTo: string;
  createdAt: string;
}

export interface SupportAgentActivity {
  id: string;
  actorName: string;
  action: string;
  entityType: string;
  description: string;
  createdAt: string;
}

export interface IncidentMetrics {
  open: number;
  inProgress: number;
  resolved: number;
  slaBreached: number;
}

export interface SupportAgentDashboardSummary {
  openIncidentsCount: number;
  waitingChatsCount: number;
  activeChatsCount: number;
  activeWorkOrders: number;
  unreadMessagesCount: number;
  incidentMetrics: IncidentMetrics;
  incidentQueue: IncidentQueueItem[];
  chatQueue: ChatQueueItem[];
  workOrderQueue: WorkOrderQueueItem[];
  recentActivity: SupportAgentActivity[];
}

// Query keys
export const supportAgentKeys = {
  all: ['supportagent'] as const,
  dashboard: () => [...supportAgentKeys.all, 'dashboard'] as const,
  incidents: () => [...supportAgentKeys.all, 'incidents'] as const,
  chats: () => [...supportAgentKeys.all, 'chats'] as const,
  workOrders: () => [...supportAgentKeys.all, 'workOrders'] as const,
  metrics: () => [...supportAgentKeys.all, 'metrics'] as const,
};

// Main dashboard summary hook
export function useSupportAgentDashboard() {
  return useQuery({
    queryKey: supportAgentKeys.dashboard(),
    queryFn: () => api.get<SupportAgentDashboardSummary>('/supportagent/dashboard'),
    staleTime: 15_000, // Consider data stale after 15 seconds (support agents need fresher data)
    refetchInterval: 30_000, // Auto refresh every 30 seconds
  });
}

// Incident queue hook
export function useIncidentQueue(limit = 20) {
  return useQuery({
    queryKey: [...supportAgentKeys.incidents(), limit],
    queryFn: () => api.get<{ items: IncidentQueueItem[] }>('/supportagent/incidents', { limit }),
  });
}

// Chat queue hook
export function useChatQueue(limit = 20) {
  return useQuery({
    queryKey: [...supportAgentKeys.chats(), limit],
    queryFn: () => api.get<{ items: ChatQueueItem[] }>('/supportagent/chats', { limit }),
    refetchInterval: 10_000, // Refresh chat queue frequently
  });
}

// Work order queue hook
export function useWorkOrderQueue(limit = 20) {
  return useQuery({
    queryKey: [...supportAgentKeys.workOrders(), limit],
    queryFn: () => api.get<{ items: WorkOrderQueueItem[] }>('/supportagent/work-orders', { limit }),
  });
}

// Incident metrics hook
export function useIncidentMetrics() {
  return useQuery({
    queryKey: supportAgentKeys.metrics(),
    queryFn: () => api.get<IncidentMetrics>('/supportagent/metrics'),
  });
}
