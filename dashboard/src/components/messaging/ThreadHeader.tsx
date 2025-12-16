import { useState, useRef, useEffect } from 'react';
import { MoreVertical, X, Archive, RefreshCw, Users, AlertTriangle } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { useCloseThread, useReopenThread } from '@/hooks/useMessages';
import type { MessageThread, ThreadParticipant } from '@/types/messaging';

interface ThreadHeaderProps {
  thread: MessageThread;
  participants?: ThreadParticipant[];
  onClose?: () => void;
}

export function ThreadHeader({ thread, participants, onClose }: ThreadHeaderProps) {
  const [menuOpen, setMenuOpen] = useState(false);
  const menuRef = useRef<HTMLDivElement>(null);
  const closeThread = useCloseThread();
  const reopenThread = useReopenThread();

  // Close menu when clicking outside
  useEffect(() => {
    const handleClickOutside = (e: MouseEvent) => {
      if (menuRef.current && !menuRef.current.contains(e.target as Node)) {
        setMenuOpen(false);
      }
    };
    if (menuOpen) {
      document.addEventListener('mousedown', handleClickOutside);
    }
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, [menuOpen]);

  const handleCloseThread = () => {
    closeThread.mutate(thread.id);
    setMenuOpen(false);
  };

  const handleReopenThread = () => {
    reopenThread.mutate(thread.id);
    setMenuOpen(false);
  };

  const getThreadTypeBadge = () => {
    switch (thread.threadType) {
      case 'incident':
        return (
          <Badge variant="outline" className="text-amber-600 border-amber-200 bg-amber-50">
            <AlertTriangle className="h-3 w-3 mr-1" />
            Incident
          </Badge>
        );
      case 'livechat':
        return (
          <Badge variant="outline" className="text-green-600 border-green-200 bg-green-50">
            Live Chat
          </Badge>
        );
      default:
        return null;
    }
  };

  return (
    <div className="flex items-center justify-between p-4 border-b border-gray-200 bg-white">
      <div className="flex-1 min-w-0">
        <div className="flex items-center gap-2">
          <h2 className="text-lg font-semibold text-gray-900 truncate">
            {thread.subject}
          </h2>
          {getThreadTypeBadge()}
          {thread.status === 'closed' && (
            <Badge variant="secondary">Closed</Badge>
          )}
        </div>

        <div className="flex items-center gap-3 mt-1 text-sm text-gray-500">
          {thread.schoolName && (
            <span>{thread.schoolName}</span>
          )}
          {participants && participants.length > 0 && (
            <span className="flex items-center gap-1">
              <Users className="h-3.5 w-3.5" />
              {participants.length} participant{participants.length !== 1 ? 's' : ''}
            </span>
          )}
          <span>{thread.messageCount} messages</span>
        </div>
      </div>

      <div className="flex items-center gap-2">
        {/* Simple dropdown menu */}
        <div className="relative" ref={menuRef}>
          <Button
            variant="ghost"
            size="sm"
            onClick={() => setMenuOpen(!menuOpen)}
            className="h-8 w-8 p-0"
          >
            <MoreVertical className="h-4 w-4" />
          </Button>

          {menuOpen && (
            <div className="absolute right-0 top-full mt-1 w-48 bg-white rounded-md shadow-lg border border-gray-200 py-1 z-50">
              {thread.status === 'open' ? (
                <button
                  onClick={handleCloseThread}
                  className="flex items-center w-full px-3 py-2 text-sm text-gray-700 hover:bg-gray-100"
                >
                  <Archive className="h-4 w-4 mr-2" />
                  Close conversation
                </button>
              ) : (
                <button
                  onClick={handleReopenThread}
                  className="flex items-center w-full px-3 py-2 text-sm text-gray-700 hover:bg-gray-100"
                >
                  <RefreshCw className="h-4 w-4 mr-2" />
                  Reopen conversation
                </button>
              )}
              {thread.incidentId && (
                <>
                  <div className="border-t border-gray-100 my-1" />
                  <button className="flex items-center w-full px-3 py-2 text-sm text-gray-700 hover:bg-gray-100">
                    <AlertTriangle className="h-4 w-4 mr-2" />
                    View linked incident
                  </button>
                </>
              )}
            </div>
          )}
        </div>

        {onClose && (
          <Button variant="ghost" size="sm" onClick={onClose} className="md:hidden h-8 w-8 p-0">
            <X className="h-4 w-4" />
          </Button>
        )}
      </div>
    </div>
  );
}
