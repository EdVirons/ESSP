import { Card, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { formatDistanceToNow } from 'date-fns';
import {
  BookOpen,
  AlertTriangle,
  Bug,
  CheckSquare,
  FileText,
  User,
  Calendar,
} from 'lucide-react';
import type { KBArticle, KBContentType } from '@/types';
import { contentTypeLabels, moduleLabels, statusLabels } from '@/types';

interface KBArticleListProps {
  articles: KBArticle[];
  isLoading: boolean;
  onArticleClick: (article: KBArticle) => void;
}

const contentTypeIcons: Record<KBContentType, React.ElementType> = {
  runbook: BookOpen,
  troubleshooting: AlertTriangle,
  kedb: Bug,
  checklist: CheckSquare,
  sop: FileText,
};

const statusColors: Record<string, string> = {
  draft: 'bg-amber-100 text-amber-800',
  published: 'bg-green-100 text-green-800',
  archived: 'bg-gray-100 text-gray-800',
};

function LoadingSkeleton() {
  return (
    <div className="space-y-4">
      {[1, 2, 3].map((i) => (
        <Card key={i}>
          <CardContent className="p-4">
            <div className="h-6 w-3/4 bg-gray-200 rounded animate-pulse mb-2" />
            <div className="h-4 w-full bg-gray-200 rounded animate-pulse mb-4" />
            <div className="flex gap-2">
              <div className="h-5 w-20 bg-gray-200 rounded animate-pulse" />
              <div className="h-5 w-16 bg-gray-200 rounded animate-pulse" />
            </div>
          </CardContent>
        </Card>
      ))}
    </div>
  );
}

export function KBArticleList({ articles, isLoading, onArticleClick }: KBArticleListProps) {
  if (isLoading) {
    return <LoadingSkeleton />;
  }

  if (articles.length === 0) {
    return (
      <Card>
        <CardContent className="flex flex-col items-center justify-center py-12">
          <BookOpen className="h-12 w-12 text-gray-300 mb-4" />
          <h3 className="text-lg font-medium text-gray-900 mb-1">No articles found</h3>
          <p className="text-sm text-gray-500">
            Try adjusting your filters or create a new article.
          </p>
        </CardContent>
      </Card>
    );
  }

  return (
    <div className="space-y-4">
      {articles.map((article) => {
        const Icon = contentTypeIcons[article.contentType] || FileText;
        return (
          <Card
            key={article.id}
            className="cursor-pointer hover:shadow-md transition-shadow"
            onClick={() => onArticleClick(article)}
          >
            <CardContent className="p-4">
              <div className="flex items-start gap-4">
                <div className="rounded-lg bg-emerald-100 p-2.5">
                  <Icon className="h-5 w-5 text-emerald-600" />
                </div>
                <div className="flex-1 min-w-0">
                  <div className="flex items-center gap-2 mb-1">
                    <h3 className="font-semibold text-gray-900 truncate">{article.title}</h3>
                    <Badge className={statusColors[article.status]} variant="secondary">
                      {statusLabels[article.status]}
                    </Badge>
                  </div>
                  <p className="text-sm text-gray-500 line-clamp-2 mb-3">
                    {article.summary || 'No summary available'}
                  </p>
                  <div className="flex flex-wrap items-center gap-4 text-xs text-gray-400">
                    <Badge variant="outline">{contentTypeLabels[article.contentType]}</Badge>
                    <Badge variant="outline">{moduleLabels[article.module]}</Badge>
                    <span className="flex items-center gap-1">
                      <User className="h-3 w-3" />
                      {article.updatedByName || 'Unknown'}
                    </span>
                    <span className="flex items-center gap-1">
                      <Calendar className="h-3 w-3" />
                      {formatDistanceToNow(new Date(article.updatedAt), { addSuffix: true })}
                    </span>
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>
        );
      })}
    </div>
  );
}
