import { Sparkles, AlertCircle, Lightbulb, MessageCircle, Loader2 } from 'lucide-react';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { Badge } from '@/components/ui/badge';
import { cn } from '@/lib/utils';
import type { EdTechProfile, FollowUpQuestion, AIRecommendation } from '@/types/edtech';

interface AIFollowUpStepProps {
  profile: EdTechProfile;
  followUpResponses: Record<string, string>;
  onResponseChange: (questionId: string, response: string) => void;
  isGenerating?: boolean;
}

const priorityColors: Record<string, string> = {
  high: 'bg-red-100 text-red-800 border-red-200',
  medium: 'bg-yellow-100 text-yellow-800 border-yellow-200',
  low: 'bg-green-100 text-green-800 border-green-200',
};

const categoryIcons: Record<string, string> = {
  infrastructure: '🏗️',
  training: '📚',
  software: '💻',
  security: '🔒',
  support: '🤝',
};

export function AIFollowUpStep({
  profile,
  followUpResponses,
  onResponseChange,
  isGenerating,
}: AIFollowUpStepProps) {
  if (isGenerating) {
    return (
      <div className="flex flex-col items-center justify-center py-16">
        <div className="relative">
          <div className="absolute inset-0 bg-purple-400 rounded-full blur-xl opacity-30 animate-pulse" />
          <Sparkles className="h-16 w-16 text-purple-600 relative z-10 animate-pulse" />
        </div>
        <h3 className="mt-6 text-lg font-medium text-gray-900">Analyzing Your Profile</h3>
        <p className="mt-2 text-sm text-gray-500 text-center max-w-sm">
          Our AI is reviewing your assessment to generate a personalized summary and recommendations...
        </p>
        <Loader2 className="mt-4 h-6 w-6 text-purple-600 animate-spin" />
      </div>
    );
  }

  if (!profile.aiSummary && !profile.followUpQuestions?.length) {
    return (
      <div className="flex flex-col items-center justify-center py-16">
        <AlertCircle className="h-12 w-12 text-amber-500" />
        <h3 className="mt-4 text-lg font-medium text-gray-900">AI Analysis Not Available</h3>
        <p className="mt-2 text-sm text-gray-500 text-center max-w-sm">
          Click "Generate AI Analysis" to receive personalized recommendations based on your assessment.
        </p>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* AI Summary */}
      {profile.aiSummary && (
        <div className="rounded-lg border border-purple-200 bg-gradient-to-br from-purple-50 to-indigo-50 p-4">
          <div className="flex items-center gap-2 mb-3">
            <Sparkles className="h-5 w-5 text-purple-600" />
            <h3 className="font-medium text-gray-900">Your EdTech Profile Summary</h3>
          </div>
          <p className="text-sm text-gray-700 leading-relaxed">{profile.aiSummary}</p>
        </div>
      )}

      {/* AI Recommendations */}
      {profile.aiRecommendations && profile.aiRecommendations.length > 0 && (
        <div className="rounded-lg border border-indigo-100 bg-indigo-50/50 p-4">
          <div className="flex items-center gap-2 mb-4">
            <Lightbulb className="h-5 w-5 text-indigo-600" />
            <h3 className="font-medium text-gray-900">Recommendations</h3>
          </div>

          <div className="space-y-3">
            {profile.aiRecommendations.map((rec, index) => (
              <RecommendationCard key={index} recommendation={rec} />
            ))}
          </div>
        </div>
      )}

      {/* Follow-up Questions */}
      {profile.followUpQuestions && profile.followUpQuestions.length > 0 && (
        <div className="rounded-lg border border-blue-100 bg-blue-50/50 p-4">
          <div className="flex items-center gap-2 mb-4">
            <MessageCircle className="h-5 w-5 text-blue-600" />
            <h3 className="font-medium text-gray-900">Follow-up Questions</h3>
          </div>
          <p className="text-sm text-gray-600 mb-4">
            Please answer these questions to help us better understand your needs:
          </p>

          <div className="space-y-4">
            {profile.followUpQuestions.map((question) => (
              <FollowUpQuestionCard
                key={question.id}
                question={question}
                response={followUpResponses[question.id] || ''}
                onResponseChange={(response) => onResponseChange(question.id, response)}
              />
            ))}
          </div>
        </div>
      )}

      {/* Completion Message */}
      <div className="rounded-lg border border-green-200 bg-green-50 p-4">
        <div className="flex items-start gap-3">
          <div className="flex-shrink-0 mt-0.5">
            <svg className="h-5 w-5 text-green-600" fill="currentColor" viewBox="0 0 20 20">
              <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
            </svg>
          </div>
          <div>
            <h4 className="text-sm font-medium text-green-800">Almost Done!</h4>
            <p className="mt-1 text-sm text-green-700">
              {profile.followUpQuestions?.length
                ? 'Answer the follow-up questions above, then click "Complete Assessment" to save your EdTech profile.'
                : 'Click "Complete Assessment" to save your EdTech profile. You can update it anytime from your dashboard.'}
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}

function RecommendationCard({ recommendation }: { recommendation: AIRecommendation }) {
  return (
    <div className="bg-white rounded-lg border border-gray-100 p-3 shadow-sm">
      <div className="flex items-start justify-between gap-2">
        <div className="flex items-center gap-2">
          <span className="text-lg">{categoryIcons[recommendation.category] || '💡'}</span>
          <h4 className="font-medium text-gray-900 text-sm">{recommendation.title}</h4>
        </div>
        <Badge className={cn('text-xs', priorityColors[recommendation.priority])}>
          {recommendation.priority}
        </Badge>
      </div>
      <p className="mt-2 text-sm text-gray-600 pl-7">{recommendation.description}</p>
    </div>
  );
}

function FollowUpQuestionCard({
  question,
  response,
  onResponseChange,
}: {
  question: FollowUpQuestion;
  response: string;
  onResponseChange: (response: string) => void;
}) {
  return (
    <div className="bg-white rounded-lg border border-blue-100 p-3">
      <Label htmlFor={`question-${question.id}`} className="text-sm font-medium text-gray-900">
        {question.question}
      </Label>
      {question.context && (
        <p className="text-xs text-gray-500 mt-1 mb-2">{question.context}</p>
      )}
      <Textarea
        id={`question-${question.id}`}
        value={response}
        onChange={(e) => onResponseChange(e.target.value)}
        placeholder="Your answer..."
        rows={2}
        className="mt-2"
      />
    </div>
  );
}
