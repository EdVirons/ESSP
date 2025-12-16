import type { ProjectType, PhaseType, ProjectTypeConfig } from '@/types';

// Phase labels for all phase types
export const phaseTypeLabels: Record<PhaseType, string> = {
  // Full Installation phases
  demo: 'Demo',
  survey: 'Site Survey',
  procurement: 'Procurement',
  install: 'Installation',
  integrate: 'Integration',
  commission: 'Commissioning',
  ops: 'Operations',
  // Device Refresh phases
  assessment: 'Assessment',
  deployment: 'Deployment',
  verification: 'Verification',
  // Support phases
  onboarding: 'Onboarding',
  active: 'Active',
  renewal: 'Renewal',
  // Repair phases
  intake: 'Intake',
  diagnosis: 'Diagnosis',
  repair: 'Repair',
  testing: 'Testing',
  handover: 'Handover',
  // Training phases
  planning: 'Planning',
  delivery: 'Delivery',
  certification: 'Certification',
};

// Phase status colors
export const phaseStatusColors: Record<string, string> = {
  pending: 'bg-gray-100 text-gray-800',
  in_progress: 'bg-blue-100 text-blue-800',
  blocked: 'bg-red-100 text-red-800',
  done: 'bg-green-100 text-green-800',
};

// Project type configurations
export const projectTypeConfigs: Record<ProjectType, ProjectTypeConfig> = {
  full_installation: {
    type: 'full_installation',
    label: 'Full Installation',
    description: 'Complete school technology installation',
    phases: ['demo', 'survey', 'procurement', 'install', 'integrate', 'commission', 'ops'],
    defaultPhase: 'demo',
  },
  device_refresh: {
    type: 'device_refresh',
    label: 'Device Refresh',
    description: 'Upgrade or replace existing devices',
    phases: ['assessment', 'procurement', 'deployment', 'verification'],
    defaultPhase: 'assessment',
  },
  support: {
    type: 'support',
    label: 'Support',
    description: 'Ongoing technical support contract',
    phases: ['onboarding', 'active', 'renewal'],
    defaultPhase: 'onboarding',
  },
  repair: {
    type: 'repair',
    label: 'Repair',
    description: 'Device repair workflow',
    phases: ['intake', 'diagnosis', 'repair', 'testing', 'handover'],
    defaultPhase: 'intake',
  },
  training: {
    type: 'training',
    label: 'Training',
    description: 'Staff training program',
    phases: ['planning', 'delivery', 'assessment', 'certification'],
    defaultPhase: 'planning',
  },
};

// Project type order for tabs
export const projectTypeOrder: ProjectType[] = [
  'full_installation',
  'device_refresh',
  'support',
  'repair',
  'training',
];

// Project type colors for badges and icons
export const projectTypeColors: Record<ProjectType, { bg: string; text: string; icon: string }> = {
  full_installation: { bg: 'bg-purple-50', text: 'text-purple-600', icon: 'purple' },
  device_refresh: { bg: 'bg-blue-50', text: 'text-blue-600', icon: 'blue' },
  support: { bg: 'bg-green-50', text: 'text-green-600', icon: 'green' },
  repair: { bg: 'bg-orange-50', text: 'text-orange-600', icon: 'orange' },
  training: { bg: 'bg-indigo-50', text: 'text-indigo-600', icon: 'indigo' },
};

// Helper to get phases for a project type
export function getPhasesForType(projectType: ProjectType): PhaseType[] {
  return projectTypeConfigs[projectType]?.phases || projectTypeConfigs.full_installation.phases;
}

// Helper to validate phase for project type
export function isValidPhaseForType(projectType: ProjectType, phase: PhaseType): boolean {
  const config = projectTypeConfigs[projectType];
  return config?.phases.includes(phase) ?? false;
}

// Helper to get project type label
export function getProjectTypeLabel(projectType: ProjectType): string {
  return projectTypeConfigs[projectType]?.label || projectType;
}

// Helper to get phase label
export function getPhaseLabel(phase: PhaseType): string {
  return phaseTypeLabels[phase] || phase;
}
