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
import {
  mkbContentTypeLabels,
  mkbPersonaLabels,
  mkbContextTagLabels,
  mkbStatusLabels,
  type MKBContentType,
  type MKBPersona,
  type MKBContextTag,
  type MKBArticleStatus,
} from '@/types';

interface MKBFiltersProps {
  search: string;
  onSearchChange: (value: string) => void;
  contentType: string;
  onContentTypeChange: (value: string) => void;
  persona: string;
  onPersonaChange: (value: string) => void;
  contextTag: string;
  onContextTagChange: (value: string) => void;
  status: string;
  onStatusChange: (value: string) => void;
  onClearFilters: () => void;
  hasFilters: boolean;
}

export function MKBFilters({
  search,
  onSearchChange,
  contentType,
  onContentTypeChange,
  persona,
  onPersonaChange,
  contextTag,
  onContextTagChange,
  status,
  onStatusChange,
  onClearFilters,
  hasFilters,
}: MKBFiltersProps) {
  return (
    <div className="flex flex-wrap items-center gap-3">
      <div className="relative flex-1 min-w-[200px] max-w-sm">
        <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-400" />
        <Input
          placeholder="Search content..."
          value={search}
          onChange={(e) => onSearchChange(e.target.value)}
          className="pl-9"
        />
      </div>

      <Select value={contentType} onValueChange={onContentTypeChange}>
        <SelectTrigger className="w-[140px]">
          <SelectValue placeholder="Type" />
        </SelectTrigger>
        <SelectContent>
          <SelectItem value="all">All Types</SelectItem>
          {(Object.keys(mkbContentTypeLabels) as MKBContentType[]).map((type) => (
            <SelectItem key={type} value={type}>
              {mkbContentTypeLabels[type]}
            </SelectItem>
          ))}
        </SelectContent>
      </Select>

      <Select value={persona} onValueChange={onPersonaChange}>
        <SelectTrigger className="w-[150px]">
          <SelectValue placeholder="Persona" />
        </SelectTrigger>
        <SelectContent>
          <SelectItem value="all">All Personas</SelectItem>
          {(Object.keys(mkbPersonaLabels) as MKBPersona[]).map((p) => (
            <SelectItem key={p} value={p}>
              {mkbPersonaLabels[p]}
            </SelectItem>
          ))}
        </SelectContent>
      </Select>

      <Select value={contextTag} onValueChange={onContextTagChange}>
        <SelectTrigger className="w-[150px]">
          <SelectValue placeholder="Context" />
        </SelectTrigger>
        <SelectContent>
          <SelectItem value="all">All Contexts</SelectItem>
          {(Object.keys(mkbContextTagLabels) as MKBContextTag[]).map((tag) => (
            <SelectItem key={tag} value={tag}>
              {mkbContextTagLabels[tag]}
            </SelectItem>
          ))}
        </SelectContent>
      </Select>

      <Select value={status} onValueChange={onStatusChange}>
        <SelectTrigger className="w-[130px]">
          <SelectValue placeholder="Status" />
        </SelectTrigger>
        <SelectContent>
          <SelectItem value="all">All Status</SelectItem>
          {(Object.keys(mkbStatusLabels) as MKBArticleStatus[]).map((s) => (
            <SelectItem key={s} value={s}>
              {mkbStatusLabels[s]}
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
