import { Link } from 'react-router-dom';
import {
  AlertTriangle,
  MessageCircle,
  Wrench,
  Clock,
  CheckCircle2,
  ArrowRight,
  Activity,
  Headphones,
  FileText,
  Users,
  AlertOctagon,
  Phone,
  Search,
} from 'lucide-react';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { useAuth } from '@/contexts/AuthContext';
import { cn } from '@/lib/utils';
import { useSupportAgentDashboard } from '@/hooks/useSupportAgentDashboard';
import { useUnreadCounts } from '@/hooks/useMessages';
import { formatDistanceToNow } from 'date-fns';

export function SupportAgentDashboard() {
  const { user } = useAuth();
  const displayName = user?.displayName || user?.username || 'Support Agent';

  // Fetch dashboard data
  const { data: dashboardData, isLoading } = useSupportAgentDashboard();
  const { data: unreadCounts } = useUnreadCounts();

  const unreadMessages = unreadCounts?.messages || dashboardData?.unreadMessagesCount || 0;

  // Default values when loading
  const openIncidentsCount = dashboardData?.openIncidentsCount || 0;
  const waitingChatsCount = dashboardData?.waitingChatsCount || 0;
  const activeChatsCount = dashboardData?.activeChatsCount || 0;
  const activeWorkOrders = dashboardData?.activeWorkOrders || 0;
  const incidentMetrics = dashboardData?.incidentMetrics || { open: 0, inProgress: 0, resolved: 0, slaBreached: 0 };
  const incidentQueue = dashboardData?.incidentQueue || [];
  const chatQueue = dashboardData?.chatQueue || [];
  const workOrderQueue = dashboardData?.workOrderQueue || [];
  const recentActivity = dashboardData?.recentActivity || [];

  // Format wait time as human readable
  const formatWaitTime = (seconds: number) => {
    if (seconds < 60) return `${seconds}s`;
    if (seconds < 3600) return `${Math.floor(seconds / 60)}m`;
    return `${Math.floor(seconds / 3600)}h ${Math.floor((seconds % 3600) / 60)}m`;
  };

  return (
    <div className="space-y-6">
      {/* Welcome Header - Amber/Orange Theme (Support Agent - Incident-Focused) */}
      <div className="relative overflow-hidden rounded-2xl bg-gradient-to-r from-amber-600 via-orange-600 to-red-600 p-6 text-white shadow-lg">
        {/* Decorative elements */}
        <div className="absolute top-0 right-0 w-64 h-64 bg-white/5 rounded-full -translate-y-32 translate-x-32" />
        <div className="absolute bottom-0 left-0 w-48 h-48 bg-white/5 rounded-full translate-y-24 -translate-x-24" />

        <div className="relative">
          <div className="flex items-center gap-3 mb-2">
            <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-white/20 backdrop-blur">
              <Headphones className="h-6 w-6" />
            </div>
            <div>
              <h1 className="text-2xl md:text-3xl font-bold">Welcome, {displayName}!</h1>
              <p className="text-amber-100">Support Agent Dashboard</p>
            </div>
          </div>

          {/* Summary Stats */}
          <div className="mt-6 grid grid-cols-2 md:grid-cols-5 gap-4">
            <div className="bg-white/10 backdrop-blur rounded-xl p-3">
              <p className="text-2xl font-bold">{openIncidentsCount}</p>
              <p className="text-amber-200 text-sm">Open Incidents</p>
            </div>
            <div className="bg-white/10 backdrop-blur rounded-xl p-3">
              <p className="text-2xl font-bold">{waitingChatsCount}</p>
              <p className="text-amber-200 text-sm">Waiting Chats</p>
            </div>
            <div className="bg-white/10 backdrop-blur rounded-xl p-3">
              <p className="text-2xl font-bold">{activeChatsCount}</p>
              <p className="text-amber-200 text-sm">Active Chats</p>
            </div>
            <div className="bg-white/10 backdrop-blur rounded-xl p-3">
              <p className="text-2xl font-bold">{activeWorkOrders}</p>
              <p className="text-amber-200 text-sm">Work Orders</p>
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
          to="/incidents/new"
          className="flex items-center gap-4 p-4 rounded-xl bg-gradient-to-r from-amber-500 to-orange-500 text-white shadow-md hover:shadow-lg transition-all hover:-translate-y-0.5"
        >
          <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-white/20">
            <FileText className="h-6 w-6" />
          </div>
          <div>
            <p className="font-semibold text-lg">New Incident</p>
            <p className="text-amber-100 text-sm">Log a ticket</p>
          </div>
        </Link>

        <Link
          to="/livechat"
          className="flex items-center gap-4 p-4 rounded-xl bg-gradient-to-r from-orange-500 to-red-500 text-white shadow-md hover:shadow-lg transition-all hover:-translate-y-0.5"
        >
          <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-white/20">
            <MessageCircle className="h-6 w-6" />
          </div>
          <div>
            <p className="font-semibold text-lg">Live Chat</p>
            <p className="text-orange-100 text-sm">Accept chats</p>
          </div>
          {waitingChatsCount > 0 && (
            <Badge className="bg-white text-orange-600 ml-auto">
              {waitingChatsCount}
            </Badge>
          )}
        </Link>

        <Link
          to="/devices"
          className="flex items-center gap-4 p-4 rounded-xl bg-gradient-to-r from-red-500 to-rose-500 text-white shadow-md hover:shadow-lg transition-all hover:-translate-y-0.5"
        >
          <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-white/20">
            <Search className="h-6 w-6" />
          </div>
          <div>
            <p className="font-semibold text-lg">Device Lookup</p>
            <p className="text-red-100 text-sm">Search devices</p>
          </div>
        </Link>

        <Link
          to="/messages"
          className="flex items-center gap-4 p-4 rounded-xl bg-gradient-to-r from-rose-500 to-pink-500 text-white shadow-md hover:shadow-lg transition-all hover:-translate-y-0.5"
        >
          <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-white/20">
            <Phone className="h-6 w-6" />
          </div>
          <div>
            <p className="font-semibold text-lg">Messages</p>
            <p className="text-rose-100 text-sm">School comms</p>
          </div>
          {unreadMessages > 0 && (
            <Badge className="bg-white text-rose-600 ml-auto">
              {unreadMessages}
            </Badge>
          )}
        </Link>
      </div>

      {/* Incident Metrics */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <AlertTriangle className="h-5 w-5 text-amber-600" />
            Incident Status Overview
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
            <div className="p-4 rounded-lg bg-amber-50 text-center">
              <p className="text-3xl font-bold text-amber-600">{incidentMetrics.open}</p>
              <p className="text-sm text-amber-700">Open</p>
            </div>
            <div className="p-4 rounded-lg bg-blue-50 text-center">
              <p className="text-3xl font-bold text-blue-600">{incidentMetrics.inProgress}</p>
              <p className="text-sm text-blue-700">In Progress</p>
            </div>
            <div className="p-4 rounded-lg bg-green-50 text-center">
              <p className="text-3xl font-bold text-green-600">{incidentMetrics.resolved}</p>
              <p className="text-sm text-green-700">Resolved</p>
            </div>
            <div className="p-4 rounded-lg bg-red-50 text-center">
              <p className="text-3xl font-bold text-red-600">{incidentMetrics.slaBreached}</p>
              <p className="text-sm text-red-700">SLA Breached</p>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Two Column Layout: Incident Queue + Chat Queue */}
      <div className="grid gap-6 lg:grid-cols-2">
        {/* Incident Queue */}
        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <CardTitle className="flex items-center gap-2">
                <AlertTriangle className="h-5 w-5 text-amber-600" />
                Incident Queue
              </CardTitle>
              {openIncidentsCount > 0 && (
                <Badge className="bg-amber-100 text-amber-800">
                  {openIncidentsCount} open
                </Badge>
              )}
            </div>
          </CardHeader>
          <CardContent>
            {isLoading ? (
              <div className="text-center py-8 text-gray-500">Loading...</div>
            ) : incidentQueue.length === 0 ? (
              <div className="text-center py-8 text-gray-500">
                <CheckCircle2 className="h-12 w-12 mx-auto mb-3 text-green-500" />
                <p className="font-medium">All caught up!</p>
                <p className="text-sm">No open incidents at this time.</p>
              </div>
            ) : (
              <div className="space-y-3">
                {incidentQueue.slice(0, 5).map((incident) => (
                  <Link
                    key={incident.id}
                    to={`/incidents/${incident.id}`}
                    className={cn(
                      "block p-3 rounded-lg border transition-colors",
                      incident.slaBreached
                        ? "border-red-200 bg-red-50/50 hover:border-red-300"
                        : "border-gray-100 hover:border-amber-200 hover:bg-amber-50/50"
                    )}
                  >
                    <div className="flex items-start justify-between">
                      <div className="flex-1 min-w-0">
                        <div className="flex items-center gap-2">
                          {incident.slaBreached && (
                            <AlertOctagon className="h-4 w-4 text-red-500 flex-shrink-0" />
                          )}
                          <p className="font-medium text-gray-900 truncate">{incident.title}</p>
                        </div>
                        <p className="text-sm text-gray-500">{incident.schoolName}</p>
                      </div>
                      <Badge
                        className={cn(
                          "flex-shrink-0 ml-2",
                          incident.severity === 'critical'
                            ? 'bg-red-100 text-red-800'
                            : incident.severity === 'high'
                              ? 'bg-orange-100 text-orange-800'
                              : incident.severity === 'medium'
                                ? 'bg-yellow-100 text-yellow-800'
                                : 'bg-green-100 text-green-800'
                        )}
                      >
                        {incident.severity}
                      </Badge>
                    </div>
                    <div className="mt-2 flex items-center gap-4 text-xs text-gray-500">
                      <span className="flex items-center gap-1">
                        <Clock className="h-3 w-3" />
                        {formatDistanceToNow(new Date(incident.createdAt), { addSuffix: true })}
                      </span>
                      <span className="capitalize">{incident.category}</span>
                    </div>
                  </Link>
                ))}
                {incidentQueue.length > 5 && (
                  <Link to="/incidents">
                    <Button variant="outline" className="w-full mt-2">
                      View All Incidents ({openIncidentsCount})
                      <ArrowRight className="h-4 w-4 ml-2" />
                    </Button>
                  </Link>
                )}
              </div>
            )}
          </CardContent>
        </Card>

        {/* Chat Queue */}
        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <CardTitle className="flex items-center gap-2">
                <MessageCircle className="h-5 w-5 text-orange-600" />
                Live Chat Queue
              </CardTitle>
              {waitingChatsCount > 0 && (
                <Badge className="bg-orange-100 text-orange-800 animate-pulse">
                  {waitingChatsCount} waiting
                </Badge>
              )}
            </div>
          </CardHeader>
          <CardContent>
            {isLoading ? (
              <div className="text-center py-8 text-gray-500">Loading...</div>
            ) : chatQueue.length === 0 ? (
              <div className="text-center py-8 text-gray-500">
                <MessageCircle className="h-12 w-12 mx-auto mb-3 text-gray-300" />
                <p>No active chat sessions</p>
              </div>
            ) : (
              <div className="space-y-3">
                {chatQueue.slice(0, 5).map((chat) => (
                  <Link
                    key={chat.id}
                    to={`/livechat/${chat.id}`}
                    className={cn(
                      "block p-3 rounded-lg border transition-colors",
                      chat.status === 'waiting'
                        ? "border-orange-200 bg-orange-50/50 hover:border-orange-300"
                        : "border-gray-100 hover:border-green-200 hover:bg-green-50/50"
                    )}
                  >
                    <div className="flex items-start justify-between">
                      <div>
                        <p className="font-medium text-gray-900">{chat.contactName || 'Unknown Contact'}</p>
                        <p className="text-sm text-gray-500">{chat.schoolName}</p>
                      </div>
                      <Badge
                        className={cn(
                          chat.status === 'waiting'
                            ? 'bg-orange-100 text-orange-800'
                            : 'bg-green-100 text-green-800'
                        )}
                      >
                        {chat.status}
                      </Badge>
                    </div>
                    <div className="mt-2 flex items-center gap-4 text-xs text-gray-500">
                      <span className="flex items-center gap-1">
                        <Clock className="h-3 w-3" />
                        {formatWaitTime(chat.waitTimeSeconds)} wait
                      </span>
                      {chat.assignedAgentName && (
                        <span className="flex items-center gap-1">
                          <Users className="h-3 w-3" />
                          {chat.assignedAgentName}
                        </span>
                      )}
                      {chat.queuePosition && (
                        <span>Queue #{chat.queuePosition}</span>
                      )}
                    </div>
                  </Link>
                ))}
                <Link to="/livechat">
                  <Button variant="outline" className="w-full mt-2">
                    Open Live Chat Console
                    <ArrowRight className="h-4 w-4 ml-2" />
                  </Button>
                </Link>
              </div>
            )}
          </CardContent>
        </Card>
      </div>

      {/* Two Column Layout: Work Orders + Recent Activity */}
      <div className="grid gap-6 lg:grid-cols-2">
        {/* Active Work Orders */}
        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <CardTitle className="flex items-center gap-2">
                <Wrench className="h-5 w-5 text-blue-600" />
                Active Work Orders
              </CardTitle>
              {activeWorkOrders > 0 && (
                <Badge className="bg-blue-100 text-blue-800">
                  {activeWorkOrders} active
                </Badge>
              )}
            </div>
          </CardHeader>
          <CardContent>
            {isLoading ? (
              <div className="text-center py-8 text-gray-500">Loading...</div>
            ) : workOrderQueue.length === 0 ? (
              <div className="text-center py-8 text-gray-500">
                <Wrench className="h-12 w-12 mx-auto mb-3 text-gray-300" />
                <p>No active work orders</p>
              </div>
            ) : (
              <div className="space-y-3">
                {workOrderQueue.slice(0, 5).map((wo) => (
                  <Link
                    key={wo.id}
                    to={`/work-orders/${wo.id}`}
                    className="block p-3 rounded-lg border border-gray-100 hover:border-blue-200 hover:bg-blue-50/50 transition-colors"
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
                            : wo.status === 'assigned'
                              ? 'bg-purple-100 text-purple-800'
                              : 'bg-gray-100 text-gray-800'
                        )}
                      >
                        {wo.status.replace('_', ' ')}
                      </Badge>
                    </div>
                    <div className="mt-2 flex items-center gap-4 text-xs text-gray-500">
                      {wo.taskType && <span className="capitalize">{wo.taskType}</span>}
                      {wo.assignedTo && (
                        <span className="flex items-center gap-1">
                          <Users className="h-3 w-3" />
                          {wo.assignedTo}
                        </span>
                      )}
                    </div>
                  </Link>
                ))}
                {workOrderQueue.length > 5 && (
                  <Link to="/work-orders">
                    <Button variant="outline" className="w-full mt-2">
                      View All Work Orders
                      <ArrowRight className="h-4 w-4 ml-2" />
                    </Button>
                  </Link>
                )}
              </div>
            )}
          </CardContent>
        </Card>

        {/* Recent Activity */}
        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <CardTitle className="flex items-center gap-2">
                <Activity className="h-5 w-5 text-emerald-600" />
                Recent Activity
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
                <Activity className="h-12 w-12 mx-auto mb-3 text-gray-300" />
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
                        activity.entityType === 'incident'
                          ? 'bg-amber-100'
                          : activity.entityType === 'work_order'
                            ? 'bg-blue-100'
                            : activity.entityType === 'chat_session'
                              ? 'bg-orange-100'
                              : 'bg-gray-100'
                      )}
                    >
                      {activity.entityType === 'incident' && (
                        <AlertTriangle className="h-5 w-5 text-amber-600" />
                      )}
                      {activity.entityType === 'work_order' && (
                        <Wrench className="h-5 w-5 text-blue-600" />
                      )}
                      {activity.entityType === 'chat_session' && (
                        <MessageCircle className="h-5 w-5 text-orange-600" />
                      )}
                      {activity.entityType === 'message' && (
                        <Phone className="h-5 w-5 text-gray-600" />
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
    </div>
  );
}
