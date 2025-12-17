import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import api from './client';
import { hrApi } from './hr-client';
import type {
  PersonSnapshot,
  TeamSnapshot,
  OrgUnitSnapshot,
  OrgTreeNode,
  TeamMembershipSnapshot,
  PersonFilters,
  TeamFilters,
  OrgUnitFilters,
  TeamMembershipFilters,
  CreatePersonInput,
  CreateTeamInput,
  CreateOrgUnitInput,
  CreateTeamMembershipInput,
} from '../types/hr';

interface PaginatedResponse<T> {
  items: T[];
  total?: number;
  limit: number;
  offset: number;
}

const PEOPLE_KEY = 'ssot-people';
const TEAMS_KEY = 'ssot-teams';
const ORG_UNITS_KEY = 'ssot-org-units';
const ORG_TREE_KEY = 'ssot-org-tree';
const TEAM_MEMBERSHIPS_KEY = 'ssot-team-memberships';

// ==================== People ====================

export function usePeople(filters: PersonFilters = {}) {
  return useQuery({
    queryKey: [PEOPLE_KEY, filters],
    queryFn: () => api.get<PaginatedResponse<PersonSnapshot>>('/ssot/people', filters),
    staleTime: 60_000,
  });
}

export function usePerson(personId: string) {
  return useQuery({
    queryKey: [PEOPLE_KEY, personId],
    queryFn: () => api.get<PersonSnapshot>(`/ssot/people/${personId}`),
    enabled: !!personId,
  });
}

export function usePersonByEmail(email: string) {
  return useQuery({
    queryKey: [PEOPLE_KEY, 'email', email],
    queryFn: () => api.get<PersonSnapshot>('/ssot/people', { email }),
    enabled: !!email,
  });
}

export function useSyncPeople() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: () => api.post<{ ok: boolean; synced: number }>('/ssot/sync/people', {}),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [PEOPLE_KEY] });
    },
  });
}

// ==================== Teams ====================

export function useTeams(filters: TeamFilters = {}) {
  return useQuery({
    queryKey: [TEAMS_KEY, filters],
    queryFn: () => api.get<PaginatedResponse<TeamSnapshot>>('/ssot/teams', filters),
    staleTime: 60_000,
  });
}

export function useTeam(teamId: string) {
  return useQuery({
    queryKey: [TEAMS_KEY, teamId],
    queryFn: () => api.get<TeamSnapshot>(`/ssot/teams/${teamId}`),
    enabled: !!teamId,
  });
}

export function useTeamByKey(key: string) {
  return useQuery({
    queryKey: [TEAMS_KEY, 'key', key],
    queryFn: () => api.get<TeamSnapshot>('/ssot/teams', { key }),
    enabled: !!key,
  });
}

export function useSyncTeams() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: () => api.post<{ ok: boolean; synced: number }>('/ssot/sync/teams', {}),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [TEAMS_KEY] });
    },
  });
}

// ==================== Org Units ====================

export function useOrgUnits(filters: OrgUnitFilters = {}) {
  return useQuery({
    queryKey: [ORG_UNITS_KEY, filters],
    queryFn: () => api.get<PaginatedResponse<OrgUnitSnapshot>>('/ssot/org-units', filters),
    staleTime: 300_000, // 5 min cache for org structure
  });
}

export function useOrgUnit(orgUnitId: string) {
  return useQuery({
    queryKey: [ORG_UNITS_KEY, orgUnitId],
    queryFn: () => api.get<OrgUnitSnapshot>(`/ssot/org-units/${orgUnitId}`),
    enabled: !!orgUnitId,
  });
}

export function useOrgTree() {
  return useQuery({
    queryKey: [ORG_TREE_KEY],
    queryFn: () => api.get<OrgTreeNode[]>('/ssot/org-units/tree'),
    staleTime: 300_000, // 5 min cache
  });
}

export function useSyncOrgUnits() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: () => api.post<{ ok: boolean; synced: number }>('/ssot/sync/org-units', {}),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [ORG_UNITS_KEY] });
      queryClient.invalidateQueries({ queryKey: [ORG_TREE_KEY] });
    },
  });
}

// ==================== Combined Sync ====================

export function useSyncAllHR() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async () => {
      const [people, teams, orgUnits] = await Promise.all([
        api.post<{ ok: boolean; synced: number }>('/ssot/sync/people', {}),
        api.post<{ ok: boolean; synced: number }>('/ssot/sync/teams', {}),
        api.post<{ ok: boolean; synced: number }>('/ssot/sync/org-units', {}),
      ]);
      return { people, teams, orgUnits };
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [PEOPLE_KEY] });
      queryClient.invalidateQueries({ queryKey: [TEAMS_KEY] });
      queryClient.invalidateQueries({ queryKey: [ORG_UNITS_KEY] });
      queryClient.invalidateQueries({ queryKey: [ORG_TREE_KEY] });
    },
  });
}

// ==================== Create Operations ====================

export function useCreatePerson() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: CreatePersonInput) =>
      hrApi.post<{ id: string; ok: boolean }>('/people', data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [PEOPLE_KEY] });
    },
  });
}

export function useUpdatePerson() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, ...data }: { id: string } & Partial<CreatePersonInput>) =>
      hrApi.patch<{ ok: boolean }>(`/people/${id}`, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [PEOPLE_KEY] });
    },
  });
}

export function useDeletePerson() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => hrApi.delete<{ ok: boolean }>(`/people/${id}`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [PEOPLE_KEY] });
    },
  });
}

export function useCreateTeam() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: CreateTeamInput) =>
      hrApi.post<{ id: string; ok: boolean }>('/teams', data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [TEAMS_KEY] });
    },
  });
}

export function useUpdateTeam() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, ...data }: { id: string } & Partial<CreateTeamInput>) =>
      hrApi.patch<{ ok: boolean }>(`/teams/${id}`, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [TEAMS_KEY] });
    },
  });
}

export function useDeleteTeam() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => hrApi.delete<{ ok: boolean }>(`/teams/${id}`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [TEAMS_KEY] });
    },
  });
}

export function useCreateOrgUnit() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: CreateOrgUnitInput) =>
      hrApi.post<{ id: string; ok: boolean }>('/org-units', data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [ORG_UNITS_KEY] });
      queryClient.invalidateQueries({ queryKey: [ORG_TREE_KEY] });
    },
  });
}

export function useUpdateOrgUnit() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, ...data }: { id: string } & Partial<CreateOrgUnitInput>) =>
      hrApi.patch<{ ok: boolean }>(`/org-units/${id}`, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [ORG_UNITS_KEY] });
      queryClient.invalidateQueries({ queryKey: [ORG_TREE_KEY] });
    },
  });
}

export function useDeleteOrgUnit() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => hrApi.delete<{ ok: boolean }>(`/org-units/${id}`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [ORG_UNITS_KEY] });
      queryClient.invalidateQueries({ queryKey: [ORG_TREE_KEY] });
    },
  });
}

// ==================== Team Memberships ====================

export function useTeamMemberships(filters: TeamMembershipFilters = {}) {
  return useQuery({
    queryKey: [TEAM_MEMBERSHIPS_KEY, filters],
    queryFn: () => hrApi.get<PaginatedResponse<TeamMembershipSnapshot>>('/team-memberships', filters as Record<string, unknown>),
    staleTime: 60_000,
  });
}

export function useTeamMembers(teamId: string) {
  return useQuery({
    queryKey: [TEAM_MEMBERSHIPS_KEY, 'team', teamId],
    queryFn: () => hrApi.get<PaginatedResponse<TeamMembershipSnapshot>>('/team-memberships', { teamId }),
    enabled: !!teamId,
    staleTime: 60_000,
  });
}

export function usePersonTeams(personId: string) {
  return useQuery({
    queryKey: [TEAM_MEMBERSHIPS_KEY, 'person', personId],
    queryFn: () => hrApi.get<PaginatedResponse<TeamMembershipSnapshot>>('/team-memberships', { personId }),
    enabled: !!personId,
    staleTime: 60_000,
  });
}

export function useCreateTeamMembership() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: CreateTeamMembershipInput) =>
      hrApi.post<{ id: string; ok: boolean }>('/team-memberships', data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [TEAM_MEMBERSHIPS_KEY] });
    },
  });
}

export function useDeleteTeamMembership() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => hrApi.delete<{ ok: boolean }>(`/team-memberships/${id}`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [TEAM_MEMBERSHIPS_KEY] });
    },
  });
}

export function useSyncTeamMemberships() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: () => api.post<{ ok: boolean; synced: number }>('/ssot/sync/team-memberships', {}),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [TEAM_MEMBERSHIPS_KEY] });
    },
  });
}
