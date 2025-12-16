import { Activity, CheckCircle, AlertCircle, XCircle } from 'lucide-react';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { cn } from '@/lib/utils';
import type { ServiceHealth } from '@/types';

interface ServiceHealthCardProps {
  services?: ServiceHealth[];
  isLoading?: boolean;
}

const statusConfig = {
  healthy: {
    icon: CheckCircle,
    color: 'text-green-500',
    bgColor: 'bg-green-50',
    label: 'Healthy',
  },
  degraded: {
    icon: AlertCircle,
    color: 'text-yellow-500',
    bgColor: 'bg-yellow-50',
    label: 'Degraded',
  },
  unhealthy: {
    icon: XCircle,
    color: 'text-red-500',
    bgColor: 'bg-red-50',
    label: 'Unhealthy',
  },
};

// Mock data for demo
const mockServices: ServiceHealth[] = [
  { name: 'ims-api', status: 'healthy', latencyMs: 12, lastCheck: new Date().toISOString() },
  { name: 'ssot-school', status: 'healthy', latencyMs: 8, lastCheck: new Date().toISOString() },
  { name: 'ssot-devices', status: 'healthy', latencyMs: 15, lastCheck: new Date().toISOString() },
  { name: 'ssot-parts', status: 'healthy', latencyMs: 10, lastCheck: new Date().toISOString() },
  { name: 'sync-worker', status: 'healthy', latencyMs: 5, lastCheck: new Date().toISOString() },
];

export function ServiceHealthCard({ services = mockServices, isLoading }: ServiceHealthCardProps) {
  const allHealthy = services.every((s) => s.status === 'healthy');
  const overallStatus = allHealthy ? 'healthy' : services.some((s) => s.status === 'unhealthy') ? 'unhealthy' : 'degraded';

  if (isLoading) {
    return (
      <Card>
        <CardHeader className="pb-2">
          <CardTitle className="flex items-center gap-2 text-base">
            <Activity className="h-5 w-5" />
            Service Health
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-3">
            {[1, 2, 3, 4, 5].map((i) => (
              <div key={i} className="skeleton h-8 rounded" />
            ))}
          </div>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card>
      <CardHeader className="pb-2">
        <div className="flex items-center justify-between">
          <CardTitle className="flex items-center gap-2 text-base">
            <Activity className="h-5 w-5" />
            Service Health
          </CardTitle>
          <div
            className={cn(
              'flex items-center gap-1.5 rounded-full px-2 py-0.5 text-xs font-medium',
              statusConfig[overallStatus].bgColor,
              statusConfig[overallStatus].color
            )}
          >
            {(() => {
              const Icon = statusConfig[overallStatus].icon;
              return <Icon className="h-3 w-3" />;
            })()}
            {statusConfig[overallStatus].label}
          </div>
        </div>
      </CardHeader>
      <CardContent>
        <div className="space-y-2">
          {services.map((service) => {
            const config = statusConfig[service.status];
            const StatusIcon = config.icon;

            return (
              <div
                key={service.name}
                className="flex items-center justify-between rounded-lg border border-gray-100 p-2"
              >
                <div className="flex items-center gap-3">
                  <StatusIcon className={cn('h-4 w-4', config.color)} />
                  <span className="text-sm font-medium text-gray-900">
                    {service.name}
                  </span>
                </div>
                <div className="flex items-center gap-3 text-xs text-gray-500">
                  <span>{service.latencyMs}ms</span>
                </div>
              </div>
            );
          })}
        </div>
      </CardContent>
    </Card>
  );
}
