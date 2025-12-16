import { useMemo } from 'react';
import { Sun, Moon, CloudSun, Sparkles } from 'lucide-react';
import { useAuth } from '@/contexts/AuthContext';

function getGreeting(): { text: string; icon: React.ElementType } {
  const hour = new Date().getHours();
  if (hour < 12) {
    return { text: 'Good morning', icon: Sun };
  } else if (hour < 17) {
    return { text: 'Good afternoon', icon: CloudSun };
  } else {
    return { text: 'Good evening', icon: Moon };
  }
}

function formatDate(): string {
  return new Date().toLocaleDateString('en-US', {
    weekday: 'long',
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  });
}

// Role display names
const roleDisplayNames: Record<string, string> = {
  ssp_admin: 'System Administrator',
  ssp_school_contact: 'School Contact',
  ssp_support_agent: 'Support Agent',
  ssp_lead_tech: 'Lead Technician',
  ssp_field_tech: 'Field Technician',
  ssp_warehouse_manager: 'Warehouse Manager',
  ssp_demo_team: 'Demo Team',
  ssp_sales: 'Sales Representative',
};

interface QuickStat {
  value: string;
  label: string;
  permission?: string;
  roles?: string[];
}

// Stats configuration with permissions
const allStats: QuickStat[] = [
  { value: '12', label: 'Open Incidents', permission: 'incident:read' },
  { value: '8', label: 'My Work Orders', permission: 'workorder:read' },
  { value: '156', label: 'Active Schools', permission: 'school:read' },
  { value: '2.4K', label: 'Total Devices', permission: 'device:read' },
];

export function WelcomeHeader() {
  const { user, hasPermission, hasRole } = useAuth();
  const greeting = getGreeting();
  const Icon = greeting.icon;

  // Get display name
  const displayName = user?.displayName || user?.username || 'User';

  // Get role display name
  const userRole = user?.roles?.[0];
  const roleDisplay = userRole ? roleDisplayNames[userRole] || userRole.replace('ssp_', '').replace('_', ' ') : '';

  // Filter stats based on permissions
  const visibleStats = useMemo(() => {
    // Admin sees all stats
    if (hasRole('ssp_admin')) return allStats;

    return allStats.filter((stat) => {
      if (!stat.permission && !stat.roles) return true;
      if (stat.permission && hasPermission(stat.permission)) return true;
      if (stat.roles?.some((r) => hasRole(r))) return true;
      return false;
    });
  }, [hasPermission, hasRole]);

  return (
    <div className="relative overflow-hidden rounded-2xl bg-gradient-to-r from-cyan-600 via-teal-600 to-cyan-700 p-6 text-white shadow-lg">
      {/* Decorative background elements */}
      <div className="absolute top-0 right-0 w-64 h-64 bg-white/5 rounded-full -translate-y-32 translate-x-32" />
      <div className="absolute bottom-0 left-0 w-48 h-48 bg-white/5 rounded-full translate-y-24 -translate-x-24" />
      <div className="absolute top-1/2 right-1/4 w-24 h-24 bg-cyan-400/20 rounded-full blur-xl" />

      <div className="relative flex items-center justify-between">
        <div>
          <div className="flex items-center gap-3 mb-2">
            <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-white/20 backdrop-blur">
              <Icon className="h-5 w-5" />
            </div>
            <div>
              <h1 className="text-2xl md:text-3xl font-bold">
                {greeting.text}, {displayName}!
              </h1>
              {roleDisplay && (
                <p className="text-cyan-200 text-sm">{roleDisplay}</p>
              )}
            </div>
          </div>
          <p className="text-cyan-100 ml-13">{formatDate()}</p>
        </div>

        <div className="hidden md:flex items-center gap-4">
          <div className="text-right">
            <p className="text-sm text-cyan-200">Welcome to</p>
            <div className="flex items-center gap-2">
              <Sparkles className="h-5 w-5 text-cyan-300" />
              <p className="text-xl font-bold">ESSP Dashboard</p>
            </div>
          </div>
        </div>
      </div>

      {/* Quick stats row - only show if user has visible stats */}
      {visibleStats.length > 0 && (
        <div className={`relative mt-6 grid grid-cols-2 md:grid-cols-${Math.min(visibleStats.length, 4)} gap-4`}>
          {visibleStats.map((stat) => (
            <div key={stat.label} className="bg-white/10 backdrop-blur rounded-xl p-3">
              <p className="text-2xl font-bold">{stat.value}</p>
              <p className="text-cyan-200 text-sm">{stat.label}</p>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
