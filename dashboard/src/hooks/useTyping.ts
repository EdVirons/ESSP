import { useState, useCallback, useRef, useEffect } from 'react';
import { useWebSocket } from './useWebSocket';

interface TypingState {
  [userId: string]: {
    userName: string;
    isTyping: boolean;
    timestamp: number;
  };
}

interface UseTypingOptions {
  threadId: string;
  userId: string;
  userName: string;
  debounceMs?: number;
  timeoutMs?: number;
}

interface UseTypingReturn {
  typingUsers: { userId: string; userName: string }[];
  setTyping: (isTyping: boolean) => void;
  handleInputChange: () => void;
}

export function useTyping({
  threadId,
  userId,
  userName,
  debounceMs = 300,
  timeoutMs = 3000,
}: UseTypingOptions): UseTypingReturn {
  const { send, isConnected } = useWebSocket({
    onMessage: (message) => {
      if (message.type === 'chat_typing') {
        const payload = message.payload as {
          threadId: string;
          userId: string;
          userName: string;
          isTyping: boolean;
        };

        if (payload.threadId !== threadId || payload.userId === userId) {
          return;
        }

        setTypingState((prev) => {
          if (payload.isTyping) {
            return {
              ...prev,
              [payload.userId]: {
                userName: payload.userName,
                isTyping: true,
                timestamp: Date.now(),
              },
            };
          } else {
            // eslint-disable-next-line @typescript-eslint/no-unused-vars
            const { [payload.userId]: _, ...rest } = prev;
            return rest;
          }
        });
      }
    },
  });

  const [typingState, setTypingState] = useState<TypingState>({});
  const [isTyping, setIsTypingState] = useState(false);
  const debounceRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const stopTypingRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  // Clean up stale typing indicators
  useEffect(() => {
    const interval = setInterval(() => {
      const now = Date.now();
      setTypingState((prev) => {
        const updated = { ...prev };
        let changed = false;
        for (const key in updated) {
          if (now - updated[key].timestamp > timeoutMs) {
            delete updated[key];
            changed = true;
          }
        }
        return changed ? updated : prev;
      });
    }, 1000);

    return () => clearInterval(interval);
  }, [timeoutMs]);

  // Send typing indicator
  const setTyping = useCallback(
    (typing: boolean) => {
      if (!isConnected) return;

      setIsTypingState(typing);
      send({
        action: 'typing',
        threadId,
        userId,
        userName,
        isTyping: typing,
      });
    },
    [isConnected, send, threadId, userId, userName]
  );

  // Handle input change with debounce
  const handleInputChange = useCallback(() => {
    // Clear existing debounce
    if (debounceRef.current) {
      clearTimeout(debounceRef.current);
    }

    // Clear existing stop typing timeout
    if (stopTypingRef.current) {
      clearTimeout(stopTypingRef.current);
    }

    // If not currently typing, send typing = true
    if (!isTyping) {
      setTyping(true);
    }

    // Set debounced stop typing
    debounceRef.current = setTimeout(() => {
      // Set timeout to stop typing after inactivity
      stopTypingRef.current = setTimeout(() => {
        setTyping(false);
      }, timeoutMs - debounceMs);
    }, debounceMs);
  }, [isTyping, setTyping, debounceMs, timeoutMs]);

  // Cleanup on unmount
  useEffect(() => {
    return () => {
      if (debounceRef.current) clearTimeout(debounceRef.current);
      if (stopTypingRef.current) clearTimeout(stopTypingRef.current);
      if (isTyping) {
        setTyping(false);
      }
    };
  }, [isTyping, setTyping]);

  // Convert typing state to array
  const typingUsers = Object.entries(typingState).map(([userId, state]) => ({
    userId,
    userName: state.userName,
  }));

  return {
    typingUsers,
    setTyping,
    handleInputChange,
  };
}
