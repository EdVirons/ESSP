import React, { type ReactElement } from 'react';
import { render, type RenderOptions } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { BrowserRouter } from 'react-router-dom';

// Create a custom render function that includes providers
interface ProvidersProps {
  children: React.ReactNode;
}

function createTestQueryClient() {
  return new QueryClient({
    defaultOptions: {
      queries: {
        retry: false,
        gcTime: 0,
      },
      mutations: {
        retry: false,
      },
    },
  });
}

function AllTheProviders({ children }: ProvidersProps) {
  const queryClient = createTestQueryClient();

  return (
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        {children}
      </BrowserRouter>
    </QueryClientProvider>
  );
}

// Custom render function
const customRender = (
  ui: ReactElement,
  options?: Omit<RenderOptions, 'wrapper'>
) => render(ui, { wrapper: AllTheProviders, ...options });

// Re-export everything
export * from '@testing-library/react';
export { customRender as render };

// Helper to wait for async operations
export const waitForLoadingToFinish = () =>
  new Promise((resolve) => setTimeout(resolve, 0));

// Mock API response helper
export function mockApiResponse<T>(data: T, delay = 0): Promise<T> {
  return new Promise((resolve) => {
    setTimeout(() => resolve(data), delay);
  });
}

// Mock API error helper
export function mockApiError(message: string, status = 500): Promise<never> {
  return Promise.reject({
    response: {
      status,
      data: { message },
    },
  });
}

// Create mock incident
export function createMockIncident(overrides = {}) {
  return {
    id: 'inc-123',
    type: 'damage',
    status: 'open',
    priority: 'medium',
    description: 'Test incident',
    schoolId: 'school-1',
    deviceSerial: 'DEV-001',
    reportedBy: 'Test User',
    createdAt: new Date().toISOString(),
    ...overrides,
  };
}

// Create mock work order
export function createMockWorkOrder(overrides = {}) {
  return {
    id: 'wo-123',
    title: 'Test Work Order',
    status: 'assigned',
    priority: 'medium',
    schoolId: 'school-1',
    assignedTo: 'tech-1',
    createdAt: new Date().toISOString(),
    ...overrides,
  };
}

// Create mock user
export function createMockUser(overrides = {}) {
  return {
    username: 'testuser',
    email: 'test@example.com',
    displayName: 'Test User',
    roles: ['ssp_support_agent'],
    ...overrides,
  };
}
