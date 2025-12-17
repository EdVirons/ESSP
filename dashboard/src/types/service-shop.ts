// Service Shop types
export interface ServiceShop {
  id: string;
  tenantId: string;
  countyCode: string;
  countyName: string;
  subCountyCode: string;
  subCountyName: string;
  coverageLevel: string;
  name: string;
  location: string;
  active: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface CreateServiceShopRequest {
  countyCode: string;
  countyName?: string;
  subCountyCode?: string;
  subCountyName?: string;
  coverageLevel?: string;
  name: string;
  location?: string;
  active?: boolean;
}

// Service Staff types
export type StaffRole = 'lead_technician' | 'assistant_technician' | 'storekeeper';

export interface ServiceStaff {
  id: string;
  tenantId: string;
  serviceShopId: string;
  userId: string;
  role: StaffRole;
  phone: string;
  active: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface CreateServiceStaffRequest {
  serviceShopId: string;
  userId: string;
  role: StaffRole;
  phone?: string;
  active?: boolean;
}

export interface UpdateServiceStaffRequest {
  serviceShopId?: string;
  role?: StaffRole;
  phone?: string;
  active?: boolean;
}

export interface ServiceStaffStats {
  total: number;
  active: number;
  inactive: number;
  byRole: Record<StaffRole, number>;
}
