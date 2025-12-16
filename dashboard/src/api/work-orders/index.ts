// Query keys
export { WORK_ORDERS_KEY } from './keys';

// Core work order operations
export {
  useWorkOrders,
  useWorkOrder,
  useCreateWorkOrder,
  useUpdateWorkOrderStatus,
} from './core';

// Update operations
export { useUpdateWorkOrder } from './update';

// Rejection/Rework operations
export {
  useRejectWorkOrder,
  useWorkOrderReworkHistory,
  REWORK_HISTORY_KEY,
} from './rejection';

// Bulk operations
export {
  useBulkStatusUpdate,
  useBulkAssignment,
  useBulkApproval,
} from './bulk';

// Bill of Materials (BOM) operations
export {
  useWorkOrderBOM,
  useAddBOMItem,
  useConsumeBOMItem,
  useReleaseBOMItem,
} from './bom';

// Scheduling operations
export {
  useWorkOrderSchedules,
  useCreateSchedule,
} from './scheduling';

// Deliverables operations
export {
  useWorkOrderDeliverables,
  useAddDeliverable,
  useSubmitDeliverable,
  useReviewDeliverable,
} from './deliverables';

// Approval operations
export {
  useRequestApproval,
  useDecideApproval,
} from './approvals';
