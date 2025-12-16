import { Input } from '@/components/ui/input';
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
import type { GoalsStepData, EdTechFormOptions } from '@/types/edtech';

interface GoalsStepProps {
  data: GoalsStepData;
  options: EdTechFormOptions | undefined;
  onChange: (data: GoalsStepData) => void;
  currentQuestion: number;
  totalQuestions: number;
}

export function GoalsStep({
  data,
  options,
  onChange,
  currentQuestion,
  totalQuestions,
}: GoalsStepProps) {
  const updateField = <K extends keyof GoalsStepData>(
    field: K,
    value: GoalsStepData[K]
  ) => {
    onChange({ ...data, [field]: value });
  };

  const toggleGoal = (goal: string) => {
    const current = data.strategicGoals || [];
    const currentRanking = data.priorityRanking || [];

    if (current.includes(goal)) {
      updateField('strategicGoals', current.filter((g) => g !== goal));
      updateField('priorityRanking', currentRanking.filter((g) => g !== goal));
    } else {
      updateField('strategicGoals', [...current, goal]);
      updateField('priorityRanking', [...currentRanking, goal]);
    }
  };

  const moveGoalUp = (index: number) => {
    if (index <= 0) return;
    const newRanking = [...(data.priorityRanking || [])];
    [newRanking[index - 1], newRanking[index]] = [newRanking[index], newRanking[index - 1]];
    updateField('priorityRanking', newRanking);
  };

  const moveGoalDown = (index: number) => {
    const ranking = data.priorityRanking || [];
    if (index >= ranking.length - 1) return;
    const newRanking = [...ranking];
    [newRanking[index], newRanking[index + 1]] = [newRanking[index + 1], newRanking[index]];
    updateField('priorityRanking', newRanking);
  };

  const handleDecisionMakersChange = (value: string) => {
    const makers = value.split(',').map((s) => s.trim()).filter(Boolean);
    updateField('decisionMakers', makers);
  };

  const renderQuestion = () => {
    switch (currentQuestion) {
      case 1:
        return (
          <QuizQuestion
            questionNumber={1}
            totalQuestions={totalQuestions}
            question="What EdTech initiatives is your school interested in pursuing?"
            description="Select all the goals and initiatives your school wants to achieve."
          >
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
              {(options?.strategicGoals || []).map((goal) => (
                <label
                  key={goal}
                  className={cn(
                    'flex items-start gap-3 p-3 rounded-lg border cursor-pointer transition-colors',
                    data.strategicGoals?.includes(goal)
                      ? 'border-emerald-400 bg-emerald-50'
                      : 'border-gray-200 bg-white hover:border-emerald-300 hover:bg-emerald-50/50'
                  )}
                >
                  <Checkbox
                    id={`goal-${goal}`}
                    checked={data.strategicGoals?.includes(goal)}
                    onCheckedChange={() => toggleGoal(goal)}
                    className="mt-0.5"
                  />
                  <span className="text-sm leading-tight">{goal}</span>
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
            question="Please rank your selected goals by priority"
            description="Use the arrows to reorder your goals from highest to lowest priority."
          >
            {(data.strategicGoals?.length || 0) <= 1 ? (
              <div className="text-center py-8 text-gray-500">
                <p>
                  {(data.strategicGoals?.length || 0) === 0
                    ? 'Please go back and select your goals first.'
                    : 'Only one goal selected - no ranking needed. Click Next to continue.'}
                </p>
              </div>
            ) : (
              <div className="space-y-2">
                {(data.priorityRanking || []).map((goal, index) => (
                  <div
                    key={goal}
                    className="flex items-center gap-3 p-3 bg-white border border-violet-200 rounded-lg"
                  >
                    <span className="w-8 h-8 rounded-full bg-violet-600 text-white text-sm flex items-center justify-center font-bold">
                      {index + 1}
                    </span>
                    <span className="flex-1 text-sm">{goal}</span>
                    <div className="flex gap-1">
                      <button
                        type="button"
                        onClick={() => moveGoalUp(index)}
                        disabled={index === 0}
                        className={cn(
                          'p-2 rounded-lg hover:bg-violet-100 transition-colors',
                          index === 0 ? 'opacity-30 cursor-not-allowed' : ''
                        )}
                      >
                        <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 15l7-7 7 7" />
                        </svg>
                      </button>
                      <button
                        type="button"
                        onClick={() => moveGoalDown(index)}
                        disabled={index === (data.priorityRanking?.length || 0) - 1}
                        className={cn(
                          'p-2 rounded-lg hover:bg-violet-100 transition-colors',
                          index === (data.priorityRanking?.length || 0) - 1 ? 'opacity-30 cursor-not-allowed' : ''
                        )}
                      >
                        <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                        </svg>
                      </button>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </QuizQuestion>
        );

      case 3:
        return (
          <QuizQuestion
            questionNumber={3}
            totalQuestions={totalQuestions}
            question="What is your available budget for EdTech investments?"
            description="Select the budget range your school has allocated for technology initiatives."
          >
            <Select
              value={data.budgetRange}
              onValueChange={(v) => updateField('budgetRange', v)}
            >
              <SelectTrigger className="w-full max-w-md">
                <SelectValue placeholder="Select budget range" />
              </SelectTrigger>
              <SelectContent>
                {(options?.budgetRange || []).map((budget) => (
                  <SelectItem key={budget} value={budget}>
                    {budget.charAt(0).toUpperCase() + budget.slice(1)}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </QuizQuestion>
        );

      case 4:
        return (
          <QuizQuestion
            questionNumber={4}
            totalQuestions={totalQuestions}
            question="What is your ideal implementation timeline?"
            description="When would you like to start implementing these initiatives?"
          >
            <Select
              value={data.timeline}
              onValueChange={(v) => updateField('timeline', v)}
            >
              <SelectTrigger className="w-full max-w-md">
                <SelectValue placeholder="Select timeline" />
              </SelectTrigger>
              <SelectContent>
                {(options?.timeline || []).map((time) => (
                  <SelectItem key={time} value={time}>
                    {time.replace(/_/g, ' ').replace(/\b\w/g, (l) => l.toUpperCase())}
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
            question="Describe any specific expansion or improvement plans"
            description="Share details about your planned technology expansions or improvements."
          >
            <Textarea
              id="expansionPlans"
              value={data.expansionPlans || ''}
              onChange={(e) => updateField('expansionPlans', e.target.value)}
              placeholder="e.g., We plan to expand the computer lab, implement a 1:1 device program for Grade 5-8, upgrade our network infrastructure..."
              rows={4}
              className="w-full"
            />
          </QuizQuestion>
        );

      case 6:
        return (
          <QuizQuestion
            questionNumber={6}
            totalQuestions={totalQuestions}
            question="Who else is involved in technology decisions at your school?"
            description="List the key decision makers for technology purchases and implementations."
          >
            <div className="max-w-md">
              <Input
                id="decisionMakers"
                value={data.decisionMakers?.join(', ') || ''}
                onChange={(e) => handleDecisionMakersChange(e.target.value)}
                placeholder="e.g., Principal, ICT Coordinator, Board Chair"
              />
              <p className="text-sm text-gray-500 mt-2">
                Separate multiple roles with commas
              </p>
            </div>
          </QuizQuestion>
        );

      default:
        return null;
    }
  };

  return renderQuestion();
}
