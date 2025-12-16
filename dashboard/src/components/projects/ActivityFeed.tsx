import { useState } from 'react';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import { Textarea } from '@/components/ui/textarea';
import { Avatar, AvatarFallback } from '@/components/ui/avatar';
import { formatDistanceToNow } from 'date-fns';
import {
  useProjectActivities,
  useCreateActivity,
  useDeleteActivity,
  useToggleActivityPin,
} from '@/api/projects';
import type { ProjectActivity, ActivityType } from '@/types';
import {
  MessageSquare,
  FileText,
  Upload,
  ArrowRight,
  Users,
  Wrench,
  RefreshCw,
  AtSign,
  Pin,
  PinOff,
  Trash2,
  Send,
} from 'lucide-react';

const activityConfig: Record<
  ActivityType,
  { icon: React.ReactNode; label: string; color: string }
> = {
  comment: {
    icon: <MessageSquare className="h-4 w-4" />,
    label: 'Comment',
    color: 'text-blue-600',
  },
  note: {
    icon: <FileText className="h-4 w-4" />,
    label: 'Note',
    color: 'text-purple-600',
  },
  file_upload: {
    icon: <Upload className="h-4 w-4" />,
    label: 'File Upload',
    color: 'text-green-600',
  },
  status_change: {
    icon: <RefreshCw className="h-4 w-4" />,
    label: 'Status Change',
    color: 'text-orange-600',
  },
  assignment: {
    icon: <Users className="h-4 w-4" />,
    label: 'Assignment',
    color: 'text-indigo-600',
  },
  work_order: {
    icon: <Wrench className="h-4 w-4" />,
    label: 'Work Order',
    color: 'text-gray-600',
  },
  phase_transition: {
    icon: <ArrowRight className="h-4 w-4" />,
    label: 'Phase Transition',
    color: 'text-teal-600',
  },
  mention: {
    icon: <AtSign className="h-4 w-4" />,
    label: 'Mention',
    color: 'text-pink-600',
  },
};

interface ActivityFeedProps {
  projectId: string;
  phaseId?: string;
  canEdit?: boolean;
}

export function ActivityFeed({ projectId, phaseId, canEdit = true }: ActivityFeedProps) {
  const [newComment, setNewComment] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);

  const { data, isLoading, error } = useProjectActivities(projectId, {
    phaseId,
    limit: 50,
  });
  const createActivity = useCreateActivity();
  const deleteActivity = useDeleteActivity();
  const togglePin = useToggleActivityPin();

  const handleSubmitComment = async () => {
    if (!newComment.trim() || isSubmitting) return;
    setIsSubmitting(true);
    try {
      await createActivity.mutateAsync({
        projectId,
        data: {
          content: newComment.trim(),
          phaseId,
          activityType: 'comment',
        },
      });
      setNewComment('');
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleDelete = async (activityId: string) => {
    if (!confirm('Delete this activity?')) return;
    await deleteActivity.mutateAsync({ activityId, projectId });
  };

  const handleTogglePin = async (activityId: string) => {
    await togglePin.mutateAsync({ activityId, projectId });
  };

  if (isLoading) {
    return (
      <div className="animate-pulse space-y-3">
        {[1, 2, 3].map((i) => (
          <div key={i} className="h-24 bg-gray-100 rounded-lg" />
        ))}
      </div>
    );
  }

  if (error) {
    return (
      <div className="text-center py-8 text-red-500">
        Failed to load activities
      </div>
    );
  }

  const activities = data?.items || [];
  const pinnedActivities = activities.filter((a) => a.isPinned);
  const regularActivities = activities.filter((a) => !a.isPinned);

  return (
    <div className="space-y-4">
      {/* Comment Input */}
      {canEdit && (
        <Card>
          <CardContent className="p-4">
            <Textarea
              placeholder="Add a comment..."
              value={newComment}
              onChange={(e) => setNewComment(e.target.value)}
              className="mb-2 resize-none"
              rows={3}
            />
            <div className="flex justify-end">
              <Button
                size="sm"
                onClick={handleSubmitComment}
                disabled={!newComment.trim() || isSubmitting}
              >
                <Send className="h-4 w-4 mr-1" />
                {isSubmitting ? 'Posting...' : 'Post Comment'}
              </Button>
            </div>
          </CardContent>
        </Card>
      )}

      {/* Pinned Activities */}
      {pinnedActivities.length > 0 && (
        <div className="space-y-2">
          <h4 className="text-xs font-medium text-gray-500 uppercase tracking-wide flex items-center gap-1">
            <Pin className="h-3 w-3" />
            Pinned
          </h4>
          {pinnedActivities.map((activity) => (
            <ActivityItem
              key={activity.id}
              activity={activity}
              projectId={projectId}
              canEdit={canEdit}
              onDelete={() => handleDelete(activity.id)}
              onTogglePin={() => handleTogglePin(activity.id)}
            />
          ))}
        </div>
      )}

      {/* Regular Activities */}
      <div className="space-y-2">
        {regularActivities.length === 0 && pinnedActivities.length === 0 ? (
          <div className="text-center py-8 text-gray-500">
            <MessageSquare className="h-12 w-12 mx-auto mb-2 text-gray-300" />
            <p>No activity yet</p>
            <p className="text-sm">
              Be the first to add a comment!
            </p>
          </div>
        ) : (
          regularActivities.map((activity) => (
            <ActivityItem
              key={activity.id}
              activity={activity}
              projectId={projectId}
              canEdit={canEdit}
              onDelete={() => handleDelete(activity.id)}
              onTogglePin={() => handleTogglePin(activity.id)}
            />
          ))
        )}
      </div>
    </div>
  );
}

interface ActivityItemProps {
  activity: ProjectActivity;
  projectId: string;
  canEdit: boolean;
  onDelete: () => void;
  onTogglePin: () => void;
}

function ActivityItem({
  activity,
  canEdit,
  onDelete,
  onTogglePin,
}: ActivityItemProps) {
  const config = activityConfig[activity.activityType] || activityConfig.comment;
  const initials = activity.actorName
    .split(' ')
    .map((n) => n[0])
    .join('')
    .toUpperCase()
    .slice(0, 2);

  const renderContent = () => {
    switch (activity.activityType) {
      case 'status_change':
        const meta = activity.metadata as { from?: string; to?: string; reason?: string };
        return (
          <p className="text-sm text-gray-600">
            Changed status from <strong>{meta.from || '?'}</strong> to{' '}
            <strong>{meta.to || '?'}</strong>
            {meta.reason && <span className="text-gray-400"> - {meta.reason}</span>}
          </p>
        );

      case 'assignment':
        const assignMeta = activity.metadata as {
          userName?: string;
          action?: string;
          role?: string;
        };
        return (
          <p className="text-sm text-gray-600">
            {assignMeta.action === 'added' ? 'Added' : 'Removed'}{' '}
            <strong>{assignMeta.userName || 'a team member'}</strong>
            {assignMeta.role && <span> as {assignMeta.role}</span>}
          </p>
        );

      case 'work_order':
        const woMeta = activity.metadata as {
          workOrderId?: string;
          action?: string;
          status?: string;
        };
        return (
          <p className="text-sm text-gray-600">
            {woMeta.action === 'created' ? 'Created' : 'Updated'} work order{' '}
            <span className="font-mono text-xs bg-gray-100 px-1 rounded">
              {woMeta.workOrderId}
            </span>
          </p>
        );

      case 'phase_transition':
        const phaseMeta = activity.metadata as {
          phaseType?: string;
          from?: string;
          to?: string;
        };
        return (
          <p className="text-sm text-gray-600">
            Phase <strong>{phaseMeta.phaseType}</strong> moved from{' '}
            <strong>{phaseMeta.from || '?'}</strong> to{' '}
            <strong>{phaseMeta.to || '?'}</strong>
          </p>
        );

      default:
        return (
          activity.content && (
            <p className="text-sm text-gray-700 whitespace-pre-wrap">
              {activity.content}
            </p>
          )
        );
    }
  };

  return (
    <Card className={activity.isPinned ? 'border-yellow-300 bg-yellow-50' : ''}>
      <CardContent className="p-4">
        <div className="flex items-start gap-3">
          <Avatar className="h-8 w-8">
            <AvatarFallback className="bg-gray-100 text-gray-600 text-xs">
              {initials || '?'}
            </AvatarFallback>
          </Avatar>

          <div className="flex-1 min-w-0">
            <div className="flex items-center gap-2 mb-1">
              <span className="font-medium text-gray-900 text-sm">
                {activity.actorName || 'Unknown'}
              </span>
              <span className={`${config.color}`}>{config.icon}</span>
              <Badge variant="outline" className="text-xs">
                {config.label}
              </Badge>
              {activity.isPinned && (
                <Pin className="h-3 w-3 text-yellow-600" />
              )}
              <span className="text-xs text-gray-400 ml-auto">
                {formatDistanceToNow(new Date(activity.createdAt), {
                  addSuffix: true,
                })}
              </span>
            </div>

            {renderContent()}

            {activity.editedAt && (
              <span className="text-xs text-gray-400">(edited)</span>
            )}
          </div>

          {canEdit && (activity.activityType === 'comment' || activity.activityType === 'note') && (
            <div className="flex items-center gap-1">
              <Button
                variant="ghost"
                size="sm"
                className="h-8 w-8 p-0 text-gray-400 hover:text-yellow-600"
                onClick={onTogglePin}
                title={activity.isPinned ? 'Unpin' : 'Pin'}
              >
                {activity.isPinned ? (
                  <PinOff className="h-4 w-4" />
                ) : (
                  <Pin className="h-4 w-4" />
                )}
              </Button>
              <Button
                variant="ghost"
                size="sm"
                className="h-8 w-8 p-0 text-gray-400 hover:text-red-500"
                onClick={onDelete}
                title="Delete"
              >
                <Trash2 className="h-4 w-4" />
              </Button>
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  );
}
