import { Card, CardContent } from '@/components/ui/card';
import { Wrench, Clock, CheckCircle2, Package } from 'lucide-react';
import type { WorkOrder } from '@/types';

interface WorkOrdersStatsProps {
  workOrders: WorkOrder[];
}

export function WorkOrdersStats({ workOrders }: WorkOrdersStatsProps) {
  const totalCount = workOrders.length;
  const inProgressCount = workOrders.filter(
    (wo) => ['assigned', 'in_repair', 'qa'].includes(wo.status)
  ).length;
  const completedCount = workOrders.filter(
    (wo) => wo.status === 'completed' || wo.status === 'approved'
  ).length;
  const draftCount = workOrders.filter((wo) => wo.status === 'draft').length;

  return (
    <div className="grid grid-cols-2 gap-3 md:gap-4 lg:grid-cols-4">
      <Card className="border-0 shadow-md overflow-hidden">
        <CardContent className="p-0">
          <div className="flex items-center gap-4 p-4">
            <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-gradient-to-br from-blue-500 to-cyan-600 shadow-lg shadow-blue-500/20">
              <Wrench className="h-6 w-6 text-white" />
            </div>
            <div>
              <div className="text-2xl font-bold text-gray-900">{totalCount}</div>
              <div className="text-sm text-gray-500">Total Work Orders</div>
            </div>
          </div>
          <div className="h-1 bg-gradient-to-r from-blue-400 to-cyan-500" />
        </CardContent>
      </Card>
      <Card className="border-0 shadow-md overflow-hidden">
        <CardContent className="p-0">
          <div className="flex items-center gap-4 p-4">
            <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-gradient-to-br from-amber-500 to-yellow-600 shadow-lg shadow-amber-500/20">
              <Clock className="h-6 w-6 text-white" />
            </div>
            <div>
              <div className="text-2xl font-bold text-gray-900">{inProgressCount}</div>
              <div className="text-sm text-gray-500">In Progress</div>
            </div>
          </div>
          <div className="h-1 bg-gradient-to-r from-amber-400 to-yellow-500" />
        </CardContent>
      </Card>
      <Card className="border-0 shadow-md overflow-hidden">
        <CardContent className="p-0">
          <div className="flex items-center gap-4 p-4">
            <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-gradient-to-br from-emerald-500 to-green-600 shadow-lg shadow-emerald-500/20">
              <CheckCircle2 className="h-6 w-6 text-white" />
            </div>
            <div>
              <div className="text-2xl font-bold text-gray-900">{completedCount}</div>
              <div className="text-sm text-gray-500">Completed</div>
            </div>
          </div>
          <div className="h-1 bg-gradient-to-r from-emerald-400 to-green-500" />
        </CardContent>
      </Card>
      <Card className="border-0 shadow-md overflow-hidden">
        <CardContent className="p-0">
          <div className="flex items-center gap-4 p-4">
            <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-gradient-to-br from-slate-400 to-gray-500 shadow-lg shadow-slate-500/20">
              <Package className="h-6 w-6 text-white" />
            </div>
            <div>
              <div className="text-2xl font-bold text-gray-900">{draftCount}</div>
              <div className="text-sm text-gray-500">Draft</div>
            </div>
          </div>
          <div className="h-1 bg-gradient-to-r from-slate-300 to-gray-400" />
        </CardContent>
      </Card>
    </div>
  );
}
