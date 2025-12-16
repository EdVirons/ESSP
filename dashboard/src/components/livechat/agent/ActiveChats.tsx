import { MessageSquare, AlertCircle } from 'lucide-react';
import { useActiveChats } from '@/hooks/useLivechat';
import { useLivechatContext } from '@/contexts/LivechatContext';
import { ActiveChatItem } from './ActiveChatItem';

export function ActiveChats() {
  const { data: activeChats, isLoading, error } = useActiveChats();
  const { activeSessionId, setActiveSessionId } = useLivechatContext();

  if (isLoading) {
    return (
      <div className="p-4">
        <div className="animate-pulse space-y-3">
          <div className="h-4 bg-gray-200 rounded w-1/3" />
          <div className="h-20 bg-gray-200 rounded" />
          <div className="h-20 bg-gray-200 rounded" />
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="p-4">
        <div className="flex items-center gap-2 p-3 rounded-lg bg-red-50 border border-red-200 text-red-700">
          <AlertCircle className="w-5 h-5" />
          <span className="text-sm">Failed to load active chats</span>
        </div>
      </div>
    );
  }

  const items = activeChats?.items ?? [];
  const totalActive = activeChats?.total ?? 0;

  return (
    <div className="flex flex-col h-full">
      {/* Header */}
      <div className="flex items-center justify-between p-4 border-b">
        <div className="flex items-center gap-2">
          <MessageSquare className="w-5 h-5 text-cyan-600" />
          <h3 className="font-semibold text-gray-900">Active Chats</h3>
          {totalActive > 0 && (
            <span className="px-2 py-0.5 text-xs font-medium bg-cyan-100 text-cyan-700 rounded-full">
              {totalActive}
            </span>
          )}
        </div>
      </div>

      {/* Active Chats List */}
      <div className="flex-1 overflow-y-auto p-3 space-y-2">
        {items.length === 0 ? (
          <div className="flex flex-col items-center justify-center py-8 text-center">
            <MessageSquare className="w-12 h-12 text-gray-300 mb-3" />
            <p className="text-sm text-gray-500">No active chats</p>
            <p className="text-xs text-gray-400 mt-1">
              Accept a chat from the queue to start
            </p>
          </div>
        ) : (
          items.map((item) => (
            <ActiveChatItem
              key={item.session.id}
              item={item}
              isSelected={activeSessionId === item.session.id}
              onClick={() => setActiveSessionId(item.session.id)}
            />
          ))
        )}
      </div>
    </div>
  );
}
