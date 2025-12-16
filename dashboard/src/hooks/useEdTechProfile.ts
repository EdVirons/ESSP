import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { edtechApi } from '@/api/edtech';
import type {
  SaveProfileRequest,
  SubmitFollowUpRequest,
} from '@/types/edtech';

// Query keys
export const edtechKeys = {
  all: ['edtech'] as const,
  options: () => [...edtechKeys.all, 'options'] as const,
  profile: (schoolId: string) => [...edtechKeys.all, 'profile', schoolId] as const,
  history: (schoolId: string) => [...edtechKeys.all, 'history', schoolId] as const,
};

// Get form options
export function useEdTechOptions() {
  return useQuery({
    queryKey: edtechKeys.options(),
    queryFn: () => edtechApi.getOptions(),
    staleTime: 24 * 60 * 60 * 1000, // Options don't change often - cache for 24 hours
  });
}

// Get profile by school ID
export function useEdTechProfile(schoolId: string | undefined, enabled = true) {
  return useQuery({
    queryKey: edtechKeys.profile(schoolId || ''),
    queryFn: () => edtechApi.getBySchoolId(schoolId!),
    enabled: enabled && !!schoolId,
  });
}

// Save profile (create or update)
export function useSaveEdTechProfile() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: SaveProfileRequest) => edtechApi.saveProfile(data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: edtechKeys.profile(variables.schoolId) });
    },
  });
}

// Generate AI analysis
export function useGenerateAI() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (params: { profileId: string; schoolId: string }) =>
      edtechApi.generateAI(params.profileId),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: edtechKeys.profile(variables.schoolId) });
    },
  });
}

// Submit follow-up responses
export function useSubmitFollowUp() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (params: { profileId: string; schoolId: string; data: SubmitFollowUpRequest }) =>
      edtechApi.submitFollowUp(params.profileId, params.data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: edtechKeys.profile(variables.schoolId) });
    },
  });
}

// Complete profile
export function useCompleteProfile() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (params: { profileId: string; schoolId: string }) =>
      edtechApi.complete(params.profileId),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: edtechKeys.profile(variables.schoolId) });
    },
  });
}

// Get profile history
export function useEdTechHistory(schoolId: string | undefined, limit?: number, enabled = true) {
  return useQuery({
    queryKey: edtechKeys.history(schoolId || ''),
    queryFn: () => edtechApi.getHistory(schoolId!, limit),
    enabled: enabled && !!schoolId,
  });
}
