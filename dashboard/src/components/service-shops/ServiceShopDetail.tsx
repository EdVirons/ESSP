import { MapPin, Wrench, Calendar, Clock, Building2 } from 'lucide-react';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Sheet, SheetHeader, SheetBody, SheetFooter } from '@/components/ui/sheet';
import { Tabs, TabsList, TabsTrigger, TabsContent } from '@/components/ui/tabs';
import { formatDate } from '@/lib/utils';
import { StaffList } from './StaffList';
import { InventoryList } from './InventoryList';
import type { ServiceShop, ServiceStaff } from '@/types';

interface InventoryItem {
  id: string;
  partName: string;
  partPuk: string;
  qtyOnHand: number;
  qtyAvailable: number;
  reorderLevel: number;
}

interface WorkOrder {
  id: string;
  title: string;
  status: string;
  priority: string;
  createdAt: string;
}

interface ServiceShopDetailProps {
  shop: ServiceShop | null;
  open: boolean;
  onClose: () => void;
  detailTab: string;
  onDetailTabChange: (tab: string) => void;
  staff: ServiceStaff[];
  inventory: InventoryItem[];
  workOrders?: WorkOrder[];
  onEditClick?: () => void;
  onAddStaffClick?: () => void;
  onAddInventoryClick?: () => void;
}

const workOrderStatusColors: Record<string, string> = {
  open: 'bg-blue-100 text-blue-800',
  in_progress: 'bg-yellow-100 text-yellow-800',
  completed: 'bg-green-100 text-green-800',
  cancelled: 'bg-gray-100 text-gray-800',
};

const priorityColors: Record<string, string> = {
  low: 'bg-gray-100 text-gray-600',
  medium: 'bg-blue-100 text-blue-600',
  high: 'bg-orange-100 text-orange-600',
  urgent: 'bg-red-100 text-red-600',
};

export function ServiceShopDetail({
  shop,
  open,
  onClose,
  detailTab,
  onDetailTabChange,
  staff,
  inventory,
  workOrders = [],
  onEditClick,
  onAddStaffClick,
  onAddInventoryClick,
}: ServiceShopDetailProps) {
  if (!shop) return null;

  const lowStockItems = inventory.filter((i) => i.qtyAvailable <= i.reorderLevel);
  const activeStaff = staff.filter((s) => s.active);

  return (
    <Sheet open={open} onClose={onClose} side="right" className="w-full max-w-xl">
      <SheetHeader onClose={onClose}>Service Shop Details</SheetHeader>
      <SheetBody className="p-0">
        <div className="h-full flex flex-col">
          {/* Header Info */}
          <div className="p-6 border-b border-gray-200 bg-gradient-to-r from-orange-50 to-amber-50">
            <div className="flex items-start justify-between mb-3">
              <div className="flex items-center gap-2">
                <Badge
                  className={
                    shop.active
                      ? 'bg-green-100 text-green-800'
                      : 'bg-gray-100 text-gray-800'
                  }
                >
                  {shop.active ? 'Active' : 'Inactive'}
                </Badge>
                <Badge variant="outline" className="capitalize">
                  {shop.coverageLevel.replace('_', ' ')}
                </Badge>
              </div>
              <div className="flex h-12 w-12 items-center justify-center rounded-full bg-orange-100">
                <Building2 className="h-6 w-6 text-orange-600" />
              </div>
            </div>
            <h2 className="text-xl font-bold text-gray-900 mb-1">{shop.name}</h2>
            <p className="text-sm text-gray-600 flex items-center gap-1">
              <MapPin className="h-4 w-4" />
              {shop.location || 'No location specified'}
            </p>

            {/* Quick Stats */}
            <div className="grid grid-cols-3 gap-4 mt-4">
              <div className="bg-white rounded-lg p-3 shadow-sm">
                <div className="text-2xl font-bold text-gray-900">{activeStaff.length}</div>
                <div className="text-xs text-gray-500">Staff</div>
              </div>
              <div className="bg-white rounded-lg p-3 shadow-sm">
                <div className="text-2xl font-bold text-gray-900">{inventory.length}</div>
                <div className="text-xs text-gray-500">Items</div>
              </div>
              <div className={`rounded-lg p-3 shadow-sm ${lowStockItems.length > 0 ? 'bg-red-50' : 'bg-white'}`}>
                <div className={`text-2xl font-bold ${lowStockItems.length > 0 ? 'text-red-600' : 'text-gray-900'}`}>
                  {lowStockItems.length}
                </div>
                <div className="text-xs text-gray-500">Low Stock</div>
              </div>
            </div>
          </div>

          {/* Tabs */}
          <Tabs
            value={detailTab}
            onValueChange={onDetailTabChange}
            className="flex-1 flex flex-col"
          >
            <div className="border-b border-gray-200 px-6">
              <TabsList className="bg-transparent -mb-px">
                <TabsTrigger value="staff">Staff ({staff.length})</TabsTrigger>
                <TabsTrigger value="inventory">
                  Inventory ({inventory.length})
                </TabsTrigger>
                <TabsTrigger value="work-orders">
                  Work Orders ({workOrders.length})
                </TabsTrigger>
                <TabsTrigger value="details">Details</TabsTrigger>
              </TabsList>
            </div>

            <div className="flex-1 overflow-auto">
              <TabsContent value="staff" className="p-6 m-0">
                <StaffList staff={staff} onAddClick={onAddStaffClick} />
              </TabsContent>

              <TabsContent value="inventory" className="p-6 m-0">
                <InventoryList inventory={inventory} onAddClick={onAddInventoryClick} />
              </TabsContent>

              <TabsContent value="work-orders" className="p-6 m-0">
                <div className="space-y-4">
                  <div className="flex items-center justify-between">
                    <h3 className="font-medium text-gray-900">Work Orders</h3>
                  </div>
                  {workOrders.length === 0 ? (
                    <div className="text-center py-8 text-gray-500">
                      <Wrench className="h-12 w-12 mx-auto mb-2 text-gray-300" />
                      <p>No work orders assigned</p>
                    </div>
                  ) : (
                    <div className="space-y-2">
                      {workOrders.map((wo) => (
                        <div
                          key={wo.id}
                          className="flex items-center justify-between p-3 bg-gray-50 rounded-lg hover:bg-gray-100 cursor-pointer"
                        >
                          <div>
                            <div className="font-medium text-gray-900">{wo.title}</div>
                            <div className="text-sm text-gray-500 flex items-center gap-2">
                              <Clock className="h-3 w-3" />
                              {formatDate(wo.createdAt)}
                            </div>
                          </div>
                          <div className="flex items-center gap-2">
                            <Badge className={priorityColors[wo.priority] || 'bg-gray-100'}>
                              {wo.priority}
                            </Badge>
                            <Badge className={workOrderStatusColors[wo.status] || 'bg-gray-100'}>
                              {wo.status.replace('_', ' ')}
                            </Badge>
                          </div>
                        </div>
                      ))}
                    </div>
                  )}
                </div>
              </TabsContent>

              <TabsContent value="details" className="p-6 m-0">
                <div className="space-y-6">
                  <div>
                    <h3 className="text-sm font-medium text-gray-500 mb-3 flex items-center gap-2">
                      <MapPin className="h-4 w-4" />
                      Location
                    </h3>
                    <div className="grid grid-cols-2 gap-4">
                      <div className="bg-gray-50 rounded-lg p-3">
                        <div className="text-xs text-gray-500 mb-1">County</div>
                        <div className="font-medium text-gray-900">{shop.countyName}</div>
                        <div className="text-sm text-gray-500">{shop.countyCode}</div>
                      </div>
                      <div className="bg-gray-50 rounded-lg p-3">
                        <div className="text-xs text-gray-500 mb-1">Sub-County</div>
                        <div className="font-medium text-gray-900">{shop.subCountyName || '-'}</div>
                        <div className="text-sm text-gray-500">{shop.subCountyCode || '-'}</div>
                      </div>
                    </div>
                  </div>

                  <div>
                    <h3 className="text-sm font-medium text-gray-500 mb-3 flex items-center gap-2">
                      <Building2 className="h-4 w-4" />
                      Address
                    </h3>
                    <div className="bg-gray-50 rounded-lg p-3">
                      <p className="text-gray-900">
                        {shop.location || 'Not specified'}
                      </p>
                    </div>
                  </div>

                  <div>
                    <h3 className="text-sm font-medium text-gray-500 mb-3 flex items-center gap-2">
                      <Calendar className="h-4 w-4" />
                      Timeline
                    </h3>
                    <div className="space-y-2">
                      <div className="flex justify-between items-center py-2 border-b border-gray-100">
                        <span className="text-gray-500">Created</span>
                        <span className="font-medium text-gray-900">
                          {formatDate(shop.createdAt)}
                        </span>
                      </div>
                      <div className="flex justify-between items-center py-2">
                        <span className="text-gray-500">Last Updated</span>
                        <span className="font-medium text-gray-900">
                          {formatDate(shop.updatedAt)}
                        </span>
                      </div>
                    </div>
                  </div>

                  {/* Coverage Info */}
                  <div>
                    <h3 className="text-sm font-medium text-gray-500 mb-3">Coverage Area</h3>
                    <div className="bg-orange-50 rounded-lg p-4 border border-orange-100">
                      <div className="flex items-center gap-2 text-orange-800">
                        <MapPin className="h-5 w-5" />
                        <span className="font-medium capitalize">
                          {shop.coverageLevel.replace('_', ' ')} Level Coverage
                        </span>
                      </div>
                      <p className="text-sm text-orange-600 mt-1">
                        This shop serves all schools within its {shop.coverageLevel.replace('_', ' ')} area.
                      </p>
                    </div>
                  </div>
                </div>
              </TabsContent>
            </div>
          </Tabs>
        </div>
      </SheetBody>
      <SheetFooter>
        <Button variant="outline" onClick={onClose}>
          Close
        </Button>
        <Button onClick={onEditClick}>Edit Shop</Button>
      </SheetFooter>
    </Sheet>
  );
}
