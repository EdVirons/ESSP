import { Link } from 'react-router-dom';
import {
  Target,
  Layers,
  School,
  ArrowRight,
  Calendar,
  Clock,
  MapPin,
  Users,
  ClipboardList,
  Camera,
} from 'lucide-react';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { useAuth } from '@/contexts/AuthContext';

export function DemoTeamDashboard() {
  const { user } = useAuth();
  const displayName = user?.displayName || user?.username || 'Demo Specialist';

  // Mock data - in production, fetch from API
  const stats = {
    scheduledDemos: 5,
    completedThisWeek: 8,
    pendingSurveys: 3,
    activeProjects: 12,
    schoolsVisited: 156,
    avgDemoRating: 4.7,
  };

  const todaySchedule = [
    { id: 1, school: 'Starehe Boys Centre', time: '09:00 AM', type: 'Full Demo', location: 'Nairobi' },
    { id: 2, school: 'Loreto Convent', time: '11:30 AM', type: 'Follow-up', location: 'Nairobi' },
    { id: 3, school: 'Moi Forces Academy', time: '02:00 PM', type: 'Survey', location: 'Nairobi' },
  ];

  const pendingTasks = [
    { type: 'survey', school: 'Kenya High School', dueDate: 'Tomorrow', priority: 'high' },
    { type: 'report', school: 'Alliance Girls', dueDate: 'In 2 days', priority: 'medium' },
    { type: 'boq', school: 'Nairobi School', dueDate: 'In 3 days', priority: 'low' },
  ];

  const recentProjects = [
    { name: 'Nairobi County Rollout', phase: 'Survey', schools: 25, progress: 60 },
    { name: 'Mombasa Pilot', phase: 'Installation', schools: 10, progress: 80 },
    { name: 'Kisumu Expansion', phase: 'Planning', schools: 15, progress: 20 },
  ];

  return (
    <div className="space-y-6">
      {/* Welcome Header - Purple/Violet Theme (Demo Team - Project Focused) */}
      <div className="relative overflow-hidden rounded-2xl bg-gradient-to-r from-violet-600 via-purple-600 to-fuchsia-600 p-6 text-white shadow-lg">
        <div className="absolute top-0 right-0 w-64 h-64 bg-white/5 rounded-full -translate-y-32 translate-x-32" />
        <div className="absolute bottom-0 left-0 w-48 h-48 bg-white/5 rounded-full translate-y-24 -translate-x-24" />

        <div className="relative">
          <div className="flex items-center gap-3 mb-2">
            <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-white/20 backdrop-blur">
              <Target className="h-6 w-6" />
            </div>
            <div>
              <h1 className="text-2xl md:text-3xl font-bold">Welcome, {displayName}!</h1>
              <p className="text-violet-100">Demo Team Dashboard</p>
            </div>
          </div>

          {/* Summary Stats */}
          <div className="mt-6 grid grid-cols-2 md:grid-cols-6 gap-4">
            <div className="bg-white/10 backdrop-blur rounded-xl p-3">
              <p className="text-2xl font-bold">{stats.scheduledDemos}</p>
              <p className="text-violet-200 text-sm">Scheduled Demos</p>
            </div>
            <div className="bg-white/10 backdrop-blur rounded-xl p-3">
              <p className="text-2xl font-bold">{stats.completedThisWeek}</p>
              <p className="text-violet-200 text-sm">Done This Week</p>
            </div>
            <div className="bg-white/10 backdrop-blur rounded-xl p-3">
              <p className="text-2xl font-bold">{stats.pendingSurveys}</p>
              <p className="text-violet-200 text-sm">Pending Surveys</p>
            </div>
            <div className="bg-white/10 backdrop-blur rounded-xl p-3">
              <p className="text-2xl font-bold">{stats.activeProjects}</p>
              <p className="text-violet-200 text-sm">Active Projects</p>
            </div>
            <div className="bg-white/10 backdrop-blur rounded-xl p-3">
              <p className="text-2xl font-bold">{stats.schoolsVisited}</p>
              <p className="text-violet-200 text-sm">Schools Visited</p>
            </div>
            <div className="bg-white/10 backdrop-blur rounded-xl p-3">
              <p className="text-2xl font-bold">{stats.avgDemoRating}</p>
              <p className="text-violet-200 text-sm">Avg Rating</p>
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
                <p className="text-sm text-violet-600">View schedule</p>
              </div>
              <ArrowRight className="h-5 w-5 text-violet-400" />
            </CardContent>
          </Card>
        </Link>

        <Link to="/projects">
          <Card className="hover:shadow-md transition-shadow cursor-pointer border-purple-200 bg-purple-50">
            <CardContent className="p-4 flex items-center gap-3">
              <div className="h-10 w-10 rounded-lg bg-purple-100 flex items-center justify-center">
                <Layers className="h-5 w-5 text-purple-600" />
              </div>
              <div className="flex-1">
                <p className="font-medium text-purple-900">Projects</p>
                <p className="text-sm text-purple-600">Manage rollouts</p>
              </div>
              <ArrowRight className="h-5 w-5 text-purple-400" />
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
                <p className="font-medium text-green-900">Schools</p>
                <p className="text-sm text-green-600">School directory</p>
              </div>
              <ArrowRight className="h-5 w-5 text-green-400" />
            </CardContent>
          </Card>
        </Link>

        <Link to="/messages">
          <Card className="hover:shadow-md transition-shadow cursor-pointer border-blue-200 bg-blue-50">
            <CardContent className="p-4 flex items-center gap-3">
              <div className="h-10 w-10 rounded-lg bg-blue-100 flex items-center justify-center">
                <Users className="h-5 w-5 text-blue-600" />
              </div>
              <div className="flex-1">
                <p className="font-medium text-blue-900">Messages</p>
                <p className="text-sm text-blue-600">Team communication</p>
              </div>
              <ArrowRight className="h-5 w-5 text-blue-400" />
            </CardContent>
          </Card>
        </Link>
      </div>

      {/* Main Content Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Today's Schedule */}
        <Card>
          <CardHeader className="pb-3">
            <div className="flex items-center justify-between">
              <CardTitle className="text-lg flex items-center gap-2">
                <Calendar className="h-5 w-5 text-violet-600" />
                Today's Schedule
              </CardTitle>
              <Link to="/demo-pipeline">
                <Button variant="ghost" size="sm">View All</Button>
              </Link>
            </div>
          </CardHeader>
          <CardContent>
            <div className="space-y-3">
              {todaySchedule.map((item) => (
                <div
                  key={item.id}
                  className="flex items-center justify-between p-3 rounded-lg bg-gray-50 hover:bg-gray-100 transition-colors"
                >
                  <div className="flex items-center gap-3">
                    <div className="h-10 w-10 rounded-lg bg-violet-100 flex items-center justify-center">
                      {item.type === 'Survey' ? (
                        <ClipboardList className="h-5 w-5 text-violet-600" />
                      ) : (
                        <Camera className="h-5 w-5 text-violet-600" />
                      )}
                    </div>
                    <div>
                      <p className="font-medium text-gray-900">{item.school}</p>
                      <div className="flex items-center gap-2 text-sm text-gray-500">
                        <Clock className="h-3 w-3" />
                        {item.time}
                        <span className="mx-1">•</span>
                        <MapPin className="h-3 w-3" />
                        {item.location}
                      </div>
                    </div>
                  </div>
                  <Badge variant="outline" className="bg-violet-50 text-violet-700 border-violet-200">
                    {item.type}
                  </Badge>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>

        {/* Pending Tasks */}
        <Card>
          <CardHeader className="pb-3">
            <div className="flex items-center justify-between">
              <CardTitle className="text-lg flex items-center gap-2">
                <ClipboardList className="h-5 w-5 text-amber-600" />
                Pending Tasks
              </CardTitle>
            </div>
          </CardHeader>
          <CardContent>
            <div className="space-y-3">
              {pendingTasks.map((task, idx) => (
                <div
                  key={idx}
                  className="flex items-center justify-between p-3 rounded-lg bg-gray-50"
                >
                  <div className="flex items-center gap-3">
                    <div className={`h-2 w-2 rounded-full ${
                      task.priority === 'high' ? 'bg-red-500' :
                      task.priority === 'medium' ? 'bg-amber-500' : 'bg-green-500'
                    }`} />
                    <div>
                      <p className="text-sm text-gray-900">
                        {task.type === 'survey' && 'Complete survey for '}
                        {task.type === 'report' && 'Submit demo report for '}
                        {task.type === 'boq' && 'Review BOQ for '}
                        <span className="font-medium">{task.school}</span>
                      </p>
                      <p className="text-xs text-gray-500">Due: {task.dueDate}</p>
                    </div>
                  </div>
                  <Button variant="outline" size="sm">Start</Button>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Active Projects */}
      <Card>
        <CardHeader className="pb-3">
          <div className="flex items-center justify-between">
            <CardTitle className="text-lg flex items-center gap-2">
              <Layers className="h-5 w-5 text-purple-600" />
              Active Projects
            </CardTitle>
            <Link to="/projects">
              <Button variant="ghost" size="sm">View All</Button>
            </Link>
          </div>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            {recentProjects.map((project, idx) => (
              <div key={idx} className="p-4 rounded-xl bg-gray-50 border border-gray-100">
                <div className="flex items-center justify-between mb-3">
                  <h4 className="font-medium text-gray-900">{project.name}</h4>
                  <Badge variant="outline" className="text-xs">
                    {project.phase}
                  </Badge>
                </div>
                <div className="flex items-center gap-2 text-sm text-gray-600 mb-3">
                  <School className="h-4 w-4" />
                  {project.schools} schools
                </div>
                <div className="relative pt-1">
                  <div className="flex items-center justify-between mb-1">
                    <span className="text-xs text-gray-500">Progress</span>
                    <span className="text-xs font-medium text-gray-700">{project.progress}%</span>
                  </div>
                  <div className="h-2 bg-gray-200 rounded-full overflow-hidden">
                    <div
                      className="h-full bg-gradient-to-r from-violet-500 to-purple-500 rounded-full transition-all"
                      style={{ width: `${project.progress}%` }}
                    />
                  </div>
                </div>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
