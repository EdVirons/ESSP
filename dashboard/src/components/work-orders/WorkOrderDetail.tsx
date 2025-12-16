import { Wrench, Laptop, User, DollarSign, FileText, Clock } from 'lucide-react';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Sheet, SheetHeader, SheetBody, SheetFooter } from '@/components/ui/sheet';
import { Tabs, TabsList, TabsTrigger, TabsContent } from '@/components/ui/tabs';
import { formatDate, formatStatus, formatCurrency, cn } from '@/lib/utils';
import { WorkOrderBOMTab } from './WorkOrderBOMTab';
import { WorkOrderScheduleTab } from './WorkOrderScheduleTab';
import { WorkOrderDeliverablesTab } from './WorkOrderDeliverablesTab';
import type { WorkOrder, WorkOrderStatus, WorkOrderPart, WorkOrderSchedule, WorkOrderDeliverable } from '@/types';

const statusWorkflow: Record<WorkOrderStatus, WorkOrderStatus[]> = {
  draft: ['assigned'],
  assigned: ['in_repair'],
  in_repair: ['qa'],
  qa: ['completed', 'in_repair'],
  completed: ['approved', 'qa'],
  approved: [],
};

const statusColors: Record<WorkOrderStatus, string> = {
  draft: 'bg-gray-100 text-gray-700 border border-gray-200',
  assigned: 'bg-blue-100 text-blue-700 border border-blue-200',
  in_repair: 'bg-amber-100 text-amber-700 border border-amber-200',
  qa: 'bg-purple-100 text-purple-700 border border-purple-200',
  completed: 'bg-emerald-100 text-emerald-700 border border-emerald-200',
  approved: 'bg-cyan-100 text-cyan-700 border border-cyan-200',
};

interface WorkOrderDetailProps {
  workOrder: WorkOrder | null;
  open: boolean;
  onClose: () => void;
  detailTab: string;
  onDetailTabChange: (tab: string) => void;
  bomItems: WorkOrderPart[];
  schedules: WorkOrderSchedule[];
  deliverables: WorkOrderDeliverable[];
  onStatusUpdate: (workOrder: WorkOrder, newStatus: WorkOrderStatus) => void;
}

export function WorkOrderDetail({
  workOrder,
  open,
  onClose,
  detailTab,
  onDetailTabChange,
  bomItems,
  schedules,
  deliverables,
  onStatusUpdate,
}: WorkOrderDetailProps) {
  if (!workOrder) return null;

  return (
    <Sheet open={open} onClose={onClose} side="right">
      <SheetHeader onClose={onClose}>Work Order Details</SheetHeader>
      <SheetBody className="p-0">
        <div className="h-full flex flex-col">
          {/* Header Info */}
          <div className="rounded-xl bg-gradient-to-r from-blue-600 via-cyan-600 to-blue-700 p-6 m-4 text-white shadow-lg">
            <div className="flex items-center gap-2 mb-3">
              <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-white/20 backdrop-blur-sm">
                <Wrench className="h-5 w-5 text-white" />
              </div>
              <Badge className={cn('capitalize shadow-sm', statusColors[workOrder.status])}>
                {formatStatus(workOrder.status)}
              </Badge>
              <span className="text-sm text-blue-100 capitalize">{workOrder.taskType}</span>
            </div>
            <h2 className="text-lg font-bold">
              {workOrder.schoolName}
            </h2>
            <p className="text-sm text-blue-100">
              {workOrder.deviceMake} {workOrder.deviceModel}
            </p>
          </div>

          {/* Tabs */}
          <Tabs
            value={detailTab}
            onValueChange={onDetailTabChange}
            className="flex-1 flex flex-col"
          >
            <div className="border-b border-gray-200 px-6">
              <TabsList className="bg-transparent -mb-px">
                <TabsTrigger value="details">Details</TabsTrigger>
                <TabsTrigger value="bom">BOM ({bomItems.length})</TabsTrigger>
                <TabsTrigger value="schedule">
                  Schedule ({schedules.length})
                </TabsTrigger>
                <TabsTrigger value="deliverables">
                  Deliverables ({deliverables.length})
                </TabsTrigger>
              </TabsList>
            </div>

            <div className="flex-1 overflow-auto">
              <TabsContent value="details" className="p-6 m-0">
                <div className="space-y-6">
                  {/* Device Info */}
                  <div className="grid grid-cols-2 gap-4">
                    <div className="rounded-xl border border-gray-100 bg-white p-4 shadow-sm">
                      <div className="flex items-center gap-2 mb-2">
                        <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-gradient-to-br from-blue-500 to-cyan-600">
                          <Laptop className="h-4 w-4 text-white" />
                        </div>
                        <h3 className="text-sm font-semibold text-gray-700">Device</h3>
                      </div>
                      <p className="text-gray-900 font-medium">
                        {workOrder.deviceMake} {workOrder.deviceModel}
                      </p>
                      <p className="text-sm text-gray-500 font-mono">
                        {workOrder.deviceSerial}
                      </p>
                    </div>
                    <div className="rounded-xl border border-gray-100 bg-white p-4 shadow-sm">
                      <div className="flex items-center gap-2 mb-2">
                        <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-gradient-to-br from-teal-500 to-emerald-600">
                          <User className="h-4 w-4 text-white" />
                        </div>
                        <h3 className="text-sm font-semibold text-gray-700">Contact</h3>
                      </div>
                      <p className="text-gray-900 font-medium">{workOrder.contactName}</p>
                      <p className="text-sm text-gray-500">
                        {workOrder.contactPhone}
                      </p>
                    </div>
                  </div>

                  {/* Cost */}
                  <div className="rounded-xl border border-gray-100 bg-white p-4 shadow-sm">
                    <div className="flex items-center gap-2 mb-2">
                      <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-gradient-to-br from-emerald-500 to-green-600">
                        <DollarSign className="h-4 w-4 text-white" />
                      </div>
                      <h3 className="text-sm font-semibold text-gray-700">Estimated Cost</h3>
                    </div>
                    <p className="text-xl font-bold text-gray-900">
                      {workOrder.costEstimateCents
                        ? formatCurrency(workOrder.costEstimateCents)
                        : 'Not estimated'}
                    </p>
                  </div>

                  {/* Notes */}
                  {workOrder.notes && (
                    <div className="rounded-xl border border-gray-100 bg-white p-4 shadow-sm">
                      <div className="flex items-center gap-2 mb-2">
                        <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-gradient-to-br from-slate-400 to-gray-500">
                          <FileText className="h-4 w-4 text-white" />
                        </div>
                        <h3 className="text-sm font-semibold text-gray-700">Notes</h3>
                      </div>
                      <p className="text-gray-700">{workOrder.notes}</p>
                    </div>
                  )}

                  {/* Timeline */}
                  <div className="rounded-xl border border-gray-100 bg-white p-4 shadow-sm">
                    <div className="flex items-center gap-2 mb-3">
                      <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-gradient-to-br from-purple-500 to-violet-600">
                        <Clock className="h-4 w-4 text-white" />
                      </div>
                      <h3 className="text-sm font-semibold text-gray-700">Timeline</h3>
                    </div>
                    <div className="space-y-2 text-sm">
                      <div className="flex justify-between items-center py-1.5 border-b border-gray-50">
                        <span className="text-gray-500">Created</span>
                        <span className="text-gray-900 font-medium bg-gray-100 px-2 py-0.5 rounded">
                          {formatDate(workOrder.createdAt)}
                        </span>
                      </div>
                      <div className="flex justify-between items-center py-1.5">
                        <span className="text-gray-500">Last Updated</span>
                        <span className="text-gray-900 font-medium bg-gray-100 px-2 py-0.5 rounded">
                          {formatDate(workOrder.updatedAt)}
                        </span>
                      </div>
                    </div>
                  </div>

                  {/* Status Workflow */}
                  {statusWorkflow[workOrder.status].length > 0 && (
                    <div className="rounded-xl bg-gradient-to-r from-blue-50 to-cyan-50 p-4 border border-blue-100">
                      <h3 className="text-sm font-semibold text-gray-700 mb-3">
                        Update Status
                      </h3>
                      <div className="flex flex-wrap gap-2">
                        {statusWorkflow[workOrder.status].map((nextStatus) => (
                          <Button
                            key={nextStatus}
                            variant="outline"
                            size="sm"
                            className="border-blue-200 hover:bg-blue-100 hover:text-blue-700"
                            onClick={() => onStatusUpdate(workOrder, nextStatus)}
                          >
                            {formatStatus(nextStatus)}
                          </Button>
                        ))}
                      </div>
                    </div>
                  )}
                </div>
              </TabsContent>

              <TabsContent value="bom" className="p-6 m-0">
                <WorkOrderBOMTab bomItems={bomItems} />
              </TabsContent>

              <TabsContent value="schedule" className="p-6 m-0">
                <WorkOrderScheduleTab schedules={schedules} />
              </TabsContent>

              <TabsContent value="deliverables" className="p-6 m-0">
                <WorkOrderDeliverablesTab deliverables={deliverables} />
              </TabsContent>
            </div>
          </Tabs>
        </div>
      </SheetBody>
      <SheetFooter>
        <Button variant="outline" onClick={onClose}>
          Close
        </Button>
      </SheetFooter>
    </Sheet>
  );
}
