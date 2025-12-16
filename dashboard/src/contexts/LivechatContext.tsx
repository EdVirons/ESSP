import React, { createContext, useContext, useState, useCallback, useEffect } from 'react';
import { useQueryClient } from '@tanstack/react-query';
import { useWebSocket } from '@/hooks/useWebSocket';
import { livechatKeys } from '@/hooks/useLivechat';
import { messagingKeys } from '@/hooks/useMessages';
import type { ChatSession, ChatSessionStatus } from '@/types/livechat';

interface LivechatContextValue {
  // Current session (for school contacts)
  currentSession: ChatSession | null;
  setCurrentSession: (session: ChatSession | null) => void;

  // Widget state
  isWidgetOpen: boolean;
  openWidget: () => void;
  closeWidget: () => void;
  toggleWidget: () => void;

  // Queue position
  queuePosition: number | null;

  // Agent state
  isAgentAvailable: boolean;
  setAgentAvailable: (available: boolean) => void;

  // Real-time updates
  onSessionUpdate: (sessionId: string, status: ChatSessionStatus, agentId?: string, agentName?: string) => void;

  // Active chat for agents
  activeSessionId: string | null;
  setActiveSessionId: (id: string | null) => void;

  // Typing
  typingInSession: Map<string, { userName: string; timestamp: number }>;

  // Notifications
  newChatNotification: { sessionId: string; contactName: string; subject: string } | null;
  clearNewChatNotification: () => void;
}

const LivechatContext = createContext<LivechatContextValue | undefined>(undefined);

export function useLivechatContext() {
  const context = useContext(LivechatContext);
  if (!context) {
    throw new Error('useLivechatContext must be used within a LivechatProvider');
  }
  return context;
}

interface LivechatProviderProps {
  children: React.ReactNode;
}

export function LivechatProvider({ children }: LivechatProviderProps) {
  const queryClient = useQueryClient();

  const [currentSession, setCurrentSession] = useState<ChatSession | null>(null);
  const [isWidgetOpen, setIsWidgetOpen] = useState(false);
  const [queuePosition, setQueuePosition] = useState<number | null>(null);
  const [isAgentAvailable, setAgentAvailable] = useState(false);
  const [activeSessionId, setActiveSessionId] = useState<string | null>(null);
  const [typingInSession, setTypingInSession] = useState<Map<string, { userName: string; timestamp: number }>>(new Map());
  const [newChatNotification, setNewChatNotification] = useState<{ sessionId: string; contactName: string; subject: string } | null>(null);

  // WebSocket handler for real-time updates
  const handleWebSocketMessage = useCallback((wsMessage: any) => {
    switch (wsMessage.type) {
      case 'chat_session_update': {
        const { sessionId, status, agentId, agentName, queuePosition: newPosition } = wsMessage.payload;

        // Update current session if it matches
        if (currentSession?.id === sessionId) {
          setCurrentSession((prev) => {
            if (!prev) return prev;
            return {
              ...prev,
              status,
              assignedAgentId: agentId,
              assignedAgentName: agentName,
              queuePosition: newPosition,
            };
          });

          if (newPosition !== undefined) {
            setQueuePosition(newPosition);
          }
        }

        // Invalidate queries
        queryClient.invalidateQueries({ queryKey: livechatKeys.queue() });
        queryClient.invalidateQueries({ queryKey: livechatKeys.activeChats() });
        break;
      }

      case 'new_chat_waiting': {
        const { sessionId, schoolContactName, subject } = wsMessage.payload;

        // Show notification for agents
        setNewChatNotification({
          sessionId,
          contactName: schoolContactName,
          subject,
        });

        // Invalidate queue
        queryClient.invalidateQueries({ queryKey: livechatKeys.queue() });
        break;
      }

      case 'chat_message': {
        const { threadId, message } = wsMessage.payload;

        // If this is for the current session's thread, update it
        if (currentSession?.threadId === threadId) {
          queryClient.setQueryData(messagingKeys.thread(threadId), (old: any) => {
            if (!old) return old;
            return {
              ...old,
              messages: [...old.messages, message],
            };
          });
        }

        // Invalidate active chats for agents
        queryClient.invalidateQueries({ queryKey: livechatKeys.activeChats() });
        break;
      }

      case 'chat_typing': {
        const { threadId, userId, userName, isTyping } = wsMessage.payload;

        // Check if this is for the current session
        if (currentSession?.threadId === threadId) {
          setTypingInSession((prev) => {
            const next = new Map(prev);
            if (isTyping) {
              next.set(userId, { userName, timestamp: Date.now() });
            } else {
              next.delete(userId);
            }
            return next;
          });
        }
        break;
      }

      case 'presence_update': {
        // Handle agent presence updates
        queryClient.invalidateQueries({ queryKey: livechatKeys.availability() });
        break;
      }
    }
  }, [currentSession, queryClient]);

  useWebSocket({
    onMessage: handleWebSocketMessage,
  });

  // Clean up stale typing indicators
  useEffect(() => {
    const interval = setInterval(() => {
      const now = Date.now();
      setTypingInSession((prev) => {
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

  const openWidget = useCallback(() => {
    setIsWidgetOpen(true);
  }, []);

  const closeWidget = useCallback(() => {
    setIsWidgetOpen(false);
  }, []);

  const toggleWidget = useCallback(() => {
    setIsWidgetOpen((prev) => !prev);
  }, []);

  const onSessionUpdate = useCallback((
    sessionId: string,
    status: ChatSessionStatus,
    agentId?: string,
    agentName?: string
  ) => {
    if (currentSession?.id === sessionId) {
      setCurrentSession((prev) => {
        if (!prev) return prev;
        return {
          ...prev,
          status,
          assignedAgentId: agentId,
          assignedAgentName: agentName,
        };
      });
    }
  }, [currentSession]);

  const clearNewChatNotification = useCallback(() => {
    setNewChatNotification(null);
  }, []);

  const value: LivechatContextValue = {
    currentSession,
    setCurrentSession,
    isWidgetOpen,
    openWidget,
    closeWidget,
    toggleWidget,
    queuePosition,
    isAgentAvailable,
    setAgentAvailable,
    onSessionUpdate,
    activeSessionId,
    setActiveSessionId,
    typingInSession,
    newChatNotification,
    clearNewChatNotification,
  };

  return (
    <LivechatContext.Provider value={value}>
      {children}
    </LivechatContext.Provider>
  );
}
