import { Card, CardContent } from '@/components/ui/card';
import { Users, UserCheck, UserX, Wrench, HardHat, Package } from 'lucide-react';
import type { ServiceStaffStats } from '@/types';

interface StaffStatsProps {
  stats?: ServiceStaffStats;
  isLoading?: boolean;
}

export function StaffStats({ stats, isLoading }: StaffStatsProps) {
  const total = stats?.total || 0;
  const active = stats?.active || 0;
  const inactive = stats?.inactive || 0;
  const leadTechs = stats?.byRole?.lead_technician || 0;
  const fieldTechs = stats?.byRole?.assistant_technician || 0;
  const storekeepers = stats?.byRole?.storekeeper || 0;

  if (isLoading) {
    return (
      <div className="grid gap-4 md:grid-cols-3 lg:grid-cols-6">
        {[...Array(6)].map((_, i) => (
          <Card key={i}>
            <CardContent className="p-4">
              <div className="animate-pulse">
                <div className="h-10 w-10 rounded-full bg-gray-200 mb-2" />
                <div className="h-6 w-12 bg-gray-200 rounded mb-1" />
                <div className="h-4 w-20 bg-gray-200 rounded" />
              </div>
            </CardContent>
          </Card>
        ))}
      </div>
    );
  }

  return (
    <div className="grid gap-4 md:grid-cols-3 lg:grid-cols-6">
      <Card>
        <CardContent className="p-4">
          <div className="flex items-center gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-full bg-blue-50">
              <Users className="h-5 w-5 text-blue-600" />
            </div>
            <div>
              <div className="text-2xl font-bold">{total}</div>
              <div className="text-sm text-gray-500">Total Staff</div>
            </div>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardContent className="p-4">
          <div className="flex items-center gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-full bg-green-50">
              <UserCheck className="h-5 w-5 text-green-600" />
            </div>
            <div>
              <div className="text-2xl font-bold">{active}</div>
              <div className="text-sm text-gray-500">Active</div>
            </div>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardContent className="p-4">
          <div className="flex items-center gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-full bg-gray-100">
              <UserX className="h-5 w-5 text-gray-600" />
            </div>
            <div>
              <div className="text-2xl font-bold">{inactive}</div>
              <div className="text-sm text-gray-500">Inactive</div>
            </div>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardContent className="p-4">
          <div className="flex items-center gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-full bg-purple-50">
              <Wrench className="h-5 w-5 text-purple-600" />
            </div>
            <div>
              <div className="text-2xl font-bold">{leadTechs}</div>
              <div className="text-sm text-gray-500">Lead Techs</div>
            </div>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardContent className="p-4">
          <div className="flex items-center gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-full bg-blue-50">
              <HardHat className="h-5 w-5 text-blue-600" />
            </div>
            <div>
              <div className="text-2xl font-bold">{fieldTechs}</div>
              <div className="text-sm text-gray-500">Field Techs</div>
            </div>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardContent className="p-4">
          <div className="flex items-center gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-full bg-green-50">
              <Package className="h-5 w-5 text-green-600" />
            </div>
            <div>
              <div className="text-2xl font-bold">{storekeepers}</div>
              <div className="text-sm text-gray-500">Storekeepers</div>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
