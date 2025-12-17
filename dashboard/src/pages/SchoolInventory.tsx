import * as React from 'react';
import { RefreshCw, Laptop, MapPin, Users, Plus } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import {
  useSchoolInventory,
  useGroups,
  useRegisterDevice,
  useCreateLocation,
  useUpdateLocation,
  useAssignDevice,
  useCreateGroup,
  useAddGroupMembers,
  useRemoveGroupMember,
} from '@/api/inventory';
import { useSchools } from '@/api/ssot';
import {
  InventoryStats,
  DeviceList,
  LocationList,
  AddDeviceModal,
  LocationModal,
  DeviceDetailModal,
  AssignDeviceModal,
  GroupManageModal,
} from '@/components/inventory';
import { cn } from '@/lib/utils';
import type {
  InventoryDevice,
  Location,
  RegisterDeviceRequest,
  CreateLocationRequest,
  UpdateLocationRequest,
  AssignDeviceRequest,
  CreateGroupRequest,
} from '@/types';
import { toast } from 'sonner';

export function SchoolInventory() {
  // Get current school from localStorage or default
  const [selectedSchoolId, setSelectedSchoolId] = React.useState(() => {
    return localStorage.getItem('school_id') || '';
  });

  // Active tab state
  const [activeTab, setActiveTab] = React.useState('devices');

  // Modal states
  const [addDeviceOpen, setAddDeviceOpen] = React.useState(false);
  const [locationModalOpen, setLocationModalOpen] = React.useState(false);
  const [deviceDetailOpen, setDeviceDetailOpen] = React.useState(false);
  const [assignDeviceOpen, setAssignDeviceOpen] = React.useState(false);
  const [groupManageOpen, setGroupManageOpen] = React.useState(false);
  const [assignMode, setAssignMode] = React.useState<'location' | 'user'>('location');

  // Selected items for modals
  const [selectedDevice, setSelectedDevice] = React.useState<InventoryDevice | null>(null);
  const [selectedLocation, setSelectedLocation] = React.useState<Location | null>(null);

  // Fetch schools for selector
  const { data: schoolsData } = useSchools({ limit: 100 });
  const schools = schoolsData?.items || [];

  // Fetch inventory for selected school
  const {
    data: inventoryData,
    isLoading,
    refetch,
    isRefetching,
  } = useSchoolInventory(selectedSchoolId);

  // Fetch groups for selected school
  const { data: groupsData, refetch: refetchGroups } = useGroups(selectedSchoolId);

  // Mutations
  const registerDevice = useRegisterDevice();
  const createLocation = useCreateLocation();
  const updateLocation = useUpdateLocation();
  const assignDevice = useAssignDevice();
  const createGroup = useCreateGroup();
  const addGroupMembers = useAddGroupMembers();
  const removeGroupMember = useRemoveGroupMember();

  // Handle device click
  const handleDeviceClick = (device: InventoryDevice) => {
    setSelectedDevice(device);
    setDeviceDetailOpen(true);
  };

  // Handle location click
  const handleLocationClick = (location: Location) => {
    setSelectedLocation(location);
    setLocationModalOpen(true);
  };

  // Handle add location
  const handleAddLocation = () => {
    setSelectedLocation(null);
    setLocationModalOpen(true);
  };

  // Handle school change
  const handleSchoolChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    setSelectedSchoolId(e.target.value);
  };

  // Handle register device
  const handleRegisterDevice = (data: RegisterDeviceRequest) => {
    registerDevice.mutate(
      { schoolId: selectedSchoolId, data },
      {
        onSuccess: () => {
          toast.success('Device registered successfully');
          setAddDeviceOpen(false);
        },
        onError: (error) => {
          toast.error('Failed to register device: ' + (error as Error).message);
        },
      }
    );
  };

  // Handle create/update location
  const handleLocationSubmit = (data: CreateLocationRequest | UpdateLocationRequest) => {
    if (selectedLocation) {
      // Update existing
      updateLocation.mutate(
        { schoolId: selectedSchoolId, id: selectedLocation.id, data: data as UpdateLocationRequest },
        {
          onSuccess: () => {
            toast.success('Location updated successfully');
            setLocationModalOpen(false);
            setSelectedLocation(null);
          },
          onError: (error) => {
            toast.error('Failed to update location: ' + (error as Error).message);
          },
        }
      );
    } else {
      // Create new
      createLocation.mutate(
        { schoolId: selectedSchoolId, data: data as CreateLocationRequest },
        {
          onSuccess: () => {
            toast.success('Location created successfully');
            setLocationModalOpen(false);
          },
          onError: (error) => {
            toast.error('Failed to create location: ' + (error as Error).message);
          },
        }
      );
    }
  };

  // Handle assign device
  const handleAssignDevice = (deviceId: string, data: AssignDeviceRequest) => {
    assignDevice.mutate(
      { deviceId, data },
      {
        onSuccess: () => {
          toast.success('Device assignment updated');
          setAssignDeviceOpen(false);
          setDeviceDetailOpen(false);
          setSelectedDevice(null);
        },
        onError: (error) => {
          toast.error('Failed to assign device: ' + (error as Error).message);
        },
      }
    );
  };

  // Handle device detail actions
  const handleAssignLocation = (device: InventoryDevice) => {
    setSelectedDevice(device);
    setAssignMode('location');
    setAssignDeviceOpen(true);
  };

  const handleAssignUser = (device: InventoryDevice) => {
    setSelectedDevice(device);
    setAssignMode('user');
    setAssignDeviceOpen(true);
  };

  const handleManageGroups = (device: InventoryDevice) => {
    setSelectedDevice(device);
    setGroupManageOpen(true);
  };

  // Handle group operations
  const handleCreateGroup = (data: CreateGroupRequest) => {
    createGroup.mutate(
      { ...data, schoolId: selectedSchoolId },
      {
        onSuccess: () => {
          toast.success('Group created successfully');
          refetchGroups();
        },
        onError: (error) => {
          toast.error('Failed to create group: ' + (error as Error).message);
        },
      }
    );
  };

  const handleAddToGroup = (groupId: string, deviceId: string) => {
    addGroupMembers.mutate(
      { groupId, data: { deviceIds: [deviceId] } },
      {
        onSuccess: () => {
          toast.success('Device added to group');
          refetchGroups();
        },
        onError: (error) => {
          toast.error('Failed to add device to group: ' + (error as Error).message);
        },
      }
    );
  };

  const handleRemoveFromGroup = (groupId: string, deviceId: string) => {
    removeGroupMember.mutate(
      { groupId, deviceId },
      {
        onSuccess: () => {
          toast.success('Device removed from group');
          refetchGroups();
        },
        onError: (error) => {
          toast.error('Failed to remove device from group: ' + (error as Error).message);
        },
      }
    );
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-semibold text-gray-900">School Inventory</h1>
          <p className="text-sm text-gray-500 mt-1">
            View and manage device inventory for your school
          </p>
        </div>
        <div className="flex items-center gap-3">
          {/* School selector */}
          <select
            value={selectedSchoolId}
            onChange={handleSchoolChange}
            className="rounded-md border border-gray-300 bg-white px-3 py-2 text-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
          >
            <option value="">Select a school</option>
            {schools.map((school) => (
              <option key={school.schoolId} value={school.schoolId}>
                {school.name}
              </option>
            ))}
          </select>

          {selectedSchoolId && (
            <Button onClick={() => setAddDeviceOpen(true)} className="gap-2">
              <Plus className="h-4 w-4" />
              Add Device
            </Button>
          )}

          <Button
            variant="outline"
            onClick={() => refetch()}
            disabled={isRefetching || !selectedSchoolId}
            className="gap-2"
          >
            <RefreshCw className={cn('h-4 w-4', isRefetching && 'animate-spin')} />
            Refresh
          </Button>
        </div>
      </div>

      {/* No school selected */}
      {!selectedSchoolId && (
        <Card>
          <CardContent className="py-12 text-center">
            <Laptop className="h-12 w-12 mx-auto text-gray-300 mb-4" />
            <h3 className="text-lg font-medium text-gray-900 mb-2">Select a School</h3>
            <p className="text-gray-500">
              Choose a school from the dropdown above to view its device inventory.
            </p>
          </CardContent>
        </Card>
      )}

      {/* Loading */}
      {selectedSchoolId && isLoading && (
        <div className="space-y-6">
          <InventoryStats summary={{ totalDevices: 0, byStatus: {}, byLocation: { assigned: 0, unassigned: 0 } }} loading />
          <Card>
            <CardContent className="py-12 text-center">
              <RefreshCw className="h-8 w-8 mx-auto text-blue-500 animate-spin mb-4" />
              <p className="text-gray-500">Loading inventory...</p>
            </CardContent>
          </Card>
        </div>
      )}

      {/* Inventory data */}
      {selectedSchoolId && !isLoading && inventoryData && (
        <>
          {/* School name */}
          <div className="bg-blue-50 border border-blue-100 rounded-lg p-4">
            <h2 className="text-lg font-semibold text-blue-900">
              {inventoryData.school.name}
            </h2>
            <p className="text-sm text-blue-700">
              {inventoryData.summary.totalDevices} devices total
            </p>
          </div>

          {/* Stats */}
          <InventoryStats summary={inventoryData.summary} />

          {/* Tabs */}
          <Tabs value={activeTab} onValueChange={setActiveTab} className="space-y-4">
            <TabsList>
              <TabsTrigger value="devices" className="gap-2">
                <Laptop className="h-4 w-4" />
                Devices ({inventoryData.devices.length})
              </TabsTrigger>
              <TabsTrigger value="locations" className="gap-2">
                <MapPin className="h-4 w-4" />
                Locations ({inventoryData.locations.length})
              </TabsTrigger>
              <TabsTrigger value="groups" className="gap-2">
                <Users className="h-4 w-4" />
                Groups ({groupsData?.items?.length || 0})
              </TabsTrigger>
            </TabsList>

            {/* Devices Tab */}
            <TabsContent value="devices">
              <Card>
                <CardHeader>
                  <CardTitle>Device Inventory</CardTitle>
                </CardHeader>
                <CardContent>
                  <DeviceList
                    devices={inventoryData.devices}
                    onDeviceClick={handleDeviceClick}
                  />
                </CardContent>
              </Card>
            </TabsContent>

            {/* Locations Tab */}
            <TabsContent value="locations">
              <Card>
                <CardHeader>
                  <CardTitle>Locations</CardTitle>
                </CardHeader>
                <CardContent>
                  <LocationList
                    locations={inventoryData.locations}
                    onLocationClick={handleLocationClick}
                    onAddLocation={handleAddLocation}
                  />
                </CardContent>
              </Card>
            </TabsContent>

            {/* Groups Tab */}
            <TabsContent value="groups">
              <Card>
                <CardHeader className="flex flex-row items-center justify-between">
                  <CardTitle>Device Groups</CardTitle>
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => {
                      setSelectedDevice(null);
                      setGroupManageOpen(true);
                    }}
                    className="gap-2"
                  >
                    <Plus className="h-4 w-4" />
                    Create Group
                  </Button>
                </CardHeader>
                <CardContent>
                  {groupsData?.items && groupsData.items.length > 0 ? (
                    <div className="space-y-2">
                      {groupsData.items.map((group) => (
                        <div
                          key={group.id}
                          className="flex items-center justify-between p-3 bg-gray-50 rounded-lg"
                        >
                          <div>
                            <p className="font-medium text-gray-900">{group.name}</p>
                            {group.description && (
                              <p className="text-sm text-gray-500">{group.description}</p>
                            )}
                          </div>
                          <div className="flex items-center gap-2">
                            <span className="text-sm text-gray-500">
                              {group.groupType}
                            </span>
                            {group.memberCount !== undefined && (
                              <span className="text-sm font-medium text-blue-600">
                                {group.memberCount} devices
                              </span>
                            )}
                          </div>
                        </div>
                      ))}
                    </div>
                  ) : (
                    <div className="text-center py-8">
                      <Users className="h-12 w-12 mx-auto text-gray-300 mb-3" />
                      <p className="text-gray-500">No device groups defined</p>
                      <Button
                        variant="outline"
                        className="mt-4"
                        onClick={() => {
                          setSelectedDevice(null);
                          setGroupManageOpen(true);
                        }}
                      >
                        <Plus className="h-4 w-4 mr-2" />
                        Create First Group
                      </Button>
                    </div>
                  )}
                </CardContent>
              </Card>
            </TabsContent>
          </Tabs>
        </>
      )}

      {/* Modals */}
      <AddDeviceModal
        open={addDeviceOpen}
        onClose={() => setAddDeviceOpen(false)}
        onSubmit={handleRegisterDevice}
        isLoading={registerDevice.isPending}
        locations={inventoryData?.locations || []}
      />

      <LocationModal
        open={locationModalOpen}
        onClose={() => {
          setLocationModalOpen(false);
          setSelectedLocation(null);
        }}
        onSubmit={handleLocationSubmit}
        isLoading={createLocation.isPending || updateLocation.isPending}
        location={selectedLocation}
        locations={inventoryData?.locations || []}
      />

      <DeviceDetailModal
        open={deviceDetailOpen}
        onClose={() => {
          setDeviceDetailOpen(false);
          setSelectedDevice(null);
        }}
        device={selectedDevice}
        groups={groupsData?.items || []}
        onAssignLocation={handleAssignLocation}
        onAssignUser={handleAssignUser}
        onManageGroups={handleManageGroups}
      />

      <AssignDeviceModal
        open={assignDeviceOpen}
        onClose={() => {
          setAssignDeviceOpen(false);
          setSelectedDevice(null);
        }}
        onSubmit={handleAssignDevice}
        isLoading={assignDevice.isPending}
        device={selectedDevice}
        locations={inventoryData?.locations || []}
        mode={assignMode}
      />

      <GroupManageModal
        open={groupManageOpen}
        onClose={() => {
          setGroupManageOpen(false);
          setSelectedDevice(null);
        }}
        device={selectedDevice}
        groups={groupsData?.items || []}
        deviceGroupIds={selectedDevice?.groups || []}
        onAddToGroup={handleAddToGroup}
        onRemoveFromGroup={handleRemoveFromGroup}
        onCreateGroup={handleCreateGroup}
        isLoading={createGroup.isPending || addGroupMembers.isPending || removeGroupMember.isPending}
      />
    </div>
  );
}
