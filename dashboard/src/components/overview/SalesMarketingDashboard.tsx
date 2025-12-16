import { Link } from 'react-router-dom';
import {
  TrendingUp,
  Target,
  Users,
  School,
  ArrowRight,
  Presentation,
  BarChart3,
  Calendar,
  CheckCircle2,
  Clock,
  FileText,
  DollarSign,
} from 'lucide-react';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { useAuth } from '@/contexts/AuthContext';

export function SalesMarketingDashboard() {
  const { user } = useAuth();
  const displayName = user?.displayName || user?.username || 'Sales Rep';

  // Mock data - in production, fetch from API
  const stats = {
    totalLeads: 124,
    activeDeals: 18,
    demoScheduled: 7,
    closedThisMonth: 12,
    pipelineValue: 2450000,
    conversionRate: 23.5,
  };

  const upcomingDemos = [
    { id: 1, school: 'Nairobi Primary School', date: '2024-01-15 10:00', stage: 'Demo', contact: 'John Kamau' },
    { id: 2, school: 'Mombasa Girls High', date: '2024-01-16 14:00', stage: 'Follow-up', contact: 'Mary Njeri' },
    { id: 3, school: 'Kisumu Academy', date: '2024-01-17 09:00', stage: 'Proposal', contact: 'Peter Ochieng' },
  ];

  const recentActivity = [
    { type: 'demo_completed', school: 'Eldoret Boys', time: '2 hours ago' },
    { type: 'proposal_sent', school: 'Nakuru Mixed', time: '4 hours ago' },
    { type: 'deal_closed', school: 'Thika Primary', time: '1 day ago' },
    { type: 'lead_added', school: 'Nyeri Girls', time: '2 days ago' },
  ];

  return (
    <div className="space-y-6">
      {/* Welcome Header - Emerald/Green Theme (Sales - Growth Focused) */}
      <div className="relative overflow-hidden rounded-2xl bg-gradient-to-r from-emerald-600 via-green-600 to-teal-600 p-6 text-white shadow-lg">
        <div className="absolute top-0 right-0 w-64 h-64 bg-white/5 rounded-full -translate-y-32 translate-x-32" />
        <div className="absolute bottom-0 left-0 w-48 h-48 bg-white/5 rounded-full translate-y-24 -translate-x-24" />

        <div className="relative">
          <div className="flex items-center gap-3 mb-2">
            <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-white/20 backdrop-blur">
              <TrendingUp className="h-6 w-6" />
            </div>
            <div>
              <h1 className="text-2xl md:text-3xl font-bold">Welcome, {displayName}!</h1>
              <p className="text-emerald-100">Sales & Marketing Dashboard</p>
            </div>
          </div>

          {/* Summary Stats */}
          <div className="mt-6 grid grid-cols-2 md:grid-cols-6 gap-4">
            <div className="bg-white/10 backdrop-blur rounded-xl p-3">
              <p className="text-2xl font-bold">{stats.totalLeads}</p>
              <p className="text-emerald-200 text-sm">Total Leads</p>
            </div>
            <div className="bg-white/10 backdrop-blur rounded-xl p-3">
              <p className="text-2xl font-bold">{stats.activeDeals}</p>
              <p className="text-emerald-200 text-sm">Active Deals</p>
            </div>
            <div className="bg-white/10 backdrop-blur rounded-xl p-3">
              <p className="text-2xl font-bold">{stats.demoScheduled}</p>
              <p className="text-emerald-200 text-sm">Demos Scheduled</p>
            </div>
            <div className="bg-white/10 backdrop-blur rounded-xl p-3">
              <p className="text-2xl font-bold">{stats.closedThisMonth}</p>
              <p className="text-emerald-200 text-sm">Closed (Month)</p>
            </div>
            <div className="bg-white/10 backdrop-blur rounded-xl p-3">
              <p className="text-2xl font-bold">KES {(stats.pipelineValue / 1000000).toFixed(1)}M</p>
              <p className="text-emerald-200 text-sm">Pipeline Value</p>
            </div>
            <div className="bg-white/10 backdrop-blur rounded-xl p-3">
              <p className="text-2xl font-bold">{stats.conversionRate}%</p>
              <p className="text-emerald-200 text-sm">Conversion Rate</p>
            </div>
          </div>
        </div>
      </div>

      {/* Quick Actions */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <Link to="/demo-pipeline">
          <Card className="hover:shadow-md transition-shadow cursor-pointer border-violet-200 bg-violet-50">
            <CardContent className="p-4 flex items-center gap-3">
              <div className="h-10 w-10 rounded-lg bg-violet-100 flex items-center justify-center">
                <Target className="h-5 w-5 text-violet-600" />
              </div>
              <div className="flex-1">
                <p className="font-medium text-violet-900">Demo Pipeline</p>
                <p className="text-sm text-violet-600">Manage opportunities</p>
              </div>
              <ArrowRight className="h-5 w-5 text-violet-400" />
            </CardContent>
          </Card>
        </Link>

        <Link to="/schools">
          <Card className="hover:shadow-md transition-shadow cursor-pointer border-green-200 bg-green-50">
            <CardContent className="p-4 flex items-center gap-3">
              <div className="h-10 w-10 rounded-lg bg-green-100 flex items-center justify-center">
                <School className="h-5 w-5 text-green-600" />
              </div>
              <div className="flex-1">
                <p className="font-medium text-green-900">Schools Directory</p>
                <p className="text-sm text-green-600">Browse all schools</p>
              </div>
              <ArrowRight className="h-5 w-5 text-green-400" />
            </CardContent>
          </Card>
        </Link>

        <Link to="/presentations">
          <Card className="hover:shadow-md transition-shadow cursor-pointer border-pink-200 bg-pink-50">
            <CardContent className="p-4 flex items-center gap-3">
              <div className="h-10 w-10 rounded-lg bg-pink-100 flex items-center justify-center">
                <Presentation className="h-5 w-5 text-pink-600" />
              </div>
              <div className="flex-1">
                <p className="font-medium text-pink-900">Presentations</p>
                <p className="text-sm text-pink-600">Sales materials</p>
              </div>
              <ArrowRight className="h-5 w-5 text-pink-400" />
            </CardContent>
          </Card>
        </Link>

        <Link to="/reports">
          <Card className="hover:shadow-md transition-shadow cursor-pointer border-teal-200 bg-teal-50">
            <CardContent className="p-4 flex items-center gap-3">
              <div className="h-10 w-10 rounded-lg bg-teal-100 flex items-center justify-center">
                <BarChart3 className="h-5 w-5 text-teal-600" />
              </div>
              <div className="flex-1">
                <p className="font-medium text-teal-900">Sales Reports</p>
                <p className="text-sm text-teal-600">View analytics</p>
              </div>
              <ArrowRight className="h-5 w-5 text-teal-400" />
            </CardContent>
          </Card>
        </Link>
      </div>

      {/* Main Content Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Upcoming Demos */}
        <Card>
          <CardHeader className="pb-3">
            <div className="flex items-center justify-between">
              <CardTitle className="text-lg flex items-center gap-2">
                <Calendar className="h-5 w-5 text-violet-600" />
                Upcoming Demos & Meetings
              </CardTitle>
              <Link to="/demo-pipeline">
                <Button variant="ghost" size="sm">View All</Button>
              </Link>
            </div>
          </CardHeader>
          <CardContent>
            <div className="space-y-3">
              {upcomingDemos.map((demo) => (
                <div
                  key={demo.id}
                  className="flex items-center justify-between p-3 rounded-lg bg-gray-50 hover:bg-gray-100 transition-colors"
                >
                  <div className="flex items-center gap-3">
                    <div className="h-10 w-10 rounded-lg bg-violet-100 flex items-center justify-center">
                      <School className="h-5 w-5 text-violet-600" />
                    </div>
                    <div>
                      <p className="font-medium text-gray-900">{demo.school}</p>
                      <p className="text-sm text-gray-500">{demo.contact} • {demo.date}</p>
                    </div>
                  </div>
                  <Badge variant="outline" className="bg-violet-50 text-violet-700 border-violet-200">
                    {demo.stage}
                  </Badge>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>

        {/* Recent Activity */}
        <Card>
          <CardHeader className="pb-3">
            <div className="flex items-center justify-between">
              <CardTitle className="text-lg flex items-center gap-2">
                <Clock className="h-5 w-5 text-emerald-600" />
                Recent Activity
              </CardTitle>
            </div>
          </CardHeader>
          <CardContent>
            <div className="space-y-3">
              {recentActivity.map((activity, idx) => (
                <div
                  key={idx}
                  className="flex items-center gap-3 p-3 rounded-lg bg-gray-50"
                >
                  <div className={`h-8 w-8 rounded-full flex items-center justify-center ${
                    activity.type === 'deal_closed' ? 'bg-green-100' :
                    activity.type === 'demo_completed' ? 'bg-violet-100' :
                    activity.type === 'proposal_sent' ? 'bg-blue-100' : 'bg-gray-100'
                  }`}>
                    {activity.type === 'deal_closed' && <CheckCircle2 className="h-4 w-4 text-green-600" />}
                    {activity.type === 'demo_completed' && <Presentation className="h-4 w-4 text-violet-600" />}
                    {activity.type === 'proposal_sent' && <FileText className="h-4 w-4 text-blue-600" />}
                    {activity.type === 'lead_added' && <Users className="h-4 w-4 text-gray-600" />}
                  </div>
                  <div className="flex-1">
                    <p className="text-sm text-gray-900">
                      {activity.type === 'deal_closed' && 'Deal closed with '}
                      {activity.type === 'demo_completed' && 'Demo completed at '}
                      {activity.type === 'proposal_sent' && 'Proposal sent to '}
                      {activity.type === 'lead_added' && 'New lead added: '}
                      <span className="font-medium">{activity.school}</span>
                    </p>
                    <p className="text-xs text-gray-500">{activity.time}</p>
                  </div>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Pipeline Summary */}
      <Card>
        <CardHeader className="pb-3">
          <CardTitle className="text-lg flex items-center gap-2">
            <DollarSign className="h-5 w-5 text-emerald-600" />
            Pipeline by Stage
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-2 md:grid-cols-5 gap-4">
            {[
              { stage: 'Lead', count: 45, value: '450K', color: 'bg-gray-500' },
              { stage: 'Qualified', count: 28, value: '680K', color: 'bg-blue-500' },
              { stage: 'Demo', count: 15, value: '520K', color: 'bg-violet-500' },
              { stage: 'Proposal', count: 8, value: '480K', color: 'bg-amber-500' },
              { stage: 'Negotiation', count: 4, value: '320K', color: 'bg-emerald-500' },
            ].map((stage) => (
              <div key={stage.stage} className="text-center p-4 rounded-xl bg-gray-50">
                <div className={`w-3 h-3 rounded-full ${stage.color} mx-auto mb-2`} />
                <p className="text-2xl font-bold text-gray-900">{stage.count}</p>
                <p className="text-sm text-gray-600">{stage.stage}</p>
                <p className="text-xs text-gray-400 mt-1">KES {stage.value}</p>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
