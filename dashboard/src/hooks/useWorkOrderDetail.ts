import * as React from 'react';
import type { WorkOrder, WorkOrderStatus } from '@/types';

interface StatusUpdate {
  workOrder: WorkOrder;
  newStatus: WorkOrderStatus;
}

export function useWorkOrderDetail() {
  const [selectedWorkOrder, setSelectedWorkOrder] = React.useState<WorkOrder | null>(null);
  const [showDetail, setShowDetail] = React.useState(false);
  const [detailTab, setDetailTab] = React.useState('details');
  const [statusUpdate, setStatusUpdate] = React.useState<StatusUpdate | null>(null);

  const openDetail = React.useCallback((workOrder: WorkOrder) => {
    setSelectedWorkOrder(workOrder);
    setShowDetail(true);
    setDetailTab('details');
  }, []);

  const closeDetail = React.useCallback(() => {
    setShowDetail(false);
  }, []);

  const handleStatusUpdateRequest = React.useCallback(
    (workOrder: WorkOrder, newStatus: WorkOrderStatus) => {
      setStatusUpdate({ workOrder, newStatus });
    },
    []
  );

  const clearStatusUpdate = React.useCallback(() => {
    setStatusUpdate(null);
  }, []);

  const updateSelectedWorkOrderStatus = React.useCallback(
    (newStatus: WorkOrderStatus) => {
      if (selectedWorkOrder) {
        setSelectedWorkOrder({ ...selectedWorkOrder, status: newStatus });
      }
    },
    [selectedWorkOrder]
  );

  return {
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
  };
}
