import api from '@/api/client';
import type {
  DemoLead,
  DemoLeadWithActivities,
  DemoLeadActivity,
  DemoSchedule,
  PipelineSummary,
  CreateDemoLeadRequest,
  UpdateDemoLeadRequest,
  UpdateLeadStageRequest,
  CreateDemoScheduleRequest,
  DemoLeadFilters,
} from '@/types/sales';

export interface LeadsListResponse {
  leads: DemoLead[];
  total: number;
}

export interface ActivitiesListResponse {
  activities: DemoLeadActivity[];
}

export const demoPipelineApi = {
  // List leads with optional filters
  listLeads: (filters?: DemoLeadFilters): Promise<LeadsListResponse> => {
    const params = new URLSearchParams();
    if (filters?.stage) params.set('stage', filters.stage);
    if (filters?.assignedTo) params.set('assignedTo', filters.assignedTo);
    if (filters?.source) params.set('source', filters.source);
    if (filters?.search) params.set('search', filters.search);
    if (filters?.limit) params.set('limit', String(filters.limit));
    if (filters?.offset) params.set('offset', String(filters.offset));

    return api.get<LeadsListResponse>(`/demo-pipeline/leads?${params.toString()}`);
  },

  // Get single lead with activities
  getLead: (id: string): Promise<DemoLeadWithActivities> => {
    return api.get<DemoLeadWithActivities>(`/demo-pipeline/leads/${id}`);
  },

  // Create new lead
  createLead: (data: CreateDemoLeadRequest): Promise<DemoLead> => {
    return api.post<DemoLead>('/demo-pipeline/leads', data);
  },

  // Update lead
  updateLead: (id: string, data: UpdateDemoLeadRequest): Promise<DemoLead> => {
    return api.put<DemoLead>(`/demo-pipeline/leads/${id}`, data);
  },

  // Update lead stage
  updateStage: (id: string, data: UpdateLeadStageRequest): Promise<DemoLead> => {
    return api.put<DemoLead>(`/demo-pipeline/leads/${id}/stage`, data);
  },

  // Delete lead
  deleteLead: (id: string): Promise<void> => {
    return api.delete(`/demo-pipeline/leads/${id}`);
  },

  // Get lead activities
  getActivities: (leadId: string, limit?: number): Promise<ActivitiesListResponse> => {
    const params = limit ? `?limit=${limit}` : '';
    return api.get<ActivitiesListResponse>(`/demo-pipeline/leads/${leadId}/activities${params}`);
  },

  // Add note to lead
  addNote: (leadId: string, note: string): Promise<DemoLeadActivity> => {
    return api.post<DemoLeadActivity>(`/demo-pipeline/leads/${leadId}/notes`, { note });
  },

  // Schedule demo
  scheduleDemo: (leadId: string, data: CreateDemoScheduleRequest): Promise<DemoSchedule> => {
    return api.post<DemoSchedule>(`/demo-pipeline/leads/${leadId}/schedule-demo`, data);
  },

  // Get pipeline summary
  getPipelineSummary: (): Promise<PipelineSummary> => {
    return api.get<PipelineSummary>('/demo-pipeline/summary');
  },

  // Get recent activities across all leads
  getRecentActivities: (limit?: number): Promise<ActivitiesListResponse> => {
    const params = limit ? `?limit=${limit}` : '';
    return api.get<ActivitiesListResponse>(`/demo-pipeline/activities${params}`);
  },
};

export default demoPipelineApi;
