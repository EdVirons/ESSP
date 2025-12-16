import { Laptop, Package, Wrench } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import type { SSOTDeviceStats as SSOTDeviceStatsType } from '@/types/device';

interface SSOTDeviceStatsProps {
  stats: SSOTDeviceStatsType;
}

export function SSOTDeviceStats({ stats }: SSOTDeviceStatsProps) {
  return (
    <div className="grid gap-4 md:grid-cols-4">
      <Card>
        <CardHeader className="flex flex-row items-center justify-between pb-2">
          <CardTitle className="text-sm font-medium">Total Devices</CardTitle>
          <Laptop className="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">{stats.total}</div>
          <p className="text-xs text-muted-foreground">
            {stats.uniqueModels} unique models
          </p>
        </CardContent>
      </Card>

      <Card>
        <CardHeader className="flex flex-row items-center justify-between pb-2">
          <CardTitle className="text-sm font-medium">Deployed</CardTitle>
          <Laptop className="h-4 w-4 text-blue-500" />
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold text-blue-600">
            {stats.byStatus.deployed || 0}
          </div>
          <p className="text-xs text-muted-foreground">Active in schools</p>
        </CardContent>
      </Card>

      <Card>
        <CardHeader className="flex flex-row items-center justify-between pb-2">
          <CardTitle className="text-sm font-medium">In Stock</CardTitle>
          <Package className="h-4 w-4 text-green-500" />
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold text-green-600">
            {stats.byStatus.in_stock || 0}
          </div>
          <p className="text-xs text-muted-foreground">Available for deployment</p>
        </CardContent>
      </Card>

      <Card>
        <CardHeader className="flex flex-row items-center justify-between pb-2">
          <CardTitle className="text-sm font-medium">In Repair</CardTitle>
          <Wrench className="h-4 w-4 text-yellow-500" />
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold text-yellow-600">
            {stats.byStatus.repair || 0}
          </div>
          <p className="text-xs text-muted-foreground">
            {stats.byStatus.retired || 0} retired
          </p>
        </CardContent>
      </Card>
    </div>
  );
}
