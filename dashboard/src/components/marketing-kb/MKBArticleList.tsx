import { Card, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { formatDistanceToNow } from 'date-fns';
import {
  MessageSquare,
  BookOpen,
  Presentation,
  Shield,
  Calculator,
  Calendar,
  TrendingUp,
  Plus,
} from 'lucide-react';
import type { MKBArticle, MKBContentType } from '@/types';
import {
  mkbPersonaLabels,
  mkbStatusColors,
  mkbStatusLabels,
} from '@/types';

interface MKBArticleListProps {
  articles: MKBArticle[];
  isLoading: boolean;
  onArticleClick: (article: MKBArticle) => void;
  onAddToKit?: (article: MKBArticle) => void;
  selectedIds?: string[];
}

const contentTypeIcons: Record<MKBContentType, React.ElementType> = {
  messaging: MessageSquare,
  case_study: BookOpen,
  deck: Presentation,
  objection: Shield,
  roi: Calculator,
};

const contentTypeColors: Record<MKBContentType, string> = {
  messaging: 'bg-blue-100 text-blue-600',
  case_study: 'bg-purple-100 text-purple-600',
  deck: 'bg-pink-100 text-pink-600',
  objection: 'bg-amber-100 text-amber-600',
  roi: 'bg-green-100 text-green-600',
};

function LoadingSkeleton() {
  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      {[1, 2, 3, 4, 5, 6].map((i) => (
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

export function MKBArticleList({
  articles,
  isLoading,
  onArticleClick,
  onAddToKit,
  selectedIds = [],
}: MKBArticleListProps) {
  if (isLoading) {
    return <LoadingSkeleton />;
  }

  if (articles.length === 0) {
    return (
      <Card>
        <CardContent className="flex flex-col items-center justify-center py-12">
          <MessageSquare className="h-12 w-12 text-gray-300 mb-4" />
          <h3 className="text-lg font-medium text-gray-900 mb-1">No content found</h3>
          <p className="text-sm text-gray-500">
            Try adjusting your filters or create new content.
          </p>
        </CardContent>
      </Card>
    );
  }

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      {articles.map((article) => {
        const Icon = contentTypeIcons[article.contentType] || MessageSquare;
        const iconColor = contentTypeColors[article.contentType] || 'bg-gray-100 text-gray-600';
        const isSelected = selectedIds.includes(article.id);

        return (
          <Card
            key={article.id}
            className={`cursor-pointer hover:shadow-md transition-shadow ${isSelected ? 'ring-2 ring-orange-500' : ''}`}
            onClick={() => onArticleClick(article)}
          >
            <CardContent className="p-4">
              <div className="flex items-start gap-3">
                <div className={`rounded-lg p-2 ${iconColor}`}>
                  <Icon className="h-5 w-5" />
                </div>
                <div className="flex-1 min-w-0">
                  <div className="flex items-center gap-2 mb-1">
                    <h3 className="font-semibold text-gray-900 truncate">{article.title}</h3>
                    <Badge className={mkbStatusColors[article.status]} variant="secondary">
                      {mkbStatusLabels[article.status]}
                    </Badge>
                  </div>
                  <p className="text-sm text-gray-500 line-clamp-2 mb-3">
                    {article.summary || 'No summary available'}
                  </p>

                  {/* Personas */}
                  <div className="flex flex-wrap gap-1 mb-2">
                    {article.personas.slice(0, 3).map((persona) => (
                      <Badge key={persona} variant="outline" className="text-xs">
                        {mkbPersonaLabels[persona as keyof typeof mkbPersonaLabels] || persona}
                      </Badge>
                    ))}
                    {article.personas.length > 3 && (
                      <Badge variant="outline" className="text-xs">
                        +{article.personas.length - 3}
                      </Badge>
                    )}
                  </div>

                  {/* Context Tags */}
                  <div className="flex flex-wrap gap-1 mb-3">
                    {article.contextTags.slice(0, 3).map((tag) => (
                      <Badge key={tag} variant="secondary" className="text-xs bg-gray-100">
                        {tag}
                      </Badge>
                    ))}
                    {article.contextTags.length > 3 && (
                      <Badge variant="secondary" className="text-xs bg-gray-100">
                        +{article.contextTags.length - 3}
                      </Badge>
                    )}
                  </div>

                  {/* Meta */}
                  <div className="flex items-center justify-between text-xs text-gray-400">
                    <div className="flex items-center gap-3">
                      <span className="flex items-center gap-1">
                        <TrendingUp className="h-3 w-3" />
                        {article.usageCount} uses
                      </span>
                      <span className="flex items-center gap-1">
                        <Calendar className="h-3 w-3" />
                        {formatDistanceToNow(new Date(article.updatedAt), { addSuffix: true })}
                      </span>
                    </div>
                    <span>v{article.version}</span>
                  </div>
                </div>
              </div>

              {/* Add to Kit Button */}
              {onAddToKit && (
                <div className="mt-3 pt-3 border-t">
                  <Button
                    variant={isSelected ? 'secondary' : 'outline'}
                    size="sm"
                    className="w-full"
                    onClick={(e) => {
                      e.stopPropagation();
                      onAddToKit(article);
                    }}
                  >
                    <Plus className="h-4 w-4 mr-2" />
                    {isSelected ? 'In Kit' : 'Add to Kit'}
                  </Button>
                </div>
              )}
            </CardContent>
          </Card>
        );
      })}
    </div>
  );
}
