import { formatDistanceToNow } from 'date-fns';
import { AlertTriangle, Wrench, Layers, Store, Laptop, FileText } from 'lucide-react';
import { cn } from '@/lib/utils';
import type { Notification } from '@/types/notification';

interface NotificationItemProps {
  notification: Notification;
  onClick?: () => void;
}

const typeIcons: Record<string, React.ElementType> = {
  incident: AlertTriangle,
  work_order: Wrench,
  program: Layers,
  service_shop: Store,
  device: Laptop,
};

const typeColors: Record<string, string> = {
  incident: 'text-red-500 bg-red-50',
  work_order: 'text-blue-500 bg-blue-50',
  program: 'text-purple-500 bg-purple-50',
  service_shop: 'text-green-500 bg-green-50',
  device: 'text-orange-500 bg-orange-50',
};

const actionColors: Record<string, string> = {
  create: 'text-green-600',
  update: 'text-blue-600',
  delete: 'text-red-600',
};

export function NotificationItem({ notification, onClick }: NotificationItemProps) {
  const Icon = typeIcons[notification.type] || FileText;
  const iconColor = typeColors[notification.type] || 'text-gray-500 bg-gray-50';

  const timeAgo = formatDistanceToNow(new Date(notification.timestamp), { addSuffix: true });

  return (
    <button
      onClick={onClick}
      className={cn(
        'w-full flex items-start gap-3 px-4 py-3 text-left hover:bg-gray-50 transition-colors',
        !notification.read && 'bg-blue-50/50'
      )}
    >
      <div className={cn('flex-shrink-0 p-2 rounded-lg', iconColor)}>
        <Icon className="h-4 w-4" />
      </div>
      <div className="flex-1 min-w-0">
        <p className="text-sm text-gray-900 line-clamp-2">{notification.summary}</p>
        <div className="flex items-center gap-2 mt-1">
          <span className={cn('text-xs font-medium capitalize', actionColors[notification.action])}>
            {notification.action}
          </span>
          <span className="text-xs text-gray-400">by {notification.actor || 'System'}</span>
          <span className="text-xs text-gray-400">{timeAgo}</span>
        </div>
      </div>
      {!notification.read && (
        <div className="flex-shrink-0">
          <div className="h-2 w-2 rounded-full bg-blue-500" />
        </div>
      )}
    </button>
  );
}
