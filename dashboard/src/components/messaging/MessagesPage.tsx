import { useState } from 'react';
import { MessageSquare } from 'lucide-react';
import { ThreadList } from './ThreadList';
import { ThreadDetail } from './ThreadDetail';
import { NewThreadModal } from './NewThreadModal';
import { MessagingProvider } from '@/contexts/MessagingContext';
import { cn } from '@/lib/utils';

function MessagesContent() {
  const [selectedThreadId, setSelectedThreadId] = useState<string | null>(null);
  const [isNewThreadModalOpen, setIsNewThreadModalOpen] = useState(false);
  const [isMobileThreadOpen, setIsMobileThreadOpen] = useState(false);

  const handleSelectThread = (threadId: string) => {
    setSelectedThreadId(threadId);
    setIsMobileThreadOpen(true);
  };

  const handleCloseThread = () => {
    setIsMobileThreadOpen(false);
  };

  const handleNewThread = () => {
    setIsNewThreadModalOpen(true);
  };

  const handleThreadCreated = (threadId: string) => {
    setSelectedThreadId(threadId);
    setIsMobileThreadOpen(true);
  };

  return (
    <div className="h-[calc(100vh-8rem)] bg-gray-50 rounded-xl overflow-hidden border border-gray-200 shadow-sm">
      <div className="flex h-full">
        {/* Thread list - sidebar */}
        <div
          className={cn(
            'w-full md:w-80 lg:w-96 border-r border-gray-200 flex-shrink-0',
            isMobileThreadOpen && 'hidden md:block'
          )}
        >
          <ThreadList
            selectedThreadId={selectedThreadId}
            onSelectThread={handleSelectThread}
            onNewThread={handleNewThread}
          />
        </div>

        {/* Thread detail - main area */}
        <div
          className={cn(
            'flex-1 flex flex-col',
            !isMobileThreadOpen && 'hidden md:flex'
          )}
        >
          {selectedThreadId ? (
            <ThreadDetail
              threadId={selectedThreadId}
              onClose={handleCloseThread}
            />
          ) : (
            <div className="flex-1 flex items-center justify-center bg-white">
              <div className="text-center text-gray-500">
                <MessageSquare className="h-16 w-16 mx-auto mb-4 text-gray-300" />
                <h3 className="text-lg font-medium text-gray-700 mb-2">
                  Select a conversation
                </h3>
                <p className="text-sm">
                  Choose a conversation from the list to view messages
                </p>
              </div>
            </div>
          )}
        </div>
      </div>

      {/* New thread modal */}
      <NewThreadModal
        isOpen={isNewThreadModalOpen}
        onClose={() => setIsNewThreadModalOpen(false)}
        onThreadCreated={handleThreadCreated}
      />
    </div>
  );
}

export function MessagesPage() {
  return (
    <MessagingProvider>
      <MessagesContent />
    </MessagingProvider>
  );
}
