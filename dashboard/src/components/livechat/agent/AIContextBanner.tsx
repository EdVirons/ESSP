import { Bot, ChevronDown, ChevronUp, AlertTriangle, Info } from 'lucide-react';
import { useState } from 'react';
import { Badge } from '@/components/ui/badge';
import { useAIContext } from '@/hooks/useLivechat';

interface AIContextBannerProps {
  sessionId: string;
}

export function AIContextBanner({ sessionId }: AIContextBannerProps) {
  const [isExpanded, setIsExpanded] = useState(false);
  const { data: context, isLoading } = useAIContext(sessionId);

  if (isLoading || !context) {
    return null;
  }

  const { turnCount, category, severity, escalationReason, summary, collectedInfo } = context;

  const severityColor = {
    low: 'bg-green-100 text-green-700 border-green-200',
    medium: 'bg-amber-100 text-amber-700 border-amber-200',
    high: 'bg-orange-100 text-orange-700 border-orange-200',
    critical: 'bg-red-100 text-red-700 border-red-200',
  }[severity || 'medium'] || 'bg-gray-100 text-gray-700 border-gray-200';

  return (
    <div className="mx-4 mb-4 rounded-lg bg-purple-50 border border-purple-200 overflow-hidden">
      {/* Header */}
      <button
        onClick={() => setIsExpanded(!isExpanded)}
        className="w-full flex items-center justify-between p-3 hover:bg-purple-100 transition-colors"
      >
        <div className="flex items-center gap-2">
          <Bot className="w-5 h-5 text-purple-600" />
          <span className="font-medium text-purple-900">AI Handoff Summary</span>
          <Badge variant="outline" className="bg-purple-100 text-purple-700 border-purple-200">
            {turnCount} turns
          </Badge>
        </div>
        {isExpanded ? (
          <ChevronUp className="w-5 h-5 text-purple-600" />
        ) : (
          <ChevronDown className="w-5 h-5 text-purple-600" />
        )}
      </button>

      {/* Collapsed Preview */}
      {!isExpanded && (
        <div className="px-3 pb-3 flex items-center gap-2 flex-wrap">
          {category && (
            <Badge variant="outline" className="text-xs">
              {category}
            </Badge>
          )}
          {severity && (
            <Badge variant="outline" className={`text-xs ${severityColor}`}>
              {severity}
            </Badge>
          )}
          {escalationReason && (
            <span className="text-sm text-purple-700 truncate">
              {escalationReason}
            </span>
          )}
        </div>
      )}

      {/* Expanded Content */}
      {isExpanded && (
        <div className="px-3 pb-3 space-y-3 border-t border-purple-200 pt-3">
          {/* Category and Severity */}
          <div className="flex items-center gap-2 flex-wrap">
            {category && (
              <div className="flex items-center gap-1">
                <span className="text-xs text-gray-500">Category:</span>
                <Badge variant="outline">{category}</Badge>
              </div>
            )}
            {severity && (
              <div className="flex items-center gap-1">
                <span className="text-xs text-gray-500">Severity:</span>
                <Badge variant="outline" className={severityColor}>{severity}</Badge>
              </div>
            )}
          </div>

          {/* Escalation Reason */}
          {escalationReason && (
            <div className="flex items-start gap-2 p-2 rounded bg-white">
              <AlertTriangle className="w-4 h-4 text-amber-500 flex-shrink-0 mt-0.5" />
              <div>
                <p className="text-xs font-medium text-gray-700">Escalation Reason</p>
                <p className="text-sm text-gray-600">{escalationReason}</p>
              </div>
            </div>
          )}

          {/* Collected Info */}
          {collectedInfo && Object.keys(collectedInfo).length > 0 && (
            <div className="p-2 rounded bg-white">
              <div className="flex items-center gap-1 mb-2">
                <Info className="w-4 h-4 text-blue-500" />
                <p className="text-xs font-medium text-gray-700">Collected Information</p>
              </div>
              <dl className="grid grid-cols-2 gap-2 text-sm">
                {Object.entries(collectedInfo).map(([key, value]) => (
                  <div key={key}>
                    <dt className="text-xs text-gray-500 capitalize">
                      {key.replace(/_/g, ' ')}
                    </dt>
                    <dd className="text-gray-900">{String(value) || '-'}</dd>
                  </div>
                ))}
              </dl>
            </div>
          )}

          {/* AI Summary */}
          {summary && Object.keys(summary).length > 0 && (
            <div className="p-2 rounded bg-white text-sm">
              <p className="text-xs font-medium text-gray-700 mb-1">AI Summary</p>
              <pre className="text-xs text-gray-600 whitespace-pre-wrap">
                {JSON.stringify(summary, null, 2)}
              </pre>
            </div>
          )}
        </div>
      )}
    </div>
  );
}
