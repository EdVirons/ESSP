import { Link } from 'react-router-dom';
import {
  Package,
  AlertTriangle,
  Wrench,
  ArrowRight,
  TrendingUp,
  TrendingDown,
  Warehouse,
  RefreshCw,
  ClipboardCheck,
  ArrowRightLeft,
  MessageSquare,
  BarChart3,
} from 'lucide-react';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { useAuth } from '@/contexts/AuthContext';
import { cn } from '@/lib/utils';
import { useWarehouseDashboard } from '@/hooks/useWarehouseDashboard';
import { useUnreadCounts } from '@/hooks/useMessages';
import { formatDistanceToNow } from 'date-fns';

export function WarehouseManagerDashboard() {
  const { user } = useAuth();
  const displayName = user?.displayName || user?.username || 'Manager';

  // Fetch dashboard data
  const { data: dashboardData, isLoading } = useWarehouseDashboard();
  const { data: unreadCounts } = useUnreadCounts();

  const unreadMessages = unreadCounts?.messages || 0;

  // Default values when loading
  const lowStockCount = dashboardData?.lowStockCount || 0;
  const pendingWorkOrders = dashboardData?.pendingWorkOrders || 0;
  const todayMovements = dashboardData?.todayMovements || 0;
  const totalParts = dashboardData?.totalParts || 0;
  const lowStockItems = dashboardData?.lowStockItems || [];
  const pendingPartIssues = dashboardData?.pendingPartIssues || [];
  const recentActivity = dashboardData?.recentActivity || [];
  const partsCategories = dashboardData?.partsCategories || {};

  return (
    <div className="space-y-6">
      {/* Welcome Header - Warehouse Theme (amber/orange gradient) */}
      <div className="relative overflow-hidden rounded-2xl bg-gradient-to-r from-amber-600 via-orange-600 to-red-600 p-6 text-white shadow-lg">
        {/* Decorative elements */}
        <div className="absolute top-0 right-0 w-64 h-64 bg-white/5 rounded-full -translate-y-32 translate-x-32" />
        <div className="absolute bottom-0 left-0 w-48 h-48 bg-white/5 rounded-full translate-y-24 -translate-x-24" />

        <div className="relative">
          <div className="flex items-center gap-3 mb-2">
            <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-white/20 backdrop-blur">
              <Warehouse className="h-6 w-6" />
            </div>
            <div>
              <h1 className="text-2xl md:text-3xl font-bold">Welcome, {displayName}!</h1>
              <p className="text-amber-100">Warehouse Manager Dashboard</p>
            </div>
          </div>

          {/* Summary Stats */}
          <div className="mt-6 grid grid-cols-2 md:grid-cols-4 gap-4">
            <div className="bg-white/10 backdrop-blur rounded-xl p-3">
              <p className="text-2xl font-bold">{lowStockCount}</p>
              <p className="text-amber-200 text-sm">Low Stock Alerts</p>
            </div>
            <div className="bg-white/10 backdrop-blur rounded-xl p-3">
              <p className="text-2xl font-bold">{pendingWorkOrders}</p>
              <p className="text-amber-200 text-sm">Pending Parts</p>
            </div>
            <div className="bg-white/10 backdrop-blur rounded-xl p-3">
              <p className="text-2xl font-bold">{todayMovements}</p>
              <p className="text-amber-200 text-sm">Today's Movements</p>
            </div>
            <div className="bg-white/10 backdrop-blur rounded-xl p-3">
              <p className="text-2xl font-bold">{unreadMessages}</p>
              <p className="text-amber-200 text-sm">Messages</p>
            </div>
          </div>
        </div>
      </div>

      {/* Quick Actions */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <Link
          to="/parts-catalog"
          className="flex items-center gap-4 p-4 rounded-xl bg-gradient-to-r from-amber-500 to-orange-500 text-white shadow-md hover:shadow-lg transition-all hover:-translate-y-0.5"
        >
          <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-white/20">
            <Package className="h-6 w-6" />
          </div>
          <div>
            <p className="font-semibold text-lg">Parts Catalog</p>
            <p className="text-amber-100 text-sm">Manage parts</p>
          </div>
        </Link>

        <Link
          to="/service-shops"
          className="flex items-center gap-4 p-4 rounded-xl bg-gradient-to-r from-rose-500 to-pink-500 text-white shadow-md hover:shadow-lg transition-all hover:-translate-y-0.5"
        >
          <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-white/20">
            <ClipboardCheck className="h-6 w-6" />
          </div>
          <div>
            <p className="font-semibold text-lg">Service Shops</p>
            <p className="text-rose-100 text-sm">View inventory</p>
          </div>
        </Link>

        <Link
          to="/devices"
          className="flex items-center gap-4 p-4 rounded-xl bg-gradient-to-r from-blue-500 to-cyan-500 text-white shadow-md hover:shadow-lg transition-all hover:-translate-y-0.5"
        >
          <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-white/20">
            <ArrowRightLeft className="h-6 w-6" />
          </div>
          <div>
            <p className="font-semibold text-lg">Devices</p>
            <p className="text-blue-100 text-sm">Device inventory</p>
          </div>
        </Link>

        <Link
          to="/work-orders"
          className="flex items-center gap-4 p-4 rounded-xl bg-gradient-to-r from-purple-500 to-indigo-500 text-white shadow-md hover:shadow-lg transition-all hover:-translate-y-0.5"
        >
          <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-white/20">
            <Wrench className="h-6 w-6" />
          </div>
          <div className="flex-1">
            <p className="font-semibold text-lg">Work Orders</p>
            <p className="text-purple-100 text-sm">Parts needed</p>
          </div>
          {pendingWorkOrders > 0 && (
            <Badge className="bg-white text-purple-600">
              {pendingWorkOrders}
            </Badge>
          )}
        </Link>
      </div>

      {/* Two Column Layout: Low Stock Alerts + Pending Work Orders */}
      <div className="grid gap-6 lg:grid-cols-2">
        {/* Low Stock Alerts */}
        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <CardTitle className="flex items-center gap-2">
                <AlertTriangle className="h-5 w-5 text-amber-600" />
                Low Stock Alerts
              </CardTitle>
              {lowStockCount > 0 && (
                <Badge className="bg-amber-100 text-amber-800">
                  {lowStockCount} items
                </Badge>
              )}
            </div>
          </CardHeader>
          <CardContent>
            {isLoading ? (
              <div className="text-center py-8 text-gray-500">Loading...</div>
            ) : lowStockItems.length === 0 ? (
              <div className="text-center py-8 text-gray-500">
                <Package className="h-12 w-12 mx-auto mb-3 text-green-500" />
                <p className="font-medium">All stock levels healthy!</p>
                <p className="text-sm">No items below reorder threshold.</p>
              </div>
            ) : (
              <div className="space-y-3">
                {lowStockItems.slice(0, 5).map((item) => (
                  <div
                    key={item.id}
                    className="flex items-center justify-between p-3 rounded-lg border border-amber-100 bg-amber-50/50"
                  >
                    <div className="flex-1">
                      <p className="font-medium text-gray-900">{item.partName}</p>
                      <p className="text-sm text-gray-500">{item.shopName}</p>
                    </div>
                    <div className="text-right">
                      <p className="text-lg font-bold text-amber-600">
                        {item.qtyAvailable}
                      </p>
                      <p className="text-xs text-gray-500">
                        of {item.reorderThreshold} min
                      </p>
                    </div>
                  </div>
                ))}
                {lowStockCount > 5 && (
                  <Link to="/parts-catalog?lowStock=true">
                    <Button variant="outline" className="w-full mt-2">
                      View All Low Stock ({lowStockCount})
                      <ArrowRight className="h-4 w-4 ml-2" />
                    </Button>
                  </Link>
                )}
              </div>
            )}
          </CardContent>
        </Card>

        {/* Pending Work Orders Needing Parts */}
        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <CardTitle className="flex items-center gap-2">
                <Wrench className="h-5 w-5 text-purple-600" />
                Work Orders Needing Parts
              </CardTitle>
            </div>
          </CardHeader>
          <CardContent>
            {isLoading ? (
              <div className="text-center py-8 text-gray-500">Loading...</div>
            ) : pendingPartIssues.length === 0 ? (
              <div className="text-center py-8 text-gray-500">
                <Wrench className="h-12 w-12 mx-auto mb-3 text-gray-300" />
                <p>No pending part issues</p>
              </div>
            ) : (
              <div className="space-y-3">
                {pendingPartIssues.slice(0, 5).map((wo) => (
                  <Link
                    key={wo.id}
                    to={`/work-orders/${wo.id}`}
                    className="block p-3 rounded-lg border border-gray-100 hover:border-purple-200 hover:bg-purple-50/50 transition-colors"
                  >
                    <div className="flex items-start justify-between">
                      <div>
                        <span className="text-sm font-mono text-gray-500">
                          {wo.id.slice(0, 12)}...
                        </span>
                        <p className="font-medium text-gray-900">{wo.title}</p>
                        <p className="text-sm text-gray-500">{wo.schoolName}</p>
                      </div>
                      <Badge
                        className={cn(
                          wo.priority === 'critical'
                            ? 'bg-red-100 text-red-800'
                            : wo.priority === 'high'
                              ? 'bg-orange-100 text-orange-800'
                              : wo.priority === 'medium'
                                ? 'bg-yellow-100 text-yellow-800'
                                : 'bg-green-100 text-green-800'
                        )}
                      >
                        {wo.priority}
                      </Badge>
                    </div>
                    <div className="mt-2 text-sm text-purple-600">
                      {wo.partsNeeded} part(s) required
                    </div>
                  </Link>
                ))}
                {pendingWorkOrders > 5 && (
                  <Link to="/work-orders?status=pending_parts">
                    <Button variant="outline" className="w-full mt-2">
                      View All Pending ({pendingWorkOrders})
                      <ArrowRight className="h-4 w-4 ml-2" />
                    </Button>
                  </Link>
                )}
              </div>
            )}
          </CardContent>
        </Card>
      </div>

      {/* Recent Activity */}
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <CardTitle className="flex items-center gap-2">
              <RefreshCw className="h-5 w-5 text-blue-600" />
              Recent Inventory Activity
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
          {isLoading ? (
            <div className="text-center py-8 text-gray-500">Loading...</div>
          ) : recentActivity.length === 0 ? (
            <div className="text-center py-8 text-gray-500">
              <RefreshCw className="h-12 w-12 mx-auto mb-3 text-gray-300" />
              <p>No recent activity</p>
            </div>
          ) : (
            <div className="space-y-3">
              {recentActivity.slice(0, 5).map((activity) => (
                <div
                  key={activity.id}
                  className="flex items-center gap-4 p-3 rounded-lg border border-gray-100"
                >
                  <div
                    className={cn(
                      'flex h-10 w-10 items-center justify-center rounded-full',
                      activity.type === 'receipt'
                        ? 'bg-green-100'
                        : activity.type === 'issue'
                          ? 'bg-blue-100'
                          : activity.type === 'adjustment'
                            ? 'bg-amber-100'
                            : 'bg-gray-100'
                    )}
                  >
                    {activity.type === 'receipt' && (
                      <TrendingUp className="h-5 w-5 text-green-600" />
                    )}
                    {activity.type === 'issue' && (
                      <TrendingDown className="h-5 w-5 text-blue-600" />
                    )}
                    {activity.type === 'adjustment' && (
                      <RefreshCw className="h-5 w-5 text-amber-600" />
                    )}
                    {!['receipt', 'issue', 'adjustment'].includes(activity.type) && (
                      <Package className="h-5 w-5 text-gray-600" />
                    )}
                  </div>
                  <div className="flex-1">
                    <p className="font-medium text-gray-900">{activity.description}</p>
                    <p className="text-sm text-gray-500">
                      {activity.partName} - {activity.actorName}
                    </p>
                  </div>
                  <span className="text-xs text-gray-400">
                    {formatDistanceToNow(new Date(activity.createdAt), { addSuffix: true })}
                  </span>
                </div>
              ))}
            </div>
          )}
        </CardContent>
      </Card>

      {/* Parts by Category Stats */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <BarChart3 className="h-5 w-5 text-emerald-600" />
            Parts by Category ({totalParts} total)
          </CardTitle>
        </CardHeader>
        <CardContent>
          {isLoading ? (
            <div className="text-center py-8 text-gray-500">Loading...</div>
          ) : Object.keys(partsCategories).length === 0 ? (
            <div className="text-center py-8 text-gray-500">
              <BarChart3 className="h-12 w-12 mx-auto mb-3 text-gray-300" />
              <p>No parts data available</p>
            </div>
          ) : (
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
              {Object.entries(partsCategories).map(([category, count]) => (
                <div key={category} className="p-4 rounded-lg bg-gray-50 text-center">
                  <p className="text-2xl font-bold text-gray-900">{count}</p>
                  <p className="text-sm text-gray-500 capitalize">{category}</p>
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
            Check your messages for updates on parts orders and work order requests.
          </p>
        </CardContent>
      </Card>
    </div>
  );
}
