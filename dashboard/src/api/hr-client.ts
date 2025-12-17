import axios, { type AxiosInstance, type AxiosError, type InternalAxiosRequestConfig } from 'axios';

// HR SSOT API client - for direct CRUD operations to ssot-hr service
// Read operations go through IMS-API (/v1/ssot/*), write operations go directly to ssot-hr (/hr-api/*)

const hrApiClient: AxiosInstance = axios.create({
  baseURL: '/hr-api',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor to add tenant header
hrApiClient.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const tenantId = localStorage.getItem('tenant_id') || 'demo-tenant';
    config.headers['X-Tenant-Id'] = tenantId;

    // Add auth token if available
    const token = localStorage.getItem('auth_token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }

    return config;
  },
  (error) => Promise.reject(error)
);

// Response interceptor for error handling
hrApiClient.interceptors.response.use(
  (response) => response,
  (error: AxiosError) => {
    const message = getHrErrorMessage(error);
    console.error('HR API Error:', message);

    // Don't show toast here - let the calling code handle it
    // This prevents double toasts when mutations already show errors

    return Promise.reject(new Error(message));
  }
);

function getHrErrorMessage(error: AxiosError): string {
  if (error.response) {
    const data = error.response.data;
    if (typeof data === 'string') return data;
    if (typeof data === 'object' && data !== null) {
      if ('error' in data) return (data as { error: string }).error;
      if ('message' in data) return (data as { message: string }).message;
    }
    return `Error ${error.response.status}`;
  }
  if (error.request) {
    return 'Network error - please check your connection';
  }
  return error.message || 'An unexpected error occurred';
}

// Convenience methods for HR SSOT operations
export const hrApi = {
  get: <T>(url: string, params?: Record<string, unknown>) =>
    hrApiClient.get<T>(url, { params }).then((res) => res.data),

  post: <T>(url: string, data?: unknown) =>
    hrApiClient.post<T>(url, data).then((res) => res.data),

  patch: <T>(url: string, data?: unknown) =>
    hrApiClient.patch<T>(url, data).then((res) => res.data),

  delete: <T>(url: string) =>
    hrApiClient.delete<T>(url).then((res) => res.data),
};

export { hrApiClient };
export default hrApi;
