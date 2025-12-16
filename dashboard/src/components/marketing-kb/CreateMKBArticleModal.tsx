import { useState } from 'react';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from '@/components/ui/dialog';
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
import { Loader2, X } from 'lucide-react';
import type { CreateMKBArticleRequest, MKBContentType, MKBPersona, MKBContextTag } from '@/types';
import { mkbContentTypeLabels, mkbPersonaLabels, mkbContextTagLabels } from '@/types';

interface CreateMKBArticleModalProps {
  open: boolean;
  onClose: () => void;
  onSubmit: (data: CreateMKBArticleRequest) => void;
  isLoading: boolean;
}

export function CreateMKBArticleModal({
  open,
  onClose,
  onSubmit,
  isLoading,
}: CreateMKBArticleModalProps) {
  const [form, setForm] = useState<CreateMKBArticleRequest>({
    title: '',
    summary: '',
    content: '',
    contentType: 'messaging',
    personas: [],
    contextTags: [],
    tags: [],
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!form.title || !form.content) return;
    onSubmit(form);
  };

  const togglePersona = (persona: string) => {
    const current = form.personas || [];
    const updated = current.includes(persona)
      ? current.filter((p) => p !== persona)
      : [...current, persona];
    setForm({ ...form, personas: updated });
  };

  const toggleContextTag = (tag: string) => {
    const current = form.contextTags || [];
    const updated = current.includes(tag)
      ? current.filter((t) => t !== tag)
      : [...current, tag];
    setForm({ ...form, contextTags: updated });
  };

  const handleClose = () => {
    setForm({
      title: '',
      summary: '',
      content: '',
      contentType: 'messaging',
      personas: [],
      contextTags: [],
      tags: [],
    });
    onClose();
  };

  return (
    <Dialog open={open} onOpenChange={handleClose}>
      <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>Create Marketing Content</DialogTitle>
          <DialogDescription>
            Add new sales enablement content to the knowledge base.
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-2">
            <Label>Title *</Label>
            <Input
              value={form.title}
              onChange={(e) => setForm({ ...form, title: e.target.value })}
              placeholder="e.g., CBC Value Proposition for Directors"
              required
            />
          </div>

          <div className="space-y-2">
            <Label>Content Type</Label>
            <Select
              value={form.contentType}
              onValueChange={(value) => setForm({ ...form, contentType: value as MKBContentType })}
            >
              <SelectTrigger>
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                {(Object.keys(mkbContentTypeLabels) as MKBContentType[]).map((type) => (
                  <SelectItem key={type} value={type}>
                    {mkbContentTypeLabels[type]}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          <div className="space-y-2">
            <Label>Summary</Label>
            <Textarea
              value={form.summary}
              onChange={(e) => setForm({ ...form, summary: e.target.value })}
              placeholder="Brief description of this content..."
              rows={2}
            />
          </div>

          <div className="space-y-2">
            <Label>Target Personas</Label>
            <div className="flex flex-wrap gap-2">
              {(Object.keys(mkbPersonaLabels) as MKBPersona[]).map((persona) => (
                <Badge
                  key={persona}
                  variant={form.personas?.includes(persona) ? 'default' : 'outline'}
                  className="cursor-pointer"
                  onClick={() => togglePersona(persona)}
                >
                  {mkbPersonaLabels[persona]}
                  {form.personas?.includes(persona) && (
                    <X className="h-3 w-3 ml-1" />
                  )}
                </Badge>
              ))}
            </div>
          </div>

          <div className="space-y-2">
            <Label>School Context</Label>
            <div className="flex flex-wrap gap-2">
              {(Object.keys(mkbContextTagLabels) as MKBContextTag[]).map((tag) => (
                <Badge
                  key={tag}
                  variant={form.contextTags?.includes(tag) ? 'default' : 'outline'}
                  className="cursor-pointer"
                  onClick={() => toggleContextTag(tag)}
                >
                  {mkbContextTagLabels[tag]}
                  {form.contextTags?.includes(tag) && (
                    <X className="h-3 w-3 ml-1" />
                  )}
                </Badge>
              ))}
            </div>
          </div>

          <div className="space-y-2">
            <Label>Content *</Label>
            <Textarea
              value={form.content}
              onChange={(e) => setForm({ ...form, content: e.target.value })}
              placeholder="Enter your content here..."
              rows={10}
              className="font-mono text-sm"
              required
            />
          </div>

          <DialogFooter>
            <Button type="button" variant="outline" onClick={handleClose}>
              Cancel
            </Button>
            <Button type="submit" disabled={isLoading || !form.title || !form.content}>
              {isLoading && <Loader2 className="h-4 w-4 mr-2 animate-spin" />}
              Create Content
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
