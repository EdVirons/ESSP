import { Link } from 'react-router-dom';
import {
  Bell,
  AlertTriangle,
  Clock,
  CheckCircle,
  Info,
  ArrowRight,
  X,
} from 'lucide-react';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { cn } from '@/lib/utils';

interface Notification {
  id: string;
  title: string;
  message: string;
  type: 'warning' | 'info' | 'success' | 'alert';
  timestamp: string;
  read: boolean;
  href?: string;
}

// Mock data - in production, this would come from API
const mockNotifications: Notification[] = [
  {
    id: 'n1',
    title: 'SLA Warning',
    message: 'Work order WO-001 is approaching SLA deadline (2 hours remaining)',
    type: 'warning',
    timestamp: '10 min ago',
    read: false,
    href: '/work-orders/wo-001',
  },
  {
    id: 'n2',
    title: 'Sync Completed',
    message: 'Device inventory sync completed successfully. 234 devices updated.',
    type: 'success',
    timestamp: '1 hour ago',
    read: false,
  },
  {
    id: 'n3',
    title: 'Low Stock Alert',
    message: '5 parts are below minimum stock level',
    type: 'alert',
    timestamp: '2 hours ago',
    read: false,
    href: '/parts-catalog?filter=low-stock',
  },
  {
    id: 'n4',
    title: 'New Assignment',
    message: 'You have been assigned to incident INC-042',
    type: 'info',
    timestamp: '3 hours ago',
    read: true,
    href: '/incidents/inc-042',
  },
  {
    id: 'n5',
    title: 'Approval Required',
    message: 'Work order WO-015 requires your approval',
    type: 'info',
    timestamp: '5 hours ago',
    read: true,
    href: '/approvals',
  },
];

const typeConfig = {
  warning: {
    icon: Clock,
    color: 'text-amber-600',
    bgColor: 'bg-amber-50',
    badgeColor: 'bg-amber-100 text-amber-800',
  },
  info: {
    icon: Info,
    color: 'text-blue-600',
    bgColor: 'bg-blue-50',
    badgeColor: 'bg-blue-100 text-blue-800',
  },
  success: {
    icon: CheckCircle,
    color: 'text-green-600',
    bgColor: 'bg-green-50',
    badgeColor: 'bg-green-100 text-green-800',
  },
  alert: {
    icon: AlertTriangle,
    color: 'text-red-600',
    bgColor: 'bg-red-50',
    badgeColor: 'bg-red-100 text-red-800',
  },
};

interface NotificationsPanelProps {
  notifications?: Notification[];
  isLoading?: boolean;
  onDismiss?: (id: string) => void;
  onMarkAllRead?: () => void;
}

export function NotificationsPanel({
  notifications = mockNotifications,
  isLoading,
  onDismiss,
  onMarkAllRead,
}: NotificationsPanelProps) {
  const unreadCount = notifications.filter((n) => !n.read).length;

  if (isLoading) {
    return (
      <Card>
        <CardHeader className="pb-3">
          <CardTitle className="flex items-center gap-2 text-base">
            <Bell className="h-5 w-5" />
            Notifications
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-3">
            {[1, 2, 3].map((i) => (
              <div key={i} className="skeleton h-16 rounded" />
            ))}
          </div>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card>
      <CardHeader className="pb-3">
        <div className="flex items-center justify-between">
          <CardTitle className="flex items-center gap-2 text-base">
            <Bell className="h-5 w-5" />
            Notifications
            {unreadCount > 0 && (
              <Badge className="bg-red-100 text-red-800 ml-1">
                {unreadCount} new
              </Badge>
            )}
          </CardTitle>
          {unreadCount > 0 && onMarkAllRead && (
            <Button
              variant="ghost"
              size="sm"
              onClick={onMarkAllRead}
              className="text-xs text-gray-500 hover:text-gray-700"
            >
              Mark all read
            </Button>
          )}
        </div>
      </CardHeader>
      <CardContent>
        {notifications.length === 0 ? (
          <div className="text-center py-8 text-gray-500">
            <Bell className="h-12 w-12 mx-auto mb-3 text-gray-300" />
            <p className="font-medium">All caught up!</p>
            <p className="text-sm">No new notifications.</p>
          </div>
        ) : (
          <div className="space-y-2">
            {notifications.slice(0, 5).map((notification) => {
              const config = typeConfig[notification.type];
              const Icon = config.icon;

              const content = (
                <div
                  className={cn(
                    'flex items-start gap-3 p-3 rounded-lg transition-colors',
                    notification.read
                      ? 'bg-gray-50 hover:bg-gray-100'
                      : 'bg-white border border-gray-200 hover:border-gray-300',
                    notification.href && 'cursor-pointer'
                  )}
                >
                  <div className={cn('rounded-full p-1.5 mt-0.5', config.bgColor)}>
                    <Icon className={cn('h-4 w-4', config.color)} />
                  </div>
                  <div className="flex-1 min-w-0">
                    <div className="flex items-start justify-between gap-2">
                      <div className="min-w-0">
                        <p
                          className={cn(
                            'text-sm truncate',
                            notification.read
                              ? 'font-medium text-gray-700'
                              : 'font-semibold text-gray-900'
                          )}
                        >
                          {notification.title}
                        </p>
                        <p className="text-xs text-gray-500 line-clamp-2 mt-0.5">
                          {notification.message}
                        </p>
                      </div>
                      {!notification.read && (
                        <div className="h-2 w-2 rounded-full bg-blue-500 mt-1.5 flex-shrink-0" />
                      )}
                    </div>
                    <p className="text-xs text-gray-400 mt-1">
                      {notification.timestamp}
                    </p>
                  </div>
                  {onDismiss && (
                    <button
                      onClick={(e) => {
                        e.preventDefault();
                        e.stopPropagation();
                        onDismiss(notification.id);
                      }}
                      className="text-gray-400 hover:text-gray-600 p-1"
                    >
                      <X className="h-4 w-4" />
                    </button>
                  )}
                </div>
              );

              if (notification.href) {
                return (
                  <Link key={notification.id} to={notification.href}>
                    {content}
                  </Link>
                );
              }

              return <div key={notification.id}>{content}</div>;
            })}

            {notifications.length > 5 && (
              <Link to="/notifications">
                <Button variant="outline" className="w-full mt-2">
                  View all {notifications.length} notifications
                  <ArrowRight className="h-4 w-4 ml-2" />
                </Button>
              </Link>
            )}
          </div>
        )}
      </CardContent>
    </Card>
  );
}
