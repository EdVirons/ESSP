import {
  Laptop,
  Package,
  CheckCircle2,
  Wrench,
  Archive,
  Layers,
} from 'lucide-react';
import { Card, CardContent } from '@/components/ui/card';
import { cn } from '@/lib/utils';
import type { DeviceStats } from '@/types/device';

interface DevicesStatsProps {
  stats?: DeviceStats;
  isLoading?: boolean;
}

interface StatCardProps {
  title: string;
  value: number;
  icon: React.ReactNode;
  colorClass: string;
  bgClass: string;
}

function StatCard({ title, value, icon, colorClass, bgClass }: StatCardProps) {
  return (
    <Card>
      <CardContent className="pt-6">
        <div className="flex items-center justify-between">
          <div>
            <p className={cn('text-sm', colorClass)}>{title}</p>
            <p className={cn('text-2xl font-bold', colorClass)}>
              {value.toLocaleString()}
            </p>
          </div>
          <div
            className={cn(
              'h-12 w-12 rounded-full flex items-center justify-center',
              bgClass
            )}
          >
            {icon}
          </div>
        </div>
      </CardContent>
    </Card>
  );
}

function StatCardSkeleton() {
  return (
    <Card>
      <CardContent className="pt-6">
        <div className="flex items-center justify-between">
          <div className="space-y-2">
            <div className="h-4 w-16 bg-gray-200 rounded animate-pulse" />
            <div className="h-8 w-12 bg-gray-200 rounded animate-pulse" />
          </div>
          <div className="h-12 w-12 rounded-full bg-gray-200 animate-pulse" />
        </div>
      </CardContent>
    </Card>
  );
}

export function DevicesStats({ stats, isLoading }: DevicesStatsProps) {
  if (isLoading || !stats) {
    return (
      <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-4">
        {Array.from({ length: 6 }).map((_, i) => (
          <StatCardSkeleton key={i} />
        ))}
      </div>
    );
  }

  const statCards: StatCardProps[] = [
    {
      title: 'Total Devices',
      value: stats.total,
      icon: <Laptop className="h-6 w-6 text-gray-600" />,
      colorClass: 'text-gray-600',
      bgClass: 'bg-gray-100',
    },
    {
      title: 'In Stock',
      value: stats.byLifecycle.in_stock || 0,
      icon: <Package className="h-6 w-6 text-green-600" />,
      colorClass: 'text-green-600',
      bgClass: 'bg-green-50',
    },
    {
      title: 'Deployed',
      value: stats.byLifecycle.deployed || 0,
      icon: <CheckCircle2 className="h-6 w-6 text-blue-600" />,
      colorClass: 'text-blue-600',
      bgClass: 'bg-blue-50',
    },
    {
      title: 'In Repair',
      value: stats.byLifecycle.repair || 0,
      icon: <Wrench className="h-6 w-6 text-yellow-600" />,
      colorClass: 'text-yellow-600',
      bgClass: 'bg-yellow-50',
    },
    {
      title: 'Retired',
      value: stats.byLifecycle.retired || 0,
      icon: <Archive className="h-6 w-6 text-red-600" />,
      colorClass: 'text-red-600',
      bgClass: 'bg-red-50',
    },
    {
      title: 'Device Models',
      value: stats.modelsCount,
      icon: <Layers className="h-6 w-6 text-purple-600" />,
      colorClass: 'text-purple-600',
      bgClass: 'bg-purple-50',
    },
  ];

  return (
    <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-4">
      {statCards.map((card) => (
        <StatCard key={card.title} {...card} />
      ))}
    </div>
  );
}
