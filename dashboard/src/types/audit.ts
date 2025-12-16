// Audit Log types
export type AuditAction = 'create' | 'update' | 'delete';

export interface AuditLog {
  id: string;
  tenantId: string;
  userId: string;
  userEmail: string;
  action: AuditAction;
  entityType: string;
  entityId: string;
  beforeState: Record<string, unknown> | null;
  afterState: Record<string, unknown> | null;
  ipAddress: string;
  userAgent: string;
  requestId: string;
  createdAt: string;
}

// Attachment types
export type AttachmentEntityType = 'incident' | 'work_order';

export interface Attachment {
  id: string;
  tenantId: string;
  schoolId: string;
  entityType: AttachmentEntityType;
  entityId: string;
  fileName: string;
  contentType: string;
  sizeBytes: number;
  objectKey: string;
  createdAt: string;
}

export interface CreateAttachmentRequest {
  entityType?: AttachmentEntityType;
  entityId: string;
  fileName: string;
  contentType?: string;
  sizeBytes?: number;
}
