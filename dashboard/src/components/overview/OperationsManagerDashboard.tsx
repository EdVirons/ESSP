import * as React from 'react';
import { Link } from 'react-router-dom';
import {
  Building2,
  CheckCircle2,
  ArrowRight,
  ClipboardCheck,
  Package,
  Activity,
  Users,
  AlertTriangle,
  MessageSquare,
  BarChart3,
  MapPin,
  TrendingUp,
  RefreshCw,
} from 'lucide-react';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { useAuth } from '@/contexts/AuthContext';
import { cn } from '@/lib/utils';
import { useUnreadCounts } from '@/hooks/useMessages';

// Operations Manager Dashboard - Global field operations overview
export function OperationsManagerDashboard() {
  const { user } = useAuth();
  const displayName = user?.displayName || user?.username || 'Operations Manager';
  const { data: unreadCounts } = useUnreadCounts();
  const unreadMessages = unreadCounts?.messages || 0;

  // TODO: Replace with real API hook (useOperationsManagerDashboard)
  // For now, using placeholder data to demonstrate the dashboard structure
  const isLoading = false;

  // Global metrics across all service shops
  const globalMetrics = {
    totalShops: 5,
    activeWorkOrders: 47,
    pendingApprovals: 12,
    staffOnDuty: 23,
  };

  // Shop performance summary
  const shopPerformance = [
    { id: '1', name: 'Nairobi Central', activeWOs: 15, completedToday: 8, staff: 6, status: 'optimal' },
    { id: '2', name: 'Mombasa Hub', activeWOs: 12, completedToday: 5, staff: 5, status: 'optimal' },
    { id: '3', name: 'Kisumu Branch', activeWOs: 8, completedToday: 3, staff: 4, status: 'warning' },
    { id: '4', name: 'Nakuru Center', activeWOs: 7, completedToday: 4, staff: 4, status: 'optimal' },
    { id: '5', name: 'Eldoret Office', activeWOs: 5, completedToday: 2, staff: 4, status: 'critical' },
  ];

  // Capture initial time once on mount to avoid impure Date.now() calls during render
  const [mountTime] = React.useState(() => Date.now());

  // Pending approvals from all shops
  const pendingApprovals = [
    { id: '1', title: 'Laptop Screen Replacement', shopName: 'Nairobi Central', schoolName: 'Greenwood Academy', priority: 'high', requestedAt: new Date(mountTime - 3600000).toISOString() },
    { id: '2', title: 'Server Rack Installation', shopName: 'Mombasa Hub', schoolName: 'Coastal International', priority: 'critical', requestedAt: new Date(mountTime - 7200000).toISOString() },
    { id: '3', title: 'Network Switch Upgrade', shopName: 'Kisumu Branch', schoolName: 'Lakeside School', priority: 'medium', requestedAt: new Date(mountTime - 10800000).toISOString() },
  ];

  // Inventory alerts across all shops
  const inventoryAlerts = [
    { id: '1', shopName: 'Eldoret Office', partName: 'Laptop Screens 15.6"', currentStock: 2, minStock: 10, severity: 'critical' },
    { id: '2', shopName: 'Kisumu Branch', partName: 'Keyboard Replacements', currentStock: 5, minStock: 15, severity: 'warning' },
    { id: '3', shopName: 'Nairobi Central', partName: 'Power Adapters', currentStock: 8, minStock: 20, severity: 'warning' },
  ];

  // Recent cross-shop activity
  const recentActivity = [
    { id: '1', description: 'Work order transferred from Nairobi to Mombasa', actor: 'System', createdAt: new Date(mountTime - 1800000).toISOString(), type: 'transfer' },
    { id: '2', description: 'New technician onboarded at Kisumu Branch', actor: 'Admin', createdAt: new Date(mountTime - 3600000).toISOString(), type: 'staff' },
    { id: '3', description: 'Inventory replenishment approved for Eldoret', actor: 'Ops Manager', createdAt: new Date(mountTime - 5400000).toISOString(), type: 'inventory' },
  ];

  // Work order status breakdown across all shops
  const workOrderMetrics = {
    inProgress: 23,
    pending: 15,
    scheduled: 19,
    completed: 156,
    overdue: 3,
  };

  const formatTimeAgo = React.useCallback((dateStr: string) => {
    const diff = mountTime - new Date(dateStr).getTime();
    const hours = Math.floor(diff / 3600000);
    if (hours < 1) return 'Just now';
    if (hours === 1) return '1 hour ago';
    if (hours < 24) return `${hours} hours ago`;
    return `${Math.floor(hours / 24)} days ago`;
  }, [mountTime]);

  return (
    <div className="space-y-6">
      {/* Welcome Header - Orange/Amber Theme (Operations Manager) */}
      <div className="relative overflow-hidden rounded-2xl bg-gradient-to-r from-orange-600 via-amber-600 to-yellow-600 p-6 text-white shadow-lg">
        {/* Decorative elements */}
        <div className="absolute top-0 right-0 w-64 h-64 bg-white/5 rounded-full -translate-y-32 translate-x-32" />
        <div className="absolute bottom-0 left-0 w-48 h-48 bg-white/5 rounded-full translate-y-24 -translate-x-24" />

        <div className="relative">
          <div className="flex items-center gap-3 mb-2">
            <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-white/20 backdrop-blur">
              <Building2 className="h-6 w-6" />
            </div>
            <div>
              <h1 className="text-2xl md:text-3xl font-bold">Welcome, {displayName}!</h1>
              <p className="text-orange-100">Operations Manager Dashboard</p>
            </div>
          </div>

          {/* Global Summary Stats */}
          <div className="mt-6 grid grid-cols-2 md:grid-cols-4 gap-4">
            <div className="bg-white/10 backdrop-blur rounded-xl p-3">
              <p className="text-2xl font-bold">{globalMetrics.totalShops}</p>
              <p className="text-orange-200 text-sm">Service Shops</p>
            </div>
            <div className="bg-white/10 backdrop-blur rounded-xl p-3">
              <p className="text-2xl font-bold">{globalMetrics.activeWorkOrders}</p>
              <p className="text-orange-200 text-sm">Active Work Orders</p>
            </div>
            <div className="bg-white/10 backdrop-blur rounded-xl p-3">
              <p className="text-2xl font-bold">{globalMetrics.pendingApprovals}</p>
              <p className="text-orange-200 text-sm">Pending Approvals</p>
            </div>
            <div className="bg-white/10 backdrop-blur rounded-xl p-3">
              <p className="text-2xl font-bold">{globalMetrics.staffOnDuty}</p>
              <p className="text-orange-200 text-sm">Staff On Duty</p>
            </div>
          </div>
        </div>
      </div>

      {/* Quick Actions - Global Operations */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <Link
          to="/service-shops"
          className="flex items-center gap-4 p-4 rounded-xl bg-gradient-to-r from-orange-500 to-amber-500 text-white shadow-md hover:shadow-lg transition-all hover:-translate-y-0.5"
        >
          <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-white/20">
            <Building2 className="h-6 w-6" />
          </div>
          <div>
            <p className="font-semibold text-lg">Manage Shops</p>
            <p className="text-orange-100 text-sm">All locations</p>
          </div>
        </Link>

        <Link
          to="/work-orders?status=pending_approval"
          className="flex items-center gap-4 p-4 rounded-xl bg-gradient-to-r from-amber-500 to-yellow-500 text-white shadow-md hover:shadow-lg transition-all hover:-translate-y-0.5"
        >
          <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-white/20">
            <ClipboardCheck className="h-6 w-6" />
          </div>
          <div>
            <p className="font-semibold text-lg">Approvals</p>
            <p className="text-amber-100 text-sm">Review requests</p>
          </div>
          {globalMetrics.pendingApprovals > 0 && (
            <Badge className="bg-white text-amber-600 ml-auto">
              {globalMetrics.pendingApprovals}
            </Badge>
          )}
        </Link>

        <Link
          to="/staff"
          className="flex items-center gap-4 p-4 rounded-xl bg-gradient-to-r from-yellow-500 to-lime-500 text-white shadow-md hover:shadow-lg transition-all hover:-translate-y-0.5"
        >
          <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-white/20">
            <Users className="h-6 w-6" />
          </div>
          <div>
            <p className="font-semibold text-lg">Staff</p>
            <p className="text-yellow-100 text-sm">All technicians</p>
          </div>
        </Link>

        <Link
          to="/reports"
          className="flex items-center gap-4 p-4 rounded-xl bg-gradient-to-r from-lime-500 to-green-500 text-white shadow-md hover:shadow-lg transition-all hover:-translate-y-0.5"
        >
          <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-white/20">
            <BarChart3 className="h-6 w-6" />
          </div>
          <div>
            <p className="font-semibold text-lg">Reports</p>
            <p className="text-lime-100 text-sm">Global analytics</p>
          </div>
        </Link>
      </div>

      {/* Service Shop Performance Overview */}
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <CardTitle className="flex items-center gap-2">
              <MapPin className="h-5 w-5 text-orange-600" />
              Service Shop Performance
            </CardTitle>
            <Link to="/service-shops">
              <Button variant="outline" size="sm">
                View All
                <ArrowRight className="ml-2 h-4 w-4" />
              </Button>
            </Link>
          </div>
        </CardHeader>
        <CardContent>
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-5">
            {shopPerformance.map((shop) => (
              <Link
                key={shop.id}
                to={`/service-shops/${shop.id}`}
                className={cn(
                  "p-4 rounded-xl border-2 transition-all hover:-translate-y-0.5",
                  shop.status === 'optimal' && "border-green-200 bg-green-50/50 hover:border-green-300",
                  shop.status === 'warning' && "border-amber-200 bg-amber-50/50 hover:border-amber-300",
                  shop.status === 'critical' && "border-red-200 bg-red-50/50 hover:border-red-300"
                )}
              >
                <div className="flex items-center justify-between mb-3">
                  <h3 className="font-semibold text-gray-900 truncate">{shop.name}</h3>
                  <div className={cn(
                    "h-3 w-3 rounded-full",
                    shop.status === 'optimal' && "bg-green-500",
                    shop.status === 'warning' && "bg-amber-500",
                    shop.status === 'critical' && "bg-red-500"
                  )} />
                </div>
                <div className="space-y-1 text-sm">
                  <div className="flex justify-between">
                    <span className="text-gray-500">Active WOs:</span>
                    <span className="font-medium">{shop.activeWOs}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-gray-500">Completed:</span>
                    <span className="font-medium text-green-600">{shop.completedToday}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-gray-500">Staff:</span>
                    <span className="font-medium">{shop.staff}</span>
                  </div>
                </div>
              </Link>
            ))}
          </div>
        </CardContent>
      </Card>

      {/* Two Column: Pending Approvals + Inventory Alerts */}
      <div className="grid gap-6 lg:grid-cols-2">
        {/* Pending Approvals (All Shops) */}
        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <CardTitle className="flex items-center gap-2">
                <ClipboardCheck className="h-5 w-5 text-amber-600" />
                Pending Approvals
              </CardTitle>
              <Badge className="bg-amber-100 text-amber-800">
                {pendingApprovals.length} pending
              </Badge>
            </div>
          </CardHeader>
          <CardContent>
            {isLoading ? (
              <div className="text-center py-8 text-gray-500">Loading...</div>
            ) : pendingApprovals.length === 0 ? (
              <div className="text-center py-8 text-gray-500">
                <CheckCircle2 className="h-12 w-12 mx-auto mb-3 text-green-500" />
                <p className="font-medium">All caught up!</p>
                <p className="text-sm">No pending approvals across any shop.</p>
              </div>
            ) : (
              <div className="space-y-3">
                {pendingApprovals.map((approval) => (
                  <Link
                    key={approval.id}
                    to={`/work-orders/${approval.id}`}
                    className="block p-3 rounded-lg border border-gray-100 hover:border-amber-200 hover:bg-amber-50/50 transition-colors"
                  >
                    <div className="flex items-start justify-between">
                      <div>
                        <p className="font-medium text-gray-900">{approval.title}</p>
                        <p className="text-sm text-gray-500">{approval.schoolName}</p>
                        <p className="text-xs text-gray-400 mt-1">
                          <MapPin className="h-3 w-3 inline mr-1" />
                          {approval.shopName}
                        </p>
                      </div>
                      <Badge
                        className={cn(
                          approval.priority === 'critical'
                            ? 'bg-red-100 text-red-800'
                            : approval.priority === 'high'
                              ? 'bg-orange-100 text-orange-800'
                              : approval.priority === 'medium'
                                ? 'bg-yellow-100 text-yellow-800'
                                : 'bg-green-100 text-green-800'
                        )}
                      >
                        {approval.priority}
                      </Badge>
                    </div>
                    <p className="mt-2 text-xs text-gray-400">
                      Requested {formatTimeAgo(approval.requestedAt)}
                    </p>
                  </Link>
                ))}
                <Link to="/work-orders?status=pending_approval">
                  <Button variant="outline" className="w-full mt-2">
                    View All Approvals
                    <ArrowRight className="h-4 w-4 ml-2" />
                  </Button>
                </Link>
              </div>
            )}
          </CardContent>
        </Card>

        {/* Inventory Alerts (All Shops) */}
        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <CardTitle className="flex items-center gap-2">
                <AlertTriangle className="h-5 w-5 text-red-600" />
                Inventory Alerts
              </CardTitle>
              {inventoryAlerts.length > 0 && (
                <Badge className="bg-red-100 text-red-800">
                  {inventoryAlerts.length} alerts
                </Badge>
              )}
            </div>
          </CardHeader>
          <CardContent>
            {isLoading ? (
              <div className="text-center py-8 text-gray-500">Loading...</div>
            ) : inventoryAlerts.length === 0 ? (
              <div className="text-center py-8 text-gray-500">
                <Package className="h-12 w-12 mx-auto mb-3 text-green-500" />
                <p className="font-medium">Inventory looks good!</p>
                <p className="text-sm">No low stock alerts at this time.</p>
              </div>
            ) : (
              <div className="space-y-3">
                {inventoryAlerts.map((alert) => (
                  <div
                    key={alert.id}
                    className={cn(
                      "p-3 rounded-lg border",
                      alert.severity === 'critical'
                        ? "border-red-200 bg-red-50/50"
                        : "border-amber-200 bg-amber-50/50"
                    )}
                  >
                    <div className="flex items-start justify-between">
                      <div>
                        <p className="font-medium text-gray-900">{alert.partName}</p>
                        <p className="text-sm text-gray-500">
                          <MapPin className="h-3 w-3 inline mr-1" />
                          {alert.shopName}
                        </p>
                      </div>
                      <Badge
                        className={cn(
                          alert.severity === 'critical'
                            ? 'bg-red-100 text-red-800'
                            : 'bg-amber-100 text-amber-800'
                        )}
                      >
                        {alert.severity}
                      </Badge>
                    </div>
                    <div className="mt-2 flex items-center gap-2 text-sm">
                      <span className={cn(
                        "font-medium",
                        alert.severity === 'critical' ? "text-red-600" : "text-amber-600"
                      )}>
                        {alert.currentStock} / {alert.minStock} min
                      </span>
                      <span className="text-gray-400">in stock</span>
                    </div>
                  </div>
                ))}
                <Link to="/parts-catalog">
                  <Button variant="outline" className="w-full mt-2">
                    Manage Inventory
                    <ArrowRight className="h-4 w-4 ml-2" />
                  </Button>
                </Link>
              </div>
            )}
          </CardContent>
        </Card>
      </div>

      {/* Global Work Order Metrics */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <TrendingUp className="h-5 w-5 text-blue-600" />
            Global Work Order Status
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-2 md:grid-cols-5 gap-4">
            <div className="p-4 rounded-lg bg-blue-50 text-center">
              <p className="text-3xl font-bold text-blue-600">{workOrderMetrics.inProgress}</p>
              <p className="text-sm text-blue-700">In Progress</p>
            </div>
            <div className="p-4 rounded-lg bg-amber-50 text-center">
              <p className="text-3xl font-bold text-amber-600">{workOrderMetrics.pending}</p>
              <p className="text-sm text-amber-700">Pending</p>
            </div>
            <div className="p-4 rounded-lg bg-purple-50 text-center">
              <p className="text-3xl font-bold text-purple-600">{workOrderMetrics.scheduled}</p>
              <p className="text-sm text-purple-700">Scheduled</p>
            </div>
            <div className="p-4 rounded-lg bg-green-50 text-center">
              <p className="text-3xl font-bold text-green-600">{workOrderMetrics.completed}</p>
              <p className="text-sm text-green-700">Completed</p>
            </div>
            <div className="p-4 rounded-lg bg-red-50 text-center">
              <p className="text-3xl font-bold text-red-600">{workOrderMetrics.overdue}</p>
              <p className="text-sm text-red-700">Overdue</p>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Recent Cross-Shop Activity */}
      <div className="grid gap-6 lg:grid-cols-2">
        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <CardTitle className="flex items-center gap-2">
                <Activity className="h-5 w-5 text-emerald-600" />
                Recent Operations Activity
              </CardTitle>
              <Link to="/audit-logs">
                <Button variant="outline" size="sm">
                  View All
                  <ArrowRight className="ml-2 h-4 w-4" />
                </Button>
              </Link>
            </div>
          </CardHeader>
          <CardContent>
            {recentActivity.length === 0 ? (
              <div className="text-center py-8 text-gray-500">
                <Activity className="h-12 w-12 mx-auto mb-3 text-gray-300" />
                <p>No recent activity</p>
              </div>
            ) : (
              <div className="space-y-3">
                {recentActivity.map((activity) => (
                  <div
                    key={activity.id}
                    className="flex items-center gap-4 p-3 rounded-lg border border-gray-100"
                  >
                    <div
                      className={cn(
                        'flex h-10 w-10 items-center justify-center rounded-full',
                        activity.type === 'transfer' && 'bg-blue-100',
                        activity.type === 'staff' && 'bg-green-100',
                        activity.type === 'inventory' && 'bg-purple-100'
                      )}
                    >
                      {activity.type === 'transfer' && (
                        <RefreshCw className="h-5 w-5 text-blue-600" />
                      )}
                      {activity.type === 'staff' && (
                        <Users className="h-5 w-5 text-green-600" />
                      )}
                      {activity.type === 'inventory' && (
                        <Package className="h-5 w-5 text-purple-600" />
                      )}
                    </div>
                    <div className="flex-1">
                      <p className="font-medium text-gray-900">{activity.description}</p>
                      <p className="text-sm text-gray-500">{activity.actor}</p>
                    </div>
                    <span className="text-xs text-gray-400">
                      {formatTimeAgo(activity.createdAt)}
                    </span>
                  </div>
                ))}
              </div>
            )}
          </CardContent>
        </Card>

        {/* Messages Quick Access */}
        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <CardTitle className="flex items-center gap-2">
                <MessageSquare className="h-5 w-5 text-blue-600" />
                Messages
                {unreadMessages > 0 && (
                  <Badge className="bg-blue-100 text-blue-800 ml-2">
                    {unreadMessages} unread
                  </Badge>
                )}
              </CardTitle>
              <Link to="/messages">
                <Button variant="outline" size="sm">
                  View All
                  <ArrowRight className="ml-2 h-4 w-4" />
                </Button>
              </Link>
            </div>
          </CardHeader>
          <CardContent>
            <p className="text-gray-500 text-center py-4">
              Coordinate with lead technicians across all service shops and manage global operations communications.
            </p>
            <div className="flex gap-2 justify-center">
              <Link to="/messages?compose=true">
                <Button variant="outline" size="sm">
                  <MessageSquare className="h-4 w-4 mr-2" />
                  New Message
                </Button>
              </Link>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
