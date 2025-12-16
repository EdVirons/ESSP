import { Card, CardContent } from '@/components/ui/card';
import { Store, CheckCircle2, Users, MapPin, AlertTriangle, Package } from 'lucide-react';
import type { ServiceShop } from '@/types';

interface ServiceShopsStatsProps {
  shops: ServiceShop[];
  totalStaff?: number;
  lowStockCount?: number;
}

export function ServiceShopsStats({ shops, totalStaff = 0, lowStockCount = 0 }: ServiceShopsStatsProps) {
  const totalCount = shops.length;
  const activeCount = shops.filter((s) => s.active).length;
  const inactiveCount = totalCount - activeCount;
  const countiesCovered = new Set(shops.map((s) => s.countyCode)).size;

  return (
    <div className="grid gap-4 md:grid-cols-3 lg:grid-cols-6">
      <Card>
        <CardContent className="p-4">
          <div className="flex items-center gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-full bg-orange-50">
              <Store className="h-5 w-5 text-orange-600" />
            </div>
            <div>
              <div className="text-2xl font-bold">{totalCount}</div>
              <div className="text-sm text-gray-500">Total Shops</div>
            </div>
          </div>
        </CardContent>
      </Card>
      <Card>
        <CardContent className="p-4">
          <div className="flex items-center gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-full bg-green-50">
              <CheckCircle2 className="h-5 w-5 text-green-600" />
            </div>
            <div>
              <div className="text-2xl font-bold">{activeCount}</div>
              <div className="text-sm text-gray-500">Active</div>
            </div>
          </div>
        </CardContent>
      </Card>
      <Card>
        <CardContent className="p-4">
          <div className="flex items-center gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-full bg-gray-100">
              <Store className="h-5 w-5 text-gray-500" />
            </div>
            <div>
              <div className="text-2xl font-bold">{inactiveCount}</div>
              <div className="text-sm text-gray-500">Inactive</div>
            </div>
          </div>
        </CardContent>
      </Card>
      <Card>
        <CardContent className="p-4">
          <div className="flex items-center gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-full bg-blue-50">
              <Users className="h-5 w-5 text-blue-600" />
            </div>
            <div>
              <div className="text-2xl font-bold">{totalStaff}</div>
              <div className="text-sm text-gray-500">Total Staff</div>
            </div>
          </div>
        </CardContent>
      </Card>
      <Card>
        <CardContent className="p-4">
          <div className="flex items-center gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-full bg-purple-50">
              <MapPin className="h-5 w-5 text-purple-600" />
            </div>
            <div>
              <div className="text-2xl font-bold">{countiesCovered}</div>
              <div className="text-sm text-gray-500">Counties</div>
            </div>
          </div>
        </CardContent>
      </Card>
      <Card>
        <CardContent className="p-4">
          <div className="flex items-center gap-3">
            <div className={`flex h-10 w-10 items-center justify-center rounded-full ${lowStockCount > 0 ? 'bg-red-50' : 'bg-green-50'}`}>
              {lowStockCount > 0 ? (
                <AlertTriangle className="h-5 w-5 text-red-600" />
              ) : (
                <Package className="h-5 w-5 text-green-600" />
              )}
            </div>
            <div>
              <div className="text-2xl font-bold">{lowStockCount}</div>
              <div className="text-sm text-gray-500">Low Stock</div>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
