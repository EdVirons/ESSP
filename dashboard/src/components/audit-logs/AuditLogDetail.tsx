import { User, Clock, Globe, Server, FileCode } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Sheet, SheetHeader, SheetBody, SheetFooter } from '@/components/ui/sheet';
import { formatDate, cn } from '@/lib/utils';
import type { AuditLog } from '@/types';
import { actionColors } from './columns';

interface AuditLogDetailProps {
  log: AuditLog | null;
  open: boolean;
  onClose: () => void;
}

export function AuditLogDetail({ log, open, onClose }: AuditLogDetailProps) {
  return (
    <Sheet open={open} onClose={onClose} side="right">
      <SheetHeader onClose={onClose}>Audit Log Details</SheetHeader>
      <SheetBody>
        {log && (
          <div className="space-y-6">
            {/* Header */}
            <div className="rounded-xl bg-gradient-to-r from-slate-100 to-gray-100 p-4 border border-slate-200">
              <div className="flex items-center gap-2 mb-3">
                <Badge className={cn('capitalize shadow-sm', actionColors[log.action])}>
                  {log.action}
                </Badge>
                <span className="text-sm text-gray-600 capitalize font-medium">
                  {log.entityType.replace(/_/g, ' ')}
                </span>
              </div>
              <h2 className="text-lg font-bold text-gray-900">
                {log.action.charAt(0).toUpperCase() + log.action.slice(1)} Operation
              </h2>
              <div className="flex items-center gap-2 mt-2 text-sm text-gray-500">
                <Clock className="h-4 w-4" />
                {formatDate(log.createdAt)}
              </div>
            </div>

            {/* User Info */}
            <div className="flex items-center gap-3 p-4 bg-gradient-to-r from-cyan-50 to-teal-50 border border-cyan-100 rounded-xl">
              <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-gradient-to-br from-cyan-500 to-teal-600 shadow-lg shadow-cyan-500/20">
                <User className="h-6 w-6 text-white" />
              </div>
              <div>
                <div className="font-semibold text-gray-900">{log.userEmail}</div>
                <div className="text-sm text-gray-500">User ID: <span className="font-mono text-xs">{log.userId}</span></div>
              </div>
            </div>

            {/* Entity Info */}
            <div className="rounded-xl border border-gray-100 bg-white p-4 shadow-sm">
              <h3 className="text-sm font-semibold text-gray-700 mb-3 flex items-center gap-2">
                <FileCode className="h-4 w-4 text-slate-500" />
                Entity Information
              </h3>
              <div className="space-y-2 text-sm">
                <div className="flex justify-between items-center py-1.5 border-b border-gray-50">
                  <span className="text-gray-500">Type</span>
                  <span className="text-gray-900 font-medium capitalize bg-gray-100 px-2 py-0.5 rounded">
                    {log.entityType.replace(/_/g, ' ')}
                  </span>
                </div>
                <div className="flex justify-between items-center py-1.5">
                  <span className="text-gray-500">ID</span>
                  <span className="text-gray-900 font-mono text-xs bg-gray-100 px-2 py-0.5 rounded">{log.entityId}</span>
                </div>
              </div>
            </div>

            {/* Request Info */}
            <div className="rounded-xl border border-gray-100 bg-white p-4 shadow-sm">
              <h3 className="text-sm font-semibold text-gray-700 mb-3 flex items-center gap-2">
                <Server className="h-4 w-4 text-slate-500" />
                Request Information
              </h3>
              <div className="space-y-2 text-sm">
                <div className="flex justify-between items-center py-1.5 border-b border-gray-50">
                  <span className="text-gray-500 flex items-center gap-1.5">
                    <Globe className="h-3 w-3" /> IP Address
                  </span>
                  <span className="text-gray-900 font-mono bg-gray-100 px-2 py-0.5 rounded">{log.ipAddress || '-'}</span>
                </div>
                <div className="flex justify-between items-center py-1.5 border-b border-gray-50">
                  <span className="text-gray-500">Request ID</span>
                  <span className="text-gray-900 font-mono text-xs bg-gray-100 px-2 py-0.5 rounded">{log.requestId || '-'}</span>
                </div>
                {log.userAgent && (
                  <div className="pt-1.5">
                    <span className="text-gray-500 block mb-1.5 text-xs font-medium">User Agent</span>
                    <span className="text-gray-700 text-xs break-all bg-gray-50 p-2 rounded block">{log.userAgent}</span>
                  </div>
                )}
              </div>
            </div>

            {/* State Changes */}
            {(log.beforeState || log.afterState) && (
              <div className="rounded-xl border border-gray-100 bg-white p-4 shadow-sm">
                <h3 className="text-sm font-semibold text-gray-700 mb-3">State Changes</h3>
                <div className="space-y-3">
                  {log.beforeState && (
                    <div>
                      <span className="text-xs font-semibold text-rose-600 block mb-1.5 flex items-center gap-1">
                        <span className="w-2 h-2 rounded-full bg-rose-500"></span>
                        Before
                      </span>
                      <pre className="text-xs bg-gradient-to-r from-rose-50 to-red-50 text-rose-800 p-3 rounded-lg overflow-auto max-h-40 border border-rose-100">
                        {JSON.stringify(log.beforeState, null, 2)}
                      </pre>
                    </div>
                  )}
                  {log.afterState && (
                    <div>
                      <span className="text-xs font-semibold text-emerald-600 block mb-1.5 flex items-center gap-1">
                        <span className="w-2 h-2 rounded-full bg-emerald-500"></span>
                        After
                      </span>
                      <pre className="text-xs bg-gradient-to-r from-emerald-50 to-green-50 text-emerald-800 p-3 rounded-lg overflow-auto max-h-40 border border-emerald-100">
                        {JSON.stringify(log.afterState, null, 2)}
                      </pre>
                    </div>
                  )}
                </div>
              </div>
            )}
          </div>
        )}
      </SheetBody>
      <SheetFooter>
        <Button variant="outline" onClick={onClose}>
          Close
        </Button>
      </SheetFooter>
    </Sheet>
  );
}
