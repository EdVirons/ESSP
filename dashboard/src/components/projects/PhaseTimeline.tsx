import { Plus, CheckCircle2 } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { formatDate, formatStatus, cn } from '@/lib/utils';
import { getPhasesForType, phaseTypeLabels, phaseStatusColors } from '@/lib/projectTypes';
import type { SchoolServiceProject, ServicePhase } from '@/types';

interface PhaseTimelineProps {
  project: SchoolServiceProject;
  phases: ServicePhase[];
}

export function PhaseTimeline({ project, phases }: PhaseTimelineProps) {
  // Get phases specific to this project's type
  const phaseTypes = getPhasesForType(project.projectType);

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h3 className="font-medium text-gray-900">Project Phases</h3>
        <Button size="sm">
          <Plus className="h-4 w-4" />
          Add Phase
        </Button>
      </div>

      <div className="space-y-3">
        {phaseTypes.map((phaseType, index) => {
          const phase = phases.find((p) => p.phaseType === phaseType);
          const isCurrentPhase = project.currentPhase === phaseType;

          return (
            <div
              key={phaseType}
              className={cn(
                'relative flex items-start gap-4 p-3 rounded-lg border',
                isCurrentPhase
                  ? 'border-blue-200 bg-blue-50'
                  : 'border-gray-200 bg-white'
              )}
            >
              {/* Step indicator */}
              <div
                className={cn(
                  'flex h-8 w-8 shrink-0 items-center justify-center rounded-full text-sm font-medium',
                  phase?.status === 'done'
                    ? 'bg-green-100 text-green-800'
                    : phase?.status === 'in_progress'
                    ? 'bg-blue-100 text-blue-800'
                    : 'bg-gray-100 text-gray-600'
                )}
              >
                {phase?.status === 'done' ? (
                  <CheckCircle2 className="h-5 w-5" />
                ) : (
                  index + 1
                )}
              </div>

              <div className="flex-1 min-w-0">
                <div className="flex items-center gap-2">
                  <span className="font-medium text-gray-900">
                    {phaseTypeLabels[phaseType]}
                  </span>
                  {phase && (
                    <Badge className={cn('text-xs', phaseStatusColors[phase.status])}>
                      {formatStatus(phase.status)}
                    </Badge>
                  )}
                </div>
                {phase && (
                  <div className="text-sm text-gray-500 mt-1">
                    {phase.startDate && (
                      <span>Started: {formatDate(phase.startDate)}</span>
                    )}
                    {phase.endDate && (
                      <span className="ml-4">Ended: {formatDate(phase.endDate)}</span>
                    )}
                  </div>
                )}
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
}
