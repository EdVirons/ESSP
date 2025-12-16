import * as React from 'react';
import { type ColumnDef } from '@tanstack/react-table';
import { Eye, Plus, Pencil, Trash2, User } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { SortableHeader } from '@/components/ui/data-table';
import { formatDate, formatRelativeTime, cn } from '@/lib/utils';
import type { AuditLog, AuditAction } from '@/types';

export const actionOptions = [
  { value: '', label: 'All Actions' },
  { value: 'create', label: 'Create' },
  { value: 'update', label: 'Update' },
  { value: 'delete', label: 'Delete' },
];

export const actionColors: Record<AuditAction, string> = {
  create: 'bg-emerald-100 text-emerald-700 border border-emerald-200',
  update: 'bg-cyan-100 text-cyan-700 border border-cyan-200',
  delete: 'bg-rose-100 text-rose-700 border border-rose-200',
};

export const actionIcons: Record<AuditAction, React.ComponentType<{ className?: string }>> = {
  create: Plus,
  update: Pencil,
  delete: Trash2,
};

interface CreateColumnsOptions {
  onViewDetail: (log: AuditLog) => void;
}

export function createAuditLogColumns({ onViewDetail }: CreateColumnsOptions): ColumnDef<AuditLog>[] {
  return [
    {
      accessorKey: 'createdAt',
      header: ({ column }) => <SortableHeader column={column}>Timestamp</SortableHeader>,
      cell: ({ row }) => (
        <div className="text-sm">
          <div className="font-medium text-gray-900">
            {formatRelativeTime(row.original.createdAt)}
          </div>
          <div className="text-gray-500">
            {formatDate(row.original.createdAt)}
          </div>
        </div>
      ),
    },
    {
      accessorKey: 'action',
      header: 'Action',
      cell: ({ row }) => {
        const ActionIcon = actionIcons[row.original.action];
        return (
          <div className="flex items-center gap-2">
            <div className={cn(
              'flex h-8 w-8 items-center justify-center rounded-lg shadow-sm',
              row.original.action === 'create' && 'bg-gradient-to-br from-emerald-500 to-green-600',
              row.original.action === 'update' && 'bg-gradient-to-br from-cyan-500 to-blue-600',
              row.original.action === 'delete' && 'bg-gradient-to-br from-rose-500 to-red-600'
            )}>
              <ActionIcon className="h-4 w-4 text-white" />
            </div>
            <Badge className={cn('capitalize', actionColors[row.original.action])}>
              {row.original.action}
            </Badge>
          </div>
        );
      },
    },
    {
      accessorKey: 'entityType',
      header: 'Entity',
      cell: ({ row }) => (
        <div className="text-sm">
          <div className="font-medium text-gray-900 capitalize">
            {row.original.entityType.replace(/_/g, ' ')}
          </div>
          <div className="text-gray-500 font-mono text-xs truncate max-w-[150px]">
            {row.original.entityId}
          </div>
        </div>
      ),
    },
    {
      accessorKey: 'userEmail',
      header: 'User',
      cell: ({ row }) => (
        <div className="flex items-center gap-2">
          <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-gradient-to-br from-slate-100 to-gray-200 shadow-sm">
            <User className="h-4 w-4 text-slate-600" />
          </div>
          <span className="text-sm text-gray-900 truncate max-w-[150px] font-medium">
            {row.original.userEmail}
          </span>
        </div>
      ),
    },
    {
      accessorKey: 'ipAddress',
      header: 'IP Address',
      cell: ({ row }) => (
        <span className="text-sm text-gray-500 font-mono">
          {row.original.ipAddress || '-'}
        </span>
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
            onViewDetail(row.original);
          }}
        >
          <Eye className="h-4 w-4" />
        </Button>
      ),
    },
  ];
}
