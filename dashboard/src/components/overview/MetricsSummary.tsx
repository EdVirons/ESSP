import { AlertTriangle, Wrench, Layers, Package } from 'lucide-react';
import { Card, CardContent } from '@/components/ui/card';
import { cn } from '@/lib/utils';
import type { DashboardMetrics } from '@/types';
import { Link } from 'react-router-dom';

interface MetricCardProps {
  title: string;
  value: number;
  subtitle?: string;
  icon: React.ElementType;
  href: string;
  trend?: 'up' | 'down' | 'neutral';
  trendValue?: string;
  iconColor: string;
  iconBgColor: string;
}

function MetricCard({
  title,
  value,
  subtitle,
  icon: Icon,
  href,
  iconColor,
  iconBgColor,
}: MetricCardProps) {
  return (
    <Link to={href}>
      <Card className="transition-shadow hover:shadow-md">
        <CardContent className="flex items-center gap-4 p-4">
          <div className={cn('rounded-lg p-3', iconBgColor)}>
            <Icon className={cn('h-6 w-6', iconColor)} />
          </div>
          <div className="flex-1">
            <p className="text-sm font-medium text-gray-500">{title}</p>
            <p className="text-2xl font-bold text-gray-900">{value.toLocaleString()}</p>
            {subtitle && (
              <p className="text-xs text-gray-500">{subtitle}</p>
            )}
          </div>
        </CardContent>
      </Card>
    </Link>
  );
}

interface MetricsSummaryProps {
  metrics?: DashboardMetrics;
  isLoading?: boolean;
}

// Mock data for demo
const mockMetrics: DashboardMetrics = {
  incidents: { total: 156, open: 23, slaBreached: 3 },
  workOrders: { total: 89, inProgress: 15, completedToday: 8 },
  programs: { active: 12, pending: 5 },
};

export function MetricsSummary({ metrics = mockMetrics, isLoading }: MetricsSummaryProps) {
  if (isLoading) {
    return (
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        {[1, 2, 3, 4].map((i) => (
          <Card key={i}>
            <CardContent className="p-4">
              <div className="skeleton h-20 rounded" />
            </CardContent>
          </Card>
        ))}
      </div>
    );
  }

  return (
    <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
      <MetricCard
        title="Open Incidents"
        value={metrics.incidents.open}
        subtitle={`${metrics.incidents.slaBreached} SLA breached`}
        icon={AlertTriangle}
        href="/incidents?status=new,acknowledged,in_progress"
        iconColor="text-amber-600"
        iconBgColor="bg-amber-50"
      />
      <MetricCard
        title="Active Work Orders"
        value={metrics.workOrders.inProgress}
        subtitle={`${metrics.workOrders.completedToday} completed today`}
        icon={Wrench}
        href="/work-orders?status=in_repair"
        iconColor="text-blue-600"
        iconBgColor="bg-blue-50"
      />
      <MetricCard
        title="Active Projects"
        value={metrics.programs.active}
        subtitle={`${metrics.programs.pending} pending`}
        icon={Layers}
        href="/projects?status=active"
        iconColor="text-purple-600"
        iconBgColor="bg-purple-50"
      />
      <MetricCard
        title="Total Incidents"
        value={metrics.incidents.total}
        subtitle="All time"
        icon={Package}
        href="/incidents"
        iconColor="text-gray-600"
        iconBgColor="bg-gray-100"
      />
    </div>
  );
}
