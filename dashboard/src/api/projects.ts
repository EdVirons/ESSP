import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import api from './client';
import type {
  SchoolServiceProject,
  CreateProjectRequest,
  ServicePhase,
  SiteSurvey,
  ProjectStatus,
  ProjectType,
  ProjectTypeConfig,
  PhaseStatus,
  PaginatedResponse,
  ProjectTeamMember,
  AddTeamMemberRequest,
  UpdateTeamMemberRequest,
  ProjectActivity,
  CreateActivityRequest,
  UpdateActivityRequest,
  ProjectAttachment,
  PhaseUserAssignment,
  UserNotification,
  MarkMultipleReadRequest,
} from '@/types';
import type { WorkOrder } from '@/types/work-order';

const PROJECTS_KEY = 'projects';
const PHASES_KEY = 'phases';
const SURVEYS_KEY = 'surveys';
const TEAM_KEY = 'team';
const ACTIVITIES_KEY = 'activities';
const WORK_ORDERS_KEY = 'workOrders';
const NOTIFICATIONS_KEY = 'userNotifications';

// Projects
export function useProjects(filters: {
  projectType?: ProjectType;
  status?: ProjectStatus;
  limit?: number;
  cursor?: string;
} = {}) {
  return useQuery({
    queryKey: [PROJECTS_KEY, filters],
    queryFn: () => api.get<PaginatedResponse<SchoolServiceProject>>('/projects', filters),
    staleTime: 30_000,
  });
}

// Project type counts for tabs
export function useProjectTypeCounts() {
  return useQuery({
    queryKey: [PROJECTS_KEY, 'counts'],
    queryFn: () => api.get<Record<ProjectType, number>>('/projects/counts'),
    staleTime: 30_000,
  });
}

// Project types configuration (from API)
export function useProjectTypes() {
  return useQuery({
    queryKey: [PROJECTS_KEY, 'types'],
    queryFn: () => api.get<{ items: ProjectTypeConfig[] }>('/projects/types'),
    staleTime: 300_000, // 5 minutes - rarely changes
  });
}

export function useProject(id: string) {
  return useQuery({
    queryKey: [PROJECTS_KEY, id],
    queryFn: () => api.get<SchoolServiceProject>(`/projects/${id}`),
    enabled: !!id,
  });
}

export function useCreateProject() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateProjectRequest) =>
      api.post<SchoolServiceProject>('/projects', data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [PROJECTS_KEY] });
    },
  });
}

// Phases
export function usePhases(projectId: string) {
  return useQuery({
    queryKey: [PROJECTS_KEY, projectId, PHASES_KEY],
    queryFn: () => api.get<PaginatedResponse<ServicePhase>>(`/projects/${projectId}/phases`),
    enabled: !!projectId,
  });
}

export function useCreatePhase() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      projectId,
      phaseType,
      ownerRole,
      startDate,
      endDate,
      notes,
    }: {
      projectId: string;
      phaseType: string;
      ownerRole?: string;
      startDate?: string;
      endDate?: string;
      notes?: string;
    }) =>
      api.post<ServicePhase>(`/projects/${projectId}/phases`, {
        phaseType,
        ownerRole,
        startDate,
        endDate,
        notes,
      }),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: [PROJECTS_KEY, variables.projectId, PHASES_KEY],
      });
    },
  });
}

export function useUpdatePhaseStatus() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      phaseId,
      status,
    }: {
      phaseId: string;
      projectId: string;
      status: PhaseStatus;
    }) =>
      api.patch<ServicePhase>(`/phases/${phaseId}/status`, { status }),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: [PROJECTS_KEY, variables.projectId, PHASES_KEY],
      });
    },
  });
}

// Surveys
export function useSurveys(projectId: string) {
  return useQuery({
    queryKey: [PROJECTS_KEY, projectId, SURVEYS_KEY],
    queryFn: () => api.get<PaginatedResponse<SiteSurvey>>(`/projects/${projectId}/surveys`),
    enabled: !!projectId,
  });
}

export function useSurvey(surveyId: string) {
  return useQuery({
    queryKey: [SURVEYS_KEY, surveyId],
    queryFn: () => api.get<SiteSurvey>(`/surveys/${surveyId}`),
    enabled: !!surveyId,
  });
}

export function useCreateSurvey() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      projectId,
      summary,
      risks,
    }: {
      projectId: string;
      summary?: string;
      risks?: string;
    }) =>
      api.post<SiteSurvey>(`/projects/${projectId}/surveys`, {
        summary,
        risks,
      }),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: [PROJECTS_KEY, variables.projectId, SURVEYS_KEY],
      });
    },
  });
}

export function useAddSurveyRoom() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      surveyId,
      name,
      roomType,
      floor,
      powerNotes,
      networkNotes,
    }: {
      surveyId: string;
      name: string;
      roomType?: string;
      floor?: string;
      powerNotes?: string;
      networkNotes?: string;
    }) =>
      api.post(`/surveys/${surveyId}/rooms`, {
        name,
        roomType,
        floor,
        powerNotes,
        networkNotes,
      }),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: [SURVEYS_KEY, variables.surveyId],
      });
    },
  });
}

// ==================== PROJECT TEAM ====================

export function useProjectTeam(projectId: string) {
  return useQuery({
    queryKey: [PROJECTS_KEY, projectId, TEAM_KEY],
    queryFn: () =>
      api.get<{ members: ProjectTeamMember[]; total: number }>(
        `/projects/${projectId}/team`
      ),
    enabled: !!projectId,
  });
}

export function useMyProjects() {
  return useQuery({
    queryKey: [PROJECTS_KEY, 'my-projects'],
    queryFn: () =>
      api.get<{
        memberships: ProjectTeamMember[];
        projectIds: string[];
        total: number;
      }>('/users/me/projects'),
  });
}

export function useAddTeamMember() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      projectId,
      data,
    }: {
      projectId: string;
      data: AddTeamMemberRequest;
    }) => api.post<ProjectTeamMember>(`/projects/${projectId}/team`, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: [PROJECTS_KEY, variables.projectId, TEAM_KEY],
      });
    },
  });
}

export function useUpdateTeamMember() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      projectId,
      memberId,
      data,
    }: {
      projectId: string;
      memberId: string;
      data: UpdateTeamMemberRequest;
    }) =>
      api.patch<ProjectTeamMember>(
        `/projects/${projectId}/team/${memberId}`,
        data
      ),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: [PROJECTS_KEY, variables.projectId, TEAM_KEY],
      });
    },
  });
}

export function useRemoveTeamMember() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      projectId,
      memberId,
    }: {
      projectId: string;
      memberId: string;
    }) => api.delete(`/projects/${projectId}/team/${memberId}`),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: [PROJECTS_KEY, variables.projectId, TEAM_KEY],
      });
    },
  });
}

// Phase Assignments
export function usePhaseAssignments(phaseId: string) {
  return useQuery({
    queryKey: [PHASES_KEY, phaseId, 'assignments'],
    queryFn: () =>
      api.get<{ assignments: PhaseUserAssignment[]; total: number }>(
        `/phases/${phaseId}/assignments`
      ),
    enabled: !!phaseId,
  });
}

export function useAddPhaseAssignment() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      phaseId,
      userId,
      userEmail,
      userName,
      assignmentType,
    }: {
      phaseId: string;
      userId: string;
      userEmail?: string;
      userName?: string;
      assignmentType?: string;
    }) =>
      api.post<PhaseUserAssignment>(`/phases/${phaseId}/assignments`, {
        userId,
        userEmail,
        userName,
        assignmentType,
      }),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: [PHASES_KEY, variables.phaseId, 'assignments'],
      });
    },
  });
}

export function useRemovePhaseAssignment() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      phaseId,
      assignmentId,
    }: {
      phaseId: string;
      assignmentId: string;
    }) => api.delete(`/phases/${phaseId}/assignments/${assignmentId}`),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: [PHASES_KEY, variables.phaseId, 'assignments'],
      });
    },
  });
}

// ==================== PROJECT ACTIVITIES ====================

export function useProjectActivities(
  projectId: string,
  filters: {
    phaseId?: string;
    type?: string;
    userId?: string;
    pinned?: boolean;
    limit?: number;
    cursor?: string;
  } = {}
) {
  return useQuery({
    queryKey: [PROJECTS_KEY, projectId, ACTIVITIES_KEY, filters],
    queryFn: () =>
      api.get<{ items: ProjectActivity[]; nextCursor: string }>(
        `/projects/${projectId}/activities`,
        filters
      ),
    enabled: !!projectId,
  });
}

export function useCreateActivity() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      projectId,
      data,
    }: {
      projectId: string;
      data: CreateActivityRequest;
    }) => api.post<ProjectActivity>(`/projects/${projectId}/activities`, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: [PROJECTS_KEY, variables.projectId, ACTIVITIES_KEY],
      });
    },
  });
}

export function useUpdateActivity() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      activityId,
      data,
    }: {
      activityId: string;
      projectId: string;
      data: UpdateActivityRequest;
    }) => api.patch<ProjectActivity>(`/activities/${activityId}`, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: [PROJECTS_KEY, variables.projectId, ACTIVITIES_KEY],
      });
    },
  });
}

export function useDeleteActivity() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      activityId,
    }: {
      activityId: string;
      projectId: string;
    }) => api.delete(`/activities/${activityId}`),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: [PROJECTS_KEY, variables.projectId, ACTIVITIES_KEY],
      });
    },
  });
}

export function useToggleActivityPin() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      activityId,
    }: {
      activityId: string;
      projectId: string;
    }) => api.post<{ id: string; isPinned: boolean }>(`/activities/${activityId}/pin`),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: [PROJECTS_KEY, variables.projectId, ACTIVITIES_KEY],
      });
    },
  });
}

// Attachments
export function useProjectAttachments(projectId: string, limit?: number) {
  return useQuery({
    queryKey: [PROJECTS_KEY, projectId, 'attachments'],
    queryFn: () =>
      api.get<{ items: ProjectAttachment[]; total: number }>(
        `/projects/${projectId}/attachments`,
        { limit }
      ),
    enabled: !!projectId,
  });
}

export function useDeleteAttachment() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      attachmentId,
    }: {
      attachmentId: string;
      projectId: string;
    }) => api.delete(`/attachments/${attachmentId}`),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: [PROJECTS_KEY, variables.projectId, 'attachments'],
      });
      queryClient.invalidateQueries({
        queryKey: [PROJECTS_KEY, variables.projectId, ACTIVITIES_KEY],
      });
    },
  });
}

// ==================== PROJECT WORK ORDERS ====================

export function useProjectWorkOrders(
  projectId: string,
  filters: {
    phaseId?: string;
    status?: string;
    limit?: number;
    cursor?: string;
  } = {}
) {
  return useQuery({
    queryKey: [PROJECTS_KEY, projectId, WORK_ORDERS_KEY, filters],
    queryFn: () =>
      api.get<{ items: WorkOrder[]; nextCursor: string }>(
        `/projects/${projectId}/work-orders`,
        filters
      ),
    enabled: !!projectId,
  });
}

export function useCreateProjectWorkOrder() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      projectId,
      data,
    }: {
      projectId: string;
      data: {
        phaseId?: string;
        deviceId?: string;
        taskType: string;
        serviceShopId?: string;
        assignedStaffId?: string;
        repairLocation?: string;
        assignedTo?: string;
        costEstimateCents?: number;
        notes?: string;
      };
    }) => api.post<WorkOrder>(`/projects/${projectId}/work-orders`, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: [PROJECTS_KEY, variables.projectId, WORK_ORDERS_KEY],
      });
      queryClient.invalidateQueries({
        queryKey: [PROJECTS_KEY, variables.projectId, ACTIVITIES_KEY],
      });
    },
  });
}

// ==================== USER NOTIFICATIONS ====================

export function useUserNotifications(filters: {
  unread?: boolean;
  projectId?: string;
  limit?: number;
  cursor?: string;
} = {}) {
  return useQuery({
    queryKey: [NOTIFICATIONS_KEY, filters],
    queryFn: () =>
      api.get<{ items: UserNotification[]; nextCursor: string }>(
        '/users/me/notifications',
        filters
      ),
  });
}

export function useUnreadNotificationCount() {
  return useQuery({
    queryKey: [NOTIFICATIONS_KEY, 'unread-count'],
    queryFn: () =>
      api.get<{ unreadCount: number }>('/users/me/notifications/unread-count'),
    refetchInterval: 30_000, // Poll every 30 seconds
  });
}

export function useMarkNotificationAsRead() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (notificationId: string) =>
      api.post<{ id: string; isRead: boolean }>(
        `/users/me/notifications/${notificationId}/read`
      ),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [NOTIFICATIONS_KEY] });
    },
  });
}

export function useMarkAllNotificationsAsRead() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: () =>
      api.post<{ markedCount: number }>('/users/me/notifications/read-all'),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [NOTIFICATIONS_KEY] });
    },
  });
}

export function useMarkMultipleNotificationsAsRead() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: MarkMultipleReadRequest) =>
      api.post<{ markedCount: number }>('/users/me/notifications/read', data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [NOTIFICATIONS_KEY] });
    },
  });
}
