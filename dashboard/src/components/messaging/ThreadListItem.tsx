import { formatDistanceToNow } from 'date-fns';
import { MessageSquare, Tag, AlertTriangle, MessageCircle } from 'lucide-react';
import { Badge } from '@/components/ui/badge';
import { cn } from '@/lib/utils';
import type { MessageThread } from '@/types/messaging';

interface ThreadListItemProps {
  thread: MessageThread;
  isSelected: boolean;
  onClick: () => void;
  isSchoolContact: boolean;
}

export function ThreadListItem({
  thread,
  isSelected,
  onClick,
  isSchoolContact,
}: ThreadListItemProps) {
  const unreadCount = isSchoolContact
    ? thread.unreadCountSchool
    : thread.unreadCountSupport;

  const getThreadTypeIcon = () => {
    switch (thread.threadType) {
      case 'incident':
        return <AlertTriangle className="h-4 w-4 text-amber-500" />;
      case 'livechat':
        return <MessageCircle className="h-4 w-4 text-green-500" />;
      default:
        return <MessageSquare className="h-4 w-4 text-blue-500" />;
    }
  };

  const getStatusBadge = () => {
    if (thread.status === 'closed') {
      return (
        <Badge variant="secondary" className="text-xs">
          Closed
        </Badge>
      );
    }
    return null;
  };

  const lastMessageTime = thread.lastMessageAt
    ? formatDistanceToNow(new Date(thread.lastMessageAt), { addSuffix: true })
    : formatDistanceToNow(new Date(thread.createdAt), { addSuffix: true });

  return (
    <div
      onClick={onClick}
      className={cn(
        'p-3 border-b border-gray-100 cursor-pointer transition-colors hover:bg-gray-50',
        isSelected && 'bg-cyan-50 border-l-2 border-l-cyan-500',
        unreadCount > 0 && !isSelected && 'bg-blue-50/50'
      )}
    >
      <div className="flex items-start gap-3">
        <div className="flex-shrink-0 mt-1">
          {getThreadTypeIcon()}
        </div>

        <div className="flex-1 min-w-0">
          <div className="flex items-center justify-between gap-2">
            <h4
              className={cn(
                'text-sm truncate',
                unreadCount > 0 ? 'font-semibold text-gray-900' : 'font-medium text-gray-700'
              )}
            >
              {thread.subject}
            </h4>
            {unreadCount > 0 && (
              <Badge className="bg-cyan-500 text-white text-xs px-1.5 py-0.5 min-w-[20px] text-center">
                {unreadCount}
              </Badge>
            )}
          </div>

          {/* Last message preview */}
          {thread.lastMessage && (
            <p className="text-xs text-gray-500 truncate mt-1">
              <span className="font-medium">{thread.lastMessage.senderName}:</span>{' '}
              {thread.lastMessage.contentType === 'system'
                ? thread.lastMessage.content
                : thread.lastMessage.content.length > 50
                ? thread.lastMessage.content.substring(0, 50) + '...'
                : thread.lastMessage.content}
            </p>
          )}

          {/* Meta info */}
          <div className="flex items-center gap-2 mt-1.5">
            <span className="text-xs text-gray-400">{lastMessageTime}</span>
            {!isSchoolContact && thread.schoolName && (
              <>
                <span className="text-gray-300">|</span>
                <span className="text-xs text-gray-400 truncate">{thread.schoolName}</span>
              </>
            )}
            {getStatusBadge()}
          </div>

          {/* Thread type badges */}
          {thread.threadType === 'incident' && thread.incidentId && (
            <div className="flex items-center gap-1 mt-1">
              <Tag className="h-3 w-3 text-gray-400" />
              <span className="text-xs text-gray-400">{thread.incidentId}</span>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
