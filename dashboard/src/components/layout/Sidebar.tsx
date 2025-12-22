import { useMemo, useState, useCallback } from 'react';
import { NavLink, useLocation } from 'react-router-dom';
import { ChevronLeft, ChevronRight, ChevronDown } from 'lucide-react';
import { cn } from '@/lib/utils';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { useAuth } from '@/contexts/AuthContext';
import { navGroups, type NavItem, type NavGroup } from '@/config/navigation';
import { useUnreadCounts } from '@/hooks/useMessages';

interface SidebarProps {
  collapsed: boolean;
  onToggle: () => void;
}

export function Sidebar({ collapsed, onToggle }: SidebarProps) {
  const { hasPermission, hasRole } = useAuth();
  const { data: unreadCounts } = useUnreadCounts();
  const location = useLocation();
  const [expandedGroups, setExpandedGroups] = useState<Set<string>>(new Set(['support']));

  // Filter nav item visibility
  const isItemVisible = useCallback((item: NavItem): boolean => {
    // If item is restricted to school_contact only, hide from admin
    if (item.roles?.length === 1 && item.roles[0] === 'ssp_school_contact') {
      return hasRole('ssp_school_contact');
    }
    if (hasRole('ssp_admin')) return true;
    if (!item.permissions && !item.roles) return true;
    if (item.permissions?.length && item.permissions.some((p) => hasPermission(p))) return true;
    if (item.roles?.length && item.roles.some((r) => hasRole(r))) return true;
    return false;
  }, [hasPermission, hasRole]);

  // Filter groups and their items based on user permissions and roles
  const visibleGroups = useMemo(() => {
    return navGroups
      .filter((group) => {
        // Hide groups that are only for school_contact from admins
        if (group.roles?.length === 1 && group.roles[0] === 'ssp_school_contact') {
          return hasRole('ssp_school_contact');
        }
        return true;
      })
      .map((group) => ({
        ...group,
        items: group.items.filter(isItemVisible),
      }))
      .filter((group) => group.items.length > 0);
  }, [hasRole, isItemVisible]);

  // Check if any item in group is active
  const isGroupActive = (group: NavGroup): boolean => {
    return group.items.some((item) => {
      // Extract pathname from href (remove query params)
      const hrefPath = item.href.split('?')[0];
      return location.pathname.startsWith(hrefPath);
    });
  };

  // Check if a specific nav item is active (considering query params)
  const isNavItemActive = (item: NavItem): boolean => {
    const hrefPath = item.href.split('?')[0];
    const hrefSearch = item.href.includes('?') ? item.href.split('?')[1] : '';

    // Pathname must match
    if (location.pathname !== hrefPath) return false;

    // If href has action=create, never show as active (it's a trigger, not a page state)
    if (hrefSearch.includes('action=create')) {
      return false;
    }

    // If href has no query params, it's active when pathname matches and no status param
    if (!hrefSearch) {
      return !location.search.includes('status=');
    }

    // If href has query params, check if they match
    return location.search.includes(hrefSearch);
  };

  const toggleGroup = (groupId: string) => {
    setExpandedGroups((prev) => {
      if (prev.has(groupId)) {
        // Clicking expanded group collapses it
        return new Set();
      } else {
        // Clicking collapsed group expands it (and closes others)
        return new Set([groupId]);
      }
    });
  };

  const renderNavItem = (item: NavItem) => {
    const unreadCount = item.href === '/messages' ? unreadCounts?.total : undefined;
    const isActive = isNavItemActive(item);

    return (
      <NavLink
        key={item.href}
        to={item.href}
        className={cn(
          'flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium transition-all duration-200',
          isActive
            ? 'bg-gradient-to-r from-cyan-50 to-teal-50 text-cyan-700 shadow-sm border border-cyan-100'
            : 'text-gray-600 hover:bg-gray-50 hover:text-gray-900',
          collapsed && 'justify-center px-2'
        )}
        title={collapsed ? item.title : undefined}
      >
        <div className="relative">
          <div
            className={cn(
              'flex h-7 w-7 items-center justify-center rounded-lg transition-all',
              isActive ? item.bgColor : 'bg-gray-100 group-hover:bg-gray-200'
            )}
          >
            <item.icon
              className={cn(
                'h-4 w-4 flex-shrink-0',
                isActive ? item.color : 'text-gray-500'
              )}
            />
          </div>
          {collapsed && unreadCount && unreadCount > 0 && (
            <span className="absolute -top-1 -right-1 w-4 h-4 bg-red-500 text-white text-xs rounded-full flex items-center justify-center">
              {unreadCount > 9 ? '9+' : unreadCount}
            </span>
          )}
        </div>
        {!collapsed && (
          <>
            <span className="flex-1">{item.title}</span>
            {unreadCount && unreadCount > 0 && (
              <Badge className="bg-red-500 text-white text-xs px-1.5 py-0.5 min-w-[20px] flex items-center justify-center">
                {unreadCount > 99 ? '99+' : unreadCount}
              </Badge>
            )}
          </>
        )}
      </NavLink>
    );
  };

  return (
    <aside
      className={cn(
        'fixed left-0 top-16 z-40 h-[calc(100vh-4rem)] border-r border-gray-200 bg-white transition-all duration-300 shadow-sm',
        collapsed ? 'w-16' : 'w-64'
      )}
    >
      <div className="flex h-full flex-col">
        <nav className="flex-1 px-2 py-3 overflow-y-auto">
          {visibleGroups.map((group) => {
            const isExpanded = expandedGroups.has(group.id);
            const groupActive = isGroupActive(group);

            // For "Main" group, don't show header - just render items
            if (group.id === 'main') {
              return (
                <div key={group.id} className="mb-2">
                  {group.items.map(renderNavItem)}
                </div>
              );
            }

            return (
              <div key={group.id} className="mb-1">
                {/* Group Header */}
                {!collapsed ? (
                  <button
                    onClick={() => toggleGroup(group.id)}
                    className={cn(
                      'w-full flex items-center gap-2 px-3 py-2 text-xs font-semibold uppercase tracking-wider rounded-lg transition-colors',
                      groupActive
                        ? 'text-cyan-700 bg-cyan-50'
                        : 'text-gray-500 hover:text-gray-700 hover:bg-gray-50'
                    )}
                  >
                    <group.icon className={cn('h-4 w-4', groupActive ? group.color : 'text-gray-400')} />
                    <span className="flex-1 text-left">{group.title}</span>
                    <ChevronDown
                      className={cn(
                        'h-4 w-4 transition-transform',
                        isExpanded ? 'rotate-0' : '-rotate-90'
                      )}
                    />
                  </button>
                ) : (
                  <div className="flex justify-center py-1 mb-1">
                    <div className={cn('w-8 h-0.5 rounded-full', groupActive ? 'bg-cyan-400' : 'bg-gray-200')} />
                  </div>
                )}

                {/* Group Items */}
                {(collapsed || isExpanded) && (
                  <div className={cn('space-y-0.5', !collapsed && 'ml-2 mt-1')}>
                    {group.items.map(renderNavItem)}
                  </div>
                )}
              </div>
            );
          })}
        </nav>

        <div className="border-t border-gray-200 p-2">
          <Button
            variant="ghost"
            size="sm"
            onClick={onToggle}
            className="w-full justify-center text-gray-500 hover:text-cyan-600 hover:bg-cyan-50"
          >
            {collapsed ? (
              <ChevronRight className="h-4 w-4" />
            ) : (
              <ChevronLeft className="h-4 w-4" />
            )}
          </Button>
        </div>
      </div>
    </aside>
  );
}
