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
  CheckCircle,
  Trash2,
  Calendar,
  User,
  Tag,
  Loader2,
  TrendingUp,
} from 'lucide-react';
import { formatDistanceToNow, format } from 'date-fns';
import type {
  MKBArticle,
  MKBContentType,
  UpdateMKBArticleRequest,
} from '@/types';
import {
  mkbContentTypeLabels,
  mkbPersonaLabels,
  mkbContextTagLabels,
  mkbStatusLabels,
  mkbStatusColors,
} from '@/types';
import { useAuth } from '@/contexts/AuthContext';

interface MKBArticleDetailProps {
  article: MKBArticle | null;
  open: boolean;
  onClose: () => void;
  onUpdate: (id: string, data: UpdateMKBArticleRequest) => void;
  onSubmitForReview: (id: string) => void;
  onApprove: (id: string) => void;
  onDelete: (id: string) => void;
  isUpdating: boolean;
}

export function MKBArticleDetail({
  article,
  open,
  onClose,
  onUpdate,
  onSubmitForReview,
  onApprove,
  onDelete,
  isUpdating,
}: MKBArticleDetailProps) {
  const { hasPermission } = useAuth();
  const canEdit = hasPermission('mkb:update');
  const canDelete = hasPermission('mkb:delete');
  const canApprove = hasPermission('mkb:approve');

  const [isEditing, setIsEditing] = useState(false);
  const [editForm, setEditForm] = useState<UpdateMKBArticleRequest>({});

  const handleStartEdit = () => {
    if (!article) return;
    setEditForm({
      title: article.title,
      slug: article.slug,
      summary: article.summary,
      content: article.content,
      contentType: article.contentType,
      personas: article.personas,
      contextTags: article.contextTags,
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

  const handleSubmitForReview = () => {
    if (!article) return;
    onSubmitForReview(article.id);
  };

  const handleApprove = () => {
    if (!article) return;
    onApprove(article.id);
  };

  const handleDelete = () => {
    if (!article) return;
    if (window.confirm('Are you sure you want to archive this content?')) {
      onDelete(article.id);
      onClose();
    }
  };

  if (!article) return null;

  return (
    <Sheet open={open} onClose={onClose} side="right" className="sm:max-w-2xl">
      <SheetHeader onClose={onClose}>
        <div className="flex items-center gap-3">
          <span>{isEditing ? 'Edit Content' : article.title}</span>
          {!isEditing && (
            <Badge className={mkbStatusColors[article.status]}>
              {mkbStatusLabels[article.status]}
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

              <div className="space-y-2">
                <Label>Content Type</Label>
                <Select
                  value={editForm.contentType}
                  onValueChange={(value) => setEditForm({ ...editForm, contentType: value as MKBContentType })}
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
                <span className="flex items-center gap-1">
                  <TrendingUp className="h-4 w-4" />
                  {article.usageCount} uses
                </span>
                <span className="text-xs">
                  Updated {formatDistanceToNow(new Date(article.updatedAt), { addSuffix: true })}
                </span>
              </div>

              {/* Type */}
              <div className="flex flex-wrap gap-2">
                <Badge variant="outline">{mkbContentTypeLabels[article.contentType]}</Badge>
                <Badge variant="outline">v{article.version}</Badge>
              </div>

              {/* Personas */}
              {article.personas && article.personas.length > 0 && (
                <div>
                  <Label className="text-xs text-gray-500">Target Personas</Label>
                  <div className="flex items-center gap-2 flex-wrap mt-1">
                    {article.personas.map((persona) => (
                      <Badge key={persona} variant="secondary">
                        {mkbPersonaLabels[persona as keyof typeof mkbPersonaLabels] || persona}
                      </Badge>
                    ))}
                  </div>
                </div>
              )}

              {/* Context Tags */}
              {article.contextTags && article.contextTags.length > 0 && (
                <div>
                  <Label className="text-xs text-gray-500">School Context</Label>
                  <div className="flex items-center gap-2 flex-wrap mt-1">
                    {article.contextTags.map((tag) => (
                      <Badge key={tag} variant="secondary" className="bg-gray-100">
                        {mkbContextTagLabels[tag as keyof typeof mkbContextTagLabels] || tag}
                      </Badge>
                    ))}
                  </div>
                </div>
              )}

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

              {/* Approval Info */}
              {article.approvedAt && (
                <div className="p-3 bg-green-50 rounded-lg">
                  <p className="text-sm text-green-700">
                    Approved by {article.approvedByName} on{' '}
                    {format(new Date(article.approvedAt), 'PPP')}
                  </p>
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
              <Button variant="outline" onClick={handleSubmitForReview} disabled={isUpdating}>
                {isUpdating && <Loader2 className="h-4 w-4 mr-2 animate-spin" />}
                <Send className="h-4 w-4 mr-2" />
                Submit for Review
              </Button>
            )}
            {canApprove && article.status === 'review' && (
              <Button onClick={handleApprove} disabled={isUpdating}>
                {isUpdating && <Loader2 className="h-4 w-4 mr-2 animate-spin" />}
                <CheckCircle className="h-4 w-4 mr-2" />
                Approve
              </Button>
            )}
          </>
        )}
      </SheetFooter>
    </Sheet>
  );
}
