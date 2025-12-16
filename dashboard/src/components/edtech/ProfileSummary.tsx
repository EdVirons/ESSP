import {
  Laptop,
  AlertTriangle,
  Target,
  Sparkles,
  Calendar,
  CheckCircle2,
} from 'lucide-react';
import { Badge } from '@/components/ui/badge';
import { cn } from '@/lib/utils';
import type { EdTechProfile } from '@/types/edtech';

interface ProfileSummaryProps {
  profile: EdTechProfile;
}

const priorityColors: Record<string, string> = {
  high: 'bg-red-100 text-red-800',
  medium: 'bg-yellow-100 text-yellow-800',
  low: 'bg-green-100 text-green-800',
};

export function ProfileSummary({ profile }: ProfileSummaryProps) {
  return (
    <div className="space-y-6">
      {/* Profile Status */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2">
          {profile.status === 'completed' ? (
            <>
              <CheckCircle2 className="h-5 w-5 text-green-600" />
              <span className="text-sm font-medium text-green-800">Profile Completed</span>
            </>
          ) : (
            <>
              <div className="h-5 w-5 rounded-full border-2 border-amber-500 border-dashed" />
              <span className="text-sm font-medium text-amber-800">Draft</span>
            </>
          )}
        </div>
        {profile.completedAt && (
          <span className="text-xs text-gray-500">
            Last updated: {new Date(profile.completedAt).toLocaleDateString()}
          </span>
        )}
      </div>

      {/* AI Summary */}
      {profile.aiSummary && (
        <div className="rounded-lg bg-gradient-to-br from-purple-50 to-indigo-50 border border-purple-100 p-4">
          <div className="flex items-center gap-2 mb-2">
            <Sparkles className="h-4 w-4 text-purple-600" />
            <span className="text-sm font-medium text-purple-900">AI Summary</span>
          </div>
          <p className="text-sm text-gray-700">{profile.aiSummary}</p>
        </div>
      )}

      {/* Infrastructure Overview */}
      <div className="rounded-lg border border-gray-100 p-4">
        <div className="flex items-center gap-2 mb-3">
          <Laptop className="h-4 w-4 text-indigo-600" />
          <span className="text-sm font-medium">Infrastructure</span>
        </div>
        <div className="grid grid-cols-2 sm:grid-cols-4 gap-3">
          <div className="bg-gray-50 rounded-lg p-2 text-center">
            <div className="text-xl font-bold text-indigo-600">{profile.totalDevices || 0}</div>
            <div className="text-xs text-gray-500">Total Devices</div>
          </div>
          <div className="bg-gray-50 rounded-lg p-2 text-center">
            <div className="text-xl font-bold text-blue-600">{profile.networkQuality || '-'}</div>
            <div className="text-xs text-gray-500">Network Quality</div>
          </div>
          <div className="bg-gray-50 rounded-lg p-2 text-center">
            <div className="text-xl font-bold text-green-600">{profile.itStaffCount || 0}</div>
            <div className="text-xs text-gray-500">IT Staff</div>
          </div>
          <div className="bg-gray-50 rounded-lg p-2 text-center">
            <div className="text-lg font-medium text-purple-600 truncate">
              {profile.lmsPlatform || 'None'}
            </div>
            <div className="text-xs text-gray-500">LMS Platform</div>
          </div>
        </div>

        {/* Device breakdown */}
        {profile.deviceTypes && (
          <div className="mt-3 pt-3 border-t border-gray-100">
            <div className="text-xs text-gray-500 mb-2">Device Breakdown</div>
            <div className="flex flex-wrap gap-2">
              {Object.entries(profile.deviceTypes).map(([type, count]) =>
                count > 0 ? (
                  <Badge key={type} variant="secondary" className="text-xs">
                    {type}: {count}
                  </Badge>
                ) : null
              )}
            </div>
          </div>
        )}
      </div>

      {/* Pain Points */}
      {(profile.biggestChallenges?.length || 0) > 0 && (
        <div className="rounded-lg border border-amber-100 bg-amber-50/50 p-4">
          <div className="flex items-center gap-2 mb-3">
            <AlertTriangle className="h-4 w-4 text-amber-600" />
            <span className="text-sm font-medium">Top Challenges</span>
          </div>
          <div className="space-y-2">
            {profile.biggestChallenges?.slice(0, 3).map((challenge, index) => (
              <div key={index} className="flex items-center gap-2">
                <span className="w-5 h-5 rounded-full bg-amber-200 text-amber-800 text-xs flex items-center justify-center font-medium">
                  {index + 1}
                </span>
                <span className="text-sm text-gray-700">{challenge}</span>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Strategic Goals */}
      {(profile.priorityRanking?.length || 0) > 0 && (
        <div className="rounded-lg border border-emerald-100 bg-emerald-50/50 p-4">
          <div className="flex items-center gap-2 mb-3">
            <Target className="h-4 w-4 text-emerald-600" />
            <span className="text-sm font-medium">Priority Goals</span>
          </div>
          <div className="space-y-2">
            {profile.priorityRanking?.slice(0, 3).map((goal, index) => (
              <div key={index} className="flex items-center gap-2">
                <span className="w-5 h-5 rounded-full bg-emerald-200 text-emerald-800 text-xs flex items-center justify-center font-medium">
                  {index + 1}
                </span>
                <span className="text-sm text-gray-700">{goal}</span>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Budget & Timeline */}
      <div className="flex gap-4">
        {profile.budgetRange && (
          <div className="flex items-center gap-2 text-sm">
            <span className="text-gray-500">Budget:</span>
            <Badge variant="secondary">{profile.budgetRange}</Badge>
          </div>
        )}
        {profile.timeline && (
          <div className="flex items-center gap-2 text-sm">
            <Calendar className="h-4 w-4 text-gray-400" />
            <span className="text-gray-500">Timeline:</span>
            <Badge variant="secondary">{profile.timeline.replace(/_/g, ' ')}</Badge>
          </div>
        )}
      </div>

      {/* AI Recommendations */}
      {profile.aiRecommendations && profile.aiRecommendations.length > 0 && (
        <div className="rounded-lg border border-indigo-100 bg-indigo-50/50 p-4">
          <div className="flex items-center gap-2 mb-3">
            <Sparkles className="h-4 w-4 text-indigo-600" />
            <span className="text-sm font-medium">AI Recommendations</span>
          </div>
          <div className="space-y-2">
            {profile.aiRecommendations.map((rec, index) => (
              <div
                key={index}
                className="flex items-start gap-2 bg-white rounded-lg p-2 border border-indigo-100"
              >
                <Badge className={cn('text-xs mt-0.5', priorityColors[rec.priority])}>
                  {rec.priority}
                </Badge>
                <div className="flex-1 min-w-0">
                  <div className="text-sm font-medium text-gray-900">{rec.title}</div>
                  <div className="text-xs text-gray-500 truncate">{rec.description}</div>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}
