import { useNavigate, Link } from 'react-router-dom';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import {
  ArrowLeft,
  ArrowRight,
  Wrench,
  AlertTriangle,
  Package,
  School,
  ShieldCheck,
  TrendingUp,
  TrendingDown,
} from 'lucide-react';
import { useExecutiveDashboard } from '@/api/reports';

function KPICard({
  title,
  value,
  subtitle,
  icon: Icon,
  trend,
  trendValue,
  color = 'default',
}: {
  title: string;
  value: string | number;
  subtitle?: string;
  icon: React.ElementType;
  trend?: 'up' | 'down' | 'neutral';
  trendValue?: string;
  color?: 'default' | 'success' | 'warning' | 'danger' | 'info';
}) {
  const colorClasses = {
    default: 'text-gray-900',
    success: 'text-green-600',
    warning: 'text-amber-600',
    danger: 'text-red-600',
    info: 'text-blue-600',
  };

  const bgColorClasses = {
    default: 'bg-gray-100',
    success: 'bg-green-100',
    warning: 'bg-amber-100',
    danger: 'bg-red-100',
    info: 'bg-blue-100',
  };

  return (
    <Card>
      <CardContent className="p-6">
        <div className="flex items-start justify-between">
          <div className={`p-3 rounded-lg ${bgColorClasses[color]}`}>
            <Icon className={`h-6 w-6 ${colorClasses[color]}`} />
          </div>
          {trend && trendValue && (
            <div
              className={`flex items-center gap-1 text-sm ${
                trend === 'up' ? 'text-green-600' : trend === 'down' ? 'text-red-600' : 'text-gray-500'
              }`}
            >
              {trend === 'up' ? (
                <TrendingUp className="h-4 w-4" />
              ) : trend === 'down' ? (
                <TrendingDown className="h-4 w-4" />
              ) : null}
              {trendValue}
            </div>
          )}
        </div>
        <div className="mt-4">
          <p className="text-sm text-gray-500 font-medium">{title}</p>
          <p className={`text-3xl font-bold mt-1 ${colorClasses[color]}`}>
            {typeof value === 'number' ? value.toLocaleString() : value}
          </p>
          {subtitle && <p className="text-xs text-gray-400 mt-1">{subtitle}</p>}
        </div>
      </CardContent>
    </Card>
  );
}

function SectionCard({
  title,
  icon: Icon,
  color,
  href,
  children,
}: {
  title: string;
  icon: React.ElementType;
  color: string;
  href: string;
  children: React.ReactNode;
}) {
  return (
    <Card className="overflow-hidden">
      <CardHeader className={`${color} text-white py-4`}>
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-3">
            <Icon className="h-5 w-5" />
            <CardTitle className="text-lg text-white">{title}</CardTitle>
          </div>
          <Link to={href}>
            <Button variant="ghost" size="sm" className="text-white hover:bg-white/20">
              View Report
              <ArrowRight className="ml-2 h-4 w-4" />
            </Button>
          </Link>
        </div>
      </CardHeader>
      <CardContent className="p-6">{children}</CardContent>
    </Card>
  );
}

export function ExecutiveDashboard() {
  const navigate = useNavigate();
  const { data, isLoading } = useExecutiveDashboard();

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900" />
      </div>
    );
  }

  if (!data) {
    return (
      <div className="text-center py-12">
        <p className="text-gray-500">Failed to load dashboard data</p>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-4">
        <Button variant="ghost" size="icon" onClick={() => navigate('/reports')}>
          <ArrowLeft className="h-5 w-5" />
        </Button>
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Executive Dashboard</h1>
          <p className="text-sm text-gray-500 mt-1">
            High-level KPIs and metrics across all platform domains
          </p>
        </div>
      </div>

      {/* Top Level KPIs */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <KPICard
          title="Total Work Orders"
          value={data.workOrders.total}
          subtitle={`${data.workOrders.completed} completed`}
          icon={Wrench}
          color="info"
        />
        <KPICard
          title="Open Incidents"
          value={data.incidents.open}
          subtitle={`${data.incidents.critical} critical`}
          icon={AlertTriangle}
          color={data.incidents.critical > 0 ? 'warning' : 'default'}
        />
        <KPICard
          title="SLA Compliance"
          value={`${data.incidents.slaCompliance.toFixed(1)}%`}
          subtitle="Target: 95%"
          icon={ShieldCheck}
          color={data.incidents.slaCompliance >= 95 ? 'success' : 'warning'}
        />
        <KPICard
          title="Low Stock Items"
          value={data.inventory.lowStock}
          subtitle={`${data.inventory.outOfStock} out of stock`}
          icon={Package}
          color={data.inventory.lowStock > 0 ? 'warning' : 'success'}
        />
      </div>

      {/* Section Cards */}
      <div className="grid gap-6 lg:grid-cols-2">
        {/* Work Orders Section */}
        <SectionCard title="Work Orders" icon={Wrench} color="bg-blue-600" href="/reports/work-orders">
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-1">
              <p className="text-sm text-gray-500">Total</p>
              <p className="text-2xl font-bold">{data.workOrders.total.toLocaleString()}</p>
            </div>
            <div className="space-y-1">
              <p className="text-sm text-gray-500">Completed</p>
              <p className="text-2xl font-bold text-green-600">
                {data.workOrders.completed.toLocaleString()}
              </p>
            </div>
            <div className="space-y-1">
              <p className="text-sm text-gray-500">In Progress</p>
              <p className="text-2xl font-bold text-amber-600">
                {data.workOrders.inProgress.toLocaleString()}
              </p>
            </div>
            <div className="space-y-1">
              <p className="text-sm text-gray-500">Completion Rate</p>
              <p className="text-2xl font-bold">{data.workOrders.completionRate.toFixed(1)}%</p>
            </div>
          </div>
          <div className="mt-4 pt-4 border-t">
            <div className="flex items-center justify-between text-sm">
              <span className="text-gray-500">Avg. Completion Time</span>
              <span className="font-medium">{data.workOrders.avgCompletionDays.toFixed(1)} days</span>
            </div>
          </div>
        </SectionCard>

        {/* Incidents Section */}
        <SectionCard title="Incidents" icon={AlertTriangle} color="bg-amber-600" href="/reports/incidents">
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-1">
              <p className="text-sm text-gray-500">Total</p>
              <p className="text-2xl font-bold">{data.incidents.total.toLocaleString()}</p>
            </div>
            <div className="space-y-1">
              <p className="text-sm text-gray-500">Open</p>
              <p className="text-2xl font-bold text-amber-600">
                {data.incidents.open.toLocaleString()}
              </p>
            </div>
            <div className="space-y-1">
              <p className="text-sm text-gray-500">Resolved</p>
              <p className="text-2xl font-bold text-green-600">
                {data.incidents.resolved.toLocaleString()}
              </p>
            </div>
            <div className="space-y-1">
              <p className="text-sm text-gray-500">Critical</p>
              <p className={`text-2xl font-bold ${data.incidents.critical > 0 ? 'text-red-600' : 'text-gray-900'}`}>
                {data.incidents.critical.toLocaleString()}
              </p>
            </div>
          </div>
          <div className="mt-4 pt-4 border-t">
            <div className="flex items-center justify-between text-sm">
              <span className="text-gray-500">SLA Compliance</span>
              <Badge className={data.incidents.slaCompliance >= 95 ? 'bg-green-100 text-green-700' : 'bg-amber-100 text-amber-700'}>
                {data.incidents.slaCompliance.toFixed(1)}%
              </Badge>
            </div>
          </div>
        </SectionCard>

        {/* Inventory Section */}
        <SectionCard title="Inventory" icon={Package} color="bg-purple-600" href="/reports/inventory">
          <div className="grid grid-cols-3 gap-4">
            <div className="space-y-1">
              <p className="text-sm text-gray-500">Total Parts</p>
              <p className="text-2xl font-bold">{data.inventory.totalParts.toLocaleString()}</p>
            </div>
            <div className="space-y-1">
              <p className="text-sm text-gray-500">Low Stock</p>
              <p className={`text-2xl font-bold ${data.inventory.lowStock > 0 ? 'text-amber-600' : 'text-gray-900'}`}>
                {data.inventory.lowStock.toLocaleString()}
              </p>
            </div>
            <div className="space-y-1">
              <p className="text-sm text-gray-500">Out of Stock</p>
              <p className={`text-2xl font-bold ${data.inventory.outOfStock > 0 ? 'text-red-600' : 'text-gray-900'}`}>
                {data.inventory.outOfStock.toLocaleString()}
              </p>
            </div>
          </div>
          <div className="mt-4 pt-4 border-t">
            <div className="flex items-center justify-between text-sm">
              <span className="text-gray-500">Stock Health</span>
              <Badge
                className={
                  data.inventory.lowStock === 0 && data.inventory.outOfStock === 0
                    ? 'bg-green-100 text-green-700'
                    : data.inventory.outOfStock > 0
                      ? 'bg-red-100 text-red-700'
                      : 'bg-amber-100 text-amber-700'
                }
              >
                {data.inventory.lowStock === 0 && data.inventory.outOfStock === 0
                  ? 'Healthy'
                  : data.inventory.outOfStock > 0
                    ? 'Critical'
                    : 'Needs Attention'}
              </Badge>
            </div>
          </div>
        </SectionCard>

        {/* Schools Section */}
        <SectionCard title="Schools & Devices" icon={School} color="bg-green-600" href="/reports/schools">
          <div className="grid grid-cols-3 gap-4">
            <div className="space-y-1">
              <p className="text-sm text-gray-500">Schools</p>
              <p className="text-2xl font-bold">{data.schools.totalSchools.toLocaleString()}</p>
            </div>
            <div className="space-y-1">
              <p className="text-sm text-gray-500">Devices</p>
              <p className="text-2xl font-bold">{data.schools.totalDevices.toLocaleString()}</p>
            </div>
            <div className="space-y-1">
              <p className="text-sm text-gray-500">Active Projects</p>
              <p className="text-2xl font-bold text-blue-600">
                {data.schools.activeProjects.toLocaleString()}
              </p>
            </div>
          </div>
          <div className="mt-4 pt-4 border-t">
            <div className="flex items-center justify-between text-sm">
              <span className="text-gray-500">Avg. Devices per School</span>
              <span className="font-medium">
                {data.schools.totalSchools > 0
                  ? (data.schools.totalDevices / data.schools.totalSchools).toFixed(1)
                  : '0'}
              </span>
            </div>
          </div>
        </SectionCard>
      </div>
    </div>
  );
}
