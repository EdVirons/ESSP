import { useState, useEffect, useCallback } from 'react';
import { useWebSocket } from './useWebSocket';

export type PresenceStatus = 'online' | 'away' | 'offline';

interface UserPresence {
  userId: string;
  status: PresenceStatus;
  lastSeen: number;
}

interface PresenceState {
  [userId: string]: UserPresence;
}

interface UsePresenceOptions {
  enabled?: boolean;
  heartbeatIntervalMs?: number;
  staleTimeoutMs?: number;
}

interface UsePresenceReturn {
  presenceMap: PresenceState;
  getPresence: (userId: string) => PresenceStatus;
  setStatus: (status: PresenceStatus) => void;
  isOnline: (userId: string) => boolean;
  onlineCount: number;
}

export function usePresence(options: UsePresenceOptions = {}): UsePresenceReturn {
  const {
    enabled = true,
    heartbeatIntervalMs = 30000,
    staleTimeoutMs = 120000,
  } = options;

  const [presenceMap, setPresenceMap] = useState<PresenceState>({});
  const [currentStatus, setCurrentStatus] = useState<PresenceStatus>('online');

  const { send, isConnected } = useWebSocket({
    enabled,
    onMessage: (message) => {
      if (message.type === 'presence_update') {
        const payload = message.payload as {
          userId: string;
          status: PresenceStatus;
        };

        setPresenceMap((prev) => ({
          ...prev,
          [payload.userId]: {
            userId: payload.userId,
            status: payload.status,
            lastSeen: Date.now(),
          },
        }));
      }
    },
  });

  // Send heartbeat
  useEffect(() => {
    if (!enabled || !isConnected) return;

    const sendHeartbeat = () => {
      send({
        action: 'presence',
        status: currentStatus,
      });
    };

    // Send initial presence
    sendHeartbeat();

    // Set up heartbeat interval
    const interval = setInterval(sendHeartbeat, heartbeatIntervalMs);

    return () => {
      clearInterval(interval);
      // Send offline status on unmount
      send({
        action: 'presence',
        status: 'offline',
      });
    };
  }, [enabled, isConnected, currentStatus, heartbeatIntervalMs, send]);

  // Clean up stale presence entries
  useEffect(() => {
    const interval = setInterval(() => {
      const now = Date.now();
      setPresenceMap((prev) => {
        const updated = { ...prev };
        let changed = false;

        for (const userId in updated) {
          if (now - updated[userId].lastSeen > staleTimeoutMs) {
            updated[userId] = { ...updated[userId], status: 'offline' };
            changed = true;
          }
        }

        return changed ? updated : prev;
      });
    }, 30000);

    return () => clearInterval(interval);
  }, [staleTimeoutMs]);

  // Get presence for a user
  const getPresence = useCallback(
    (userId: string): PresenceStatus => {
      return presenceMap[userId]?.status || 'offline';
    },
    [presenceMap]
  );

  // Check if user is online
  const isOnline = useCallback(
    (userId: string): boolean => {
      return presenceMap[userId]?.status === 'online';
    },
    [presenceMap]
  );

  // Set current user's status
  const setStatus = useCallback(
    (status: PresenceStatus) => {
      setCurrentStatus(status);
      if (isConnected) {
        send({
          action: 'presence',
          status,
        });
      }
    },
    [isConnected, send]
  );

  // Count online users
  const onlineCount = Object.values(presenceMap).filter(
    (p) => p.status === 'online'
  ).length;

  return {
    presenceMap,
    getPresence,
    setStatus,
    isOnline,
    onlineCount,
  };
}
