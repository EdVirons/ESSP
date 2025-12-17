import type { ColumnDef } from '@tanstack/react-table';
import { MoreHorizontal, Eye, Pencil, Trash2 } from 'lucide-react';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import type { ServiceStaff, StaffRole } from '@/types';
import { cn } from '@/lib/utils';

// Role labels - assistant_technician displays as "Field Technician"
export const staffRoleLabels: Record<StaffRole, string> = {
  lead_technician: 'Lead Technician',
  assistant_technician: 'Field Technician',
  storekeeper: 'Storekeeper',
};

export const staffRoleColors: Record<StaffRole, string> = {
  lead_technician: 'bg-purple-100 text-purple-800',
  assistant_technician: 'bg-blue-100 text-blue-800',
  storekeeper: 'bg-green-100 text-green-800',
};

export const roleOptions = [
  { value: '', label: 'All Roles' },
  { value: 'lead_technician', label: 'Lead Technician' },
  { value: 'assistant_technician', label: 'Field Technician' },
  { value: 'storekeeper', label: 'Storekeeper' },
];

export const statusOptions = [
  { value: '', label: 'All Status' },
  { value: 'true', label: 'Active' },
  { value: 'false', label: 'Inactive' },
];

interface CreateStaffColumnsOptions {
  shopLookup: Map<string, string>;
  onEdit: (staff: ServiceStaff) => void;
  onDelete: (staff: ServiceStaff) => void;
  onViewDetail: (staff: ServiceStaff) => void;
}

export function createStaffColumns({
  shopLookup,
  onEdit,
  onDelete,
  onViewDetail,
}: CreateStaffColumnsOptions): ColumnDef<ServiceStaff>[] {
  return [
    {
      accessorKey: 'userId',
      header: 'Staff ID / Name',
      cell: ({ row }) => (
        <div className="font-medium">{row.original.userId}</div>
      ),
    },
    {
      accessorKey: 'serviceShopId',
      header: 'Service Shop',
      cell: ({ row }) => (
        <div className="text-gray-600">
          {shopLookup.get(row.original.serviceShopId) || row.original.serviceShopId}
        </div>
      ),
    },
    {
      accessorKey: 'role',
      header: 'Role',
      cell: ({ row }) => (
        <Badge className={cn('text-xs', staffRoleColors[row.original.role])}>
          {staffRoleLabels[row.original.role]}
        </Badge>
      ),
    },
    {
      accessorKey: 'phone',
      header: 'Phone',
      cell: ({ row }) => (
        <div className="text-gray-600">{row.original.phone || '-'}</div>
      ),
    },
    {
      accessorKey: 'active',
      header: 'Status',
      cell: ({ row }) => (
        <Badge variant={row.original.active ? 'default' : 'secondary'}>
          {row.original.active ? 'Active' : 'Inactive'}
        </Badge>
      ),
    },
    {
      id: 'actions',
      header: '',
      cell: ({ row }) => (
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="ghost" className="h-8 w-8 p-0">
              <span className="sr-only">Open menu</span>
              <MoreHorizontal className="h-4 w-4" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            <DropdownMenuItem onClick={() => onViewDetail(row.original)}>
              <Eye className="mr-2 h-4 w-4" />
              View Details
            </DropdownMenuItem>
            <DropdownMenuItem onClick={() => onEdit(row.original)}>
              <Pencil className="mr-2 h-4 w-4" />
              Edit
            </DropdownMenuItem>
            <DropdownMenuItem
              onClick={() => onDelete(row.original)}
              className="text-red-600"
            >
              <Trash2 className="mr-2 h-4 w-4" />
              Remove
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      ),
    },
  ];
}
