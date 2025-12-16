import type {
  DeviceModel,
  DeviceCategory,
  CreateDeviceModelInput,
  UpdateDeviceModelInput,
  DeviceModelSpecs,
} from '@/types/device';

export interface DeviceModelManagerProps {
  open: boolean;
  onClose: () => void;
  models: DeviceModel[];
  isLoading: boolean;
  onCreate: (data: CreateDeviceModelInput) => Promise<void>;
  onUpdate: (id: string, data: UpdateDeviceModelInput) => Promise<void>;
  onDelete: (id: string) => Promise<void>;
}

export interface ModelFormData {
  make: string;
  model: string;
  category: DeviceCategory;
  specs: DeviceModelSpecs;
}

export const initialFormData: ModelFormData = {
  make: '',
  model: '',
  category: 'laptop',
  specs: {},
};

export const specLabels: Array<{ key: string; label: string }> = [
  { key: 'processor', label: 'Processor' },
  { key: 'ram', label: 'RAM' },
  { key: 'storage', label: 'Storage' },
  { key: 'display', label: 'Display' },
  { key: 'os', label: 'Operating System' },
  { key: 'battery', label: 'Battery' },
  { key: 'weight', label: 'Weight' },
];

export type { DeviceModel, DeviceCategory, CreateDeviceModelInput, UpdateDeviceModelInput };
