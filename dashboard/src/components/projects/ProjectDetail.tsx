import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Sheet, SheetHeader, SheetBody, SheetFooter } from '@/components/ui/sheet';
import { Tabs, TabsList, TabsTrigger, TabsContent } from '@/components/ui/tabs';
import { formatDate, cn } from '@/lib/utils';
import { phaseTypeLabels, projectTypeConfigs } from '@/lib/projectTypes';
import { PhaseTimeline } from './PhaseTimeline';
import { SurveysList } from './SurveysList';
import { TeamPanel } from './TeamPanel';
import { ActivityFeed } from './ActivityFeed';
import { useProjectTeam, useProjectActivities } from '@/api/projects';
import type { SchoolServiceProject, ProjectStatus, ServicePhase, SiteSurvey } from '@/types';

const statusColors: Record<ProjectStatus, string> = {
  active: 'bg-green-100 text-green-800',
  paused: 'bg-yellow-100 text-yellow-800',
  completed: 'bg-gray-100 text-gray-800',
};

interface ProjectDetailProps {
  project: SchoolServiceProject | null;
  open: boolean;
  onClose: () => void;
  detailTab: string;
  onDetailTabChange: (tab: string) => void;
  phases: ServicePhase[];
  surveys: SiteSurvey[];
}

export function ProjectDetail({
  project,
  open,
  onClose,
  detailTab,
  onDetailTabChange,
  phases,
  surveys,
}: ProjectDetailProps) {
  // Fetch team and activity counts for tab badges
  const { data: teamData } = useProjectTeam(project?.id || '');
  const { data: activitiesData } = useProjectActivities(project?.id || '', { limit: 1 });

  const teamCount = teamData?.total || 0;
  const hasActivities = (activitiesData?.items?.length || 0) > 0;

  if (!project) return null;

  return (
    <Sheet open={open} onClose={onClose} side="right">
      <SheetHeader onClose={onClose}>Project Details</SheetHeader>
      <SheetBody className="p-0">
        <div className="h-full flex flex-col">
          {/* Header Info */}
          <div className="p-6 border-b border-gray-200">
            <div className="flex items-center gap-2 mb-2">
              <Badge className={cn('capitalize', statusColors[project.status])}>
                {project.status}
              </Badge>
              <Badge variant="outline">
                {projectTypeConfigs[project.projectType]?.label || project.projectType}
              </Badge>
            </div>
            <h2 className="text-lg font-semibold text-gray-900">
              {project.schoolId}
            </h2>
            <p className="text-sm text-gray-500">
              Current Phase: {phaseTypeLabels[project.currentPhase]}
            </p>
          </div>

          {/* Tabs */}
          <Tabs
            value={detailTab}
            onValueChange={onDetailTabChange}
            className="flex-1 flex flex-col"
          >
            <div className="border-b border-gray-200 px-6">
              <TabsList className="bg-transparent -mb-px">
                <TabsTrigger value="phases">Phases ({phases.length})</TabsTrigger>
                <TabsTrigger value="team">
                  Team {teamCount > 0 && `(${teamCount})`}
                </TabsTrigger>
                <TabsTrigger value="activity">
                  Activity {hasActivities && '●'}
                </TabsTrigger>
                <TabsTrigger value="surveys">Surveys ({surveys.length})</TabsTrigger>
                <TabsTrigger value="details">Details</TabsTrigger>
              </TabsList>
            </div>

            <div className="flex-1 overflow-auto">
              <TabsContent value="phases" className="p-6 m-0">
                <PhaseTimeline project={project} phases={phases} />
              </TabsContent>

              <TabsContent value="team" className="p-6 m-0">
                <TeamPanel projectId={project.id} canEdit={true} />
              </TabsContent>

              <TabsContent value="activity" className="p-6 m-0">
                <ActivityFeed projectId={project.id} canEdit={true} />
              </TabsContent>

              <TabsContent value="surveys" className="p-6 m-0">
                <SurveysList surveys={surveys} />
              </TabsContent>

              <TabsContent value="details" className="p-6 m-0">
                <div className="space-y-6">
                  <div className="grid grid-cols-2 gap-4">
                    <div>
                      <h3 className="text-sm font-medium text-gray-500 mb-1">
                        Start Date
                      </h3>
                      <p className="text-gray-900">
                        {project.startDate
                          ? formatDate(project.startDate)
                          : 'Not set'}
                      </p>
                    </div>
                    <div>
                      <h3 className="text-sm font-medium text-gray-500 mb-1">
                        Go-Live Date
                      </h3>
                      <p className="text-gray-900">
                        {project.goLiveDate
                          ? formatDate(project.goLiveDate)
                          : 'Not set'}
                      </p>
                    </div>
                  </div>

                  {project.notes && (
                    <div>
                      <h3 className="text-sm font-medium text-gray-500 mb-1">
                        Notes
                      </h3>
                      <p className="text-gray-900">{project.notes}</p>
                    </div>
                  )}

                  <div>
                    <h3 className="text-sm font-medium text-gray-500 mb-2">
                      Timeline
                    </h3>
                    <div className="space-y-2 text-sm">
                      <div className="flex justify-between">
                        <span className="text-gray-500">Created</span>
                        <span className="text-gray-900">
                          {formatDate(project.createdAt)}
                        </span>
                      </div>
                      <div className="flex justify-between">
                        <span className="text-gray-500">Last Updated</span>
                        <span className="text-gray-900">
                          {formatDate(project.updatedAt)}
                        </span>
                      </div>
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
      </SheetFooter>
    </Sheet>
  );
}
