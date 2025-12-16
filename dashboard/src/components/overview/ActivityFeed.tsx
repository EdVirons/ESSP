import { Activity, AlertTriangle, Wrench, Layers, User, Clock } from 'lucide-react';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { formatRelativeTime } from '@/lib/utils';
import type { ActivityEvent } from '@/types';

interface ActivityFeedProps {
  activities?: ActivityEvent[];
  isLoading?: boolean;
}

const activityIcons: Record<string, React.ElementType> = {
  incident: AlertTriangle,
  workorder: Wrench,
  program: Layers,
  user: User,
  default: Activity,
};

const activityColors: Record<string, string> = {
  incident: 'text-amber-500 bg-amber-50',
  workorder: 'text-blue-500 bg-blue-50',
  program: 'text-purple-500 bg-purple-50',
  user: 'text-green-500 bg-green-50',
  default: 'text-gray-500 bg-gray-50',
};

// Mock data for demo
const mockActivities: ActivityEvent[] = [
  {
    id: '1',
    type: 'incident',
    action: 'created',
    actor: 'John Doe',
    target: 'INC-2024-001',
    timestamp: new Date(Date.now() - 2 * 60 * 1000).toISOString(),
    metadata: { severity: 'high' },
  },
  {
    id: '2',
    type: 'workorder',
    action: 'status_changed',
    actor: 'Jane Smith',
    target: 'WO-2024-015',
    timestamp: new Date(Date.now() - 15 * 60 * 1000).toISOString(),
    metadata: { from: 'assigned', to: 'in_repair' },
  },
  {
    id: '3',
    type: 'program',
    action: 'phase_completed',
    actor: 'Mike Johnson',
    target: 'PRG-2024-003',
    timestamp: new Date(Date.now() - 45 * 60 * 1000).toISOString(),
    metadata: { phase: 'survey' },
  },
  {
    id: '4',
    type: 'workorder',
    action: 'deliverable_submitted',
    actor: 'Sarah Williams',
    target: 'WO-2024-012',
    timestamp: new Date(Date.now() - 2 * 60 * 60 * 1000).toISOString(),
    metadata: {},
  },
  {
    id: '5',
    type: 'incident',
    action: 'resolved',
    actor: 'Alex Brown',
    target: 'INC-2024-089',
    timestamp: new Date(Date.now() - 3 * 60 * 60 * 1000).toISOString(),
    metadata: {},
  },
];

function getActivityMessage(activity: ActivityEvent): string {
  const messages: Record<string, Record<string, string>> = {
    incident: {
      created: `created incident`,
      updated: `updated incident`,
      resolved: `resolved incident`,
      closed: `closed incident`,
    },
    workorder: {
      created: `created work order`,
      status_changed: `changed status of work order`,
      deliverable_submitted: `submitted deliverable for`,
      completed: `completed work order`,
    },
    program: {
      created: `created program`,
      phase_completed: `completed phase in program`,
      updated: `updated program`,
    },
  };

  return messages[activity.type]?.[activity.action] || `performed action on`;
}

export function ActivityFeed({ activities = mockActivities, isLoading }: ActivityFeedProps) {
  if (isLoading) {
    return (
      <Card>
        <CardHeader className="pb-2">
          <CardTitle className="flex items-center gap-2 text-base">
            <Clock className="h-5 w-5" />
            Recent Activity
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            {[1, 2, 3, 4, 5].map((i) => (
              <div key={i} className="flex items-start gap-3">
                <div className="skeleton h-8 w-8 rounded-full" />
                <div className="flex-1 space-y-2">
                  <div className="skeleton h-4 w-3/4 rounded" />
                  <div className="skeleton h-3 w-1/4 rounded" />
                </div>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card>
      <CardHeader className="pb-2">
        <CardTitle className="flex items-center gap-2 text-base">
          <Clock className="h-5 w-5" />
          Recent Activity
        </CardTitle>
      </CardHeader>
      <CardContent>
        <div className="space-y-4">
          {activities.map((activity) => {
            const Icon = activityIcons[activity.type] || activityIcons.default;
            const colorClass = activityColors[activity.type] || activityColors.default;

            return (
              <div key={activity.id} className="flex items-start gap-3">
                <div className={`rounded-full p-2 ${colorClass}`}>
                  <Icon className="h-4 w-4" />
                </div>
                <div className="flex-1 min-w-0">
                  <p className="text-sm text-gray-900">
                    <span className="font-medium">{activity.actor}</span>{' '}
                    {getActivityMessage(activity)}{' '}
                    <span className="font-medium">{activity.target}</span>
                  </p>
                  <p className="text-xs text-gray-500">
                    {formatRelativeTime(activity.timestamp)}
                  </p>
                </div>
              </div>
            );
          })}
        </div>
      </CardContent>
    </Card>
  );
}
