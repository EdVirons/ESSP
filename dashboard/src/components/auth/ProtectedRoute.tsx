import { Navigate, useLocation } from 'react-router-dom';
import { useAuth } from '@/contexts/AuthContext';
import { AccessDenied } from './AccessDenied';

interface ProtectedRouteProps {
  children: React.ReactNode;
  // Permissions required to access this route (OR logic - user needs ANY of these)
  permissions?: string[];
  // Roles required to access this route (OR logic - user needs ANY of these)
  roles?: string[];
}

/**
 * ProtectedRoute - Wraps routes that require authentication and optionally authorization
 *
 * Usage examples:
 *
 * // Just require authentication
 * <ProtectedRoute>
 *   <Dashboard />
 * </ProtectedRoute>
 *
 * // Require specific permission
 * <ProtectedRoute permissions={['incident:read']}>
 *   <IncidentsPage />
 * </ProtectedRoute>
 *
 * // Require specific role
 * <ProtectedRoute roles={['ssp_admin']}>
 *   <SettingsPage />
 * </ProtectedRoute>
 *
 * // Require permission OR role (user needs at least one)
 * <ProtectedRoute permissions={['workorder:read']} roles={['ssp_lead_tech']}>
 *   <WorkOrdersPage />
 * </ProtectedRoute>
 */
export function ProtectedRoute({
  children,
  permissions,
  roles,
}: ProtectedRouteProps) {
  const { isAuthenticated, isLoading, hasPermission, hasRole } = useAuth();
  const location = useLocation();

  // Show loading state while checking authentication
  if (isLoading) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-gray-50">
        <div className="flex flex-col items-center gap-4">
          <div className="h-8 w-8 animate-spin rounded-full border-4 border-cyan-600 border-t-transparent"></div>
          <p className="text-sm text-gray-500">Loading...</p>
        </div>
      </div>
    );
  }

  // Redirect to login if not authenticated
  if (!isAuthenticated) {
    return <Navigate to="/login" state={{ from: location }} replace />;
  }

  // Admin bypasses all authorization checks
  if (hasRole('ssp_admin')) {
    return <>{children}</>;
  }

  // Check authorization if permissions or roles are specified
  if (permissions?.length || roles?.length) {
    // Check permissions (OR logic)
    const hasRequiredPermission =
      !permissions?.length || permissions.some((p) => hasPermission(p));

    // Check roles (OR logic)
    const hasRequiredRole = !roles?.length || roles.some((r) => hasRole(r));

    // User needs to satisfy at least one of the conditions
    // If both are specified, either one can grant access
    const hasAccess =
      (permissions?.length && hasRequiredPermission) ||
      (roles?.length && hasRequiredRole) ||
      (!permissions?.length && !roles?.length);

    if (!hasAccess) {
      return <AccessDenied />;
    }
  }

  return <>{children}</>;
}
