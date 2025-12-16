import { useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { Label } from '@/components/ui/label';
import { Badge } from '@/components/ui/badge';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import {
  Briefcase,
  Save,
  X,
  ChevronUp,
  ChevronDown,
  Loader2,
  Package,
} from 'lucide-react';
import type { MKBArticle, MKBPersona, MKBContextTag, CreatePitchKitRequest } from '@/types';
import { mkbPersonaLabels, mkbContextTagLabels, mkbContentTypeLabels } from '@/types';

interface PitchKitBuilderProps {
  selectedArticles: MKBArticle[];
  onRemoveArticle: (id: string) => void;
  onClearAll: () => void;
  onReorder: (fromIndex: number, toIndex: number) => void;
  onSave: (data: CreatePitchKitRequest) => void;
  isSaving: boolean;
}

export function PitchKitBuilder({
  selectedArticles,
  onRemoveArticle,
  onClearAll,
  onReorder,
  onSave,
  isSaving,
}: PitchKitBuilderProps) {
  const [isExpanded, setIsExpanded] = useState(true);
  const [kitName, setKitName] = useState('');
  const [kitDescription, setKitDescription] = useState('');
  const [targetPersona, setTargetPersona] = useState<MKBPersona>('director');
  const [contextTags, setContextTags] = useState<MKBContextTag[]>([]);

  const toggleContextTag = (tag: MKBContextTag) => {
    setContextTags((prev) =>
      prev.includes(tag) ? prev.filter((t) => t !== tag) : [...prev, tag]
    );
  };

  const handleSave = () => {
    if (!kitName || selectedArticles.length === 0) return;
    onSave({
      name: kitName,
      description: kitDescription,
      targetPersona,
      contextTags,
      articleIds: selectedArticles.map((a) => a.id),
      isTemplate: false,
    });
    // Reset form after save
    setKitName('');
    setKitDescription('');
    setTargetPersona('director');
    setContextTags([]);
  };

  const moveUp = (index: number) => {
    if (index > 0) {
      onReorder(index, index - 1);
    }
  };

  const moveDown = (index: number) => {
    if (index < selectedArticles.length - 1) {
      onReorder(index, index + 1);
    }
  };

  if (selectedArticles.length === 0) {
    return (
      <Card className="h-full">
        <CardHeader className="pb-3">
          <CardTitle className="text-lg flex items-center gap-2">
            <Briefcase className="h-5 w-5" />
            Pitch Kit Builder
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex flex-col items-center justify-center py-8 text-center">
            <Package className="h-12 w-12 text-gray-300 mb-4" />
            <p className="text-sm text-gray-500">
              Select content from the list to build your pitch kit.
            </p>
            <p className="text-xs text-gray-400 mt-1">
              Click "Add to Kit" on any article card.
            </p>
          </div>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card className="h-full">
      <CardHeader className="pb-3">
        <div className="flex items-center justify-between">
          <CardTitle className="text-lg flex items-center gap-2">
            <Briefcase className="h-5 w-5" />
            Pitch Kit Builder
            <Badge variant="secondary">{selectedArticles.length}</Badge>
          </CardTitle>
          <div className="flex items-center gap-2">
            <Button variant="ghost" size="sm" onClick={onClearAll}>
              Clear All
            </Button>
            <Button
              variant="ghost"
              size="icon"
              onClick={() => setIsExpanded(!isExpanded)}
            >
              {isExpanded ? (
                <ChevronUp className="h-4 w-4" />
              ) : (
                <ChevronDown className="h-4 w-4" />
              )}
            </Button>
          </div>
        </div>
      </CardHeader>

      {isExpanded && (
        <CardContent className="space-y-4">
          {/* Kit Details */}
          <div className="space-y-3">
            <div className="space-y-1">
              <Label className="text-xs">Kit Name *</Label>
              <Input
                value={kitName}
                onChange={(e) => setKitName(e.target.value)}
                placeholder="e.g., Director Pitch - Rural Schools"
                className="h-8"
              />
            </div>

            <div className="space-y-1">
              <Label className="text-xs">Description</Label>
              <Textarea
                value={kitDescription}
                onChange={(e) => setKitDescription(e.target.value)}
                placeholder="Brief description..."
                rows={2}
                className="text-sm"
              />
            </div>

            <div className="space-y-1">
              <Label className="text-xs">Target Persona</Label>
              <Select
                value={targetPersona}
                onValueChange={(v) => setTargetPersona(v as MKBPersona)}
              >
                <SelectTrigger className="h-8">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  {(Object.keys(mkbPersonaLabels) as MKBPersona[]).map((p) => (
                    <SelectItem key={p} value={p}>
                      {mkbPersonaLabels[p]}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>

            <div className="space-y-1">
              <Label className="text-xs">Context</Label>
              <div className="flex flex-wrap gap-1">
                {(Object.keys(mkbContextTagLabels) as MKBContextTag[]).map((tag) => (
                  <Badge
                    key={tag}
                    variant={contextTags.includes(tag) ? 'default' : 'outline'}
                    className="cursor-pointer text-xs"
                    onClick={() => toggleContextTag(tag)}
                  >
                    {mkbContextTagLabels[tag]}
                  </Badge>
                ))}
              </div>
            </div>
          </div>

          {/* Selected Articles */}
          <div className="border-t pt-3">
            <Label className="text-xs text-gray-500 mb-2 block">
              Selected Content ({selectedArticles.length})
            </Label>
            <div className="space-y-2 max-h-[300px] overflow-y-auto">
              {selectedArticles.map((article, index) => (
                <div
                  key={article.id}
                  className="flex items-center gap-2 p-2 bg-gray-50 rounded-lg"
                >
                  <div className="flex flex-col gap-1">
                    <Button
                      variant="ghost"
                      size="icon"
                      className="h-5 w-5"
                      onClick={() => moveUp(index)}
                      disabled={index === 0}
                    >
                      <ChevronUp className="h-3 w-3" />
                    </Button>
                    <Button
                      variant="ghost"
                      size="icon"
                      className="h-5 w-5"
                      onClick={() => moveDown(index)}
                      disabled={index === selectedArticles.length - 1}
                    >
                      <ChevronDown className="h-3 w-3" />
                    </Button>
                  </div>
                  <div className="flex-1 min-w-0">
                    <p className="text-sm font-medium truncate">{article.title}</p>
                    <p className="text-xs text-gray-500">
                      {mkbContentTypeLabels[article.contentType]}
                    </p>
                  </div>
                  <Button
                    variant="ghost"
                    size="icon"
                    className="h-6 w-6 text-gray-400 hover:text-red-500"
                    onClick={() => onRemoveArticle(article.id)}
                  >
                    <X className="h-3 w-3" />
                  </Button>
                </div>
              ))}
            </div>
          </div>

          {/* Save Button */}
          <Button
            className="w-full"
            onClick={handleSave}
            disabled={!kitName || selectedArticles.length === 0 || isSaving}
          >
            {isSaving ? (
              <Loader2 className="h-4 w-4 mr-2 animate-spin" />
            ) : (
              <Save className="h-4 w-4 mr-2" />
            )}
            Save Pitch Kit
          </Button>
        </CardContent>
      )}
    </Card>
  );
}
