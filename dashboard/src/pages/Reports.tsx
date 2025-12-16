import { Link } from 'react-router-dom';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import {
  AlertTriangle,
  Package,
  School,
  LayoutDashboard,
  ArrowRight,
  Wrench,
} from 'lucide-react';

const reportCategories = [
  {
    title: 'Work Orders Report',
    description: 'Track work order status, completion rates, and repair timelines',
    icon: Wrench,
    href: '/reports/work-orders',
    color: 'text-blue-600',
    bgColor: 'bg-blue-100',
  },
  {
    title: 'Incidents Report',
    description: 'Analyze incidents, SLA compliance, and resolution times',
    icon: AlertTriangle,
    href: '/reports/incidents',
    color: 'text-amber-600',
    bgColor: 'bg-amber-100',
  },
  {
    title: 'Inventory Report',
    description: 'Monitor stock levels, low stock alerts, and parts distribution',
    icon: Package,
    href: '/reports/inventory',
    color: 'text-purple-600',
    bgColor: 'bg-purple-100',
  },
  {
    title: 'Schools Report',
    description: 'View schools, device counts, and incident history by location',
    icon: School,
    href: '/reports/schools',
    color: 'text-green-600',
    bgColor: 'bg-green-100',
  },
  {
    title: 'Executive Dashboard',
    description: 'High-level KPIs and metrics across all domains',
    icon: LayoutDashboard,
    href: '/reports/executive',
    color: 'text-teal-600',
    bgColor: 'bg-teal-100',
  },
];

export function Reports() {
  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-gray-900">Reports</h1>
        <p className="text-sm text-gray-500 mt-1">
          Access detailed reports and analytics across all platform domains
        </p>
      </div>

      <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
        {reportCategories.map((report) => {
          const Icon = report.icon;
          return (
            <Link key={report.href} to={report.href}>
              <Card className="h-full hover:shadow-md transition-shadow cursor-pointer">
                <CardHeader className="flex flex-row items-start gap-4 pb-2">
                  <div className={`p-3 rounded-lg ${report.bgColor}`}>
                    <Icon className={`h-6 w-6 ${report.color}`} />
                  </div>
                  <div className="flex-1">
                    <CardTitle className="text-lg">{report.title}</CardTitle>
                  </div>
                </CardHeader>
                <CardContent>
                  <p className="text-sm text-gray-500">{report.description}</p>
                  <div className="mt-4 flex items-center text-sm font-medium text-gray-900">
                    View Report
                    <ArrowRight className="ml-2 h-4 w-4" />
                  </div>
                </CardContent>
              </Card>
            </Link>
          );
        })}
      </div>
    </div>
  );
}
