import { useState, useRef, useCallback } from 'react';
import { Send, Paperclip, X, Loader2 } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Textarea } from '@/components/ui/textarea';
import { useSendMessage } from '@/hooks/useMessages';
import { messagingApi } from '@/api/messaging';
import { cn } from '@/lib/utils';

interface MessageComposerProps {
  threadId: string;
  disabled?: boolean;
  placeholder?: string;
  onTyping?: () => void;
}

export function MessageComposer({
  threadId,
  disabled,
  placeholder = 'Type a message...',
  onTyping,
}: MessageComposerProps) {
  const [content, setContent] = useState('');
  const [attachments, setAttachments] = useState<File[]>([]);
  const textareaRef = useRef<HTMLTextAreaElement>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const sendMessage = useSendMessage(threadId);

  const [isUploading, setIsUploading] = useState(false);

  const uploadAttachments = async (files: File[]): Promise<string[]> => {
    const attachmentRefs: string[] = [];

    for (const file of files) {
      // Get upload URL from the server
      const { uploadUrl, attachmentRef } = await messagingApi.getUploadUrl({
        fileName: file.name,
        contentType: file.type || 'application/octet-stream',
        sizeBytes: file.size,
      });

      // Upload file to the pre-signed URL
      await fetch(uploadUrl, {
        method: 'PUT',
        body: file,
        headers: {
          'Content-Type': file.type || 'application/octet-stream',
        },
      });

      attachmentRefs.push(attachmentRef);
    }

    return attachmentRefs;
  };

  const handleSubmit = useCallback(async () => {
    if (!content.trim() && attachments.length === 0) return;
    if (sendMessage.isPending || isUploading) return;

    try {
      setIsUploading(true);

      // Upload attachments if any
      let attachmentRefs: string[] = [];
      if (attachments.length > 0) {
        attachmentRefs = await uploadAttachments(attachments);
      }

      // Send message with attachment references
      await sendMessage.mutateAsync({
        content: content.trim(),
        attachments: attachmentRefs.length > 0 ? attachmentRefs : undefined,
      });

      setContent('');
      setAttachments([]);

      // Focus textarea after sending
      textareaRef.current?.focus();
    } catch (error) {
      console.error('Failed to send message:', error);
    } finally {
      setIsUploading(false);
    }
  }, [content, attachments, sendMessage, isUploading]);

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSubmit();
    }
  };

  const handleInputChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    setContent(e.target.value);
    onTyping?.();

    // Auto-resize textarea
    if (textareaRef.current) {
      textareaRef.current.style.height = 'auto';
      textareaRef.current.style.height = Math.min(textareaRef.current.scrollHeight, 150) + 'px';
    }
  };

  const handleFileSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
    const files = Array.from(e.target.files || []);
    setAttachments((prev) => [...prev, ...files]);
  };

  const removeAttachment = (index: number) => {
    setAttachments((prev) => prev.filter((_, i) => i !== index));
  };

  const canSend = (content.trim() || attachments.length > 0) && !disabled && !isUploading;

  return (
    <div className="border-t border-gray-200 bg-white p-4">
      {/* Attachments preview */}
      {attachments.length > 0 && (
        <div className="flex flex-wrap gap-2 mb-3">
          {attachments.map((file, index) => (
            <div
              key={index}
              className="flex items-center gap-2 bg-gray-100 rounded-lg px-3 py-1.5 text-sm"
            >
              <Paperclip className="h-4 w-4 text-gray-500" />
              <span className="truncate max-w-[150px]">{file.name}</span>
              <button
                onClick={() => removeAttachment(index)}
                className="text-gray-400 hover:text-gray-600"
              >
                <X className="h-4 w-4" />
              </button>
            </div>
          ))}
        </div>
      )}

      {/* Input area */}
      <div className="flex items-end gap-2">
        {/* Attachment button */}
        <Button
          type="button"
          variant="ghost"
          size="icon"
          onClick={() => fileInputRef.current?.click()}
          disabled={disabled || isUploading}
          className="flex-shrink-0"
          aria-label="Add attachment"
        >
          <Paperclip className="h-5 w-5 text-gray-500" />
        </Button>
        <input
          ref={fileInputRef}
          type="file"
          multiple
          className="hidden"
          onChange={handleFileSelect}
          accept="image/*,.pdf,.doc,.docx,.xls,.xlsx,.txt"
        />

        {/* Message input */}
        <div className="flex-1 relative">
          <Textarea
            ref={textareaRef}
            value={content}
            onChange={handleInputChange}
            onKeyDown={handleKeyDown}
            placeholder={placeholder}
            disabled={disabled}
            className={cn(
              'min-h-[44px] max-h-[150px] resize-none pr-4',
              disabled && 'bg-gray-50'
            )}
            rows={1}
          />
        </div>

        {/* Send button */}
        <Button
          onClick={handleSubmit}
          disabled={!canSend || sendMessage.isPending || isUploading}
          className="flex-shrink-0 bg-cyan-600 hover:bg-cyan-700"
          aria-label={isUploading ? 'Uploading attachments' : sendMessage.isPending ? 'Sending message' : 'Send message'}
        >
          {sendMessage.isPending || isUploading ? (
            <Loader2 className="h-5 w-5 animate-spin" />
          ) : (
            <Send className="h-5 w-5" />
          )}
        </Button>
      </div>

      {/* Disabled message */}
      {disabled && (
        <p className="text-xs text-gray-500 mt-2">
          This conversation is closed. Reopen it to send messages.
        </p>
      )}
    </div>
  );
}
