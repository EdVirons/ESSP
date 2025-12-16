import { useQuery, useMutation, useQueryClient, useInfiniteQuery } from '@tanstack/react-query';
import { messagingApi } from '@/api/messaging';
import type {
  CreateThreadRequest,
  CreateMessageRequest,
  ThreadsListParams,
} from '@/types/messaging';

// Query keys
export const messagingKeys = {
  all: ['messaging'] as const,
  threads: (params?: ThreadsListParams) => [...messagingKeys.all, 'threads', params] as const,
  thread: (id: string) => [...messagingKeys.all, 'thread', id] as const,
  unread: () => [...messagingKeys.all, 'unread'] as const,
  search: (query: string) => [...messagingKeys.all, 'search', query] as const,
  analytics: (from?: string, to?: string) => [...messagingKeys.all, 'analytics', from, to] as const,
};

// List threads with cursor-based pagination
export function useThreads(params?: ThreadsListParams) {
  return useQuery({
    queryKey: messagingKeys.threads(params),
    queryFn: () => messagingApi.listThreads(params),
  });
}

// Infinite loading threads
export function useInfiniteThreads(params?: Omit<ThreadsListParams, 'cursor'>) {
  return useInfiniteQuery({
    queryKey: messagingKeys.threads(params),
    queryFn: ({ pageParam }) => messagingApi.listThreads({ ...params, cursor: pageParam }),
    getNextPageParam: (lastPage) => lastPage.nextCursor,
    initialPageParam: undefined as string | undefined,
  });
}

// Get single thread with messages
export function useThread(threadId: string | undefined) {
  return useQuery({
    queryKey: messagingKeys.thread(threadId || ''),
    queryFn: () => messagingApi.getThread(threadId!),
    enabled: !!threadId,
  });
}

// Create new thread
export function useCreateThread() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateThreadRequest) => messagingApi.createThread(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: messagingKeys.threads() });
      queryClient.invalidateQueries({ queryKey: messagingKeys.unread() });
    },
  });
}

// Send message to thread
export function useSendMessage(threadId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateMessageRequest) => messagingApi.sendMessage(threadId, data),
    onSuccess: (newMessage) => {
      // Optimistically update the thread cache
      queryClient.setQueryData(messagingKeys.thread(threadId), (old: any) => {
        if (!old) return old;
        return {
          ...old,
          messages: [...old.messages, newMessage],
          thread: {
            ...old.thread,
            messageCount: old.thread.messageCount + 1,
            lastMessageAt: newMessage.createdAt,
            lastMessage: newMessage,
          },
        };
      });

      // Invalidate threads list to update order
      queryClient.invalidateQueries({ queryKey: messagingKeys.threads() });
    },
  });
}

// Close thread
export function useCloseThread() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (threadId: string) => messagingApi.closeThread(threadId),
    onSuccess: (_, threadId) => {
      queryClient.invalidateQueries({ queryKey: messagingKeys.thread(threadId) });
      queryClient.invalidateQueries({ queryKey: messagingKeys.threads() });
    },
  });
}

// Reopen thread
export function useReopenThread() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (threadId: string) => messagingApi.reopenThread(threadId),
    onSuccess: (_, threadId) => {
      queryClient.invalidateQueries({ queryKey: messagingKeys.thread(threadId) });
      queryClient.invalidateQueries({ queryKey: messagingKeys.threads() });
    },
  });
}

// Mark thread as read
export function useMarkRead(threadId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (lastMessageId: string) => messagingApi.markRead(threadId, lastMessageId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: messagingKeys.thread(threadId) });
      queryClient.invalidateQueries({ queryKey: messagingKeys.unread() });
    },
  });
}

// Get unread counts
export function useUnreadCounts() {
  return useQuery({
    queryKey: messagingKeys.unread(),
    queryFn: () => messagingApi.getUnreadCounts(),
    refetchInterval: 30000, // Refresh every 30 seconds
    staleTime: 10000, // Consider data stale after 10 seconds
  });
}

// Search messages
export function useSearchMessages(query: string, enabled = true) {
  return useQuery({
    queryKey: messagingKeys.search(query),
    queryFn: () => messagingApi.searchMessages(query),
    enabled: enabled && query.length >= 2,
  });
}

// Get messaging analytics (admin only)
export function useMessagingAnalytics(from?: string, to?: string) {
  return useQuery({
    queryKey: messagingKeys.analytics(from, to),
    queryFn: () => messagingApi.getAnalytics(from, to),
  });
}
