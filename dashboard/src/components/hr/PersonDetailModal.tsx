import * as React from 'react';
import { Mail, Phone, Building2, BadgeCheck, Calendar, Edit2, Trash2 } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Modal, ModalHeader, ModalBody, ModalFooter } from '@/components/ui/modal';
import { Input } from '@/components/ui/input';
import { OrgUnitSelect } from './OrgUnitSelect';
import { PersonStatusSelect } from './PersonStatusSelect';
import { getPersonStatusColor } from '@/lib/hr-constants';
import type { PersonSnapshot, OrgUnitSnapshot, CreatePersonInput } from '@/types/hr';

interface PersonDetailModalProps {
  open: boolean;
  onClose: () => void;
  person: PersonSnapshot | null;
  orgUnits?: OrgUnitSnapshot[];
  onUpdate?: (id: string, data: Partial<CreatePersonInput>) => void;
  onDelete?: (id: string) => void;
  isUpdating?: boolean;
  isDeleting?: boolean;
}

export function PersonDetailModal({
  open,
  onClose,
  person,
  orgUnits = [],
  onUpdate,
  onDelete,
  isUpdating,
  isDeleting,
}: PersonDetailModalProps) {
  const [isEditing, setIsEditing] = React.useState(false);
  const [formData, setFormData] = React.useState<Partial<CreatePersonInput>>({});
  const [showDeleteConfirm, setShowDeleteConfirm] = React.useState(false);

  React.useEffect(() => {
    if (person) {
      setFormData({
        givenName: person.givenName,
        familyName: person.familyName,
        email: person.email,
        phone: person.phone || '',
        title: person.title || '',
        status: person.status,
        orgUnitId: person.orgUnitId || '',
      });
    }
  }, [person]);

  const handleClose = () => {
    setIsEditing(false);
    setShowDeleteConfirm(false);
    onClose();
  };

  const handleSave = () => {
    if (person && onUpdate) {
      onUpdate(person.personId, formData);
      setIsEditing(false);
    }
  };

  const handleDelete = () => {
    if (person && onDelete) {
      onDelete(person.personId);
    }
  };

  const updateField = <K extends keyof CreatePersonInput>(key: K, value: CreatePersonInput[K]) => {
    setFormData((prev) => ({ ...prev, [key]: value }));
  };

  const getOrgUnitName = (orgUnitId: string) => {
    const unit = orgUnits.find((u) => u.orgUnitId === orgUnitId);
    return unit?.name || 'Not assigned';
  };

  if (!person) return null;

  return (
    <Modal open={open} onClose={handleClose} className="max-w-lg">
      <ModalHeader onClose={handleClose}>
        <div className="flex items-center gap-3">
          <div className="h-12 w-12 rounded-full bg-violet-100 flex items-center justify-center text-violet-600 font-semibold text-lg">
            {person.givenName?.[0]}{person.familyName?.[0]}
          </div>
          <div>
            <div className="font-semibold">{person.fullName || `${person.givenName} ${person.familyName}`}</div>
            {person.title && <div className="text-sm text-gray-500">{person.title}</div>}
          </div>
        </div>
      </ModalHeader>
      <ModalBody>
        {showDeleteConfirm ? (
          <div className="text-center py-4">
            <p className="text-gray-700 mb-4">
              Are you sure you want to delete <strong>{person.fullName}</strong>?
            </p>
            <p className="text-sm text-gray-500">This action cannot be undone.</p>
          </div>
        ) : isEditing ? (
          <div className="space-y-4">
            <div className="grid grid-cols-2 gap-4">
              <div>
                <label htmlFor="editGivenName" className="block text-sm font-medium text-gray-700 mb-1">First Name</label>
                <Input
                  id="editGivenName"
                  value={formData.givenName || ''}
                  onChange={(e) => updateField('givenName', e.target.value)}
                />
              </div>
              <div>
                <label htmlFor="editFamilyName" className="block text-sm font-medium text-gray-700 mb-1">Last Name</label>
                <Input
                  id="editFamilyName"
                  value={formData.familyName || ''}
                  onChange={(e) => updateField('familyName', e.target.value)}
                />
              </div>
            </div>
            <div>
              <label htmlFor="editEmail" className="block text-sm font-medium text-gray-700 mb-1">Email</label>
              <Input
                id="editEmail"
                type="email"
                value={formData.email || ''}
                onChange={(e) => updateField('email', e.target.value)}
              />
            </div>
            <div className="grid grid-cols-2 gap-4">
              <div>
                <label htmlFor="editPhone" className="block text-sm font-medium text-gray-700 mb-1">Phone</label>
                <Input
                  id="editPhone"
                  value={formData.phone || ''}
                  onChange={(e) => updateField('phone', e.target.value)}
                />
              </div>
              <div>
                <label htmlFor="editTitle" className="block text-sm font-medium text-gray-700 mb-1">Title</label>
                <Input
                  id="editTitle"
                  value={formData.title || ''}
                  onChange={(e) => updateField('title', e.target.value)}
                />
              </div>
            </div>
            <div className="grid grid-cols-2 gap-4">
              <div>
                <label htmlFor="editStatus" className="block text-sm font-medium text-gray-700 mb-1">Status</label>
                <PersonStatusSelect
                  id="editStatus"
                  value={formData.status || 'active'}
                  onChange={(value) => updateField('status', value)}
                />
              </div>
              <div>
                <label htmlFor="editOrgUnit" className="block text-sm font-medium text-gray-700 mb-1">Organization</label>
                <OrgUnitSelect
                  id="editOrgUnit"
                  value={formData.orgUnitId || ''}
                  onChange={(value) => updateField('orgUnitId', value)}
                  orgUnits={orgUnits}
                />
              </div>
            </div>
          </div>
        ) : (
          <div className="space-y-4">
            <div className="flex items-center gap-2">
              <Badge className={getPersonStatusColor(person.status)} aria-label={`Status: ${person.status}`}>
                <BadgeCheck className="h-3 w-3 mr-1" aria-hidden="true" />
                {person.status}
              </Badge>
            </div>

            <div className="space-y-3">
              <div className="flex items-center gap-3 text-gray-600">
                <Mail className="h-4 w-4 text-gray-400" aria-hidden="true" />
                <a href={`mailto:${person.email}`} className="text-violet-600 hover:underline">
                  {person.email}
                </a>
              </div>

              {person.phone && (
                <div className="flex items-center gap-3 text-gray-600">
                  <Phone className="h-4 w-4 text-gray-400" aria-hidden="true" />
                  <a href={`tel:${person.phone}`} className="hover:underline">
                    {person.phone}
                  </a>
                </div>
              )}

              <div className="flex items-center gap-3 text-gray-600">
                <Building2 className="h-4 w-4 text-gray-400" aria-hidden="true" />
                <span>{getOrgUnitName(person.orgUnitId)}</span>
              </div>

              <div className="flex items-center gap-3 text-gray-600">
                <Calendar className="h-4 w-4 text-gray-400" aria-hidden="true" />
                <span>Updated {new Date(person.updatedAt).toLocaleDateString()}</span>
              </div>
            </div>

            <div className="pt-4 border-t">
              <h4 className="text-sm font-medium text-gray-700 mb-2">Details</h4>
              <dl className="grid grid-cols-2 gap-2 text-sm">
                <dt className="text-gray-500">Person ID</dt>
                <dd className="text-gray-900 font-mono text-xs">{person.personId}</dd>
                <dt className="text-gray-500">Tenant ID</dt>
                <dd className="text-gray-900 font-mono text-xs">{person.tenantId}</dd>
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
              {isDeleting ? 'Deleting...' : 'Delete Person'}
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
    </Modal>
  );
}
