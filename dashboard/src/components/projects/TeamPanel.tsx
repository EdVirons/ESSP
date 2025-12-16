import { useState } from 'react';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import { Avatar, AvatarFallback } from '@/components/ui/avatar';
import { formatDate } from '@/lib/utils';
import { useProjectTeam, useRemoveTeamMember } from '@/api/projects';
import type { ProjectTeamMember, TeamMemberRole, PhaseType } from '@/types';
import { Plus, X, User, Crown, Eye, UserCog } from 'lucide-react';
import { phaseTypeLabels } from '@/lib/projectTypes';

const roleConfig: Record<
  TeamMemberRole,
  { label: string; color: string; icon: React.ReactNode }
> = {
  owner: {
    label: 'Owner',
    color: 'bg-yellow-100 text-yellow-800',
    icon: <Crown className="h-3 w-3" />,
  },
  collaborator: {
    label: 'Collaborator',
    color: 'bg-blue-100 text-blue-800',
    icon: <UserCog className="h-3 w-3" />,
  },
  viewer: {
    label: 'Viewer',
    color: 'bg-gray-100 text-gray-800',
    icon: <Eye className="h-3 w-3" />,
  },
};

interface TeamPanelProps {
  projectId: string;
  onAddMember?: () => void;
  canEdit?: boolean;
}

export function TeamPanel({ projectId, onAddMember, canEdit = true }: TeamPanelProps) {
  const { data, isLoading, error } = useProjectTeam(projectId);
  const removeMember = useRemoveTeamMember();
  const [removingId, setRemovingId] = useState<string | null>(null);

  const handleRemove = async (memberId: string) => {
    if (!confirm('Are you sure you want to remove this team member?')) return;
    setRemovingId(memberId);
    try {
      await removeMember.mutateAsync({ projectId, memberId });
    } finally {
      setRemovingId(null);
    }
  };

  if (isLoading) {
    return (
      <div className="animate-pulse space-y-3">
        {[1, 2, 3].map((i) => (
          <div key={i} className="h-20 bg-gray-100 rounded-lg" />
        ))}
      </div>
    );
  }

  if (error) {
    return (
      <div className="text-center py-8 text-red-500">
        Failed to load team members
      </div>
    );
  }

  const members = data?.members || [];

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h3 className="text-sm font-medium text-gray-700">
          Team Members ({members.length})
        </h3>
        {canEdit && onAddMember && (
          <Button size="sm" onClick={onAddMember}>
            <Plus className="h-4 w-4 mr-1" />
            Add Member
          </Button>
        )}
      </div>

      {members.length === 0 ? (
        <div className="text-center py-8 text-gray-500">
          <User className="h-12 w-12 mx-auto mb-2 text-gray-300" />
          <p>No team members yet</p>
          {canEdit && onAddMember && (
            <Button
              variant="outline"
              size="sm"
              className="mt-2"
              onClick={onAddMember}
            >
              Add the first member
            </Button>
          )}
        </div>
      ) : (
        <div className="space-y-3">
          {members.map((member) => (
            <TeamMemberCard
              key={member.id}
              member={member}
              canEdit={canEdit}
              isRemoving={removingId === member.id}
              onRemove={() => handleRemove(member.id)}
            />
          ))}
        </div>
      )}
    </div>
  );
}

interface TeamMemberCardProps {
  member: ProjectTeamMember;
  canEdit: boolean;
  isRemoving: boolean;
  onRemove: () => void;
}

function TeamMemberCard({
  member,
  canEdit,
  isRemoving,
  onRemove,
}: TeamMemberCardProps) {
  const role = roleConfig[member.role] || roleConfig.collaborator;
  const initials = member.userName
    .split(' ')
    .map((n) => n[0])
    .join('')
    .toUpperCase()
    .slice(0, 2);

  return (
    <Card>
      <CardContent className="p-4">
        <div className="flex items-start gap-3">
          <Avatar className="h-10 w-10">
            <AvatarFallback className="bg-blue-100 text-blue-600 text-sm">
              {initials || '?'}
            </AvatarFallback>
          </Avatar>

          <div className="flex-1 min-w-0">
            <div className="flex items-center gap-2">
              <span className="font-medium text-gray-900 truncate">
                {member.userName || 'Unknown User'}
              </span>
              <Badge className={role.color} variant="secondary">
                <span className="mr-1">{role.icon}</span>
                {role.label}
              </Badge>
            </div>

            {member.userEmail && (
              <p className="text-sm text-gray-500 truncate">{member.userEmail}</p>
            )}

            {member.responsibility && (
              <p className="text-sm text-gray-600 mt-1">{member.responsibility}</p>
            )}

            {member.assignedPhases.length > 0 && (
              <div className="flex flex-wrap gap-1 mt-2">
                {member.assignedPhases.map((phase) => (
                  <Badge key={phase} variant="outline" className="text-xs">
                    {phaseTypeLabels[phase as PhaseType] || phase}
                  </Badge>
                ))}
              </div>
            )}

            <p className="text-xs text-gray-400 mt-2">
              Added {formatDate(member.assignedAt)}
            </p>
          </div>

          {canEdit && (
            <Button
              variant="ghost"
              size="sm"
              className="text-gray-400 hover:text-red-500"
              onClick={onRemove}
              disabled={isRemoving}
            >
              <X className="h-4 w-4" />
            </Button>
          )}
        </div>
      </CardContent>
    </Card>
  );
}
