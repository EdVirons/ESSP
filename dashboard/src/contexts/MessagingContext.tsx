import React, { createContext, useContext, useState, useCallback, useEffect } from 'react';
import { useQueryClient } from '@tanstack/react-query';
import { useWebSocket } from '@/hooks/useWebSocket';
import { messagingKeys } from '@/hooks/useMessages';
import type { Message } from '@/types/messaging';

interface MessagingContextValue {
  // Selected thread
  selectedThreadId: string | null;
  setSelectedThreadId: (id: string | null) => void;

  // Real-time message updates
  onNewMessage: (threadId: string, message: Message) => void;

  // Unread counts
  unreadCount: number;
  refreshUnreadCount: () => void;

  // Typing indicators
  typingUsers: Map<string, { threadId: string; userName: string; timestamp: number }>;

  // New thread modal
  isNewThreadModalOpen: boolean;
  openNewThreadModal: () => void;
  closeNewThreadModal: () => void;

  // Search
  searchQuery: string;
  setSearchQuery: (query: string) => void;
  isSearching: boolean;
}

const MessagingContext = createContext<MessagingContextValue | undefined>(undefined);

export function useMessaging() {
  const context = useContext(MessagingContext);
  if (!context) {
    throw new Error('useMessaging must be used within a MessagingProvider');
  }
  return context;
}

interface MessagingProviderProps {
  children: React.ReactNode;
}

export function MessagingProvider({ children }: MessagingProviderProps) {
  const queryClient = useQueryClient();
  const [selectedThreadId, setSelectedThreadId] = useState<string | null>(null);
  const [unreadCount] = useState(0);
  const [typingUsers, setTypingUsers] = useState<Map<string, { threadId: string; userName: string; timestamp: number }>>(new Map());
  const [isNewThreadModalOpen, setIsNewThreadModalOpen] = useState(false);
  const [searchQuery, setSearchQuery] = useState('');

  // WebSocket handler for real-time updates
  const handleWebSocketMessage = useCallback((wsMessage: any) => {
    switch (wsMessage.type) {
      case 'chat_message': {
        const { threadId, message } = wsMessage.payload;

        // Update thread cache with new message
        queryClient.setQueryData(messagingKeys.thread(threadId), (old: any) => {
          if (!old) return old;
          return {
            ...old,
            messages: [...old.messages, message],
            thread: {
              ...old.thread,
              messageCount: old.thread.messageCount + 1,
              lastMessageAt: message.createdAt,
              lastMessage: message,
            },
          };
        });

        // Invalidate threads list
        queryClient.invalidateQueries({ queryKey: messagingKeys.threads() });
        queryClient.invalidateQueries({ queryKey: messagingKeys.unread() });
        break;
      }

      case 'chat_typing': {
        const { threadId, userId, userName, isTyping } = wsMessage.payload;

        setTypingUsers((prev) => {
          const next = new Map(prev);
          if (isTyping) {
            next.set(userId, { threadId, userName, timestamp: Date.now() });
          } else {
            next.delete(userId);
          }
          return next;
        });
        break;
      }

      case 'chat_read': {
        const { threadId } = wsMessage.payload;
        queryClient.invalidateQueries({ queryKey: messagingKeys.thread(threadId) });
        queryClient.invalidateQueries({ queryKey: messagingKeys.unread() });
        break;
      }
    }
  }, [queryClient]);

  useWebSocket({
    onMessage: handleWebSocketMessage,
  });

  // Clean up stale typing indicators
  useEffect(() => {
    const interval = setInterval(() => {
      const now = Date.now();
      setTypingUsers((prev) => {
        const next = new Map(prev);
        let changed = false;
        for (const [userId, data] of prev) {
          if (now - data.timestamp > 5000) {
            next.delete(userId);
            changed = true;
          }
        }
        return changed ? next : prev;
      });
    }, 1000);

    return () => clearInterval(interval);
  }, []);

  const onNewMessage = useCallback((threadId: string, message: Message) => {
    // Update cache
    queryClient.setQueryData(messagingKeys.thread(threadId), (old: any) => {
      if (!old) return old;
      return {
        ...old,
        messages: [...old.messages, message],
      };
    });
  }, [queryClient]);

  const refreshUnreadCount = useCallback(() => {
    queryClient.invalidateQueries({ queryKey: messagingKeys.unread() });
  }, [queryClient]);

  const openNewThreadModal = useCallback(() => {
    setIsNewThreadModalOpen(true);
  }, []);

  const closeNewThreadModal = useCallback(() => {
    setIsNewThreadModalOpen(false);
  }, []);

  const value: MessagingContextValue = {
    selectedThreadId,
    setSelectedThreadId,
    onNewMessage,
    unreadCount,
    refreshUnreadCount,
    typingUsers,
    isNewThreadModalOpen,
    openNewThreadModal,
    closeNewThreadModal,
    searchQuery,
    setSearchQuery,
    isSearching: searchQuery.length >= 2,
  };

  return (
    <MessagingContext.Provider value={value}>
      {children}
    </MessagingContext.Provider>
  );
}
