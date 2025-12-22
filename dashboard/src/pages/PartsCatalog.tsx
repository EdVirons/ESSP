import * as React from 'react';
import {
  Package,
  Search,
  Plus,
  Upload,
  Download,
  Layers,
  DollarSign,
  CheckCircle,
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardContent } from '@/components/ui/card';
import { Select } from '@/components/ui/select';
import { DataTable } from '@/components/ui/data-table';
import { ConfirmDialog } from '@/components/ui/modal';
import { useParts, usePartsStats, usePartsCategories } from '@/api/parts';
import {
  CreatePartModal,
  EditPartModal,
  ImportPartsModal,
  createPartsColumns,
  usePartsCRUD,
} from '@/components/parts';
import type { PartFilters } from '@/types';

export function PartsCatalog() {
  const [filters, setFilters] = React.useState<PartFilters>({
    limit: 50,
  });
  const [searchQuery, setSearchQuery] = React.useState('');

  // API hooks
  const { data, isLoading, refetch } = useParts(filters);
  const { data: stats } = usePartsStats();
  const { data: categoriesData } = usePartsCategories();

  // CRUD operations hook
  const {
    createModalOpen,
    setCreateModalOpen,
    editModalOpen,
    deleteModalOpen,
    importModalOpen,
    setImportModalOpen,
    selectedPart,
    isCreating,
    isUpdating,
    isDeleting,
    isImporting,
    isExporting,
    handleCreate,
    handleEdit,
    handleUpdate,
    handleDeleteClick,
    handleDelete,
    handleImport,
    handleExport,
    closeEditModal,
    closeDeleteModal,
  } = usePartsCRUD({ onRefetch: refetch });

  const categories = React.useMemo(() => {
    const items = categoriesData?.items;
    if (!items) return [];
    return items.map((cat) => ({
      value: cat,
      label: cat.charAt(0).toUpperCase() + cat.slice(1),
    }));
  }, [categoriesData?.items]);

  const handleSearch = () => {
    setFilters((prev) => ({ ...prev, q: searchQuery || undefined, cursor: undefined }));
  };

  // Table columns
  const columns = React.useMemo(
    () =>
      createPartsColumns({
        onEdit: handleEdit,
        onDelete: handleDeleteClick,
      }),
    [handleEdit, handleDeleteClick]
  );

  const items = data?.items || [];
  const total = stats?.total || items.length;
  const categoryCount = stats?.byCategory ? Object.keys(stats.byCategory).length : 0;
  const activeCount = items.filter((p) => p.active).length;

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Parts Catalog </h1>
          <p className="text-gray-500">Manage your parts inventory</p>
        </div>
        <div className="flex gap-2">
          <Button variant="outline" onClick={() => setImportModalOpen(true)}>
            <Upload className="mr-2 h-4 w-4" />
            Import
          </Button>
          <Button variant="outline" onClick={handleExport} disabled={isExporting}>
            <Download className="mr-2 h-4 w-4" />
            {isExporting ? 'Exporting...' : 'Export'}
          </Button>
          <Button onClick={() => setCreateModalOpen(true)}>
            <Plus className="mr-2 h-4 w-4" />
            Create Part
          </Button>
        </div>
      </div>

      {/* Stats */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        <Card>
          <CardContent className="pt-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-gray-500">Total Parts</p>
                <p className="text-2xl font-bold text-gray-900">{total}</p>
              </div>
              <div className="h-12 w-12 rounded-full bg-amber-50 flex items-center justify-center">
                <Package className="h-6 w-6 text-amber-600" />
              </div>
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="pt-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-gray-500">Categories</p>
                <p className="text-2xl font-bold text-gray-900">{categoryCount}</p>
              </div>
              <div className="h-12 w-12 rounded-full bg-purple-50 flex items-center justify-center">
                <Layers className="h-6 w-6 text-purple-600" />
              </div>
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="pt-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-gray-500">Active</p>
                <p className="text-2xl font-bold text-gray-900">{activeCount}</p>
              </div>
              <div className="h-12 w-12 rounded-full bg-green-50 flex items-center justify-center">
                <CheckCircle className="h-6 w-6 text-green-600" />
              </div>
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="pt-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-gray-500">Showing</p>
                <p className="text-2xl font-bold text-gray-900">{items.length}</p>
              </div>
              <div className="h-12 w-12 rounded-full bg-blue-50 flex items-center justify-center">
                <DollarSign className="h-6 w-6 text-blue-600" />
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Filters */}
      <div className="flex flex-col sm:flex-row gap-4">
        <div className="flex-1 flex gap-2">
          <div className="relative flex-1 max-w-md">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-gray-400" />
            <Input
              placeholder="Search by name, SKU, or supplier..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              onKeyDown={(e) => e.key === 'Enter' && handleSearch()}
              className="pl-9"
            />
          </div>
          <Button variant="outline" onClick={handleSearch}>
            Search
          </Button>
        </div>
        <div className="flex gap-2">
          <Select
            value={filters.category || ''}
            onChange={(value) =>
              setFilters((prev) => ({ ...prev, category: value || undefined, cursor: undefined }))
            }
            options={[{ value: '', label: 'All Categories' }, ...categories]}
            className="w-40"
          />
          <Select
            value={filters.active === undefined ? '' : filters.active ? 'active' : 'inactive'}
            onChange={(value) =>
              setFilters((prev) => ({
                ...prev,
                active: value === '' ? undefined : value === 'active',
                cursor: undefined,
              }))
            }
            options={[
              { value: '', label: 'All Status' },
              { value: 'active', label: 'Active' },
              { value: 'inactive', label: 'Inactive' },
            ]}
            className="w-32"
          />
        </div>
      </div>

      {/* Data Table */}
      <DataTable
        columns={columns}
        data={items}
        isLoading={isLoading}
        emptyMessage="No parts found. Create your first part to get started."
      />

      {/* Load More */}
      {data?.nextCursor && (
        <div className="flex justify-center">
          <Button
            variant="outline"
            onClick={() =>
              setFilters((prev) => ({ ...prev, cursor: data.nextCursor || undefined }))
            }
          >
            Load More
          </Button>
        </div>
      )}

      {/* Create Modal */}
      <CreatePartModal
        open={createModalOpen}
        onClose={() => setCreateModalOpen(false)}
        onSubmit={handleCreate}
        isLoading={isCreating}
        categories={categoriesData?.items}
      />

      {/* Edit Modal */}
      <EditPartModal
        part={selectedPart}
        open={editModalOpen}
        onClose={closeEditModal}
        onSubmit={handleUpdate}
        isLoading={isUpdating}
        categories={categoriesData?.items}
      />

      {/* Delete Confirm Modal */}
      <ConfirmDialog
        open={deleteModalOpen}
        onClose={closeDeleteModal}
        onConfirm={handleDelete}
        title="Delete Part"
        description={`Are you sure you want to delete "${selectedPart?.name}"? This action cannot be undone.`}
        confirmText="Delete"
        isLoading={isDeleting}
        variant="destructive"
      />

      {/* Import Modal */}
      <ImportPartsModal
        open={importModalOpen}
        onClose={() => setImportModalOpen(false)}
        onImport={handleImport}
        isLoading={isImporting}
      />
    </div>
  );
}
