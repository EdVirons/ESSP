import { type ColumnDef } from '@tanstack/react-table';
import { Store, ChevronRight, Pencil, Trash2, Users, Package } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { SortableHeader, createSelectColumn } from '@/components/ui/data-table';
import { formatDate } from '@/lib/utils';
import type { ServiceShop, ServiceStaff } from '@/types';

export const statusOptions = [
  { value: '', label: 'All Status' },
  { value: 'true', label: 'Active' },
  { value: 'false', label: 'Inactive' },
];

export const coverageOptions = [
  { value: '', label: 'All Coverage' },
  { value: 'county', label: 'County' },
  { value: 'sub_county', label: 'Sub-County' },
  { value: 'region', label: 'Region' },
];

// Minimal type for inventory items used in column rendering
interface InventoryItem {
  serviceShopId: string;
}

interface CreateColumnsOptions {
  allStaff?: ServiceStaff[];
  allInventory?: InventoryItem[];
  onEdit: (shop: ServiceShop) => void;
  onDelete: (shop: ServiceShop) => void;
  onViewDetail: (shop: ServiceShop) => void;
}

export function createServiceShopColumns({
  allStaff = [],
  allInventory = [],
  onEdit,
  onDelete,
  onViewDetail,
}: CreateColumnsOptions): ColumnDef<ServiceShop>[] {
  return [
    createSelectColumn<ServiceShop>(),
    {
      accessorKey: 'name',
      header: ({ column }) => <SortableHeader column={column}>Shop Name</SortableHeader>,
      cell: ({ row }) => (
        <div className="flex items-center gap-3">
          <div className="flex h-10 w-10 items-center justify-center rounded-full bg-orange-50">
            <Store className="h-5 w-5 text-orange-600" />
          </div>
          <div className="min-w-0">
            <div className="font-medium text-gray-900 truncate max-w-[200px]">
              {row.original.name}
            </div>
            <div className="text-sm text-gray-500">{row.original.location || 'No location'}</div>
          </div>
        </div>
      ),
    },
    {
      accessorKey: 'active',
      header: 'Status',
      cell: ({ row }) => (
        <Badge
          className={
            row.original.active ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-800'
          }
        >
          {row.original.active ? 'Active' : 'Inactive'}
        </Badge>
      ),
    },
    {
      accessorKey: 'countyName',
      header: 'County',
      cell: ({ row }) => (
        <div className="text-sm">
          <div className="font-medium">{row.original.countyName}</div>
          <div className="text-gray-500">{row.original.subCountyName || '-'}</div>
        </div>
      ),
    },
    {
      accessorKey: 'coverageLevel',
      header: 'Coverage',
      cell: ({ row }) => (
        <Badge variant="outline" className="capitalize">
          {row.original.coverageLevel.replace('_', ' ')}
        </Badge>
      ),
    },
    {
      id: 'stats',
      header: 'Resources',
      cell: ({ row }) => {
        const shopStaff = allStaff.filter((s) => s.serviceShopId === row.original.id);
        const shopInventory = allInventory.filter((i) => i.serviceShopId === row.original.id);
        return (
          <div className="flex items-center gap-3 text-sm">
            <span className="flex items-center gap-1 text-gray-600">
              <Users className="h-4 w-4" />
              {shopStaff.length}
            </span>
            <span className="flex items-center gap-1 text-gray-600">
              <Package className="h-4 w-4" />
              {shopInventory.length}
            </span>
          </div>
        );
      },
    },
    {
      accessorKey: 'createdAt',
      header: ({ column }) => <SortableHeader column={column}>Created</SortableHeader>,
      cell: ({ row }) => (
        <div className="text-sm text-gray-500">{formatDate(row.original.createdAt)}</div>
      ),
    },
    {
      id: 'actions',
      cell: ({ row }) => (
        <div className="flex items-center gap-1">
          <Button
            variant="ghost"
            size="sm"
            onClick={(e) => {
              e.stopPropagation();
              onEdit(row.original);
            }}
          >
            <Pencil className="h-4 w-4" />
          </Button>
          <Button
            variant="ghost"
            size="sm"
            onClick={(e) => {
              e.stopPropagation();
              onDelete(row.original);
            }}
          >
            <Trash2 className="h-4 w-4 text-red-500" />
          </Button>
          <Button
            variant="ghost"
            size="sm"
            onClick={(e) => {
              e.stopPropagation();
              onViewDetail(row.original);
            }}
          >
            <ChevronRight className="h-4 w-4" />
          </Button>
        </div>
      ),
    },
  ];
}
