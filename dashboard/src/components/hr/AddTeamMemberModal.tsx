import * as React from 'react';
import { UserPlus, Search } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Modal, ModalHeader, ModalBody, ModalFooter } from '@/components/ui/modal';
import { usePeople, useCreateTeamMembership, useTeamMembers } from '@/api/hr';
import type { PersonSnapshot, CreateTeamMembershipInput } from '@/types/hr';

interface AddTeamMemberModalProps {
  open: boolean;
  onClose: () => void;
  teamId: string;
  teamName: string;
}

const ROLES = [
  { value: 'member', label: 'Member' },
  { value: 'lead', label: 'Lead' },
  { value: 'observer', label: 'Observer' },
];

export function AddTeamMemberModal({ open, onClose, teamId, teamName }: AddTeamMemberModalProps) {
  const [searchQuery, setSearchQuery] = React.useState('');
  const [selectedPerson, setSelectedPerson] = React.useState<PersonSnapshot | null>(null);
  const [selectedRole, setSelectedRole] = React.useState('member');

  const { data: peopleData, isLoading: loadingPeople } = usePeople({ limit: 500 });
  const { data: membershipsData } = useTeamMembers(teamId);
  const createMembership = useCreateTeamMembership();

  // Get existing member person IDs
  const existingMemberIds = React.useMemo(() => {
    return new Set(membershipsData?.items?.map((m) => m.personId) ?? []);
  }, [membershipsData]);

  // Filter people who are not already members
  const availablePeople = React.useMemo(() => {
    const people = peopleData?.items ?? [];
    return people.filter((p) => {
      if (existingMemberIds.has(p.personId)) return false;
      if (p.status !== 'active') return false;
      if (!searchQuery) return true;
      const query = searchQuery.toLowerCase();
      return (
        p.fullName?.toLowerCase().includes(query) ||
        p.email?.toLowerCase().includes(query) ||
        p.title?.toLowerCase().includes(query)
      );
    });
  }, [peopleData, existingMemberIds, searchQuery]);

  const handleClose = () => {
    setSearchQuery('');
    setSelectedPerson(null);
    setSelectedRole('member');
    onClose();
  };

  const handleSubmit = async () => {
    if (!selectedPerson) return;

    const data: CreateTeamMembershipInput = {
      teamId,
      personId: selectedPerson.personId,
      role: selectedRole,
      status: 'active',
    };

    try {
      await createMembership.mutateAsync(data);
      handleClose();
    } catch (error) {
      console.error('Failed to add team member:', error);
    }
  };

  return (
    <Modal open={open} onClose={handleClose} className="max-w-md">
      <ModalHeader onClose={handleClose}>
        <div className="flex items-center gap-3">
          <div className="p-3 bg-blue-100 rounded-lg">
            <UserPlus className="h-6 w-6 text-blue-600" aria-hidden="true" />
          </div>
          <div>
            <div className="font-semibold">Add Team Member</div>
            <div className="text-sm text-gray-500">{teamName}</div>
          </div>
        </div>
      </ModalHeader>
      <ModalBody>
        <div className="space-y-4">
          {/* Search Input */}
          <div>
            <label htmlFor="searchPeople" className="block text-sm font-medium text-gray-700 mb-1">
              Search People
            </label>
            <div className="relative">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-gray-400" aria-hidden="true" />
              <Input
                id="searchPeople"
                type="text"
                placeholder="Search by name, email, or title..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="pl-9"
              />
            </div>
          </div>

          {/* People List */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Select Person
            </label>
            <div className="border rounded-lg max-h-48 overflow-y-auto">
              {loadingPeople ? (
                <div className="p-4 text-center text-gray-500">Loading...</div>
              ) : availablePeople.length === 0 ? (
                <div className="p-4 text-center text-gray-500">
                  {searchQuery ? 'No people found matching your search' : 'No available people to add'}
                </div>
              ) : (
                <ul role="listbox" className="divide-y divide-gray-100">
                  {availablePeople.slice(0, 10).map((person) => (
                    <li
                      key={person.personId}
                      role="option"
                      aria-selected={selectedPerson?.personId === person.personId}
                      className={`p-3 cursor-pointer hover:bg-gray-50 ${
                        selectedPerson?.personId === person.personId ? 'bg-blue-50 border-l-2 border-blue-500' : ''
                      }`}
                      onClick={() => setSelectedPerson(person)}
                    >
                      <div className="flex items-center gap-3">
                        <div className="h-8 w-8 rounded-full bg-blue-100 flex items-center justify-center text-blue-600 text-sm font-medium">
                          {person.givenName?.[0] ?? '?'}
                          {person.familyName?.[0] ?? ''}
                        </div>
                        <div className="flex-1 min-w-0">
                          <div className="text-sm font-medium text-gray-900 truncate">
                            {person.fullName}
                          </div>
                          <div className="text-xs text-gray-500 truncate">
                            {person.email}
                            {person.title && ` - ${person.title}`}
                          </div>
                        </div>
                      </div>
                    </li>
                  ))}
                </ul>
              )}
            </div>
            {availablePeople.length > 10 && (
              <p className="text-xs text-gray-500 mt-1">
                Showing 10 of {availablePeople.length} results. Use search to narrow down.
              </p>
            )}
          </div>

          {/* Role Selection */}
          <div>
            <label htmlFor="memberRole" className="block text-sm font-medium text-gray-700 mb-1">
              Role
            </label>
            <select
              id="memberRole"
              value={selectedRole}
              onChange={(e) => setSelectedRole(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
            >
              {ROLES.map((role) => (
                <option key={role.value} value={role.value}>
                  {role.label}
                </option>
              ))}
            </select>
          </div>
        </div>
      </ModalBody>
      <ModalFooter>
        <Button variant="outline" onClick={handleClose}>
          Cancel
        </Button>
        <Button
          onClick={handleSubmit}
          disabled={!selectedPerson || createMembership.isPending}
        >
          {createMembership.isPending ? 'Adding...' : 'Add Member'}
        </Button>
      </ModalFooter>
    </Modal>
  );
}
