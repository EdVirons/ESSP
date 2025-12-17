import { useState } from 'react';
import { Outlet } from 'react-router-dom';
import { Header } from './Header';
import { Sidebar } from './Sidebar';
import { ImpersonationBanner } from '@/components/ImpersonationBanner';
import { cn } from '@/lib/utils';
import { useNotificationContext } from '@/contexts/NotificationContext';
import { useImpersonation } from '@/contexts/ImpersonationContext';

export function Layout() {
  const [sidebarCollapsed, setSidebarCollapsed] = useState(false);
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);
  const { isConnected } = useNotificationContext();
  const { isImpersonating } = useImpersonation();

  const toggleSidebar = () => {
    setSidebarCollapsed(!sidebarCollapsed);
  };

  const toggleMobileMenu = () => {
    setMobileMenuOpen(!mobileMenuOpen);
  };

  return (
    <div className={cn("min-h-screen bg-gray-50", isImpersonating && "border-t-4 border-orange-500")}>
      {/* Impersonation banner - shown when ops manager is acting as school contact */}
      <ImpersonationBanner />

      <Header onMenuClick={toggleMobileMenu} />

      {/* Mobile sidebar overlay */}
      {mobileMenuOpen && (
        <div
          className="fixed inset-0 z-40 bg-black/50 lg:hidden"
          onClick={() => setMobileMenuOpen(false)}
        />
      )}

      {/* Sidebar - hidden on mobile unless menu open */}
      <div
        className={cn(
          'lg:block',
          mobileMenuOpen ? 'block' : 'hidden'
        )}
      >
        <Sidebar collapsed={sidebarCollapsed} onToggle={toggleSidebar} />
      </div>

      {/* Main content */}
      <main
        className={cn(
          'pt-16 transition-all duration-300',
          sidebarCollapsed ? 'lg:pl-16' : 'lg:pl-64'
        )}
      >
        <div className="p-6 mx-auto max-w-[1600px]">
          <Outlet />
        </div>
      </main>

      {/* Status bar */}
      <footer
        className={cn(
          'fixed bottom-0 left-0 right-0 z-30 flex h-8 items-center justify-between border-t border-gray-200 bg-white px-4 text-xs text-gray-500 transition-all duration-300',
          sidebarCollapsed ? 'lg:pl-16' : 'lg:pl-64'
        )}
      >
        <div className="flex items-center gap-4">
          <span className="flex items-center gap-1">
            <span
              className={cn(
                'h-2 w-2 rounded-full',
                isConnected ? 'bg-green-500' : 'bg-yellow-500 animate-pulse'
              )}
            />
            {isConnected ? 'Connected' : 'Reconnecting...'}
          </span>
          <span>Real-time updates {isConnected ? 'enabled' : 'paused'}</span>
        </div>
        <div>
          <span>ESSP Dashboard v1.0.0</span>
        </div>
      </footer>
    </div>
  );
}
