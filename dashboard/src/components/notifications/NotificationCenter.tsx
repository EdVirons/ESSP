import * as React from 'react';
import { Bell, Check, RefreshCw } from 'lucide-react';
import { useNavigate } from 'react-router-dom';
import { Button } from '@/components/ui/button';
import { Popover } from '@/components/ui/popover';
import { NotificationItem } from './NotificationItem';
import { useNotifications, useUnreadCount, useMarkNotificationsRead } from '@/api/notifications';
import { cn } from '@/lib/utils';
import type { Notification } from '@/types/notification';

export function NotificationCenter() {
  const [open, setOpen] = React.useState(false);
  const navigate = useNavigate();

  const { data: notificationsData, isLoading, refetch } = useNotifications({ limit: 20 });
  const { data: unreadData } = useUnreadCount();
  const markRead = useMarkNotificationsRead();

  const notifications = notificationsData?.items ?? [];
  const unreadCount = unreadData?.count ?? 0;

  const handleNotificationClick = (notification: Notification) => {
    // Navigate to the relevant page based on notification type
    switch (notification.type) {
      case 'incident':
        navigate(`/incidents/${notification.target}`);
        break;
      case 'work_order':
        navigate(`/work-orders/${notification.target}`);
        break;
      case 'project':
        navigate(`/projects/${notification.target}`);
        break;
      case 'service_shop':
        navigate(`/service-shops/${notification.target}`);
        break;
      default:
        // Stay on current page
        break;
    }
    setOpen(false);
  };

  const handleMarkAllRead = () => {
    markRead.mutate('all');
  };

  return (
    <Popover
      open={open}
      onClose={() => setOpen(false)}
      align="end"
      trigger={
        <Button
          variant="ghost"
          size="icon"
          onClick={() => setOpen(!open)}
          className="relative"
        >
          <Bell className="h-5 w-5" />
          {unreadCount > 0 && (
            <span className="absolute -top-1 -right-1 flex h-5 w-5 items-center justify-center rounded-full bg-red-500 text-[10px] font-medium text-white">
              {unreadCount > 99 ? '99+' : unreadCount}
            </span>
          )}
        </Button>
      }
    >
      <div className="w-[380px] max-h-[480px] flex flex-col">
        {/* Header */}
        <div className="flex items-center justify-between px-4 py-3 border-b border-gray-100">
          <h3 className="font-semibold text-gray-900">Notifications</h3>
          <div className="flex items-center gap-1">
            <Button
              variant="ghost"
              size="sm"
              onClick={() => refetch()}
              disabled={isLoading}
              className="h-8 px-2"
            >
              <RefreshCw className={cn('h-4 w-4', isLoading && 'animate-spin')} />
            </Button>
            {unreadCount > 0 && (
              <Button
                variant="ghost"
                size="sm"
                onClick={handleMarkAllRead}
                className="h-8 text-xs"
              >
                <Check className="h-3 w-3 mr-1" />
                Mark all read
              </Button>
            )}
          </div>
        </div>

        {/* Notifications list */}
        <div className="flex-1 overflow-y-auto">
          {isLoading ? (
            <div className="flex items-center justify-center py-8">
              <RefreshCw className="h-5 w-5 animate-spin text-gray-400" />
            </div>
          ) : notifications.length === 0 ? (
            <div className="flex flex-col items-center justify-center py-8 text-gray-500">
              <Bell className="h-8 w-8 mb-2 text-gray-300" />
              <p className="text-sm">No notifications yet</p>
            </div>
          ) : (
            <div className="divide-y divide-gray-100">
              {notifications.map((notification) => (
                <NotificationItem
                  key={notification.id}
                  notification={notification}
                  onClick={() => handleNotificationClick(notification)}
                />
              ))}
            </div>
          )}
        </div>

        {/* Footer */}
        {notifications.length > 0 && (
          <div className="px-4 py-2 border-t border-gray-100">
            <Button
              variant="ghost"
              size="sm"
              className="w-full text-xs text-gray-500 hover:text-gray-700"
              onClick={() => {
                navigate('/audit-logs');
                setOpen(false);
              }}
            >
              View all activity
            </Button>
          </div>
        )}
      </div>
    </Popover>
  );
}
