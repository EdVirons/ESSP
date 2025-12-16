import { api } from '@/lib/api';
import type {
  AgentAvailability,
  StartSessionRequest,
  StartSessionResponse,
  EndSessionRequest,
  QueuePositionResponse,
  AcceptChatResponse,
  TransferChatRequest,
  SetAvailabilityRequest,
  ChatQueueResponse,
  ActiveChatsResponse,
  ChatMetrics,
  AIChatMessageRequest,
  AIChatMessageResponse,
  AIEscalationRequest,
  AIEscalationResponse,
  AIConversationContext,
} from '@/types/livechat';

export const livechatApi = {
  // Session operations (for school contacts)
  startSession: async (data?: StartSessionRequest): Promise<StartSessionResponse> => {
    const response = await api.post('/chat/sessions', data || {});
    return response.data;
  },

  getQueuePosition: async (sessionId: string): Promise<QueuePositionResponse> => {
    const response = await api.get(`/chat/sessions/${sessionId}/queue`);
    return response.data;
  },

  endSession: async (sessionId: string, data?: EndSessionRequest): Promise<void> => {
    await api.post(`/chat/sessions/${sessionId}/end`, data || {});
  },

  // Agent operations
  acceptChat: async (): Promise<AcceptChatResponse> => {
    const response = await api.post('/chat/accept');
    return response.data;
  },

  transferChat: async (sessionId: string, data: TransferChatRequest): Promise<void> => {
    await api.post(`/chat/sessions/${sessionId}/transfer`, data);
  },

  getQueue: async (): Promise<ChatQueueResponse> => {
    const response = await api.get('/chat/queue');
    return response.data;
  },

  getActiveChats: async (): Promise<ActiveChatsResponse> => {
    const response = await api.get('/chat/active');
    return response.data;
  },

  // Availability
  setAvailability: async (data: SetAvailabilityRequest): Promise<void> => {
    await api.put('/chat/availability', data);
  },

  getAvailability: async (): Promise<AgentAvailability> => {
    const response = await api.get('/chat/availability');
    return response.data;
  },

  // Metrics (admin only)
  getMetrics: async (from?: string, to?: string): Promise<ChatMetrics> => {
    const response = await api.get('/chat/metrics', {
      params: { from, to },
    });
    return response.data;
  },

  // AI Chat operations
  sendAIMessage: async (sessionId: string, data: AIChatMessageRequest): Promise<AIChatMessageResponse> => {
    const response = await api.post(`/chat/ai/sessions/${sessionId}/message`, data);
    return response.data;
  },

  requestEscalation: async (sessionId: string, data?: AIEscalationRequest): Promise<AIEscalationResponse> => {
    const response = await api.post(`/chat/ai/sessions/${sessionId}/escalate`, data || {});
    return response.data;
  },

  getAIContext: async (sessionId: string): Promise<AIConversationContext> => {
    const response = await api.get(`/chat/ai/sessions/${sessionId}/context`);
    return response.data;
  },
};

export default livechatApi;
