// Project types
export type ProjectType =
  | 'full_installation'
  | 'device_refresh'
  | 'support'
  | 'repair'
  | 'training';

export type ProjectStatus = 'active' | 'paused' | 'completed';

// All possible phases across all project types
export type PhaseType =
  // Full Installation phases
  | 'demo'
  | 'survey'
  | 'procurement'
  | 'install'
  | 'integrate'
  | 'commission'
  | 'ops'
  // Device Refresh phases
  | 'assessment'
  | 'deployment'
  | 'verification'
  // Support phases
  | 'onboarding'
  | 'active'
  | 'renewal'
  // Repair phases
  | 'intake'
  | 'diagnosis'
  | 'repair'
  | 'testing'
  | 'handover'
  // Training phases
  | 'planning'
  | 'delivery'
  | 'certification';

export type PhaseStatus = 'pending' | 'in_progress' | 'blocked' | 'done';

// Project type configuration
export interface ProjectTypeConfig {
  type: ProjectType;
  label: string;
  description: string;
  phases: PhaseType[];
  defaultPhase: PhaseType;
}

export interface SchoolServiceProject {
  id: string;
  tenantId: string;
  schoolId: string;
  projectType: ProjectType;
  status: ProjectStatus;
  currentPhase: PhaseType;
  startDate: string;
  goLiveDate: string;
  accountManagerUserId: string;
  notes: string;
  createdAt: string;
  updatedAt: string;
}

export interface CreateProjectRequest {
  schoolId: string;
  projectType?: ProjectType;
  startDate?: string;
  goLiveDate?: string;
  accountManagerUserId?: string;
  notes?: string;
}

export interface ServicePhase {
  id: string;
  tenantId: string;
  projectId: string;
  phaseType: PhaseType;
  status: PhaseStatus;
  ownerRole: string;
  ownerUserId: string;
  ownerUserName: string;
  statusChangedAt?: string;
  statusChangedByUserId: string;
  statusChangedByUserName: string;
  startDate: string;
  endDate: string;
  notes: string;
  createdAt: string;
  updatedAt: string;
}

// Team management types
export type TeamMemberRole = 'owner' | 'collaborator' | 'viewer';

export interface ProjectTeamMember {
  id: string;
  tenantId: string;
  projectId: string;
  userId: string;
  userEmail: string;
  userName: string;
  role: TeamMemberRole;
  assignedPhases: PhaseType[];
  responsibility: string;
  assignedByUserId: string;
  assignedAt: string;
  createdAt: string;
  updatedAt: string;
}

export interface AddTeamMemberRequest {
  userId: string;
  userEmail?: string;
  userName?: string;
  role?: TeamMemberRole;
  assignedPhases?: PhaseType[];
  responsibility?: string;
}

export interface UpdateTeamMemberRequest {
  role?: TeamMemberRole;
  assignedPhases?: PhaseType[];
  responsibility?: string;
}

export type PhaseAssignmentType = 'lead' | 'collaborator' | 'reviewer';

export interface PhaseUserAssignment {
  id: string;
  tenantId: string;
  phaseId: string;
  projectId: string;
  userId: string;
  userEmail: string;
  userName: string;
  assignmentType: PhaseAssignmentType;
  assignedByUserId: string;
  assignedAt: string;
  createdAt: string;
}

// Activity feed types
export type ActivityType =
  | 'comment'
  | 'note'
  | 'file_upload'
  | 'status_change'
  | 'assignment'
  | 'work_order'
  | 'phase_transition'
  | 'mention';

export type ActivityVisibility = 'team' | 'public' | 'private';

export interface ProjectActivity {
  id: string;
  tenantId: string;
  projectId: string;
  phaseId?: string;
  workOrderId?: string;
  activityType: ActivityType;
  actorUserId: string;
  actorEmail: string;
  actorName: string;
  content: string;
  metadata: Record<string, unknown>;
  attachmentIds: string[];
  visibility: ActivityVisibility;
  isPinned: boolean;
  editedAt?: string;
  createdAt: string;
}

export interface CreateActivityRequest {
  phaseId?: string;
  activityType?: ActivityType;
  content: string;
  visibility?: ActivityVisibility;
  mentions?: string[];
}

export interface UpdateActivityRequest {
  content: string;
}

export interface ProjectAttachment {
  id: string;
  tenantId: string;
  projectId: string;
  phaseId?: string;
  activityId?: string;
  fileName: string;
  contentType: string;
  sizeBytes: number;
  objectKey: string;
  uploadedByUserId: string;
  createdAt: string;
}

// User notification types
export type UserNotificationType =
  | 'assignment'
  | 'mention'
  | 'status_change'
  | 'comment';

export interface UserNotification {
  id: string;
  tenantId: string;
  userId: string;
  notificationType: UserNotificationType;
  entityType: string;
  entityId: string;
  projectId?: string;
  title: string;
  body: string;
  metadata: Record<string, unknown>;
  isRead: boolean;
  readAt?: string;
  createdAt: string;
}

export interface MarkMultipleReadRequest {
  notificationIds: string[];
}

// Survey types
export type SurveyStatus = 'draft' | 'submitted' | 'approved';

export interface SiteSurvey {
  id: string;
  tenantId: string;
  projectId: string;
  status: SurveyStatus;
  conductedByUserId: string;
  conductedAt: string | null;
  summary: string;
  risks: string;
  createdAt: string;
  updatedAt: string;
}

export interface SurveyRoom {
  id: string;
  tenantId: string;
  surveyId: string;
  name: string;
  roomType: string;
  floor: string;
  powerNotes: string;
  networkNotes: string;
  createdAt: string;
}

export interface SurveyPhoto {
  id: string;
  tenantId: string;
  surveyId: string;
  roomId: string;
  attachmentId: string;
  caption: string;
  createdAt: string;
}
