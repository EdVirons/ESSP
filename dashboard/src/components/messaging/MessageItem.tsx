import { format } from 'date-fns';
import { Paperclip, Download, CheckCheck } from 'lucide-react';
import { cn } from '@/lib/utils';
import type { Message } from '@/types/messaging';

interface MessageItemProps {
  message: Message;
  isOwnMessage: boolean;
  showAvatar?: boolean;
}

export function MessageItem({ message, isOwnMessage, showAvatar = true }: MessageItemProps) {
  // System messages
  if (message.contentType === 'system') {
    return (
      <div className="flex justify-center my-3">
        <div className="bg-gray-100 text-gray-600 text-xs px-3 py-1.5 rounded-full">
          {message.content}
        </div>
      </div>
    );
  }

  const formattedTime = format(new Date(message.createdAt), 'HH:mm');

  const getRoleBadgeColor = (role: string) => {
    switch (role) {
      case 'ssp_admin':
        return 'bg-purple-100 text-purple-700';
      case 'ssp_support_agent':
        return 'bg-blue-100 text-blue-700';
      case 'ssp_school_contact':
        return 'bg-green-100 text-green-700';
      default:
        return 'bg-gray-100 text-gray-700';
    }
  };

  const getRoleLabel = (role: string) => {
    switch (role) {
      case 'ssp_admin':
        return 'Admin';
      case 'ssp_support_agent':
        return 'Support';
      case 'ssp_school_contact':
        return 'School';
      default:
        return role;
    }
  };

  return (
    <div
      className={cn(
        'flex gap-3 mb-4',
        isOwnMessage ? 'flex-row-reverse' : 'flex-row'
      )}
    >
      {/* Avatar */}
      {showAvatar && (
        <div
          className={cn(
            'flex-shrink-0 w-8 h-8 rounded-full flex items-center justify-center text-sm font-medium',
            isOwnMessage ? 'bg-cyan-100 text-cyan-700' : 'bg-gray-100 text-gray-700'
          )}
        >
          {message.senderName.charAt(0).toUpperCase()}
        </div>
      )}

      {/* Message content */}
      <div
        className={cn(
          'flex flex-col max-w-[70%]',
          isOwnMessage ? 'items-end' : 'items-start'
        )}
      >
        {/* Sender info */}
        <div className="flex items-center gap-2 mb-1">
          <span className="text-xs font-medium text-gray-700">
            {message.senderName}
          </span>
          <span
            className={cn(
              'text-xs px-1.5 py-0.5 rounded',
              getRoleBadgeColor(message.senderRole)
            )}
          >
            {getRoleLabel(message.senderRole)}
          </span>
        </div>

        {/* Message bubble */}
        <div
          className={cn(
            'rounded-2xl px-4 py-2.5',
            isOwnMessage
              ? 'bg-cyan-600 text-white rounded-br-md'
              : 'bg-gray-100 text-gray-900 rounded-bl-md'
          )}
        >
          <p className="text-sm whitespace-pre-wrap break-words">{message.content}</p>

          {/* Attachments */}
          {message.attachments && message.attachments.length > 0 && (
            <div className="mt-2 space-y-1">
              {message.attachments.map((attachment) => (
                <a
                  key={attachment.id}
                  href={attachment.downloadUrl}
                  target="_blank"
                  rel="noopener noreferrer"
                  className={cn(
                    'flex items-center gap-2 text-xs p-2 rounded',
                    isOwnMessage
                      ? 'bg-cyan-500 hover:bg-cyan-400'
                      : 'bg-gray-200 hover:bg-gray-300'
                  )}
                >
                  <Paperclip className="h-3.5 w-3.5" />
                  <span className="truncate flex-1">{attachment.fileName}</span>
                  <Download className="h-3.5 w-3.5" />
                </a>
              ))}
            </div>
          )}
        </div>

        {/* Time and status */}
        <div className="flex items-center gap-1 mt-1">
          <span className="text-xs text-gray-400">{formattedTime}</span>
          {message.editedAt && (
            <span className="text-xs text-gray-400">(edited)</span>
          )}
          {isOwnMessage && (
            <CheckCheck className="h-3.5 w-3.5 text-gray-400" />
          )}
        </div>
      </div>
    </div>
  );
}
