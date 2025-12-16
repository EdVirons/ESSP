import { useMemo } from 'react';
import { Link } from 'react-router-dom';
import {
  AlertTriangle,
  Wrench,
  School,
  Package,
  Laptop,
  FileText,
  Plus,
  ArrowRight,
} from 'lucide-react';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { useAuth } from '@/contexts/AuthContext';

interface QuickAction {
  label: string;
  description: string;
  icon: React.ElementType;
  href: string;
  gradient: string;
  iconColor: string;
  hoverBg: string;
  permission?: string;
  roles?: string[];
}

const quickActions: QuickAction[] = [
  {
    label: 'New Incident',
    description: 'Report a new issue',
    icon: AlertTriangle,
    href: '/incidents?action=create',
    gradient: 'from-amber-500 to-orange-500',
    iconColor: 'text-amber-600',
    hoverBg: 'hover:bg-amber-50',
    permission: 'incident:create',
    // Field techs don't create incidents
    roles: ['ssp_admin', 'ssp_support_agent', 'ssp_school_contact'],
  },
  {
    label: 'Create Work Order',
    description: 'Start a repair task',
    icon: Wrench,
    href: '/work-orders?action=create',
    gradient: 'from-blue-500 to-cyan-500',
    iconColor: 'text-blue-600',
    hoverBg: 'hover:bg-blue-50',
    permission: 'workorder:create',
    // Field techs can't create work orders, only work on them
    roles: ['ssp_admin', 'ssp_support_agent', 'ssp_lead_tech'],
  },
  {
    label: 'View Schools',
    description: 'Browse school list',
    icon: School,
    href: '/schools',
    gradient: 'from-emerald-500 to-teal-500',
    iconColor: 'text-emerald-600',
    hoverBg: 'hover:bg-emerald-50',
    permission: 'school:read',
    // Exclude field_tech (sees school info within work orders) and school_contact
    roles: ['ssp_admin', 'ssp_support_agent', 'ssp_lead_tech', 'ssp_demo_team', 'ssp_sales', 'ssp_contractor'],
  },
  {
    label: 'Device Inventory',
    description: 'Check device status',
    icon: Laptop,
    href: '/devices',
    gradient: 'from-purple-500 to-indigo-500',
    iconColor: 'text-purple-600',
    hoverBg: 'hover:bg-purple-50',
    permission: 'device:read',
    // Only support and warehouse see devices
    roles: ['ssp_admin', 'ssp_support_agent', 'ssp_warehouse_manager'],
  },
  {
    label: 'Parts Catalog',
    description: 'Manage spare parts',
    icon: Package,
    href: '/parts-catalog',
    gradient: 'from-rose-500 to-pink-500',
    iconColor: 'text-rose-600',
    hoverBg: 'hover:bg-rose-50',
    permission: 'parts:read',
    // Parts for support, lead tech, supplier, warehouse
    roles: ['ssp_admin', 'ssp_support_agent', 'ssp_lead_tech', 'ssp_supplier', 'ssp_warehouse_manager'],
  },
  {
    label: 'Reports',
    description: 'View analytics',
    icon: FileText,
    href: '/audit-logs',
    gradient: 'from-slate-500 to-gray-600',
    iconColor: 'text-slate-600',
    hoverBg: 'hover:bg-slate-50',
    roles: ['ssp_admin'],
  },
];

export function QuickActions() {
  const { hasPermission, hasRole } = useAuth();

  // Filter actions based on permissions and roles
  const visibleActions = useMemo(() => {
    // Admin sees everything
    if (hasRole('ssp_admin')) return quickActions;

    return quickActions.filter((action) => {
      // No restrictions = visible to all
      if (!action.permission && !action.roles) return true;

      // If both permission and roles are specified, BOTH must be satisfied
      if (action.permission && action.roles) {
        return hasPermission(action.permission) && action.roles.some((r) => hasRole(r));
      }

      // Only permission specified - check permission
      if (action.permission && !action.roles) {
        return hasPermission(action.permission);
      }

      // Only roles specified - check roles
      if (action.roles && !action.permission) {
        return action.roles.some((r) => hasRole(r));
      }

      return false;
    });
  }, [hasPermission, hasRole]);

  if (visibleActions.length === 0) return null;
  return (
    <Card className="border-0 shadow-md bg-white/80 backdrop-blur">
      <CardHeader className="pb-3">
        <CardTitle className="flex items-center gap-2 text-base">
          <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-gradient-to-br from-cyan-500 to-teal-600">
            <Plus className="h-4 w-4 text-white" />
          </div>
          Quick Actions
        </CardTitle>
      </CardHeader>
      <CardContent>
        <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-3">
          {visibleActions.map((action) => {
            const Icon = action.icon;
            return (
              <Link
                key={action.label}
                to={action.href}
                className={`group relative flex flex-col items-center justify-center p-4 rounded-xl bg-white border border-gray-100 transition-all duration-200 ${action.hoverBg} hover:shadow-md hover:border-transparent hover:-translate-y-0.5`}
              >
                {/* Gradient top border on hover */}
                <div className={`absolute top-0 left-0 right-0 h-1 rounded-t-xl bg-gradient-to-r ${action.gradient} opacity-0 group-hover:opacity-100 transition-opacity`} />

                <div className={`flex h-12 w-12 items-center justify-center rounded-xl bg-gradient-to-br ${action.gradient} mb-3 shadow-sm group-hover:shadow-md transition-shadow`}>
                  <Icon className="h-6 w-6 text-white" />
                </div>

                <span className="text-sm font-medium text-gray-900 text-center">
                  {action.label}
                </span>
                <span className="text-xs text-gray-500 text-center hidden md:block mt-1">
                  {action.description}
                </span>

                {/* Hover arrow */}
                <ArrowRight className="absolute bottom-2 right-2 h-4 w-4 text-gray-300 opacity-0 group-hover:opacity-100 transition-opacity" />
              </Link>
            );
          })}
        </div>
      </CardContent>
    </Card>
  );
}
