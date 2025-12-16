import { Link } from 'react-router-dom';
import { Search, Menu, Sparkles, MessageSquare } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { NotificationCenter } from '@/components/notifications';
import { UserMenu } from '@/components/profile';
import { useUnreadCounts } from '@/hooks/useMessages';

interface HeaderProps {
  onMenuClick: () => void;
}

export function Header({ onMenuClick }: HeaderProps) {
  const { data: messageCounts } = useUnreadCounts();
  const unreadMessages = messageCounts?.total ?? 0;

  return (
    <header className="fixed top-0 left-0 right-0 z-50 flex h-16 items-center header-gradient px-4 shadow-lg">
      {/* Left section - Logo and menu */}
      <div className="flex items-center gap-4">
        <Button
          variant="ghost"
          size="icon"
          className="lg:hidden text-white/80 hover:text-white hover:bg-white/10"
          onClick={onMenuClick}
        >
          <Menu className="h-5 w-5" />
        </Button>

        <Link to="/" className="flex items-center gap-3">
          <div className="flex h-9 w-9 items-center justify-center rounded-lg bg-gradient-to-br from-cyan-400 to-cyan-600 shadow-md">
            <Sparkles className="h-5 w-5 text-white" />
          </div>
          <div className="hidden sm:block">
            <span className="font-bold text-white text-lg tracking-tight">
              ESSP
            </span>
            <span className="text-cyan-200 text-sm ml-2 font-medium">
              Dashboard
            </span>
          </div>
        </Link>
      </div>

      {/* Center section - Search */}
      <div className="mx-4 flex-1 max-w-xl hidden md:block">
        <div className="relative">
          <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-cyan-300" />
          <Input
            type="search"
            placeholder="Search incidents, work orders, schools..."
            className="pl-10 bg-white/10 border-white/20 text-white placeholder:text-cyan-200/60 focus:bg-white/15 focus:border-cyan-300"
          />
        </div>
      </div>

      {/* Right section - Messages, Notifications, and user menu */}
      <div className="flex items-center gap-2 ml-auto">
        {/* Message notifications */}
        <Button
          variant="ghost"
          size="icon"
          asChild
          className="relative text-white/80 hover:text-white hover:bg-white/10"
        >
          <Link to="/messages">
            <MessageSquare className="h-5 w-5" />
            {unreadMessages > 0 && (
              <span className="absolute -top-1 -right-1 flex h-5 w-5 items-center justify-center rounded-full bg-red-500 text-[10px] font-medium text-white">
                {unreadMessages > 99 ? '99+' : unreadMessages}
              </span>
            )}
          </Link>
        </Button>

        {/* System notifications */}
        <NotificationCenter />

        <div className="ml-2 border-l border-white/20 pl-4">
          <UserMenu />
        </div>
      </div>
    </header>
  );
}
