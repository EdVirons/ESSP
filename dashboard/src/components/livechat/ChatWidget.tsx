import { useState, useCallback, useRef, useEffect } from 'react';
import { MessageCircle, X, Minimize2, Send, Loader2, Bot, User } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Textarea } from '@/components/ui/textarea';
import { useAuth } from '@/contexts/AuthContext';
import { useStartSession, useEndSession, useSendAIMessage, useRequestEscalation } from '@/hooks/useLivechat';
import { useSendMessage, useThread } from '@/hooks/useMessages';
import { useWebSocket } from '@/hooks/useWebSocket';
import type { ChatSession, ChatSessionStatus } from '@/types/livechat';
import type { WSMessage } from '@/types/notification';
import { cn } from '@/lib/utils';
import { formatDistanceToNow } from 'date-fns';

interface ChatWidgetProps {
  className?: string;
}

export function ChatWidget({ className }: ChatWidgetProps) {
  const [isOpen, setIsOpen] = useState(false);
  const [isMinimized, setIsMinimized] = useState(false);
  const [session, setSession] = useState<ChatSession | null>(null);
  const [unreadCount, setUnreadCount] = useState(0);
  const [messageInput, setMessageInput] = useState('');
  const [aiTyping, setAiTyping] = useState(false);
  const messagesEndRef = useRef<HTMLDivElement>(null);

  const { user } = useAuth();
  const startSession = useStartSession();
  const endSession = useEndSession();
  const sendAIMessage = useSendAIMessage();
  const requestEscalation = useRequestEscalation();

  // Fetch thread messages if we have a session
  const { data: threadData, refetch: refetchThread } = useThread(session?.threadId);
  const sendMessage = useSendMessage(session?.threadId || '');

  // Auto-scroll to bottom when messages change
  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [threadData?.messages]);

  // WebSocket for real-time updates
  useWebSocket({
    onMessage: useCallback((wsMessage: WSMessage) => {
      if (wsMessage.type === 'chat_message' && session?.threadId) {
        // Refetch thread to get new messages
        refetchThread();
        // Increment unread if minimized or closed
        if (isMinimized || !isOpen) {
          setUnreadCount((prev) => prev + 1);
        }
      }

      if (wsMessage.type === 'typing_indicator') {
        const payload = wsMessage.payload as {
          threadId: string;
          userId: string;
          isTyping: boolean;
        };
        if (payload.threadId === session?.threadId && payload.userId === 'ai_assistant') {
          setAiTyping(payload.isTyping);
        }
      }

      if (wsMessage.type === 'chat_session_update') {
        const payload = wsMessage.payload as {
          sessionId: string;
          status: string;
          agentId?: string;
          agentName?: string;
        };
        if (payload.sessionId === session?.id) {
          setSession((prev) => {
            if (!prev) return null;
            return {
              ...prev,
              status: payload.status as ChatSessionStatus,
              assignedAgentId: payload.agentId,
              assignedAgentName: payload.agentName,
            };
          });
        }
      }
    }, [session?.id, session?.threadId, refetchThread, isMinimized, isOpen]),
  });

  const handleOpen = () => {
    setIsOpen(true);
    setIsMinimized(false);
    setUnreadCount(0);
  };

  const handleMinimize = () => {
    setIsMinimized(true);
  };

  const handleClose = async () => {
    if (session && session.status !== 'ended') {
      try {
        await endSession.mutateAsync({ sessionId: session.id });
      } catch (error) {
        console.error('Failed to end session:', error);
      }
    }
    setIsOpen(false);
    setIsMinimized(false);
    setSession(null);
    setUnreadCount(0);
  };

  const handleStartChat = async () => {
    if (!user) return;

    try {
      const result = await startSession.mutateAsync({
        subject: 'Live Chat Support',
      });
      setSession(result.session);
    } catch (error) {
      console.error('Failed to start chat session:', error);
    }
  };

  const handleSendMessage = async () => {
    if (!messageInput.trim() || !session?.threadId) return;

    const content = messageInput.trim();
    setMessageInput('');

    try {
      // If in AI mode, send to AI endpoint
      if (session.status === 'ai_active') {
        setAiTyping(true);
        const response = await sendAIMessage.mutateAsync({
          sessionId: session.id,
          content,
        });
        setAiTyping(false);

        // Update session status if changed
        if (response.sessionStatus !== session.status) {
          setSession((prev) => prev ? { ...prev, status: response.sessionStatus } : null);
        }
      } else {
        // Regular message to human agent
        await sendMessage.mutateAsync({ content });
      }
    } catch (error) {
      console.error('Failed to send message:', error);
      setAiTyping(false);
    }
  };

  const handleRequestHuman = async () => {
    if (!session) return;

    try {
      const response = await requestEscalation.mutateAsync({
        sessionId: session.id,
        reason: 'User requested human agent',
      });
      setSession(response.session);
    } catch (error) {
      console.error('Failed to request human agent:', error);
    }
  };

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSendMessage();
    }
  };

  const messages = threadData?.messages || [];
  const isAIMode = session?.status === 'ai_active';

  // Floating button when closed
  if (!isOpen) {
    return (
      <button
        onClick={handleOpen}
        className={cn(
          'fixed bottom-6 right-6 w-14 h-14 rounded-full bg-cyan-600 text-white shadow-lg hover:bg-cyan-700 transition-all hover:scale-105 flex items-center justify-center z-50',
          className
        )}
        aria-label="Open chat"
      >
        <MessageCircle className="h-6 w-6" />
        {unreadCount > 0 && (
          <span className="absolute -top-1 -right-1 w-5 h-5 bg-red-500 rounded-full text-xs flex items-center justify-center">
            {unreadCount > 9 ? '9+' : unreadCount}
          </span>
        )}
      </button>
    );
  }

  // Minimized state
  if (isMinimized) {
    return (
      <div
        className={cn(
          'fixed bottom-6 right-6 w-72 rounded-lg shadow-lg cursor-pointer z-50',
          isAIMode ? 'bg-purple-600' : 'bg-cyan-600',
          'text-white',
          className
        )}
        onClick={() => {
          setIsMinimized(false);
          setUnreadCount(0);
        }}
      >
        <div className="flex items-center justify-between p-3">
          <div className="flex items-center gap-2">
            {isAIMode ? <Bot className="h-5 w-5" /> : <MessageCircle className="h-5 w-5" />}
            <span className="font-medium">
              {isAIMode ? 'AI Assistant' : 'Support Chat'}
            </span>
            {unreadCount > 0 && (
              <span className="w-5 h-5 bg-red-500 rounded-full text-xs flex items-center justify-center">
                {unreadCount > 9 ? '9+' : unreadCount}
              </span>
            )}
          </div>
          <button
            onClick={(e) => {
              e.stopPropagation();
              handleClose();
            }}
            className="hover:bg-white/20 rounded p-1"
          >
            <X className="h-4 w-4" />
          </button>
        </div>
      </div>
    );
  }

  // Full chat panel
  return (
    <div
      className={cn(
        'fixed bottom-6 right-6 w-96 h-[500px] max-h-[80vh] bg-white rounded-lg shadow-2xl flex flex-col z-50 overflow-hidden',
        className
      )}
    >
      {/* Header */}
      <div className={cn(
        'flex items-center justify-between p-4 text-white',
        isAIMode ? 'bg-purple-600' : 'bg-cyan-600'
      )}>
        <div className="flex items-center gap-2">
          {isAIMode ? <Bot className="h-5 w-5" /> : <MessageCircle className="h-5 w-5" />}
          <span className="font-medium">
            {isAIMode ? 'AI Assistant' : 'Support Chat'}
          </span>
        </div>
        <div className="flex items-center gap-1">
          <button
            onClick={handleMinimize}
            className="hover:bg-white/20 rounded p-1"
          >
            <Minimize2 className="h-4 w-4" />
          </button>
          <button
            onClick={handleClose}
            className="hover:bg-white/20 rounded p-1"
          >
            <X className="h-4 w-4" />
          </button>
        </div>
      </div>

      {/* Content */}
      {!session ? (
        // Start chat screen
        <div className="flex-1 flex flex-col items-center justify-center p-6 text-center">
          <Bot className="h-12 w-12 text-purple-600 mb-4" />
          <h3 className="text-lg font-semibold text-gray-900 mb-2">
            Need Help?
          </h3>
          <p className="text-sm text-gray-600 mb-6">
            Our AI assistant is ready to help. You can also request a human agent anytime.
          </p>
          <Button
            onClick={handleStartChat}
            disabled={startSession.isPending}
            className="bg-purple-600 hover:bg-purple-700"
          >
            {startSession.isPending ? (
              <>
                <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                Starting...
              </>
            ) : (
              <>
                <Bot className="h-4 w-4 mr-2" />
                Start Chat
              </>
            )}
          </Button>
        </div>
      ) : (
        <>
          {/* Session status banner */}
          {session.status === 'ai_active' && (
            <div className="px-4 py-2 bg-purple-50 border-b border-purple-100 flex items-center justify-between">
              <div className="flex items-center gap-2">
                <Bot className="h-4 w-4 text-purple-600" />
                <span className="text-sm text-purple-700">AI Assistant</span>
              </div>
              <Button
                variant="outline"
                size="sm"
                onClick={handleRequestHuman}
                disabled={requestEscalation.isPending}
                className="text-xs h-7 border-purple-300 text-purple-700 hover:bg-purple-100"
              >
                {requestEscalation.isPending ? (
                  <Loader2 className="h-3 w-3 animate-spin" />
                ) : (
                  <>
                    <User className="h-3 w-3 mr-1" />
                    Talk to Human
                  </>
                )}
              </Button>
            </div>
          )}
          {session.status === 'waiting' && (
            <div className="px-4 py-2 bg-amber-50 border-b border-amber-100 text-center">
              <p className="text-sm text-amber-700">
                <Loader2 className="h-3 w-3 animate-spin inline mr-2" />
                Connecting you with an agent...
                {session.queuePosition && session.queuePosition > 0 && (
                  <span className="ml-1">(Position: {session.queuePosition})</span>
                )}
              </p>
            </div>
          )}
          {session.status === 'active' && session.assignedAgentName && (
            <div className="px-4 py-2 bg-green-50 border-b border-green-100 text-center">
              <p className="text-sm text-green-700">
                Connected with {session.assignedAgentName}
              </p>
            </div>
          )}

          {/* Messages */}
          <div className="flex-1 overflow-y-auto p-4 space-y-4">
            {messages.length === 0 ? (
              <p className="text-center text-gray-500 text-sm">
                No messages yet. Start the conversation!
              </p>
            ) : (
              messages.map((msg) => {
                const isOwnMessage = msg.senderId === user?.username;
                const isSystem = msg.contentType === 'system';
                const isAI = msg.senderRole === 'ai' || msg.senderId === 'ai_assistant';

                if (isSystem) {
                  return (
                    <div key={msg.id} className="text-center">
                      <span className="text-xs text-gray-500 bg-gray-100 px-2 py-1 rounded">
                        {msg.content}
                      </span>
                    </div>
                  );
                }

                return (
                  <div
                    key={msg.id}
                    className={cn(
                      'flex flex-col max-w-[80%]',
                      isOwnMessage ? 'ml-auto items-end' : 'items-start'
                    )}
                  >
                    <div className="flex items-center gap-1 mb-1">
                      {isAI && <Bot className="h-3 w-3 text-purple-600" />}
                      <span className={cn(
                        'text-xs',
                        isAI ? 'text-purple-600' : 'text-gray-500'
                      )}>
                        {msg.senderName}
                      </span>
                    </div>
                    <div
                      className={cn(
                        'rounded-lg px-3 py-2',
                        isOwnMessage
                          ? 'bg-cyan-600 text-white'
                          : isAI
                          ? 'bg-purple-100 text-gray-900'
                          : 'bg-gray-100 text-gray-900'
                      )}
                    >
                      <p className="text-sm whitespace-pre-wrap">{msg.content}</p>
                    </div>
                    <span className="text-xs text-gray-400 mt-1">
                      {formatDistanceToNow(new Date(msg.createdAt), { addSuffix: true })}
                    </span>
                  </div>
                );
              })
            )}

            {/* AI Typing indicator */}
            {aiTyping && (
              <div className="flex items-start gap-2">
                <div className="flex items-center gap-1">
                  <Bot className="h-3 w-3 text-purple-600" />
                  <span className="text-xs text-purple-600">AI Assistant</span>
                </div>
                <div className="bg-purple-100 rounded-lg px-3 py-2">
                  <div className="flex gap-1">
                    <span className="w-2 h-2 bg-purple-400 rounded-full animate-bounce" style={{ animationDelay: '0ms' }} />
                    <span className="w-2 h-2 bg-purple-400 rounded-full animate-bounce" style={{ animationDelay: '150ms' }} />
                    <span className="w-2 h-2 bg-purple-400 rounded-full animate-bounce" style={{ animationDelay: '300ms' }} />
                  </div>
                </div>
              </div>
            )}

            <div ref={messagesEndRef} />
          </div>

          {/* Input */}
          {session.status !== 'ended' && (
            <div className="p-4 border-t border-gray-200">
              <div className="flex gap-2">
                <Textarea
                  value={messageInput}
                  onChange={(e) => setMessageInput(e.target.value)}
                  onKeyPress={handleKeyPress}
                  placeholder={isAIMode ? "Describe your issue..." : "Type a message..."}
                  className="flex-1 min-h-[40px] max-h-[100px] resize-none"
                  rows={1}
                  disabled={aiTyping}
                />
                <Button
                  onClick={handleSendMessage}
                  disabled={!messageInput.trim() || sendMessage.isPending || sendAIMessage.isPending || aiTyping}
                  size="sm"
                  className={cn(
                    'h-10',
                    isAIMode ? 'bg-purple-600 hover:bg-purple-700' : 'bg-cyan-600 hover:bg-cyan-700'
                  )}
                >
                  {(sendMessage.isPending || sendAIMessage.isPending) ? (
                    <Loader2 className="h-4 w-4 animate-spin" />
                  ) : (
                    <Send className="h-4 w-4" />
                  )}
                </Button>
              </div>
            </div>
          )}
        </>
      )}
    </div>
  );
}
