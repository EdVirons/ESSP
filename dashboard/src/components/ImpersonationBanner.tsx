import { useImpersonation } from '@/contexts/ImpersonationContext';
import { UserCircle2, X } from 'lucide-react';
import { Button } from '@/components/ui/button';

export function ImpersonationBanner() {
  const { isImpersonating, targetUser, reason, stopImpersonation } = useImpersonation();

  if (!isImpersonating || !targetUser) {
    return null;
  }

  return (
    <div className="bg-orange-500 text-white px-4 py-2 flex items-center justify-between shadow-md z-50">
      <div className="flex items-center gap-3">
        <UserCircle2 className="h-5 w-5" />
        <div>
          <span className="font-medium">Acting as:</span>{' '}
          <span className="font-semibold">{targetUser.name}</span>
          <span className="text-orange-100 ml-2">({targetUser.email})</span>
          {reason && (
            <span className="text-orange-200 ml-3 text-sm italic">
              Reason: {reason}
            </span>
          )}
        </div>
      </div>
      <div className="flex items-center gap-2">
        <span className="text-orange-100 text-sm">
          {targetUser.schools.length} school{targetUser.schools.length !== 1 ? 's' : ''}
        </span>
        <Button
          variant="ghost"
          size="sm"
          onClick={stopImpersonation}
          className="text-white hover:bg-orange-600 hover:text-white gap-1"
        >
          <X className="h-4 w-4" />
          Stop Impersonating
        </Button>
      </div>
    </div>
  );
}
