import { useState, useCallback } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import {
  Target,
  Plus,
  Search,
  Filter,
  Calendar,
  School,
  User,
  MoreVertical,
  Loader2,
} from 'lucide-react';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { demoPipelineApi } from '@/api/demo-pipeline';
import { CreateLeadModal } from '@/components/sales/CreateLeadModal';
import { LeadDetailSheet } from '@/components/sales/LeadDetailSheet';
import { ScheduleDemoModal } from '@/components/sales/ScheduleDemoModal';
import { toast } from '@/lib/toast';
import type {
  DemoLead,
  DemoLeadWithActivities,
  DemoLeadStage,
  CreateDemoLeadRequest,
  UpdateLeadStageRequest,
  CreateDemoScheduleRequest,
} from '@/types/sales';
import { stageLabels, stageColors } from '@/types/sales';

// Map API stages to display stages
const displayStages: DemoLeadStage[] = [
  'new_lead',
  'contacted',
  'demo_scheduled',
  'demo_completed',
  'proposal_sent',
  'negotiation',
  'won',
];

function formatCurrency(amount: number): string {
  return new Intl.NumberFormat('en-KE', {
    style: 'currency',
    currency: 'KES',
    minimumFractionDigits: 0,
  }).format(amount);
}

interface PipelineCardProps {
  lead: DemoLead;
  onClick: () => void;
  onScheduleDemo: () => void;
  onAddNote: () => void;
}

function PipelineCard({ lead, onClick, onScheduleDemo, onAddNote }: PipelineCardProps) {
  return (
    <Card className="mb-3 cursor-pointer hover:shadow-md transition-shadow" onClick={onClick}>
      <CardContent className="p-4">
        <div className="flex items-start justify-between mb-2">
          <div className="flex items-center gap-2">
            <School className="h-4 w-4 text-gray-400" />
            <span className="font-medium text-sm truncate">{lead.schoolName}</span>
          </div>
          <DropdownMenu>
            <DropdownMenuTrigger asChild onClick={(e) => e.stopPropagation()}>
              <Button variant="ghost" size="sm" className="h-8 w-8 p-0">
                <MoreVertical className="h-4 w-4" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuItem onClick={(e) => { e.stopPropagation(); onClick(); }}>
                View Details
              </DropdownMenuItem>
              <DropdownMenuItem onClick={(e) => { e.stopPropagation(); onScheduleDemo(); }}>
                Schedule Demo
              </DropdownMenuItem>
              <DropdownMenuItem onClick={(e) => { e.stopPropagation(); onAddNote(); }}>
                Add Note
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>

        <div className="space-y-2 text-xs text-gray-500">
          <div className="flex items-center gap-2">
            <User className="h-3 w-3" />
            <span className="truncate">{lead.contactName || 'No contact'}</span>
          </div>
          {lead.expectedCloseDate && (
            <div className="flex items-center gap-2">
              <Calendar className="h-3 w-3" />
              <span>{new Date(lead.expectedCloseDate).toLocaleDateString()}</span>
            </div>
          )}
        </div>

        <div className="mt-3 pt-3 border-t flex items-center justify-between">
          <span className="font-semibold text-sm text-gray-900">
            {lead.estimatedValue ? formatCurrency(lead.estimatedValue) : '-'}
          </span>
          <Badge className={stageColors[lead.stage]}>
            {stageLabels[lead.stage]}
          </Badge>
        </div>
      </CardContent>
    </Card>
  );
}

export function DemoPipeline() {
  const queryClient = useQueryClient();
  const [searchTerm, setSearchTerm] = useState('');
  const [selectedStage, setSelectedStage] = useState<DemoLeadStage | 'all'>('all');
  const [createModalOpen, setCreateModalOpen] = useState(false);
  const [scheduleModalOpen, setScheduleModalOpen] = useState(false);
  const [scheduleLeadId, setScheduleLeadId] = useState<string | null>(null);
  const [scheduleLeadName, setScheduleLeadName] = useState<string>('');
  const [selectedLead, setSelectedLead] = useState<DemoLeadWithActivities | null>(null);
  const [detailSheetOpen, setDetailSheetOpen] = useState(false);

  // Fetch leads
  const { data, isLoading, error } = useQuery({
    queryKey: ['demo-leads', searchTerm, selectedStage],
    queryFn: () => demoPipelineApi.listLeads({
      search: searchTerm || undefined,
      stage: selectedStage !== 'all' ? selectedStage : undefined,
      limit: 100,
    }),
  });

  const leads = data?.leads || [];

  // Create lead mutation
  const createMutation = useMutation({
    mutationFn: (data: CreateDemoLeadRequest) => demoPipelineApi.createLead(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['demo-leads'] });
      setCreateModalOpen(false);
      toast.success('Lead created', 'New lead has been added to the pipeline');
    },
    onError: () => {
      toast.error('Failed to create lead', 'Please try again');
    },
  });

  // Update stage mutation
  const updateStageMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateLeadStageRequest }) =>
      demoPipelineApi.updateStage(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['demo-leads'] });
      if (selectedLead) {
        fetchLeadDetails(selectedLead.id);
      }
      toast.success('Stage updated', 'Lead stage has been updated');
    },
    onError: () => {
      toast.error('Failed to update stage', 'Please try again');
    },
  });

  // Add note mutation
  const addNoteMutation = useMutation({
    mutationFn: ({ leadId, note }: { leadId: string; note: string }) =>
      demoPipelineApi.addNote(leadId, note),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ['demo-leads'] });
      if (selectedLead && selectedLead.id === variables.leadId) {
        fetchLeadDetails(variables.leadId);
      }
      toast.success('Note added', 'Note has been added to the lead');
    },
    onError: () => {
      toast.error('Failed to add note', 'Please try again');
    },
  });

  // Schedule demo mutation
  const scheduleDemoMutation = useMutation({
    mutationFn: ({ leadId, data }: { leadId: string; data: CreateDemoScheduleRequest }) =>
      demoPipelineApi.scheduleDemo(leadId, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['demo-leads'] });
      setScheduleModalOpen(false);
      setScheduleLeadId(null);
      setScheduleLeadName('');
      toast.success('Demo scheduled', 'The demo has been scheduled successfully');
    },
    onError: () => {
      toast.error('Failed to schedule demo', 'Please try again');
    },
  });

  const fetchLeadDetails = useCallback(async (id: string) => {
    try {
      const lead = await demoPipelineApi.getLead(id);
      setSelectedLead(lead);
    } catch {
      toast.error('Failed to load lead details', 'Please try again');
    }
  }, []);

  const handleLeadClick = async (lead: DemoLead) => {
    await fetchLeadDetails(lead.id);
    setDetailSheetOpen(true);
  };

  const handleUpdateStage = (id: string, data: UpdateLeadStageRequest) => {
    updateStageMutation.mutate({ id, data });
  };

  const handleAddNote = (leadId: string, note: string) => {
    addNoteMutation.mutate({ leadId, note });
  };

  const handleScheduleDemo = (leadId: string, leadName?: string) => {
    setScheduleLeadId(leadId);
    setScheduleLeadName(leadName || '');
    setScheduleModalOpen(true);
  };

  const handleScheduleDemoSubmit = (data: CreateDemoScheduleRequest) => {
    if (scheduleLeadId) {
      scheduleDemoMutation.mutate({ leadId: scheduleLeadId, data });
    }
  };

  const getLeadsByStage = (stage: DemoLeadStage) =>
    leads.filter((lead) => lead.stage === stage);

  const getTotalValueByStage = (stage: DemoLeadStage) =>
    getLeadsByStage(stage).reduce((sum, lead) => sum + (lead.estimatedValue || 0), 0);

  if (error) {
    return (
      <div className="flex items-center justify-center h-64">
        <p className="text-red-500">Failed to load pipeline data</p>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Demo Pipeline</h1>
          <p className="text-sm text-gray-500 mt-1">
            Manage your sales pipeline and track demo progress
          </p>
        </div>
        <Button onClick={() => setCreateModalOpen(true)}>
          <Plus className="h-4 w-4 mr-2" />
          Add Lead
        </Button>
      </div>

      {/* Filters */}
      <div className="flex items-center gap-4">
        <div className="relative flex-1 max-w-sm">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-gray-400" />
          <Input
            placeholder="Search schools or contacts..."
            className="pl-10"
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
          />
        </div>
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="outline">
              <Filter className="h-4 w-4 mr-2" />
              {selectedStage === 'all' ? 'All Stages' : stageLabels[selectedStage]}
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent>
            <DropdownMenuItem onClick={() => setSelectedStage('all')}>
              All Stages
            </DropdownMenuItem>
            {displayStages.map((stage) => (
              <DropdownMenuItem key={stage} onClick={() => setSelectedStage(stage)}>
                {stageLabels[stage]}
              </DropdownMenuItem>
            ))}
          </DropdownMenuContent>
        </DropdownMenu>
      </div>

      {/* Loading State */}
      {isLoading && (
        <div className="flex items-center justify-center h-64">
          <Loader2 className="h-8 w-8 animate-spin text-gray-400" />
        </div>
      )}

      {/* Pipeline Kanban - Desktop View */}
      {!isLoading && (
        <div className="hidden lg:grid grid-cols-7 gap-4 overflow-x-auto pb-4">
          {displayStages.map((stage) => {
            const stageLeads = getLeadsByStage(stage);
            const totalValue = getTotalValueByStage(stage);

            return (
              <div key={stage} className="min-w-[240px]">
                <div className={`rounded-t-lg px-4 py-3 ${stageColors[stage].replace('text-', 'bg-').replace('800', '100')}`}>
                  <div className="flex items-center justify-between">
                    <h3 className={`font-semibold text-sm ${stageColors[stage].split(' ')[1]}`}>
                      {stageLabels[stage]}
                    </h3>
                    <Badge variant="secondary" className="bg-white">
                      {stageLeads.length}
                    </Badge>
                  </div>
                  <p className="text-xs text-gray-500 mt-1">
                    {formatCurrency(totalValue)}
                  </p>
                </div>

                <div className="bg-gray-50 rounded-b-lg p-3 min-h-[400px]">
                  {stageLeads.length === 0 ? (
                    <div className="text-center py-8 text-gray-400 text-sm">
                      No leads
                    </div>
                  ) : (
                    stageLeads.map((lead) => (
                      <PipelineCard
                        key={lead.id}
                        lead={lead}
                        onClick={() => handleLeadClick(lead)}
                        onScheduleDemo={() => handleScheduleDemo(lead.id, lead.schoolName)}
                        onAddNote={() => {
                          handleLeadClick(lead);
                        }}
                      />
                    ))
                  )}
                </div>
              </div>
            );
          })}
        </div>
      )}

      {/* Pipeline List - Mobile View */}
      {!isLoading && (
        <div className="lg:hidden space-y-4">
          {displayStages.map((stage) => {
            const stageLeads = getLeadsByStage(stage);
            const totalValue = getTotalValueByStage(stage);

            return (
              <Card key={stage} className="overflow-hidden">
                <div className={`px-4 py-3 ${stageColors[stage].replace('text-', 'bg-').replace('800', '100')}`}>
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-3">
                      <h3 className={`font-semibold ${stageColors[stage].split(' ')[1]}`}>
                        {stageLabels[stage]}
                      </h3>
                      <Badge variant="secondary" className="bg-white">
                        {stageLeads.length}
                      </Badge>
                    </div>
                    <p className="text-sm font-medium text-gray-700">
                      {formatCurrency(totalValue)}
                    </p>
                  </div>
                </div>
                <CardContent className="p-3">
                  {stageLeads.length === 0 ? (
                    <p className="text-center py-4 text-gray-400 text-sm">No leads</p>
                  ) : (
                    <div className="space-y-2">
                      {stageLeads.slice(0, 3).map((lead) => (
                        <div
                          key={lead.id}
                          onClick={() => handleLeadClick(lead)}
                          className="flex items-center justify-between p-3 bg-gray-50 rounded-lg cursor-pointer hover:bg-gray-100"
                        >
                          <div>
                            <p className="font-medium text-sm">{lead.schoolName}</p>
                            <p className="text-xs text-gray-500">{lead.contactName || 'No contact'}</p>
                          </div>
                          <p className="font-semibold text-sm">
                            {lead.estimatedValue ? formatCurrency(lead.estimatedValue) : '-'}
                          </p>
                        </div>
                      ))}
                      {stageLeads.length > 3 && (
                        <p className="text-center text-xs text-gray-500 py-2">
                          +{stageLeads.length - 3} more leads
                        </p>
                      )}
                    </div>
                  )}
                </CardContent>
              </Card>
            );
          })}
        </div>
      )}

      {/* Summary Card */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Target className="h-5 w-5" />
            Pipeline Summary
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-6">
            <div>
              <p className="text-sm text-gray-500">Total Leads</p>
              <p className="text-2xl font-bold">{leads.length}</p>
            </div>
            <div>
              <p className="text-sm text-gray-500">Total Pipeline Value</p>
              <p className="text-2xl font-bold">
                {formatCurrency(leads.reduce((sum, lead) => sum + (lead.estimatedValue || 0), 0))}
              </p>
            </div>
            <div>
              <p className="text-sm text-gray-500">Average Deal Size</p>
              <p className="text-2xl font-bold">
                {formatCurrency(
                  leads.length > 0
                    ? leads.reduce((sum, lead) => sum + (lead.estimatedValue || 0), 0) / leads.length
                    : 0
                )}
              </p>
            </div>
            <div>
              <p className="text-sm text-gray-500">Won This Month</p>
              <p className="text-2xl font-bold text-green-600">
                {getLeadsByStage('won').length}
              </p>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Create Lead Modal */}
      <CreateLeadModal
        open={createModalOpen}
        onClose={() => setCreateModalOpen(false)}
        onSubmit={(data) => createMutation.mutate(data)}
        isLoading={createMutation.isPending}
      />

      {/* Lead Detail Sheet */}
      <LeadDetailSheet
        lead={selectedLead}
        open={detailSheetOpen}
        onClose={() => {
          setDetailSheetOpen(false);
          setSelectedLead(null);
        }}
        onUpdateStage={handleUpdateStage}
        onAddNote={handleAddNote}
        onScheduleDemo={(leadId) => handleScheduleDemo(leadId, selectedLead?.schoolName)}
        isUpdating={updateStageMutation.isPending || addNoteMutation.isPending}
      />

      {/* Schedule Demo Modal */}
      <ScheduleDemoModal
        open={scheduleModalOpen}
        onClose={() => {
          setScheduleModalOpen(false);
          setScheduleLeadId(null);
          setScheduleLeadName('');
        }}
        onSubmit={handleScheduleDemoSubmit}
        isLoading={scheduleDemoMutation.isPending}
        leadName={scheduleLeadName}
      />
    </div>
  );
}
