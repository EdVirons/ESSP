import { useState } from 'react';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
} from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { Label } from '@/components/ui/label';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Loader2 } from 'lucide-react';
import type {
  CreateKBArticleRequest,
  KBContentType,
  KBModule,
  KBLifecycleStage,
} from '@/types';
import {
  contentTypeLabels,
  moduleLabels,
  lifecycleStageLabels,
} from '@/types';

interface CreateArticleModalProps {
  open: boolean;
  onClose: () => void;
  onSubmit: (data: CreateKBArticleRequest) => void;
  isLoading: boolean;
}

const defaultForm: CreateKBArticleRequest = {
  title: '',
  slug: '',
  summary: '',
  content: '',
  contentType: 'runbook',
  module: 'general',
  lifecycleStage: 'support',
  tags: [],
};

export function CreateArticleModal({
  open,
  onClose,
  onSubmit,
  isLoading,
}: CreateArticleModalProps) {
  const [form, setForm] = useState<CreateKBArticleRequest>(defaultForm);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!form.title || !form.content) return;
    onSubmit(form);
  };

  const handleClose = () => {
    setForm(defaultForm);
    onClose();
  };

  const generateSlug = () => {
    const slug = form.title
      .toLowerCase()
      .replace(/[^a-z0-9\s-]/g, '')
      .replace(/\s+/g, '-')
      .replace(/-+/g, '-')
      .replace(/^-|-$/g, '');
    setForm({ ...form, slug });
  };

  return (
    <Dialog open={open} onOpenChange={(open) => !open && handleClose()}>
      <DialogContent className="sm:max-w-2xl max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>Create New Article</DialogTitle>
          <DialogDescription>
            Add a new knowledge base article. Articles are created as drafts and can be published later.
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit} className="space-y-4 mt-4">
          <div className="space-y-2">
            <Label htmlFor="title">Title *</Label>
            <Input
              id="title"
              value={form.title}
              onChange={(e) => setForm({ ...form, title: e.target.value })}
              placeholder="Enter article title"
              required
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="slug">
              Slug
              <Button
                type="button"
                variant="link"
                size="sm"
                className="ml-2 h-auto p-0 text-xs"
                onClick={generateSlug}
              >
                Generate from title
              </Button>
            </Label>
            <Input
              id="slug"
              value={form.slug}
              onChange={(e) => setForm({ ...form, slug: e.target.value })}
              placeholder="article-url-slug"
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="summary">Summary</Label>
            <Textarea
              id="summary"
              value={form.summary}
              onChange={(e) => setForm({ ...form, summary: e.target.value })}
              placeholder="Brief description of the article"
              rows={2}
            />
          </div>

          <div className="grid grid-cols-3 gap-4">
            <div className="space-y-2">
              <Label>Content Type</Label>
              <Select
                value={form.contentType}
                onValueChange={(value) => setForm({ ...form, contentType: value as KBContentType })}
              >
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  {(Object.keys(contentTypeLabels) as KBContentType[]).map((type) => (
                    <SelectItem key={type} value={type}>
                      {contentTypeLabels[type]}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>

            <div className="space-y-2">
              <Label>Module</Label>
              <Select
                value={form.module}
                onValueChange={(value) => setForm({ ...form, module: value as KBModule })}
              >
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  {(Object.keys(moduleLabels) as KBModule[]).map((mod) => (
                    <SelectItem key={mod} value={mod}>
                      {moduleLabels[mod]}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>

            <div className="space-y-2">
              <Label>Lifecycle Stage</Label>
              <Select
                value={form.lifecycleStage}
                onValueChange={(value) => setForm({ ...form, lifecycleStage: value as KBLifecycleStage })}
              >
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  {(Object.keys(lifecycleStageLabels) as KBLifecycleStage[]).map((stage) => (
                    <SelectItem key={stage} value={stage}>
                      {lifecycleStageLabels[stage]}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
          </div>

          <div className="space-y-2">
            <Label htmlFor="content">Content *</Label>
            <Textarea
              id="content"
              value={form.content}
              onChange={(e) => setForm({ ...form, content: e.target.value })}
              placeholder="Write your article content here..."
              rows={10}
              className="font-mono text-sm"
              required
            />
          </div>

          <div className="flex justify-end gap-3 pt-4">
            <Button type="button" variant="outline" onClick={handleClose}>
              Cancel
            </Button>
            <Button type="submit" disabled={isLoading || !form.title || !form.content}>
              {isLoading && <Loader2 className="h-4 w-4 mr-2 animate-spin" />}
              Create Article
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
}
