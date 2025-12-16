import { Card, CardContent } from '@/components/ui/card';
import { Layers, PlayCircle, PauseCircle, CheckCircle2 } from 'lucide-react';
import { projectTypeConfigs } from '@/lib/projectTypes';
import type { SchoolServiceProject, ProjectType } from '@/types';

interface ProjectsStatsProps {
  projects: SchoolServiceProject[];
  projectType?: ProjectType;
}

export function ProjectsStats({ projects, projectType }: ProjectsStatsProps) {
  const config = projectType ? projectTypeConfigs[projectType] : null;
  const totalCount = projects.length;
  const activeCount = projects.filter((p) => p.status === 'active').length;
  const pausedCount = projects.filter((p) => p.status === 'paused').length;
  const completedCount = projects.filter((p) => p.status === 'completed').length;

  return (
    <div className="grid gap-4 md:grid-cols-4">
      <Card>
        <CardContent className="p-4">
          <div className="flex items-center gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-full bg-purple-50">
              <Layers className="h-5 w-5 text-purple-600" />
            </div>
            <div>
              <div className="text-2xl font-bold">{totalCount}</div>
              <div className="text-sm text-gray-500">
                {config ? config.label : 'Total Projects'}
              </div>
            </div>
          </div>
        </CardContent>
      </Card>
      <Card>
        <CardContent className="p-4">
          <div className="flex items-center gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-full bg-green-50">
              <PlayCircle className="h-5 w-5 text-green-600" />
            </div>
            <div>
              <div className="text-2xl font-bold">{activeCount}</div>
              <div className="text-sm text-gray-500">Active</div>
            </div>
          </div>
        </CardContent>
      </Card>
      <Card>
        <CardContent className="p-4">
          <div className="flex items-center gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-full bg-yellow-50">
              <PauseCircle className="h-5 w-5 text-yellow-600" />
            </div>
            <div>
              <div className="text-2xl font-bold">{pausedCount}</div>
              <div className="text-sm text-gray-500">Paused</div>
            </div>
          </div>
        </CardContent>
      </Card>
      <Card>
        <CardContent className="p-4">
          <div className="flex items-center gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-full bg-gray-100">
              <CheckCircle2 className="h-5 w-5 text-gray-600" />
            </div>
            <div>
              <div className="text-2xl font-bold">{completedCount}</div>
              <div className="text-sm text-gray-500">Completed</div>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
