import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Sheet, SheetHeader, SheetBody, SheetFooter } from '@/components/ui/sheet';
import { formatDate, formatStatus, cn } from '@/lib/utils';
import { useAuth } from '@/contexts/AuthContext';
import type { Incident, IncidentStatus, Severity } from '@/types';

const statusWorkflow: Record<IncidentStatus, IncidentStatus[]> = {
  new: ['acknowledged', 'in_progress', 'escalated'],
  acknowledged: ['in_progress', 'escalated'],
  in_progress: ['resolved', 'escalated'],
  escalated: ['in_progress', 'resolved'],
  resolved: ['closed', 'in_progress'],
  closed: [],
};

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

interface IncidentDetailProps {
  incident: Incident | null;
  open: boolean;
  onClose: () => void;
  onStatusUpdate: (incident: Incident, newStatus: IncidentStatus) => void;
}

export function IncidentDetail({
  incident,
  open,
  onClose,
  onStatusUpdate,
}: IncidentDetailProps) {
  const { hasRole } = useAuth();

  // Only admin and lead tech can create work orders from incidents
  const canCreateWorkOrder = hasRole('ssp_admin') || hasRole('ssp_lead_tech');

  if (!incident) return null;

  return (
    <Sheet open={open} onClose={onClose} side="right">
      <SheetHeader onClose={onClose}>Incident Details</SheetHeader>
      <SheetBody>
        <div className="space-y-6">
          {/* Header */}
          <div>
            <div className="flex items-center gap-2 mb-2">
              <Badge className={cn('capitalize', severityColors[incident.severity])}>
                {incident.severity}
              </Badge>
              <Badge className={cn('capitalize', statusColors[incident.status])}>
                {formatStatus(incident.status)}
              </Badge>
              {incident.slaBreached && (
                <Badge variant="destructive">SLA Breached</Badge>
              )}
            </div>
            <h2 className="text-xl font-semibold text-gray-900">
              {incident.title}
            </h2>
            <p className="text-sm text-gray-500 mt-1">ID: {incident.id}</p>
          </div>

          {/* Description */}
          {incident.description && (
            <div>
              <h3 className="text-sm font-medium text-gray-500 mb-1">
                Description
              </h3>
              <p className="text-gray-900">{incident.description}</p>
            </div>
          )}

          {/* School Info */}
          <div className="grid grid-cols-2 gap-4">
            <div>
              <h3 className="text-sm font-medium text-gray-500 mb-1">School</h3>
              <p className="text-gray-900">{incident.schoolName}</p>
              <p className="text-sm text-gray-500">
                {incident.countyName}, {incident.subCountyName}
              </p>
            </div>
            <div>
              <h3 className="text-sm font-medium text-gray-500 mb-1">Contact</h3>
              <p className="text-gray-900">{incident.contactName}</p>
              <p className="text-sm text-gray-500">{incident.contactPhone}</p>
            </div>
          </div>

          {/* Device Info */}
          <div className="grid grid-cols-2 gap-4">
            <div>
              <h3 className="text-sm font-medium text-gray-500 mb-1">Device</h3>
              <p className="text-gray-900">
                {incident.deviceMake} {incident.deviceModel}
              </p>
              <p className="text-sm text-gray-500">{incident.deviceCategory}</p>
            </div>
            <div>
              <h3 className="text-sm font-medium text-gray-500 mb-1">
                Serial / Asset Tag
              </h3>
              <p className="text-gray-900">{incident.deviceSerial}</p>
              <p className="text-sm text-gray-500">{incident.deviceAssetTag}</p>
            </div>
          </div>

          {/* Timeline */}
          <div>
            <h3 className="text-sm font-medium text-gray-500 mb-2">Timeline</h3>
            <div className="space-y-2 text-sm">
              <div className="flex justify-between">
                <span className="text-gray-500">Created</span>
                <span className="text-gray-900">
                  {formatDate(incident.createdAt)}
                </span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-500">Last Updated</span>
                <span className="text-gray-900">
                  {formatDate(incident.updatedAt)}
                </span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-500">SLA Due</span>
                <span
                  className={cn(
                    incident.slaBreached
                      ? 'text-red-600 font-medium'
                      : 'text-gray-900'
                  )}
                >
                  {formatDate(incident.slaDueAt)}
                </span>
              </div>
            </div>
          </div>

          {/* Status Workflow */}
          {statusWorkflow[incident.status].length > 0 && (
            <div>
              <h3 className="text-sm font-medium text-gray-500 mb-2">
                Update Status
              </h3>
              <div className="flex flex-wrap gap-2">
                {statusWorkflow[incident.status].map((nextStatus) => (
                  <Button
                    key={nextStatus}
                    variant="outline"
                    size="sm"
                    onClick={() => onStatusUpdate(incident, nextStatus)}
                  >
                    {formatStatus(nextStatus)}
                  </Button>
                ))}
              </div>
            </div>
          )}
        </div>
      </SheetBody>
      <SheetFooter>
        <Button variant="outline" onClick={onClose}>
          Close
        </Button>
        {canCreateWorkOrder && (
          <Button>Create Work Order</Button>
        )}
      </SheetFooter>
    </Sheet>
  );
}
