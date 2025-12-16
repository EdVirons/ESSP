import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { formatDistanceToNow } from 'date-fns';
import {
  Briefcase,
  FileText,
  Trash2,
  ExternalLink,
  User,
} from 'lucide-react';
import type { PitchKit } from '@/types';
import { mkbPersonaLabels, mkbContextTagLabels } from '@/types';

interface PitchKitListProps {
  kits: PitchKit[];
  isLoading: boolean;
  onSelect: (kit: PitchKit) => void;
  onDelete: (id: string) => void;
}

function LoadingSkeleton() {
  return (
    <div className="space-y-2">
      {[1, 2, 3].map((i) => (
        <div key={i} className="p-3 border rounded-lg">
          <div className="h-4 w-3/4 bg-gray-200 rounded animate-pulse mb-2" />
          <div className="h-3 w-1/2 bg-gray-200 rounded animate-pulse" />
        </div>
      ))}
    </div>
  );
}

export function PitchKitList({
  kits,
  isLoading,
  onSelect,
  onDelete,
}: PitchKitListProps) {
  if (isLoading) {
    return (
      <Card>
        <CardHeader className="pb-3">
          <CardTitle className="text-lg flex items-center gap-2">
            <Briefcase className="h-5 w-5" />
            Saved Pitch Kits
          </CardTitle>
        </CardHeader>
        <CardContent>
          <LoadingSkeleton />
        </CardContent>
      </Card>
    );
  }

  if (kits.length === 0) {
    return (
      <Card>
        <CardHeader className="pb-3">
          <CardTitle className="text-lg flex items-center gap-2">
            <Briefcase className="h-5 w-5" />
            Saved Pitch Kits
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex flex-col items-center justify-center py-6 text-center">
            <Briefcase className="h-8 w-8 text-gray-300 mb-2" />
            <p className="text-sm text-gray-500">No saved pitch kits yet.</p>
            <p className="text-xs text-gray-400">
              Build a kit and save it to reuse later.
            </p>
          </div>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card>
      <CardHeader className="pb-3">
        <CardTitle className="text-lg flex items-center gap-2">
          <Briefcase className="h-5 w-5" />
          Saved Pitch Kits
          <Badge variant="secondary">{kits.length}</Badge>
        </CardTitle>
      </CardHeader>
      <CardContent>
        <div className="space-y-2 max-h-[400px] overflow-y-auto">
          {kits.map((kit) => (
            <div
              key={kit.id}
              className="p-3 border rounded-lg hover:bg-gray-50 transition-colors"
            >
              <div className="flex items-start justify-between gap-2">
                <div className="flex-1 min-w-0">
                  <button
                    onClick={() => onSelect(kit)}
                    className="text-left w-full"
                  >
                    <h4 className="font-medium text-sm truncate hover:text-orange-600">
                      {kit.name}
                    </h4>
                  </button>
                  {kit.description && (
                    <p className="text-xs text-gray-500 line-clamp-1 mt-0.5">
                      {kit.description}
                    </p>
                  )}
                </div>
                <div className="flex items-center gap-1">
                  <Button
                    variant="ghost"
                    size="icon"
                    className="h-6 w-6"
                    onClick={() => onSelect(kit)}
                  >
                    <ExternalLink className="h-3 w-3" />
                  </Button>
                  <Button
                    variant="ghost"
                    size="icon"
                    className="h-6 w-6 text-gray-400 hover:text-red-500"
                    onClick={() => {
                      if (window.confirm('Delete this pitch kit?')) {
                        onDelete(kit.id);
                      }
                    }}
                  >
                    <Trash2 className="h-3 w-3" />
                  </Button>
                </div>
              </div>

              <div className="flex flex-wrap items-center gap-2 mt-2">
                <Badge variant="outline" className="text-xs">
                  <User className="h-3 w-3 mr-1" />
                  {mkbPersonaLabels[kit.targetPersona as keyof typeof mkbPersonaLabels] || kit.targetPersona}
                </Badge>
                <Badge variant="secondary" className="text-xs">
                  <FileText className="h-3 w-3 mr-1" />
                  {kit.articleIds.length} items
                </Badge>
              </div>

              {kit.contextTags && kit.contextTags.length > 0 && (
                <div className="flex flex-wrap gap-1 mt-2">
                  {kit.contextTags.slice(0, 3).map((tag) => (
                    <Badge
                      key={tag}
                      variant="secondary"
                      className="text-xs bg-gray-100"
                    >
                      {mkbContextTagLabels[tag as keyof typeof mkbContextTagLabels] || tag}
                    </Badge>
                  ))}
                  {kit.contextTags.length > 3 && (
                    <Badge variant="secondary" className="text-xs bg-gray-100">
                      +{kit.contextTags.length - 3}
                    </Badge>
                  )}
                </div>
              )}

              <div className="flex items-center gap-2 mt-2 text-xs text-gray-400">
                <span>{kit.createdByName}</span>
                <span>·</span>
                <span>
                  {formatDistanceToNow(new Date(kit.updatedAt), { addSuffix: true })}
                </span>
              </div>
            </div>
          ))}
        </div>
      </CardContent>
    </Card>
  );
}
