import { Card, CardContent } from '@/components/ui/card';
import { FileText, CheckCircle, PenLine, BookOpen } from 'lucide-react';
import type { KBStats as KBStatsType } from '@/types';

interface KBStatsProps {
  stats: KBStatsType | undefined;
  isLoading: boolean;
}

export function KBStats({ stats, isLoading }: KBStatsProps) {
  const cards = [
    {
      label: 'Total Articles',
      value: stats?.total ?? 0,
      icon: FileText,
      color: 'text-blue-600 bg-blue-100',
    },
    {
      label: 'Published',
      value: stats?.published ?? 0,
      icon: CheckCircle,
      color: 'text-green-600 bg-green-100',
    },
    {
      label: 'Drafts',
      value: stats?.draft ?? 0,
      icon: PenLine,
      color: 'text-amber-600 bg-amber-100',
    },
    {
      label: 'Runbooks',
      value: stats?.byContentType?.runbook ?? 0,
      icon: BookOpen,
      color: 'text-purple-600 bg-purple-100',
    },
  ];

  return (
    <div className="grid grid-cols-2 gap-3 md:gap-4 lg:grid-cols-4">
      {cards.map((card) => (
        <Card key={card.label}>
          <CardContent className="flex items-center p-4">
            <div className={`rounded-lg p-3 ${card.color.split(' ')[1]}`}>
              <card.icon className={`h-5 w-5 ${card.color.split(' ')[0]}`} />
            </div>
            <div className="ml-4">
              <p className="text-sm font-medium text-gray-500">{card.label}</p>
              {isLoading ? (
                <div className="h-7 w-16 bg-gray-200 rounded animate-pulse mt-1" />
              ) : (
                <p className="text-2xl font-bold text-gray-900">{card.value}</p>
              )}
            </div>
          </CardContent>
        </Card>
      ))}
    </div>
  );
}
