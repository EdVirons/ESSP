import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import {
  Presentation as PresentationIcon,
  Search,
  Download,
  ExternalLink,
  FileText,
  BarChart,
  Calculator,
  BookOpen,
  Video,
  Image,
  Clock,
  Eye,
  Plus,
  Loader2,
  Star,
} from 'lucide-react';
import { presentationsApi } from '@/api/presentations';
import { UploadPresentationModal } from '@/components/presentations/UploadPresentationModal';
import { toast } from '@/lib/toast';
import type {
  Presentation,
  PresentationType,
  PresentationCategory,
  CreatePresentationRequest,
} from '@/types/sales';
import { presentationTypeLabels, presentationCategoryLabels } from '@/types/sales';

const typeConfig: Record<PresentationType, { icon: React.ElementType; color: string }> = {
  presentation: { icon: PresentationIcon, color: 'text-blue-600' },
  case_study: { icon: BookOpen, color: 'text-green-600' },
  roi_calculator: { icon: Calculator, color: 'text-purple-600' },
  brochure: { icon: FileText, color: 'text-amber-600' },
  video: { icon: Video, color: 'text-red-600' },
  template: { icon: FileText, color: 'text-indigo-600' },
  other: { icon: Image, color: 'text-pink-600' },
};

const categoryLabels: Record<PresentationCategory | 'all', string> = {
  all: 'All',
  general: 'General',
  product_overview: 'Product',
  technical: 'Technical',
  pricing: 'Pricing',
  onboarding: 'Onboarding',
  training: 'Training',
};

const categories: (PresentationCategory | 'all')[] = [
  'all',
  'general',
  'product_overview',
  'technical',
  'pricing',
  'onboarding',
  'training',
];

function PresentationCard({
  item,
  onDownload,
}: {
  item: Presentation;
  onDownload: () => void;
}) {
  const config = typeConfig[item.type] || typeConfig.other;
  const Icon = config.icon;

  return (
    <Card className="hover:shadow-md transition-shadow">
      <CardContent className="p-6">
        <div className="flex items-start gap-4">
          <div className={`p-3 rounded-lg bg-gray-100 ${config.color}`}>
            <Icon className="h-6 w-6" />
          </div>
          <div className="flex-1 min-w-0">
            <div className="flex items-start justify-between gap-2">
              <div>
                <div className="flex items-center gap-2">
                  <h3 className="font-semibold text-gray-900 line-clamp-1">{item.title}</h3>
                  {item.isFeatured && (
                    <Star className="h-4 w-4 text-yellow-500 fill-yellow-500" />
                  )}
                </div>
                <Badge variant="secondary" className="mt-1">
                  {presentationTypeLabels[item.type]}
                </Badge>
              </div>
              <Badge variant="outline">{presentationCategoryLabels[item.category]}</Badge>
            </div>
            <p className="text-sm text-gray-500 mt-2 line-clamp-2">{item.description}</p>
            <div className="flex items-center gap-4 mt-4 text-xs text-gray-400">
              <span className="flex items-center gap-1">
                <Clock className="h-3 w-3" />
                Updated {new Date(item.updatedAt).toLocaleDateString()}
              </span>
              <span className="flex items-center gap-1">
                <Eye className="h-3 w-3" />
                {item.viewCount} views
              </span>
              <span className="flex items-center gap-1">
                <Download className="h-3 w-3" />
                {item.downloadCount} downloads
              </span>
            </div>
            <div className="flex items-center gap-2 mt-4">
              <Button size="sm" variant="outline" onClick={onDownload}>
                <Download className="h-4 w-4 mr-1" />
                Download
              </Button>
              {item.previewUrl && (
                <Button size="sm" variant="ghost" asChild>
                  <a href={item.previewUrl} target="_blank" rel="noopener noreferrer">
                    <ExternalLink className="h-4 w-4 mr-1" />
                    Preview
                  </a>
                </Button>
              )}
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}

export function Presentations() {
  const queryClient = useQueryClient();
  const [searchTerm, setSearchTerm] = useState('');
  const [selectedCategory, setSelectedCategory] = useState<PresentationCategory | 'all'>('all');
  const [selectedType, setSelectedType] = useState<PresentationType | 'all'>('all');
  const [uploadModalOpen, setUploadModalOpen] = useState(false);
  const [uploadProgress, setUploadProgress] = useState(0);

  // Fetch presentations
  const { data, isLoading, error } = useQuery({
    queryKey: ['presentations', searchTerm, selectedCategory, selectedType],
    queryFn: () =>
      presentationsApi.list({
        search: searchTerm || undefined,
        category: selectedCategory !== 'all' ? selectedCategory : undefined,
        type: selectedType !== 'all' ? selectedType : undefined,
        active: true,
        limit: 100,
      }),
  });

  const presentations = data?.presentations || [];

  // Upload mutation
  const uploadMutation = useMutation({
    mutationFn: async ({
      data,
      file,
    }: {
      data: CreatePresentationRequest;
      file: File | null;
    }) => {
      // Create presentation record and get upload URL
      const response = await presentationsApi.create(data);

      // If there's a file and upload URL, upload the file
      if (file && response.uploadUrl) {
        setUploadProgress(30);
        await presentationsApi.uploadFile(response.uploadUrl, file);
        setUploadProgress(100);
      }

      return response.presentation;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['presentations'] });
      setUploadModalOpen(false);
      setUploadProgress(0);
      toast.success('Presentation uploaded', 'Your presentation has been added');
    },
    onError: () => {
      setUploadProgress(0);
      toast.error('Upload failed', 'Please try again');
    },
  });

  const handleDownload = async (presentation: Presentation) => {
    try {
      const { url, fileName } = await presentationsApi.getDownloadUrl(presentation.id);
      // Create a temporary link and click it
      const link = document.createElement('a');
      link.href = url;
      link.download = fileName;
      link.target = '_blank';
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
    } catch {
      toast.error('Download failed', 'Could not generate download link');
    }
  };

  const handleUpload = (data: CreatePresentationRequest, file: File | null) => {
    uploadMutation.mutate({ data, file });
  };

  // Calculate type counts
  const typeCounts = Object.keys(typeConfig).map((type) => ({
    type: type as PresentationType,
    label: presentationTypeLabels[type as PresentationType],
    icon: typeConfig[type as PresentationType].icon,
    color: typeConfig[type as PresentationType].color,
    count: presentations.filter((p) => p.type === type).length,
  }));

  if (error) {
    return (
      <div className="flex items-center justify-center h-64">
        <p className="text-red-500">Failed to load presentations</p>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Sales Presentations</h1>
          <p className="text-sm text-gray-500 mt-1">
            Access sales materials, case studies, and marketing content
          </p>
        </div>
        <Button onClick={() => setUploadModalOpen(true)}>
          <Plus className="h-4 w-4 mr-2" />
          Upload
        </Button>
      </div>

      {/* Type Overview Cards */}
      <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-7 gap-4">
        {typeCounts.map(({ type, label, icon: Icon, color, count }) => (
          <Card
            key={type}
            className={`cursor-pointer transition-all ${
              selectedType === type ? 'ring-2 ring-primary' : 'hover:shadow-md'
            }`}
            onClick={() => setSelectedType(selectedType === type ? 'all' : type)}
          >
            <CardContent className="p-4 text-center">
              <div className={`mx-auto w-fit p-2 rounded-lg bg-gray-100 ${color} mb-2`}>
                <Icon className="h-5 w-5" />
              </div>
              <p className="text-xs text-gray-500">{label}</p>
              <p className="text-lg font-bold">{count}</p>
            </CardContent>
          </Card>
        ))}
      </div>

      {/* Filters */}
      <div className="flex items-center gap-4 flex-wrap">
        <div className="relative flex-1 max-w-sm">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-gray-400" />
          <Input
            placeholder="Search presentations..."
            className="pl-10"
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
          />
        </div>
        <div className="flex gap-2 flex-wrap">
          {categories.map((category) => (
            <Button
              key={category}
              variant={selectedCategory === category ? 'default' : 'outline'}
              size="sm"
              onClick={() => setSelectedCategory(category)}
            >
              {categoryLabels[category]}
            </Button>
          ))}
        </div>
        {selectedType !== 'all' && (
          <Button variant="ghost" size="sm" onClick={() => setSelectedType('all')}>
            Clear Type Filter
          </Button>
        )}
      </div>

      {/* Loading State */}
      {isLoading && (
        <div className="flex items-center justify-center h-64">
          <Loader2 className="h-8 w-8 animate-spin text-gray-400" />
        </div>
      )}

      {/* Presentations Grid */}
      {!isLoading && presentations.length === 0 ? (
        <Card>
          <CardContent className="py-12 text-center">
            <PresentationIcon className="h-12 w-12 mx-auto text-gray-300 mb-4" />
            <h3 className="text-lg font-semibold text-gray-900">No presentations found</h3>
            <p className="text-sm text-gray-500 mt-1">
              {searchTerm || selectedCategory !== 'all' || selectedType !== 'all'
                ? 'Try adjusting your search or filter criteria'
                : 'Upload your first presentation to get started'}
            </p>
            {!searchTerm && selectedCategory === 'all' && selectedType === 'all' && (
              <Button className="mt-4" onClick={() => setUploadModalOpen(true)}>
                <Plus className="h-4 w-4 mr-2" />
                Upload Presentation
              </Button>
            )}
          </CardContent>
        </Card>
      ) : (
        !isLoading && (
          <div className="grid gap-4 md:grid-cols-2">
            {presentations.map((presentation) => (
              <PresentationCard
                key={presentation.id}
                item={presentation}
                onDownload={() => handleDownload(presentation)}
              />
            ))}
          </div>
        )
      )}

      {/* Stats Summary */}
      {presentations.length > 0 && (
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <BarChart className="h-5 w-5" />
              Content Statistics
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-2 md:grid-cols-4 gap-6">
              <div>
                <p className="text-sm text-gray-500">Total Materials</p>
                <p className="text-2xl font-bold">{presentations.length}</p>
              </div>
              <div>
                <p className="text-sm text-gray-500">Total Views</p>
                <p className="text-2xl font-bold">
                  {presentations
                    .reduce((sum, item) => sum + item.viewCount, 0)
                    .toLocaleString()}
                </p>
              </div>
              <div>
                <p className="text-sm text-gray-500">Total Downloads</p>
                <p className="text-2xl font-bold">
                  {presentations
                    .reduce((sum, item) => sum + item.downloadCount, 0)
                    .toLocaleString()}
                </p>
              </div>
              <div>
                <p className="text-sm text-gray-500">Featured Items</p>
                <p className="text-2xl font-bold">
                  {presentations.filter((p) => p.isFeatured).length}
                </p>
              </div>
            </div>
          </CardContent>
        </Card>
      )}

      {/* Upload Modal */}
      <UploadPresentationModal
        open={uploadModalOpen}
        onClose={() => setUploadModalOpen(false)}
        onSubmit={handleUpload}
        isLoading={uploadMutation.isPending}
        uploadProgress={uploadProgress}
      />
    </div>
  );
}
