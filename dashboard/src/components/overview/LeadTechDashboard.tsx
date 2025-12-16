import { Link } from 'react-router-dom';
import {
  Users,
  Clock,
  CheckCircle2,
  ArrowRight,
  Calendar,
  ClipboardCheck,
  Package,
  Activity,
  Shield,
  Wrench,
  AlertTriangle,
  MessageSquare,
} from 'lucide-react';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { useAuth } from '@/contexts/AuthContext';
import { cn } from '@/lib/utils';
import { useLeadTechDashboard } from '@/hooks/useLeadTechDashboard';
import { useUnreadCounts } from '@/hooks/useMessages';
import { formatDistanceToNow } from 'date-fns';

export function LeadTechDashboard() {
  const { user } = useAuth();
  const displayName = user?.displayName || user?.username || 'Lead Tech';

  // Fetch dashboard data
  const { data: dashboardData, isLoading } = useLeadTechDashboard();
  const { data: unreadCounts } = useUnreadCounts();

  const unreadMessages = unreadCounts?.messages || 0;

  // Default values when loading
  const pendingApprovalsCount = dashboardData?.pendingApprovalsCount || 0;
  const todaysScheduledCount = dashboardData?.todaysScheduledCount || 0;
  const teamMetrics = dashboardData?.teamMetrics || { inProgress: 0, completed: 0, pending: 0, scheduled: 0 };
  const pendingApprovals = dashboardData?.pendingApprovals || [];
  const todaysSchedule = dashboardData?.todaysSchedule || [];
  const bomReadiness = dashboardData?.bomReadiness || [];
  const recentTeamActivity = dashboardData?.recentTeamActivity || [];

  return (
    <div className="space-y-6">
      {/* Welcome Header - Teal/Cyan Theme (Lead Tech) */}
      <div className="relative overflow-hidden rounded-2xl bg-gradient-to-r from-teal-600 via-cyan-600 to-blue-600 p-6 text-white shadow-lg">
        {/* Decorative elements */}
        <div className="absolute top-0 right-0 w-64 h-64 bg-white/5 rounded-full -translate-y-32 translate-x-32" />
        <div className="absolute bottom-0 left-0 w-48 h-48 bg-white/5 rounded-full translate-y-24 -translate-x-24" />

        <div className="relative">
          <div className="flex items-center gap-3 mb-2">
            <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-white/20 backdrop-blur">
              <Users className="h-6 w-6" />
            </div>
            <div>
              <h1 className="text-2xl md:text-3xl font-bold">Welcome, {displayName}!</h1>
              <p className="text-teal-100">Lead Technician Dashboard</p>
            </div>
          </div>

          {/* Summary Stats */}
          <div className="mt-6 grid grid-cols-2 md:grid-cols-4 gap-4">
            <div className="bg-white/10 backdrop-blur rounded-xl p-3">
              <p className="text-2xl font-bold">{pendingApprovalsCount}</p>
              <p className="text-teal-200 text-sm">Pending Approvals</p>
            </div>
            <div className="bg-white/10 backdrop-blur rounded-xl p-3">
              <p className="text-2xl font-bold">{todaysScheduledCount}</p>
              <p className="text-teal-200 text-sm">Today's Schedule</p>
            </div>
            <div className="bg-white/10 backdrop-blur rounded-xl p-3">
              <p className="text-2xl font-bold">{teamMetrics.inProgress}</p>
              <p className="text-teal-200 text-sm">In Progress</p>
            </div>
            <div className="bg-white/10 backdrop-blur rounded-xl p-3">
              <p className="text-2xl font-bold">{unreadMessages}</p>
              <p className="text-teal-200 text-sm">Messages</p>
            </div>
          </div>
        </div>
      </div>

      {/* Quick Actions */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <Link
          to="/work-orders?status=pending_approval"
          className="flex items-center gap-4 p-4 rounded-xl bg-gradient-to-r from-teal-500 to-cyan-500 text-white shadow-md hover:shadow-lg transition-all hover:-translate-y-0.5"
        >
          <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-white/20">
            <Shield className="h-6 w-6" />
          </div>
          <div>
            <p className="font-semibold text-lg">Review Approvals</p>
            <p className="text-teal-100 text-sm">Pending reviews</p>
          </div>
          {pendingApprovalsCount > 0 && (
            <Badge className="bg-white text-teal-600 ml-auto">
              {pendingApprovalsCount}
            </Badge>
          )}
        </Link>

        <Link
          to="/work-orders?action=schedule"
          className="flex items-center gap-4 p-4 rounded-xl bg-gradient-to-r from-cyan-500 to-blue-500 text-white shadow-md hover:shadow-lg transition-all hover:-translate-y-0.5"
        >
          <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-white/20">
            <Calendar className="h-6 w-6" />
          </div>
          <div>
            <p className="font-semibold text-lg">Schedule Work</p>
            <p className="text-cyan-100 text-sm">Plan assignments</p>
          </div>
        </Link>

        <Link
          to="/service-shops"
          className="flex items-center gap-4 p-4 rounded-xl bg-gradient-to-r from-blue-500 to-indigo-500 text-white shadow-md hover:shadow-lg transition-all hover:-translate-y-0.5"
        >
          <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-white/20">
            <Users className="h-6 w-6" />
          </div>
          <div>
            <p className="font-semibold text-lg">Service Shops</p>
            <p className="text-blue-100 text-sm">View team locations</p>
          </div>
        </Link>

        <Link
          to="/parts-catalog"
          className="flex items-center gap-4 p-4 rounded-xl bg-gradient-to-r from-indigo-500 to-purple-500 text-white shadow-md hover:shadow-lg transition-all hover:-translate-y-0.5"
        >
          <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-white/20">
            <Package className="h-6 w-6" />
          </div>
          <div>
            <p className="font-semibold text-lg">Parts Status</p>
            <p className="text-indigo-100 text-sm">Check inventory</p>
          </div>
        </Link>
      </div>

      {/* Two Column Layout: Pending Approvals + Today's Schedule */}
      <div className="grid gap-6 lg:grid-cols-2">
        {/* Pending Approvals */}
        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <CardTitle className="flex items-center gap-2">
                <ClipboardCheck className="h-5 w-5 text-teal-600" />
                Pending Approvals
              </CardTitle>
              {pendingApprovalsCount > 0 && (
                <Badge className="bg-teal-100 text-teal-800">
                  {pendingApprovalsCount} pending
                </Badge>
              )}
            </div>
          </CardHeader>
          <CardContent>
            {isLoading ? (
              <div className="text-center py-8 text-gray-500">Loading...</div>
            ) : pendingApprovals.length === 0 ? (
              <div className="text-center py-8 text-gray-500">
                <CheckCircle2 className="h-12 w-12 mx-auto mb-3 text-green-500" />
                <p className="font-medium">All caught up!</p>
                <p className="text-sm">No pending approvals at this time.</p>
              </div>
            ) : (
              <div className="space-y-3">
                {pendingApprovals.slice(0, 5).map((approval) => (
                  <Link
                    key={approval.id}
                    to={`/work-orders/${approval.workOrderId}`}
                    className="block p-3 rounded-lg border border-gray-100 hover:border-teal-200 hover:bg-teal-50/50 transition-colors"
                  >
                    <div className="flex items-start justify-between">
                      <div>
                        <p className="font-medium text-gray-900">{approval.title}</p>
                        <p className="text-sm text-gray-500">{approval.schoolName}</p>
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
                      Requested {formatDistanceToNow(new Date(approval.requestedAt), { addSuffix: true })}
                    </p>
                  </Link>
                ))}
                {pendingApprovalsCount > 5 && (
                  <Link to="/work-orders?status=pending_approval">
                    <Button variant="outline" className="w-full mt-2">
                      View All Pending ({pendingApprovalsCount})
                      <ArrowRight className="h-4 w-4 ml-2" />
                    </Button>
                  </Link>
                )}
              </div>
            )}
          </CardContent>
        </Card>

        {/* Today's Schedule */}
        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <CardTitle className="flex items-center gap-2">
                <Calendar className="h-5 w-5 text-cyan-600" />
                Today's Scheduled Work
              </CardTitle>
              {todaysScheduledCount > 0 && (
                <Badge className="bg-cyan-100 text-cyan-800">
                  {todaysScheduledCount} jobs
                </Badge>
              )}
            </div>
          </CardHeader>
          <CardContent>
            {isLoading ? (
              <div className="text-center py-8 text-gray-500">Loading...</div>
            ) : todaysSchedule.length === 0 ? (
              <div className="text-center py-8 text-gray-500">
                <Calendar className="h-12 w-12 mx-auto mb-3 text-gray-300" />
                <p>No scheduled work for today</p>
              </div>
            ) : (
              <div className="space-y-3">
                {todaysSchedule.slice(0, 5).map((wo) => (
                  <Link
                    key={wo.id}
                    to={`/work-orders/${wo.id}`}
                    className="block p-3 rounded-lg border border-gray-100 hover:border-cyan-200 hover:bg-cyan-50/50 transition-colors"
                  >
                    <div className="flex items-start justify-between">
                      <div>
                        <p className="font-medium text-gray-900">{wo.title}</p>
                        <p className="text-sm text-gray-500">{wo.schoolName}</p>
                      </div>
                      <Badge
                        className={cn(
                          wo.status === 'in_progress' || wo.status === 'in_repair'
                            ? 'bg-blue-100 text-blue-800'
                            : wo.status === 'completed'
                              ? 'bg-green-100 text-green-800'
                              : 'bg-gray-100 text-gray-800'
                        )}
                      >
                        {wo.status.replace('_', ' ')}
                      </Badge>
                    </div>
                    <div className="mt-2 flex items-center gap-4 text-xs text-gray-500">
                      {wo.scheduledStart && (
                        <span className="flex items-center gap-1">
                          <Clock className="h-3 w-3" />
                          {new Date(wo.scheduledStart).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
                        </span>
                      )}
                      {wo.assignedTo && (
                        <span className="flex items-center gap-1">
                          <Users className="h-3 w-3" />
                          {wo.assignedTo}
                        </span>
                      )}
                    </div>
                  </Link>
                ))}
                {todaysScheduledCount > 5 && (
                  <Link to="/work-orders?scheduled=today">
                    <Button variant="outline" className="w-full mt-2">
                      View Full Schedule ({todaysScheduledCount})
                      <ArrowRight className="h-4 w-4 ml-2" />
                    </Button>
                  </Link>
                )}
              </div>
            )}
          </CardContent>
        </Card>
      </div>

      {/* Team Work Order Metrics */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Wrench className="h-5 w-5 text-blue-600" />
            Team Work Order Status
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
            <div className="p-4 rounded-lg bg-blue-50 text-center">
              <p className="text-3xl font-bold text-blue-600">{teamMetrics.inProgress}</p>
              <p className="text-sm text-blue-700">In Progress</p>
            </div>
            <div className="p-4 rounded-lg bg-amber-50 text-center">
              <p className="text-3xl font-bold text-amber-600">{teamMetrics.pending}</p>
              <p className="text-sm text-amber-700">Pending</p>
            </div>
            <div className="p-4 rounded-lg bg-purple-50 text-center">
              <p className="text-3xl font-bold text-purple-600">{teamMetrics.scheduled}</p>
              <p className="text-sm text-purple-700">Scheduled</p>
            </div>
            <div className="p-4 rounded-lg bg-green-50 text-center">
              <p className="text-3xl font-bold text-green-600">{teamMetrics.completed}</p>
              <p className="text-sm text-green-700">Completed</p>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Two Column: BOM Readiness + Recent Activity */}
      <div className="grid gap-6 lg:grid-cols-2">
        {/* BOM Readiness */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Package className="h-5 w-5 text-purple-600" />
              Parts Readiness
            </CardTitle>
          </CardHeader>
          <CardContent>
            {isLoading ? (
              <div className="text-center py-8 text-gray-500">Loading...</div>
            ) : bomReadiness.length === 0 ? (
              <div className="text-center py-8 text-gray-500">
                <Package className="h-12 w-12 mx-auto mb-3 text-gray-300" />
                <p>No scheduled jobs with BOM requirements</p>
              </div>
            ) : (
              <div className="space-y-3">
                {bomReadiness.slice(0, 5).map((item) => (
                  <div
                    key={item.workOrderId}
                    className={cn(
                      "p-3 rounded-lg border",
                      item.isReady
                        ? "border-green-100 bg-green-50/50"
                        : "border-amber-100 bg-amber-50/50"
                    )}
                  >
                    <div className="flex items-center justify-between">
                      <div>
                        <p className="font-medium text-gray-900">{item.workOrderTitle}</p>
                        <p className="text-sm text-gray-500">
                          {item.availableParts}/{item.totalParts} parts available
                        </p>
                      </div>
                      {item.isReady ? (
                        <CheckCircle2 className="h-5 w-5 text-green-600" />
                      ) : (
                        <AlertTriangle className="h-5 w-5 text-amber-600" />
                      )}
                    </div>
                    {!item.isReady && (
                      <p className="mt-1 text-xs text-amber-600">
                        {item.missingParts} parts missing
                      </p>
                    )}
                  </div>
                ))}
              </div>
            )}
          </CardContent>
        </Card>

        {/* Recent Team Activity */}
        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <CardTitle className="flex items-center gap-2">
                <Activity className="h-5 w-5 text-emerald-600" />
                Recent Team Activity
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
            ) : recentTeamActivity.length === 0 ? (
              <div className="text-center py-8 text-gray-500">
                <Activity className="h-12 w-12 mx-auto mb-3 text-gray-300" />
                <p>No recent activity</p>
              </div>
            ) : (
              <div className="space-y-3">
                {recentTeamActivity.slice(0, 5).map((activity) => (
                  <div
                    key={activity.id}
                    className="flex items-center gap-4 p-3 rounded-lg border border-gray-100"
                  >
                    <div
                      className={cn(
                        'flex h-10 w-10 items-center justify-center rounded-full',
                        activity.entityType === 'work_order'
                          ? 'bg-blue-100'
                          : activity.entityType === 'bom_item'
                            ? 'bg-purple-100'
                            : 'bg-teal-100'
                      )}
                    >
                      {activity.entityType === 'work_order' && (
                        <Wrench className="h-5 w-5 text-blue-600" />
                      )}
                      {activity.entityType === 'bom_item' && (
                        <Package className="h-5 w-5 text-purple-600" />
                      )}
                      {activity.entityType === 'work_order_approval' && (
                        <Shield className="h-5 w-5 text-teal-600" />
                      )}
                    </div>
                    <div className="flex-1">
                      <p className="font-medium text-gray-900">{activity.description}</p>
                      <p className="text-sm text-gray-500">{activity.actorName}</p>
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
      </div>

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
            Coordinate with your team and receive updates on work order approvals and schedules.
          </p>
        </CardContent>
      </Card>
    </div>
  );
}
