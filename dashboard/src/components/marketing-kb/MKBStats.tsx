import { Card, CardContent } from '@/components/ui/card';
import {
  CheckCircle,
  Clock,
  FileText,
  MessageSquare,
  BookOpen,
  Presentation,
  Shield,
  Calculator,
} from 'lucide-react';
import type { MKBStats as MKBStatsType } from '@/types';

interface MKBStatsProps {
  stats: MKBStatsType | undefined;
  isLoading: boolean;
}

function StatCard({
  label,
  value,
  icon: Icon,
  color,
}: {
  label: string;
  value: number;
  icon: React.ElementType;
  color: string;
}) {
  return (
    <Card>
      <CardContent className="flex items-center gap-3 p-4">
        <div className={`rounded-lg ${color} p-2`}>
          <Icon className="h-5 w-5" />
        </div>
        <div>
          <p className="text-2xl font-bold">{value}</p>
          <p className="text-xs text-gray-500">{label}</p>
        </div>
      </CardContent>
    </Card>
  );
}

function LoadingSkeleton() {
  return (
    <div className="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-4">
      {[1, 2, 3, 4, 5, 6].map((i) => (
        <Card key={i}>
          <CardContent className="flex items-center gap-3 p-4">
            <div className="h-9 w-9 bg-gray-200 rounded-lg animate-pulse" />
            <div>
              <div className="h-6 w-12 bg-gray-200 rounded animate-pulse mb-1" />
              <div className="h-3 w-16 bg-gray-200 rounded animate-pulse" />
            </div>
          </CardContent>
        </Card>
      ))}
    </div>
  );
}

const contentTypeIcons: Record<string, { icon: React.ElementType; color: string }> = {
  messaging: { icon: MessageSquare, color: 'bg-blue-100 text-blue-600' },
  case_study: { icon: BookOpen, color: 'bg-purple-100 text-purple-600' },
  deck: { icon: Presentation, color: 'bg-pink-100 text-pink-600' },
  objection: { icon: Shield, color: 'bg-amber-100 text-amber-600' },
  roi: { icon: Calculator, color: 'bg-green-100 text-green-600' },
};

export function MKBStats({ stats, isLoading }: MKBStatsProps) {
  if (isLoading) {
    return <LoadingSkeleton />;
  }

  if (!stats) {
    return null;
  }

  return (
    <div className="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-4">
      <StatCard
        label="Total"
        value={stats.total}
        icon={FileText}
        color="bg-gray-100 text-gray-600"
      />
      <StatCard
        label="Approved"
        value={stats.approved}
        icon={CheckCircle}
        color="bg-green-100 text-green-600"
      />
      <StatCard
        label="In Review"
        value={stats.inReview}
        icon={Clock}
        color="bg-blue-100 text-blue-600"
      />
      {Object.entries(stats.byContentType || {}).slice(0, 3).map(([type, count]) => {
        const config = contentTypeIcons[type] || { icon: FileText, color: 'bg-gray-100 text-gray-600' };
        return (
          <StatCard
            key={type}
            label={type.replace('_', ' ')}
            value={count}
            icon={config.icon}
            color={config.color}
          />
        );
      })}
    </div>
  );
}
