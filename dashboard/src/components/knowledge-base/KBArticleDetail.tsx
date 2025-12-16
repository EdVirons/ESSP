import { useState } from 'react';
import { Sheet, SheetHeader, SheetBody, SheetFooter } from '@/components/ui/sheet';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
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
import {
  Pencil,
  X,
  Send,
  Trash2,
  Calendar,
  User,
  Tag,
  Loader2,
} from 'lucide-react';
import { formatDistanceToNow, format } from 'date-fns';
import type {
  KBArticle,
  KBContentType,
  KBModule,
  KBLifecycleStage,
  UpdateKBArticleRequest,
} from '@/types';
import {
  contentTypeLabels,
  moduleLabels,
  lifecycleStageLabels,
  statusLabels,
} from '@/types';
import { useAuth } from '@/contexts/AuthContext';

interface KBArticleDetailProps {
  article: KBArticle | null;
  open: boolean;
  onClose: () => void;
  onUpdate: (id: string, data: UpdateKBArticleRequest) => void;
  onPublish: (id: string) => void;
  onDelete: (id: string) => void;
  isUpdating: boolean;
}

const statusColors: Record<string, string> = {
  draft: 'bg-amber-100 text-amber-800',
  published: 'bg-green-100 text-green-800',
  archived: 'bg-gray-100 text-gray-800',
};

export function KBArticleDetail({
  article,
  open,
  onClose,
  onUpdate,
  onPublish,
  onDelete,
  isUpdating,
}: KBArticleDetailProps) {
  const { hasPermission } = useAuth();
  const canEdit = hasPermission('kb:update');
  const canDelete = hasPermission('kb:delete');

  const [isEditing, setIsEditing] = useState(false);
  const [editForm, setEditForm] = useState<UpdateKBArticleRequest>({});

  const handleStartEdit = () => {
    if (!article) return;
    setEditForm({
      title: article.title,
      slug: article.slug,
      summary: article.summary,
      content: article.content,
      contentType: article.contentType,
      module: article.module,
      lifecycleStage: article.lifecycleStage,
      tags: article.tags,
    });
    setIsEditing(true);
  };

  const handleCancelEdit = () => {
    setIsEditing(false);
    setEditForm({});
  };

  const handleSave = () => {
    if (!article) return;
    onUpdate(article.id, editForm);
    setIsEditing(false);
  };

  const handlePublish = () => {
    if (!article) return;
    onPublish(article.id);
  };

  const handleDelete = () => {
    if (!article) return;
    if (window.confirm('Are you sure you want to archive this article?')) {
      onDelete(article.id);
      onClose();
    }
  };

  if (!article) return null;

  return (
    <Sheet open={open} onClose={onClose} side="right" className="sm:max-w-2xl">
      <SheetHeader onClose={onClose}>
        <div className="flex items-center gap-3">
          <span>{isEditing ? 'Edit Article' : article.title}</span>
          {!isEditing && (
            <Badge className={statusColors[article.status]}>
              {statusLabels[article.status]}
            </Badge>
          )}
        </div>
      </SheetHeader>

      <SheetBody>
        <div className="space-y-6">
          {isEditing ? (
            // Edit Form
            <div className="space-y-4">
              <div className="space-y-2">
                <Label>Title</Label>
                <Input
                  value={editForm.title || ''}
                  onChange={(e) => setEditForm({ ...editForm, title: e.target.value })}
                />
              </div>

              <div className="space-y-2">
                <Label>Slug</Label>
                <Input
                  value={editForm.slug || ''}
                  onChange={(e) => setEditForm({ ...editForm, slug: e.target.value })}
                />
              </div>

              <div className="space-y-2">
                <Label>Summary</Label>
                <Textarea
                  value={editForm.summary || ''}
                  onChange={(e) => setEditForm({ ...editForm, summary: e.target.value })}
                  rows={2}
                />
              </div>

              <div className="grid grid-cols-3 gap-4">
                <div className="space-y-2">
                  <Label>Content Type</Label>
                  <Select
                    value={editForm.contentType}
                    onValueChange={(value) => setEditForm({ ...editForm, contentType: value as KBContentType })}
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
                    value={editForm.module}
                    onValueChange={(value) => setEditForm({ ...editForm, module: value as KBModule })}
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
                    value={editForm.lifecycleStage}
                    onValueChange={(value) => setEditForm({ ...editForm, lifecycleStage: value as KBLifecycleStage })}
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
                <Label>Content</Label>
                <Textarea
                  value={editForm.content || ''}
                  onChange={(e) => setEditForm({ ...editForm, content: e.target.value })}
                  rows={12}
                  className="font-mono text-sm"
                />
              </div>
            </div>
          ) : (
            // View Mode
            <>
              {/* Summary */}
              {article.summary && (
                <p className="text-gray-600">{article.summary}</p>
              )}

              {/* Meta info */}
              <div className="flex flex-wrap gap-4 text-sm text-gray-500">
                <span className="flex items-center gap-1">
                  <User className="h-4 w-4" />
                  {article.createdByName}
                </span>
                <span className="flex items-center gap-1">
                  <Calendar className="h-4 w-4" />
                  {format(new Date(article.createdAt), 'PPP')}
                </span>
                <span className="text-xs">
                  Updated {formatDistanceToNow(new Date(article.updatedAt), { addSuffix: true })}
                </span>
              </div>

              {/* Categories */}
              <div className="flex flex-wrap gap-2">
                <Badge variant="outline">{contentTypeLabels[article.contentType]}</Badge>
                <Badge variant="outline">{moduleLabels[article.module]}</Badge>
                <Badge variant="outline">{lifecycleStageLabels[article.lifecycleStage]}</Badge>
              </div>

              {/* Tags */}
              {article.tags && article.tags.length > 0 && (
                <div className="flex items-center gap-2 flex-wrap">
                  <Tag className="h-4 w-4 text-gray-400" />
                  {article.tags.map((tag) => (
                    <Badge key={tag} variant="secondary" className="text-xs">
                      {tag}
                    </Badge>
                  ))}
                </div>
              )}

              {/* Content */}
              <div className="prose prose-sm max-w-none">
                <div className="bg-gray-50 rounded-lg p-4 whitespace-pre-wrap font-mono text-sm">
                  {article.content}
                </div>
              </div>
            </>
          )}
        </div>
      </SheetBody>

      <SheetFooter>
        {isEditing ? (
          <>
            <Button variant="ghost" onClick={handleCancelEdit}>
              <X className="h-4 w-4 mr-2" />
              Cancel
            </Button>
            <Button onClick={handleSave} disabled={isUpdating}>
              {isUpdating && <Loader2 className="h-4 w-4 mr-2 animate-spin" />}
              Save Changes
            </Button>
          </>
        ) : (
          <>
            {canDelete && article.status !== 'archived' && (
              <Button variant="ghost" size="sm" className="text-red-600 mr-auto" onClick={handleDelete}>
                <Trash2 className="h-4 w-4 mr-2" />
                Archive
              </Button>
            )}
            {canEdit && (
              <Button variant="outline" onClick={handleStartEdit}>
                <Pencil className="h-4 w-4 mr-2" />
                Edit
              </Button>
            )}
            {canEdit && article.status === 'draft' && (
              <Button onClick={handlePublish} disabled={isUpdating}>
                {isUpdating && <Loader2 className="h-4 w-4 mr-2 animate-spin" />}
                <Send className="h-4 w-4 mr-2" />
                Publish
              </Button>
            )}
          </>
        )}
      </SheetFooter>
    </Sheet>
  );
}
