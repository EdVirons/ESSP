import {
  LayoutDashboard,
  AlertTriangle,
  Wrench,
  Layers,
  Store,
  School,
  Laptop,
  Package,
  RefreshCw,
  FileText,
  Settings,
  MessageSquare,
  Headphones,
  BookOpen,
  TrendingUp,
  Target,
  Presentation,
  BarChart3,
  Boxes,
  Users,
  Shield,
  PlusCircle,
  Clock,
  CheckCircle,
  List,
  Monitor,
  Contact,
  type LucideIcon,
} from 'lucide-react';

export interface NavItem {
  title: string;
  href: string;
  icon: LucideIcon;
  color: string;
  bgColor: string;
  // Permissions required to see this item (OR logic - user needs ANY of these)
  permissions?: string[];
  // Roles that can see this item (OR logic - user needs ANY of these, fallback if no permissions)
  roles?: string[];
}

export interface NavGroup {
  id: string;
  title: string;
  icon: LucideIcon;
  color: string;
  items: NavItem[];
  // Roles that can see this group (if any item is visible, group is visible)
  roles?: string[];
}

// Flat navigation items (for backward compatibility)
export const navItems: NavItem[] = [
  {
    title: 'Overview',
    href: '/overview',
    icon: LayoutDashboard,
    color: 'text-cyan-600',
    bgColor: 'bg-cyan-100',
  },
  {
    title: 'Incidents',
    href: '/incidents',
    icon: AlertTriangle,
    color: 'text-amber-600',
    bgColor: 'bg-amber-100',
    permissions: ['incident:read', 'incident:create'],
    roles: ['ssp_admin', 'ssp_support_agent', 'ssp_school_contact'],
  },
  {
    title: 'Messages',
    href: '/messages',
    icon: MessageSquare,
    color: 'text-blue-600',
    bgColor: 'bg-blue-100',
    permissions: ['messages:read'],
  },
  {
    title: 'Live Chat',
    href: '/live-chat',
    icon: Headphones,
    color: 'text-cyan-600',
    bgColor: 'bg-cyan-100',
    permissions: ['chat:accept'],
    roles: ['ssp_admin', 'ssp_support_agent'],
  },
  {
    title: 'Work Orders',
    href: '/work-orders',
    icon: Wrench,
    color: 'text-blue-600',
    bgColor: 'bg-blue-100',
    // Note: permissions removed so school_contact and sales_marketing don't see this
    // even though they have workorder:read for API access
    roles: [
      'ssp_admin', 'ssp_support_agent', 'ssp_field_tech', 'ssp_lead_tech',
      'ssp_contractor', 'ssp_supplier', 'ssp_warehouse_manager',
    ],
  },
  {
    title: 'Sales Dashboard',
    href: '/sales',
    icon: TrendingUp,
    color: 'text-emerald-600',
    bgColor: 'bg-emerald-100',
    permissions: ['demo:manage', 'reporting:sales'],
    roles: ['ssp_admin', 'ssp_sales_marketing'],
  },
  {
    title: 'Demo Pipeline',
    href: '/demo-pipeline',
    icon: Target,
    color: 'text-violet-600',
    bgColor: 'bg-violet-100',
    permissions: ['demo:pipeline'],
    roles: ['ssp_admin', 'ssp_sales_marketing', 'ssp_demo_team'],
  },
  {
    title: 'Presentations',
    href: '/presentations',
    icon: Presentation,
    color: 'text-pink-600',
    bgColor: 'bg-pink-100',
    permissions: ['presentations:view'],
    roles: ['ssp_admin', 'ssp_sales_marketing'],
  },
  {
    title: 'Marketing KB',
    href: '/marketing-kb',
    icon: BookOpen,
    color: 'text-orange-600',
    bgColor: 'bg-orange-100',
    permissions: ['mkb:read'],
    roles: ['ssp_admin', 'ssp_sales_marketing'],
  },
  {
    title: 'Projects',
    href: '/projects',
    icon: Layers,
    color: 'text-purple-600',
    bgColor: 'bg-purple-100',
    permissions: ['project:read'],
    roles: ['ssp_admin', 'ssp_demo_team', 'ssp_sales_marketing'],
  },
  {
    title: 'Service Shops',
    href: '/service-shops',
    icon: Store,
    color: 'text-orange-600',
    bgColor: 'bg-orange-100',
    permissions: ['serviceshop:read'],
    roles: ['ssp_admin', 'ssp_support_agent', 'ssp_lead_tech', 'ssp_warehouse_manager'],
  },
  {
    title: 'Staff',
    href: '/staff',
    icon: Users,
    color: 'text-indigo-600',
    bgColor: 'bg-indigo-100',
    permissions: ['ops:manage_staff'],
    roles: ['ssp_admin', 'ssp_ops_manager'],
  },
  {
    title: 'HR Directory',
    href: '/hr',
    icon: Contact,
    color: 'text-violet-600',
    bgColor: 'bg-violet-100',
    permissions: ['hr:read', 'people:read'],
    roles: ['ssp_admin', 'ssp_ops_manager'],
  },
  {
    title: 'Schools',
    href: '/schools',
    icon: School,
    color: 'text-green-600',
    bgColor: 'bg-green-100',
    // Note: school:read permission is NOT used here because school contacts have it
    // for their own school info, but should not see the full directory
    roles: [
      'ssp_admin', 'ssp_support_agent', 'ssp_lead_tech', 'ssp_demo_team',
      'ssp_sales_marketing', 'ssp_contractor',
    ],
  },
  {
    title: 'Devices',
    href: '/devices',
    icon: Laptop,
    color: 'text-indigo-600',
    bgColor: 'bg-indigo-100',
    permissions: ['device:read'],
    roles: ['ssp_admin', 'ssp_support_agent', 'ssp_lead_tech', 'ssp_warehouse_manager'],
  },
  {
    title: 'School Inventory',
    href: '/school-inventory',
    icon: Monitor,
    color: 'text-teal-600',
    bgColor: 'bg-teal-100',
    permissions: ['inventory:read'],
    roles: ['ssp_school_contact'],
  },
  {
    title: 'Parts Catalog',
    href: '/parts-catalog',
    icon: Package,
    color: 'text-rose-600',
    bgColor: 'bg-rose-100',
    permissions: ['parts:read'],
    roles: ['ssp_admin', 'ssp_support_agent', 'ssp_lead_tech', 'ssp_supplier', 'ssp_warehouse_manager'],
  },
  {
    title: 'Knowledge Base',
    href: '/knowledge-base',
    icon: BookOpen,
    color: 'text-emerald-600',
    bgColor: 'bg-emerald-100',
    permissions: ['kb:read'],
    roles: ['ssp_admin', 'ssp_field_tech', 'ssp_lead_tech', 'ssp_support_agent'],
  },
  {
    title: 'SSOT Sync',
    href: '/ssot-sync',
    icon: RefreshCw,
    color: 'text-teal-600',
    bgColor: 'bg-teal-100',
    permissions: ['ssot:read'],
    roles: ['ssp_admin'],
  },
  {
    title: 'Reports',
    href: '/reports',
    icon: BarChart3,
    color: 'text-teal-600',
    bgColor: 'bg-teal-100',
    permissions: ['reports:read'],
    roles: ['ssp_admin', 'ssp_support_agent', 'ssp_lead_tech', 'ssp_warehouse_manager', 'ssp_sales_marketing'],
  },
  {
    title: 'Audit Logs',
    href: '/audit-logs',
    icon: FileText,
    color: 'text-slate-600',
    bgColor: 'bg-slate-100',
    roles: ['ssp_admin'],
  },
  {
    title: 'Settings',
    href: '/settings',
    icon: Settings,
    color: 'text-gray-600',
    bgColor: 'bg-gray-100',
    roles: ['ssp_admin'],
  },
];

// Grouped navigation structure - organized by role functionality
export const navGroups: NavGroup[] = [
  {
    id: 'main',
    title: 'Main',
    icon: LayoutDashboard,
    color: 'text-cyan-600',
    items: [
      {
        title: 'Overview',
        href: '/overview',
        icon: LayoutDashboard,
        color: 'text-cyan-600',
        bgColor: 'bg-cyan-100',
      },
    ],
  },
  {
    id: 'my-incidents',
    title: 'My Incidents',
    icon: AlertTriangle,
    color: 'text-amber-600',
    roles: ['ssp_school_contact'],
    items: [
      {
        title: 'All Incidents',
        href: '/incidents',
        icon: List,
        color: 'text-amber-600',
        bgColor: 'bg-amber-100',
        roles: ['ssp_school_contact'],
      },
      {
        title: 'New',
        href: '/incidents?action=create',
        icon: PlusCircle,
        color: 'text-blue-600',
        bgColor: 'bg-blue-100',
        roles: ['ssp_school_contact'],
      },
      {
        title: 'Open',
        href: '/incidents?status=open',
        icon: Clock,
        color: 'text-yellow-600',
        bgColor: 'bg-yellow-100',
        roles: ['ssp_school_contact'],
      },
      {
        title: 'Resolved',
        href: '/incidents?status=resolved',
        icon: CheckCircle,
        color: 'text-green-600',
        bgColor: 'bg-green-100',
        roles: ['ssp_school_contact'],
      },
    ],
  },
  {
    id: 'my-school',
    title: 'My School',
    icon: School,
    color: 'text-teal-600',
    roles: ['ssp_school_contact'],
    items: [
      {
        title: 'Device Inventory',
        href: '/school-inventory',
        icon: Monitor,
        color: 'text-teal-600',
        bgColor: 'bg-teal-100',
        permissions: ['inventory:read'],
        roles: ['ssp_school_contact'],
      },
    ],
  },
  {
    id: 'support',
    title: 'Support & Operations',
    icon: Headphones,
    color: 'text-amber-600',
    roles: ['ssp_admin', 'ssp_support_agent', 'ssp_field_tech', 'ssp_lead_tech'],
    items: [
      {
        title: 'Incidents',
        href: '/incidents',
        icon: AlertTriangle,
        color: 'text-amber-600',
        bgColor: 'bg-amber-100',
        permissions: ['incident:read', 'incident:create'],
        roles: ['ssp_admin', 'ssp_support_agent'],
      },
      {
        title: 'Work Orders',
        href: '/work-orders',
        icon: Wrench,
        color: 'text-blue-600',
        bgColor: 'bg-blue-100',
        // Note: permissions removed so school_contact and sales_marketing don't see this
        roles: [
          'ssp_admin', 'ssp_support_agent', 'ssp_field_tech', 'ssp_lead_tech',
          'ssp_contractor', 'ssp_supplier', 'ssp_warehouse_manager',
        ],
      },
      {
        title: 'Live Chat',
        href: '/live-chat',
        icon: Headphones,
        color: 'text-cyan-600',
        bgColor: 'bg-cyan-100',
        permissions: ['chat:accept'],
        roles: ['ssp_admin', 'ssp_support_agent'],
      },
      {
        title: 'Messages',
        href: '/messages',
        icon: MessageSquare,
        color: 'text-blue-600',
        bgColor: 'bg-blue-100',
        permissions: ['messages:read'],
      },
    ],
  },
  {
    id: 'sales',
    title: 'Sales & Projects',
    icon: TrendingUp,
    color: 'text-emerald-600',
    roles: ['ssp_admin', 'ssp_sales_marketing', 'ssp_demo_team'],
    items: [
      {
        title: 'Sales Dashboard',
        href: '/sales',
        icon: TrendingUp,
        color: 'text-emerald-600',
        bgColor: 'bg-emerald-100',
        permissions: ['demo:manage', 'reporting:sales'],
        roles: ['ssp_admin', 'ssp_sales_marketing'],
      },
      {
        title: 'Demo Pipeline',
        href: '/demo-pipeline',
        icon: Target,
        color: 'text-violet-600',
        bgColor: 'bg-violet-100',
        permissions: ['demo:pipeline'],
        roles: ['ssp_admin', 'ssp_sales_marketing', 'ssp_demo_team'],
      },
      {
        title: 'Presentations',
        href: '/presentations',
        icon: Presentation,
        color: 'text-pink-600',
        bgColor: 'bg-pink-100',
        permissions: ['presentations:view'],
        roles: ['ssp_admin', 'ssp_sales_marketing'],
      },
      {
        title: 'Marketing KB',
        href: '/marketing-kb',
        icon: BookOpen,
        color: 'text-orange-600',
        bgColor: 'bg-orange-100',
        permissions: ['mkb:read'],
        roles: ['ssp_admin', 'ssp_sales_marketing'],
      },
      {
        title: 'Projects',
        href: '/projects',
        icon: Layers,
        color: 'text-purple-600',
        bgColor: 'bg-purple-100',
        permissions: ['project:read'],
        roles: ['ssp_admin', 'ssp_demo_team', 'ssp_sales_marketing'],
      },
    ],
  },
  {
    id: 'inventory',
    title: 'Inventory & Catalog',
    icon: Boxes,
    color: 'text-rose-600',
    roles: ['ssp_admin', 'ssp_ops_manager', 'ssp_warehouse_manager', 'ssp_supplier', 'ssp_lead_tech', 'ssp_support_agent'],
    items: [
      {
        title: 'Parts Catalog',
        href: '/parts-catalog',
        icon: Package,
        color: 'text-rose-600',
        bgColor: 'bg-rose-100',
        permissions: ['parts:read'],
        roles: ['ssp_admin', 'ssp_support_agent', 'ssp_lead_tech', 'ssp_supplier', 'ssp_warehouse_manager'],
      },
      {
        title: 'Devices',
        href: '/devices',
        icon: Laptop,
        color: 'text-indigo-600',
        bgColor: 'bg-indigo-100',
        permissions: ['device:read'],
        roles: ['ssp_admin', 'ssp_support_agent', 'ssp_lead_tech', 'ssp_warehouse_manager'],
      },
      {
        title: 'Service Shops',
        href: '/service-shops',
        icon: Store,
        color: 'text-orange-600',
        bgColor: 'bg-orange-100',
        permissions: ['serviceshop:read'],
        roles: ['ssp_admin', 'ssp_support_agent', 'ssp_lead_tech', 'ssp_warehouse_manager'],
      },
      {
        title: 'Staff',
        href: '/staff',
        icon: Users,
        color: 'text-indigo-600',
        bgColor: 'bg-indigo-100',
        permissions: ['ops:manage_staff'],
        roles: ['ssp_admin', 'ssp_ops_manager'],
      },
    ],
  },
  {
    id: 'directory',
    title: 'Directory',
    icon: Users,
    color: 'text-green-600',
    roles: ['ssp_admin', 'ssp_ops_manager', 'ssp_support_agent', 'ssp_lead_tech', 'ssp_demo_team', 'ssp_sales_marketing', 'ssp_contractor'],
    items: [
      {
        title: 'Schools',
        href: '/schools',
        icon: School,
        color: 'text-green-600',
        bgColor: 'bg-green-100',
        // Note: school:read permission is NOT used here because school contacts have it
        // for their own school info, but should not see the full directory
        roles: [
          'ssp_admin', 'ssp_support_agent', 'ssp_lead_tech', 'ssp_demo_team',
          'ssp_sales_marketing', 'ssp_contractor',
        ],
      },
      {
        title: 'HR Directory',
        href: '/hr',
        icon: Contact,
        color: 'text-violet-600',
        bgColor: 'bg-violet-100',
        permissions: ['hr:read', 'people:read'],
        roles: ['ssp_admin', 'ssp_ops_manager'],
      },
    ],
  },
  {
    id: 'knowledge',
    title: 'Knowledge & Reports',
    icon: BookOpen,
    color: 'text-teal-600',
    roles: ['ssp_admin', 'ssp_support_agent', 'ssp_lead_tech', 'ssp_field_tech', 'ssp_warehouse_manager', 'ssp_sales_marketing'],
    items: [
      {
        title: 'Knowledge Base',
        href: '/knowledge-base',
        icon: BookOpen,
        color: 'text-emerald-600',
        bgColor: 'bg-emerald-100',
        permissions: ['kb:read'],
        roles: ['ssp_admin', 'ssp_field_tech', 'ssp_lead_tech', 'ssp_support_agent'],
      },
      {
        title: 'Reports',
        href: '/reports',
        icon: BarChart3,
        color: 'text-teal-600',
        bgColor: 'bg-teal-100',
        permissions: ['reports:read'],
        roles: ['ssp_admin', 'ssp_support_agent', 'ssp_lead_tech', 'ssp_warehouse_manager', 'ssp_sales_marketing'],
      },
    ],
  },
  {
    id: 'admin',
    title: 'Administration',
    icon: Shield,
    color: 'text-slate-600',
    roles: ['ssp_admin'],
    items: [
      {
        title: 'SSOT Sync',
        href: '/ssot-sync',
        icon: RefreshCw,
        color: 'text-teal-600',
        bgColor: 'bg-teal-100',
        permissions: ['ssot:read'],
        roles: ['ssp_admin'],
      },
      {
        title: 'Audit Logs',
        href: '/audit-logs',
        icon: FileText,
        color: 'text-slate-600',
        bgColor: 'bg-slate-100',
        roles: ['ssp_admin'],
      },
      {
        title: 'Settings',
        href: '/settings',
        icon: Settings,
        color: 'text-gray-600',
        bgColor: 'bg-gray-100',
        roles: ['ssp_admin'],
      },
    ],
  },
];

// Role-based default routes
export const roleDefaultRoutes: Record<string, string> = {
  ssp_school_contact: '/incidents',
  ssp_support_agent: '/incidents',
  ssp_lead_tech: '/work-orders',
  ssp_field_tech: '/work-orders',
  ssp_warehouse_manager: '/parts-catalog',
  ssp_demo_team: '/projects',
  ssp_sales_marketing: '/sales',
  ssp_supplier: '/work-orders',
  ssp_contractor: '/work-orders',
};

// Get the default route for a user based on their roles
export function getDefaultRouteForRoles(roles: string[]): string {
  // Check roles in priority order
  for (const [role, route] of Object.entries(roleDefaultRoutes)) {
    if (roles.includes(role)) {
      return route;
    }
  }
  // Default fallback
  return '/overview';
}
