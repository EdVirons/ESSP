import * as React from 'react';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Checkbox } from '@/components/ui/checkbox';
import { Loader2 } from 'lucide-react';
import type { ServiceShop, ServiceStaff, StaffRole, UpdateServiceStaffRequest } from '@/types';

interface EditStaffModalProps {
  staff: ServiceStaff | null;
  open: boolean;
  onClose: () => void;
  onSubmit: (id: string, data: UpdateServiceStaffRequest) => void;
  isLoading: boolean;
  shops: ServiceShop[];
}

export function EditStaffModal({
  staff,
  open,
  onClose,
  onSubmit,
  isLoading,
  shops,
}: EditStaffModalProps) {
  const [formData, setFormData] = React.useState<UpdateServiceStaffRequest>({});

  React.useEffect(() => {
    if (staff) {
      setFormData({
        serviceShopId: staff.serviceShopId,
        role: staff.role,
        phone: staff.phone,
        active: staff.active,
      });
    }
  }, [staff]);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (staff) {
      onSubmit(staff.id, formData);
    }
  };

  const handleClose = () => {
    setFormData({});
    onClose();
  };

  if (!staff) return null;

  return (
    <Dialog open={open} onOpenChange={handleClose}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Edit Staff Member</DialogTitle>
          <DialogDescription>
            Update information for {staff.userId}.
          </DialogDescription>
        </DialogHeader>
        <form onSubmit={handleSubmit}>
          <div className="grid gap-4 py-4">
            <div className="grid gap-2">
              <Label htmlFor="serviceShopId">Service Shop</Label>
              <Select
                value={formData.serviceShopId}
                onValueChange={(value) =>
                  setFormData((prev) => ({ ...prev, serviceShopId: value }))
                }
              >
                <SelectTrigger>
                  <SelectValue placeholder="Select a shop" />
                </SelectTrigger>
                <SelectContent>
                  {shops.map((shop) => (
                    <SelectItem key={shop.id} value={shop.id}>
                      {shop.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>

            <div className="grid gap-2">
              <Label htmlFor="userId">Staff Name / ID</Label>
              <Input
                id="userId"
                value={staff.userId}
                disabled
                className="bg-gray-50"
              />
              <p className="text-xs text-gray-500">Staff ID cannot be changed</p>
            </div>

            <div className="grid gap-2">
              <Label htmlFor="role">Role</Label>
              <Select
                value={formData.role}
                onValueChange={(value) =>
                  setFormData((prev) => ({ ...prev, role: value as StaffRole }))
                }
              >
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="lead_technician">Lead Technician</SelectItem>
                  <SelectItem value="assistant_technician">Field Technician</SelectItem>
                  <SelectItem value="storekeeper">Storekeeper</SelectItem>
                </SelectContent>
              </Select>
            </div>

            <div className="grid gap-2">
              <Label htmlFor="phone">Phone</Label>
              <Input
                id="phone"
                value={formData.phone || ''}
                onChange={(e) =>
                  setFormData((prev) => ({ ...prev, phone: e.target.value }))
                }
                placeholder="Enter phone number"
              />
            </div>

            <div className="flex items-center gap-2">
              <Checkbox
                id="active"
                checked={formData.active}
                onCheckedChange={(checked: boolean) =>
                  setFormData((prev) => ({ ...prev, active: checked }))
                }
              />
              <Label htmlFor="active">Active</Label>
            </div>
          </div>

          <DialogFooter>
            <Button type="button" variant="outline" onClick={handleClose}>
              Cancel
            </Button>
            <Button type="submit" disabled={isLoading}>
              {isLoading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              Save Changes
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
