import { User, School, ArrowRightLeft, X, MoreVertical } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import type { ChatSession } from '@/types/livechat';

interface ChatWindowHeaderProps {
  session: ChatSession;
  onTransfer: () => void;
  onEndChat: () => void;
}

export function ChatWindowHeader({ session, onTransfer, onEndChat }: ChatWindowHeaderProps) {
  const statusColor = {
    ai_active: 'bg-purple-100 text-purple-700 border-purple-200',
    waiting: 'bg-amber-100 text-amber-700 border-amber-200',
    active: 'bg-green-100 text-green-700 border-green-200',
    ended: 'bg-gray-100 text-gray-700 border-gray-200',
  }[session.status];

  const startedAt = new Date(session.startedAt).toLocaleTimeString([], {
    hour: '2-digit',
    minute: '2-digit',
  });

  return (
    <div className="flex items-center justify-between px-4 py-3 border-b bg-white">
      <div className="flex items-center gap-3 min-w-0">
        <div className="flex-shrink-0 w-10 h-10 rounded-full bg-cyan-100 flex items-center justify-center">
          <User className="w-5 h-5 text-cyan-600" />
        </div>
        <div className="min-w-0">
          <div className="flex items-center gap-2">
            <h3 className="font-semibold text-gray-900 truncate">
              {session.schoolContactName}
            </h3>
            <Badge variant="outline" className={statusColor}>
              {session.status.replace('_', ' ')}
            </Badge>
          </div>
          <div className="flex items-center gap-2 text-sm text-gray-500">
            <School className="w-4 h-4" />
            <span className="truncate">School #{session.schoolId}</span>
            <span className="text-gray-300">|</span>
            <span>Started {startedAt}</span>
          </div>
        </div>
      </div>

      <div className="flex items-center gap-2">
        <Button
          variant="outline"
          size="sm"
          onClick={onTransfer}
          className="hidden sm:flex"
        >
          <ArrowRightLeft className="w-4 h-4 mr-1" />
          Transfer
        </Button>
        <Button
          variant="outline"
          size="sm"
          onClick={onEndChat}
          className="text-red-600 hover:text-red-700 hover:bg-red-50 hidden sm:flex"
        >
          <X className="w-4 h-4 mr-1" />
          End Chat
        </Button>

        {/* Mobile dropdown */}
        <DropdownMenu>
          <DropdownMenuTrigger asChild className="sm:hidden">
            <Button variant="ghost" size="icon">
              <MoreVertical className="w-5 h-5" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            <DropdownMenuItem onClick={onTransfer}>
              <ArrowRightLeft className="w-4 h-4 mr-2" />
              Transfer Chat
            </DropdownMenuItem>
            <DropdownMenuItem onClick={onEndChat} className="text-red-600">
              <X className="w-4 h-4 mr-2" />
              End Chat
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    </div>
  );
}
