import { useState } from 'react';
import { Plus, MessageSquare, Loader2 } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { ThreadListItem } from './ThreadListItem';
import { ThreadFilters } from './ThreadFilters';
import { useThreads } from '@/hooks/useMessages';
import { useAuth } from '@/contexts/AuthContext';

interface ThreadListProps {
  selectedThreadId: string | null;
  onSelectThread: (threadId: string) => void;
  onNewThread: () => void;
}

export function ThreadList({
  selectedThreadId,
  onSelectThread,
  onNewThread,
}: ThreadListProps) {
  const { hasRole } = useAuth();
  const [searchQuery, setSearchQuery] = useState('');
  const [statusFilter, setStatusFilter] = useState('all');

  const isSchoolContact = hasRole('ssp_school_contact');

  const { data, isLoading, error } = useThreads({
    status: statusFilter === 'all' ? undefined : statusFilter,
    q: searchQuery || undefined,
  });

  const threads = data?.items || [];

  return (
    <div className="flex flex-col h-full bg-white">
      {/* Header */}
      <div className="p-4 border-b border-gray-200">
        <div className="flex items-center justify-between">
          <h2 className="text-lg font-semibold text-gray-900">Messages</h2>
          <Button
            size="sm"
            onClick={onNewThread}
            className="bg-cyan-600 hover:bg-cyan-700"
          >
            <Plus className="h-4 w-4 mr-1" />
            New
          </Button>
        </div>
      </div>

      {/* Filters */}
      <ThreadFilters
        searchQuery={searchQuery}
        onSearchChange={setSearchQuery}
        statusFilter={statusFilter}
        onStatusChange={setStatusFilter}
      />

      {/* Thread list */}
      <div className="flex-1 overflow-y-auto">
        {isLoading ? (
          <div className="flex items-center justify-center py-8">
            <Loader2 className="h-6 w-6 animate-spin text-cyan-600" />
          </div>
        ) : error ? (
          <div className="p-4 text-center text-red-500">
            Failed to load conversations
          </div>
        ) : threads.length === 0 ? (
          <div className="p-8 text-center">
            <MessageSquare className="h-12 w-12 mx-auto text-gray-300 mb-3" />
            <p className="text-gray-500">No conversations yet</p>
            <Button
              variant="outline"
              size="sm"
              onClick={onNewThread}
              className="mt-3"
            >
              <Plus className="h-4 w-4 mr-1" />
              Start a conversation
            </Button>
          </div>
        ) : (
          <div>
            {threads.map((thread) => (
              <ThreadListItem
                key={thread.id}
                thread={thread}
                isSelected={thread.id === selectedThreadId}
                onClick={() => onSelectThread(thread.id)}
                isSchoolContact={isSchoolContact}
              />
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
