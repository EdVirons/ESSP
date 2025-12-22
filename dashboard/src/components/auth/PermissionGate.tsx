import type { ReactNode } from 'react';
import { useAuth } from '@/contexts/AuthContext';

interface PermissionGateProps {
  children: ReactNode;
  // User needs ANY of these permissions (OR logic)
  permissions?: string[];
  // User needs ANY of these roles (OR logic)
  roles?: string[];
  // If true, requires ALL permissions/roles instead of ANY
  requireAll?: boolean;
  // What to render if access is denied (defaults to null/nothing)
  fallback?: ReactNode;
}

/**
 * PermissionGate - Conditionally renders children based on user permissions/roles
 *
 * Usage examples:
 *
 * // Show only if user can create incidents
 * <PermissionGate permissions={['incident:create']}>
 *   <CreateIncidentButton />
 * </PermissionGate>
 *
 * // Show only if user is admin or lead tech
 * <PermissionGate roles={['ssp_admin', 'ssp_lead_tech']}>
 *   <AdminPanel />
 * </PermissionGate>
 *
 * // Show with fallback for unauthorized users
 * <PermissionGate permissions={['workorder:update']} fallback={<ReadOnlyView />}>
 *   <EditableView />
 * </PermissionGate>
 *
 * // Require ALL permissions (AND logic)
 * <PermissionGate permissions={['bom:read', 'bom:update']} requireAll>
 *   <BOMEditor />
 * </PermissionGate>
 */
export function PermissionGate({
  children,
  permissions,
  roles,
  requireAll = false,
  fallback = null,
}: PermissionGateProps) {
  const { hasPermission, hasRole } = useAuth();

  // Admin bypasses all permission checks
  if (hasRole('ssp_admin')) {
    return <>{children}</>;
  }

  // If no restrictions specified, allow access
  if (!permissions?.length && !roles?.length) {
    return <>{children}</>;
  }

  const checkFn = requireAll
    ? (arr: string[], checker: (v: string) => boolean) => arr.every(checker)
    : (arr: string[], checker: (v: string) => boolean) => arr.some(checker);

  // Check permissions
  const hasRequiredPermissions = !permissions?.length || checkFn(permissions, hasPermission);

  // Check roles
  const hasRequiredRoles = !roles?.length || checkFn(roles, hasRole);

  // Both conditions must be satisfied (if both are specified)
  if (hasRequiredPermissions && hasRequiredRoles) {
    return <>{children}</>;
  }

  return <>{fallback}</>;
}

/**
 * Hook version for programmatic permission checking
 */
// eslint-disable-next-line react-refresh/only-export-components
export function usePermissionCheck(
  permissions?: string[],
  roles?: string[],
  requireAll = false
): boolean {
  const { hasPermission, hasRole } = useAuth();

  // Admin bypasses all checks
  if (hasRole('ssp_admin')) return true;

  // No restrictions = allowed
  if (!permissions?.length && !roles?.length) return true;

  const checkFn = requireAll
    ? (arr: string[], checker: (v: string) => boolean) => arr.every(checker)
    : (arr: string[], checker: (v: string) => boolean) => arr.some(checker);

  const hasRequiredPermissions = !permissions?.length || checkFn(permissions, hasPermission);
  const hasRequiredRoles = !roles?.length || checkFn(roles, hasRole);

  return hasRequiredPermissions && hasRequiredRoles;
}
