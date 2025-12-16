import * as React from 'react';
import {
  School,
  User,
  Mail,
  Phone,
  Calendar,
  DollarSign,
  Monitor,
  MessageSquare,
  Clock,
  Tag,
  Target,
  ArrowRight,
  MapPin,
} from 'lucide-react';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Sheet, SheetHeader, SheetBody, SheetFooter } from '@/components/ui/sheet';
import { Tabs, TabsList, TabsTrigger, TabsContent } from '@/components/ui/tabs';
import { Textarea } from '@/components/ui/textarea';
import { formatDate, formatRelativeTime, cn } from '@/lib/utils';
import type {
  DemoLeadWithActivities,
  DemoLeadActivity,
  DemoLeadStage,
  UpdateLeadStageRequest,
} from '@/types/sales';
import { stageLabels, stageColors } from '@/types/sales';

interface LeadDetailSheetProps {
  lead: DemoLeadWithActivities | null;
  open: boolean;
  onClose: () => void;
  onUpdateStage: (id: string, data: UpdateLeadStageRequest) => void;
  onAddNote: (leadId: string, note: string) => void;
  onScheduleDemo: (leadId: string) => void;
  isUpdating?: boolean;
}

function InfoRow({ icon, label, value }: { icon: React.ReactNode; label: string; value: React.ReactNode }) {
  return (
    <div className="flex items-start gap-3 py-2">
      <div className="text-gray-400 mt-0.5">{icon}</div>
      <div className="flex-1 min-w-0">
        <div className="text-sm text-gray-500">{label}</div>
        <div className="font-medium text-gray-900">{value || '-'}</div>
      </div>
    </div>
  );
}

function formatCurrency(amount: number): string {
  return new Intl.NumberFormat('en-KE', {
    style: 'currency',
    currency: 'KES',
    minimumFractionDigits: 0,
  }).format(amount);
}

// Stage progression order
const stageOrder: DemoLeadStage[] = [
  'new_lead',
  'contacted',
  'demo_scheduled',
  'demo_completed',
  'proposal_sent',
  'negotiation',
  'won',
];

function getNextStage(current: DemoLeadStage): DemoLeadStage | null {
  const currentIndex = stageOrder.indexOf(current);
  if (currentIndex === -1 || currentIndex >= stageOrder.length - 1) return null;
  return stageOrder[currentIndex + 1];
}

function ActivityItem({ activity }: { activity: DemoLeadActivity }) {
  const getActivityIcon = () => {
    switch (activity.activityType) {
      case 'note':
        return <MessageSquare className="h-4 w-4 text-blue-500" />;
      case 'call':
        return <Phone className="h-4 w-4 text-green-500" />;
      case 'email':
        return <Mail className="h-4 w-4 text-purple-500" />;
      case 'meeting':
        return <Calendar className="h-4 w-4 text-orange-500" />;
      case 'demo':
        return <Monitor className="h-4 w-4 text-indigo-500" />;
      case 'stage_change':
        return <ArrowRight className="h-4 w-4 text-yellow-500" />;
      default:
        return <Clock className="h-4 w-4 text-gray-400" />;
    }
  };

  return (
    <div className="flex gap-3 py-3 border-b last:border-b-0">
      <div className="mt-0.5">{getActivityIcon()}</div>
      <div className="flex-1">
        <p className="text-sm text-gray-900">{activity.description}</p>
        <p className="text-xs text-gray-500 mt-1">
          {formatRelativeTime(activity.createdAt)}
        </p>
      </div>
    </div>
  );
}

export function LeadDetailSheet({
  lead,
  open,
  onClose,
  onUpdateStage,
  onAddNote,
  onScheduleDemo,
  isUpdating,
}: LeadDetailSheetProps) {
  const [activeTab, setActiveTab] = React.useState('details');
  const [newNote, setNewNote] = React.useState('');
  const [lostReason, setLostReason] = React.useState('');

  if (!lead) return null;

  const nextStage = getNextStage(lead.stage);
  const canProgress = lead.stage !== 'won' && lead.stage !== 'lost';

  const handleProgressStage = () => {
    if (nextStage) {
      onUpdateStage(lead.id, { stage: nextStage });
    }
  };

  const handleMarkWon = () => {
    onUpdateStage(lead.id, { stage: 'won' });
  };

  const handleMarkLost = () => {
    onUpdateStage(lead.id, { stage: 'lost', lostReason, lostNotes: '' });
    setLostReason('');
  };

  const handleAddNote = () => {
    if (newNote.trim()) {
      onAddNote(lead.id, newNote.trim());
      setNewNote('');
    }
  };

  return (
    <Sheet open={open} onClose={onClose} side="right" className="max-w-lg">
      <SheetHeader onClose={onClose}>Lead Details</SheetHeader>
      <SheetBody className="p-0">
        <div className="h-full flex flex-col">
          {/* Lead Header */}
          <div className="p-6 border-b border-gray-200 bg-gray-50">
            <div className="flex items-start gap-4">
              <div className="flex h-14 w-14 items-center justify-center rounded-xl bg-blue-100">
                <School className="h-7 w-7 text-blue-600" />
              </div>
              <div className="flex-1 min-w-0">
                <h2 className="text-lg font-semibold text-gray-900 truncate">
                  {lead.schoolName}
                </h2>
                <p className="text-sm text-gray-500">{lead.contactName}</p>
                <Badge className={cn('mt-2', stageColors[lead.stage])}>
                  {stageLabels[lead.stage]}
                </Badge>
              </div>
            </div>
          </div>

          {/* Tabs */}
          <Tabs value={activeTab} onValueChange={setActiveTab} className="flex-1 flex flex-col">
            <TabsList className="border-b px-6">
              <TabsTrigger value="details">Details</TabsTrigger>
              <TabsTrigger value="activity">Activity</TabsTrigger>
              <TabsTrigger value="actions">Actions</TabsTrigger>
            </TabsList>

            <TabsContent value="details" className="flex-1 overflow-auto p-6">
              <div className="space-y-1">
                <InfoRow
                  icon={<User className="h-4 w-4" />}
                  label="Contact"
                  value={lead.contactName}
                />
                <InfoRow
                  icon={<Mail className="h-4 w-4" />}
                  label="Email"
                  value={lead.contactEmail}
                />
                <InfoRow
                  icon={<Phone className="h-4 w-4" />}
                  label="Phone"
                  value={lead.contactPhone}
                />
                <InfoRow
                  icon={<Tag className="h-4 w-4" />}
                  label="Role"
                  value={lead.contactRole}
                />
                {(lead.countyName || lead.subCountyName) && (
                  <>
                    <div className="border-t my-4" />
                    <InfoRow
                      icon={<MapPin className="h-4 w-4" />}
                      label="County"
                      value={lead.countyName || '-'}
                    />
                    <InfoRow
                      icon={<MapPin className="h-4 w-4" />}
                      label="Sub-County"
                      value={lead.subCountyName || '-'}
                    />
                  </>
                )}
                <div className="border-t my-4" />
                <InfoRow
                  icon={<DollarSign className="h-4 w-4" />}
                  label="Estimated Value"
                  value={lead.estimatedValue ? formatCurrency(lead.estimatedValue) : '-'}
                />
                <InfoRow
                  icon={<Monitor className="h-4 w-4" />}
                  label="Estimated Devices"
                  value={lead.estimatedDevices?.toString()}
                />
                <InfoRow
                  icon={<Target className="h-4 w-4" />}
                  label="Win Probability"
                  value={`${lead.probability}%`}
                />
                <div className="border-t my-4" />
                <InfoRow
                  icon={<Calendar className="h-4 w-4" />}
                  label="Created"
                  value={formatDate(lead.createdAt)}
                />
                {lead.expectedCloseDate && (
                  <InfoRow
                    icon={<Calendar className="h-4 w-4" />}
                    label="Expected Close"
                    value={formatDate(lead.expectedCloseDate)}
                  />
                )}
                {lead.nextDemo && (
                  <InfoRow
                    icon={<Monitor className="h-4 w-4" />}
                    label="Next Demo"
                    value={`${formatDate(lead.nextDemo.scheduledDate)} at ${lead.nextDemo.scheduledTime}`}
                  />
                )}
              </div>

              {lead.notes && (
                <div className="mt-6">
                  <h4 className="font-medium text-gray-900 mb-2">Notes</h4>
                  <p className="text-sm text-gray-600 whitespace-pre-wrap">{lead.notes}</p>
                </div>
              )}

              {lead.tags && lead.tags.length > 0 && (
                <div className="mt-6">
                  <h4 className="font-medium text-gray-900 mb-2">Tags</h4>
                  <div className="flex flex-wrap gap-2">
                    {lead.tags.map((tag, i) => (
                      <Badge key={i} variant="secondary">
                        {tag}
                      </Badge>
                    ))}
                  </div>
                </div>
              )}
            </TabsContent>

            <TabsContent value="activity" className="flex-1 overflow-auto p-6">
              {/* Add Note */}
              <div className="mb-6">
                <Textarea
                  placeholder="Add a note..."
                  value={newNote}
                  onChange={(e) => setNewNote(e.target.value)}
                  rows={2}
                />
                <Button
                  size="sm"
                  className="mt-2"
                  onClick={handleAddNote}
                  disabled={!newNote.trim() || isUpdating}
                >
                  Add Note
                </Button>
              </div>

              {/* Activity Timeline */}
              <div>
                <h4 className="font-medium text-gray-900 mb-3">Recent Activity</h4>
                {lead.recentActivities && lead.recentActivities.length > 0 ? (
                  <div className="divide-y">
                    {lead.recentActivities.map((activity) => (
                      <ActivityItem key={activity.id} activity={activity} />
                    ))}
                  </div>
                ) : (
                  <p className="text-sm text-gray-500">No activity recorded yet.</p>
                )}
              </div>
            </TabsContent>

            <TabsContent value="actions" className="flex-1 overflow-auto p-6">
              <div className="space-y-6">
                {/* Stage Progression */}
                {canProgress && nextStage && (
                  <div>
                    <h4 className="font-medium text-gray-900 mb-3">Progress Lead</h4>
                    <Button
                      className="w-full"
                      onClick={handleProgressStage}
                      disabled={isUpdating}
                    >
                      Move to {stageLabels[nextStage]}
                    </Button>
                  </div>
                )}

                {/* Schedule Demo */}
                {lead.stage !== 'won' && lead.stage !== 'lost' && (
                  <div>
                    <h4 className="font-medium text-gray-900 mb-3">Schedule Demo</h4>
                    <Button
                      variant="outline"
                      className="w-full"
                      onClick={() => onScheduleDemo(lead.id)}
                    >
                      <Calendar className="h-4 w-4 mr-2" />
                      Schedule Demo
                    </Button>
                  </div>
                )}

                {/* Quick Actions */}
                {canProgress && (
                  <div className="border-t pt-6">
                    <h4 className="font-medium text-gray-900 mb-3">Quick Actions</h4>
                    <div className="grid grid-cols-2 gap-3">
                      <Button
                        variant="outline"
                        className="bg-green-50 border-green-200 text-green-700 hover:bg-green-100"
                        onClick={handleMarkWon}
                        disabled={isUpdating}
                      >
                        Mark as Won
                      </Button>
                      <Button
                        variant="outline"
                        className="bg-red-50 border-red-200 text-red-700 hover:bg-red-100"
                        onClick={handleMarkLost}
                        disabled={isUpdating || !lostReason}
                      >
                        Mark as Lost
                      </Button>
                    </div>
                    <Textarea
                      placeholder="Lost reason (required to mark as lost)"
                      className="mt-3"
                      value={lostReason}
                      onChange={(e) => setLostReason(e.target.value)}
                      rows={2}
                    />
                  </div>
                )}

                {/* Won/Lost Status */}
                {lead.stage === 'won' && (
                  <div className="bg-green-50 p-4 rounded-lg">
                    <p className="text-green-800 font-medium">This deal was won!</p>
                  </div>
                )}
                {lead.stage === 'lost' && (
                  <div className="bg-red-50 p-4 rounded-lg">
                    <p className="text-red-800 font-medium">This deal was lost</p>
                    {lead.lostReason && (
                      <p className="text-red-600 text-sm mt-1">Reason: {lead.lostReason}</p>
                    )}
                  </div>
                )}
              </div>
            </TabsContent>
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
