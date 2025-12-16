// Demo Lead Types
export type DemoLeadStage =
  | 'new_lead'
  | 'contacted'
  | 'demo_scheduled'
  | 'demo_completed'
  | 'proposal_sent'
  | 'negotiation'
  | 'won'
  | 'lost';

export type DemoLeadSource =
  | 'website'
  | 'referral'
  | 'event'
  | 'cold_outreach'
  | 'inbound';

export type DemoActivityType =
  | 'note'
  | 'call'
  | 'email'
  | 'meeting'
  | 'demo'
  | 'stage_change'
  | 'created'
  | 'updated';

export type DemoScheduleStatus =
  | 'scheduled'
  | 'completed'
  | 'cancelled'
  | 'rescheduled';

export interface DemoLead {
  id: string;
  tenantId: string;
  schoolId?: string;
  schoolName: string;
  contactName: string;
  contactEmail: string;
  contactPhone: string;
  contactRole: string;
  countyCode: string;
  countyName: string;
  subCountyCode: string;
  subCountyName: string;
  stage: DemoLeadStage;
  stageChangedAt: string;
  estimatedValue?: number;
  estimatedDevices?: number;
  probability: number;
  expectedCloseDate?: string;
  leadSource: DemoLeadSource;
  assignedTo: string;
  notes: string;
  tags: string[];
  lostReason?: string;
  lostNotes?: string;
  createdBy: string;
  createdAt: string;
  updatedAt: string;
}

export interface DemoLeadActivity {
  id: string;
  tenantId: string;
  leadId: string;
  activityType: DemoActivityType;
  description: string;
  fromStage?: DemoLeadStage;
  toStage?: DemoLeadStage;
  scheduledAt?: string;
  completedAt?: string;
  createdBy: string;
  createdAt: string;
}

export interface DemoAttendee {
  name: string;
  email: string;
  role: string;
}

export interface DemoSchedule {
  id: string;
  tenantId: string;
  leadId: string;
  scheduledDate: string;
  scheduledTime: string;
  durationMinutes: number;
  location: string;
  meetingLink: string;
  attendees: DemoAttendee[];
  status: DemoScheduleStatus;
  outcome: string;
  outcomeNotes: string;
  reminderSent: boolean;
  createdBy: string;
  createdAt: string;
  updatedAt: string;
}

export interface DemoLeadWithActivities extends DemoLead {
  recentActivities: DemoLeadActivity[];
  nextDemo?: DemoSchedule;
}

export interface PipelineStageCount {
  stage: DemoLeadStage;
  count: number;
  totalValue: number;
}

export interface PipelineSummary {
  stages: PipelineStageCount[];
  totalLeads: number;
  totalValue: number;
  averageValue: number;
  conversionRate: number;
}

// Request types
export interface CreateDemoLeadRequest {
  schoolId?: string;
  schoolName: string;
  contactName?: string;
  contactEmail?: string;
  contactPhone?: string;
  contactRole?: string;
  countyCode?: string;
  countyName?: string;
  subCountyCode?: string;
  subCountyName?: string;
  estimatedValue?: number;
  estimatedDevices?: number;
  leadSource?: string;
  notes?: string;
  tags?: string[];
}

export interface UpdateDemoLeadRequest {
  schoolName?: string;
  contactName?: string;
  contactEmail?: string;
  contactPhone?: string;
  contactRole?: string;
  countyCode?: string;
  countyName?: string;
  subCountyCode?: string;
  subCountyName?: string;
  estimatedValue?: number;
  estimatedDevices?: number;
  probability?: number;
  expectedCloseDate?: string;
  assignedTo?: string;
  notes?: string;
  tags?: string[];
}

export interface UpdateLeadStageRequest {
  stage: DemoLeadStage;
  lostReason?: string;
  lostNotes?: string;
}

export interface CreateDemoScheduleRequest {
  scheduledDate: string;
  scheduledTime?: string;
  durationMinutes?: number;
  location?: string;
  meetingLink?: string;
  attendees?: DemoAttendee[];
}

export interface DemoLeadFilters {
  stage?: DemoLeadStage;
  assignedTo?: string;
  source?: DemoLeadSource;
  search?: string;
  limit?: number;
  offset?: number;
}

// Presentation Types
export type PresentationType =
  | 'presentation'
  | 'brochure'
  | 'case_study'
  | 'video'
  | 'roi_calculator'
  | 'template'
  | 'other';

export type PresentationCategory =
  | 'general'
  | 'product_overview'
  | 'technical'
  | 'pricing'
  | 'onboarding'
  | 'training';

export interface Presentation {
  id: string;
  tenantId: string;
  title: string;
  description: string;
  type: PresentationType;
  category: PresentationCategory;
  fileKey: string;
  fileName: string;
  fileSize: number;
  fileType: string;
  thumbnailKey: string;
  previewType: string;
  tags: string[];
  version: number;
  isActive: boolean;
  isFeatured: boolean;
  viewCount: number;
  downloadCount: number;
  lastViewedAt?: string;
  createdBy: string;
  updatedBy?: string;
  createdAt: string;
  updatedAt: string;
  downloadUrl?: string;
  previewUrl?: string;
}

export interface CreatePresentationRequest {
  title: string;
  description?: string;
  type: string;
  category?: string;
  tags?: string[];
  isFeatured?: boolean;
  fileName?: string;
  fileSize?: number;
  fileType?: string;
}

export interface UpdatePresentationRequest {
  title?: string;
  description?: string;
  type?: string;
  category?: string;
  tags?: string[];
  isActive?: boolean;
  isFeatured?: boolean;
}

export interface PresentationFilters {
  type?: PresentationType;
  category?: PresentationCategory;
  featured?: boolean;
  active?: boolean;
  search?: string;
  limit?: number;
  offset?: number;
}

export interface PresentationUploadResponse {
  presentation: Presentation;
  uploadUrl: string;
  thumbnailUploadUrl: string;
}

// Sales Metrics Types
export interface SalesMetricsSummary {
  totalLeads: number;
  newLeadsThisPeriod: number;
  demosScheduled: number;
  demosCompleted: number;
  proposalsSent: number;
  dealsWon: number;
  dealsLost: number;
  totalPipelineValue: number;
  wonValueThisPeriod: number;
  conversionRate: number;
  winRate: number;
  averageDealSize: number;
  totalActivities: number;
}

export interface RecentActivity {
  id: string;
  type: string;
  description: string;
  leadId: string;
  leadName: string;
  userId: string;
  userName?: string;
  createdAt: string;
}

export interface SchoolsByRegion {
  region: string;
  count: number;
  value: number;
}

export interface SalesDashboardResponse {
  metrics: SalesMetricsSummary;
  pipelineStages: PipelineStageCount[];
  recentActivities: RecentActivity[];
  schoolsByRegion: SchoolsByRegion[];
  period: {
    startDate: string;
    endDate: string;
    days: number;
  };
}

// Stage label and color mappings
export const stageLabels: Record<DemoLeadStage, string> = {
  new_lead: 'New Lead',
  contacted: 'Contacted',
  demo_scheduled: 'Demo Scheduled',
  demo_completed: 'Demo Completed',
  proposal_sent: 'Proposal Sent',
  negotiation: 'Negotiation',
  won: 'Won',
  lost: 'Lost',
};

export const stageColors: Record<DemoLeadStage, string> = {
  new_lead: 'bg-blue-100 text-blue-800',
  contacted: 'bg-yellow-100 text-yellow-800',
  demo_scheduled: 'bg-purple-100 text-purple-800',
  demo_completed: 'bg-indigo-100 text-indigo-800',
  proposal_sent: 'bg-orange-100 text-orange-800',
  negotiation: 'bg-pink-100 text-pink-800',
  won: 'bg-green-100 text-green-800',
  lost: 'bg-red-100 text-red-800',
};

export const sourceLabels: Record<DemoLeadSource, string> = {
  website: 'Website',
  referral: 'Referral',
  event: 'Event',
  cold_outreach: 'Cold Outreach',
  inbound: 'Inbound',
};

export const presentationTypeLabels: Record<PresentationType, string> = {
  presentation: 'Presentation',
  brochure: 'Brochure',
  case_study: 'Case Study',
  video: 'Video',
  roi_calculator: 'ROI Calculator',
  template: 'Template',
  other: 'Other',
};

export const presentationCategoryLabels: Record<PresentationCategory, string> = {
  general: 'General',
  product_overview: 'Product Overview',
  technical: 'Technical',
  pricing: 'Pricing',
  onboarding: 'Onboarding',
  training: 'Training',
};
