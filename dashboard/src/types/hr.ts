// HR SSOT Types

export interface PersonSnapshot {
  tenantId: string;
  personId: string;
  orgUnitId: string;
  status: string;
  givenName: string;
  familyName: string;
  fullName: string;
  email: string;
  phone?: string;
  title?: string;
  updatedAt: string;
}

export interface TeamSnapshot {
  tenantId: string;
  teamId: string;
  orgUnitId: string;
  key: string;
  name: string;
  description?: string;
  updatedAt: string;
}

export interface OrgUnitSnapshot {
  tenantId: string;
  orgUnitId: string;
  parentId?: string;
  code: string;
  name: string;
  kind: string;
  updatedAt: string;
}

export interface TeamMembershipSnapshot {
  tenantId: string;
  membershipId: string;
  teamId: string;
  personId: string;
  role: string;
  status: string;
  startedAt?: string;
  endedAt?: string;
  updatedAt: string;
}

export interface PersonFilters {
  q?: string;
  orgUnitId?: string;
  status?: string;
  email?: string;
  limit?: number;
  offset?: number;
}

export interface TeamFilters {
  q?: string;
  orgUnitId?: string;
  key?: string;
  limit?: number;
  offset?: number;
}

export interface OrgUnitFilters {
  parentId?: string;
  kind?: string;
  limit?: number;
  offset?: number;
}

export interface TeamMembershipFilters {
  teamId?: string;
  personId?: string;
  role?: string;
  status?: string;
  limit?: number;
  offset?: number;
}

export interface OrgTreeNode extends OrgUnitSnapshot {
  children?: OrgTreeNode[];
}

// Create input types for CRUD operations
export interface CreatePersonInput {
  orgUnitId?: string;
  status?: string;
  givenName: string;
  familyName: string;
  email: string;
  phone?: string;
  title?: string;
  attributes?: Record<string, unknown>;
}

export interface CreateTeamInput {
  orgUnitId?: string;
  key: string;
  name: string;
  description?: string;
  metadata?: Record<string, unknown>;
}

export interface CreateOrgUnitInput {
  parentId?: string;
  code: string;
  name: string;
  kind?: string;
  metadata?: Record<string, unknown>;
}

export interface CreateTeamMembershipInput {
  teamId: string;
  personId: string;
  role?: string;
  status?: string;
}
