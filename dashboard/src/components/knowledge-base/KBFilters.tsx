import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Search, X } from 'lucide-react';
import type { KBContentType, KBModule, KBLifecycleStage, KBArticleStatus } from '@/types';
import { contentTypeLabels, moduleLabels, lifecycleStageLabels, statusLabels } from '@/types';

interface KBFiltersProps {
  search: string;
  onSearchChange: (value: string) => void;
  contentType: string;
  onContentTypeChange: (value: string) => void;
  module: string;
  onModuleChange: (value: string) => void;
  lifecycleStage: string;
  onLifecycleStageChange: (value: string) => void;
  status: string;
  onStatusChange: (value: string) => void;
  onClearFilters: () => void;
  hasFilters: boolean;
}

export function KBFilters({
  search,
  onSearchChange,
  contentType,
  onContentTypeChange,
  module,
  onModuleChange,
  lifecycleStage,
  onLifecycleStageChange,
  status,
  onStatusChange,
  onClearFilters,
  hasFilters,
}: KBFiltersProps) {
  return (
    <div className="flex flex-wrap items-center gap-4">
      <div className="relative flex-1 min-w-[200px] max-w-sm">
        <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-gray-400" />
        <Input
          placeholder="Search articles..."
          className="pl-10"
          value={search}
          onChange={(e) => onSearchChange(e.target.value)}
        />
      </div>

      <Select value={contentType} onValueChange={onContentTypeChange}>
        <SelectTrigger className="w-[160px]">
          <SelectValue placeholder="Content Type" />
        </SelectTrigger>
        <SelectContent>
          <SelectItem value="all">All Types</SelectItem>
          {(Object.keys(contentTypeLabels) as KBContentType[]).map((type) => (
            <SelectItem key={type} value={type}>
              {contentTypeLabels[type]}
            </SelectItem>
          ))}
        </SelectContent>
      </Select>

      <Select value={module} onValueChange={onModuleChange}>
        <SelectTrigger className="w-[160px]">
          <SelectValue placeholder="Module" />
        </SelectTrigger>
        <SelectContent>
          <SelectItem value="all">All Modules</SelectItem>
          {(Object.keys(moduleLabels) as KBModule[]).map((mod) => (
            <SelectItem key={mod} value={mod}>
              {moduleLabels[mod]}
            </SelectItem>
          ))}
        </SelectContent>
      </Select>

      <Select value={lifecycleStage} onValueChange={onLifecycleStageChange}>
        <SelectTrigger className="w-[160px]">
          <SelectValue placeholder="Stage" />
        </SelectTrigger>
        <SelectContent>
          <SelectItem value="all">All Stages</SelectItem>
          {(Object.keys(lifecycleStageLabels) as KBLifecycleStage[]).map((stage) => (
            <SelectItem key={stage} value={stage}>
              {lifecycleStageLabels[stage]}
            </SelectItem>
          ))}
        </SelectContent>
      </Select>

      <Select value={status} onValueChange={onStatusChange}>
        <SelectTrigger className="w-[140px]">
          <SelectValue placeholder="Status" />
        </SelectTrigger>
        <SelectContent>
          <SelectItem value="all">All Status</SelectItem>
          {(Object.keys(statusLabels) as KBArticleStatus[]).map((s) => (
            <SelectItem key={s} value={s}>
              {statusLabels[s]}
            </SelectItem>
          ))}
        </SelectContent>
      </Select>

      {hasFilters && (
        <Button variant="ghost" size="sm" onClick={onClearFilters}>
          <X className="h-4 w-4 mr-1" />
          Clear
        </Button>
      )}
    </div>
  );
}
