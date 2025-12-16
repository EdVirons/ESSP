import { useQuery } from '@tanstack/react-query';
import api from './client';
import type { AuditLog, AuditLogFilters, PaginatedResponse } from '@/types';

const AUDIT_LOGS_KEY = 'audit-logs';

export function useAuditLogs(filters: AuditLogFilters = {}) {
  return useQuery({
    queryKey: [AUDIT_LOGS_KEY, filters],
    queryFn: () => api.get<PaginatedResponse<AuditLog>>('/audit-logs', filters),
    staleTime: 30_000,
  });
}

export function useAuditLog(id: string) {
  return useQuery({
    queryKey: [AUDIT_LOGS_KEY, id],
    queryFn: () => api.get<AuditLog>(`/audit-logs/${id}`),
    enabled: !!id,
  });
}

// Helper to get unique entity types for filtering
export function useAuditLogEntityTypes() {
  return useQuery({
    queryKey: [AUDIT_LOGS_KEY, 'entity-types'],
    queryFn: async () => {
      // This returns a static list - in production could be fetched from API
      return [
        { value: 'incident', label: 'Incidents' },
        { value: 'work_order', label: 'Work Orders' },
        { value: 'program', label: 'Programs' },
        { value: 'phase', label: 'Phases' },
        { value: 'service_shop', label: 'Service Shops' },
        { value: 'service_staff', label: 'Service Staff' },
        { value: 'inventory', label: 'Inventory' },
        { value: 'school', label: 'Schools' },
        { value: 'device', label: 'Devices' },
        { value: 'attachment', label: 'Attachments' },
      ];
    },
    staleTime: Infinity,
  });
}

// Export audit logs to CSV
export async function exportAuditLogs(filters: AuditLogFilters): Promise<Blob> {
  const response = await api.get<PaginatedResponse<AuditLog>>('/audit-logs', {
    ...filters,
    limit: 1000, // Get more records for export
  });

  // Convert to CSV
  const headers = [
    'ID',
    'Timestamp',
    'User',
    'Action',
    'Entity Type',
    'Entity ID',
    'IP Address',
    'Request ID',
  ];

  const rows = response.items.map((log) => [
    log.id,
    log.createdAt,
    log.userEmail,
    log.action,
    log.entityType,
    log.entityId,
    log.ipAddress,
    log.requestId,
  ]);

  const csvContent = [
    headers.join(','),
    ...rows.map((row) =>
      row.map((cell) => `"${String(cell).replace(/"/g, '""')}"`).join(',')
    ),
  ].join('\n');

  return new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
}
