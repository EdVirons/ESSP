import { FileText, Plus, Pencil, Trash2 } from 'lucide-react';
import { Card, CardContent } from '@/components/ui/card';
import type { AuditLog } from '@/types';

interface AuditLogStatsProps {
  logs: AuditLog[];
}

export function AuditLogStats({ logs }: AuditLogStatsProps) {
  const createCount = logs.filter((l) => l.action === 'create').length;
  const updateCount = logs.filter((l) => l.action === 'update').length;
  const deleteCount = logs.filter((l) => l.action === 'delete').length;

  return (
    <div className="grid gap-4 md:grid-cols-4">
      <Card className="border-0 shadow-md overflow-hidden">
        <CardContent className="p-0">
          <div className="flex items-center gap-4 p-4">
            <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-gradient-to-br from-slate-500 to-slate-700 shadow-lg shadow-slate-500/20">
              <FileText className="h-6 w-6 text-white" />
            </div>
            <div>
              <div className="text-2xl font-bold text-gray-900">{logs.length}</div>
              <div className="text-sm text-gray-500">Total Events</div>
            </div>
          </div>
          <div className="h-1 bg-gradient-to-r from-slate-400 to-slate-600" />
        </CardContent>
      </Card>
      <Card className="border-0 shadow-md overflow-hidden">
        <CardContent className="p-0">
          <div className="flex items-center gap-4 p-4">
            <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-gradient-to-br from-emerald-500 to-green-600 shadow-lg shadow-emerald-500/20">
              <Plus className="h-6 w-6 text-white" />
            </div>
            <div>
              <div className="text-2xl font-bold text-gray-900">{createCount}</div>
              <div className="text-sm text-gray-500">Creates</div>
            </div>
          </div>
          <div className="h-1 bg-gradient-to-r from-emerald-400 to-green-500" />
        </CardContent>
      </Card>
      <Card className="border-0 shadow-md overflow-hidden">
        <CardContent className="p-0">
          <div className="flex items-center gap-4 p-4">
            <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-gradient-to-br from-cyan-500 to-blue-600 shadow-lg shadow-cyan-500/20">
              <Pencil className="h-6 w-6 text-white" />
            </div>
            <div>
              <div className="text-2xl font-bold text-gray-900">{updateCount}</div>
              <div className="text-sm text-gray-500">Updates</div>
            </div>
          </div>
          <div className="h-1 bg-gradient-to-r from-cyan-400 to-blue-500" />
        </CardContent>
      </Card>
      <Card className="border-0 shadow-md overflow-hidden">
        <CardContent className="p-0">
          <div className="flex items-center gap-4 p-4">
            <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-gradient-to-br from-rose-500 to-red-600 shadow-lg shadow-rose-500/20">
              <Trash2 className="h-6 w-6 text-white" />
            </div>
            <div>
              <div className="text-2xl font-bold text-gray-900">{deleteCount}</div>
              <div className="text-sm text-gray-500">Deletes</div>
            </div>
          </div>
          <div className="h-1 bg-gradient-to-r from-rose-400 to-red-500" />
        </CardContent>
      </Card>
    </div>
  );
}
