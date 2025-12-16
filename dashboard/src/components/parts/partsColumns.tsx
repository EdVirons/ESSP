import { type ColumnDef } from '@tanstack/react-table';
import {
  Package,
  Pencil,
  Trash2,
  Layers,
  DollarSign,
  CheckCircle,
  XCircle,
  Building2,
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { SortableHeader } from '@/components/ui/data-table';
import { formatDate, formatCurrency } from '@/lib/utils';
import type { Part } from '@/types';

export const categoryColors: Record<string, string> = {
  electronics: 'bg-blue-100 text-blue-800',
  mechanical: 'bg-orange-100 text-orange-800',
  consumable: 'bg-green-100 text-green-800',
  accessory: 'bg-purple-100 text-purple-800',
  replacement: 'bg-yellow-100 text-yellow-800',
  cables: 'bg-cyan-100 text-cyan-800',
  hardware: 'bg-amber-100 text-amber-800',
};

interface CreateColumnsOptions {
  onEdit: (part: Part) => void;
  onDelete: (part: Part) => void;
}

export function createPartsColumns({ onEdit, onDelete }: CreateColumnsOptions): ColumnDef<Part>[] {
  return [
    {
      accessorKey: 'name',
      header: ({ column }) => <SortableHeader column={column}>Part Name</SortableHeader>,
      cell: ({ row }) => (
        <div className="flex items-center gap-3">
          <div className="flex h-8 w-8 items-center justify-center rounded-full bg-amber-50">
            <Package className="h-4 w-4 text-amber-600" />
          </div>
          <div className="min-w-0">
            <div className="font-medium text-gray-900 truncate max-w-[200px]">
              {row.original.name}
            </div>
            <div className="text-sm text-gray-500 font-mono">{row.original.sku}</div>
          </div>
        </div>
      ),
    },
    {
      accessorKey: 'category',
      header: 'Category',
      cell: ({ row }) => {
        const category = row.original.category || 'unknown';
        return (
          <Badge
            className={categoryColors[category.toLowerCase()] || 'bg-gray-100 text-gray-800'}
          >
            <span className="flex items-center gap-1">
              <Layers className="h-3 w-3" />
              {category.charAt(0).toUpperCase() + category.slice(1)}
            </span>
          </Badge>
        );
      },
    },
    {
      accessorKey: 'unitCostCents',
      header: ({ column }) => <SortableHeader column={column}>Unit Price</SortableHeader>,
      cell: ({ row }) => (
        <div className="flex items-center gap-1">
          <DollarSign className="h-4 w-4 text-gray-400" />
          <span className="font-medium">
            {formatCurrency(row.original.unitCostCents / 100)}
          </span>
        </div>
      ),
    },
    {
      accessorKey: 'supplier',
      header: 'Supplier',
      cell: ({ row }) => (
        <div className="min-w-0">
          {row.original.supplier ? (
            <div className="flex items-center gap-2">
              <Building2 className="h-4 w-4 text-gray-400" />
              <div>
                <div className="text-sm text-gray-900 truncate max-w-[150px]">
                  {row.original.supplier}
                </div>
                {row.original.supplierSku && (
                  <div className="text-xs text-gray-500 font-mono">
                    {row.original.supplierSku}
                  </div>
                )}
              </div>
            </div>
          ) : (
            <span className="text-gray-400">-</span>
          )}
        </div>
      ),
    },
    {
      accessorKey: 'active',
      header: 'Status',
      cell: ({ row }) => (
        <Badge
          className={
            row.original.active ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-600'
          }
        >
          {row.original.active ? (
            <span className="flex items-center gap-1">
              <CheckCircle className="h-3 w-3" />
              Active
            </span>
          ) : (
            <span className="flex items-center gap-1">
              <XCircle className="h-3 w-3" />
              Inactive
            </span>
          )}
        </Badge>
      ),
    },
    {
      accessorKey: 'updatedAt',
      header: ({ column }) => <SortableHeader column={column}>Updated</SortableHeader>,
      cell: ({ row }) => (
        <span className="text-gray-600 text-sm">{formatDate(row.original.updatedAt)}</span>
      ),
    },
    {
      id: 'actions',
      header: '',
      cell: ({ row }) => (
        <div className="flex items-center gap-1">
          <Button variant="ghost" size="sm" onClick={() => onEdit(row.original)}>
            <Pencil className="h-4 w-4" />
          </Button>
          <Button variant="ghost" size="sm" onClick={() => onDelete(row.original)}>
            <Trash2 className="h-4 w-4 text-red-500" />
          </Button>
        </div>
      ),
    },
  ];
}
