import api from '@/api/client';
import type {
  Presentation,
  PresentationType,
  PresentationCategory,
  CreatePresentationRequest,
  UpdatePresentationRequest,
  PresentationFilters,
  PresentationUploadResponse,
} from '@/types/sales';

export interface PresentationsListResponse {
  presentations: Presentation[];
  total: number;
}

export interface PresentationTypesResponse {
  types: PresentationType[];
  categories: PresentationCategory[];
}

export interface DownloadUrlResponse {
  url: string;
  fileName: string;
  fileType: string;
  expiresInS: number;
}

export const presentationsApi = {
  // List presentations with optional filters
  list: (filters?: PresentationFilters): Promise<PresentationsListResponse> => {
    const params = new URLSearchParams();
    if (filters?.type) params.set('type', filters.type);
    if (filters?.category) params.set('category', filters.category);
    if (filters?.featured !== undefined) params.set('featured', String(filters.featured));
    if (filters?.active !== undefined) params.set('active', String(filters.active));
    if (filters?.search) params.set('search', filters.search);
    if (filters?.limit) params.set('limit', String(filters.limit));
    if (filters?.offset) params.set('offset', String(filters.offset));

    return api.get<PresentationsListResponse>(`/presentations?${params.toString()}`);
  },

  // Get single presentation
  getById: (id: string): Promise<Presentation> => {
    return api.get<Presentation>(`/presentations/${id}`);
  },

  // Create presentation (returns upload URLs)
  create: (data: CreatePresentationRequest): Promise<PresentationUploadResponse> => {
    return api.post<PresentationUploadResponse>('/presentations', data);
  },

  // Update presentation
  update: (id: string, data: UpdatePresentationRequest): Promise<Presentation> => {
    return api.put<Presentation>(`/presentations/${id}`, data);
  },

  // Delete presentation
  delete: (id: string): Promise<void> => {
    return api.delete(`/presentations/${id}`);
  },

  // Get download URL
  getDownloadUrl: (id: string): Promise<DownloadUrlResponse> => {
    return api.get<DownloadUrlResponse>(`/presentations/${id}/download`);
  },

  // Record view event
  recordView: (id: string, context?: string, durationSeconds?: number): Promise<void> => {
    return api.post(`/presentations/${id}/view`, { context, durationSeconds });
  },

  // Get available types and categories
  getTypes: (): Promise<PresentationTypesResponse> => {
    return api.get<PresentationTypesResponse>('/presentations/types');
  },

  // Upload file directly to the presigned URL
  uploadFile: async (uploadUrl: string, file: File): Promise<void> => {
    await fetch(uploadUrl, {
      method: 'PUT',
      body: file,
      headers: {
        'Content-Type': file.type,
      },
    });
  },
};

export default presentationsApi;
