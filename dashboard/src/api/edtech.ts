import { api } from '@/lib/api';
import type {
  EdTechProfileResponse,
  EdTechFormOptions,
  SaveProfileRequest,
  SubmitFollowUpRequest,
  EdTechHistoryResponse,
} from '@/types/edtech';

export const edtechApi = {
  // Get form options
  getOptions: async (): Promise<EdTechFormOptions> => {
    const response = await api.get('/edtech-profiles/options');
    return response.data;
  },

  // Get profile by school ID
  getBySchoolId: async (schoolId: string): Promise<EdTechProfileResponse> => {
    const response = await api.get(`/edtech-profiles/school/${schoolId}`);
    return response.data;
  },

  // Create or update profile
  saveProfile: async (data: SaveProfileRequest): Promise<EdTechProfileResponse> => {
    const response = await api.post('/edtech-profiles', data);
    return response.data;
  },

  // Generate AI analysis
  generateAI: async (profileId: string): Promise<EdTechProfileResponse> => {
    const response = await api.post(`/edtech-profiles/${profileId}/generate-ai`);
    return response.data;
  },

  // Submit follow-up responses
  submitFollowUp: async (profileId: string, data: SubmitFollowUpRequest): Promise<EdTechProfileResponse> => {
    const response = await api.post(`/edtech-profiles/${profileId}/submit-followup`, data);
    return response.data;
  },

  // Mark profile as complete
  complete: async (profileId: string): Promise<EdTechProfileResponse> => {
    const response = await api.post(`/edtech-profiles/${profileId}/complete`);
    return response.data;
  },

  // Get profile version history
  getHistory: async (schoolId: string, limit?: number): Promise<EdTechHistoryResponse> => {
    const response = await api.get(`/edtech-profiles/school/${schoolId}/history`, {
      params: { limit },
    });
    return response.data;
  },
};
