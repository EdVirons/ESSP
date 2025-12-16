// Incident types
export type IncidentStatus =
  | 'new'
  | 'acknowledged'
  | 'in_progress'
  | 'escalated'
  | 'resolved'
  | 'closed';

export type Severity = 'low' | 'medium' | 'high' | 'critical';

export interface Incident {
  id: string;
  tenantId: string;
  schoolId: string;
  deviceId: string;
  schoolName: string;
  countyId: string;
  countyName: string;
  subCountyId: string;
  subCountyName: string;
  contactName: string;
  contactPhone: string;
  contactEmail: string;
  deviceSerial: string;
  deviceAssetTag: string;
  deviceModelId: string;
  deviceMake: string;
  deviceModel: string;
  deviceCategory: string;
  category: string;
  severity: Severity;
  status: IncidentStatus;
  title: string;
  description: string;
  reportedBy: string;
  slaDueAt: string;
  slaBreached: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface CreateIncidentRequest {
  deviceId: string;
  category?: string;
  severity?: Severity;
  title: string;
  description?: string;
  reportedBy?: string;
}
