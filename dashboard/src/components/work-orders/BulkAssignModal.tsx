import { useState } from 'react';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { Label } from '@/components/ui/label';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Loader2 } from 'lucide-react';
import { useBulkAssignment } from '@/api/work-orders';

interface BulkAssignModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  selectedIds: string[];
  onSuccess: () => void;
}

// These would typically come from an API
const staffOptions = [
  { value: '', label: 'No change' },
  { value: 'staff-1', label: 'John Doe' },
  { value: 'staff-2', label: 'Jane Smith' },
  { value: 'staff-3', label: 'Bob Wilson' },
];

const shopOptions = [
  { value: '', label: 'No change' },
  { value: 'shop-1', label: 'Main Service Shop' },
  { value: 'shop-2', label: 'Regional Shop A' },
  { value: 'shop-3', label: 'Regional Shop B' },
];

export function BulkAssignModal({
  open,
  onOpenChange,
  selectedIds,
  onSuccess,
}: BulkAssignModalProps) {
  const bulkAssign = useBulkAssignment();
  const [staffId, setStaffId] = useState('');
  const [shopId, setShopId] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!staffId && !shopId) return;

    try {
      await bulkAssign.mutateAsync({
        workOrderIds: selectedIds,
        assignedStaffId: staffId || undefined,
        serviceShopId: shopId || undefined,
      });
      onOpenChange(false);
      setStaffId('');
      setShopId('');
      onSuccess();
    } catch {
      // Error handled by mutation
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[425px]">
        <form onSubmit={handleSubmit}>
          <DialogHeader>
            <DialogTitle>Bulk Assignment</DialogTitle>
            <DialogDescription>
              Assign {selectedIds.length} work order{selectedIds.length > 1 ? 's' : ''} to a
              technician or service shop.
            </DialogDescription>
          </DialogHeader>

          <div className="grid gap-4 py-4">
            <div className="grid gap-2">
              <Label htmlFor="staff">Assign to Technician</Label>
              <Select value={staffId} onValueChange={setStaffId}>
                <SelectTrigger id="staff">
                  <SelectValue placeholder="Select technician" />
                </SelectTrigger>
                <SelectContent>
                  {staffOptions.map((option) => (
                    <SelectItem key={option.value || 'none'} value={option.value}>
                      {option.label}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>

            <div className="grid gap-2">
              <Label htmlFor="shop">Assign to Service Shop</Label>
              <Select value={shopId} onValueChange={setShopId}>
                <SelectTrigger id="shop">
                  <SelectValue placeholder="Select service shop" />
                </SelectTrigger>
                <SelectContent>
                  {shopOptions.map((option) => (
                    <SelectItem key={option.value || 'none'} value={option.value}>
                      {option.label}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
          </div>

          <DialogFooter>
            <Button
              type="button"
              variant="outline"
              onClick={() => onOpenChange(false)}
              disabled={bulkAssign.isPending}
            >
              Cancel
            </Button>
            <Button type="submit" disabled={(!staffId && !shopId) || bulkAssign.isPending}>
              {bulkAssign.isPending ? (
                <>
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  Assigning...
                </>
              ) : (
                `Assign ${selectedIds.length} Work Order${selectedIds.length > 1 ? 's' : ''}`
              )}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
