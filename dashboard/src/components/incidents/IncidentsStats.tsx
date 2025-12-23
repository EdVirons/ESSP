import { Card, CardContent } from '@/components/ui/card';
import { AlertTriangle, Clock, AlertCircle } from 'lucide-react';
import type { Incident } from '@/types';

interface IncidentsStatsProps {
  incidents: Incident[];
}

export function IncidentsStats({ incidents }: IncidentsStatsProps) {
  const totalCount = incidents.length;
  const openCount = incidents.filter(
    (i) => ['new', 'acknowledged', 'in_progress'].includes(i.status)
  ).length;
  const slaBreachedCount = incidents.filter((i) => i.slaBreached).length;
  const highPriorityCount = incidents.filter(
    (i) => i.severity === 'critical' || i.severity === 'high'
  ).length;

  return (
    <div className="grid grid-cols-2 gap-3 md:gap-4 lg:grid-cols-4">
      <Card>
        <CardContent className="p-4">
          <div className="flex items-center gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-full bg-blue-50">
              <AlertTriangle className="h-5 w-5 text-blue-600" />
            </div>
            <div>
              <div className="text-2xl font-bold">{totalCount}</div>
              <div className="text-sm text-gray-500">Total Incidents</div>
            </div>
          </div>
        </CardContent>
      </Card>
      <Card>
        <CardContent className="p-4">
          <div className="flex items-center gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-full bg-yellow-50">
              <Clock className="h-5 w-5 text-yellow-600" />
            </div>
            <div>
              <div className="text-2xl font-bold">{openCount}</div>
              <div className="text-sm text-gray-500">Open</div>
            </div>
          </div>
        </CardContent>
      </Card>
      <Card>
        <CardContent className="p-4">
          <div className="flex items-center gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-full bg-red-50">
              <AlertCircle className="h-5 w-5 text-red-600" />
            </div>
            <div>
              <div className="text-2xl font-bold">{slaBreachedCount}</div>
              <div className="text-sm text-gray-500">SLA Breached</div>
            </div>
          </div>
        </CardContent>
      </Card>
      <Card>
        <CardContent className="p-4">
          <div className="flex items-center gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-full bg-orange-50">
              <AlertTriangle className="h-5 w-5 text-orange-600" />
            </div>
            <div>
              <div className="text-2xl font-bold">{highPriorityCount}</div>
              <div className="text-sm text-gray-500">High Priority</div>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
