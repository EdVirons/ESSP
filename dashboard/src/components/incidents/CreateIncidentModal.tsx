import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Modal, ModalHeader, ModalBody, ModalFooter } from '@/components/ui/modal';
import { Select } from '@/components/ui/select';
import { Textarea } from '@/components/ui/textarea';
import type { Severity } from '@/types';

interface CreateIncidentFormData {
  deviceId: string;
  title: string;
  description: string;
  category: string;
  severity: Severity;
}

interface CreateIncidentModalProps {
  open: boolean;
  onClose: () => void;
  formData: CreateIncidentFormData;
  onFormChange: (data: CreateIncidentFormData) => void;
  onSubmit: () => void;
  isLoading: boolean;
}

export function CreateIncidentModal({
  open,
  onClose,
  formData,
  onFormChange,
  onSubmit,
  isLoading,
}: CreateIncidentModalProps) {
  return (
    <Modal open={open} onClose={onClose} className="max-w-lg">
      <ModalHeader onClose={onClose}>Create New Incident</ModalHeader>
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
              Title *
            </label>
            <Input
              value={formData.title}
              onChange={(e) =>
                onFormChange({ ...formData, title: e.target.value })
              }
              placeholder="Brief description of the issue"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Description
            </label>
            <Textarea
              value={formData.description}
              onChange={(e) =>
                onFormChange({ ...formData, description: e.target.value })
              }
              placeholder="Detailed description of the issue"
              rows={3}
            />
          </div>
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Category
              </label>
              <Select
                value={formData.category}
                onChange={(value) =>
                  onFormChange({ ...formData, category: value })
                }
                options={[
                  { value: 'hardware', label: 'Hardware' },
                  { value: 'software', label: 'Software' },
                  { value: 'network', label: 'Network' },
                  { value: 'other', label: 'Other' },
                ]}
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Severity
              </label>
              <Select
                value={formData.severity}
                onChange={(value) =>
                  onFormChange({ ...formData, severity: value as Severity })
                }
                options={[
                  { value: 'low', label: 'Low' },
                  { value: 'medium', label: 'Medium' },
                  { value: 'high', label: 'High' },
                  { value: 'critical', label: 'Critical' },
                ]}
              />
            </div>
          </div>
        </div>
      </ModalBody>
      <ModalFooter>
        <Button variant="outline" onClick={onClose}>
          Cancel
        </Button>
        <Button
          onClick={onSubmit}
          disabled={!formData.deviceId || !formData.title || isLoading}
        >
          {isLoading ? 'Creating...' : 'Create Incident'}
        </Button>
      </ModalFooter>
    </Modal>
  );
}
