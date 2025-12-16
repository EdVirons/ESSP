// Parts Catalog Types

export interface Part {
  id: string;
  tenantId: string;
  sku: string;
  name: string;
  category: string;
  description: string;
  unitCostCents: number;
  supplier: string;
  supplierSku: string;
  active: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface CreatePartRequest {
  sku: string;
  name: string;
  category?: string;
  description?: string;
  unitCostCents?: number;
  supplier?: string;
  supplierSku?: string;
}

export interface UpdatePartRequest {
  name?: string;
  category?: string;
  description?: string;
  unitCostCents?: number;
  supplier?: string;
  supplierSku?: string;
  active?: boolean;
}

export interface PartFilters {
  q?: string;
  category?: string;
  active?: boolean;
  limit?: number;
  cursor?: string;
}

export interface PartStats {
  total: number;
  byCategory: Record<string, number>;
}

export interface ImportPartsResult {
  created: number;
  failed: number;
  errors: string[];
}
