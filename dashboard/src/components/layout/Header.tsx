import { useState, useRef, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { Search, Menu, Sparkles, MessageSquare, X } from 'lucide-react';
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
  const [mobileSearchOpen, setMobileSearchOpen] = useState(false);
  const mobileSearchRef = useRef<HTMLInputElement>(null);

  // Focus input when mobile search opens
  useEffect(() => {
    if (mobileSearchOpen && mobileSearchRef.current) {
      mobileSearchRef.current.focus();
    }
  }, [mobileSearchOpen]);

  return (
    <header className="fixed top-0 left-0 right-0 z-50 flex h-16 items-center header-gradient px-4 shadow-lg">
      {/* Left section - Logo and menu */}
      <div className="flex items-center gap-4">
        <Button
          variant="ghost"
          size="icon"
          className="lg:hidden text-white/80 hover:text-white hover:bg-white/10"
          onClick={onMenuClick}
          aria-label="Toggle menu"
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

      {/* Center section - Search (Desktop) */}
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

      {/* Mobile Search Button */}
      <div className="md:hidden flex-1 flex justify-end mr-2">
        <Button
          variant="ghost"
          size="icon"
          className="text-white/80 hover:text-white hover:bg-white/10"
          onClick={() => setMobileSearchOpen(true)}
          aria-label="Open search"
        >
          <Search className="h-5 w-5" />
        </Button>
      </div>

      {/* Mobile Search Overlay */}
      {mobileSearchOpen && (
        <div className="fixed inset-x-0 top-0 z-50 bg-gradient-to-r from-slate-800 via-slate-900 to-slate-800 p-4 shadow-lg md:hidden">
          <div className="flex items-center gap-3">
            <Button
              variant="ghost"
              size="icon"
              className="text-white/80 hover:text-white hover:bg-white/10 shrink-0"
              onClick={() => setMobileSearchOpen(false)}
              aria-label="Close search"
            >
              <X className="h-5 w-5" />
            </Button>
            <div className="relative flex-1">
              <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-cyan-300" />
              <Input
                ref={mobileSearchRef}
                type="search"
                placeholder="Search..."
                className="pl-10 bg-white/10 border-white/20 text-white placeholder:text-cyan-200/60 focus:bg-white/15 focus:border-cyan-300"
                onKeyDown={(e) => {
                  if (e.key === 'Escape') {
                    setMobileSearchOpen(false);
                  }
                }}
              />
            </div>
          </div>
        </div>
      )}

      {/* Right section - Messages, Notifications, and user menu */}
      <div className="flex items-center gap-2 ml-auto">
        {/* Message notifications */}
        <Button
          variant="ghost"
          size="icon"
          asChild
          className="relative text-white/80 hover:text-white hover:bg-white/10"
          aria-label={unreadMessages > 0 ? `Messages (${unreadMessages} unread)` : 'Messages'}
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
