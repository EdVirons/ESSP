import type { IncidentStatus, Severity } from './incident';
import type { WorkOrderStatus } from './work-order';
import type { AuditAction } from './audit';

// Filter types
export interface IncidentFilters {
  status?: IncidentStatus;
  severity?: Severity;
  deviceId?: string;
  q?: string;
  limit?: number;
  cursor?: string;
}

export interface WorkOrderFilters {
  status?: WorkOrderStatus;
  deviceId?: string;
  incidentId?: string;
  limit?: number;
  cursor?: string;
}

export interface AuditLogFilters {
  entityType?: string;
  entityId?: string;
  userId?: string;
  action?: AuditAction;
  startDate?: string;
  endDate?: string;
  limit?: number;
  cursor?: string;
}
