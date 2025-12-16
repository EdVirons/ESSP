import { useQuery } from '@tanstack/react-query';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import {
  TrendingUp,
  Target,
  Users,
  School,
  Calendar,
  FileText,
  ArrowRight,
  BarChart3,
  Presentation,
  Loader2,
  DollarSign,
  CheckCircle,
} from 'lucide-react';
import { Link } from 'react-router-dom';
import { salesApi } from '@/api/sales';
import { stageLabels } from '@/types/sales';
import type { PipelineStageCount, RecentActivity } from '@/types/sales';

function formatCurrency(amount: number): string {
  if (amount >= 1000000) {
    return `KES ${(amount / 1000000).toFixed(1)}M`;
  }
  if (amount >= 1000) {
    return `KES ${(amount / 1000).toFixed(0)}K`;
  }
  return `KES ${amount.toFixed(0)}`;
}

function formatRelativeTime(dateString: string): string {
  const date = new Date(dateString);
  const now = new Date();
  const diffMs = now.getTime() - date.getTime();
  const diffMins = Math.floor(diffMs / (1000 * 60));
  const diffHours = Math.floor(diffMs / (1000 * 60 * 60));
  const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24));

  if (diffMins < 60) return `${diffMins} min ago`;
  if (diffHours < 24) return `${diffHours} hours ago`;
  if (diffDays === 1) return 'Yesterday';
  if (diffDays < 7) return `${diffDays} days ago`;
  return date.toLocaleDateString();
}

function getActivityIcon(type: string) {
  switch (type) {
    case 'demo':
      return <Users className="h-4 w-4 text-blue-600" />;
    case 'stage_change':
      return <ArrowRight className="h-4 w-4 text-yellow-600" />;
    case 'note':
      return <FileText className="h-4 w-4 text-purple-600" />;
    case 'call':
      return <Calendar className="h-4 w-4 text-amber-600" />;
    case 'created':
      return <Target className="h-4 w-4 text-gray-600" />;
    default:
      return <TrendingUp className="h-4 w-4 text-green-600" />;
  }
}

const stageColorMap: Record<string, string> = {
  new_lead: 'bg-blue-400',
  contacted: 'bg-yellow-400',
  demo_scheduled: 'bg-purple-400',
  demo_completed: 'bg-indigo-400',
  proposal_sent: 'bg-orange-400',
  negotiation: 'bg-pink-400',
  won: 'bg-green-400',
  lost: 'bg-red-400',
};

export function SalesDashboard() {
  const { data, isLoading, error } = useQuery({
    queryKey: ['sales-dashboard'],
    queryFn: () => salesApi.getDashboard(30),
    refetchInterval: 60000, // Refresh every minute
  });

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <Loader2 className="h-8 w-8 animate-spin text-gray-400" />
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex items-center justify-center h-64">
        <p className="text-red-500">Failed to load dashboard data</p>
      </div>
    );
  }

  const metrics = data?.metrics || {
    totalLeads: 0,
    newLeadsThisPeriod: 0,
    demosScheduled: 0,
    demosCompleted: 0,
    proposalsSent: 0,
    dealsWon: 0,
    dealsLost: 0,
    totalPipelineValue: 0,
    wonValueThisPeriod: 0,
    conversionRate: 0,
    winRate: 0,
    averageDealSize: 0,
  };

  const pipelineStages = data?.pipelineStages || [];
  const recentActivities = data?.recentActivities || [];
  const schoolsByRegion = data?.schoolsByRegion || [];

  // Filter out lost stage for display
  const displayStages = pipelineStages.filter(
    (s: PipelineStageCount) => s.stage !== 'lost'
  );

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Sales Dashboard</h1>
          <p className="text-sm text-gray-500 mt-1">
            Track your sales pipeline and demo performance
          </p>
        </div>
        <div className="flex gap-2">
          <Button variant="outline" asChild>
            <Link to="/presentations">
              <Presentation className="h-4 w-4 mr-2" />
              Presentations
            </Link>
          </Button>
          <Button asChild>
            <Link to="/demo-pipeline">
              <Target className="h-4 w-4 mr-2" />
              View Pipeline
            </Link>
          </Button>
        </div>
      </div>

      {/* Metrics Cards */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Demos Completed</CardTitle>
            <Calendar className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{metrics.demosCompleted}</div>
            <p className="text-xs text-gray-500">
              {metrics.demosScheduled} scheduled
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Pipeline Value</CardTitle>
            <TrendingUp className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {formatCurrency(metrics.totalPipelineValue)}
            </div>
            <p className="text-xs text-green-600">
              {formatCurrency(metrics.wonValueThisPeriod)} won this period
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Win Rate</CardTitle>
            <BarChart3 className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{metrics.winRate.toFixed(0)}%</div>
            <p className="text-xs text-gray-500">
              {metrics.dealsWon} won / {metrics.dealsLost} lost
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Leads</CardTitle>
            <School className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{metrics.totalLeads}</div>
            <p className="text-xs text-green-600">
              +{metrics.newLeadsThisPeriod} this period
            </p>
          </CardContent>
        </Card>
      </div>

      {/* Pipeline Overview */}
      <Card>
        <CardHeader className="flex flex-row items-center justify-between">
          <CardTitle>Demo Pipeline</CardTitle>
          <Button variant="ghost" size="sm" asChild>
            <Link to="/demo-pipeline">
              View All <ArrowRight className="ml-2 h-4 w-4" />
            </Link>
          </Button>
        </CardHeader>
        <CardContent>
          {displayStages.length === 0 ? (
            <div className="text-center py-8 text-gray-500">
              No pipeline data available
            </div>
          ) : (
            <div className="flex items-center justify-between gap-2">
              {displayStages.map((stage: PipelineStageCount) => (
                <div key={stage.stage} className="flex-1 text-center relative">
                  <div
                    className={`h-16 ${
                      stageColorMap[stage.stage] || 'bg-gray-400'
                    } rounded-lg flex items-center justify-center mb-2`}
                  >
                    <span className="text-2xl font-bold text-white">{stage.count}</span>
                  </div>
                  <p className="text-xs text-gray-600 truncate">
                    {stageLabels[stage.stage] || stage.stage}
                  </p>
                  <p className="text-xs text-gray-400">
                    {formatCurrency(stage.totalValue)}
                  </p>
                </div>
              ))}
            </div>
          )}
        </CardContent>
      </Card>

      {/* Additional Metrics Row */}
      <div className="grid gap-4 md:grid-cols-3">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Proposals Sent</CardTitle>
            <FileText className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{metrics.proposalsSent}</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Deals Won</CardTitle>
            <CheckCircle className="h-4 w-4 text-green-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-green-600">{metrics.dealsWon}</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Average Deal Size</CardTitle>
            <DollarSign className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {formatCurrency(metrics.averageDealSize)}
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Two Column Layout */}
      <div className="grid gap-6 lg:grid-cols-2">
        {/* Recent Activities */}
        <Card>
          <CardHeader>
            <CardTitle>Recent Activities</CardTitle>
          </CardHeader>
          <CardContent>
            {recentActivities.length === 0 ? (
              <div className="text-center py-8 text-gray-500">
                No recent activities
              </div>
            ) : (
              <div className="space-y-4">
                {recentActivities.map((activity: RecentActivity) => (
                  <div key={activity.id} className="flex items-start gap-3">
                    <div className="p-2 rounded-full bg-gray-100">
                      {getActivityIcon(activity.type)}
                    </div>
                    <div className="flex-1 min-w-0">
                      <p className="text-sm font-medium text-gray-900 truncate">
                        {activity.leadName || 'Unknown Lead'}
                      </p>
                      <p className="text-sm text-gray-500">{activity.description}</p>
                    </div>
                    <span className="text-xs text-gray-400">
                      {formatRelativeTime(activity.createdAt)}
                    </span>
                  </div>
                ))}
              </div>
            )}
          </CardContent>
        </Card>

        {/* Schools by Region */}
        <Card>
          <CardHeader className="flex flex-row items-center justify-between">
            <CardTitle>Leads by Region</CardTitle>
            <Button variant="ghost" size="sm" asChild>
              <Link to="/demo-pipeline">
                View All <ArrowRight className="ml-2 h-4 w-4" />
              </Link>
            </Button>
          </CardHeader>
          <CardContent>
            {schoolsByRegion.length === 0 ? (
              <div className="text-center py-8 text-gray-500">
                No regional data available
              </div>
            ) : (
              <div className="space-y-4">
                {schoolsByRegion.map((region) => (
                  <div key={region.region} className="flex items-center gap-4">
                    <div className="w-24 text-sm font-medium text-gray-900">
                      {region.region}
                    </div>
                    <div className="flex-1">
                      <div className="flex h-4 overflow-hidden rounded-full bg-gray-100">
                        <div
                          className="bg-blue-500"
                          style={{
                            width: `${Math.min(
                              100,
                              (region.count /
                                Math.max(...schoolsByRegion.map((r) => r.count))) *
                                100
                            )}%`,
                          }}
                        />
                      </div>
                    </div>
                    <div className="text-sm text-gray-600 w-16 text-right">
                      {region.count} leads
                    </div>
                  </div>
                ))}
              </div>
            )}
          </CardContent>
        </Card>
      </div>

      {/* Quick Actions */}
      <Card>
        <CardHeader>
          <CardTitle>Quick Actions</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid gap-4 md:grid-cols-4">
            <Button variant="outline" className="h-20 flex-col" asChild>
              <Link to="/demo-pipeline">
                <Calendar className="h-6 w-6 mb-2" />
                Schedule Demo
              </Link>
            </Button>
            <Button variant="outline" className="h-20 flex-col" asChild>
              <Link to="/schools">
                <School className="h-6 w-6 mb-2" />
                View Schools
              </Link>
            </Button>
            <Button variant="outline" className="h-20 flex-col" asChild>
              <Link to="/presentations">
                <Presentation className="h-6 w-6 mb-2" />
                Presentations
              </Link>
            </Button>
            <Button variant="outline" className="h-20 flex-col" asChild>
              <Link to="/projects">
                <FileText className="h-6 w-6 mb-2" />
                View Projects
              </Link>
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
