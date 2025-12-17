import * as React from 'react';
import { Search, MapPin, Wifi } from 'lucide-react';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import type { InventoryDevice } from '@/types';
import { LIFECYCLE_STATUS_COLORS } from '@/types';

interface DeviceListProps {
  devices: InventoryDevice[];
  loading?: boolean;
  onDeviceClick?: (device: InventoryDevice) => void;
}

export function DeviceList({ devices, loading, onDeviceClick }: DeviceListProps) {
  const [searchQuery, setSearchQuery] = React.useState('');

  const filteredDevices = React.useMemo(() => {
    if (!searchQuery) return devices;
    const query = searchQuery.toLowerCase();
    return devices.filter(
      (d) =>
        d.serial.toLowerCase().includes(query) ||
        d.assetTag.toLowerCase().includes(query) ||
        d.model.toLowerCase().includes(query) ||
        d.locationPath?.toLowerCase().includes(query)
    );
  }, [devices, searchQuery]);

  if (loading) {
    return (
      <div className="space-y-3">
        {[1, 2, 3, 4, 5].map((i) => (
          <div key={i} className="h-12 animate-pulse bg-gray-100 rounded" />
        ))}
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {/* Search */}
      <div className="relative">
        <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-gray-400" />
        <Input
          placeholder="Search by serial, asset tag, model, or location..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className="pl-10"
        />
      </div>

      {/* Table */}
      <div className="border rounded-lg">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Device</TableHead>
              <TableHead>Model</TableHead>
              <TableHead>Location</TableHead>
              <TableHead>Status</TableHead>
              <TableHead>MAC Address</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {filteredDevices.length === 0 ? (
              <TableRow>
                <TableCell colSpan={5} className="text-center py-8 text-gray-500">
                  {searchQuery ? 'No devices match your search' : 'No devices found'}
                </TableCell>
              </TableRow>
            ) : (
              filteredDevices.map((device) => (
                <TableRow
                  key={device.id}
                  className="cursor-pointer hover:bg-gray-50"
                  onClick={() => onDeviceClick?.(device)}
                >
                  <TableCell>
                    <div>
                      <p className="font-medium text-gray-900">{device.assetTag || device.serial}</p>
                      <p className="text-sm text-gray-500">SN: {device.serial}</p>
                    </div>
                  </TableCell>
                  <TableCell>
                    <div>
                      <p className="text-gray-900">{device.model}</p>
                      <p className="text-sm text-gray-500">{device.make}</p>
                    </div>
                  </TableCell>
                  <TableCell>
                    {device.locationPath ? (
                      <div className="flex items-center gap-2">
                        <MapPin className="h-4 w-4 text-gray-400" />
                        <span className="text-gray-700">{device.locationPath}</span>
                      </div>
                    ) : (
                      <span className="text-gray-400 italic">Unassigned</span>
                    )}
                  </TableCell>
                  <TableCell>
                    <Badge
                      variant="secondary"
                      className={LIFECYCLE_STATUS_COLORS[device.lifecycle as keyof typeof LIFECYCLE_STATUS_COLORS] || 'bg-gray-100 text-gray-600'}
                    >
                      {device.lifecycle}
                    </Badge>
                  </TableCell>
                  <TableCell>
                    {device.macAddresses && device.macAddresses.length > 0 ? (
                      <div className="flex items-center gap-2">
                        <Wifi className="h-4 w-4 text-gray-400" />
                        <span className="font-mono text-sm text-gray-600">
                          {device.macAddresses[0]}
                          {device.macAddresses.length > 1 && (
                            <span className="text-gray-400"> +{device.macAddresses.length - 1}</span>
                          )}
                        </span>
                      </div>
                    ) : (
                      <span className="text-gray-400">—</span>
                    )}
                  </TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </div>

      {/* Count */}
      <p className="text-sm text-gray-500">
        Showing {filteredDevices.length} of {devices.length} devices
      </p>
    </div>
  );
}
