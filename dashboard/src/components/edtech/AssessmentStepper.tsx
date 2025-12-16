import { Check, Laptop, AlertTriangle, Target, Sparkles } from 'lucide-react';
import { cn } from '@/lib/utils';

export interface Step {
  id: number;
  title: string;
  description: string;
}

const steps: Step[] = [
  { id: 1, title: 'Infrastructure', description: 'Devices, network, software' },
  { id: 2, title: 'Pain Points', description: 'Challenges & needs' },
  { id: 3, title: 'Goals', description: 'Priorities & timeline' },
  { id: 4, title: 'AI Analysis', description: 'Summary & recommendations' },
];

const stepIcons = [Laptop, AlertTriangle, Target, Sparkles];

interface AssessmentStepperProps {
  currentStep: number;
  onStepClick?: (step: number) => void;
  allowNavigation?: boolean;
}

export function AssessmentStepper({
  currentStep,
  onStepClick,
  allowNavigation = false,
}: AssessmentStepperProps) {
  return (
    <nav aria-label="Progress" className="mb-8">
      <ol className="flex items-center justify-between">
        {steps.map((step, index) => {
          const Icon = stepIcons[index];
          const isCompleted = step.id < currentStep;
          const isCurrent = step.id === currentStep;
          const isClickable = allowNavigation && step.id <= currentStep;

          return (
            <li key={step.id} className="relative flex-1">
              {/* Connector line */}
              {index < steps.length - 1 && (
                <div
                  className={cn(
                    'absolute top-5 left-1/2 w-full h-0.5 -z-10',
                    isCompleted ? 'bg-indigo-600' : 'bg-gray-200'
                  )}
                />
              )}

              <button
                type="button"
                onClick={() => isClickable && onStepClick?.(step.id)}
                disabled={!isClickable}
                className={cn(
                  'group flex flex-col items-center w-full',
                  isClickable ? 'cursor-pointer' : 'cursor-default'
                )}
              >
                {/* Step circle */}
                <span
                  className={cn(
                    'flex h-10 w-10 items-center justify-center rounded-full border-2 transition-colors',
                    isCompleted
                      ? 'border-indigo-600 bg-indigo-600 text-white'
                      : isCurrent
                      ? 'border-indigo-600 bg-white text-indigo-600'
                      : 'border-gray-300 bg-white text-gray-400'
                  )}
                >
                  {isCompleted ? (
                    <Check className="h-5 w-5" />
                  ) : (
                    <Icon className="h-5 w-5" />
                  )}
                </span>

                {/* Step text */}
                <span
                  className={cn(
                    'mt-2 text-sm font-medium',
                    isCurrent ? 'text-indigo-600' : isCompleted ? 'text-gray-900' : 'text-gray-500'
                  )}
                >
                  {step.title}
                </span>
                <span
                  className={cn(
                    'text-xs hidden sm:block',
                    isCurrent ? 'text-indigo-500' : 'text-gray-400'
                  )}
                >
                  {step.description}
                </span>
              </button>
            </li>
          );
        })}
      </ol>
    </nav>
  );
}
