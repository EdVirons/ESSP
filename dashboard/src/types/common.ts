// Common types
export interface PaginatedResponse<T> {
  items: T[];
  nextCursor: string | null;
}

export interface User {
  id: string;
  email: string;
  name: string;
  roles: string[];
  permissions: string[];
  tenantId: string;
  schoolId?: string;
}

// Dashboard specific types
export interface ServiceHealth {
  name: string;
  status: 'healthy' | 'degraded' | 'unhealthy';
  latencyMs: number;
  lastCheck: string;
}

export interface DashboardMetrics {
  incidents: {
    total: number;
    open: number;
    slaBreached: number;
  };
  workOrders: {
    total: number;
    inProgress: number;
    completedToday: number;
  };
  programs: {
    active: number;
    pending: number;
  };
}

export interface ActivityEvent {
  id: string;
  type: string;
  action: string;
  actor: string;
  target: string;
  timestamp: string;
  metadata: Record<string, unknown>;
}
