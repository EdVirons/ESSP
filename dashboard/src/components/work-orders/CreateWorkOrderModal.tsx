import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Modal, ModalHeader, ModalBody, ModalFooter } from '@/components/ui/modal';
import { Select } from '@/components/ui/select';
import { Textarea } from '@/components/ui/textarea';

interface CreateWorkOrderFormData {
  deviceId: string;
  taskType: string;
  incidentId: string;
  notes: string;
}

interface CreateWorkOrderModalProps {
  open: boolean;
  onClose: () => void;
  formData: CreateWorkOrderFormData;
  onFormChange: (data: CreateWorkOrderFormData) => void;
  onSubmit: () => void;
  isLoading: boolean;
}

export function CreateWorkOrderModal({
  open,
  onClose,
  formData,
  onFormChange,
  onSubmit,
  isLoading,
}: CreateWorkOrderModalProps) {
  return (
    <Modal open={open} onClose={onClose} className="max-w-lg">
      <ModalHeader onClose={onClose}>Create New Work Order</ModalHeader>
      <ModalBody>
        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Device ID *
            </label>
            <Input
              value={formData.deviceId}
              onChange={(e) =>
                onFormChange({ ...formData, deviceId: e.target.value })
              }
              placeholder="Enter device ID"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Task Type *
            </label>
            <Select
              value={formData.taskType}
              onChange={(value) =>
                onFormChange({ ...formData, taskType: value })
              }
              options={[
                { value: 'repair', label: 'Repair' },
                { value: 'maintenance', label: 'Maintenance' },
                { value: 'installation', label: 'Installation' },
                { value: 'inspection', label: 'Inspection' },
              ]}
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Related Incident ID
            </label>
            <Input
              value={formData.incidentId}
              onChange={(e) =>
                onFormChange({ ...formData, incidentId: e.target.value })
              }
              placeholder="Optional - link to an incident"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Notes
            </label>
            <Textarea
              value={formData.notes}
              onChange={(e) =>
                onFormChange({ ...formData, notes: e.target.value })
              }
              placeholder="Additional notes"
              rows={3}
            />
          </div>
        </div>
      </ModalBody>
      <ModalFooter>
        <Button variant="outline" onClick={onClose}>
          Cancel
        </Button>
        <Button onClick={onSubmit} disabled={!formData.deviceId || isLoading}>
          {isLoading ? 'Creating...' : 'Create Work Order'}
        </Button>
      </ModalFooter>
    </Modal>
  );
}
