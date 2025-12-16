import { Link } from 'react-router-dom';
import {
  Wrench,
  Clock,
  CheckCircle2,
  ArrowRight,
  MessageSquare,
  BookOpen,
  AlertCircle,
  MapPin,
  Calendar,
  Timer,
  TrendingUp,
} from 'lucide-react';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { useAuth } from '@/contexts/AuthContext';
import { cn } from '@/lib/utils';
import { useThreads, useUnreadCounts } from '@/hooks/useMessages';
import { formatDistanceToNow } from 'date-fns';

// Helper function to format time ago
function formatTimeAgo(date: Date): string {
  return formatDistanceToNow(date, { addSuffix: true });
}

// Work order status configuration
const workOrderStatusConfig = {
  assigned: { label: 'Assigned', color: 'text-blue-600', bgColor: 'bg-blue-100', icon: Clock },
  in_progress: { label: 'In Progress', color: 'text-purple-600', bgColor: 'bg-purple-100', icon: Wrench },
  pending_parts: { label: 'Pending Parts', color: 'text-yellow-600', bgColor: 'bg-yellow-100', icon: Timer },
  completed: { label: 'Completed', color: 'text-green-600', bgColor: 'bg-green-100', icon: CheckCircle2 },
};

// Mock data for field tech's work orders
const myWorkOrders = [
  {
    id: 'WO-2024-0892',
    title: 'Laptop screen replacement',
    status: 'in_progress',
    priority: 'high',
    schoolName: 'Greenwood Primary School',
    schoolAddress: '123 Oak Street, Nairobi',
    scheduledDate: '2024-12-14',
    scheduledTime: '09:00 AM',
    deviceTag: 'DEV-LP-0042',
    estimatedDuration: '2 hours',
  },
  {
    id: 'WO-2024-0895',
    title: 'Battery replacement - Dell Latitude',
    status: 'assigned',
    priority: 'medium',
    schoolName: 'Mombasa Secondary School',
    schoolAddress: '456 Palm Road, Mombasa',
    scheduledDate: '2024-12-14',
    scheduledTime: '02:00 PM',
    deviceTag: 'DEV-LP-0156',
    estimatedDuration: '1 hour',
  },
  {
    id: 'WO-2024-0878',
    title: 'Keyboard replacement',
    status: 'pending_parts',
    priority: 'low',
    schoolName: 'Kisumu Academy',
    schoolAddress: '789 Lake View, Kisumu',
    scheduledDate: '2024-12-16',
    deviceTag: 'DEV-LP-0089',
    estimatedDuration: '45 mins',
    partsNote: 'Waiting for keyboard delivery',
  },
  {
    id: 'WO-2024-0865',
    title: 'Motherboard diagnostic',
    status: 'completed',
    priority: 'high',
    schoolName: 'Nakuru High School',
    completedAt: '2024-12-13T16:30:00Z',
    deviceTag: 'DEV-LP-0201',
    resolution: 'Replaced faulty RAM module',
  },
];

// Mock messages
const recentMessages = [
  {
    id: 'msg-1',
    from: 'Support Team',
    subject: 'Parts ready for pickup - WO-2024-0878',
    preview: 'The keyboard for work order WO-2024-0878 is now available...',
    timestamp: '1 hour ago',
    unread: true,
  },
  {
    id: 'msg-2',
    from: 'Lead Technician',
    subject: 'Schedule change for tomorrow',
    preview: 'Please note the updated schedule for tomorrow morning...',
    timestamp: '3 hours ago',
    unread: false,
  },
];

export function FieldTechDashboard() {
  const { user } = useAuth();
  const displayName = user?.displayName || user?.username || 'Technician';

  // Fetch real message data
  const { data: threads } = useThreads({ limit: 5 });
  const { data: unreadCounts } = useUnreadCounts();

  // Calculate stats
  const todaysOrders = myWorkOrders.filter(
    (wo) => wo.status !== 'completed' && wo.scheduledDate === '2024-12-14'
  ).length;
  const inProgress = myWorkOrders.filter((wo) => wo.status === 'in_progress').length;
  const completedToday = myWorkOrders.filter((wo) => wo.status === 'completed').length;
  const unreadMessages = unreadCounts?.total || recentMessages.filter((m) => m.unread).length;

  // Convert threads to message display format
  const displayMessages =
    threads?.items?.slice(0, 3).map((thread) => ({
      id: thread.id,
      from: thread.createdByName || 'Support Team',
      subject: thread.subject,
      preview: thread.lastMessage?.content || '',
      timestamp: formatTimeAgo(new Date(thread.lastMessageAt || thread.updatedAt)),
      unread: thread.unreadCountSchool > 0,
    })) || recentMessages;

  // Get active work orders (not completed)
  const activeWorkOrders = myWorkOrders.filter((wo) => wo.status !== 'completed');
  const completedWorkOrders = myWorkOrders.filter((wo) => wo.status === 'completed');

  return (
    <div className="space-y-6">
      {/* Welcome Header for Field Tech */}
      <div className="relative overflow-hidden rounded-2xl bg-gradient-to-r from-blue-600 via-indigo-600 to-purple-600 p-6 text-white shadow-lg">
        <div className="absolute top-0 right-0 w-64 h-64 bg-white/5 rounded-full -translate-y-32 translate-x-32" />
        <div className="absolute bottom-0 left-0 w-48 h-48 bg-white/5 rounded-full translate-y-24 -translate-x-24" />

        <div className="relative">
          <div className="flex items-center gap-3 mb-2">
            <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-white/20 backdrop-blur">
              <Wrench className="h-6 w-6" />
            </div>
            <div>
              <h1 className="text-2xl md:text-3xl font-bold">Welcome, {displayName}!</h1>
              <p className="text-blue-100">Field Technician</p>
            </div>
          </div>

          <div className="mt-6 grid grid-cols-2 md:grid-cols-4 gap-4">
            <div className="bg-white/10 backdrop-blur rounded-xl p-3">
              <p className="text-2xl font-bold">{todaysOrders}</p>
              <p className="text-blue-200 text-sm">Today's Jobs</p>
            </div>
            <div className="bg-white/10 backdrop-blur rounded-xl p-3">
              <p className="text-2xl font-bold">{inProgress}</p>
              <p className="text-blue-200 text-sm">In Progress</p>
            </div>
            <div className="bg-white/10 backdrop-blur rounded-xl p-3">
              <p className="text-2xl font-bold">{completedToday}</p>
              <p className="text-blue-200 text-sm">Completed</p>
            </div>
            <div className="bg-white/10 backdrop-blur rounded-xl p-3">
              <p className="text-2xl font-bold">{unreadMessages}</p>
              <p className="text-blue-200 text-sm">Messages</p>
            </div>
          </div>
        </div>
      </div>

      {/* Quick Actions for Field Tech */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <Link
          to="/work-orders?assignee=me"
          className="flex items-center gap-4 p-4 rounded-xl bg-gradient-to-r from-blue-500 to-cyan-500 text-white shadow-md hover:shadow-lg transition-all hover:-translate-y-0.5"
        >
          <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-white/20">
            <Wrench className="h-6 w-6" />
          </div>
          <div>
            <p className="font-semibold text-lg">My Work Orders</p>
            <p className="text-blue-100 text-sm">View assigned jobs</p>
          </div>
          {activeWorkOrders.length > 0 && (
            <Badge className="ml-auto bg-white text-blue-600">{activeWorkOrders.length}</Badge>
          )}
        </Link>

        <Link
          to="/messages"
          className="flex items-center gap-4 p-4 rounded-xl bg-gradient-to-r from-purple-500 to-pink-500 text-white shadow-md hover:shadow-lg transition-all hover:-translate-y-0.5"
        >
          <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-white/20">
            <MessageSquare className="h-6 w-6" />
          </div>
          <div>
            <p className="font-semibold text-lg">Messages</p>
            <p className="text-purple-100 text-sm">Contact support team</p>
          </div>
          {unreadMessages > 0 && (
            <Badge className="ml-auto bg-white text-purple-600">{unreadMessages} new</Badge>
          )}
        </Link>

        <Link
          to="/knowledge-base"
          className="flex items-center gap-4 p-4 rounded-xl bg-gradient-to-r from-emerald-500 to-teal-500 text-white shadow-md hover:shadow-lg transition-all hover:-translate-y-0.5"
        >
          <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-white/20">
            <BookOpen className="h-6 w-6" />
          </div>
          <div>
            <p className="font-semibold text-lg">Knowledge Base</p>
            <p className="text-emerald-100 text-sm">Repair guides & docs</p>
          </div>
          <ArrowRight className="ml-auto h-5 w-5" />
        </Link>
      </div>

      {/* Today's Schedule */}
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <CardTitle className="flex items-center gap-2">
              <Calendar className="h-5 w-5 text-blue-600" />
              Today's Schedule
            </CardTitle>
            <Link to="/work-orders?assignee=me">
              <Button variant="outline" size="sm">
                View All
                <ArrowRight className="ml-2 h-4 w-4" />
              </Button>
            </Link>
          </div>
        </CardHeader>
        <CardContent>
          {activeWorkOrders.length === 0 ? (
            <div className="text-center py-8 text-gray-500">
              <CheckCircle2 className="h-12 w-12 mx-auto mb-3 text-green-500" />
              <p className="font-medium">All caught up!</p>
              <p className="text-sm">No pending work orders assigned to you.</p>
            </div>
          ) : (
            <div className="space-y-4">
              {activeWorkOrders.map((wo) => {
                const status = workOrderStatusConfig[wo.status as keyof typeof workOrderStatusConfig];
                const StatusIcon = status?.icon || Clock;

                return (
                  <Link
                    key={wo.id}
                    to={`/work-orders/${wo.id}`}
                    className="block p-4 rounded-lg border border-gray-100 hover:border-gray-200 hover:bg-gray-50/50 transition-colors"
                  >
                    <div className="flex items-start justify-between gap-4">
                      <div className="flex-1 min-w-0">
                        <div className="flex items-center gap-2 mb-1">
                          <span className="text-sm font-mono text-gray-500">{wo.id}</span>
                          <Badge
                            className={cn(
                              wo.priority === 'high'
                                ? 'bg-red-100 text-red-800'
                                : wo.priority === 'medium'
                                ? 'bg-yellow-100 text-yellow-800'
                                : 'bg-green-100 text-green-800'
                            )}
                          >
                            {wo.priority}
                          </Badge>
                          <Badge className={cn(status?.bgColor, status?.color)}>
                            <StatusIcon className="h-3 w-3 mr-1" />
                            {status?.label}
                          </Badge>
                        </div>

                        <h4 className="font-medium text-gray-900">{wo.title}</h4>

                        <div className="mt-2 space-y-1">
                          <div className="flex items-center gap-2 text-sm text-gray-600">
                            <MapPin className="h-4 w-4 text-gray-400" />
                            <span>{wo.schoolName}</span>
                          </div>
                          {wo.scheduledTime && (
                            <div className="flex items-center gap-2 text-sm text-gray-600">
                              <Clock className="h-4 w-4 text-gray-400" />
                              <span>
                                {wo.scheduledTime} • Est. {wo.estimatedDuration}
                              </span>
                            </div>
                          )}
                          {wo.deviceTag && (
                            <div className="flex items-center gap-2 text-sm text-gray-500">
                              <span className="font-mono text-xs bg-gray-100 px-2 py-0.5 rounded">
                                {wo.deviceTag}
                              </span>
                            </div>
                          )}
                        </div>

                        {wo.status === 'pending_parts' && wo.partsNote && (
                          <div className="mt-2 p-2 bg-yellow-50 border border-yellow-100 rounded text-sm text-yellow-800">
                            <AlertCircle className="h-4 w-4 inline mr-1" />
                            {wo.partsNote}
                          </div>
                        )}
                      </div>
                    </div>
                  </Link>
                );
              })}
            </div>
          )}
        </CardContent>
      </Card>

      {/* Two Column: Messages and Recent Completions */}
      <div className="grid gap-6 lg:grid-cols-2">
        {/* Messages */}
        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <CardTitle className="flex items-center gap-2">
                <MessageSquare className="h-5 w-5 text-purple-600" />
                Messages
              </CardTitle>
              {unreadMessages > 0 && (
                <Badge className="bg-purple-100 text-purple-800">{unreadMessages} unread</Badge>
              )}
            </div>
          </CardHeader>
          <CardContent>
            {displayMessages.length === 0 ? (
              <div className="text-center py-8 text-gray-500">
                <MessageSquare className="h-12 w-12 mx-auto mb-3 text-gray-300" />
                <p>No messages yet</p>
              </div>
            ) : (
              <div className="space-y-3">
                {displayMessages.map((message) => (
                  <Link
                    key={message.id}
                    to={`/messages/${message.id}`}
                    className={cn(
                      'block p-3 rounded-lg border transition-colors',
                      message.unread
                        ? 'border-purple-200 bg-purple-50/50 hover:bg-purple-50'
                        : 'border-gray-100 hover:border-gray-200 hover:bg-gray-50'
                    )}
                  >
                    <div className="flex items-start justify-between gap-3">
                      <div className="min-w-0 flex-1">
                        <div className="flex items-center gap-2">
                          <span className="font-medium text-gray-900">{message.from}</span>
                          {message.unread && <span className="w-2 h-2 bg-purple-600 rounded-full" />}
                        </div>
                        <p className="font-medium text-sm text-gray-700 truncate">{message.subject}</p>
                        <p className="text-sm text-gray-500 truncate">{message.preview}</p>
                      </div>
                      <span className="text-xs text-gray-400 whitespace-nowrap">{message.timestamp}</span>
                    </div>
                  </Link>
                ))}
                <Link to="/messages">
                  <Button variant="outline" className="w-full mt-2">
                    View All Messages
                    <ArrowRight className="h-4 w-4 ml-2" />
                  </Button>
                </Link>
              </div>
            )}
          </CardContent>
        </Card>

        {/* Recent Completions */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <TrendingUp className="h-5 w-5 text-green-600" />
              Recent Completions
            </CardTitle>
          </CardHeader>
          <CardContent>
            {completedWorkOrders.length === 0 ? (
              <div className="text-center py-8 text-gray-500">
                <CheckCircle2 className="h-12 w-12 mx-auto mb-3 text-gray-300" />
                <p>No completed work orders yet</p>
              </div>
            ) : (
              <div className="space-y-3">
                {completedWorkOrders.slice(0, 3).map((wo) => (
                  <div
                    key={wo.id}
                    className="p-3 rounded-lg border border-green-100 bg-green-50/50"
                  >
                    <div className="flex items-start justify-between">
                      <div>
                        <div className="flex items-center gap-2 mb-1">
                          <span className="text-sm font-mono text-gray-500">{wo.id}</span>
                          <Badge className="bg-green-100 text-green-800">
                            <CheckCircle2 className="h-3 w-3 mr-1" />
                            Completed
                          </Badge>
                        </div>
                        <p className="font-medium text-gray-900">{wo.title}</p>
                        <p className="text-sm text-gray-600">{wo.schoolName}</p>
                        {wo.resolution && (
                          <p className="text-xs text-green-700 mt-1">
                            Resolution: {wo.resolution}
                          </p>
                        )}
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </CardContent>
        </Card>
      </div>

      {/* Knowledge Base Quick Access */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <BookOpen className="h-5 w-5 text-emerald-600" />
            Knowledge Base
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex items-center justify-between p-4 bg-gradient-to-r from-emerald-50 to-teal-50 rounded-lg">
            <div>
              <p className="font-medium text-gray-900">Need repair guidance?</p>
              <p className="text-sm text-gray-600">
                Access troubleshooting guides, repair manuals, and technical documentation.
              </p>
            </div>
            <Link to="/knowledge-base">
              <Button className="bg-emerald-600 hover:bg-emerald-700">
                <BookOpen className="h-4 w-4 mr-2" />
                Open KB
              </Button>
            </Link>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
