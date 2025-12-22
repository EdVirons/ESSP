import { useEffect, useRef, useCallback, useState } from 'react';
import { WS_BASE_URL } from '@/lib/constants';
import type { WSMessage } from '@/types/notification';

interface UseWebSocketOptions {
  onMessage?: (message: WSMessage) => void;
  onConnect?: () => void;
  onDisconnect?: () => void;
  reconnectInterval?: number;
  maxReconnectAttempts?: number;
  enabled?: boolean;
}

interface UseWebSocketReturn {
  isConnected: boolean;
  reconnectAttempt: number;
  send: (message: object) => void;
  disconnect: () => void;
  connect: () => void;
}

export function useWebSocket(options: UseWebSocketOptions = {}): UseWebSocketReturn {
  const {
    onMessage,
    onConnect,
    onDisconnect,
    reconnectInterval = 3000,
    maxReconnectAttempts = 10,
    enabled = true,
  } = options;

  const [isConnected, setIsConnected] = useState(false);
  const [reconnectAttempt, setReconnectAttempt] = useState(0);
  const wsRef = useRef<WebSocket | null>(null);
  const reconnectTimeoutRef = useRef<number | null>(null);
  const mountedRef = useRef(true);

  const connectRef = useRef<(() => void) | null>(null);

  const doConnect = useCallback(() => {
    if (!enabled || wsRef.current?.readyState === WebSocket.OPEN) {
      return;
    }

    // Get tenant from localStorage
    const tenantId = localStorage.getItem('tenant_id') || 'demo-tenant';
    const userId = localStorage.getItem('user_id') || 'anonymous';

    // Build WebSocket URL with query params
    const url = `${WS_BASE_URL}?tenant=${encodeURIComponent(tenantId)}&user=${encodeURIComponent(userId)}`;

    try {
      const ws = new WebSocket(url);

      ws.onopen = () => {
        if (!mountedRef.current) return;
        setIsConnected(true);
        setReconnectAttempt(0);
        onConnect?.();
        console.log('[WebSocket] Connected');
      };

      ws.onclose = (event) => {
        if (!mountedRef.current) return;
        setIsConnected(false);
        onDisconnect?.();
        console.log('[WebSocket] Disconnected', event.code, event.reason);

        // Attempt to reconnect with exponential backoff
        setReconnectAttempt((currentAttempt) => {
          if (currentAttempt < maxReconnectAttempts && enabled) {
            const delay = Math.min(
              reconnectInterval * Math.pow(1.5, currentAttempt),
              30000 // Max 30 seconds
            );
            console.log(`[WebSocket] Reconnecting in ${delay}ms (attempt ${currentAttempt + 1})`);

            reconnectTimeoutRef.current = window.setTimeout(() => {
              if (mountedRef.current) {
                connectRef.current?.();
              }
            }, delay);
            return currentAttempt + 1;
          }
          return currentAttempt;
        });
      };

      ws.onerror = (error) => {
        console.error('[WebSocket] Error:', error);
      };

      ws.onmessage = (event) => {
        if (!mountedRef.current) return;
        try {
          const message: WSMessage = JSON.parse(event.data);
          onMessage?.(message);
        } catch {
          console.error('[WebSocket] Failed to parse message:', event.data);
        }
      };

      wsRef.current = ws;
    } catch (error) {
      console.error('[WebSocket] Failed to connect:', error);
    }
  }, [enabled, onMessage, onConnect, onDisconnect, reconnectInterval, maxReconnectAttempts]);

  // Update ref in effect to avoid updating during render
  useEffect(() => {
    connectRef.current = doConnect;
  }, [doConnect]);

  const connect = useCallback(() => {
    connectRef.current?.();
  }, []);

  const disconnect = useCallback(() => {
    if (reconnectTimeoutRef.current) {
      clearTimeout(reconnectTimeoutRef.current);
      reconnectTimeoutRef.current = null;
    }
    if (wsRef.current) {
      wsRef.current.close();
      wsRef.current = null;
    }
    setIsConnected(false);
    setReconnectAttempt(0);
  }, []);

  const send = useCallback((message: object) => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify(message));
    } else {
      console.warn('[WebSocket] Cannot send message - not connected');
    }
  }, []);

  // Connect on mount
  useEffect(() => {
    mountedRef.current = true;
    if (enabled) {
      connect();
    }

    return () => {
      mountedRef.current = false;
      disconnect();
    };
  }, [enabled]); // eslint-disable-line react-hooks/exhaustive-deps

  return { isConnected, reconnectAttempt, send, disconnect, connect };
}
