import { useState, useCallback } from 'react';
import { Button } from '@/components/ui/button';
import { BookOpen, Plus, Loader2 } from 'lucide-react';
import {
  useMKBArticles,
  useMKBStats,
  useCreateMKBArticle,
  useUpdateMKBArticle,
  useDeleteMKBArticle,
  useApproveMKBArticle,
  useSubmitForReview,
  usePitchKits,
  useCreatePitchKit,
  useDeletePitchKit,
} from '@/api/marketing-kb';
import {
  MKBStats,
  MKBFilters,
  MKBArticleList,
  MKBArticleDetail,
  CreateMKBArticleModal,
  PitchKitBuilder,
  PitchKitList,
} from '@/components/marketing-kb';
import { toast } from '@/lib/toast';
import { useAuth } from '@/contexts/AuthContext';
import type {
  MKBArticle,
  MKBArticleFilters,
  CreateMKBArticleRequest,
  UpdateMKBArticleRequest,
  CreatePitchKitRequest,
  PitchKit,
} from '@/types';

export function MarketingKB() {
  const { hasPermission } = useAuth();
  const canCreate = hasPermission('mkb:create');

  // Filters state
  const [search, setSearch] = useState('');
  const [contentType, setContentType] = useState('all');
  const [persona, setPersona] = useState('all');
  const [contextTag, setContextTag] = useState('all');
  const [status, setStatus] = useState('all');

  // Modal state
  const [createModalOpen, setCreateModalOpen] = useState(false);
  const [selectedArticle, setSelectedArticle] = useState<MKBArticle | null>(null);
  const [detailSheetOpen, setDetailSheetOpen] = useState(false);

  // Pitch Kit state
  const [selectedKitArticles, setSelectedKitArticles] = useState<MKBArticle[]>([]);

  // Build filters object
  const filters: MKBArticleFilters = {
    q: search || undefined,
    contentType: contentType !== 'all' ? (contentType as MKBArticleFilters['contentType']) : undefined,
    persona: persona !== 'all' ? (persona as MKBArticleFilters['persona']) : undefined,
    contextTag: contextTag !== 'all' ? (contextTag as MKBArticleFilters['contextTag']) : undefined,
    status: status !== 'all' ? (status as MKBArticleFilters['status']) : undefined,
    limit: 50,
  };

  const hasFilters = search !== '' || contentType !== 'all' || persona !== 'all' || contextTag !== 'all' || status !== 'all';

  // Queries
  const { data: articlesData, isLoading: isLoadingArticles } = useMKBArticles(filters);
  const { data: stats, isLoading: isLoadingStats } = useMKBStats();
  const { data: pitchKitsData, isLoading: isLoadingPitchKits } = usePitchKits({ limit: 20 });

  // Mutations
  const createMutation = useCreateMKBArticle();
  const updateMutation = useUpdateMKBArticle();
  const deleteMutation = useDeleteMKBArticle();
  const approveMutation = useApproveMKBArticle();
  const submitForReviewMutation = useSubmitForReview();
  const createPitchKitMutation = useCreatePitchKit();
  const deletePitchKitMutation = useDeletePitchKit();

  const articles = articlesData?.items || [];
  const pitchKits = pitchKitsData?.items || [];

  const handleClearFilters = () => {
    setSearch('');
    setContentType('all');
    setPersona('all');
    setContextTag('all');
    setStatus('all');
  };

  const handleArticleClick = (article: MKBArticle) => {
    setSelectedArticle(article);
    setDetailSheetOpen(true);
  };

  const handleCreate = useCallback((data: CreateMKBArticleRequest) => {
    createMutation.mutate(data, {
      onSuccess: () => {
        setCreateModalOpen(false);
        toast.success('Content created', 'Your content has been saved as a draft');
      },
      onError: () => {
        toast.error('Failed to create content', 'Please try again');
      },
    });
  }, [createMutation]);

  const handleUpdate = useCallback((id: string, data: UpdateMKBArticleRequest) => {
    updateMutation.mutate({ id, data }, {
      onSuccess: (updatedArticle) => {
        setSelectedArticle(updatedArticle);
        toast.success('Content updated', 'Your changes have been saved');
      },
      onError: () => {
        toast.error('Failed to update content', 'Please try again');
      },
    });
  }, [updateMutation]);

  const handleSubmitForReview = useCallback((id: string) => {
    submitForReviewMutation.mutate(id, {
      onSuccess: (updatedArticle) => {
        setSelectedArticle(updatedArticle);
        toast.success('Submitted for review', 'Content is now pending approval');
      },
      onError: () => {
        toast.error('Failed to submit for review', 'Please try again');
      },
    });
  }, [submitForReviewMutation]);

  const handleApprove = useCallback((id: string) => {
    approveMutation.mutate(id, {
      onSuccess: (updatedArticle) => {
        setSelectedArticle(updatedArticle);
        toast.success('Content approved', 'Content is now ready for use');
      },
      onError: () => {
        toast.error('Failed to approve content', 'Please try again');
      },
    });
  }, [approveMutation]);

  const handleDelete = useCallback((id: string) => {
    deleteMutation.mutate(id, {
      onSuccess: () => {
        setDetailSheetOpen(false);
        setSelectedArticle(null);
        toast.success('Content archived', 'The content has been archived');
      },
      onError: () => {
        toast.error('Failed to archive content', 'Please try again');
      },
    });
  }, [deleteMutation]);

  // Pitch Kit handlers
  const handleAddToKit = (article: MKBArticle) => {
    setSelectedKitArticles((prev) => {
      const exists = prev.some((a) => a.id === article.id);
      if (exists) {
        return prev.filter((a) => a.id !== article.id);
      }
      return [...prev, article];
    });
  };

  const handleRemoveFromKit = (id: string) => {
    setSelectedKitArticles((prev) => prev.filter((a) => a.id !== id));
  };

  const handleClearKit = () => {
    setSelectedKitArticles([]);
  };

  const handleReorderKit = (fromIndex: number, toIndex: number) => {
    setSelectedKitArticles((prev) => {
      const updated = [...prev];
      const [removed] = updated.splice(fromIndex, 1);
      updated.splice(toIndex, 0, removed);
      return updated;
    });
  };

  const handleSavePitchKit = useCallback((data: CreatePitchKitRequest) => {
    createPitchKitMutation.mutate(data, {
      onSuccess: () => {
        setSelectedKitArticles([]);
        toast.success('Pitch kit saved', 'Your pitch kit has been saved');
      },
      onError: () => {
        toast.error('Failed to save pitch kit', 'Please try again');
      },
    });
  }, [createPitchKitMutation]);

  const handleDeletePitchKit = useCallback((id: string) => {
    deletePitchKitMutation.mutate(id, {
      onSuccess: () => {
        toast.success('Pitch kit deleted', 'The pitch kit has been removed');
      },
      onError: () => {
        toast.error('Failed to delete pitch kit', 'Please try again');
      },
    });
  }, [deletePitchKitMutation]);

  const handleSelectPitchKit = (kit: PitchKit) => {
    // If the kit has articles populated, use them
    // Otherwise we'd need to fetch the full kit with articles
    if (kit.articles && kit.articles.length > 0) {
      setSelectedKitArticles(kit.articles);
      toast.success('Kit loaded', `Loaded "${kit.name}" with ${kit.articles.length} items`);
    } else {
      // For now, just show a message
      toast.info('Kit selected', `Selected "${kit.name}". Articles will be loaded.`);
    }
  };

  const isUpdating = updateMutation.isPending || approveMutation.isPending || deleteMutation.isPending || submitForReviewMutation.isPending;

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="rounded-lg bg-orange-100 p-2">
            <BookOpen className="h-6 w-6 text-orange-600" />
          </div>
          <div>
            <h1 className="text-2xl font-bold text-gray-900">Marketing Knowledge Base</h1>
            <p className="text-sm text-gray-500">
              Sales enablement content organized by persona and context
            </p>
          </div>
        </div>
        {canCreate && (
          <Button onClick={() => setCreateModalOpen(true)}>
            <Plus className="h-4 w-4 mr-2" />
            New Content
          </Button>
        )}
      </div>

      {/* Stats */}
      <MKBStats stats={stats} isLoading={isLoadingStats} />

      {/* Filters */}
      <MKBFilters
        search={search}
        onSearchChange={setSearch}
        contentType={contentType}
        onContentTypeChange={setContentType}
        persona={persona}
        onPersonaChange={setPersona}
        contextTag={contextTag}
        onContextTagChange={setContextTag}
        status={status}
        onStatusChange={setStatus}
        onClearFilters={handleClearFilters}
        hasFilters={hasFilters}
      />

      {/* Main Content Area */}
      <div className="grid grid-cols-1 lg:grid-cols-4 gap-6">
        {/* Articles List (3 columns) */}
        <div className="lg:col-span-3">
          <MKBArticleList
            articles={articles}
            isLoading={isLoadingArticles}
            onArticleClick={handleArticleClick}
            onAddToKit={handleAddToKit}
            selectedIds={selectedKitArticles.map((a) => a.id)}
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
        </div>

        {/* Pitch Kit Builder Sidebar (1 column) */}
        <div className="lg:col-span-1 space-y-4">
          <PitchKitBuilder
            selectedArticles={selectedKitArticles}
            onRemoveArticle={handleRemoveFromKit}
            onClearAll={handleClearKit}
            onReorder={handleReorderKit}
            onSave={handleSavePitchKit}
            isSaving={createPitchKitMutation.isPending}
          />

          <PitchKitList
            kits={pitchKits}
            isLoading={isLoadingPitchKits}
            onSelect={handleSelectPitchKit}
            onDelete={handleDeletePitchKit}
          />
        </div>
      </div>

      {/* Create Modal */}
      <CreateMKBArticleModal
        open={createModalOpen}
        onClose={() => setCreateModalOpen(false)}
        onSubmit={handleCreate}
        isLoading={createMutation.isPending}
      />

      {/* Article Detail Sheet */}
      <MKBArticleDetail
        article={selectedArticle}
        open={detailSheetOpen}
        onClose={() => {
          setDetailSheetOpen(false);
          setSelectedArticle(null);
        }}
        onUpdate={handleUpdate}
        onSubmitForReview={handleSubmitForReview}
        onApprove={handleApprove}
        onDelete={handleDelete}
        isUpdating={isUpdating}
      />
    </div>
  );
}
