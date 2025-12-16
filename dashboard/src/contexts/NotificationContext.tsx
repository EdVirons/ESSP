import * as React from 'react';
import { useQueryClient } from '@tanstack/react-query';
import { useWebSocket } from '@/hooks/useWebSocket';
import { toast } from '@/lib/toast';
import type { Notification, WSMessage } from '@/types/notification';

interface NotificationContextValue {
  notifications: Notification[];
  unreadCount: number;
  isConnected: boolean;
  addNotification: (notification: Notification) => void;
  clearNotifications: () => void;
}

const NotificationContext = React.createContext<NotificationContextValue | null>(null);

const MAX_NOTIFICATIONS = 50;

export function NotificationProvider({ children }: { children: React.ReactNode }) {
  const queryClient = useQueryClient();
  const [notifications, setNotifications] = React.useState<Notification[]>([]);

  const unreadCount = React.useMemo(
    () => notifications.filter((n) => !n.read).length,
    [notifications]
  );

  const addNotification = React.useCallback((notification: Notification) => {
    setNotifications((prev) => {
      // Check if notification already exists
      if (prev.some((n) => n.id === notification.id)) {
        return prev;
      }
      // Add to front, keep max notifications
      return [notification, ...prev].slice(0, MAX_NOTIFICATIONS);
    });
  }, []);

  const clearNotifications = React.useCallback(() => {
    setNotifications([]);
  }, []);

  const handleWSMessage = React.useCallback(
    (message: WSMessage) => {
      if (message.type === 'notification') {
        const payload = message.payload as Notification;

        // Add to local notifications
        addNotification({
          ...payload,
          read: false,
        });

        // Show toast for important events
        const entityLabel = getEntityLabel(payload.type);
        if (payload.action === 'create') {
          toast.info(`New ${entityLabel}`, payload.summary);
        }

        // Invalidate relevant queries to refresh data
        switch (payload.type) {
          case 'incident':
            queryClient.invalidateQueries({ queryKey: ['incidents'] });
            break;
          case 'work_order':
            queryClient.invalidateQueries({ queryKey: ['work-orders'] });
            break;
          case 'project':
            queryClient.invalidateQueries({ queryKey: ['projects'] });
            break;
          case 'service_shop':
            queryClient.invalidateQueries({ queryKey: ['service-shops'] });
            break;
          default:
            break;
        }

        // Also invalidate notifications query
        queryClient.invalidateQueries({ queryKey: ['notifications'] });
      }
    },
    [queryClient, addNotification]
  );

  const { isConnected } = useWebSocket({
    onMessage: handleWSMessage,
    onConnect: () => {
      console.log('[NotificationContext] WebSocket connected');
    },
    onDisconnect: () => {
      console.log('[NotificationContext] WebSocket disconnected');
    },
    enabled: true,
  });

  const value = React.useMemo(
    () => ({
      notifications,
      unreadCount,
      isConnected,
      addNotification,
      clearNotifications,
    }),
    [notifications, unreadCount, isConnected, addNotification, clearNotifications]
  );

  return (
    <NotificationContext.Provider value={value}>
      {children}
    </NotificationContext.Provider>
  );
}

export function useNotificationContext() {
  const context = React.useContext(NotificationContext);
  if (!context) {
    throw new Error('useNotificationContext must be used within NotificationProvider');
  }
  return context;
}

function getEntityLabel(type: string): string {
  switch (type) {
    case 'incident':
      return 'Incident';
    case 'work_order':
      return 'Work Order';
    case 'project':
      return 'Project';
    case 'service_shop':
      return 'Service Shop';
    case 'device':
      return 'Device';
    default:
      return type;
  }
}
