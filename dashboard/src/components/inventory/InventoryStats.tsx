import { Card, CardContent } from '@/components/ui/card';
import { Laptop, MapPin, CheckCircle, AlertCircle } from 'lucide-react';
import type { InventorySummary } from '@/types';

interface InventoryStatsProps {
  summary: InventorySummary;
  loading?: boolean;
}

export function InventoryStats({ summary, loading }: InventoryStatsProps) {
  if (loading) {
    return (
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        {[1, 2, 3, 4].map((i) => (
          <Card key={i}>
            <CardContent className="p-6">
              <div className="h-16 animate-pulse bg-gray-100 rounded" />
            </CardContent>
          </Card>
        ))}
      </div>
    );
  }

  const stats = [
    {
      label: 'Total Devices',
      value: summary.totalDevices,
      icon: Laptop,
      color: 'text-blue-600',
      bgColor: 'bg-blue-50',
    },
    {
      label: 'Assigned',
      value: summary.byLocation.assigned,
      icon: MapPin,
      color: 'text-green-600',
      bgColor: 'bg-green-50',
    },
    {
      label: 'Unassigned',
      value: summary.byLocation.unassigned,
      icon: AlertCircle,
      color: 'text-yellow-600',
      bgColor: 'bg-yellow-50',
    },
    {
      label: 'Active',
      value: summary.byStatus['active'] || summary.byStatus['deployed'] || 0,
      icon: CheckCircle,
      color: 'text-emerald-600',
      bgColor: 'bg-emerald-50',
    },
  ];

  return (
    <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
      {stats.map((stat) => (
        <Card key={stat.label}>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-gray-500">{stat.label}</p>
                <p className="text-3xl font-bold text-gray-900 mt-1">
                  {stat.value.toLocaleString()}
                </p>
              </div>
              <div className={`p-3 rounded-full ${stat.bgColor}`}>
                <stat.icon className={`h-6 w-6 ${stat.color}`} />
              </div>
            </div>
          </CardContent>
        </Card>
      ))}
    </div>
  );
}
