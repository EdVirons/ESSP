import { useEffect, useRef } from 'react';
import { User, Bot } from 'lucide-react';
import { cn } from '@/lib/utils';
import type { Message } from '@/types/messaging';
import { useAuth } from '@/contexts/AuthContext';

interface AgentMessageListProps {
  messages: Message[];
  typingUser?: string;
}

export function AgentMessageList({ messages, typingUser }: AgentMessageListProps) {
  const { user } = useAuth();
  const bottomRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages, typingUser]);

  const formatTime = (dateStr: string) => {
    return new Date(dateStr).toLocaleTimeString([], {
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  const isOwnMessage = (msg: Message) => msg.senderId === user?.username;
  const isAIMessage = (msg: Message) => msg.senderRole === 'ai_assistant' || msg.senderName === 'AI Assistant';
  const isSystemMessage = (msg: Message) => msg.contentType === 'system';

  return (
    <div className="flex-1 overflow-y-auto p-4 space-y-4">
      {messages.map((msg) => {
        if (isSystemMessage(msg)) {
          return (
            <div key={msg.id} className="flex justify-center">
              <span className="px-3 py-1 text-xs text-gray-500 bg-gray-100 rounded-full">
                {msg.content}
              </span>
            </div>
          );
        }

        const isOwn = isOwnMessage(msg);
        const isAI = isAIMessage(msg);

        return (
          <div
            key={msg.id}
            className={cn('flex items-end gap-2', isOwn ? 'flex-row-reverse' : 'flex-row')}
          >
            {/* Avatar */}
            <div
              className={cn(
                'flex-shrink-0 w-8 h-8 rounded-full flex items-center justify-center',
                isOwn
                  ? 'bg-cyan-100'
                  : isAI
                  ? 'bg-purple-100'
                  : 'bg-gray-100'
              )}
            >
              {isAI ? (
                <Bot className="w-4 h-4 text-purple-600" />
              ) : (
                <User
                  className={cn('w-4 h-4', isOwn ? 'text-cyan-600' : 'text-gray-600')}
                />
              )}
            </div>

            {/* Message Bubble */}
            <div
              className={cn(
                'max-w-[70%] rounded-2xl px-4 py-2',
                isOwn
                  ? 'bg-cyan-600 text-white rounded-br-md'
                  : isAI
                  ? 'bg-purple-100 text-purple-900 rounded-bl-md'
                  : 'bg-gray-100 text-gray-900 rounded-bl-md'
              )}
            >
              {/* Sender Name */}
              {!isOwn && (
                <p
                  className={cn(
                    'text-xs font-medium mb-1',
                    isAI ? 'text-purple-700' : 'text-gray-600'
                  )}
                >
                  {msg.senderName}
                </p>
              )}

              {/* Content */}
              <p className="text-sm whitespace-pre-wrap break-words">{msg.content}</p>

              {/* Timestamp */}
              <p
                className={cn(
                  'text-xs mt-1',
                  isOwn ? 'text-cyan-100' : isAI ? 'text-purple-500' : 'text-gray-400'
                )}
              >
                {formatTime(msg.createdAt)}
              </p>
            </div>
          </div>
        );
      })}

      {/* Typing Indicator */}
      {typingUser && (
        <div className="flex items-end gap-2">
          <div className="w-8 h-8 rounded-full bg-gray-100 flex items-center justify-center">
            <User className="w-4 h-4 text-gray-600" />
          </div>
          <div className="bg-gray-100 rounded-2xl rounded-bl-md px-4 py-3">
            <div className="flex items-center gap-1">
              <span className="text-xs text-gray-500">{typingUser} is typing</span>
              <div className="flex gap-1 ml-2">
                <span className="w-2 h-2 bg-gray-400 rounded-full animate-bounce" />
                <span
                  className="w-2 h-2 bg-gray-400 rounded-full animate-bounce"
                  style={{ animationDelay: '0.1s' }}
                />
                <span
                  className="w-2 h-2 bg-gray-400 rounded-full animate-bounce"
                  style={{ animationDelay: '0.2s' }}
                />
              </div>
            </div>
          </div>
        </div>
      )}

      <div ref={bottomRef} />
    </div>
  );
}
