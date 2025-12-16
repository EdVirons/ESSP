import { Plus, Users, CheckCircle2, XCircle } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { cn } from '@/lib/utils';
import type { ServiceStaff, StaffRole } from '@/types';

const staffRoleLabels: Record<StaffRole, string> = {
  lead_technician: 'Lead Technician',
  assistant_technician: 'Assistant Technician',
  storekeeper: 'Storekeeper',
};

const staffRoleColors: Record<StaffRole, string> = {
  lead_technician: 'bg-purple-100 text-purple-800',
  assistant_technician: 'bg-blue-100 text-blue-800',
  storekeeper: 'bg-green-100 text-green-800',
};

interface StaffListProps {
  staff: ServiceStaff[];
  onAddClick?: () => void;
}

export function StaffList({ staff, onAddClick }: StaffListProps) {
  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h3 className="font-medium text-gray-900">Staff Members</h3>
        <Button size="sm" onClick={onAddClick}>
          <Plus className="h-4 w-4" />
          Add Staff
        </Button>
      </div>
      {staff.length === 0 ? (
        <div className="text-center py-8 text-gray-500">
          <Users className="h-12 w-12 mx-auto mb-2 text-gray-300" />
          <p>No staff members assigned</p>
        </div>
      ) : (
        <div className="space-y-2">
          {staff.map((member) => (
            <div
              key={member.id}
              className="flex items-center justify-between p-3 bg-gray-50 rounded-lg"
            >
              <div className="flex items-center gap-3">
                <div className="flex h-10 w-10 items-center justify-center rounded-full bg-gray-200">
                  <Users className="h-5 w-5 text-gray-600" />
                </div>
                <div>
                  <div className="font-medium text-gray-900">{member.userId}</div>
                  <div className="text-sm text-gray-500">{member.phone}</div>
                </div>
              </div>
              <div className="flex items-center gap-2">
                <Badge className={cn('text-xs', staffRoleColors[member.role])}>
                  {staffRoleLabels[member.role]}
                </Badge>
                {member.active ? (
                  <CheckCircle2 className="h-4 w-4 text-green-500" />
                ) : (
                  <XCircle className="h-4 w-4 text-gray-400" />
                )}
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

export { staffRoleLabels, staffRoleColors };
