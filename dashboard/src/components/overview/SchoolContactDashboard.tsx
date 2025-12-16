import { Link } from 'react-router-dom';
import {
  AlertTriangle,
  Clock,
  CheckCircle2,
  ArrowRight,
  MessageSquare,
  Plus,
  School,
  TrendingUp,
  Timer,
  Wrench,
  Headphones,
} from 'lucide-react';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { useAuth } from '@/contexts/AuthContext';
import { cn } from '@/lib/utils';
import { ChatWidget } from '@/components/livechat';
import { EdTechProfileCard } from '@/components/edtech';
import { useThreads, useUnreadCounts } from '@/hooks/useMessages';
import { formatDistanceToNow } from 'date-fns';

// Helper function to format time ago
function formatTimeAgo(date: Date): string {
  return formatDistanceToNow(date, { addSuffix: true });
}

// Incident status configuration with timeline info
const incidentStatusFlow = [
  { status: 'new', label: 'Reported', icon: AlertTriangle, color: 'text-blue-600', bgColor: 'bg-blue-100' },
  { status: 'acknowledged', label: 'Acknowledged', icon: CheckCircle2, color: 'text-yellow-600', bgColor: 'bg-yellow-100' },
  { status: 'in_progress', label: 'In Progress', icon: Wrench, color: 'text-purple-600', bgColor: 'bg-purple-100' },
  { status: 'resolved', label: 'Resolved', icon: CheckCircle2, color: 'text-green-600', bgColor: 'bg-green-100' },
];

// Mock data for school contact's incidents
const myIncidents = [
  {
    id: 'INC-2024-0156',
    title: 'Laptop not turning on - Lab Room 3',
    status: 'in_progress',
    priority: 'high',
    createdAt: '2024-12-10T09:30:00Z',
    acknowledgedAt: '2024-12-10T10:15:00Z',
    inProgressAt: '2024-12-11T08:00:00Z',
    assignedTo: 'James Wilson',
    lastUpdate: 'Technician dispatched, expected arrival tomorrow',
    estimatedResolution: '2024-12-15',
  },
  {
    id: 'INC-2024-0142',
    title: 'Projector display issues - Main Hall',
    status: 'acknowledged',
    priority: 'medium',
    createdAt: '2024-12-08T14:20:00Z',
    acknowledgedAt: '2024-12-08T15:00:00Z',
    assignedTo: 'Support Team',
    lastUpdate: 'Reviewing issue, will assign technician soon',
  },
  {
    id: 'INC-2024-0128',
    title: 'Network connectivity issue',
    status: 'resolved',
    priority: 'high',
    createdAt: '2024-12-01T08:00:00Z',
    acknowledgedAt: '2024-12-01T08:30:00Z',
    inProgressAt: '2024-12-01T10:00:00Z',
    resolvedAt: '2024-12-02T14:00:00Z',
    resolution: 'Router replaced, network restored',
    resolutionTime: '30 hours',
  },
];

// Mock messages from support
const recentMessages = [
  {
    id: 'msg-1',
    from: 'Support Team',
    subject: 'Update on INC-2024-0156',
    preview: 'Technician has been assigned and will arrive tomorrow morning...',
    timestamp: '2 hours ago',
    unread: true,
  },
  {
    id: 'msg-2',
    from: 'James Wilson',
    subject: 'Parts ordered for your laptop repair',
    preview: 'The replacement screen has been ordered and should arrive by...',
    timestamp: '1 day ago',
    unread: false,
  },
];

function IncidentTimeline({ incident }: { incident: typeof myIncidents[0] }) {
  const currentStatusIndex = incidentStatusFlow.findIndex(s => s.status === incident.status);

  return (
    <div className="flex items-center gap-1 mt-3">
      {incidentStatusFlow.map((step, index) => {
        const isCompleted = index < currentStatusIndex;
        const isCurrent = index === currentStatusIndex;
        const Icon = step.icon;

        return (
          <div key={step.status} className="flex items-center">
            <div
              className={cn(
                'flex items-center justify-center w-8 h-8 rounded-full transition-all',
                isCompleted ? 'bg-green-100' : isCurrent ? step.bgColor : 'bg-gray-100'
              )}
            >
              <Icon
                className={cn(
                  'w-4 h-4',
                  isCompleted ? 'text-green-600' : isCurrent ? step.color : 'text-gray-400'
                )}
              />
            </div>
            {index < incidentStatusFlow.length - 1 && (
              <div
                className={cn(
                  'w-8 h-1 mx-1',
                  isCompleted ? 'bg-green-300' : 'bg-gray-200'
                )}
              />
            )}
          </div>
        );
      })}
    </div>
  );
}

export function SchoolContactDashboard() {
  const { user } = useAuth();
  const displayName = user?.displayName || user?.username || 'User';

  // In production, this would come from the user's profile
  const schoolName = 'Greenwood Primary School';
  const schoolId = user?.schoolId || 'demo-school-001'; // Would come from user profile

  // Fetch real message data
  const { data: threads } = useThreads({ limit: 5 });
  const { data: unreadCounts } = useUnreadCounts();

  const openIncidents = myIncidents.filter(i => i.status !== 'resolved').length;
  const unreadMessages = unreadCounts?.total || recentMessages.filter(m => m.unread).length;

  // Convert threads to message display format
  const displayMessages = threads?.items?.slice(0, 2).map((thread) => ({
    id: thread.id,
    from: thread.createdByName || 'Support Team',
    subject: thread.subject,
    preview: thread.lastMessage?.content || '',
    timestamp: formatTimeAgo(new Date(thread.lastMessageAt || thread.updatedAt)),
    unread: thread.unreadCountSchool > 0,
  })) || recentMessages;

  return (
    <div className="space-y-6">
      {/* Welcome Header for School Contact */}
      <div className="relative overflow-hidden rounded-2xl bg-gradient-to-r from-green-600 via-emerald-600 to-teal-600 p-6 text-white shadow-lg">
        <div className="absolute top-0 right-0 w-64 h-64 bg-white/5 rounded-full -translate-y-32 translate-x-32" />
        <div className="absolute bottom-0 left-0 w-48 h-48 bg-white/5 rounded-full translate-y-24 -translate-x-24" />

        <div className="relative">
          <div className="flex items-center gap-3 mb-2">
            <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-white/20 backdrop-blur">
              <School className="h-6 w-6" />
            </div>
            <div>
              <h1 className="text-2xl md:text-3xl font-bold">Welcome, {displayName}!</h1>
              <p className="text-green-100">{schoolName}</p>
            </div>
          </div>

          <div className="mt-6 grid grid-cols-2 md:grid-cols-3 gap-4">
            <div className="bg-white/10 backdrop-blur rounded-xl p-3">
              <p className="text-2xl font-bold">{openIncidents}</p>
              <p className="text-green-200 text-sm">Open Incidents</p>
            </div>
            <div className="bg-white/10 backdrop-blur rounded-xl p-3">
              <p className="text-2xl font-bold">{unreadMessages}</p>
              <p className="text-green-200 text-sm">New Messages</p>
            </div>
            <div className="bg-white/10 backdrop-blur rounded-xl p-3 hidden md:block">
              <p className="text-2xl font-bold">~24h</p>
              <p className="text-green-200 text-sm">Avg. Response Time</p>
            </div>
          </div>
        </div>
      </div>

      {/* Quick Actions - School Contact Specific */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <Link
          to="/incidents?action=create"
          className="flex items-center gap-4 p-4 rounded-xl bg-gradient-to-r from-amber-500 to-orange-500 text-white shadow-md hover:shadow-lg transition-all hover:-translate-y-0.5"
        >
          <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-white/20">
            <Plus className="h-6 w-6" />
          </div>
          <div>
            <p className="font-semibold text-lg">Report New Incident</p>
            <p className="text-amber-100 text-sm">Submit a new issue for your school</p>
          </div>
          <ArrowRight className="ml-auto h-5 w-5" />
        </Link>

        <Link
          to="/messages"
          className="flex items-center gap-4 p-4 rounded-xl bg-gradient-to-r from-blue-500 to-cyan-500 text-white shadow-md hover:shadow-lg transition-all hover:-translate-y-0.5"
        >
          <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-white/20">
            <MessageSquare className="h-6 w-6" />
          </div>
          <div>
            <p className="font-semibold text-lg">Contact Support</p>
            <p className="text-blue-100 text-sm">Send a message to the support team</p>
          </div>
          {unreadMessages > 0 && (
            <Badge className="ml-auto bg-white text-blue-600">{unreadMessages} new</Badge>
          )}
        </Link>
      </div>

      {/* EdTech Profile Assessment */}
      <EdTechProfileCard schoolId={schoolId} />

      {/* My Incidents with Lifecycle Tracking */}
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <CardTitle className="flex items-center gap-2">
              <AlertTriangle className="h-5 w-5 text-amber-600" />
              My Reported Incidents
            </CardTitle>
            <Link to="/incidents">
              <Button variant="outline" size="sm">
                View All
                <ArrowRight className="ml-2 h-4 w-4" />
              </Button>
            </Link>
          </div>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            {myIncidents.map((incident) => (
              <div
                key={incident.id}
                className="p-4 rounded-lg border border-gray-100 hover:border-gray-200 hover:bg-gray-50/50 transition-colors"
              >
                <div className="flex items-start justify-between gap-4">
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-2 mb-1">
                      <span className="text-sm font-mono text-gray-500">{incident.id}</span>
                      <Badge
                        className={cn(
                          incident.priority === 'high'
                            ? 'bg-red-100 text-red-800'
                            : incident.priority === 'medium'
                            ? 'bg-yellow-100 text-yellow-800'
                            : 'bg-green-100 text-green-800'
                        )}
                      >
                        {incident.priority}
                      </Badge>
                    </div>
                    <h4 className="font-medium text-gray-900">{incident.title}</h4>

                    {/* Status Timeline */}
                    <IncidentTimeline incident={incident} />

                    {/* Latest Update */}
                    <div className="mt-3 p-3 bg-gray-50 rounded-lg">
                      <div className="flex items-center gap-2 text-sm text-gray-600 mb-1">
                        <Clock className="h-4 w-4" />
                        <span>Latest Update</span>
                        {incident.assignedTo && (
                          <>
                            <span>•</span>
                            <span>Assigned to {incident.assignedTo}</span>
                          </>
                        )}
                      </div>
                      <p className="text-sm text-gray-700">
                        {incident.status === 'resolved'
                          ? `Resolved: ${incident.resolution}`
                          : incident.lastUpdate}
                      </p>
                      {incident.estimatedResolution && incident.status !== 'resolved' && (
                        <p className="text-xs text-gray-500 mt-1 flex items-center gap-1">
                          <Timer className="h-3 w-3" />
                          Expected resolution: {new Date(incident.estimatedResolution).toLocaleDateString()}
                        </p>
                      )}
                      {incident.resolutionTime && (
                        <p className="text-xs text-green-600 mt-1 flex items-center gap-1">
                          <TrendingUp className="h-3 w-3" />
                          Resolved in {incident.resolutionTime}
                        </p>
                      )}
                    </div>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>

      {/* Recent Messages from Support */}
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <CardTitle className="flex items-center gap-2">
              <MessageSquare className="h-5 w-5 text-blue-600" />
              Messages from Support
            </CardTitle>
            {unreadMessages > 0 && (
              <Badge className="bg-blue-100 text-blue-800">{unreadMessages} unread</Badge>
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
                      ? 'border-blue-200 bg-blue-50/50 hover:bg-blue-50'
                      : 'border-gray-100 hover:border-gray-200 hover:bg-gray-50'
                  )}
                >
                  <div className="flex items-start justify-between gap-3">
                    <div className="min-w-0 flex-1">
                      <div className="flex items-center gap-2">
                        <span className="font-medium text-gray-900">{message.from}</span>
                        {message.unread && (
                          <span className="w-2 h-2 bg-blue-600 rounded-full" />
                        )}
                      </div>
                      <p className="font-medium text-sm text-gray-700 truncate">{message.subject}</p>
                      <p className="text-sm text-gray-500 truncate">{message.preview}</p>
                    </div>
                    <span className="text-xs text-gray-400 whitespace-nowrap">{message.timestamp}</span>
                  </div>
                </Link>
              ))}
            </div>
          )}
        </CardContent>
      </Card>

      {/* Live Chat Support Card */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Headphones className="h-5 w-5 text-cyan-600" />
            Need Immediate Help?
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex items-center justify-between p-4 bg-gradient-to-r from-cyan-50 to-blue-50 rounded-lg">
            <div>
              <p className="font-medium text-gray-900">Live Chat Support</p>
              <p className="text-sm text-gray-600">
                Connect with a support agent in real-time for immediate assistance.
              </p>
            </div>
            <div className="text-cyan-600">
              <Headphones className="h-8 w-8" />
            </div>
          </div>
          <p className="text-xs text-gray-500 mt-3 text-center">
            Click the chat button in the bottom right corner to start a conversation.
          </p>
        </CardContent>
      </Card>

      {/* Chat Widget - Floating button for live chat */}
      <ChatWidget />
    </div>
  );
}
