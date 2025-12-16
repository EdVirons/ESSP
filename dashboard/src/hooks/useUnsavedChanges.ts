import { useEffect, useCallback } from 'react';
import { useBlocker, type Location } from 'react-router-dom';

/**
 * Hook to warn users about unsaved changes when navigating away.
 * Shows browser's native confirmation dialog when trying to close/refresh the page.
 * Shows a custom blocker when navigating within the app (React Router).
 *
 * @param isDirty - Whether there are unsaved changes
 * @param message - Custom message to show (optional, defaults to generic message)
 *
 * @example
 * ```tsx
 * const [formData, setFormData] = useState(initialData);
 * const isDirty = JSON.stringify(formData) !== JSON.stringify(initialData);
 *
 * useUnsavedChanges(isDirty);
 * ```
 */
export function useUnsavedChanges(
  isDirty: boolean,
  message = 'You have unsaved changes. Are you sure you want to leave?'
) {
  // Handle browser close/refresh
  useEffect(() => {
    const handleBeforeUnload = (e: BeforeUnloadEvent) => {
      if (isDirty) {
        e.preventDefault();
        // Modern browsers ignore custom messages, but we still set returnValue
        e.returnValue = message;
        return message;
      }
    };

    window.addEventListener('beforeunload', handleBeforeUnload);
    return () => window.removeEventListener('beforeunload', handleBeforeUnload);
  }, [isDirty, message]);

  // Handle React Router navigation
  const blocker = useBlocker(
    useCallback(
      ({ currentLocation, nextLocation }: { currentLocation: Location; nextLocation: Location }) => {
        return isDirty && currentLocation.pathname !== nextLocation.pathname;
      },
      [isDirty]
    )
  );

  // Reset blocker when user confirms navigation
  const confirmNavigation = useCallback(() => {
    if (blocker.state === 'blocked') {
      blocker.proceed();
    }
  }, [blocker]);

  // Cancel navigation
  const cancelNavigation = useCallback(() => {
    if (blocker.state === 'blocked') {
      blocker.reset();
    }
  }, [blocker]);

  return {
    isBlocked: blocker.state === 'blocked',
    confirmNavigation,
    cancelNavigation,
    message,
  };
}

/**
 * Simpler version that only warns on browser close/refresh.
 * Use this when you don't need to block React Router navigation.
 *
 * @param isDirty - Whether there are unsaved changes
 *
 * @example
 * ```tsx
 * useBeforeUnload(hasUnsavedChanges);
 * ```
 */
export function useBeforeUnload(isDirty: boolean) {
  useEffect(() => {
    const handleBeforeUnload = (e: BeforeUnloadEvent) => {
      if (isDirty) {
        e.preventDefault();
        e.returnValue = '';
      }
    };

    window.addEventListener('beforeunload', handleBeforeUnload);
    return () => window.removeEventListener('beforeunload', handleBeforeUnload);
  }, [isDirty]);
}
