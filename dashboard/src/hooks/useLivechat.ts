import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { livechatApi } from '@/api/livechat';
import type {
  StartSessionRequest,
  EndSessionRequest,
  TransferChatRequest,
  SetAvailabilityRequest,
  AIChatMessageRequest,
  AIEscalationRequest,
} from '@/types/livechat';

// Query keys
export const livechatKeys = {
  all: ['livechat'] as const,
  queue: () => [...livechatKeys.all, 'queue'] as const,
  activeChats: () => [...livechatKeys.all, 'active'] as const,
  availability: () => [...livechatKeys.all, 'availability'] as const,
  queuePosition: (sessionId: string) => [...livechatKeys.all, 'queuePosition', sessionId] as const,
  metrics: (from?: string, to?: string) => [...livechatKeys.all, 'metrics', from, to] as const,
};

// Start a new chat session (for school contacts)
export function useStartSession() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data?: StartSessionRequest) => livechatApi.startSession(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: livechatKeys.queue() });
    },
  });
}

// Get queue position
export function useQueuePosition(sessionId: string | undefined, enabled = true) {
  return useQuery({
    queryKey: livechatKeys.queuePosition(sessionId || ''),
    queryFn: () => livechatApi.getQueuePosition(sessionId!),
    enabled: enabled && !!sessionId,
    refetchInterval: 5000, // Poll every 5 seconds
  });
}

// End a chat session
export function useEndSession() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ sessionId, data }: { sessionId: string; data?: EndSessionRequest }) =>
      livechatApi.endSession(sessionId, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: livechatKeys.queue() });
      queryClient.invalidateQueries({ queryKey: livechatKeys.activeChats() });
    },
  });
}

// Accept next chat from queue (for agents)
export function useAcceptChat() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: () => livechatApi.acceptChat(),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: livechatKeys.queue() });
      queryClient.invalidateQueries({ queryKey: livechatKeys.activeChats() });
      queryClient.invalidateQueries({ queryKey: livechatKeys.availability() });
    },
  });
}

// Transfer chat to another agent
export function useTransferChat() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ sessionId, data }: { sessionId: string; data: TransferChatRequest }) =>
      livechatApi.transferChat(sessionId, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: livechatKeys.activeChats() });
    },
  });
}

// Get chat queue
export function useChatQueue(enabled = true) {
  return useQuery({
    queryKey: livechatKeys.queue(),
    queryFn: () => livechatApi.getQueue(),
    enabled,
    refetchInterval: 10000, // Refresh every 10 seconds
  });
}

// Get active chats for current agent
export function useActiveChats(enabled = true) {
  return useQuery({
    queryKey: livechatKeys.activeChats(),
    queryFn: () => livechatApi.getActiveChats(),
    enabled,
    refetchInterval: 10000, // Refresh every 10 seconds
  });
}

// Get/set agent availability
export function useAgentAvailability(enabled = true) {
  return useQuery({
    queryKey: livechatKeys.availability(),
    queryFn: () => livechatApi.getAvailability(),
    enabled,
  });
}

export function useSetAvailability() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: SetAvailabilityRequest) => livechatApi.setAvailability(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: livechatKeys.availability() });
    },
  });
}

// Get chat metrics (admin only)
export function useChatMetrics(from?: string, to?: string) {
  return useQuery({
    queryKey: livechatKeys.metrics(from, to),
    queryFn: () => livechatApi.getMetrics(from, to),
  });
}

// AI Chat hooks

// Send message to AI
export function useSendAIMessage() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ sessionId, content, deviceSerial }: { sessionId: string } & AIChatMessageRequest) =>
      livechatApi.sendAIMessage(sessionId, { content, deviceSerial }),
    onSuccess: () => {
      // Refresh thread messages
      queryClient.invalidateQueries({ queryKey: ['messages'] });
    },
  });
}

// Request escalation to human agent
export function useRequestEscalation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ sessionId, reason }: { sessionId: string } & AIEscalationRequest) =>
      livechatApi.requestEscalation(sessionId, { reason }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: livechatKeys.queue() });
      queryClient.invalidateQueries({ queryKey: ['messages'] });
    },
  });
}

// Get AI conversation context for agent handoff
export function useAIContext(sessionId: string | undefined, enabled = true) {
  return useQuery({
    queryKey: [...livechatKeys.all, 'ai-context', sessionId] as const,
    queryFn: () => livechatApi.getAIContext(sessionId!),
    enabled: enabled && !!sessionId,
  });
}
