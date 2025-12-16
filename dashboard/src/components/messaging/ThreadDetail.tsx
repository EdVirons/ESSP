import { useEffect } from 'react';
import { MessageSquare, Loader2 } from 'lucide-react';
import { ThreadHeader } from './ThreadHeader';
import { MessageList } from './MessageList';
import { MessageComposer } from './MessageComposer';
import { useThread, useMarkRead } from '@/hooks/useMessages';
import { useAuth } from '@/contexts/AuthContext';
import { useMessaging } from '@/contexts/MessagingContext';

interface ThreadDetailProps {
  threadId: string;
  onClose?: () => void;
}

export function ThreadDetail({ threadId, onClose }: ThreadDetailProps) {
  const { user } = useAuth();
  const { typingUsers } = useMessaging();

  const { data, isLoading, error } = useThread(threadId);
  const markRead = useMarkRead(threadId);

  // Mark thread as read when viewing
  useEffect(() => {
    if (data?.messages && data.messages.length > 0) {
      const lastMessage = data.messages[data.messages.length - 1];
      markRead.mutate(lastMessage.id);
    }
  }, [data?.messages?.length]);

  // Get typing users for this thread
  const threadTypingUsers = Array.from(typingUsers.entries())
    .filter(([_, data]) => data.threadId === threadId)
    .map(([userId, data]) => ({ userId, userName: data.userName }));

  if (isLoading) {
    return (
      <div className="flex-1 flex items-center justify-center bg-white">
        <Loader2 className="h-8 w-8 animate-spin text-cyan-600" />
      </div>
    );
  }

  if (error || !data) {
    return (
      <div className="flex-1 flex items-center justify-center bg-white">
        <div className="text-center text-gray-500">
          <MessageSquare className="h-12 w-12 mx-auto mb-3 text-gray-300" />
          <p>Failed to load conversation</p>
        </div>
      </div>
    );
  }

  const { thread, messages, participants } = data;
  const currentUserId = user?.username || '';

  return (
    <div className="flex flex-col h-full bg-white">
      <ThreadHeader
        thread={thread}
        participants={participants}
        onClose={onClose}
      />

      <MessageList
        messages={messages}
        currentUserId={currentUserId}
        typingUsers={threadTypingUsers}
      />

      <MessageComposer
        threadId={threadId}
        disabled={thread.status === 'closed'}
      />
    </div>
  );
}
