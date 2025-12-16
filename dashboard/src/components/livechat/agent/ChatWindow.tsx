import { useState, useMemo } from 'react';
import { MessageSquare } from 'lucide-react';
import { useActiveChats } from '@/hooks/useLivechat';
import { useLivechatContext } from '@/contexts/LivechatContext';
import { useThread, useSendMessage } from '@/hooks/useMessages';
import { AIContextBanner } from './AIContextBanner';
import { ChatWindowHeader } from './ChatWindowHeader';
import { AgentMessageList } from './AgentMessageList';
import { AgentMessageComposer } from './AgentMessageComposer';
import { ChatTransferModal } from './ChatTransferModal';
import { EndChatModal } from './EndChatModal';

export function ChatWindow() {
  const { activeSessionId, typingInSession } = useLivechatContext();
  const { data: activeChats } = useActiveChats();
  const [transferOpen, setTransferOpen] = useState(false);
  const [endChatOpen, setEndChatOpen] = useState(false);

  // Find the active session
  const activeChat = useMemo(() => {
    return activeChats?.items?.find((item) => item.session.id === activeSessionId);
  }, [activeChats, activeSessionId]);

  const session = activeChat?.session;
  const threadId = session?.threadId;

  // Fetch thread messages
  const { data: threadData, isLoading: loadingThread } = useThread(threadId);
  const sendMessage = useSendMessage(threadId || '');

  // Get typing user for this thread
  const typingUser = useMemo(() => {
    if (!threadId) return undefined;
    const typing = typingInSession.get(threadId);
    return typing?.userName;
  }, [threadId, typingInSession]);

  const handleSendMessage = (content: string) => {
    if (!threadId) return;
    sendMessage.mutate({ content });
  };

  // Empty state - no chat selected
  if (!activeSessionId || !session) {
    return (
      <div className="flex flex-col items-center justify-center h-full bg-gray-50">
        <MessageSquare className="w-16 h-16 text-gray-300 mb-4" />
        <h3 className="text-lg font-medium text-gray-600 mb-2">No chat selected</h3>
        <p className="text-sm text-gray-500 text-center max-w-xs">
          Select an active chat from the sidebar or accept a new chat from the queue
        </p>
      </div>
    );
  }

  // Loading state
  if (loadingThread) {
    return (
      <div className="flex flex-col h-full bg-white">
        <div className="p-4 border-b animate-pulse">
          <div className="h-10 bg-gray-200 rounded w-1/3" />
        </div>
        <div className="flex-1 p-4 space-y-4">
          <div className="h-16 bg-gray-100 rounded w-2/3" />
          <div className="h-16 bg-gray-100 rounded w-1/2 ml-auto" />
          <div className="h-16 bg-gray-100 rounded w-2/3" />
        </div>
      </div>
    );
  }

  const messages = threadData?.messages ?? [];

  return (
    <div className="flex flex-col h-full bg-white">
      {/* Header */}
      <ChatWindowHeader
        session={session}
        onTransfer={() => setTransferOpen(true)}
        onEndChat={() => setEndChatOpen(true)}
      />

      {/* AI Context Banner - Show if escalated from AI */}
      {session.aiHandled && !session.aiResolved && (
        <AIContextBanner sessionId={session.id} />
      )}

      {/* Messages */}
      <AgentMessageList messages={messages} typingUser={typingUser} />

      {/* Composer */}
      <AgentMessageComposer
        onSend={handleSendMessage}
        disabled={sendMessage.isPending || session.status === 'ended'}
        placeholder={
          session.status === 'ended'
            ? 'This chat has ended'
            : 'Type a message...'
        }
      />

      {/* Modals */}
      <ChatTransferModal
        open={transferOpen}
        onOpenChange={setTransferOpen}
        sessionId={session.id}
      />
      <EndChatModal
        open={endChatOpen}
        onOpenChange={setEndChatOpen}
        sessionId={session.id}
        contactName={session.schoolContactName}
      />
    </div>
  );
}
