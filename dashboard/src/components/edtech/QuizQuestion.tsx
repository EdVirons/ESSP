import type { ReactNode } from 'react';

interface QuizQuestionProps {
  questionNumber: number;
  totalQuestions: number;
  question: string;
  description?: string;
  children: ReactNode;
}

export function QuizQuestion({
  questionNumber,
  totalQuestions,
  question,
  description,
  children,
}: QuizQuestionProps) {
  const progress = (questionNumber / totalQuestions) * 100;

  return (
    <div className="space-y-6">
      {/* Progress indicator */}
      <div className="space-y-2">
        <div className="flex items-center justify-between text-sm">
          <span className="font-medium text-indigo-600">
            Question {questionNumber} of {totalQuestions}
          </span>
          <span className="text-gray-500">{Math.round(progress)}% complete</span>
        </div>
        <div className="h-2 w-full rounded-full bg-gray-100">
          <div
            className="h-2 rounded-full bg-gradient-to-r from-indigo-500 to-purple-500 transition-all duration-300"
            style={{ width: `${progress}%` }}
          />
        </div>
      </div>

      {/* Question */}
      <div className="rounded-xl border border-indigo-100 bg-gradient-to-br from-indigo-50/50 to-purple-50/50 p-6">
        <h3 className="text-xl font-semibold text-gray-900 mb-2">{question}</h3>
        {description && (
          <p className="text-gray-600 mb-6">{description}</p>
        )}
        <div className="mt-6">{children}</div>
      </div>
    </div>
  );
}
