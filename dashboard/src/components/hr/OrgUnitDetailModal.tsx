import * as React from 'react';
import { Building2, Calendar, Edit2, Trash2, GitBranch, Tag } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Modal, ModalHeader, ModalBody, ModalFooter } from '@/components/ui/modal';
import { Input } from '@/components/ui/input';
import { OrgUnitSelect } from './OrgUnitSelect';
import { OrgUnitKindSelect } from './OrgUnitKindSelect';
import { getOrgUnitKindColor, getOrgUnitKindLabel } from '@/lib/hr-constants';
import type { OrgUnitSnapshot, CreateOrgUnitInput } from '@/types/hr';

interface OrgUnitDetailModalProps {
  open: boolean;
  onClose: () => void;
  orgUnit: OrgUnitSnapshot | null;
  orgUnits?: OrgUnitSnapshot[];
  onUpdate?: (id: string, data: Partial<CreateOrgUnitInput>) => void;
  onDelete?: (id: string) => void;
  isUpdating?: boolean;
  isDeleting?: boolean;
}

export function OrgUnitDetailModal({
  open,
  onClose,
  orgUnit,
  orgUnits = [],
  onUpdate,
  onDelete,
  isUpdating,
  isDeleting,
}: OrgUnitDetailModalProps) {
  const [isEditing, setIsEditing] = React.useState(false);
  const [formData, setFormData] = React.useState<Partial<CreateOrgUnitInput>>({});
  const [showDeleteConfirm, setShowDeleteConfirm] = React.useState(false);

  React.useEffect(() => {
    if (orgUnit) {
      setFormData({
        code: orgUnit.code,
        name: orgUnit.name,
        kind: orgUnit.kind || 'department',
        parentId: orgUnit.parentId || '',
      });
    }
  }, [orgUnit]);

  const handleClose = () => {
    setIsEditing(false);
    setShowDeleteConfirm(false);
    onClose();
  };

  const handleSave = () => {
    if (orgUnit && onUpdate) {
      onUpdate(orgUnit.orgUnitId, formData);
      setIsEditing(false);
    }
  };

  const handleDelete = () => {
    if (orgUnit && onDelete) {
      onDelete(orgUnit.orgUnitId);
    }
  };

  const updateField = <K extends keyof CreateOrgUnitInput>(key: K, value: CreateOrgUnitInput[K]) => {
    setFormData((prev) => ({ ...prev, [key]: value }));
  };

  const getParentName = (parentId: string | undefined) => {
    if (!parentId) return 'Top Level';
    const parent = orgUnits.find((u) => u.orgUnitId === parentId);
    return parent?.name || 'Unknown';
  };

  if (!orgUnit) return null;

  return (
    <Modal open={open} onClose={handleClose} className="max-w-lg">
      <ModalHeader onClose={handleClose}>
        <div className="flex items-center gap-3">
          <div className="p-3 bg-green-100 rounded-lg">
            <Building2 className="h-6 w-6 text-green-600" aria-hidden="true" />
          </div>
          <div>
            <div className="font-semibold">{orgUnit.name}</div>
            <div className="flex items-center gap-2 mt-1">
              <Badge variant="outline" className="text-xs">
                {orgUnit.code}
              </Badge>
              <Badge className={`text-xs ${getOrgUnitKindColor(orgUnit.kind)}`} aria-label={`Type: ${getOrgUnitKindLabel(orgUnit.kind)}`}>
                {getOrgUnitKindLabel(orgUnit.kind)}
              </Badge>
            </div>
          </div>
        </div>
      </ModalHeader>
      <ModalBody>
        {showDeleteConfirm ? (
          <div className="text-center py-4">
            <p className="text-gray-700 mb-4">
              Are you sure you want to delete the org unit <strong>{orgUnit.name}</strong>?
            </p>
            <p className="text-sm text-gray-500">This action cannot be undone.</p>
          </div>
        ) : isEditing ? (
          <div className="space-y-4">
            <div>
              <label htmlFor="editOrgName" className="block text-sm font-medium text-gray-700 mb-1">Name</label>
              <Input
                id="editOrgName"
                value={formData.name || ''}
                onChange={(e) => updateField('name', e.target.value)}
              />
            </div>
            <div>
              <label htmlFor="editOrgCode" className="block text-sm font-medium text-gray-700 mb-1">Code</label>
              <Input
                id="editOrgCode"
                value={formData.code || ''}
                onChange={(e) => updateField('code', e.target.value)}
              />
              <p className="text-xs text-gray-500 mt-1">Unique code for the org unit</p>
            </div>
            <div className="grid grid-cols-2 gap-4">
              <div>
                <label htmlFor="editOrgKind" className="block text-sm font-medium text-gray-700 mb-1">Type</label>
                <OrgUnitKindSelect
                  id="editOrgKind"
                  value={formData.kind || 'department'}
                  onChange={(value) => updateField('kind', value)}
                />
              </div>
              <div>
                <label htmlFor="editOrgParent" className="block text-sm font-medium text-gray-700 mb-1">Parent Unit</label>
                <OrgUnitSelect
                  id="editOrgParent"
                  value={formData.parentId || ''}
                  onChange={(value) => updateField('parentId', value)}
                  orgUnits={orgUnits}
                  excludeId={orgUnit.orgUnitId}
                  noneLabel="None (Top Level)"
                />
              </div>
            </div>
          </div>
        ) : (
          <div className="space-y-4">
            <div className="space-y-3">
              <div className="flex items-center gap-3 text-gray-600">
                <Tag className="h-4 w-4 text-gray-400" aria-hidden="true" />
                <span>Code: <code className="bg-gray-100 px-2 py-0.5 rounded text-sm">{orgUnit.code}</code></span>
              </div>

              <div className="flex items-center gap-3 text-gray-600">
                <GitBranch className="h-4 w-4 text-gray-400" aria-hidden="true" />
                <span>Parent: {getParentName(orgUnit.parentId)}</span>
              </div>

              <div className="flex items-center gap-3 text-gray-600">
                <Calendar className="h-4 w-4 text-gray-400" aria-hidden="true" />
                <span>Updated {new Date(orgUnit.updatedAt).toLocaleDateString()}</span>
              </div>
            </div>

            <div className="pt-4 border-t">
              <h4 className="text-sm font-medium text-gray-700 mb-2">Details</h4>
              <dl className="grid grid-cols-2 gap-2 text-sm">
                <dt className="text-gray-500">Org Unit ID</dt>
                <dd className="text-gray-900 font-mono text-xs">{orgUnit.orgUnitId}</dd>
                <dt className="text-gray-500">Tenant ID</dt>
                <dd className="text-gray-900 font-mono text-xs">{orgUnit.tenantId}</dd>
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
              {isDeleting ? 'Deleting...' : 'Delete Org Unit'}
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
