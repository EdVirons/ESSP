// Knowledge Base Types

export type KBContentType = 'runbook' | 'troubleshooting' | 'kedb' | 'checklist' | 'sop';
export type KBModule = 'learning_portal' | 'mdm' | 'sso' | 'devices' | 'inventory' | 'general';
export type KBLifecycleStage = 'demo' | 'install' | 'commission' | 'support';
export type KBArticleStatus = 'draft' | 'published' | 'archived';

export interface KBArticle {
  id: string;
  tenantId: string;
  title: string;
  slug: string;
  summary: string;
  content: string;
  contentType: KBContentType;
  module: KBModule;
  lifecycleStage: KBLifecycleStage;
  tags: string[];
  version: number;
  status: KBArticleStatus;
  createdById: string;
  createdByName: string;
  updatedById: string;
  updatedByName: string;
  publishedAt?: string;
  createdAt: string;
  updatedAt: string;
}

export interface CreateKBArticleRequest {
  title: string;
  slug?: string;
  summary?: string;
  content: string;
  contentType?: KBContentType;
  module?: KBModule;
  lifecycleStage?: KBLifecycleStage;
  tags?: string[];
}

export interface UpdateKBArticleRequest {
  title?: string;
  slug?: string;
  summary?: string;
  content?: string;
  contentType?: KBContentType;
  module?: KBModule;
  lifecycleStage?: KBLifecycleStage;
  tags?: string[];
}

export interface KBArticleFilters {
  q?: string;
  contentType?: KBContentType;
  module?: KBModule;
  lifecycleStage?: KBLifecycleStage;
  status?: KBArticleStatus;
  limit?: number;
  cursor?: string;
}

export interface KBStats {
  total: number;
  published: number;
  draft: number;
  byContentType: Record<string, number>;
  byModule: Record<string, number>;
  byLifecycleStage: Record<string, number>;
}

// Display labels for enum values
export const contentTypeLabels: Record<KBContentType, string> = {
  runbook: 'Runbook',
  troubleshooting: 'Troubleshooting',
  kedb: 'Known Issues',
  checklist: 'Checklist',
  sop: 'SOP',
};

export const moduleLabels: Record<KBModule, string> = {
  learning_portal: 'Learning Portal',
  mdm: 'MDM',
  sso: 'SSO',
  devices: 'Devices',
  inventory: 'Inventory',
  general: 'General',
};

export const lifecycleStageLabels: Record<KBLifecycleStage, string> = {
  demo: 'Demo',
  install: 'Install',
  commission: 'Commission',
  support: 'Support',
};

export const statusLabels: Record<KBArticleStatus, string> = {
  draft: 'Draft',
  published: 'Published',
  archived: 'Archived',
};
