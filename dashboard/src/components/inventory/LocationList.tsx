import * as React from 'react';
import { Building2, Layers, DoorOpen, Monitor, Warehouse, Briefcase, ChevronRight, Plus } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import type { Location, LocationType } from '@/types';

const LOCATION_ICONS: Record<LocationType, React.ElementType> = {
  block: Building2,
  floor: Layers,
  room: DoorOpen,
  lab: Monitor,
  storage: Warehouse,
  office: Briefcase,
};

interface LocationListProps {
  locations: Location[];
  loading?: boolean;
  onLocationClick?: (location: Location) => void;
  onAddLocation?: () => void;
}

export function LocationList({ locations, loading, onLocationClick, onAddLocation }: LocationListProps) {
  if (loading) {
    return (
      <div className="space-y-2">
        {[1, 2, 3, 4].map((i) => (
          <div key={i} className="h-12 animate-pulse bg-gray-100 rounded" />
        ))}
      </div>
    );
  }

  if (locations.length === 0) {
    return (
      <div className="text-center py-8">
        <Building2 className="h-12 w-12 mx-auto text-gray-300 mb-3" />
        <p className="text-gray-500 mb-4">No locations defined yet</p>
        {onAddLocation && (
          <Button variant="outline" onClick={onAddLocation} className="gap-2">
            <Plus className="h-4 w-4" />
            Add Location
          </Button>
        )}
      </div>
    );
  }

  // Group by type
  const byType = locations.reduce((acc, loc) => {
    if (!acc[loc.locationType]) acc[loc.locationType] = [];
    acc[loc.locationType].push(loc);
    return acc;
  }, {} as Record<string, Location[]>);

  return (
    <div className="space-y-4">
      {onAddLocation && (
        <div className="flex justify-end">
          <Button variant="outline" size="sm" onClick={onAddLocation} className="gap-2">
            <Plus className="h-4 w-4" />
            Add Location
          </Button>
        </div>
      )}

      {Object.entries(byType).map(([type, locs]) => {
        const Icon = LOCATION_ICONS[type as LocationType] || Building2;
        return (
          <div key={type} className="space-y-2">
            <h4 className="text-sm font-medium text-gray-500 uppercase tracking-wider flex items-center gap-2">
              <Icon className="h-4 w-4" />
              {type}s ({locs.length})
            </h4>
            <div className="space-y-1">
              {locs.map((location) => (
                <div
                  key={location.id}
                  className="flex items-center justify-between p-3 bg-gray-50 rounded-lg hover:bg-gray-100 cursor-pointer transition-colors"
                  onClick={() => onLocationClick?.(location)}
                >
                  <div className="flex items-center gap-3">
                    <Icon className="h-5 w-5 text-gray-400" />
                    <div>
                      <p className="font-medium text-gray-900">{location.name}</p>
                      {location.code && (
                        <p className="text-sm text-gray-500">Code: {location.code}</p>
                      )}
                    </div>
                  </div>
                  <div className="flex items-center gap-3">
                    {location.deviceCount !== undefined && location.deviceCount > 0 && (
                      <Badge variant="secondary" className="bg-blue-50 text-blue-700">
                        {location.deviceCount} devices
                      </Badge>
                    )}
                    {location.capacity > 0 && (
                      <span className="text-sm text-gray-500">
                        Cap: {location.capacity}
                      </span>
                    )}
                    <ChevronRight className="h-4 w-4 text-gray-400" />
                  </div>
                </div>
              ))}
            </div>
          </div>
        );
      })}
    </div>
  );
}
