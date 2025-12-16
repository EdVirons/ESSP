import { Link } from 'react-router-dom';
import {
  ClipboardList,
  Clock,
  AlertCircle,
  CheckCircle2,
  ArrowRight,
  Wrench,
} from 'lucide-react';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { cn } from '@/lib/utils';

interface Task {
  id: string;
  title: string;
  type: 'work_order' | 'approval' | 'incident';
  status: 'urgent' | 'pending' | 'in_progress';
  dueDate?: string;
  school?: string;
}

// Mock data - in production, this would come from API
const mockTasks: Task[] = [
  {
    id: 'wo-001',
    title: 'Laptop screen replacement',
    type: 'work_order',
    status: 'urgent',
    dueDate: '2024-12-14',
    school: 'Nairobi Primary',
  },
  {
    id: 'wo-002',
    title: 'Battery replacement - Dell Latitude',
    type: 'work_order',
    status: 'in_progress',
    dueDate: '2024-12-16',
    school: 'Mombasa Secondary',
  },
  {
    id: 'appr-001',
    title: 'Work order approval pending',
    type: 'approval',
    status: 'pending',
    school: 'Kisumu Academy',
  },
  {
    id: 'inc-001',
    title: 'Device not turning on',
    type: 'incident',
    status: 'pending',
    school: 'Nakuru High',
  },
];

const statusConfig = {
  urgent: {
    label: 'Urgent',
    color: 'bg-red-100 text-red-800',
    icon: AlertCircle,
  },
  pending: {
    label: 'Pending',
    color: 'bg-yellow-100 text-yellow-800',
    icon: Clock,
  },
  in_progress: {
    label: 'In Progress',
    color: 'bg-blue-100 text-blue-800',
    icon: Wrench,
  },
};

const typeConfig = {
  work_order: { label: 'Work Order', href: '/work-orders' },
  approval: { label: 'Approval', href: '/approvals' },
  incident: { label: 'Incident', href: '/incidents' },
};

interface MyTasksProps {
  tasks?: Task[];
  isLoading?: boolean;
}

export function MyTasks({ tasks = mockTasks, isLoading }: MyTasksProps) {
  const urgentCount = tasks.filter((t) => t.status === 'urgent').length;
  const pendingCount = tasks.filter((t) => t.status === 'pending').length;

  if (isLoading) {
    return (
      <Card>
        <CardHeader className="pb-3">
          <CardTitle className="flex items-center gap-2 text-base">
            <ClipboardList className="h-5 w-5" />
            My Tasks
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-3">
            {[1, 2, 3, 4].map((i) => (
              <div key={i} className="skeleton h-16 rounded" />
            ))}
          </div>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card>
      <CardHeader className="pb-3">
        <div className="flex items-center justify-between">
          <CardTitle className="flex items-center gap-2 text-base">
            <ClipboardList className="h-5 w-5" />
            My Tasks
          </CardTitle>
          <div className="flex gap-2">
            {urgentCount > 0 && (
              <Badge className="bg-red-100 text-red-800">
                {urgentCount} urgent
              </Badge>
            )}
            {pendingCount > 0 && (
              <Badge className="bg-yellow-100 text-yellow-800">
                {pendingCount} pending
              </Badge>
            )}
          </div>
        </div>
      </CardHeader>
      <CardContent>
        {tasks.length === 0 ? (
          <div className="text-center py-8 text-gray-500">
            <CheckCircle2 className="h-12 w-12 mx-auto mb-3 text-green-500" />
            <p className="font-medium">All caught up!</p>
            <p className="text-sm">No pending tasks assigned to you.</p>
          </div>
        ) : (
          <div className="space-y-3">
            {tasks.slice(0, 5).map((task) => {
              const status = statusConfig[task.status];
              const type = typeConfig[task.type];
              const StatusIcon = status.icon;

              return (
                <Link
                  key={task.id}
                  to={`${type.href}/${task.id}`}
                  className="block p-3 rounded-lg border border-gray-100 hover:border-gray-200 hover:bg-gray-50 transition-colors"
                >
                  <div className="flex items-start justify-between gap-3">
                    <div className="flex items-start gap-3 min-w-0">
                      <div
                        className={cn(
                          'rounded-full p-1.5 mt-0.5',
                          task.status === 'urgent'
                            ? 'bg-red-100'
                            : task.status === 'pending'
                            ? 'bg-yellow-100'
                            : 'bg-blue-100'
                        )}
                      >
                        <StatusIcon
                          className={cn(
                            'h-4 w-4',
                            task.status === 'urgent'
                              ? 'text-red-600'
                              : task.status === 'pending'
                              ? 'text-yellow-600'
                              : 'text-blue-600'
                          )}
                        />
                      </div>
                      <div className="min-w-0">
                        <p className="font-medium text-gray-900 truncate">
                          {task.title}
                        </p>
                        <div className="flex items-center gap-2 mt-1 text-xs text-gray-500">
                          <span>{type.label}</span>
                          {task.school && (
                            <>
                              <span>•</span>
                              <span>{task.school}</span>
                            </>
                          )}
                          {task.dueDate && (
                            <>
                              <span>•</span>
                              <span>Due {task.dueDate}</span>
                            </>
                          )}
                        </div>
                      </div>
                    </div>
                    <Badge className={status.color}>{status.label}</Badge>
                  </div>
                </Link>
              );
            })}
            {tasks.length > 5 && (
              <Link to="/work-orders?assignee=me">
                <Button variant="outline" className="w-full">
                  View all {tasks.length} tasks
                  <ArrowRight className="h-4 w-4 ml-2" />
                </Button>
              </Link>
            )}
          </div>
        )}
      </CardContent>
    </Card>
  );
}
