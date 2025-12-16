import { Users, Bell, AlertCircle } from 'lucide-react';
import { useChatQueue, useAcceptChat } from '@/hooks/useLivechat';
import { useLivechatContext } from '@/contexts/LivechatContext';
import { QueueItem } from './QueueItem';
import { Button } from '@/components/ui/button';

interface ChatQueueProps {
  onChatAccepted?: (sessionId: string) => void;
}

export function ChatQueue({ onChatAccepted }: ChatQueueProps) {
  const { data: queue, isLoading, error } = useChatQueue();
  const acceptChat = useAcceptChat();
  const { newChatNotification, clearNewChatNotification, setActiveSessionId } = useLivechatContext();

  const handleAccept = () => {
    acceptChat.mutate(undefined, {
      onSuccess: (data) => {
        setActiveSessionId(data.session.id);
        onChatAccepted?.(data.session.id);
        clearNewChatNotification();
      },
    });
  };

  if (isLoading) {
    return (
      <div className="p-4">
        <div className="animate-pulse space-y-3">
          <div className="h-4 bg-gray-200 rounded w-1/3" />
          <div className="h-16 bg-gray-200 rounded" />
          <div className="h-16 bg-gray-200 rounded" />
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="p-4">
        <div className="flex items-center gap-2 p-3 rounded-lg bg-red-50 border border-red-200 text-red-700">
          <AlertCircle className="w-5 h-5" />
          <span className="text-sm">Failed to load queue</span>
        </div>
      </div>
    );
  }

  const items = queue?.items ?? [];
  const totalWaiting = queue?.totalWaiting ?? 0;

  return (
    <div className="flex flex-col h-full">
      {/* Header */}
      <div className="flex items-center justify-between p-4 border-b">
        <div className="flex items-center gap-2">
          <Users className="w-5 h-5 text-amber-600" />
          <h3 className="font-semibold text-gray-900">Queue</h3>
          {totalWaiting > 0 && (
            <span className="px-2 py-0.5 text-xs font-medium bg-amber-100 text-amber-700 rounded-full">
              {totalWaiting}
            </span>
          )}
        </div>
      </div>

      {/* New Chat Notification */}
      {newChatNotification && (
        <div className="p-3 mx-3 mt-3 rounded-lg bg-amber-100 border border-amber-300">
          <div className="flex items-start gap-2">
            <Bell className="w-5 h-5 text-amber-600 flex-shrink-0 mt-0.5" />
            <div className="flex-1 min-w-0">
              <p className="text-sm font-medium text-amber-900">
                New chat waiting
              </p>
              <p className="text-sm text-amber-700 truncate">
                {newChatNotification.contactName}
              </p>
            </div>
            <Button
              size="sm"
              variant="ghost"
              onClick={clearNewChatNotification}
              className="text-amber-700 hover:text-amber-900"
            >
              Dismiss
            </Button>
          </div>
        </div>
      )}

      {/* Queue List */}
      <div className="flex-1 overflow-y-auto p-3 space-y-2">
        {items.length === 0 ? (
          <div className="flex flex-col items-center justify-center py-8 text-center">
            <Users className="w-12 h-12 text-gray-300 mb-3" />
            <p className="text-sm text-gray-500">No chats waiting</p>
            <p className="text-xs text-gray-400 mt-1">
              New chats will appear here
            </p>
          </div>
        ) : (
          items.map((item) => (
            <QueueItem
              key={item.sessionId}
              item={item}
              onAccept={handleAccept}
              isAccepting={acceptChat.isPending}
            />
          ))
        )}
      </div>
    </div>
  );
}
