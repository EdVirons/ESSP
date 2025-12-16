import { useState } from 'react';
import {
  Sparkles,
  ArrowRight,
  Laptop,
  Target,
  CheckCircle2,
  Pencil,
  Loader2,
} from 'lucide-react';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { useEdTechProfile } from '@/hooks/useEdTechProfile';
import { EdTechAssessmentModal } from './EdTechAssessmentModal';
import { ProfileSummary } from './ProfileSummary';

interface EdTechProfileCardProps {
  schoolId: string;
}

export function EdTechProfileCard({ schoolId }: EdTechProfileCardProps) {
  const [modalOpen, setModalOpen] = useState(false);
  const [showFullProfile, setShowFullProfile] = useState(false);
  const { data, isLoading } = useEdTechProfile(schoolId);

  const profile = data?.profile;
  const hasProfile = !!profile;
  const isComplete = profile?.status === 'completed';

  if (isLoading) {
    return (
      <Card>
        <CardContent className="py-8">
          <div className="flex flex-col items-center justify-center">
            <Loader2 className="h-8 w-8 text-indigo-600 animate-spin" />
            <p className="mt-2 text-sm text-gray-500">Loading EdTech profile...</p>
          </div>
        </CardContent>
      </Card>
    );
  }

  // No profile yet - show CTA to create one
  if (!hasProfile) {
    return (
      <>
        <Card className="overflow-hidden">
          <div className="bg-gradient-to-r from-indigo-500 via-purple-500 to-indigo-600 p-6 text-white">
            <div className="flex items-start gap-4">
              <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-white/20 backdrop-blur">
                <Sparkles className="h-6 w-6" />
              </div>
              <div className="flex-1">
                <h3 className="text-lg font-semibold">Complete Your EdTech Assessment</h3>
                <p className="mt-1 text-sm text-indigo-100">
                  Help us understand your school's technology landscape to provide better support
                  and personalized recommendations.
                </p>
              </div>
            </div>

            <div className="mt-6 flex flex-wrap gap-4">
              <div className="flex items-center gap-2 text-sm text-indigo-100">
                <Laptop className="h-4 w-4" />
                <span>Infrastructure</span>
              </div>
              <div className="flex items-center gap-2 text-sm text-indigo-100">
                <Target className="h-4 w-4" />
                <span>Goals & Priorities</span>
              </div>
              <div className="flex items-center gap-2 text-sm text-indigo-100">
                <Sparkles className="h-4 w-4" />
                <span>AI Recommendations</span>
              </div>
            </div>

            <Button
              onClick={() => setModalOpen(true)}
              className="mt-6 bg-white text-indigo-600 hover:bg-indigo-50"
            >
              Start Assessment
              <ArrowRight className="ml-2 h-4 w-4" />
            </Button>
          </div>
        </Card>

        <EdTechAssessmentModal
          open={modalOpen}
          onClose={() => setModalOpen(false)}
          schoolId={schoolId}
        />
      </>
    );
  }

  // Has profile - show summary or full profile
  return (
    <>
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <CardTitle className="flex items-center gap-2">
              <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-gradient-to-br from-indigo-500 to-purple-600 text-white">
                <Sparkles className="h-4 w-4" />
              </div>
              EdTech Profile
            </CardTitle>
            <div className="flex items-center gap-2">
              {isComplete ? (
                <Badge className="bg-green-100 text-green-800 border-green-200">
                  <CheckCircle2 className="h-3 w-3 mr-1" />
                  Complete
                </Badge>
              ) : (
                <Badge className="bg-amber-100 text-amber-800 border-amber-200">
                  Draft
                </Badge>
              )}
              <Button
                variant="outline"
                size="sm"
                onClick={() => setModalOpen(true)}
              >
                <Pencil className="h-3 w-3 mr-1" />
                {isComplete ? 'Update' : 'Continue'}
              </Button>
            </div>
          </div>
        </CardHeader>

        <CardContent>
          {showFullProfile ? (
            <>
              <ProfileSummary profile={profile} />
              <Button
                variant="ghost"
                size="sm"
                className="mt-4 w-full"
                onClick={() => setShowFullProfile(false)}
              >
                Show Less
              </Button>
            </>
          ) : (
            <>
              {/* Quick Stats */}
              <div className="grid grid-cols-3 gap-3 mb-4">
                <div className="bg-indigo-50 rounded-lg p-3 text-center">
                  <div className="text-2xl font-bold text-indigo-600">
                    {profile.totalDevices || 0}
                  </div>
                  <div className="text-xs text-gray-500">Devices</div>
                </div>
                <div className="bg-amber-50 rounded-lg p-3 text-center">
                  <div className="text-2xl font-bold text-amber-600">
                    {profile.biggestChallenges?.length || 0}
                  </div>
                  <div className="text-xs text-gray-500">Challenges</div>
                </div>
                <div className="bg-emerald-50 rounded-lg p-3 text-center">
                  <div className="text-2xl font-bold text-emerald-600">
                    {profile.strategicGoals?.length || 0}
                  </div>
                  <div className="text-xs text-gray-500">Goals</div>
                </div>
              </div>

              {/* AI Summary snippet */}
              {profile.aiSummary && (
                <div className="bg-purple-50 rounded-lg p-3 mb-4">
                  <div className="flex items-center gap-1 text-xs font-medium text-purple-700 mb-1">
                    <Sparkles className="h-3 w-3" />
                    AI Summary
                  </div>
                  <p className="text-sm text-gray-700 line-clamp-2">{profile.aiSummary}</p>
                </div>
              )}

              {/* Top Priority */}
              {profile.priorityRanking?.[0] && (
                <div className="flex items-center gap-2 text-sm">
                  <Target className="h-4 w-4 text-emerald-600" />
                  <span className="text-gray-500">Top Priority:</span>
                  <span className="font-medium text-gray-900">{profile.priorityRanking[0]}</span>
                </div>
              )}

              <Button
                variant="ghost"
                size="sm"
                className="mt-4 w-full"
                onClick={() => setShowFullProfile(true)}
              >
                View Full Profile
                <ArrowRight className="ml-1 h-3 w-3" />
              </Button>
            </>
          )}
        </CardContent>
      </Card>

      <EdTechAssessmentModal
        open={modalOpen}
        onClose={() => setModalOpen(false)}
        schoolId={schoolId}
      />
    </>
  );
}
