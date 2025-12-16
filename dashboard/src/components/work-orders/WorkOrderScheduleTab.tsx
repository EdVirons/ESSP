import { Plus, Calendar } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { formatDate } from '@/lib/utils';
import type { WorkOrderSchedule } from '@/types';

interface WorkOrderScheduleTabProps {
  schedules: WorkOrderSchedule[];
}

export function WorkOrderScheduleTab({ schedules }: WorkOrderScheduleTabProps) {
  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h3 className="font-medium text-gray-900">Schedule</h3>
        <Button size="sm">
          <Plus className="h-4 w-4" />
          Add Schedule
        </Button>
      </div>
      {schedules.length === 0 ? (
        <div className="text-center py-8 text-gray-500">
          <Calendar className="h-12 w-12 mx-auto mb-2 text-gray-300" />
          <p>No schedules set</p>
        </div>
      ) : (
        <div className="space-y-2">
          {schedules.map((schedule) => (
            <div key={schedule.id} className="p-3 bg-gray-50 rounded-lg">
              <div className="flex items-center gap-2 mb-1">
                <Calendar className="h-4 w-4 text-gray-400" />
                <span className="font-medium text-gray-900">
                  {schedule.scheduledStart
                    ? formatDate(schedule.scheduledStart)
                    : 'Not scheduled'}
                </span>
              </div>
              {schedule.notes && (
                <p className="text-sm text-gray-500">{schedule.notes}</p>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
