import { BookOpen, Wrench } from 'lucide-react';

export function KnowledgeBase() {
  return (
    <div className="flex min-h-[60vh] flex-col items-center justify-center text-center">
      <div className="space-y-4">
        <div className="mx-auto flex h-20 w-20 items-center justify-center rounded-full bg-emerald-100">
          <BookOpen className="h-10 w-10 text-emerald-600" />
        </div>
        <h1 className="text-2xl font-bold text-gray-900">Knowledge Base</h1>
        <p className="text-gray-500 max-w-md">
          The Knowledge Base is coming soon! This will be your go-to resource for
          troubleshooting guides, repair manuals, and technical documentation.
        </p>
        <div className="flex items-center justify-center gap-2 pt-4 text-sm text-gray-400">
          <Wrench className="h-4 w-4" />
          <span>Under Construction</span>
        </div>
      </div>
    </div>
  );
}
