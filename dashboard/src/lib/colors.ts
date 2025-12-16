/**
 * Centralized color definitions for consistent styling across the application.
 *
 * Usage:
 * ```tsx
 * import { incidentSeverityColors, incidentStatusColors } from '@/lib/colors';
 *
 * <Badge className={incidentSeverityColors[severity]}>{severity}</Badge>
 * ```
 */

// ============================================================================
// Incident Colors
// ============================================================================

export type IncidentSeverity = 'low' | 'medium' | 'high' | 'critical';
export type IncidentStatus = 'new' | 'acknowledged' | 'in_progress' | 'escalated' | 'resolved' | 'closed';

export const incidentSeverityColors: Record<IncidentSeverity, string> = {
  low: 'bg-green-100 text-green-800',
  medium: 'bg-yellow-100 text-yellow-800',
  high: 'bg-orange-100 text-orange-800',
  critical: 'bg-red-100 text-red-800',
};

export const incidentStatusColors: Record<IncidentStatus, string> = {
  new: 'bg-blue-100 text-blue-800',
  acknowledged: 'bg-yellow-100 text-yellow-800',
  in_progress: 'bg-purple-100 text-purple-800',
  escalated: 'bg-red-100 text-red-800',
  resolved: 'bg-green-100 text-green-800',
  closed: 'bg-gray-100 text-gray-800',
};

// ============================================================================
// Work Order Colors
// ============================================================================

export type WorkOrderStatus = 'draft' | 'assigned' | 'in_repair' | 'qa' | 'completed' | 'cancelled';

export const workOrderStatusColors: Record<WorkOrderStatus, string> = {
  draft: 'bg-gray-100 text-gray-700 border border-gray-200',
  assigned: 'bg-blue-100 text-blue-700 border border-blue-200',
  in_repair: 'bg-amber-100 text-amber-700 border border-amber-200',
  qa: 'bg-purple-100 text-purple-700 border border-purple-200',
  completed: 'bg-emerald-100 text-emerald-700 border border-emerald-200',
  cancelled: 'bg-red-100 text-red-700 border border-red-200',
};

// ============================================================================
// Device Colors
// ============================================================================

export type DeviceStatus = 'active' | 'inactive' | 'repair' | 'disposed';

export const deviceStatusColors: Record<DeviceStatus, string> = {
  active: 'bg-green-100 text-green-800',
  inactive: 'bg-gray-100 text-gray-800',
  repair: 'bg-yellow-100 text-yellow-800',
  disposed: 'bg-red-100 text-red-800',
};

// ============================================================================
// Project Colors
// ============================================================================

export type ProjectStatus = 'active' | 'paused' | 'completed';

export const projectStatusColors: Record<ProjectStatus, string> = {
  active: 'bg-green-100 text-green-800',
  paused: 'bg-yellow-100 text-yellow-800',
  completed: 'bg-gray-100 text-gray-800',
};

// ============================================================================
// Knowledge Base Colors
// ============================================================================

export type KBArticleStatus = 'draft' | 'published' | 'archived';

export const kbArticleStatusColors: Record<KBArticleStatus, string> = {
  draft: 'bg-amber-100 text-amber-800',
  published: 'bg-green-100 text-green-800',
  archived: 'bg-gray-100 text-gray-800',
};

// ============================================================================
// Marketing KB Colors
// ============================================================================

export type MKBContentType = 'talking_point' | 'objection_handler' | 'case_study' | 'faq' | 'script' | 'template';
export type MKBArticleStatus = 'draft' | 'in_review' | 'approved' | 'archived';

export const mkbContentTypeColors: Record<MKBContentType, string> = {
  talking_point: 'bg-blue-100 text-blue-800',
  objection_handler: 'bg-orange-100 text-orange-800',
  case_study: 'bg-purple-100 text-purple-800',
  faq: 'bg-green-100 text-green-800',
  script: 'bg-indigo-100 text-indigo-800',
  template: 'bg-cyan-100 text-cyan-800',
};

export const mkbArticleStatusColors: Record<MKBArticleStatus, string> = {
  draft: 'bg-amber-100 text-amber-800',
  in_review: 'bg-blue-100 text-blue-800',
  approved: 'bg-green-100 text-green-800',
  archived: 'bg-gray-100 text-gray-800',
};

// ============================================================================
// Demo Pipeline / Sales Colors (re-exported from types/sales.ts for convenience)
// ============================================================================

export type DemoLeadStage =
  | 'new_lead'
  | 'contacted'
  | 'demo_scheduled'
  | 'demo_completed'
  | 'proposal_sent'
  | 'negotiation'
  | 'won'
  | 'lost';

export const demoLeadStageColors: Record<DemoLeadStage, string> = {
  new_lead: 'bg-blue-100 text-blue-800',
  contacted: 'bg-yellow-100 text-yellow-800',
  demo_scheduled: 'bg-purple-100 text-purple-800',
  demo_completed: 'bg-indigo-100 text-indigo-800',
  proposal_sent: 'bg-orange-100 text-orange-800',
  negotiation: 'bg-pink-100 text-pink-800',
  won: 'bg-green-100 text-green-800',
  lost: 'bg-red-100 text-red-800',
};

// ============================================================================
// Generic / Common Colors
// ============================================================================

export const priorityColors = {
  low: 'bg-green-100 text-green-800',
  medium: 'bg-yellow-100 text-yellow-800',
  high: 'bg-orange-100 text-orange-800',
  critical: 'bg-red-100 text-red-800',
  urgent: 'bg-red-100 text-red-800',
} as const;

export const booleanColors = {
  true: 'bg-green-100 text-green-800',
  false: 'bg-gray-100 text-gray-800',
  yes: 'bg-green-100 text-green-800',
  no: 'bg-gray-100 text-gray-800',
} as const;

// ============================================================================
// Helper Functions
// ============================================================================

/**
 * Get a status color with a fallback for unknown values.
 */
export function getStatusColor(
  colors: Record<string, string>,
  status: string,
  fallback = 'bg-gray-100 text-gray-800'
): string {
  return colors[status] || fallback;
}
