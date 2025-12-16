import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import api from './client';
import { toast } from '@/lib/toast';
import type {
  Incident,
  CreateIncidentRequest,
  IncidentFilters,
  IncidentStatus,
  PaginatedResponse,
} from '@/types';

const INCIDENTS_KEY = 'incidents';

// List incidents
export function useIncidents(filters: IncidentFilters = {}) {
  return useQuery({
    queryKey: [INCIDENTS_KEY, filters],
    queryFn: () => api.get<PaginatedResponse<Incident>>('/incidents', filters),
    staleTime: 30_000, // 30 seconds
  });
}

// Get single incident
export function useIncident(id: string) {
  return useQuery({
    queryKey: [INCIDENTS_KEY, id],
    queryFn: () => api.get<Incident>(`/incidents/${id}`),
    enabled: !!id,
  });
}

// Create incident
export function useCreateIncident() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateIncidentRequest) =>
      api.post<Incident>('/incidents', data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [INCIDENTS_KEY] });
      toast.success('Incident Created', 'The incident has been reported successfully');
    },
  });
}

// Update incident status
export function useUpdateIncidentStatus() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, status }: { id: string; status: IncidentStatus }) =>
      api.patch<Incident>(`/incidents/${id}/status`, { status }),
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: [INCIDENTS_KEY] });
      queryClient.setQueryData([INCIDENTS_KEY, data.id], data);
      toast.success('Status Updated', `Incident status changed to ${data.status.replace('_', ' ')}`);
    },
  });
}

// Prefetch incident
export function usePrefetchIncident() {
  const queryClient = useQueryClient();

  return (id: string) => {
    queryClient.prefetchQuery({
      queryKey: [INCIDENTS_KEY, id],
      queryFn: () => api.get<Incident>(`/incidents/${id}`),
    });
  };
}
