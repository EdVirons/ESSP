import { User, MessageSquare, Bot } from 'lucide-react';
import { Badge } from '@/components/ui/badge';
import { cn } from '@/lib/utils';
import type { ActiveChatItem as ActiveChatItemType } from '@/types/livechat';

interface ActiveChatItemProps {
  item: ActiveChatItemType;
  isSelected: boolean;
  onClick: () => void;
}

export function ActiveChatItem({ item, isSelected, onClick }: ActiveChatItemProps) {
  const { session, unreadCount, lastMessageAt } = item;
  const lastMessageTime = lastMessageAt
    ? new Date(lastMessageAt).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
    : '';

  return (
    <button
      onClick={onClick}
      className={cn(
        'w-full text-left p-3 rounded-lg border transition-colors',
        isSelected
          ? 'bg-cyan-50 border-cyan-300 ring-1 ring-cyan-300'
          : 'bg-white border-gray-200 hover:bg-gray-50'
      )}
    >
      <div className="flex items-start gap-3">
        <div className={cn(
          'flex-shrink-0 w-10 h-10 rounded-full flex items-center justify-center',
          isSelected ? 'bg-cyan-200' : 'bg-gray-200'
        )}>
          <User className={cn(
            'w-5 h-5',
            isSelected ? 'text-cyan-700' : 'text-gray-600'
          )} />
        </div>
        <div className="flex-1 min-w-0">
          <div className="flex items-center justify-between">
            <p className={cn(
              'font-medium truncate',
              isSelected ? 'text-cyan-900' : 'text-gray-900'
            )}>
              {session.schoolContactName}
            </p>
            {unreadCount > 0 && (
              <Badge className="bg-red-500 text-white text-xs ml-2">
                {unreadCount}
              </Badge>
            )}
          </div>
          <div className="flex items-center gap-2 mt-1">
            {session.aiHandled && !session.aiResolved && (
              <Badge variant="outline" className="text-xs bg-purple-50 text-purple-700 border-purple-200">
                <Bot className="w-3 h-3 mr-1" />
                AI Escalated
              </Badge>
            )}
            {session.issueCategory && (
              <Badge variant="outline" className="text-xs">
                {session.issueCategory}
              </Badge>
            )}
          </div>
          <div className="flex items-center justify-between mt-2">
            <div className="flex items-center gap-1 text-xs text-gray-500">
              <MessageSquare className="w-3 h-3" />
              <span>{session.totalMessages} messages</span>
            </div>
            {lastMessageTime && (
              <span className="text-xs text-gray-400">{lastMessageTime}</span>
            )}
          </div>
        </div>
      </div>
    </button>
  );
}
