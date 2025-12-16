import { Star } from 'lucide-react';
import { Textarea } from '@/components/ui/textarea';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Checkbox } from '@/components/ui/checkbox';
import { cn } from '@/lib/utils';
import { QuizQuestion } from './QuizQuestion';
import type { PainPointsStepData, EdTechFormOptions } from '@/types/edtech';

interface PainPointsStepProps {
  data: PainPointsStepData;
  options: EdTechFormOptions | undefined;
  onChange: (data: PainPointsStepData) => void;
  currentQuestion: number;
  totalQuestions: number;
}

export function PainPointsStep({
  data,
  options,
  onChange,
  currentQuestion,
  totalQuestions,
}: PainPointsStepProps) {
  const updateField = <K extends keyof PainPointsStepData>(
    field: K,
    value: PainPointsStepData[K]
  ) => {
    onChange({ ...data, [field]: value });
  };

  const togglePainPoint = (point: string) => {
    const current = data.painPoints || [];
    if (current.includes(point)) {
      updateField('painPoints', current.filter((p) => p !== point));
      // Also remove from biggest challenges if present
      const challenges = data.biggestChallenges || [];
      if (challenges.includes(point)) {
        updateField('biggestChallenges', challenges.filter((c) => c !== point));
      }
    } else {
      updateField('painPoints', [...current, point]);
    }
  };

  const toggleChallenge = (challenge: string) => {
    const current = data.biggestChallenges || [];
    if (current.includes(challenge)) {
      updateField('biggestChallenges', current.filter((c) => c !== challenge));
    } else if (current.length < 3) {
      updateField('biggestChallenges', [...current, challenge]);
    }
  };

  const renderQuestion = () => {
    switch (currentQuestion) {
      case 1:
        return (
          <QuizQuestion
            questionNumber={1}
            totalQuestions={totalQuestions}
            question="What technology challenges does your school currently face?"
            description="Select all the challenges that apply to your school."
          >
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
              {(options?.painPoints || []).map((point) => (
                <label
                  key={point}
                  className={cn(
                    'flex items-start gap-3 p-3 rounded-lg border cursor-pointer transition-colors',
                    data.painPoints?.includes(point)
                      ? 'border-amber-400 bg-amber-50'
                      : 'border-gray-200 bg-white hover:border-amber-300 hover:bg-amber-50/50'
                  )}
                >
                  <Checkbox
                    id={`pain-${point}`}
                    checked={data.painPoints?.includes(point)}
                    onCheckedChange={() => togglePainPoint(point)}
                    className="mt-0.5"
                  />
                  <span className="text-sm leading-tight">{point}</span>
                </label>
              ))}
            </div>
          </QuizQuestion>
        );

      case 2:
        return (
          <QuizQuestion
            questionNumber={2}
            totalQuestions={totalQuestions}
            question="Which are your TOP 3 biggest challenges?"
            description={`Select up to 3 challenges from the ones you selected. ${data.biggestChallenges?.length || 0}/3 selected.`}
          >
            {(data.painPoints?.length || 0) === 0 ? (
              <div className="text-center py-8 text-gray-500">
                <p>Please go back and select your challenges first.</p>
              </div>
            ) : (
              <div className="space-y-2">
                {(data.painPoints || []).map((point) => {
                  const isSelected = data.biggestChallenges?.includes(point);
                  const rank = isSelected ? (data.biggestChallenges?.indexOf(point) ?? -1) + 1 : null;

                  return (
                    <button
                      key={point}
                      type="button"
                      onClick={() => toggleChallenge(point)}
                      disabled={!isSelected && (data.biggestChallenges?.length || 0) >= 3}
                      className={cn(
                        'w-full flex items-center gap-3 p-3 rounded-lg border text-left transition-colors',
                        isSelected
                          ? 'border-red-400 bg-red-50'
                          : (data.biggestChallenges?.length || 0) >= 3
                          ? 'border-gray-100 bg-gray-50 opacity-50 cursor-not-allowed'
                          : 'border-gray-200 bg-white hover:border-red-300 hover:bg-red-50/50'
                      )}
                    >
                      <span
                        className={cn(
                          'w-8 h-8 rounded-full flex items-center justify-center text-sm font-bold',
                          isSelected
                            ? 'bg-red-600 text-white'
                            : 'bg-gray-100 text-gray-400'
                        )}
                      >
                        {rank || '-'}
                      </span>
                      <span className="text-sm">{point}</span>
                    </button>
                  );
                })}
              </div>
            )}
          </QuizQuestion>
        );

      case 3:
        return (
          <QuizQuestion
            questionNumber={3}
            totalQuestions={totalQuestions}
            question="How satisfied are you with your current tech support?"
            description="Rate your satisfaction from 1 (very dissatisfied) to 5 (very satisfied)."
          >
            <div className="flex justify-center gap-4">
              {[1, 2, 3, 4, 5].map((rating) => (
                <button
                  key={rating}
                  type="button"
                  onClick={() => updateField('supportSatisfaction', rating)}
                  className={cn(
                    'w-14 h-14 rounded-full flex items-center justify-center transition-all',
                    data.supportSatisfaction === rating
                      ? 'bg-yellow-500 text-white scale-110 shadow-lg'
                      : 'bg-white border-2 border-gray-200 hover:border-yellow-400 hover:scale-105'
                  )}
                >
                  <Star
                    className={cn(
                      'w-7 h-7',
                      data.supportSatisfaction && rating <= data.supportSatisfaction
                        ? 'fill-current'
                        : ''
                    )}
                  />
                </button>
              ))}
            </div>
            <div className="flex justify-between text-sm text-gray-500 mt-4 max-w-sm mx-auto">
              <span>Very Dissatisfied</span>
              <span>Very Satisfied</span>
            </div>
          </QuizQuestion>
        );

      case 4:
        return (
          <QuizQuestion
            questionNumber={4}
            totalQuestions={totalQuestions}
            question="How often does your school need tech support?"
            description="Select the frequency that best describes your support needs."
          >
            <Select
              value={data.supportFrequency}
              onValueChange={(v) => updateField('supportFrequency', v)}
            >
              <SelectTrigger className="w-full max-w-md">
                <SelectValue placeholder="Select frequency" />
              </SelectTrigger>
              <SelectContent>
                {(options?.supportFrequency || []).map((freq) => (
                  <SelectItem key={freq} value={freq}>
                    {freq.charAt(0).toUpperCase() + freq.slice(1)}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </QuizQuestion>
        );

      case 5:
        return (
          <QuizQuestion
            questionNumber={5}
            totalQuestions={totalQuestions}
            question="How long does it typically take to resolve tech issues?"
            description="Select the average time it takes to resolve technical problems."
          >
            <Select
              value={data.avgResolutionTime}
              onValueChange={(v) => updateField('avgResolutionTime', v)}
            >
              <SelectTrigger className="w-full max-w-md">
                <SelectValue placeholder="Select resolution time" />
              </SelectTrigger>
              <SelectContent>
                {(options?.resolutionTime || []).map((time) => (
                  <SelectItem key={time} value={time}>
                    {time.replace(/_/g, ' ').replace(/\b\w/g, (l) => l.toUpperCase())}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </QuizQuestion>
        );

      case 6:
        return (
          <QuizQuestion
            questionNumber={6}
            totalQuestions={totalQuestions}
            question="What is your biggest frustration with technology at your school?"
            description="Share your main pain point or frustration in your own words."
          >
            <Textarea
              id="biggestFrustration"
              value={data.biggestFrustration || ''}
              onChange={(e) => updateField('biggestFrustration', e.target.value)}
              placeholder="Describe your biggest frustration..."
              rows={4}
              className="w-full"
            />
          </QuizQuestion>
        );

      case 7:
        return (
          <QuizQuestion
            questionNumber={7}
            totalQuestions={totalQuestions}
            question="If you could have one tech-related wish granted, what would it be?"
            description="What single improvement would make the biggest difference for your school?"
          >
            <Textarea
              id="wishList"
              value={data.wishList || ''}
              onChange={(e) => updateField('wishList', e.target.value)}
              placeholder="What would make the biggest difference for your school?"
              rows={4}
              className="w-full"
            />
          </QuizQuestion>
        );

      default:
        return null;
    }
  };

  return renderQuestion();
}
