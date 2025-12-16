import { Plus, ClipboardList } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { formatDate } from '@/lib/utils';
import type { SiteSurvey } from '@/types';

interface SurveysListProps {
  surveys: SiteSurvey[];
}

export function SurveysList({ surveys }: SurveysListProps) {
  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h3 className="font-medium text-gray-900">Site Surveys</h3>
        <Button size="sm">
          <Plus className="h-4 w-4" />
          New Survey
        </Button>
      </div>
      {surveys.length === 0 ? (
        <div className="text-center py-8 text-gray-500">
          <ClipboardList className="h-12 w-12 mx-auto mb-2 text-gray-300" />
          <p>No surveys completed yet</p>
        </div>
      ) : (
        <div className="space-y-2">
          {surveys.map((survey) => (
            <div key={survey.id} className="p-3 bg-gray-50 rounded-lg">
              <div className="flex items-center justify-between mb-1">
                <span className="font-medium text-gray-900">
                  Survey {survey.id.slice(0, 8)}
                </span>
                <Badge
                  variant={
                    survey.status === 'approved'
                      ? 'success'
                      : survey.status === 'submitted'
                      ? 'default'
                      : 'outline'
                  }
                >
                  {survey.status}
                </Badge>
              </div>
              {survey.summary && (
                <p className="text-sm text-gray-500">{survey.summary}</p>
              )}
              {survey.conductedAt && (
                <p className="text-xs text-gray-400 mt-1">
                  Conducted: {formatDate(survey.conductedAt)}
                </p>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
