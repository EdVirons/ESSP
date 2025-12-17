import * as React from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Modal, ModalHeader, ModalBody, ModalFooter } from '@/components/ui/modal';
import { useImpersonation } from '@/contexts/ImpersonationContext';
import { useImpersonatableUsers } from '@/api/inventory';
import { Search, UserCircle2, Building2 } from 'lucide-react';

interface ImpersonationModalProps {
  open: boolean;
  onClose: () => void;
}

export function ImpersonationModal({ open, onClose }: ImpersonationModalProps) {
  const [searchQuery, setSearchQuery] = React.useState('');
  const [selectedUserId, setSelectedUserId] = React.useState<string | null>(null);
  const [reason, setReason] = React.useState('');
  const [isStarting, setIsStarting] = React.useState(false);

  const { startImpersonation } = useImpersonation();
  const { data: usersData, isLoading } = useImpersonatableUsers();

  const users = usersData?.items || [];

  // Filter users by search query
  const filteredUsers = React.useMemo(() => {
    if (!searchQuery.trim()) return users;
    const query = searchQuery.toLowerCase();
    return users.filter(
      (user) =>
        user.name.toLowerCase().includes(query) ||
        user.email.toLowerCase().includes(query)
    );
  }, [users, searchQuery]);

  const handleStart = async () => {
    if (!selectedUserId) return;

    setIsStarting(true);
    try {
      const success = await startImpersonation(selectedUserId, reason);
      if (success) {
        onClose();
        // Reset state
        setSearchQuery('');
        setSelectedUserId(null);
        setReason('');
      }
    } finally {
      setIsStarting(false);
    }
  };

  const handleClose = () => {
    setSearchQuery('');
    setSelectedUserId(null);
    setReason('');
    onClose();
  };

  const selectedUser = users.find((u) => u.userId === selectedUserId);

  return (
    <Modal open={open} onClose={handleClose} className="max-w-lg">
      <ModalHeader onClose={handleClose}>
        <div className="flex items-center gap-2">
          <UserCircle2 className="h-5 w-5 text-orange-500" />
          Impersonate School Contact
        </div>
      </ModalHeader>

      <ModalBody>
        <div className="space-y-4">
          <p className="text-sm text-gray-500">
            Select a school contact to act on their behalf. All actions will be
            logged with your identity as the actor.
          </p>

          {/* Search */}
          <div className="relative">
            <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-400" />
            <Input
              placeholder="Search by name or email..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="pl-10"
            />
          </div>

          {/* User list */}
          <div className="max-h-64 overflow-y-auto rounded-lg border border-gray-200">
            {isLoading ? (
              <div className="p-4 text-center text-gray-500">
                Loading school contacts...
              </div>
            ) : filteredUsers.length === 0 ? (
              <div className="p-4 text-center text-gray-500">
                {searchQuery ? 'No matching contacts found' : 'No school contacts available'}
              </div>
            ) : (
              <div className="divide-y divide-gray-100">
                {filteredUsers.map((user) => (
                  <button
                    key={user.userId}
                    onClick={() => setSelectedUserId(user.userId)}
                    className={`w-full flex items-start gap-3 p-3 text-left hover:bg-gray-50 transition-colors ${
                      selectedUserId === user.userId ? 'bg-orange-50 border-l-4 border-orange-500' : ''
                    }`}
                  >
                    <div className="p-2 bg-gray-100 rounded-full">
                      <UserCircle2 className="h-5 w-5 text-gray-600" />
                    </div>
                    <div className="flex-1 min-w-0">
                      <p className="font-medium text-gray-900 truncate">{user.name}</p>
                      <p className="text-sm text-gray-500 truncate">{user.email}</p>
                      <div className="flex items-center gap-1 mt-1 text-xs text-gray-400">
                        <Building2 className="h-3 w-3" />
                        <span>
                          {user.schools.length} school{user.schools.length !== 1 ? 's' : ''}
                        </span>
                      </div>
                    </div>
                  </button>
                ))}
              </div>
            )}
          </div>

          {/* Selected user summary */}
          {selectedUser && (
            <div className="bg-orange-50 border border-orange-200 rounded-lg p-3">
              <p className="text-sm font-medium text-orange-900">
                Selected: {selectedUser.name}
              </p>
              <p className="text-xs text-orange-700 mt-1">
                You will have access to {selectedUser.schools.length} school
                {selectedUser.schools.length !== 1 ? 's' : ''}.
              </p>
            </div>
          )}

          {/* Reason field */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Reason for impersonation (optional)
            </label>
            <Input
              placeholder="e.g., Helping register devices"
              value={reason}
              onChange={(e) => setReason(e.target.value)}
            />
            <p className="mt-1 text-xs text-gray-500">
              This will be logged for audit purposes.
            </p>
          </div>
        </div>
      </ModalBody>

      <ModalFooter>
        <Button variant="outline" onClick={handleClose}>
          Cancel
        </Button>
        <Button
          onClick={handleStart}
          disabled={!selectedUserId || isStarting}
          className="bg-orange-500 hover:bg-orange-600"
        >
          {isStarting ? 'Starting...' : 'Start Impersonation'}
        </Button>
      </ModalFooter>
    </Modal>
  );
}
