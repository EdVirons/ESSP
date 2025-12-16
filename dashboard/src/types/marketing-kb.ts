// Marketing Knowledge Base Types

export type MKBContentType = 'messaging' | 'case_study' | 'deck' | 'objection' | 'roi';
export type MKBPersona = 'director' | 'principal' | 'teacher' | 'parent' | 'it_admin' | 'county_official';
export type MKBContextTag = 'rural' | 'urban' | 'low_connectivity' | 'no_isp' | 'connected' | 'private' | 'public' | 'cbc' | 'igcse' | '8-4-4';
export type MKBArticleStatus = 'draft' | 'review' | 'approved' | 'archived';

export interface MKBArticle {
  id: string;
  tenantId: string;
  title: string;
  slug: string;
  summary: string;
  content: string;
  contentType: MKBContentType;
  personas: string[];
  contextTags: string[];
  tags: string[];
  version: number;
  status: MKBArticleStatus;
  usageCount: number;
  lastUsedAt?: string;
  createdById: string;
  createdByName: string;
  updatedById: string;
  updatedByName: string;
  approvedAt?: string;
  approvedById?: string;
  approvedByName?: string;
  createdAt: string;
  updatedAt: string;
}

export interface PitchKit {
  id: string;
  tenantId: string;
  name: string;
  description: string;
  targetPersona: string;
  contextTags: string[];
  articleIds: string[];
  articles?: MKBArticle[];
  isTemplate: boolean;
  createdById: string;
  createdByName: string;
  updatedById: string;
  updatedByName: string;
  createdAt: string;
  updatedAt: string;
}

export interface CreateMKBArticleRequest {
  title: string;
  slug?: string;
  summary?: string;
  content: string;
  contentType?: MKBContentType;
  personas?: string[];
  contextTags?: string[];
  tags?: string[];
}

export interface UpdateMKBArticleRequest {
  title?: string;
  slug?: string;
  summary?: string;
  content?: string;
  contentType?: MKBContentType;
  personas?: string[];
  contextTags?: string[];
  tags?: string[];
  status?: 'draft' | 'review';
}

export interface CreatePitchKitRequest {
  name: string;
  description?: string;
  targetPersona?: string;
  contextTags?: string[];
  articleIds?: string[];
  isTemplate?: boolean;
}

export interface UpdatePitchKitRequest {
  name?: string;
  description?: string;
  targetPersona?: string;
  contextTags?: string[];
  articleIds?: string[];
  isTemplate?: boolean;
}

export interface MKBArticleFilters {
  q?: string;
  contentType?: MKBContentType;
  persona?: MKBPersona;
  contextTag?: MKBContextTag;
  status?: MKBArticleStatus;
  limit?: number;
  cursor?: string;
}

export interface PitchKitFilters {
  targetPersona?: MKBPersona;
  isTemplate?: boolean;
  limit?: number;
  cursor?: string;
}

export interface MKBStats {
  total: number;
  approved: number;
  inReview: number;
  draft: number;
  byContentType: Record<string, number>;
  byPersona: Record<string, number>;
  byContextTag: Record<string, number>;
}

// Display labels for enum values
export const mkbContentTypeLabels: Record<MKBContentType, string> = {
  messaging: 'Messaging',
  case_study: 'Case Study',
  deck: 'Deck',
  objection: 'Objection Handler',
  roi: 'ROI Calculator',
};

export const mkbContentTypeIcons: Record<MKBContentType, string> = {
  messaging: 'MessageSquare',
  case_study: 'BookOpen',
  deck: 'Presentation',
  objection: 'Shield',
  roi: 'Calculator',
};

export const mkbPersonaLabels: Record<MKBPersona, string> = {
  director: 'Director',
  principal: 'Principal',
  teacher: 'Teacher',
  parent: 'Parent',
  it_admin: 'IT Admin',
  county_official: 'County Official',
};

export const mkbContextTagLabels: Record<MKBContextTag, string> = {
  rural: 'Rural',
  urban: 'Urban',
  low_connectivity: 'Low Connectivity',
  no_isp: 'No ISP',
  connected: 'Connected',
  private: 'Private',
  public: 'Public',
  cbc: 'CBC',
  igcse: 'IGCSE',
  '8-4-4': '8-4-4',
};

export const mkbStatusLabels: Record<MKBArticleStatus, string> = {
  draft: 'Draft',
  review: 'In Review',
  approved: 'Approved',
  archived: 'Archived',
};

export const mkbStatusColors: Record<MKBArticleStatus, string> = {
  draft: 'bg-amber-100 text-amber-800',
  review: 'bg-blue-100 text-blue-800',
  approved: 'bg-green-100 text-green-800',
  archived: 'bg-gray-100 text-gray-800',
};
