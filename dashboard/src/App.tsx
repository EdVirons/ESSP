import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { ReactQueryDevtools } from '@tanstack/react-query-devtools';

import { Layout } from '@/components/layout/Layout';
import { Toaster } from '@/components/ui/toaster';
import { NotificationProvider } from '@/contexts/NotificationContext';
import { AuthProvider } from '@/contexts/AuthContext';
import { ProtectedRoute } from '@/components/auth/ProtectedRoute';
import { Login } from '@/pages/Login';
import { Overview } from '@/pages/Overview';
import { Incidents } from '@/pages/Incidents';
import { WorkOrders } from '@/pages/WorkOrders';
import { Projects } from '@/pages/Projects';
import { ServiceShops } from '@/pages/ServiceShops';
import { AuditLogs } from '@/pages/AuditLogs';
import { Settings } from '@/pages/Settings';
import { Schools } from '@/pages/Schools';
import { DevicesPage } from '@/pages/DevicesPage';
import { PartsCatalog } from '@/pages/PartsCatalog';
import { SSOTSync } from '@/pages/SSOTSync';
import { Profile } from '@/pages/Profile';
import { Messages } from '@/pages/Messages';
import LiveChat from '@/pages/LiveChat';
import { KnowledgeBase } from '@/pages/KnowledgeBase';
import { SalesDashboard } from '@/pages/SalesDashboard';
import { DemoPipeline } from '@/pages/DemoPipeline';
import { Presentations } from '@/pages/Presentations';
import { Reports } from '@/pages/Reports';
import {
  WorkOrdersReport,
  IncidentsReport,
  InventoryReport,
  SchoolsReport,
  ExecutiveDashboard,
} from '@/pages/reports';
import { NotFound } from '@/pages/NotFound';

// Create a client
const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 30_000, // 30 seconds
      retry: 1,
      refetchOnWindowFocus: false,
    },
  },
});

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <BrowserRouter basename="/">
        <AuthProvider>
          <NotificationProvider>
            <Routes>
              {/* Redirect root to overview */}
              <Route path="/" element={<Navigate to="/overview" replace />} />

              {/* Login route - public */}
              <Route path="/login" element={<Login />} />

              {/* Main layout routes - protected */}
              <Route
                element={
                  <ProtectedRoute>
                    <Layout />
                  </ProtectedRoute>
                }
              >
                <Route path="/overview" element={<Overview />} />
                <Route path="/incidents" element={<Incidents />} />
                <Route path="/incidents/:id" element={<Incidents />} />
                <Route path="/work-orders" element={<WorkOrders />} />
                <Route path="/work-orders/:id" element={<WorkOrders />} />
                <Route path="/sales" element={<SalesDashboard />} />
                <Route path="/demo-pipeline" element={<DemoPipeline />} />
                <Route path="/presentations" element={<Presentations />} />
                <Route path="/projects" element={<Projects />} />
                <Route path="/projects/:id" element={<Projects />} />
                <Route path="/service-shops" element={<ServiceShops />} />
                <Route path="/service-shops/:id" element={<ServiceShops />} />
                <Route path="/schools" element={<Schools />} />
                <Route path="/devices" element={<DevicesPage />} />
                <Route path="/devices/:id" element={<DevicesPage />} />
                <Route path="/parts-catalog" element={<PartsCatalog />} />
                <Route path="/ssot-sync" element={<SSOTSync />} />
                <Route path="/messages" element={<Messages />} />
                <Route path="/messages/:id" element={<Messages />} />
                <Route path="/live-chat" element={<LiveChat />} />
                <Route path="/knowledge-base" element={<KnowledgeBase />} />
                <Route path="/reports" element={<Reports />} />
                <Route path="/reports/work-orders" element={<WorkOrdersReport />} />
                <Route path="/reports/incidents" element={<IncidentsReport />} />
                <Route path="/reports/inventory" element={<InventoryReport />} />
                <Route path="/reports/schools" element={<SchoolsReport />} />
                <Route path="/reports/executive" element={<ExecutiveDashboard />} />
                <Route path="/audit-logs" element={<AuditLogs />} />
                <Route path="/settings" element={<Settings />} />
                <Route path="/profile" element={<Profile />} />
                <Route path="*" element={<NotFound />} />
              </Route>
            </Routes>
            <Toaster />
          </NotificationProvider>
        </AuthProvider>
      </BrowserRouter>
      <ReactQueryDevtools initialIsOpen={false} />
    </QueryClientProvider>
  );
}

export default App;
