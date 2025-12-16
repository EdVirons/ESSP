import { Link } from 'react-router-dom';
import { ShieldX, Home, ArrowLeft } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { useAuth } from '@/contexts/AuthContext';
import { getDefaultRouteForRoles } from '@/config/navigation';

interface AccessDeniedProps {
  title?: string;
  message?: string;
  showHomeButton?: boolean;
  showBackButton?: boolean;
}

/**
 * AccessDenied - 403 Forbidden page component
 *
 * Displayed when a user tries to access a page they don't have permission for.
 */
export function AccessDenied({
  title = 'Access Denied',
  message = "You don't have permission to access this page.",
  showHomeButton = true,
  showBackButton = true,
}: AccessDeniedProps) {
  const { profile, user } = useAuth();
  const roles = profile?.roles || user?.roles || [];
  const defaultRoute = getDefaultRouteForRoles(roles);

  return (
    <div className="flex min-h-[60vh] flex-col items-center justify-center px-4 text-center">
      <div className="rounded-2xl bg-gradient-to-br from-rose-50 to-red-100 p-6 shadow-lg">
        <div className="flex h-20 w-20 items-center justify-center rounded-xl bg-gradient-to-br from-rose-500 to-red-600 shadow-lg shadow-rose-500/30">
          <ShieldX className="h-10 w-10 text-white" />
        </div>
      </div>

      <h1 className="mt-8 text-3xl font-bold text-gray-900">{title}</h1>
      <p className="mt-3 max-w-md text-gray-600">{message}</p>

      <div className="mt-4 rounded-lg bg-gray-50 px-4 py-2 text-sm text-gray-500">
        <span className="font-medium">Your roles:</span>{' '}
        {roles.length > 0 ? (
          roles.map((role, i) => (
            <span key={role}>
              <code className="rounded bg-gray-200 px-1.5 py-0.5 text-xs">
                {role.replace('ssp_', '')}
              </code>
              {i < roles.length - 1 && ', '}
            </span>
          ))
        ) : (
          <span className="italic">No roles assigned</span>
        )}
      </div>

      <div className="mt-8 flex gap-3">
        {showBackButton && (
          <Button
            variant="outline"
            onClick={() => window.history.back()}
            className="border-gray-300"
          >
            <ArrowLeft className="h-4 w-4" />
            Go Back
          </Button>
        )}
        {showHomeButton && (
          <Link to={defaultRoute}>
            <Button className="bg-gradient-to-r from-cyan-600 to-teal-600 hover:from-cyan-700 hover:to-teal-700">
              <Home className="h-4 w-4" />
              Go to Dashboard
            </Button>
          </Link>
        )}
      </div>

      <p className="mt-8 text-sm text-gray-400">
        If you believe this is an error, please contact your administrator.
      </p>
    </div>
  );
}
