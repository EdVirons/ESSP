import { useAuth } from '@/contexts/AuthContext';
import { getDefaultRouteForRoles } from '@/config/navigation';

/**
 * useDefaultRoute - Returns the default route based on user's roles
 *
 * This hook determines the best landing page for a user based on their role.
 * It helps redirect users to the most relevant page after login.
 *
 * Role priorities:
 * - ssp_school_contact → /incidents (they create incidents)
 * - ssp_support_agent → /incidents (they handle the queue)
 * - ssp_lead_tech → /work-orders (they manage repairs)
 * - ssp_field_tech → /work-orders (they do repairs)
 * - ssp_warehouse_manager → /parts-catalog (they manage inventory)
 * - ssp_demo_team → /projects (they handle demos)
 * - ssp_sales → /projects (they handle sales)
 * - Default → /overview (general dashboard)
 *
 * Usage:
 * const defaultRoute = useDefaultRoute();
 * navigate(defaultRoute);
 */
export function useDefaultRoute(): string {
  const { profile, user } = useAuth();
  const roles = profile?.roles || user?.roles || [];

  return getDefaultRouteForRoles(roles);
}
