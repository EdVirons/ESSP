import * as React from 'react';
import { type ColumnDef } from '@tanstack/react-table';
import {
  Plus,
  Layers,
  Search,
  ChevronRight,
  RefreshCw,
  Headphones,
  Wrench,
  GraduationCap,
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Select } from '@/components/ui/select';
import { Tabs, TabsList, TabsTrigger, TabsContent } from '@/components/ui/tabs';
import { DataTable, SortableHeader, createSelectColumn } from '@/components/ui/data-table';
import {
  ProjectsStats,
  ProjectDetail,
  CreateProjectModal,
} from '@/components/projects';
import { useProjects, useCreateProject, usePhases, useSurveys, useProjectTypeCounts } from '@/api/projects';
import { formatDate, cn } from '@/lib/utils';
import {
  projectTypeConfigs,
  projectTypeOrder,
  projectTypeColors,
  phaseTypeLabels,
} from '@/lib/projectTypes';
import type { SchoolServiceProject, ProjectStatus, ProjectType } from '@/types';

const statusOptions = [
  { value: '', label: 'All Statuses' },
  { value: 'active', label: 'Active' },
  { value: 'paused', label: 'Paused' },
  { value: 'completed', label: 'Completed' },
];

const statusColors: Record<ProjectStatus, string> = {
  active: 'bg-green-100 text-green-800',
  paused: 'bg-yellow-100 text-yellow-800',
  completed: 'bg-gray-100 text-gray-800',
};

// Icon mapping for project types
const projectTypeIcons: Record<ProjectType, React.ElementType> = {
  full_installation: Layers,
  device_refresh: RefreshCw,
  support: Headphones,
  repair: Wrench,
  training: GraduationCap,
};

export function Projects() {
  // Active project type tab
  const [activeProjectType, setActiveProjectType] = React.useState<ProjectType>('full_installation');

  // Filters state
  const [filters, setFilters] = React.useState<{ status?: ProjectStatus; limit?: number }>({
    limit: 50,
  });
  const [searchQuery, setSearchQuery] = React.useState('');

  // Selected project for detail view
  const [selectedProject, setSelectedProject] = React.useState<SchoolServiceProject | null>(null);
  const [showDetail, setShowDetail] = React.useState(false);
  const [detailTab, setDetailTab] = React.useState('phases');

  // Create project modal
  const [showCreateModal, setShowCreateModal] = React.useState(false);
  const [createForm, setCreateForm] = React.useState({
    schoolId: '',
    projectType: 'full_installation' as ProjectType,
    startDate: null as Date | null,
    goLiveDate: null as Date | null,
    notes: '',
  });

  // API hooks - filtered by project type
  const { data, isLoading } = useProjects({
    ...filters,
    projectType: activeProjectType,
  });
  const { data: countsData } = useProjectTypeCounts();
  const createProject = useCreateProject();

  // Detail view data
  const { data: phasesData } = usePhases(selectedProject?.id || '');
  const { data: surveysData } = useSurveys(selectedProject?.id || '');

  // Table columns
  const columns: ColumnDef<SchoolServiceProject>[] = React.useMemo(
    () => [
      createSelectColumn<SchoolServiceProject>(),
      {
        accessorKey: 'schoolId',
        header: ({ column }) => <SortableHeader column={column}>School</SortableHeader>,
        cell: ({ row }) => {
          const colors = projectTypeColors[row.original.projectType];
          const Icon = projectTypeIcons[row.original.projectType] || Layers;
          return (
            <div className="flex items-center gap-3">
              <div className={cn('flex h-8 w-8 items-center justify-center rounded-full', colors?.bg || 'bg-purple-50')}>
                <Icon className={cn('h-4 w-4', colors?.text || 'text-purple-600')} />
              </div>
              <div className="min-w-0">
                <div className="font-medium text-gray-900 truncate max-w-[200px]">
                  {row.original.schoolId}
                </div>
                <div className="text-sm text-gray-500">
                  Current: {phaseTypeLabels[row.original.currentPhase]}
                </div>
              </div>
            </div>
          );
        },
      },
      {
        accessorKey: 'status',
        header: 'Status',
        cell: ({ row }) => (
          <Badge className={cn('capitalize', statusColors[row.original.status])}>
            {row.original.status}
          </Badge>
        ),
      },
      {
        accessorKey: 'currentPhase',
        header: 'Current Phase',
        cell: ({ row }) => (
          <div className="text-sm">
            <div className="font-medium">{phaseTypeLabels[row.original.currentPhase]}</div>
          </div>
        ),
      },
      {
        accessorKey: 'startDate',
        header: 'Start Date',
        cell: ({ row }) => (
          <div className="text-sm text-gray-500">
            {row.original.startDate ? formatDate(row.original.startDate) : '-'}
          </div>
        ),
      },
      {
        accessorKey: 'goLiveDate',
        header: 'Go-Live Date',
        cell: ({ row }) => (
          <div className="text-sm text-gray-500">
            {row.original.goLiveDate ? formatDate(row.original.goLiveDate) : '-'}
          </div>
        ),
      },
      {
        id: 'actions',
        cell: ({ row }) => (
          <Button
            variant="ghost"
            size="sm"
            onClick={(e) => {
              e.stopPropagation();
              setSelectedProject(row.original);
              setShowDetail(true);
              setDetailTab('phases');
            }}
          >
            <ChevronRight className="h-4 w-4" />
          </Button>
        ),
      },
    ],
    []
  );

  // Handle create project
  const handleCreateProject = async () => {
    try {
      await createProject.mutateAsync({
        schoolId: createForm.schoolId,
        projectType: createForm.projectType,
        startDate: createForm.startDate?.toISOString(),
        goLiveDate: createForm.goLiveDate?.toISOString(),
        notes: createForm.notes || undefined,
      });
      setShowCreateModal(false);
      setCreateForm({
        schoolId: '',
        projectType: activeProjectType,
        startDate: null,
        goLiveDate: null,
        notes: '',
      });
    } catch (err) {
      console.error('Failed to create project:', err);
    }
  };

  // When opening create modal, default to current tab's project type
  const handleOpenCreateModal = () => {
    setCreateForm((prev) => ({ ...prev, projectType: activeProjectType }));
    setShowCreateModal(true);
  };

  const projects = data?.items || [];
  const phases = phasesData?.items || [];
  const surveys = surveysData?.items || [];
  const counts: Record<ProjectType, number> = countsData || {} as Record<ProjectType, number>;

  return (
    <div className="space-y-6">
      {/* Page Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Projects</h1>
          <p className="text-sm text-gray-500">
            Manage school service projects and phases
          </p>
        </div>
        <Button onClick={handleOpenCreateModal}>
          <Plus className="h-4 w-4" />
          Create Project
        </Button>
      </div>

      {/* Project Type Tabs */}
      <Tabs value={activeProjectType} onValueChange={(v) => setActiveProjectType(v as ProjectType)}>
        <TabsList className="bg-white border border-gray-200 p-1 flex-wrap h-auto gap-1">
          {projectTypeOrder.map((type) => {
            const config = projectTypeConfigs[type];
            const colors = projectTypeColors[type];
            const Icon = projectTypeIcons[type];
            const count = counts[type] || 0;
            return (
              <TabsTrigger key={type} value={type} className="gap-2 data-[state=active]:bg-gray-100">
                <div className={cn('flex h-6 w-6 items-center justify-center rounded-full', colors.bg)}>
                  <Icon className={cn('h-3.5 w-3.5', colors.text)} />
                </div>
                <span className="hidden sm:inline">{config.label}</span>
                <Badge variant="secondary" className="ml-1 text-xs h-5 px-1.5">
                  {count}
                </Badge>
              </TabsTrigger>
            );
          })}
        </TabsList>

        {projectTypeOrder.map((type) => (
          <TabsContent key={type} value={type} className="space-y-6">
            {/* Stats Cards */}
            <ProjectsStats projects={projects} projectType={type} />

            {/* Filters */}
            <Card>
              <CardContent className="p-4">
                <div className="flex flex-wrap items-center gap-4">
                  <div className="relative flex-1 min-w-[200px] max-w-md">
                    <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-400" />
                    <Input
                      placeholder="Search projects..."
                      value={searchQuery}
                      onChange={(e) => setSearchQuery(e.target.value)}
                      className="pl-9"
                    />
                  </div>
                  <Select
                    value={filters.status || ''}
                    onChange={(value) => setFilters((prev) => ({ ...prev, status: value as ProjectStatus || undefined }))}
                    options={statusOptions}
                    placeholder="Status"
                    className="w-40"
                  />
                </div>
              </CardContent>
            </Card>

            {/* Projects Table */}
            <Card>
              <CardContent className="p-0">
                <DataTable
                  columns={columns}
                  data={projects}
                  isLoading={isLoading}
                  searchKey="schoolId"
                  searchPlaceholder="Search by school..."
                  showRowSelection
                  showColumnVisibility
                  onRowClick={(row) => {
                    setSelectedProject(row);
                    setShowDetail(true);
                    setDetailTab('phases');
                  }}
                  emptyMessage={`No ${projectTypeConfigs[type].label.toLowerCase()} projects found`}
                />
              </CardContent>
            </Card>
          </TabsContent>
        ))}
      </Tabs>

      {/* Project Detail Sheet */}
      <ProjectDetail
        project={selectedProject}
        open={showDetail}
        onClose={() => setShowDetail(false)}
        detailTab={detailTab}
        onDetailTabChange={setDetailTab}
        phases={phases}
        surveys={surveys}
      />

      {/* Create Project Modal */}
      <CreateProjectModal
        open={showCreateModal}
        onClose={() => setShowCreateModal(false)}
        formData={createForm}
        onFormChange={setCreateForm}
        onSubmit={handleCreateProject}
        isLoading={createProject.isPending}
      />
    </div>
  );
}
