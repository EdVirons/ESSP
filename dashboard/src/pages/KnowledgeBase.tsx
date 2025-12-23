import { useState, useCallback } from 'react';
import { Button } from '@/components/ui/button';
import { BookOpen, Plus, Loader2 } from 'lucide-react';
import {
  useKBArticles,
  useKBStats,
  useCreateKBArticle,
  useUpdateKBArticle,
  usePublishKBArticle,
  useDeleteKBArticle,
} from '@/api/knowledge-base';
import {
  KBStats,
  KBFilters,
  KBArticleList,
  KBArticleDetail,
  CreateArticleModal,
} from '@/components/knowledge-base';
import { toast } from '@/lib/toast';
import { useAuth } from '@/contexts/AuthContext';
import type { KBArticle, KBArticleFilters, CreateKBArticleRequest, UpdateKBArticleRequest } from '@/types';

export function KnowledgeBase() {
  const { hasPermission } = useAuth();
  const canCreate = hasPermission('kb:create');

  // Filters state
  const [search, setSearch] = useState('');
  const [contentType, setContentType] = useState('all');
  const [module, setModule] = useState('all');
  const [lifecycleStage, setLifecycleStage] = useState('all');
  const [status, setStatus] = useState('all');

  // Modal state
  const [createModalOpen, setCreateModalOpen] = useState(false);
  const [selectedArticle, setSelectedArticle] = useState<KBArticle | null>(null);
  const [detailSheetOpen, setDetailSheetOpen] = useState(false);

  // Build filters object
  const filters: KBArticleFilters = {
    q: search || undefined,
    contentType: contentType !== 'all' ? (contentType as KBArticleFilters['contentType']) : undefined,
    module: module !== 'all' ? (module as KBArticleFilters['module']) : undefined,
    lifecycleStage: lifecycleStage !== 'all' ? (lifecycleStage as KBArticleFilters['lifecycleStage']) : undefined,
    status: status !== 'all' ? (status as KBArticleFilters['status']) : undefined,
    limit: 50,
  };

  const hasFilters = search !== '' || contentType !== 'all' || module !== 'all' || lifecycleStage !== 'all' || status !== 'all';

  // Queries
  const { data: articlesData, isLoading: isLoadingArticles } = useKBArticles(filters);
  const { data: stats, isLoading: isLoadingStats } = useKBStats();

  // Mutations
  const createMutation = useCreateKBArticle();
  const updateMutation = useUpdateKBArticle();
  const publishMutation = usePublishKBArticle();
  const deleteMutation = useDeleteKBArticle();

  const articles = articlesData?.items || [];

  const handleClearFilters = () => {
    setSearch('');
    setContentType('all');
    setModule('all');
    setLifecycleStage('all');
    setStatus('all');
  };

  const handleArticleClick = (article: KBArticle) => {
    setSelectedArticle(article);
    setDetailSheetOpen(true);
  };

  const handleCreate = useCallback((data: CreateKBArticleRequest) => {
    createMutation.mutate(data, {
      onSuccess: () => {
        setCreateModalOpen(false);
        toast.success('Article created', 'Your article has been saved as a draft');
      },
      onError: () => {
        toast.error('Failed to create article', 'Please try again');
      },
    });
  }, [createMutation]);

  const handleUpdate = useCallback((id: string, data: UpdateKBArticleRequest) => {
    updateMutation.mutate({ id, data }, {
      onSuccess: (updatedArticle) => {
        setSelectedArticle(updatedArticle);
        toast.success('Article updated', 'Your changes have been saved');
      },
      onError: () => {
        toast.error('Failed to update article', 'Please try again');
      },
    });
  }, [updateMutation]);

  const handlePublish = useCallback((id: string) => {
    publishMutation.mutate(id, {
      onSuccess: (updatedArticle) => {
        setSelectedArticle(updatedArticle);
        toast.success('Article published', 'Your article is now live');
      },
      onError: () => {
        toast.error('Failed to publish article', 'Please try again');
      },
    });
  }, [publishMutation]);

  const handleDelete = useCallback((id: string) => {
    deleteMutation.mutate(id, {
      onSuccess: () => {
        setDetailSheetOpen(false);
        setSelectedArticle(null);
        toast.success('Article archived', 'The article has been archived');
      },
      onError: () => {
        toast.error('Failed to archive article', 'Please try again');
      },
    });
  }, [deleteMutation]);

  const isUpdating = updateMutation.isPending || publishMutation.isPending || deleteMutation.isPending;

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <div className="flex items-center gap-3">
          <div className="rounded-lg bg-emerald-100 p-2">
            <BookOpen className="h-5 w-5 sm:h-6 sm:w-6 text-emerald-600" />
          </div>
          <div>
            <h1 className="text-xl sm:text-2xl font-bold text-gray-900">Knowledge Base</h1>
            <p className="text-sm text-gray-500">
              Technical documentation for field operations
            </p>
          </div>
        </div>
        {canCreate && (
          <Button onClick={() => setCreateModalOpen(true)} className="w-full sm:w-auto">
            <Plus className="h-4 w-4 mr-2" />
            Create Article
          </Button>
        )}
      </div>

      {/* Stats */}
      <KBStats stats={stats} isLoading={isLoadingStats} />

      {/* Filters */}
      <KBFilters
        search={search}
        onSearchChange={setSearch}
        contentType={contentType}
        onContentTypeChange={setContentType}
        module={module}
        onModuleChange={setModule}
        lifecycleStage={lifecycleStage}
        onLifecycleStageChange={setLifecycleStage}
        status={status}
        onStatusChange={setStatus}
        onClearFilters={handleClearFilters}
        hasFilters={hasFilters}
      />

      {/* Articles List */}
      <KBArticleList
        articles={articles}
        isLoading={isLoadingArticles}
        onArticleClick={handleArticleClick}
      />

      {/* Load More (if pagination needed) */}
      {articlesData?.nextCursor && (
        <div className="flex justify-center pt-4">
          <Button variant="outline" disabled={isLoadingArticles}>
            {isLoadingArticles && <Loader2 className="h-4 w-4 mr-2 animate-spin" />}
            Load More
          </Button>
        </div>
      )}

      {/* Create Modal */}
      <CreateArticleModal
        open={createModalOpen}
        onClose={() => setCreateModalOpen(false)}
        onSubmit={handleCreate}
        isLoading={createMutation.isPending}
      />

      {/* Article Detail Sheet */}
      <KBArticleDetail
        article={selectedArticle}
        open={detailSheetOpen}
        onClose={() => {
          setDetailSheetOpen(false);
          setSelectedArticle(null);
        }}
        onUpdate={handleUpdate}
        onPublish={handlePublish}
        onDelete={handleDelete}
        isUpdating={isUpdating}
      />
    </div>
  );
}
