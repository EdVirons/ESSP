import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { User, Settings, LogOut, ChevronDown, Shield, UserCircle2 } from 'lucide-react';
import { useAuth } from '@/contexts/AuthContext';
import { useImpersonation } from '@/contexts/ImpersonationContext';
import { Button } from '@/components/ui/button';
import { Popover } from '@/components/ui/popover';
import { UserAvatar } from './UserAvatar';
import { ImpersonationModal } from '@/components/ImpersonationModal';

export function UserMenu() {
  const [isOpen, setIsOpen] = useState(false);
  const [impersonationModalOpen, setImpersonationModalOpen] = useState(false);
  const { user, profile, logout, hasPermission } = useAuth();
  const { isImpersonating } = useImpersonation();
  const navigate = useNavigate();

  const canImpersonate = hasPermission('impersonate:user') || hasPermission('*');

  const handleLogout = async () => {
    setIsOpen(false);
    await logout();
    navigate('/login');
  };

  const displayName = profile?.displayName || user?.displayName || user?.username || 'User';
  const email = profile?.email || user?.email;
  const avatarUrl = profile?.avatarUrl || user?.avatarUrl;
  const roles = profile?.roles || user?.roles || [];

  return (
    <>
    <Popover
      open={isOpen}
      onClose={() => setIsOpen(false)}
      align="end"
      trigger={
        <Button
          variant="ghost"
          onClick={() => setIsOpen(!isOpen)}
          className="flex items-center gap-2 px-2"
        >
          <UserAvatar
            src={avatarUrl}
            fallback={displayName}
            size="sm"
          />
          <div className="hidden text-left sm:block">
            <p className="text-sm font-medium text-gray-900 truncate max-w-[120px]">
              {displayName}
            </p>
          </div>
          <ChevronDown className="h-4 w-4 text-gray-500" />
        </Button>
      }
    >
      <div className="min-w-[240px]">
        {/* User info header */}
        <div className="border-b border-gray-100 p-4">
          <div className="flex items-center gap-3">
            <UserAvatar
              src={avatarUrl}
              fallback={displayName}
              size="lg"
            />
            <div className="min-w-0 flex-1">
              <p className="truncate font-medium text-gray-900">
                {displayName}
              </p>
              {email && (
                <p className="truncate text-sm text-gray-500">
                  {email}
                </p>
              )}
              {!email && user?.username && (
                <p className="truncate text-sm text-gray-500">
                  @{user.username}
                </p>
              )}
            </div>
          </div>
          {roles.length > 0 && (
            <div className="mt-3 flex flex-wrap gap-1">
              {roles.map((role) => (
                <span
                  key={role}
                  className="inline-flex items-center gap-1 rounded-full bg-blue-100 px-2 py-0.5 text-xs font-medium text-blue-800"
                >
                  <Shield className="h-3 w-3" />
                  {role.replace('ssp_', '').replace('_', ' ')}
                </span>
              ))}
            </div>
          )}
        </div>

        {/* Menu items */}
        <div className="p-2">
          <Link
            to="/profile"
            onClick={() => setIsOpen(false)}
            className="flex w-full items-center gap-3 rounded-md px-3 py-2 text-sm text-gray-700 hover:bg-gray-100"
          >
            <User className="h-4 w-4" />
            View Profile
          </Link>
          <Link
            to="/settings"
            onClick={() => setIsOpen(false)}
            className="flex w-full items-center gap-3 rounded-md px-3 py-2 text-sm text-gray-700 hover:bg-gray-100"
          >
            <Settings className="h-4 w-4" />
            Settings
          </Link>
          {/* Impersonation option - only for ops managers/admins */}
          {canImpersonate && !isImpersonating && (
            <button
              onClick={() => {
                setIsOpen(false);
                setImpersonationModalOpen(true);
              }}
              className="flex w-full items-center gap-3 rounded-md px-3 py-2 text-sm text-orange-600 hover:bg-orange-50"
            >
              <UserCircle2 className="h-4 w-4" />
              Impersonate User
            </button>
          )}
        </div>

        {/* Logout */}
        <div className="border-t border-gray-100 p-2">
          <button
            onClick={handleLogout}
            className="flex w-full items-center gap-3 rounded-md px-3 py-2 text-sm text-red-600 hover:bg-red-50"
          >
            <LogOut className="h-4 w-4" />
            Sign Out
          </button>
        </div>
      </div>
    </Popover>

    {/* Impersonation modal */}
    <ImpersonationModal
      open={impersonationModalOpen}
      onClose={() => setImpersonationModalOpen(false)}
    />
    </>
  );
}
