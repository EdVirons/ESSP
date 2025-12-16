// API Configuration
export const API_BASE_URL = '/v1';
export const ADMIN_API_BASE_URL = '/v1';
// Use ws:// for localhost/http, wss:// for https
export const WS_BASE_URL = (() => {
  const isSecure = window.location.protocol === 'https:';
  const wsProtocol = isSecure ? 'wss:' : 'ws:';
  return `${wsProtocol}//${window.location.host}/ws`;
})();

// Pagination
export const DEFAULT_PAGE_SIZE = 50;
export const MAX_PAGE_SIZE = 200;

// Refresh intervals (ms)
export const HEALTH_CHECK_INTERVAL = 30_000;
export const METRICS_REFRESH_INTERVAL = 60_000;
export const ACTIVITY_REFRESH_INTERVAL = 15_000;

// Status colors
export const STATUS_COLORS = {
  // Incident statuses
  new: 'bg-blue-100 text-blue-800',
  acknowledged: 'bg-yellow-100 text-yellow-800',
  in_progress: 'bg-purple-100 text-purple-800',
  escalated: 'bg-red-100 text-red-800',
  resolved: 'bg-green-100 text-green-800',
  closed: 'bg-gray-100 text-gray-800',

  // Work order statuses
  draft: 'bg-gray-100 text-gray-800',
  assigned: 'bg-blue-100 text-blue-800',
  in_repair: 'bg-yellow-100 text-yellow-800',
  qa: 'bg-purple-100 text-purple-800',
  completed: 'bg-green-100 text-green-800',
  approved: 'bg-emerald-100 text-emerald-800',

  // Severity colors
  low: 'bg-green-100 text-green-800',
  medium: 'bg-yellow-100 text-yellow-800',
  high: 'bg-orange-100 text-orange-800',
  critical: 'bg-red-100 text-red-800',

  // Health statuses
  healthy: 'bg-green-100 text-green-800',
  degraded: 'bg-yellow-100 text-yellow-800',
  unhealthy: 'bg-red-100 text-red-800',
} as const;

// Navigation items
export const NAV_ITEMS = [
  {
    title: 'Overview',
    href: '/overview',
    icon: 'LayoutDashboard',
  },
  {
    title: 'Incidents',
    href: '/incidents',
    icon: 'AlertTriangle',
  },
  {
    title: 'Work Orders',
    href: '/work-orders',
    icon: 'Wrench',
  },
  {
    title: 'Projects',
    href: '/projects',
    icon: 'Layers',
  },
  {
    title: 'Service Shops',
    href: '/service-shops',
    icon: 'Store',
  },
  {
    title: 'Schools',
    href: '/schools',
    icon: 'School',
  },
  {
    title: 'Devices',
    href: '/devices',
    icon: 'Laptop',
  },
  {
    title: 'Parts Catalog',
    href: '/parts-catalog',
    icon: 'Package',
  },
  {
    title: 'SSOT Sync',
    href: '/ssot-sync',
    icon: 'RefreshCw',
  },
  {
    title: 'Audit Logs',
    href: '/audit-logs',
    icon: 'FileText',
  },
  {
    title: 'Settings',
    href: '/settings',
    icon: 'Settings',
  },
] as const;
