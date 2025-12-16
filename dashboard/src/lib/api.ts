import axios from 'axios';
import type { AxiosError, AxiosInstance } from 'axios';

// API base URLs
export const API_BASE_URL = '/v1';
export const ADMIN_API_BASE_URL = '/v1';

// Create axios instance for admin API
export const adminApi: AxiosInstance = axios.create({
  baseURL: ADMIN_API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
  withCredentials: true, // Include cookies in requests
});

// Create axios instance for regular API
export const api: AxiosInstance = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
  withCredentials: true, // Include cookies in requests
});

// Response interceptor for handling auth errors
const handleAuthError = (error: AxiosError) => {
  if (error.response?.status === 401) {
    // Redirect to login page if unauthorized
    const currentPath = window.location.pathname;
    if (!currentPath.includes('/login')) {
      window.location.href = '/login';
    }
  }
  return Promise.reject(error);
};

// Apply interceptor to both API instances
adminApi.interceptors.response.use((response) => response, handleAuthError);
api.interceptors.response.use((response) => response, handleAuthError);

// Auth API types
export interface LoginRequest {
  username: string;
  password: string;
}

export interface User {
  username: string;
  roles: string[];
  email?: string;
  displayName?: string;
  avatarUrl?: string;
  tenantId?: string;
  schoolId?: string;
}

export interface LoginResponse {
  success: boolean;
  message?: string;
  user?: User;
}

export interface MeResponse {
  authenticated: boolean;
  user?: User;
}

// SSO User Profile types for Edvirons ecosystem
export interface Organization {
  id: string;
  name: string;
  displayName?: string;
  type?: string; // e.g., "district", "school", "service_provider"
  logoUrl?: string;
}

export interface NotificationPreferences {
  emailEnabled: boolean;
  browserEnabled: boolean;
  incidentAlerts: boolean;
  workOrderAlerts: boolean;
}

export interface UserPreferences {
  theme?: string;
  language?: string;
  timezone?: string;
  sidebarCollapsed?: boolean;
  notifications?: NotificationPreferences;
}

export interface SSOUserProfile {
  id: string;
  username: string;
  email: string;
  displayName: string;
  firstName?: string;
  lastName?: string;
  avatarUrl?: string;
  organization?: Organization;
  tenantId: string;
  roles: string[];
  permissions?: string[];
  ssoProvider?: string;
  ssoSubject?: string;
  emailVerified: boolean;
  lastLoginAt?: string;
  createdAt?: string;
  preferences?: UserPreferences;
}

export interface ProfileResponse {
  profile: SSOUserProfile;
}

// Auth API functions
export const authApi = {
  login: async (credentials: LoginRequest): Promise<LoginResponse> => {
    const response = await adminApi.post<LoginResponse>('/auth/login', credentials);
    return response.data;
  },

  logout: async (): Promise<void> => {
    await adminApi.post('/auth/logout');
  },

  me: async (): Promise<MeResponse> => {
    const response = await adminApi.get<MeResponse>('/auth/me');
    return response.data;
  },

  refresh: async (): Promise<void> => {
    await adminApi.post('/auth/refresh');
  },

  // Get full SSO user profile
  profile: async (): Promise<ProfileResponse> => {
    const response = await adminApi.get<ProfileResponse>('/auth/profile');
    return response.data;
  },
};

// Export default api instance
export default api;
