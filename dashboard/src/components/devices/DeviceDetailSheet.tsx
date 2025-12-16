import * as React from 'react';
import {
  Laptop,
  Tag,
  School,
  User,
  Calendar,
  Clock,
  Shield,
  FileText,
  History,
  Settings,
  Wrench,
  ArrowRightLeft,
} from 'lucide-react';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Sheet, SheetHeader, SheetBody, SheetFooter } from '@/components/ui/sheet';
import { Tabs, TabsList, TabsTrigger, TabsContent } from '@/components/ui/tabs';
import { formatDate, formatRelativeTime, formatStatus, cn } from '@/lib/utils';
import type { Device, DeviceHistoryEntry, DeviceLifecycleStatus } from '@/types/device';
import {
  LIFECYCLE_STATUS_COLORS,
  ENROLLMENT_STATUS_COLORS,
  LIFECYCLE_TRANSITIONS,
} from '@/types/device';

interface DeviceDetailSheetProps {
  device: Device | null;
  open: boolean;
  onClose: () => void;
  activeTab: string;
  onTabChange: (tab: string) => void;
  history?: DeviceHistoryEntry[];
  historyLoading?: boolean;
  onStatusChange: (status: DeviceLifecycleStatus) => void;
  onEdit: () => void;
  onDelete: () => void;
  onCreateWorkOrder: () => void;
}

interface InfoRowProps {
  icon: React.ReactNode;
  label: string;
  value: React.ReactNode;
}

function InfoRow({ icon, label, value }: InfoRowProps) {
  return (
    <div className="flex items-start gap-3 py-2">
      <div className="text-gray-400 mt-0.5">{icon}</div>
      <div className="flex-1 min-w-0">
        <div className="text-sm text-gray-500">{label}</div>
        <div className="font-medium text-gray-900">{value || '-'}</div>
      </div>
    </div>
  );
}

function SpecRow({ label, value }: { label: string; value?: string }) {
  if (!value) return null;
  return (
    <div className="flex justify-between py-1.5 text-sm">
      <span className="text-gray-500">{label}</span>
      <span className="text-gray-900 font-medium">{value}</span>
    </div>
  );
}

export function DeviceDetailSheet({
  device,
  open,
  onClose,
  activeTab,
  onTabChange,
  history = [],
  historyLoading,
  onStatusChange,
  onEdit,
  onDelete,
  onCreateWorkOrder,
}: DeviceDetailSheetProps) {
  if (!device) return null;

  const availableTransitions = LIFECYCLE_TRANSITIONS[device.lifecycle] || [];
  const model = device.model;
  const specs = model?.specs || {};

  return (
    <Sheet open={open} onClose={onClose} side="right" className="max-w-lg">
      <SheetHeader onClose={onClose}>Device Details</SheetHeader>
      <SheetBody className="p-0">
        <div className="h-full flex flex-col">
          {/* Device Header */}
          <div className="p-6 border-b border-gray-200 bg-gray-50">
            <div className="flex items-start gap-4">
              <div className="flex h-14 w-14 items-center justify-center rounded-xl bg-indigo-100">
                <Laptop className="h-7 w-7 text-indigo-600" />
              </div>
              <div className="flex-1 min-w-0">
                <div className="flex items-center gap-2 mb-1">
                  <Badge className={cn('capitalize', LIFECYCLE_STATUS_COLORS[device.lifecycle])}>
                    {formatStatus(device.lifecycle)}
                  </Badge>
                  <Badge className={cn('capitalize', ENROLLMENT_STATUS_COLORS[device.enrolled])}>
                    {formatStatus(device.enrolled)}
                  </Badge>
                </div>
                <h2 className="text-lg font-semibold text-gray-900 font-mono">
                  {device.serial}
                </h2>
                <div className="flex items-center gap-1 text-sm text-gray-500">
                  <Tag className="h-3.5 w-3.5" />
                  {device.assetTag || 'No asset tag'}
                </div>
              </div>
            </div>
          </div>

          {/* Tabs */}
          <Tabs
            value={activeTab}
            onValueChange={onTabChange}
            className="flex-1 flex flex-col"
          >
            <div className="border-b border-gray-200 px-6">
              <TabsList className="bg-transparent -mb-px">
                <TabsTrigger value="info">
                  <FileText className="h-4 w-4 mr-1.5" />
                  Info
                </TabsTrigger>
                <TabsTrigger value="history">
                  <History className="h-4 w-4 mr-1.5" />
                  History
                </TabsTrigger>
                <TabsTrigger value="actions">
                  <Settings className="h-4 w-4 mr-1.5" />
                  Actions
                </TabsTrigger>
              </TabsList>
            </div>

            <div className="flex-1 overflow-auto">
              {/* Info Tab */}
              <TabsContent value="info" className="p-6 m-0">
                <div className="space-y-6">
                  {/* Device Model */}
                  <section>
                    <h3 className="text-sm font-semibold text-gray-900 uppercase tracking-wide mb-3">
                      Device Model
                    </h3>
                    <div className="bg-gray-50 rounded-lg p-4">
                      <div className="font-medium text-gray-900 text-lg">
                        {model ? `${model.make} ${model.model}` : 'Unknown Model'}
                      </div>
                      <div className="text-sm text-gray-500 capitalize">
                        {model?.category || 'Unknown category'}
                      </div>

                      {Object.keys(specs).length > 0 && (
                        <div className="mt-4 pt-4 border-t border-gray-200 space-y-1">
                          <SpecRow
                            label="Processor"
                            value={specs.processor}
                          />
                          <SpecRow label="RAM" value={specs.ram} />
                          <SpecRow label="Storage" value={specs.storage} />
                          <SpecRow label="Display" value={specs.display} />
                          <SpecRow label="OS" value={specs.os} />
                        </div>
                      )}
                    </div>
                  </section>

                  {/* Assignment */}
                  <section>
                    <h3 className="text-sm font-semibold text-gray-900 uppercase tracking-wide mb-3">
                      Assignment
                    </h3>
                    <div className="space-y-1">
                      <InfoRow
                        icon={<School className="h-4 w-4" />}
                        label="School"
                        value={device.schoolName || device.schoolId}
                      />
                      <InfoRow
                        icon={<User className="h-4 w-4" />}
                        label="Assigned To"
                        value={device.assignedToName || device.assignedTo}
                      />
                      {device.assignedAt && (
                        <InfoRow
                          icon={<Calendar className="h-4 w-4" />}
                          label="Assigned Date"
                          value={formatDate(device.assignedAt)}
                        />
                      )}
                    </div>
                  </section>

                  {/* Dates */}
                  <section>
                    <h3 className="text-sm font-semibold text-gray-900 uppercase tracking-wide mb-3">
                      Dates
                    </h3>
                    <div className="space-y-1">
                      <InfoRow
                        icon={<Calendar className="h-4 w-4" />}
                        label="Purchase Date"
                        value={device.purchaseDate ? formatDate(device.purchaseDate) : undefined}
                      />
                      <InfoRow
                        icon={<Shield className="h-4 w-4" />}
                        label="Warranty Expiry"
                        value={device.warrantyExpiry ? formatDate(device.warrantyExpiry) : undefined}
                      />
                      <InfoRow
                        icon={<Clock className="h-4 w-4" />}
                        label="Last Seen"
                        value={device.lastSeen ? formatRelativeTime(device.lastSeen) : 'Never'}
                      />
                      <InfoRow
                        icon={<Calendar className="h-4 w-4" />}
                        label="Created"
                        value={formatDate(device.createdAt)}
                      />
                    </div>
                  </section>

                  {/* Notes */}
                  {device.notes && (
                    <section>
                      <h3 className="text-sm font-semibold text-gray-900 uppercase tracking-wide mb-3">
                        Notes
                      </h3>
                      <p className="text-gray-700 text-sm whitespace-pre-wrap bg-gray-50 rounded-lg p-4">
                        {device.notes}
                      </p>
                    </section>
                  )}
                </div>
              </TabsContent>

              {/* History Tab */}
              <TabsContent value="history" className="p-6 m-0">
                {historyLoading ? (
                  <div className="flex items-center justify-center py-12">
                    <div className="h-6 w-6 animate-spin rounded-full border-2 border-gray-300 border-t-blue-600" />
                  </div>
                ) : history.length === 0 ? (
                  <div className="text-center py-12">
                    <History className="h-12 w-12 text-gray-300 mx-auto mb-3" />
                    <p className="text-gray-500">No history available</p>
                  </div>
                ) : (
                  <div className="relative">
                    {/* Timeline line */}
                    <div className="absolute left-4 top-0 bottom-0 w-0.5 bg-gray-200" />

                    <div className="space-y-6">
                      {history.map((entry) => (
                        <div key={entry.id} className="relative pl-10">
                          {/* Timeline dot */}
                          <div className="absolute left-2.5 top-1.5 w-3 h-3 rounded-full bg-blue-500 border-2 border-white" />

                          <div className="bg-gray-50 rounded-lg p-4">
                            <div className="flex items-center justify-between mb-1">
                              <span className="font-medium text-gray-900 capitalize">
                                {entry.action.replace(/_/g, ' ')}
                              </span>
                              <span className="text-xs text-gray-500">
                                {formatRelativeTime(entry.timestamp)}
                              </span>
                            </div>
                            {entry.field && (
                              <div className="text-sm text-gray-600">
                                <span className="text-gray-500">{entry.field}:</span>{' '}
                                {entry.oldValue && (
                                  <>
                                    <span className="line-through text-gray-400">
                                      {entry.oldValue}
                                    </span>
                                    {' -> '}
                                  </>
                                )}
                                <span className="font-medium">{entry.newValue}</span>
                              </div>
                            )}
                            <div className="text-xs text-gray-500 mt-1">
                              by {entry.actorName || entry.actor}
                            </div>
                          </div>
                        </div>
                      ))}
                    </div>
                  </div>
                )}
              </TabsContent>

              {/* Actions Tab */}
              <TabsContent value="actions" className="p-6 m-0">
                <div className="space-y-6">
                  {/* Status Changes */}
                  {availableTransitions.length > 0 && (
                    <section>
                      <h3 className="text-sm font-semibold text-gray-900 uppercase tracking-wide mb-3">
                        Change Status
                      </h3>
                      <div className="flex flex-wrap gap-2">
                        {availableTransitions.map((status) => (
                          <Button
                            key={status}
                            variant="outline"
                            size="sm"
                            onClick={() => onStatusChange(status)}
                            className="gap-2"
                          >
                            <ArrowRightLeft className="h-4 w-4" />
                            Set to {formatStatus(status)}
                          </Button>
                        ))}
                      </div>
                    </section>
                  )}

                  {/* Work Order */}
                  <section>
                    <h3 className="text-sm font-semibold text-gray-900 uppercase tracking-wide mb-3">
                      Repair
                    </h3>
                    <Button
                      variant="outline"
                      onClick={onCreateWorkOrder}
                      className="gap-2"
                    >
                      <Wrench className="h-4 w-4" />
                      Create Work Order
                    </Button>
                  </section>

                  {/* Danger Zone */}
                  <section className="pt-4 border-t border-gray-200">
                    <h3 className="text-sm font-semibold text-red-600 uppercase tracking-wide mb-3">
                      Danger Zone
                    </h3>
                    <Button
                      variant="outline"
                      onClick={onDelete}
                      className="gap-2 text-red-600 hover:text-red-700 hover:bg-red-50 border-red-200"
                    >
                      Delete Device
                    </Button>
                  </section>
                </div>
              </TabsContent>
            </div>
          </Tabs>
        </div>
      </SheetBody>
      <SheetFooter>
        <Button variant="outline" onClick={onClose}>
          Close
        </Button>
        <Button onClick={onEdit}>Edit Device</Button>
      </SheetFooter>
    </Sheet>
  );
}
