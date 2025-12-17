// HR Constants - Single source of truth for status, kind mappings and styling

export const PERSON_STATUSES = [
  { value: 'active', label: 'Active' },
  { value: 'inactive', label: 'Inactive' },
  { value: 'onboarding', label: 'Onboarding' },
  { value: 'offboarding', label: 'Offboarding' },
] as const;

export const PERSON_STATUS_COLORS: Record<string, string> = {
  active: 'bg-green-100 text-green-800',
  inactive: 'bg-gray-100 text-gray-800',
  onboarding: 'bg-blue-100 text-blue-800',
  offboarding: 'bg-orange-100 text-orange-800',
};

export const ORG_UNIT_KINDS = [
  { value: 'company', label: 'Company' },
  { value: 'division', label: 'Division' },
  { value: 'department', label: 'Department' },
  { value: 'team', label: 'Team' },
  { value: 'region', label: 'Region' },
  { value: 'office', label: 'Office' },
  { value: 'site', label: 'Site' },
] as const;

export const ORG_UNIT_KIND_COLORS: Record<string, string> = {
  company: 'bg-purple-100 text-purple-800',
  division: 'bg-blue-100 text-blue-800',
  department: 'bg-green-100 text-green-800',
  team: 'bg-yellow-100 text-yellow-800',
  region: 'bg-orange-100 text-orange-800',
  office: 'bg-teal-100 text-teal-800',
  site: 'bg-pink-100 text-pink-800',
};

export function getPersonStatusColor(status: string): string {
  return PERSON_STATUS_COLORS[status] || 'bg-gray-100 text-gray-800';
}

export function getOrgUnitKindColor(kind: string): string {
  return ORG_UNIT_KIND_COLORS[kind] || 'bg-gray-100 text-gray-800';
}

export function getOrgUnitKindLabel(kind: string): string {
  const found = ORG_UNIT_KINDS.find((k) => k.value === kind);
  return found?.label || kind;
}

export type PersonStatus = typeof PERSON_STATUSES[number]['value'];
export type OrgUnitKind = typeof ORG_UNIT_KINDS[number]['value'];
