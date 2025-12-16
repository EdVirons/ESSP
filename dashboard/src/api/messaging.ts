import { api } from '@/lib/api';
import type {
  Message,
  ThreadsListResponse,
  ThreadDetailResponse,
  CreateThreadRequest,
  CreateThreadResponse,
  CreateMessageRequest,
  UnreadCounts,
  SearchResponse,
  MessagingAnalytics,
  ThreadsListParams,
  UploadURLRequest,
  UploadURLResponse,
} from '@/types/messaging';

export const messagingApi = {
  // Thread operations
  listThreads: async (params?: ThreadsListParams): Promise<ThreadsListResponse> => {
    const response = await api.get('/messages/threads', { params });
    return response.data;
  },

  getThread: async (threadId: string): Promise<ThreadDetailResponse> => {
    const response = await api.get(`/messages/threads/${threadId}`);
    return response.data;
  },

  createThread: async (data: CreateThreadRequest): Promise<CreateThreadResponse> => {
    const response = await api.post('/messages/threads', data);
    return response.data;
  },

  updateThreadStatus: async (threadId: string, status: 'open' | 'closed'): Promise<void> => {
    await api.patch(`/messages/threads/${threadId}/status`, { status });
  },

  closeThread: async (threadId: string): Promise<void> => {
    await api.patch(`/messages/threads/${threadId}/status`, { status: 'closed' });
  },

  reopenThread: async (threadId: string): Promise<void> => {
    await api.patch(`/messages/threads/${threadId}/status`, { status: 'open' });
  },

  // Message operations
  sendMessage: async (threadId: string, data: CreateMessageRequest): Promise<Message> => {
    const response = await api.post(`/messages/threads/${threadId}/messages`, data);
    return response.data.message;
  },

  markRead: async (threadId: string, lastMessageId: string): Promise<void> => {
    await api.post(`/messages/threads/${threadId}/read`, { lastMessageId });
  },

  // Unread counts
  getUnreadCounts: async (): Promise<UnreadCounts> => {
    const response = await api.get('/messages/unread');
    return response.data;
  },

  // Search
  searchMessages: async (query: string, limit?: number): Promise<SearchResponse> => {
    const response = await api.get('/messages/search', {
      params: { q: query, limit: limit || 20 },
    });
    return response.data;
  },

  // Analytics (admin only)
  getAnalytics: async (from?: string, to?: string): Promise<MessagingAnalytics> => {
    const response = await api.get('/messages/analytics', {
      params: { from, to },
    });
    return response.data;
  },

  // Attachments
  getUploadUrl: async (data: UploadURLRequest): Promise<UploadURLResponse> => {
    const response = await api.post('/messages/attachments/upload-url', data);
    return response.data;
  },
};

export default messagingApi;
