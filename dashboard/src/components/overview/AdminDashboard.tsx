import { Link } from 'react-router-dom';
import {
  Shield,
  Users,
  School,
  Wrench,
  AlertTriangle,
  Activity,
  Server,
  Database,
  Settings,
  CheckCircle2,
  XCircle,
  TrendingUp,
  FileText,
  RefreshCw,
  BarChart3,
} from 'lucide-react';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { useAuth } from '@/contexts/AuthContext';

export function AdminDashboard() {
  const { user } = useAuth();
  const displayName = user?.displayName || user?.username || 'Administrator';

  // Mock data - in production, fetch from APIs
  const systemStats = {
    totalUsers: 156,
    activeSchools: 1248,
    totalDevices: 15420,
    openIncidents: 23,
    activeWorkOrders: 87,
    pendingApprovals: 12,
  };

  const serviceHealth = [
    { name: 'IMS API', status: 'healthy', latency: '45ms', uptime: '99.9%' },
    { name: 'SSOT School', status: 'healthy', latency: '32ms', uptime: '99.8%' },
    { name: 'SSOT Devices', status: 'healthy', latency: '28ms', uptime: '99.9%' },
    { name: 'Sync Worker', status: 'healthy', latency: '120ms', uptime: '99.7%' },
    { name: 'PostgreSQL', status: 'healthy', latency: '12ms', uptime: '100%' },
    { name: 'Redis', status: 'healthy', latency: '3ms', uptime: '100%' },
  ];

  const recentActivity = [
    { user: 'John Kamau', action: 'Created work order', target: 'WO-2024-001234', time: '5 min ago' },
    { user: 'Mary Wanjiku', action: 'Resolved incident', target: 'INC-2024-000891', time: '12 min ago' },
    { user: 'System', action: 'SSOT sync completed', target: '1,248 schools', time: '30 min ago' },
    { user: 'Peter Ochieng', action: 'Updated device status', target: '15 devices', time: '1 hour ago' },
  ];

  const quickStats = [
    { label: 'Today\'s Incidents', value: 8, change: -15, color: 'text-amber-600', bgColor: 'bg-amber-100' },
    { label: 'Today\'s Work Orders', value: 24, change: +12, color: 'text-blue-600', bgColor: 'bg-blue-100' },
    { label: 'Resolved Today', value: 31, change: +8, color: 'text-green-600', bgColor: 'bg-green-100' },
    { label: 'SLA Breaches', value: 2, change: -50, color: 'text-red-600', bgColor: 'bg-red-100' },
  ];

  return (
    <div className="space-y-6">
      {/* Welcome Header - Slate/Gray Theme (Admin - System Overview) */}
      <div className="relative overflow-hidden rounded-2xl bg-gradient-to-r from-slate-700 via-slate-800 to-slate-900 p-6 text-white shadow-lg">
        <div className="absolute top-0 right-0 w-64 h-64 bg-white/5 rounded-full -translate-y-32 translate-x-32" />
        <div className="absolute bottom-0 left-0 w-48 h-48 bg-white/5 rounded-full translate-y-24 -translate-x-24" />

        <div className="relative">
          <div className="flex items-center gap-3 mb-2">
            <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-white/20 backdrop-blur">
              <Shield className="h-6 w-6" />
            </div>
            <div>
              <h1 className="text-2xl md:text-3xl font-bold">Welcome, {displayName}!</h1>
              <p className="text-slate-300">System Administrator Dashboard</p>
            </div>
          </div>

          {/* Summary Stats */}
          <div className="mt-6 grid grid-cols-2 md:grid-cols-6 gap-4">
            <div className="bg-white/10 backdrop-blur rounded-xl p-3">
              <p className="text-2xl font-bold">{systemStats.totalUsers}</p>
              <p className="text-slate-300 text-sm">Total Users</p>
            </div>
            <div className="bg-white/10 backdrop-blur rounded-xl p-3">
              <p className="text-2xl font-bold">{systemStats.activeSchools.toLocaleString()}</p>
              <p className="text-slate-300 text-sm">Active Schools</p>
            </div>
            <div className="bg-white/10 backdrop-blur rounded-xl p-3">
              <p className="text-2xl font-bold">{(systemStats.totalDevices / 1000).toFixed(1)}K</p>
              <p className="text-slate-300 text-sm">Total Devices</p>
            </div>
            <div className="bg-white/10 backdrop-blur rounded-xl p-3">
              <p className="text-2xl font-bold">{systemStats.openIncidents}</p>
              <p className="text-slate-300 text-sm">Open Incidents</p>
            </div>
            <div className="bg-white/10 backdrop-blur rounded-xl p-3">
              <p className="text-2xl font-bold">{systemStats.activeWorkOrders}</p>
              <p className="text-slate-300 text-sm">Work Orders</p>
            </div>
            <div className="bg-white/10 backdrop-blur rounded-xl p-3">
              <p className="text-2xl font-bold">{systemStats.pendingApprovals}</p>
              <p className="text-slate-300 text-sm">Pending Approvals</p>
            </div>
          </div>
        </div>
      </div>

      {/* Quick Stats */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        {quickStats.map((stat) => (
          <Card key={stat.label}>
            <CardContent className="p-4">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-gray-500">{stat.label}</p>
                  <p className={`text-2xl font-bold ${stat.color}`}>{stat.value}</p>
                </div>
                <div className={`flex items-center gap-1 text-sm ${stat.change >= 0 ? 'text-green-600' : 'text-red-600'}`}>
                  <TrendingUp className={`h-4 w-4 ${stat.change < 0 ? 'rotate-180' : ''}`} />
                  {Math.abs(stat.change)}%
                </div>
              </div>
            </CardContent>
          </Card>
        ))}
      </div>

      {/* Quick Actions */}
      <div className="grid grid-cols-2 md:grid-cols-6 gap-4">
        <Link to="/incidents">
          <Card className="hover:shadow-md transition-shadow cursor-pointer border-amber-200 bg-amber-50">
            <CardContent className="p-4 flex flex-col items-center gap-2 text-center">
              <AlertTriangle className="h-6 w-6 text-amber-600" />
              <p className="font-medium text-amber-900 text-sm">Incidents</p>
            </CardContent>
          </Card>
        </Link>
        <Link to="/work-orders">
          <Card className="hover:shadow-md transition-shadow cursor-pointer border-blue-200 bg-blue-50">
            <CardContent className="p-4 flex flex-col items-center gap-2 text-center">
              <Wrench className="h-6 w-6 text-blue-600" />
              <p className="font-medium text-blue-900 text-sm">Work Orders</p>
            </CardContent>
          </Card>
        </Link>
        <Link to="/schools">
          <Card className="hover:shadow-md transition-shadow cursor-pointer border-green-200 bg-green-50">
            <CardContent className="p-4 flex flex-col items-center gap-2 text-center">
              <School className="h-6 w-6 text-green-600" />
              <p className="font-medium text-green-900 text-sm">Schools</p>
            </CardContent>
          </Card>
        </Link>
        <Link to="/reports">
          <Card className="hover:shadow-md transition-shadow cursor-pointer border-teal-200 bg-teal-50">
            <CardContent className="p-4 flex flex-col items-center gap-2 text-center">
              <BarChart3 className="h-6 w-6 text-teal-600" />
              <p className="font-medium text-teal-900 text-sm">Reports</p>
            </CardContent>
          </Card>
        </Link>
        <Link to="/audit-logs">
          <Card className="hover:shadow-md transition-shadow cursor-pointer border-slate-200 bg-slate-50">
            <CardContent className="p-4 flex flex-col items-center gap-2 text-center">
              <FileText className="h-6 w-6 text-slate-600" />
              <p className="font-medium text-slate-900 text-sm">Audit Logs</p>
            </CardContent>
          </Card>
        </Link>
        <Link to="/settings">
          <Card className="hover:shadow-md transition-shadow cursor-pointer border-gray-200 bg-gray-50">
            <CardContent className="p-4 flex flex-col items-center gap-2 text-center">
              <Settings className="h-6 w-6 text-gray-600" />
              <p className="font-medium text-gray-900 text-sm">Settings</p>
            </CardContent>
          </Card>
        </Link>
      </div>

      {/* Main Content Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Service Health */}
        <Card>
          <CardHeader className="pb-3">
            <div className="flex items-center justify-between">
              <CardTitle className="text-lg flex items-center gap-2">
                <Server className="h-5 w-5 text-slate-600" />
                Service Health
              </CardTitle>
              <Link to="/ssot-sync">
                <Button variant="ghost" size="sm">
                  <RefreshCw className="h-4 w-4 mr-1" />
                  Sync Status
                </Button>
              </Link>
            </div>
          </CardHeader>
          <CardContent>
            <div className="space-y-2">
              {serviceHealth.map((service) => (
                <div
                  key={service.name}
                  className="flex items-center justify-between p-3 rounded-lg bg-gray-50"
                >
                  <div className="flex items-center gap-3">
                    {service.status === 'healthy' ? (
                      <CheckCircle2 className="h-5 w-5 text-green-500" />
                    ) : (
                      <XCircle className="h-5 w-5 text-red-500" />
                    )}
                    <div>
                      <p className="font-medium text-gray-900">{service.name}</p>
                      <p className="text-xs text-gray-500">{service.latency} latency</p>
                    </div>
                  </div>
                  <div className="text-right">
                    <Badge variant="outline" className={`${
                      service.status === 'healthy' ? 'bg-green-50 text-green-700 border-green-200' : 'bg-red-50 text-red-700 border-red-200'
                    }`}>
                      {service.uptime} uptime
                    </Badge>
                  </div>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>

        {/* Recent System Activity */}
        <Card>
          <CardHeader className="pb-3">
            <div className="flex items-center justify-between">
              <CardTitle className="text-lg flex items-center gap-2">
                <Activity className="h-5 w-5 text-slate-600" />
                Recent Activity
              </CardTitle>
              <Link to="/audit-logs">
                <Button variant="ghost" size="sm">View All</Button>
              </Link>
            </div>
          </CardHeader>
          <CardContent>
            <div className="space-y-3">
              {recentActivity.map((activity, idx) => (
                <div
                  key={idx}
                  className="flex items-start gap-3 p-3 rounded-lg bg-gray-50"
                >
                  <div className="h-8 w-8 rounded-full bg-slate-200 flex items-center justify-center flex-shrink-0">
                    <Users className="h-4 w-4 text-slate-600" />
                  </div>
                  <div className="flex-1 min-w-0">
                    <p className="text-sm text-gray-900">
                      <span className="font-medium">{activity.user}</span>
                      {' '}{activity.action}
                    </p>
                    <p className="text-xs text-gray-500 truncate">{activity.target}</p>
                  </div>
                  <span className="text-xs text-gray-400 flex-shrink-0">{activity.time}</span>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      </div>

      {/* System Metrics */}
      <Card>
        <CardHeader className="pb-3">
          <CardTitle className="text-lg flex items-center gap-2">
            <Database className="h-5 w-5 text-slate-600" />
            System Metrics Overview
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-2 md:grid-cols-5 gap-4">
            {[
              { label: 'API Requests (24h)', value: '125.4K', trend: '+5%' },
              { label: 'Avg Response Time', value: '45ms', trend: '-12%' },
              { label: 'Error Rate', value: '0.02%', trend: '-8%' },
              { label: 'Active Sessions', value: '234', trend: '+15%' },
              { label: 'DB Connections', value: '45/100', trend: '0%' },
            ].map((metric) => (
              <div key={metric.label} className="text-center p-4 rounded-xl bg-gray-50">
                <p className="text-2xl font-bold text-gray-900">{metric.value}</p>
                <p className="text-sm text-gray-600">{metric.label}</p>
                <p className={`text-xs mt-1 ${
                  metric.trend.startsWith('+') ? 'text-green-600' :
                  metric.trend.startsWith('-') ? 'text-red-600' : 'text-gray-400'
                }`}>
                  {metric.trend}
                </p>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
