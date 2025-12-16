// Work Order types
export type WorkOrderStatus =
  | 'draft'
  | 'assigned'
  | 'in_repair'
  | 'qa'
  | 'completed'
  | 'approved';

export type RepairLocation = 'service_shop' | 'on_site';

export interface WorkOrder {
  id: string;
  incidentId: string;
  tenantId: string;
  schoolId: string;
  deviceId: string;
  schoolName: string;
  contactName: string;
  contactPhone: string;
  contactEmail: string;
  deviceSerial: string;
  deviceAssetTag: string;
  deviceModelId: string;
  deviceMake: string;
  deviceModel: string;
  deviceCategory: string;
  status: WorkOrderStatus;
  serviceShopId: string;
  assignedStaffId: string;
  repairLocation: RepairLocation;
  assignedTo: string;
  taskType: string;
  programId: string;
  phaseId: string;
  onsiteContactId: string;
  approvalStatus: string;
  costEstimateCents: number;
  notes: string;
  // Rework tracking
  reworkCount: number;
  lastReworkAt?: string;
  lastReworkReason: string;
  createdAt: string;
  updatedAt: string;
}

export interface CreateWorkOrderRequest {
  incidentId?: string;
  deviceId: string;
  taskType: string;
  serviceShopId?: string;
  assignedStaffId?: string;
  repairLocation?: RepairLocation;
  assignedTo?: string;
  costEstimateCents?: number;
  notes?: string;
}

// Work Order Operations
export interface WorkOrderSchedule {
  id: string;
  tenantId: string;
  schoolId: string;
  workOrderId: string;
  phaseId: string;
  scheduledStart: string | null;
  scheduledEnd: string | null;
  timezone: string;
  notes: string;
  createdByUserId: string;
  createdAt: string;
}

export type DeliverableStatus = 'pending' | 'submitted' | 'approved' | 'rejected';

export interface WorkOrderDeliverable {
  id: string;
  tenantId: string;
  schoolId: string;
  workOrderId: string;
  phaseId: string;
  title: string;
  description: string;
  status: DeliverableStatus;
  evidenceAttachmentId: string;
  submittedByUserId: string;
  submittedAt: string | null;
  reviewedByUserId: string;
  reviewedAt: string | null;
  reviewNotes: string;
  createdAt: string;
  updatedAt: string;
}

export type ApprovalStatus = 'pending' | 'approved' | 'rejected';

export interface WorkOrderApproval {
  id: string;
  tenantId: string;
  schoolId: string;
  workOrderId: string;
  phaseId: string;
  approvalType: string;
  requestedByUserId: string;
  requestedAt: string;
  status: ApprovalStatus;
  decidedByUserId: string;
  decidedAt: string | null;
  decisionNotes: string;
}

// BOM types
export interface WorkOrderPart {
  id: string;
  tenantId: string;
  schoolId: string;
  workOrderId: string;
  serviceShopId: string;
  partId: string;
  partName: string;
  partPuk: string;
  partCategory: string;
  deviceModelId: string;
  isCompatible: boolean;
  qtyPlanned: number;
  qtyUsed: number;
  createdAt: string;
  updatedAt: string;
}

// Rejection/Rework types
export type RejectionCategory = 'quality' | 'incomplete' | 'wrong_parts' | 'safety' | 'other';

export interface WorkOrderReworkHistory {
  id: string;
  tenantId: string;
  schoolId: string;
  workOrderId: string;
  fromStatus: WorkOrderStatus;
  toStatus: WorkOrderStatus;
  rejectionReason: string;
  rejectionCategory: RejectionCategory;
  rejectedByUserId: string;
  rejectedByName: string;
  reworkSequence: number;
  createdAt: string;
}

export interface RejectWorkOrderRequest {
  targetStatus: WorkOrderStatus;
  reason: string;
  category: RejectionCategory;
}

export interface RejectWorkOrderResponse {
  workOrder: WorkOrder;
  reworkHistory: WorkOrderReworkHistory;
}

// Update types
export interface UpdateWorkOrderRequest {
  assignedStaffId?: string;
  serviceShopId?: string;
  costEstimateCents?: number;
  notes?: string;
  repairLocation?: RepairLocation;
  onsiteContactId?: string;
}

export interface UpdateWorkOrderResponse {
  workOrder: WorkOrder;
  updatedFields: string[];
  previousValues?: Record<string, unknown>;
}

// Bulk operation types
export interface BulkStatusUpdateRequest {
  workOrderIds: string[];
  status: WorkOrderStatus;
}

export interface BulkAssignmentRequest {
  workOrderIds: string[];
  assignedStaffId?: string;
  serviceShopId?: string;
}

export interface BulkApprovalRequest {
  workOrderIds: string[];
  decision: 'approved' | 'rejected';
  notes?: string;
}

export interface BulkOperationError {
  id: string;
  message: string;
  code?: string;
}

export interface BulkOperationResult {
  operationId: string;
  succeeded: string[];
  failed: BulkOperationError[];
  totalCount: number;
  successCount: number;
  failureCount: number;
}
