// API client
export { api, apiClient } from './client';

// Incidents
export {
  useIncidents,
  useIncident,
  useCreateIncident,
  useUpdateIncidentStatus,
  usePrefetchIncident,
} from './incidents';

// Work Orders
export {
  useWorkOrders,
  useWorkOrder,
  useCreateWorkOrder,
  useUpdateWorkOrderStatus,
  useWorkOrderBOM,
  useAddBOMItem,
  useConsumeBOMItem,
  useReleaseBOMItem,
  useWorkOrderSchedules,
  useCreateSchedule,
  useWorkOrderDeliverables,
  useAddDeliverable,
  useSubmitDeliverable,
  useReviewDeliverable,
  useRequestApproval,
  useDecideApproval,
} from './work-orders';

// Projects
export {
  useProjects,
  useProject,
  useCreateProject,
  usePhases,
  useCreatePhase,
  useUpdatePhaseStatus,
  useSurveys,
  useSurvey,
  useCreateSurvey,
  useAddSurveyRoom,
} from './projects';

// Service Shops
export {
  useServiceShops,
  useServiceShop,
  useCreateServiceShop,
  useServiceStaff,
  useServiceStaffMember,
  useCreateServiceStaff,
  useInventory,
  useUpsertInventory,
  useParts,
  usePart,
} from './service-shops';

// Audit Logs
export {
  useAuditLogs,
  useAuditLog,
  useAuditLogEntityTypes,
  exportAuditLogs,
} from './audit-logs';

// Health
export {
  useServiceHealth,
  useDashboardMetrics,
  useActivityFeed,
  useHealthCheck,
  useReadinessCheck,
} from './health';

// SSOT (Single Source of Truth)
export {
  useSchools,
  useDevices,
  usePartsSnapshot,
  useSyncStatus,
  useSyncSchools,
  useSyncDevices,
  useSyncParts,
} from './ssot';

export type {
  SchoolSnapshot,
  DeviceSnapshot,
  PartSnapshot,
  SSOTSyncStatus,
  SchoolFilters,
  DeviceFilters,
  PartFilters,
} from './ssot';

// Sales / Demo Pipeline / Presentations
export { demoPipelineApi } from './demo-pipeline';
export { presentationsApi } from './presentations';
export { salesApi } from './sales';
