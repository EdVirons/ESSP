import axios, { type AxiosInstance, type AxiosError, type InternalAxiosRequestConfig } from 'axios';
import { API_BASE_URL } from '@/lib/constants';
import { toast } from '@/lib/toast';

// Create axios instance
const apiClient: AxiosInstance = axios.create({
  baseURL: API_BASE_URL,
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor to add auth token and tenant headers
apiClient.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    // Get token from localStorage
    const token = localStorage.getItem('auth_token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }

    // Get tenant and school from localStorage or context
    const tenantId = localStorage.getItem('tenant_id') || 'demo-tenant';
    const schoolId = localStorage.getItem('school_id') || 'demo-school';

    config.headers['X-Tenant-ID'] = tenantId;
    config.headers['X-School-ID'] = schoolId;

    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Response interceptor for error handling
apiClient.interceptors.response.use(
  (response) => response,
  async (error: AxiosError) => {
    const originalRequest = error.config;

    // Handle 401 - token expired
    if (error.response?.status === 401 && originalRequest) {
      // Try to refresh token
      const refreshToken = localStorage.getItem('refresh_token');
      if (refreshToken) {
        try {
          const response = await axios.post(`${API_BASE_URL}/auth/refresh`, {
            refreshToken,
          });
          const { accessToken } = response.data;
          localStorage.setItem('auth_token', accessToken);

          // Retry original request
          originalRequest.headers.Authorization = `Bearer ${accessToken}`;
          return apiClient(originalRequest);
        } catch {
          // Refresh failed, redirect to login
          localStorage.removeItem('auth_token');
          localStorage.removeItem('refresh_token');
          window.location.href = '/login';
        }
      } else {
        // No refresh token, redirect to login
        localStorage.removeItem('auth_token');
        window.location.href = '/login';
      }
    }

    // Handle other errors
    const message = getErrorMessage(error);
    console.error('API Error:', message);

    // Show error toast (skip for 401 handled above)
    if (error.response?.status !== 401) {
      toast.error('Request Failed', message);
    }

    return Promise.reject(error);
  }
);

function getErrorMessage(error: AxiosError): string {
  if (error.response) {
    // Server responded with error
    const data = error.response.data;
    if (typeof data === 'string') return data;
    if (typeof data === 'object' && data !== null && 'message' in data) {
      return (data as { message: string }).message;
    }
    return `Error ${error.response.status}`;
  }
  if (error.request) {
    // No response received
    return 'Network error - please check your connection';
  }
  return error.message || 'An unexpected error occurred';
}

export { apiClient };

// Convenience methods
export const api = {
  get: <T>(url: string, params?: object) =>
    apiClient.get<T>(url, { params }).then((res) => res.data),

  post: <T>(url: string, data?: unknown) =>
    apiClient.post<T>(url, data).then((res) => res.data),

  put: <T>(url: string, data?: unknown) =>
    apiClient.put<T>(url, data).then((res) => res.data),

  patch: <T>(url: string, data?: unknown) =>
    apiClient.patch<T>(url, data).then((res) => res.data),

  delete: <T>(url: string) =>
    apiClient.delete<T>(url).then((res) => res.data),
};

export default api;
