import { Sheet, SheetHeader, SheetBody, SheetFooter } from '@/components/ui/sheet';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Users, Phone, MapPin, Calendar, Pencil } from 'lucide-react';
import type { ServiceStaff } from '@/types';
import { staffRoleLabels, staffRoleColors } from './columns';
import { cn } from '@/lib/utils';
import { format } from 'date-fns';

interface StaffDetailProps {
  staff: ServiceStaff | null;
  shopName?: string;
  open: boolean;
  onClose: () => void;
  onEditClick: () => void;
}

export function StaffDetail({
  staff,
  shopName,
  open,
  onClose,
  onEditClick,
}: StaffDetailProps) {
  if (!staff) return null;

  return (
    <Sheet open={open} onClose={onClose} side="right" className="w-full max-w-lg">
      <SheetHeader onClose={onClose}>Staff Details</SheetHeader>
      <SheetBody>
        <div className="space-y-6">
          {/* Header */}
          <div className="flex items-start justify-between">
            <div className="flex items-center gap-3">
              <div className="flex h-12 w-12 items-center justify-center rounded-full bg-blue-100">
                <Users className="h-6 w-6 text-blue-600" />
              </div>
              <div>
                <h2 className="text-xl font-semibold text-gray-900">{staff.userId}</h2>
                <Badge className={cn('text-xs mt-1', staffRoleColors[staff.role])}>
                  {staffRoleLabels[staff.role]}
                </Badge>
              </div>
            </div>
          </div>

          {/* Status */}
          <div className="flex items-center justify-between py-3 border-b border-gray-100">
            <span className="text-sm text-gray-500">Status</span>
            <Badge variant={staff.active ? 'default' : 'secondary'}>
              {staff.active ? 'Active' : 'Inactive'}
            </Badge>
          </div>

          {/* Details */}
          <div className="space-y-4">
            <h3 className="font-medium text-gray-900">Details</h3>

            <div className="flex items-center gap-3 text-sm">
              <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-gray-100">
                <MapPin className="h-4 w-4 text-gray-600" />
              </div>
              <div>
                <div className="text-gray-500">Service Shop</div>
                <div className="font-medium">{shopName || staff.serviceShopId}</div>
              </div>
            </div>

            {staff.phone && (
              <div className="flex items-center gap-3 text-sm">
                <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-gray-100">
                  <Phone className="h-4 w-4 text-gray-600" />
                </div>
                <div>
                  <div className="text-gray-500">Phone</div>
                  <div className="font-medium">{staff.phone}</div>
                </div>
              </div>
            )}

            <div className="flex items-center gap-3 text-sm">
              <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-gray-100">
                <Calendar className="h-4 w-4 text-gray-600" />
              </div>
              <div>
                <div className="text-gray-500">Added</div>
                <div className="font-medium">
                  {format(new Date(staff.createdAt), 'PPP')}
                </div>
              </div>
            </div>

            <div className="flex items-center gap-3 text-sm">
              <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-gray-100">
                <Calendar className="h-4 w-4 text-gray-600" />
              </div>
              <div>
                <div className="text-gray-500">Last Updated</div>
                <div className="font-medium">
                  {format(new Date(staff.updatedAt), 'PPP')}
                </div>
              </div>
            </div>
          </div>

          {/* System Info */}
          <div className="space-y-2 pt-4 border-t border-gray-100">
            <h3 className="font-medium text-gray-900">System Information</h3>
            <div className="text-sm">
              <div className="flex justify-between py-1">
                <span className="text-gray-500">Staff ID</span>
                <span className="font-mono text-xs">{staff.id}</span>
              </div>
              <div className="flex justify-between py-1">
                <span className="text-gray-500">Shop ID</span>
                <span className="font-mono text-xs">{staff.serviceShopId}</span>
              </div>
            </div>
          </div>
        </div>
      </SheetBody>
      <SheetFooter>
        <Button variant="outline" onClick={onClose}>
          Close
        </Button>
        <Button onClick={onEditClick}>
          <Pencil className="h-4 w-4 mr-1" />
          Edit
        </Button>
      </SheetFooter>
    </Sheet>
  );
}
