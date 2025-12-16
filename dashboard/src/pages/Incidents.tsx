import * as React from 'react';
import { useSearchParams } from 'react-router-dom';
import { type ColumnDef } from '@tanstack/react-table';
import { Plus, AlertTriangle, ChevronRight, AlertCircle } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { DataTable, SortableHeader, createSelectColumn } from '@/components/ui/data-table';
import { ConfirmDialog } from '@/components/ui/modal';
import {
  IncidentsStats,
  IncidentsFilters,
  IncidentDetail,
  CreateIncidentModal,
} from '@/components/incidents';
import { useIncidentFilters } from '@/hooks';
import { useIncidents, useCreateIncident, useUpdateIncidentStatus } from '@/api/incidents';
import { formatRelativeTime, formatStatus, cn } from '@/lib/utils';
import type { Incident, IncidentStatus, Severity } from '@/types';

const severityColors: Record<Severity, string> = {
  low: 'bg-green-100 text-green-800',
  medium: 'bg-yellow-100 text-yellow-800',
  high: 'bg-orange-100 text-orange-800',
  critical: 'bg-red-100 text-red-800',
};

const statusColors: Record<IncidentStatus, string> = {
  new: 'bg-blue-100 text-blue-800',
  acknowledged: 'bg-yellow-100 text-yellow-800',
  in_progress: 'bg-purple-100 text-purple-800',
  escalated: 'bg-red-100 text-red-800',
  resolved: 'bg-green-100 text-green-800',
  closed: 'bg-gray-100 text-gray-800',
};

export function Incidents() {
  const [searchParams, setSearchParams] = useSearchParams();

  // Filters
  const { filters, searchQuery, setSearchQuery, handleFilterChange, resetFilters } =
    useIncidentFilters();

  // Selected incident for detail view
  const [selectedIncident, setSelectedIncident] = React.useState<Incident | null>(null);
  const [showDetail, setShowDetail] = React.useState(false);

  // Create incident modal
  const [showCreateModal, setShowCreateModal] = React.useState(false);
  const [createForm, setCreateForm] = React.useState({
    deviceId: '',
    title: '',
    description: '',
    category: 'hardware',
    severity: 'medium' as Severity,
  });

  // Status update confirmation
  const [statusUpdate, setStatusUpdate] = React.useState<{
    incident: Incident;
    newStatus: IncidentStatus;
  } | null>(null);

  // API hooks
  const { data, isLoading, error } = useIncidents({
    ...filters,
    q: searchQuery || undefined,
  });
  const createIncident = useCreateIncident();
  const updateStatus = useUpdateIncidentStatus();

  // Handle action=create from URL (e.g., nav link "New")
  React.useEffect(() => {
    if (searchParams.get('action') === 'create') {
      setShowCreateModal(true);
      // Clear the action param from URL
      searchParams.delete('action');
      setSearchParams(searchParams, { replace: true });
    }
  }, [searchParams, setSearchParams]);

  // Table columns
  const columns: ColumnDef<Incident>[] = React.useMemo(
    () => [
      createSelectColumn<Incident>(),
      {
        accessorKey: 'title',
        header: ({ column }) => <SortableHeader column={column}>Title</SortableHeader>,
        cell: ({ row }) => (
          <div className="flex items-center gap-3">
            <div className="flex h-8 w-8 items-center justify-center rounded-full bg-amber-50">
              <AlertTriangle className="h-4 w-4 text-amber-600" />
            </div>
            <div className="min-w-0">
              <div className="font-medium text-gray-900 truncate max-w-[300px]">
                {row.original.title}
              </div>
              <div className="text-sm text-gray-500">{row.original.schoolName}</div>
            </div>
          </div>
        ),
      },
      {
        accessorKey: 'severity',
        header: 'Severity',
        cell: ({ row }) => (
          <Badge className={cn('capitalize', severityColors[row.original.severity])}>
            {row.original.severity}
          </Badge>
        ),
      },
      {
        accessorKey: 'status',
        header: 'Status',
        cell: ({ row }) => (
          <div className="flex items-center gap-2">
            <Badge className={cn('capitalize', statusColors[row.original.status])}>
              {formatStatus(row.original.status)}
            </Badge>
            {row.original.slaBreached && (
              <Badge variant="destructive" className="text-xs">
                SLA Breached
              </Badge>
            )}
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
              setSelectedIncident(row.original);
              setShowDetail(true);
            }}
          >
            <ChevronRight className="h-4 w-4" />
          </Button>
        ),
      },
    ],
    []
  );

  // Handle create incident
  const handleCreateIncident = async () => {
    try {
      await createIncident.mutateAsync({
        deviceId: createForm.deviceId,
        title: createForm.title,
        description: createForm.description,
        category: createForm.category,
        severity: createForm.severity,
      });
      setShowCreateModal(false);
      setCreateForm({
        deviceId: '',
        title: '',
        description: '',
        category: 'hardware',
        severity: 'medium',
      });
    } catch (err) {
      console.error('Failed to create incident:', err);
    }
  };

  // Handle status update
  const handleStatusUpdate = async () => {
    if (!statusUpdate) return;
    try {
      await updateStatus.mutateAsync({
        id: statusUpdate.incident.id,
        status: statusUpdate.newStatus,
      });
      setStatusUpdate(null);
      if (selectedIncident?.id === statusUpdate.incident.id) {
        setSelectedIncident({
          ...selectedIncident,
          status: statusUpdate.newStatus,
        });
      }
    } catch (err) {
      console.error('Failed to update status:', err);
    }
  };

  const incidents = data?.items || [];

  return (
    <div className="space-y-6">
      {/* Page Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Incidents</h1>
          <p className="text-sm text-gray-500">Manage device incidents and issues</p>
        </div>
        <Button onClick={() => setShowCreateModal(true)}>
          <Plus className="h-4 w-4" />
          Create Incident
        </Button>
      </div>

      {/* Stats Cards */}
      <IncidentsStats incidents={incidents} />

      {/* Filters */}
      <IncidentsFilters
        filters={filters}
        searchQuery={searchQuery}
        onSearchChange={setSearchQuery}
        onFilterChange={handleFilterChange}
        onClearFilters={resetFilters}
      />

      {/* Incidents Table */}
      <Card>
        <CardContent className="p-0">
          {error ? (
            <div className="p-8 text-center">
              <AlertCircle className="h-12 w-12 text-red-400 mx-auto mb-4" />
              <p className="text-gray-500">Failed to load incidents. Please try again.</p>
            </div>
          ) : (
            <DataTable
              columns={columns}
              data={incidents}
              isLoading={isLoading}
              searchKey="title"
              searchPlaceholder="Search by title..."
              showRowSelection
              showColumnVisibility
              onRowClick={(row) => {
                setSelectedIncident(row);
                setShowDetail(true);
              }}
              emptyMessage="No incidents found"
            />
          )}
        </CardContent>
      </Card>

      {/* Incident Detail Sheet */}
      <IncidentDetail
        incident={selectedIncident}
        open={showDetail}
        onClose={() => setShowDetail(false)}
        onStatusUpdate={(incident, newStatus) =>
          setStatusUpdate({ incident, newStatus })
        }
      />

      {/* Create Incident Modal */}
      <CreateIncidentModal
        open={showCreateModal}
        onClose={() => setShowCreateModal(false)}
        formData={createForm}
        onFormChange={setCreateForm}
        onSubmit={handleCreateIncident}
        isLoading={createIncident.isPending}
      />

      {/* Status Update Confirmation */}
      <ConfirmDialog
        open={!!statusUpdate}
        onClose={() => setStatusUpdate(null)}
        onConfirm={handleStatusUpdate}
        title="Update Incident Status"
        description={
          statusUpdate
            ? `Are you sure you want to change the status from "${formatStatus(statusUpdate.incident.status)}" to "${formatStatus(statusUpdate.newStatus)}"?`
            : ''
        }
        confirmText="Update Status"
        isLoading={updateStatus.isPending}
      />
    </div>
  );
}
