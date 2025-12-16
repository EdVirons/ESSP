import * as React from 'react';
import {
  RefreshCw,
  School,
  Laptop,
  Package,
  CheckCircle2,
  Clock,
  Database,
  ArrowRight,
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import {
  useSyncStatus,
  useSyncSchools,
  useSyncDevices,
  useSyncParts,
} from '@/api/ssot';
import { formatDate } from '@/lib/utils';

interface SyncCardProps {
  title: string;
  icon: React.ReactNode;
  count: number;
  lastSyncAt: string;
  onSync: () => void;
  isSyncing: boolean;
  color: string;
}

function SyncCard({ title, icon, count, lastSyncAt, onSync, isSyncing, color }: SyncCardProps) {
  const hasSync = lastSyncAt && lastSyncAt !== '0001-01-01T00:00:00Z';

  return (
    <Card className="overflow-hidden">
      <div className={`h-2 ${color}`} />
      <CardContent className="pt-6">
        <div className="flex items-start justify-between">
          <div className="flex items-center gap-4">
            <div className={`h-14 w-14 rounded-xl ${color.replace('bg-', 'bg-opacity-20 bg-')} flex items-center justify-center`}>
              {icon}
            </div>
            <div>
              <h3 className="font-semibold text-gray-900 text-lg">{title}</h3>
              <div className="flex items-center gap-2 mt-1">
                <Database className="h-4 w-4 text-gray-400" />
                <span className="text-2xl font-bold text-gray-900">{count.toLocaleString()}</span>
                <span className="text-gray-500">records</span>
              </div>
            </div>
          </div>
          <Button
            variant="outline"
            size="sm"
            onClick={onSync}
            disabled={isSyncing}
            className="shrink-0"
          >
            <RefreshCw className={`h-4 w-4 mr-2 ${isSyncing ? 'animate-spin' : ''}`} />
            {isSyncing ? 'Syncing...' : 'Sync Now'}
          </Button>
        </div>

        <div className="mt-6 pt-4 border-t border-gray-100">
          <div className="flex items-center justify-between text-sm">
            <div className="flex items-center gap-2 text-gray-500">
              <Clock className="h-4 w-4" />
              <span>Last synced:</span>
            </div>
            <div className="flex items-center gap-2">
              {hasSync ? (
                <>
                  <CheckCircle2 className="h-4 w-4 text-green-500" />
                  <span className="text-gray-700">{formatDate(lastSyncAt)}</span>
                </>
              ) : (
                <span className="text-gray-400">Never synced</span>
              )}
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}

export function SSOTSync() {
  const { data: syncStatus, isLoading, refetch } = useSyncStatus();
  const syncSchools = useSyncSchools();
  const syncDevices = useSyncDevices();
  const syncParts = useSyncParts();

  const handleSyncAll = async () => {
    await Promise.all([
      syncSchools.mutateAsync(),
      syncDevices.mutateAsync(),
      syncParts.mutateAsync(),
    ]);
    refetch();
  };

  const isSyncingAll = syncSchools.isPending || syncDevices.isPending || syncParts.isPending;

  const totalRecords =
    (syncStatus?.schools?.count || 0) +
    (syncStatus?.devices?.count || 0) +
    (syncStatus?.parts?.count || 0);

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">SSOT Sync Control</h1>
          <p className="text-gray-500">Manage synchronization of Single Source of Truth data</p>
        </div>
        <Button onClick={handleSyncAll} disabled={isSyncingAll}>
          <RefreshCw className={`mr-2 h-4 w-4 ${isSyncingAll ? 'animate-spin' : ''}`} />
          {isSyncingAll ? 'Syncing All...' : 'Sync All'}
        </Button>
      </div>

      {/* Overview Stats */}
      <Card>
        <CardContent className="py-6">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-6">
              <div className="h-16 w-16 rounded-2xl bg-gradient-to-br from-blue-500 to-purple-600 flex items-center justify-center">
                <Database className="h-8 w-8 text-white" />
              </div>
              <div>
                <p className="text-sm text-gray-500">Total SSOT Records</p>
                <p className="text-4xl font-bold text-gray-900">{totalRecords.toLocaleString()}</p>
              </div>
            </div>
            <div className="hidden md:flex items-center gap-8">
              <div className="text-center">
                <div className="flex items-center gap-2 text-blue-600">
                  <School className="h-5 w-5" />
                  <span className="font-semibold">{syncStatus?.schools?.count || 0}</span>
                </div>
                <p className="text-xs text-gray-500 mt-1">Schools</p>
              </div>
              <ArrowRight className="h-4 w-4 text-gray-300" />
              <div className="text-center">
                <div className="flex items-center gap-2 text-indigo-600">
                  <Laptop className="h-5 w-5" />
                  <span className="font-semibold">{syncStatus?.devices?.count || 0}</span>
                </div>
                <p className="text-xs text-gray-500 mt-1">Devices</p>
              </div>
              <ArrowRight className="h-4 w-4 text-gray-300" />
              <div className="text-center">
                <div className="flex items-center gap-2 text-amber-600">
                  <Package className="h-5 w-5" />
                  <span className="font-semibold">{syncStatus?.parts?.count || 0}</span>
                </div>
                <p className="text-xs text-gray-500 mt-1">Parts</p>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Sync Cards */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <SyncCard
          title="Schools"
          icon={<School className="h-7 w-7 text-blue-600" />}
          count={syncStatus?.schools?.count || 0}
          lastSyncAt={syncStatus?.schools?.lastSyncAt || ''}
          onSync={() => syncSchools.mutate(undefined, { onSuccess: () => refetch() })}
          isSyncing={syncSchools.isPending}
          color="bg-blue-500"
        />
        <SyncCard
          title="Devices"
          icon={<Laptop className="h-7 w-7 text-indigo-600" />}
          count={syncStatus?.devices?.count || 0}
          lastSyncAt={syncStatus?.devices?.lastSyncAt || ''}
          onSync={() => syncDevices.mutate(undefined, { onSuccess: () => refetch() })}
          isSyncing={syncDevices.isPending}
          color="bg-indigo-500"
        />
        <SyncCard
          title="Parts"
          icon={<Package className="h-7 w-7 text-amber-600" />}
          count={syncStatus?.parts?.count || 0}
          lastSyncAt={syncStatus?.parts?.lastSyncAt || ''}
          onSync={() => syncParts.mutate(undefined, { onSuccess: () => refetch() })}
          isSyncing={syncParts.isPending}
          color="bg-amber-500"
        />
      </div>

      {/* Information */}
      <Card>
        <CardContent className="py-6">
          <h3 className="font-semibold text-gray-900 mb-4">About SSOT Sync</h3>
          <div className="space-y-3 text-sm text-gray-600">
            <p>
              <strong>Single Source of Truth (SSOT)</strong> data is synchronized from external
              master data services to ensure the IMS has up-to-date information about schools,
              devices, and parts.
            </p>
            <p>
              Syncing pulls data from the following SSOT services:
            </p>
            <ul className="list-disc list-inside space-y-1 ml-4">
              <li><strong>Schools</strong> - School registry with location information (ssot-school:8081)</li>
              <li><strong>Devices</strong> - Device inventory with assignments and status (ssot-devices:8082)</li>
              <li><strong>Parts</strong> - Parts catalog with PUK codes and categories (ssot-parts:8083)</li>
            </ul>
            <p className="text-gray-500">
              Sync runs incrementally using cursors to minimize data transfer. Only changed records
              since the last sync are fetched.
            </p>
          </div>
        </CardContent>
      </Card>

      {/* Loading State */}
      {isLoading && (
        <div className="flex items-center justify-center py-12">
          <RefreshCw className="h-8 w-8 animate-spin text-gray-400" />
        </div>
      )}
    </div>
  );
}
