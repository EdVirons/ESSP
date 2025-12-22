import { MapPin, User, Users, Clock, ChevronRight } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Modal, ModalHeader, ModalBody, ModalFooter } from '@/components/ui/modal';
import type { InventoryDevice, DeviceGroup } from '@/types';
import { formatDistanceToNow } from 'date-fns';

interface DeviceDetailModalProps {
  open: boolean;
  onClose: () => void;
  device: InventoryDevice | null;
  groups?: DeviceGroup[];
  onAssignLocation: (device: InventoryDevice) => void;
  onAssignUser: (device: InventoryDevice) => void;
  onManageGroups: (device: InventoryDevice) => void;
}

export function DeviceDetailModal({
  open,
  onClose,
  device,
  onAssignLocation,
  onAssignUser,
  onManageGroups,
}: DeviceDetailModalProps) {
  if (!device) return null;

  return (
    <Modal open={open} onClose={onClose} className="max-w-lg">
      <ModalHeader onClose={onClose}>Device Details</ModalHeader>
      <ModalBody>
        <div className="space-y-6">
          {/* Device Info */}
          <div className="bg-gray-50 rounded-lg p-4">
            <div className="grid grid-cols-2 gap-4">
              <div>
                <p className="text-xs text-gray-500 uppercase tracking-wide">Serial Number</p>
                <p className="font-mono font-medium text-gray-900">{device.serial}</p>
              </div>
              {device.assetTag && (
                <div>
                  <p className="text-xs text-gray-500 uppercase tracking-wide">Asset Tag</p>
                  <p className="font-medium text-gray-900">{device.assetTag}</p>
                </div>
              )}
              <div>
                <p className="text-xs text-gray-500 uppercase tracking-wide">Model</p>
                <p className="font-medium text-gray-900">{device.model}</p>
              </div>
              {device.make && (
                <div>
                  <p className="text-xs text-gray-500 uppercase tracking-wide">Make</p>
                  <p className="font-medium text-gray-900">{device.make}</p>
                </div>
              )}
              <div>
                <p className="text-xs text-gray-500 uppercase tracking-wide">Status</p>
                <span className={`inline-flex items-center px-2 py-0.5 rounded text-xs font-medium ${
                  device.lifecycle === 'active' ? 'bg-green-100 text-green-800' :
                  device.lifecycle === 'repair' ? 'bg-yellow-100 text-yellow-800' :
                  'bg-gray-100 text-gray-800'
                }`}>
                  {device.lifecycle}
                </span>
              </div>
              {device.lastSeenAt && (
                <div>
                  <p className="text-xs text-gray-500 uppercase tracking-wide">Last Seen</p>
                  <p className="text-sm text-gray-600">
                    {formatDistanceToNow(new Date(device.lastSeenAt), { addSuffix: true })}
                  </p>
                </div>
              )}
            </div>

            {device.macAddresses && device.macAddresses.length > 0 && (
              <div className="mt-4 pt-4 border-t border-gray-200">
                <p className="text-xs text-gray-500 uppercase tracking-wide mb-1">MAC Addresses</p>
                <div className="flex flex-wrap gap-2">
                  {device.macAddresses.map((mac) => (
                    <span key={mac} className="font-mono text-xs bg-gray-200 px-2 py-1 rounded">
                      {mac}
                    </span>
                  ))}
                </div>
              </div>
            )}
          </div>

          {/* Current Assignment */}
          <div>
            <h4 className="text-sm font-medium text-gray-900 mb-3">Current Assignment</h4>
            <div className="space-y-2">
              {/* Location */}
              <button
                onClick={() => onAssignLocation(device)}
                className="w-full flex items-center justify-between p-3 bg-white border border-gray-200 rounded-lg hover:border-blue-300 hover:bg-blue-50 transition-colors"
              >
                <div className="flex items-center gap-3">
                  <div className="p-2 bg-blue-100 rounded-lg">
                    <MapPin className="h-4 w-4 text-blue-600" />
                  </div>
                  <div className="text-left">
                    <p className="text-sm font-medium text-gray-900">Location</p>
                    <p className="text-xs text-gray-500">
                      {device.location ? device.locationPath || device.location.name : 'Not assigned'}
                    </p>
                  </div>
                </div>
                <ChevronRight className="h-4 w-4 text-gray-400" />
              </button>

              {/* User Assignment */}
              <button
                onClick={() => onAssignUser(device)}
                className="w-full flex items-center justify-between p-3 bg-white border border-gray-200 rounded-lg hover:border-blue-300 hover:bg-blue-50 transition-colors"
              >
                <div className="flex items-center gap-3">
                  <div className="p-2 bg-purple-100 rounded-lg">
                    <User className="h-4 w-4 text-purple-600" />
                  </div>
                  <div className="text-left">
                    <p className="text-sm font-medium text-gray-900">Assigned User</p>
                    <p className="text-xs text-gray-500">Not assigned to a user</p>
                  </div>
                </div>
                <ChevronRight className="h-4 w-4 text-gray-400" />
              </button>

              {/* Groups */}
              <button
                onClick={() => onManageGroups(device)}
                className="w-full flex items-center justify-between p-3 bg-white border border-gray-200 rounded-lg hover:border-blue-300 hover:bg-blue-50 transition-colors"
              >
                <div className="flex items-center gap-3">
                  <div className="p-2 bg-green-100 rounded-lg">
                    <Users className="h-4 w-4 text-green-600" />
                  </div>
                  <div className="text-left">
                    <p className="text-sm font-medium text-gray-900">Device Groups</p>
                    <p className="text-xs text-gray-500">
                      {device.groups && device.groups.length > 0
                        ? `Member of ${device.groups.length} group(s)`
                        : 'Not in any groups'}
                    </p>
                  </div>
                </div>
                <ChevronRight className="h-4 w-4 text-gray-400" />
              </button>
            </div>
          </div>

          {/* Updated Info */}
          {device.updatedAt && (
            <div className="flex items-center gap-2 text-xs text-gray-500">
              <Clock className="h-3 w-3" />
              Last updated {formatDistanceToNow(new Date(device.updatedAt), { addSuffix: true })}
            </div>
          )}
        </div>
      </ModalBody>
      <ModalFooter>
        <Button variant="outline" onClick={onClose}>
          Close
        </Button>
      </ModalFooter>
    </Modal>
  );
}
