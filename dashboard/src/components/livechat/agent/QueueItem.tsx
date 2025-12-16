import { User, Clock, Bot } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import type { ChatQueueItem } from '@/types/livechat';

interface QueueItemProps {
  item: ChatQueueItem;
  onAccept: () => void;
  isAccepting?: boolean;
}

export function QueueItem({ item, onAccept, isAccepting }: QueueItemProps) {
  const waitMinutes = Math.round(item.waitingTimeSeconds / 60);

  return (
    <div className="flex items-center justify-between p-3 rounded-lg bg-amber-50 border border-amber-200 hover:bg-amber-100 transition-colors">
      <div className="flex items-center gap-3 min-w-0">
        <div className="flex-shrink-0 w-10 h-10 rounded-full bg-amber-200 flex items-center justify-center">
          <User className="w-5 h-5 text-amber-700" />
        </div>
        <div className="min-w-0">
          <p className="font-medium text-gray-900 truncate">
            {item.schoolContactName}
          </p>
          <p className="text-sm text-gray-600 truncate">
            {item.subject || 'New chat'}
          </p>
          <div className="flex items-center gap-2 mt-1">
            <Badge variant="outline" className="text-xs bg-white">
              <Clock className="w-3 h-3 mr-1" />
              {waitMinutes}m waiting
            </Badge>
            <Badge variant="outline" className="text-xs bg-purple-50 text-purple-700 border-purple-200">
              <Bot className="w-3 h-3 mr-1" />
              AI Escalated
            </Badge>
          </div>
        </div>
      </div>
      <Button
        size="sm"
        onClick={onAccept}
        disabled={isAccepting}
        className="flex-shrink-0 bg-amber-600 hover:bg-amber-700"
      >
        {isAccepting ? 'Accepting...' : 'Accept'}
      </Button>
    </div>
  );
}
