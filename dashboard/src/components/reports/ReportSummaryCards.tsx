import { Card, CardContent } from '@/components/ui/card';
import type { LucideIcon } from 'lucide-react';

interface SummaryCard {
  title: string;
  value: string | number;
  subtitle?: string;
  icon?: LucideIcon;
  color?: 'default' | 'success' | 'warning' | 'danger' | 'info';
}

interface ReportSummaryCardsProps {
  cards: SummaryCard[];
}

const colorClasses = {
  default: 'text-gray-900',
  success: 'text-green-600',
  warning: 'text-amber-600',
  danger: 'text-red-600',
  info: 'text-blue-600',
};

const bgColorClasses = {
  default: 'bg-gray-100',
  success: 'bg-green-100',
  warning: 'bg-amber-100',
  danger: 'bg-red-100',
  info: 'bg-blue-100',
};

export function ReportSummaryCards({ cards }: ReportSummaryCardsProps) {
  return (
    <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-4">
      {cards.map((card, idx) => {
        const Icon = card.icon;
        const color = card.color || 'default';
        return (
          <Card key={idx}>
            <CardContent className="p-4">
              <div className="flex items-start justify-between">
                <div className="flex-1">
                  <p className="text-xs text-gray-500 font-medium uppercase tracking-wide">
                    {card.title}
                  </p>
                  <p className={`text-2xl font-bold mt-1 ${colorClasses[color]}`}>
                    {typeof card.value === 'number' ? card.value.toLocaleString() : card.value}
                  </p>
                  {card.subtitle && (
                    <p className="text-xs text-gray-400 mt-1">{card.subtitle}</p>
                  )}
                </div>
                {Icon && (
                  <div className={`p-2 rounded-lg ${bgColorClasses[color]}`}>
                    <Icon className={`h-5 w-5 ${colorClasses[color]}`} />
                  </div>
                )}
              </div>
            </CardContent>
          </Card>
        );
      })}
    </div>
  );
}
