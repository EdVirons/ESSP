import * as React from 'react';
import { UsersRound, Building2, Calendar, Edit2, Trash2, Key } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Modal, ModalHeader, ModalBody, ModalFooter } from '@/components/ui/modal';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { OrgUnitSelect } from './OrgUnitSelect';
import { TeamMembersSection } from './TeamMembersSection';
import { AddTeamMemberModal } from './AddTeamMemberModal';
import type { TeamSnapshot, OrgUnitSnapshot, CreateTeamInput } from '@/types/hr';

interface TeamDetailModalProps {
  open: boolean;
  onClose: () => void;
  team: TeamSnapshot | null;
  orgUnits?: OrgUnitSnapshot[];
  onUpdate?: (id: string, data: Partial<CreateTeamInput>) => void;
  onDelete?: (id: string) => void;
  isUpdating?: boolean;
  isDeleting?: boolean;
}

export function TeamDetailModal({
  open,
  onClose,
  team,
  orgUnits = [],
  onUpdate,
  onDelete,
  isUpdating,
  isDeleting,
}: TeamDetailModalProps) {
  const [isEditing, setIsEditing] = React.useState(false);
  const [formData, setFormData] = React.useState<Partial<CreateTeamInput>>({});
  const [showDeleteConfirm, setShowDeleteConfirm] = React.useState(false);
  const [showAddMember, setShowAddMember] = React.useState(false);

  React.useEffect(() => {
    if (team) {
      setFormData({
        key: team.key,
        name: team.name,
        description: team.description || '',
        orgUnitId: team.orgUnitId || '',
      });
    }
  }, [team]);

  const handleClose = () => {
    setIsEditing(false);
    setShowDeleteConfirm(false);
    setShowAddMember(false);
    onClose();
  };

  const handleSave = () => {
    if (team && onUpdate) {
      onUpdate(team.teamId, formData);
      setIsEditing(false);
    }
  };

  const handleDelete = () => {
    if (team && onDelete) {
      onDelete(team.teamId);
    }
  };

  const updateField = <K extends keyof CreateTeamInput>(key: K, value: CreateTeamInput[K]) => {
    setFormData((prev) => ({ ...prev, [key]: value }));
  };

  const getOrgUnitName = (orgUnitId: string) => {
    const unit = orgUnits.find((u) => u.orgUnitId === orgUnitId);
    return unit?.name || 'Not assigned';
  };

  if (!team) return null;

  return (
    <Modal open={open} onClose={handleClose} className="max-w-lg">
      <ModalHeader onClose={handleClose}>
        <div className="flex items-center gap-3">
          <div className="p-3 bg-blue-100 rounded-lg">
            <UsersRound className="h-6 w-6 text-blue-600" aria-hidden="true" />
          </div>
          <div>
            <div className="font-semibold">{team.name}</div>
            <Badge variant="outline" className="text-xs mt-1">
              <Key className="h-3 w-3 mr-1" aria-hidden="true" />
              {team.key}
            </Badge>
          </div>
        </div>
      </ModalHeader>
      <ModalBody>
        {showDeleteConfirm ? (
          <div className="text-center py-4">
            <p className="text-gray-700 mb-4">
              Are you sure you want to delete the team <strong>{team.name}</strong>?
            </p>
            <p className="text-sm text-gray-500">This action cannot be undone.</p>
          </div>
        ) : isEditing ? (
          <div className="space-y-4">
            <div>
              <label htmlFor="editTeamName" className="block text-sm font-medium text-gray-700 mb-1">Team Name</label>
              <Input
                id="editTeamName"
                value={formData.name || ''}
                onChange={(e) => updateField('name', e.target.value)}
              />
            </div>
            <div>
              <label htmlFor="editTeamKey" className="block text-sm font-medium text-gray-700 mb-1">Team Key</label>
              <Input
                id="editTeamKey"
                value={formData.key || ''}
                onChange={(e) => updateField('key', e.target.value)}
              />
              <p className="text-xs text-gray-500 mt-1">Unique identifier for the team</p>
            </div>
            <div>
              <label htmlFor="editTeamDesc" className="block text-sm font-medium text-gray-700 mb-1">Description</label>
              <Textarea
                id="editTeamDesc"
                value={formData.description || ''}
                onChange={(e) => updateField('description', e.target.value)}
                rows={3}
              />
            </div>
            <div>
              <label htmlFor="editTeamOrg" className="block text-sm font-medium text-gray-700 mb-1">Organization Unit</label>
              <OrgUnitSelect
                id="editTeamOrg"
                value={formData.orgUnitId || ''}
                onChange={(value) => updateField('orgUnitId', value)}
                orgUnits={orgUnits}
              />
            </div>
          </div>
        ) : (
          <div className="space-y-4">
            {team.description && (
              <div>
                <h4 className="text-sm font-medium text-gray-700 mb-1">Description</h4>
                <p className="text-gray-600">{team.description}</p>
              </div>
            )}

            <div className="space-y-3">
              <div className="flex items-center gap-3 text-gray-600">
                <Building2 className="h-4 w-4 text-gray-400" aria-hidden="true" />
                <span>{getOrgUnitName(team.orgUnitId)}</span>
              </div>

              <div className="flex items-center gap-3 text-gray-600">
                <Calendar className="h-4 w-4 text-gray-400" aria-hidden="true" />
                <span>Updated {new Date(team.updatedAt).toLocaleDateString()}</span>
              </div>
            </div>

            {/* Team Members Section */}
            <div className="pt-4 border-t">
              <TeamMembersSection
                teamId={team.teamId}
                onAddMember={() => setShowAddMember(true)}
                canManageMembers={!!onUpdate}
              />
            </div>

            <div className="pt-4 border-t">
              <h4 className="text-sm font-medium text-gray-700 mb-2">Details</h4>
              <dl className="grid grid-cols-2 gap-2 text-sm">
                <dt className="text-gray-500">Team ID</dt>
                <dd className="text-gray-900 font-mono text-xs">{team.teamId}</dd>
                <dt className="text-gray-500">Tenant ID</dt>
                <dd className="text-gray-900 font-mono text-xs">{team.tenantId}</dd>
              </dl>
            </div>
          </div>
        )}
      </ModalBody>
      <ModalFooter>
        {showDeleteConfirm ? (
          <>
            <Button variant="outline" onClick={() => setShowDeleteConfirm(false)}>
              Cancel
            </Button>
            <Button variant="destructive" onClick={handleDelete} disabled={isDeleting}>
              {isDeleting ? 'Deleting...' : 'Delete Team'}
            </Button>
          </>
        ) : isEditing ? (
          <>
            <Button variant="outline" onClick={() => setIsEditing(false)}>
              Cancel
            </Button>
            <Button onClick={handleSave} disabled={isUpdating}>
              {isUpdating ? 'Saving...' : 'Save Changes'}
            </Button>
          </>
        ) : (
          <>
            {onDelete && (
              <Button variant="ghost" className="text-red-600 hover:text-red-700 hover:bg-red-50" onClick={() => setShowDeleteConfirm(true)}>
                <Trash2 className="h-4 w-4 mr-2" aria-hidden="true" />
                Delete
              </Button>
            )}
            <div className="flex-1" />
            {onUpdate && (
              <Button variant="outline" onClick={() => setIsEditing(true)}>
                <Edit2 className="h-4 w-4 mr-2" aria-hidden="true" />
                Edit
              </Button>
            )}
            <Button onClick={handleClose}>Close</Button>
          </>
        )}
      </ModalFooter>

      {/* Add Team Member Modal */}
      <AddTeamMemberModal
        open={showAddMember}
        onClose={() => setShowAddMember(false)}
        teamId={team.teamId}
        teamName={team.name}
      />
    </Modal>
  );
}
