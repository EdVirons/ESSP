import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Modal, ModalHeader, ModalBody, ModalFooter } from '@/components/ui/modal';
import { DatePicker } from '@/components/ui/date-picker';
import { Textarea } from '@/components/ui/textarea';
import { Select } from '@/components/ui/select';
import { projectTypeConfigs, projectTypeOrder } from '@/lib/projectTypes';
import type { ProjectType } from '@/types';

interface CreateProjectFormData {
  schoolId: string;
  projectType: ProjectType;
  startDate: Date | null;
  goLiveDate: Date | null;
  notes: string;
}

interface CreateProjectModalProps {
  open: boolean;
  onClose: () => void;
  formData: CreateProjectFormData;
  onFormChange: (data: CreateProjectFormData) => void;
  onSubmit: () => void;
  isLoading: boolean;
}

const projectTypeOptions = projectTypeOrder.map((type) => ({
  value: type,
  label: projectTypeConfigs[type].label,
}));

export function CreateProjectModal({
  open,
  onClose,
  formData,
  onFormChange,
  onSubmit,
  isLoading,
}: CreateProjectModalProps) {
  const selectedConfig = projectTypeConfigs[formData.projectType];

  return (
    <Modal open={open} onClose={onClose} className="max-w-lg">
      <ModalHeader onClose={onClose}>Create New Project</ModalHeader>
      <ModalBody>
        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Project Type *
            </label>
            <Select
              value={formData.projectType}
              onChange={(value) =>
                onFormChange({ ...formData, projectType: value as ProjectType })
              }
              options={projectTypeOptions}
              placeholder="Select project type"
            />
            {selectedConfig && (
              <p className="text-xs text-gray-500 mt-1">
                {selectedConfig.description}
              </p>
            )}
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              School ID *
            </label>
            <Input
              value={formData.schoolId}
              onChange={(e) =>
                onFormChange({ ...formData, schoolId: e.target.value })
              }
              placeholder="Enter school ID"
            />
          </div>
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Start Date
              </label>
              <DatePicker
                value={formData.startDate}
                onChange={(date) =>
                  onFormChange({ ...formData, startDate: date })
                }
                placeholder="Select start date"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Go-Live Date
              </label>
              <DatePicker
                value={formData.goLiveDate}
                onChange={(date) =>
                  onFormChange({ ...formData, goLiveDate: date })
                }
                placeholder="Select go-live date"
              />
            </div>
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
        <Button
          onClick={onSubmit}
          disabled={!formData.schoolId || !formData.projectType || isLoading}
        >
          {isLoading ? 'Creating...' : 'Create Project'}
        </Button>
      </ModalFooter>
    </Modal>
  );
}
