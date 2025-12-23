import * as React from 'react';
import { type ColumnDef } from '@tanstack/react-table';
import { Plus, Wrench, ChevronRight } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { DataTable, SortableHeader, createSelectColumn } from '@/components/ui/data-table';
import { ConfirmDialog } from '@/components/ui/modal';
import { PermissionGate } from '@/components/auth';
import {
  WorkOrdersStats,
  WorkOrdersFilters,
  WorkOrderDetail,
  CreateWorkOrderModal,
} from '@/components/work-orders';
import { useWorkOrderFilters, useWorkOrderDetail } from '@/hooks';
import {
  useWorkOrders,
  useCreateWorkOrder,
  useUpdateWorkOrderStatus,
  useWorkOrderBOM,
  useWorkOrderSchedules,
  useWorkOrderDeliverables,
} from '@/api/work-orders';
import { formatRelativeTime, formatStatus, formatCurrency, cn } from '@/lib/utils';
import type { WorkOrder, WorkOrderStatus } from '@/types';

const statusColors: Record<WorkOrderStatus, string> = {
  draft: 'bg-gray-100 text-gray-700 border border-gray-200',
  assigned: 'bg-blue-100 text-blue-700 border border-blue-200',
  in_repair: 'bg-amber-100 text-amber-700 border border-amber-200',
  qa: 'bg-purple-100 text-purple-700 border border-purple-200',
  completed: 'bg-emerald-100 text-emerald-700 border border-emerald-200',
  approved: 'bg-cyan-100 text-cyan-700 border border-cyan-200',
};

export function WorkOrders() {
  // Filters
  const { filters, searchQuery, setSearchQuery, handleFilterChange } = useWorkOrderFilters();

  // Detail view state
  const {
    selectedWorkOrder,
    showDetail,
    detailTab,
    statusUpdate,
    openDetail,
    closeDetail,
    setDetailTab,
    handleStatusUpdateRequest,
    clearStatusUpdate,
    updateSelectedWorkOrderStatus,
  } = useWorkOrderDetail();

  // Create work order modal
  const [showCreateModal, setShowCreateModal] = React.useState(false);
  const [createForm, setCreateForm] = React.useState({
    deviceId: '',
    taskType: 'repair',
    incidentId: '',
    notes: '',
  });

  // API hooks
  const { data, isLoading } = useWorkOrders(filters);
  const createWorkOrder = useCreateWorkOrder();
  const updateStatus = useUpdateWorkOrderStatus();

  // Detail view data
  const { data: bomData } = useWorkOrderBOM(selectedWorkOrder?.id || '');
  const { data: schedulesData } = useWorkOrderSchedules(selectedWorkOrder?.id || '');
  const { data: deliverablesData } = useWorkOrderDeliverables(selectedWorkOrder?.id || '');

  // Table columns
  const columns: ColumnDef<WorkOrder>[] = React.useMemo(
    () => [
      createSelectColumn<WorkOrder>(),
      {
        accessorKey: 'id',
        header: ({ column }) => <SortableHeader column={column}>Work Order</SortableHeader>,
        cell: ({ row }) => (
          <div className="flex items-center gap-3">
            <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-gradient-to-br from-blue-500 to-cyan-600 shadow-md shadow-blue-500/20">
              <Wrench className="h-4 w-4 text-white" />
            </div>
            <div className="min-w-0">
              <div className="font-medium text-gray-900 truncate max-w-[200px]">
                {row.original.id.slice(0, 8)}...
              </div>
              <div className="text-sm text-gray-500">{row.original.taskType}</div>
            </div>
          </div>
        ),
      },
      {
        accessorKey: 'status',
        header: 'Status',
        cell: ({ row }) => (
          <Badge className={cn('capitalize', statusColors[row.original.status])}>
            {formatStatus(row.original.status)}
          </Badge>
        ),
      },
      {
        accessorKey: 'schoolName',
        header: 'School',
        cell: ({ row }) => (
          <div className="text-sm">
            <div className="font-medium">{row.original.schoolName}</div>
            <div className="text-gray-500">{row.original.contactName}</div>
          </div>
        ),
      },
      {
        accessorKey: 'deviceModel',
        header: 'Device',
        cell: ({ row }) => (
          <div className="text-sm">
            <div className="font-medium">
              {row.original.deviceMake} {row.original.deviceModel}
            </div>
            <div className="text-gray-500">{row.original.deviceSerial}</div>
          </div>
        ),
      },
      {
        accessorKey: 'costEstimateCents',
        header: 'Est. Cost',
        cell: ({ row }) => (
          <div className="text-sm font-medium">
            {row.original.costEstimateCents
              ? formatCurrency(row.original.costEstimateCents)
              : '-'}
          </div>
        ),
      },
      {
        accessorKey: 'createdAt',
        header: ({ column }) => <SortableHeader column={column}>Created</SortableHeader>,
        cell: ({ row }) => (
          <div className="text-sm text-gray-500">
            {formatRelativeTime(row.original.createdAt)}
          </div>
        ),
      },
      {
        id: 'actions',
        cell: ({ row }) => (
          <Button
            variant="ghost"
            size="sm"
            onClick={(e) => {
              e.stopPropagation();
              openDetail(row.original);
            }}
          >
            <ChevronRight className="h-4 w-4" />
          </Button>
        ),
      },
    ],
    [openDetail]
  );

  // Handle create work order
  const handleCreateWorkOrder = async () => {
    try {
      await createWorkOrder.mutateAsync({
        deviceId: createForm.deviceId,
        taskType: createForm.taskType,
        incidentId: createForm.incidentId || undefined,
        notes: createForm.notes || undefined,
      });
      setShowCreateModal(false);
      setCreateForm({
        deviceId: '',
        taskType: 'repair',
        incidentId: '',
        notes: '',
      });
    } catch (err) {
      console.error('Failed to create work order:', err);
    }
  };

  // Handle status update
  const handleStatusUpdate = async () => {
    if (!statusUpdate) return;
    try {
      await updateStatus.mutateAsync({
        id: statusUpdate.workOrder.id,
        status: statusUpdate.newStatus,
      });
      if (selectedWorkOrder?.id === statusUpdate.workOrder.id) {
        updateSelectedWorkOrderStatus(statusUpdate.newStatus);
      }
      clearStatusUpdate();
    } catch (err) {
      console.error('Failed to update status:', err);
    }
  };

  const workOrders = data?.items || [];
  const bomItems = bomData?.items || [];
  const schedules = schedulesData?.items || [];
  const deliverables = deliverablesData?.items || [];

  return (
    <div className="space-y-6">
      {/* Page Header */}
      <div className="rounded-xl bg-gradient-to-r from-blue-600 via-cyan-600 to-blue-700 p-4 sm:p-6 text-white shadow-lg">
        <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
          <div className="flex items-center gap-3 sm:gap-4">
            <div className="flex h-12 w-12 sm:h-14 sm:w-14 items-center justify-center rounded-xl bg-white/20 backdrop-blur-sm">
              <Wrench className="h-6 w-6 sm:h-7 sm:w-7 text-white" />
            </div>
            <div>
              <h1 className="text-xl sm:text-2xl font-bold">Work Orders</h1>
              <p className="text-sm sm:text-base text-blue-100">
                Manage repair work orders and assignments
              </p>
            </div>
          </div>
          <PermissionGate permissions={['workorder:create']}>
            <Button onClick={() => setShowCreateModal(true)} className="w-full sm:w-auto bg-white text-blue-700 hover:bg-blue-50">
              <Plus className="h-4 w-4" />
              Create Work Order
            </Button>
          </PermissionGate>
        </div>
      </div>

      {/* Stats Cards */}
      <WorkOrdersStats workOrders={workOrders} />

      {/* Filters */}
      <WorkOrdersFilters
        filters={filters}
        searchQuery={searchQuery}
        onSearchChange={setSearchQuery}
        onFilterChange={handleFilterChange}
      />

      {/* Work Orders Table */}
      <Card>
        <CardContent className="p-0">
          <DataTable
            columns={columns}
            data={workOrders}
            isLoading={isLoading}
            searchKey="id"
            searchPlaceholder="Search by ID..."
            showRowSelection
            showColumnVisibility
            onRowClick={openDetail}
            emptyMessage="No work orders found"
          />
        </CardContent>
      </Card>

      {/* Work Order Detail Sheet */}
      <WorkOrderDetail
        workOrder={selectedWorkOrder}
        open={showDetail}
        onClose={closeDetail}
        detailTab={detailTab}
        onDetailTabChange={setDetailTab}
        bomItems={bomItems}
        schedules={schedules}
        deliverables={deliverables}
        onStatusUpdate={handleStatusUpdateRequest}
      />

      {/* Create Work Order Modal */}
      <CreateWorkOrderModal
        open={showCreateModal}
        onClose={() => setShowCreateModal(false)}
        formData={createForm}
        onFormChange={setCreateForm}
        onSubmit={handleCreateWorkOrder}
        isLoading={createWorkOrder.isPending}
      />

      {/* Status Update Confirmation */}
      <ConfirmDialog
        open={!!statusUpdate}
        onClose={clearStatusUpdate}
        onConfirm={handleStatusUpdate}
        title="Update Work Order Status"
        description={
          statusUpdate
            ? `Are you sure you want to change the status from "${formatStatus(statusUpdate.workOrder.status)}" to "${formatStatus(statusUpdate.newStatus)}"?`
            : ''
        }
        confirmText="Update Status"
        isLoading={updateStatus.isPending}
      />
    </div>
  );
}
