import * as React from 'react';
import { Users, UserPlus, Trash2, Shield, Crown } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { useTeamMembers, useDeleteTeamMembership, usePeople } from '@/api/hr';
import type { PersonSnapshot } from '@/types/hr';

interface TeamMembersSectionProps {
  teamId: string;
  onAddMember?: () => void;
  canManageMembers?: boolean;
}

function getRoleBadgeVariant(role: string): 'default' | 'secondary' | 'outline' {
  switch (role) {
    case 'lead':
      return 'default';
    case 'member':
      return 'secondary';
    default:
      return 'outline';
  }
}

function getRoleIcon(role: string) {
  switch (role) {
    case 'lead':
      return <Crown className="h-3 w-3 mr-1" aria-hidden="true" />;
    case 'observer':
      return <Shield className="h-3 w-3 mr-1" aria-hidden="true" />;
    default:
      return null;
  }
}

export function TeamMembersSection({ teamId, onAddMember, canManageMembers = true }: TeamMembersSectionProps) {
  const { data: membershipsData, isLoading: loadingMemberships } = useTeamMembers(teamId);
  const { data: peopleData } = usePeople({ limit: 500 });
  const deleteMembership = useDeleteTeamMembership();

  const peopleMap = React.useMemo(() => {
    const map = new Map<string, PersonSnapshot>();
    for (const person of peopleData?.items ?? []) {
      map.set(person.personId, person);
    }
    return map;
  }, [peopleData]);

  const handleRemoveMember = async (membershipId: string) => {
    if (!confirm('Are you sure you want to remove this team member?')) return;
    try {
      await deleteMembership.mutateAsync(membershipId);
    } catch (error) {
      console.error('Failed to remove member:', error);
    }
  };

  const sortedMemberships = React.useMemo(() => {
    const memberships = membershipsData?.items ?? [];
    return [...memberships].sort((a, b) => {
      // Sort by role: lead first, then member, then observer
      const roleOrder = { lead: 0, member: 1, observer: 2 };
      const aOrder = roleOrder[a.role as keyof typeof roleOrder] ?? 3;
      const bOrder = roleOrder[b.role as keyof typeof roleOrder] ?? 3;
      return aOrder - bOrder;
    });
  }, [membershipsData?.items]);

  if (loadingMemberships) {
    return (
      <div className="py-4 text-center text-gray-500">
        Loading team members...
      </div>
    );
  }

  return (
    <div className="space-y-3">
      <div className="flex items-center justify-between">
        <h4 className="text-sm font-medium text-gray-700 flex items-center gap-2">
          <Users className="h-4 w-4" aria-hidden="true" />
          Team Members ({sortedMemberships.length})
        </h4>
        {canManageMembers && onAddMember && (
          <Button variant="outline" size="sm" onClick={onAddMember}>
            <UserPlus className="h-4 w-4 mr-1" aria-hidden="true" />
            Add Member
          </Button>
        )}
      </div>

      {sortedMemberships.length === 0 ? (
        <div className="py-6 text-center text-gray-500 bg-gray-50 rounded-lg">
          <Users className="h-8 w-8 mx-auto mb-2 text-gray-400" aria-hidden="true" />
          <p>No members in this team yet</p>
          {canManageMembers && onAddMember && (
            <Button variant="link" size="sm" onClick={onAddMember} className="mt-2">
              Add the first member
            </Button>
          )}
        </div>
      ) : (
        <ul className="divide-y divide-gray-100" role="list">
          {sortedMemberships.map((membership) => {
            const person = peopleMap.get(membership.personId);
            return (
              <li key={membership.membershipId} className="py-3 flex items-center justify-between">
                <div className="flex items-center gap-3">
                  <div className="h-8 w-8 rounded-full bg-blue-100 flex items-center justify-center text-blue-600 text-sm font-medium">
                    {person?.givenName?.[0] ?? '?'}
                    {person?.familyName?.[0] ?? ''}
                  </div>
                  <div>
                    <div className="text-sm font-medium text-gray-900">
                      {person?.fullName ?? person?.email ?? membership.personId}
                    </div>
                    {person?.title && (
                      <div className="text-xs text-gray-500">{person.title}</div>
                    )}
                  </div>
                </div>
                <div className="flex items-center gap-2">
                  <Badge variant={getRoleBadgeVariant(membership.role)} className="text-xs">
                    {getRoleIcon(membership.role)}
                    {membership.role}
                  </Badge>
                  {membership.status !== 'active' && (
                    <Badge variant="outline" className="text-xs text-gray-500">
                      {membership.status}
                    </Badge>
                  )}
                  {canManageMembers && (
                    <Button
                      variant="ghost"
                      size="sm"
                      className="h-7 w-7 p-0 text-gray-400 hover:text-red-600"
                      onClick={() => handleRemoveMember(membership.membershipId)}
                      disabled={deleteMembership.isPending}
                      aria-label={`Remove ${person?.fullName ?? 'member'} from team`}
                    >
                      <Trash2 className="h-4 w-4" aria-hidden="true" />
                    </Button>
                  )}
                </div>
              </li>
            );
          })}
        </ul>
      )}
    </div>
  );
}
