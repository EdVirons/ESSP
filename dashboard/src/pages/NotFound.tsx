import { Link } from 'react-router-dom';
import { Home, ArrowLeft } from 'lucide-react';
import { Button } from '@/components/ui/button';

export function NotFound() {
  return (
    <div className="flex min-h-[60vh] flex-col items-center justify-center text-center">
      <div className="space-y-4">
        <div className="text-6xl font-bold text-gray-200">404</div>
        <h1 className="text-2xl font-bold text-gray-900">Page Not Found</h1>
        <p className="text-gray-500 max-w-md">
          The page you are looking for does not exist or has been moved.
        </p>
        <div className="flex items-center justify-center gap-4 pt-4">
          <Button variant="outline" asChild>
            <Link to={-1 as unknown as string}>
              <ArrowLeft className="h-4 w-4" />
              Go Back
            </Link>
          </Button>
          <Button asChild>
            <Link to="/overview">
              <Home className="h-4 w-4" />
              Dashboard
            </Link>
          </Button>
        </div>
      </div>
    </div>
  );
}
