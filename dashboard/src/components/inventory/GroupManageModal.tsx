import * as React from 'react';
import { Plus, Minus, Users } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { Modal, ModalHeader, ModalBody, ModalFooter } from '@/components/ui/modal';
import type { InventoryDevice, DeviceGroup, CreateGroupRequest, GroupType } from '@/types';

const GROUP_TYPES: Array<{ value: GroupType; label: string; description: string }> = [
  { value: 'manual', label: 'Manual', description: 'Manually add devices to this group' },
  { value: 'location', label: 'Location-based', description: 'Auto-include devices at a location' },
  { value: 'dynamic', label: 'Dynamic', description: 'Auto-include devices matching criteria' },
];

interface GroupManageModalProps {
  open: boolean;
  onClose: () => void;
  device?: InventoryDevice | null;
  groups: DeviceGroup[];
  deviceGroupIds?: string[]; // Groups the device is currently in
  onAddToGroup: (groupId: string, deviceId: string) => void;
  onRemoveFromGroup: (groupId: string, deviceId: string) => void;
  onCreateGroup: (data: CreateGroupRequest) => void;
  isLoading: boolean;
}

export function GroupManageModal({
  open,
  onClose,
  device,
  groups,
  deviceGroupIds = [],
  onAddToGroup,
  onRemoveFromGroup,
  onCreateGroup,
  isLoading,
}: GroupManageModalProps) {
  const [mode, setMode] = React.useState<'manage' | 'create'>('manage');
  const [newGroup, setNewGroup] = React.useState<CreateGroupRequest>({
    name: '',
    description: '',
    groupType: 'manual',
  });

  // Reset when modal opens
  React.useEffect(() => {
    if (open) {
      setMode('manage');
      setNewGroup({
        name: '',
        description: '',
        groupType: 'manual',
      });
    }
  }, [open]);

  const handleCreateGroup = (e: React.FormEvent) => {
    e.preventDefault();
    onCreateGroup(newGroup);
    setMode('manage');
    setNewGroup({ name: '', description: '', groupType: 'manual' });
  };

  const activeGroups = groups.filter(g => g.active);
  const isInGroup = (groupId: string) => deviceGroupIds.includes(groupId);

  return (
    <Modal open={open} onClose={onClose} className="max-w-lg">
      <ModalHeader onClose={onClose}>
        {mode === 'create' ? 'Create New Group' : 'Manage Device Groups'}
      </ModalHeader>

      {mode === 'create' ? (
        <form onSubmit={handleCreateGroup}>
          <ModalBody>
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Group Name *
                </label>
                <Input
                  value={newGroup.name}
                  onChange={(e) => setNewGroup({ ...newGroup, name: e.target.value })}
                  placeholder="e.g., Lab A Devices"
                  required
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Group Type
                </label>
                <select
                  value={newGroup.groupType}
                  onChange={(e) => setNewGroup({ ...newGroup, groupType: e.target.value as GroupType })}
                  className="w-full rounded-md border border-gray-300 bg-white px-3 py-2 text-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
                >
                  {GROUP_TYPES.map((type) => (
                    <option key={type.value} value={type.value}>
                      {type.label}
                    </option>
                  ))}
                </select>
                <p className="mt-1 text-xs text-gray-500">
                  {GROUP_TYPES.find(t => t.value === newGroup.groupType)?.description}
                </p>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Description
                </label>
                <Textarea
                  value={newGroup.description}
                  onChange={(e) => setNewGroup({ ...newGroup, description: e.target.value })}
                  placeholder="What is this group for?"
                  rows={2}
                />
              </div>
            </div>
          </ModalBody>
          <ModalFooter>
            <Button type="button" variant="outline" onClick={() => setMode('manage')}>
              Back
            </Button>
            <Button type="submit" disabled={!newGroup.name || isLoading}>
              {isLoading ? 'Creating...' : 'Create Group'}
            </Button>
          </ModalFooter>
        </form>
      ) : (
        <>
          <ModalBody>
            <div className="space-y-4">
              {/* Device info */}
              {device && (
                <div className="bg-gray-50 rounded-lg p-3">
                  <p className="text-xs text-gray-500">Device</p>
                  <p className="font-medium text-gray-900">
                    {device.model} - {device.serial}
                  </p>
                </div>
              )}

              {/* Groups list */}
              {activeGroups.length > 0 ? (
                <div className="space-y-2">
                  <p className="text-sm font-medium text-gray-700">Available Groups</p>
                  {activeGroups.map((group) => {
                    const inGroup = isInGroup(group.id);
                    return (
                      <div
                        key={group.id}
                        className={`flex items-center justify-between p-3 rounded-lg border ${
                          inGroup
                            ? 'bg-blue-50 border-blue-200'
                            : 'bg-white border-gray-200'
                        }`}
                      >
                        <div className="flex items-center gap-3">
                          <div className={`p-2 rounded-lg ${
                            inGroup ? 'bg-blue-100' : 'bg-gray-100'
                          }`}>
                            <Users className={`h-4 w-4 ${
                              inGroup ? 'text-blue-600' : 'text-gray-500'
                            }`} />
                          </div>
                          <div>
                            <p className="font-medium text-gray-900">{group.name}</p>
                            {group.description && (
                              <p className="text-xs text-gray-500">{group.description}</p>
                            )}
                            <p className="text-xs text-gray-400">
                              {group.groupType} {group.memberCount !== undefined && `- ${group.memberCount} devices`}
                            </p>
                          </div>
                        </div>
                        {device && group.groupType === 'manual' && (
                          <Button
                            variant={inGroup ? 'outline' : 'default'}
                            size="sm"
                            onClick={() => {
                              if (inGroup) {
                                onRemoveFromGroup(group.id, device.id);
                              } else {
                                onAddToGroup(group.id, device.id);
                              }
                            }}
                            disabled={isLoading}
                          >
                            {inGroup ? (
                              <><Minus className="h-3 w-3 mr-1" /> Remove</>
                            ) : (
                              <><Plus className="h-3 w-3 mr-1" /> Add</>
                            )}
                          </Button>
                        )}
                      </div>
                    );
                  })}
                </div>
              ) : (
                <div className="text-center py-6">
                  <Users className="h-10 w-10 mx-auto text-gray-300 mb-2" />
                  <p className="text-gray-500">No groups created yet</p>
                </div>
              )}
            </div>
          </ModalBody>
          <ModalFooter>
            <Button variant="outline" onClick={() => setMode('create')}>
              <Plus className="h-4 w-4 mr-2" />
              Create New Group
            </Button>
            <Button onClick={onClose}>Done</Button>
          </ModalFooter>
        </>
      )}
    </Modal>
  );
}
