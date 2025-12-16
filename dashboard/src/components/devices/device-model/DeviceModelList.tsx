import {
  Search,
  Plus,
  Edit,
  Trash2,
  Laptop,
  ChevronDown,
  ChevronRight,
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import { cn } from '@/lib/utils';
import { DEVICE_CATEGORY_OPTIONS } from '@/types/device';
import type { DeviceModel } from './types';
import { categoryIcons } from './categoryIcons';

interface DeviceModelListProps {
  isLoading: boolean;
  searchQuery: string;
  onSearchChange: (query: string) => void;
  categoryFilter: string;
  onCategoryFilterChange: (category: string) => void;
  sortedMakes: string[];
  groupedModels: Record<string, DeviceModel[]>;
  expandedMakes: Set<string>;
  onToggleMake: (make: string) => void;
  onCreateClick: () => void;
  onEditClick: (model: DeviceModel) => void;
  onDeleteClick: (model: DeviceModel) => void;
}

export function DeviceModelList({
  isLoading,
  searchQuery,
  onSearchChange,
  categoryFilter,
  onCategoryFilterChange,
  sortedMakes,
  groupedModels,
  expandedMakes,
  onToggleMake,
  onCreateClick,
  onEditClick,
  onDeleteClick,
}: DeviceModelListProps) {
  return (
    <div className="h-[500px] flex flex-col">
      {/* Search and Filters */}
      <div className="p-4 border-b border-gray-200 space-y-3">
        <div className="flex gap-3">
          <div className="relative flex-1">
            <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-400" />
            <Input
              placeholder="Search models..."
              value={searchQuery}
              onChange={(e) => onSearchChange(e.target.value)}
              className="pl-9"
            />
          </div>
          <Button onClick={onCreateClick}>
            <Plus className="h-4 w-4" />
            Add Model
          </Button>
        </div>

        {/* Category filter tabs */}
        <div className="flex gap-2">
          <button
            type="button"
            onClick={() => onCategoryFilterChange('')}
            className={cn(
              'px-3 py-1.5 text-sm rounded-md transition-colors',
              !categoryFilter
                ? 'bg-blue-100 text-blue-700'
                : 'text-gray-600 hover:bg-gray-100'
            )}
          >
            All
          </button>
          {DEVICE_CATEGORY_OPTIONS.map((cat) => (
            <button
              key={cat.value}
              type="button"
              onClick={() => onCategoryFilterChange(cat.value)}
              className={cn(
                'px-3 py-1.5 text-sm rounded-md transition-colors flex items-center gap-1.5',
                categoryFilter === cat.value
                  ? 'bg-blue-100 text-blue-700'
                  : 'text-gray-600 hover:bg-gray-100'
              )}
            >
              {categoryIcons[cat.value]}
              {cat.label}
            </button>
          ))}
        </div>
      </div>

      {/* Models list */}
      <div className="flex-1 overflow-auto p-4">
        {isLoading ? (
          <div className="flex items-center justify-center h-full">
            <div className="h-8 w-8 animate-spin rounded-full border-2 border-gray-300 border-t-blue-600" />
          </div>
        ) : sortedMakes.length === 0 ? (
          <div className="flex flex-col items-center justify-center h-full text-gray-500">
            <Laptop className="h-12 w-12 mb-3 text-gray-300" />
            <p>No models found</p>
            {searchQuery && (
              <button
                type="button"
                onClick={() => onSearchChange('')}
                className="text-blue-600 hover:text-blue-800 mt-2"
              >
                Clear search
              </button>
            )}
          </div>
        ) : (
          <div className="space-y-3">
            {sortedMakes.map((make) => {
              const isExpanded = expandedMakes.has(make);
              const makeModels = groupedModels[make];

              return (
                <div
                  key={make}
                  className="border border-gray-200 rounded-lg overflow-hidden"
                >
                  {/* Make header */}
                  <button
                    type="button"
                    onClick={() => onToggleMake(make)}
                    className="w-full flex items-center justify-between px-4 py-3 bg-gray-50 hover:bg-gray-100 transition-colors"
                  >
                    <div className="flex items-center gap-2">
                      {isExpanded ? (
                        <ChevronDown className="h-4 w-4 text-gray-500" />
                      ) : (
                        <ChevronRight className="h-4 w-4 text-gray-500" />
                      )}
                      <span className="font-medium text-gray-900">{make}</span>
                      <Badge variant="secondary" className="ml-2">
                        {makeModels.length}
                      </Badge>
                    </div>
                  </button>

                  {/* Models list */}
                  {isExpanded && (
                    <div className="divide-y divide-gray-100">
                      {makeModels.map((model) => (
                        <div
                          key={model.id}
                          className="flex items-center justify-between px-4 py-3 hover:bg-gray-50"
                        >
                          <div className="flex items-center gap-3">
                            <div className="text-gray-400">
                              {categoryIcons[model.category]}
                            </div>
                            <div>
                              <div className="font-medium text-gray-900">
                                {model.model}
                              </div>
                              <div className="text-sm text-gray-500 capitalize">
                                {model.category}
                              </div>
                            </div>
                          </div>
                          <div className="flex items-center gap-1">
                            <Button
                              variant="ghost"
                              size="sm"
                              onClick={() => onEditClick(model)}
                              className="h-8 w-8 p-0"
                            >
                              <Edit className="h-4 w-4" />
                            </Button>
                            <Button
                              variant="ghost"
                              size="sm"
                              onClick={() => onDeleteClick(model)}
                              className="h-8 w-8 p-0 text-red-600 hover:text-red-700 hover:bg-red-50"
                            >
                              <Trash2 className="h-4 w-4" />
                            </Button>
                          </div>
                        </div>
                      ))}
                    </div>
                  )}
                </div>
              );
            })}
          </div>
        )}
      </div>
    </div>
  );
}
