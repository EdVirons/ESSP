import { ServiceHealthCard } from '@/components/overview/ServiceHealthCard';
import { MetricsSummary } from '@/components/overview/MetricsSummary';
import { ActivityFeed } from '@/components/overview/ActivityFeed';
import {
  WelcomeHeader,
  QuickActions,
  MyTasks,
  NotificationsPanel,
  SchoolContactDashboard,
  FieldTechDashboard,
  LeadTechDashboard,
  WarehouseManagerDashboard,
  SupportAgentDashboard,
  SalesMarketingDashboard,
  DemoTeamDashboard,
  AdminDashboard,
} from '@/components/overview';
import { PermissionGate } from '@/components/auth';
import { useAuth } from '@/contexts/AuthContext';

export function Overview() {
  const { hasRole } = useAuth();

  // Admin gets the system admin dashboard
  if (hasRole('ssp_admin')) {
    return <AdminDashboard />;
  }

  // School contacts get their own focused dashboard
  if (hasRole('ssp_school_contact')) {
    return <SchoolContactDashboard />;
  }

  // Field technicians get their work-order focused dashboard
  if (hasRole('ssp_field_tech')) {
    return <FieldTechDashboard />;
  }

  // Lead technicians get their team management dashboard
  if (hasRole('ssp_lead_tech')) {
    return <LeadTechDashboard />;
  }

  // Warehouse managers get their inventory-focused dashboard
  if (hasRole('ssp_warehouse_manager')) {
    return <WarehouseManagerDashboard />;
  }

  // Support agents get their incident-focused dashboard
  if (hasRole('ssp_support_agent')) {
    return <SupportAgentDashboard />;
  }

  // Sales/Marketing team gets their sales dashboard
  if (hasRole('ssp_sales_marketing')) {
    return <SalesMarketingDashboard />;
  }

  // Demo team gets their project-focused dashboard
  if (hasRole('ssp_demo_team')) {
    return <DemoTeamDashboard />;
  }

  // Default fallback for other roles (contractors, suppliers)
  return (
    <div className="space-y-6">
      {/* Welcome Header - shows for all, but content is role-filtered */}
      <WelcomeHeader />

      {/* Quick Actions - actions are role-filtered inside the component */}
      <QuickActions />

      {/* My Tasks and Notifications - Two column layout (visible to all) */}
      <div className="grid gap-6 lg:grid-cols-2">
        <MyTasks />
        <NotificationsPanel />
      </div>

      {/* System Metrics Summary - Admin and managers only */}
      <PermissionGate roles={['ssp_admin', 'ssp_lead_tech', 'ssp_support_agent']}>
        <div>
          <h2 className="text-lg font-semibold text-gray-900 mb-4">
            System Metrics
          </h2>
          <MetricsSummary />
        </div>
      </PermissionGate>

      {/* Service Health and Activity - Admin only */}
      <PermissionGate roles={['ssp_admin']}>
        <div className="grid gap-6 lg:grid-cols-2">
          <ServiceHealthCard />
          <ActivityFeed />
        </div>
      </PermissionGate>
    </div>
  );
}
